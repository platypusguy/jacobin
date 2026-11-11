/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"strings"
	"testing"
)

func loadAllJavaMath() {
	globals.InitGlobals("test")
	Load_Math_Big_Decimal()
	Load_Math_Big_Integer()
	Load_Math_Math_Context()
	Load_Math_Rounding_Mode()
}

const classPrefix = "java/math/"

func TestJavaMathParamSlotsConsistency(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	loadAllJavaMath()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, classPrefix) {
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
	t.Logf("Total %s signatures checked: %d, mismatches: %d", classPrefix, len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
