/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
)

func Load_Lang_Short() {

	ghelpers.MethodSignatures["java/lang/Short.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Short.byteValue()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Short.compare(SS)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  shortCompare,
		}

	ghelpers.MethodSignatures["java/lang/Short.compareUnsigned(SS)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  shortCompareUnsigned,
		}

	ghelpers.MethodSignatures["java/lang/Short.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Short.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortEquals,
		}

	ghelpers.MethodSignatures["java/lang/Short.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortFloatValue,
		}

	ghelpers.MethodSignatures["java/lang/Short.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortIntValue,
		}

	ghelpers.MethodSignatures["java/lang/Short.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortLongValue,
		}

	ghelpers.MethodSignatures["java/lang/Short.parseShort(Ljava/lang/String;)S"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortParseShort,
		}

	ghelpers.MethodSignatures["java/lang/Short.parseShort(Ljava/lang/String;I)S"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  shortParseShortRadix,
		}

	ghelpers.MethodSignatures["java/lang/Short.reverseBytes(S)S"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortReverseBytes,
		}

	ghelpers.MethodSignatures["java/lang/Short.shortValue()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortShortValue,
		}

	ghelpers.MethodSignatures["java/lang/Short.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  shortToString,
		}

	ghelpers.MethodSignatures["java/lang/Short.toString(S)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortToStringS,
		}

	ghelpers.MethodSignatures["java/lang/Short.toUnsignedInt(S)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortToUnsignedInt,
		}

	ghelpers.MethodSignatures["java/lang/Short.toUnsignedLong(S)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortToUnsignedLong,
		}

	ghelpers.MethodSignatures["java/lang/Short.valueOf(S)Ljava/lang/Short;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  shortValueOf,
		}

}

// "java/lang/Short.byteValue()B"
func shortByteValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	return int64(int8(ii))
}

// "java/lang/Short.intValue()I"
func shortIntValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	return int64(int32(ii))
}

// "java/lang/Short.longValue()J"
func shortLongValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	return ii
}

// "java/lang/Short.shortValue()S"
func shortShortValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	return int64(int16(ii))
}

// "java/lang/Short.equals(Ljava/lang/Object;)Z"
func shortEquals(params []interface{}) interface{} {
	shortObj, ok1 := params[0].(*object.Object)
	otherObj, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 {
		return types.JavaBoolFalse
	}

	shortValue, exists1 := shortObj.FieldTable["value"]
	otherValue, exists2 := otherObj.FieldTable["value"]

	if !exists1 || shortValue.Ftype != types.Short || !exists2 || otherValue.Ftype != types.Short {
		return types.JavaBoolFalse
	}

	if shortValue.Fvalue == otherValue.Fvalue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Short.floatValue()F"
func shortFloatValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(ii)
}

// "java/lang/Short.doubleValue()D"
func shortDoubleValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(ii)
}

// "java/lang/Short.compare(SS)I"
func shortCompare(params []interface{}) interface{} {
	x := int16(params[0].(int64)) // interpret as signed short
	y := int16(params[1].(int64))
	return int64(int32(x) - int32(y)) // return difference exactly like HotSpot
}

// "java/lang/Short.compareUnsigned(SS)I"
func shortCompareUnsigned(params []interface{}) interface{} {
	x := uint16(params[0].(int64)) // interpret as unsigned short
	y := uint16(params[1].(int64))
	return int64(int32(x) - int32(y))
}

// "java/lang/Short.parseShort(Ljava/lang/String;)S"
func shortParseShort(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	output, err := strconv.ParseInt(strArg, 10, 16)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, err.Error())
	}
	return output
}

// "java/lang/Short.parseShort(Ljava/lang/String;I)S"
func shortParseShortRadix(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	radix := int(params[1].(int64))
	output, err := strconv.ParseInt(strArg, radix, 16)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, err.Error())
	}
	return output
}

// "java/lang/Short.reverseBytes(S)S"
func shortReverseBytes(params []interface{}) interface{} {
	i := uint16(params[0].(int64))
	res := (i << 8) | (i >> 8)
	return int64(int16(res))
}

// "java/lang/Short.toString()Ljava/lang/String;"
func shortToString(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	ii := parmObj.FieldTable["value"].Fvalue.(int64)
	str := strconv.FormatInt(ii, 10)
	return object.StringObjectFromGoString(str)
}

// "java/lang/Short.toString(S)Ljava/lang/String;"
func shortToStringS(params []interface{}) interface{} {
	ii := params[0].(int64)
	str := strconv.FormatInt(ii, 10)
	return object.StringObjectFromGoString(str)
}

// "java/lang/Short.toUnsignedInt(S)I"
func shortToUnsignedInt(params []interface{}) interface{} {
	i := uint16(params[0].(int64))
	return int64(i)
}

// "java/lang/Short.toUnsignedLong(S)J"
func shortToUnsignedLong(params []interface{}) interface{} {
	i := uint16(params[0].(int64))
	return int64(i)
}

// "java/lang/Short.valueOf(S)Ljava/lang/Short;"
func shortValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return object.MakePrimitiveObject("java/lang/Short", types.Short, int64Value)
}
