/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"testing"
)

func TestCompare_Int(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_INT, 3)
	v1 := arr1.FieldTable["value"].Fvalue.([]int64)
	v1[0], v1[1], v1[2] = 1, 2, 3

	arr2 := object.Make1DimArray(object.T_INT, 3)
	v2 := arr2.FieldTable["value"].Fvalue.([]int64)
	v2[0], v2[1], v2[2] = 1, 2, 3

	// Test equal
	res := utilArraysCompare([]interface{}{arr1, arr2})
	if res.(int64) != 0 {
		t.Errorf("Expected 0, got %v", res)
	}

	// Test different
	v2[2] = 4
	res = utilArraysCompare([]interface{}{arr1, arr2})
	if res.(int64) >= 0 {
		t.Errorf("Expected negative, got %v", res)
	}

	// Test different length
	arr3 := object.Make1DimArray(object.T_INT, 2)
	v3 := arr3.FieldTable["value"].Fvalue.([]int64)
	v3[0], v3[1] = 1, 2
	res = utilArraysCompare([]interface{}{arr1, arr3})
	if res.(int64) <= 0 {
		t.Errorf("Expected positive, got %v", res)
	}
}

func TestCompare_Boolean(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_BOOLEAN, 2)
	v1 := arr1.FieldTable["value"].Fvalue.([]types.JavaByte)
	v1[0], v1[1] = types.JavaByte(types.JavaBoolFalse), types.JavaByte(types.JavaBoolTrue)

	arr2 := object.Make1DimArray(object.T_BOOLEAN, 2)
	v2 := arr2.FieldTable["value"].Fvalue.([]types.JavaByte)
	v2[0], v2[1] = types.JavaByte(types.JavaBoolFalse), types.JavaByte(types.JavaBoolFalse)

	// false < true, so arr1[1] > arr2[1]
	res := utilArraysCompare([]interface{}{arr1, arr2})
	if res.(int64) <= 0 {
		t.Errorf("Expected positive, got %v", res)
	}
}

func TestCompare_Range(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_INT, 5)
	v1 := arr1.FieldTable["value"].Fvalue.([]int64)
	v1[0], v1[1], v1[2], v1[3], v1[4] = 0, 1, 2, 3, 0

	arr2 := object.Make1DimArray(object.T_INT, 5)
	v2 := arr2.FieldTable["value"].Fvalue.([]int64)
	v2[0], v2[1], v2[2], v2[3], v2[4] = 9, 1, 2, 3, 9

	// Compare range [1, 4) of both
	res := utilArraysCompare([]interface{}{arr1, int64(1), int64(4), arr2, int64(1), int64(4)})
	if err, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("utilArraysCompare returned error: %v: %s", err.ExceptionType, err.ErrMsg)
	}
	if res.(int64) != 0 {
		t.Errorf("Expected 0, got %v", res)
	}
}

func TestCompare_Nulls(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_INT, 3)

	// Both null
	res := utilArraysCompare([]interface{}{nil, nil})
	if res.(int64) != 0 {
		t.Errorf("Expected 0")
	}

	// First null
	res = utilArraysCompare([]interface{}{nil, arr1})
	if res.(int64) != -1 {
		t.Errorf("Expected -1")
	}

	// Second null
	res = utilArraysCompare([]interface{}{arr1, nil})
	if res.(int64) != 1 {
		t.Errorf("Expected 1")
	}
}
