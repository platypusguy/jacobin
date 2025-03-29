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
	"jacobin/stringPool"
	"jacobin/types"
	"math"
	"strconv"
	"unsafe"
)

func Load_Lang_Double() {

	MethodSignatures["java/lang/Double.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Double.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleByteValue,
		}

	MethodSignatures["java/lang/Double.compare(DD)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  doubleCompare,
		}

	MethodSignatures["java/lang/Double.compareTo(Ljava/lang/Double;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleCompareTo,
		}

	MethodSignatures["java/lang/Double.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleDoubleValue,
		}

	MethodSignatures["java/lang/Double.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleEquals,
		}

	MethodSignatures["java/lang/Double.parseDouble(Ljava/lang/String;)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleParseDouble,
		}

	MethodSignatures["java/lang/Double.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleToString,
		}

	// Native functions or caller to native functions

	MethodSignatures["java/lang/Double.doubleToLongBits(D)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleToLongBits,
		}

	MethodSignatures["java/lang/Double.doubleToRawLongBits(D)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleToRawLongBits,
		}

	MethodSignatures["java/lang/Double.longBitsToDouble(J)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longBitsToDouble,
		}

	MethodSignatures["java/lang/Double.valueOf(Ljava/lang/String;)Ljava/lang/Double;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleValueOf,
		}

	MethodSignatures["java/lang/Double.valueOf(D)Ljava/lang/Double;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleValueOf,
		}

}

// "java/lang/Double.byteValue()B"
func doubleByteValue(params []interface{}) interface{} {
	var dd float64
	parmObj := params[0].(*object.Object)
	dd = parmObj.FieldTable["value"].Fvalue.(float64)
	return int64(byte(dd))
}

// "java/lang/Double.compare(DD)I"
func doubleCompare(params []interface{}) interface{} {
	dd1 := params[0].(float64)
	dd2 := params[1].(float64)
	if dd1 == dd2 {
		return int64(0)
	}
	if dd1 < dd2 {
		return int64(-1)
	}
	return int(1)
}

// "java/lang/Double.compareTo(Ljava/lang/Double;)I"
func doubleCompareTo(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	dd1 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Get the actual Java Double parameter
	parmObj = params[1].(*object.Object)
	dd2 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Now, its just like doubleCompare.
	if dd1 == dd2 {
		return int64(0)
	}
	if dd1 < dd2 {
		return int64(-1)
	}
	return int64(1)
}

// "java/lang/Double.doubleValue()D"
func doubleDoubleValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	return parmObj.FieldTable["value"].Fvalue.(float64)
}

// "java/lang/Double.equals(Ljava/lang/Object;)Z"
func doubleEquals(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	dd1 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Get the actual Java Object parameter
	parmObj = params[1].(*object.Object)
	if parmObj.KlassName == types.InvalidStringIndex {
		return int64(0)
	}
	if *(stringPool.GetStringPointer(parmObj.KlassName)) != "java/lang/Double" {
		return int64(0)
	}
	dd2 = parmObj.FieldTable["value"].Fvalue.(float64)

	// If equal, return true; else return false.
	// fmt.Printf("DEBUG doubleEquals dd1=%f, dd2=%f\n", dd1, dd2)
	if dd1 == dd2 {
		return int64(1)
	}
	return int64(0)
}

// "java/lang/Double.parseDouble(Ljava/lang/String;)D"
func doubleParseDouble(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return getGErrBlk(excNames.NumberFormatException, "doubleParseDouble: String length is zero")
	}

	// Compute output.
	output, err := strconv.ParseFloat(strArg, 64)
	if err != nil {
		errMsg := fmt.Sprintf("doubleParseDouble: strconv.ParseFloat(%s) failed, reason: %s", strArg, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return output

}

// "java/lang/Double.toString()Ljava/lang/String;"
func doubleToString(params []interface{}) interface{} {
	var dd float64
	parmObj := params[0].(*object.Object)
	dd = parmObj.FieldTable["value"].Fvalue.(float64)
	str := fmt.Sprintf("%f", dd)
	objPtr := object.StringObjectFromGoString(str)
	return objPtr
}

// Simulating doubleToRawLongBits in Go
// "java/lang/Double.doubleToRawLongBits(D)J"
func doubleToRawLongBits(params []interface{}) interface{} {
	value := params[0].(float64)
	return *(*int64)(unsafe.Pointer(&value))
}

// Simulating doubleToLongBits in Go
// "java/lang/Double.doubleToLongBits(D)J"
func doubleToLongBits(params []interface{}) interface{} {
	value := params[0].(float64)
	if !math.IsNaN(value) {
		return int64(math.Float64bits(value))
	}
	return 0x7ff8000000000000 // equivalent to Java's 0x7ff8000000000000L
}

// Simulating longBitsToDouble in Go
// "java/lang/Double.longBitsToDouble(J)D"
func longBitsToDouble(params []interface{}) interface{} {
	bits := params[0].(int64)
	return math.Float64frombits(uint64(bits))
}

func doubleValueOf(params []interface{}) interface{} {
	className := "java/lang/Double"
	obj := object.MakeEmptyObjectWithClassName(&className)
	switch params[0].(type) {
	case float64:
		obj.FieldTable["value"] = object.Field{Ftype: types.Double, Fvalue: params[0].(float64)}
	case *object.Object:
		that := params[0].(*object.Object)
		str := object.GoStringFromJavaByteArray(that.FieldTable["value"].Fvalue.([]types.JavaByte))
		if len(str) < 1 {
			return getGErrBlk(excNames.NullPointerException, "doubleValueOf: nil string argument")
		}
		dbl, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return getGErrBlk(excNames.NumberFormatException, "doubleValueOf: "+err.Error())
		}
		obj.FieldTable["value"] = object.Field{Ftype: types.Double, Fvalue: dbl}
	default:
		errMsg := fmt.Sprintf("doubleValueOf: unrecognizable parameter type: %T", params[0])
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return obj
}
