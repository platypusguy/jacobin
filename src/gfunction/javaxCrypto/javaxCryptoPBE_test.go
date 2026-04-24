/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"bytes"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestPBEWithMD5AndDES(t *testing.T) {
	globals.InitGlobals("test")

	// 1. Setup Password and Salt
	password := "password"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	iterations := int64(1000)

	passwordChars := make([]int64, len(password))
	for i, c := range password {
		passwordChars[i] = int64(c)
	}
	passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)
	saltObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(salt))

	// 2. Derive Key using SecretKeyFactory
	skfAlgo := "PBEWithMD5AndDES"
	skfAlgoObj := object.StringObjectFromGoString(skfAlgo)

	// Mock Provider
	providerObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)
	providerObj.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(types.SecurityProviderName)}
	ghelpers.DefaultSecurityProvider = providerObj

	skf := secretKeyFactoryGetInstance([]any{skfAlgoObj}).(*object.Object)

	pbeKeySpecClassName := "javax/crypto/spec/PBEKeySpec"
	pbeKeySpec := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	genKey := secretKeyFactoryGenerateSecret([]any{skf, pbeKeySpec}).(*object.Object)

	// 3. Setup Cipher
	trans := "PBEWithMD5AndDES"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// PBEParameterSpec
	pbeParamSpecClass := "javax/crypto/spec/PBEParameterSpec"
	pbeParamSpec := object.MakeEmptyObjectWithClassName(&pbeParamSpecClass)
	pbeParamSpec.FieldTable["salt"] = object.Field{Ftype: types.ByteArray, Fvalue: salt}
	pbeParamSpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}

	// 4. Encrypt
	cipherInit([]any{cipherObj, int64(1), genKey, pbeParamSpec}) // ENCRYPT_MODE with PBEParameterSpec

	plaintext := []byte("This is a secret message")
	inputObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(plaintext))

	resEnc := cipherDoFinal([]any{cipherObj, inputObj})
	if errBlk, ok := resEnc.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Encryption failed: %s", errBlk.ErrMsg)
	}
	ciphertextObj := resEnc.(*object.Object)
	ciphertext := object.GoByteArrayFromJavaByteArray(ciphertextObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("Ciphertext should not match plaintext")
	}

	// 5. Decrypt
	cipherInit([]any{cipherObj, int64(2), genKey, pbeParamSpec}) // DECRYPT_MODE
	resDec := cipherDoFinal([]any{cipherObj, ciphertextObj})
	if errBlk, ok := resDec.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Decryption failed: %s", errBlk.ErrMsg)
	}
	decryptedObj := resDec.(*object.Object)
	decrypted := object.GoByteArrayFromJavaByteArray(decryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted message mismatch. Expected %s, got %s", string(plaintext), string(decrypted))
	}
}

func TestPBEWithSHA1AndRC2(t *testing.T) {
	globals.InitGlobals("test")

	// 1. Setup Password and Salt
	password := "password"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	iterations := int64(1000)

	passwordChars := make([]int64, len(password))
	for i, c := range password {
		passwordChars[i] = int64(c)
	}
	passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)
	saltObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(salt))

	// 2. Derive Key using SecretKeyFactory
	skfAlgo := "PBEWithSHA1AndRC2_128"
	skfAlgoObj := object.StringObjectFromGoString(skfAlgo)

	// Mock Provider
	providerObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)
	providerObj.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(types.SecurityProviderName)}
	ghelpers.DefaultSecurityProvider = providerObj

	skf := secretKeyFactoryGetInstance([]any{skfAlgoObj}).(*object.Object)

	pbeKeySpecClassName := "javax/crypto/spec/PBEKeySpec"
	pbeKeySpec := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	genKey := secretKeyFactoryGenerateSecret([]any{skf, pbeKeySpec}).(*object.Object)

	// 3. Setup Cipher
	trans := "PBEWithSHA1AndRC2_128"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// PBEParameterSpec
	pbeParamSpecClass := "javax/crypto/spec/PBEParameterSpec"
	pbeParamSpec := object.MakeEmptyObjectWithClassName(&pbeParamSpecClass)
	pbeParamSpec.FieldTable["salt"] = object.Field{Ftype: types.ByteArray, Fvalue: salt}
	pbeParamSpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}

	// 4. Encrypt
	cipherInit([]any{cipherObj, int64(1), genKey, pbeParamSpec}) // ENCRYPT_MODE

	plaintext := []byte("This is a secret message")
	inputObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(plaintext))

	resEnc := cipherDoFinal([]any{cipherObj, inputObj})
	if errBlk, ok := resEnc.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Encryption failed: %s", errBlk.ErrMsg)
	}
	ciphertextObj := resEnc.(*object.Object)
	ciphertext := object.GoByteArrayFromJavaByteArray(ciphertextObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("Ciphertext should not match plaintext")
	}

	// 5. Decrypt
	cipherInit([]any{cipherObj, int64(2), genKey, pbeParamSpec}) // DECRYPT_MODE
	resDec := cipherDoFinal([]any{cipherObj, ciphertextObj})
	if errBlk, ok := resDec.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Decryption failed: %s", errBlk.ErrMsg)
	}
	decryptedObj := resDec.(*object.Object)
	decrypted := object.GoByteArrayFromJavaByteArray(decryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted message mismatch. Expected %s, got %s", string(plaintext), string(decrypted))
	}
}

func TestPBEWithSHA1AndDESede(t *testing.T) {
	globals.InitGlobals("test")

	// 1. Setup Password and Salt
	password := "password"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	iterations := int64(1000)

	passwordChars := make([]int64, len(password))
	for i, c := range password {
		passwordChars[i] = int64(c)
	}
	passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)
	saltObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(salt))

	// 2. Derive Key using SecretKeyFactory
	skfAlgo := "PBEWithSHA1AndDESede"
	skfAlgoObj := object.StringObjectFromGoString(skfAlgo)

	// Mock Provider
	providerObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)
	providerObj.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(types.SecurityProviderName)}
	ghelpers.DefaultSecurityProvider = providerObj

	skf := secretKeyFactoryGetInstance([]any{skfAlgoObj}).(*object.Object)

	pbeKeySpecClassName := "javax/crypto/spec/PBEKeySpec"
	pbeKeySpec := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	genKey := secretKeyFactoryGenerateSecret([]any{skf, pbeKeySpec}).(*object.Object)

	// 3. Setup Cipher
	trans := "PBEWithSHA1AndDESede"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// PBEParameterSpec
	pbeParamSpecClass := "javax/crypto/spec/PBEParameterSpec"
	pbeParamSpec := object.MakeEmptyObjectWithClassName(&pbeParamSpecClass)
	pbeParamSpec.FieldTable["salt"] = object.Field{Ftype: types.ByteArray, Fvalue: salt}
	pbeParamSpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}

	// 4. Encrypt
	cipherInit([]any{cipherObj, int64(1), genKey, pbeParamSpec}) // ENCRYPT_MODE

	plaintext := []byte("TripleDES PBE test message")
	inputObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(plaintext))

	resEnc := cipherDoFinal([]any{cipherObj, inputObj})
	if errBlk, ok := resEnc.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Encryption failed: %s", errBlk.ErrMsg)
	}
	ciphertextObj := resEnc.(*object.Object)

	// 5. Decrypt
	cipherInit([]any{cipherObj, int64(2), genKey, pbeParamSpec}) // DECRYPT_MODE
	resDec := cipherDoFinal([]any{cipherObj, ciphertextObj})
	if errBlk, ok := resDec.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Decryption failed: %s", errBlk.ErrMsg)
	}
	decryptedObj := resDec.(*object.Object)
	decrypted := object.GoByteArrayFromJavaByteArray(decryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted message mismatch. Expected %s, got %s", string(plaintext), string(decrypted))
	}
}

func runPBEWithHmacTest(t *testing.T, algo string) {
	globals.InitGlobals("test")

	// 1. Setup Password and Salt
	password := "password"
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	iterations := int64(1000)

	passwordChars := make([]int64, len(password))
	for i, c := range password {
		passwordChars[i] = int64(c)
	}
	passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)
	saltObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(salt))

	// 2. Derive Key using SecretKeyFactory
	skfAlgoObj := object.StringObjectFromGoString(algo)

	// Mock Provider
	providerObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)
	providerObj.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(types.SecurityProviderName)}
	ghelpers.DefaultSecurityProvider = providerObj

	skf := secretKeyFactoryGetInstance([]any{skfAlgoObj}).(*object.Object)

	pbeKeySpecClassName := "javax/crypto/spec/PBEKeySpec"
	pbeKeySpec := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
	pbeKeySpec.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}
	pbeKeySpec.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: saltObj}
	pbeKeySpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}
	pbeKeySpec.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)} // Inferred

	genKey := secretKeyFactoryGenerateSecret([]any{skf, pbeKeySpec}).(*object.Object)

	// 3. Setup Cipher
	transObj := object.StringObjectFromGoString(algo)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// IV
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ivSpecClass := "javax/crypto/spec/IvParameterSpec"
	ivSpec := object.MakeEmptyObjectWithClassName(&ivSpecClass)
	ivSpec.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(iv))}

	// 4. Encrypt
	cipherInit([]any{cipherObj, int64(1), genKey, ivSpec}) // ENCRYPT_MODE

	plaintext := []byte("AES PBE test message for " + algo)
	inputObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(plaintext))

	resEnc := cipherDoFinal([]any{cipherObj, inputObj})
	if errBlk, ok := resEnc.(*ghelpers.GErrBlk); ok {
		t.Fatalf("[%s] Encryption failed: %s", algo, errBlk.ErrMsg)
	}
	ciphertextObj := resEnc.(*object.Object)

	// 5. Decrypt
	cipherInit([]any{cipherObj, int64(2), genKey, ivSpec}) // DECRYPT_MODE
	resDec := cipherDoFinal([]any{cipherObj, ciphertextObj})
	if errBlk, ok := resDec.(*ghelpers.GErrBlk); ok {
		t.Fatalf("[%s] Decryption failed: %s", algo, errBlk.ErrMsg)
	}
	decryptedObj := resDec.(*object.Object)
	decrypted := object.GoByteArrayFromJavaByteArray(decryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("[%s] Decrypted message mismatch. Expected %s, got %s", algo, string(plaintext), string(decrypted))
	}
}

func TestPBEWithHmac(t *testing.T) {
	algos := []string{
		"PBEWithHmacSHA1AndAES_128",
		"PBEWithHmacSHA224AndAES_128",
		"PBEWithHmacSHA256AndAES_128",
		"PBEWithHmacSHA384AndAES_128",
		"PBEWithHmacSHA512AndAES_128",
		"PBEWithHmacSHA1AndAES_256",
		"PBEWithHmacSHA224AndAES_256",
		"PBEWithHmacSHA256AndAES_256",
		"PBEWithHmacSHA384AndAES_256",
		"PBEWithHmacSHA512AndAES_256",
		"PBEWithHmacSHA512/224AndAES_128",
		"PBEWithHmacSHA512/224AndAES_256",
		"PBEWithHmacSHA512/256AndAES_128",
		"PBEWithHmacSHA512/256AndAES_256",
	}

	for _, algo := range algos {
		t.Run(algo, func(t *testing.T) {
			runPBEWithHmacTest(t, algo)
		})
	}
}
