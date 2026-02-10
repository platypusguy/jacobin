package javaSecurity

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Helper to create RSAPrivateKey and RSAPublicKey objects similar to generateKeyPair
func makeRSAKeyObjects(t *testing.T, bits int) (*object.Object, *object.Object) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		t.Fatalf("rsa.GenerateKey failed: %v", err)
	}

	// Public key object
	pubObj := NewGoRuntimeService("RSAPublicKey", "RSA", types.ClassNameRSAPublicKey)
	pubObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: &priv.PublicKey}

	// Private key object
	prvObj := NewGoRuntimeService("RSAPrivateKey", "RSA", types.ClassNameRSAPrivateKey)
	prvObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

	return prvObj, pubObj
}

func TestLoad_Security_Interfaces_RSA_Keys(t *testing.T) {
	globals.InitGlobals("test")

	// Clear and load
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Security_Interfaces_RSA_Keys()

	expected := []string{
		"java/security/interfaces/RSAKey.<clinit>()V",
		"java/security/interfaces/RSAKey.<init>()V",
		"java/security/interfaces/RSAKey.getModulus()Ljava/math/BigInteger;",
		"java/security/interfaces/RSAKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;",
		"java/security/interfaces/RSAPrivateKey.<clinit>()V",
		"java/security/interfaces/RSAPrivateKey.<init>()V",
		"java/security/interfaces/RSAPrivateKey.getPrivateExponent()Ljava/math/BigInteger;",
		"java/security/interfaces/RSAPublicKey.<clinit>()V",
		"java/security/interfaces/RSAPublicKey.<init>()V",
		"java/security/interfaces/RSAPublicKey.getModulus()Ljava/math/BigInteger;",
		"java/security/interfaces/RSAPublicKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;",
	}

	for _, sig := range expected {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("missing registration for %s", sig)
		}
	}
}

func TestRSAKeyGetModulus_PublicKey(t *testing.T) {
	globals.InitGlobals("test")
	prvObj, pubObj := makeRSAKeyObjects(t, 1024)

	// Silence unused
	_ = prvObj

	// Call rsaKeyGetModulus on public key object
	out := rsaKeyGetModulus([]any{pubObj})
	bi, ok := out.(*object.Object)
	if !ok {
		t.Fatalf("expected BigInteger object, got %T", out)
	}

	// The BigInteger object stores a *big.Int in value
	N, ok := bi.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		t.Fatalf("expected BigInteger underlying *big.Int, got %T", bi.FieldTable["value"].Fvalue)
	}
	if N.Cmp(pubObj.FieldTable["value"].Fvalue.(*rsa.PublicKey).N) != 0 {
		t.Errorf("modulus mismatch for public key")
	}
}

func TestRSAKeyGetModulus_PrivateKey(t *testing.T) {
	globals.InitGlobals("test")
	prvObj, _ := makeRSAKeyObjects(t, 1024)

	out := rsaKeyGetModulus([]any{prvObj})
	bi, ok := out.(*object.Object)
	if !ok {
		t.Fatalf("expected BigInteger object, got %T", out)
	}
	N, ok := bi.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		t.Fatalf("expected *big.Int in BigInteger value, got %T", bi.FieldTable["value"].Fvalue)
	}
	if N.Cmp(prvObj.FieldTable["value"].Fvalue.(*rsa.PrivateKey).N) != 0 {
		t.Errorf("modulus mismatch for private key")
	}
}

func TestRSAPrivateKeyGetPrivateExponent(t *testing.T) {
	globals.InitGlobals("test")
	prvObj, _ := makeRSAKeyObjects(t, 1024)

	out := rsaprivateGetExponent([]any{prvObj})
	bi, ok := out.(*object.Object)
	if !ok {
		t.Fatalf("expected BigInteger object, got %T", out)
	}
	D, ok := bi.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		t.Fatalf("expected *big.Int in BigInteger value, got %T", bi.FieldTable["value"].Fvalue)
	}
	if D.Cmp(prvObj.FieldTable["value"].Fvalue.(*rsa.PrivateKey).D) != 0 {
		t.Errorf("private exponent mismatch")
	}
}

func TestRSAInterfaces_InvalidParams(t *testing.T) {
	globals.InitGlobals("test")

	// Missing params
	if _, ok := rsaKeyGetModulus([]any{}).(*ghelpers.GErrBlk); !ok {
		t.Error("expected error for missing params in rsaKeyGetModulus")
	}

	// Wrong this type
	if _, ok := rsaKeyGetModulus([]any{123}).(*ghelpers.GErrBlk); !ok {
		t.Error("expected error for wrong this type in rsaKeyGetModulus")
	}

	// Wrong value type in this
	bogus := object.MakeEmptyObjectWithClassName(&types.ClassNameRSAPublicKey)
	bogus.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: "not-a-key"}
	if _, ok := rsaKeyGetModulus([]any{bogus}).(*ghelpers.GErrBlk); !ok {
		t.Error("expected error for wrong value type in rsaKeyGetModulus")
	}

	// Missing this for rsaprivateGetExponent
	if _, ok := rsaprivateGetExponent([]any{}).(*ghelpers.GErrBlk); !ok {
		t.Error("expected error for missing params in rsaprivateGetExponent")
	}
}
