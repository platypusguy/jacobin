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

func makeByteArrayObject(b []byte) *object.Object {
	jBytes := object.JavaByteArrayFromGoByteArray(b)
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}

func TestGCMParameterSpec(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/spec/GCMParameterSpec"
	iv := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}
	ivObj := makeByteArrayObject(iv)
	tLen := int64(128)

	// Test GCMParameterSpec(int tLen, byte[] src)
	specObj1 := object.MakeEmptyObjectWithClassName(&className)
	params1 := []any{specObj1, tLen, ivObj}
	result1 := gcmParameterSpecInit(params1)
	if result1 != nil {
		t.Errorf("Expected nil result for GCMParameterSpecInit, got %v", result1)
	}

	// Test getTLen()
	resTLen := gcmParameterSpecGetTLen([]any{specObj1})
	if resTLen.(int64) != tLen {
		t.Errorf("Expected tLen %d, got %v", tLen, resTLen)
	}

	// Test getIV()
	resIV1 := gcmParameterSpecGetIV([]any{specObj1})
	resIV1Obj, ok := resIV1.(*object.Object)
	if !ok || resIV1Obj == nil {
		t.Fatalf("Expected IV object for getIV(), got %v", resIV1)
	}
	resIV1JBytes := resIV1Obj.FieldTable["value"].Fvalue.([]types.JavaByte)
	resIV1Bytes := object.GoByteArrayFromJavaByteArray(resIV1JBytes)
	if !bytes.Equal(resIV1Bytes, iv) {
		t.Errorf("Expected IV %v, got %v", iv, resIV1Bytes)
	}

	// Verify getIV() returns a copy
	resIV1Bytes[0] ^= 0xFF
	resIV1Again := gcmParameterSpecGetIV([]any{specObj1})
	resIV1AgainBytes := object.GoByteArrayFromJavaByteArray(resIV1Again.(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte))
	if bytes.Equal(resIV1AgainBytes, resIV1Bytes) {
		t.Error("getIV() should return a copy, but modifying the result affected the original")
	}

	// Test GCMParameterSpec(int tLen, byte[] src, int offset, int len)
	specObj2 := object.MakeEmptyObjectWithClassName(&className)
	offset := int64(2)
	length := int64(8)
	params2 := []any{specObj2, tLen, ivObj, offset, length}
	result2 := gcmParameterSpecInit(params2)
	if result2 != nil {
		t.Errorf("Expected nil result for GCMParameterSpecInit with offset/len, got %v", result2)
	}

	expectedSubIv := iv[offset : offset+length]
	resIV2 := gcmParameterSpecGetIV([]any{specObj2})
	resIV2Bytes := object.GoByteArrayFromJavaByteArray(resIV2.(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(resIV2Bytes, expectedSubIv) {
		t.Errorf("Expected subset IV %v, got %v", expectedSubIv, resIV2Bytes)
	}

	// Error Case: null IV
	specObj3 := object.MakeEmptyObjectWithClassName(&className)
	params3 := []any{specObj3, tLen, object.Null}
	result3 := gcmParameterSpecInit(params3)
	errBlk, ok := result3.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null IV, got %v", result3)
	}

	// Error Case: empty IV
	emptyIvObj := makeByteArrayObject([]byte{})
	specObj4 := object.MakeEmptyObjectWithClassName(&className)
	params4 := []any{specObj4, tLen, emptyIvObj}
	result4 := gcmParameterSpecInit(params4)
	errBlk, ok = result4.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for empty IV, got %v", result4)
	}

	// Error Case: invalid offset/len
	specObj5 := object.MakeEmptyObjectWithClassName(&className)
	params5 := []any{specObj5, tLen, ivObj, int64(10), int64(5)}
	result5 := gcmParameterSpecInit(params5)
	errBlk, ok = result5.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid offset/len, got %v", result5)
	}

	// Error Case: wrong number of parameters
	specObj6 := object.MakeEmptyObjectWithClassName(&className)
	params6 := []any{specObj6, tLen, ivObj, int64(1)}
	result6 := gcmParameterSpecInit(params6)
	errBlk, ok = result6.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for wrong number of params, got %v", result6)
	}
}

func TestGCMParameterSpecMethodsErrors(t *testing.T) {
	globals.InitGlobals("test")

	// Test missing 'this' for GetIV
	res := gcmParameterSpecGetIV([]any{})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for missing 'this', got %v", res)
	}

	// Test missing 'this' for GetTLen
	res = gcmParameterSpecGetTLen([]any{})
	errBlk, ok = res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for missing 'this', got %v", res)
	}

	// Test invalid 'this' for Init
	res = gcmParameterSpecInit([]any{nil, int64(128), object.Null})
	errBlk, ok = res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for invalid 'this', got %v", res)
	}
}
