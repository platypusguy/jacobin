/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
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

/*
	func TestStringBuilderAppend(t *testing.T) {
		// Test case: append string
		sbParams := []any{object.NewStringObject(), "Yes, "}
		sb := stringBuilderInitString(sbParams)
		strObj := object.StringObjectFromGoString("Hello")
		params := []any{strObj, sb}
		result := stringBuilderAppend(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result.(*object.Object).FieldTable["value"].Fvalue != "Yes, Hello" {
			t.Errorf("Expected 'Yes, Hello', got %s", result)
		}

		// Test case: append int64
		sbParams = []any{object.MakeEmptyObject(), "It's "}
		sb = stringBuilderInitString(sbParams)
		params = []any{sb, int64(123)}
		result = stringBuilderAppend(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result != "It's 123" {
			t.Errorf("Expected 'It's 123', got %s", result)
		}

		// Test case: append float64
		sbParams = []any{object.MakeEmptyObject(), "It's "}
		sb = stringBuilderInitString(sbParams)
		params = []any{sb, 123.45}
		result = stringBuilderAppend(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result != "It's 123.45" {
			t.Errorf("Expected 'It's 123', got %s", result)
		}

		// Test case: invalid parameter type
		params = []any{object.MakeEmptyObject(), true}
		result = stringBuilderAppend(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderAppendBoolean(t *testing.T) {
		// Test case: append true
		params := []any{object.NewStringObject(), int64(types.JavaBoolTrue)}
		result := stringBuilderAppendBoolean(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result != "true" {
			t.Errorf("Expected 'true', got %s", result)
		}

		// Test case: append false
		params = []any{object.NewStringObject(), int64(types.JavaBoolFalse)}
		result = stringBuilderAppendBoolean(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result != "false" {
			t.Errorf("Expected 'false', got %s", result)
		}

		// Test case: invalid parameter type
		params = []any{object.NewStringObject(), "true"}
		result = stringBuilderAppendBoolean(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderAppendChar(t *testing.T) {
		// Test case: append char
		params := []any{object.NewStringObject(), int64('A')}
		result := stringBuilderAppendChar(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: invalid parameter type
		params = []any{object.NewStringObject(), "A"}
		result = stringBuilderAppendChar(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderCharAt(t *testing.T) {
		// Test case: valid index
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj, int64(1)}
		result := stringBuilderCharAt(params)
		if result != int64('e') {
			t.Errorf("Expected 'e', got %v", result)
		}

		// Test case: invalid index
		params = []any{obj, int64(10)}
		result = stringBuilderCharAt(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderDelete(t *testing.T) {
		// Test case: delete range
		obj := object.NewStringObject()
		obj.FieldTable["value"] =
			object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj, int64(1), int64(3)}
		result := stringBuilderDelete(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: delete char at
		params = []any{obj, int64(1)}
		result = stringBuilderDelete(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: invalid start index
		params = []any{obj, int64(10), int64(3)}
		result = stringBuilderDelete(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}

		// Test case: invalid end index
		params = []any{obj, int64(1), int64(10)}
		result = stringBuilderDelete(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderInsert(t *testing.T) {
		// Test case: insert string
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		strObj := object.StringObjectFromGoString("World")
		params := []any{obj, int64(1), strObj}
		result := stringBuilderInsert(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: insert int64
		params = []any{obj, int64(1), int64(123)}
		result = stringBuilderInsert(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: insert float64
		params = []any{obj, int64(1), float64(123.45)}
		result = stringBuilderInsert(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: invalid parameter type
		params = []any{obj, int64(1), true}
		result = stringBuilderInsert(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderInsertBoolean(t *testing.T) {
		// Test case: insert true
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj, int64(1), int64(types.JavaBoolTrue)}
		result := stringBuilderInsertBoolean(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: insert false
		params = []any{obj, int64(1), int64(types.JavaBoolFalse)}
		result = stringBuilderInsertBoolean(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: invalid parameter type
		params = []any{obj, int64(1), "true"}
		result = stringBuilderInsertBoolean(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderInsertChar(t *testing.T) {
		// Test case: insert char
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj, int64(1), int64('A')}
		result := stringBuilderInsertChar(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		// Test case: invalid parameter type
		params = []any{obj, int64(1), "A"}
		result = stringBuilderInsertChar(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderReplace(t *testing.T) {
		// Test case: replace substring
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		strObj := object.StringObjectFromGoString("World")
		params := []any{obj, int64(1), int64(3), strObj}
		result := stringBuilderReplace(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result != "HWorldo" {
			t.Errorf("Expected 'HWorldo', got %s", result)
		}

		// Test case: invalid start index
		params = []any{obj, int64(10), int64(3), strObj}
		result = stringBuilderReplace(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.StringIndexOutOfBoundsException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}

		// Test case: invalid end index
		params = []any{obj, int64(1), int64(10), strObj}
		result = stringBuilderReplace(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.StringIndexOutOfBoundsException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderReverse(t *testing.T) {
		// Test case: reverse string
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj}
		result := stringBuilderReverse(params)
		if result == nil {
			t.Errorf("Expected non-nil result, got nil")
		}

		if result != "olleH" {
			t.Errorf("Expected 'olleH', got %s", result)
		}
	}

	func TestStringBuilderSetCharAt(t *testing.T) {
		// Test case: set char at valid index
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj, int64(1), int64('A')}
		result := stringBuilderSetCharAt(params)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}

		// Test case: set char at invalid index
		params = []any{obj, int64(10), int64('A')}
		result = stringBuilderSetCharAt(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IndexOutOfBoundsException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}

	func TestStringBuilderSetLength(t *testing.T) {
		// Test case: set length greater than current length
		obj := object.NewStringObject()
		obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'e', 'l', 'l', 'o'}}
		params := []any{obj, int64(10)}
		result := stringBuilderSetLength(params)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}

		// Test case: set length less than current length
		params = []any{obj, int64(3)}
		result = stringBuilderSetLength(params)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}

		// Test case: set negative length
		params = []any{obj, int64(-1)}
		result = stringBuilderSetLength(params)
		if result == nil || result.(*GErrBlk).ExceptionType != excNames.IndexOutOfBoundsException {
			t.Errorf("Expected IllegalArgumentException, got %s",
				excNames.JVMexceptionNames[result.(*GErrBlk).ExceptionType])
		}
	}
*/
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
