/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"fmt"
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

// test that pointer/reference in the params is handled correctly
// especially, that references (start with L and with ;) are correctly
// parsed and represented in the output

// Parameter-driven checker
func checker(t *testing.T, caller int, methType string, expCount int, expString string) {
	res := ParseIncomingParamsFromMethTypeString(methType)
	if len(res) != expCount { // short, byte and int all become 'I'
		t.Errorf("[%d] Expected %d parsed parameters, got %d", caller, expCount, len(res))
		for ii := 0; ii < len(res); ii++ {
			fmt.Printf("Parameter %d: %v\n", ii, res[ii])
		}
	}

	var paramString string
	for ii := 0; ii < len(res); ii++ {
		paramString += res[ii]
	}
	if paramString != expString {
		t.Errorf("[%d] Expected param string of '%s', got: '%s'", caller, expString, paramString)
	}
}

// Individual tests for ParseIncomingParamsFromMethTypeString

func TestParseIncomingReferenceParamsFromMethType1(t *testing.T) {
	checker(t, 1, "(LString;Ljava/lang/Integer;JJ)V", 4, "LLJJ")
}

func TestParseIncomingReferenceParamsFromMethType2(t *testing.T) {
	checker(t, 2, "(Ljava/lang/String;Ljava/lang/String;)Ljava/nio/file/Path;", 2, "LL")
}

func TestParseIncomingReferenceParamsFromMethType3(t *testing.T) {
	checker(t, 3, "(Ljava/lang/String;[Ljava/lang/String;)Ljava/nio/file/Path;", 2, "L[L")
}

func TestParseIncomingReferenceParamsFromMethType4(t *testing.T) {
	checker(t, 4, "([Ljava/lang/String;)V", 1, "[L")
}

func TestParseIncomingReferenceParamsFromMethType5(t *testing.T) {
	checker(t, 5, "([Ljava/lang/String;JLjava/lang/String;)V", 3, "[LJL")
}

func TestParseIncomingReferenceParamsFromMethType6(t *testing.T) {
	checker(t, 6, "([Ljava/lang/String;J[Ljava/lang/String;)V", 3, "[LJ[L")
}

func TestParseIncomingReferenceParamsFromMethType7(t *testing.T) {
	checker(t, 7, "(F[Ljava/lang/String;J[Ljava/lang/String;D)V", 5, "F[LJ[LD")
}

func TestParseIncomingReferenceParamsFromMethType8(t *testing.T) {
	checker(t, 8, "(F[[[[[Ljava/lang/String;J[Ljava/lang/String;D)V", 5, "F[[[[[LJ[LD")
}

func TestParseIncomingReferenceParamsFromMethType9(t *testing.T) {
	checker(t, 9, "(Labc)V", 0, "")
}

func TestParseIncomingReferenceParamsFromMethType10(t *testing.T) {
	checker(t, 10, "([Labc)V", 0, "")
}

func TestParseIncomingReferenceParamsFromMethType11(t *testing.T) {
	checker(t, 11, "([[[Labc)V", 0, "")
}

func TestParseIncomingReferenceParamsFromMethType12(t *testing.T) {
	checker(t, 12, "([[[[)V", 0, "")
}

func TestParseIncomingReferenceParamsFromMethType13(t *testing.T) {
	checker(t, 13, "([[[[I)V", 1, "[[[[I")
}

func TestParseIncomingReferenceParamsFromMethType14(t *testing.T) {
	checker(t, 14, "(JF[[[[[Ljava/lang/String;[[[J)V", 4, "JF[[[[[L[[[J")
}

func TestParseIncomingReferenceParamsFromMethType15(t *testing.T) {
	checker(t, 15, "(JD[[[[[Ljava/lang/String;[[[J%)V", 0, "")
}

func TestParseIncomingReferenceParamsFromMethType16(t *testing.T) {
	checker(t, 16, "(JD[[I[[[Ljava/lang/String;[[[J)V", 5, "JD[[I[[[L[[[J")
}

func TestParseIncomingReferenceParamsFromMethType17(t *testing.T) {
	checker(t, 17, "(JD[I[F[[[Ljava/lang/String;%[[[J)V", 0, "")
}

func TestParseIncomingReferenceParamsFromMethType18(t *testing.T) {
	checker(t, 18, "(JD[I[F[[[Ljava/lang/String;[[[J)V", 6, "JD[I[F[[[L[[[J")
}
