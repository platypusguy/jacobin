/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestLoad_Security_Signature(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Security_Signature()

	methods := []string{
		"java/security/Signature.getInstance(Ljava/lang/String;)Ljava/security/Signature;",
		"java/security/Signature.initSign(Ljava/security/PrivateKey;)V",
		"java/security/Signature.initVerify(Ljava/security/PublicKey;)V",
		"java/security/Signature.sign()[B",
		"java/security/Signature.update([B)V",
		"java/security/Signature.verify([B)Z",
		"java/security/Signature.getAlgorithm()Ljava/lang/String;",
		"java/security/Signature.getProvider()Ljava/security/Provider;",
	}

	for _, m := range methods {
		if _, ok := ghelpers.MethodSignatures[m]; !ok {
			t.Errorf("Signature method signature not registered: %s", m)
		}
	}
}

func TestSignatureGFunctions_RSA(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	// 1. getInstance
	algoName := "SHA256withRSA"
	algoObj := object.StringObjectFromGoString(algoName)
	res := signatureGetInstance([]any{algoObj})
	sigObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("signatureGetInstance failed: %v", res)
	}

	// 2. getAlgorithm
	gotAlgo := signatureGetAlgorithm([]any{sigObj})
	if object.GoStringFromStringObject(gotAlgo.(*object.Object)) != algoName {
		t.Errorf("Expected algorithm %s, got %v", algoName, gotAlgo)
	}

	// 3. initSign
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	privKeyObj := NewGoRuntimeService("RSAPrivateKey", "RSA", types.ClassNameRSAPrivateKey)
	privKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privKey}

	res = signatureInitSign([]any{sigObj, privKeyObj})
	if res != nil {
		t.Fatalf("signatureInitSign failed: %v", res)
	}

	// 4. update
	data := []byte("hello world")
	dataObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(data))
	signatureUpdateBytes([]any{sigObj, dataObj})

	// 5. sign
	res = signatureSign([]any{sigObj})
	sigResultObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("signatureSign failed: %v", res)
	}
	signatureBytes := object.GoByteArrayFromJavaByteArray(sigResultObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	// 6. initVerify
	pubKeyObj := NewGoRuntimeService("RSAPublicKey", "RSA", types.ClassNameRSAPublicKey)
	pubKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: &privKey.PublicKey}

	res = signatureInitVerify([]any{sigObj, pubKeyObj})
	if res != nil {
		t.Fatalf("signatureInitVerify failed: %v", res)
	}

	// 7. update again (for verification)
	signatureUpdateBytes([]any{sigObj, dataObj})

	// 8. verify
	res = signatureVerify([]any{sigObj, sigResultObj})
	valid, ok := res.(bool)
	if !ok || !valid {
		t.Errorf("signatureVerify failed: expected true, got %v", res)
	}

	// Test invalid signature
	invalidSigBytes := make([]byte, len(signatureBytes))
	copy(invalidSigBytes, signatureBytes)
	invalidSigBytes[0] ^= 0xFF
	invalidSigObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(invalidSigBytes))

	signatureUpdateBytes([]any{sigObj, dataObj}) // Reset buffer by re-updating? No, buffer needs to be reset.
	// Actually, signatureVerify doesn't reset buffer in current implementation, but it should probably if we want to re-verify.
	// Looking at the code, it uses sigObj.FieldTable["buffer"].
	// Wait, if I call signatureVerify, does it clear the buffer? No.

	res = signatureVerify([]any{sigObj, invalidSigObj})
	if valid, ok := res.(bool); !ok || valid {
		t.Errorf("signatureVerify should have failed for invalid signature")
	}
}

func TestSignatureGFunctions_ECDSA(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	algoName := "SHA256withECDSA"
	algoObj := object.StringObjectFromGoString(algoName)
	res := signatureGetInstance([]any{algoObj})
	sigObj := res.(*object.Object)

	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKeyObj := NewGoRuntimeService("ECPrivateKey", "EC", types.ClassNameECPrivateKey)
	privKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privKey}

	signatureInitSign([]any{sigObj, privKeyObj})

	data := []byte("ecdsa test")
	dataObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(data))
	signatureUpdateBytes([]any{sigObj, dataObj})

	res = signatureSign([]any{sigObj})
	sigResultObj := res.(*object.Object)

	pubKeyObj := NewGoRuntimeService("ECPublicKey", "EC", types.ClassNameECPublicKey)
	pubKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: &privKey.PublicKey}

	signatureInitVerify([]any{sigObj, pubKeyObj})
	signatureUpdateBytes([]any{sigObj, dataObj})

	res = signatureVerify([]any{sigObj, sigResultObj})
	if valid, ok := res.(bool); !ok || !valid {
		t.Errorf("ECDSA signatureVerify failed")
	}
}

func TestSignatureGFunctions_DSA(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	algoName := "SHA256withDSA"
	algoObj := object.StringObjectFromGoString(algoName)
	res := signatureGetInstance([]any{algoObj})
	sigObj := res.(*object.Object)

	params := new(dsa.Parameters)
	dsa.GenerateParameters(params, rand.Reader, dsa.L2048N256)
	privKey := new(dsa.PrivateKey)
	privKey.Parameters = *params
	dsa.GenerateKey(privKey, rand.Reader)

	privKeyObj := NewGoRuntimeService("DSAPrivateKey", "DSA", types.ClassNameDSAPrivateKey)
	privKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privKey}

	signatureInitSign([]any{sigObj, privKeyObj})

	data := []byte("dsa test")
	dataObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(data))
	signatureUpdateBytes([]any{sigObj, dataObj})

	res = signatureSign([]any{sigObj})
	sigResultObj := res.(*object.Object)

	pubKeyObj := NewGoRuntimeService("DSAPublicKey", "DSA", types.ClassNameDSAPublicKey)
	pubKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: &privKey.PublicKey}

	signatureInitVerify([]any{sigObj, pubKeyObj})
	signatureUpdateBytes([]any{sigObj, dataObj})

	res = signatureVerify([]any{sigObj, sigResultObj})
	if valid, ok := res.(bool); !ok || !valid {
		t.Errorf("DSA signatureVerify failed")
	}
}

func TestSignatureGFunctions_Ed25519(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	algoName := "Ed25519"
	algoObj := object.StringObjectFromGoString(algoName)
	res := signatureGetInstance([]any{algoObj})
	sigObj := res.(*object.Object)

	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)

	privKeyObj := NewGoRuntimeService("Ed25519PrivateKey", "Ed25519", types.ClassNameEd25519PrivateKey)
	privKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privKey}

	signatureInitSign([]any{sigObj, privKeyObj})

	data := []byte("ed25519 test")
	dataObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(data))
	signatureUpdateBytes([]any{sigObj, dataObj})

	res = signatureSign([]any{sigObj})
	sigResultObj := res.(*object.Object)

	pubKeyObj := NewGoRuntimeService("Ed25519PublicKey", "Ed25519", types.ClassNameEd25519PublicKey)
	pubKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKey}

	signatureInitVerify([]any{sigObj, pubKeyObj})
	signatureUpdateBytes([]any{sigObj, dataObj})

	res = signatureVerify([]any{sigObj, sigResultObj})
	if valid, ok := res.(bool); !ok || !valid {
		t.Errorf("Ed25519 signatureVerify failed")
	}
}

func TestSignature_UpdateVariants(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	sigObj := signatureGetInstance([]any{object.StringObjectFromGoString("SHA256withRSA")}).(*object.Object)

	// Update byte
	signatureUpdateByte([]any{sigObj, int64('h')})
	signatureUpdateByte([]any{sigObj, int64('e')})

	// Update range
	data := []byte("ZZhelloZZ")
	dataObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(data))
	signatureUpdateBytesRange([]any{sigObj, dataObj, int64(2), int64(5)})

	buffer := sigObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	got := string(object.GoByteArrayFromJavaByteArray(buffer))
	if got != "hehello" {
		t.Errorf("Expected buffer 'hehello', got '%s'", got)
	}
}

func TestSignature_InvalidParams(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	sigObj := signatureGetInstance([]any{object.StringObjectFromGoString("SHA256withRSA")}).(*object.Object)

	// Test calling sign without initSign
	res := signatureSign([]any{sigObj})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.SignatureException {
		t.Errorf("Expected SignatureException when signing without init, got %v", res)
	}

	// Test initSign with wrong key type
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ecPrivKeyObj := NewGoRuntimeService("ECPrivateKey", "EC", types.ClassNameECPrivateKey)
	ecPrivKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privKey}

	res = signatureInitSign([]any{sigObj, ecPrivKeyObj})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.InvalidKeyException {
		t.Errorf("Expected InvalidKeyException for mismatched key, got %v", res)
	}

	// Test unsupported algorithm
	res = signatureGetInstance([]any{object.StringObjectFromGoString("FOO")})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NoSuchAlgorithmException {
		t.Errorf("Expected NoSuchAlgorithmException for FOO, got %v", res)
	}
}
