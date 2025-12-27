/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

func Load_Lang_Boolean() {

	MethodSignatures["java/lang/Boolean.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Boolean.<init>(Z)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Boolean.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Boolean.booleanValue()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  booleanBooleanValue,
		}

	MethodSignatures["java/lang/Boolean.compare(ZZ)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  booleanCompare,
		}

	MethodSignatures["java/lang/Boolean.compareTo(Ljava/lang/Boolean;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanCompareTo,
		}

	MethodSignatures["java/lang/Boolean.describeConstable()Ljava.util.Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Boolean.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanEquals,
		}

	MethodSignatures["java/lang/Boolean.getBoolean(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanGetBoolean,
		}

	MethodSignatures["java/lang/Boolean.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  booleanHashCode,
		}

	MethodSignatures["java/lang/Boolean.hashCode(Z)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanHashCodeStatic,
		}

	MethodSignatures["java/lang/Boolean.logicalAnd(ZZ)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  booleanLogicalAnd,
		}

	MethodSignatures["java/lang/Boolean.logicalOr(ZZ)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  booleanLogicalOr,
		}

	MethodSignatures["java/lang/Boolean.logicalXor(ZZ)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  booleanLogicalXor,
		}

	MethodSignatures["java/lang/Boolean.parseBoolean(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanParseBoolean,
		}

	MethodSignatures["java/lang/Boolean.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  booleanToString,
		}

	MethodSignatures["java/lang/Boolean.toString(Z)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanToStringStatic,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Ljava/lang/String;)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOfString,
		}

}

// Return the value of this Boolean object as a boolean primitive.
func booleanBooleanValue(params []interface{}) interface{} {
	// Try for a Boolean object.
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := fmt.Sprintf("booleanBooleanValue: Boolean parameter is neither nil or an invalid object: {%T, %v}",
			params[0], params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the value field.
	fld, ok := obj.FieldTable["value"]
	if !ok {
		errMsg := "booleanBooleanValue: Missing the \"value\" field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var zz int64
	switch fld.Fvalue.(type) {
	case int64:
		zz = fld.Fvalue.(int64)
	default:
		errMsg := fmt.Sprintf("booleanBooleanValue: The \"value\" field has the wrong type: %T", fld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// inObj should be a String object.
	switch zz {
	case types.JavaBoolTrue, types.JavaBoolFalse:
		return zz
	}

	// Return exception.
	errMsg := fmt.Sprintf("booleanBooleanValue: The \"value\" field value is neither true nor false: %d", zz)
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)

}

// Returns true if and only if the system property named by the argument exists
// and is equal to, ignoring case, the string "true".
// Else, return false.
func booleanGetBoolean(params []interface{}) interface{} {
	// Get the property name and validate it.
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil || !object.IsStringObject(obj) {
		return types.JavaBoolFalse
	}

	propName := object.GoStringFromStringObject(obj)
	if propName == "" {
		return types.JavaBoolFalse
	}

	// Call systemGetProperty() to get the property value.
	propValue := globals.GetSystemProperty(propName)
	if propValue == "" {
		return types.JavaBoolFalse
	}

	if strings.EqualFold(propValue, "true") {
		return types.JavaBoolTrue
	}

	return types.JavaBoolFalse
}

// Returns a Boolean object instance, based on the parameter: boolean.
func booleanValueOf(params []interface{}) interface{} {
	inBool, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("booleanValueOf: expected int64 (boolean), got %T", params[0]))
	}
	return object.MakePrimitiveObject("java/lang/Boolean", types.Bool, inBool)
}

// booleanValueOfString returns a Boolean object instance, based on the parameter: String.
func booleanValueOfString(params []interface{}) interface{} {
	strObj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		return getGErrBlk(excNames.IllegalArgumentException, "booleanValueOfString: expected String object")
	}

	strValue := object.GoStringFromStringObject(strObj)
	if strings.EqualFold(strValue, "true") {
		return object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolTrue)
	}

	return object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolFalse)
}

// booleanParseBoolean parses the string argument as a boolean.
func booleanParseBoolean(params []interface{}) interface{} {
	strObj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		return types.JavaBoolFalse
	}

	strValue := object.GoStringFromStringObject(strObj)
	if strings.EqualFold(strValue, "true") {
		return types.JavaBoolTrue
	}

	return types.JavaBoolFalse
}

// booleanCompare compares two boolean values.
func booleanCompare(params []interface{}) interface{} {
	x := params[0].(int64)
	y := params[1].(int64)

	if x == y {
		return int64(0)
	}
	if x == types.JavaBoolTrue {
		return int64(1)
	}
	return int64(-1)
}

// booleanCompareTo compares this Boolean instance with another.
func booleanCompareTo(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)

	thisVal := thisObj.FieldTable["value"].Fvalue.(int64)
	otherVal := otherObj.FieldTable["value"].Fvalue.(int64)

	return booleanCompare([]interface{}{thisVal, otherVal})
}

// booleanEquals returns true if the specified object is a Boolean and has the same value.
func booleanEquals(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj, ok := params[1].(*object.Object)
	if !ok || otherObj == nil {
		return types.JavaBoolFalse
	}

	if object.GoStringFromStringPoolIndex(otherObj.KlassName) != "java/lang/Boolean" {
		return types.JavaBoolFalse
	}

	thisVal := thisObj.FieldTable["value"].Fvalue.(int64)
	otherVal := otherObj.FieldTable["value"].Fvalue.(int64)

	if thisVal == otherVal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// booleanHashCode returns a hash code for this Boolean object.
func booleanHashCode(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	val := obj.FieldTable["value"].Fvalue.(int64)
	return booleanHashCodeStatic([]interface{}{val})
}

// booleanHashCodeStatic returns a hash code for a boolean value.
func booleanHashCodeStatic(params []interface{}) interface{} {
	val := params[0].(int64)
	if val == types.JavaBoolTrue {
		return int64(1231)
	}
	return int64(1237)
}

// booleanLogicalAnd returns the result of applying the logical AND operator to the specified boolean operands.
func booleanLogicalAnd(params []interface{}) interface{} {
	a := params[0].(int64)
	b := params[1].(int64)
	if a == types.JavaBoolTrue && b == types.JavaBoolTrue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// booleanLogicalOr returns the result of applying the logical OR operator to the specified boolean operands.
func booleanLogicalOr(params []interface{}) interface{} {
	a := params[0].(int64)
	b := params[1].(int64)
	if a == types.JavaBoolTrue || b == types.JavaBoolTrue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// booleanLogicalXor returns the result of applying the logical XOR operator to the specified boolean operands.
func booleanLogicalXor(params []interface{}) interface{} {
	a := params[0].(int64)
	b := params[1].(int64)
	if a != b {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// booleanToString returns a String object representing this Boolean's value.
func booleanToString(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	val := obj.FieldTable["value"].Fvalue.(int64)
	return booleanToStringStatic([]interface{}{val})
}

// booleanToStringStatic returns a String object representing the specified boolean.
func booleanToStringStatic(params []interface{}) interface{} {
	val := params[0].(int64)
	if val == types.JavaBoolTrue {
		return object.StringObjectFromGoString("true")
	}
	return object.StringObjectFromGoString("false")
}
