/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import "testing"

// tests for runUtils.go. Note that most functions are tested inside the tests for run.go,
// but several benefit from standalone testing. Those are tested here

func TestConvertBoolByteToInt64(t *testing.T) {
	var bite = byte(0x01)
	res := convertInterfaceToInt64(bite)
	if res != 1 {
		t.Errorf("convertBoolByteToInt64(byte), expected = 1, got %d", res)
	}

	yesNo := true
	if convertInterfaceToInt64(yesNo) != 1 {
		t.Errorf("convertBoolByteToInt64(bool) != 1 (true), got %d", res)
	}
}

func TestConvertFloatToInt64RoundDown(t *testing.T) {
	f := float64(5432.10)
	val := convertInterfaceToUint64(f)
	if val != 5432 {
		t.Errorf("convertFloatToInt64(float64), expected = 5432, got %d", val)
	}
}

func TestConvertFloatToInt64RoundUp(t *testing.T) {
	f := float64(5432.501)
	val := convertInterfaceToUint64(f)
	if val != 5433 {
		t.Errorf("convertFloatToInt64(float64), expected = 5433, got %d", val)
	}
}

// golang bytes are unsigned 8-bit fields. However, when a byte is part of a
// larger number (i.e., a 32-bit field) the most significant bit can indeed
// represent a sign. This test makes sure we convert such a data byte to a
// negative number.
func TestDataByteToInt64(t *testing.T) {
	b := byte(0xA0)
	val := byteToInt64(b)
	if !(val < 0) {
		t.Errorf("dataByteToInt64(byte), expected value < 0, got %d", val)
	}
}
