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

func Load_Lang_Integer() {

	MethodSignatures["java/lang/Integer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Integer.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerByteValue,
		}

	MethodSignatures["java/lang/Integer.decode(Ljava/lang/String;)Ljava/lang/Integer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerDecode,
		}

	MethodSignatures["java/lang/Integer.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerFloatDoubleValue,
		}

	MethodSignatures["java/lang/Integer.floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerFloatDoubleValue,
		}

	MethodSignatures["java/lang/Integer.intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerIntLongValue,
		}

	MethodSignatures["java/lang/Integer.longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerIntLongValue,
		}

	MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  integerParseInt,
		}

	MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  integerParseIntRadix,
		}

	MethodSignatures["java/lang/Integer.valueOf(I)Ljava/lang/Integer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerValueOf,
		}

	MethodSignatures["java/lang/Integer.toBinaryString(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerToBinaryString,
		}

	MethodSignatures["java/lang/Integer.toHexString(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerToHexString,
		}

	MethodSignatures["java/lang/Integer.toOctalString(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerToOctalString,
		}

	MethodSignatures["java/lang/Integer.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerToString,
		}

	MethodSignatures["java/lang/Integer.toString(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerToStringI,
		}

	MethodSignatures["java/lang/Integer.toUnsignedString(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerToUnsignedString,
		}

	MethodSignatures["java/lang/Integer.toUnsignedString(II)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  integerToUnsignedStringRadix,
		}

}

// "java/lang/Integer.byteValue()B"
func integerByteValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return ii
}

// "java/lang/Integer.decode(Ljava/lang/String;)Ljava/lang/Integer;"
func integerDecode(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return getGErrBlk(excNames.NumberFormatException, "Byte array length is zero")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Parse the input integer.
	int64Value, err := strconv.ParseInt(strArg, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("strconv.ParseInt(%s,10,64) failed, failed, reason: %s", strArg, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Create Integer object.
	return populator("java/lang/Integer", types.Int, int64Value)
}

// "java/lang/Integer.doubleValue()D"
// "java/lang/Integer.floatValue()F"
func integerFloatDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(ii)
}

// "java/lang/Integer.longValue()J"
func integerIntLongValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return ii
}

// "java/lang/Integer.parseInt(Ljava/lang/String;)I"
// Radix = 10
func integerParseInt(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return getGErrBlk(excNames.NumberFormatException, "String length is zero")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Compute output.
	output, err := strconv.ParseInt(strArg, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("strconv.ParseInt(%s,10,64) failed, reason: %s", strArg, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Return computed value.
	return output
}

// "java/lang/Integer.parseInt(Ljava/lang/String;I)I"
func integerParseIntRadix(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return getGErrBlk(excNames.NumberFormatException, "String length is zero")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		return getGErrBlk(excNames.NumberFormatException, "Radix is not an integer")
	}
	rdx := params[1].(int64)
	if rdx < MinRadix || rdx > MaxRadix {
		errMsg := fmt.Sprintf("Invalid radix value (%d)", rdx)
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Compute output.
	output, err := strconv.ParseInt(strArg, int(rdx), 64)
	if err != nil {
		errMsg := fmt.Sprintf("strconv.ParseInt(%s,%d,64) failed, reason: %s", strArg, rdx, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Check Integer boundaries.
	if output > MaxIntValue {
		return getGErrBlk(excNames.NumberFormatException, "Computed integer exceeds upper limit")
	}
	if output < MinIntValue {
		return getGErrBlk(excNames.NumberFormatException, "Computed integer is less than lower limit")
	}

	// Return computed value.
	return output
}

// "java/lang/Integer.valueOf(I)Ljava/lang/Integer;"
func integerValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Integer", types.Int, int64Value)
}

// "java/lang/Integer.toString()Ljava/lang/String;"
func integerToString(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	argInt64 := obj1.FieldTable["value"].Fvalue.(int64)
	str := fmt.Sprintf("%d", argInt64)
	obj2 := object.StringObjectFromGoString(str)
	return obj2
}

// "java/lang/Integer.toString(I)Ljava/lang/String;"
func integerToStringI(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	str := fmt.Sprintf("%d", argInt64)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toUnsignedString(I)Ljava/lang/String;"
func integerToUnsignedString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	if argInt64 < 0 {
		argInt64 &= 0x00000000FFFFFFFF
	}
	str := fmt.Sprintf("%d", argInt64)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toUnsignedString(II)Ljava/lang/String;"
func integerToUnsignedStringRadix(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	if argInt64 < 0 {
		argInt64 &= 0x00000000FFFFFFFF
	}
	// fmt.Printf("DEBUG integerToUnsignedStringRadix %d - %08x\n", argInt64, argInt64)

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		errMsg := fmt.Sprintf("Invalid radix (%v) format", params[1])
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	rdx := params[1].(int64)
	if rdx < MinRadix || rdx > MaxRadix {
		errMsg := fmt.Sprintf("Invalid radix value (%d)", rdx)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	str := strconv.FormatInt(argInt64, int(rdx))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toBinaryString(I)Ljava/lang/String;"
func integerToBinaryString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	str := strconv.FormatInt(argInt64, 2)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toOctalString(I)Ljava/lang/String;"
func integerToOctalString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	str := strconv.FormatInt(argInt64, 8)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toHexString(I)Ljava/lang/String;"
func integerToHexString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	str := strconv.FormatInt(argInt64, 16)
	obj := object.StringObjectFromGoString(str)
	return obj
}
