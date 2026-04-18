/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestSecretKeyFactory(t *testing.T) {
	globals.InitGlobals("test")

	// Initialize DefaultSecurityProvider
	providerObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)
	nameObj := object.StringObjectFromGoString(types.SecurityProviderName)
	providerObj.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: nameObj}
	ghelpers.DefaultSecurityProvider = providerObj

	algo := "AES/ECB/PKCS5Padding"
	algoObj := object.StringObjectFromGoString(algo)

	// Test getInstance(String algorithm)
	res := secretKeyFactoryGetInstance([]any{algoObj})
	skf, ok := res.(*object.Object)
	if !ok {
		if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
			t.Fatalf("Expected SecretKeyFactory object, got error: %v", errBlk.ErrMsg)
		}
		t.Fatalf("Expected SecretKeyFactory object, got %T", res)
	}

	skfAlgoObj := skf.FieldTable["algorithm"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(skfAlgoObj) != algo {
		t.Errorf("Expected algorithm %s, got %s", algo, object.GoStringFromStringObject(skfAlgoObj))
	}

	// Test getInstance with provider name
	providerNameObj := object.StringObjectFromGoString(types.SecurityProviderName)
	res = secretKeyFactoryGetInstance([]any{algoObj, providerNameObj})
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("Expected SecretKeyFactory object with provider name, got %v", res)
	}

	// Test getInstance with an invalid provider name
	invalidProviderNameObj := object.StringObjectFromGoString("InvalidProvider")
	res = secretKeyFactoryGetInstance([]any{algoObj, invalidProviderNameObj})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.ProviderNotFoundException {
		t.Errorf("Expected ProviderNotFoundException for invalid provider, got %v", res)
	}

	// Test generateSecret with SecretKeySpec
	key := []byte("0123456789abcdef")
	specClassName := "javax/crypto/spec/SecretKeySpec"
	specObj := object.MakeEmptyObjectWithClassName(&specClassName)
	specObj.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key}
	specObj.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algoObj}

	res = secretKeyFactoryGenerateSecret([]any{skf, specObj})
	genKey, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected SecretKey object, got %v", res)
	}
	if genKey != specObj {
		t.Errorf("Expected generated key to be the same as specObj for SecretKeySpec")
	}

	// Test generateSecret with algorithm mismatch
	otherAlgoObj := object.StringObjectFromGoString("DES")
	specObjMismatch := object.MakeEmptyObjectWithClassName(&specClassName)
	specObjMismatch.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: key}
	specObjMismatch.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: otherAlgoObj}

	res = secretKeyFactoryGenerateSecret([]any{skf, specObjMismatch})
	errBlk, ok = res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.InvalidKeyException {
		t.Errorf("Expected InvalidKeyException for algorithm mismatch, got %v", res)
	}

	// Test PBKDF2WithHmacSHA1 getInstance
	pbkdf2AlgoObj := object.StringObjectFromGoString("PBKDF2WithHmacSHA1")
	res = secretKeyFactoryGetInstance([]any{pbkdf2AlgoObj})
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("Expected SecretKeyFactory object for PBKDF2WithHmacSHA1, got %v", res)
	}

	// Test PBKDF2WithHmacSHA1 generateSecret
	password := "password"
	salt := []byte("salt")
	iterations := int64(1000)
	keyLength := int64(128)

	passwordChars := make([]int64, len(password))
	for i, c := range password {
		passwordChars[i] = int64(c)
	}
	passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)

	saltJavaBytes := object.JavaByteArrayFromGoByteArray(salt)
	saltObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, saltJavaBytes)

	pbeKeySpecClassName := "javax/crypto/spec/PBEKeySpec"
	pbeKeySpec := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: keyLength}

	// Re-get SKF for PBKDF2WithHmacSHA1
	res = secretKeyFactoryGetInstance([]any{pbkdf2AlgoObj})
	skfPbkdf2 := res.(*object.Object)

	res = secretKeyFactoryGenerateSecret([]any{skfPbkdf2, pbeKeySpec})
	genKey, ok = res.(*object.Object)
	if !ok {
		if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
			t.Fatalf("Expected generated key, got error: %v", errBlk.ErrMsg)
		}
		t.Fatalf("Expected generated key object, got %T", res)
	}

	derivedKey := genKey.FieldTable["value"].Fvalue.([]byte)
	if len(derivedKey) != int(keyLength/8) {
		t.Errorf("Expected key length %d, got %d", keyLength/8, len(derivedKey))
	}

	// Test PBEWithHmacSHA1AndAES_256 getInstance
	pbeAesAlgo := "PBEWithHmacSHA1AndAES_256"
	pbeAesAlgoObj := object.StringObjectFromGoString(pbeAesAlgo)
	res = secretKeyFactoryGetInstance([]any{pbeAesAlgoObj})
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("Expected SecretKeyFactory object for %s, got %v", pbeAesAlgo, res)
	}

	// Test PBEWithHmacSHA1AndAES_256 generateSecret
	skfPbeAes := res.(*object.Object)
	pbeKeySpec2 := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec2.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec2.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec2.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec2.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)} // Should be inferred

	res = secretKeyFactoryGenerateSecret([]any{skfPbeAes, pbeKeySpec2})
	genKey, ok = res.(*object.Object)
	if !ok {
		if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
			t.Fatalf("Expected generated key for %s, got error: %v", pbeAesAlgo, errBlk.ErrMsg)
		}
		t.Fatalf("Expected generated key object for %s, got %T", pbeAesAlgo, res)
	}

	derivedKey = genKey.FieldTable["value"].Fvalue.([]byte)
	if len(derivedKey) != 256/8 {
		t.Errorf("Expected key length %d, got %d for %s", 256/8, len(derivedKey), pbeAesAlgo)
	}

	// Test PBEWithMD5AndDES getInstance
	pbeLegacyAlgo := "PBEWithMD5AndDES"
	pbeLegacyAlgoObj := object.StringObjectFromGoString(pbeLegacyAlgo)
	res = secretKeyFactoryGetInstance([]any{pbeLegacyAlgoObj})
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("Expected SecretKeyFactory object for %s, got %v", pbeLegacyAlgo, res)
	}

	// Test PBEWithMD5AndDES generateSecret
	skfPbeLegacy := res.(*object.Object)
	pbeKeySpec3 := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec3.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec3.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec3.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec3.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)} // Should be inferred

	res = secretKeyFactoryGenerateSecret([]any{skfPbeLegacy, pbeKeySpec3})
	genKey, ok = res.(*object.Object)
	if !ok {
		if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
			t.Fatalf("Expected generated key for %s, got error: %v", pbeLegacyAlgo, errBlk.ErrMsg)
		}
		t.Fatalf("Expected generated key object for %s, got %T", pbeLegacyAlgo, res)
	}

	derivedKey = genKey.FieldTable["value"].Fvalue.([]byte)
	if len(derivedKey) != 64/8 {
		t.Errorf("Expected key length %d, got %d for %s", 64/8, len(derivedKey), pbeLegacyAlgo)
	}
}
