/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func TestTemporalAdjusterGeneral(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_DayOfWeek()
	Load_Time_Month()
	Load_Time_Year()
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_LocalDateTime()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
	Load_Time_Temporal_TemporalAdjuster()

	// Null checks
	err := temporalAdjusterAdjustInto(nil)
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil params, got %v", err)
	}

	err = temporalAdjusterAdjustInto([]any{nil, nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self, got %v", err)
	}

	date := CreateLocalDate(2026, 5, 10)
	err = temporalAdjusterAdjustInto([]any{date, nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil target, got %v", err)
	}

	// Adjusting LocalDate into LocalDateTime
	dObj := CreateLocalDate(2020, 1, 1)
	tObj := CreateLocalTime(12, 30, 45, 0)
	dateTime := CreateLocalDateTime(dObj, tObj)
	newDate := CreateLocalDate(2026, 11, 8)
	res := temporalAdjusterAdjustInto([]any{newDate, dateTime})
	resObj, ok := res.(*object.Object)
	if !ok || resObj == nil {
		t.Fatalf("expected LocalDateTime object, got %v", res)
	}
	y, m, d, h, min, s, n := getLocalDateTimeYMDHMSN(resObj)
	if y != 2026 || m != 11 || d != 8 || h != 12 || min != 30 || s != 45 || n != 0 {
		t.Errorf("adjusting LocalDate into LocalDateTime: got %d-%d-%d %d:%d:%d.%d", y, m, d, h, min, s, n)
	}

	// Adjusting LocalTime into LocalDateTime
	newTime := CreateLocalTime(18, 45, 15, 500)
	res = temporalAdjusterAdjustInto([]any{newTime, dateTime})
	resObj = res.(*object.Object)
	y, m, d, h, min, s, n = getLocalDateTimeYMDHMSN(resObj)
	if y != 2020 || m != 1 || d != 1 || h != 18 || min != 45 || s != 15 || n != 500 {
		t.Errorf("adjusting LocalTime into LocalDateTime: got %d-%d-%d %d:%d:%d.%d", y, m, d, h, min, s, n)
	}

	// Adjusting DayOfWeek into LocalDate
	dowWed := dayOfWeekOf([]any{int64(3)}).(*object.Object) // Wednesday
	// 2026-11-08 is Sunday (7). Adjusting to Wednesday in the same week should change day
	res = temporalAdjusterAdjustInto([]any{dowWed, newDate})
	resObj = res.(*object.Object)
	y, m, d, okDate := getLocalDateYMD(resObj)
	if !okDate || y != 2026 || m != 11 || d != 4 {
		t.Errorf("adjusting DayOfWeek(Wed) into 2026-11-08: got %d-%d-%d, want 2026-11-04", y, m, d)
	}

	// Adjusting Month into LocalDate
	monthDec := monthOf([]any{int64(12)}).(*object.Object) // December
	res = temporalAdjusterAdjustInto([]any{monthDec, newDate})
	resObj = res.(*object.Object)
	y, m, d, okDate = getLocalDateYMD(resObj)
	if !okDate || y != 2026 || m != 12 || d != 8 {
		t.Errorf("adjusting Month(Dec) into 2026-11-08: got %d-%d-%d, want 2026-12-08", y, m, d)
	}

	// Adjusting Year into LocalDate
	year2030 := yearOf([]any{int64(2030)}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{year2030, newDate})
	resObj = res.(*object.Object)
	y, m, d, okDate = getLocalDateYMD(resObj)
	if !okDate || y != 2030 || m != 11 || d != 8 {
		t.Errorf("adjusting Year(2030) into 2026-11-08: got %d-%d-%d, want 2030-11-08", y, m, d)
	}
}
