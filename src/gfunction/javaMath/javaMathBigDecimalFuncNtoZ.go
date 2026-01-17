/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// bigdecimalNegate returns a BigDecimal with value = -this, preserving scale and precision.
func bigdecimalNegate(params []interface{}) interface{} {
	// Implements BigDecimal.negate()
	bd := params[0].(*object.Object)

	// Extract the BigInteger intVal field
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	origScale := bd.FieldTable["scale"].Fvalue.(int64)
	origPrecision := bd.FieldTable["precision"].Fvalue.(int64)

	// Convert BigInteger object to big.Int and negate
	orig := intValObj.FieldTable["value"].Fvalue.(*big.Int)
	negatedValue := new(big.Int).Neg(orig)

	// Return new BigDecimal with same scale and precision
	return bigDecimalObjectFromBigInt(negatedValue, origPrecision, origScale)
}

// bigdecimalPlus returns a new BigDecimal identical to the input (unary plus).
// This effectively clones the BigDecimal, preserving precision and scale.
func bigdecimalPlus(params []interface{}) interface{} {
	bd := params[0].(*object.Object)

	// Clone intVal
	intVal := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	value := intVal.FieldTable["value"].Fvalue.(*big.Int)
	newBigInt := new(big.Int).Set(value)

	// Extract precision and scale
	precision := bd.FieldTable["precision"].Fvalue.(int64)
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	return bigDecimalObjectFromBigInt(newBigInt, precision, scale)
}

// bigdecimalPow computes this^exponent for non-negative exponents.
// For a BigDecimal with unscaled value u and scale s, (u * 10^{-s})^n = (u^n) * 10^{-s*n}.
// Therefore, result unscaled = u^n and result scale = s*n. Negative n -> ArithmeticException.
func bigdecimalPow(params []interface{}) interface{} {
	// Implements BigDecimal.pow(int exponent)
	bd := params[0].(*object.Object)
	exponent := params[1].(int64)

	if exponent < 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "bigdecimalPow: negative exponent")
	}
	// Special case: x^0 = 1 with scale 0
	if exponent == 0 {
		one := big.NewInt(1)
		return bigDecimalObjectFromBigInt(one, 1, 0)
	}

	// Extract BigInteger intVal and current scale
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	u := intValObj.FieldTable["value"].Fvalue.(*big.Int)
	s := bd.FieldTable["scale"].Fvalue.(int64)

	// Compute u^n
	resUnscaled := new(big.Int).Exp(u, big.NewInt(exponent), nil)
	resScale := s * exponent
	prec := precisionFromBigInt(resUnscaled)
	return bigDecimalObjectFromBigInt(resUnscaled, prec, resScale)
}

// bigdecimalPrecision simply returns the precision stored in the BigDecimal's FieldTable.
func bigdecimalPrecision(params []interface{}) interface{} {
	// Implements BigDecimal.precision()
	bd := params[0].(*object.Object)

	// Retrieve the precision field from the FieldTable
	precision := bd.FieldTable["precision"].Fvalue.(int64)

	return precision
}

// bigdecimalRemainder computes this - this.divideToIntegralValue(divisor) * divisor
func bigdecimalRemainder(params []interface{}) interface{} {
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)

	// Division by zero check
	dr := divisor.FieldTable["intVal"].Fvalue.(*object.Object)
	drBigInt := dr.FieldTable["value"].Fvalue.(*big.Int)
	if drBigInt.Sign() == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "bigdecimalRemainder: divide by zero")
	}
	q := bigdecimalDivideToIntegralValue([]interface{}{dividend, divisor})
	if blk, ok := q.(*ghelpers.GErrBlk); ok {
		return blk
	}
	qbd := q.(*object.Object)
	prod := bigdecimalMultiply([]interface{}{qbd, divisor})
	if blk, ok := prod.(*ghelpers.GErrBlk); ok {
		return blk
	}
	res := bigdecimalSubtract([]interface{}{dividend, prod.(*object.Object)})
	return res
}

// bigdecimalScale returns the current scale of the BigDecimal.
// Scale represents the number of digits to the right of the decimal point.
func bigdecimalScale(params []interface{}) interface{} {
	// Implements BigDecimal.scale()
	bd := params[0].(*object.Object)

	// Retrieve the scale from the FieldTable
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	return scale
}

// bigdecimalScaleByPowerOfTen adjusts the scale by subtracting 'num'.
// This corresponds to shifting the decimal point to the right by 'num' places.
func bigdecimalScaleByPowerOfTen(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	num := params[1].(int64)

	// Get current unscaled value and scale
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	precision := bd.FieldTable["precision"].Fvalue.(int64)

	// Adjust scale: newScale = scale - num
	newScale := scale - num

	return bigDecimalObjectFromBigInt(bigInt, precision, newScale)
}

// bigdecimalSetScale changes the scale of the BigDecimal to 'newScale'.
// If increasing scale, multiply unscaled value by 10^(newScale - oldScale).
// If decreasing, divide unscaled value by 10^(oldScale - newScale).
// Recomputes precision based on the new unscaled value.
func bigdecimalSetScale(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	newScale := params[1].(int64)

	intVal := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	oldBigInt := intVal.FieldTable["value"].Fvalue.(*big.Int)
	oldScale := bd.FieldTable["scale"].Fvalue.(int64)

	// If newScale is equal to current, return original
	if newScale == oldScale {
		return bd
	}

	diff := newScale - oldScale
	var newBigInt *big.Int

	if diff > 0 {
		// Scale increased: multiply by 10^diff
		multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil)
		newBigInt = new(big.Int).Mul(oldBigInt, multiplier)
	} else {
		// Scale decreased: divide by 10^(-diff)
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(-diff), nil)
		newBigInt = new(big.Int).Div(oldBigInt, divisor)
	}

	// Precision is recomputed from digit count
	precision := int64(len(strings.TrimLeft(newBigInt.String(), "-0")))
	if precision == 0 {
		precision = 1
	}

	return bigDecimalObjectFromBigInt(newBigInt, precision, newScale)
}

// bigdecimalSetScaleRoundingMode changes the scale using the provided RoundingMode
// Behavior:
// - If newScale == oldScale: return this
// - If newScale > oldScale: multiply unscaled by 10^(diff)
// - If newScale < oldScale: divide unscaled by 10^(oldScale-newScale) and round per RoundingMode
func bigdecimalSetScaleRoundingMode(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	newScale := params[1].(int64)
	rmodeObj := params[2].(*object.Object)

	if object.IsNull(rmodeObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalSetScale: RoundingMode is null")
	}

	intVal := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	oldBigInt := new(big.Int).Set(intVal.FieldTable["value"].Fvalue.(*big.Int))
	oldScale := bd.FieldTable["scale"].Fvalue.(int64)

	// Early return if scales match
	if newScale == oldScale {
		return bd
	}

	diff := newScale - oldScale
	// Increasing scale: multiply
	if diff > 0 {
		multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil)
		newBigInt := new(big.Int).Mul(oldBigInt, multiplier)
		precision := precisionFromBigInt(newBigInt)
		return bigDecimalObjectFromBigInt(newBigInt, precision, newScale)
	}

	// Decreasing scale: need rounding
	steps := new(big.Int).Exp(big.NewInt(10), big.NewInt(-diff), nil) // 10^(oldScale-newScale)
	abs := new(big.Int).Abs(oldBigInt)
	q := new(big.Int).Quo(abs, steps)
	r := new(big.Int).Mod(abs, steps)

	if r.Sign() == 0 {
		// exact; UNNECESSARY is fine since no rounding needed
		if oldBigInt.Sign() < 0 && q.Sign() != 0 {
			q.Neg(q)
		}
		precision := precisionFromBigInt(q)
		return bigDecimalObjectFromBigInt(q, precision, newScale)
	}

	ord, ok := extractRoundingModeOrdinal(rmodeObj)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "bigdecimalSetScale: invalid RoundingMode")
	}
	// UNNECESSARY with non-zero remainder -> ArithmeticException
	if ord == 7 { // UNNECESSARY
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "bigdecimalSetScale: rounding necessary")
	}

	increment := false
	positive := oldBigInt.Sign() >= 0
	D := steps
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
		twiceR := new(big.Int).Lsh(r, 1)
		cmp := twiceR.Cmp(D)
		if cmp > 0 {
			increment = true
		} else if cmp < 0 {
			increment = false
		} else { // exactly half
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
		// Fallback to HALF_UP-like
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
	return bigDecimalObjectFromBigInt(q, precision, newScale)
}

// bigdecimalShortValueExact converts the BigDecimal to an int16 exactly.
// Returns ArithmeticException if the value is out of the int16 range.
func bigdecimalShortValueExact(params []interface{}) interface{} {
	// Implements BigDecimal.shortValueExact()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	biObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger to big.Int
	dvBigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)

	// Check if the value fits in a short (16-bit signed integer)
	if dvBigInt.Cmp(big.NewInt(int64(math.MinInt16))) < 0 || dvBigInt.Cmp(big.NewInt(int64(math.MaxInt16))) > 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "bigdecimalShortValueExact: value out of range for short")
	}

	// Convert the BigInt to a short (int16)
	shortValue := int16(dvBigInt.Int64())

	return int64(shortValue)
}

// bigdecimalSignum returns the sign of the BigDecimal unscaled value:
// -1 if negative, 0 if zero, 1 if positive.
func bigdecimalSignum(params []interface{}) interface{} {
	// Implements BigDecimal.signum()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	biObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger to big.Int
	dvBigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)

	// Determine the sign of the value
	sign := dvBigInt.Sign()

	return int64(sign) // -1, 0, or 1
}

// bigdecimalStripTrailingZeros removes trailing zeros from the unscaled value
// and adjusts the scale accordingly. Updates precision to match.
func bigdecimalStripTrailingZeros(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	oldScale := bd.FieldTable["scale"].Fvalue.(int64)
	biObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	oldBigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)

	stripTrailingZeros := func(unscaled *big.Int, scale int64) (*big.Int, int64) {
		if unscaled.Sign() == 0 {
			//println("stripTrailingZeros: input unscaled is zero")
			return big.NewInt(0), 0
		}

		ten := big.NewInt(10)
		mod := new(big.Int)
		u := new(big.Int).Set(unscaled)
		s := scale

		for {
			u.QuoRem(u, ten, mod)
			//println("stripTrailingZeros: mod =", mod.String(), "scale =", s)
			if mod.Sign() != 0 {
				// Restore u to its previous value since the division left a remainder
				u.Mul(u, ten).Add(u, mod)
				break
			}
			s--
		}
		return u, s
	}

	precisionFromBigInt := func(bi *big.Int) int64 {
		if bi.Sign() == 0 {
			//println("precisionFromBigInt: input is zero")
			return 1
		}
		str := bi.String()
		str = strings.TrimLeft(str, "-0")
		if len(str) == 0 {
			//println("precisionFromBigInt: trimmed string is empty")
			return 1
		}
		return int64(len(str))
	}

	newBigInt, newScale := stripTrailingZeros(oldBigInt, oldScale)
	//println("bigdecimalStripTrailingZeros: newBigInt =", newBigInt.String(), "newScale =", newScale)

	newPrecision := precisionFromBigInt(newBigInt)
	//println("bigdecimalStripTrailingZeros: newPrecision =", newPrecision)

	newBD := bigDecimalObjectFromBigInt(newBigInt, newPrecision, newScale)

	return newBD
}

// bigdecimalSubtract subtracts the specified BigDecimal from this one.
// Aligns scales to the maximum of the two before subtracting, per BigDecimal semantics.
func bigdecimalSubtract(params []interface{}) interface{} {
	// Implements BigDecimal.subtract(BigDecimal subtrahend)
	minuendBD := params[0].(*object.Object)
	subtrahendBD := params[1].(*object.Object)

	// Extract unscaled values and scales
	minuendBI := minuendBD.FieldTable["intVal"].Fvalue.(*object.Object)
	subtrahendBI := subtrahendBD.FieldTable["intVal"].Fvalue.(*object.Object)
	val1 := new(big.Int).Set(minuendBI.FieldTable["value"].Fvalue.(*big.Int))
	val2 := new(big.Int).Set(subtrahendBI.FieldTable["value"].Fvalue.(*big.Int))
	s1 := minuendBD.FieldTable["scale"].Fvalue.(int64)
	s2 := subtrahendBD.FieldTable["scale"].Fvalue.(int64)

	// Align scales to s = max(s1, s2)
	s := s1
	if s2 > s {
		s = s2
	}
	if s > s1 {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-s1), nil)
		val1.Mul(val1, mul)
	}
	if s > s2 {
		mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(s-s2), nil)
		val2.Mul(val2, mul)
	}

	// Perform subtraction on aligned unscaled values: val1 - val2
	resultBigInt := new(big.Int).Sub(val1, val2)
	precision := precisionFromBigInt(resultBigInt)
	return bigDecimalObjectFromBigInt(resultBigInt, precision, s)
}

// bigdecimalToBigInteger returns the integer part of this BigDecimal as a BigInteger,
// truncating toward zero (i.e., ignoring the fractional part implied by scale).
func bigdecimalToBigInteger(params []interface{}) interface{} {
	// Implements BigDecimal.toBigInteger()
	bd := params[0].(*object.Object)

	// Extract unscaled BigInteger and scale
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	unscaled := intValObj.FieldTable["value"].Fvalue.(*big.Int)
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	// If no fractional digits, return the unscaled value directly.
	if scale <= 0 {
		return makeBigIntegerFromBigInt(unscaled)
	}

	// Compute quotient = unscaled / 10^scale with truncation toward zero
	div := pow10(scale)
	q := new(big.Int).Quo(new(big.Int).Set(unscaled), div)
	return makeBigIntegerFromBigInt(q)
}

// bigdecimalToBigIntegerExact converts the BigDecimal to BigInteger exactly.
// Throws ArithmeticException if BigDecimal has a fractional part (non-zero scale).
func bigdecimalToBigIntegerExact(params []interface{}) interface{} {
	bd := params[0].(*object.Object)

	// Extract the BigInteger intVal field from BigDecimal
	bi := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	bigInt := bi.FieldTable["value"].Fvalue.(*big.Int)

	// Check for any fractional part (scale != 0)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	if scale != 0 {
		// If scale is non-zero, the BigDecimal has a fractional part and cannot be converted exactly
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "bigdecimalToBigIntegerExact: non-zero fractional part")
	}

	// Return the BigInteger value
	return makeBigIntegerFromBigInt(bigInt)
}

// bigdecimalToString returns a string representation of the BigDecimal,
// properly inserting a decimal point based on the scale.
func BigdecimalToString(params []interface{}) interface{} {
	// Implements BigDecimal.toString()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Retrieve the big.Int from the BigInteger object
	bigIntValue := intValObj.FieldTable["value"].Fvalue.(*big.Int)

	// Retrieve the scale
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	// Format the string representation, including the scale.
	// Handle the scale to produce decimal point position if necessary.
	decimalString := FormatDecimalString(bigIntValue, scale)
	return object.StringObjectFromGoString(decimalString)
}

// bigdecimalUnscaledValue returns the unscaled BigInteger value underlying this BigDecimal.
func bigdecimalUnscaledValue(params []interface{}) interface{} {
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field.
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Retrieve the big.Int from the BigInteger object.
	bigIntValue := intValObj.FieldTable["value"].Fvalue.(*big.Int)

	// Return the unscaled value as a new BigInteger.
	biObj := makeBigIntegerFromBigInt(bigIntValue)
	return biObj
}

// bigdecimalValueOfDouble constructs a BigDecimal from a double (float64)
// following Java semantics: BigDecimal.valueOf(double) == new BigDecimal(Double.toString(val)).
// In particular, integral finite doubles like 6.0 must produce a decimal string with ".0" so the
// resulting BigDecimal has a non-zero scale (e.g., 6.0 -> scale 1).
func bigdecimalValueOfDouble(params []interface{}) interface{} {
	// Implements BigDecimal.valueOf(double val)
	value := params[0].(float64)

	// Build a Java-like string for the double. Java's Double.toString(6.0) -> "6.0".
	// Go's FormatFloat with 'g' may emit "6"; ensure a decimal point for integral values.
	s := strconv.FormatFloat(value, 'g', -1, 64)
	if !strings.ContainsAny(s, ".eE") {
		// No decimal point or exponent -> append .0 to match Java
		s = s + ".0"
	}

	// Delegate to string-based constructor logic to preserve fractional zeros in scale
	bd := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
	strObj := object.StringObjectFromGoString(s)
	if res := BigdecimalInitString([]interface{}{bd, strObj}); res != nil {
		// pass through error block if any
		if _, ok := res.(*ghelpers.GErrBlk); ok {
			return res
		}
	}
	return bd
}

// bigdecimalValueOfLong creates a BigDecimal with scale 0 from an int64.
func bigdecimalValueOfLong(params []interface{}) interface{} {
	// Implements BigDecimal.valueOf(long val)
	val := params[0].(int64)

	// Create BigInteger object from the long value
	bigIntObj := BigIntegerFromInt64(val)

	// Calculate precision
	precision := int64(len(strconv.FormatInt(val, 10)))
	if val < 0 {
		precision -= 1
	}

	// Create BigDecimal object with scale 0
	bd := bigDecimalObjectFromBigInt(bigIntObj.FieldTable["value"].Fvalue.(*big.Int), precision, 0)

	return bd
}

// bigdecimalValueOfLongInt creates a BigDecimal with a specified scale from an int64.
func bigdecimalValueOfLongInt(params []interface{}) interface{} {
	// Implements BigDecimal.valueOf(long val, int scale)
	val := params[0].(int64)
	scale := params[1].(int64)

	// Create BigInteger object from the long value
	bigIntObj := BigIntegerFromInt64(val)

	// Calculate precision
	precision := int64(len(strconv.FormatInt(val, 10)))
	if val < 0 {
		precision -= 1
	}

	// Create BigDecimal object with the provided scale
	bd := bigDecimalObjectFromBigInt(bigIntObj.FieldTable["value"].Fvalue.(*big.Int), precision, int64(scale))

	return bd
}

// bigdecimalNegateContext implements BigDecimal.negate(MathContext)
// Minimal: NPE if MathContext is null; else delegate to negate()
func bigdecimalNegateContext(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	mc := params[1].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalNegate: MathContext is null")
	}
	return bigdecimalNegate([]interface{}{bd})
}

// bigdecimalPlusContext implements BigDecimal.plus(MathContext)
// Minimal: NPE if MathContext is null; else delegate to plus()
func bigdecimalPlusContext(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	mc := params[1].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalPlus: MathContext is null")
	}
	return bigdecimalPlus([]interface{}{bd})
}

// bigdecimalPowContext implements BigDecimal.pow(int, MathContext)
// Minimal: NPE if MathContext is null; else delegate to pow(int)
func bigdecimalPowContext(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	exponent := params[1].(int64)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalPow: MathContext is null")
	}
	return bigdecimalPow([]interface{}{bd, exponent})
}

// bigdecimalRemainderContext implements BigDecimal.remainder(BigDecimal, MathContext)
// Minimal: NPE if MathContext is null; else delegate to remainder(BigDecimal)
func bigdecimalRemainderContext(params []interface{}) interface{} {
	dividend := params[0].(*object.Object)
	divisor := params[1].(*object.Object)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalRemainder: MathContext is null")
	}
	return bigdecimalRemainder([]interface{}{dividend, divisor})
}

// bigdecimalRoundContext implements BigDecimal.round(MathContext)
// Minimal: NPE if MathContext is null; returns this unchanged
func bigdecimalRoundContext(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	mc := params[1].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalRound: MathContext is null")
	}
	// Could implement precision-based rounding; minimal behavior returns bd unchanged
	return bd
}

// bigdecimalSqrtContext implements BigDecimal.sqrt(MathContext)
// Minimal: NPE if MathContext is null; ArithmeticException for negative values; else sqrt via float
func bigdecimalSqrtContext(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	mc := params[1].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalSqrt: MathContext is null")
	}
	// Check for negative value using doubleValue
	dv := bigdecimalDoubleValue([]interface{}{bd}).(float64)
	if dv < 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "bigdecimalSqrt: square root of negative value")
	}
	res := math.Sqrt(dv)
	return bigdecimalValueOfDouble([]interface{}{res})
}

// bigdecimalSubtractContext implements BigDecimal.subtract(BigDecimal, MathContext)
// Minimal: NPE if MathContext is null; else delegate to subtract(BigDecimal)
func bigdecimalSubtractContext(params []interface{}) interface{} {
	bd1 := params[0].(*object.Object)
	bd2 := params[1].(*object.Object)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigdecimalSubtract: MathContext is null")
	}
	return bigdecimalSubtract([]interface{}{bd1, bd2})
}

// bigdecimalUlp implements BigDecimal.ulp()
// Returns a BigDecimal equal to 1 scaled by this.scale (i.e., 10^-scale)
func bigdecimalUlp(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	one := big.NewInt(1)
	return bigDecimalObjectFromBigInt(one, 1, scale)
}
