/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package ghelpers

import (
	"container/list"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestConvertArgsToParams_ZeroArgs(t *testing.T) {
	result := ConvertArgsToParams()
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

func TestConvertArgsToParams_OneArg(t *testing.T) {
	result := ConvertArgsToParams(42)
	if len(result) != 1 {
		t.Errorf("Expected length 1, got %d", len(result))
	}
	if result[0] != 42 {
		t.Errorf("Expected 42, got %v", result[0])
	}
}

func TestConvertArgsToParams_OneArgNil(t *testing.T) {
	result := ConvertArgsToParams(nil)
	if len(result) != 1 {
		t.Errorf("Expected length 1, got %d", len(result))
	}
	if result[0] != nil {
		t.Errorf("Expected nil, got %v", result[0])
	}
}

func TestConvertArgsToParams_TwoArgs(t *testing.T) {
	result := ConvertArgsToParams("hello", 100)
	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
	if result[0] != "hello" {
		t.Errorf("Expected 'hello', got %v", result[0])
	}
	if result[1] != 100 {
		t.Errorf("Expected 100, got %v", result[1])
	}
}

func TestConvertArgsToParams_TwoArgsWithNil(t *testing.T) {
	result := ConvertArgsToParams(nil, "world")
	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
	if result[0] != nil {
		t.Errorf("Expected nil at index 0, got %v", result[0])
	}
	if result[1] != "world" {
		t.Errorf("Expected 'world', got %v", result[1])
	}
}

func TestConvertArgsToParams_ThreeArgs(t *testing.T) {
	result := ConvertArgsToParams(1, 2.5, true)
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
	if result[0] != 1 {
		t.Errorf("Expected 1, got %v", result[0])
	}
	if result[1] != 2.5 {
		t.Errorf("Expected 2.5, got %v", result[1])
	}
	if result[2] != true {
		t.Errorf("Expected true, got %v", result[2])
	}
}

func TestConvertArgsToParams_ThreeArgsWithNil(t *testing.T) {
	result := ConvertArgsToParams("first", nil, "third")
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
	if result[0] != "first" {
		t.Errorf("Expected 'first', got %v", result[0])
	}
	if result[1] != nil {
		t.Errorf("Expected nil at index 1, got %v", result[1])
	}
	if result[2] != "third" {
		t.Errorf("Expected 'third', got %v", result[2])
	}
}

func TestConvertArgsToParams_FiveArgs(t *testing.T) {
	result := ConvertArgsToParams(10, 20, 30, 40, 50)
	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}
	for i := 0; i < 5; i++ {
		expected := (i + 1) * 10
		if result[i] != expected {
			t.Errorf("Expected %d at index %d, got %v", expected, i, result[i])
		}
	}
}

func TestConvertArgsToParams_FiveArgsWithNil(t *testing.T) {
	result := ConvertArgsToParams(nil, "two", nil, 4, nil)
	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}
	if result[0] != nil {
		t.Errorf("Expected nil at index 0, got %v", result[0])
	}
	if result[1] != "two" {
		t.Errorf("Expected 'two', got %v", result[1])
	}
	if result[2] != nil {
		t.Errorf("Expected nil at index 2, got %v", result[2])
	}
	if result[3] != 4 {
		t.Errorf("Expected 4, got %v", result[3])
	}
	if result[4] != nil {
		t.Errorf("Expected nil at index 4, got %v", result[4])
	}
}

func TestGErrBlk(t *testing.T) {
	var nilErr *GErrBlk
	if nilErr.Error() != "<nil>" {
		t.Errorf("nil GErrBlk Error() = %q, want <nil>", nilErr.Error())
	}
	if nilErr.IsMethodNotFound() {
		t.Error("nil GErrBlk IsMethodNotFound() returned true, want false")
	}

	blk := GetGErrBlk(excNames.IllegalArgumentException, "test error message")
	if blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("blk.ExceptionType = %d, want %d", blk.ExceptionType, excNames.IllegalArgumentException)
	}
	if blk.Error() != "test error message" {
		t.Errorf("blk.Error() = %q, want 'test error message'", blk.Error())
	}
	if blk.IsMethodNotFound() {
		t.Error("blk.IsMethodNotFound() = true, want false")
	}

	notFoundBlk := &GErrBlk{
		ExceptionType:  excNames.NoSuchMethodException,
		ErrMsg:         "method not found",
		methodNotFound: true,
	}
	if !notFoundBlk.IsMethodNotFound() {
		t.Error("notFoundBlk.IsMethodNotFound() = false, want true")
	}
}

func TestGenericReturnFunctions(t *testing.T) {
	if ClinitGeneric(nil) != nil {
		t.Error("ClinitGeneric did not return nil")
	}
	if JustReturn(nil) != nil {
		t.Error("JustReturn did not return nil")
	}
	if ReturnNull(nil) != object.Null {
		t.Error("ReturnNull did not return object.Null")
	}
	if ReturnFalse(nil) != types.JavaBoolFalse {
		t.Errorf("ReturnFalse did not return JavaBoolFalse, got %v", ReturnFalse(nil))
	}
	if ReturnTrue(nil) != types.JavaBoolTrue {
		t.Errorf("ReturnTrue did not return JavaBoolTrue, got %v", ReturnTrue(nil))
	}

	csObj := ReturnCharsetName(nil).(*object.Object)
	if cs := object.GoStringFromStringObject(csObj); cs != globals.GetCharsetName() {
		t.Errorf("ReturnCharsetName = %q, want %q", cs, globals.GetCharsetName())
	}
}

func TestFileEofHelpers(t *testing.T) {
	obj := object.MakeEmptyObject()
	if EofGet(obj) {
		t.Error("EofGet on empty object returned true, want false")
	}
	EofSet(obj, true)
	if !EofGet(obj) {
		t.Error("EofGet after EofSet(true) returned false, want true")
	}
	EofSet(obj, false)
	if EofGet(obj) {
		t.Error("EofGet after EofSet(false) returned true, want false")
	}

	obj.FieldTable[FileAtEOF] = object.Field{Ftype: types.Int, Fvalue: int64(123)}
	if EofGet(obj) {
		t.Error("EofGet with invalid non-bool field value returned true, want false")
	}
}

func TestInitBigIntegerField(t *testing.T) {
	tests := []struct {
		val         int64
		wantSignum  int64
	}{
		{val: 100, wantSignum: 1},
		{val: -50, wantSignum: -1},
		{val: 0, wantSignum: 0},
	}

	for _, tt := range tests {
		obj := object.MakeEmptyObject()
		InitBigIntegerField(obj, tt.val)

		valField, ok := obj.FieldTable["value"]
		if !ok || valField.Ftype != types.Ref {
			t.Fatalf("missing or invalid 'value' field for val %d", tt.val)
		}
		bigInt, ok := valField.Fvalue.(*big.Int)
		if !ok || bigInt.Int64() != tt.val {
			t.Errorf("bigInt.Int64() = %v, want %d", valField.Fvalue, tt.val)
		}

		signField, ok := obj.FieldTable["signum"]
		if !ok || signField.Ftype != types.Ref {
			t.Fatalf("missing or invalid 'signum' field for val %d", tt.val)
		}
		if signVal, ok := signField.Fvalue.(int64); !ok || signVal != tt.wantSignum {
			t.Errorf("signum = %v, want %d", signField.Fvalue, tt.wantSignum)
		}
	}
}

func TestReturnRandomLong(t *testing.T) {
	r1 := ReturnRandomLong(nil).(int64)
	r2 := ReturnRandomLong(nil).(int64)
	_ = r1
	_ = r2
}

func TestGetDefaultSecurityProvider(t *testing.T) {
	orig := DefaultSecurityProvider
	defer func() { DefaultSecurityProvider = orig }()

	DefaultSecurityProvider = object.MakeEmptyObject()
	if GetDefaultSecurityProvider() != DefaultSecurityProvider {
		t.Error("GetDefaultSecurityProvider did not return DefaultSecurityProvider")
	}
}

func TestGetLinkedListFromObject(t *testing.T) {
	emptyObj := object.MakeEmptyObject()
	listVal, err := GetLinkedListFromObject(emptyObj)
	if listVal != nil || err == nil {
		t.Fatal("expected error on empty object missing 'value' field")
	}
	if errBlk, ok := err.(*GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %+v", err)
	}

	invalidObj := object.MakeEmptyObject()
	invalidObj.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: "not a list"}
	listVal, err = GetLinkedListFromObject(invalidObj)
	if listVal != nil || err == nil {
		t.Fatal("expected error on object with non-list 'value' field")
	}
	if errBlk, ok := err.(*GErrBlk); !ok || errBlk.ExceptionType != excNames.VirtualMachineError {
		t.Errorf("expected VirtualMachineError, got %+v", err)
	}

	validObj := object.MakeEmptyObject()
	l := list.New()
	validObj.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: l}
	listVal, err = GetLinkedListFromObject(validObj)
	if err != nil || listVal != l {
		t.Errorf("expected list %v and nil error, got (%v, %v)", l, listVal, err)
	}
}

func TestInvoke(t *testing.T) {
	res := Invoke("non/existent/Class.method()V", nil)
	errBlk, ok := res.(*GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NoSuchMethodException || !errBlk.IsMethodNotFound() {
		t.Fatalf("Invoke non-existent method expected NoSuchMethod GErrBlk, got %+v", res)
	}

	testFuncSig := "test/Invoke/Method.test(I)I"
	MethodSignatures[testFuncSig] = GMeth{
		ParamSlots: 1,
		GFunction: func(params []any) any {
			if len(params) > 0 {
				return params[0].(int) * 2
			}
			return 0
		},
	}
	defer delete(MethodSignatures, testFuncSig)

	res = Invoke(testFuncSig, []any{21})
	if res != 42 {
		t.Errorf("Invoke returned %v, want 42", res)
	}
}

func TestJavaFloatingPointToString(t *testing.T) {
	tests := []struct {
		val     float64
		bitSize int
		want    string
	}{
		{val: math.NaN(), bitSize: 64, want: "NaN"},
		{val: math.Inf(1), bitSize: 64, want: "Infinity"},
		{val: math.Inf(-1), bitSize: 64, want: "-Infinity"},
		{val: 0.0, bitSize: 64, want: "0.0"},
		{val: math.Copysign(0.0, -1.0), bitSize: 64, want: "-0.0"},
		{val: 1.0, bitSize: 64, want: "1.0"},
		{val: -1.0, bitSize: 64, want: "-1.0"},
		{val: 123.456, bitSize: 64, want: "123.456"},
		{val: -123.456, bitSize: 64, want: "-123.456"},
		{val: 1e-3, bitSize: 64, want: "0.001"},
		{val: 1e-4, bitSize: 64, want: "1.0E-4"},
		{val: 1e7, bitSize: 64, want: "1.0E7"},
		{val: 1.5e7, bitSize: 64, want: "1.5E7"},
		{val: 100.0, bitSize: 64, want: "100.0"},
		{val: float64(float32(1.5)), bitSize: 32, want: "1.5"},
	}

	for _, tt := range tests {
		got := JavaFloatingPointToString(tt.val, tt.bitSize)
		if got != tt.want {
			t.Errorf("JavaFloatingPointToString(%v, %d) = %q, want %q", tt.val, tt.bitSize, got, tt.want)
		}
	}
}

func TestWrongArgType(t *testing.T) {
	res := WrongArgType("testFunc", "int", "invalidString")
	errBlk, ok := res.(*GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("WrongArgType expected IllegalArgumentException, got %+v", res)
	}
	if !strings.Contains(errBlk.ErrMsg, "testFunc") || !strings.Contains(errBlk.ErrMsg, "expected int") {
		t.Errorf("WrongArgType message = %q", errBlk.ErrMsg)
	}
}

func TestGetInt64Val(t *testing.T) {
	tests := []struct {
		input   any
		wantVal int64
		wantOk  bool
	}{
		{input: int64(12345), wantVal: 12345, wantOk: true},
		{input: int32(678), wantVal: 678, wantOk: true},
		{input: int(90), wantVal: 90, wantOk: true},
		{input: float64(12.3), wantVal: 0, wantOk: false},
		{input: "123", wantVal: 0, wantOk: false},
		{input: nil, wantVal: 0, wantOk: false},
	}

	for _, tt := range tests {
		val, ok := GetInt64Val(tt.input)
		if ok != tt.wantOk || val != tt.wantVal {
			t.Errorf("GetInt64Val(%v) = (%d, %t), want (%d, %t)", tt.input, val, ok, tt.wantVal, tt.wantOk)
		}
	}
}
