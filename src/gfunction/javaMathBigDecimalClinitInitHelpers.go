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
	"jacobin/statics"
	"jacobin/types"
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
	bigInt, scale, gerr := parseBigDecimalString(valStr)
	if gerr != nil {
		return gerr
	}

	// Create BigInteger from string
	bigIntObj := makeBigIntegerFromBigInt(bigInt)

	// Estimate precision
	precision := int64(len(strings.ReplaceAll(valStr, ".", "")))

	// Set fields
	setupBasicFields(self, bigIntObj, precision, scale)

	return nil
}

// bigdecimalInitIntLong: Set up a BigDecimal object based on an integer or long argument.
func bigdecimalInitIntLong(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	valInt64 := params[1].(int64)

	// Create a BigInteger object representing the value.
	bigIntObj := bigIntegerFromInt64(valInt64)

	// Compute precision: number of decimal digits in value.
	precision := int64(len(strconv.FormatInt(valInt64, 10)))

	// Assign fields to the BigDecimal object.
	setupBasicFields(self, bigIntObj, precision, int64(0))

	return nil
}

// bigdecimalInitString: Set up a BigDecimal object based on a string object argument.
func bigdecimalInitString(params []interface{}) interface{} {
	self := params[0].(*object.Object)   // BigDecimal
	strObj := params[1].(*object.Object) // String

	if !object.IsStringObject(strObj) {
		return getGErrBlk(excNames.IllegalArgumentException, "bigdecimalInitString: argument is not a string")
	}
	str := object.GoStringFromStringObject(strObj)

	// Parse the string into a *big.Int (unscaled value) and a scale.
	bigInt, scale, gerr := parseBigDecimalString(str)
	if gerr != nil {
		return gerr
	}

	// Compute precision: number of decimal digits in unscaled value.
	precision := int64(len(bigInt.Text(10)))

	// Create BigInteger object for field intVal.
	bigIntObj := object.MakeEmptyObjectWithClassName(&classNameBigInteger)
	setBigIntegerFields(bigIntObj, bigInt)

	// Set fields into the BigDecimal object.
	setupBasicFields(self, bigIntObj, precision, scale)

	return nil
}

func bigdecimalInitBigInteger(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	biObj := params[1].(*object.Object)
	bigInt := biObj.FieldTable["value"].Fvalue.(*big.Int)
	precision := int64(len(bigInt.Text(10)))
	scale := int64(0)
	setupBasicFields(self, biObj, precision, scale)

	return nil
}

/*
Helper Functions
*/

func loadStaticsBigDouble() {
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

// Parse a String into decimal precision and scale values.
func parseBigDecimalString(argStr string) (*big.Int, int64, interface{}) {
	// Check for empty string.
	argStr = strings.TrimSpace(argStr)
	if argStr == "" {
		argStr = "0"
	}

	// Set up negative flag.
	negative := false
	if argStr[0] == '+' {
		argStr = argStr[1:]
	} else if argStr[0] == '-' {
		negative = true
		argStr = argStr[1:]
	}

	// wholePart = left of '.' substring.
	// fracPart = right of '.' substring.
	parts := strings.SplitN(argStr, ".", 2)
	wholePart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	// Remove leading zeros from whole part, keeping at least one digit.
	wholePart = strings.TrimLeft(wholePart, "0")
	if wholePart == "" {
		wholePart = "0"
	}

	// Form the precision string.
	precisionStr := wholePart + fracPart
	if precisionStr == "" {
		precisionStr = "0"
	}
	if negative {
		precisionStr = "-" + precisionStr
	}

	// Compute the precision.
	precision := new(big.Int)
	_, ok := precision.SetString(precisionStr, 10)
	if !ok {
		errMsg := fmt.Sprintf("bigdecimalObjectFromBigDecimal: invalid digits detected: %s", argStr)
		return nil, int64(0), getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Compute the scale.
	scale := int64(len(fracPart))

	return precision, scale, nil
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
