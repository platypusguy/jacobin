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
	"jacobin/src/types"
	"testing"
)

func TestFill_Int(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimArray(object.T_INT, 5)
	arr := arrObj.FieldTable["value"].Fvalue.([]int64)

	// Fill all
	utilArraysFill([]interface{}{arrObj, int64(42)})
	for i, v := range arr {
		if v != 42 {
			t.Errorf("Expected 42 at index %d, got %v", i, v)
		}
	}

	// Fill range
	utilArraysFill([]interface{}{arrObj, int64(1), int64(4), int64(7)})
	if arr[0] != 42 || arr[1] != 7 || arr[2] != 7 || arr[3] != 7 || arr[4] != 42 {
		t.Errorf("Range fill failed: %v", arr)
	}
}

func TestFill_Object(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimRefArray("Ljava/lang/Object;", 3)
	arr := arrObj.FieldTable["value"].Fvalue.([]*object.Object)

	val := object.StringObjectFromGoString("hello")
	utilArraysFill([]interface{}{arrObj, val})
	for i, v := range arr {
		if v != val {
			t.Errorf("Expected 'hello' at index %d", i)
		}
	}
}

func TestFill_Boolean(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	arrObj := object.Make1DimArray(object.T_BOOLEAN, 3)
	arr := arrObj.FieldTable["value"].Fvalue.([]types.JavaByte)

	utilArraysFill([]interface{}{arrObj, types.JavaBoolTrue})
	for i, v := range arr {
		if v != types.JavaByte(types.JavaBoolTrue) {
			t.Errorf("Expected true at index %d", i)
		}
	}
}
