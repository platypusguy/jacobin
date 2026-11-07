/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestChronoField(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_ChronoField()

	// Run clinit
	chronoFieldClinit(nil)

	// Test values()
	valuesResult := chronoFieldValues(nil)
	valuesArray, okArray := valuesResult.(*object.Object)
	if !okArray || object.IsNull(valuesArray) {
		t.Fatalf("Expected non-null array from values(), got %v", valuesResult)
	}

	// Test valueOf
	dayOfMonthName := object.StringObjectFromGoString("DAY_OF_MONTH")
	fieldResult := chronoFieldValueOf([]any{dayOfMonthName})
	fieldObject, okField := fieldResult.(*object.Object)
	if !okField || object.IsNull(fieldObject) {
		t.Fatalf("Expected DAY_OF_MONTH object, got %v", fieldResult)
	}

	nameResult := chronoFieldName([]any{fieldObject})
	nameObject := nameResult.(*object.Object)
	if object.GoStringFromStringObject(nameObject) != "DAY_OF_MONTH" {
		t.Errorf("Expected name 'DAY_OF_MONTH', got '%s'", object.GoStringFromStringObject(nameObject))
	}

	isDateBased := chronoFieldIsDateBased([]any{fieldObject})
	if !isJavaBoolTrue(isDateBased) {
		t.Errorf("Expected DAY_OF_MONTH isDateBased to be true")
	}

	isTimeBased := chronoFieldIsTimeBased([]any{fieldObject})
	if isJavaBoolTrue(isTimeBased) {
		t.Errorf("Expected DAY_OF_MONTH isTimeBased to be false")
	}

	rangeResult := chronoFieldRange([]any{fieldObject})
	rangeObject, okRange := rangeResult.(*object.Object)
	if !okRange || object.IsNull(rangeObject) {
		t.Fatalf("Expected non-null ValueRange for DAY_OF_MONTH")
	}

	// Test checkValidValue and checkValidIntValue
	checkedValue := chronoFieldCheckValidValue([]any{fieldObject, int64(15)})
	if checkedValue.(int64) != 15 {
		t.Errorf("Expected checkValidValue to return 15, got %v", checkedValue)
	}

	invalidCheck := chronoFieldCheckValidValue([]any{fieldObject, int64(32)})
	if !isGErrBlk(invalidCheck) {
		t.Errorf("Expected error for day 32, got %v", invalidCheck)
	}

	checkedIntValue := chronoFieldCheckValidIntValue([]any{fieldObject, int64(15)})
	if checkedIntValue.(int64) != 15 {
		t.Errorf("Expected checkValidIntValue to return 15, got %v", checkedIntValue)
	}

	// Test time-based field: NANO_OF_SECOND
	nanoName := object.StringObjectFromGoString("NANO_OF_SECOND")
	nanoResult := chronoFieldValueOf([]any{nanoName})
	nanoObject := nanoResult.(*object.Object)
	if !isJavaBoolTrue(chronoFieldIsTimeBased([]any{nanoObject})) {
		t.Errorf("Expected NANO_OF_SECOND isTimeBased to be true")
	}

	// Test toString
	stringResult := chronoFieldToString([]any{fieldObject})
	stringRepresentation := object.GoStringFromStringObject(stringResult.(*object.Object))
	if stringRepresentation != "DayOfMonth" {
		t.Errorf("Expected 'DayOfMonth', got '%s'", stringRepresentation)
	}

	// Test compareTo and equals
	compareResult := chronoFieldCompareTo([]any{fieldObject, fieldObject})
	if compareResult.(int64) != 0 {
		t.Errorf("Expected compareTo with self to be 0, got %v", compareResult)
	}

	equalResult := chronoFieldEquals([]any{fieldObject, fieldObject})
	if !isJavaBoolTrue(equalResult) {
		t.Errorf("Expected field to equal itself")
	}

	// Test getBaseUnit and getRangeUnit
	baseUnit := chronoFieldGetBaseUnit([]any{fieldObject})
	if object.IsNull(baseUnit.(*object.Object)) {
		t.Errorf("Expected non-null base unit")
	}

	rangeUnit := chronoFieldGetRangeUnit([]any{fieldObject})
	if object.IsNull(rangeUnit.(*object.Object)) {
		t.Errorf("Expected non-null range unit")
	}

	// Test getFrom with mock TemporalAccessor
	classNameTemporalAccessor := "java/time/temporal/TemporalAccessor"
	mockAccessor := object.MakeEmptyObjectWithClassName(&classNameTemporalAccessor)
	mockAccessor.FieldTable["DAY_OF_MONTH"] = object.Field{
		Ftype:  types.Long,
		Fvalue: int64(15),
	}

	fromValue := chronoFieldGetFrom([]any{fieldObject, mockAccessor})
	if fromValue.(int64) != 15 {
		t.Errorf("Expected getFrom to return 15, got %v", fromValue)
	}
}
