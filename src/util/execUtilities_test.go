/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"jacobin/src/types"
	"testing"
)

// Table-driven, OS-independent tests for ParseIncomingParamsFromMethTypeString
func TestParseIncomingParamsFromMethType_Basics(t *testing.T) {
	tests := []struct {
		name      string
		meth      string
		expCount  int
		expConcat string
	}{
		{"empty", "", 0, ""},
		{"primitives_group", "(SBI)", 3, "III"},
		{"mix_with_array", "(S[BI)I", 3, "I[BI"},
		{"single_int", "(I)V", 1, types.Int},
	}
	for _, tc := range tests {
		res := ParseIncomingParamsFromMethTypeString(tc.meth)
		if len(res) != tc.expCount {
			t.Fatalf("%s: expected %d parsed parameters, got %d", tc.name, tc.expCount, len(res))
		}
		var got string
		for _, r := range res {
			got += r
		}
		if got != tc.expConcat {
			t.Fatalf("%s: expected concat '%s', got '%s'", tc.name, tc.expConcat, got)
		}
	}
}

func TestParseIncomingParamsFromMethType_Table(t *testing.T) {
	tests := []struct {
		name      string
		meth      string
		expCount  int
		expConcat string
	}{
		{"refs_and_longs", "(LString;Ljava/lang/Integer;JJ)V", 4, "LLJJ"},
		{"two_strings", "(Ljava/lang/String;Ljava/lang/String;)Ljava/nio/file/Path;", 2, "LL"},
		{"string_and_array_of_strings", "(Ljava/lang/String;[Ljava/lang/String;)Ljava/nio/file/Path;", 2, "L[L"},
		{"array_of_strings_only", "([Ljava/lang/String;)V", 1, "[L"},
		{"array_string_long_string", "([Ljava/lang/String;JLjava/lang/String;)V", 3, "[LJL"},
		{"array_string_long_array_string", "([Ljava/lang/String;J[Ljava/lang/String;)V", 3, "[LJ[L"},
		{"float_arrays_string_long_array_string_double", "(F[Ljava/lang/String;J[Ljava/lang/String;D)V", 5, "F[LJ[LD"},
		{"deep_array_strings", "(F[[[[[Ljava/lang/String;J[Ljava/lang/String;D)V", 5, "F[[[[[LJ[LD"},
		{"invalid_missing_semicolon_ref", "(Labc)V", 0, ""},
		{"invalid_missing_semicolon_array_ref", "([Labc)V", 0, ""},
		{"invalid_missing_semicolon_multi_array_ref", "([[[Labc)V", 0, ""},
		{"invalid_unended_arrays", "([[[[)V", 0, ""},
		{"valid_4d_int_array", "([[[[I)V", 1, "[[[[I"},
		{"mix_with_deep_arrays", "(JF[[[[[Ljava/lang/String;[[[J)V", 4, "JF[[[[[L[[[J"},
		{"invalid_illegal_char", "(JD[[[[[Ljava/lang/String;[[[J%)V", 0, ""},
		{"valid_mixed_arrays_refs", "(JD[[I[[[Ljava/lang/String;[[[J)V", 5, "JD[[I[[[L[[[J"},
		{"invalid_illegal_char_mid", "(JD[I[F[[[Ljava/lang/String;%[[[J)V", 0, ""},
		{"valid_mixed_prims_and_refs", "(JD[I[F[[[Ljava/lang/String;[[[J)V", 6, "JD[I[F[[[L[[[J"},
	}
	for _, tc := range tests {
		res := ParseIncomingParamsFromMethTypeString(tc.meth)
		if len(res) != tc.expCount {
			t.Errorf("%s: expected %d parsed parameters, got %d", tc.name, tc.expCount, len(res))
			continue
		}
		var got string
		for _, r := range res {
			got += r
		}
		if got != tc.expConcat {
			t.Errorf("%s: expected concat '%s', got '%s'", tc.name, tc.expConcat, got)
		}
	}
}
