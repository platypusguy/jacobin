/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"testing"
)

func TestSort_Int(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimArray(object.T_INT, 5)
	arr := arrObj.FieldTable["value"].Fvalue.([]int64)
	arr[0], arr[1], arr[2], arr[3], arr[4] = 50, 10, 40, 20, 30

	utilArraysSort([]interface{}{arrObj})
	if arr[0] != 10 || arr[1] != 20 || arr[2] != 30 || arr[3] != 40 || arr[4] != 50 {
		t.Errorf("Sort failed: %v", arr)
	}
}

func TestSort_Range(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimArray(object.T_INT, 5)
	arr := arrObj.FieldTable["value"].Fvalue.([]int64)
	arr[0], arr[1], arr[2], arr[3], arr[4] = 50, 40, 30, 20, 10

	// Sort range [1, 4) -> index 1, 2, 3. Values: 40, 30, 20
	utilArraysSort([]interface{}{arrObj, int64(1), int64(4)})
	if arr[0] != 50 || arr[1] != 20 || arr[2] != 30 || arr[3] != 40 || arr[4] != 10 {
		t.Errorf("Range sort failed: %v", arr)
	}
}

func TestSort_String(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimRefArray("Ljava/lang/String;", 3)
	arr := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	arr[0] = object.StringObjectFromGoString("cherry")
	arr[1] = object.StringObjectFromGoString("apple")
	arr[2] = object.StringObjectFromGoString("banana")

	utilArraysSort([]interface{}{arrObj})
	if object.GoStringFromStringObject(arr[0]) != "apple" ||
		object.GoStringFromStringObject(arr[1]) != "banana" ||
		object.GoStringFromStringObject(arr[2]) != "cherry" {
		t.Errorf("String sort failed")
	}
}
