/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/classloader"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"reflect"
	"testing"
)

func TestLoad_TestGfunctions_RegistersMethods(t *testing.T) {
	// Save and restore ghelpers.TestMethodSignatures to avoid cross-test pollution
	saved := ghelpers.TestMethodSignatures
	defer func() { ghelpers.TestMethodSignatures = saved }()
	ghelpers.TestMethodSignatures = make(map[string]ghelpers.GMeth)

	Load_TestGfunctions()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"jacobin/src/test/Object.test()D", 0, vd},
		{"jacobin/src/test/Object.test()Ljava/lang/Object;", 0, vl},
		{"jacobin/src/test/Object.test(I)V", 1, iv},
		{"jacobin/src/test/Object.test(D)V", 1, dv},
		{"jacobin/src/test/Object.test(Ljava/lang/Object;)V", 1, lv},
		{"jacobin/src/test/Object.test(I)I", 1, ii},
		{"jacobin/src/test/Object.test(I)D", 1, id},
		{"jacobin/src/test/Object.test(I)Ljava/lang/Object;", 1, il},
		{"jacobin/src/test/Object.test(Ljava/lang/Object;)I", 1, li},
		{"jacobin/src/test/Object.test(Ljava/lang/Object;)Ljava/lang/Object;", 1, ll},
		{"jacobin/src/test/Object.test(Ljava/lang/Object;)D", 1, ld},
		{"jacobin/src/test/Object.test(D)E", 1, ie},
	}

	for _, c := range checks {
		got, ok := ghelpers.TestMethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing ghelpers.TestMethodSignatures entry for %s", c.key)
		}
		if got.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
		}
		if got.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}

func TestCheckTestGfunctionsLoaded_PopulatesMTable(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea() // required by CheckTestGfunctionsLoaded -> MethAreaInsert

	// clean MTable and ghelpers.TestMethodSignatures
	classloader.MTable = make(map[string]classloader.MTentry)
	saved := ghelpers.TestMethodSignatures
	defer func() { ghelpers.TestMethodSignatures = saved }()
	ghelpers.TestMethodSignatures = make(map[string]ghelpers.GMeth)

	// Load test gfunctions into MTable
	CheckTestGfunctionsLoaded()

	// A couple of representative keys should be in the classloader MTable now
	keys := []string{
		"jacobin/src/test/Object.test()D",
		"jacobin/src/test/Object.test(I)I",
		"jacobin/src/test/Object.test(D)E",
	}

	for _, k := range keys {
		mte, ok := classloader.MTable[k]
		if !ok {
			t.Fatalf("MTable missing key %s after CheckTestGfunctionsLoaded", k)
		}
		if mte.MType != 'G' {
			t.Fatalf("MType for %s expected 'G', got %c", k, mte.MType)
		}
		// Sanity check: Meth must be a ghelpers.GMeth with a non-nil GFunction
		gm, ok := mte.Meth.(ghelpers.GMeth)
		if !ok {
			t.Fatalf("MTable entry %s not a ghelpers.GMeth", k)
		}
		if gm.GFunction == nil {
			t.Fatalf("MTable entry %s has nil GFunction", k)
		}
	}
}

func TestTestGfunctions_BasicBehaviors(t *testing.T) {
	// simple direct calls to confirm contracts
	if v, ok := vd(nil).(float64); !ok || v == 0 {
		t.Fatalf("vd expected non-zero float64, got %v (%T)", v, v)
	}
	if v := ii(nil); v != int64(43) {
		t.Fatalf("ii expected 43, got %v", v)
	}
	if v := id(nil); v != float64(43.43) {
		t.Fatalf("id expected 43.43, got %v", v)
	}
	if v := li(nil); v != 44 {
		t.Fatalf("li expected 44, got %v", v)
	}
	if v := ld(nil); v != float64(44.44) {
		t.Fatalf("ld expected 44.44, got %v", v)
	}

	// vl/il/ll return *object.Object pointing to bare java/lang/Object
	if o, ok := vl(nil).(*object.Object); !ok || o == nil || o.KlassName != types.StringPoolObjectIndex {
		t.Fatalf("vl returned wrong object")
	}
	if o, ok := il(nil).(*object.Object); !ok || o == nil || o.KlassName != types.StringPoolObjectIndex {
		t.Fatalf("il returned wrong object")
	}
	if o, ok := ll(nil).(*object.Object); !ok || o == nil || o.KlassName != types.StringPoolObjectIndex {
		t.Fatalf("ll returned wrong object")
	}
}
