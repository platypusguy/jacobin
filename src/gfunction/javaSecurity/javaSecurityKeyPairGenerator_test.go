/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestLoad_KeyPairGenerator_Detail(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_KeyPairGenerator()

	expectedSignatures := []string{
		"java/security/KeyPairGenerator.<init>(Ljava/lang/String;)V",
		"java/security/KeyPairGeneratorSpi.<init>()V",
		"java/security/KeyPairGeneratorSpi.generateKeyPair()Ljava/security/KeyPair;",
		"java/security/KeyPairGeneratorSpi.initialize(I)V",
		"java/security/KeyPairGeneratorSpi.initialize(ILjava/security/SecureRandom;)V",
		"java/security/KeyPairGenerator.getInstance(Ljava/lang/String;)Ljava/security/KeyPairGenerator;",
		"java/security/KeyPairGenerator.initialize(I)V",
		"java/security/KeyPairGenerator.initialize(ILjava/security/SecureRandom;)V",
		"java/security/KeyPairGenerator.initialize(Ljava/security/spec/AlgorithmParameterSpec;)V",
		"java/security/KeyPairGenerator.generateKeyPair()Ljava/security/KeyPair;",
		"java/security/KeyPairGenerator.genKeyPair()Ljava/security/KeyPair;",
		"java/security/KeyPairGenerator.getAlgorithm()Ljava/lang/String;",
		"java/security/KeyPairGenerator.getProvider()Ljava/security/Provider;",
		"java/security/KeyPairGenerator.getKeySize()I",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected method signature %s not registered", sig)
		}
	}
}

func TestKeyPairGeneratorGetInstance_Detail(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.DefaultSecurityProvider = InitDefaultSecurityProvider()

	// Test successful getInstance
	algoName := "RSA"
	params := []any{object.StringObjectFromGoString(algoName)}
	result := keypairgeneratorGetInstance(params)

	kpgObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	// Verify algorithm
	algObj := kpgObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(algObj) != algoName {
		t.Errorf("Expected algorithm %s, got %s", algoName, object.GoStringFromStringObject(algObj))
	}

	// Verify provider
	provider := kpgObj.FieldTable["provider"].Fvalue.(*object.Object)
	if provider == nil {
		t.Error("Provider not correctly set")
	}

	// Verify service
	svc := kpgObj.FieldTable["service"].Fvalue.(*object.Object)
	if svc == nil {
		t.Error("Service not correctly set")
	}

	// Test unsupported algorithm
	params = []any{object.StringObjectFromGoString("INVALID_ALGO")}
	result = keypairgeneratorGetInstance(params)
	if errBlk, ok := result.(*ghelpers.GErrBlk); ok {
		if errBlk.ExceptionType != excNames.IllegalArgumentException {
			t.Errorf("Expected IllegalArgumentException, got %v", errBlk.ExceptionType)
		}
	} else {
		t.Errorf("Expected GErrBlk for invalid algorithm, got %T", result)
	}
}

func TestKeyPairGeneratorInitialize(t *testing.T) {
	globals.InitGlobals("test")

	kpgObj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)
	keySize := int64(2048)

	params := []any{kpgObj, keySize}
	result := keypairgeneratorInitialize(params)

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	if val, ok := kpgObj.FieldTable["keySize"].Fvalue.(int64); !ok || val != keySize {
		t.Errorf("Expected keySize %d, got %v", keySize, kpgObj.FieldTable["keySize"].Fvalue)
	}
}

func TestKeyPairGeneratorInitializeWithRandom(t *testing.T) {
	globals.InitGlobals("test")

	kpgObj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)
	keySize := int64(1024)
	randomObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecureRandom)

	params := []any{kpgObj, keySize, randomObj}
	result := keypairgeneratorInitializeWithRandom(params)

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	if val, ok := kpgObj.FieldTable["keySize"].Fvalue.(int64); !ok || val != keySize {
		t.Errorf("Expected keySize %d, got %v", keySize, kpgObj.FieldTable["keySize"].Fvalue)
	}

	if val, ok := kpgObj.FieldTable["random"].Fvalue.(*object.Object); !ok || val != randomObj {
		t.Errorf("Expected random object, got %v", kpgObj.FieldTable["random"].Fvalue)
	}
}

func TestKeyPairGeneratorInitializeParmSpec(t *testing.T) {
	globals.InitGlobals("test")

	kpgObj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)
	algoName := "EC"
	kpgObj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algoName),
	}

	specObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECGenParameterSpec)
	curveName := "secp256r1"
	specObj.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(curveName),
	}

	params := []any{kpgObj, specObj}
	result := keypairgeneratorInitializeWithParmSpec(params)

	if result != nil {
		if err, ok := result.(*ghelpers.GErrBlk); ok {
			t.Fatalf("keypairgeneratorInitializeWithParmSpec failed: %s", err.ErrMsg)
		}
		t.Errorf("Expected nil result, got %v", result)
	}

	if val, ok := kpgObj.FieldTable["paramSpec"].Fvalue.(*object.Object); !ok || val != specObj {
		t.Errorf("Expected paramSpec object, got %v", kpgObj.FieldTable["paramSpec"].Fvalue)
	}

	if val, ok := kpgObj.FieldTable["keySize"].Fvalue.(int64); !ok || val != 256 {
		t.Errorf("Expected keySize 256 for secp256r1, got %v", kpgObj.FieldTable["keySize"].Fvalue)
	}
}

func TestKeyPairGeneratorGetters(t *testing.T) {
	globals.InitGlobals("test")

	kpgObj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)
	algoName := "EC"
	kpgObj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algoName),
	}
	providerObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)
	kpgObj.FieldTable["provider"] = object.Field{
		Ftype:  types.ClassNameSecurityProvider,
		Fvalue: providerObj,
	}
	keySize := int64(256)
	kpgObj.FieldTable["keySize"] = object.Field{Ftype: types.Int, Fvalue: keySize}

	// Test getAlgorithm
	result := keypairgeneratorGetAlgorithm([]any{kpgObj})
	algObj, ok := result.(*object.Object)
	if !ok || object.GoStringFromStringObject(algObj) != algoName {
		t.Errorf("Expected algorithm %s, got %v", algoName, result)
	}

	// Test getProvider
	result = keypairgeneratorGetProvider([]any{kpgObj})
	if result != providerObj {
		t.Errorf("Expected provider object, got %v", result)
	}

	// Test getKeySize
	result = keypairgeneratorGetKeySize([]any{kpgObj})
	if val, ok := result.(int64); !ok || val != keySize {
		t.Errorf("Expected keySize %d, got %v", keySize, result)
	}
}

func TestKeyPairGeneratorInvalidParams(t *testing.T) {
	globals.InitGlobals("test")

	// Test keypairgeneratorInitialize with wrong param count
	result := keypairgeneratorInitialize([]any{})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected error for empty params in keypairgeneratorInitialize")
	}

	// Test keypairgeneratorInitialize with wrong types
	result = keypairgeneratorInitialize([]any{"not an object", int64(1024)})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected error for invalid first param in keypairgeneratorInitialize")
	}

	result = keypairgeneratorInitialize([]any{object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator), "not an int"})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected error for invalid second param in keypairgeneratorInitialize")
	}

	// Test keypairgeneratorGetInstance with null algorithm
	result = keypairgeneratorGetInstance([]any{"not an object"})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Error("Expected error for non-object algorithm in keypairgeneratorGetInstance")
	}
}
