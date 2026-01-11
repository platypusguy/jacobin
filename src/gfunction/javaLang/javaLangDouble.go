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

func Load_Lang_Double() {

	ghelpers.MethodSignatures["java/lang/Double.<init>(D)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/lang/Double.byteValue()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleByteValue,
		}

	ghelpers.MethodSignatures["java/lang/Double.compare(DD)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  doubleCompare,
		}

	ghelpers.MethodSignatures["java/lang/Double.compareTo(Ljava/lang/Double;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleCompareTo,
		}

	ghelpers.MethodSignatures["java/lang/Double.describeConstable()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Double.doubleToLongBits(D)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleToLongBits,
		}

	ghelpers.MethodSignatures["java/lang/Double.doubleToRawLongBits(D)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleToRawLongBits,
		}

	ghelpers.MethodSignatures["java/lang/Double.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleDoubleValue,
		}

	ghelpers.MethodSignatures["java/lang/Double.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleEquals,
		}

	ghelpers.MethodSignatures["java/lang/Double.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleFloatValue,
		}

	ghelpers.MethodSignatures["java/lang/Double.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleHashCode,
		}

	ghelpers.MethodSignatures["java/lang/Double.hashCode(D)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleHashCodeStatic,
		}

	ghelpers.MethodSignatures["java/lang/Double.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleIntValue,
		}

	ghelpers.MethodSignatures["java/lang/Double.isFinite(D)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleIsFiniteStatic,
		}

	ghelpers.MethodSignatures["java/lang/Double.isInfinite()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleIsInfinite,
		}

	ghelpers.MethodSignatures["java/lang/Double.isInfinite(D)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleIsInfiniteStatic,
		}

	ghelpers.MethodSignatures["java/lang/Double.isNaN()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleIsNaN,
		}

	ghelpers.MethodSignatures["java/lang/Double.isNaN(D)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleIsNaNStatic,
		}

	ghelpers.MethodSignatures["java/lang/Double.longBitsToDouble(J)D"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleLongBitsToDouble,
		}

	ghelpers.MethodSignatures["java/lang/Double.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleLongValue,
		}

	ghelpers.MethodSignatures["java/lang/Double.max(DD)D"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  doubleMax,
		}

	ghelpers.MethodSignatures["java/lang/Double.min(DD)D"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  doubleMin,
		}

	ghelpers.MethodSignatures["java/lang/Double.parseDouble(Ljava/lang/String;)D"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleParseDouble,
		}

	ghelpers.MethodSignatures["java/lang/Double.resolveConstantDesc(Ljava/lang/invoke/MethodHandles$Lookup;)Ljava/lang/Double;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Double.shortValue()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleShortValue,
		}

	ghelpers.MethodSignatures["java/lang/Double.sum(DD)D"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  doubleSum,
		}

	ghelpers.MethodSignatures["java/lang/Double.toHexString(D)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleToHexString,
		}

	ghelpers.MethodSignatures["java/lang/Double.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  doubleToString,
		}

	ghelpers.MethodSignatures["java/lang/Double.toString(D)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleToStringStatic,
		}

	ghelpers.MethodSignatures["java/lang/Double.valueOf(D)Ljava/lang/Double;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleValueOf,
		}

	ghelpers.MethodSignatures["java/lang/Double.valueOf(Ljava/lang/String;)Ljava/lang/Double;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  doubleValueOfString,
		}
}

var classNameDouble = "java/lang/Double"

// getFloat64ValueFromObject - Extract a float64 from a Double object.
func getFloat64ValueFromObject(obj *object.Object) (float64, bool) {
	field := obj.FieldTable["value"]
	if field.Ftype != types.Double {
		return math.NaN(), false
	}
	fvalue, ok := field.Fvalue.(float64)
	if ok {
		return fvalue, true
	}

	return math.NaN(), false
}

// Method: byteValue
func doubleByteValue(params []interface{}) interface{} {
	var dd float64
	self := params[0].(*object.Object)
	dd = self.FieldTable["value"].Fvalue.(float64)
	return int64(int8(dd))
}

// Method: compare (DD)I
func doubleCompare(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Incorrect number of arguments")
	}

	d1, ok1 := params[0].(float64)
	d2, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Invalid argument types")
	}

	if d1 < d2 {
		return int64(-1)
	}
	if d1 > d2 {
		return int64(1)
	}

	// Handle special cases: NaN and +/-0.0
	// Java's Double.compare(d1, d2) behavior:
	// - NaN is equal to itself and greater than any other value (including +Inf)
	// - +0.0 is greater than -0.0
	d1bits := math.Float64bits(d1)
	d2bits := math.Float64bits(d2)
	if d1bits == d2bits {
		return int64(0)
	}
	// For all other cases, we can compare the bit patterns as SIGNED integers to get Java behavior
	// except that we need to handle NaN specially if we want it to be "greater" than Inf.
	// Actually, if we use the bit patterns and treat them as signed, it works for everything EXCEPT NaN.
	// In IEEE 754, NaNs have bit patterns larger than Inf.
	// If we compare as SIGNED 64-bit:
	// - Positive numbers: larger bits mean larger values.
	// - Negative numbers: larger bits (more negative in two's complement) mean smaller values.
	// This matches Java's requirement for 0.0 > -0.0.
	// However, Java wants NaN > +Inf.
	// In bit patterns: +Inf is 0x7FF0000000000000, NaNs are 0x7FF0...1 to 0x7FF...F.
	// So NaN bits are numerically greater than +Inf bits.
	// For negative NaN: bits are 0xFFF... and they should be > everything else too.

	// Let's use a simpler logic that is explicitly correct:
	if math.IsNaN(d1) {
		if math.IsNaN(d2) {
			return int64(0)
		}
		return int64(1)
	} else if math.IsNaN(d2) {
		return int64(-1)
	}

	// Not NaN, handle zeros
	if d1 == 0 && d2 == 0 {
		if d1bits == d2bits {
			return int64(0)
		}
		if d1bits < d2bits { // -0.0 has sign bit set, so it's larger as uint64
			return int64(1) // +0.0 > -0.0
		}
		return int64(-1)
	}

	// Normal comparison
	if d1 < d2 {
		return int64(-1)
	}
	if d1 > d2 {
		return int64(1)
	}
	return int64(0)
}

// Method: compareTo (Ljava/lang/Double;)I
func doubleCompareTo(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleCompareTo: Incorrect number of arguments")
	}

	self, ok1 := params[0].(*object.Object)
	other, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleCompareTo: Invalid argument types")
	}

	d1 := self.FieldTable["value"].Fvalue.(float64)
	d2 := other.FieldTable["value"].Fvalue.(float64)

	return doubleCompare([]interface{}{d1, d2})
}

// Method: doubleToLongBits (D)J
func doubleToLongBits(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToLongBits: Incorrect number of arguments")
	}
	arg, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToLongBits: Invalid float64 argument")
	}

	if math.IsNaN(arg) {
		return int64(0x7ff8000000000000)
	}

	return int64(math.Float64bits(arg))
}

// Method: doubleToRawLongBits (D)J
func doubleToRawLongBits(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToRawLongBits: Incorrect number of arguments")
	}
	arg, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToRawLongBits: Invalid float64 argument")
	}

	return int64(math.Float64bits(arg))
}

// Method: equals (Ljava/lang/Object;)Z
func doubleEquals(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Incorrect number of arguments")
	}
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Invalid self object")
	}

	other, ok := params[1].(*object.Object)
	if !ok || object.IsNull(other) {
		return types.JavaBoolFalse
	}

	if object.GoStringFromStringPoolIndex(other.KlassName) != classNameDouble {
		return types.JavaBoolFalse
	}

	selfValue := self.FieldTable["value"].Fvalue.(float64)
	otherValue := other.FieldTable["value"].Fvalue.(float64)

	if math.Float64bits(selfValue) == math.Float64bits(otherValue) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: floatValue ()F
func doubleFloatValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["value"].Fvalue.(float64)
}

// doubleHashCode ()I
func doubleHashCode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	f := self.FieldTable["value"].Fvalue.(float64)
	return doubleHashCodeStatic([]interface{}{f})
}

// doubleHashCodeStatic (D)I
func doubleHashCodeStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	bits := math.Float64bits(f)
	return int64(int32(bits ^ (bits >> 32)))
}

// Method: intValue ()I
func doubleIntValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return int64(int32(selfValue))
}

// doubleIsFiniteStatic (D)Z
func doubleIsFiniteStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	if !math.IsNaN(f) && !math.IsInf(f, 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isInfinite ()Z
func doubleIsInfinite(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	if math.IsInf(selfValue, 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// doubleIsInfiniteStatic (D)Z
func doubleIsInfiniteStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	if math.IsInf(f, 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isNaN ()Z
func doubleIsNaN(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	if math.IsNaN(selfValue) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// doubleIsNaNStatic (D)Z
func doubleIsNaNStatic(params []interface{}) interface{} {
	f := params[0].(float64)
	if math.IsNaN(f) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: longBitsToDouble (J)D
func doubleLongBitsToDouble(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleLongBitsToDouble: Incorrect number of arguments")
	}
	lb, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleLongBitsToDouble: Invalid argument type")
	}
	return math.Float64frombits(uint64(lb))
}

// Method: longValue ()J
func doubleLongValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return int64(selfValue)
}

// Method: max (DD)D
func doubleMax(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleMax: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleMax: Invalid argument types")
	}
	return math.Max(a, b)
}

// Method: min (DD)D
func doubleMin(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleMin: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleMin: Invalid argument types")
	}
	return math.Min(a, b)
}

// Method: parseDouble (Ljava/lang/String;)D
func doubleParseDouble(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleParseDouble: Incorrect number of arguments")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleParseDouble: Invalid argument type")
	}
	str := object.GoStringFromStringObject(obj)
	if len(str) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "doubleParseDouble: Argument string is null")
	}
	dd, err := strconv.ParseFloat(str, 64)
	if err != nil {
		errMsg := fmt.Sprintf("doubleParseDouble: Failed to parse %s to a float64 value", str)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return dd
}

// Method: shortValue ()S
func doubleShortValue(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return int64(int16(selfValue))
}

// Method: sum (DD)D
func doubleSum(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleSum: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleSum: Invalid argument types")
	}
	return a + b
}

// Method: toHexString (D)Ljava/lang/String;
func doubleToHexString(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToHexString: Incorrect number of arguments")
	}
	f, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToHexString: Invalid argument type")
	}

	str := strconv.FormatFloat(f, 'x', -1, 64)
	return object.StringObjectFromGoString(str)
}

// Method: toString ()Ljava/lang/String;
func doubleToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	selfValue := self.FieldTable["value"].Fvalue.(float64)
	return doubleToStringStatic([]interface{}{selfValue})
}

// Method: toString (D)Ljava/lang/String;
func doubleToStringStatic(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToStringStatic: Incorrect number of arguments")
	}
	dd, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleToStringStatic: Invalid argument type")
	}

	str := strconv.FormatFloat(dd, 'g', -1, 64)
	return object.StringObjectFromGoString(str)
}

// Method: valueOf (D)Ljava/lang/Double;
func doubleValueOf(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleValueOf: Incorrect number of arguments")
	}
	dd, ok := params[0].(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleValueOf: Invalid argument type")
	}

	// Create a new Double object with the given value and return it
	return object.MakePrimitiveObject(classNameDouble, types.Double, dd)
}

// Method: valueOf (Ljava/lang/String;)D
func doubleValueOfString(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleValueOfString: Incorrect number of arguments")
	}
	// The first parameter is a string (to convert to double)
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleValueOfString: Invalid argument, expected String object")
	}

	// Convert the string to a double.
	strValue := object.GoStringFromStringObject(obj)
	dd, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.NumberFormatException,
			fmt.Sprintf("doubleValueOfString: Invalid string format for double: %v", strValue))
	}

	// Create a new Double object with the given value and return it
	return object.MakePrimitiveObject(classNameDouble, types.Double, dd)
}

// Method: doubleValue ()D
func doubleDoubleValue(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleDoubleValue: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleDoubleValue: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "doubleDoubleValue: Failed to retrieve value from self Double object")
	}

	// Return the double value
	return selfValue
}
