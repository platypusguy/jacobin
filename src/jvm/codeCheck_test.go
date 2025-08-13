/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"io"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/opcodes"
	"os"
	"strings"
	"testing"
)

// GETFIELD : get field from object -- here testing for CP not pointing to a field ref
func TestCodeCheckForGetfield(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0} // should be a field ref

	f.CP = &CP

	af := classloader.AccessFlags{}

	err := classloader.CheckCodeValidity(&f.Meth, f.CP.(*classloader.CPool), 5, af)
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

// INVOKEVIRTUAL : invoke method -- here testing for error
func TestNewInvokevirtualInvalidMethRef2(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0} // should be a method ref
	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		ClName:  "testClass",
		FldName: "testField",
		FldType: "I",
	}

	f.CP = &CP

	af := classloader.AccessFlags{}
	err := classloader.CheckCodeValidity(&f.Meth, f.CP.(*classloader.CPool), 5, af)
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

// === JUnie-generated tests ===
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
func createBasicCP() classloader.CPool {
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	return CP
}

// Helper function to create a constant pool with specific entry types
func createCPWithEntry(index int, entryType int) classloader.CPool {
	CP := createBasicCP()
	if index < len(CP.CpIndex) {
		CP.CpIndex[index] = classloader.CpEntry{Type: uint16(entryType), Slot: 0}
	}
	return CP
}

// ==================== CheckCodeValidity Main Function Tests ====================

func TestCheckCodeValidity_NilCodePointer(t *testing.T) {
	globals.InitGlobals("test")

	cp := createBasicCP()
	af := classloader.AccessFlags{}

	err := classloader.CheckCodeValidity(nil, &cp, 5, af)
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
	af := classloader.AccessFlags{ClassIsAbstract: false}

	err := classloader.CheckCodeValidity(&code, &cp, 5, af)
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
	af := classloader.AccessFlags{ClassIsAbstract: true}

	err := classloader.CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("Expected no error for empty code in abstract class, but got: %s", err.Error())
	}
}

func TestCheckCodeValidity_NilConstantPool(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.NOP}

	err := classloader.CheckCodeValidity(&code, nil, 5, classloader.AccessFlags{})
	if err == nil {
		t.Errorf("Expected error for nil constant pool, but got none")
	}
	if !strings.Contains(err.Error(), "ptr to constant pool is nil") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestCheckCodeValidity_EmptyConstantPool(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.NOP}
	cp := classloader.CPool{} // empty CP
	af := classloader.AccessFlags{}

	err := classloader.CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for empty constant pool, but got none")
	}
	if !strings.Contains(err.Error(), "empty constant pool") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

func TestCheckCodeValidity_ValidCode(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{opcodes.NOP, opcodes.ACONST_NULL, opcodes.RETURN}
	cp := createBasicCP()
	af := classloader.AccessFlags{}

	err := classloader.CheckCodeValidity(&code, &cp, 5, af)
	if err != nil {
		t.Errorf("Expected no error for valid code, but got: %s", err.Error())
	}
}

func TestCheckCodeValidity_InvalidBytecodeLength(t *testing.T) {
	globals.InitGlobals("test")

	// BIPUSH requires 1 additional byte but we don't provide it
	code := []byte{opcodes.BIPUSH}
	cp := createBasicCP()
	af := classloader.AccessFlags{}

	err := classloader.CheckCodeValidity(&code, &cp, 5, af)
	if err == nil {
		t.Errorf("Expected error for invalid bytecode length, but got none")
	}
	if !strings.Contains(err.Error(), "Invalid bytecode or argument") {
		t.Errorf("Expected specific error message, got: %s", err.Error())
	}
}

// ==================== Arithmetic Operations Tests ====================

func TestArith_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	// Set up global variables that Arith() uses
	classloader.StackEntries = 5

	result := classloader.Arith() // We'll need to expose this function for testing

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 4 {
		t.Errorf("Expected StackEntries to decrease by 1, got: %d", classloader.StackEntries)
	}
}

// ==================== Constant Loading Operations Tests ====================

func TestCheckAconstnull_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	classloader.StackEntries = 0
	result := classloader.TestCheckAconstnull()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
}

func TestPushFloat_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	classloader.StackEntries = 0
	result := classloader.TestPushFloat()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
}

func TestPushInt_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	classloader.StackEntries = 0
	result := classloader.TestPushInt()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
}

// ==================== Bytecode with Arguments Tests ====================

func TestCheckBipush_ValidLength(t *testing.T) {
	globals.InitGlobals("test")

	classloader.Code = []byte{opcodes.BIPUSH, 0x42}
	classloader.PC = 0
	classloader.StackEntries = 0

	result := classloader.TestCheckBipush()

	if result != 2 {
		t.Errorf("Expected return value 2, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
}

func TestCheckBipush_InsufficientLength(t *testing.T) {
	globals.InitGlobals("test")

	classloader.Code = []byte{opcodes.BIPUSH} // Missing the required byte argument
	classloader.PC = 0
	classloader.StackEntries = 0

	result := classloader.TestCheckBipush()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED, got: %d", result)
	}
}

func TestCheckSipush_ValidLength(t *testing.T) {
	globals.InitGlobals("test")

	classloader.Code = []byte{opcodes.SIPUSH, 0x12, 0x34}
	classloader.PC = 0
	classloader.StackEntries = 0

	result := classloader.TestCheckSipush()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
}

func TestCheckSipush_InsufficientLength(t *testing.T) {
	globals.InitGlobals("test")

	classloader.Code = []byte{opcodes.SIPUSH, 0x12} // Missing second byte
	classloader.PC = 0
	classloader.StackEntries = 0

	result := classloader.TestCheckSipush()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED, got: %d", result)
	}
}

// ==================== Stack Manipulation Tests ====================

func TestDup1_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	classloader.StackEntries = 1
	result := classloader.CheckDup1()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
}

func TestDup2_LongDoubleOperation(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is for long/double
	classloader.Code = []byte{opcodes.DUP2, opcodes.LADD}
	classloader.PC = 0
	classloader.StackEntries = 1

	result := classloader.CheckDup2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 2 {
		t.Errorf("Expected StackEntries to increase by 1, got: %d", classloader.StackEntries)
	}
	// Check that DUP2 was converted to DUP
	if classloader.Code[0] != opcodes.DUP {
		t.Errorf("Expected DUP2 to be converted to DUP, got: 0x%x", classloader.Code[0])
	}
}

func TestDup2_RegularOperation(t *testing.T) {
	globals.InitGlobals("test")

	// Create code where next bytecode is NOT for long/double
	classloader.Code = []byte{opcodes.DUP2, opcodes.IADD}
	classloader.PC = 0
	classloader.StackEntries = 2

	result := classloader.CheckDup2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 4 {
		t.Errorf("Expected StackEntries to increase by 2, got: %d", classloader.StackEntries)
	}
}

func TestCheckPop_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	classloader.StackEntries = 3
	result := classloader.CheckPop()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 2 {
		t.Errorf("Expected StackEntries to decrease by 1, got: %d", classloader.StackEntries)
	}
}

func TestCheckPop2_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	classloader.StackEntries = 3
	result := classloader.CheckPop2()

	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to decrease by 2, got: %d", classloader.StackEntries)
	}
}

// ==================== Field Access Operations Tests ====================

func TestCheckGetfield_ValidFieldRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.FieldRef)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.GETFIELD, 0x00, 0x01}
	classloader.PC = 0

	result := classloader.CheckGetfield()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

func TestCheckGetfield_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	cp := createBasicCP()
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.GETFIELD, 0xFF, 0xFF} // Invalid CP slot
	classloader.PC = 0

	result := classloader.CheckGetfield()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED, got: %d", result)
	}
}

// ==================== Control Flow Operations Tests ====================

func TestCheckGoto_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	// Create a code array with valid jump target
	classloader.Code = []byte{opcodes.GOTO, 0x00, 0x02, opcodes.NOP, opcodes.NOP, opcodes.NOP}
	classloader.PC = 0

	result := classloader.CheckGoto()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

func TestCheckGoto_InvalidJumpNegative(t *testing.T) {
	globals.InitGlobals("test")

	errMsg := captureStderr(t, func() {
		classloader.Code = []byte{opcodes.GOTO, 0xFF, 0xFE} // Jump to negative location
		classloader.PC = 0
		classloader.CheckGoto()
	})

	if !strings.Contains(errMsg, "GOTO") || !strings.Contains(errMsg, "illegal jump") {
		t.Errorf("Expected GOTO error message, got: %s", errMsg)
	}
}

func TestCheckGoto_InvalidJumpOutOfBounds(t *testing.T) {
	globals.InitGlobals("test")

	errMsg := captureStderr(t, func() {
		classloader.Code = []byte{opcodes.GOTO, 0x00, 0xFF} // Jump beyond code bounds
		classloader.PC = 0
		classloader.CheckGoto()
	})

	if !strings.Contains(errMsg, "GOTO") || !strings.Contains(errMsg, "illegal jump") {
		t.Errorf("Expected GOTO error message, got: %s", errMsg)
	}
}

// ==================== Conditional Operations Tests ====================

func TestCheckIf_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	classloader.Code = []byte{opcodes.IF_ICMPEQ, 0x00, 0x02, opcodes.NOP, opcodes.NOP, opcodes.NOP}
	classloader.PC = 0

	result := classloader.CheckIf()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

func TestCheckIfZero_ValidJump(t *testing.T) {
	globals.InitGlobals("test")

	classloader.Code = []byte{opcodes.IFEQ, 0x00, 0x02, opcodes.NOP, opcodes.NOP, opcodes.NOP}
	classloader.PC = 0
	classloader.StackEntries = 2

	result := classloader.CheckIfzero()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
	if classloader.StackEntries != 1 {
		t.Errorf("Expected StackEntries to decrease by 1, got: %d", classloader.StackEntries)
	}
}

// ==================== Method Invocation Operations Tests ====================

func TestCheckInvokeinterface_ValidInterface(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.Interface)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x02, 0x00} // count=2, zero=0
	classloader.PC = 0

	result := classloader.CheckInvokeinterface()

	if result != 4 {
		t.Errorf("Expected return value 4, got: %d", result)
	}
}

func TestCheckInvokeinterface_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	cp := createBasicCP()
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.INVOKEINTERFACE, 0xFF, 0xFF, 0x02, 0x00}
	classloader.PC = 0

	result := classloader.CheckInvokeinterface()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED, got: %d", result)
	}
}

func TestCheckInvokeinterface_ZeroCountByte(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.Interface)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x00, 0x00} // count=0
	classloader.PC = 0

	result := classloader.CheckInvokeinterface()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED for zero count byte, got: %d", result)
	}
}

func TestCheckInvokeinterface_NonZeroZeroByte(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.Interface)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.INVOKEINTERFACE, 0x00, 0x01, 0x02, 0x01} // zero byte != 0
	classloader.PC = 0

	result := classloader.CheckInvokeinterface()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED for non-zero zero byte, got: %d", result)
	}
}

func TestCheckInvokevirtual_ValidMethodRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.MethodRef)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.INVOKEVIRTUAL, 0x00, 0x01}
	classloader.PC = 0

	result := classloader.CheckInvokevirtual()

	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

func TestCheckInvokevirtual_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	cp := createBasicCP()
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.INVOKEVIRTUAL, 0xFF, 0xFF}
	classloader.PC = 0

	result := classloader.CheckInvokevirtual()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED, got: %d", result)
	}
}

// ==================== Utility Functions Tests ====================

func TestReturn1(t *testing.T) {
	result := classloader.Return1()
	if result != 1 {
		t.Errorf("Expected return value 1, got: %d", result)
	}
}

func TestReturn2(t *testing.T) {
	result := classloader.Return2()
	if result != 2 {
		t.Errorf("Expected return value 2, got: %d", result)
	}
}

func TestReturn3(t *testing.T) {
	result := classloader.Return3()
	if result != 3 {
		t.Errorf("Expected return value 3, got: %d", result)
	}
}

func TestReturn4(t *testing.T) {
	result := classloader.Return4()
	if result != 4 {
		t.Errorf("Expected return value 4, got: %d", result)
	}
}

func TestReturn5(t *testing.T) {
	result := classloader.Return5()
	if result != 5 {
		t.Errorf("Expected return value 5, got: %d", result)
	}
}

// ==================== ByteCodeIsForLongOrDouble Tests ====================

func TestByteCodeIsForLongOrDouble_LongDoubleCodes(t *testing.T) {
	globals.InitGlobals("test")

	longDoubleCodes := []byte{
		0x09, 0x0A, 0x0E, 0x0F, 0x14, 0x16, 0x1E, 0x1F,
		0x20, 0x21, 0x26, 0x27, 0x28, 0x29, 0x3F, 0x40,
		0x41, 0x42, 0x47, 0x48, 0x49, 0x4A, 0x63, 0x65,
		0x67, 0x69, 0x6B, 0x6D, 0x6F, 0x71, 0x73, 0x75,
		0x77, 0x79, 0x7B, 0x7D, 0x7F, 0x81, 0x83, 0x90,
		0x94, 0x98, // LCMP, DCMPG
	}

	for _, code := range longDoubleCodes {
		result := classloader.BytecodeIsForLongOrDouble(code)
		if !result {
			t.Errorf("Expected true for bytecode 0x%02x, got false", code)
		}
	}
}

func TestByteCodeIsForLongOrDouble_OtherCodes(t *testing.T) {
	globals.InitGlobals("test")

	otherCodes := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, // NOP, ACONST_NULL, ICONST_*
		0x60, 0x64, 0x68, 0x6C, 0x70, 0x74, 0x78, 0x7C, // IADD, ISUB, IMUL, IDIV, IREM, INEG, ISHL, IUSHR
	}

	for _, code := range otherCodes {
		result := classloader.BytecodeIsForLongOrDouble(code)
		if result {
			t.Errorf("Expected false for bytecode 0x%02x, got true", code)
		}
	}
}

// ==================== Additional Switch Operation Tests ====================

func TestCheckTableswitch_ValidRange(t *testing.T) {
	globals.InitGlobals("test")

	// Create a valid tableswitch bytecode sequence
	// tableswitch with padding, default, low=1, high=3, and 3 offsets
	classloader.Code = []byte{
		opcodes.TABLESWITCH,
		0x00, 0x00, 0x00, // padding to align to 4-byte boundary
		0x00, 0x00, 0x00, 0x10, // default offset
		0x00, 0x00, 0x00, 0x01, // low = 1
		0x00, 0x00, 0x00, 0x03, // high = 3
		0x00, 0x00, 0x00, 0x14, // offset for case 1
		0x00, 0x00, 0x00, 0x18, // offset for case 2
		0x00, 0x00, 0x00, 0x1C, // offset for case 3
	}
	classloader.PC = 0

	result := classloader.CheckTableSwitch()

	// Should return the total number of bytes consumed
	if result <= 0 {
		t.Errorf("Expected positive return value, got: %d", result)
	}
}

func TestCheckTableswitch_InvalidRange(t *testing.T) {
	globals.InitGlobals("test")

	// Create tableswitch with low > high (invalid)
	classloader.Code = []byte{
		opcodes.TABLESWITCH,
		0x00, 0x00, 0x00, // padding
		0x00, 0x00, 0x00, 0x10, // default
		0x00, 0x00, 0x00, 0x05, // low = 5
		0x00, 0x00, 0x00, 0x03, // high = 3 (invalid: low > high)
	}
	classloader.PC = 0

	result := classloader.CheckTableSwitch()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED for invalid range, got: %d", result)
	}
}

func TestCheckMultianewarray_ValidClassRef(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.ClassRef)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.MULTIANEWARRAY, 0x00, 0x01, 0x02} // dimensions = 2
	classloader.PC = 0

	result := classloader.CheckMultianewarray()

	if result != 4 {
		t.Errorf("Expected return value 4, got: %d", result)
	}
}

func TestCheckMultianewarray_ZeroDimensions(t *testing.T) {
	globals.InitGlobals("test")

	cp := createCPWithEntry(1, classloader.ClassRef)
	classloader.CP = &cp
	classloader.Code = []byte{opcodes.MULTIANEWARRAY, 0x00, 0x01, 0x00} // dimensions = 0
	classloader.PC = 0

	result := classloader.CheckMultianewarray()

	if result != classloader.ERROR_OCCURRED {
		t.Errorf("Expected ERROR_OCCURRED for zero dimensions, got: %d", result)
	}
}

// ==================== Existing Tests (Keep these as they were) ====================

// GETFIELD : get field from object -- here testing for CP not pointing to a field ref
func TestCodeCheckForGetfield2(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0} // should be a field ref

	f.CP = &CP

	af := classloader.AccessFlags{}

	err := classloader.CheckCodeValidity(&f.Meth, f.CP.(*classloader.CPool), 5, af)
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

// INVOKEVIRTUAL : invoke method -- here testing for error
func TestNewInvokevirtualInvalidMethRef(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0} // should be a method ref
	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		ClName:  "testClass",
		FldName: "testField",
		FldType: "I",
	}

	f.CP = &CP

	af := classloader.AccessFlags{}
	err := classloader.CheckCodeValidity(&f.Meth, f.CP.(*classloader.CPool), 5, af)
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
