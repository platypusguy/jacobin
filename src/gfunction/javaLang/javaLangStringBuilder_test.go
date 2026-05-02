/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestStringBuilderInit(t *testing.T) {
	// Test case: without capacity
	params := []any{object.MakeEmptyObject()}
	result := stringBuilderInit(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	// Test case: with capacity
	obj := object.MakeEmptyObject()
	params = []any{obj, int64(32)}
	result = stringBuilderInit(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	if obj.FieldTable["capacity"].Fvalue != int64(32) {
		t.Errorf("Expected 32, got %v", result.(*object.Object).FieldTable["capacity"].Fvalue)
	}
}

func TestStringBuilderInitString(t *testing.T) {
	// Test case: valid string object
	obj := object.MakeEmptyObject()
	strObj := object.StringObjectFromGoString("Hello")
	params := []any{obj, strObj}
	result := stringBuilderInitString(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	if object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte)) !=
		object.GoStringFromJavaByteArray(strObj.FieldTable["value"].Fvalue.([]types.JavaByte)) {
		t.Errorf("Expected %v, got %v", strObj.FieldTable["value"].Fvalue,
			result.(*object.Object).FieldTable["value"].Fvalue)
	}

	// Test case: invalid parameter type
	params = []any{obj, 123}
	result = stringBuilderInitString(params)
	if result == nil || result.(*ghelpers.GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException, got %s",
			excNames.JVMexceptionNames[result.(*ghelpers.GErrBlk).ExceptionType])
	}
}

func TestStringBuilderToString(t *testing.T) {
	// Test case: convert to string
	obj := object.NewStringObject()
	obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
	params := []any{obj}
	result := stringBuilderToString(params)
	if result == nil {
		t.Errorf("Expected non-nil result, got nil")
	}
	if result == "Hello" {
		t.Errorf("Expected 'Hello', got %s", result)
	}
}

func TestStringBuilderCapacity(t *testing.T) {
	// Test case: get capacity
	obj := object.MakeEmptyObject()
	obj.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: int64(16)}
	params := []any{obj}
	result := stringBuilderCapacity(params)
	if result != int64(16) {
		t.Errorf("Expected 16, got %v", result)
	}
}

func TestStringBuilderLength(t *testing.T) {
	// Test case: get length
	obj := object.MakeEmptyObject()
	obj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(5)}
	params := []any{obj}
	result := stringBuilderLength(params)
	if result != int64(5) {
		t.Errorf("Expected 5, got %v", result)
	}
}

func TestStringBuilder_AppendNullByte_ToString_GetBytes(t *testing.T) {
	// 1. Create StringBuilder
	sbObj := object.MakeEmptyObject()
	sbObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{}}
	sbObj.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: int64(16)}
	sbObj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	// 2. Append empty string (no-op)
	strObjEmpty := object.StringObjectFromGoString("")
	stringBuilderAppend([]any{sbObj, strObjEmpty})

	// 3. Append "\000"
	strObjNull := object.StringObjectFromGoString("\x00")
	stringBuilderAppend([]any{sbObj, strObjNull})

	// 4. sb.toString()
	resToString := stringBuilderToString([]any{sbObj})
	strObjRes := resToString.(*object.Object)

	// 5. getBytes("UTF-8")
	resBytes := getBytesFromString([]interface{}{strObjRes, object.StringObjectFromGoString("UTF-8")})
	arrObj := resBytes.(*object.Object)
	gotBytes := bytesFromByteArrayObject(arrObj)

	if len(gotBytes) != 1 {
		t.Fatalf("expected byte array of length 1, got %d", len(gotBytes))
	}
	if gotBytes[0] != 0 {
		t.Fatalf("expected byte 0x00, got 0x%02x", gotBytes[0])
	}

	// 6. Test with capacity > count
	sbObj2 := object.MakeEmptyObject()
	sbObj2.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: make([]types.JavaByte, 0, 16)}
	sbObj2.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: int64(16)}
	sbObj2.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	stringBuilderAppend([]any{sbObj2, object.StringObjectFromGoString("\x00")})

	resToString2 := stringBuilderToString([]any{sbObj2})
	strObjRes2 := resToString2.(*object.Object)

	// String length should be 1
	resLength := stringLength([]interface{}{strObjRes2}).(int64)
	if resLength != 1 {
		t.Errorf("expected string length 1, got %d", resLength)
	}

	resBytes2 := getBytesFromString([]interface{}{strObjRes2, object.StringObjectFromGoString("UTF-8")})
	gotBytes2 := bytesFromByteArrayObject(resBytes2.(*object.Object))
	if len(gotBytes2) != 1 {
		t.Errorf("expected byte array of length 1, got %d", len(gotBytes2))
	}
}

func TestString_UserReportedIssue(t *testing.T) {
	globals.InitStringPool()

	// String password = "";
	passwordStr := ""
	passwordObj := object.StringObjectFromGoString(passwordStr)

	// password.getBytes()
	res1 := getBytesFromString([]any{passwordObj})
	gotBytes1 := bytesFromByteArrayObject(res1.(*object.Object))
	if len(gotBytes1) != 0 {
		t.Errorf("DEBUG password: expected 0 bytes, got %d", len(gotBytes1))
	}

	// StringBuilder sb = new StringBuilder();
	sbObj := object.MakeEmptyObject()
	sbObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{}}
	sbObj.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: int64(16)}
	sbObj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	// sb.append(password);
	stringBuilderAppend([]any{sbObj, passwordObj})

	// sb.append("\000");
	stringBuilderAppend([]any{sbObj, object.StringObjectFromGoString("\x00")})

	// byte passwordb[] = sb.toString().getBytes("UTF-8");
	resToString := stringBuilderToString([]any{sbObj})
	strObjRes := resToString.(*object.Object)

	resBytes := getBytesFromString([]any{strObjRes, object.StringObjectFromGoString("UTF-8")})
	passwordb := bytesFromByteArrayObject(resBytes.(*object.Object))

	// Jacobin output (incorrect): length: 2, c0 80
	// Hotspot output: length: 1, 00
	if len(passwordb) != 1 {
		t.Errorf("DEBUG passwordb: expected length 1, got %d", len(passwordb))
	}
	if len(passwordb) > 0 && passwordb[0] != 0 {
		t.Errorf("DEBUG passwordb[0]: expected 00, got %02x", passwordb[0])
	}

	// Verify sb.toString() length
	resLen := stringLength([]any{strObjRes}).(int64)
	if resLen != 1 {
		t.Errorf("DEBUG sb.toString() length: expected 1, got %d", resLen)
	}

	// Verify Go string representation for logging/printing
	goStr := object.GoStringFromStringObject(strObjRes)
	if len(goStr) != 1 {
		t.Errorf("Go string length: expected 1, got %d (value: %x)", len(goStr), goStr)
	}

	// Double-check GoStringFromJavaByteArray directly
	jbarr := []types.JavaByte{0}
	goStrFromJBA := object.GoStringFromJavaByteArray(jbarr)
	if len(goStrFromJBA) != 1 || goStrFromJBA[0] != 0 {
		t.Errorf("GoStringFromJavaByteArray(0) failed: len=%d, hex=%x", len(goStrFromJBA), goStrFromJBA)
	}
}

func TestStringBuilder_AppendMultiByte(t *testing.T) {
	globals.InitStringPool()
	sbObj := object.MakeEmptyObject()
	sbObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{}}
	sbObj.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: int64(16)}
	sbObj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	// Append '€' as char (0x20AC)
	stringBuilderAppendChar([]any{sbObj, int64(0x20AC)})

	resToString := stringBuilderToString([]any{sbObj})
	strObjRes := resToString.(*object.Object)

	goStr := object.GoStringFromStringObject(strObjRes)
	if goStr != "€" {
		t.Errorf("expected '€', got %q", goStr)
	}

	resBytes := getBytesFromString([]any{strObjRes, object.StringObjectFromGoString("UTF-8")})
	gotBytes := bytesFromByteArrayObject(resBytes.(*object.Object))
	if len(gotBytes) != 3 {
		t.Errorf("expected 3 bytes for €, got %d", len(gotBytes))
	}
}

func TestStringBuilder_CharAt_Negative(t *testing.T) {
	globals.InitStringPool()
	sbObj := object.MakeEmptyObject()
	// Create a StringBuilder with a byte that has the sign bit set (0x80 = 128)
	sbObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{types.JavaByte(-128)}}
	sbObj.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: int64(16)}
	sbObj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}

	res := stringBuilderCharAt([]any{sbObj, int64(0)})
	charVal, ok := res.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", res)
	}

	if charVal != 128 {
		t.Errorf("expected char value 128, got %d", charVal)
	}
}
