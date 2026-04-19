/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"bytes"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestIvParameterSpec(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/spec/IvParameterSpec"
	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}
	ivObj := makeByteArrayObject(iv)

	// Test IvParameterSpec(byte[] iv)
	spec1 := object.MakeEmptyObjectWithClassName(&className)
	res1 := ivParameterSpecInit([]any{spec1, ivObj})
	if res1 != nil {
		t.Errorf("Expected nil result, got %v", res1)
	}

	ivStored := spec1.FieldTable["iv"].Fvalue.([]byte)
	if !bytes.Equal(ivStored, iv) {
		t.Errorf("Expected IV %v, got %v", iv, ivStored)
	}

	// Test getIV()
	resIV := ivParameterSpecGetIV([]any{spec1})
	resIVObj, ok := resIV.(*object.Object)
	if !ok || resIVObj == nil {
		t.Fatalf("Expected IV object, got %v", resIV)
	}
	resIVJBytes := resIVObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	resIVGoBytes := object.GoByteArrayFromJavaByteArray(resIVJBytes)
	if !bytes.Equal(resIVGoBytes, iv) {
		t.Errorf("Expected getIV() %v, got %v", iv, resIVGoBytes)
	}

	// Verify getIV returns a copy
	resIVGoBytes[0] ^= 0xFF
	if bytes.Equal(spec1.FieldTable["iv"].Fvalue.([]byte), resIVGoBytes) {
		t.Error("getIV() should return a copy")
	}

	// Test IvParameterSpec(byte[] iv, int offset, int len)
	spec2 := object.MakeEmptyObjectWithClassName(&className)
	offset := int64(4)
	length := int64(8)
	res2 := ivParameterSpecInit([]any{spec2, ivObj, offset, length})
	if res2 != nil {
		t.Errorf("Expected nil result, got %v", res2)
	}
	ivStored2 := spec2.FieldTable["iv"].Fvalue.([]byte)
	expected2 := iv[offset : offset+length]
	if !bytes.Equal(ivStored2, expected2) {
		t.Errorf("Expected subset IV %v, got %v", expected2, ivStored2)
	}

	// Test null IV error
	spec3 := object.MakeEmptyObjectWithClassName(&className)
	res3 := ivParameterSpecInit([]any{spec3, object.Null})
	errBlk, ok := res3.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null IV, got %v", res3)
	}

	// Test empty IV error
	spec4 := object.MakeEmptyObjectWithClassName(&className)
	emptyIvObj := makeByteArrayObject([]byte{})
	res4 := ivParameterSpecInit([]any{spec4, emptyIvObj})
	errBlk, ok = res4.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for empty IV, got %v", res4)
	}

	// Test invalid offset/length error
	spec5 := object.MakeEmptyObjectWithClassName(&className)
	res5 := ivParameterSpecInit([]any{spec5, ivObj, int64(10), int64(10)})
	errBlk, ok = res5.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid offset/length, got %v", res5)
	}
}
