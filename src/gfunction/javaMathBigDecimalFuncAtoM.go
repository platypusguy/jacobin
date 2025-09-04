/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"math/big"
)

// bigdecimalAbs returns the absolute value of the BigDecimal
func bigdecimalAbs(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	intVal := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	val := intVal.FieldTable["value"].Fvalue.(*big.Int)

	// Calculate the absolute value
	absVal := new(big.Int).Abs(val)

	// Create and return a new BigDecimal object with the absolute value
	return bigDecimalObjectFromBigInt(absVal, int64(len(absVal.String())), bd.FieldTable["scale"].Fvalue.(int64))
}

// bigdecimalAdd returns the result of adding this BigDecimal to the specified one
func bigdecimalAdd(params []interface{}) interface{} {
	// Extract BigDecimal objects
	bd1 := params[0].(*object.Object)
	bd2 := params[1].(*object.Object)

	// Extract unscaled values and scales
	intVal1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	intVal2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)
	val1 := new(big.Int).Set(intVal1.FieldTable["value"].Fvalue.(*big.Int))
	val2 := new(big.Int).Set(intVal2.FieldTable["value"].Fvalue.(*big.Int))
	s1 := bd1.FieldTable["scale"].Fvalue.(int64)
	s2 := bd2.FieldTable["scale"].Fvalue.(int64)

	// Align scales: use s = max(s1, s2)
	s := s1
	if s2 > s {
		s = s2
	}
	if s > s1 {
		// scale up val1 by 10^(s - s1)
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-s1), nil)
		val1.Mul(val1, mul)
	}
	if s > s2 {
		// scale up val2 by 10^(s - s2)
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-s2), nil)
		val2.Mul(val2, mul)
	}

	// Perform the addition on aligned unscaled values
	sum := new(big.Int).Add(val1, val2)
	precision := precisionFromBigInt(sum)
	return bigDecimalObjectFromBigInt(sum, precision, s)
}

// bigdecimalByteValueExact returns the exact byte value of this BigDecimal
func bigdecimalByteValueExact(params []interface{}) interface{} {
	// Extract BigDecimal object
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field.
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger object to *big.Int, then to int64.
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)
	i64 := bigInt.Int64()

	// Check if the value fits in a Java byte.
	if i64 > 127 || i64 < -128 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalByteValueExact: out of range for byte value")
	}

	// Return the exact byte value.
	return types.JavaByte(bigInt.Int64())
}

// bigdecimalCompareTo compares this BigDecimal to the specified BigDecimal.
// Returns a negative integer if this BigDecimal is less than the specified BigDecimal,
// zero if they are equal, and a positive integer if this BigDecimal is greater.
func bigdecimalCompareTo(params []interface{}) interface{} {
	// Extract BigDecimal objects
	bd1 := params[0].(*object.Object)
	bd2 := params[1].(*object.Object)

	// Extract BigInteger intVal fields
	intVal1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	intVal2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger objects to big.Int
	val1 := intVal1.FieldTable["value"].Fvalue.(*big.Int)
	val2 := intVal2.FieldTable["value"].Fvalue.(*big.Int)

	// Compare the two values
	return int64(val1.Cmp(val2))
}

// bigdecimalDivide returns the exact quotient of this BigDecimal divided by the specified one.
// Behavior per JDK:
// - If the exact quotient has a non-terminating decimal expansion, throw ArithmeticException.
// - Otherwise, return the exact result with the minimal scale needed to represent it.
func bigdecimalDivide(params []interface{}) interface{} {
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)

	// Extract unscaled values and scales
	dv := dividend.FieldTable["intVal"].Fvalue.(*object.Object)
	dr := divisor.FieldTable["intVal"].Fvalue.(*object.Object)
	a := dv.FieldTable["value"].Fvalue.(*big.Int)
	b := dr.FieldTable["value"].Fvalue.(*big.Int)
	sa := dividend.FieldTable["scale"].Fvalue.(int64)
	sb := divisor.FieldTable["scale"].Fvalue.(int64)

	// Division by zero
	if b.Sign() == 0 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivide: divide by zero")
	}

	// Work with absolute values and track sign
	neg := (a.Sign() < 0) != (b.Sign() < 0)
	absA := new(big.Int).Abs(a)
	absB := new(big.Int).Abs(b)

	// Form fraction N/D corresponding to (a * 10^(-sa)) / (b * 10^(-sb))
	// N = |a| * 10^sb ; D = |b| * 10^sa
	powNum := new(big.Int).Exp(big.NewInt(10), big.NewInt(sb), nil)
	N := new(big.Int).Mul(absA, powNum)
	powDen := new(big.Int).Exp(big.NewInt(10), big.NewInt(sa), nil)
	D := new(big.Int).Mul(absB, powDen)

	// Reduce fraction by GCD
	g := new(big.Int).GCD(nil, nil, N, D)
	if g.Sign() != 0 && g.Cmp(big.NewInt(1)) != 0 {
		N.Quo(N, g)
		D.Quo(D, g)
	}

	// Factor denominator into 2^x * 5^y * rest
	countFactors := func(n *big.Int, p int64) (int64, *big.Int) {
		zero := big.NewInt(0)
		pp := big.NewInt(p)
		cnt := int64(0)
		rem := new(big.Int)
		q := new(big.Int).Set(n)
		for {
			q, rem = new(big.Int).QuoRem(q, pp, rem)
			if rem.Cmp(zero) != 0 {
				break
			}
			cnt++
			n = q
		}
		return cnt, n
	}

	x, after2 := countFactors(new(big.Int).Set(D), 2)
	y, after5 := countFactors(after2, 5)
	if after5.Cmp(big.NewInt(1)) != 0 {
		// Denominator has prime factors other than 2 or 5 -> non-terminating
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivide: non-terminating decimal expansion; no rounding specified")
	}

	// Make denominator 1 by multiplying numerator by 5^x or 2^y accordingly; scale increases by max(x,y)
	addScale := x
	if y > x {
		addScale = y
	}
	if x > y {
		// multiply by 5^(x-y)
		mul := new(big.Int).Exp(big.NewInt(5), big.NewInt(x-y), nil)
		N.Mul(N, mul)
	} else if y > x {
		// multiply by 2^(y-x)
		mul := new(big.Int).Exp(big.NewInt(2), big.NewInt(y-x), nil)
		N.Mul(N, mul)
	}
	// Denominator now effectively 10^addScale, so final unscaled is N, final scale is addScale
	if neg && N.Sign() != 0 {
		N.Neg(N)
	}

	precision := precisionFromBigInt(N)
	return bigDecimalObjectFromBigInt(N, precision, addScale)
}

// bigdecimalDivideAndRemainder returns both the quotient and remainder of division
func bigdecimalDivideAndRemainder(params []interface{}) interface{} {
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)

	// Extract BigInteger intVal fields
	dv := dividend.FieldTable["intVal"].Fvalue.(*object.Object)
	dr := divisor.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger objects to big.Int
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)
	drBigInt := dr.FieldTable["value"].Fvalue.(*big.Int)

	// Check for division by zero
	if drBigInt.Sign() == 0 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivideAndRemainder: divide by zero")
	}

	// Perform division and remainder
	quotient := new(big.Int).Div(dvBigInt, drBigInt)
	remainder := new(big.Int).Mod(dvBigInt, drBigInt)

	// Create BigDecimal objects for the results
	quotObj := bigDecimalObjectFromBigInt(quotient, int64(len(quotient.String())), int64(0))
	remObj := bigDecimalObjectFromBigInt(remainder, int64(len(remainder.String())), int64(0))

	arrObj := makeArray2ElemsOfBigDecimal(quotObj, remObj)
	return arrObj
}

// bigdecimalDivideToIntegralValue returns the quotient of this BigDecimal divided by the divisor, truncating the result
func bigdecimalDivideToIntegralValue(params []interface{}) interface{} {
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)

	// Extract BigInteger intVal fields
	dv := dividend.FieldTable["intVal"].Fvalue.(*object.Object)
	dr := divisor.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger objects to big.Int
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)
	drBigInt := dr.FieldTable["value"].Fvalue.(*big.Int)

	// Check for division by zero
	if drBigInt.Sign() == 0 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivideToIntegralValue: divide by zero")
	}

	// Do integer division
	quotient := new(big.Int).Div(dvBigInt, drBigInt)

	// Create result BigDecimal object
	result := bigDecimalObjectFromBigInt(quotient, int64(len(quotient.String())), int64(0))

	return result
}

// bigdecimalDoubleValue returns the BigDecimal as a float64
func bigdecimalDoubleValue(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	val := bi.FieldTable["value"].Fvalue.(*big.Int)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	f := new(big.Float).SetInt(val)
	divisor := new(big.Float).SetFloat64(math.Pow10(int(scale)))
	f.Quo(f, divisor)
	result, _ := f.Float64()
	return result
}

// bigdecimalEquals checks if two BigDecimal values are equal
func bigdecimalEquals(params []interface{}) interface{} {
	bd1 := params[0].(*object.Object)
	bd2 := params[1].(*object.Object)

	// Compare scales first
	scale1 := bd1.FieldTable["scale"].Fvalue.(int64)
	scale2 := bd2.FieldTable["scale"].Fvalue.(int64)

	if scale1 != scale2 {
		return types.JavaBoolFalse // different scales means not equal
	}

	// Compare unscaled values (bigInt)
	bi1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	bi2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)

	unscaled1 := bi1.FieldTable["value"].Fvalue.(*big.Int)
	unscaled2 := bi2.FieldTable["value"].Fvalue.(*big.Int)

	if unscaled1.Cmp(unscaled2) != 0 {
		return types.JavaBoolFalse // different unscaled values means not equal
	}

	// If both scale and unscaled value are the same, they are equal
	return types.JavaBoolTrue
}

// bigdecimalFloatValue returns the BigDecimal as a float64 (same as doubleValue)
func bigdecimalFloatValue(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	val := bi.FieldTable["value"].Fvalue.(*big.Int)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	f := new(big.Float).SetInt(val)
	divisor := new(big.Float).SetFloat64(math.Pow10(int(scale)))
	f.Quo(f, divisor)
	result, _ := f.Float32()
	return float64(result)
}

// bigdecimalIntValue returns the BigDecimal as an int64.
func bigdecimalIntValue(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)
	return bigInt.Int64()
}

// bigdecimalIntValueExact returns int64 if value fits, else ArithmeticException
func bigdecimalIntValueExact(params []interface{}) interface{} {
	// Get BigDecimal object and scale value (must be 0).
	bd := params[0].(*object.Object)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	if scale != int64(0) {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalIntValueExact: scale is non-zero")
	}

	// Get intValue as an int64.
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)
	int64Value := bigInt.Int64()

	// Make sure that we are within int boundaries.
	if int64Value < math.MinInt32 || int64Value > math.MaxInt32 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalIntValueExact: value out of int range")
	}

	return int64Value
}

// bigdecimalLongValue returns the BigDecimal as an int64.
func bigdecimalLongValue(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	val := bi.FieldTable["value"].Fvalue.(*big.Int)
	return val.Int64()
}

// bigdecimalLongValueExact returns int64 if value fits, else ArithmeticException
func bigdecimalLongValueExact(params []interface{}) interface{} {
	// Get BigDecimal object and scale value (must be 0).
	bd := params[0].(*object.Object)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	if scale != int64(0) {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalLongValueExact: scale is non-zero")
	}

	// Get intValue as an int64.
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)
	int64Value := bigInt.Int64()

	// Make sure that we are within long boundaries.
	if int64Value < math.MinInt64 || int64Value > math.MaxInt64 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalLongValueExact: value out of long range")
	}

	return int64Value
}

// bigdecimalMax returns the greater of two BigDecimals
func bigdecimalMax(params []interface{}) interface{} {
	cmp := bigdecimalCompareTo(params)
	if cmp.(int64) >= 0 {
		return params[0]
	}
	return params[1]
}

// bigdecimalMin returns the lesser of two BigDecimals
func bigdecimalMin(params []interface{}) interface{} {
	cmp := bigdecimalCompareTo(params)
	if cmp.(int64) <= 0 {
		return params[0]
	}
	return params[1]
}

// bigdecimalMovePointLeft shifts the decimal point to the left by n
func bigdecimalMovePointLeft(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	num := params[1].(int64)

	// Extract intVal and scale
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	// Get the underlying *big.Int value
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)

	// New scale is original scale + num
	newScale := scale + num

	// Precision is length of digits in unscaled value
	precision := precisionFromBigInt(bigInt)

	// Create new BigDecimal object
	newBigInt := new(big.Int).Set(bigInt)
	newBD := bigDecimalObjectFromBigInt(newBigInt, precision, newScale)

	return newBD
}

// bigdecimalMovePointRight shifts the decimal point to the right by n
func bigdecimalMovePointRight(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	num := params[1].(int64)

	intVal := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	// Extract the unscaled big.Int
	bigInt := intVal.FieldTable["value"].Fvalue.(*big.Int)

	var newBigInt *big.Int
	var newScale int64

	if num <= scale {
		// Just reduce the scale
		newBigInt = new(big.Int).Set(bigInt)
		newScale = scale - num
	} else {
		// Shift the decimal point right by multiplying
		shift := num - scale
		factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(shift), nil)
		newBigInt = new(big.Int).Mul(bigInt, factor)
		newScale = 0
	}

	// Compute precision
	precision := precisionFromBigInt(newBigInt)

	// Construct and return a new BigDecimal object
	return bigDecimalObjectFromBigInt(newBigInt, precision, newScale)
}

// bigdecimalMultiply returns the result of multiplying two BigDecimals
func bigdecimalMultiply(params []interface{}) interface{} {
	bd1 := params[0].(*object.Object)
	bd2 := params[1].(*object.Object)
	intVal1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	intVal2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)
	val1 := intVal1.FieldTable["value"].Fvalue.(*big.Int)
	val2 := intVal2.FieldTable["value"].Fvalue.(*big.Int)
	result := new(big.Int).Mul(val1, val2)
	scale := bd1.FieldTable["scale"].Fvalue.(int64) + bd2.FieldTable["scale"].Fvalue.(int64)

	return bigDecimalObjectFromBigInt(result, int64(len(result.String())), scale)
}

// divide with scale and rounding mode: compute correct unscaled result and apply rounding
func bigdecimalDivideScaleRoundingMode(params []interface{}) interface{} {
	// params: this(BigDecimal), divisor(BigDecimal), scale(int), roundingMode(RoundingMode)
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)
	scaleParam := params[2].(int64)
	rmodeObj := params[3].(*object.Object)

	// Validate inputs
	if scaleParam < 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "bigdecimalDivide: negative scale")
	}
	// Per JDK, roundingMode must not be null even if no rounding would be required
	if object.IsNull(rmodeObj) {
		return getGErrBlk(excNames.NullPointerException, "bigdecimalDivide: RoundingMode is null")
	}
	// Extract unscaled and scales
	dv := dividend.FieldTable["intVal"].Fvalue.(*object.Object)
	dr := divisor.FieldTable["intVal"].Fvalue.(*object.Object)
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)
	drBigInt := dr.FieldTable["value"].Fvalue.(*big.Int)
	sa := dividend.FieldTable["scale"].Fvalue.(int64)
	sb := divisor.FieldTable["scale"].Fvalue.(int64)

	if drBigInt.Sign() == 0 {
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivide: divide by zero")
	}

	// Compute N/D where A=a*10^-sa, B=b*10^-sb, desired result scale = scaleParam
	// Target unscaled result ur approximates (a * 10^(sb+scaleParam)) / (b * 10^sa)
	// So N = |a| * 10^(sb+scaleParam); D = |b| * 10^sa; apply sign at the end.
	absA := new(big.Int).Abs(dvBigInt)
	absB := new(big.Int).Abs(drBigInt)

	// Build power-of-ten factor for numerator: sb + scaleParam
	powNum := new(big.Int).Exp(big.NewInt(10), big.NewInt(sb+scaleParam), nil)
	N := new(big.Int).Mul(absA, powNum)
	// Build power-of-ten factor for denominator: sa
	powDen := new(big.Int).Exp(big.NewInt(10), big.NewInt(sa), nil)
	D := new(big.Int).Mul(absB, powDen)

	// Integer division and remainder with positive values
	q := new(big.Int).Quo(N, D)
	r := new(big.Int).Mod(N, D)

	// If remainder is zero, exact result; check UNNECESSARY
	if r.Sign() == 0 {
		// If rounding mode is UNNECESSARY, it's fine because no rounding is needed.
		// Apply sign and return.
		if (dvBigInt.Sign() < 0) != (drBigInt.Sign() < 0) {
			q.Neg(q)
		}
		precision := precisionFromBigInt(q)
		return bigDecimalObjectFromBigInt(q, precision, scaleParam)
	}

	// Determine rounding behavior
	// Resolve rounding mode ordinal (must be an enum object)
	if object.IsNull(rmodeObj) {
		return getGErrBlk(excNames.NullPointerException, "bigdecimalDivide: RoundingMode is null")
	}
	ord, ok := extractRoundingModeOrdinal(rmodeObj)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "bigdecimalDivide: invalid RoundingMode")
	}
	// UNNECESSARY with non-zero remainder -> ArithmeticException
	if ord == 7 { // UNNECESSARY
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivide: rounding necessary")
	}

	increment := false
	// Sign of the true result
	positive := (dvBigInt.Sign() >= 0) == (drBigInt.Sign() >= 0)

	switch ord {
	case 0: // UP
		increment = true
	case 1: // DOWN
		increment = false
	case 2: // CEILING
		increment = positive
	case 3: // FLOOR
		increment = !positive
	case 4, 5, 6: // HALF_UP, HALF_DOWN, HALF_EVEN
		// Compare 2*remainder with denominator
		twiceR := new(big.Int).Lsh(r, 1) // r*2
		cmp := twiceR.Cmp(D)
		if cmp > 0 {
			increment = true
		} else if cmp < 0 {
			increment = false
		} else {          // exactly half
			if ord == 4 { // HALF_UP
				increment = true
			} else if ord == 5 { // HALF_DOWN
				increment = false
			} else { // HALF_EVEN
				// increment iff q is odd
				if q.Bit(0) == 1 {
					increment = true
				}
			}
		}
	default:
		// Fallback: behave like HALF_UP
		twiceR := new(big.Int).Lsh(r, 1)
		if twiceR.Cmp(D) >= 0 {
			increment = true
		}
	}

	if increment {
		q.Add(q, big.NewInt(1))
	}
	if !positive && q.Sign() != 0 {
		q.Neg(q)
	}

	precision := precisionFromBigInt(q)
	return bigDecimalObjectFromBigInt(q, precision, scaleParam)
}

// extractRoundingModeOrdinal tries to read the ordinal field from a RoundingMode enum object
func extractRoundingModeOrdinal(rmode *object.Object) (int, bool) {
	if object.IsNull(rmode) {
		return 0, false
	}
	fld, ok := rmode.FieldTable["ordinal"]
	if ok {
		if ival, ok := fld.Fvalue.(int64); ok {
			return int(ival), true
		}
	}
	// Fallback: if a name field exists, derive ordinal from it.
	if nameFld, ok := rmode.FieldTable["name"]; ok {
		if sObj, ok := nameFld.Fvalue.(*object.Object); ok {
			name := object.GoStringFromStringObject(sObj)
			for i, nm := range rmodeNames {
				if nm == name {
					return i, true
				}
			}
		}
	}
	return 0, false
}

// divide with rounding mode only: choose scale based on dividend's scale
// Rule inferred from JDK behavior: use scaleParam = this.scale (e.g., 6.0/x -> 1 decimal)
func bigdecimalDivideRoundingMode(params []interface{}) interface{} {
	// params: this(BigDecimal), divisor(BigDecimal), roundingMode(RoundingMode)
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)
	rmodeObj := params[2].(*object.Object)

	// Use the dividend's scale as the target scale
	sa := dividend.FieldTable["scale"].Fvalue.(int64)
	scaleParam := sa

	// Delegate to the precise divide-with-scale implementation
	return bigdecimalDivideScaleRoundingMode([]interface{}{dividend, divisor, scaleParam, rmodeObj})
}

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
