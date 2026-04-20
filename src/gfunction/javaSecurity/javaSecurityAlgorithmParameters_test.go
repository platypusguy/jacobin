/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"bytes"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestAlgorithmParameters(t *testing.T) {
	globals.InitGlobals("test")

	// Pre-requisite: Need a provider with the service
	InitDefaultSecurityProvider()

	algo := "AES"
	algoObj := object.StringObjectFromGoString(algo)

	// 1. Test getInstance
	res := AlgparamsGetInstance([]any{nil, algoObj})
	algParams, ok := res.(*object.Object)
	if !ok || algParams == nil {
		t.Fatalf("AlgparamsGetInstance failed: %v", res)
	}

	if algParams.KlassName != object.StringPoolIndexFromGoString(types.ClassNameAlgorithmParameters) {
		t.Errorf("Expected class %s, got index %d", types.ClassNameAlgorithmParameters, algParams.KlassName)
	}

	// 2. Test getAlgorithm
	resAlgo := AlgparamsGetAlgorithm([]any{algParams})
	resAlgoObj, ok := resAlgo.(*object.Object)
	if !ok || object.GoStringFromStringObject(resAlgoObj) != algo {
		t.Errorf("Expected algorithm %s, got %v", algo, resAlgo)
	}

	// 3. Test init with byte[]
	iv := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	ivObj := makeByteArrayObject(iv)

	resInit := AlgparamsInit([]any{algParams, ivObj})
	if resInit != nil {
		t.Errorf("AlgparamsInit failed: %v", resInit)
	}

	if algParams.FieldTable["initialized"].Fvalue.(bool) != true {
		t.Error("algParams should be marked as initialized")
	}

	// 4. Test getEncoded
	resEncoded := AlgparamsGetEncoded([]any{algParams})
	resEncodedObj, ok := resEncoded.(*object.Object)
	if !ok {
		t.Fatalf("AlgparamsGetEncoded failed: %v", resEncoded)
	}
	encodedJBytes := resEncodedObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	encoded := object.GoByteArrayFromJavaByteArray(encodedJBytes)
	if !bytes.Equal(encoded, iv) {
		t.Errorf("Expected encoded %v, got %v", iv, encoded)
	}

	// 5. Test getParameterSpec (IvParameterSpec)
	resSpec := AlgparamsGetParameterSpec([]any{algParams})
	specObj, ok := resSpec.(*object.Object)
	if !ok || specObj == nil {
		t.Fatalf("AlgparamsGetParameterSpec failed: %v", resSpec)
	}

	// Verify it's an IvParameterSpec
	if specObj.FieldTable["iv"].Fvalue == nil {
		t.Error("specObj should have an iv field")
	}

	// 5.5 Test init with AlgorithmParameterSpec
	algParams2, _ := AlgparamsGetInstance([]any{nil, algoObj}).(*object.Object)
	resInit3 := AlgparamsInit([]any{algParams2, specObj})
	if resInit3 != nil {
		t.Errorf("AlgparamsInit with spec failed: %v", resInit3)
	}
	if algParams2.FieldTable["initialized"].Fvalue.(bool) != true {
		t.Error("algParams2 should be marked as initialized")
	}
	resSpec2 := AlgparamsGetParameterSpec([]any{algParams2})
	if resSpec2 != specObj {
		t.Error("AlgparamsGetParameterSpec should return the same spec object")
	}

	// 6. Test already initialized error
	resInit2 := AlgparamsInit([]any{algParams, ivObj})
	errBlk, ok := resInit2.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException for double init, got %v", resInit2)
	}
}

func TestAlgorithmParametersErrors(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	// Test getInstance with unknown algorithm
	badAlgoObj := object.StringObjectFromGoString("UNKNOWN")
	res := AlgparamsGetInstance([]any{nil, badAlgoObj})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NoSuchAlgorithmException {
		t.Errorf("Expected NoSuchAlgorithmException, got %v", res)
	}

	// Test getEncoded before init
	algParams := object.MakeEmptyObjectWithClassName(&types.ClassNameAlgorithmParameters)
	algParams.FieldTable["initialized"] = object.Field{Ftype: types.Bool, Fvalue: false}

	resEnc := AlgparamsGetEncoded([]any{algParams})
	errBlk, ok = resEnc.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException for getEncoded before init, got %v", resEnc)
	}
}

// Helper (copy of what's often in crypto tests)
func makeByteArrayObject(data []byte) *object.Object {
	jBytes := object.JavaByteArrayFromGoByteArray(data)
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}
