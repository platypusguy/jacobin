/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
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

	// Extract BigInteger intVal fields
	intVal1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	intVal2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger objects to big.Int
	val1 := intVal1.FieldTable["value"].Fvalue.(*big.Int)
	val2 := intVal2.FieldTable["value"].Fvalue.(*big.Int)

	// Perform the addition
	bigInt := new(big.Int).Add(val1, val2)

	// Calculate the precision and scale (same scale as the first BigDecimal)
	precision := precisionFromBigInt(bigInt)
	scale := bd1.FieldTable["scale"].Fvalue.(int64)

	// Create a new BigDecimal object with the result
	return bigDecimalObjectFromBigInt(bigInt, precision, scale)
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

// bigdecimalDivide returns the result of dividing this BigDecimal by the specified one
func bigdecimalDivide(params []interface{}) interface{} {
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
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalDivide: divide by zero")
	}

	// Do integer division
	quotient := new(big.Int).Div(dvBigInt, drBigInt)

	// Create result BigDecimal object
	result := bigDecimalObjectFromBigInt(quotient, int64(len(quotient.String())), int64(0))

	return result
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
		return object.BooleanObjectFromGoBoolean(false) // different scales means not equal
	}

	// Compare unscaled values (bigInt)
	bi1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	bi2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)

	unscaled1 := bi1.FieldTable["value"].Fvalue.(*big.Int)
	unscaled2 := bi2.FieldTable["value"].Fvalue.(*big.Int)

	if unscaled1.Cmp(unscaled2) != 0 {
		return object.BooleanObjectFromGoBoolean(false) // different unscaled values means not equal
	}

	// If both scale and unscaled value are the same, they are equal
	return object.BooleanObjectFromGoBoolean(true)
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
