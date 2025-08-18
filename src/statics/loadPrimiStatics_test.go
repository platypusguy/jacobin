/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package statics

import (
    "math"
    "testing"

    "jacobin/src/types"
)

// helper to isolate Statics map between tests
func withFreshStatics(t *testing.T, fn func()) {
    t.Helper()
    saved := Statics
    Statics = make(map[string]Static)
    defer func() { Statics = saved }()
    fn()
}

func TestLoadStaticsByte(t *testing.T) {
    withFreshStatics(t, func() {
        LoadStaticsByte()
        // Expected entries
        cases := []struct{
            key  string
            typ  string
            val  interface{}
        }{
            {"java/lang/Byte.BYTES", types.Int, int64(1)},
            {"java/lang/Byte.MAX_VALUE", types.Byte, int64(0x7f)},
            {"java/lang/Byte.MIN_VALUE", types.Byte, int64(0x80)},
            {"java/lang/Byte.SIZE", types.Int, int64(8)},
        }
        for _, c := range cases {
            st, ok := Statics[c.key]
            if !ok {
                t.Fatalf("missing static %s", c.key)
            }
            if st.Type != c.typ {
                t.Fatalf("%s type expected %s, got %s", c.key, c.typ, st.Type)
            }
            if st.Value != c.val {
                t.Fatalf("%s value expected %v, got %v", c.key, c.val, st.Value)
            }
        }
    })
}

func TestLoadStaticsCharacter_SampleSet(t *testing.T) {
    withFreshStatics(t, func() {
        LoadStaticsCharacter()
        // Basic size/bytes and value bounds
        if st, ok := Statics["java/lang/Character.BYTES"]; !ok || st.Type != types.Int || st.Value != int64(2) {
            t.Fatalf("Character.BYTES wrong: %+v", st)
        }
        if st, ok := Statics["java/lang/Character.SIZE"]; !ok || st.Type != types.Int || st.Value != int64(16) {
            t.Fatalf("Character.SIZE wrong: %+v", st)
        }
        if st, ok := Statics["java/lang/Character.MAX_CODE_POINT"]; !ok || st.Type != types.Int || st.Value != int64(1114111) {
            t.Fatalf("Character.MAX_CODE_POINT wrong: %+v", st)
        }
        if st, ok := Statics["java/lang/Character.MAX_VALUE"]; !ok || st.Type != types.Char || st.Value != rune(65535) {
            t.Fatalf("Character.MAX_VALUE wrong: %+v", st)
        }
        // Spot-check a couple category/directionality constants
        if st, ok := Statics["java/lang/Character.UPPERCASE_LETTER"]; !ok || st.Type != types.Byte || st.Value != int64(0x1) {
            t.Fatalf("Character.UPPERCASE_LETTER wrong: %+v", st)
        }
        if st, ok := Statics["java/lang/Character.DIRECTIONALITY_RIGHT_TO_LEFT"]; !ok || st.Type != types.Byte || st.Value != int64(0x1) {
            t.Fatalf("Character.DIRECTIONALITY_RIGHT_TO_LEFT wrong: %+v", st)
        }
    })
}

func TestLoadStaticsDouble(t *testing.T) {
    withFreshStatics(t, func() {
        LoadStaticsDouble()
        // simple integer fields
        if st := Statics["java/lang/Double.BYTES"]; st.Type != types.Int || st.Value != int64(8) {
            t.Fatalf("Double.BYTES wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.MAX_EXPONENT"]; st.Type != types.Int || st.Value != int64(1023) {
            t.Fatalf("Double.MAX_EXPONENT wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.MIN_EXPONENT"]; st.Type != types.Int || st.Value != int64(-1022) {
            t.Fatalf("Double.MIN_EXPONENT wrong: %+v", st)
        }
        // floats: MIN_NORMAL, MIN_VALUE, MAX_VALUE
        if st := Statics["java/lang/Double.MIN_NORMAL"]; st.Type != types.Double || st.Value.(float64) != float64(2.2250738585072014e-308) {
            t.Fatalf("Double.MIN_NORMAL wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.MIN_VALUE"]; st.Type != types.Double || st.Value.(float64) != float64(4.9e-324) {
            t.Fatalf("Double.MIN_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.MAX_VALUE"]; st.Type != types.Double || st.Value.(float64) != float64(1.7976931348623157e308) {
            t.Fatalf("Double.MAX_VALUE wrong: %+v", st)
        }
        // NaN/Inf checks
        if st := Statics["java/lang/Double.NaN"]; st.Type != types.Double || !math.IsNaN(st.Value.(float64)) {
            t.Fatalf("Double.NaN wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.NEGATIVE_INFINITY"]; st.Type != types.Double || !math.IsInf(st.Value.(float64), -1) {
            t.Fatalf("Double.NEGATIVE_INFINITY wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.POSITIVE_INFINITY"]; st.Type != types.Double || !math.IsInf(st.Value.(float64), +1) {
            t.Fatalf("Double.POSITIVE_INFINITY wrong: %+v", st)
        }
        if st := Statics["java/lang/Double.SIZE"]; st.Type != types.Int || st.Value != int64(64) {
            t.Fatalf("Double.SIZE wrong: %+v", st)
        }
    })
}

func TestLoadStaticsFloat(t *testing.T) {
    withFreshStatics(t, func() {
        LoadStaticsFloat()
        if st := Statics["java/lang/Float.BYTES"]; st.Type != types.Int || st.Value != int64(4) {
            t.Fatalf("Float.BYTES wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.MAX_EXPONENT"]; st.Type != types.Int || st.Value != int64(127) {
            t.Fatalf("Float.MAX_EXPONENT wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.MIN_EXPONENT"]; st.Type != types.Int || st.Value != int64(-126) {
            t.Fatalf("Float.MIN_EXPONENT wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.MIN_NORMAL"]; st.Type != types.Float || st.Value.(float64) != float64(1.1754943508222875e-38) {
            t.Fatalf("Float.MIN_NORMAL wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.MIN_VALUE"]; st.Type != types.Float || st.Value.(float64) != float64(1.401298464324817e-45) {
            t.Fatalf("Float.MIN_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.NaN"]; st.Type != types.Float || !math.IsNaN(st.Value.(float64)) {
            t.Fatalf("Float.NaN wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.NEGATIVE_INFINITY"]; st.Type != types.Float || !math.IsInf(st.Value.(float64), -1) {
            t.Fatalf("Float.NEGATIVE_INFINITY wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.POSITIVE_INFINITY"]; st.Type != types.Float || !math.IsInf(st.Value.(float64), +1) {
            t.Fatalf("Float.POSITIVE_INFINITY wrong: %+v", st)
        }
        if st := Statics["java/lang/Float.SIZE"]; st.Type != types.Int || st.Value != int64(32) {
            t.Fatalf("Float.SIZE wrong: %+v", st)
        }
    })
}

func TestLoadStaticsInteger_Long_Short(t *testing.T) {
    withFreshStatics(t, func() {
        LoadStaticsInteger()
        if st := Statics["java/lang/Integer.BYTES"]; st.Type != types.Int || st.Value != int64(4) {
            t.Fatalf("Integer.BYTES wrong: %+v", st)
        }
        if st := Statics["java/lang/Integer.MAX_VALUE"]; st.Type != types.Int || st.Value != int64(2147483647) {
            t.Fatalf("Integer.MAX_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Integer.MIN_VALUE"]; st.Type != types.Int || st.Value != int64(-2147483648) {
            t.Fatalf("Integer.MIN_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Integer.SIZE"]; st.Type != types.Int || st.Value != int64(32) {
            t.Fatalf("Integer.SIZE wrong: %+v", st)
        }
    })

    withFreshStatics(t, func() {
        LoadStaticsLong()
        if st := Statics["java/lang/Long.BYTES"]; st.Type != types.Int || st.Value != int64(8) {
            t.Fatalf("Long.BYTES wrong: %+v", st)
        }
        if st := Statics["java/lang/Long.MAX_VALUE"]; st.Type != types.Long || st.Value != int64(9223372036854775807) {
            t.Fatalf("Long.MAX_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Long.MIN_VALUE"]; st.Type != types.Long || st.Value != int64(-9223372036854775808) {
            t.Fatalf("Long.MIN_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Long.SIZE"]; st.Type != types.Int || st.Value != int64(64) {
            t.Fatalf("Long.SIZE wrong: %+v", st)
        }
    })

    withFreshStatics(t, func() {
        LoadStaticsShort()
        if st := Statics["java/lang/Short.BYTES"]; st.Type != types.Int || st.Value != int64(2) {
            t.Fatalf("Short.BYTES wrong: %+v", st)
        }
        if st := Statics["java/lang/Short.MAX_VALUE"]; st.Type != types.Short || st.Value != int64(32767) {
            t.Fatalf("Short.MAX_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Short.MIN_VALUE"]; st.Type != types.Short || st.Value != int64(-32768) {
            t.Fatalf("Short.MIN_VALUE wrong: %+v", st)
        }
        if st := Statics["java/lang/Short.SIZE"]; st.Type != types.Int || st.Value != int64(16) {
            t.Fatalf("Short.SIZE wrong: %+v", st)
        }
    })
}

func TestLoadStaticsMath_StrictMath(t *testing.T) {
    withFreshStatics(t, func() {
        LoadStaticsMath()
        if st := Statics["java/lang/Math.E"]; st.Type != types.Double || st.Value.(float64) != float64(2.718281828459045) {
            t.Fatalf("Math.E wrong: %+v", st)
        }
        if st := Statics["java/lang/Math.PI"]; st.Type != types.Double || st.Value.(float64) != float64(3.141592653589793) {
            t.Fatalf("Math.PI wrong: %+v", st)
        }
    })

    withFreshStatics(t, func() {
        LoadStaticsStrictMath()
        if st := Statics["java/lang/StrictMath.E"]; st.Type != types.Double || st.Value.(float64) != float64(2.718281828459045) {
            t.Fatalf("StrictMath.E wrong: %+v", st)
        }
        if st := Statics["java/lang/StrictMath.PI"]; st.Type != types.Double || st.Value.(float64) != float64(3.141592653589793) {
            t.Fatalf("StrictMath.PI wrong: %+v", st)
        }
    })
}
