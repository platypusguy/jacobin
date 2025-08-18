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
	"math"
	"strconv"
)

func Load_Lang_Float() {

	MethodSignatures["java/lang/Float.<init>(F)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Float.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatByteValue,
		}

	MethodSignatures["java/lang/Float.compare(FF)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  floatCompare,
		}

	MethodSignatures["java/lang/Float.compareTo(Ljava/lang/Float;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatCompareTo,
		}

	MethodSignatures["java/lang/Float.describeConstable()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Float.floatToIntBits(F)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatFloatToIntBits,
		}

	MethodSignatures["java/lang/Float.floatToRawIntBits(F)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatFloatToIntBits,
		}

	MethodSignatures["java/lang/Float.float16ToFloat(S)F"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatFloat16ToFloat,
		}

	MethodSignatures["java/lang/Float.floatToFloat16(F)S"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatFloatToFloat16,
		}

	MethodSignatures["java/lang/Float.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatDoubleValue,
		}

	MethodSignatures["java/lang/Float.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatEquals,
		}

	MethodSignatures["java/lang/Float.floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatFloatValue,
		}

	MethodSignatures["java/lang/Float.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Float.hashCode(F)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Float.intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatIntValue,
		}

	MethodSignatures["java/lang/Float.isFinite(F)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatIsFinite,
		}

	MethodSignatures["java/lang/Float.isInfinite()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatIsInfinite,
		}

	MethodSignatures["java/lang/Float.isInfinite(F)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatIsInfinite,
		}

	MethodSignatures["java/lang/Float.isNaN()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatIsNaN,
		}

	MethodSignatures["java/lang/Float.isNaN(F)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatIsNaN,
		}

	MethodSignatures["java/lang/Float.intBitsToFloat(I)F"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatIntBitsToFloat,
		}

	MethodSignatures["java/lang/Float.longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatLongValue,
		}

	MethodSignatures["java/lang/Float.max(FF)F"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  floatMax,
		}

	MethodSignatures["java/lang/Float.min(FF)F"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  floatMin,
		}

	MethodSignatures["java/lang/Float.parseFloat(Ljava/lang/String;)F"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatParseFloat,
		}

	MethodSignatures["java/lang/Float.resolveConstantDesc(Ljava/lang/invoke/MethodHandles$Lookup;)Ljava/lang/Float;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Float.shortValue()S"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatShortValue,
		}

	MethodSignatures["java/lang/Float.sum(FF)F"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  floatSum,
		}

	MethodSignatures["java/lang/Float.toHexString(F)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatToHexString,
		}

	MethodSignatures["java/lang/Float.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  floatToString,
		}

	MethodSignatures["java/lang/Float.toString(F)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatToStringStatic,
		}

	MethodSignatures["java/lang/Float.valueOf(F)Ljava/lang/Float;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatValueOf,
		}

	MethodSignatures["java/lang/Float.valueOf(Ljava/lang/String;)Ljava/lang/Float;"] =
		GMeth{
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
	return int64(byte(ff))
}

// Method: compare (FF)I
func floatCompare(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatCompare: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatCompare: Invalid self object, expected Double object")
	}

	// The second parameter is the other object to compare
	other, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatCompare: Invalid other object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatCompare: Failed to retrieve value from self Double object")
	}

	// Retrieve the value of the other Double object
	otherValue, ok := getFloat64ValueFromObject(other)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatCompare: Failed to retrieve value from other Double object")
	}

	// Java's compareTo method for Double (return 0 if equal, 1 if greater, -1 if smaller)
	if selfValue < otherValue {
		return int64(-1)
	} else if selfValue > otherValue {
		return int64(1)
	}
	return int64(0)
}

// Method: compareTo (Ljava/lang/Float;)I
func floatCompareTo(params []interface{}) interface{} {
	var ff1, ff2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	ff1 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Get the actual Java Double parameter
	parmObj = params[1].(*object.Object)
	ff2 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Now, its just like doubleCompare.
	if ff1 == ff2 {
		return int64(0)
	}
	if ff1 < ff2 {
		return int64(-1)
	}
	return int64(1)
}

// Method: floatToIntBits (F)I
func floatFloatToIntBits(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatFloatToIntBits: expected 1 float argument"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	ff, ok := args[0].(float64) // Java float maps to Go float64 in your setup
	if !ok {
		errMsg := "floatFloatToIntBits: argument is not a float64"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Convert float64 to float32 before using Float32bits
	bits := math.Float32bits(float32(ff))
	return int64(bits)
}

// Method: equals (Ljava/lang/Object;)Z
func floatEquals(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatEquals: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatEquals: Invalid self object, expected Double object")
	}

	// The second parameter is the other object to compare
	other, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatEquals: Invalid other object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatEquals: Failed to retrieve value from self Double object")
	}

	// Retrieve the value of the other Double object
	otherValue, ok := getFloat64ValueFromObject(other)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatEquals: Failed to retrieve value from other Double object")
	}

	// Check if the values are equal (Java's == for primitive doubles)
	if selfValue == otherValue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: floatValue ()F
func floatFloatValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatFloatValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatFloatValue: Failed to retrieve value from self Double object")
	}
	return selfValue
}

// Method: intValue ()I
func floatIntValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIntValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIntValue: Failed to retrieve value from self Double object")
	}
	return int64(int32(selfValue))
}

// Method: isFinite (F)Z
func floatIsFinite(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsFinite: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsFinite: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsFinite: Failed to retrieve value from self Double object")
	}

	// Check if the value is finite (i.e., not NaN or Infinity)
	isFinite := !math.IsNaN(selfValue) && !math.IsInf(selfValue, 0)

	// Return the result as Java boolean (JavaBoolTrue or JavaBoolFalse)
	if isFinite {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isInfinite (F)Z
func floatIsInfinite(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsInfinite: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsInfinite: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsInfinite: Failed to retrieve value from self Double object")
	}

	// Check if the value is infinite (positive or negative infinity)
	isInfinite := math.IsInf(selfValue, 0)

	// Return the result as Java boolean (JavaBoolTrue or JavaBoolFalse)
	if isInfinite {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isNaN ()
func floatIsNaN(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsInfinite: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsInfinite: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatIsInfinite: Failed to retrieve value from self Double object")
	}

	// Check if the value is infinite (positive or negative infinity)
	isNaN := selfValue == math.NaN()

	// Return the result as Java boolean (JavaBoolTrue or JavaBoolFalse)
	if isNaN {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: intBitsToFloat (I)F
func floatIntBitsToFloat(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "floatIntBitsToFloat: expected 1 int argument"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	bits, ok := args[0].(int64) // Java int maps to Go int64
	if !ok {
		errMsg := "floatIntBitsToFloat: argument is not an int64"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	ff := math.Float32frombits(uint32(bits))
	return float64(ff)
}

// Method: longValue ()J
func floatLongValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatLongValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatLongValue: Failed to retrieve value from self Double object")
	}
	return int64(selfValue)
}

// Method: max (FF)F
func floatMax(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatMax: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatMax: Invalid argument types")
	}
	if a > b {
		return a
	}
	return b
}

// Method: min (FF)F
func floatMin(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatMin: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatMin: Invalid argument types")
	}
	if a < b {
		return a
	}
	return b
}

// Method: parseFloat (Ljava/lang/String;)F
func floatParseFloat(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatParseFloat: Incorrect number of arguments")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return getGErrBlk(excNames.IllegalArgumentException, "floatParseFloat: Invalid argument type")
	}
	str := object.GoStringFromStringObject(obj)
	if len(str) == 0 {
		return getGErrBlk(excNames.NullPointerException, "floatParseFloat: Argument string is null")
	}
	ff, err := strconv.ParseFloat(str, 32)
	if err != nil {
		errMsg := fmt.Sprintf("floatParseFloat: Failed to parse %s to a float32 value", str)
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return ff
}

// Method: shortValue ()S
func floatShortValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatShortValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatShortValue: Failed to retrieve value from self Double object")
	}
	return int16(selfValue)
}

// Method: sum (FF)F
func floatSum(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatSum: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatSum: Invalid argument types")
	}
	return a + b
}

// Method: toHexString (F)Ljava/lang/String;
func floatToHexString(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToHexString: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToHexString: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToHexString: Failed to retrieve value from self Double object")
	}

	// Get the raw bits of the float value.
	rawBits := math.Float64bits(selfValue)

	// Format the raw bits as a hexadecimal string
	hexString := fmt.Sprintf("0x%016X", rawBits)

	// Return the result as a Java String
	return object.StringObjectFromGoString(hexString)

}

// Method: toString (F)Ljava/lang/String;
func floatToString(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToString: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToString: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToString: Failed to retrieve value from self Double object")
	}

	// Convert the float value to string
	strValue := fmt.Sprintf("%g", selfValue) // %g is the format for general floating-point notation
	return object.StringObjectFromGoString(strValue)
}

// Method: toString (F)Ljava/lang/String;
func floatToStringStatic(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToStringStatic: Incorrect number of arguments")
	}
	ff, ok := params[0].(float64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatToStringStatic: Invalid argument type")
	}
	// Return string representation of the double.
	return object.StringObjectFromGoString(fmt.Sprintf("%f", ff))
}

// Method: valueOf (F)Ljava/lang/Float;
func floatValueOf(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatValueOf: Incorrect number of arguments")
	}
	ff, ok := params[0].(float64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatValueOf: Invalid argument type")
	}

	// Create a new Double object with the given value and return it.
	return object.MakePrimitiveObject(classNameFloat, types.Float, ff)
}

// Method: valueOf (Ljava/lang/String;)F
func floatValueOfString(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatValueOfString: Incorrect number of arguments")
	}
	// The first parameter is a string (to convert to double)
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return getGErrBlk(excNames.IllegalArgumentException, "floatValueOfString: Invalid argument, expected String object")
	}

	// Convert the string to a double.
	strValue := object.GoStringFromStringObject(obj)
	ff, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return getGErrBlk(excNames.NumberFormatException,
			fmt.Sprintf("floatValueOfString: Invalid string format for double: %v", strValue))
	}

	// Create a new Double object with the given value and return it.
	return object.MakePrimitiveObject(classNameFloat, types.Float, ff)
}

// Method: doubleValue ()D
func floatDoubleValue(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "floatDoubleValue: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatDoubleValue: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "floatDoubleValue: Failed to retrieve value from self Double object")
	}

	// Return the float value.
	return selfValue
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	raw, ok := args[0].(int64)
	if !ok {
		errMsg := "floatFloat16ToFloat: argument is not an int64"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	val, ok := args[0].(float64)
	if !ok {
		errMsg := "floatFloatToFloat16: argument is not a float64"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
