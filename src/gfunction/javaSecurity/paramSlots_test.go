/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"strings"
	"testing"
)

func loadAllJavaSecurity() {
	globals.InitGlobals("test")
	Load_ECFieldAndPoint()
	Load_Security_Interfaces_EC_Keys()
	Load_ECParameterSpec()
	Load_EllipticCurve()
	Load_KeyPairGenerator()
	Load_PublicAndPrivateKeys()
	Load_Security()
	Load_Security_Interfaces_DSA_Keys()
	Load_Security_Interfaces_RSA_Keys()
	Load_Security_Key()
	Load_Security_SecretKey()
	Load_Security_KeyPair()
	Load_Security_MessageDigest()
	Load_Security_Provider()
	Load_Security_Provider_Service()
	Load_Security_SecureRandom()
	Load_Security_Signature()
	Load_Security_AlgorithmParameters()
	Load_Security_Spec_NamedParameterSpec()
	Load_Security_Spec_AlgorithmParameterSpec()
}

func TestJavaSecurityParamSlotsConsistency(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	loadAllJavaSecurity()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, "java/security/") && !strings.HasPrefix(k, "javax/crypto/") {
			continue
		}
		_, _, methType := util.ParseFQN(k)
		params := util.ParseIncomingParamsFromMethTypeString(methType)
		expectedSlots := len(params)
		if v.ParamSlots != expectedSlots {
			t.Errorf("Signature mismatch for %s: ParamSlots is %d, expected %d (params: %v)", k, v.ParamSlots, expectedSlots, params)
			mismatches = append(mismatches, k)
		}
	}
	t.Logf("Total signatures checked: %d, mismatches: %d", len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
