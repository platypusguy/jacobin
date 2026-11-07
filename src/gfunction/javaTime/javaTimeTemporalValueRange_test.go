/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func isGErrBlk(value any) bool {
	_, ok := value.(*ghelpers.GErrBlk)
	return ok
}

func isJavaBoolTrue(value any) bool {
	if boolVal, ok := value.(int64); ok {
		return boolVal == types.JavaBoolTrue
	}
	return false
}

func TestValueRange(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_ValueRange()

	// Test of(min, max)
	rangeResult := valueRangeOf2([]any{int64(1), int64(12)})
	rangeObject, okRange := rangeResult.(*object.Object)
	if !okRange || object.IsNull(rangeObject) {
		t.Fatalf("Expected non-null ValueRange object, got %v", rangeResult)
	}

	minimumValue := valueRangeGetMinimum([]any{rangeObject})
	if minimumValue.(int64) != 1 {
		t.Errorf("Expected minimum 1, got %v", minimumValue)
	}

	maximumValue := valueRangeGetMaximum([]any{rangeObject})
	if maximumValue.(int64) != 12 {
		t.Errorf("Expected maximum 12, got %v", maximumValue)
	}

	largestMinimumValue := valueRangeGetLargestMinimum([]any{rangeObject})
	if largestMinimumValue.(int64) != 1 {
		t.Errorf("Expected largest minimum 1, got %v", largestMinimumValue)
	}

	smallestMaximumValue := valueRangeGetSmallestMaximum([]any{rangeObject})
	if smallestMaximumValue.(int64) != 12 {
		t.Errorf("Expected smallest maximum 12, got %v", smallestMaximumValue)
	}

	isFixedValue := valueRangeIsFixed([]any{rangeObject})
	if !isJavaBoolTrue(isFixedValue) {
		t.Errorf("Expected isFixed to be true")
	}

	isIntValueResult := valueRangeIsIntValue([]any{rangeObject})
	if !isJavaBoolTrue(isIntValueResult) {
		t.Errorf("Expected isIntValue to be true")
	}

	// Test of(min, minSmallestMax, max)
	range3Result := valueRangeOf3([]any{int64(1), int64(28), int64(31)})
	range3Object := range3Result.(*object.Object)
	if isJavaBoolTrue(valueRangeIsFixed([]any{range3Object})) {
		t.Errorf("Expected isFixed to be false for range (1, 28, 31)")
	}

	// Test isValidValue and isValidIntValue
	if !isJavaBoolTrue(valueRangeIsValidValue([]any{rangeObject, int64(6)})) {
		t.Errorf("Expected 6 to be valid value in 1..12")
	}
	if isJavaBoolTrue(valueRangeIsValidValue([]any{rangeObject, int64(13)})) {
		t.Errorf("Expected 13 to be invalid value in 1..12")
	}

	// Test checkValidValue
	checkedValue := valueRangeCheckValidValue([]any{rangeObject, int64(5), object.Null})
	if checkedValue.(int64) != 5 {
		t.Errorf("Expected checkValidValue to return 5, got %v", checkedValue)
	}

	invalidCheck := valueRangeCheckValidValue([]any{rangeObject, int64(20), object.Null})
	if !isGErrBlk(invalidCheck) {
		t.Errorf("Expected error for out of range checkValidValue, got %v", invalidCheck)
	}

	// Test checkValidIntValue
	checkedIntValue := valueRangeCheckValidIntValue([]any{rangeObject, int64(5), object.Null})
	if checkedIntValue.(int64) != 5 {
		t.Errorf("Expected checkValidIntValue to return 5, got %v", checkedIntValue)
	}

	// Test equals and hashCode
	equalResult := valueRangeEquals([]any{rangeObject, rangeObject})
	if !isJavaBoolTrue(equalResult) {
		t.Errorf("Expected range to equal itself")
	}

	hashResult := valueRangeHashCode([]any{rangeObject})
	if _, okHash := hashResult.(int64); !okHash {
		t.Errorf("Expected int64 hashCode, got %T", hashResult)
	}

	// Test toString
	stringResult := valueRangeToStringMethod([]any{rangeObject})
	stringObject, okString := stringResult.(*object.Object)
	if !okString || !object.IsStringObject(stringObject) {
		t.Fatalf("Expected String object from toString, got %v", stringResult)
	}
	stringRepresentation := object.GoStringFromStringObject(stringObject)
	if stringRepresentation != "1 - 12" {
		t.Errorf("Expected '1 - 12', got '%s'", stringRepresentation)
	}

	// Test invalid factory: min > max
	invalidRangeResult := valueRangeOf2([]any{int64(10), int64(5)})
	if !isGErrBlk(invalidRangeResult) {
		t.Errorf("Expected error when min > max, got %v", invalidRangeResult)
	}
}

func TestValueRangeOf4(t *testing.T) {
	globals.InitGlobals("test")

	rangeResult := valueRangeOf4([]any{int64(1), int64(2), int64(28), int64(31)})
	rangeObject, okRange := rangeResult.(*object.Object)
	if !okRange || object.IsNull(rangeObject) {
		t.Fatalf("Expected non-null ValueRange object, got %v", rangeResult)
	}

	stringResult := valueRangeToStringMethod([]any{rangeObject})
	stringRepresentation := object.GoStringFromStringObject(stringResult.(*object.Object))
	if stringRepresentation != "1/2 - 28/31" {
		t.Errorf("Expected '1/2 - 28/31', got '%s'", stringRepresentation)
	}

	// Invalid of4 cases
	if !isGErrBlk(valueRangeOf4([]any{int64(2), int64(1), int64(28), int64(31)})) {
		t.Errorf("Expected error when minSmallest > minLargest")
	}
	if !isGErrBlk(valueRangeOf4([]any{int64(1), int64(2), int64(32), int64(31)})) {
		t.Errorf("Expected error when maxSmallest > maxLargest")
	}
	if !isGErrBlk(valueRangeOf4([]any{int64(1), int64(50), int64(28), int64(31)})) {
		t.Errorf("Expected error when minLargest > maxSmallest")
	}
}
