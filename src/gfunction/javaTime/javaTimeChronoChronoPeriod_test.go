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

func TestChronoPeriod(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Chrono_Chronology()
	Load_Time_Chrono_ChronoPeriod()
	Load_Time_Temporal_ChronoUnit()

	chronoUnitClinit(nil)

	isoName := object.StringObjectFromGoString("ISO")
	chronoObj := chronologyOf([]any{isoName}).(*object.Object)

	periodObj := chronologyPeriod([]any{chronoObj, int64(1), int64(2), int64(3)}).(*object.Object)

	// isZero and isNegative
	if isJavaBoolTrue(chronoPeriodIsZero([]any{periodObj})) {
		t.Errorf("Expected P1Y2M3D isZero to be false")
	}
	if isJavaBoolTrue(chronoPeriodIsNegative([]any{periodObj})) {
		t.Errorf("Expected P1Y2M3D isNegative to be false")
	}

	// get(TemporalUnit)
	yearsName := object.StringObjectFromGoString("YEARS")
	yearsUnit := chronoUnitValueOf([]any{yearsName}).(*object.Object)
	yearsVal := chronoPeriodGet([]any{periodObj, yearsUnit})
	if yearsVal.(int64) != 1 {
		t.Errorf("Expected 1 year, got %v", yearsVal)
	}

	monthsName := object.StringObjectFromGoString("MONTHS")
	monthsUnit := chronoUnitValueOf([]any{monthsName}).(*object.Object)
	monthsVal := chronoPeriodGet([]any{periodObj, monthsUnit})
	if monthsVal.(int64) != 2 {
		t.Errorf("Expected 2 months, got %v", monthsVal)
	}

	daysName := object.StringObjectFromGoString("DAYS")
	daysUnit := chronoUnitValueOf([]any{daysName}).(*object.Object)
	daysVal := chronoPeriodGet([]any{periodObj, daysUnit})
	if daysVal.(int64) != 3 {
		t.Errorf("Expected 3 days, got %v", daysVal)
	}

	// multipliedBy, negated, normalized
	multPeriod := chronoPeriodMultipliedBy([]any{periodObj, int64(2)}).(*object.Object)
	if chronoPeriodGet([]any{multPeriod, yearsUnit}).(int64) != 2 {
		t.Errorf("Expected 2 years after multipliedBy 2, got %v", chronoPeriodGet([]any{multPeriod, yearsUnit}))
	}

	negPeriod := chronoPeriodNegated([]any{periodObj}).(*object.Object)
	if !isJavaBoolTrue(chronoPeriodIsNegative([]any{negPeriod})) {
		t.Errorf("Expected negated period to be negative")
	}

	// toString
	strResult := chronoPeriodToString([]any{periodObj})
	if object.GoStringFromStringObject(strResult.(*object.Object)) != "ISO P1Y2M3D" {
		t.Errorf("Expected 'ISO P1Y2M3D', got '%s'", object.GoStringFromStringObject(strResult.(*object.Object)))
	}
}
