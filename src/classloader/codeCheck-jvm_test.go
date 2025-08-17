/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

// Migrated tests from jvm/codeCheck_test.go for use in classloader package
package classloader

import (
	"io"
	"jacobin/globals"
	"jacobin/opcodes"
	"os"
	"strings"
	"testing"
)

// Helper function to redirect stderr and capture error messages
func captureStderr(t *testing.T, testFunc func()) string {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	testFunc()

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	return string(msg)
}

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

// ==================== CheckCodeValidity Main Function Tests ====================

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

// ==================== Individual function tests in alphabetical order of the instructions ========

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

func TestArith_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x60} // IADD
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestPushFloat_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x0B} // FCONST_0
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestPushInt_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x03} // ICONST_0
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// BIPUSH
func TestCheckBipush_ValidLength(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.BIPUSH, 0x42}
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckBipush_ValidLength2(t *testing.T) {
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

	code := []byte{0x10} // BIPUSH with missing byte
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for insufficient BIPUSH length, but got none")
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

func TestDup1_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x59} // DUP
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestDup2_LongDoubleOperation(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x09, 0x5C, 0x00} // LCONST_0, DUP2, NOP (extra byte for DUP2 to check next instruction)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestDup2_RegularOperation(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x03, 0x5C} // ICONST_0, DUP2
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckPop_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x57} // POP
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckPop2_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x58} // POP2
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckGetfield_ValidFieldRef(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB4, 0x00, 0x01} // GETFIELD with CP index 1
	cp := createCPWithEntry(1, int(FieldRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckGetfield_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB4, 0x00, 0xFF} // GETFIELD with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid GETFIELD CP slot, but got none")
	}
}

func TestCheckGoto_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xA7, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00} // GOTO +3
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckGoto_InvalidJumpNegative(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xA7, 0xFF, 0xFE} // GOTO -2 (invalid)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid GOTO jump, but got none")
	}
}

func TestCheckGoto_InvalidJumpOutOfBounds(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xA7, 0x00, 0x10} // GOTO +16 (out of bounds)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for out-of-bounds GOTO jump, but got none")
	}
}

func TestCheckIf_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xA5, 0x00, 0x03, 0x00, 0x00, 0x00} // IF_ACMPEQ +3
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckIfZero_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x99, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00} // IFEQ +3
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckInvokeinterface_ValidInterface(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB9, 0x00, 0x01, 0x02, 0x00} // INVOKEINTERFACE with CP index 1, count 2, zero
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckInvokeinterface_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB9, 0x00, 0xFF, 0x02, 0x00} // INVOKEINTERFACE with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid INVOKEINTERFACE CP slot, but got none")
	}
}

func TestCheckInvokeinterface_ZeroCountByte(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB9, 0x00, 0x01, 0x00, 0x00} // INVOKEINTERFACE with count 0
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for zero count byte in INVOKEINTERFACE, but got none")
	}
}

func TestCheckInvokeinterface_NonZeroZeroByte(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB9, 0x00, 0x01, 0x02, 0x01} // INVOKEINTERFACE with non-zero zero byte
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for non-zero zero byte in INVOKEINTERFACE, but got none")
	}
}

func TestCheckInvokevirtual_ValidMethodRef(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB6, 0x00, 0x01} // INVOKEVIRTUAL with CP index 1
	cp := createCPWithEntry(1, int(MethodRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckInvokevirtual_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB6, 0x00, 0xFF} // INVOKEVIRTUAL with invalid CP index
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid INVOKEVIRTUAL CP slot, but got none")
	}
}

func TestReturn1(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x00} // NOP (uses Return1)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestReturn2(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x15, 0x01} // ILOAD with index (uses Return2)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestReturn3(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x13, 0x00, 0x01} // LDC_W (uses Return3)
	cp := createCPWithEntry(1, int(IntConst))
	cp.IntConsts = make([]int32, 1)
	cp.IntConsts[0] = 42
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestReturn4(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xB9, 0x00, 0x01, 0x02, 0x00} // INVOKEINTERFACE (uses Return4)
	cp := createCPWithEntry(1, int(Interface))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestReturn5(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xC8, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00} // GOTO_W (uses Return5)
	cp := createBasicCP()
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
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

func TestCheckMultianewarray_ValidClassRef(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xC5, 0x00, 0x01, 0x02} // MULTIANEWARRAY with CP index 1, dimensions 2
	cp := createCPWithEntry(1, int(ClassRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

func TestCheckMultianewarray_ZeroDimensions(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0xC5, 0x00, 0x01, 0x00} // MULTIANEWARRAY with 0 dimensions
	cp := createCPWithEntry(1, int(ClassRef))
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for zero dimensions in MULTIANEWARRAY, but got none")
	}
}

// Original tests that use newFrame need to be converted

func TestCodeCheckForGetfield2(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	code := []byte{0xB4, 0x00, 0x01} // GETFIELD pointing to slot 1
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

func TestNewInvokevirtualInvalidMethRef(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	code := []byte{0xB6, 0x00, 0x01} // INVOKEVIRTUAL pointing to slot 1
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
