/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoadJavaUtilZipZipException_RegistersMethods(t *testing.T) {
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Util_Zip_ZipException()

	expected := []string{
		"java/util/zip/ZipException.<init>()V",
		"java/util/zip/ZipException.<init>(Ljava/lang/String;)V",
	}

	for _, sig := range expected {
		gmeth, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("expected method signature %q to be registered", sig)
			continue
		}
		if gmeth.GFunction == nil {
			t.Errorf("method signature %q has a nil GFunction", sig)
		}
		if !gmeth.NeedsContext {
			t.Errorf("method signature %q should need context", sig)
		}
	}
}
