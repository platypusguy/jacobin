/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"testing"
)

/**
 * This test verifies BigDecimal.sqrt(MathContext) minimal behavior implemented in
 * bigdecimalSqrtContext. Our implementation requires a non-null MathContext,
 * throws ArithmeticException for negative values, and otherwise computes the
 * result using double math with sqrt, returning a BigDecimal via valueOf(double).
 *
 * We DO NOT attempt a huge Gauss–Legendre (Brent–Salamin) computation here; that
 * comment is illustrative only. Instead, we validate correctness and errors.
 */
func TestBigDecimalSqrt(t *testing.T) {
	bdutInit() // initialize globals/statics used by BigDecimal

	// Helper to make a MathContext with given precision using <init>(int)
	newMC := func(prec int64) *object.Object {
		className := "java/math/MathContext"
		mc := object.MakeEmptyObjectWithClassName(&className)
		if res := mconInitInt([]interface{}{mc, prec}); res != nil {
			t.Fatalf("mconInitInt returned error: %v", res)
		}
		return mc
	}

	t.Run("Null_MathContext_Yields_NPE", func(t *testing.T) {
		// sqrt with null MathContext -> NullPointerException
		// Make operand 4
		self := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
		if res := BigdecimalInitString([]interface{}{self, object.StringObjectFromGoString("4")}); res != nil {
			t.Fatalf("init BigDecimal(4) failed: %v", res)
		}
		ret := bigdecimalSqrtContext([]interface{}{self, object.Null})
		if blk, ok := ret.(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.NullPointerException {
			t.Fatalf("expected NullPointerException, got %T (%v)", ret, ret)
		}
	})

	t.Run("Negative_Value_Yields_ArithmeticException", func(t *testing.T) {
		self := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
		if res := BigdecimalInitString([]interface{}{self, object.StringObjectFromGoString("-1")}); res != nil {
			t.Fatalf("init BigDecimal(-1) failed: %v", res)
		}
		mc := newMC(10)
		ret := bigdecimalSqrtContext([]interface{}{self, mc})
		if blk, ok := ret.(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.ArithmeticException {
			t.Fatalf("expected ArithmeticException for sqrt(-1), got %T (%v)", ret, ret)
		}
	})

	t.Run("Sqrt_Positive_Basic_Values", func(t *testing.T) {
		cases := []struct {
			in   string
			want float64
		}{
			{"0", 0},
			{"1", 1},
			{"4", 2},
			{"2", math.Sqrt(2)},
			{"123.456", math.Sqrt(123.456)},
		}
		for _, c := range cases {
			self := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
			if res := BigdecimalInitString([]interface{}{self, object.StringObjectFromGoString(c.in)}); res != nil {
				t.Fatalf("init BigDecimal(%s) failed: %v", c.in, res)
			}
			mc := newMC(16)
			out := bigdecimalSqrtContext([]interface{}{self, mc})
			if blk, ok := out.(*ghelpers.GErrBlk); ok {
				t.Fatalf("unexpected error for sqrt(%s): %d %s", c.in, blk.ExceptionType, blk.ErrMsg)
			}
			bd := out.(*object.Object)
			got := bigdecimalDoubleValue([]interface{}{bd}).(float64)
			if math.IsNaN(got) || math.IsInf(got, 0) {
				t.Fatalf("sqrt(%s) produced invalid double: %v", c.in, got)
			}
			// Allow small tolerance due to double conversion and toBigDecimal from double
			if math.Abs(got-c.want) > 1e-12 {
				t.Fatalf("sqrt(%s) mismatch: got %.15g want %.15g", c.in, got, c.want)
			}
		}
	})
}
