/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
	"testing"
)

func TestLoad_EllipticCurve(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_EllipticCurve()

	expectedSignatures := []string{
		"java/security/spec/EllipticCurve.<init>(Ljava/security/spec/ECField;Ljava/math/BigInteger;Ljava/math/BigInteger;)V",
		"java/security/spec/EllipticCurve.getField()Ljava/security/spec/ECField;",
		"java/security/spec/EllipticCurve.getA()Ljava/math/BigInteger;",
		"java/security/spec/EllipticCurve.getB()Ljava/math/BigInteger;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s not registered", sig)
		}
	}
}

func TestEllipticCurve(t *testing.T) {
	globals.InitGlobals("test")

	// Create dependencies
	fieldObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPoint) // Using a dummy class name
	aObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(123))
	bObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(456))

	// Create EllipticCurve object
	ecObj := object.MakeEmptyObjectWithClassName(&types.ClassNameEllipticCurve)

	// Test Init
	params := []any{ecObj, fieldObj, aObj, bObj}
	result := ellipticCurveInit(params)

	if result != nil {
		if errBlk, ok := result.(*ghelpers.GErrBlk); ok {
			t.Fatalf("ellipticCurveInit failed: %s", errBlk.ErrMsg)
		} else {
			t.Fatalf("ellipticCurveInit returned unexpected result: %v", result)
		}
	}

	// Verify fields
	if ecObj.FieldTable["field"].Fvalue != fieldObj {
		t.Error("field not correctly set")
	}
	if ecObj.FieldTable["a"].Fvalue != aObj {
		t.Error("a not correctly set")
	}
	if ecObj.FieldTable["b"].Fvalue != bObj {
		t.Error("b not correctly set")
	}

	// Test Getters
	// getField
	resField := ellipticCurveGetField([]any{ecObj})
	if resField != fieldObj {
		t.Errorf("ellipticCurveGetField returned %v, expected %v", resField, fieldObj)
	}

	// getA
	resA := ellipticCurveGetA([]any{ecObj})
	if resA != aObj {
		t.Errorf("ellipticCurveGetA returned %v, expected %v", resA, aObj)
	}

	// getB
	resB := ellipticCurveGetB([]any{ecObj})
	if resB != bObj {
		t.Errorf("ellipticCurveGetB returned %v, expected %v", resB, bObj)
	}
}

func TestEllipticCurveInvalidParams(t *testing.T) {
	globals.InitGlobals("test")

	ecObj := object.MakeEmptyObjectWithClassName(&types.ClassNameEllipticCurve)
	fieldObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPoint)
	aObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(123))
	// bObj missing

	// Test Init with missing parameter
	params := []any{ecObj, fieldObj, aObj}
	result := ellipticCurveInit(params)
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("ellipticCurveInit should have failed with GErrBlk for missing parameter")
	}

	// Test Init with wrong type
	params = []any{ecObj, fieldObj, aObj, "not an object"}
	result = ellipticCurveInit(params)
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("ellipticCurveInit should have failed with GErrBlk for wrong parameter type")
	}

	// Test Getters with missing 'this'
	result = ellipticCurveGetField([]any{})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("ellipticCurveGetField should have failed with GErrBlk for missing 'this'")
	}

	// Test Getters with null/invalid fields
	ecObjEmpty := object.MakeEmptyObjectWithClassName(&types.ClassNameEllipticCurve)
	result = ellipticCurveGetA([]any{ecObjEmpty})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("ellipticCurveGetA should have failed with GErrBlk for missing field")
	}
}
