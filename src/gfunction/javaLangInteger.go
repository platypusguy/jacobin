/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"strconv"
	"strings"
)

func Load_Lang_Integer() map[string]GMeth {

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

	MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  integerParseInt,
		}

	MethodSignatures["java/lang/Integer.valueOf(I)Ljava/lang/Integer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerValueOf,
		}

	return MethodSignatures
}

func integerByteValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return ii
}

func integerDecode(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := string(parmObj.FieldTable["value"].Fvalue.([]byte))
	if len(strArg) < 1 {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerDecode: byte array length < 1")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Parse the input integer.
	int64Value, err := strconv.ParseInt(strArg, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("javaPrimitives.integerDecode: arg=%s, err: %s", strArg, err.Error())
		return getGErrBlk(exceptions.NumberFormatException, errMsg)
	}

	// Create Integer object.
	return populator("java/lang/Integer", types.Int, int64Value)
}

func integerFloatDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(ii)
}

func integerIntLongValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return ii
}

func integerParseInt(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := string(parmObj.FieldTable["value"].Fvalue.([]byte))
	if len(strArg) < 1 {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: string length < 1")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: radix is not an integer")
	}
	rdx := params[1].(int64)
	if rdx < MinRadix || rdx > MaxRadix {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: invalid radix")
	}

	// Compute output.
	output, err := strconv.ParseInt(strArg, int(rdx), 64)
	if err != nil {
		errMsg := fmt.Sprintf("javaPrimitives.integerParseInt: arg=%s, radix=%d, err: %s", strArg, rdx, err.Error())
		return getGErrBlk(exceptions.NumberFormatException, errMsg)
	}

	// Check Integer boundaries.
	if output > MaxIntValue {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: upper limit is Integer.MAX_VALUE")
	}
	if output < MinIntValue {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: lower limit is Integer.MIN_VALUE")
	}

	// Return computed value.
	return output
}

func integerValueOf(params []interface{}) interface{} {
	// fmt.Printf("DEBUG integerValueOf at entry params[0]: (%T) %v\n", params[0], params[0])
	int64Value := params[0].(int64)
	return populator("java/lang/Integer", types.Int, int64Value)
}
