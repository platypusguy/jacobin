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

func TestLoad_Lang_Byte_RegistersMethods(t *testing.T) {
    saved := MethodSignatures
    defer func() { MethodSignatures = saved }()
    MethodSignatures = make(map[string]GMeth)

    Load_Lang_Byte()

    checks := []struct{
        key   string
        slots int
        fn    func([]interface{}) interface{}
    }{
        {"java/lang/Byte.<clinit>()V", 0, clinitGeneric},
        {"java/lang/Byte.decode(Ljava/lang/String;)Ljava/lang/Byte;", 1, byteDecode},
        {"java/lang/Byte.doubleValue()D", 0, byteDoubleValue},
        {"java/lang/Byte.toString()Ljava/lang/String;", 0, byteToString},
        {"java/lang/Byte.valueOf(B)Ljava/lang/Byte;", 1, byteValueOf},
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

func TestByteDecode_Various(t *testing.T) {
    globals.InitGlobals("test")

    // valid with leading #
    s := object.StringObjectFromGoString("#0a")
    ret := byteDecode([]interface{}{s})
    obj := ret.(*object.Object)
    if obj.FieldTable["value"].Fvalue.(int64) != 10 {
        t.Fatalf("decode #0a expected 10, got %v", obj.FieldTable["value"].Fvalue)
    }

    // valid with 0x
    s = object.StringObjectFromGoString("0x2f")
    ret = byteDecode([]interface{}{s})
    obj = ret.(*object.Object)
    if obj.FieldTable["value"].Fvalue.(int64) != 47 {
        t.Fatalf("decode 0x2f expected 47, got %v", obj.FieldTable["value"].Fvalue)
    }

    // too large
    s = object.StringObjectFromGoString("1ff") // 511
    if blk, ok := byteDecode([]interface{}{s}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
        t.Fatalf("decode too-large expected NFE, got %T", blk)
    }

    // invalid hex
    s = object.StringObjectFromGoString("zz")
    if blk, ok := byteDecode([]interface{}{s}).(*GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
        t.Fatalf("decode invalid hex expected NFE, got %T", blk)
    }
}

func TestByteDoubleValue_ToString_ValueOf(t *testing.T) {
    globals.InitGlobals("test")

    b := Populator("java/lang/Byte", types.Byte, int64(127))

    // doubleValue
    if v := byteDoubleValue([]interface{}{b}); v.(float64) != float64(127) {
        t.Fatalf("doubleValue wrong")
    }

    // toString
    s := byteToString([]interface{}{b}).(*object.Object)
    if str := object.GoStringFromStringObject(s); str != "127" {
        t.Fatalf("toString wrong: %q", str)
    }

    // valueOf
    vobj := byteValueOf([]interface{}{int64(5)}).(*object.Object)
    if v := vobj.FieldTable["value"].Fvalue.(int64); v != 5 {
        t.Fatalf("valueOf 5 wrong: %v", v)
    }
}
