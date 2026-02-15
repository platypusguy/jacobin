/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"slices"
	"testing"
)

func TestValidateCipherTransformation(t *testing.T) {
	tests := []struct {
		name           string
		transformation string
		wantExists     bool
		wantEnabled    bool
	}{
		{"Valid AES", "AES/CBC/PKCS5Padding", true, true},
		{"Valid ChaCha20", "ChaCha20", true, true},
		{"Invalid transformation", "Invalid/Algo", false, false},
		{"Empty transformation", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ValidateCipherTransformation(tt.transformation)
			if ok != tt.wantExists {
				t.Errorf("ValidateCipherTransformation() ok = %v, want %v", ok, tt.wantExists)
			}
			if ok && got.Enabled != tt.wantEnabled {
				t.Errorf("ValidateCipherTransformation() enabled = %v, want %v", got.Enabled, tt.wantEnabled)
			}
			if ok && got.Name != tt.transformation {
				t.Errorf("ValidateCipherTransformation() name = %v, want %v", got.Name, tt.transformation)
			}
		})
	}
}

func TestGetRequiredParameters(t *testing.T) {
	tests := []struct {
		name           string
		transformation CipherTransformation
		want           []string
	}{
		{
			name: "AES CBC",
			transformation: CipherTransformation{
				NeedsIV: true,
			},
			want: []string{"key", "iv"},
		},
		{
			name: "PBE MD5 DES",
			transformation: CipherTransformation{
				KeyDerivation:   true,
				NeedsSalt:       true,
				NeedsIterations: true,
			},
			want: []string{"password", "salt", "iterations"},
		},
		{
			name: "AES GCM",
			transformation: CipherTransformation{
				NeedsIV:        true,
				NeedsTagLength: true,
			},
			want: []string{"key", "iv", "tagLength"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.transformation.GetRequiredParameters()
			if !slices.Equal(got, tt.want) {
				t.Errorf("GetRequiredParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnabledTransformations(t *testing.T) {
	enabled := GetEnabledTransformations()
	if len(enabled) == 0 {
		t.Error("GetEnabledTransformations() returned empty map")
	}
	for name, config := range enabled {
		if !config.Enabled {
			t.Errorf("Transformation %s is in enabled map but Enabled field is false", name)
		}
	}
}

func TestDisableEnableTransformation(t *testing.T) {
	name := "AES/CBC/PKCS5Padding"

	// Ensure it starts enabled (or at least exists)
	if _, exists := CipherConfigTable[name]; !exists {
		t.Fatalf("Test transformation %s does not exist in CipherConfigTable", name)
	}

	// Test Disable
	ok := DisableTransformation(name)
	if !ok {
		t.Errorf("DisableTransformation(%s) returned false", name)
	}
	if CipherConfigTable[name].Enabled {
		t.Errorf("Transformation %s should be disabled", name)
	}

	// Test Enable
	ok = EnableTransformation(name)
	if !ok {
		t.Errorf("EnableTransformation(%s) returned false", name)
	}
	if !CipherConfigTable[name].Enabled {
		t.Errorf("Transformation %s should be enabled", name)
	}

	// Test non-existent
	if DisableTransformation("NonExistent") {
		t.Error("DisableTransformation(NonExistent) should return false")
	}
	if EnableTransformation("NonExistent") {
		t.Error("EnableTransformation(NonExistent) should return false")
	}
}
