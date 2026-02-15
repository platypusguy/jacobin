/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"math/big"
	"testing"

	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

func TestGenerateDHKeyPair_Success(t *testing.T) {
	globals.InitGlobals("test")

	// 1. Success with l = 0
	pVal := int64(23)
	gVal := int64(5)
	kpg := makeDHKeyPairGenerator(pVal, gVal, 0)
	privKey, pubKey, err := generateDHKeyPair(kpg)
	if err != nil {
		t.Fatalf("generateDHKeyPair failed: %v", err)
	}

	if privKey == nil || pubKey == nil {
		t.Fatal("Expected non-nil keys")
	}

	// Verify types
	if privKey.KlassName != stringPool.GetStringIndex(&types.ClassNameDHPrivateKey) {
		t.Errorf("Expected DHPrivateKey, got %s", object.GoStringFromStringPoolIndex(privKey.KlassName))
	}
	if pubKey.KlassName != stringPool.GetStringIndex(&types.ClassNameDHPublicKey) {
		t.Errorf("Expected DHPublicKey, got %s", object.GoStringFromStringPoolIndex(pubKey.KlassName))
	}

	// Verify fields in private key
	xBI := privKey.FieldTable["x"].Fvalue.(*big.Int)
	if xBI == nil {
		t.Error("Private key 'x' should be *big.Int")
	}
	pBI := privKey.FieldTable["p"].Fvalue.(*big.Int)
	if pBI.Int64() != pVal {
		t.Error("Private key 'p' mismatch")
	}
	gBI := privKey.FieldTable["g"].Fvalue.(*big.Int)
	if gBI.Int64() != gVal {
		t.Error("Private key 'g' mismatch")
	}
	if privKey.FieldTable["l"].Fvalue.(int64) != 0 {
		t.Error("Private key 'l' mismatch")
	}

	// Verify fields in public key
	yBI := pubKey.FieldTable["y"].Fvalue.(*big.Int)
	if yBI == nil {
		t.Error("Public key 'y' should be *big.Int")
	}
	if pubKey.FieldTable["p"].Fvalue.(*big.Int) != pBI {
		t.Error("Public key 'p' should match private key 'p'")
	}
	if pubKey.FieldTable["g"].Fvalue.(*big.Int) != gBI {
		t.Error("Public key 'g' should match private key 'g'")
	}
	if pubKey.FieldTable["l"].Fvalue.(int64) != 0 {
		t.Error("Public key 'l' mismatch")
	}

	// 2. Success with l > 0
	kpgL := makeDHKeyPairGenerator(23, 5, 10)
	privKeyL, pubKeyL, err := generateDHKeyPair(kpgL)
	if err != nil {
		t.Fatalf("generateDHKeyPair failed with l=10: %v", err)
	}
	if privKeyL == nil || pubKeyL == nil {
		t.Fatal("Expected non-nil keys with l=10")
	}
}

func TestGenerateDHKeyPair_MissingFields(t *testing.T) {
	globals.InitGlobals("test")

	tests := []struct {
		name      string
		removeKey string
	}{
		{"missing p", "p"},
		{"missing g", "g"},
		{"missing l", "l"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kpg := makeDHKeyPairGenerator(23, 5, 0)
			delete(kpg.FieldTable, tt.removeKey)
			_, _, err := generateDHKeyPair(kpg)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestRandomBigInt(t *testing.T) {
	t.Run("invalid bit size", func(t *testing.T) {
		_, err := randomBigInt(0)
		if err == nil {
			t.Error("Expected error for 0 bits")
		}
		_, err = randomBigInt(-1)
		if err == nil {
			t.Error("Expected error for -1 bits")
		}
	})

	t.Run("valid bit size", func(t *testing.T) {
		bits := 10
		for range 100 {
			val, err := randomBigInt(bits)
			if err != nil {
				t.Fatalf("randomBigInt failed: %v", err)
			}
			if val.Sign() <= 0 {
				t.Errorf("Expected positive value, got %v", val)
			}
			max := new(big.Int).Lsh(big.NewInt(1), uint(bits))
			if val.Cmp(max) >= 0 {
				t.Errorf("Value %v too large for %d bits", val, bits)
			}
		}
	})
}
