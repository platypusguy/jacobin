/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaMath

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"reflect"
	"testing"
)

// resetRMode clears lazy-init state to ensure tests don't influence each other.
func resetRMode() {
	rmodeOnceInitialized = false
	rmodeInstances = nil
}

func TestLoad_Math_Rounding_Mode_RegistersCoreMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Math_Rounding_Mode()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/math/RoundingMode.valueOf(I)Ljava/math/RoundingMode;", 1, rmodeValueOfInt},
		{"java/math/RoundingMode.valueOf(Ljava/lang/String;)Ljava/math/RoundingMode;", 1, rmodeValueOfString},
		{"java/math/RoundingMode.values()[Ljava/math/RoundingMode;", 0, rmodeValues},
	}

	for _, c := range checks {
		got, ok := ghelpers.MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing ghelpers.MethodSignatures entry for %s", c.key)
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

func TestRoundingMode_values_ReturnsEightConstantsInOrder(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()

	ret := rmodeValues(nil)
	arr, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("values() expected *object.Object array, got %T", ret)
	}
	// Verify array type and length
	if arr.FieldTable["value"].Ftype != types.RefArray+"Ljava/math/RoundingMode;" {
		t.Fatalf("values() wrong array ftype: %s", arr.FieldTable["value"].Ftype)
	}
	slot := arr.FieldTable["value"].Fvalue.([]*object.Object)
	if len(slot) != 8 {
		t.Fatalf("values() expected 8 elements, got %d", len(slot))
	}

	for i, obj := range slot {
		if obj == nil {
			t.Fatalf("values()[%d] is nil", i)
		}
		// Check class name
		if cn := object.GoStringFromStringPoolIndex(obj.KlassName); cn != rmodeClassName {
			t.Fatalf("values()[%d] wrong class: %s", i, cn)
		}
		// Check ordinal and name fields
		ordFld, ok := obj.FieldTable["ordinal"]
		if !ok || ordFld.Ftype != types.Int {
			t.Fatalf("values()[%d] missing/invalid ordinal field", i)
		}
		if ordFld.Fvalue.(int64) != int64(i) {
			t.Fatalf("values()[%d] ordinal expected %d, got %v", i, i, ordFld.Fvalue)
		}
		nameFld, ok := obj.FieldTable["name"]
		if !ok || nameFld.Ftype != types.StringClassRef {
			t.Fatalf("values()[%d] missing/invalid name field", i)
		}
		if nm := object.GoStringFromStringObject(nameFld.Fvalue.(*object.Object)); nm != rmodeNames[i] {
			t.Fatalf("values()[%d] name expected %s, got %s", i, rmodeNames[i], nm)
		}
	}
}

func TestRoundingMode_valueOfInt_MappingsAndErrors(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()

	// Valid mappings 0..7
	for i := 0; i < len(rmodeNames); i++ {
		ret := rmodeValueOfInt([]interface{}{int64(i)})
		obj, ok := ret.(*object.Object)
		if !ok {
			t.Fatalf("valueOf(int %d) expected *object.Object, got %T", i, ret)
		}
		// Check ordinal and name
		if obj.FieldTable["ordinal"].Fvalue.(int64) != int64(i) {
			t.Fatalf("valueOf(%d) wrong ordinal: %v", i, obj.FieldTable["ordinal"].Fvalue)
		}
		nm := object.GoStringFromStringObject(obj.FieldTable["name"].Fvalue.(*object.Object))
		if nm != rmodeNames[i] {
			t.Fatalf("valueOf(%d) wrong name: %s", i, nm)
		}
	}

	// Invalid code (negative)
	if blk, ok := rmodeValueOfInt([]interface{}{int64(-1)}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("valueOf(-1) expected IAE, got %T", blk)
	}
	// Invalid code (too large)
	if blk, ok := rmodeValueOfInt([]interface{}{int64(8)}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("valueOf(8) expected IAE, got %T", blk)
	}
	// Wrong type
	if blk, ok := rmodeValueOfInt([]interface{}{object.StringObjectFromGoString("1")}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("valueOf(wrong type) expected IAE, got %T", blk)
	}
}

func TestRoundingMode_valueOfString_MappingsAndErrors(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()

	for i, nm := range rmodeNames {
		s := object.StringObjectFromGoString(nm)
		ret := rmodeValueOfString([]interface{}{s})
		obj, ok := ret.(*object.Object)
		if !ok {
			t.Fatalf("valueOf(String %s) expected *object.Object, got %T", nm, ret)
		}
		// Confirm it's the right one by ordinal and name
		if obj.FieldTable["ordinal"].Fvalue.(int64) != int64(i) {
			t.Fatalf("valueOf(%s) wrong ordinal: %v", nm, obj.FieldTable["ordinal"].Fvalue)
		}
		gotName := object.GoStringFromStringObject(obj.FieldTable["name"].Fvalue.(*object.Object))
		if gotName != nm {
			t.Fatalf("valueOf(%s) wrong name: %s", nm, gotName)
		}
	}

	// Null argument -> NPE
	if blk, ok := rmodeValueOfString([]interface{}{object.Null}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("valueOfString(null) expected NPE, got %T", blk)
	}
	// Wrong type -> IAE
	if blk, ok := rmodeValueOfString([]interface{}{object.MakeEmptyObject()}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("valueOfString(wrong type) expected IAE, got %T", blk)
	}
	// Invalid name -> IAE
	if blk, ok := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("ROUND_UP")}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("valueOfString(invalid) expected IAE, got %T", blk)
	}
}

func TestRoundingMode_StaticConstants_NonNilAndCorrect(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()
	// Trigger init
	_ = rmodeValues(nil)

	// HALF_UP
	hu := statics.GetStaticValue("java/math/RoundingMode", "HALF_UP")
	if hu == nil {
		t.Fatalf("static RoundingMode.HALF_UP is nil")
	}
	obj, ok := hu.(*object.Object)
	if !ok {
		t.Fatalf("static HALF_UP not *object.Object; got %T", hu)
	}
	ord, ok2 := obj.FieldTable["ordinal"]
	if !ok2 || ord.Fvalue.(int64) != 4 {
		t.Fatalf("static HALF_UP wrong ordinal; got %v", ord)
	}
	nm := object.GoStringFromStringObject(obj.FieldTable["name"].Fvalue.(*object.Object))
	if nm != "HALF_UP" {
		t.Fatalf("static HALF_UP wrong name: %s", nm)
	}

	// All constants are present and distinct
	for i, name := range rmodeNames {
		v := statics.GetStaticValue("java/math/RoundingMode", name)
		if v == nil {
			t.Fatalf("static %s is nil", name)
		}
		o := v.(*object.Object)
		if o.FieldTable["ordinal"].Fvalue.(int64) != int64(i) {
			t.Fatalf("%s ordinal mismatch", name)
		}
	}
}

func TestRoundingMode_name_ReturnsEnumName(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()
	// Obtain HALF_UP instance
	ret := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("HALF_UP")})
	if blk, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("failed to get HALF_UP: %v", blk)
	}
	rm := ret.(*object.Object)
	res := rmodeName([]interface{}{rm})
	s, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("name() expected *object.Object (String), got %T", res)
	}
	if got := object.GoStringFromStringObject(s); got != "HALF_UP" {
		t.Fatalf("name() expected HALF_UP, got %s", got)
	}
}

func TestRoundingMode_name_NullReceiverThrowsNPE(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()
	res := rmodeName([]interface{}{object.Null})
	if blk, ok := res.(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.NullPointerException {
		t.Fatalf("name(null) expected NPE, got %T", res)
	}
}

func TestRoundingMode_equals_Behavior(t *testing.T) {
	globals.InitGlobals("test")
	resetRMode()
	// Prepare a couple of constants
	up := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("UP")}).(*object.Object)
	down := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("DOWN")}).(*object.Object)

	// Reflexive and same-constant
	if ret := rmodeEquals([]interface{}{up, up}); ret != types.JavaBoolTrue {
		t.Fatalf("UP.equals(UP) expected true, got %v", ret)
	}
	up2 := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("UP")}).(*object.Object)
	if ret := rmodeEquals([]interface{}{up, up2}); ret != types.JavaBoolTrue {
		t.Fatalf("UP.equals(UP2) expected true, got %v", ret)
	}

	// Different constant
	if ret := rmodeEquals([]interface{}{up, down}); ret != types.JavaBoolFalse {
		t.Fatalf("UP.equals(DOWN) expected false, got %v", ret)
	}

	// Null argument -> false
	if ret := rmodeEquals([]interface{}{up, object.Null}); ret != types.JavaBoolFalse {
		t.Fatalf("UP.equals(null) expected false, got %v", ret)
	}

	// Different type -> false
	if ret := rmodeEquals([]interface{}{up, object.MakeEmptyObject()}); ret != types.JavaBoolFalse {
		t.Fatalf("UP.equals(new Object()) expected false, got %v", ret)
	}
}
