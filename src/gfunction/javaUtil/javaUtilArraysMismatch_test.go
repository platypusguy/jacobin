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
	"testing"
)

func TestMismatch_Int(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_INT, 3)
	v1 := arr1.FieldTable["value"].Fvalue.([]int64)
	v1[0], v1[1], v1[2] = 1, 2, 3

	arr2 := object.Make1DimArray(object.T_INT, 3)
	v2 := arr2.FieldTable["value"].Fvalue.([]int64)
	v2[0], v2[1], v2[2] = 1, 2, 3

	// Test no mismatch
	res := utilArraysMismatch([]interface{}{arr1, arr2})
	if res.(int64) != -1 {
		t.Errorf("Expected -1, got %v", res)
	}

	// Test mismatch at index 1
	v2[1] = 9
	res = utilArraysMismatch([]interface{}{arr1, arr2})
	if res.(int64) != 1 {
		t.Errorf("Expected 1, got %v", res)
	}

	// Test different lengths (prefix)
	arr3 := object.Make1DimArray(object.T_INT, 2)
	v3 := arr3.FieldTable["value"].Fvalue.([]int64)
	v3[0], v3[1] = 1, 2
	res = utilArraysMismatch([]interface{}{arr1, arr3})
	if res.(int64) != 2 {
		t.Errorf("Expected 2 (length of shorter array), got %v", res)
	}
}

func TestMismatch_Range(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arr1 := object.Make1DimArray(object.T_INT, 5)
	v1 := arr1.FieldTable["value"].Fvalue.([]int64)
	v1[0], v1[1], v1[2], v1[3], v1[4] = 0, 1, 2, 3, 0

	arr2 := object.Make1DimArray(object.T_INT, 5)
	v2 := arr2.FieldTable["value"].Fvalue.([]int64)
	v2[0], v2[1], v2[2], v2[3], v2[4] = 9, 1, 2, 3, 9

	// Compare range [1, 4) of both
	res := utilArraysMismatch([]interface{}{arr1, int64(1), int64(4), arr2, int64(1), int64(4)})
	if err, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("utilArraysMismatch returned error: %v: %s", err.ExceptionType, err.ErrMsg)
	}
	if res.(int64) != -1 {
		t.Errorf("Expected -1, got %v", res)
	}
}
