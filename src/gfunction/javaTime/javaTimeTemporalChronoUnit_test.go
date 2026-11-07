/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func TestChronoUnit(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Temporal_ChronoUnit()

	// Run clinit
	chronoUnitClinit(nil)

	// Test values()
	valuesResult := chronoUnitValues(nil)
	valuesArray, okArray := valuesResult.(*object.Object)
	if !okArray || object.IsNull(valuesArray) {
		t.Fatalf("Expected non-null array from values(), got %v", valuesResult)
	}

	// Test valueOf("DAYS")
	daysName := object.StringObjectFromGoString("DAYS")
	daysResult := chronoUnitValueOf([]any{daysName})
	daysObject, okDays := daysResult.(*object.Object)
	if !okDays || object.IsNull(daysObject) {
		t.Fatalf("Expected DAYS object, got %v", daysResult)
	}

	daysUnitName := chronoUnitName([]any{daysObject})
	if object.GoStringFromStringObject(daysUnitName.(*object.Object)) != "DAYS" {
		t.Errorf("Expected 'DAYS', got '%s'", object.GoStringFromStringObject(daysUnitName.(*object.Object)))
	}

	if !isJavaBoolTrue(chronoUnitIsDateBased([]any{daysObject})) {
		t.Errorf("Expected DAYS isDateBased to be true")
	}

	if isJavaBoolTrue(chronoUnitIsTimeBased([]any{daysObject})) {
		t.Errorf("Expected DAYS isTimeBased to be false")
	}

	if !isJavaBoolTrue(chronoUnitIsDurationEstimated([]any{daysObject})) {
		t.Errorf("Expected DAYS isDurationEstimated to be true")
	}

	// Test duration
	durationResult := chronoUnitGetDuration([]any{daysObject})
	durationObject, okDuration := durationResult.(*object.Object)
	if !okDuration || object.IsNull(durationObject) {
		t.Fatalf("Expected non-null Duration object for DAYS")
	}

	// Test valueOf("SECONDS")
	secondsName := object.StringObjectFromGoString("SECONDS")
	secondsResult := chronoUnitValueOf([]any{secondsName})
	secondsObject := secondsResult.(*object.Object)

	if !isJavaBoolTrue(chronoUnitIsTimeBased([]any{secondsObject})) {
		t.Errorf("Expected SECONDS isTimeBased to be true")
	}

	if isJavaBoolTrue(chronoUnitIsDurationEstimated([]any{secondsObject})) {
		t.Errorf("Expected SECONDS isDurationEstimated to be false")
	}

	// Test toString
	stringResult := chronoUnitToString([]any{daysObject})
	stringRepresentation := object.GoStringFromStringObject(stringResult.(*object.Object))
	if stringRepresentation != "Days" {
		t.Errorf("Expected 'Days', got '%s'", stringRepresentation)
	}

	// Test compareTo and equals
	compareResult := chronoUnitCompareTo([]any{daysObject, daysObject})
	if compareResult.(int64) != 0 {
		t.Errorf("Expected compareTo with self to be 0, got %v", compareResult)
	}

	equalResult := chronoUnitEquals([]any{daysObject, daysObject})
	if !isJavaBoolTrue(equalResult) {
		t.Errorf("Expected DAYS to equal itself")
	}
}
