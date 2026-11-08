/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"strconv"
	"strings"
)

func Load_Time_Temporal_ValueRange() {
	ghelpers.MethodSignatures["java/time/temporal/ValueRange.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.of(JJ)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  valueRangeOf2,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.of(JJJ)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  valueRangeOf3,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.of(JJJJ)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  valueRangeOf4,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.checkValidIntValue(JLjava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  valueRangeCheckValidIntValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.checkValidValue(JLjava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  valueRangeCheckValidValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueRangeEquals,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.getLargestMinimum()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeGetLargestMinimum,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.getMaximum()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeGetMaximum,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.getMinimum()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeGetMinimum,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.getSmallestMaximum()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeGetSmallestMaximum,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeHashCode,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.isFixed()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeIsFixed,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.isIntValue()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeIsIntValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.isValidIntValue(J)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueRangeIsValidIntValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.isValidValue(J)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  valueRangeIsValidValue,
		}

	ghelpers.MethodSignatures["java/time/temporal/ValueRange.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  valueRangeToStringMethod,
		}
}

func extractValueRangeBounds(rangeObject *object.Object) (int64, int64, int64, int64) {
	if rangeObject == nil || rangeObject.FieldTable == nil {
		return 0, 0, 0, 0
	}
	var minSmallestVal int64
	var minLargestVal int64
	var maxSmallestVal int64
	var maxLargestVal int64

	if fieldVal, exists := rangeObject.FieldTable["minSmallest"]; exists {
		if parsedVal, ok := getInt64Val(fieldVal.Fvalue); ok {
			minSmallestVal = parsedVal
		}
	}
	if fieldVal, exists := rangeObject.FieldTable["minLargest"]; exists {
		if parsedVal, ok := getInt64Val(fieldVal.Fvalue); ok {
			minLargestVal = parsedVal
		}
	}
	if fieldVal, exists := rangeObject.FieldTable["maxSmallest"]; exists {
		if parsedVal, ok := getInt64Val(fieldVal.Fvalue); ok {
			maxSmallestVal = parsedVal
		}
	}
	if fieldVal, exists := rangeObject.FieldTable["maxLargest"]; exists {
		if parsedVal, ok := getInt64Val(fieldVal.Fvalue); ok {
			maxLargestVal = parsedVal
		}
	}
	return minSmallestVal, minLargestVal, maxSmallestVal, maxLargestVal
}

func valueRangeToStringHelper(rangeObject *object.Object) string {
	minSmallestVal, minLargestVal, maxSmallestVal, maxLargestVal := extractValueRangeBounds(rangeObject)
	var builder strings.Builder
	builder.WriteString(strconv.FormatInt(minSmallestVal, 10))
	if minSmallestVal != minLargestVal {
		builder.WriteByte('/')
		builder.WriteString(strconv.FormatInt(minLargestVal, 10))
	}
	builder.WriteString(" - ")
	builder.WriteString(strconv.FormatInt(maxSmallestVal, 10))
	if maxSmallestVal != maxLargestVal {
		builder.WriteByte('/')
		builder.WriteString(strconv.FormatInt(maxLargestVal, 10))
	}
	return builder.String()
}

func valueRangeOf2(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ValueRange.of: insufficient parameters")
	}
	minimumValue, okMin := getInt64Val(params[0])
	maximumValue, okMax := getInt64Val(params[1])
	if !okMin || !okMax {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ValueRange.of: invalid parameter type")
	}
	if minimumValue > maximumValue {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Minimum value must be less than maximum value")
	}
	return createValueRange(minimumValue, minimumValue, maximumValue, maximumValue)
}

func valueRangeOf3(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ValueRange.of: insufficient parameters")
	}
	minimumValue, okMin := getInt64Val(params[0])
	maxSmallestValue, okMaxSmallest := getInt64Val(params[1])
	maxLargestValue, okMaxLargest := getInt64Val(params[2])
	if !okMin || !okMaxSmallest || !okMaxLargest {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ValueRange.of: invalid parameter type")
	}
	return valueRangeOf4([]any{minimumValue, minimumValue, maxSmallestValue, maxLargestValue})
}

func valueRangeOf4(params []any) any {
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ValueRange.of: insufficient parameters")
	}
	minSmallestVal, okMinSmallest := getInt64Val(params[0])
	minLargestVal, okMinLargest := getInt64Val(params[1])
	maxSmallestVal, okMaxSmallest := getInt64Val(params[2])
	maxLargestVal, okMaxLargest := getInt64Val(params[3])
	if !okMinSmallest || !okMinLargest || !okMaxSmallest || !okMaxLargest {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ValueRange.of: invalid parameter type")
	}
	if minSmallestVal > minLargestVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Smallest minimum value must be less than largest minimum value")
	}
	if maxSmallestVal > maxLargestVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Smallest maximum value must be less than largest maximum value")
	}
	if minLargestVal > maxLargestVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Minimum value must be less than maximum value")
	}
	if minSmallestVal > maxSmallestVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Minimum value must be less than maximum value")
	}
	return createValueRange(minSmallestVal, minLargestVal, maxSmallestVal, maxLargestVal)
}

func valueRangeGetMinimum(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.getMinimum: self is null")
	}
	selfObject := params[0].(*object.Object)
	minSmallestVal, _, _, _ := extractValueRangeBounds(selfObject)
	return minSmallestVal
}

func valueRangeGetLargestMinimum(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.getLargestMinimum: self is null")
	}
	selfObject := params[0].(*object.Object)
	_, minLargestVal, _, _ := extractValueRangeBounds(selfObject)
	return minLargestVal
}

func valueRangeGetSmallestMaximum(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.getSmallestMaximum: self is null")
	}
	selfObject := params[0].(*object.Object)
	_, _, maxSmallestVal, _ := extractValueRangeBounds(selfObject)
	return maxSmallestVal
}

func valueRangeGetMaximum(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.getMaximum: self is null")
	}
	selfObject := params[0].(*object.Object)
	_, _, _, maxLargestVal := extractValueRangeBounds(selfObject)
	return maxLargestVal
}

func valueRangeIsFixed(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.isFixed: self is null")
	}
	selfObject := params[0].(*object.Object)
	minSmallestVal, minLargestVal, maxSmallestVal, maxLargestVal := extractValueRangeBounds(selfObject)
	if minSmallestVal == minLargestVal && maxSmallestVal == maxLargestVal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func valueRangeIsIntValue(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.isIntValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	minSmallestVal, _, _, maxLargestVal := extractValueRangeBounds(selfObject)
	if minSmallestVal >= math.MinInt32 && maxLargestVal <= math.MaxInt32 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func valueRangeIsValidValue(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.isValidValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	testedValue, okVal := getInt64Val(params[1])
	if !okVal {
		return types.JavaBoolFalse
	}
	minSmallestVal, _, _, maxLargestVal := extractValueRangeBounds(selfObject)
	if testedValue >= minSmallestVal && testedValue <= maxLargestVal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func valueRangeIsValidIntValue(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.isValidIntValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	testedValue, okVal := getInt64Val(params[1])
	if !okVal {
		return types.JavaBoolFalse
	}
	if valueRangeIsIntValue([]any{selfObject}) == types.JavaBoolTrue &&
		valueRangeIsValidValue([]any{selfObject, testedValue}) == types.JavaBoolTrue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func valueRangeCheckValidValue(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.checkValidValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	testedValue, okVal := getInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "ValueRange.checkValidValue: invalid value type")
	}
	var fieldObject *object.Object
	if len(params) >= 3 && params[2] != nil {
		if candidateObject, okObj := params[2].(*object.Object); okObj && !object.IsNull(candidateObject) {
			fieldObject = candidateObject
		}
	}

	if valueRangeIsValidValue([]any{selfObject, testedValue}) != types.JavaBoolTrue {
		fieldName := getTemporalFieldName(fieldObject)
		var errorMessage string
		if fieldName != "" {
			errorMessage = fmt.Sprintf("Invalid value for %s (valid values %s): %d", fieldName, valueRangeToStringHelper(selfObject), testedValue)
		} else {
			errorMessage = fmt.Sprintf("Invalid value (valid values %s): %d", valueRangeToStringHelper(selfObject), testedValue)
		}
		return ghelpers.GetGErrBlk(excNames.DateTimeException, errorMessage)
	}
	return testedValue
}

func valueRangeCheckValidIntValue(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.checkValidIntValue: self is null")
	}
	selfObject := params[0].(*object.Object)
	testedValue, okVal := getInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "ValueRange.checkValidIntValue: invalid value type")
	}
	var fieldObject *object.Object
	if len(params) >= 3 && params[2] != nil {
		if candidateObject, okObj := params[2].(*object.Object); okObj && !object.IsNull(candidateObject) {
			fieldObject = candidateObject
		}
	}

	if valueRangeIsValidIntValue([]any{selfObject, testedValue}) != types.JavaBoolTrue {
		fieldName := getTemporalFieldName(fieldObject)
		var errorMessage string
		if valueRangeIsValidValue([]any{selfObject, testedValue}) == types.JavaBoolTrue {
			if fieldName != "" {
				errorMessage = fmt.Sprintf("Invalid int value for %s: %d", fieldName, testedValue)
			} else {
				errorMessage = fmt.Sprintf("Invalid int value: %d", testedValue)
			}
		} else {
			if fieldName != "" {
				errorMessage = fmt.Sprintf("Invalid value for %s (valid values %s): %d", fieldName, valueRangeToStringHelper(selfObject), testedValue)
			} else {
				errorMessage = fmt.Sprintf("Invalid value (valid values %s): %d", valueRangeToStringHelper(selfObject), testedValue)
			}
		}
		return ghelpers.GetGErrBlk(excNames.DateTimeException, errorMessage)
	}
	return testedValue
}

func valueRangeEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	selfObject := params[0].(*object.Object)
	if params[1] == nil {
		return types.JavaBoolFalse
	}
	otherObject, okObj := params[1].(*object.Object)
	if !okObj || object.IsNull(otherObject) {
		return types.JavaBoolFalse
	}
	if selfObject == otherObject {
		return types.JavaBoolTrue
	}
	if object.GoStringFromStringPoolIndex(otherObject.KlassName) != classNameValueRange && otherObject.FieldTable == nil {
		return types.JavaBoolFalse
	}
	minSmallestA, minLargestA, maxSmallestA, maxLargestA := extractValueRangeBounds(selfObject)
	minSmallestB, minLargestB, maxSmallestB, maxLargestB := extractValueRangeBounds(otherObject)
	if minSmallestA == minSmallestB && minLargestA == minLargestB &&
		maxSmallestA == maxSmallestB && maxLargestA == maxLargestB {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func valueRangeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return int64(0)
	}
	selfObject := params[0].(*object.Object)
	minSmallestVal, minLargestVal, maxSmallestVal, maxLargestVal := extractValueRangeBounds(selfObject)
	totalValue := minSmallestVal +
		(minLargestVal << 16) + (minLargestVal >> 48) +
		(maxSmallestVal << 32) + (maxSmallestVal >> 32) +
		(maxLargestVal << 48) + (maxLargestVal >> 16)
	hashComputed := int64(int32(totalValue ^ int64(uint64(totalValue)>>32)))
	return hashComputed
}

func valueRangeToStringMethod(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ValueRange.toString: self is null")
	}
	selfObject := params[0].(*object.Object)
	return object.StringObjectFromGoString(valueRangeToStringHelper(selfObject))
}
