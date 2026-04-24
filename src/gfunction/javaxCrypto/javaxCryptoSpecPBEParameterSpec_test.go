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

func TestPBEParameterSpec(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/spec/PBEParameterSpec"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	saltObj := makeByteArrayObject(salt)
	iterationCount := int64(1000)

	// Test Constructor and Getters
	specObj := object.MakeEmptyObjectWithClassName(&className)
	params := []any{specObj, saltObj, iterationCount}
	res := pbeParameterSpecInit(params)
	if res != nil {
		t.Fatalf("pbeParameterSpecInit failed: %v", res)
	}

	// Test getIterationCount()
	resIC := pbeParameterSpecGetIterationCount([]any{specObj})
	if resIC.(int64) != iterationCount {
		t.Errorf("Expected iterationCount %d, got %d", iterationCount, resIC)
	}

	// Test getSalt()
	resSalt := pbeParameterSpecGetSalt([]any{specObj})
	resSaltObj, ok := resSalt.(*object.Object)
	if !ok || resSaltObj == nil {
		t.Fatalf("Expected salt object, got %v", resSalt)
	}
	observedSalt := object.GoByteArrayFromJavaByteArray(resSaltObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(observedSalt, salt) {
		t.Errorf("Expected salt %v, got %v", salt, observedSalt)
	}

	// Verify getSalt() returns a copy
	observedSalt[0] ^= 0xFF
	resSalt2 := pbeParameterSpecGetSalt([]any{specObj})
	observedSalt2 := object.GoByteArrayFromJavaByteArray(resSalt2.(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte))
	if bytes.Equal(observedSalt2, observedSalt) {
		t.Error("getSalt() should return a copy, but modifying the result affected the original")
	}
}

func TestPBEParameterSpecErrors(t *testing.T) {
	globals.InitGlobals("test")
	className := "javax/crypto/spec/PBEParameterSpec"
	specObj := object.MakeEmptyObjectWithClassName(&className)

	// Test null salt
	resNullSalt := pbeParameterSpecInit([]any{specObj, object.Null, int64(1000)})
	errBlk, ok := resNullSalt.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null salt, got %v", resNullSalt)
	}

	// Test negative iteration count
	saltObj := makeByteArrayObject([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	resNegIC := pbeParameterSpecInit([]any{specObj, saltObj, int64(-1)})
	errBlk, ok = resNegIC.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for negative iterationCount, got %v", resNegIC)
	}

	// Test this=null
	resNullThis := pbeParameterSpecGetSalt([]any{object.Null})
	errBlk, ok = resNullThis.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for null 'this', got %v", resNullThis)
	}
}

func TestPBEParameterSpecWithSpec(t *testing.T) {
	globals.InitGlobals("test")

	className := "javax/crypto/spec/PBEParameterSpec"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	saltObj := makeByteArrayObject(salt)
	iterationCount := int64(1000)

	iv := []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ivSpecClass := "javax/crypto/spec/IvParameterSpec"
	ivSpec := object.MakeEmptyObjectWithClassName(&ivSpecClass)
	ivSpec.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: makeByteArrayObject(iv)}

	// 1. Test constructor with AlgorithmParameterSpec
	specObj := object.MakeEmptyObjectWithClassName(&className)
	res := pbeParameterSpecInitWithSpec([]any{specObj, saltObj, iterationCount, ivSpec})
	if res != nil {
		t.Fatalf("pbeParameterSpecInitWithSpec failed: %v", res)
	}

	// Verify iterationCount
	resIC := pbeParameterSpecGetIterationCount([]any{specObj})
	if resIC.(int64) != iterationCount {
		t.Errorf("Expected iterationCount %d, got %d", iterationCount, resIC)
	}

	// Verify salt
	resSalt := pbeParameterSpecGetSalt([]any{specObj})
	resSaltObj := resSalt.(*object.Object)
	observedSalt := object.GoByteArrayFromJavaByteArray(resSaltObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if !bytes.Equal(observedSalt, salt) {
		t.Errorf("Expected salt %v, got %v", salt, observedSalt)
	}

	// Verify paramSpec
	resParamSpec := pbeParameterSpecGetParameterSpec([]any{specObj})
	if resParamSpec != ivSpec {
		t.Errorf("Expected paramSpec %v, got %v", ivSpec, resParamSpec)
	}

	// 2. Test with null paramSpec
	specObj2 := object.MakeEmptyObjectWithClassName(&className)
	res = pbeParameterSpecInitWithSpec([]any{specObj2, saltObj, iterationCount, object.Null})
	if res != nil {
		t.Fatalf("pbeParameterSpecInitWithSpec failed with null spec: %v", res)
	}
	resParamSpec2 := pbeParameterSpecGetParameterSpec([]any{specObj2})
	if !object.IsNull(resParamSpec2) {
		t.Errorf("Expected null paramSpec, got %v", resParamSpec2)
	}
}
