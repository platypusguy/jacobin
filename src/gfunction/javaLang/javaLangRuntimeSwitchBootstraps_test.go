/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Lang_Runtime_SwitchBootstraps_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Runtime_SwitchBootstraps()

	if len(ghelpers.MethodSignatures) == 0 {
		t.Fatalf("expected method signatures to be registered")
	}

	clinit, ok := ghelpers.MethodSignatures["java/lang/runtime/SwitchBootstraps.<clinit>()V"]
	if !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if clinit.ParamSlots != 0 {
		t.Errorf("expected <clinit> ParamSlots 0, got %d", clinit.ParamSlots)
	}

	for key, gm := range ghelpers.MethodSignatures {
		if key == "java/lang/runtime/SwitchBootstraps.<clinit>()V" {
			continue
		}
		if gm.GFunction == nil {
			t.Errorf("entry %q has nil GFunction", key)
		}
		if gm.ParamSlots != 4 {
			t.Errorf("entry %q expected ParamSlots 4, got %d", key, gm.ParamSlots)
		}
	}

	if _, ok := ghelpers.MethodSignatures["java/lang/runtime/SwitchBootstraps.typeSwitch(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"]; !ok {
		t.Fatalf("expected typeSwitch() to be registered")
	}

	if _, ok := ghelpers.MethodSignatures["java/lang/runtime/SwitchBootstraps.enumSwitch(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"]; !ok {
		t.Fatalf("expected enumSwitch() to be registered")
	}
}
