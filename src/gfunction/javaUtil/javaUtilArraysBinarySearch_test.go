/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
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

func TestBinarySearch_Int(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// Create a mock int array [I
	arrObj := object.Make1DimArray(object.T_INT, 5)
	arr := arrObj.FieldTable["value"].Fvalue.([]int64)
	arr[0], arr[1], arr[2], arr[3], arr[4] = 10, 20, 30, 40, 50

	// Test found
	res := utilArraysBinarySearch([]interface{}{arrObj, int64(30)})
	if res.(int64) != 2 {
		t.Errorf("Expected 2, got %v", res)
	}

	// Test not found (insertion point 1, returns -1-1 = -2)
	res = utilArraysBinarySearch([]interface{}{arrObj, int64(15)})
	if res.(int64) != -2 {
		t.Errorf("Expected -2, got %v", res)
	}

	// Test range search
	res = utilArraysBinarySearch([]interface{}{arrObj, int64(1), int64(4), int64(30)})
	if res.(int64) != 2 {
		t.Errorf("Expected 2, got %v", res)
	}

	res = utilArraysBinarySearch([]interface{}{arrObj, int64(1), int64(2), int64(30)})
	if res.(int64) >= 0 {
		t.Errorf("Expected negative, got %v", res)
	}
}

func TestBinarySearch_Byte(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimArray(object.T_BYTE, 3)
	arr := arrObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	arr[0], arr[1], arr[2] = 1, 5, 10

	res := utilArraysBinarySearch([]interface{}{arrObj, int64(5)})
	if res.(int64) != 1 {
		t.Errorf("Expected 1, got %v", res)
	}
}

func TestBinarySearch_String(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimRefArray("Ljava/lang/String;", 3)
	arr := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	arr[0] = object.StringObjectFromGoString("apple")
	arr[1] = object.StringObjectFromGoString("banana")
	arr[2] = object.StringObjectFromGoString("cherry")

	key := object.StringObjectFromGoString("banana")
	res := utilArraysBinarySearch([]interface{}{arrObj, key})
	if res.(int64) != 1 {
		t.Errorf("Expected 1, got %v", res)
	}

	key2 := object.StringObjectFromGoString("date")
	res = utilArraysBinarySearch([]interface{}{arrObj, key2})
	if res.(int64) >= 0 {
		t.Errorf("Expected negative, got %v", res)
	}
}

func TestBinarySearch_Errors(t *testing.T) {
	globals.InitGlobals("test")

	// Too few arguments
	res := utilArraysBinarySearch([]interface{}{})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException")
	}

	// Null array
	res = utilArraysBinarySearch([]interface{}{nil, int64(1)})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException")
	}
}
