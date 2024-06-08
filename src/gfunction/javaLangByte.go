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
	"strconv"
	"strings"
)

func Load_Lang_Byte() {

	MethodSignatures["java/lang/Byte.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Byte.decode(Ljava/lang/String;)Ljava/lang/Byte;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  byteDecode,
		}

	MethodSignatures["java/lang/Byte.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  byteDoubleValue,
		}

	MethodSignatures["java/lang/Byte.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  byteToString,
		}

	MethodSignatures["java/lang/Byte.valueOf(B)Ljava/lang/Byte;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  byteValueOf,
		}

}

// "java/lang/Byte.decode(Ljava/lang/String;)Ljava/lang/Byte;"
func byteDecode(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return getGErrBlk(excNames.NumberFormatException, "Byte array length < 1")
	}

	// Strip off a leading "#" or "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "", 1)
	}
	if strings.HasPrefix(strArg, "0x") {
		strArg = strings.Replace(strArg, "0x", "", 1)
	}

	// Parse the input integer.
	int64Value, err := strconv.ParseInt(strArg, 16, 64)
	if err != nil {
		errMsg := fmt.Sprintf("strconv.ParseInt(%s) failed, reason: %s", strArg, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	if int64Value > 255 {
		errMsg := fmt.Sprintf("Byte value too large: %d", int64Value)
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Create Byte object.
	return populator("java/lang/Byte", types.Byte, int64Value)
}

// "java/lang/Byte.doubleValue()D"
func byteDoubleValue(params []interface{}) interface{} {
	var bb int64
	parmObj := params[0].(*object.Object)
	bb = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(bb)
}

// "java/lang/Byte.toString()Ljava/lang/String;"
func byteToString(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	str := fmt.Sprintf("%d", ii)
	outObjPtr := object.StringObjectFromGoString(str)
	return outObjPtr
}

// "java/lang/Byte.valueOf(B)Ljava/lang/Byte;"
func byteValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Byte", types.Byte, int64Value)
}
