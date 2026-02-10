/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"math/big"
	"testing"
)

func TestLoad_Security_Interfaces_EC_Keys(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Security_Interfaces_EC_Keys()

	expectedSignatures := []string{
		"java/security/interfaces/ECKey.<clinit>()V",
		"java/security/interfaces/ECKey.<init>()V",
		"java/security/interfaces/ECPrivateKey.<clinit>()V",
		"java/security/interfaces/ECPrivateKey.<init>()V",
		"java/security/interfaces/ECPublicKey.<clinit>()V",
		"java/security/interfaces/ECPublicKey.<init>()V",
		"java/security/interfaces/ECPrivateKey.getParams()Ljava/security/spec/ECParameterSpec;",
		"java/security/interfaces/ECPrivateKey.getS()Ljava/math/BigInteger;",
		"java/security/interfaces/ECPublicKey.getParams()Ljava/security/spec/ECParameterSpec;",
		"java/security/interfaces/ECPublicKey.getW()Ljava/security/spec/ECPoint;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s not found", sig)
		}
	}

	if len(ghelpers.MethodSignatures) != len(expectedSignatures) {
		t.Errorf("Expected %d signatures, got %d", len(expectedSignatures), len(ghelpers.MethodSignatures))
	}
}

func TestECPrivateKeyGetParams(t *testing.T) {
	globals.InitGlobals("test")

	paramsObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECParameterSpec)
	privKeyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPrivateKey)
	privKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: paramsObj}

	res := ecPrivateKeyGetParams([]any{privKeyObj})
	if res != paramsObj {
		t.Errorf("Expected params object, got %v", res)
	}

	// Negative test: missing params
	delete(privKeyObj.FieldTable, "params")
	res = ecPrivateKeyGetParams([]any{privKeyObj})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk for missing params, got %T", res)
	}
}

func TestECPrivateKeyGetS(t *testing.T) {
	globals.InitGlobals("test")

	curve := elliptic.P256()
	priv, _ := ecdsa.GenerateKey(curve, rand.Reader)

	privKeyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPrivateKey)
	privKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

	res := ecPrivateKeyGetS([]any{privKeyObj})
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected object.Object, got %T", res)
	}

	sVal := obj.FieldTable["value"].Fvalue.(*big.Int)
	if sVal.Cmp(priv.D) != 0 {
		t.Errorf("Expected S value %v, got %v", priv.D, sVal)
	}

	// Negative test: invalid value type
	privKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: "not a key"}
	res = ecPrivateKeyGetS([]any{privKeyObj})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk for invalid value type, got %T", res)
	}
}

func TestECPublicKeyGetParams(t *testing.T) {
	globals.InitGlobals("test")

	paramsObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECParameterSpec)
	pubKeyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPublicKey)
	pubKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: paramsObj}

	res := ecPublicKeyGetParams([]any{pubKeyObj})
	if res != paramsObj {
		t.Errorf("Expected params object, got %v", res)
	}

	// Negative test: missing params
	delete(pubKeyObj.FieldTable, "params")
	res = ecPublicKeyGetParams([]any{pubKeyObj})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk for missing params, got %T", res)
	}
}

func TestECPublicKeyGetW(t *testing.T) {
	globals.InitGlobals("test")

	// Case 1: Manual "w" field
	pointObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPoint)
	pubKeyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPublicKey)
	pubKeyObj.FieldTable["w"] = object.Field{Ftype: types.ECPoint, Fvalue: pointObj}

	res := ecPublicKeyGetW([]any{pubKeyObj})
	if res != pointObj {
		t.Errorf("Expected manual point object, got %v", res)
	}

	// Case 2: Extraction from "value" (*ecdsa.PublicKey)
	delete(pubKeyObj.FieldTable, "w")
	curve := elliptic.P256()
	priv, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: &priv.PublicKey}

	res = ecPublicKeyGetW([]any{pubKeyObj})
	resObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected object.Object from extraction, got %T", res)
	}
	if resObj.KlassName != stringPool.GetStringIndex(&types.ClassNameECPoint) {
		t.Errorf("Expected ECPoint class, got %v", resObj.KlassName)
	}

	xVal := resObj.FieldTable["x"].Fvalue.(*object.Object).FieldTable["value"].Fvalue.(*big.Int)
	if xVal.Cmp(priv.PublicKey.X) != 0 {
		t.Errorf("Expected X value %v, got %v", priv.PublicKey.X, xVal)
	}

	// Negative test: missing both
	delete(pubKeyObj.FieldTable, "value")
	res = ecPublicKeyGetW([]any{pubKeyObj})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk for missing data, got %T", res)
	}
}
