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

// Table-driven tests for ParseIncomingParamsFromMethTypeStringRaw
func TestParseIncomingParamsFromMethTypeStringRaw_Basics(t *testing.T) {
	tests := []struct {
		name     string
		meth     string
		expected []string
	}{
		{"empty", "", []string{}},
		{"no_params_void", "()V", []string{}},
		{"no_params_int_return", "()I", []string{}},
		{"no_params_ref_return", "()Ljava/lang/String;", []string{}},
		{"single_byte", "(B)V", []string{types.Byte}},
		{"single_char", "(C)V", []string{types.Char}},
		{"single_int", "(I)V", []string{types.Int}},
		{"single_short", "(S)V", []string{types.Short}},
		{"single_float", "(F)V", []string{types.Float}},
		{"single_long", "(J)V", []string{types.Long}},
		{"single_double", "(D)V", []string{types.Double}},
		{"single_bool", "(Z)V", []string{types.Bool}},
		{"single_ref", "(Ljava/lang/String;)V", []string{types.Ref}},
		{"primitives_distinct_unreduced", "(SBI)", []string{types.Short, types.Byte, types.Int}},
		{"byte_int_pair", "(BI)Ljava/lang/String;", []string{types.Byte, types.Int}},
		{"mix_with_array", "(S[BI)I", []string{types.Short, "[B", types.Int}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseIncomingParamsFromMethTypeStringRaw(tc.meth)
			if len(res) != len(tc.expected) {
				t.Fatalf("%s: expected %d params %v, got %d: %v", tc.name, len(tc.expected), tc.expected, len(res), res)
			}
			for i := range res {
				if res[i] != tc.expected[i] {
					t.Errorf("%s: param[%d] expected '%s', got '%s'", tc.name, i, tc.expected[i], res[i])
				}
			}
		})
	}
}

func TestParseIncomingParamsFromMethTypeStringRaw_Table(t *testing.T) {
	tests := []struct {
		name     string
		meth     string
		expected []string
	}{
		{
			name:     "all_primitive_types_unreduced",
			meth:     "(BCISFJDZ)V",
			expected: []string{types.Byte, types.Char, types.Int, types.Short, types.Float, types.Long, types.Double, types.Bool},
		},
		{
			name:     "refs_and_longs",
			meth:     "(LString;Ljava/lang/Integer;JJ)V",
			expected: []string{types.Ref, types.Ref, types.Long, types.Long},
		},
		{
			name:     "two_strings",
			meth:     "(Ljava/lang/String;Ljava/lang/String;)Ljava/nio/file/Path;",
			expected: []string{types.Ref, types.Ref},
		},
		{
			name:     "string_and_array_of_strings",
			meth:     "(Ljava/lang/String;[Ljava/lang/String;)Ljava/nio/file/Path;",
			expected: []string{types.Ref, "[L"},
		},
		{
			name:     "array_of_strings_only",
			meth:     "([Ljava/lang/String;)V",
			expected: []string{"[L"},
		},
		{
			name:     "all_1d_primitive_arrays",
			meth:     "([B[C[I[S[F[J[D[Z)V",
			expected: []string{"[B", "[C", "[I", "[S", "[F", "[J", "[D", "[Z"},
		},
		{
			name:     "multi_dimensional_primitive_arrays",
			meth:     "([[B[[[C[[[[I[[[[[S[[F[[[J[[[[D[[[[[Z)V",
			expected: []string{"[[B", "[[[C", "[[[[I", "[[[[[S", "[[F", "[[[J", "[[[[D", "[[[[[Z"},
		},
		{
			name:     "multi_dimensional_reference_arrays",
			meth:     "([[Ljava/lang/Object;[[[[[Ljava/lang/String;)V",
			expected: []string{"[[L", "[[[[[L"},
		},
		{
			name:     "array_string_long_string",
			meth:     "([Ljava/lang/String;JLjava/lang/String;)V",
			expected: []string{"[L", types.Long, types.Ref},
		},
		{
			name:     "array_string_long_array_string",
			meth:     "([Ljava/lang/String;J[Ljava/lang/String;)V",
			expected: []string{"[L", types.Long, "[L"},
		},
		{
			name:     "float_arrays_string_long_array_string_double",
			meth:     "(F[Ljava/lang/String;J[Ljava/lang/String;D)V",
			expected: []string{types.Float, "[L", types.Long, "[L", types.Double},
		},
		{
			name:     "deep_array_strings",
			meth:     "(F[[[[[Ljava/lang/String;J[Ljava/lang/String;D)V",
			expected: []string{types.Float, "[[[[[L", types.Long, "[L", types.Double},
		},
		{
			name:     "mix_with_deep_arrays",
			meth:     "(JF[[[[[Ljava/lang/String;[[[J)V",
			expected: []string{types.Long, types.Float, "[[[[[L", "[[[J"},
		},
		{
			name:     "valid_mixed_arrays_refs",
			meth:     "(JD[[I[[[Ljava/lang/String;[[[J)V",
			expected: []string{types.Long, types.Double, "[[I", "[[[L", "[[[J"},
		},
		{
			name:     "valid_mixed_prims_and_refs",
			meth:     "(JD[I[F[[[Ljava/lang/String;[[[J)V",
			expected: []string{types.Long, types.Double, "[I", "[F", "[[[L", "[[[J"},
		},
		{
			name:     "comprehensive_all_variants_combined",
			meth:     "(BCISFJDZLjava/lang/String;[B[C[I[S[F[J[D[Z[Ljava/lang/Object;[[I[[[Ljava/lang/String;)Ljava/lang/Object;",
			expected: []string{types.Byte, types.Char, types.Int, types.Short, types.Float, types.Long, types.Double, types.Bool, types.Ref, "[B", "[C", "[I", "[S", "[F", "[J", "[D", "[Z", "[L", "[[I", "[[[L"},
		},
		// Invalid / Malformed cases
		{
			name:     "invalid_missing_semicolon_ref",
			meth:     "(Labc)V",
			expected: []string{},
		},
		{
			name:     "invalid_missing_semicolon_array_ref",
			meth:     "([Labc)V",
			expected: []string{},
		},
		{
			name:     "invalid_missing_semicolon_multi_array_ref",
			meth:     "([[[Labc)V",
			expected: []string{},
		},
		{
			name:     "invalid_unended_arrays_before_close_paren",
			meth:     "([[[[)V",
			expected: []string{},
		},
		{
			name:     "invalid_single_array_bracket_before_close_paren",
			meth:     "([)V",
			expected: []string{},
		},
		{
			name:     "invalid_unended_arrays_at_end_of_string",
			meth:     "([[[",
			expected: []string{},
		},
		{
			name:     "invalid_unended_ref_at_end_of_string",
			meth:     "(Ljava/lang/String",
			expected: []string{},
		},
		{
			name:     "invalid_unended_array_ref_at_end_of_string",
			meth:     "([Ljava/lang/String",
			expected: []string{},
		},
		{
			name:     "invalid_illegal_char_at_end",
			meth:     "(JD[[[[[Ljava/lang/String;[[[J%)V",
			expected: []string{},
		},
		{
			name:     "invalid_illegal_char_in_middle",
			meth:     "(JD[I[F[[[Ljava/lang/String;%[[[J)V",
			expected: []string{},
		},
		{
			name:     "invalid_illegal_char_start",
			meth:     "(%I)V",
			expected: []string{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := ParseIncomingParamsFromMethTypeStringRaw(tc.meth)
			if len(res) != len(tc.expected) {
				t.Fatalf("%s: expected %d params %v, got %d: %v", tc.name, len(tc.expected), tc.expected, len(res), res)
			}
			for i := range res {
				if res[i] != tc.expected[i] {
					t.Errorf("%s: param[%d] expected '%s', got '%s'", tc.name, i, tc.expected[i], res[i])
				}
			}
		})
	}
}

// TestParseFQN_Standard tests a standard FQN with class, method, and signature.
func TestParseFQN_Standard(t *testing.T) {
	fqn := "java/lang/String.valueOf(I)Ljava/lang/String;"
	expectedClass := "java/lang/String"
	expectedMethod := "valueOf"
	expectedType := "(I)Ljava/lang/String;"

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_Standard: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_Standard: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_Standard: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_Constructor tests an FQN for a constructor method.
func TestParseFQN_Constructor(t *testing.T) {
	fqn := "java/lang/Object.<init>()V"
	expectedClass := "java/lang/Object"
	expectedMethod := "<init>"
	expectedType := "()V"

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_Constructor: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_Constructor: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_Constructor: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_StaticInitializer tests an FQN for a static initializer.
func TestParseFQN_StaticInitializer(t *testing.T) {
	fqn := "com/example/MyClass.<clinit>()V"
	expectedClass := "com/example/MyClass"
	expectedMethod := "<clinit>"
	expectedType := "()V"

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_StaticInitializer: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_StaticInitializer: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_StaticInitializer: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_NoSignature tests an FQN where no method signature is present.
func TestParseFQN_NoSignature(t *testing.T) {
	fqn := "com/example/MyClass.myMethod"
	expectedClass := "com/example/MyClass"
	expectedMethod := "myMethod"
	expectedType := ""

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_NoSignature: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_NoSignature: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_NoSignature: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_NoClassWithSignature tests an FQN without a class prefix but with a signature.
func TestParseFQN_NoClassWithSignature(t *testing.T) {
	fqn := "main([Ljava/lang/String;)V"
	expectedClass := ""
	expectedMethod := "main"
	expectedType := "([Ljava/lang/String;)V"

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_NoClassWithSignature: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_NoClassWithSignature: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_NoClassWithSignature: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_NoClassNoSignature tests an FQN with neither class nor signature.
func TestParseFQN_NoClassNoSignature(t *testing.T) {
	fqn := "justAMethodName"
	expectedClass := ""
	expectedMethod := "justAMethodName"
	expectedType := ""

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_NoClassNoSignature: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_NoClassNoSignature: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_NoClassNoSignature: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_EmptyString tests an empty input string.
func TestParseFQN_EmptyString(t *testing.T) {
	fqn := ""
	expectedClass := ""
	expectedMethod := ""
	expectedType := ""

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_EmptyString: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_EmptyString: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_EmptyString: Expected type '%s', got '%s'", expectedType, methodType)
	}
}

// TestParseFQN_OnlyClass tests a string that looks like only a class name.
func TestParseFQN_OnlyClass(t *testing.T) {
	fqn := "java/lang/String"
	expectedClass := "" // No dot for method, so it's treated as a method name
	expectedMethod := "java/lang/String"
	expectedType := ""

	className, methodName, methodType := ParseFQN(fqn)

	if className != expectedClass {
		t.Errorf("TestParseFQN_OnlyClass: Expected class '%s', got '%s'", expectedClass, className)
	}
	if methodName != expectedMethod {
		t.Errorf("TestParseFQN_OnlyClass: Expected method '%s', got '%s'", expectedMethod, methodName)
	}
	if methodType != expectedType {
		t.Errorf("TestParseFQN_OnlyClass: Expected type '%s', got '%s'", expectedType, methodType)
	}
}
