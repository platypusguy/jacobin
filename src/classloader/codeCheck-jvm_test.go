/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

/*
 * Unit tests for functions in classloader/codeCheck.go that are not covered by jvm/codeCheck_test.go
 * These tests are placed here to test internal/unexported functions that can't be tested from jvm package
 */

package classloader

import (
	"jacobin/globals"
	"os"
	"testing"
)

// Test PushFloatRet2 via FLOAD bytecode (which uses PushFloatRet2)
func TestPushFloatRet2_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with FLOAD (0x17) which calls PushFloatRet2
	code := []byte{0x17, 0x01} // FLOAD with local variable index 1

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

// Test storeFloat via FSTORE_0 bytecode (which uses storeFloat)
func TestStoreFloat_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with FSTORE_0 (0x43) which calls storeFloat
	code := []byte{0x43} // FSTORE_0

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

// Test storeFloatRet2 via FSTORE bytecode (which uses storeFloatRet2)
func TestStoreFloatRet2_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with FSTORE (0x38) which calls storeFloatRet2
	code := []byte{0x38, 0x01} // FSTORE with local variable index 1

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

// Test PushIntRet3 via LDC_W bytecode (which uses PushIntRet3)
func TestPushIntRet3_StackIncrement(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with LDC_W (0x13) which calls PushIntRet3
	code := []byte{0x13, 0x00, 0x01} // LDC_W with CP index 1

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

// Test CheckGotow via GOTO_W bytecode with valid jump
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

// Test CheckGotow via GOTO_W bytecode with invalid jump
func TestCheckGotow_InvalidJumpFOrward(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with GOTO_W (0xC8) with valid 4-byte offset
	code := []byte{0xC8, 0x00, 0x00, 0x00, 0x05} // GOTO_W, jump +5 bytes forward, which is outside the code

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

// Test CheckGotow via GOTO_W bytecode with invalid jump
func TestCheckGotow_InvalidJumpNegative(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with GOTO_W (0xC8) with invalid negative jump beyond start
	code := []byte{0x00, 0xC8, 0xFF, 0xFF, 0xFF, 0xFE} // NOP then GOTO_W, jump beyond start

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

// Test checkInvokespecial via INVOKESPECIAL bytecode with valid method ref
func TestCheckInvokespecial_ValidMethodRef(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESPECIAL (0xB7)
	code := []byte{0xB7, 0x00, 0x01} // INVOKESPECIAL with CP index 1

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

// Test checkInvokespecial via INVOKESPECIAL bytecode with invalid CP slot
func TestCheckInvokespecial_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESPECIAL (0xB7) with invalid CP index
	code := []byte{0xB7, 0x00, 0xFF} // INVOKESPECIAL with invalid CP index

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

// Test checkInvokestatic via INVOKESTATIC bytecode with valid method ref
func TestCheckInvokestatic_ValidMethodRef(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESTATIC (0xB8)
	code := []byte{0xB8, 0x00, 0x01} // INVOKESTATIC with CP index 1

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

// Test checkInvokestatic via INVOKESTATIC bytecode with invalid CP slot
func TestCheckInvokestatic_InvalidCPSlot(t *testing.T) {
	globals.InitGlobals("test")

	// Create bytecode with INVOKESTATIC (0xB8) with invalid CP index
	code := []byte{0xB8, 0x00, 0xFF} // INVOKESTATIC with invalid CP index

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

// Test storeInt via ISTORE_0 bytecode
func TestStoreInt_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x3B} // ISTORE_0
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// Test storeIntRet2 via ISTORE bytecode
func TestStoreIntRet2_StackDecrement(t *testing.T) {
	globals.InitGlobals("test")

	code := []byte{0x36, 0x01} // ISTORE with index
	cp := CPool{}
	cp.CpIndex = make([]CpEntry, 10)
	cp.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	af := AccessFlags{}

	err := CheckCodeValidity(&code, &cp, 10, af)
	if err != nil {
		t.Errorf("CheckCodeValidity failed: %v", err)
	}
}

// Test BytecodePushes32BitValue function directly (it's exported)
func TestBytecodePushes32BitValue_True(t *testing.T) {
	globals.InitGlobals("test")

	// Test with ICONST_0 which should return true
	result := BytecodePushes32BitValue(0x03)

	if !result {
		t.Errorf("Expected true for ICONST_0 (0x03), got false")
	}

	// Test with BIPUSH which should return true
	result = BytecodePushes32BitValue(0x10)

	if !result {
		t.Errorf("Expected true for BIPUSH (0x10), got false")
	}
}

func TestBytecodePushes32BitValue_False(t *testing.T) {
	globals.InitGlobals("test")

	// Test with LCONST_0 which should return false (long/double)
	result := BytecodePushes32BitValue(0x09)

	if result {
		t.Errorf("Expected false for LCONST_0 (0x09), got true")
	}

	// Test with a random bytecode not in the list
	result = BytecodePushes32BitValue(0xFF)

	if result {
		t.Errorf("Expected false for unknown bytecode (0xFF), got true")
	}
}
