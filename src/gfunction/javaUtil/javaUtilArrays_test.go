/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"testing"
)

func TestCopyOfObjectPointers_TooFewArguments(t *testing.T) {
	result := *(utilArraysCopyOf([]interface{}{}).(*ghelpers.GErrBlk))
	if result.ExceptionType != excNames.IllegalArgumentException || result.ErrMsg != "utilArraysCopyOf: too few arguments" {
		t.Errorf("Expected IllegalArgumentException for too few arguments")
	}
}

func TestCopyOfObjectPointers_NullArray(t *testing.T) {
	result := *(utilArraysCopyOf([]interface{}{nil, int64(5)}).(*ghelpers.GErrBlk))
	if result.ExceptionType != excNames.NullPointerException || result.ErrMsg != "utilArraysCopyOf: null array argument" {
		t.Errorf("Expected NullPointerException for null array argument")
	}
}

func TestCopyOfObjectPointers_NegativeLength(t *testing.T) {
	obj := object.MakeEmptyObject()
	result := *(utilArraysCopyOf([]interface{}{obj, int64(-1)}).(*ghelpers.GErrBlk))
	if result.ExceptionType != excNames.NegativeArraySizeException || result.ErrMsg != "utilArraysCopyOf: negative array length" {
		// if result != ghelpers.GetGErrBlk(excNames.NegativeArraySizeException, "copyOf: negative array length") {
		t.Errorf("Expected NegativeArraySizeException for negative array length")
	}
}

func TestCopyOfObjectPointers_CopyArray(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool() // needed for creation of array

	// Create a mock array object
	oldArray := object.Make1DimRefArray("Ljava/lang/Object;", 2)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]*object.Object)
	rawOldArray[0] = object.StringObjectFromGoString("foo")
	rawOldArray[1] = object.StringObjectFromGoString("bar")

	// Test copying to a larger array
	result := utilArraysCopyOf([]interface{}{oldArray, int64(4)})
	newArray := result.(*object.Object).FieldTable["value"].Fvalue.([]*object.Object)
	if len(newArray) != 4 {
		t.Errorf("Expected new array length of 4, got %d", len(newArray))
	}

	if len(newArray) != 4 {
		t.Errorf("Expected new array length of 4, got %d", len(newArray))
	}

	if object.GoStringFromStringObject(newArray[0]) != "foo" || object.GoStringFromStringObject(newArray[1]) != "bar" {
		t.Errorf("Array elements not copied correctly")
	}
}

func TestCopyOfPrimitive_Int(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// Create a mock int array object [I
	oldArray := object.Make1DimArray(object.T_INT, 3)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]int64)
	rawOldArray[0] = 10
	rawOldArray[1] = 20
	rawOldArray[2] = 30

	// Test copying to a larger array
	result := utilArraysCopyOfPrimitive([]interface{}{oldArray, int64(5)})
	resObj := result.(*object.Object)
	newArray := resObj.FieldTable["value"].Fvalue.([]int64)

	if len(newArray) != 5 {
		t.Errorf("Expected new array length of 5, got %d", len(newArray))
	}
	if newArray[0] != 10 || newArray[1] != 20 || newArray[2] != 30 || newArray[3] != 0 || newArray[4] != 0 {
		t.Errorf("Array elements not copied correctly: %v", newArray)
	}

	// Verify type preservation
	if resObj.FieldTable["value"].Ftype != types.IntArray {
		t.Errorf("Expected type %s, got %s", types.IntArray, resObj.FieldTable["value"].Ftype)
	}
}

func TestCopyOfPrimitive_Byte(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// Create a mock byte array object [B
	oldArray := object.Make1DimArray(object.T_BYTE, 2)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]types.JavaByte)
	rawOldArray[0] = 1
	rawOldArray[1] = 2

	// Test copying to a smaller array
	result := utilArraysCopyOfPrimitive([]interface{}{oldArray, int64(1)})
	resObj := result.(*object.Object)
	newArray := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)

	if len(newArray) != 1 {
		t.Errorf("Expected new array length of 1, got %d", len(newArray))
	}
	if newArray[0] != 1 {
		t.Errorf("Array elements not copied correctly: %v", newArray)
	}

	// Verify type preservation
	if resObj.FieldTable["value"].Ftype != types.ByteArray {
		t.Errorf("Expected type %s, got %s", types.ByteArray, resObj.FieldTable["value"].Ftype)
	}
}

func TestCopyOfPrimitive_Boolean(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// Create a mock boolean array object [Z
	oldArray := object.Make1DimArray(object.T_BOOLEAN, 2)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]types.JavaByte)
	rawOldArray[0] = 1 // true
	rawOldArray[1] = 0 // false

	// Test copying to a larger array
	result := utilArraysCopyOfPrimitive([]interface{}{oldArray, int64(3)})
	resObj := result.(*object.Object)
	newArray := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)

	if len(newArray) != 3 {
		t.Errorf("Expected new array length of 3, got %d", len(newArray))
	}
	if newArray[0] != 1 || newArray[1] != 0 || newArray[2] != 0 {
		t.Errorf("Array elements not copied correctly: %v", newArray)
	}

	// Verify type preservation
	if resObj.FieldTable["value"].Ftype != types.BoolArray {
		t.Errorf("Expected type %s, got %s", types.BoolArray, resObj.FieldTable["value"].Ftype)
	}
}

func TestCopyOfPrimitive_Float(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// Create a mock float array object [F
	oldArray := object.Make1DimArray(object.T_FLOAT, 2)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]float64)
	rawOldArray[0] = 1.5
	rawOldArray[1] = 2.5

	// Test copying
	result := utilArraysCopyOfPrimitive([]interface{}{oldArray, int64(2)})
	resObj := result.(*object.Object)
	newArray := resObj.FieldTable["value"].Fvalue.([]float64)

	if len(newArray) != 2 {
		t.Errorf("Expected new array length of 2, got %d", len(newArray))
	}
	if newArray[0] != 1.5 || newArray[1] != 2.5 {
		t.Errorf("Array elements not copied correctly: %v", newArray)
	}

	// Verify type preservation
	if resObj.FieldTable["value"].Ftype != types.FloatArray {
		t.Errorf("Expected type %s, got %s", types.FloatArray, resObj.FieldTable["value"].Ftype)
	}
}

func TestCopyOfObjectWithClass(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// Create a mock array object [Ljava/lang/Object;
	oldArray := object.Make1DimRefArray("Ljava/lang/Object;", 2)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]*object.Object)
	rawOldArray[0] = object.StringObjectFromGoString("one")
	rawOldArray[1] = object.StringObjectFromGoString("two")

	// Dummy class object
	classClassName := "java/lang/Class"
	classObj := object.MakeEmptyObjectWithClassName(&classClassName)

	// Test copying with 3 arguments
	result := utilArraysCopyOfObjectWithClass([]interface{}{oldArray, int64(3), classObj})
	resObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	newArray := resObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(newArray) != 3 {
		t.Errorf("Expected length 3, got %d", len(newArray))
	}
	if object.GoStringFromStringObject(newArray[0]) != "one" || object.GoStringFromStringObject(newArray[1]) != "two" || newArray[2] != nil {
		t.Errorf("Elements not copied correctly")
	}
}

func TestCopyOf_UnsupportedType(t *testing.T) {
	globals.InitGlobals("test")

	// Create an object that is NOT an array (missing "value" field or wrong type in value field)
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: "invalid", Fvalue: "not a slice"}

	result := utilArraysCopyOf([]interface{}{obj, int64(5)})
	err, ok := result.(*ghelpers.GErrBlk)
	if !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for unsupported array type, got %v", result)
	}
}

func TestArraysAsList(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// 1. Create a mock array object [Ljava/lang/Object;
	arrayObj := object.Make1DimRefArray("Ljava/lang/Object;", 3)
	rawArray := arrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")
	rawArray[0] = s1
	rawArray[1] = s2
	rawArray[2] = nil // Arrays.asList allows nulls

	// 2. Call utilArraysAsList
	res := utilArraysAsList([]interface{}{arrayObj})
	listObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}

	// 3. Verify it's an ArrayList
	if className := object.GoStringFromStringPoolIndex(listObj.KlassName); className != "java/util/ArrayList" {
		t.Errorf("expected className java/util/ArrayList, got %s", className)
	}

	// 4. Verify elements
	al, err := GetArrayListFromObject(listObj)
	if err != nil {
		t.Fatalf("getArrayListFromObject failed: %v", err)
	}

	if len(al) != 3 {
		t.Errorf("expected size 3, got %d", len(al))
	}

	if al[0] != s1 || al[1] != s2 || al[2] != object.Null {
		t.Errorf("elements mismatch")
	}
}

func TestArraysEquals(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_INT, 3)
	v1 := arr1.FieldTable["value"].Fvalue.([]int64)
	v1[0], v1[1], v1[2] = 1, 2, 3

	arr2 := object.Make1DimArray(object.T_INT, 3)
	v2 := arr2.FieldTable["value"].Fvalue.([]int64)
	v2[0], v2[1], v2[2] = 1, 2, 3

	// Equal
	if utilArraysEquals([]interface{}{arr1, arr2}) != types.JavaBoolTrue {
		t.Errorf("Expected true for equal arrays")
	}

	// Not equal
	v2[2] = 9
	if utilArraysEquals([]interface{}{arr1, arr2}) != types.JavaBoolFalse {
		t.Errorf("Expected false for unequal arrays")
	}

	// Different types
	arr3 := object.Make1DimArray(object.T_BYTE, 3)
	if utilArraysEquals([]interface{}{arr1, arr3}) != types.JavaBoolFalse {
		t.Errorf("Expected false for different array types")
	}
}

func TestArraysToString(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimArray(object.T_INT, 3)
	arr := arrObj.FieldTable["value"].Fvalue.([]int64)
	arr[0], arr[1], arr[2] = 1, 2, 3

	res := utilArraysToString([]interface{}{arrObj})
	strObj := res.(*object.Object)
	goStr := object.GoStringFromStringObject(strObj)
	if goStr != "[1, 2, 3]" {
		t.Errorf("Expected '[1, 2, 3]', got '%s'", goStr)
	}

	// Null
	res = utilArraysToString([]interface{}{nil})
	if object.GoStringFromStringObject(res.(*object.Object)) != "null" {
		t.Errorf("Expected 'null'")
	}
}
