/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestPBEKeySpec(t *testing.T) {
	globals.InitGlobals("test")

	password := "secret"
	passwordChars := make([]int64, len(password))
	for i, c := range password {
		passwordChars[i] = int64(c)
	}
	passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)

	className := "javax/crypto/spec/PBEKeySpec"
	specObj := object.MakeEmptyObjectWithClassName(&className)

	// Test Init
	params := []any{specObj, passwordCharsObj}
	pbeKeySpecInit(params)

	if specObj.FieldTable["password"].Fvalue != passwordCharsObj {
		t.Errorf("Expected password object to be stored")
	}

	// Test getPassword
	res := pbeKeySpecGetPassword([]any{specObj})
	if res != passwordCharsObj {
		t.Errorf("getPassword failed")
	}

	// Test clearPassword
	pbeKeySpecClearPassword([]any{specObj})

	// Verify field is null
	if !object.IsNull(specObj.FieldTable["password"].Fvalue) {
		t.Errorf("clearPassword should set password field to null")
	}

	// Verify original array is zeroed
	for i, val := range passwordChars {
		if val != 0 {
			t.Errorf("Password char at index %d should be 0, got %d", i, val)
		}
	}
}
