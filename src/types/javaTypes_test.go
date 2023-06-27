/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package types

import "testing"

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

func TestTheIsFunctions(t *testing.T) {
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
		t.Errorf("IsFloatingPoint() returned false for double, should be true")
	}

	if !UsesTwoSlots(Double) {
		t.Errorf("UsesTwoSlots() returned false for double, should be true")
	}
}
