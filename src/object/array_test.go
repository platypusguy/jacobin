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
	bArr := Make1DimArray(REF, 10)
	bArrType := stringPool.GetStringPointer(bArr.KlassName)
	if *bArrType != "[L" {
		t.Errorf("did not get expected Jacobin type for byte array, got %s", *bArrType)
	}

	rawArray := bArr.FieldTable["value"].Fvalue.([]*Object)
	if len(rawArray) != 10 {
		t.Errorf("Expecting 10 elements in ref array, got %d", len(rawArray))
	}

	if rawArray[0] != nil {
		t.Errorf("Expecting ref[0] ==  nil, got %v", rawArray[0])
	}
}
