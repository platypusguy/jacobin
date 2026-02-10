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

func TestLoad_Security_Interfaces_DSA_Keys(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Security_Interfaces_DSA_Keys()

	expectedSignatures := []string{
		"java/security/interfaces/DSAKey.<clinit>()V",
		"java/security/interfaces/DSAKey.<init>()V",
		"java/security/interfaces/DSAKey.getParams()Ljava/security/interfaces/DSAParams;",
		"java/security/interfaces/DSAParams.<clinit>()V",
		"java/security/interfaces/DSAParams.<init>()V",
		"java/security/interfaces/DSAParams.getG()Ljava/math/BigInteger;",
		"java/security/interfaces/DSAParams.getP()Ljava/math/BigInteger;",
		"java/security/interfaces/DSAParams.getQ()Ljava/math/BigInteger;",
		"java/security/spec/DSAParameterSpec.<clinit>()V",
		"java/security/spec/DSAParameterSpec.<init>()V",
		"java/security/spec/DSAParameterSpec.getG()Ljava/math/BigInteger;",
		"java/security/spec/DSAParameterSpec.getP()Ljava/math/BigInteger;",
		"java/security/spec/DSAParameterSpec.getQ()Ljava/math/BigInteger;",
		"java/security/interfaces/DSAPrivateKey.<clinit>()V",
		"java/security/interfaces/DSAPrivateKey.<init>()V",
		"java/security/interfaces/DSAPrivateKey.getX()Ljava/math/BigInteger;",
		"java/security/interfaces/DSAPrivateKey.getParams()Ljava/security/interfaces/DSAParams;",
		"java/security/interfaces/DSAPublicKey.<clinit>()V",
		"java/security/interfaces/DSAPublicKey.<init>()V",
		"java/security/interfaces/DSAPublicKey.getY()Ljava/math/BigInteger;",
		"java/security/interfaces/DSAPublicKey.getParams()Ljava/security/interfaces/DSAParams;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s not registered", sig)
		}
	}
}

func TestDSAPrivateKeyGetX(t *testing.T) {
	globals.InitGlobals("test")
	xVal := big.NewInt(12345)

	// Positive test
	keyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDSAPrivateKey)
	xObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, xVal)
	keyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: xObj}

	result := dsaPrivateGetX([]any{keyObj})
	if result != xObj {
		t.Errorf("Expected same BigInteger object, got %p vs %p", xObj, result)
	}
	resObj := result.(*object.Object)

	if resObj.KlassName != object.StringPoolIndexFromGoString(types.ClassNameBigInteger) {
		t.Errorf("Expected BigInteger, got class index %d", resObj.KlassName)
	}

	resVal := resObj.FieldTable["value"].Fvalue.(*big.Int)
	if resVal.Cmp(xVal) != 0 {
		t.Errorf("Expected %v, got %v", xVal, resVal)
	}

	// Negative tests
	if err := dsaPrivateGetX([]any{}); !isGErrBlk(err) {
		t.Error("Expected error for empty params")
	}
	if err := dsaPrivateGetX([]any{"not an object"}); !isGErrBlk(err) {
		t.Error("Expected error for non-object param")
	}
	keyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: "not a big.Int"}
	if err := dsaPrivateGetX([]any{keyObj}); !isGErrBlk(err) {
		t.Error("Expected error for invalid value type")
	}
}

func TestDSAPublicKeyGetY(t *testing.T) {
	globals.InitGlobals("test")
	yVal := big.NewInt(67890)

	// Positive test
	keyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDSAPublicKey)
	yObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, yVal)
	keyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: yObj}

	result := dsaPublicKeyGetY([]any{keyObj})
	if result != yObj {
		t.Errorf("Expected same BigInteger object, got %p vs %p", yObj, result)
	}
	resObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	resVal := resObj.FieldTable["value"].Fvalue.(*big.Int)
	if resVal.Cmp(yVal) != 0 {
		t.Errorf("Expected %v, got %v", yVal, resVal)
	}

	// Negative tests
	if err := dsaPublicKeyGetY([]any{}); !isGErrBlk(err) {
		t.Error("Expected error for empty params")
	}
	keyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: nil}
	if err := dsaPublicKeyGetY([]any{keyObj}); !isGErrBlk(err) {
		t.Error("Expected error for missing value")
	}
}

func TestDSAParamsGetters(t *testing.T) {
	globals.InitGlobals("test")
	p := big.NewInt(11)
	q := big.NewInt(12)
	g := big.NewInt(13)

	paramsObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDSAParameterSpec)
	pObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, p)
	qObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, q)
	gObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, g)
	paramsObj.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: pObj}
	paramsObj.FieldTable["q"] = object.Field{Ftype: types.BigInteger, Fvalue: qObj}
	paramsObj.FieldTable["g"] = object.Field{Ftype: types.BigInteger, Fvalue: gObj}

	// Test getP
	resP := dsaParamsGetP([]any{paramsObj})
	if resP != pObj {
		t.Errorf("getP failed: expected same object %p, got %p", pObj, resP)
	}
	if val := resP.(*object.Object).FieldTable["value"].Fvalue.(*big.Int); val.Cmp(p) != 0 {
		t.Errorf("getP failed: expected %v, got %v", p, val)
	}

	// Test getQ
	resQ := dsaParamsGetQ([]any{paramsObj})
	if resQ != qObj {
		t.Errorf("getQ failed: expected same object %p, got %p", qObj, resQ)
	}
	if val := resQ.(*object.Object).FieldTable["value"].Fvalue.(*big.Int); val.Cmp(q) != 0 {
		t.Errorf("getQ failed: expected %v, got %v", q, val)
	}

	// Test getG
	resG := dsaParamsGetG([]any{paramsObj})
	if resG != gObj {
		t.Errorf("getG failed: expected same object %p, got %p", gObj, resG)
	}
	if val := resG.(*object.Object).FieldTable["value"].Fvalue.(*big.Int); val.Cmp(g) != 0 {
		t.Errorf("getG failed: expected %v, got %v", g, val)
	}
}

func TestDSAKeysGetParams(t *testing.T) {
	globals.InitGlobals("test")
	paramsObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDSAParameterSpec)

	keyObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDSAPublicKey)
	keyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: paramsObj}

	result := dsaKeysGetParams([]any{keyObj})
	if result != paramsObj {
		t.Error("getParams failed to return the expected params object")
	}

	// Negative test
	keyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: "not an object"}
	if err := dsaKeysGetParams([]any{keyObj}); !isGErrBlk(err) {
		t.Error("Expected error for invalid params field")
	}
}

func isGErrBlk(v any) bool {
	_, ok := v.(*ghelpers.GErrBlk)
	return ok
}
