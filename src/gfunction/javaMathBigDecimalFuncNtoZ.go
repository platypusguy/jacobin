/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// bigdecimalNegate returns a BigDecimal with value = -this
// Extracts the internal unscaled BigInteger, negates it,
// and creates a new BigDecimal with scale = 0 and recalculated precision.
func bigdecimalNegate(params []interface{}) interface{} {
	// Implements BigDecimal.negate()
	bd := params[0].(*object.Object)

	// Extract the BigInteger intVal field
	dv := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger object to big.Int
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)

	// Negate the value
	negatedValue := new(big.Int).Neg(dvBigInt)

	// Create result BigDecimal object for the negated value
	result := bigDecimalObjectFromBigInt(negatedValue, int64(len(negatedValue.String())), int64(0))

	return result
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
// Uses big.Int.Exp for exponentiation on the unscaled value,
// and sets scale = 0 because pow affects unscaled value directly.
func bigdecimalPow(params []interface{}) interface{} {
	// Implements BigDecimal.pow(int exponent)
	bd := params[0].(*object.Object)
	exponent := params[1].(int64)

	// Extract BigInteger intVal field
	dv := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger object to big.Int
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)

	// Perform exponentiation
	resultValue := new(big.Int).Exp(dvBigInt, big.NewInt(exponent), nil)

	// Create result BigDecimal object for the power
	result := bigDecimalObjectFromBigInt(resultValue, int64(len(resultValue.String())), int64(0))

	return result
}

// bigdecimalPrecision simply returns the precision stored in the BigDecimal's FieldTable.
func bigdecimalPrecision(params []interface{}) interface{} {
	// Implements BigDecimal.precision()
	bd := params[0].(*object.Object)

	// Retrieve the precision field from the FieldTable
	precision := bd.FieldTable["precision"].Fvalue.(int64)

	return precision
}

// bigdecimalRemainder computes this % divisor.
// If divisor == 0, returns ArithmeticException to avoid division by zero.
// Result scale is 0 (remainder is integral).
func bigdecimalRemainder(params []interface{}) interface{} {
	// Implements BigDecimal.remainder(BigDecimal divisor)
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
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalRemainder: divide by zero")
	}

	// Perform Java-like remainder operation.
	remainder := javaLikeRemainder(dvBigInt, drBigInt)

	// Create result BigDecimal object for the remainder
	remObj := bigDecimalObjectFromBigInt(remainder, int64(len(remainder.String())), int64(0))

	return remObj
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
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalShortValueExact: value out of range for short")
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

// bigdecimalToBigInteger returns the floor of this BigDecimal as a BigInteger.
// Simply returns the unscaled BigInteger (ignoring scale).
func bigdecimalToBigInteger(params []interface{}) interface{} {
	// Implements BigDecimal.toBigInteger()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Retrieve the big.Int from the BigInteger object
	bigIntValue := intValObj.FieldTable["value"].Fvalue.(*big.Int)

	// Return the BigInteger object (as an *object.Object)
	biObj := makeBigIntegerFromBigInt(bigIntValue)
	return biObj
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
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalToBigIntegerExact: non-zero fractional part")
	}

	// Return the BigInteger value
	return makeBigIntegerFromBigInt(bigInt)
}

// bigdecimalToString returns a string representation of the BigDecimal,
// properly inserting a decimal point based on the scale.
func bigdecimalToString(params []interface{}) interface{} {
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
	decimalString := formatDecimalString(bigIntValue, scale)
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
	bd := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	strObj := object.StringObjectFromGoString(s)
	if res := bigdecimalInitString([]interface{}{bd, strObj}); res != nil {
		// pass through error block if any
		if _, ok := res.(*GErrBlk); ok {
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
	bigIntObj := bigIntegerFromInt64(val)

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
	bigIntObj := bigIntegerFromInt64(val)

	// Calculate precision
	precision := int64(len(strconv.FormatInt(val, 10)))
	if val < 0 {
		precision -= 1
	}

	// Create BigDecimal object with the provided scale
	bd := bigDecimalObjectFromBigInt(bigIntObj.FieldTable["value"].Fvalue.(*big.Int), precision, int64(scale))

	return bd
}
