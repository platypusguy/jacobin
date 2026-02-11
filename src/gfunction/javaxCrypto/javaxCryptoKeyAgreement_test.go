/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaSecurity"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestLoad_Crypto_KeyAgreement(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Crypto_KeyAgreement()

	methods := []string{
		"javax/crypto/KeyAgreement.getInstance(Ljava/lang/String;)Ljavax/crypto/KeyAgreement;",
		"javax/crypto/KeyAgreement.init(Ljava/security/Key;)V",
		"javax/crypto/KeyAgreement.doPhase(Ljava/security/Key;Z)Ljava/security/Key;",
		"javax/crypto/KeyAgreement.generateSecret()[B",
		"javax/crypto/KeyAgreement.getAlgorithm()Ljava/lang/String;",
	}

	for _, m := range methods {
		if _, ok := ghelpers.MethodSignatures[m]; !ok {
			t.Errorf("KeyAgreement method signature not registered: %s", m)
		}
	}
}

func TestKeyAgreement_ECDH(t *testing.T) {
	globals.InitGlobals("test")
	javaSecurity.InitDefaultSecurityProvider()

	// 1. Get Instance
	algoObj := object.StringObjectFromGoString("ECDH")
	result := keyagreementGetInstance([]any{algoObj})
	kaObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	// 2. Generate EC Keys
	curve := elliptic.P256()
	privA, _ := ecdsa.GenerateKey(curve, rand.Reader)
	privB, _ := ecdsa.GenerateKey(curve, rand.Reader)

	// Wrap Keys
	privateKeyObjA := javaSecurity.NewGoRuntimeService("ECPrivateKey", "EC", types.ClassNameECPrivateKey)
	privateKeyObjA.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privA}

	publicKeyObjB := javaSecurity.NewGoRuntimeService("ECPublicKey", "EC", types.ClassNameECPublicKey)
	publicKeyObjB.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: &privB.PublicKey}

	// 3. Init
	resInit := keyagreementInit([]any{kaObj, privateKeyObjA})
	if resInit != nil {
		t.Fatalf("Init failed: %v", resInit)
	}

	// 4. DoPhase
	resPhase := keyagreementDoPhase([]any{kaObj, publicKeyObjB, types.JavaBoolTrue})
	if resPhase != object.Null {
		t.Fatalf("DoPhase failed, expected null, got %v", resPhase)
	}

	// 5. Generate Secret
	resSecret := keyagreementGenerateSecret([]any{kaObj})
	secretObj, ok := resSecret.(*object.Object)
	if !ok {
		t.Fatalf("GenerateSecret failed: %v", resSecret)
	}

	secretBytes := object.GoByteArrayFromJavaByteArray(secretObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if len(secretBytes) == 0 {
		t.Error("Generated secret is empty")
	}

	// Verify against Go's own calculation
	x, _ := curve.ScalarMult(privB.PublicKey.X, privB.PublicKey.Y, privA.D.Bytes())
	expectedSecret := x.Bytes()

	if fmt.Sprintf("%x", secretBytes) != fmt.Sprintf("%x", expectedSecret) {
		t.Errorf("Secret mismatch.\nGot: %x\nExp: %x", secretBytes, expectedSecret)
	}
}

func TestKeyAgreement_X25519(t *testing.T) {
	globals.InitGlobals("test")
	javaSecurity.InitDefaultSecurityProvider()

	// 1. Get Instance
	algoObj := object.StringObjectFromGoString("X25519")
	result := keyagreementGetInstance([]any{algoObj})
	kaObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	// 2. Mock Keys (32 bytes for X25519)
	privA := make([]byte, 32)
	_, _ = rand.Read(privA)
	pubB := make([]byte, 32)
	_, _ = rand.Read(pubB)

	privateKeyObjA := javaSecurity.NewGoRuntimeService("XDH", "XDH", types.ClassNameEdECPrivateKey)
	privateKeyObjA.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privA}

	publicKeyObjB := javaSecurity.NewGoRuntimeService("XDH", "XDH", types.ClassNameEdECPublicKey)
	publicKeyObjB.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubB}

	// 3. Init
	keyagreementInit([]any{kaObj, privateKeyObjA})

	// 4. DoPhase
	keyagreementDoPhase([]any{kaObj, publicKeyObjB, types.JavaBoolTrue})

	// 5. Generate Secret
	resSecret := keyagreementGenerateSecret([]any{kaObj})
	secretObj := resSecret.(*object.Object)
	secretBytes := object.GoByteArrayFromJavaByteArray(secretObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if len(secretBytes) != 32 {
		t.Errorf("Expected 32 byte secret for X25519, got %d", len(secretBytes))
	}
}

func TestKeyAgreement_InvalidStates(t *testing.T) {
	globals.InitGlobals("test")
	javaSecurity.InitDefaultSecurityProvider()

	kaObj := keyagreementGetInstance([]any{object.StringObjectFromGoString("ECDH")}).(*object.Object)

	// 1. doPhase before init
	res := keyagreementDoPhase([]any{kaObj, object.Null, types.JavaBoolTrue})
	if gerr, ok := res.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.IllegalStateException {
		t.Errorf("Expected IllegalStateException, got %v", res)
	}

	// 2. generateSecret before doPhase
	kaObj.FieldTable["state"] = object.Field{Ftype: types.Int, Fvalue: int64(1)} // initialized
	res = keyagreementGenerateSecret([]any{kaObj})
	if gerr, ok := res.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.IllegalStateException {
		t.Errorf("Expected IllegalStateException, got %v", res)
	}
}

func TestKeyAgreement_UnsupportedAlgorithm(t *testing.T) {
	globals.InitGlobals("test")
	res := keyagreementGetInstance([]any{object.StringObjectFromGoString("UNKNOWN")})
	if gerr, ok := res.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NoSuchAlgorithmException {
		t.Errorf("Expected NoSuchAlgorithmException, got %v", res)
	}
}
