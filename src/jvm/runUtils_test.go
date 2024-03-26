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
	res := convertIntegralValueToInt64(bite)
	if res != 1 {
		t.Errorf("convertBoolByteToInt64(byte), expected = 1, got %d", res)
	}

	yesNo := true
	if convertIntegralValueToInt64(yesNo) != 1 {
		t.Errorf("convertBoolByteToInt64(bool) != 1 (true), got %d", res)
	}
}
