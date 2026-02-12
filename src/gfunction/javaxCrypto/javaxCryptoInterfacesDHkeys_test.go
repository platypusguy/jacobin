/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"math/big"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestLoad_Security_Interfaces_DH_Keys(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Crypto_Interfaces_DH_Keys()

	expectedSignatures := []string{
		"java/security/interfaces/DHKey.<clinit>()V",
		"java/security/interfaces/DHKey.<init>()V",
		"java/security/interfaces/DHKey.getParams()Ljavax/crypto/spec/DHParameterSpec;",
		"java/security/interfaces/DHPrivateKey.<clinit>()V",
		"java/security/interfaces/DHPrivateKey.<init>()V",
		"java/security/interfaces/DHPrivateKey.getX()()Ljava/math/BigInteger;",
		"java/security/interfaces/DHPublicKey.<clinit>()V",
		"java/security/interfaces/DHPublicKey.<init>()V",
		"java/security/interfaces/DHPublicKey.getY()Ljava/math/BigInteger;",
		"java/security/interfaces/DHPublicKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s not registered", sig)
		}
	}
}

func TestDHPrivateKeyGetX(t *testing.T) {
	globals.InitGlobals("test")

	// Valid case
	x := big.NewInt(123456789)
	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDHPrivateKey)
	thisObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: x}

	res := dhPrivateGetX([]any{thisObj})
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", res)
	}

	val, ok := obj.FieldTable["value"].Fvalue.(*big.Int)
	if !ok || val.Cmp(x) != 0 {
		t.Errorf("Expected BigInteger with value %v, got %v", x, val)
	}

	// Invalid case: missing this
	res = dhPrivateGetX([]any{})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected GErrBlk for missing params")
	}

	// Invalid case: wrong type
	res = dhPrivateGetX([]any{"not an object"})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected GErrBlk for wrong type")
	}

	// Invalid case: missing value field
	thisObjEmpty := object.MakeEmptyObjectWithClassName(&types.ClassNameDHPrivateKey)
	res = dhPrivateGetX([]any{thisObjEmpty})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected GErrBlk for missing value field")
	}
}

func TestDHPublicKeyGetY(t *testing.T) {
	globals.InitGlobals("test")

	// Valid case
	y := big.NewInt(987654321)
	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDHPublicKey)
	thisObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: y}

	res := dhPublicKeyGetY([]any{thisObj})
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", res)
	}

	val, ok := obj.FieldTable["value"].Fvalue.(*big.Int)
	if !ok || val.Cmp(y) != 0 {
		t.Errorf("Expected BigInteger with value %v, got %v", y, val)
	}

	// Invalid case: missing this
	res = dhPublicKeyGetY([]any{})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected GErrBlk for missing params")
	}

	// Invalid case: wrong type
	res = dhPublicKeyGetY([]any{"not an object"})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected GErrBlk for wrong type")
	}

	// Invalid case: missing value field
	thisObjEmpty := object.MakeEmptyObjectWithClassName(&types.ClassNameDHPublicKey)
	res = dhPublicKeyGetY([]any{thisObjEmpty})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected GErrBlk for missing value field")
	}
}
