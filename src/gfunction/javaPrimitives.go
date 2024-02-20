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
	"unicode"
)

// Implementation some of the functions in Byte, Character, Integer, Long, Short, and Boolean.

// Radix boundaries:
var minRadix int64 = 2
var maxRadix int64 = 36
var MaxIntValue int64 = 2147483647
var MinIntValue int64 = -2147483648

func Load_Primitives() map[string]GMeth {

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

	MethodSignatures["java/lang/Character.valueOf(C)Ljava/lang/Character;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  characterValueOf,
		}

	MethodSignatures["java/lang/Character.charValue()C"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  charValue,
		}

	MethodSignatures["java/lang/Character.isLetter(C)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charIsLetter,
		}

	MethodSignatures["java/lang/Character.isDigit(C)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charIsDigit,
		}

	MethodSignatures["java/lang/Character.toLowerCase(C)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charToLowerCase,
		}

	MethodSignatures["java/lang/Character.toUpperCase(C)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charToUpperCase,
		}

	MethodSignatures["java/lang/Double.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleByteValue,
		}

	MethodSignatures["java/lang/Double.compare(DD)I"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  doubleCompare,
		}

	MethodSignatures["java/lang/Double.compareTo(Ljava/lang/Double;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleCompareTo,
		}

	MethodSignatures["java/lang/Double.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleEquals,
		}

	MethodSignatures["java/lang/Double.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleToString,
		}

	MethodSignatures["java/lang/Double.parseDouble(Ljava/lang/String;)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleParseDouble,
		}

	MethodSignatures["java/lang/Double.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleDoubleValue,
		}

	MethodSignatures["java/lang/Integer.valueOf(I)Ljava/lang/Integer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerValueOf,
		}

	MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  integerParseInt,
		}

	MethodSignatures["java/lang/Integer.decode(Ljava/lang/String;)Ljava/lang/Integer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerDecode,
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

	MethodSignatures["java/lang/Integer.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  integerByteValue,
		}

	MethodSignatures["java/lang/Long.valueOf(J)Ljava/lang/Long;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longValueOf,
		}

	MethodSignatures["java/lang/Long.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  longDoubleValue,
		}

	MethodSignatures["java/lang/Short.valueOf(S)Ljava/lang/Short;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  shortValueOf,
		}

	MethodSignatures["java/lang/Short.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  shortDoubleValue,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
		}

	MethodSignatures["java/lang/Byte.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Boolean.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Integer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Long.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Float.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Double.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Short.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	return MethodSignatures
}

func byteToString(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ii = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ii = parmObj.Fields[0].Fvalue.(int64)
	}
	str := fmt.Sprintf("%d", ii)
	objPtr := object.CreateCompactStringFromGoString(&str)
	return objPtr
}

func byteValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Byte", types.Byte, int64Value)
}

func byteDoubleValue(params []interface{}) interface{} {
	var bb int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		bb = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		bb = parmObj.Fields[0].Fvalue.(int64)
	}

	return float64(bb)
}

func byteDecode(params []interface{}) interface{} {
	var bptr *[]byte
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		bptr = parmObj.FieldTable["value"].Fvalue.(*[]byte)
	} else {
		bptr = parmObj.Fields[0].Fvalue.(*[]byte)
	}

	// Validate byte array.
	if bptr == nil {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.byteDecode: Nil byte array pointer")
	}
	strArg := string(*bptr)
	if len(strArg) < 1 {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.byteDecode: byte array length < 1")
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
		errMsg := fmt.Sprintf("javaPrimitives.byteDecode: arg=%s, err: %s", strArg, err.Error())
		return getGErrBlk(exceptions.NumberFormatException, errMsg)
	}
	if int64Value > 255 {
		errMsg := fmt.Sprintf("javaPrimitives.byteDecode: value too large: %d", int64Value)
		return getGErrBlk(exceptions.NumberFormatException, errMsg)
	}

	// Create Byte object.
	return populator("java/lang/Byte", types.Byte, int64Value)
}

func characterValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Character", types.Char, int64Value)
}

func charValue(params []interface{}) interface{} {
	var ch int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ch = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ch = parmObj.Fields[0].Fvalue.(int64)
	}
	return ch
}

func charIsLetter(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsLetter(rune(ii)) {
		return int64(1)
	}
	return int64(0)
}

func charIsDigit(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsDigit(rune(ii)) {
		return int64(1)
	}
	return int64(0)
}

func charToLowerCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	rr := unicode.ToLower(rune(ii))
	return int64(rr)
}

func charToUpperCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	rr := unicode.ToUpper(rune(ii))
	return int64(rr)
}

func doubleByteValue(params []interface{}) interface{} {
	var dd float64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		dd = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd = parmObj.Fields[0].Fvalue.(float64)
	}
	return int64(byte(dd))
}

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

func doubleCompareTo(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		dd1 = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd1 = parmObj.Fields[0].Fvalue.(float64)
	}

	// Get the actual Java Double parameter
	parmObj = params[1].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		dd2 = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd2 = parmObj.Fields[0].Fvalue.(float64)
	}

	// Now, its just like doubleCompare.
	if dd1 == dd2 {
		return int64(0)
	}
	if dd1 < dd2 {
		return int64(-1)
	}
	return int64(1)
}

func doubleEquals(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		dd1 = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd1 = parmObj.Fields[0].Fvalue.(float64)
	}

	// Get the actual Java Object parameter
	parmObj = params[1].(*object.Object)
	if parmObj.Klass == nil {
		return int64(0)
	}
	// fmt.Printf("DEBUG doubleEquals Klass --> %s\n", *parmObj.Klass)
	if *parmObj.Klass != "java/lang/Double" {
		return int64(0)
	}
	if len(parmObj.FieldTable) > 0 {
		dd2 = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd2 = parmObj.Fields[0].Fvalue.(float64)
	}

	// If equal, return true; else return false.
	// fmt.Printf("DEBUG doubleEquals dd1=%f, dd2=%f\n", dd1, dd2)
	if dd1 == dd2 {
		return int64(1)
	}
	return int64(0)
}

func doubleToString(params []interface{}) interface{} {
	var dd float64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		dd = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd = parmObj.Fields[0].Fvalue.(float64)
	}
	str := fmt.Sprintf("%f", dd)
	objPtr := object.CreateCompactStringFromGoString(&str)
	return objPtr
}

func doubleParseDouble(params []interface{}) interface{} {
	var bptr *[]byte
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		bptr = parmObj.FieldTable["value"].Fvalue.(*[]byte)
	} else {
		bptr = parmObj.Fields[0].Fvalue.(*[]byte)
	}
	if bptr == nil {
		parmObj.DumpObject("javaPrimitives.doubleParseDouble: Nil byte array pointer", 0)
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.doubleParseDouble: Nil byte array pointer")
	}
	strArg := string(*bptr)
	if len(strArg) < 1 {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.doubleParseDouble: string length < 1")
	}

	// Compute output.
	output, err := strconv.ParseFloat(strArg, 64)
	if err != nil {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.doubleParseDouble Error(): "+err.Error())
	}
	return output

}

func doubleDoubleValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		return parmObj.FieldTable["value"].Fvalue.(float64)
	}
	return parmObj.Fields[0].Fvalue.(float64)
}

func integerValueOf(params []interface{}) interface{} {
	// fmt.Printf("DEBUG integerValueOf at entry params[0]: (%T) %v\n", params[0], params[0])
	int64Value := params[0].(int64)
	return populator("java/lang/Integer", types.Int, int64Value)
}

func integerDecode(params []interface{}) interface{} {
	var bptr *[]byte
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		bptr = parmObj.FieldTable["value"].Fvalue.(*[]byte)
	} else {
		bptr = parmObj.Fields[0].Fvalue.(*[]byte)
	}

	// Validate byte array.
	if bptr == nil {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerDecode: Nil byte array pointer")
	}
	strArg := string(*bptr)
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

func integerParseInt(params []interface{}) interface{} {
	var bptr *[]byte
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		bptr = parmObj.FieldTable["value"].Fvalue.(*[]byte)
	} else {
		bptr = parmObj.Fields[0].Fvalue.(*[]byte)
	}
	if bptr == nil {
		return getGErrBlk(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: Nil byte array pointer")
	}
	strArg := string(*bptr)
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
	if rdx < minRadix || rdx > maxRadix {
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

func integerIntLongValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ii = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ii = parmObj.Fields[0].Fvalue.(int64)
	}

	return ii
}

func integerFloatDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ii = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ii = parmObj.Fields[0].Fvalue.(int64)
	}

	return float64(ii)
}

func integerByteValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ii = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ii = parmObj.Fields[0].Fvalue.(int64)
	}

	return ii
}

func longValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Long", types.Long, int64Value)
}

func longDoubleValue(params []interface{}) interface{} {
	var jj int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		jj = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		jj = parmObj.Fields[0].Fvalue.(int64)
	}
	return float64(jj)
}

func shortValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Short", types.Short, int64Value)
}

func shortDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ii = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ii = parmObj.Fields[0].Fvalue.(int64)
	}

	return float64(ii)
}

func booleanValueOf(params []interface{}) interface{} {
	zz := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, zz)
	return objPtr
}
