/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"testing"
)

// verify that a trailing slash in JAVA_HOME is removed
func TestParseIncomingParamsFromMethType(t *testing.T) {
	res1 := ParseIncomingParamsFromMethTypeString("(SBI)")
	if string(res1) != "III" {
		t.Errorf("Expected parse would return 3 values of 'I', got: %s", string(res1))
	}

	res2 := ParseIncomingParamsFromMethTypeString("FJDL")
	if string(res2) != "FJDL" {
		t.Errorf("Expected parse would return \"FJDL\", got: %s", string(res2))
	}

	res3 := ParseIncomingParamsFromMethTypeString("[[")
	if string(res3) != "LL" {
		t.Errorf("Expected parse would return value of \"LL\", got: %s", string(res3))
	}

	res4 := ParseIncomingParamsFromMethTypeString("")
	if len(string(res4)) != 0 {
		t.Errorf("Expected parse would return value an empty string, got: %s", string(res4))
	}
}
