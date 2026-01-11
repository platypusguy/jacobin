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
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/bits"
	"strconv"
	"strings"
)

func Load_Lang_Integer() {

	ghelpers.MethodSignatures["java/lang/Integer.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Integer.bitCount(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerBitCount,
		}

	ghelpers.MethodSignatures["java/lang/Integer.byteValue()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Integer.compare(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerCompare,
		}

	ghelpers.MethodSignatures["java/lang/Integer.compareTo(Ljava/lang/Integer;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerCompareTo,
		}

	ghelpers.MethodSignatures["java/lang/Integer.compareUnsigned(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerCompareUnsigned,
		}

	ghelpers.MethodSignatures["java/lang/Integer.compress(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerCompress,
		}

	ghelpers.MethodSignatures["java/lang/Integer.decode(Ljava/lang/String;)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerDecode,
		}

	ghelpers.MethodSignatures["java/lang/Integer.describeConstable()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Integer.divideUnsigned(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerDivideUnsigned,
		}

	ghelpers.MethodSignatures["java/lang/Integer.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerFloatDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Integer.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerEquals,
		}

	ghelpers.MethodSignatures["java/lang/Integer.expand(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerExpand,
		}

	ghelpers.MethodSignatures["java/lang/Integer.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerFloatDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Integer.getInteger(Ljava/lang/String;)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerGetInteger,
		}

	ghelpers.MethodSignatures["java/lang/Integer.getInteger(Ljava/lang/String;I)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerGetInteger,
		}

	ghelpers.MethodSignatures["java/lang/Integer.getInteger(Ljava/lang/String;Ljava/lang/Integer;)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerGetInteger,
		}

	ghelpers.MethodSignatures["java/lang/Integer.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Integer.hashCode(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Integer.highestOneBit(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerHighestOneBit,
		}

	ghelpers.MethodSignatures["java/lang/Integer.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerIntLongShortValue,
		}

	ghelpers.MethodSignatures["java/lang/Integer.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerIntLongShortValue,
		}

	ghelpers.MethodSignatures["java/lang/Integer.lowestOneBit(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerLowestOneBit,
		}

	ghelpers.MethodSignatures["java/lang/Integer.max(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerMax,
		}

	ghelpers.MethodSignatures["java/lang/Integer.min(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerMin,
		}

	ghelpers.MethodSignatures["java/lang/Integer.numberOfLeadingZeros(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerNumberOfLeadingZeros,
		}

	ghelpers.MethodSignatures["java/lang/Integer.numberOfTrailingZeros(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerNumberOfTrailingZeros,
		}

	ghelpers.MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/CharSequence;III)I"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  integerParseIntCharSequence,
		}

	ghelpers.MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerParseInt,
		}

	ghelpers.MethodSignatures["java/lang/Integer.parseInt(Ljava/lang/String;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerParseIntRadix,
		}

	ghelpers.MethodSignatures["java/lang/Integer.parseUnsignedInt(Ljava/lang/CharSequence;III)I"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  integerParseUnsignedIntCharSequence,
		}

	ghelpers.MethodSignatures["java/lang/Integer.parseUnsignedInt(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerParseUnsignedInt,
		}

	ghelpers.MethodSignatures["java/lang/Integer.parseUnsignedInt(Ljava/lang/String;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerParseUnsignedInt,
		}

	ghelpers.MethodSignatures["java/lang/Integer.remainderUnsigned(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerRemainderUnsigned,
		}

	ghelpers.MethodSignatures["java/lang/Integer.resolveConstantDesc(Ljava/lang/invoke/MethodHandles/Lookup;)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Integer.reverse(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerReverse,
		}

	ghelpers.MethodSignatures["java/lang/Integer.reverseBytes(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerReverseBytes,
		}

	ghelpers.MethodSignatures["java/lang/Integer.rotateLeft(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerRotateLeft,
		}

	ghelpers.MethodSignatures["java/lang/Integer.rotateRight(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerRotateRight,
		}

	ghelpers.MethodSignatures["java/lang/Integer.shortValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerIntLongShortValue,
		}

	ghelpers.MethodSignatures["java/lang/Integer.signum(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerSignum,
		}

	ghelpers.MethodSignatures["java/lang/Integer.sum(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerSum,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toBinaryString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerToBinaryString,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toHexString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerToHexString,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toOctalString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerToOctalString,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  integerToString,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerToStringIorII,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toString(II)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerToStringIorII,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toUnsignedLong(I)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerToUnsignedLong,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toUnsignedString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerToUnsignedString,
		}

	ghelpers.MethodSignatures["java/lang/Integer.toUnsignedString(II)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerToUnsignedStringRadix,
		}

	ghelpers.MethodSignatures["java/lang/Integer.valueOf(I)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerValueOfInt,
		}

	ghelpers.MethodSignatures["java/lang/Integer.valueOf(Ljava/lang/String;)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  integerValueOfString,
		}

	ghelpers.MethodSignatures["java/lang/Integer.valueOf(Ljava/lang/String;I)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  integerValueOfString,
		}

}

var classNameInteger = "java/lang/Integer"

// "java/lang/Integer.byteValue()B"
func integerByteValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = (parmObj.FieldTable["value"].Fvalue.(int64)) & 0xFF
	return ii
}

// "java/lang/Integer.decode(Ljava/lang/String;)Ljava/lang/Integer;"
func integerDecode(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GoStringFromStringObject(parmObj)
	if len(strArg) < 1 {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "integerDecode: Byte array length is zero")
	}

	// Replace a leading "#" with "0x" in strArg.
	wbase := 10
	if strings.HasPrefix(strArg, "#") {
		wbase = 16
		strArg = strArg[1:]
	}

	// Parse the input integer.
	int64Value, err := strconv.ParseInt(strArg, wbase, 64)
	if err != nil {
		errMsg := fmt.Sprintf("integerDecode: strconv.ParseInt(%s,10,64) failed, failed, reason: %s", strArg, err.Error())
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Create Integer object.
	return object.MakePrimitiveObject("java/lang/Integer", types.Int, int64Value)
}

// "java/lang/Integer.doubleValue()D"
// "java/lang/Integer.floatValue()F"
func integerFloatDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(ii)
}

// "java/lang/Integer.intValue()J"
// "java/lang/Integer.longValue()J"
// "java/lang/Integer.shortValue()J"
func integerIntLongShortValue(params []interface{}) interface{} {
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
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "integerParseInt: String length is zero")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Compute output.
	output, err := strconv.ParseInt(strArg, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("integerParseInt: strconv.ParseInt(%s,10,64) failed, reason: %s", strArg, err.Error())
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
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
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "integerParseIntRadix: String length is zero")
	}

	// Replace a leading "#" with "0x" in strArg.
	if strings.HasPrefix(strArg, "#") {
		strArg = strings.Replace(strArg, "#", "0x", 1)
	}

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "integerParseIntRadix: Radix is not an integer")
	}
	rdx := params[1].(int64)
	if rdx < ghelpers.MinRadix || rdx > ghelpers.MaxRadix {
		errMsg := fmt.Sprintf("integerParseIntRadix: Invalid radix value (%d)", rdx)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Compute output.
	output, err := strconv.ParseInt(strArg, int(rdx), 64)
	if err != nil {
		errMsg := fmt.Sprintf("integerParseIntRadix: strconv.ParseInt(%s,%d,64) failed, reason: %s", strArg, rdx, err.Error())
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Check Integer boundaries.
	if output > ghelpers.MaxIntValue {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "integerParseIntRadix: Computed integer exceeds upper limit")
	}
	if output < ghelpers.MinIntValue {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, "integerParseIntRadix: Computed integer is less than lower limit")
	}

	// Return computed value.
	return output
}

// "java/lang/Integer.signum(I)I"
func integerSignum(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	switch {
	case int64Value < 0:
		return int64(-1)
	case int64Value > 0:
		return int64(+1)
	default:
		return int64(0)
	}
}

// "java/lang/Integer.toString()Ljava/lang/String;"
func integerToString(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	argInt64 := obj1.FieldTable["value"].Fvalue.(int64)
	str := fmt.Sprintf("%d", argInt64)
	obj2 := object.StringObjectFromGoString(str)
	return obj2
}

// integerToStringIorII returns a string representation of the integer, optionally in the specified radix.
func integerToStringIorII(params []interface{}) interface{} {
	if len(params) < 1 || len(params) > 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerToStringIorII requires 1 or 2 arguments")
	}

	input, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerToStringIorII: First argument must be an int64")
	}

	radix := 10
	if len(params) == 2 {
		rr, ok := params[1].(int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerToStringIorII: Second argument must be an int64 representing the radix")
		}
		radix = int(rr)
		if radix < 2 || radix > 36 {
			return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("integerToStringIorII: Radix out of range: %d", radix))
		}
	}

	str := strconv.FormatInt(input, radix)
	return object.StringObjectFromGoString(str)
}

// "java/lang/Integer.toUnsignedString(I)Ljava/lang/String;"
func integerToUnsignedString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	val := uint32(argInt64)
	str := fmt.Sprintf("%d", val)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toUnsignedString(II)Ljava/lang/String;"
func integerToUnsignedStringRadix(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	val := uint32(argInt64)
	// fmt.Printf("DEBUG integerToUnsignedStringRadix %d - %08x\n", argInt64, argInt64)

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		errMsg := fmt.Sprintf("integerToUnsignedStringRadix: Invalid radix (%v) format", params[1])
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}
	rdx := params[1].(int64)
	if rdx < ghelpers.MinRadix || rdx > ghelpers.MaxRadix {
		errMsg := fmt.Sprintf("integerToUnsignedStringRadix: Invalid radix value (%d)", rdx)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	str := strconv.FormatUint(uint64(val), int(rdx))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toOctalString(I)Ljava/lang/String;"
func integerToOctalString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	str := strconv.FormatUint(uint64(uint32(argInt64)), 8)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.toHexString(I)Ljava/lang/String;"
func integerToHexString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	str := strconv.FormatUint(uint64(uint32(argInt64)), 16)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Integer.numberOfTrailingZeros(I)I"
func integerNumberOfTrailingZeros(params []interface{}) interface{} {
	arg := uint32(params[0].(int64))
	return int64(bits.TrailingZeros32(arg))
}

// "java/lang/Integer.numberOfLeadingZeros(I)I"
func integerNumberOfLeadingZeros(params []interface{}) interface{} {
	arg := uint32(params[0].(int64))
	return int64(bits.LeadingZeros32(arg))
}

// RotateLeft performs a left bitwise rotation on an integer.
func integerRotateLeft(params []interface{}) interface{} {

	input, ok1 := params[0].(int64)
	distance, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerRotateLeft: Invalid argument types")
	}

	uInput := uint32(input)
	uDistance := uint(distance & 0x1F)
	result := (uInput << uDistance) | (uInput >> (32 - uDistance))
	return int64(int32(result))
}

// RotateRight performs a right bitwise rotation on an integer.
func integerRotateRight(params []interface{}) interface{} {

	input, ok1 := params[0].(int64)
	distance, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerRotateRight: Invalid argument types")
	}

	uInput := uint32(input)
	uDistance := uint(distance & 0x1F)
	result := (uInput >> uDistance) | (uInput << (32 - uDistance))
	return int64(int32(result))
}

// BitCount returns the number of one-bits in the two’s complement binary representation of an integer.
func integerBitCount(params []interface{}) interface{} {

	input, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerBitCount: Invalid argument type")
	}

	return int64(bits.OnesCount32(uint32(input)))
}

// Compare two integer values numerically.
// Return 0 if x == y; return less than 0 if x < y; and return a value greater than 0 if x > y
func integerCompare(params []interface{}) interface{} {

	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompare requires exactly 2 arguments")
	}
	inputA, ok1 := params[0].(int64)
	inputB, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompare: Invalid argument types")
	}

	return inputA - inputB
}

// CompareUnsigned compares two integers as unsigned values.
func integerCompareUnsigned(params []interface{}) interface{} {

	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompareUnsigned requires exactly 2 arguments")
	}
	x, ok1 := params[0].(int64)
	y, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompareUnsigned: Invalid argument types")
	}

	ux, uy := uint32(x), uint32(y)
	if ux < uy {
		return int64(-1)
	} else if ux > uy {
		return int64(1)
	}
	return int64(0)
}

// integerCompress extracts bits from the input using the provided mask.
func integerCompress(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompress requires exactly 2 arguments")
	}

	inputRaw, ok1 := params[0].(int64)
	maskRaw, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompress: Invalid argument types for Compress")
	}

	input := uint32(inputRaw)
	mask := uint32(maskRaw)
	result := uint32(0)
	pos := uint32(0)
	for mask != 0 {
		if mask&1 != 0 {
			result |= (input & 1) << pos
			pos++
		}
		mask >>= 1
		input >>= 1
	}

	return int64(int32(result))
}

// integerDivideUnsigned performs unsigned integer division.
func integerDivideUnsigned(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerDivideUnsigned requires exactly 2 arguments")
	}

	dividend, ok1 := params[0].(int64)
	divisor, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerDivideUnsigned: Invalid argument types for DivideUnsigned")
	}
	if divisor == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "integerDivideUnsigned: Division by zero")
	}

	return int64(uint32(dividend) / uint32(divisor))
}

// integerEquals checks if an Integer object is equal to another object.
func integerEquals(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerEquals requires exactly 2 arguments")
	}

	integerObj, ok1 := params[0].(*object.Object)
	otherObj, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerEquals: Invalid argument types for Equals")
	}

	integerValue, exists := integerObj.FieldTable["value"]
	if !exists || integerValue.Ftype != types.Int {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerEquals: Invalid Integer object structure")
	}

	otherValue, exists := otherObj.FieldTable["value"]
	if !exists || otherValue.Ftype != types.Int {
		return types.JavaBoolFalse
	}

	if integerValue.Fvalue == otherValue.Fvalue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// integerExpand expands bits from the input using the provided mask.
func integerExpand(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerExpand requires exactly 2 arguments")
	}

	inputRaw, ok1 := params[0].(int64)
	maskRaw, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerExpand: Invalid argument types for integerExpand")
	}

	input := uint32(inputRaw)
	mask := uint32(maskRaw)
	result := uint32(0)
	pos := uint32(0)
	for mask != 0 {
		if mask&1 != 0 {
			if input&1 != 0 {
				result |= 1 << pos
			}
			input >>= 1
		}
		mask >>= 1
		pos++
	}

	return int64(int32(result))
}

// integerGetInteger retrieves the Integer object based on different types of input.
func integerGetInteger(params []interface{}) interface{} {
	if len(params) == 0 || len(params) > 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerGetInteger requires 1 to 3 arguments")
	}

	var name string
	var defaultValue int64
	hasDefault := false

	// Get the property name.
	nameObj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(nameObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerGetInteger: First parameter must be a Java String object")
	}
	name = object.GoStringFromStringObject(nameObj)
	if len(name) < 1 {
		return object.Null
	}

	// More than one parameter?
	if len(params) > 1 {
		// Try for a primitive integer default value.
		defaultValue, ok = params[1].(int64)
		if !ok {
			// Try for an integer object default value.
			thatObj, ok := params[1].(*object.Object)
			if !ok || object.GoStringFromStringPoolIndex(thatObj.KlassName) != classNameInteger {
				return object.Null
			}
			defaultValue, ok = thatObj.FieldTable["value"].Fvalue.(int64)
			if !ok {
				return object.Null
			}
		}
		hasDefault = true
	}

	// Get the System.getProperty(name) value.
	value := globals.GetSystemProperty(name)
	if value == "" {
		// If no system property by that name is available, return the default value if available.
		if hasDefault {
			return defaultValue
		}
		return object.Null
	}

	// Numeric?
	numeric, err := strconv.Atoi(value)
	if err != nil {
		return object.Null
	}

	// It is a numeric.
	return object.MakePrimitiveObject(classNameInteger, types.Int, numeric)
}

// integerHighestOneBit returns an int value with at most a single one-bit, in the position of the highest-order one-bit in the specified int value.
func integerHighestOneBit(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerHighestOneBit requires exactly 1 argument")
	}
	input, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerHighestOneBit: Invalid argument type")
	}

	if input == 0 {
		return int64(0)
	}
	return int64(1 << (31 - bits.LeadingZeros32(uint32(input))))
}

// integerLowestOneBit returns an int value with at most a single one-bit, in the position of the lowest-order one-bit in the specified int value.
func integerLowestOneBit(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerLowestOneBit requires exactly 1 argument")
	}
	input, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerLowestOneBit: Invalid argument type")
	}

	return int64(int32(input) & -int32(input))
}

// integerMax If A > B return A else return B.
func integerMax(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerMax requires exactly 2 arguments")
	}
	A, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerMax: Invalid left argument type")
	}
	B, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerMax: Invalid right argument type")
	}

	if A > B {
		return A
	}
	return B
}

// integerMin If A < B return A else return B.
func integerMin(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerMin requires exactly 2 arguments")
	}
	A, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerMin: Invalid left argument type")
	}
	B, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerMin: Invalid right argument type")
	}

	if A < B {
		return A
	}
	return B
}

// integerParseUnsignedInt parses the string argument as an unsigned integer in base 10 or the specified radix.
func integerParseUnsignedInt(params []interface{}) interface{} {
	if len(params) < 1 || len(params) > 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseUnsignedInt requires 1 or 2 arguments")
	}

	strObj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseUnsignedInt: First parameter must be a Java String object")
	}
	str := object.GoStringFromStringObject(strObj)

	radix := 10
	if len(params) == 2 {
		rr, ok := params[1].(int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseUnsignedInt: Second parameter must be an int64 representing the radix")
		}
		radix = int(rr)
	}

	value, err := strconv.ParseUint(str, radix, 32)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("integerParseUnsignedInt: Invalid unsigned integer: %v", err))
	}

	return int64(value)
}

// integerRemainderUnsigned returns the remainder of dividing two unsigned integers.
// integerParseIntCharSequence parses a CharSequence as an integer.
func integerParseIntCharSequence(params []interface{}) interface{} {
	if len(params) != 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseIntCharSequence requires exactly 4 arguments")
	}

	csObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseIntCharSequence: First parameter must be an object")
	}

	begin, ok1 := params[1].(int64)
	end, ok2 := params[2].(int64)
	radix, ok3 := params[3].(int64)
	if !ok1 || !ok2 || !ok3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseIntCharSequence: Invalid numeric parameters")
	}

	str := object.GoStringFromStringObject(csObj)
	if begin < 0 || end > int64(len(str)) || begin > end {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "integerParseIntCharSequence: bounds out of range")
	}

	subStr := str[begin:end]
	output, err := strconv.ParseInt(subStr, int(radix), 64)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("integerParseIntCharSequence: Invalid integer: %v", err))
	}

	return output
}

// integerParseUnsignedIntCharSequence parses a CharSequence as an unsigned integer.
func integerParseUnsignedIntCharSequence(params []interface{}) interface{} {
	if len(params) != 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseUnsignedIntCharSequence requires exactly 4 arguments")
	}

	csObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseUnsignedIntCharSequence: First parameter must be an object")
	}

	begin, ok1 := params[1].(int64)
	end, ok2 := params[2].(int64)
	radix, ok3 := params[3].(int64)
	if !ok1 || !ok2 || !ok3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerParseUnsignedIntCharSequence: Invalid numeric parameters")
	}

	str := object.GoStringFromStringObject(csObj)
	if begin < 0 || end > int64(len(str)) || begin > end {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "integerParseUnsignedIntCharSequence: bounds out of range")
	}

	subStr := str[begin:end]
	output, err := strconv.ParseUint(subStr, int(radix), 32)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("integerParseUnsignedIntCharSequence: Invalid unsigned integer: %v", err))
	}

	return int64(output)
}

func integerRemainderUnsigned(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerRemainderUnsigned requires exactly 2 arguments")
	}

	dividend, ok1 := params[0].(int64)
	divisor, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerRemainderUnsigned: Invalid argument types")
	}
	if divisor == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "integerRemainderUnsigned: Division by zero")
	}

	return int64(uint32(dividend) % uint32(divisor))
}

// integerReverse returns the value obtained by reversing the order of the bits in the two’s complement binary representation of the specified int value.
func integerReverse(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerReverse requires exactly 1 argument")
	}
	i, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerReverse: Invalid argument type")
	}

	return int64(bits.Reverse32(uint32(i)))
}

// integerReverseBytes returns the value obtained by reversing the order of the bytes in the two’s complement representation of the specified int value.
func integerReverseBytes(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerReverseBytes requires exactly 1 argument")
	}
	i, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerReverseBytes: Invalid argument type")
	}

	return int64(bits.ReverseBytes32(uint32(i)))
}

// integerSum returns the sum of two integers.
func integerSum(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerSum requires exactly 2 arguments")
	}
	a, ok1 := params[0].(int64)
	b, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerSum: Invalid argument types for integerSum")
	}

	return a + b
}

// integerToBinaryString returns a string representation of the unsigned integer value in binary (base 2).
func integerToBinaryString(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerToBinaryString requires exactly 1 argument")
	}
	input, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid argument type for integerToBinaryString")
	}

	binaryStr := strconv.FormatUint(uint64(uint32(input)), 2)
	return object.StringObjectFromGoString(binaryStr)
}

// integerToUnsignedLong converts the argument to a long by an unsigned conversion.
func integerToUnsignedLong(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerToUnsignedLong requires exactly 1 argument")
	}
	input, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerToUnsignedLong: Invalid argument type")
	}

	return int64(uint64(uint32(input)))
}

// "java/lang/Integer.valueOf(I)Ljava/lang/Integer;"
func integerValueOfInt(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return object.MakePrimitiveObject("java/lang/Integer", types.Int, int64Value)
}

// integerValueOf returns an Integer object for the specified string,
// parsing it as an integer using the specified radix if one is supplied.
func integerValueOfString(params []interface{}) interface{} {
	if len(params) < 1 || len(params) > 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerValueOfString requires 1 or 2 arguments")
	}

	strObj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerValueOfString: First parameter must be a Java String object")
	}
	str := object.GoStringFromStringObject(strObj)

	// Default radix is 10
	radix := 10
	if len(params) == 2 {
		rr, ok := params[1].(int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerValueOfString: Second parameter must be an int64 representing the radix")
		}
		radix = int(rr)
	}

	// Parse the string as an integer with the specified radix
	value, err := strconv.ParseInt(str, radix, 32)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, fmt.Sprintf("integerValueOfString: Invalid radix(%d): %v", radix, err))
	}

	// Create and return an Integer object
	return object.MakePrimitiveObject(classNameInteger, types.Int, value)
}

// integerCompareTo compares two Integer objects.
func integerCompareTo(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompareTo requires exactly 2 arguments")
	}

	thisObj, ok1 := params[0].(*object.Object)
	otherObj, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompareTo: Arguments must be Integer objects")
	}

	thisVal, exists1 := thisObj.FieldTable["value"]
	otherVal, exists2 := otherObj.FieldTable["value"]
	if !exists1 || thisVal.Ftype != types.Int || !exists2 || otherVal.Ftype != types.Int {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "integerCompareTo: Invalid Integer object structure")
	}

	x := thisVal.Fvalue.(int64)
	y := otherVal.Fvalue.(int64)

	switch {
	case x < y:
		return int64(-1)
	case x > y:
		return int64(1)
	default:
		return int64(0)
	}
}
