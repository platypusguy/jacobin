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

func TestInitDefaultSecurityProvider(t *testing.T) {
	globals.InitGlobals("test")

	provider := InitDefaultSecurityProvider()

	if provider == nil {
		t.Fatal("InitDefaultSecurityProvider returned nil")
	}

	// Verify provider identity
	nameObj := provider.FieldTable["name"].Fvalue.(*object.Object)
	name := object.GoStringFromStringObject(nameObj)
	if name != types.SecurityProviderName {
		t.Errorf("Expected provider name %s, got %s", types.SecurityProviderName, name)
	}

	// Verify that all services in SecurityProviderServices are registered
	services := provider.FieldTable["services"].Fvalue.(map[string]*object.Object)

	for typeStr, algos := range SecurityProviderServices {
		for algoStr := range algos {
			key := typeStr + "/" + algoStr
			if _, exists := services[key]; !exists {
				// Special case: "DH" and "DiffieHellman" might map to the same key if the provider
				// implementation normalizes algorithm names, or if it only stores one of them.
				// In javaSecurity_Helpers_ProviderServices.go, DH and DiffieHellman both return "DH" as algo.
				// securityProviderPutService uses svc.FieldTable["algorithm"] to build the key.

				// Let's check what the algorithm in the service is.
				svcInit := algos[algoStr]
				svc := svcInit()
				svcAlgo := object.GoStringFromStringObject(svc.FieldTable["algorithm"].Fvalue.(*object.Object))
				actualKey := typeStr + "/" + svcAlgo

				if _, exists2 := services[actualKey]; !exists2 {
					t.Errorf("Service %s (actual key %s) not registered in provider", key, actualKey)
				}
			}
		}
	}
}

func TestNewGoRuntimeService(t *testing.T) {
	globals.InitGlobals("test")

	// Ensure DefaultSecurityProvider is set for NewGoRuntimeService
	provider := InitDefaultSecurityProvider()
	ghelpers.DefaultSecurityProvider = provider

	typ := "MessageDigest"
	algo := "SHA-256"
	className := types.ClassNameMessageDigest

	svc := NewGoRuntimeService(typ, algo, className)

	if svc == nil {
		t.Fatal("NewGoRuntimeService returned nil")
	}

	// Verify fields
	if got := object.GoStringFromStringObject(svc.FieldTable["type"].Fvalue.(*object.Object)); got != typ {
		t.Errorf("Expected type %s, got %s", typ, got)
	}
	if got := object.GoStringFromStringObject(svc.FieldTable["algorithm"].Fvalue.(*object.Object)); got != algo {
		t.Errorf("Expected algorithm %s, got %s", algo, got)
	}
	if got := object.GoStringFromStringObject(svc.FieldTable["className"].Fvalue.(*object.Object)); got != className {
		t.Errorf("Expected className %s, got %s", className, got)
	}

	// Verify attributes
	attributes := svc.FieldTable["attributes"].Fvalue.(map[string]*object.Object)
	if got := object.GoStringFromStringObject(attributes["ImplementedIn"]); got != "Software" {
		t.Errorf("Expected ImplementedIn Software, got %s", got)
	}

	// SHA-256 should have block size 64
	if got := object.GoStringFromStringObject(attributes["blockSize"]); got != "64" {
		t.Errorf("Expected blockSize 64 for SHA-256, got %s", got)
	}

	// Verify provider
	if svc.FieldTable["provider"].Fvalue != provider {
		t.Errorf("Provider not correctly set in service")
	}
}

func TestGetBlockSizeForAlgorithm(t *testing.T) {
	tests := []struct {
		algo string
		want int
	}{
		{"MD5", 64},
		{"sha-1", 64},
		{"SHA-224", 64},
		{"SHA-256", 64},
		{"SHA-384", 128},
		{"SHA-512", 128},
		{"SHA-512/224", 128},
		{"SHA-512/256", 128},
		{"AES", 0},
		{"RSA", 0},
	}

	for _, tt := range tests {
		t.Run(tt.algo, func(t *testing.T) {
			if got := getBlockSizeForAlgorithm(tt.algo); got != tt.want {
				t.Errorf("getBlockSizeForAlgorithm(%s) = %d, want %d", tt.algo, got, tt.want)
			}
		})
	}
}
