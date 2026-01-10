/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package misc

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"testing"
)

func TestLoad_Awt_Graphics_Environment_RegistersMethods(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()

	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Awt_Graphics_Environment()

	// Expected keys
	kClinit := "java/awt/GraphicsEnvironment.<clinit>()V"
	kIsHeadless := "java/awt/GraphicsEnvironment.isHeadless()Z"
	kIsHeadlessInstance := "java/awt/GraphicsEnvironment.isHeadlessInstance()Z"

	// Check presence
	if _, ok := ghelpers.MethodSignatures[kClinit]; !ok {
		t.Fatalf("missing ghelpers.MethodSignatures entry for %s", kClinit)
	}
	if _, ok := ghelpers.MethodSignatures[kIsHeadless]; !ok {
		t.Fatalf("missing ghelpers.MethodSignatures entry for %s", kIsHeadless)
	}
	if _, ok := ghelpers.MethodSignatures[kIsHeadlessInstance]; !ok {
		t.Fatalf("missing ghelpers.MethodSignatures entry for %s", kIsHeadlessInstance)
	}

	// Validate isHeadless entries have zero params and non-nil functions
	if m := ghelpers.MethodSignatures[kIsHeadless]; m.ParamSlots != 0 {
		t.Fatalf("isHeadless ParamSlots expected 0, got %d", m.ParamSlots)
	} else if m.GFunction == nil {
		t.Fatalf("isHeadless GFunction expected non-nil")
	}

	if m := ghelpers.MethodSignatures[kIsHeadlessInstance]; m.ParamSlots != 0 {
		t.Fatalf("isHeadlessInstance ParamSlots expected 0, got %d", m.ParamSlots)
	} else if m.GFunction == nil {
		t.Fatalf("isHeadlessInstance GFunction expected non-nil")
	}
}

func TestAwtgeIsHeadless_ReflectsGlobals(t *testing.T) {
	globals.InitGlobals("test")
	glob := globals.GetGlobalRef()

	// true case
	glob.Headless = true
	if v, ok := awtgeIsHeadless(nil).(bool); !ok {
		t.Fatalf("awtgeIsHeadless did not return bool when Headless=true, got %T", awtgeIsHeadless(nil))
	} else if !v {
		t.Fatalf("awtgeIsHeadless expected true when globals.Headless=true")
	}

	// false case
	glob.Headless = false
	if v, ok := awtgeIsHeadless(nil).(bool); !ok {
		t.Fatalf("awtgeIsHeadless did not return bool when Headless=false, got %T", awtgeIsHeadless(nil))
	} else if v {
		t.Fatalf("awtgeIsHeadless expected false when globals.Headless=false")
	}
}
