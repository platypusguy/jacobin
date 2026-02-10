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
	"crypto/rsa"
	"math/big"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

// makeStringObject creates a mock Java String object
func makeStringObject(s string) *object.Object {
	return object.StringObjectFromGoString(s)
}

// makeKeyPairGenerator creates the input object for keypairgeneratorGenerateKeyPair
func makeKeyPairGenerator(algo string, keySize int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)
	obj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: makeStringObject(algo),
	}
	if keySize >= 0 {
		obj.FieldTable["keySize"] = object.Field{
			Ftype:  types.Int,
			Fvalue: keySize,
		}
	}
	return obj
}

func TestGenerateKeyPairRSA(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	kpg := makeKeyPairGenerator("RSA", 2048)
	result := keypairgeneratorGenerateKeyPair([]any{kpg})

	kpObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	if kpObj.KlassName != stringPool.GetStringIndex(&types.ClassNameKeyPair) {
		t.Errorf("Expected ClassNameKeyPair, got index %d (expected %d)", kpObj.KlassName, stringPool.GetStringIndex(&types.ClassNameKeyPair))
	}

	privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
	pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

	if privObj.FieldTable["value"].Fvalue.(*rsa.PrivateKey) == nil {
		t.Error("RSA Private key value is nil")
	}
	if pubObj.FieldTable["value"].Fvalue.(*rsa.PublicKey) == nil {
		t.Error("RSA Public key value is nil")
	}
}

func TestGenerateKeyPairDSA(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	kpg := makeKeyPairGenerator("DSA", 2048)
	result := keypairgeneratorGenerateKeyPair([]any{kpg})

	kpObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
	pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

	if privObj.FieldTable["value"].Fvalue.(*dsa.PrivateKey) == nil {
		t.Error("DSA Private key value is nil")
	}
	if pubObj.FieldTable["value"].Fvalue.(*dsa.PublicKey) == nil {
		t.Error("DSA Public key value is nil")
	}
}

func TestGenerateKeyPairEC(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	kpg := makeKeyPairGenerator("EC", 256)
	result := keypairgeneratorGenerateKeyPair([]any{kpg})

	kpObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
	pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

	if privObj.FieldTable["value"].Fvalue.(*ecdsa.PrivateKey) == nil {
		t.Error("EC Private key value is nil")
	}
	if pubObj.FieldTable["value"].Fvalue.(*ecdsa.PublicKey) == nil {
		t.Error("EC Public key value is nil")
	}
}

func TestGenerateKeyPairEd25519(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	kpg := makeKeyPairGenerator("Ed25519", -1)
	result := keypairgeneratorGenerateKeyPair([]any{kpg})

	kpObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
	pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

	if _, ok := privObj.FieldTable["value"].Fvalue.(ed25519.PrivateKey); !ok {
		t.Errorf("Expected ed25519.PrivateKey, got %T", privObj.FieldTable["value"].Fvalue)
	}
	if _, ok := pubObj.FieldTable["value"].Fvalue.(ed25519.PublicKey); !ok {
		t.Errorf("Expected ed25519.PublicKey, got %T", pubObj.FieldTable["value"].Fvalue)
	}
}

func TestGenerateKeyPairXDH(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	algorithms := []string{"XDH", "X25519"}
	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			kpg := makeKeyPairGenerator(algo, -1)
			result := keypairgeneratorGenerateKeyPair([]any{kpg})

			kpObj, ok := result.(*object.Object)
			if !ok {
				t.Fatalf("Expected *object.Object, got %T", result)
			}

			privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
			pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

			if _, ok := privObj.FieldTable["value"].Fvalue.([]byte); !ok {
				t.Errorf("Expected []byte for private key, got %T", privObj.FieldTable["value"].Fvalue)
			}
			if _, ok := pubObj.FieldTable["value"].Fvalue.([]byte); !ok {
				t.Errorf("Expected []byte for public key, got %T", pubObj.FieldTable["value"].Fvalue)
			}
		})
	}
}

func TestGenerateKeyPairDH(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	kpg := makeKeyPairGenerator("DH", -1)
	result := keypairgeneratorGenerateKeyPair([]any{kpg})

	kpObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}

	privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
	pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

	if privObj.FieldTable["value"].Fvalue.(*big.Int) == nil {
		t.Error("DH Private key value is nil")
	}
	if pubObj.FieldTable["value"].Fvalue.(*big.Int) == nil {
		t.Error("DH Public key value is nil")
	}
}

func TestGenerateKeyPairEd448X448(t *testing.T) {
	globals.InitGlobals("test")
	InitDefaultSecurityProvider()

	algorithms := []string{"Ed448", "X448"}
	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			kpg := makeKeyPairGenerator(algo, -1)
			result := keypairgeneratorGenerateKeyPair([]any{kpg})

			kpObj, ok := result.(*object.Object)
			if !ok {
				t.Fatalf("Expected *object.Object, got %T", result)
			}

			privObj := kpObj.FieldTable["private"].Fvalue.(*object.Object)
			pubObj := kpObj.FieldTable["public"].Fvalue.(*object.Object)

			if _, ok := privObj.FieldTable["value"].Fvalue.([]byte); !ok {
				t.Errorf("Expected []byte for private key, got %T", privObj.FieldTable["value"].Fvalue)
			}
			if _, ok := pubObj.FieldTable["value"].Fvalue.([]byte); !ok {
				t.Errorf("Expected []byte for public key, got %T", pubObj.FieldTable["value"].Fvalue)
			}
		})
	}
}

func TestGenerateKeyPairInvalidParams(t *testing.T) {
	globals.InitGlobals("test")

	// Missing params
	result := keypairgeneratorGenerateKeyPair([]any{})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for missing params, got %T", result)
	}

	// Wrong param type
	result = keypairgeneratorGenerateKeyPair([]any{"not an object"})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for wrong param type, got %T", result)
	}

	// Missing algorithm
	kpg := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)
	result = keypairgeneratorGenerateKeyPair([]any{kpg})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for missing algorithm, got %T", result)
	}
}
