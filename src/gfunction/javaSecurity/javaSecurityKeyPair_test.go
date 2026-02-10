/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestLoad_Security_KeyPair(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Security_KeyPair()

	expectedSignatures := []string{
		"java/security/KeyPair.<init>(Ljava/security/PublicKey;Ljava/security/PrivateKey;)V",
		"java/security/KeyPair.getPublic()Ljava/security/PublicKey;",
		"java/security/KeyPair.getPrivate()Ljava/security/PrivateKey;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s not found", sig)
		}
	}
}

func TestKeyPairGFunctions(t *testing.T) {
	globals.InitGlobals("test")

	// Create mock public and private key objects
	pubKey := object.MakeEmptyObjectWithClassName(&types.ClassNamePublicKey)
	privKey := object.MakeEmptyObjectWithClassName(&types.ClassNamePrivateKey)

	// Create KeyPair object
	keyPairObj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPair)

	// Test keypairInit
	params := []any{keyPairObj, pubKey, privKey}
	res := keypairInit(params)
	if res != nil {
		t.Fatalf("keypairInit failed: %v", res)
	}

	// Verify fields
	if keyPairObj.FieldTable["public"].Fvalue != pubKey {
		t.Errorf("Expected public key field to be %v, got %v", pubKey, keyPairObj.FieldTable["public"].Fvalue)
	}
	if keyPairObj.FieldTable["private"].Fvalue != privKey {
		t.Errorf("Expected private key field to be %v, got %v", privKey, keyPairObj.FieldTable["private"].Fvalue)
	}

	// Test keypairGetPublic
	resPub := keypairGetPublic([]any{keyPairObj})
	if resPub != pubKey {
		t.Errorf("keypairGetPublic returned %v, expected %v", resPub, pubKey)
	}

	// Test keypairGetPrivate
	resPriv := keypairGetPrivate([]any{keyPairObj})
	if resPriv != privKey {
		t.Errorf("keypairGetPrivate returned %v, expected %v", resPriv, privKey)
	}
}

func TestKeyPairInvalidParams(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("keypairInit_MissingParams", func(t *testing.T) {
		res := keypairInit([]any{nil, nil})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Error("keypairInit should return GErrBlk for missing params")
		}
	})

	t.Run("keypairInit_WrongType", func(t *testing.T) {
		res := keypairInit([]any{"not an object", nil, nil})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Error("keypairInit should return GErrBlk for wrong type param[0]")
		}
	})

	t.Run("keypairGetPublic_WrongType", func(t *testing.T) {
		res := keypairGetPublic([]any{"not an object"})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Error("keypairGetPublic should return GErrBlk for wrong type param[0]")
		}
	})

	t.Run("keypairGetPrivate_WrongType", func(t *testing.T) {
		res := keypairGetPrivate([]any{"not an object"})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Error("keypairGetPrivate should return GErrBlk for wrong type param[0]")
		}
	})
}
