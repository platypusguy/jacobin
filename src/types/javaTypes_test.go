/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package types

import "testing"

// Test that the functions return the correct results
func TestTheIsFunctionsValidate(t *testing.T) {
	if !IsIntegral(Bool) {
		t.Errorf("IsIntegral() returned false for boolean, should be true")
	}

	if !IsIntegral(Char) {
		t.Errorf("IsIntegral() returned false for char, should be true")
	}

	if !IsIntegral(Short) {
		t.Errorf("IsIntegral() returned false for short, should be true")
	}

	if !IsAddress(ByteArray) {
		t.Errorf("IsAddress() returned false for byte array, should be true")
	}

	if !IsStatic("X[B") {
		t.Errorf("IsStatic() returned false for 'X[B', should be true")
	}

	if !IsFloatingPoint(Float) {
		t.Errorf("IsFloatingPoint() returned false for float, should be true")
	}

	if !IsFloatingPoint(Double) {
		t.Errorf("IsFloatingPoint() returned false for double, should be true")
	}

	if !IsError("0") {
		t.Errorf("IsError returned false for Error, should be true")
	}
}

// Test that the functions don't return invalid results
func TestTheIsFunctionsNegatively(t *testing.T) {
	if IsIntegral(Error) {
		t.Errorf("Error incorrectly was true in IsIntegral()")
	}

	if IsFloatingPoint(Bool) {
		t.Errorf("Error: Bool incorrectly is true in IsFloatingPoint()")
	}

	if IsAddress(Int) {
		t.Errorf("Error: Int incorrectly is true in IsAddress()")
	}

	if IsStatic(ByteArray) {
		t.Errorf("Error: ByteArray incorrectly is true in IsStatic()")
	}

	if IsError(Short) {
		t.Errorf("Error: Short is incorrectly true in IsError()")
	}
}

// Test the go-to-Java conversion of booleans
func TestJavaBoolean(t *testing.T) {

	val := ConvertGoBoolToJavaBool(true)
	if val != JavaBoolTrue {
		t.Errorf("JavaBool: expected a result of 1, but got: %d", val)
	}

	val = ConvertGoBoolToJavaBool(false)
	if val != JavaBoolFalse {
		t.Errorf("JavaBool: expected a result of 0, but got: %d", val)
	}
}
