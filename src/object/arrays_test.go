/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"fmt"
	"io"
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
	"testing"
)

// This file tests array primitives. Array bytecodes are tested in
// jvm.arrayByetcodes.go

func TestGetArrayType(t *testing.T) {
	retVal := GetArrayType("[[B")
	if retVal != "B" {
		t.Errorf("did not get expected 'B' for [[B, got: %s", retVal)
	}

	retVal = GetArrayType("I")
	if retVal != "I" {
		t.Errorf("did not get expected 'I' for I, got: %s", retVal)
	}

	retVal = GetArrayType("[Ljava/lang/Object;")
	if retVal != "Ljava/lang/Object;" {
		t.Errorf("did not get expected 'Ljava/lang/Object;', got: %s", retVal)
	}

	retVal = GetArrayType("")
	if retVal != "" {
		t.Errorf("did not get expected empty string for an empty string, got %s", retVal)
	}
}

func TestMakde1DimByteArray(t *testing.T) {
	globals.InitGlobals("test")
	bArr := Make1DimArray(T_BYTE, 10)
	bArrType := stringPool.GetStringPointer(bArr.KlassName)
	if *bArrType != "[B" {
		t.Errorf("did not get expected Jacobin type for byte array, got %s", *bArrType)
	}

	rawArray := bArr.FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(rawArray) != 10 {
		t.Errorf("Expecting 10 elements in byte array, got %d", len(rawArray))
	}

	if rawArray[0] != types.JavaByte(0) {
		t.Errorf("Expecting byte[0] ==  0, got %d", rawArray[0])
	}
}

func TestMakde1DimRefArray(t *testing.T) {
	globals.InitGlobals("test")
	rArr := Make1DimArray(T_REF, 10)
	rArrType := stringPool.GetStringPointer(rArr.KlassName)
	if *rArrType != "[L" {
		t.Errorf("did not get expected Jacobin type for ref array, got %s", *rArrType)
	}

	rawArray := rArr.FieldTable["value"].Fvalue.([]*Object)
	if len(rawArray) != 10 {
		t.Errorf("Expecting 10 elements in ref array, got %d", len(rawArray))
	}

	if rawArray[0] != nil {
		t.Errorf("Expecting ref[0] ==  nil, got %v", rawArray[0])
	}
}

func TestMakde1DimIntArray(t *testing.T) {
	globals.InitGlobals("test")
	iArr := Make1DimArray(T_INT, 10)
	iArrType := stringPool.GetStringPointer(iArr.KlassName)
	if *iArrType != "[I" {
		t.Errorf("did not get expected Jacobin type for int array, got %s", *iArrType)
	}

	rawArray := iArr.FieldTable["value"].Fvalue.([]int64)
	if len(rawArray) != 10 {
		t.Errorf("Expecting 10 elements in int array, got %d", len(rawArray))
	}

	if rawArray[0] != int64(0) {
		t.Errorf("Expecting int[0] ==  0, got %d", rawArray[0])
	}
}

func TestMake1DimRefArray(t *testing.T) {
	globals.InitGlobals("test")
	rArr := Make1DimRefArray("genericObj", 11)
	size := ArrayLength(rArr)
	if size != 11 {
		t.Errorf("Expecting 11 elements in array, got %d", size)
	}
}

func TestMake2DimByteArray(t *testing.T) {
	globals.InitGlobals("test")
	b2arr, ok := Make2DimArray(10, 15, T_BYTE)
	if ok != nil {
		t.Errorf("Unexpected error in Make2DimArray: %s", ok.Error())
	}

	arrType := *(stringPool.GetStringPointer(b2arr.KlassName))
	if arrType != "[B" {
		t.Errorf("Did not get expected Jacobin type for byte array, got %s", arrType)
	}

	topArray := b2arr.FieldTable["value"]
	if topArray.Ftype != "[[B" {
		t.Errorf("Expecting top array type to be [[B, got %s)", topArray.Ftype)
	}

	// now test the size of the leaf arrays
	leafArray := topArray.Fvalue.([]*Object)[0]
	length := ArrayLength(leafArray)
	if length != 15 {
		t.Errorf("Expecting length 15, got %d", length)
	}
}

func TestMakeArrayFromRawArray(t *testing.T) {
	globals.InitGlobals("test")
	rawArray := make([]types.JavaByte, 10)
	newArr := MakeArrayFromRawArray(rawArray)
	size := ArrayLength(newArr)
	if size != 10 {
		t.Errorf("Expecting length 10, got %d", size)
	}
}

func TestMakeArrayFromRawArrayPtr(t *testing.T) {
	globals.InitGlobals("test")
	rawArray := make([]types.JavaByte, 10)
	newArr := MakeArrayFromRawArray(&rawArray)
	size := ArrayLength(newArr)
	if size != 10 {
		t.Errorf("Expecting length 10, got %d", size)
	}
}

func TestMakeArrayFromObjectPtr(t *testing.T) {
	globals.InitGlobals("test")
	o := MakeEmptyObject()
	newArr := MakeArrayFromRawArray(o)
	if newArr != o {
		t.Errorf("Expecting length output = input, got %v", newArr)
	}
}

func TestMakeArrayFromRawArrayNil(t *testing.T) {
	globals.InitGlobals("test")
	newArr := MakeArrayFromRawArray(nil)
	if newArr != nil {
		t.Errorf("Expecting nil, got %v", newArr)
	}
}

func TestMakeArrayFromRawArrayInvalidInput(t *testing.T) {
	globals.InitGlobals("test")
	globals.GetGlobalRef().FuncThrowException = announceError

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	newArr := MakeArrayFromRawArray(10)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "Error in tested function\n" {
		t.Errorf("Expecting error message 'Error in tested function\n', got %s", errMsg)
	}
	if newArr != nil {
		t.Errorf("Expecting nil, got %v", newArr)
	}
}

// used instead of throwing an exception, which creates a circularity problem
func announceError(_ int, _ string) bool {
	fmt.Fprintln(os.Stderr, "Error in tested function")
	return true
}

func TestArrayLength(t *testing.T) {
	globals.InitGlobals("test")
	iArr := Make1DimArray(T_INT, 256)
	length := ArrayLength(iArr)
	if length != 256 {
		t.Errorf("Expecting 256 elements in int array, got %d", length)
	}

	bArr := Make1DimArray(T_BYTE, 200)
	length = ArrayLength(bArr)
	if length != 200 {
		t.Errorf("Expecting 200 elements in byte array, got %d", length)
	}

	fArr := Make1DimArray(T_FLOAT, 222)
	length = ArrayLength(fArr)
	if length != 222 {
		t.Errorf("Expecting 222 elements in float array, got %d", length)
	}

	rArr := Make1DimArray(T_REF, 256)
	length = ArrayLength(rArr)
	if length != 256 {
		t.Errorf("Expecting 256 elements in ref array, got %d", length)
	}
}
