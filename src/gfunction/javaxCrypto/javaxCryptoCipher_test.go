/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"bytes"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestCipherClinit(t *testing.T) {
	globals.InitGlobals("test")
	cipherClinit(nil)

	tests := []struct {
		name  string
		value int64
	}{
		{"DECRYPT_MODE", 2},
		{"ENCRYPT_MODE", 1},
		{"PRIVATE_KEY", 2},
		{"PUBLIC_KEY", 1},
		{"SECRET_KEY", 3},
		{"UNWRAP_MODE", 4},
		{"WRAP_MODE", 3},
	}

	for _, tt := range tests {
		fqn := types.ClassNameCipher + "." + tt.name
		s, ok := statics.QueryStatic(types.ClassNameCipher, tt.name)
		if !ok {
			t.Errorf("Static %s not found", fqn)
			continue
		}
		if s.Value.(int64) != tt.value {
			t.Errorf("Static %s: expected %d, got %v", fqn, tt.value, s.Value)
		}
	}
}

func TestCipherGetInstance(t *testing.T) {
	globals.InitGlobals("test")

	// Valid transformation
	trans := "AES/CBC/PKCS5Padding"
	transObj := object.StringObjectFromGoString(trans)
	res := cipherGetInstance([]any{transObj})
	cipherObj, ok := res.(*object.Object)
	if !ok || cipherObj == nil {
		t.Fatalf("Expected cipher object, got %v", res)
	}

	observedTransObj := cipherObj.FieldTable["transformation"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(observedTransObj) != trans {
		t.Errorf("Expected transformation %s, got %s", trans, object.GoStringFromStringObject(observedTransObj))
	}

	// Invalid transformation
	invalidTransObj := object.StringObjectFromGoString("INVALID/ALGO")
	resErr := cipherGetInstance([]any{invalidTransObj})
	errBlk, ok := resErr.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NoSuchAlgorithmException {
		t.Errorf("Expected NoSuchAlgorithmException for invalid transformation, got %v", resErr)
	}

	// Null transformation
	resNull := cipherGetInstance([]any{nil})
	errBlk, ok = resNull.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for null transformation, got %v", resNull)
	}
}

func TestCipherInitAndUpdate(t *testing.T) {
	globals.InitGlobals("test")

	// Setup Cipher
	trans := "AES/CBC/PKCS5Padding"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// Setup Key
	keyBytes := []byte("1234567812345678") // 16 bytes for AES
	keyObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(keyBytes))

	// Setup IV
	ivBytes := []byte("iviviviviviviviv") // 16 bytes
	ivSpecClass := "javax/crypto/spec/IvParameterSpec"
	ivSpecObj := object.MakeEmptyObjectWithClassName(&ivSpecClass)
	ivSpecObj.FieldTable["iv"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(ivBytes)),
	}

	// cipherInit(opmode, key, spec)
	opmode := int64(1) // ENCRYPT_MODE
	cipherInit([]any{cipherObj, opmode, keyObj, ivSpecObj})

	if cipherObj.FieldTable["opmode"].Fvalue.(int64) != opmode {
		t.Errorf("Expected opmode %d, got %v", opmode, cipherObj.FieldTable["opmode"].Fvalue)
	}

	observedIvObj := cipherObj.FieldTable["iv"].Fvalue.(*object.Object)
	observedIvJBytes := observedIvObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !bytes.Equal(object.GoByteArrayFromJavaByteArray(observedIvJBytes), ivBytes) {
		t.Errorf("Expected IV %v, got %v", ivBytes, object.GoByteArrayFromJavaByteArray(observedIvJBytes))
	}

	// cipherUpdate
	input := []byte("hello world")
	inputObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(input))
	cipherUpdate([]any{cipherObj, inputObj})

	buffered := cipherObj.FieldTable["buffer"].Fvalue.([]byte)
	if !bytes.Equal(buffered, input) {
		t.Errorf("Expected buffer %v, got %v", input, buffered)
	}

	// cipherDoFinal
	res := cipherDoFinal([]any{cipherObj})
	if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("cipherDoFinal (encrypt) failed: %v", errBlk.ErrMsg)
	}
	encryptedObj := res.(*object.Object)
	encryptedJBytes := encryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	encrypted := object.GoByteArrayFromJavaByteArray(encryptedJBytes)

	if bytes.Equal(encrypted, input) {
		t.Errorf("Result should be encrypted, but matches input")
	}

	// Now decrypt
	cipherInit([]any{cipherObj, int64(2), keyObj, ivSpecObj}) // DECRYPT_MODE
	resDec := cipherDoFinal([]any{cipherObj, encryptedObj})
	if errBlk, ok := resDec.(*ghelpers.GErrBlk); ok {
		t.Fatalf("cipherDoFinal (decrypt) failed: %v", errBlk.ErrMsg)
	}
	decryptedObj := resDec.(*object.Object)
	decryptedJBytes := decryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	decrypted := object.GoByteArrayFromJavaByteArray(decryptedJBytes)

	if !bytes.Equal(decrypted, input) {
		t.Errorf("Expected decrypted result %v, got %v", input, decrypted)
	}

	// Buffer should be empty after doFinal
	bufferedPost := cipherObj.FieldTable["buffer"].Fvalue.([]byte)
	if len(bufferedPost) != 0 {
		t.Errorf("Expected empty buffer after doFinal, got len %d", len(bufferedPost))
	}
}

func TestCipherGetIV(t *testing.T) {
	globals.InitGlobals("test")

	// Setup Cipher
	trans := "AES/CBC/PKCS5Padding"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// Setup IV
	ivBytes := []byte("iviviviviviviviv") // 16 bytes
	ivSpecClass := "javax/crypto/spec/IvParameterSpec"
	ivSpecObj := object.MakeEmptyObjectWithClassName(&ivSpecClass)
	ivSpecObj.FieldTable["iv"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(ivBytes)),
	}

	// cipherInit(opmode, key, spec)
	keyBytes := []byte("1234567812345678")
	keyObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(keyBytes))
	cipherInit([]any{cipherObj, int64(1), keyObj, ivSpecObj})

	// cipherGetIV
	resIV := cipherGetIV([]any{cipherObj})
	resIVObj, ok := resIV.(*object.Object)
	if !ok || resIVObj == nil {
		t.Fatalf("Expected IV object, got %v", resIV)
	}
	resIVJBytes := resIVObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !bytes.Equal(object.GoByteArrayFromJavaByteArray(resIVJBytes), ivBytes) {
		t.Errorf("Expected IV %v, got %v", ivBytes, object.GoByteArrayFromJavaByteArray(resIVJBytes))
	}

	// No IV case
	cipherObj2 := cipherGetInstance([]any{transObj}).(*object.Object)
	// Use DECRYPT_MODE (2) to avoid auto-generating IV
	cipherInit([]any{cipherObj2, int64(2), keyObj})
	resIVNull := cipherGetIV([]any{cipherObj2})
	if !object.IsNull(resIVNull) {
		t.Errorf("Expected null IV for DECRYPT_MODE without spec, got %v", resIVNull)
	}
}

func TestCipherUpdateAndDoFinalWithOffset(t *testing.T) {
	globals.InitGlobals("test")

	// Setup Cipher
	trans := "AES/CBC/PKCS5Padding"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	keyBytes := []byte("1234567812345678")
	keyObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(keyBytes))
	cipherInit([]any{cipherObj, int64(1), keyObj})

	// cipherUpdate with offset and length
	input := []byte("0123456789ABCDEF")
	inputObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(input))
	cipherUpdate([]any{cipherObj, inputObj, int64(2), int64(4)}) // "2345"

	buffered := cipherObj.FieldTable["buffer"].Fvalue.([]byte)
	if !bytes.Equal(buffered, []byte("2345")) {
		t.Errorf("Expected buffer 2345, got %s", string(buffered))
	}

	// multiple updates
	cipherUpdate([]any{cipherObj, inputObj, int64(6), int64(2)}) // "67"
	buffered = cipherObj.FieldTable["buffer"].Fvalue.([]byte)
	if !bytes.Equal(buffered, []byte("234567")) {
		t.Errorf("Expected buffer 234567, got %s", string(buffered))
	}

	// cipherDoFinal with data, offset and length
	res := cipherDoFinal([]any{cipherObj, inputObj, int64(10), int64(3)}) // "ABC"
	if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("cipherDoFinal (encrypt) failed: %v", errBlk.ErrMsg)
	}
	// Total encrypted: "234567" + "ABC" = "234567ABC"
	encryptedObj := res.(*object.Object)

	// Now decrypt to verify
	ivObj := cipherObj.FieldTable["iv"].Fvalue.(*object.Object)
	ivSpecClass := "javax/crypto/spec/IvParameterSpec"
	ivSpecObj := object.MakeEmptyObjectWithClassName(&ivSpecClass)
	ivSpecObj.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: ivObj}

	cipherInit([]any{cipherObj, int64(2), keyObj, ivSpecObj}) // DECRYPT_MODE
	// We need the IV that was auto-generated or the one we started with.
	// Since we didn't provide one, and it's AES/CBC, it was generated during encrypt init.

	resDec := cipherDoFinal([]any{cipherObj, encryptedObj})
	if errBlk, ok := resDec.(*ghelpers.GErrBlk); ok {
		t.Fatalf("cipherDoFinal (decrypt) failed: %v", errBlk.ErrMsg)
	}
	decryptedObj := resDec.(*object.Object)
	decryptedJBytes := decryptedObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	decrypted := object.GoByteArrayFromJavaByteArray(decryptedJBytes)

	expected := []byte("234567ABC")
	if !bytes.Equal(decrypted, expected) {
		t.Errorf("Expected result %s, got %s", string(expected), string(decrypted))
	}
}

func TestCipherGetters(t *testing.T) {
	globals.InitGlobals("test")

	trans := "AES/CBC/PKCS5Padding"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	// getAlgorithm
	resAlgo := cipherGetAlgorithm([]any{cipherObj})
	if object.GoStringFromStringObject(resAlgo.(*object.Object)) != trans {
		t.Errorf("Expected algorithm %s, got %v", trans, resAlgo)
	}

	// getBlockSize (AES should be 16)
	resBS := cipherGetBlockSize([]any{cipherObj})
	if resBS.(int64) != 16 {
		t.Errorf("Expected block size 16, got %v", resBS)
	}

	// getOutputSize
	resOS := cipherGetOutputSize([]any{cipherObj, int64(10)})
	// 10 + blockSize (16) = 26
	if resOS.(int64) != 26 {
		t.Errorf("Expected output size 26, got %v", resOS)
	}
}

func TestCipherInitVariations(t *testing.T) {
	globals.InitGlobals("test")

	// Setup Cipher
	trans := "AES/CBC/PKCS5Padding"
	transObj := object.StringObjectFromGoString(trans)
	cipherObj := cipherGetInstance([]any{transObj}).(*object.Object)

	keyBytes := []byte("1234567812345678")
	keyObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(keyBytes))

	// init(opmode, key)
	cipherInit([]any{cipherObj, int64(2), keyObj}) // DECRYPT_MODE
	if cipherObj.FieldTable["opmode"].Fvalue.(int64) != 2 {
		t.Errorf("Expected opmode 2, got %v", cipherObj.FieldTable["opmode"].Fvalue)
	}
	// For DECRYPT_MODE, IV should NOT be automatically generated if not provided
	if !object.IsNull(cipherObj.FieldTable["iv"].Fvalue) {
		t.Errorf("Expected null IV for DECRYPT_MODE, got %v", cipherObj.FieldTable["iv"].Fvalue)
	}

	// init(opmode, key) - ENCRYPT_MODE, should generate IV for AES/CBC
	cipherObjGenerated := cipherGetInstance([]any{transObj}).(*object.Object)
	cipherInit([]any{cipherObjGenerated, int64(1), keyObj}) // ENCRYPT_MODE
	if object.IsNull(cipherObjGenerated.FieldTable["iv"].Fvalue) {
		t.Error("Expected automatically generated IV for AES/CBC ENCRYPT_MODE, got null")
	} else {
		ivObj := cipherObjGenerated.FieldTable["iv"].Fvalue.(*object.Object)
		ivJBytes := ivObj.FieldTable["value"].Fvalue.([]types.JavaByte)
		if len(ivJBytes) != 16 {
			t.Errorf("Expected 16-byte generated IV for AES, got %d", len(ivJBytes))
		}
	}

	// init(opmode, key) - ENCRYPT_MODE, should NOT generate IV for ECB
	transECB := "AES/ECB/PKCS5Padding"
	transECBObj := object.StringObjectFromGoString(transECB)
	cipherObjECB := cipherGetInstance([]any{transECBObj}).(*object.Object)
	cipherInit([]any{cipherObjECB, int64(1), keyObj})
	if !object.IsNull(cipherObjECB.FieldTable["iv"].Fvalue) {
		t.Errorf("Expected null IV for AES/ECB, got %v", cipherObjECB.FieldTable["iv"].Fvalue)
	}

	// init(opmode, key, IvParameterSpec) - with null IvParameterSpec
	cipherInit([]any{cipherObj, int64(1), keyObj, object.Null})
	// In the current implementation, it doesn't clear the previous IV if it was set,
	// because it checks !object.IsNull(spec).
	// But cipherInit always re-initializes opmode and buffer.
	// Let's re-verify behavior.

	// Fresh object
	cipherObj2 := cipherGetInstance([]any{transObj}).(*object.Object)
	cipherInit([]any{cipherObj2, int64(1), keyObj, object.Null})
	if !object.IsNull(cipherObj2.FieldTable["iv"].Fvalue) {
		t.Errorf("Expected null IV after init with null spec, got %v", cipherObj2.FieldTable["iv"].Fvalue)
	}

	// init(opmode, key, someOtherSpec)
	otherSpecClass := "javax/crypto/spec/SomeOtherSpec"
	otherSpecObj := object.MakeEmptyObjectWithClassName(&otherSpecClass)
	cipherInit([]any{cipherObj2, int64(1), keyObj, otherSpecObj})
	if !object.IsNull(cipherObj2.FieldTable["iv"].Fvalue) {
		t.Errorf("Expected null IV after init with non-IvParameterSpec, got %v", cipherObj2.FieldTable["iv"].Fvalue)
	}

	// Test this=null error
	resErr := cipherInit([]any{object.Null, int64(1), keyObj})
	errBlk, ok := resErr.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for null 'this', got %v", resErr)
	}
}
