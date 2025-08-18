/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/excNames"
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
	if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException, got %s",
			excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
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
