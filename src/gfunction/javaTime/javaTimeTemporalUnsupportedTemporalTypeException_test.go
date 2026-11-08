/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Time_Temporal_UnsupportedTemporalTypeException(t *testing.T) {
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Time_Temporal_UnsupportedTemporalTypeException()

	expectedSignatures := []string{
		"java/time/temporal/UnsupportedTemporalTypeException.<init>(Ljava/lang/String;)V",
		"java/time/temporal/UnsupportedTemporalTypeException.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V",
	}

	for _, signature := range expectedSignatures {
		gMethod, exists := ghelpers.MethodSignatures[signature]
		if !exists {
			t.Errorf("Expected method signature %q to be registered", signature)
			continue
		}
		if gMethod.GFunction == nil {
			t.Errorf("Method signature %q has a nil GFunction", signature)
		}
		if !gMethod.NeedsContext {
			t.Errorf("Method signature %q should have NeedsContext=true", signature)
		}
	}
}
