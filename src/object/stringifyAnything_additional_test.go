/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0). Consult jacobin.org.
 */

package object

import (
	"jacobin/src/globals"
	"jacobin/src/types"
	"testing"
)

func TestStringifyAnythingGo_SpecialClassDelegation(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	originalInvoker := global.FuncInvokeGFunction
	defer func() {
		global.FuncInvokeGFunction = originalInvoker
	}()

	tests := []struct {
		name      string
		className string
		method    string
	}{
		{
			name:      "local date time",
			className: "java/time/LocalDateTime",
			method:    "java/time/LocalDateTime.toString()Ljava/lang/String;",
		},
		{
			name:      "locale",
			className: "java/util/Locale",
			method:    "java/util/Locale.toString()Ljava/lang/String;",
		},
		{
			name:      "big decimal",
			className: types.ClassNameBigDecimal,
			method:    "java/math/BigDecimal.toString()Ljava/lang/String;",
		},
		{
			name:      "big integer",
			className: types.ClassNameBigInteger,
			method:    "java/math/BigInteger.toString()Ljava/lang/String;",
		},
		{
			name:      "thread",
			className: types.ClassNameThread,
			method:    "java/lang/Thread.toString()Ljava/lang/String;",
		},
		{
			name:      "thread state",
			className: types.ClassNameThreadState,
			method:    "java/lang/Thread$State.toString()Ljava/lang/String;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := MakeEmptyObjectWithClassName(&tt.className)
			called := false
			returnedValue := "string from " + tt.name
			global.FuncInvokeGFunction = func(method string, args []any) any {
				called = true
				if method != tt.method {
					t.Errorf("invoked method = %q, want %q", method, tt.method)
				}
				if len(args) != 1 {
					t.Fatalf("invoker received %d arguments, want exactly 1", len(args))
				}
				if args[0] != receiver {
					t.Errorf("invoker receiver = %v, want original object %v", args[0], receiver)
				}
				return StringObjectFromGoString(returnedValue)
			}

			if got := StringifyAnythingGo(receiver); got != returnedValue {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, returnedValue)
			}
			if !called {
				t.Error("global invoker was not called")
			}
		})
	}
}

func TestStringifyAnythingGo_ReferenceArraysAndNestedValues(t *testing.T) {
	globals.InitGlobals("test")

	stringArrayClass := "[Ljava/lang/String"
	stringArray := MakeEmptyObjectWithClassName(&stringArrayClass)
	stringArray.FieldTable["value"] = Field{
		Ftype:  types.Array,
		Fvalue: []any{StringObjectFromGoString("first"), (*Object)(nil), StringObjectFromGoString("third")},
	}
	if got, want := StringifyAnythingGo(stringArray), "[first, null, third]"; got != want {
		t.Errorf("StringifyAnythingGo(reference array) = %q, want %q", got, want)
	}

	nestedArrayClass := "[[I"
	nestedArray := MakeEmptyObjectWithClassName(&nestedArrayClass)
	nestedArray.FieldTable["value"] = Field{
		Ftype:  types.MultiArray,
		Fvalue: []any{[]int64{1, 2}, []int64{3}},
	}
	if got, want := StringifyAnythingGo(nestedArray), "[[1, 2], [3]]"; got != want {
		t.Errorf("StringifyAnythingGo(nested array) = %q, want %q", got, want)
	}

	fieldArray := Field{Ftype: types.Array, Fvalue: []any{"item", int64(7), nil}}
	if got, want := StringifyAnythingGo(fieldArray), "[item, 7, null]"; got != want {
		t.Errorf("StringifyAnythingGo(array field) = %q, want %q", got, want)
	}
}

func TestStringifyAnythingGo_NullFieldValues(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "nil string reference", value: Field{Ftype: types.StringClassRef, Fvalue: nil}},
		{name: "nil object reference", value: Field{Ftype: types.Ref, Fvalue: (*Object)(nil)}},
		{name: "typed nil object", value: (*Object)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringifyAnythingGo(tt.value); got != types.NullString {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, types.NullString)
			}
		})
	}
}

func TestStringifyAnythingGo_PrimitiveFallbacks(t *testing.T) {
	var nilIntPointer *int
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "signed integer", value: int8(-12), want: "-12"},
		{name: "unsigned integer", value: uint16(42), want: "42"},
		{name: "float32", value: float32(1.25), want: "1.25"},
		{name: "mixed slice", value: []any{int8(-2), uint16(3), false, "four"}, want: "[-2, 3, false, four]"},
		{name: "pointer value", value: func() *int { value := 9; return &value }(), want: "9"},
		{name: "nil pointer", value: nilIntPointer, want: types.NullString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringifyAnythingGo(tt.value); got != tt.want {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringifyAnythingGo_StringRepresentationsAndMalformedFields(t *testing.T) {
	globals.InitGlobals("test")
	stringClass := types.StringClassName
	tests := []struct {
		name  string
		field Field
		want  string
	}{
		{name: "Go byte slice", field: Field{Ftype: types.StringClassRef, Fvalue: []byte("go bytes")}, want: "go bytes"},
		{name: "string value", field: Field{Ftype: types.StringClassRef, Fvalue: "raw string"}, want: "*** ERROR, StringifyAnythingGo Field types.StringClassRef: expected Object but saw type string"},
		{name: "array field is not a slice", field: Field{Ftype: types.Array, Fvalue: "not an array"}, want: "*** ERROR, StringifyAnythingGo Field: expected array/slice but saw string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringifyAnythingGo(tt.field); got != tt.want {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, tt.want)
			}
		})
	}

	for _, value := range []any{[]byte("Go bytes"), "raw string"} {
		obj := MakeEmptyObjectWithClassName(&stringClass)
		obj.FieldTable["value"] = Field{Ftype: types.StringClassRef, Fvalue: value}
		want := "Go bytes"
		if value == "raw string" {
			want = "raw string"
		}
		if got := StringifyAnythingGo(obj); got != want {
			t.Errorf("StringifyAnythingGo(String with %T value) = %q, want %q", value, got, want)
		}
	}
}

func TestStringifyAnythingJava_NonStringInput(t *testing.T) {
	globals.InitGlobals("test")
	result := StringifyAnythingJava(int32(314))
	if result == nil {
		t.Fatal("StringifyAnythingJava() returned nil")
	}
	if got, want := StringifyAnythingGo(result), "314"; got != want {
		t.Errorf("StringifyAnythingJava() content = %q, want %q", got, want)
	}
}
