/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/bits"
	"strconv"
)

func Load_Lang_Long() {

	MethodSignatures["java/lang/Long.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Long.bitCount(J)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longBitCount,
		}

	MethodSignatures["java/lang/Long.compare(JJ)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longCompare,
		}

	MethodSignatures["java/lang/Long.compareUnsigned(JJ)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longCompareUnsigned,
		}

	MethodSignatures["java/lang/Long.compress(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longCompress,
		}

	MethodSignatures["java/lang/Long.divideUnsigned(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longDivideUnsigned,
		}

	MethodSignatures["java/lang/Long.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  longDoubleValue,
		}

	MethodSignatures["java/lang/Long.expand(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longExpand,
		}

	MethodSignatures["java/lang/Long.highestOneBit(J)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longHighestOneBit,
		}

	MethodSignatures["java/lang/Long.lowestOneBit(J)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longLowestOneBit,
		}

	MethodSignatures["java/lang/Long.max(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longMax,
		}

	MethodSignatures["java/lang/Long.min(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longMin,
		}

	MethodSignatures["java/lang/Long.numberOfLeadingZeros(J)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longNumberOfLeadingZeros,
		}

	MethodSignatures["java/lang/Long.numberOfTrailingZeros(J)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longNumberOfTrailingZeros,
		}

	MethodSignatures["java/lang/Long.parseLong(Ljava/lang/String;)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longParseLong,
		}

	MethodSignatures["java/lang/Long.remainderUnsigned(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longRemainderUnsigned,
		}

	MethodSignatures["java/lang/Long.reverse(J)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longReverse,
		}

	MethodSignatures["java/lang/Long.reverseBytes(J)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longReverseBytes,
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

	MethodSignatures["java/lang/Long.signum(J)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longSignum,
		}

	MethodSignatures["java/lang/Long.sum(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longSum,
		}

	MethodSignatures["java/lang/Long.toBinaryString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToBinaryString,
		}

	MethodSignatures["java/lang/Long.toHexString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToHexString,
		}

	MethodSignatures["java/lang/Long.toOctalString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToOctalString,
		}

	MethodSignatures["java/lang/Long.toString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToString,
		}

	MethodSignatures["java/lang/Long.toUnsignedString(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  longToUnsignedString,
		}

	MethodSignatures["java/lang/Long.toUnsignedString(JI)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longToUnsignedStringRadix,
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
	obj := params[0].(*object.Object)
	str := object.GoStringFromStringObject(obj)
	jj, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("longParseLong: strconv.ParseInt(%s,10,64), failed, reason: %s", str, err.Error())
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
	return object.MakePrimitiveObject("java/lang/Long", types.Long, int64Value)
}

// "java/lang/Long.toHexString(J)Ljava/lang/String;"
func longToHexString(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	uint64Value := uint64(int64Value)
	str := strconv.FormatUint(uint64Value, 16)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Long.toOctalString(J)Ljava/lang/String;"
func longToOctalString(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	uint64Value := uint64(int64Value)
	str := strconv.FormatUint(uint64Value, 8)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Long.toBinaryString(J)Ljava/lang/String;"
func longToBinaryString(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	uint64Value := uint64(int64Value)
	str := strconv.FormatUint(uint64Value, 2)
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

// "java/lang/Long.toUnsignedString(J)Ljava/lang/String;"
func longToUnsignedString(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	val := uint64(argInt64)
	str := fmt.Sprintf("%d", val)
	obj := object.StringObjectFromGoString(str)
	return obj
}

// "java/lang/Long.toUnsignedString(JI)Ljava/lang/String;"
func longToUnsignedStringRadix(params []interface{}) interface{} {
	argInt64 := params[0].(int64)
	val := uint64(argInt64)

	// Extract and validate the radix.
	switch params[1].(type) {
	case int64:
	default:
		errMsg := fmt.Sprintf("longToUnsignedStringRadix: Invalid radix (%v) format", params[1])
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	rdx := params[1].(int64)
	if rdx < MinRadix || rdx > MaxRadix {
		errMsg := fmt.Sprintf("longToUnsignedStringRadix: Invalid radix value (%d)", rdx)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	str := strconv.FormatUint(val, int(rdx))
	obj := object.StringObjectFromGoString(str)
	return obj
}

// longBitCount returns the number of one-bits in the two’s complement binary representation of a long.
func longBitCount(params []interface{}) interface{} {
	input, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "longBitCount: Invalid argument type")
	}
	return int64(bits.OnesCount64(uint64(input)))
}

// longCompare compares two long values numerically.
func longCompare(params []interface{}) interface{} {
	inputA, ok1 := params[0].(int64)
	inputB, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "longCompare: Invalid argument types")
	}
	if inputA < inputB {
		return int64(-1)
	} else if inputA > inputB {
		return int64(1)
	}
	return int64(0)
}

// longCompareUnsigned compares two longs as unsigned values.
func longCompareUnsigned(params []interface{}) interface{} {
	x, ok1 := params[0].(int64)
	y, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "longCompareUnsigned: Invalid argument types")
	}

	ux, uy := uint64(x), uint64(y)
	if ux < uy {
		return int64(-1)
	} else if ux > uy {
		return int64(1)
	}
	return int64(0)
}

// longDivideUnsigned performs unsigned long division.
func longDivideUnsigned(params []interface{}) interface{} {
	dividend, ok1 := params[0].(int64)
	divisor, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "longDivideUnsigned: Invalid argument types")
	}
	if divisor == 0 {
		return getGErrBlk(excNames.ArithmeticException, "longDivideUnsigned: Division by zero")
	}
	return int64(uint64(dividend) / uint64(divisor))
}

// longRemainderUnsigned returns the remainder of dividing two unsigned longs.
func longRemainderUnsigned(params []interface{}) interface{} {
	dividend, ok1 := params[0].(int64)
	divisor, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "longRemainderUnsigned: Invalid argument types")
	}
	if divisor == 0 {
		return getGErrBlk(excNames.ArithmeticException, "longRemainderUnsigned: Division by zero")
	}
	return int64(uint64(dividend) % uint64(divisor))
}

// longHighestOneBit returns a long value with at most a single one-bit, in the position of the highest-order one-bit in the specified long value.
func longHighestOneBit(params []interface{}) interface{} {
	input, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "longHighestOneBit: Invalid argument type")
	}
	if input == 0 {
		return int64(0)
	}
	return int64(uint64(1) << (63 - bits.LeadingZeros64(uint64(input))))
}

// longLowestOneBit returns a long value with at most a single one-bit, in the position of the lowest-order one-bit in the specified long value.
func longLowestOneBit(params []interface{}) interface{} {
	input, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "longLowestOneBit: Invalid argument type")
	}
	return input & -input
}

// longNumberOfLeadingZeros returns the number of zero bits preceding the highest-order ("leftmost") one-bit in the two's complement binary representation of the specified long value.
func longNumberOfLeadingZeros(params []interface{}) interface{} {
	arg := uint64(params[0].(int64))
	return int64(bits.LeadingZeros64(arg))
}

// longNumberOfTrailingZeros returns the number of zero bits following the lowest-order ("rightmost") one-bit in the two's complement binary representation of the specified long value.
func longNumberOfTrailingZeros(params []interface{}) interface{} {
	arg := uint64(params[0].(int64))
	return int64(bits.TrailingZeros64(arg))
}

// longReverse returns the value obtained by reversing the order of the bits in the two’s complement binary representation of the specified long value.
func longReverse(params []interface{}) interface{} {
	i, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "longReverse: Invalid argument type")
	}
	return int64(bits.Reverse64(uint64(i)))
}

// longReverseBytes returns the value obtained by reversing the order of the bytes in the two’s complement representation of the specified long value.
func longReverseBytes(params []interface{}) interface{} {
	i, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "longReverseBytes: Invalid argument type")
	}
	return int64(bits.ReverseBytes64(uint64(i)))
}

// longSignum returns the signum function of the specified long value.
func longSignum(params []interface{}) interface{} {
	val := params[0].(int64)
	switch {
	case val < 0:
		return int64(-1)
	case val > 0:
		return int64(1)
	default:
		return int64(0)
	}
}

// longSum returns the sum of two longs.
func longSum(params []interface{}) interface{} {
	a := params[0].(int64)
	b := params[1].(int64)
	return a + b
}

// longMax returns the greater of two long values.
func longMax(params []interface{}) interface{} {
	a := params[0].(int64)
	b := params[1].(int64)
	if a > b {
		return a
	}
	return b
}

// longMin returns the smaller of two long values.
func longMin(params []interface{}) interface{} {
	a := params[0].(int64)
	b := params[1].(int64)
	if a < b {
		return a
	}
	return b
}

// longCompress extracts bits from the input using the provided mask.
func longCompress(params []interface{}) interface{} {
	inputRaw, ok1 := params[0].(int64)
	maskRaw, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "longCompress: Invalid argument types")
	}

	input := uint64(inputRaw)
	mask := uint64(maskRaw)
	result := uint64(0)
	pos := uint32(0)
	for mask != 0 {
		if mask&1 != 0 {
			result |= (input & 1) << pos
			pos++
		}
		mask >>= 1
		input >>= 1
	}

	return int64(result)
}

// longExpand expands bits from the input using the provided mask.
func longExpand(params []interface{}) interface{} {
	inputRaw, ok1 := params[0].(int64)
	maskRaw, ok2 := params[1].(int64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "longExpand: Invalid argument types")
	}

	input := uint64(inputRaw)
	mask := uint64(maskRaw)
	result := uint64(0)
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

	return int64(result)
}
