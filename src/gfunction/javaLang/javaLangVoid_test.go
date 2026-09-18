/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/classloader"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestLoad_Lang_Void_RegistersMethod(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Void()

	gm, ok := ghelpers.MethodSignatures["java/lang/Void.<clinit>()V"]
	if !ok {
		t.Fatalf("expected java/lang/Void.<clinit>()V to be registered")
	}
	if gm.ParamSlots != 0 {
		t.Errorf("expected ParamSlots=0, got %d", gm.ParamSlots)
	}
	if gm.GFunction == nil {
		t.Errorf("expected a non-nil GFunction")
	}
}

func TestVoidClinit_Success(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	statics.Statics = make(map[string]statics.Static)

	result := voidClinit(nil)
	if result != nil {
		t.Errorf("expected nil return, got %v", result)
	}

	st := statics.GetStaticValue("java/lang/Void", "TYPE")
	if st == nil {
		t.Fatalf("expected java/lang/Void.TYPE static to be set")
	}

	k := classloader.MethAreaFetch("void")
	if k == nil || k.Data.ClassObject == nil {
		t.Fatalf("test setup failure: 'void' class not properly preloaded")
	}
	if st != k.Data.ClassObject {
		t.Errorf("expected java/lang/Void.TYPE to point to the 'void' class object")
	}
}

func TestVoidClinit_TypeFieldHasRefType(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	statics.Statics = make(map[string]statics.Static)

	voidClinit(nil)

	stat, ok := statics.Statics["java/lang/Void.TYPE"]
	if !ok {
		t.Fatalf("expected java/lang/Void.TYPE entry to exist")
	}
	if stat.Type != types.Ref {
		t.Errorf("expected Type=%s, got %s", types.Ref, stat.Type)
	}
}

func TestVoidClinit_MethAreaMissing(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	statics.Statics = make(map[string]statics.Static)

	// Remove the "void" entry from the method area to force the error branch.
	classloader.MethArea.Delete("void")

	result := voidClinit(nil)
	if result != nil {
		t.Errorf("expected nil return on error branch, got %v", result)
	}

	if _, ok := statics.Statics["java/lang/Void.TYPE"]; ok {
		t.Errorf("expected java/lang/Void.TYPE to remain unset")
	}
}

func TestVoidClinit_MethAreaClassObjectNil(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	statics.Statics = make(map[string]statics.Static)

	k := classloader.MethAreaFetch("void")
	if k == nil {
		t.Fatalf("test setup failure: 'void' class not preloaded")
	}
	k.Data.ClassObject = nil

	result := voidClinit(nil)
	if result != nil {
		t.Errorf("expected nil return on error branch, got %v", result)
	}

	if _, ok := statics.Statics["java/lang/Void.TYPE"]; ok {
		t.Errorf("expected java/lang/Void.TYPE to remain unset")
	}
}
