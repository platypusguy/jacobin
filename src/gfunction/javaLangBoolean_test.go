/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"reflect"
	"testing"

	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestLoad_Lang_Boolean_RegistersMethods(t *testing.T) {
	saved := MethodSignatures
	defer func() { MethodSignatures = saved }()
	MethodSignatures = make(map[string]GMeth)

	Load_Lang_Boolean()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/lang/Boolean.<clinit>()V", 0, booleanClinit},
		{"java/lang/Boolean.<init>(Z)V", 1, trapDeprecated},
		{"java/lang/Boolean.<init>(Ljava/lang/String;)V", 1, trapDeprecated},
		{"java/lang/Boolean.booleanValue()Z", 0, booleanBooleanValue},
		{"java/lang/Boolean.compare(ZZ)I", 2, booleanCompare},
		{"java/lang/Boolean.compareTo(Ljava/lang/Boolean;)I", 1, booleanCompareTo},
		{"java/lang/Boolean.describeConstable()Ljava.util.Optional;", 0, trapFunction},
		{"java/lang/Boolean.equals(Ljava/lang/Object;)Z", 1, booleanEquals},
		{"java/lang/Boolean.getBoolean(Ljava/lang/String;)Z", 1, booleanGetBoolean},
		{"java/lang/Boolean.hashCode()I", 0, booleanHashCode},
		{"java/lang/Boolean.hashCode(Z)I", 1, booleanHashCodeStatic},
		{"java/lang/Boolean.logicalAnd(ZZ)Z", 2, booleanLogicalAnd},
		{"java/lang/Boolean.logicalOr(ZZ)Z", 2, booleanLogicalOr},
		{"java/lang/Boolean.logicalXor(ZZ)Z", 2, booleanLogicalXor},
		{"java/lang/Boolean.parseBoolean(Ljava/lang/String;)Z", 1, booleanParseBoolean},
		{"java/lang/Boolean.toString()Ljava/lang/String;", 0, booleanToString},
		{"java/lang/Boolean.toString(Z)Ljava/lang/String;", 1, booleanToStringStatic},
		{"java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;", 1, booleanValueOf},
		{"java/lang/Boolean.valueOf(Ljava/lang/String;)Ljava/lang/Boolean;", 1, booleanValueOfString},
	}

	for _, c := range checks {
		got, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
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

func TestBooleanBooleanValue_ValidAndInvalid(t *testing.T) {
	globals.InitGlobals("test")

	btrue := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolTrue)
	ret := booleanBooleanValue([]interface{}{btrue})
	if v, ok := ret.(int64); !ok || v != types.JavaBoolTrue {
		t.Fatalf("booleanBooleanValue(true) got %v (%T)", ret, ret)
	}

	// Wrong field type -> error
	bad := object.MakeEmptyObject()
	bad.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: float64(1)}
	er := booleanBooleanValue([]interface{}{bad})
	if _, ok := er.(*GErrBlk); !ok {
		t.Fatalf("expected *GErrBlk for wrong field type, got %T", er)
	}
}

func TestBooleanValueOf_FromBooleanAndString(t *testing.T) {
	globals.InitGlobals("test")

	// From boolean primitive
	ret := booleanValueOf([]interface{}{types.JavaBoolFalse})
	obj := ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolFalse {
		t.Fatalf("valueOf(false) wrong value: %v", obj.FieldTable["value"].Fvalue)
	}

	// From string "true"
	sTrue := object.StringObjectFromGoString("true")
	ret = booleanValueOfString([]interface{}{sTrue})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolTrue {
		t.Fatalf("valueOf(\"true\") wrong value")
	}

	// From string "TRUE" (case-insensitive)
	sTrueUpper := object.StringObjectFromGoString("TRUE")
	ret = booleanValueOfString([]interface{}{sTrueUpper})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolTrue {
		t.Fatalf("valueOf(\"TRUE\") wrong value")
	}

	// From string "false"
	sFalse := object.StringObjectFromGoString("false")
	ret = booleanValueOfString([]interface{}{sFalse})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolFalse {
		t.Fatalf("valueOf(\"false\") wrong value")
	}

	// From invalid string -> should be false in Java, but valueOf(String) might be different?
	// Actually Boolean.valueOf(String) returns Boolean.TRUE if string is "true" (ignore case), else FALSE.
	sinv := object.StringObjectFromGoString("maybe")
	ret = booleanValueOfString([]interface{}{sinv})
	obj = ret.(*object.Object)
	if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolFalse {
		t.Fatalf("valueOf(\"maybe\") should be FALSE")
	}
}

func TestBooleanParseBoolean(t *testing.T) {
	globals.InitGlobals("test")

	sTrue := object.StringObjectFromGoString("true")
	if v := booleanParseBoolean([]interface{}{sTrue}); v.(int64) != types.JavaBoolTrue {
		t.Fatalf("parseBoolean(true) wrong")
	}
	sTrueUpper := object.StringObjectFromGoString("TrUe")
	if v := booleanParseBoolean([]interface{}{sTrueUpper}); v.(int64) != types.JavaBoolTrue {
		t.Fatalf("parseBoolean(TrUe) wrong")
	}
	sFalse := object.StringObjectFromGoString("false")
	if v := booleanParseBoolean([]interface{}{sFalse}); v.(int64) != types.JavaBoolFalse {
		t.Fatalf("parseBoolean(false) wrong")
	}
	sinv := object.StringObjectFromGoString("nope")
	if v := booleanParseBoolean([]interface{}{sinv}); v.(int64) != types.JavaBoolFalse {
		t.Fatalf("parseBoolean(invalid) should be false")
	}
}

func TestBooleanGetBoolean_UsesSystemProperties(t *testing.T) {
	globals.InitGlobals("test")
	globals.SetSystemProperty("my.flag", "true")
	globals.SetSystemProperty("another.flag", "false")
	globals.SetSystemProperty("upper.flag", "TRUE")
	globals.SetSystemProperty("bad.flag", "yes")

	s := object.StringObjectFromGoString("my.flag")
	if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolTrue {
		t.Fatalf("getBoolean(true) wrong")
	}

	s = object.StringObjectFromGoString("upper.flag")
	if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolTrue {
		t.Fatalf("getBoolean(TRUE) wrong")
	}

	s = object.StringObjectFromGoString("another.flag")
	if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolFalse {
		t.Fatalf("getBoolean(false) wrong")
	}

	s = object.StringObjectFromGoString("bad.flag")
	if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolFalse {
		t.Fatalf("getBoolean(bad) should be false")
	}

	s = object.StringObjectFromGoString("nonexistent")
	if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolFalse {
		t.Fatalf("getBoolean(nonexistent) should be false")
	}
}

func TestBooleanHashCode_And_Compare(t *testing.T) {
	globals.InitGlobals("test")

	// hashCode
	if v := booleanHashCodeStatic([]interface{}{types.JavaBoolTrue}); v.(int64) != 1231 {
		t.Fatalf("hashCode(true) expected 1231, got %v", v)
	}
	if v := booleanHashCodeStatic([]interface{}{types.JavaBoolFalse}); v.(int64) != 1237 {
		t.Fatalf("hashCode(false) expected 1237, got %v", v)
	}

	bTrue := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolTrue)
	if v := booleanHashCode([]interface{}{bTrue}); v.(int64) != 1231 {
		t.Fatalf("hashCode(obj true) expected 1231")
	}

	// compare
	if v := booleanCompare([]interface{}{types.JavaBoolTrue, types.JavaBoolTrue}).(int64); v != 0 {
		t.Fail()
	}
	if v := booleanCompare([]interface{}{types.JavaBoolTrue, types.JavaBoolFalse}).(int64); v <= 0 {
		t.Fail()
	}
	if v := booleanCompare([]interface{}{types.JavaBoolFalse, types.JavaBoolTrue}).(int64); v >= 0 {
		t.Fail()
	}

	// compareTo
	bFalse := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolFalse)
	if v := booleanCompareTo([]interface{}{bTrue, bFalse}).(int64); v <= 0 {
		t.Fail()
	}
}

func TestBooleanLogical_And_ToString_Equals(t *testing.T) {
	globals.InitGlobals("test")

	// Logical
	if v := booleanLogicalAnd([]interface{}{types.JavaBoolTrue, types.JavaBoolFalse}); v.(int64) != types.JavaBoolFalse {
		t.Fail()
	}
	if v := booleanLogicalOr([]interface{}{types.JavaBoolTrue, types.JavaBoolFalse}); v.(int64) != types.JavaBoolTrue {
		t.Fail()
	}
	if v := booleanLogicalXor([]interface{}{types.JavaBoolTrue, types.JavaBoolTrue}); v.(int64) != types.JavaBoolFalse {
		t.Fail()
	}

	// toString
	sObj := booleanToStringStatic([]interface{}{types.JavaBoolTrue}).(*object.Object)
	if object.GoStringFromStringObject(sObj) != "true" {
		t.Fail()
	}

	bFalse := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolFalse)
	sObj2 := booleanToString([]interface{}{bFalse}).(*object.Object)
	if object.GoStringFromStringObject(sObj2) != "false" {
		t.Fail()
	}

	// equals
	bFalse2 := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolFalse)
	if v := booleanEquals([]interface{}{bFalse, bFalse2}); v != types.JavaBoolTrue {
		t.Fail()
	}
	bTrue := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolTrue)
	if v := booleanEquals([]interface{}{bFalse, bTrue}); v != types.JavaBoolFalse {
		t.Fail()
	}
}
