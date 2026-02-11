package javaSecurity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestKeyMethods(t *testing.T) {
	globals.InitGlobals("test")

	// 1. RSA Public Key
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	pubObj := NewGoRuntimeService("RSAPublicKey", "RSA", types.ClassNameRSAPublicKey)
	pubObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}

	// Test getAlgorithm
	alg := keyGetAlgorithm([]any{pubObj})
	if object.GoStringFromStringObject(alg.(*object.Object)) != "RSA" {
		t.Errorf("Expected RSA, got %v", alg)
	}

	// Test getFormat
	format := keyGetFormat([]any{pubObj})
	if object.GoStringFromStringObject(format.(*object.Object)) != "X.509" {
		t.Errorf("Expected X.509, got %v", format)
	}

	// Test getEncoded
	encoded := keyGetEncoded([]any{pubObj}).(*object.Object)
	expectedEncoded, _ := x509.MarshalPKIXPublicKey(pub)
	if !object.JavaByteArrayEquals(encoded.FieldTable["value"].Fvalue.([]types.JavaByte), object.JavaByteArrayFromGoByteArray(expectedEncoded)) {
		t.Error("Encoded bytes mismatch for RSA Public Key")
	}

	// 2. RSA Private Key
	privObj := NewGoRuntimeService("RSAPrivateKey", "RSA", types.ClassNameRSAPrivateKey)
	privObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

	format = keyGetFormat([]any{privObj})
	if object.GoStringFromStringObject(format.(*object.Object)) != "PKCS#8" {
		t.Errorf("Expected PKCS#8, got %v", format)
	}

	encoded = keyGetEncoded([]any{privObj}).(*object.Object)
	expectedEncoded, _ = x509.MarshalPKCS8PrivateKey(priv)
	if !object.JavaByteArrayEquals(encoded.FieldTable["value"].Fvalue.([]types.JavaByte), object.JavaByteArrayFromGoByteArray(expectedEncoded)) {
		t.Error("Encoded bytes mismatch for RSA Private Key")
	}

	// 3. Ed25519
	pubEd, _, _ := ed25519.GenerateKey(rand.Reader)
	pubEdObj := NewGoRuntimeService("Ed25519", "Ed25519", types.ClassNameEdECPublicKey)
	pubEdObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubEd}

	alg = keyGetAlgorithm([]any{pubEdObj})
	if object.GoStringFromStringObject(alg.(*object.Object)) != "Ed25519" {
		t.Errorf("Expected Ed25519, got %v", alg)
	}

	format = keyGetFormat([]any{pubEdObj})
	if object.GoStringFromStringObject(format.(*object.Object)) != "X.509" {
		t.Errorf("Expected X.509, got %v", format)
	}

	encoded = keyGetEncoded([]any{pubEdObj}).(*object.Object)
	expectedEncoded, _ = x509.MarshalPKIXPublicKey(pubEd)
	if !object.JavaByteArrayEquals(encoded.FieldTable["value"].Fvalue.([]types.JavaByte), object.JavaByteArrayFromGoByteArray(expectedEncoded)) {
		t.Error("Encoded bytes mismatch for Ed25519 Public Key")
	}

	// 4. hashCode
	h1 := keyHashCode([]any{pubObj})
	h2 := keyHashCode([]any{pubObj})
	if h1 != h2 {
		t.Errorf("Hash codes should be equal for the same object, got %v and %v", h1, h2)
	}
}
