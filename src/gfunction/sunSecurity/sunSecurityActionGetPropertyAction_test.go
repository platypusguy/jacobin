/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"testing"
)

func TestLoad_Sun_Security_Action_GetPropertyAction(t *testing.T) {
	globals.InitGlobals("test")

	// Clear any existing signatures to ensure clean test
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	// Call the loader
	Load_Sun_Security_Action_GetPropertyAction()

	expectedMethods := []struct {
		sig   string
		slots int
	}{
		{"sun/security/action/GetPropertyAction.privilegedGetProperties()Ljava/util/Properties;", 0},
		{"sun/security/action/GetPropertyAction.privilegedGetProperty(Ljava/lang/String;)Ljava/lang/String;", 1},
		{"sun/security/action/GetPropertyAction.privilegedGetProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;", 2},
		{"sun/security/action/GetPropertyAction.privilegedGetTimeoutProp(Ljava/lang/String;ILsun/security/util/Debug;)I", 3},
	}

	for _, m := range expectedMethods {
		gm, ok := ghelpers.MethodSignatures[m.sig]
		if !ok {
			t.Errorf("Method signature not registered: %s", m.sig)
			continue
		}
		if gm.ParamSlots != m.slots {
			t.Errorf("%s: expected %d param slots, got %d", m.sig, m.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Errorf("%s: GFunction is nil", m.sig)
		}
	}
}
