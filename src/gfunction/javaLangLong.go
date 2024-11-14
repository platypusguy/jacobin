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
	"math/bits"
	"strconv"
)

func Load_Lang_Long() {

	MethodSignatures["java/lang/Long.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Long.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  longDoubleValue,
		}

	MethodSignatures["java/lang/Long.parseLong(Ljava/lang/String;)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longParseLong,
		}

	MethodSignatures["java/lang/Long.rotateLeft(JI)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longRotateLeft,
		}

	MethodSignatures["java/lang/Long.rotateRight(JI)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longRotateRight,
		}

	MethodSignatures["java/lang/Long.toHexString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToHexString,
		}

	MethodSignatures["java/lang/Long.toString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToString,
		}

	MethodSignatures["java/lang/Long.valueOf(J)Ljava/lang/Long;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longValueOf,
		}

}

// "java/lang/Long.doubleValue()D"
func longDoubleValue(params []interface{}) interface{} {
	var jj int64
	parmObj := params[0].(*object.Object)
	jj = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(jj)
}

// "java/lang/Long.parseLong(Ljava/lang/String;)J"
func longParseLong(params []interface{}) interface{} {
	obj := params[1].(*object.Object)
	str := object.GoStringFromStringObject(obj)
	jj, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("strconv.ParseInt(%s,10,64), failed, reason: %s", str, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return jj
}

// "java/lang/Long.rotateLeft(JI)J"
func longRotateLeft(params []interface{}) interface{} {
	jj := uint64(params[0].(int64))
	shiftLength := int(params[1].(int64))
	value := bits.RotateLeft64(jj, shiftLength)
	return int64(value)
}

// "java/lang/Long.rotateRight(JI)J"
func longRotateRight(params []interface{}) interface{} {
	jj := uint64(params[0].(int64))
	shiftLength := int(params[1].(int64))
	value := bits.RotateLeft64(jj, -shiftLength)
	return int64(value)
}

// "java/lang/Long.valueOf(J)Ljava/lang/Long;"
func longValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Long", types.Long, int64Value)
}

// "java/lang/Long.toHexString(J)Ljava/lang/String;"
func longToHexString(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	uint64Value := uint64(int64Value)
	str := fmt.Sprintf("%016x", uint64Value)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Long.toString(J)Ljava/lang/String;"
func longToString(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	str := fmt.Sprintf("%d", int64Value)
	obj := object.StringObjectFromGoString(str)
	return obj
}
