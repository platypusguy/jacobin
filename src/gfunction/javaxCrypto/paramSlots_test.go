/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"testing"
)

func loadAllJavaxCrypto() {
	globals.InitGlobals("test")
	Load_Crypto_Cipher()
	Load_Crypto_CipherInputStream()
	Load_Crypto_ExemptionMechanism()
	Load_Crypto_Interfaces_DH_Keys()
	Load_Crypto_KEM()
	Load_Crypto_KeyAgreement()
	Load_Crypto_Mac()
	Load_Crypto_SecretKeyFactory()
	Load_Crypto_Spec_DHParameterSpec()
	Load_Crypto_Spec_GCMParameterSpec()
	Load_Crypto_Spec_IvParameterSpec()
	Load_Crypto_Spec_PBEKeySpec()
	Load_Crypto_Spec_PBEParameterSpec()
	Load_Crypto_Spec_SecretKeySpec()
}

func TestJavaxCryptoParamSlotsConsistency(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	loadAllJavaxCrypto()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		_, _, methType := util.ParseFQN(k)
		params := util.ParseIncomingParamsFromMethTypeString(methType)
		expectedSlots := len(params)
		if v.ParamSlots != expectedSlots {
			t.Errorf("Signature mismatch for %s: ParamSlots is %d, expected %d (params: %v)", k, v.ParamSlots, expectedSlots, params)
			mismatches = append(mismatches, k)
		}
	}
	t.Logf("Total javax/crypto signatures checked: %d, mismatches: %d", len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
