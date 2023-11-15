/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"strconv"
	"unicode"
)

// Implementation some of the functions in Byte, Character, Integer, Long, Short, and Boolean.

// Radix boundaries:
var minRadix int64 = 2
var maxRadix int64 = 36

func Load_Primitives() map[string]GMeth {

	MethodSignatures["java/lang/Byte.valueOf(B)Ljava/lang/Byte;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  byteValueOf,
		}

	MethodSignatures["java/lang/Byte.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			ObjectRef:  true,
			GFunction:  byteDoubleValue,
		}

	MethodSignatures["java/lang/Character.valueOf(C)Ljava/lang/Character;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  characterValueOf,
		}

	MethodSignatures["java/lang/Character.charValue()C"] =
		GMeth{
			ParamSlots: 0,
			ObjectRef:  true,
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
			ObjectRef:  true,
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
			ObjectRef:  true,
			GFunction:  doubleCompareTo,
		}

	MethodSignatures["java/lang/Double.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			ObjectRef:  true,
			GFunction:  doubleEquals,
		}

	MethodSignatures["java/lang/Double.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			ObjectRef:  true,
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
			ObjectRef:  true,
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

	MethodSignatures["java/lang/Integer.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			ObjectRef:  true,
			GFunction:  integerDoubleValue,
		}

	MethodSignatures["java/lang/Long.valueOf(J)Ljava/lang/Long;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longValueOf,
		}

	MethodSignatures["java/lang/Long.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			ObjectRef:  true,
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
			ObjectRef:  true,
			GFunction:  shortDoubleValue,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
		}

	MethodSignatures["java/lang/Boolean.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  booleanJustReturn,
		}

	return MethodSignatures
}

func byteValueOf(params []interface{}) interface{} {
	bb := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Byte", types.Byte, bb)
	return objPtr
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

func characterValueOf(params []interface{}) interface{} {
	cc := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Character", types.Char, cc)
	return objPtr
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
	//fmt.Printf("DEBUG doubleEquals Klass --> %s\n", *parmObj.Klass)
	if *parmObj.Klass != "java/lang/Double" {
		return int64(0)
	}
	if len(parmObj.FieldTable) > 0 {
		dd2 = parmObj.FieldTable["value"].Fvalue.(float64)
	} else {
		dd2 = parmObj.Fields[0].Fvalue.(float64)
	}

	// If equal, return true; else return false.
	//fmt.Printf("DEBUG doubleEquals dd1=%f, dd2=%f\n", dd1, dd2)
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
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.doubleParseDouble: Nil byte array pointer")
	}
	strArg := string(*bptr)
	if len(strArg) < 1 {
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.doubleParseDouble: string length < 1")
	}

	// Compute output.
	output, err := strconv.ParseFloat(strArg, 64)
	if err != nil {
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.doubleParseDouble Error(): "+err.Error())
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
	//fmt.Printf("DEBUG integerValueOf at entry params[0]: (%T) %v\n", params[0], params[0])
	ii := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Integer", types.Int, ii)
	return objPtr
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
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: Nil byte array pointer")
	}
	strArg := string(*bptr)
	if len(strArg) < 1 {
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: string length < 1")
	}

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: radix is not an integer")
	}
	rdx := params[1].(int64)
	if rdx < minRadix || rdx > maxRadix {
		exceptions.Throw(exceptions.NumberFormatException, "javaPrimitives.integerParseInt: invalid radix")
	}

	// Compute output.
	output, err := strconv.ParseInt(strArg, int(rdx), 64)
	if err != nil {
		errMsg := fmt.Sprintf("javaPrimitives.integerParseInt: arg=%s, radix=%d, err: %s", strArg, rdx, err.Error())
		exceptions.Throw(exceptions.NumberFormatException, errMsg)
	}
	return output
}

func integerDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		ii = parmObj.FieldTable["value"].Fvalue.(int64)
	} else {
		ii = parmObj.Fields[0].Fvalue.(int64)
	}

	return float64(ii)
}

func longValueOf(params []interface{}) interface{} {
	jj := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Long", types.Long, jj)
	return objPtr
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
	ss := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Short", types.Short, ss)
	return objPtr
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

func booleanJustReturn(params []interface{}) interface{} {
	return nil
}
