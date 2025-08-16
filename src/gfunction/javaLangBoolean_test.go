/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
    "reflect"
    "testing"

    "jacobin/excNames"
    "jacobin/globals"
    "jacobin/object"
    "jacobin/types"
)

func TestLoad_Lang_Boolean_RegistersMethods(t *testing.T) {
    saved := MethodSignatures
    defer func() { MethodSignatures = saved }()
    MethodSignatures = make(map[string]GMeth)

    Load_Lang_Boolean()

    checks := []struct{
        key   string
        slots int
        fn    func([]interface{}) interface{}
    }{
        {"java/lang/Boolean.<clinit>()V", 0, clinitGeneric},
        {"java/lang/Boolean.<init>(Z)V", 1, trapDeprecated},
        {"java/lang/Boolean.<init>(Ljava/lang/String;)V", 1, trapDeprecated},
        {"java/lang/Boolean.booleanValue()Z", 0, booleanBooleanValue},
        {"java/lang/Boolean.describeConstable()Ljava.util.Optional;", 0, trapFunction},
        {"java/lang/Boolean.getBoolean(Ljava/lang/String;)Z", 1, booleanGetBoolean},
        {"java/lang/Boolean.hashCode()I", 0, booleanHashCode},
        {"java/lang/Boolean.hashCode(Z)I", 1, booleanHashCode},
        {"java/lang/Boolean.parseBoolean(Ljava/lang/String;)Z", 1, booleanParseBoolean},
        {"java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;", 1, booleanValueOf},
        {"java/lang/Boolean.valueOf(Ljava/lang/String;)Ljava/lang/Boolean;", 1, booleanValueOf},
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

    btrue := Populator("java/lang/Boolean", types.Bool, types.JavaBoolTrue)
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
    ret = booleanValueOf([]interface{}{sTrue})
    obj = ret.(*object.Object)
    if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolTrue {
        t.Fatalf("valueOf(\"true\") wrong value")
    }

    // From string "false"
    sFalse := object.StringObjectFromGoString("false")
    ret = booleanValueOf([]interface{}{sFalse})
    obj = ret.(*object.Object)
    if obj.FieldTable["value"].Fvalue.(int64) != types.JavaBoolFalse {
        t.Fatalf("valueOf(\"false\") wrong value")
    }

    // From invalid string -> error block
    sinv := object.StringObjectFromGoString("maybe")
    er := booleanValueOf([]interface{}{sinv})
    if _, ok := er.(*GErrBlk); !ok {
        t.Fatalf("expected *GErrBlk for invalid string valueOf, got %T", er)
    }
}

func TestBooleanParseBoolean(t *testing.T) {
    globals.InitGlobals("test")

    sTrue := object.StringObjectFromGoString("true")
    if v := booleanParseBoolean([]interface{}{sTrue}); v.(int64) != types.JavaBoolTrue {
        t.Fatalf("parseBoolean(true) wrong")
    }
    sFalse := object.StringObjectFromGoString("false")
    if v := booleanParseBoolean([]interface{}{sFalse}); v.(int64) != types.JavaBoolFalse {
        t.Fatalf("parseBoolean(false) wrong")
    }
    sinv := object.StringObjectFromGoString("nope")
    if blk, ok := booleanParseBoolean([]interface{}{sinv}).(*GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
        t.Fatalf("parseBoolean(invalid) expected IAE, got %T", blk)
    }
}

func TestBooleanGetBoolean_UsesSystemProperties(t *testing.T) {
    globals.InitGlobals("test")
    globals.SetSystemProperty("my.flag", "true")
    globals.SetSystemProperty("another.flag", "false")
    globals.SetSystemProperty("bad.flag", "yes")

    s := object.StringObjectFromGoString("my.flag")
    if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolTrue {
        t.Fatalf("getBoolean(true) wrong")
    }

    s = object.StringObjectFromGoString("another.flag")
    if v := booleanGetBoolean([]interface{}{s}); v.(int64) != types.JavaBoolFalse {
        t.Fatalf("getBoolean(false) wrong")
    }

    s = object.StringObjectFromGoString("bad.flag")
    if blk, ok := booleanGetBoolean([]interface{}{s}).(*GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
        t.Fatalf("getBoolean(bad) expected IAE, got %T", blk)
    }
}

func TestBooleanHashCode_FromPrimitiveInput(t *testing.T) {
    globals.InitGlobals("test")
    // booleanHashCode expects params[1] for the boolean case; pass a dummy first arg.
    if v := booleanHashCode([]interface{}{nil, types.JavaBoolTrue}); v.(int64) != 1231 {
        t.Fatalf("hashCode(true) expected 1231, got %v", v)
    }
    if v := booleanHashCode([]interface{}{nil, types.JavaBoolFalse}); v.(int64) != 1237 {
        t.Fatalf("hashCode(false) expected 1237, got %v", v)
    }
    // invalid value -> error
    if blk, ok := booleanHashCode([]interface{}{nil, int64(2)}).(*GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
        t.Fatalf("hashCode(invalid) expected IAE, got %T", blk)
    }
}
