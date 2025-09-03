/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
)

// bigdecimalDivideMathContext implements BigDecimal.divide(BigDecimal, MathContext)
// Minimal behavior:
// - Null MathContext -> NullPointerException
// - Uses MathContext.getRoundingMode(); if null -> NullPointerException
// - If precision <= 0: delegate to divide(BigDecimal, RoundingMode) using the MC's rounding mode
// - Else: delegate to divide(BigDecimal, int scale, RoundingMode) using dividend's scale as target scale
func bigdecimalDivideMathContext(params []interface{}) interface{} {
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return getGErrBlk(excNames.NullPointerException, "bigdecimalDivide: MathContext is null")
	}
	// Extract rounding mode from MathContext
	rmObj := mconGetRoundingMode([]interface{}{mc}).(*object.Object)
	if object.IsNull(rmObj) {
		return getGErrBlk(excNames.NullPointerException, "bigdecimalDivide: MathContext.roundingMode is null")
	}
	prec := mconGetPrecision([]interface{}{mc}).(int64)
	if prec <= 0 {
		// Unlimited precision: behave like divide(BigDecimal, RoundingMode)
		return bigdecimalDivideRoundingMode([]interface{}{dividend, divisor, rmObj})
	}
	// For now, choose target scale as the dividend's scale and apply rounding mode.
	sa := dividend.FieldTable["scale"].Fvalue.(int64)
	return bigdecimalDivideScaleRoundingMode([]interface{}{dividend, divisor, sa, rmObj})
}
