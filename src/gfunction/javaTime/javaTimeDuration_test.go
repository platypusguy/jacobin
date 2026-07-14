/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestDurationClinit(t *testing.T) {
	globals.InitGlobals("test")
	durationClinit(nil)

	s, ok := statics.QueryStatic("java/time/Duration", "ZERO")
	if !ok {
		t.Fatal("Duration.ZERO not found in statics")
	}
	dur := s.Value.(*object.Object)
	sec := dur.FieldTable["seconds"].Fvalue.(int64)
	nano := dur.FieldTable["nanos"].Fvalue.(int64)

	if sec != 0 || nano != 0 {
		t.Errorf("Expected Duration.ZERO to be (0, 0), got (%d, %d)", sec, nano)
	}
}

func TestDurationFactoriesAndAccessors(t *testing.T) {
	globals.InitGlobals("test")

	tests := []struct {
		name     string
		factory  func([]any) any
		input    int64
		expectedS int64
		expectedN int32
	}{
		{"ofSeconds", durationOfSeconds, 10, 10, 0},
		{"ofMillis", durationOfMillis, 1500, 1, 500000000},
		{"ofNanos", durationOfNanos, 123456789, 0, 123456789},
		{"ofMinutes", durationOfMinutes, 2, 120, 0},
		{"ofHours", durationOfHours, 1, 3600, 0},
		{"ofDays", durationOfDays, 1, 86400, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.factory([]any{tt.input})
			dur, ok := res.(*object.Object)
			if !ok {
				t.Fatalf("Expected *object.Object, got %T", res)
			}
			s := dur.FieldTable["seconds"].Fvalue.(int64)
			n := dur.FieldTable["nanos"].Fvalue.(int64)
			if s != tt.expectedS || int32(n) != tt.expectedN {
				t.Errorf("Expected (%d, %d), got (%d, %d)", tt.expectedS, tt.expectedN, s, n)
			}

			// Test accessors
			secRes := durationGetSeconds([]any{dur})
			if secRes.(int64) != tt.expectedS {
				t.Errorf("getSeconds: expected %d, got %v", tt.expectedS, secRes)
			}
			nanoRes := durationGetNano([]any{dur})
			if nanoRes.(int64) != int64(tt.expectedN) {
				t.Errorf("getNano: expected %d, got %v", tt.expectedN, nanoRes)
			}
		})
	}
}

func TestDurationArithmetic(t *testing.T) {
	globals.InitGlobals("test")

	d1 := createDuration(10, 500)
	d2 := createDuration(5, 600)

	// plus(Duration)
	resPlus := durationPlus([]any{d1, d2})
	durPlus := resPlus.(*object.Object)
	if durPlus.FieldTable["seconds"].Fvalue.(int64) != 15 || durPlus.FieldTable["nanos"].Fvalue.(int64) != 1100 {
		t.Errorf("plus: expected (15, 1100), got (%v, %v)", durPlus.FieldTable["seconds"].Fvalue, durPlus.FieldTable["nanos"].Fvalue)
	}

	// minus(Duration)
	resMinus := durationMinus([]any{d1, d2})
	durMinus := resMinus.(*object.Object)
	if durMinus.FieldTable["seconds"].Fvalue.(int64) != 4 || durMinus.FieldTable["nanos"].Fvalue.(int64) != 999999900 {
		t.Errorf("minus: expected (4, 999999900), got (%v, %v)", durMinus.FieldTable["seconds"].Fvalue, durMinus.FieldTable["nanos"].Fvalue)
	}

	// multipliedBy(long)
	resMult := durationMultipliedBy([]any{d1, int64(3)})
	durMult := resMult.(*object.Object)
	if durMult.FieldTable["seconds"].Fvalue.(int64) != 30 || durMult.FieldTable["nanos"].Fvalue.(int64) != 1500 {
		t.Errorf("multipliedBy: expected (30, 1500), got (%v, %v)", durMult.FieldTable["seconds"].Fvalue, durMult.FieldTable["nanos"].Fvalue)
	}

	// negated()
	resNeg := durationNegated([]any{d1})
	durNeg := resNeg.(*object.Object)
	if durNeg.FieldTable["seconds"].Fvalue.(int64) != -11 || durNeg.FieldTable["nanos"].Fvalue.(int64) != 999999500 {
		t.Errorf("negated: expected (-11, 999999500), got (%v, %v)", durNeg.FieldTable["seconds"].Fvalue, durNeg.FieldTable["nanos"].Fvalue)
	}

	// abs()
	resAbs := durationAbs([]any{durNeg})
	durAbs := resAbs.(*object.Object)
	if durAbs.FieldTable["seconds"].Fvalue.(int64) != 10 || durAbs.FieldTable["nanos"].Fvalue.(int64) != 500 {
		t.Errorf("abs: expected (10, 500), got (%v, %v)", durAbs.FieldTable["seconds"].Fvalue, durAbs.FieldTable["nanos"].Fvalue)
	}
}

func TestDurationOverflow(t *testing.T) {
	globals.InitGlobals("test")

	const maxLong = int64(9223372036854775807)

	t.Run("plusSeconds overflow", func(t *testing.T) {
		d := durationOfSeconds([]any{maxLong}).(*object.Object)
		res := durationPlusSeconds([]any{d, int64(1)})
		if _, ok := res.(*ghelpers.GErrBlk); ok {
			// Success - it returned an error block
		} else {
			t.Errorf("Expected ArithmeticException on overflow, but got %T: %+v", res, res)
		}
	})

	t.Run("ofDays overflow", func(t *testing.T) {
		res := durationOfDays([]any{maxLong})
		if _, ok := res.(*ghelpers.GErrBlk); ok {
			// Success
		} else {
			t.Errorf("Expected ArithmeticException on ofDays overflow, but got %T", res)
		}
	})

	t.Run("negated overflow", func(t *testing.T) {
		const minLong = int64(-9223372036854775808)
		d := durationOfSeconds([]any{minLong}).(*object.Object)
		res := durationNegated([]any{d})
		if _, ok := res.(*ghelpers.GErrBlk); ok {
			// Success
		} else {
			t.Errorf("Expected ArithmeticException on negating MinInt64 seconds, but got %T", res)
		}
	})

	t.Run("toNanos overflow", func(t *testing.T) {
		d := durationOfSeconds([]any{maxLong}).(*object.Object)
		res := durationToNanos([]any{d})
		if _, ok := res.(*ghelpers.GErrBlk); ok {
			// Success
		} else {
			t.Errorf("Expected ArithmeticException on toNanos overflow, but got %T", res)
		}
	})
}

func TestDurationConversions(t *testing.T) {
	globals.InitGlobals("test")

	dur := createDuration(3661, 500000000) // 1 hour, 1 minute, 1 second, 500ms

	if res := durationToDays([]any{dur}); res.(int64) != 0 {
		t.Errorf("toDays: expected 0, got %v", res)
	}
	if res := durationToHours([]any{dur}); res.(int64) != 1 {
		t.Errorf("toHours: expected 1, got %v", res)
	}
	if res := durationToMinutes([]any{dur}); res.(int64) != 61 {
		t.Errorf("toMinutes: expected 61, got %v", res)
	}
	if res := durationToMillis([]any{dur}); res.(int64) != 3661500 {
		t.Errorf("toMillis: expected 3661500, got %v", res)
	}
}

func TestDurationComparison(t *testing.T) {
	globals.InitGlobals("test")

	d1 := createDuration(10, 500)
	d2 := createDuration(10, 600)
	d3 := createDuration(9, 999999999)

	// equals
	if res := durationEquals([]any{d1, d1}); res.(types.JavaBool) != types.JavaBoolTrue {
		t.Errorf("equals(d1, d1) should be true")
	}
	if res := durationEquals([]any{d1, d2}); res.(types.JavaBool) != types.JavaBoolFalse {
		t.Errorf("equals(d1, d2) should be false")
	}

	// compareTo
	if res := durationCompareTo([]any{d1, d2}); res.(int64) >= 0 {
		t.Errorf("compareTo(d1, d2) should be negative, got %v", res)
	}
	if res := durationCompareTo([]any{d1, d3}); res.(int64) <= 0 {
		t.Errorf("compareTo(d1, d3) should be positive, got %v", res)
	}
	if res := durationCompareTo([]any{d1, d1}); res.(int64) != 0 {
		t.Errorf("compareTo(d1, d1) should be 0, got %v", res)
	}

	// isNegative / isZero
	if res := durationIsNegative([]any{d1}); res.(types.JavaBool) != types.JavaBoolFalse {
		t.Errorf("isNegative(d1) should be false")
	}
	neg := durationNegated([]any{d1}).(*object.Object)
	if res := durationIsNegative([]any{neg}); res.(types.JavaBool) != types.JavaBoolTrue {
		t.Errorf("isNegative(neg) should be true")
	}
	zero := createDuration(0, 0)
	if res := durationIsZero([]any{zero}); res.(types.JavaBool) != types.JavaBoolTrue {
		t.Errorf("isZero(zero) should be true")
	}
}

func TestDurationStringAndParse(t *testing.T) {
	globals.InitGlobals("test")

	tests := []struct {
		input    string
		expectedS int64
		expectedN int32
	}{
		{"PT1H", 3600, 0},
		{"PT1M", 60, 0},
		{"PT1S", 1, 0},
		{"PT1.5S", 1, 500000000},
		{"PT0.000000001S", 0, 1},
		{"-PT1S", -1, 0},
		{"PT-1S", -1, 0},
		{"PT-6H+3M", -6*3600 + 3*60, 0},
		{"PT-1M-1S", -61, 0},
		{"-PT-1M-1S", 61, 0},
		{"P1DT1H", 86400 + 3600, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			inputObj := object.StringObjectFromGoString(tt.input)
			res := durationParse([]any{inputObj})
			dur, ok := res.(*object.Object)
			if !ok {
				t.Fatalf("parse failed for %s: %v", tt.input, res)
			}
			s := dur.FieldTable["seconds"].Fvalue.(int64)
			n := dur.FieldTable["nanos"].Fvalue.(int64)
			if s != tt.expectedS || int32(n) != tt.expectedN {
				t.Errorf("Parse %s: expected (%d, %d), got (%d, %d)", tt.input, tt.expectedS, tt.expectedN, s, n)
			}

			// Test toString
			strRes := durationToString([]any{dur})
			strObj := strRes.(*object.Object)
			str := object.GoStringFromStringObject(strObj)
			
			// Re-parse the toString output and it should result in the same duration
			res2 := durationParse([]any{strObj})
			dur2 := res2.(*object.Object)
			s2 := dur2.FieldTable["seconds"].Fvalue.(int64)
			n2 := dur2.FieldTable["nanos"].Fvalue.(int64)
			if s2 != s || n2 != n {
				t.Errorf("ToString/Parse roundtrip failed for %s: original (%d, %d), toString %s, re-parsed (%d, %d)", tt.input, s, n, str, s2, n2)
			}
		})
	}
}

func TestDurationParts(t *testing.T) {
	globals.InitGlobals("test")

	// 90061 seconds = 1 day (86400) + 1 hour (3600) + 1 minute (60) + 1 second (1)
	// nanos = 123,000,000
	dur := createDuration(90061, 123000000)

	if res := durationToDaysPart([]any{dur}); res.(int64) != 1 {
		t.Errorf("toDaysPart: expected 1, got %v", res)
	}
	if res := durationToHoursPart([]any{dur}); res.(int64) != 1 {
		t.Errorf("toHoursPart: expected 1, got %v", res)
	}
	if res := durationToMinutesPart([]any{dur}); res.(int64) != 1 {
		t.Errorf("toMinutesPart: expected 1, got %v", res)
	}
	if res := durationToSecondsPart([]any{dur}); res.(int64) != 1 {
		t.Errorf("toSecondsPart: expected 1, got %v", res)
	}
	if res := durationToMillisPart([]any{dur}); res.(int64) != 123 {
		t.Errorf("toMillisPart: expected 123, got %v", res)
	}
	if res := durationToNanosPart([]any{dur}); res.(int64) != 123000000 {
		t.Errorf("toNanosPart: expected 123000000, got %v", res)
	}
}
