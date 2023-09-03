/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"jacobin/types"
	"testing"
)

// verify that a trailing slash in JAVA_HOME is removed
func TestParseIncomingParamsFromMethType(t *testing.T) {
	res := ParseIncomingParamsFromMethTypeString("(SBI)")
	if len(res) != 3 { // short, byte and int all become 'I'
		t.Errorf("Expected 3 parsed parameters, got %d", len(res))
	}

	if res[0] != types.Int || res[1] != types.Int || res[2] != types.Int {
		t.Errorf("Expected parse would return 3 values of 'I', got: %s%s%s",
			res[0], res[1], res[2])
	}

	res = ParseIncomingParamsFromMethTypeString("(S[BI)I")
	if len(res) != 3 { // short, byte and int all become 'I' aka types.Int
		t.Errorf("Expected 3 parsed parameters, got %d", len(res))
	}

	if res[0] != types.Int || res[1] != types.ByteArray || res[2] != types.Int {
		t.Errorf("Expected parse would return S [B I, got: %s %s %s",
			res[0], res[1], res[2])
	}

	res = ParseIncomingParamsFromMethTypeString("")
	if len(res) != 0 {
		t.Errorf("Expected parse would return value an empty string array, got: %s", res)
	}
}

// test that pointer/refernce in the params is handled correctly
// especially, that references (start with L and with ;) are correctly
// parsed and represeneted in the output
func TestParseIncomingReferenceParamsFromMethType(t *testing.T) {
	res := ParseIncomingParamsFromMethTypeString("(LString;Ljava/lang/Integer;JJ)")
	if len(res) != 4 { // short, byte and int all become 'I'
		t.Errorf("Expected 4 parsed parameters, got %d", len(res))
	}

	var params string = res[0] + res[1] + res[2] + res[3]
	if params != "LLJJ" {
		t.Errorf("Expected param string of 'LLJJ', got: %s", params)
	}
}
