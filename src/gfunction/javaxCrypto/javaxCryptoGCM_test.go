/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"bytes"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestPerformCipherGCM(t *testing.T) {
	globals.InitGlobals("test")

	key := []byte("1234567812345678") // 16 bytes for AES-128
	iv := []byte("123456781234")      // 12 bytes for GCM
	plaintext := []byte("Hello GCM World!")

	config, ok := CipherConfigTable["AES/GCM/NoPadding"]
	if !ok {
		t.Fatal("AES/GCM/NoPadding not found in config table")
	}

	// Encrypt
	ciphertext, err := performCipher(config, 1, key, iv, plaintext)
	if err != nil {
		t.Fatalf("GCM Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("Ciphertext should not match plaintext")
	}

	// GCM ciphertext should be plaintext length + tag length (default 16 bytes for Go)
	expectedLen := len(plaintext) + 16
	if len(ciphertext) != expectedLen {
		t.Errorf("Expected ciphertext length %d, got %d", expectedLen, len(ciphertext))
	}

	// Decrypt
	decrypted, err := performCipher(config, 2, key, iv, ciphertext)
	if err != nil {
		t.Fatalf("GCM Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted data mismatch. Expected %s, got %s", string(plaintext), string(decrypted))
	}
}

func TestCipherGCMIntegration(t *testing.T) {
	globals.InitGlobals("test")

	// Setup Cipher
	trans := "AES/GCM/NoPadding"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// Setup Key
	keyBytes := []byte("1234567812345678")
	keyObj := makeByteArrayObject(keyBytes)

	// Setup GCMParameterSpec
	ivBytes := []byte("123456781234")
	tLen := int64(128)
	specClass := "javax/crypto/spec/GCMParameterSpec"
	specObj := object.MakeEmptyObjectWithClassName(&specClass)
	gcmParameterSpecInit([]any{specObj, tLen, makeByteArrayObject(ivBytes)})

	// Initialize Cipher
	cipherInit([]any{cipherObj, int64(1), keyObj, specObj}) // ENCRYPT_MODE

	// Encrypt
	plaintext := []byte("Secret Message")
	inputObj := makeByteArrayObject(plaintext)
	resEnc := cipherDoFinal([]any{cipherObj, inputObj})
	if _, ok := resEnc.(*object.Object); !ok {
		t.Fatalf("cipherDoFinal encrypt failed: %v", resEnc)
	}
	ciphertextObj := resEnc.(*object.Object)
	_ = object.GoByteArrayFromJavaByteArray(ciphertextObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	// Initialize for Decrypt
	cipherInit([]any{cipherObj, int64(2), keyObj, specObj}) // DECRYPT_MODE
	resDec := cipherDoFinal([]any{cipherObj, ciphertextObj})
	if _, ok := resDec.(*object.Object); !ok {
		t.Fatalf("cipherDoFinal decrypt failed: %v", resDec)
	}
	decrypted := object.GoByteArrayFromJavaByteArray(resDec.(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte))

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted mismatch. Expected %s, got %s", string(plaintext), string(decrypted))
	}
}
