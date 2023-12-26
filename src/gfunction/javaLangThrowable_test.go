/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/statics"
	"testing"
)

func TestJavaLangThrowableClinit(t *testing.T) {
	statics.Statics = make(map[string]statics.Static)

	throwableClinit(nil)
	_, ok := statics.Statics["Throwable.UNASSIGNED_STACK"]
	if !ok {
		t.Error("JavaLangThrowableClinit: Throwable.UNASSIGNED_STACK not found")
	}

	_, ok = statics.Statics["Throwable.SUPPRESSED_SENTINEL"]
	if !ok {
		t.Error("JavaLangThrowableClinit: Throwable.SUPPRESSED_SENTINEL not found")
	}

	_, ok = statics.Statics["Throwable.EMPTY_THROWABLE_ARRAY"]
	if !ok {
		t.Error("Throwable.EMPTY_THROWABLE_ARRAY not found")
	}

}
