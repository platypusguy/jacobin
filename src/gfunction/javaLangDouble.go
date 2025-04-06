/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"math"
	"strconv"
	"unsafe"
)

func Load_Lang_Double() {

	MethodSignatures["java/lang/Double.<init>(D)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Double.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleByteValue,
		}

	MethodSignatures["java/lang/Double.compare(DD)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  doubleCompare,
		}

	MethodSignatures["java/lang/Double.compareTo(Ljava/lang/Double;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleCompareTo,
		}

	MethodSignatures["java/lang/Double.describeConstable()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Double.doubleToLongBits(D)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleToLongBits,
		}

	MethodSignatures["java/lang/Double.doubleToRawLongBits(D)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleToLongBits,
		}

	MethodSignatures["java/lang/Double.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleDoubleValue,
		}

	MethodSignatures["java/lang/Double.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleEquals,
		}

	MethodSignatures["java/lang/Double.floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleFloatValue,
		}

	MethodSignatures["java/lang/Double.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Double.hashCode(D)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Double.intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleIntValue,
		}

	MethodSignatures["java/lang/Double.isFinite(D)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleIsFinite,
		}

	MethodSignatures["java/lang/Double.isInfinite()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleIsInfinite,
		}

	MethodSignatures["java/lang/Double.isInfinite(D)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleIsInfinite,
		}

	MethodSignatures["java/lang/Double.isNaN()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleIsNaN,
		}

	MethodSignatures["java/lang/Double.isNaN(D)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleIsNaN,
		}

	MethodSignatures["java/lang/Double.longBitsToDouble(J)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleLongBitsToDouble,
		}

	MethodSignatures["java/lang/Double.longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleLongValue,
		}

	MethodSignatures["java/lang/Double.max(DD)D"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  doubleMax,
		}

	MethodSignatures["java/lang/Double.min(DD)D"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  doubleMin,
		}

	MethodSignatures["java/lang/Double.parseDouble(Ljava/lang/String;)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleParseDouble,
		}

	MethodSignatures["java/lang/Double.resolveConstantDesc(Ljava/lang/invoke/MethodHandles$Lookup;)Ljava/lang/Double;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Double.shortValue()S"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleShortValue,
		}

	MethodSignatures["java/lang/Double.sum(DD)D"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  doubleSum,
		}

	MethodSignatures["java/lang/Double.toHexString(D)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleToHexString,
		}

	MethodSignatures["java/lang/Double.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleToString,
		}

	MethodSignatures["java/lang/Double.toString(D)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleToStringStatic,
		}

	MethodSignatures["java/lang/Double.valueOf(D)Ljava/lang/Double;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleValueOf,
		}

	MethodSignatures["java/lang/Double.valueOf(Ljava/lang/String;)Ljava/lang/Double;"] =
		GMeth{
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
	return int64(byte(dd))
}

// Method: compare (DD)I
func doubleCompare(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Invalid self object, expected Double object")
	}

	// The second parameter is the other object to compare
	other, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Invalid other object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Failed to retrieve value from self Double object")
	}

	// Retrieve the value of the other Double object
	otherValue, ok := getFloat64ValueFromObject(other)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleCompare: Failed to retrieve value from other Double object")
	}

	// Java's compareTo method for Double (return 0 if equal, 1 if greater, -1 if smaller)
	if selfValue < otherValue {
		return int64(-1)
	} else if selfValue > otherValue {
		return int64(1)
	}
	return int64(0)
}

// Method: compareTo (Ljava/lang/Double;)I
func doubleCompareTo(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	dd1 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Get the actual Java Double parameter
	parmObj = params[1].(*object.Object)
	dd2 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Now, its just like doubleCompare.
	if dd1 == dd2 {
		return int64(0)
	}
	if dd1 < dd2 {
		return int64(-1)
	}
	return int64(1)
}

// Method: doubleToLongBits (D)J
func doubleToLongBits(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToLongBits: Incorrect number of arguments")
	}
	// The parameter is the float64 argument.
	arg, ok := params[0].(float64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToLongBits: Invalid float64 argument")
	}

	// If the value is NaN, return the special bit pattern for NaN
	if math.IsNaN(arg) {
		return int64(0x7FF8000000000000) // NaN bit pattern for double in Java
	}

	// Otherwise, convert double (float64) to raw long bits (uint64)
	rawBits := math.Float64bits(arg)

	// Return the raw bits as a Java long (represented by int64 in Go)
	return int64(rawBits)
}

// Method: equals (Ljava/lang/Object;)Z
func doubleEquals(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Invalid self object, expected Double object")
	}

	// The second parameter is the other object to compare
	other, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Invalid other object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Failed to retrieve value from self Double object")
	}

	// Retrieve the value of the other Double object
	otherValue, ok := getFloat64ValueFromObject(other)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleEquals: Failed to retrieve value from other Double object")
	}

	// Check if the values are equal (Java's == for primitive doubles)
	if selfValue == otherValue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: floatValue ()F
func doubleFloatValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleFloatValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleFloatValue: Failed to retrieve value from self Double object")
	}
	return selfValue
}

// Method: intValue ()I
func doubleIntValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIntValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIntValue: Failed to retrieve value from self Double object")
	}
	return int32(selfValue)
}

// Method: isFinite (D)Z
func doubleIsFinite(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsFinite: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsFinite: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsFinite: Failed to retrieve value from self Double object")
	}

	// Check if the value is finite (i.e., not NaN or Infinity)
	isFinite := !math.IsNaN(selfValue) && !math.IsInf(selfValue, 0)

	// Return the result as Java boolean (JavaBoolTrue or JavaBoolFalse)
	if isFinite {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: isInfinite (D)Z
func doubleIsInfinite(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsInfinite: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsInfinite: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsInfinite: Failed to retrieve value from self Double object")
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
func doubleIsNaN(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsInfinite: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsInfinite: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleIsInfinite: Failed to retrieve value from self Double object")
	}

	// Check if the value is infinite (positive or negative infinity)
	isNaN := selfValue == math.NaN()

	// Return the result as Java boolean (JavaBoolTrue or JavaBoolFalse)
	if isNaN {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Method: longBitsToDouble (J)D
func doubleLongBitsToDouble(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleLongBitsToDouble: Incorrect number of arguments")
	}
	lb, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleLongBitsToDouble: Invalid argument type")
	}
	// Convert long bits to double
	return *(*float64)(unsafe.Pointer(&lb))
}

// Method: longValue ()J
func doubleLongValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleLongValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleLongValue: Failed to retrieve value from self Double object")
	}
	return int64(selfValue)
}

// Method: max (DD)D
func doubleMax(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleMax: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleMax: Invalid argument types")
	}
	if a > b {
		return a
	}
	return b
}

// Method: min (DD)D
func doubleMin(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleMin: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleMin: Invalid argument types")
	}
	if a < b {
		return a
	}
	return b
}

// Method: parseDouble (Ljava/lang/String;)D
func doubleParseDouble(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleParseDouble: Incorrect number of arguments")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleParseDouble: Invalid argument type")
	}
	str := object.GoStringFromStringObject(obj)
	if len(str) == 0 {
		return getGErrBlk(excNames.NullPointerException, "doubleParseDouble: Argument string is null")
	}
	dd, err := strconv.ParseFloat(str, 64)
	if err != nil {
		errMsg := fmt.Sprintf("doubleParseDouble: Failed to parse %s to a float64 value", str)
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}
	return dd
}

// Method: shortValue ()S
func doubleShortValue(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleShortValue: Invalid self object, expected Double object")
	}
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleShortValue: Failed to retrieve value from self Double object")
	}
	return int16(selfValue)
}

// Method: sum (DD)D
func doubleSum(params []interface{}) interface{} {
	if len(params) != 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleSum: Incorrect number of arguments")
	}
	a, ok1 := params[0].(float64)
	b, ok2 := params[1].(float64)
	if !ok1 || !ok2 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleSum: Invalid argument types")
	}
	return a + b
}

// Method: toHexString (D)Ljava/lang/String;
func doubleToHexString(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToHexString: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToHexString: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToHexString: Failed to retrieve value from self Double object")
	}

	// Get the raw bits of the double value
	rawBits := math.Float64bits(selfValue)

	// Format the raw bits as a hexadecimal string
	hexString := fmt.Sprintf("0x%016X", rawBits)

	// Return the result as a Java String
	return object.StringObjectFromGoString(hexString)

}

// Method: toString (D)Ljava/lang/String;
func doubleToString(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToString: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToString: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToString: Failed to retrieve value from self Double object")
	}

	// Convert the double value to string
	strValue := fmt.Sprintf("%g", selfValue) // %g is the format for general floating-point notation
	return object.StringObjectFromGoString(strValue)
}

// Method: toString (D)Ljava/lang/String;
func doubleToStringStatic(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToStringStatic: Incorrect number of arguments")
	}
	dd, ok := params[0].(float64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleToStringStatic: Invalid argument type")
	}
	// Return string representation of the double (mocked for illustration)
	return object.StringObjectFromGoString(fmt.Sprintf("%f", dd))
}

// Method: valueOf (D)Ljava/lang/Double;
func doubleValueOf(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleValueOf: Incorrect number of arguments")
	}
	dd, ok := params[0].(float64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleValueOf: Invalid argument type")
	}

	// Create a new Double object with the given value and return it
	return object.MakePrimitiveObject(classNameDouble, types.Double, dd)
}

// Method: valueOf (Ljava/lang/String;)D
func doubleValueOfString(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleValueOfString: Incorrect number of arguments")
	}
	// The first parameter is a string (to convert to double)
	obj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(obj) {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleValueOfString: Invalid argument, expected String object")
	}

	// Convert the string to a double.
	strValue := object.GoStringFromStringObject(obj)
	dd, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return getGErrBlk(excNames.NumberFormatException,
			fmt.Sprintf("doubleValueOfString: Invalid string format for double: %v", strValue))
	}

	// Create a new Double object with the given value and return it
	return object.MakePrimitiveObject(classNameDouble, types.Double, dd)
}

// Method: doubleValue ()D
func doubleDoubleValue(params []interface{}) interface{} {
	if len(params) != 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleDoubleValue: Incorrect number of arguments")
	}
	// The first parameter is the self object (this)
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleDoubleValue: Invalid self object, expected Double object")
	}

	// Retrieve the value of the current Double (this object)
	selfValue, ok := getFloat64ValueFromObject(self)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "doubleDoubleValue: Failed to retrieve value from self Double object")
	}

	// Return the double value
	return selfValue
}
