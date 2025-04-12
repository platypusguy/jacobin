/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"math"
	"math/big"
	"strconv"
)

// bigdecimalNegate returns a BigDecimal whose value is the negation of the current BigDecimal.
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

// bigdecimalPlus returns a BigDecimal whose value is the sum of the current BigDecimal and the specified one.
func bigdecimalPlus(params []interface{}) interface{} {
	bd1 := params[0].(*object.Object)
	bd2 := params[1].(*object.Object)

	// Extract BigInteger intVal fields
	dv1 := bd1.FieldTable["intVal"].Fvalue.(*object.Object)
	dv2 := bd2.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger objects to big.Int
	dv1BigInt := dv1.FieldTable["value"].Fvalue.(*big.Int)
	dv2BigInt := dv2.FieldTable["value"].Fvalue.(*big.Int)

	// Perform addition
	resultValue := new(big.Int).Add(dv1BigInt, dv2BigInt)

	// Create result BigDecimal object for the sum
	result := bigDecimalObjectFromBigInt(resultValue, int64(len(resultValue.String())), int64(0))

	return result
}

// bigdecimalPow returns a BigDecimal whose value is the result of raising this BigDecimal to the specified power.
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

// bigdecimalPrecision returns the precision of this BigDecimal, i.e., the number of decimal digits.
func bigdecimalPrecision(params []interface{}) interface{} {
	// Implements BigDecimal.precision()
	bd := params[0].(*object.Object)

	// Retrieve the precision field from the FieldTable
	precision := bd.FieldTable["precision"].Fvalue.(int64)

	return precision
}

// bigdecimalRemainder returns the remainder when this BigDecimal is divided by the specified one.
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

	// Perform remainder operation
	remainder := new(big.Int).Mod(dvBigInt, drBigInt)

	// Create result BigDecimal object for the remainder
	remObj := bigDecimalObjectFromBigInt(remainder, int64(len(remainder.String())), int64(0))

	return remObj
}

// bigdecimalScale returns the scale of the BigDecimal object.
func bigdecimalScale(params []interface{}) interface{} {
	// Implements BigDecimal.scale()
	bd := params[0].(*object.Object)

	// Retrieve the scale from the FieldTable
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	return scale
}

// bigdecimalScaleByPowerOfTen scales the BigDecimal by the specified power of ten.
func bigdecimalScaleByPowerOfTen(params []interface{}) interface{} {
	// Implements BigDecimal.scaleByPowerOfTen(int n)
	bd := params[0].(*object.Object)
	num := params[1].(int64)

	// Retrieve the current scale and adjust it by the power of ten
	currentScale := bd.FieldTable["scale"].Fvalue.(int64)
	newScale := currentScale + num

	// Update the scale in the BigDecimal's FieldTable
	bd.FieldTable["scale"] = object.Field{Fvalue: newScale, Ftype: types.Int}

	return bd
}

// bigdecimalSetScale returns a new BigDecimal with the specified scale.
func bigdecimalSetScale(params []interface{}) interface{} {
	// Implements BigDecimal.setScale(int newScale)
	bd := params[0].(*object.Object)
	newScale := params[1].(int64)

	// Extract BigInteger intVal field
	biObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger to big.Int
	dvBigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)

	// Create result BigDecimal object with new scale
	// Assuming that the scale adjustment doesn't change the underlying value
	result := bigDecimalObjectFromBigInt(dvBigInt, int64(len(dvBigInt.String())), newScale)

	return result
}

// bigdecimalShortValueExact returns the exact short value of this BigDecimal.
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

// bigdecimalSignum returns the signum function of this BigDecimal.
// It returns -1, 0, or 1 depending on whether the value is negative, zero, or positive, respectively.
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

func bigdecimalStripTrailingZeros(params []interface{}) interface{} {
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	biObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	dvBigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)

	// Strip trailing zeros by dividing by 10 until the remainder is non-zero.
	for dvBigInt.BitLen() > 0 && dvBigInt.Mod(dvBigInt, big.NewInt(10)).Sign() == 0 {
		dvBigInt.Div(dvBigInt, big.NewInt(10))
	}

	// Update BigDecimal with new value and adjusted precision.
	bd.FieldTable["intVal"] = object.Field{Ftype: types.BigInteger, Fvalue: biObj}
	bd.FieldTable["precision"] = object.Field{Ftype: types.Int, Fvalue: int64(len(dvBigInt.String()))}

	return bd
}

// bigdecimalSubtract returns a BigDecimal representing the result of subtracting the specified BigDecimal from this BigDecimal.
func bigdecimalSubtract(params []interface{}) interface{} {
	// Implements BigDecimal.subtract(BigDecimal subtrahend)
	minuendBD := params[0].(*object.Object)
	subtrahendBD := params[1].(*object.Object)

	// Extract BigInteger intVal fields
	minuendBI := minuendBD.FieldTable["intVal"].Fvalue.(*object.Object)
	subtrahendBI := subtrahendBD.FieldTable["intVal"].Fvalue.(*object.Object)

	// Convert BigInteger objects to big.Int
	dvBigInt := minuendBI.FieldTable["value"].Fvalue.(*big.Int)
	drBigInt := subtrahendBI.FieldTable["value"].Fvalue.(*big.Int)

	// Perform subtraction
	resultBigInt := new(big.Int).Sub(dvBigInt, drBigInt)

	// Create a new BigDecimal object with the result
	result := bigDecimalObjectFromBigInt(resultBigInt, int64(len(resultBigInt.String())), minuendBD.FieldTable["scale"].Fvalue.(int64))

	return result
}

// bigdecimalToBigInteger returns the BigInteger value represented by this BigDecimal.
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

func bigdecimalToBigIntegerExact(params []interface{}) interface{} {
	// Implements BigDecimal.toBigIntegerExact()
	bd := params[0].(*object.Object)

	// Extract the BigInteger intVal field from BigDecimal
	dv := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)

	// Check for any fractional part (scale != 0)
	scale := bd.FieldTable["scale"].Fvalue.(int64)
	if scale != 0 {
		// If scale is non-zero, the BigDecimal has a fractional part and cannot be converted exactly
		return getGErrBlk(excNames.ArithmeticException, "bigdecimalToBigIntegerExact: non-zero fractional part")
	}

	// Return the BigInteger value
	return dvBigInt
}

// bigdecimalToEngineeringString returns the engineering string representation of this BigDecimal.
func bigdecimalToEngineeringString(params []interface{}) interface{} {
	// Implements BigDecimal.toEngineeringString()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Retrieve the big.Int from the BigInteger object
	bigIntValue := intValObj.FieldTable["value"].Fvalue.(*big.Int)

	// Convert the big.Int to an engineering string
	// In Go, the best way to get the engineering string would involve adjusting the scale
	// and then formatting the result to match the engineering representation
	// Assuming scientific formatting for now
	engineeringString := fmt.Sprintf("%e", bigIntValue)

	return object.StringObjectFromGoString(engineeringString)
}

// bigdecimalToPlainString returns the plain string representation of this BigDecimal without scientific notation.
func bigdecimalToPlainString(params []interface{}) interface{} {
	// Implements BigDecimal.toPlainString()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Retrieve the big.Int from the BigInteger object
	bigIntValue := intValObj.FieldTable["value"].Fvalue.(*big.Int)

	// Convert the big.Int to a plain string (no scientific notation)
	plainString := bigIntValue.String()

	return object.StringObjectFromGoString(plainString)
}

// bigdecimalToString returns the string representation of this BigDecimal.
func bigdecimalToString(params []interface{}) interface{} {
	// Implements BigDecimal.toString()
	bd := params[0].(*object.Object)

	// Extract BigInteger intVal field
	intValObj := bd.FieldTable["intVal"].Fvalue.(*object.Object)

	// Retrieve the big.Int from the BigInteger object
	bigIntValue := intValObj.FieldTable["value"].Fvalue.(*big.Int)

	// Retrieve the scale
	scale := bd.FieldTable["scale"].Fvalue.(int64)

	// Format the string representation, including the scale
	// Handle the scale to produce decimal point position if necessary
	decimalString := bigIntValue.String()
	if scale > 0 {
		// Add the decimal point
		if len(decimalString) <= int(scale) {
			decimalString = "0." + fmt.Sprintf("%0*s", int(scale)-len(decimalString), decimalString)
		} else {
			decimalString = decimalString[:len(decimalString)-int(scale)] + "." + decimalString[len(decimalString)-int(scale):]
		}
	}

	return object.StringObjectFromGoString(decimalString)
}

func bigdecimalUlp(params []interface{}) interface{} {
	bd := params[0].(*object.Object)

	// Extract the BigInteger intVal field from BigDecimal
	dv := bd.FieldTable["intVal"].Fvalue.(*object.Object)
	dvBigInt := dv.FieldTable["value"].Fvalue.(*big.Int)

	// Compute ULP by creating a BigDecimal with a value of 1 (smallest possible unit)
	// and subtracting the current BigDecimal value from it
	ulp := new(big.Int).Add(dvBigInt, big.NewInt(1)) // ULP is current value + 1
	if dvBigInt.Sign() < 0 {
		ulp = new(big.Int).Sub(dvBigInt, big.NewInt(1)) // If the value is negative, ULP is current value - 1
	}

	// Create a new BigDecimal object with the computed ULP value
	// Set scale to 0 and precision to 1 (since this is the smallest difference)
	ulpBigDecimal := bigDecimalObjectFromBigInt(ulp, int64(len(ulp.String())), int64(0))

	return ulpBigDecimal
}

// bigdecimalUnscaledValue returns the unscaled value of this BigDecimal as a BigInteger.
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

// bigdecimalValueOfDouble returns a BigDecimal initialized with the given double value.
func bigdecimalValueOfDouble(params []interface{}) interface{} {
	// Implements BigDecimal.valueOf(double val)
	val := params[0].(float64)

	// Convert the double value to a string, and then to a BigInteger object
	valStr := strconv.FormatFloat(val, 'f', -1, 64)
	bigIntObj, gerr := makeBigIntegerFromString(valStr)
	if gerr != nil {
		return gerr
	}

	// Extract the *big.Int from the BigInteger object
	bigIntVal := bigIntObj.FieldTable["value"].Fvalue.(*big.Int)

	// Calculate the precision: number of digits in the string representation of the value
	precision := int64(len(valStr))

	// Create a BigDecimal object with the BigInteger value, scale 0, and precision based on the string length
	bd := bigDecimalObjectFromBigInt(bigIntVal, precision, 0)

	return bd
}

// bigdecimalValueOfLong creates a BigDecimal from a long value.
func bigdecimalValueOfLong(params []interface{}) interface{} {
	// Implements BigDecimal.valueOf(long val)
	val := params[0].(int64)

	// Create BigInteger object from the long value
	bigIntObj := bigIntegerFromInt64(val)

	// Calculate precision
	precision := int64(len(strconv.FormatInt(val, 10)))

	// Create BigDecimal object with scale 0
	bd := bigDecimalObjectFromBigInt(bigIntObj.FieldTable["value"].Fvalue.(*big.Int), precision, 0)

	return bd
}

// bigdecimalValueOfLongInt creates a BigDecimal from a long and an int value.
func bigdecimalValueOfLongInt(params []interface{}) interface{} {
	// Implements BigDecimal.valueOf(long val, int scale)
	val := params[0].(int64)
	scale := params[1].(int64)

	// Create BigInteger object from the long value
	bigIntObj := bigIntegerFromInt64(val)

	// Calculate precision
	precision := int64(len(strconv.FormatInt(val, 10)))

	// Create BigDecimal object with the provided scale
	bd := bigDecimalObjectFromBigInt(bigIntObj.FieldTable["value"].Fvalue.(*big.Int), precision, int64(scale))

	return bd
}
