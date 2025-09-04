/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
    "reflect"
    "testing"

    "jacobin/src/excNames"
    "jacobin/src/globals"
    "jacobin/src/object"
)

// assertGErrBlkMC: small helper for this test file to assert a GErrBlk with a specific exception type
func assertGErrBlkMC(t *testing.T, res interface{}, wantExc int) {
    t.Helper()
    if res == nil {
        t.Fatalf("expected *GErrBlk, got nil")
    }
    blk, ok := res.(*GErrBlk)
    if !ok {
        t.Fatalf("expected *GErrBlk, got %T", res)
    }
    if blk.ExceptionType != wantExc {
        t.Fatalf("expected exception %d, got %d; msg=%s", wantExc, blk.ExceptionType, blk.ErrMsg)
    }
}

func TestLoad_Math_Math_Context_RegistersMethods(t *testing.T) {
    saved := MethodSignatures
    defer func() { MethodSignatures = saved }()
    MethodSignatures = make(map[string]GMeth)

    Load_Math_Math_Context()

    checks := []struct{
        key   string
        slots int
        fn    func([]interface{}) interface{}
    }{
        {"java/math/MathContext.<clinit>()V", 0, clinitGeneric},
        {"java/math/MathContext.<init>(I)V", 1, mconInitInt},
        {"java/math/MathContext.<init>(ILjava/math/RoundingMode;)V", 2, mconInitIntRoundingMode},
        {"java/math/MathContext.<init>(Ljava/lang/String;)V", 1, mconInitString},
        {"java/math/MathContext.getPrecision()I", 0, mconGetPrecision},
        {"java/math/MathContext.getRoundingMode()Ljava/math/RoundingMode;", 0, mconGetRoundingMode},
        {"java/math/MathContext.toString()Ljava/lang/String;", 0, mconToString},
    }

    for _, c := range checks {
        got, ok := MethodSignatures[c.key]
        if !ok {
            t.Fatalf("missing MethodSignatures entry for %s", c.key)
        }
        if got.ParamSlots != c.slots {
            t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
        }
        if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
            t.Fatalf("%s GFunction mismatch", c.key)
        }
    }
}

func Test_MathContext_InitInt_DefaultsAndErrors(t *testing.T) {
    globals.InitGlobals("test")
    className := "java/math/MathContext"

    // Valid precision -> default HALF_UP
    mc := object.MakeEmptyObjectWithClassName(&className)
    if res := mconInitInt([]interface{}{mc, int64(5)}); res != nil {
        t.Fatalf("unexpected error: %v", res)
    }
    // precision field
    if v := mconGetPrecision([]interface{}{mc}).(int64); v != 5 {
        t.Fatalf("precision expected 5, got %d", v)
    }
    // roundingMode default HALF_UP (ordinal 4)
    rmObj := mconGetRoundingMode([]interface{}{mc}).(*object.Object)
    if rmObj == nil {
        t.Fatalf("roundingMode is nil")
    }
    if ordFld, ok := rmObj.FieldTable["ordinal"]; !ok || ordFld.Fvalue.(int64) != 4 {
        t.Fatalf("default roundingMode not HALF_UP, got %+v", ordFld)
    }

    // Negative precision -> IAE
    mc2 := object.MakeEmptyObjectWithClassName(&className)
    res := mconInitInt([]interface{}{mc2, int64(-1)})
    assertGErrBlkMC(t, res, excNames.IllegalArgumentException)
}

func Test_MathContext_InitIntRoundingMode_ValidAndNull(t *testing.T) {
    globals.InitGlobals("test")
    className := "java/math/MathContext"

    // Valid with HALF_DOWN
    mc := object.MakeEmptyObjectWithClassName(&className)
    rm := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("HALF_DOWN")} )
    if blk, ok := rm.(*GErrBlk); ok {
        t.Fatalf("failed to get RoundingMode.HALF_DOWN: %v", blk)
    }
    if res := mconInitIntRoundingMode([]interface{}{mc, int64(7), rm.(*object.Object)}); res != nil {
        t.Fatalf("unexpected error: %v", res)
    }
    if v := mconGetPrecision([]interface{}{mc}).(int64); v != 7 {
        t.Fatalf("precision expected 7, got %d", v)
    }
    rmObj := mconGetRoundingMode([]interface{}{mc}).(*object.Object)
    if ord := rmObj.FieldTable["ordinal"].Fvalue.(int64); ord != 5 { // HALF_DOWN ordinal 5
        t.Fatalf("roundingMode expected HALF_DOWN ordinal 5, got %d", ord)
    }

    // Null rounding mode -> NPE
    mc2 := object.MakeEmptyObjectWithClassName(&className)
    res := mconInitIntRoundingMode([]interface{}{mc2, int64(2), object.Null})
    assertGErrBlkMC(t, res, excNames.NullPointerException)
}

func Test_MathContext_InitString_ParseAndErrors(t *testing.T) {
    globals.InitGlobals("test")
    className := "java/math/MathContext"

    // Full string with rounding mode
    mc := object.MakeEmptyObjectWithClassName(&className)
    s := object.StringObjectFromGoString("precision=3 roundingMode=HALF_EVEN")
    if res := mconInitString([]interface{}{mc, s}); res != nil {
        t.Fatalf("unexpected error: %v", res)
    }
    if p := mconGetPrecision([]interface{}{mc}).(int64); p != 3 {
        t.Fatalf("precision expected 3, got %d", p)
    }
    rmObj := mconGetRoundingMode([]interface{}{mc}).(*object.Object)
    if ord := rmObj.FieldTable["ordinal"].Fvalue.(int64); ord != 6 { // HALF_EVEN ordinal 6
        t.Fatalf("roundingMode expected HALF_EVEN ordinal 6, got %d", ord)
    }

    // Default rounding mode when omitted -> HALF_UP
    mc2 := object.MakeEmptyObjectWithClassName(&className)
    s2 := object.StringObjectFromGoString("precision=9")
    if res := mconInitString([]interface{}{mc2, s2}); res != nil {
        t.Fatalf("unexpected error: %v", res)
    }
    rm2 := mconGetRoundingMode([]interface{}{mc2}).(*object.Object)
    if ord := rm2.FieldTable["ordinal"].Fvalue.(int64); ord != 4 {
        t.Fatalf("default roundingMode expected HALF_UP ordinal 4, got %d", ord)
    }

    // Null string -> NPE
    mc3 := object.MakeEmptyObjectWithClassName(&className)
    assertGErrBlkMC(t, mconInitString([]interface{}{mc3, object.Null}), excNames.NullPointerException)

    // Missing precision -> IAE
    mc4 := object.MakeEmptyObjectWithClassName(&className)
    s4 := object.StringObjectFromGoString("roundingMode=UP")
    assertGErrBlkMC(t, mconInitString([]interface{}{mc4, s4}), excNames.IllegalArgumentException)

    // Invalid precision value -> IAE
    mc5 := object.MakeEmptyObjectWithClassName(&className)
    s5 := object.StringObjectFromGoString("precision=abc roundingMode=DOWN")
    assertGErrBlkMC(t, mconInitString([]interface{}{mc5, s5}), excNames.IllegalArgumentException)

    // Invalid rounding mode name -> IAE
    mc6 := object.MakeEmptyObjectWithClassName(&className)
    s6 := object.StringObjectFromGoString("precision=2 roundingMode=ROUND_UP")
    assertGErrBlkMC(t, mconInitString([]interface{}{mc6, s6}), excNames.IllegalArgumentException)
}

func Test_MathContext_toString_Format(t *testing.T) {
    globals.InitGlobals("test")
    className := "java/math/MathContext"

    mc := object.MakeEmptyObjectWithClassName(&className)
    rm := rmodeValueOfString([]interface{}{object.StringObjectFromGoString("FLOOR")})
    if blk, ok := rm.(*GErrBlk); ok {
        t.Fatalf("failed to get RoundingMode.FLOOR: %v", blk)
    }
    if res := mconInitIntRoundingMode([]interface{}{mc, int64(12), rm.(*object.Object)}); res != nil {
        t.Fatalf("unexpected error: %v", res)
    }
    strObj := mconToString([]interface{}{mc}).(*object.Object)
    got := object.GoStringFromStringObject(strObj)
    want := "precision=12 roundingMode=FLOOR"
    if got != want {
        t.Fatalf("toString mismatch: want %q, got %q", want, got)
    }
}
