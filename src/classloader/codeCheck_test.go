/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"io"
	"jacobin/src/globals"
	"jacobin/src/opcodes"
	"os"
	"strings"
	"testing"
)

// Helper function to create a basic constant pool
func createBasicCP() CPool {
	CP := CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	return CP
}

// Helper function to create a constant pool with specific entry types
func createCPWithEntry(index int, entryType int) CPool {
	CP := createBasicCP()
	if index < len(CP.CpIndex) {
		CP.CpIndex[index] = CpEntry{Type: uint16(entryType), Slot: 0}
	}
	return CP
}

// ==================== CheckCodeValidity Main Functions ====================

func TestCheckCodeValidity_NilCodePointer(t *testing.T) {
	globals.InitGlobals("test")

	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(nil, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for nil codePtr, but got none")
	}
	if !strings.Contains(err.Error(), "ptr to code segment is nil") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestCheckCodeValidity_EmptyCodeNonAbstract(t *testing.T) {
	globals.InitGlobals("test")

	var code []byte // nil code
	cp := createBasicCP()
	af := AccessFlags{ClassIsAbstract: false}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for empty code in non-abstract class, but got none")
	}
	if !strings.Contains(err.Error(), "Empty code segment") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestCheckCodeValidity_EmptyCodeAbstract(t *testing.T) {
	globals.InitGlobals("test")

	var code []byte // nil code
	cp := createBasicCP()
	af := AccessFlags{ClassIsAbstract: true}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("Expected no error for empty code in abstract class, but got: %s", err.Error())
	}
}

func TestCheckCodeValidity_NilConstantPool(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x00} // NOP

	err := CheckCodeValidity(&code, nil, 5, AccessFlags{})
	if err == nil {
		t.Errorf("Expected error for nil constant pool, but got none")
	}
	if !strings.Contains(err.Error(), "ptr to constant pool is nil") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestCheckCodeValidity_EmptyConstantPool(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x00} // NOP
	cp := CPool{}        // empty CP
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for empty constant pool, but got none")
	}
	if !strings.Contains(err.Error(), "empty constant pool") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestCheckCodeValidity_ValidCode(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x00, 0x01, 0xB1} // NOP, ACONST_NULL, RETURN
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("Expected no error for valid code, but got: %s", err.Error())
	}
}

func TestCheckCodeValidity_InvalidBytecodeLength(t *testing.T) {
	globals.InitGlobals("test")

	// BIPUSH requires 1 additional byte but we don't provide it
	code := []byte{0x10} // BIPUSH
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid bytecode length, but got none")
	}
	if !strings.Contains(err.Error(), "Invalid bytecode or argument") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestArith_StackDecrement(t *testing.T) { // test whether this is recognized as a 32-bit value push
	globals.InitGlobals("test")

	code := []byte{opcodes.IADD} // IADD
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// ==================== Utility Functions Tests ====================

func TestReturn1(t *testing.T) {
	result := Return1()
	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
}

func TestReturn2(t *testing.T) {
	result := Return2()
	if result != 2 {
		t.Errorf("Expected return value 2, got: %d", result)
	}
}

func TestReturn3(t *testing.T) {
	result := Return3()
	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

func TestReturn4(t *testing.T) {
	result := Return4()
	if result != 4 {
		t.Errorf("Expected return value 4, got: %d", result)
	}
}

func TestReturn5(t *testing.T) {
	result := Return5()
	if result != 5 {
		t.Errorf("Expected return value 5, got: %d", result)
	}
}

func TestByteCodeIsForLongOrDouble_LongDoubleCodes(t *testing.T) {
	globals.InitGlobals("test")

	// Test some long/double bytecodes
	longDoubleCodes := []byte{0x09, 0x0A, 0x0E, 0x0F, 0x14, 0x16, 0x18}
	for _, code := range longDoubleCodes {
		result := BytecodeIsForLongOrDouble(code)
		if !result {
			t.Errorf("Expected true for long/double bytecode 0x%02X, got false", code)
		}
	}
}

func TestByteCodeIsForLongOrDouble_OtherCodes(t *testing.T) {
	globals.InitGlobals("test")

	// Test some non-long/double bytecodes
	otherCodes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x10, 0x11}
	for _, code := range otherCodes {
		result := BytecodeIsForLongOrDouble(code)
		if result {
			t.Errorf("Expected false for non-long/double bytecode 0x%02X, got true", code)
		}
	}
}

// ==================== Individual bytecode checks in alphabetical order of the instructions ========

// ACONST_NULL

func TestCheckAconstnull_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.ACONST_NULL}
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckAconstnull_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	StackEntries = 0
	result := CheckAconstnull()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
}

// BIPUSH

func TestCheckBipush_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.BIPUSH, 0x42}
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckBipush_ValidLength(t *testing.T) {
	globals.InitGlobals("test")

	Code = []byte{opcodes.BIPUSH, 0x42}
	PC = 0
	StackEntries = 0

	result := CheckBipush()

	if result != 2 {
		t.Errorf("Expected return value 2, got: %d", result)
	}
	if StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
}

func TestCheckBipush_InsufficientLength(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.BIPUSH} // BIPUSH with missing byte
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for insufficient BIPUSH length, but got none")
	}
}

// DUP

func TestDup_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.DUP} // DUP
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestDup(t *testing.T) {
	globals.InitGlobals("test")

	StackEntries = 1
	result := CheckDup1()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
}

// DUP2

func TestDup2_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x09, 0x5C, 0x00} // LCONST_0, DUP2, NOP (extra byte for DUP2 to check next instruction)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestDup2_HighLevel2(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x03, 0x5C} // ICONST_0, DUP2
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestDup2_LongDoubleOperation(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is for long/double
	Code = []byte{opcodes.NOP, opcodes.DUP2, opcodes.LADD}
	PC = 1
	PrevPC = 0
	StackEntries = 1

	result := CheckDup2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
	// Check that DUP2 was converted to DUP
	if Code[PC] != opcodes.DUP {
		t.Errorf("Expected DUP2 to be converted to DUP, got: 0x%x", Code[PC])
	}
}

func TestDup2_RegularOperation(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is NOT for long/double
	Code = []byte{opcodes.DUP2, opcodes.IADD}
	PC = 0
	StackEntries = 2

	result := CheckDup2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 4 {
		t.Errorf("Expected StackEntries to increase by 2, got: %d", StackEntries)
	}
}

// DUP2_X1 test when the preceding bytecode is for long or double
func TestDup2_X1_Preceding(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is for long/double
	Code = []byte{opcodes.LADD, opcodes.DUP2_X1, opcodes.NOP}
	PC = 1
	PrevPC = 0
	StackEntries = 1

	result := CheckDup2x1()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
	// Check that DUP2_X1 was converted to DUP_X1
	if Code[PC] != opcodes.DUP_X1 {
		t.Errorf("Expected DUP2_X1 to be converted to DUP_X1, got: 0x%x", Code[PC])
	}
}

// DUP2_X1 test when the following bytecode is for long or double
func TestDup2_X1_Following(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is for long/double
	Code = []byte{opcodes.NOP, opcodes.DUP2_X1, opcodes.LADD}
	PC = 1
	PrevPC = 0
	StackEntries = 1

	result := CheckDup2x1()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
	// Check that DUP2_X1 was converted to DUP_X1
	if Code[PC] != opcodes.DUP_X1 {
		t.Errorf("Expected DUP2_X1 to be converted to DUP_X1, got: 0x%x", Code[PC])
	}
}

// DUP2_X2 test when the preceding bytecode is for long or double
func TestDup2_X2_Preceding(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is for long/double
	Code = []byte{opcodes.LADD, opcodes.DUP2_X2, opcodes.NOP}
	PC = 1
	PrevPC = 0
	StackEntries = 1

	result := CheckDup2x2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
	// Check that DUP2_X1 was converted to DUP_X1
	if Code[PC] != opcodes.DUP_X2 {
		t.Errorf("Expected DUP2_X2 to be converted to DUP_X2, got: 0x%x", Code[PC])
	}
}

// DUP2_X2 test when the following bytecode is for long or double
func TestDup2_X2_Following(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is for long/double
	Code = []byte{opcodes.NOP, opcodes.DUP2_X2, opcodes.LADD}
	PC = 1
	PrevPC = 0
	StackEntries = 1

	result := CheckDup2x2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
	// Check that DUP2_X1 was converted to DUP_X1
	if Code[PC] != opcodes.DUP_X2 {
		t.Errorf("Expected DUP2_X2 to be converted to DUP_X2, got: 0x%x", Code[PC])
	}
}

// FCONST_0
func TestPushFloat0_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.FCONST_0} // FCONST_0
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// FLOAD
func TestPushFloatRet2_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with FLOAD (0x17) which calls PushFloatRet2
	code := []byte{opcodes.FLOAD, 0x01} // FLOAD with local variable index 1

	// Create basic constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// FSTORE
func TestStore_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with FSTORE (0x38) which calls storeFloatRet2
	code := []byte{opcodes.FSTORE, 0x01} // FSTORE with local variable index 1

	// Create basic constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// FSTORE_0
func TestStore0_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with FSTORE_0 (0x43) which calls storeFloat
	code := []byte{opcodes.FSTORE_0}

	// Create basic constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// GETFIELD
func TestCheckGetfield_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GETFIELD, 0x00, 0x01} // GETFIELD with CP index 1
	cp := createCPWithEntry(1, int(FieldRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// GETFIELD: valid field ref
func TestCheckGetfield_ValidFieldRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, FieldRef)
	CP = &cp
	Code = []byte{opcodes.GETFIELD, 0x00, 0x01}
	PC = 0

	result := CheckGetfield()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

// GETFIELD: invalid CP slot
func TestCheckGetfield_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GETFIELD, 0x00, 0xFF} // GETFIELD with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid GETFIELD CP slot, but got none")
	}
}

// GETFIELD: invalid field ref (points to method ref)
func TestCodeCheckGetfield(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	code := []byte{opcodes.GETFIELD, 0x00, 0x01} // GETFIELD pointing to slot 1
	CP := CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = CpEntry{Type: MethodRef, Slot: 0} // should be a field ref

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &CP, 5, af)
	if err == nil {
		t.Errorf("GETFIELD: Expected error but did not get one.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "java.lang.VerifyError") || !strings.Contains(errMsg, "not a field reference") {
		t.Errorf("GETFIELD: Did not get expected error message, got: %s", errMsg)
	}
}

// GETSTATIC
func TestCheckGetstatic_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GETSTATIC, 0x00, 0x01} // GETSTATIC with CP index 1
	cp := createCPWithEntry(1, int(FieldRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("GETSTATIC CheckCodeValidity failed: %v", err)
	}
}

// GETSTATIC: valid field ref
func TestCheckGetstatic_ValidFieldRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, FieldRef)
	CP = &cp
	Code = []byte{opcodes.GETSTATIC, 0x00, 0x01}
	PC = 0

	result := CheckGetfield()

	if result != 3 {
		t.Errorf("GETSTATIC: Expected return value 3, got: %d", result)
	}
}

// GETSTATIC: invalid CP slot
func TestCheckGetstatic_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GETSTATIC, 0x00, 0xFF} // GETSTATIC with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid GETSTATIC CP slot, but got none")
	}
}

// GETSTATIC: invalid field ref (points to method ref)
func TestCodeCheckGetStatic(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	code := []byte{opcodes.GETSTATIC, 0x00, 0x01} // GETSTATIC pointing to slot 1
	CP := CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = CpEntry{Type: MethodRef, Slot: 0} // should be a field ref

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &CP, 5, af)
	if err == nil {
		t.Errorf("GETSTATIC: Expected error but did not get one.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "java.lang.VerifyError") || !strings.Contains(errMsg, "not a field reference") {
		t.Errorf("GETSTATIC: Did not get expected error message, got: %s", errMsg)
	}
}

// GOTO
func TestCheckGoto_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GOTO, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00} // GOTO +3
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// GOTO
func TestCheckGoto_InvalidJumpNegative(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GOTO, 0xFF, 0xFE} // GOTO -2 (invalid)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid GOTO jump, but got none")
	}
}

// GOTO
func TestCheckGoto_InvalidJumpOutOfBounds(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.GOTO, 0x00, 0x10} // GOTO +16 (out of bounds)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for out-of-bounds GOTO jump, but got none")
	}
}

// GOTO_W valid jump
func TestCheckGotow_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with GOTO_W (0xC8) with valid 4-byte offset
	code := []byte{0xC8, 0x00, 0x00, 0x00, 0x05, 0x00} // GOTO_W, jump +5 bytes forward

	// Create basic constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// GOTO_W: invalid jump
func TestCheckGotow_InvalidJumpFOrward(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with GOTO_W (0xC8) with valid 4-byte offset
	code := []byte{opcodes.GOTO_W, 0x00, 0x00, 0x00, 0x05} // GOTO_W, jump +5 bytes forward, which is outside the code

	// Create basic constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	// Redirect stderr to a pipe to avoid printing to console
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := CheckCodeValidity(&code, &cp, 10, af)

	_ = w.Close()
	os.Stderr = normalStderr

	if err == nil {
		t.Errorf("Expected an error because forward jump is too large")
	}
}

// GOTO_W with invalid jump
func TestCheckGotow_InvalidJumpNegative(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with GOTO_W (0xC8) with invalid negative jump beyond start
	code := []byte{0x00, opcodes.GOTO_W, 0xFF, 0xFF, 0xFF, 0xFE} // NOP then GOTO_W, jump beyond start

	// Create basic constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	// Redirect stderr to a pipe to avoid printing to console
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := CheckCodeValidity(&code, &cp, 10, af)

	_ = w.Close()
	os.Stderr = normalStderr

	if err == nil {
		t.Errorf("Expected CheckCodeValidity to fail for invalid GOTO_W jump")
	}
}

// ICONST_0
func TestIconst0_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x03} // ICONST_0
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// IF_ACMPEQ
func TestCheckIfAcmpeq_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.IF_ACMPEQ, 0x00, 0x03, 0x00, 0x00, 0x00} // IF_ACMPEQ +3
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// IFEQ
func TestCheckIfeq_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.IFEQ, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00} // IFEQ +3
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// INVOKEINTERFACE
func TestCheckInvokeinterface_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x02, 0x00} // INVOKEINTERFACE with CP index 1, count 2, zero
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// INVOKEINTERFACE: valid interface
func TestCheckInvokeinterface_ValidInterface(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, Interface)
	CP = &cp
	Code = []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x02, 0x00} // count=2, zero=0
	PC = 0

	result := CheckInvokeinterface()

	if result != 4 {
		t.Errorf("Expected return value 4, got: %d", result)
	}
}

// INVOKEINTERFACE: invalid interface
func TestCheckInvokeinterface_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.INVOKEINTERFACE, 0x00, 0xFF, 0x02, 0x00} // INVOKEINTERFACE with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid INVOKEINTERFACE CP slot, but got none")
	}
}

// INVOKEINTERFACE: invalid count byte
func TestCheckInvokeinterface_ZeroCountByte(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x00, 0x00} // INVOKEINTERFACE with count 0
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for zero count byte in INVOKEINTERFACE, but got none")
	}
}

// INVOKEINTERFACE: invalid zero byte
func TestCheckInvokeinterface_NonZeroZeroByte(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x02, 0x01} // INVOKEINTERFACE with non-zero zero byte
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for non-zero zero byte in INVOKEINTERFACE, but got none")
	}
}

// INVOKESPECIAL bytecode with valid method ref
func TestCheckInvokespecial_ValidMethodRef(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESPECIAL (0xB7)
	code := []byte{opcodes.INVOKESPECIAL, 0x00, 0x01} // INVOKESPECIAL with CP index 1

	// Create constant pool with valid method ref
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	cp.CpIndex[1] = CpEntry{Type: MethodRef, Slot: 0}

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// INVOKESPECIAL bytecode with invalid CP slot
func TestCheckInvokespecial_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESPECIAL (0xB7) with invalid CP index
	code := []byte{opcodes.INVOKESPECIAL, 0x00, 0xFF} // INVOKESPECIAL with invalid CP index

	// Create small constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	// Redirect stderr to a pipe to avoid printing to console
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := CheckCodeValidity(&code, &cp, 10, af)

	_ = w.Close()
	os.Stderr = normalStderr

	if err == nil {

		t.Errorf("Expected CheckCodeValidity to fail for invalid INVOKESPECIAL CP slot")
	}
}

// INVOKESTATIC with valid method ref
func TestCheckInvokestatic_ValidMethodRef(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESTATIC (0xB8)
	code := []byte{opcodes.INVOKESTATIC, 0x00, 0x01} // INVOKESTATIC with CP index 1

	// Create constant pool with valid method ref
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	cp.CpIndex[1] = CpEntry{Type: MethodRef, Slot: 0}

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// INVOKESTATIC bytecode with invalid CP slot
func TestCheckInvokestatic_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESTATIC (0xB8) with invalid CP index
	code := []byte{opcodes.INVOKESTATIC, 0x00, 0xFF} // INVOKESTATIC with invalid CP index

	// Create small constant pool
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}

	af := AccessFlags{}

	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	err := CheckCodeValidity(&code, &cp, 10, af)

	_ = w.Close()
	os.Stderr = normalStderr

	if err == nil {
		t.Errorf("Expected CheckCodeValidity to fail for invalid INVOKESTATIC CP slot")
	}
}

// INVOKEVIRTUAL
func TestCheckInvokevirtual_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.INVOKEVIRTUAL, 0x00, 0x01} // INVOKEVIRTUAL with CP index 1
	cp := createCPWithEntry(1, int(MethodRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// INVOKEVIRTUAL bytecode with invalid CP slot
func TestCheckInvokevirtual_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.INVOKEVIRTUAL, 0x00, 0xFF} // INVOKEVIRTUAL with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid INVOKEVIRTUAL CP slot, but got none")
	}
}

// INVOKEVIRTUAL bytecode with valid method ref
func TestCheckInvokevirtual_ValidMethodRef(t *testing.T) {}
func TestNewInvokevirtualInvalidMethRef(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	code := []byte{opcodes.INVOKEVIRTUAL, 0x00, 0x01} // INVOKEVIRTUAL pointing to slot 1
	CP := CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = CpEntry{Type: ClassRef, Slot: 0} // should be a method ref
	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = ResolvedFieldEntry{
		ClName:  "testClass",
		FldName: "testField",
		FldType: "I",
	}

	af := AccessFlags{}
	err := CheckCodeValidity(&code, &CP, 5, af)
	if err == nil {
		t.Errorf("INVOKEVIRTUAL: Expected error but did not get one.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("INVOKEVIRTUAL: Expected error but did not get one.")
	}

	if !strings.Contains(errMsg, "java.lang.VerifyError") || !strings.Contains(errMsg, "not a method reference") {
		t.Errorf("INVOKEVIRTUAL: Did not get expected error message, got:\n %s", errMsg)
	}
}

// ISTORE
func TestIstore_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.ISTORE, 0x01} // ISTORE with index
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// ISTORE_0
func TestIstore0_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.ISTORE_0}
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// LDC_W
func TestPushIntRet3_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with LDC_W (0x13) which calls PushIntRet3
	code := []byte{opcodes.LDC_W, 0x00, 0x01} // LDC_W with CP index 1

	// Create basic constant pool with an integer constant
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	cp.CpIndex[1] = CpEntry{Type: IntConst, Slot: 0}
	cp.IntConsts = make([]int32, 1)
	cp.IntConsts[0] = 42

	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// LOOKUPSWITCH with zero padding (PC+1 is already 4-byte aligned)
func TestLookupswitch_ZeroPadding(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH at position 3 (PC+1=4, which is 4-byte aligned, so no padding needed)
	code := []byte{
		0x00, 0x00, 0x00, // 3 NOPs to align LOOKUPSWITCH at position 3
		opcodes.LOOKUPSWITCH,   // 0xAB - LOOKUPSWITCH at PC=3
		0x00, 0x00, 0x00, 0x10, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x02, // npairs = 2 (4 bytes)
		0x00, 0x00, 0x00, 0x01, // match1 = 1 (4 bytes)
		0x00, 0x00, 0x00, 0x20, // offset1 = 32 (4 bytes)
		0x00, 0x00, 0x00, 0x05, // match2 = 5 (4 bytes)
		0x00, 0x00, 0x00, 0x30, // offset2 = 48 (4 bytes)
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for zero padding case: %v", err)
	}
}

// LOOKUPSWITCH with 1-byte padding
func TestLookupswitch_OneBytePadding(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH at position 0 (PC+1=1, needs 3 padding bytes to align to 4)
	code := []byte{
		opcodes.LOOKUPSWITCH, // 0xAB - LOOKUPSWITCH at PC=0
		0x00, 0x00, 0x00,     // 3 padding bytes
		0x00, 0x00, 0x00, 0x18, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x01, // npairs = 1 (4 bytes)
		0x00, 0x00, 0x00, 0x42, // match1 = 66 (4 bytes)
		0x00, 0x00, 0x00, 0x28, // offset1 = 40 (4 bytes)
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for one byte padding case: %v", err)
	}
}

// LOOKUPSWITCH with 2-byte padding
func TestLookupswitch_TwoBytePadding(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH at position 1 (PC+1=2, needs 2 padding bytes to align to 4)
	code := []byte{
		0x00,                 // NOP to position LOOKUPSWITCH at PC=1
		opcodes.LOOKUPSWITCH, // 0xAB - LOOKUPSWITCH at PC=1
		0x00, 0x00,           // 2 padding bytes
		0x00, 0x00, 0x00, 0x14, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x01, // npairs = 1 (4 bytes)
		0x00, 0x00, 0x00, 0x7B, // match1 = 123 (4 bytes)
		0x00, 0x00, 0x00, 0x24, // offset1 = 36 (4 bytes)
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for two byte padding case: %v", err)
	}
}

// LOOKUPSWITCH with 3-byte padding
func TestLookupswitch_ThreeBytePadding(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH at position 2 (PC+1=3, needs 1 padding byte to align to 4)
	code := []byte{
		0x00, 0x00, // 2 NOPs to position LOOKUPSWITCH at PC=2
		opcodes.LOOKUPSWITCH,   // 0xAB - LOOKUPSWITCH at PC=2
		0x00,                   // 1 padding byte
		0x00, 0x00, 0x00, 0x10, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x01, // npairs = 1 (4 bytes)
		0x00, 0x00, 0x00, 0x99, // match1 = 153 (4 bytes)
		0x00, 0x00, 0x00, 0x20, // offset1 = 32 (4 bytes)
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for three byte padding case: %v", err)
	}
}

// LOOKUPSWITCH with zero pairs (npairs = 0)
func TestLookupswitch_ZeroPairs(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH with 0 pairs
	code := []byte{
		0x00, 0x00, 0x00, // 3 NOPs to align LOOKUPSWITCH at position 3
		opcodes.LOOKUPSWITCH,   // 0xAB - LOOKUPSWITCH at PC=3
		0x00, 0x00, 0x00, 0x0C, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x00, // npairs = 0 (4 bytes)
		// No match/offset pairs since npairs = 0
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for zero pairs case: %v", err)
	}
}

// LOOKUPSWITCH with multiple pairs (npairs > 1)
func TestLookupswitch_MultiplePairs(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH with 3 pairs
	code := []byte{
		0x00, 0x00, 0x00, // 3 NOPs to align LOOKUPSWITCH at position 3
		opcodes.LOOKUPSWITCH,   // 0xAB - LOOKUPSWITCH at PC=3
		0x00, 0x00, 0x00, 0x24, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x03, // npairs = 3 (4 bytes)
		0x00, 0x00, 0x00, 0x0A, // match1 = 10 (4 bytes)
		0x00, 0x00, 0x00, 0x30, // offset1 = 48 (4 bytes)
		0x00, 0x00, 0x00, 0x14, // match2 = 20 (4 bytes)
		0x00, 0x00, 0x00, 0x38, // offset2 = 56 (4 bytes)
		0x00, 0x00, 0x00, 0x1E, // match3 = 30 (4 bytes)
		0x00, 0x00, 0x00, 0x40, // offset3 = 64 (4 bytes)
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for multiple pairs case: %v", err)
	}
}

// LOOKUPSWITCH with large npairs value
func TestLookupswitch_LargePairs(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH with 255 pairs (testing larger npairs)
	npairs := 255

	code := []byte{
		0x00, 0x00, 0x00, // 3 NOPs to align LOOKUPSWITCH at position 3
		opcodes.LOOKUPSWITCH,   // 0xAB - LOOKUPSWITCH at PC=3
		0x00, 0x00, 0x01, 0x08, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0xFF, // npairs = 255 (4 bytes)
	}

	// Add 255 match/offset pairs
	for i := 0; i < npairs; i++ {
		matchValue := uint32(i + 1)
		offsetValue := uint32(0x100 + i*4)

		code = append(code, byte(matchValue>>24), byte(matchValue>>16), byte(matchValue>>8), byte(matchValue))
		code = append(code, byte(offsetValue>>24), byte(offsetValue>>16), byte(offsetValue>>8), byte(offsetValue))
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed for large pairs case: %v", err)
	}
}

// Lookupswitch edge case where the calculated length exceeds available code
func TestLookupswitch_CalculatedLengthExceedsCode(t *testing.T) {
	globals.InitGlobals("test")

	// LOOKUPSWITCH claiming to have more pairs than the code contains
	code := []byte{
		0x00, 0x00, 0x00, // 3 NOPs to align LOOKUPSWITCH at position 3
		opcodes.LOOKUPSWITCH,   // 0xAB - LOOKUPSWITCH at PC=3
		0x00, 0x00, 0x00, 0x10, // default offset (4 bytes)
		0x00, 0x00, 0x00, 0x10, // npairs = 16 (but we don't provide 16 pairs)
		0x00, 0x00, 0x00, 0x01, // match1 = 1 (4 bytes)
		0x00, 0x00, 0x00, 0x20, // offset1 = 32 (4 bytes)
		// Missing 15 more pairs
	}

	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err == nil {
		t.Errorf("Expected CheckCodeValidity to fail when calculated length exceeds code, but got no error")
	}
}

// MULTIANEWARRAY
func TestCheckMultianewarray_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.MULTIANEWARRAY, 0x00, 0x01, 0x02} // MULTIANEWARRAY with CP index 1, dimensions 2
	cp := createCPWithEntry(1, int(ClassRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// MULTIANEWARRAY
func TestCheckMultianewarray_ValidClassRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, ClassRef)
	CP = &cp
	Code = []byte{opcodes.MULTIANEWARRAY, 0x00, 0x01, 0x02} // dimensions = 2
	PC = 0

	result := CheckMultianewarray()

	if result != 4 {
		t.Errorf("Expected return value 4, got: %d", result)
	}
}

// MULTIANEWARRAY
func TestCheckMultianewarray_ZeroDimensions(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.MULTIANEWARRAY, 0x00, 0x01, 0x00} // MULTIANEWARRAY with 0 dimensions
	cp := createCPWithEntry(1, int(ClassRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for zero dimensions in MULTIANEWARRAY, but got none")
	}
}

// NOP
func TestNop(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.NOP} // NOP
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// POP
func TestCheckPop_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.POP} // POP
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// POP
func TestCheckPop_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	StackEntries = 3
	result := CheckPop()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 2 {
		t.Errorf("Expected StackEntries to decrease by 1, got: %d", StackEntries)
	}
}

// POP2
func TestCheckPop2_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.POP2}
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckPop2_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	StackEntries = 3
	result := CheckPop2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if StackEntries != 1 {
		t.Errorf("Expected StackEntries to decrease by 2, got: %d", StackEntries)
	}
}

// PUTSTATIC
func TestCheckPutstatic_HighLevel(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.PUTSTATIC, 0x00, 0x01} // PUTFIELD with CP index 1
	cp := createCPWithEntry(1, int(FieldRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {

		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// PUTSTATIC with valid field ref
func TestCheckPutstatic_ValidFieldRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, FieldRef)
	CP = &cp
	Code = []byte{opcodes.PUTSTATIC, 0x00, 0x01}
	PC = 0

	result := CheckGetfield()

	if result != 3 {
		t.Errorf("Putstatic expected return value 3, got: %d", result)
	}
}

// PUTSTATIC with invalic CP slot
func TestCheckPutstatic_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.PUTSTATIC, 0x00, 0xFF} // PUTSTATIC with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid PUTSTATIC CP slot, but got none")
	}
}

// PUTSTATIC with invalid field ref
func TestCodeCheckPutfield(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	code := []byte{opcodes.PUTSTATIC, 0x00, 0x01} // PUTSTATIC pointing to slot 1
	CP = &CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = CpEntry{Type: MethodRef, Slot: 0} // should be a field ref

	af := AccessFlags{}

	err := CheckCodeValidity(&code, CP, 5, af)
	if err == nil {
		t.Errorf("PUTSTATIC: Expected error but did not get one.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "java.lang.VerifyError") || !strings.Contains(errMsg, "not a field reference") {
		t.Errorf("PUTSTATIC: Did not get expected error message, got: %s", errMsg)
	}
}

// SIPUSH
func TestCheckSipush_ValidLength(t *testing.T) {
	globals.InitGlobals("test")

	Code = []byte{opcodes.SIPUSH, 0x12, 0x34}
	PC = 0
	StackEntries = 0

	result := CheckSipush()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
	if StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", StackEntries)
	}
}

func TestCheckSipush_ValidLength2(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x11, 0x01, 0x00} // SIPUSH 256
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckSipush_InsufficientLength(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x11, 0x01} // SIPUSH with missing byte
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for insufficient SIPUSH length, but got none")
	}
}

// TABLESWITCH
func TestCheckTableswitch_ValidRange(t *testing.T) {
	globals.InitGlobals("test")

	// TABLESWITCH with proper padding and valid range
	code := []byte{
		0xAA,             // TABLESWITCH
		0x00, 0x00, 0x00, // padding to 4-byte boundary
		0x00, 0x00, 0x00, 0x0A, // default offset
		0x00, 0x00, 0x00, 0x01, // low = 1
		0x00, 0x00, 0x00, 0x03, // high = 3
		0x00, 0x00, 0x00, 0x14, // offset for 1
		0x00, 0x00, 0x00, 0x18, // offset for 2
		0x00, 0x00, 0x00, 0x1C, // offset for 3
	}
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckTableswitch_InvalidRange(t *testing.T) {
	globals.InitGlobals("test")

	// TABLESWITCH with invalid range (low > high)
	code := []byte{
		0xAA,             // TABLESWITCH
		0x00, 0x00, 0x00, // padding to 4-byte boundary
		0x00, 0x00, 0x00, 0x0A, // default offset
		0x00, 0x00, 0x00, 0x05, // low = 5
		0x00, 0x00, 0x00, 0x03, // high = 3 (invalid: low > high)
	}
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid TABLESWITCH range, but got none")
	}
}
