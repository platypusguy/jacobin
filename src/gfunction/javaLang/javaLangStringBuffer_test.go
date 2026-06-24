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

func TestStringBufferInit(t *testing.T) {
	// Test case: without capacity
	obj1 := object.MakeEmptyObject()
	params := []any{obj1}
	result := stringBufferInit(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
	if obj1.FieldTable["capacity"].Fvalue.(int64) != 16 {
		t.Errorf("Expected default capacity 16, got %v", obj1.FieldTable["capacity"].Fvalue)
	}
	if obj1.FieldTable["count"].Fvalue.(int64) != 0 {
		t.Errorf("Expected count 0, got %v", obj1.FieldTable["count"].Fvalue)
	}

	// Test case: with capacity
	obj2 := object.MakeEmptyObject()
	params = []any{obj2, int64(32)}
	result = stringBufferInit(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
	if obj2.FieldTable["capacity"].Fvalue.(int64) != 32 {
		t.Errorf("Expected capacity 32, got %v", obj2.FieldTable["capacity"].Fvalue)
	}
}

func TestStringBufferInitString(t *testing.T) {
	globals.InitStringPool()

	// Test case: valid string object
	obj := object.MakeEmptyObject()
	strObj := object.StringObjectFromGoString("Hello")
	params := []any{obj, strObj}
	result := stringBufferInitString(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	expectedValue := "Hello"
	gotValue := object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if gotValue != expectedValue {
		t.Errorf("Expected %q, got %q", expectedValue, gotValue)
	}
	if obj.FieldTable["count"].Fvalue.(int64) != 5 {
		t.Errorf("Expected count 5, got %v", obj.FieldTable["count"].Fvalue)
	}

	// Test case: invalid parameter type
	params = []any{obj, 123}
	result = stringBufferInitString(params)
	if result == nil {
		t.Fatal("Expected non-nil result for invalid parameter type")
	}
	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected *ghelpers.GErrBlk, got %T", result)
	}
	if errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException, got %v", errBlk.ExceptionType)
	}
}

func TestStringBufferOperations(t *testing.T) {
	globals.InitStringPool()
	sbClassName := "java/lang/StringBuffer"

	// Initialize a StringBuffer
	sbObj := object.MakeEmptyObjectWithClassName(&sbClassName)
	stringBufferInit([]any{sbObj, int64(16)})

	// Test append(String)
	stringBuilderAppend([]any{sbObj, object.StringObjectFromGoString("Hello")})

	// Test length()
	resLen := stringBuilderLength([]any{sbObj})
	if resLen != int64(5) {
		t.Errorf("Expected length 5, got %v", resLen)
	}

	// Test append(char)
	stringBuilderAppendChar([]any{sbObj, int64(' ')})

	// Test append(CharSequence)
	stringBuilderAppend([]any{sbObj, object.StringObjectFromGoString("World")})

	// Test toString()
	resStr := stringBuilderToString([]any{sbObj}).(*object.Object)
	if object.GoStringFromStringObject(resStr) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", object.GoStringFromStringObject(resStr))
	}

	// Test charAt(int)
	resChar := stringBuilderCharAt([]any{sbObj, int64(0)})
	if resChar.(int64) != int64('H') {
		t.Errorf("Expected 'H', got %v", resChar)
	}

	// Test setCharAt(int, char)
	stringBuilderSetCharAt([]any{sbObj, int64(0), int64('h')})
	resChar = stringBuilderCharAt([]any{sbObj, int64(0)})
	if resChar.(int64) != int64('h') {
		t.Errorf("Expected 'h', got %v", resChar)
	}

	// Test reverse()
	stringBuilderReverse([]any{sbObj})
	resStr = stringBuilderToString([]any{sbObj}).(*object.Object)
	if object.GoStringFromStringObject(resStr) != "dlroW olleh" {
		t.Errorf("Expected 'dlroW olleh', got %q", object.GoStringFromStringObject(resStr))
	}

	// Test delete(int, int)
	stringBuilderReverse([]any{sbObj}) // Reverse back to "hello World"
	stringBuilderDelete([]any{sbObj, int64(5), int64(11)})
	resStr = stringBuilderToString([]any{sbObj}).(*object.Object)
	if object.GoStringFromStringObject(resStr) != "hello" {
		t.Errorf("Expected 'hello', got %q", object.GoStringFromStringObject(resStr))
	}
}

func TestStringBuffer_Insert(t *testing.T) {
	globals.InitStringPool()
	sbClassName := "java/lang/StringBuffer"
	sbObj := object.MakeEmptyObjectWithClassName(&sbClassName)
	stringBufferInit([]any{sbObj, int64(16)})
	stringBuilderAppend([]any{sbObj, object.StringObjectFromGoString("Hello")})

	// Test insert(int, String)
	stringBuilderInsert([]any{sbObj, int64(5), object.StringObjectFromGoString(" World")})
	resStr := stringBuilderToString([]any{sbObj}).(*object.Object)
	if object.GoStringFromStringObject(resStr) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", object.GoStringFromStringObject(resStr))
	}

	// Test insert(int, char)
	stringBuilderInsertChar([]any{sbObj, int64(0), int64('!')})
	resStr = stringBuilderToString([]any{sbObj}).(*object.Object)
	if object.GoStringFromStringObject(resStr) != "!Hello World" {
		t.Errorf("Expected '!Hello World', got %q", object.GoStringFromStringObject(resStr))
	}
}

func TestStringBuffer_Traps(t *testing.T) {
	globals.InitStringPool()
	sbClassName := "java/lang/StringBuffer"
	sbObj := object.MakeEmptyObjectWithClassName(&sbClassName)
	stringBufferInit([]any{sbObj, int64(16)})

	// Test a trapped method (e.g., chars())
	res := ghelpers.TrapFunction([]any{sbObj})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected *ghelpers.GErrBlk, got %T", res)
	}
	if errBlk.ExceptionType != excNames.UnsupportedOperationException {
		t.Errorf("Expected UnsupportedOperationException, got %v", errBlk.ExceptionType)
	}
}
