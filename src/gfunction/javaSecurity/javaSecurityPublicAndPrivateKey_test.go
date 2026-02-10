/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"testing"
)

func TestLoad_PublicAndPrivateKeys(t *testing.T) {
	globals.InitGlobals("test")

	// Clear any existing signatures to ensure clean test
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	// Call the loader
	Load_PublicAndPrivateKeys()

	// Verify PublicKey <init>
	if gm, ok := ghelpers.MethodSignatures["java/security/PublicKey.<init>()V"]; ok {
		if gm.ParamSlots != 0 {
			t.Errorf("PublicKey.<init>: expected 0 param slots, got %d", gm.ParamSlots)
		}
	} else {
		t.Error("PublicKey.<init> signature not registered")
	}

	// Verify PrivateKey <init>
	if gm, ok := ghelpers.MethodSignatures["java/security/PrivateKey.<init>()V"]; ok {
		if gm.ParamSlots != 0 {
			t.Errorf("PrivateKey.<init>: expected 0 param slots, got %d", gm.ParamSlots)
		}
	} else {
		t.Error("PrivateKey.<init> signature not registered")
	}
}
