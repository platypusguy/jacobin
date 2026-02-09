/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"math/big"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestLoad_ECFieldAndPoint(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_ECFieldAndPoint()

	expectedSignatures := []string{
		"java/security/spec/ECField.<init>(I)V",
		"java/security/spec/ECField.getFieldSize()I",
		"java/security/spec/ECFieldFp.<init>(Ljava/math/BigInteger;)V",
		"java/security/spec/ECFieldFp.getP()Ljava/math/BigInteger;",
		"java/security/spec/ECPoint.<init>(Ljava/math/BigInteger;Ljava/math/BigInteger;)V",
		"java/security/spec/ECPoint.getAffineX()Ljava/math/BigInteger;",
		"java/security/spec/ECPoint.getAffineY()Ljava/math/BigInteger;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected method signature %s not found", sig)
		}
	}
}

func TestECField(t *testing.T) {
	globals.InitGlobals("test")

	// Test ecFieldInit
	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECParameterSpec) // Using a dummy class name
	fieldSize := int64(256)
	params := []any{thisObj, fieldSize}

	result := ecFieldInit(params)
	if result != nil {
		t.Fatalf("ecFieldInit returned %v, expected nil", result)
	}

	if val, ok := thisObj.FieldTable["fieldSize"]; !ok || val.Fvalue != fieldSize {
		t.Errorf("fieldSize not set correctly in FieldTable: got %v, expected %d", val.Fvalue, fieldSize)
	}

	// Test ecFieldGetFieldSize
	paramsGet := []any{thisObj}
	resultGet := ecFieldGetFieldSize(paramsGet)
	if resultGet != fieldSize {
		t.Errorf("ecFieldGetFieldSize returned %v, expected %d", resultGet, fieldSize)
	}

	// Test negative cases for ecFieldInit
	badParams := []any{thisObj}
	resultBad := ecFieldInit(badParams)
	if _, ok := resultBad.(*ghelpers.GErrBlk); !ok {
		t.Errorf("ecFieldInit with missing params should return GErrBlk, got %T", resultBad)
	}

	badParams2 := []any{thisObj, "not an int"}
	resultBad2 := ecFieldInit(badParams2)
	if _, ok := resultBad2.(*ghelpers.GErrBlk); !ok {
		t.Errorf("ecFieldInit with wrong type should return GErrBlk, got %T", resultBad2)
	}

	// Test negative cases for ecFieldGetFieldSize
	resultBadGet := ecFieldGetFieldSize([]any{})
	if _, ok := resultBadGet.(*ghelpers.GErrBlk); !ok {
		t.Errorf("ecFieldGetFieldSize with missing params should return GErrBlk, got %T", resultBadGet)
	}
}

func TestECFieldFp(t *testing.T) {
	globals.InitGlobals("test")

	pVal := big.NewInt(123456789)
	pObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, pVal)
	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECParameterSpec)

	// Test ecFieldFpInit
	params := []any{thisObj, pObj}
	result := ecFieldFpInit(params)
	if result != nil {
		t.Fatalf("ecFieldFpInit returned %v, expected nil", result)
	}

	if val, ok := thisObj.FieldTable["p"]; !ok || val.Fvalue != pObj {
		t.Errorf("p not set correctly in FieldTable")
	}

	// Test ecFieldFpGetP
	paramsGet := []any{thisObj}
	resultGet := ecFieldFpGetP(paramsGet)
	if resultGet != pObj {
		t.Errorf("ecFieldFpGetP returned %v, expected %v", resultGet, pObj)
	}

	// Test negative cases
	resultBad := ecFieldFpInit([]any{thisObj})
	if _, ok := resultBad.(*ghelpers.GErrBlk); !ok {
		t.Errorf("ecFieldFpInit with missing params should return GErrBlk")
	}

	resultBadGet := ecFieldFpGetP([]any{"not an object"})
	if _, ok := resultBadGet.(*ghelpers.GErrBlk); !ok {
		t.Errorf("ecFieldFpGetP with wrong type should return GErrBlk")
	}
}

func TestECPoint(t *testing.T) {
	globals.InitGlobals("test")

	xVal := big.NewInt(100)
	yVal := big.NewInt(200)
	xObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, xVal)
	yObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, yVal)
	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPoint)

	// Test ecPointInit
	params := []any{thisObj, xObj, yObj}
	result := ecPointInit(params)
	if result != nil {
		t.Fatalf("ecPointInit returned %v, expected nil", result)
	}

	if val, ok := thisObj.FieldTable["x"]; !ok || val.Fvalue != xObj {
		t.Errorf("x not set correctly")
	}
	if val, ok := thisObj.FieldTable["y"]; !ok || val.Fvalue != yObj {
		t.Errorf("y not set correctly")
	}

	// Test getters
	resX := ecPointGetAffineX([]any{thisObj})
	if resX != xObj {
		t.Errorf("ecPointGetAffineX failed")
	}

	resY := ecPointGetAffineY([]any{thisObj})
	if resY != yObj {
		t.Errorf("ecPointGetAffineY failed")
	}

	// Test negative cases
	resultBad := ecPointInit([]any{thisObj, xObj})
	if _, ok := resultBad.(*ghelpers.GErrBlk); !ok {
		t.Errorf("ecPointInit with missing params should return GErrBlk")
	}
}
