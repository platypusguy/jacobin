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
	bigInt, scale := parseDecimalString(valStr)

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
	bigInt, scale := parseDecimalString(str)

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
func parseDecimalString(s string) (*big.Int, int64) {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	// Remove leading zeros in intPart to avoid incorrect precision
	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" {
		intPart = "0"
	}

	scale := int64(len(fracPart))
	fullDigits := intPart + fracPart
	fullDigits = strings.TrimLeft(fullDigits, "0")

	if fullDigits == "" {
		fullDigits = "0"
		scale = 0
	}

	bi := new(big.Int)
	bi.SetString(fullDigits, 10)

	if neg {
		bi.Neg(bi)
	}

	return bi, scale
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
