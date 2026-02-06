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
	"testing"
)

// TestLoadSecurity verifies that Load_Security() properly initializes method signatures
func TestLoadSecurity(t *testing.T) {
	globals.InitGlobals("test")

	// Clear any existing signatures to ensure clean test
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	// Call Load_Security
	Load_Security()

	// Test that getProvider method signature is registered
	if _, ok := ghelpers.MethodSignatures["java/security/Security.getProvider(Ljava/lang/String;)Ljava/security/Provider;"]; !ok {
		t.Error("getProvider method signature not registered")
	}

	// Test that getProviders method signature is registered
	if _, ok := ghelpers.MethodSignatures["java/security/Security.getProviders()[Ljava/security/Provider;"]; !ok {
		t.Error("getProviders method signature not registered")
	}

	// Test that addProvider method signature is registered
	if _, ok := ghelpers.MethodSignatures["java/security/Security.addProvider(Ljava/security/Provider;)I"]; !ok {
		t.Error("addProvider method signature not registered")
	}

	// Test that insertProviderAt method signature is registered
	if _, ok := ghelpers.MethodSignatures["java/security/Security.insertProviderAt(Ljava/security/Provider;I)I"]; !ok {
		t.Error("insertProviderAt method signature not registered")
	}

	// Test that removeProvider method signature is registered
	if _, ok := ghelpers.MethodSignatures["java/security/Security.removeProvider(Ljava/lang/String;)V"]; !ok {
		t.Error("removeProvider method signature not registered")
	}
}

// TestLoadSecurityParamSlots verifies that param slots are correctly set
func TestLoadSecurityParamSlots(t *testing.T) {
	globals.InitGlobals("test")

	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Security()

	tests := []struct {
		name       string
		method     string
		paramSlots int
	}{
		{"getProvider", "java/security/Security.getProvider(Ljava/lang/String;)Ljava/security/Provider;", 1},
		{"getProviders", "java/security/Security.getProviders()[Ljava/security/Provider;", 0},
		{"addProvider", "java/security/Security.addProvider(Ljava/security/Provider;)I", 1},
		{"insertProviderAt", "java/security/Security.insertProviderAt(Ljava/security/Provider;I)I", 2},
		{"removeProvider", "java/security/Security.removeProvider(Ljava/lang/String;)V", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gmeth, ok := ghelpers.MethodSignatures[tt.method]; ok {
				if gmeth.ParamSlots != tt.paramSlots {
					t.Errorf("%s: expected %d param slots, got %d", tt.name, tt.paramSlots, gmeth.ParamSlots)
				}
			} else {
				t.Errorf("%s: method signature not found", tt.name)
			}
		})
	}
}

// TestSecurityGetProvider verifies that securityGetProvider returns DefaultSecurityProvider
func TestSecurityGetProvider(t *testing.T) {
	globals.InitGlobals("test")

	result := securityGetProvider([]any{})
	expected := ghelpers.GetDefaultSecurityProvider()

	if result != expected {
		t.Errorf("securityGetProvider should return DefaultSecurityProvider, got %v", result)
	}
}

// TestSecurityGetProviders verifies that securityGetProviders returns an array with DefaultSecurityProvider
func TestSecurityGetProviders(t *testing.T) {
	globals.InitGlobals("test")

	result := securityGetProviders([]any{})
	expected := ghelpers.GetDefaultSecurityProvider()

	obj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("securityGetProviders should return *object.Object, got %T", result)
	}
	providers := obj.FieldTable["value"].Fvalue.([]*object.Object)

	if len(providers) != 1 {
		t.Errorf("expected 1 provider, got %d", len(providers))
	}

	if len(providers) > 0 && providers[0] != expected {
		t.Errorf("expected DefaultSecurityProvider in array, got %v", providers[0])
	}
}

// TestSecurityGetProvidersWithValidProvider verifies the provider is properly initialized
func TestSecurityGetProvidersWithValidProvider(t *testing.T) {
	globals.InitGlobals("test")

	provider := ghelpers.GetDefaultSecurityProvider()
	if provider == nil {
		t.Skip("DefaultSecurityProvider is nil, skipping test")
	}
	result := securityGetProviders([]any{})

	obj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("securityGetProviders should return *object.Object, got %T", result)
	}
	providers := obj.FieldTable["value"].Fvalue.([]*object.Object)

	if len(providers) == 0 {
		t.Fatal("expected at least one provider when DefaultSecurityProvider is not nil")
	}

	// Verify it's an object
	if providers[0] == nil {
		t.Error("provider should not be nil")
	}
}

// TestSecurityGetProvidersEmptyParams tests that the function works with empty parameter slice
func TestSecurityGetProvidersEmptyParams(t *testing.T) {
	globals.InitGlobals("test")

	// Test with nil params
	result := securityGetProviders(nil)
	if _, ok := result.(*object.Object); !ok {
		t.Errorf("securityGetProviders with nil params should still return *object.Object, got %T", result)
	}

	// Test with empty slice
	result = securityGetProviders([]any{})
	if _, ok := result.(*object.Object); !ok {
		t.Errorf("securityGetProviders with empty params should return *object.Object, got %T", result)
	}
}
