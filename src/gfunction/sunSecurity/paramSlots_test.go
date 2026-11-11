/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"testing"
)

func loadAllSunSecurity() {
	globals.InitGlobals("test")
	Load_Sun_Security_Action_GetBooleanAction()
	Load_Sun_Security_Action_GetIntegerAction()
	Load_Sun_Security_Action_GetLongAction()
	Load_Sun_Security_Action_GetPropertyAction()
	Load_Sun_Security_Jca_ProviderList()
	Load_Sun_Security_Jca_Providers()
	Load_Sun_Security_Util_Stubs()
}

func TestSunSecurityParamSlotsConsistency(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	loadAllSunSecurity()

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
	t.Logf("Total sun/security signatures checked: %d, mismatches: %d", len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
