package javaMath

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"math/big"
	"strconv"
	"strings"
)

/*
<init> functions
*/

func bigdecimalInitDouble(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	valObj := params[1].(float64)

	// Convert float64 to string with full precision
	valStr := strconv.FormatFloat(valObj, 'g', -1, 64)

	// Get *big.Int value and scale.
	bigInt, scale, ok := parseDecimalString(valStr)
	if !ok {
		errMsg := fmt.Sprintf("bigdecimalInitDouble: Failed to parse '%s'", valStr)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Create BigInteger from string
	bigIntObj := makeBigIntegerFromBigInt(bigInt)

	// Estimate precision
	precision := int64(len(strings.ReplaceAll(valStr, ".", "")))
	if bigInt.Sign() < 0 {
		precision -= 1
	}

	// Set fields
	setupBasicFields(self, bigIntObj, precision, scale)

	return nil
}

// bigdecimalInitIntLong: Set up a BigDecimal object based on an integer or long argument.
func bigdecimalInitIntLong(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	valInt64 := params[1].(int64)

	// Create a BigInteger object representing the value.
	bigIntObj := BigIntegerFromInt64(valInt64)

	// Compute precision: number of decimal digits in value.
	precision := int64(len(strconv.FormatInt(valInt64, 10)))
	if valInt64 < 0 {
		precision -= 1
	}

	// Assign fields to the BigDecimal object.
	setupBasicFields(self, bigIntObj, precision, int64(0))

	return nil
}

/*
BigdecimalInitString: Set up a BigDecimal object based on a string object argument.
Handles optional leading + or - sign.
Splits on decimal point correctly.
Keeps fractional zeros for scale.
Uses the combined integer+fraction string as the unscaled value.
Computes precision as the number of digits excluding leading zeros.
Builds the intVal BigInteger object.
Returns nil on success, or a NumberFormatException on invalid input.
*/
func BigdecimalInitString(params []interface{}) interface{} {
	bd := params[0].(*object.Object)
	strObj := params[1].(*object.Object)

	s := object.GoStringFromStringObject(strObj)
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: empty string")
	}

	// Handle optional leading sign
	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}
	if len(s) == 0 { // sign only is invalid
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid number format")
	}

	// Extract optional exponent part (scientific notation)
	exponent := int64(0)
	numPart := s
	if idx := strings.IndexAny(s, "eE"); idx != -1 {
		numPart = s[:idx]
		expPart := s[idx+1:]
		if len(expPart) == 0 {
			return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid exponent")
		}
		// Parse exponent with optional sign
		if expPart[0] == '+' || expPart[0] == '-' {
			if len(expPart) == 1 {
				return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid exponent")
			}
		}
		expVal, err := strconv.ParseInt(expPart, 10, 64)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid exponent")
		}
		exponent = expVal
	}

	// numPart should contain at least one digit
	hasDigit := false
	for i := 0; i < len(numPart); i++ {
		if numPart[i] >= '0' && numPart[i] <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid number format")
	}

	// Split integer and fractional parts of mantissa
	parts := strings.SplitN(numPart, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	// Validate that intPart and fracPart contain only digits (intPart may include leading zeros or be empty)
	for i := 0; i < len(intPart); i++ {
		if intPart[i] < '0' || intPart[i] > '9' {
			return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid number format")
		}
	}
	for i := 0; i < len(fracPart); i++ {
		if fracPart[i] < '0' || fracPart[i] > '9' {
			return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid number format")
		}
	}

	// Remove leading zeros from intPart; if empty, set to "0"
	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" {
		intPart = "0"
	}

	// Build unscaled string by concatenating intPart and fracPart (keep fractional zeros)
	unscaledStr := intPart + fracPart
	if unscaledStr == "" { // should not happen due to hasDigit check, but keep safe
		unscaledStr = "0"
	}

	// Parse unscaledStr into big.Int
	unscaledBigInt := new(big.Int)
	_, ok := unscaledBigInt.SetString(unscaledStr, 10)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "BigdecimalInitString: invalid number format")
	}

	if negative {
		unscaledBigInt.Neg(unscaledBigInt)
	}

	// Scale is length of fractional part adjusted by exponent
	scale := int64(len(fracPart)) - exponent

	// Precision is number of digits in unscaled value excluding sign and leading zeros
	precStr := strings.TrimLeft(unscaledStr, "0")
	if precStr == "" {
		precStr = "0"
	}
	precision := int64(len(precStr))

	// Create BigInteger object from unscaledBigInt
	biObj := makeBigIntegerFromBigInt(unscaledBigInt)

	// Set fields on BigDecimal object
	bd.FieldTable["intVal"] = object.Field{Fvalue: biObj}
	bd.FieldTable["scale"] = object.Field{Fvalue: scale}
	bd.FieldTable["precision"] = object.Field{Fvalue: precision}
	return nil
}

func bigdecimalInitBigInteger(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	biObj := params[1].(*object.Object)
	bigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)
	precision := precisionFromBigInt(bigInt)
	scale := int64(0)
	setupBasicFields(self, biObj, precision, scale)

	return nil
}

// bigdecimalInitDoubleContext implements BigDecimal.<init>(double, MathContext)
// Minimal behavior: NPE if MathContext is null; otherwise, delegate to the double initializer
func bigdecimalInitDoubleContext(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	val := params[1].(float64)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "BigDecimal.<init>(double, MathContext): MathContext is null")
	}
	return bigdecimalInitDouble([]interface{}{self, val})
}

// bigdecimalInitStringContext implements BigDecimal.<init>(String, MathContext)
// Minimal behavior: NPE if MathContext is null; otherwise, delegate to the string initializer
func bigdecimalInitStringContext(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	strObj := params[1].(*object.Object)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "BigDecimal.<init>(String, MathContext): MathContext is null")
	}
	return BigdecimalInitString([]interface{}{self, strObj})
}

// bigdecimalInitBigIntegerScale implements BigDecimal.<init>(BigInteger, int)
// Sets intVal to the provided BigInteger and the scale to the provided int; precision derived from value.
func bigdecimalInitBigIntegerScale(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	biObj := params[1].(*object.Object)
	scale := params[2].(int64)
	bigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)
	precision := precisionFromBigInt(bigInt)
	setupBasicFields(self, biObj, precision, scale)
	return nil
}

// bigdecimalInitBigIntegerContext implements BigDecimal.<init>(BigInteger, MathContext)
// Minimal behavior: NPE if MathContext is null; otherwise delegate to BigInteger-only init.
func bigdecimalInitBigIntegerContext(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	biObj := params[1].(*object.Object)
	mc := params[2].(*object.Object)
	if object.IsNull(mc) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "BigDecimal.<init>(BigInteger, MathContext): MathContext is null")
	}
	return bigdecimalInitBigInteger([]interface{}{self, biObj})
}
