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

func TestChronoLocalDate(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Chrono_Chronology()
	Load_Time_Chrono_ChronoLocalDate()
	Load_Time_Temporal_ChronoUnit()

	chronoUnitClinit(nil)

	isoName := object.StringObjectFromGoString("ISO")
	chronoObj := chronologyOf([]any{isoName}).(*object.Object)

	dateObj1 := chronologyDateYMD([]any{chronoObj, int64(2026), int64(11), int64(5)}).(*object.Object)
	dateObj2 := chronologyDateYMD([]any{chronoObj, int64(2026), int64(11), int64(10)}).(*object.Object)

	// isBefore, isAfter, isEqual
	if !isJavaBoolTrue(chronoLocalDateIsBefore([]any{dateObj1, dateObj2})) {
		t.Errorf("Expected date1 to be before date2")
	}
	if !isJavaBoolTrue(chronoLocalDateIsAfter([]any{dateObj2, dateObj1})) {
		t.Errorf("Expected date2 to be after date1")
	}
	if !isJavaBoolTrue(chronoLocalDateIsEqual([]any{dateObj1, dateObj1})) {
		t.Errorf("Expected date1 to equal itself")
	}

	// lengthOfMonth and lengthOfYear
	if chronoLocalDateLengthOfMonth([]any{dateObj1}).(int64) != 30 {
		t.Errorf("Expected lengthOfMonth for November to be 30, got %v", chronoLocalDateLengthOfMonth([]any{dateObj1}))
	}
	if chronoLocalDateLengthOfYear([]any{dateObj1}).(int64) != 365 {
		t.Errorf("Expected lengthOfYear for 2026 to be 365, got %v", chronoLocalDateLengthOfYear([]any{dateObj1}))
	}

	// isLeapYear
	if isJavaBoolTrue(chronoLocalDateIsLeapYear([]any{dateObj1})) {
		t.Errorf("Expected 2026 not to be leap year")
	}

	// toEpochDay
	epochDay := chronoLocalDateToEpochDay([]any{dateObj1})
	if _, okDay := epochDay.(int64); !okDay {
		t.Errorf("Expected int64 from toEpochDay, got %T", epochDay)
	}

	// plus and minus
	daysName := object.StringObjectFromGoString("DAYS")
	daysUnit := chronoUnitValueOf([]any{daysName}).(*object.Object)

	plusDate := chronoLocalDatePlusUnit([]any{dateObj1, int64(5), daysUnit}).(*object.Object)
	if !isJavaBoolTrue(chronoLocalDateIsEqual([]any{plusDate, dateObj2})) {
		t.Errorf("Expected date1 + 5 days to equal date2")
	}

	minusDate := chronoLocalDateMinusUnit([]any{dateObj2, int64(5), daysUnit}).(*object.Object)
	if !isJavaBoolTrue(chronoLocalDateIsEqual([]any{minusDate, dateObj1})) {
		t.Errorf("Expected date2 - 5 days to equal date1")
	}

	// until(ChronoLocalDate) -> ChronoPeriod
	untilPeriod := chronoLocalDateUntilPeriod([]any{dateObj1, dateObj2})
	if object.IsNull(untilPeriod.(*object.Object)) {
		t.Errorf("Expected non-null ChronoPeriod from until")
	}

	// until(Temporal, TemporalUnit) -> long
	untilDays := chronoLocalDateUntilUnit([]any{dateObj1, dateObj2, daysUnit})
	if untilDays.(int64) != 5 {
		t.Errorf("Expected 5 days until, got %v", untilDays)
	}

	// toString
	strResult := chronoLocalDateToString([]any{dateObj1})
	if object.GoStringFromStringObject(strResult.(*object.Object)) != "2026-11-05" {
		t.Errorf("Expected '2026-11-05', got '%s'", object.GoStringFromStringObject(strResult.(*object.Object)))
	}
}
