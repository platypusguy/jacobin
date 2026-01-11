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
	"math"
	"strconv"
)

func Load_Lang_Float() {

	ghelpers.MethodSignatures["java/lang/Float.<init>(F)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/lang/Float.byteValue()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Float.compare(FF)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  floatCompare,
		}

	ghelpers.MethodSignatures["java/lang/Float.compareTo(Ljava/lang/Float;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatCompareTo,
		}

	ghelpers.MethodSignatures["java/lang/Float.describeConstable()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Float.floatToIntBits(F)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatFloatToIntBits,
		}

	ghelpers.MethodSignatures["java/lang/Float.floatToRawIntBits(F)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatFloatToRawIntBits,
		}

	ghelpers.MethodSignatures["java/lang/Float.float16ToFloat(S)F"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatFloat16ToFloat,
		}

	ghelpers.MethodSignatures["java/lang/Float.floatToFloat16(F)S"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatFloatToFloat16,
		}

	ghelpers.MethodSignatures["java/lang/Float.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Float.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatEquals,
		}

	ghelpers.MethodSignatures["java/lang/Float.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatFloatValue,
		}

	ghelpers.MethodSignatures["java/lang/Float.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatHashCode,
		}

	ghelpers.MethodSignatures["java/lang/Float.hashCode(F)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatHashCodeStatic,
		}

	ghelpers.MethodSignatures["java/lang/Float.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatIntValue,
		}

	ghelpers.MethodSignatures["java/lang/Float.isFinite(F)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatIsFiniteStatic,
		}

	ghelpers.MethodSignatures["java/lang/Float.isInfinite()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatIsInfinite,
		}

	ghelpers.MethodSignatures["java/lang/Float.isInfinite(F)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatIsInfiniteStatic,
		}

	ghelpers.MethodSignatures["java/lang/Float.isNaN()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatIsNaN,
		}

	ghelpers.MethodSignatures["java/lang/Float.isNaN(F)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatIsNaNStatic,
		}

	ghelpers.MethodSignatures["java/lang/Float.intBitsToFloat(I)F"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatIntBitsToFloat,
		}

	ghelpers.MethodSignatures["java/lang/Float.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatLongValue,
		}

	ghelpers.MethodSignatures["java/lang/Float.max(FF)F"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  floatMax,
		}

	ghelpers.MethodSignatures["java/lang/Float.min(FF)F"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  floatMin,
		}

	ghelpers.MethodSignatures["java/lang/Float.parseFloat(Ljava/lang/String;)F"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatParseFloat,
		}

	ghelpers.MethodSignatures["java/lang/Float.resolveConstantDesc(Ljava/lang/invoke/MethodHandles$Lookup;)Ljava/lang/Float;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Float.shortValue()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatShortValue,
		}

	ghelpers.MethodSignatures["java/lang/Float.sum(FF)F"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  floatSum,
		}

	ghelpers.MethodSignatures["java/lang/Float.toHexString(F)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatToHexString,
		}

	ghelpers.MethodSignatures["java/lang/Float.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  floatToString,
		}

	ghelpers.MethodSignatures["java/lang/Float.toString(F)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatToStringStatic,
		}

	ghelpers.MethodSignatures["java/lang/Float.valueOf(F)Ljava/lang/Float;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatValueOf,
		}

	ghelpers.MethodSignatures["java/lang/Float.valueOf(Ljava/lang/String;)Ljava/lang/Float;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  floatValueOfString,
		}
}

var classNameFloat = "java/lang/Float"

// getFloat64ValueFromObject - Extract a float64 from a Double object.
// See javaLangDouble.go.

// Method: byteValue
func floatByteValue(params []interface{}) interface{} {
	var ff float64
	self := params[0].(*object.Object)
	ff = self.FieldTable["value"].Fvalue.(float64)
	return int64(int8(ff))
}

// Method: compare (FF)I
func floatCompare(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatCompare: Incorrect number of arguments")
	}

	f1, ok1 := params[0].(float64)
	f2, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatCompare: Invalid argument types")
	}

	if f1 < f2 {
		return int64(-1) // Properly handles NaN, as NaN < f is false
	}
	if f1 > f2 {
		return int64(1)
	}

	// Handle NaN and zeros
	f1bits := math.Float32bits(float32(f1))
	f2bits := math.Float32bits(float32(f2))
	if f1bits == f2bits {
		return int64(0)
	}
	if f1bits < f2bits {
		return int64(1)
	}
	return int64(-1)
}

// Method: compareTo (Ljava/lang/Float;)I
func floatCompareTo(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatCompareTo: Incorrect number of arguments")
	}

	self, ok1 := params[0].(*object.Object)
	other, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatCompareTo: Invalid argument types")
	}

	f1 := self.FieldTable["value"].Fvalue.(float64)
	f2 := other.FieldTable["value"].Fvalue.(float64)

	return floatCompare([]interface{}{f1, f2})
}

// Method: floatToIntBits (F)I
func floatFloatToIntBits(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatFloatToIntBits: expected 1 float argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	ff, ok := args[0].(float64)
	if !ok {
		errMsg := "floatFloatToIntBits: argument is not a float64"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	bits := math.Float32bits(float32(ff))
	if (bits&0x7F800000) == 0x7F800000 && (bits&0x007FFFFF) != 0 {
		return int64(int32(0x7fc00000))
	}
	return int64(int32(bits))
}

// Method: floatToRawIntBits (F)I
func floatFloatToRawIntBits(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatFloatToRawIntBits: expected 1 float argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	ff, ok := args[0].(float64)
	if !ok {
		errMsg := "floatFloatToRawIntBits: argument is not a float64"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	bits := math.Float32bits(float32(ff))
	return int64(int32(bits))
}

// Method: equals (Ljava/lang/Object;)Z
func floatEquals(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatEquals: Incorrect number of arguments")
	}
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatEquals: Invalid self object")
	}

	other, ok := params[1].(*object.Object)
	if !ok || object.IsNull(other) {
		return types.JavaBoolFalse
	}

	if object.GoStringFromStringPoolIndex(other.KlassName) != classNameFloat {
		return types.JavaBoolFalse
	}

	selfValue := self.FieldTable["value"].Fvalue.(float64)
	otherValue := other.FieldTable["value"].Fvalue.(float64)

	if math.Float32bits(float32(selfValue)) == math.Float32bits(float32(otherValue)) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: floatValue ()F
func floatFloatValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["value"].Fvalue.(float64)
}

// Method: intValue ()I
func floatIntValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return int64(int32(selfValue))
}

// floatIsFiniteStatic (F)Z
func floatIsFiniteStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	if !math.IsNaN(f) && !math.IsInf(f, 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isInfinite ()Z
func floatIsInfinite(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	if math.IsInf(selfValue, 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// floatIsInfiniteStatic (F)Z
func floatIsInfiniteStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	if math.IsInf(f, 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isNaN ()Z
func floatIsNaN(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	if math.IsNaN(selfValue) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// floatIsNaNStatic (F)Z
func floatIsNaNStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	if math.IsNaN(f) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// floatHashCode ()I
func floatHashCode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	f := self.FieldTable["value"].Fvalue.(float64)
	return floatHashCodeStatic([]interface{}{f})
}

// floatHashCodeStatic (F)I
func floatHashCodeStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	bits := math.Float32bits(float32(f))
	return int64(int32(bits))
}

// Method: intBitsToFloat (I)F
func floatIntBitsToFloat(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatIntBitsToFloat: expected 1 int argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	bits, ok := args[0].(int64) // Java int maps to Go int64
	if !ok {
		errMsg := "floatIntBitsToFloat: argument is not an int64"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	ff := math.Float32frombits(uint32(bits))
	return float64(ff)
}

// Method: longValue ()J
func floatLongValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return int64(selfValue)
}

// Method: max (FF)F
func floatMax(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatMax: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatMax: Invalid argument types")
	}
	return math.Max(a, b)
}

// Method: min (FF)F
func floatMin(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatMin: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatMin: Invalid argument types")
	}
	return math.Min(a, b)
}

// Method: parseFloat (Ljava/lang/String;)F
func floatParseFloat(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatParseFloat: Incorrect number of arguments")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatParseFloat: Invalid argument type")
	}
	str := object.GoStringFromStringObject(obj)
	if len(str) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "floatParseFloat: Argument string is null")
	}
	ff, err := strconv.ParseFloat(str, 32)
	if err != nil {
		errMsg := fmt.Sprintf("floatParseFloat: Failed to parse %s to a float32 value", str)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return ff
}

// Method: shortValue ()S
func floatShortValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return int64(int16(selfValue))
}

// Method: sum (FF)F
func floatSum(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatSum: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatSum: Invalid argument types")
	}
	return a + b
}

// Method: toHexString (F)Ljava/lang/String;
func floatToHexString(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatToHexString: Incorrect number of arguments")
	}

	f, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatToHexString: Invalid argument type")
	}

	// Use strconv.FormatFloat with 'x' to get Java-compatible hex string.
	// Java Float.toHexString(float) returns a string like "0x1.0p0"
	str := strconv.FormatFloat(f, 'x', -1, 32)
	return object.StringObjectFromGoString(str)
}

// Method: toString (F)Ljava/lang/String;
func floatToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return floatToStringStatic([]interface{}{selfValue})
}

// Method: toString (F)Ljava/lang/String;
func floatToStringStatic(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatToStringStatic: Incorrect number of arguments")
	}
	ff, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatToStringStatic: Invalid argument type")
	}

	// Java's Float.toString(float) behavior:
	// Use 'g' for very large or very small, but usually decimal.
	// strconv.FormatFloat(f, 'g', -1, 32) is closer than fmt.Sprintf("%f", ff)
	str := strconv.FormatFloat(ff, 'g', -1, 32)
	return object.StringObjectFromGoString(str)
}

// Method: valueOf (F)Ljava/lang/Float;
func floatValueOf(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatValueOf: Incorrect number of arguments")
	}
	ff, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatValueOf: Invalid argument type")
	}

	// Create a new Double object with the given value and return it.
	return object.MakePrimitiveObject(classNameFloat, types.Float, ff)
}

// Method: valueOf (Ljava/lang/String;)F
func floatValueOfString(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatValueOfString: Incorrect number of arguments")
	}
	// The first parameter is a string (to convert to double)
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "floatValueOfString: Invalid argument, expected String object")
	}

	// Convert the string to a double.
	strValue := object.GoStringFromStringObject(obj)
	ff, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException,
			fmt.Sprintf("floatValueOfString: Invalid string format for double: %v", strValue))
	}

	// Create a new Double object with the given value and return it.
	return object.MakePrimitiveObject(classNameFloat, types.Float, ff)
}

// Method: doubleValue ()D
func floatDoubleValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["value"].Fvalue.(float64)
}

/*
floatBinary16:
Governed by IEEE 754 for conversion.
Returns the float value closest to the numerical value of the argument, a floating-point binary16 value encoded in a short. The conversion is exact; all binary16 values can be exactly represented in float. Special cases:
If the argument is zero, the result is a zero with the same sign as the argument.
If the argument is infinite, the result is an infinity with the same sign as the argument.
If the argument is a NaN, the result is a NaN.
*/
func floatFloat16ToFloat(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatFloat16ToFloat: expected 1 argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	raw, ok := args[0].(int64)
	if !ok {
		errMsg := "floatFloat16ToFloat: argument is not an int64"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	h := uint16(raw)

	sign := (h >> 15) & 0x0001
	exponent := (h >> 10) & 0x001F
	fraction := h & 0x03FF

	var result float32

	switch exponent {
	case 0:
		if fraction == 0 {
			// Zero
			result = float32(0)
		} else {
			// Subnormal number
			result = float32((1 << 23) * float64(fraction) / float64(1<<10) / float64(1<<14))
		}
	case 0x1F:
		if fraction == 0 {
			// Infinity
			result = float32(math.Inf(int(sign)))
		} else {
			// NaN
			result = float32(math.NaN())
		}
	default:
		// Normalized number
		exp := int(exponent) - 15 + 127 // adjust bias from float16 (15) to float32 (127)
		mantissa := int(fraction) << 13 // align to 23-bit mantissa
		bits := (int(sign) << 31) | (exp << 23) | mantissa
		result = math.Float32frombits(uint32(bits))
	}

	return float64(result) // return as float64 to match Java float mapping
}

/*
floatFloatToFloat16:
Governed by IEEE 754 for conversion.
Returns the floating-point binary16 value, encoded in a short, closest in value to the argument. The conversion is computed under the round to nearest even rounding mode. Special cases:
If the argument is zero, the result is a zero with the same sign as the argument.
If the argument is infinite, the result is an infinity with the same sign as the argument.
If the argument is a NaN, the result is a NaN.
*/
func floatFloatToFloat16(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatFloatToFloat16: expected 1 argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	val, ok := args[0].(float64)
	if !ok {
		errMsg := "floatFloatToFloat16: argument is not a float64"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	f := float32(val)
	bits := math.Float32bits(f)

	sign := (bits >> 31) & 0x1
	exp := (bits >> 23) & 0xFF
	mantissa := bits & 0x7FFFFF

	var h uint16

	switch exp {
	case 0:
		// Zero or subnormal
		h = uint16(sign << 15)
	case 0xFF:
		// Inf or NaN
		if mantissa == 0 {
			h = (uint16(sign) << 15) | 0x7C00 // Inf
		} else {
			h = (uint16(sign) << 15) | 0x7C00 | uint16(mantissa>>13) // NaN
		}
	default:
		newExp := int(exp) - 127 + 15
		if newExp >= 0x1F {
			// Overflow to infinity
			h = (uint16(sign) << 15) | 0x7C00
		} else if newExp <= 0 {
			// Underflow to subnormal or zero
			if newExp < -10 {
				h = uint16(sign << 15) // too small → zero
			} else {
				// Subnormal float16
				shift := uint(14 - newExp)
				subMantissa := int((mantissa | 0x800000) >> shift)
				h = (uint16(sign) << 15) | uint16(subMantissa&0x03FF)
			}
		} else {
			// Proper float16
			h = (uint16(sign) << 15) | (uint16(newExp) << 10) | uint16(mantissa>>13)
		}
	}

	return int64(h)
}
