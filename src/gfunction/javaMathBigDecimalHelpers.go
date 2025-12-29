/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

// Cached powers of ten to avoid recomputing 10^k repeatedly during BigDecimal ops
var pow10Cache = map[int64]*big.Int{}
var pow10lock = &sync.Mutex{}

// pow10 returns a cached big.Int representing 10^exp. The returned *big.Int is immutable by callers.
func pow10(exp int64) *big.Int {
	pow10lock.Lock()
	defer pow10lock.Unlock()
	if exp <= 0 {
		if exp == 0 {
			return big.NewInt(1)
		}
		// For negative exponents, callers should handle scaling differently; return 0 to avoid misuse
		return big.NewInt(0)
	}
	if v, ok := pow10Cache[exp]; ok {
		return v
	}
	v := new(big.Int).Exp(big.NewInt(10), big.NewInt(exp), nil)
	pow10Cache[exp] = v
	return v
}

/*
Helper Functions
*/

// precisionFromBigInt: From a *big.Int, compute the precision.
// Note that the absolute value of the argument *big.Int must be used.
func precisionFromBigInt(arg *big.Int) int64 {
	return int64(len((new(big.Int).Abs(arg)).Text(10)))
}

// loadStaticsBigDecimal: Load the static fields for BigDouble.
func loadStaticsBigDecimal() {
	INFLATED := int64(-9223372036854775808)
	_ = statics.AddStatic(classNameBigDecimal+".INFLATED", statics.Static{Type: types.Long, Value: INFLATED})
	addStaticBigDecimal("ZERO", int64(0))
	addStaticBigDecimal("ONE", int64(1))
	addStaticBigDecimal("TWO", int64(2))
	addStaticBigDecimal("TEN", int64(10))
}

// addStaticBigDecimal:
// * Form a BigInteger object.
// * Set the value field of the BigInteger object = argValue.
// * Add a BigDecimal static field with the supplied argName whose value is the BigInteger object.
func addStaticBigDecimal(argName string, argValue int64) {
	bd := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	bi := object.MakeEmptyObjectWithClassName(&classNameBigInteger)
	var params []interface{}
	InitBigIntegerField(bi, argValue)
	params = append(params, bd)
	params = append(params, bi)
	bigdecimalInitBigInteger(params)
	_ = statics.AddStatic(classNameBigDecimal+"."+argName, statics.Static{Type: types.BigDecimal, Value: bd})
}

func setupBasicFields(self, bigIntObj *object.Object, precision, scale int64) {
	object.ClearFieldTable(self)
	self.FieldTable["intVal"] = object.Field{Ftype: types.BigInteger, Fvalue: bigIntObj}
	self.FieldTable["scale"] = object.Field{Ftype: types.Int, Fvalue: scale}
	self.FieldTable["precision"] = object.Field{Ftype: types.Int, Fvalue: precision}
	self.FieldTable["intCompact"] = object.Field{Ftype: types.Long,
		Fvalue: statics.GetStaticValue(classNameBigDecimal, "INFLATED")}
}

// bigDecimalObjectFromBigInt: Given a *big.Int, precision, and scale, make a BigDecimal object.
func bigDecimalObjectFromBigInt(bigInt *big.Int, precision, scale int64) *object.Object {
	bdObj := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	// Create BigInteger object for field intVal.
	bigIntObj := object.MakeEmptyObjectWithClassName(&classNameBigInteger)
	setBigIntegerFields(bigIntObj, bigInt)

	// Set fields into the BigDecimal object.
	setupBasicFields(bdObj, bigIntObj, precision, scale)

	return bdObj
}

// Make a BigInteger object from an int64.
func bigIntegerFromInt64(arg int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameBigInteger)
	InitBigIntegerField(obj, arg)
	return obj
}

// parseDecimalString parses a string representation of a decimal number
// and returns the unscaled big integer, the scale, and a boolean indicating success.
func parseDecimalString(s string) (*big.Int, int64, bool) {
	// Trim any whitespace
	s = strings.TrimSpace(s)

	// Handle scientific notation (e.g., "3.1416E0" or "314.16e-2")
	exponentIndex := strings.IndexAny(s, "eE")

	// Initialize scale and exponent
	scale := int64(0)
	exponent := int64(0)

	// If there's an exponent, separate the number and exponent
	if exponentIndex != -1 {
		// Split into number part and exponent part
		numberPart := s[:exponentIndex]
		exponentPart := s[exponentIndex+1:]

		// Parse the exponent
		exp, err := strconv.ParseInt(exponentPart, 10, 64)
		if err != nil {
			return nil, 0, false // Invalid exponent
		}
		exponent = exp

		// Update the input string to the number part only (without the exponent)
		s = numberPart
	}

	// Split the number into integer and fractional parts
	parts := strings.Split(s, ".")

	// If the string does not contain at least one part (integer or fractional), it's invalid
	if len(parts) == 0 || len(parts[0]) == 0 {
		return nil, 0, false // Invalid number
	}

	// Default scale is the number of digits after the decimal
	var unscaledStr string
	if len(parts) > 1 && parts[1] != "0" {
		scale = int64(len(parts[1])) // Scale is based on the number of digits after the decimal
		unscaledStr = parts[0] + parts[1]
	} else {
		unscaledStr = parts[0]
	}

	// Now, convert the string (without the decimal) to a big integer
	unscaled := new(big.Int)
	unscaled.SetString(unscaledStr, 10)

	// If the string cannot be converted to a big integer, it's invalid
	if unscaled.Cmp(big.NewInt(0)) == 0 && s != "0" {
		return nil, 0, false // Invalid number (like "ABC")
	}

	// Adjust the scale based on the exponent
	scale -= exponent

	// Return the unscaled value, scale, and success flag
	return unscaled, scale, true
}

// setBigIntegerFields: Given the BigInteger object and the *big.Int, set the BigInteger object fields.
func setBigIntegerFields(obj *object.Object, bigInt *big.Int) {
	field := object.Field{Ftype: types.BigInteger, Fvalue: bigInt}
	obj.FieldTable["value"] = field
	fldSign := object.Field{Ftype: types.BigInteger, Fvalue: int64(bigInt.Sign())}
	obj.FieldTable["signum"] = fldSign
}

// makeBigIntegerFromBigInt: Given a *big.Int, make a BigInteger object.
func makeBigIntegerFromBigInt(bigIntValue *big.Int) *object.Object {
	biObj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, bigIntValue)
	fldSign := object.Field{Ftype: types.BigInteger, Fvalue: int64(bigIntValue.Sign())}
	biObj.FieldTable["signum"] = fldSign
	return biObj
}

// makeBigIntegerFromString: Make a BigInteger object from a Go object.
func makeBigIntegerFromString(str string) (*object.Object, *GErrBlk) {
	var zz = new(big.Int)
	_, ok := zz.SetString(str, 10)
	if !ok {
		errMsg := fmt.Sprintf("makeBigIntegerFromString: string (%s) not all numerics", str)
		return nil, getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Create BigInteger object with value set to zz.
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj, nil
}

// makeArray2ElemsOfBigDecimal: Make a 2-element array of BigDecimal objects.
func makeArray2ElemsOfBigDecimal(bd1, bd2 *object.Object) *object.Object {
	ref := "[L" + classNameBigDecimal + ";"
	arr := []*object.Object{bd1, bd2}
	obj := object.MakePrimitiveObject("["+classNameBigDecimal, ref, arr)
	return obj
}

// stripTrailingZeros removes trailing zeros from a string representation of a decimal.
// It returns a new *big.Int and the adjusted scale.
func stripTrailingZeros(unscaled *big.Int, scale int64) (*big.Int, int64) {
	// If the scale is 0, no need to strip any zeros
	if scale == 0 {
		return unscaled, scale
	}

	// Convert the big integer to a string to handle its decimal representation
	unscaledStr := unscaled.String()

	// Check for trailing zeros in the unscaled number string
	for {
		// Check if the last character is a zero (trailing zero)
		if len(unscaledStr) > 1 && unscaledStr[len(unscaledStr)-1] == '0' {
			// Remove the trailing zero and decrease the scale
			unscaledStr = unscaledStr[:len(unscaledStr)-1]
			scale--
		} else {
			break
		}
	}

	// If the string is empty or just "0", reset the unscaled value to 0
	if len(unscaledStr) == 0 || unscaledStr == "0" {
		unscaled.SetInt64(0)
		scale = 0
	} else {
		// Convert the modified string back to big.Int
		unscaled.SetString(unscaledStr, 10)
	}

	return unscaled, scale
}

// float64ToDecimalComponents converts a float64 to unscaled *big.Int, precision, and scale
func float64ToDecimalComponents(arg float64) (*big.Int, int64, int64) {
	// Handle special cases
	if math.IsNaN(arg) || math.IsInf(arg, 0) {
		return nil, 0, 0
	}

	// Format with high precision to preserve all decimal digits
	str := strconv.FormatFloat(arg, 'f', -1, 64)

	// Remove sign for analysis
	negative := strings.HasPrefix(str, "-")
	if negative {
		str = str[1:]
	}

	parts := strings.Split(str, ".")
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	// Count significant digits
	sigDigits := strings.TrimLeft(intPart, "0") + fracPart
	sigDigits = strings.TrimLeft(sigDigits, "0")
	precision := int64(len(sigDigits))

	// Calculate scale
	scale := int64(len(fracPart))

	// Form unscaled string by removing dot
	unscaledStr := intPart + fracPart
	unscaledStr = strings.TrimLeft(unscaledStr, "0")
	if unscaledStr == "" {
		unscaledStr = "0"
	}

	unscaled := new(big.Int)
	unscaled.SetString(unscaledStr, 10)
	if negative {
		unscaled.Neg(unscaled)
	}

	return unscaled, precision, scale
}

// formatDecimalString: Given a *big.Int signed & unscaled quantity and the int64 scale, produce a Go string.
func formatDecimalString(unscaled *big.Int, scale int64) string {
	isNegative := unscaled.Sign() < 0
	absStr := new(big.Int).Abs(unscaled).String()

	// Add leading zeros if absStr is too short for the scale
	if int64(len(absStr)) <= scale {
		zerosToAdd := int(scale) - len(absStr) + 1
		absStr = strings.Repeat("0", zerosToAdd) + absStr
	}

	// Preliminary intPart and fracPart.
	intPart := absStr[:len(absStr)-int(scale)]
	fracPart := absStr[len(absStr)-int(scale):]

	// If intPart is empty, replace it with "0".
	var result string
	if intPart == "" {
		intPart = "0"
	}

	// If fracPart is empty, do not include the decimal place nor the fracPart in the result.
	if fracPart == "" {
		result = intPart
	} else {
		result = intPart + "." + fracPart
	}

	// If negative, prepend the minus sign.
	if isNegative {
		result = "-" + result
	}

	return result
}

/*
javaLikeRemainder - be like Java when returning remainders (signed).

big.Int.Mod returns remainder in [0, divisor).
Java’s BigDecimal.remainder returns remainder with sign same as dividend.
To mimic Java in Go, adjust the result when dividend is negative.
*/
func javaLikeRemainder(dividend, divisor *big.Int) *big.Int {
	result := new(big.Int).Mod(dividend, divisor)
	if dividend.Sign() < 0 && result.Sign() != 0 {
		if divisor.Sign() > 0 {
			result.Sub(result, divisor) // r = r - divisor (divisor > 0)
		} else {
			result.Add(result, divisor) // r = r + divisor (divisor < 0)
		}
	}
	return result
}
