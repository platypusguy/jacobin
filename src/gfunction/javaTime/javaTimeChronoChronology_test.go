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

func TestChronology(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Chrono_Chronology()

	// of("ISO")
	isoName := object.StringObjectFromGoString("ISO")
	chronoResult := chronologyOf([]any{isoName})
	chronoObj, okChrono := chronoResult.(*object.Object)
	if !okChrono || object.IsNull(chronoObj) {
		t.Fatalf("Expected non-null ISO Chronology object, got %v", chronoResult)
	}

	idObj := chronologyGetId([]any{chronoObj})
	if object.GoStringFromStringObject(idObj.(*object.Object)) != "ISO" {
		t.Errorf("Expected id 'ISO', got '%s'", object.GoStringFromStringObject(idObj.(*object.Object)))
	}

	calTypeObj := chronologyGetCalendarType([]any{chronoObj})
	if object.GoStringFromStringObject(calTypeObj.(*object.Object)) != "iso8601" {
		t.Errorf("Expected calendarType 'iso8601', got '%s'", object.GoStringFromStringObject(calTypeObj.(*object.Object)))
	}

	// isLeapYear
	if !isJavaBoolTrue(chronologyIsLeapYear([]any{chronoObj, int64(2024)})) {
		t.Errorf("Expected 2024 to be leap year")
	}
	if isJavaBoolTrue(chronologyIsLeapYear([]any{chronoObj, int64(2023)})) {
		t.Errorf("Expected 2023 not to be leap year")
	}

	// date(year, month, day)
	dateResult := chronologyDateYMD([]any{chronoObj, int64(2026), int64(11), int64(5)})
	dateObj, okDate := dateResult.(*object.Object)
	if !okDate || object.IsNull(dateObj) {
		t.Fatalf("Expected non-null ChronoLocalDate from dateYMD")
	}

	// dateEpochDay
	epochDayDate := chronologyDateEpochDay([]any{chronoObj, int64(0)})
	if object.IsNull(epochDayDate.(*object.Object)) {
		t.Errorf("Expected non-null date from dateEpochDay(0)")
	}

	// dateYearDay
	yearDayDate := chronologyDateYearDay([]any{chronoObj, int64(2026), int64(309)})
	if object.IsNull(yearDayDate.(*object.Object)) {
		t.Errorf("Expected non-null date from dateYearDay")
	}

	// epochSecond
	epochSec := chronologyEpochSecond([]any{chronoObj, int64(2026), int64(11), int64(5), int64(12), int64(0), int64(0), object.Null})
	if _, okSec := epochSec.(int64); !okSec {
		t.Errorf("Expected int64 from epochSecond, got %T", epochSec)
	}

	// eras
	erasList := chronologyEras([]any{chronoObj})
	if object.IsNull(erasList.(*object.Object)) {
		t.Errorf("Expected non-null eras list")
	}

	// eraOf(1) -> CE
	ceEra := chronologyEraOf([]any{chronoObj, int64(1)})
	if object.IsNull(ceEra.(*object.Object)) {
		t.Errorf("Expected non-null Era for eraOf(1)")
	}

	// prolepticYear
	prolepticYear := chronologyProlepticYear([]any{chronoObj, ceEra, int64(2026)})
	if prolepticYear.(int64) != 2026 {
		t.Errorf("Expected prolepticYear 2026, got %v", prolepticYear)
	}

	// period(years, months, days)
	periodResult := chronologyPeriod([]any{chronoObj, int64(1), int64(2), int64(3)})
	if object.IsNull(periodResult.(*object.Object)) {
		t.Errorf("Expected non-null ChronoPeriod")
	}

	// toString, equals, hashCode
	strResult := chronologyToString([]any{chronoObj})
	if object.GoStringFromStringObject(strResult.(*object.Object)) != "ISO" {
		t.Errorf("Expected toString 'ISO', got '%s'", object.GoStringFromStringObject(strResult.(*object.Object)))
	}

	if !isJavaBoolTrue(chronologyEquals([]any{chronoObj, chronoObj})) {
		t.Errorf("Expected chronology to equal itself")
	}

	if chronologyCompareTo([]any{chronoObj, chronoObj}).(int64) != 0 {
		t.Errorf("Expected compareTo with self to be 0")
	}
}
