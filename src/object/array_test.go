/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/globals"
	"jacobin/stringPool"
	"testing"
)

// This file tests array primitives. Array bytecodes are tested in
// jvm.arrayByetcodes.go

func TestArrayTypeConversions(t *testing.T) {
	if JdkArrayTypeToJacobinType(T_BOOLEAN) != BYTE {
		t.Errorf("did not get expected Jacobin type BOOLEAN")
	}

	if JdkArrayTypeToJacobinType(T_LONG) != INT {
		t.Errorf("did not get expected Jacobin type for LONG")
	}

	if JdkArrayTypeToJacobinType(T_DOUBLE) != FLOAT {
		t.Errorf("did not get expected Jacobin type for DOUBLE")
	}

	if JdkArrayTypeToJacobinType(T_REF) != REF {
		t.Errorf("did not get expected Jacobin type for REF")
	}

	if JdkArrayTypeToJacobinType(99) != 0 {
		t.Errorf("did not get expected Jacobin type for invalid value")
	}
}

func TestMakde1DimByteArray(t *testing.T) {
	globals.InitGlobals("test")
	bArr := Make1DimArray(BYTE, 10)
	bArrType := stringPool.GetStringPointer(bArr.KlassName)
	if *bArrType != "[B" {
		t.Errorf("did not get expected Jacobin type for byte array, got %s", *bArrType)
	}

	rawArray := bArr.FieldTable["value"].Fvalue.([]byte)
	if len(rawArray) != 10 {
		t.Errorf("Expecting 10 elements in byte array, got %d", len(rawArray))
	}

	if rawArray[0] != byte(0) {
		t.Errorf("Expecting byte[0] ==  0, got %d", rawArray[0])
	}
}

func TestMakde1DimRefArray(t *testing.T) {
	globals.InitGlobals("test")
	rArr := Make1DimArray(REF, 10)
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
	iArr := Make1DimArray(INT, 10)
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
	objType := "genericObj"
	rArr := Make1DimRefArray(&objType, 11)
	size := ArrayLength(rArr)
	if size != 11 {
		t.Errorf("Expecting 11 elements in array, got %d", size)
	}
}

func TestMake2DimByteArray(t *testing.T) {
	globals.InitGlobals("test")
	b2arr, ok := Make2DimArray(10, 15, BYTE)
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

func TestArrayLength(t *testing.T) {
	globals.InitGlobals("test")
	iArr := Make1DimArray(INT, 256)
	length := ArrayLength(iArr)
	if length != 256 {
		t.Errorf("Expecting 256 elements in int array, got %d", length)
	}

	bArr := Make1DimArray(BYTE, 200)
	length = ArrayLength(bArr)
	if length != 200 {
		t.Errorf("Expecting 200 elements in byte array, got %d", length)
	}

	fArr := Make1DimArray(FLOAT, 222)
	length = ArrayLength(fArr)
	if length != 222 {
		t.Errorf("Expecting 222 elements in float array, got %d", length)
	}

	rArr := Make1DimArray(REF, 256)
	length = ArrayLength(rArr)
	if length != 256 {
		t.Errorf("Expecting 256 elements in ref array, got %d", length)
	}
}
