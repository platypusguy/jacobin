/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0). Consult jacobin.org.
 */

package object

import (
	"jacobin/src/globals"
	"jacobin/src/types"
	"strings"
	"testing"
)

type stringifyInvokerTestError struct {
	message    string
	methodMiss bool
}

func (e stringifyInvokerTestError) Error() string {
	return e.message
}

func (e stringifyInvokerTestError) IsMethodNotFound() bool {
	return e.methodMiss
}

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
			name:      "big integer",
			className: types.ClassNameBigInteger,
			method:    "java/math/BigInteger.toString()Ljava/lang/String;",
		},
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

			got := StringifyAnythingGo(receiver)
			if got != returnedValue {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, returnedValue)
			}
			if !called {
				t.Error("global invoker was not called")
			}
		})
	}
}

func TestStringifyAnythingGo_DelegatesAnyClassToString(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	className := "example/Widget"
	receiver := MakeEmptyObjectWithClassName(&className)
	called := false
	global.FuncInvokeGFunction = func(method string, args []any) any {
		called = true
		if want := className + ".toString()Ljava/lang/String;"; method != want {
			t.Errorf("invoked method = %q, want %q", method, want)
		}
		if len(args) != 1 || args[0] != receiver {
			t.Errorf("invoker arguments = %v, want the original receiver", args)
		}
		return StringObjectFromGoString("widget string")
	}

	if got, want := StringifyAnythingGo(receiver), "widget string"; got != want {
		t.Errorf("StringifyAnythingGo() = %q, want %q", got, want)
	}
	if !called {
		t.Error("global invoker was not called")
	}
}

func TestStringifyAnythingGo_DelegationFallbackAndFailures(t *testing.T) {
	tests := []struct {
		name   string
		result any
		want   string
	}{
		{name: "missing method", result: stringifyInvokerTestError{message: "missing", methodMiss: true}, want: "Widget{}"},
		{name: "nil result", result: nil, want: "Widget{}"},
		{name: "invocation error", result: stringifyInvokerTestError{message: "failed"}, want: "*** ERROR, StringifyAnythingGo: example/Widget.toString()Ljava/lang/String; failed: failed"},
		{name: "malformed result", result: "not a Java object", want: "*** ERROR, StringifyAnythingGo: example/Widget.toString()Ljava/lang/String; returned string, not a non-null Object"},
		{name: "null string result", result: (*Object)(nil), want: "*** ERROR, StringifyAnythingGo: example/Widget.toString()Ljava/lang/String; returned *object.Object, not a non-null Object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globals.InitGlobals("test")
			global := globals.GetGlobalRef()
			className := "example/Widget"
			receiver := MakeEmptyObjectWithClassName(&className)
			global.FuncInvokeGFunction = func(string, []any) any { return tt.result }

			if got := StringifyAnythingGo(receiver); got != tt.want {
				t.Errorf("StringifyAnythingGo() = %q, want %q", got, tt.want)
			}
		})
	}

	t.Run("nil invoker", func(t *testing.T) {
		globals.InitGlobals("test")
		global := globals.GetGlobalRef()
		global.FuncInvokeGFunction = nil
		className := "example/Widget"
		if got, want := StringifyAnythingGo(MakeEmptyObjectWithClassName(&className)), "Widget{}"; got != want {
			t.Errorf("StringifyAnythingGo() = %q, want %q", got, want)
		}
	})
}

func TestStringifyAnythingGo_DoesNotDispatchJavaLangObject(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	called := false
	global.FuncInvokeGFunction = func(string, []any) any {
		called = true
		return StringObjectFromGoString("unexpected")
	}
	className := "java/lang/Object"
	if got, want := StringifyAnythingGo(MakeEmptyObjectWithClassName(&className)), "Object{}"; got != want {
		t.Errorf("StringifyAnythingGo() = %q, want %q", got, want)
	}
	if called {
		t.Error("java/lang/Object.toString was dispatched recursively")
	}
}

func TestStringifyAnythingGo_PreservesStringAndArrayHandlingBeforeDispatch(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	called := false
	global.FuncInvokeGFunction = func(string, []any) any {
		called = true
		return StringObjectFromGoString("unexpected")
	}
	stringObj := StringObjectFromGoString("plain string")
	if got, want := StringifyAnythingGo(stringObj), "plain string"; got != want {
		t.Errorf("StringifyAnythingGo(String) = %q, want %q", got, want)
	}
	arrayClass := "[I"
	array := MakeEmptyObjectWithClassName(&arrayClass)
	array.FieldTable["value"] = Field{Ftype: types.IntArray, Fvalue: []int64{1, 2}}
	if got, want := StringifyAnythingGo(array), "[1, 2]"; got != want {
		t.Errorf("StringifyAnythingGo(array) = %q, want %q", got, want)
	}
	if called {
		t.Error("String or array handling unexpectedly invoked the global callback")
	}
}

func TestStringifyAnythingGo_InvocationErrorsAreReadable(t *testing.T) {
	globals.InitGlobals("test")
	global := globals.GetGlobalRef()
	className := "example/Widget"
	global.FuncInvokeGFunction = func(string, []any) any {
		return stringifyInvokerTestError{message: "specific failure"}
	}
	if got := StringifyAnythingGo(MakeEmptyObjectWithClassName(&className)); !strings.Contains(got, "specific failure") {
		t.Errorf("StringifyAnythingGo() = %q, want an error containing the callback failure", got)
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
