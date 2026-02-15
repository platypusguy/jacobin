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

func TestSecretKeySpecInit(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/SecretKeySpec"
	key := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
		'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
		'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	}
	keyObj := makeByteArrayObject(key)
	algo := "AES/ECB/PKCS5Padding"
	algoObj := object.StringObjectFromGoString(algo)

	// Test SecretKeySpec(byte[] key, String algorithm)
	specObj1 := object.MakeEmptyObjectWithClassName(&className)
	params1 := []any{specObj1, keyObj, algoObj}
	result1 := secretKeySpecInit(params1)
	if result1 != nil {
		t.Errorf("Expected nil result, got %v", result1)
	}

	specAlgoObj := specObj1.FieldTable["algorithm"].Fvalue.(*object.Object)
	specAlgo := object.GoStringFromStringObject(specAlgoObj)
	if specAlgo != algo {
		t.Errorf("Expected algorithm %s, got %s", algo, specAlgo)
	}

	specKey := specObj1.FieldTable["key"].Fvalue.([]byte)
	if !bytes.Equal(specKey, key) {
		t.Errorf("Expected key %v, got %v", key, specKey)
	}

	// Test SecretKeySpec(byte[] key, int offset, int len, String algorithm)
	specObj2 := object.MakeEmptyObjectWithClassName(&className)
	offset := int64(2)
	length := int64(16)
	params2 := []any{specObj2, keyObj, offset, length, algoObj}
	result2 := secretKeySpecInit(params2)
	if result2 != nil {
		t.Errorf("Expected nil result, got %v", result2)
	}

	expectedSubKey := key[offset : offset+length]
	observedSubKey := specObj2.FieldTable["key"].Fvalue.([]byte)
	if !bytes.Equal(observedSubKey, expectedSubKey) {
		t.Errorf("Expected subset key %v, got %v", expectedSubKey, observedSubKey)
	}

	// Test error: invalid offset/length
	specObj3 := object.MakeEmptyObjectWithClassName(&className)
	params3 := []any{specObj3, keyObj, int64(23), int64(8), algoObj}
	result3 := secretKeySpecInit(params3)
	errBlk, ok := result3.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.InvalidKeyException {
		t.Errorf("Expected InvalidKeyException for invalid offset/length, got %v", result3)
	}

	// Test error: empty algorithm
	specObj4 := object.MakeEmptyObjectWithClassName(&className)
	emptyAlgoObj := object.StringObjectFromGoString("")
	params4 := []any{specObj4, keyObj, emptyAlgoObj}
	result4 := secretKeySpecInit(params4)
	errBlk, ok = result4.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for empty algorithm, got %v", result4)
	}

	// Test error: VALID offset/length
	specObj5 := object.MakeEmptyObjectWithClassName(&className)
	params5 := []any{specObj5, keyObj, int64(8), int64(16), algoObj} // up to the last allowed length of 4
	result5 := secretKeySpecInit(params5)
	if result5 != nil {
		t.Errorf("Unexpected InvalidKeyException, got %v", result5)
	}
}

func TestSecretKeySpecMethods(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/SecretKeySpec"
	key := []byte("secret")
	algo := "HmacSHA256"

	specObj := object.MakeEmptyObjectWithClassName(&className)
	specObj.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key}
	specObj.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algo}

	// Test getAlgorithm()
	resAlgo := secretKeySpecGetAlgorithm([]any{specObj})
	resAlgoObj, ok := resAlgo.(*object.Object)
	if !ok || object.GoStringFromStringObject(resAlgoObj) != algo {
		t.Errorf("Expected algorithm %s, got %v", algo, resAlgo)
	}

	// Test getFormat()
	resFormat := secretKeySpecGetFormat([]any{specObj})
	resFormatObj, ok := resFormat.(*object.Object)
	if !ok || object.GoStringFromStringObject(resFormatObj) != "RAW" {
		t.Errorf("Expected format RAW, got %v", resFormat)
	}

	// Test getEncoded()
	resEncoded := secretKeySpecGetEncoded([]any{specObj})
	resEncodedObj, ok := resEncoded.(*object.Object)
	if !ok {
		t.Fatalf("Expected object result for getEncoded, got %T", resEncoded)
	}
	jBytes := resEncodedObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	goBytes := object.GoByteArrayFromJavaByteArray(jBytes)
	if !bytes.Equal(goBytes, key) {
		t.Errorf("Expected encoded key %v, got %v", key, goBytes)
	}

	// Verify getEncoded returns a copy
	goBytes[0] ^= 0xFF
	if bytes.Equal(specObj.FieldTable["key"].Fvalue.([]byte), goBytes) {
		t.Error("getEncoded() should return a copy, but modifying the result affected the original")
	}
}

func TestSecretKeySpecEqualsAndHashCode(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/SecretKeySpec"
	key1 := []byte("key1")
	key2 := []byte("key2")
	algo1 := "AES"
	algo2 := "DES"

	spec1 := object.MakeEmptyObjectWithClassName(&className)
	spec1.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key1}
	spec1.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algo1}

	spec1Clone := object.MakeEmptyObjectWithClassName(&className)
	spec1Clone.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key1}
	spec1Clone.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algo1}

	spec2 := object.MakeEmptyObjectWithClassName(&className)
	spec2.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key2}
	spec2.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algo1}

	spec3 := object.MakeEmptyObjectWithClassName(&className)
	spec3.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key1}
	spec3.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algo2}

	// Test equals
	if secretKeySpecEquals([]any{spec1, spec1Clone}) != int64(1) {
		t.Error("spec1 should be equal to its clone")
	}
	if secretKeySpecEquals([]any{spec1, spec2}) != int64(0) {
		t.Error("spec1 should not be equal to spec2 (different key)")
	}
	if secretKeySpecEquals([]any{spec1, spec3}) != int64(0) {
		t.Error("spec1 should not be equal to spec3 (different algorithm)")
	}
	if secretKeySpecEquals([]any{spec1, spec1}) != int64(1) {
		t.Error("spec1 should be equal to itself")
	}
	if secretKeySpecEquals([]any{spec1, nil}) != int64(0) {
		t.Error("spec1 should not be equal to nil")
	}

	// Test hashCode
	h1 := secretKeySpecHashCode([]any{spec1}).(int64)
	h1c := secretKeySpecHashCode([]any{spec1Clone}).(int64)
	h2 := secretKeySpecHashCode([]any{spec2}).(int64)

	if h1 != h1c {
		t.Errorf("Equal objects should have same hash code: %d != %d", h1, h1c)
	}
	if h1 == h2 {
		t.Errorf("Different objects should likely have different hash codes (though not guaranteed): %d == %d", h1, h2)
	}
}
