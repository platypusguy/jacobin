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
	"testing"
)

// Test PushFloatRet2 via FLOAD bytecode (which uses PushFloatRet2)

// Test storeInt via ISTORE_0 bytecode

// Test storeIntRet2 via ISTORE bytecode

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
