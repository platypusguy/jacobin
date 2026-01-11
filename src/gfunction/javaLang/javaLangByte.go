/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
	"strings"
)

func Load_Lang_Byte() {

	ghelpers.MethodSignatures["java/lang/Byte.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Byte.byteValue()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteIntLongShortByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Byte.compare(BB)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  byteCompare,
		}

	ghelpers.MethodSignatures["java/lang/Byte.compareUnsigned(BB)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  byteCompareUnsigned,
		}

	ghelpers.MethodSignatures["java/lang/Byte.decode(Ljava/lang/String;)Ljava/lang/Byte;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteDecode,
		}

	ghelpers.MethodSignatures["java/lang/Byte.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteFloatDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Byte.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteEquals,
		}

	ghelpers.MethodSignatures["java/lang/Byte.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteFloatDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Byte.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteHashCode,
		}

	ghelpers.MethodSignatures["java/lang/Byte.hashCode(B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteHashCodeStatic,
		}

	ghelpers.MethodSignatures["java/lang/Byte.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteIntLongShortByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Byte.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteIntLongShortByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Byte.parseByte(Ljava/lang/String;)B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteParseByte,
		}

	ghelpers.MethodSignatures["java/lang/Byte.parseByte(Ljava/lang/String;I)B"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  byteParseByteRadix,
		}

	ghelpers.MethodSignatures["java/lang/Byte.shortValue()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteIntLongShortByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Byte.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  byteToString,
		}

	ghelpers.MethodSignatures["java/lang/Byte.toString(B)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteToStringStatic,
		}

	ghelpers.MethodSignatures["java/lang/Byte.toUnsignedInt(B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteToUnsignedInt,
		}

	ghelpers.MethodSignatures["java/lang/Byte.toUnsignedLong(B)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteToUnsignedLong,
		}

	ghelpers.MethodSignatures["java/lang/Byte.valueOf(B)Ljava/lang/Byte;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteValueOf,
		}

	ghelpers.MethodSignatures["java/lang/Byte.valueOf(Ljava/lang/String;)Ljava/lang/Byte;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  byteValueOfString,
		}

	ghelpers.MethodSignatures["java/lang/Byte.valueOf(Ljava/lang/String;I)Ljava/lang/Byte;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  byteValueOfString,
		}

}

var classNameByte = "java/lang/Byte"

// "java/lang/Byte.decode(Ljava/lang/String;)Ljava/lang/Byte;"
func byteDecode(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "byteDecode: String length is zero")
	}

	// This logic should match Integer.decode
	nm := strArg
	radix := 10
	negative := false
	if strings.HasPrefix(nm, "-") {
		negative = true
		nm = nm[1:]
	} else if strings.HasPrefix(nm, "+") {
		nm = nm[1:]
	}

	if strings.HasPrefix(nm, "0x") || strings.HasPrefix(nm, "0X") {
		radix = 16
		nm = nm[2:]
	} else if strings.HasPrefix(nm, "#") {
		radix = 16
		nm = nm[1:]
	} else if strings.HasPrefix(nm, "0") && len(nm) > 1 {
		radix = 8
		nm = nm[1:]
	}

	if negative {
		nm = "-" + nm
	}

	// Parse the input integer.
	int64Value, err := strconv.ParseInt(nm, radix, 64)
	if err != nil {
		errMsg := fmt.Sprintf("byteDecode: strconv.ParseInt(%s,%d) failed, reason: %s", nm, radix, err.Error())
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}

	if int64Value < -128 || int64Value > 127 {
		errMsg := fmt.Sprintf("byteDecode: Value out of range for byte: %d", int64Value)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Create Byte object.
	return object.MakePrimitiveObject(classNameByte, types.Byte, int64Value)
}

// "java/lang/Byte.doubleValue()D"
// "java/lang/Byte.floatValue()F"
func byteFloatDoubleValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	bb := parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(bb)
}

// "java/lang/Byte.byteValue()B"
// "java/lang/Byte.intValue()I"
// "java/lang/Byte.longValue()J"
// "java/lang/Byte.shortValue()S"
func byteIntLongShortByteValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	bb := parmObj.FieldTable["value"].Fvalue.(int64)
	return bb
}

// "java/lang/Byte.compare(BB)I"
func byteCompare(params []interface{}) interface{} {
	x := params[0].(int64)
	y := params[1].(int64)
	return x - y
}

// "java/lang/Byte.compareUnsigned(BB)I"
func byteCompareUnsigned(params []interface{}) interface{} {
	x := uint8(params[0].(int64))
	y := uint8(params[1].(int64))
	if x < y {
		return int64(-1)
	} else if x > y {
		return int64(1)
	}
	return int64(0)
}

// "java/lang/Byte.equals(Ljava/lang/Object;)Z"
func byteEquals(params []interface{}) interface{} {
	byteObj, ok1 := params[0].(*object.Object)
	otherObj, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 || object.IsNull(otherObj) {
		return types.JavaBoolFalse
	}

	if object.GoStringFromStringPoolIndex(otherObj.KlassName) != classNameByte {
		return types.JavaBoolFalse
	}

	byteValue := byteObj.FieldTable["value"].Fvalue.(int64)
	otherValue := otherObj.FieldTable["value"].Fvalue.(int64)

	if byteValue == otherValue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Byte.hashCode()I"
func byteHashCode(params []interface{}) interface{} {
	byteObj := params[0].(*object.Object)
	val := byteObj.FieldTable["value"].Fvalue.(int64)
	return val
}

// "java/lang/Byte.hashCode(B)I"
func byteHashCodeStatic(params []interface{}) interface{} {
	val := params[0].(int64)
	return val
}

// "java/lang/Byte.parseByte(Ljava/lang/String;)B"
func byteParseByte(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(parmObj)
	val, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("byteParseByte: %v", err))
	}
	return val
}

// "java/lang/Byte.parseByte(Ljava/lang/String;I)B"
func byteParseByteRadix(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(parmObj)
	radix := int(params[1].(int64))
	if int64(radix) < ghelpers.MinRadix || int64(radix) > ghelpers.MaxRadix {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("byteParseByteRadix: Invalid radix %d", radix))
	}
	val, err := strconv.ParseInt(str, radix, 8)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("byteParseByteRadix: %v", err))
	}
	return val
}

// "java/lang/Byte.toString()Ljava/lang/String;"
func byteToString(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	str := strconv.FormatInt(ii, 10)
	return object.StringObjectFromGoString(str)
}

// "java/lang/Byte.toString(B)Ljava/lang/String;"
func byteToStringStatic(params []interface{}) interface{} {
	ii := params[0].(int64)
	str := strconv.FormatInt(ii, 10)
	return object.StringObjectFromGoString(str)
}

// "java/lang/Byte.toUnsignedInt(B)I"
func byteToUnsignedInt(params []interface{}) interface{} {
	val := uint8(params[0].(int64))
	return int64(val)
}

// "java/lang/Byte.toUnsignedLong(B)J"
func byteToUnsignedLong(params []interface{}) interface{} {
	val := uint8(params[0].(int64))
	return int64(val)
}

// "java/lang/Byte.valueOf(B)Ljava/lang/Byte;"
func byteValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return object.MakePrimitiveObject(classNameByte, types.Byte, int64Value)
}

// "java/lang/Byte.valueOf(Ljava/lang/String;)Ljava/lang/Byte;"
// "java/lang/Byte.valueOf(Ljava/lang/String;I)Ljava/lang/Byte;"
func byteValueOfString(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(parmObj)
	radix := 10
	if len(params) == 2 {
		radix = int(params[1].(int64))
	}
	val, err := strconv.ParseInt(str, radix, 8)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("byteValueOfString: %v", err))
	}
	return object.MakePrimitiveObject(classNameByte, types.Byte, val)
}
