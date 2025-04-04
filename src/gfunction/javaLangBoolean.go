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

	// The JVM library Boolean.compare and Boolean.compareTo functions work.

	MethodSignatures["java/lang/Boolean.describeConstable()Ljava.util.Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// The JVM library Boolean.equals function works.

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
			GFunction:  booleanHashCode,
		}

	// The JVM library Boolean.logicalAnd function works.
	// The JVM library Boolean.logicalOr function works.
	// The JVM library Boolean.logicalXor function works.

	MethodSignatures["java/lang/Boolean.parseBoolean(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanParseBoolean,
		}

	MethodSignatures["java/lang/Boolean.parseBoolean(Ljava/lang/String;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanParseBoolean,
		}

	// The JVM library Boolean.toString function works.

	MethodSignatures["java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Ljava/lang/String;)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
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
// If neither true nor false, return an exception.
func booleanGetBoolean(params []interface{}) interface{} {

	// Get the property name and validate it.
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := fmt.Sprintf("booleanGetBoolean: Boolean parameter is neither nil or an invalid object: {%T, %v}",
			params[0], params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fld, ok := obj.FieldTable["value"]
	if !ok {
		errMsg := "booleanGetBoolean: Missing the \"value\" field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var propName string
	switch fld.Fvalue.(type) {
	case []byte:
		propName = string(obj.FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		propName = object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	default:
		errMsg := fmt.Sprintf("booleanGetBoolean: The \"value\" field has invalid type: %T", fld.Ftype)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Call systemGetProperty() to get the property value.
	var sysParams []interface{}
	sysParams = append(sysParams, object.StringObjectFromGoString(propName))
	propObj := systemGetProperty(sysParams).(*object.Object)
	propStr := object.GoStringFromStringObject(propObj)
	switch propStr {
	case "true":
		return types.JavaBoolTrue
	case "false":
		return types.JavaBoolFalse
	}

	// Neither true nor false ==> exception.
	errMsg := fmt.Sprintf("booleanGetBoolean: systemGetProperty returned neither \"true\" nor \"false\": %v", propStr)
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)
}

// Returns a Boolean object instance, based on the parameter type: boolean or String.
func booleanValueOf(params []interface{}) interface{} {

	// Try for a String object.
	inObj, ok := params[0].(*object.Object)
	if !ok {
		inBool, ok := params[0].(int64)
		if !ok {
			errMsg := fmt.Sprintf("booleanValueOf: The parameter is neither String nor boolean: %T", params[0])
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		obj := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, inBool)
		return obj
	}

	// inObj should be a String object.
	zz := _booleanStringParser(inObj)
	switch zz {
	case types.JavaBoolTrue, types.JavaBoolFalse:
		return object.MakePrimitiveObject("java/lang/Boolean", types.Bool, zz)
	}

	// Return exception.
	return zz
}

// Given a Boolean object, return one of the following:
// types.JavaBoolTrue
// types.JavaBoolFalse
// an exception
func _booleanStringParser(obj *object.Object) interface{} {
	// Get the String argument value.
	fld, ok := obj.FieldTable["value"]
	if !ok {
		errMsg := "_booleanStringParser: Missing the \"value\" field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var strValue string
	switch fld.Fvalue.(type) {
	case []byte:
		strValue = string(obj.FieldTable["value"].Fvalue.([]byte))
	case []types.JavaByte:
		strValue = object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte))
	default:
		errMsg := fmt.Sprintf("_booleanStringParser: The \"value\" field has invalid type: %T", fld.Ftype)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	switch strValue {
	case "true":
		return types.JavaBoolTrue
	case "false":
		return types.JavaBoolFalse
	}

	// Neither true nor false ==> exception.
	errMsg := fmt.Sprintf("_booleanStringParser: The \"value\" field is neither \"true\" nor \"false\": %v", strValue)
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)
}

// Parse the string argument as a boolean and return true, false, or an exception.
func booleanParseBoolean(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	return _booleanStringParser(obj)
}

// Return a hash code for this Boolean object.
func booleanHashCode(params []interface{}) interface{} {
	var zz interface{}
	if len(params) == 0 {
		obj := params[0].(*object.Object)
		zz = _booleanStringParser(obj)
	} else {
		zz = params[1].(int64)
		switch zz {
		case types.JavaBoolTrue, types.JavaBoolFalse:
		default:
			errMsg := fmt.Sprintf("booleanHashCode: The argument is neither \"true\" nor \"false\": %v", zz)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}
	switch zz {
	case types.JavaBoolTrue:
		return int64(1231)
	case types.JavaBoolFalse:
		return int64(1237)
	}

	// Return exception from _booleanStringParser.
	return zz
}
