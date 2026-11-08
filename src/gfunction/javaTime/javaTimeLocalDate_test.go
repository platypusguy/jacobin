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
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func localDateTestInit() {
	globals.InitGlobals("test")
	localDateMutex.Lock()
	localDateInitialized = false
	localDateMutex.Unlock()
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_LocalDateTime()
	Load_Time_Period()
	Load_Time_DayOfWeek()
	Load_Time_Month()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
}

func TestLocalDateClinitAndConstants(t *testing.T) {
	localDateTestInit()

	if res := localDateClinit(nil); res != nil {
		t.Errorf("localDateClinit returned %v", res)
	}

	for _, name := range []string{"MIN", "MAX", "EPOCH"} {
		st, ok := statics.QueryStatic(classNameLocalDateVal, name)
		if !ok || st.Value == nil {
			t.Fatalf("LocalDate.%s static not found", name)
		}
	}

	epochObj := statics.GetStaticValue(classNameLocalDateVal, "EPOCH").(*object.Object)
	y, m, d, ok := getLocalDateYMD(epochObj)
	if !ok || y != 1970 || m != 1 || d != 1 {
		t.Errorf("LocalDate.EPOCH = %d-%d-%d, want 1970-1-1", y, m, d)
	}
}

func TestLocalDateFactories(t *testing.T) {
	localDateTestInit()

	// of(Y, M, D)
	d1 := localDateOfYMD([]any{int64(2026), int64(11), int64(8)}).(*object.Object)
	y, m, d, ok := getLocalDateYMD(d1)
	if !ok || y != 2026 || m != 11 || d != 8 {
		t.Errorf("LocalDate.of(2026, 11, 8) = %d-%d-%d", y, m, d)
	}

	// of invalid month/day
	err := localDateOfYMD([]any{int64(2026), int64(13), int64(1)})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeException {
		t.Errorf("of(2026, 13, 1) expected DateTimeException, got %v", err)
	}
	err = localDateOfYMD([]any{int64(2026), int64(2), int64(29)}) // 2026 not leap
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeException {
		t.Errorf("of(2026, 2, 29) expected DateTimeException, got %v", err)
	}

	// of with Month object
	novObj := monthOf([]any{int64(11)}).(*object.Object)
	d2 := localDateOfYMonthD([]any{int64(2026), novObj, int64(8)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(d2); y != 2026 || m != 11 || d != 8 {
		t.Errorf("of(2026, NOVEMBER, 8) = %d-%d-%d", y, m, d)
	}

	// ofYearDay
	dLeap := localDateOfYearDay([]any{int64(2024), int64(60)}).(*object.Object) // 2024 leap -> Feb 29
	if y, m, d, _ := getLocalDateYMD(dLeap); y != 2024 || m != 2 || d != 29 {
		t.Errorf("ofYearDay(2024, 60) = %d-%d-%d, want 2024-2-29", y, m, d)
	}

	// ofEpochDay
	dEpoch := localDateOfEpochDay([]any{int64(0)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(dEpoch); y != 1970 || m != 1 || d != 1 {
		t.Errorf("ofEpochDay(0) = %d-%d-%d, want 1970-1-1", y, m, d)
	}

	// parse
	parseD := localDateParse([]any{object.StringObjectFromGoString("2026-11-08")}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(parseD); y != 2026 || m != 11 || d != 8 {
		t.Errorf("parse('2026-11-08') = %d-%d-%d", y, m, d)
	}
	err = localDateParse([]any{object.StringObjectFromGoString("INVALID_DATE")})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeParseException {
		t.Errorf("parse('INVALID') expected DateTimeParseException, got %v", err)
	}

	// now
	nowRes := localDateNow(nil)
	if nowObj, ok := nowRes.(*object.Object); !ok || nowObj == nil {
		t.Errorf("now() returned %v", nowRes)
	}
}

func TestLocalDateGettersAndQueries(t *testing.T) {
	localDateTestInit()

	date := CreateLocalDate(2026, 11, 8) // Sunday

	if gy := localDateGetYear([]any{date}).(int64); gy != 2026 {
		t.Errorf("getYear() = %d, want 2026", gy)
	}
	if gm := localDateGetMonthValue([]any{date}).(int64); gm != 11 {
		t.Errorf("getMonthValue() = %d, want 11", gm)
	}
	if gd := localDateGetDayOfMonth([]any{date}).(int64); gd != 8 {
		t.Errorf("getDayOfMonth() = %d, want 8", gd)
	}
	if gdoy := localDateGetDayOfYear([]any{date}).(int64); gdoy != 312 {
		t.Errorf("getDayOfYear() = %d, want 312", gdoy)
	}

	dowObj := localDateGetDayOfWeek([]any{date}).(*object.Object)
	if dowVal := dayOfWeekGetValue([]any{dowObj}).(int64); dowVal != 7 {
		t.Errorf("getDayOfWeek() = %d, want 7 (SUNDAY)", dowVal)
	}

	mObj := localDateGetMonth([]any{date}).(*object.Object)
	if mVal := monthGetValue([]any{mObj}).(int64); mVal != 11 {
		t.Errorf("getMonth() = %d, want 11 (NOVEMBER)", mVal)
	}

	if isLeap := localDateIsLeapYear([]any{date}); isLeap != types.JavaBoolFalse && isLeap != false {
		t.Errorf("2026.isLeapYear() = %v, want false", isLeap)
	}
	date2024 := CreateLocalDate(2024, 1, 1)
	if isLeap := localDateIsLeapYear([]any{date2024}); isLeap != types.JavaBoolTrue && isLeap != true {
		t.Errorf("2024.isLeapYear() = %v, want true", isLeap)
	}

	if lm := localDateLengthOfMonth([]any{date}).(int64); lm != 30 {
		t.Errorf("lengthOfMonth(Nov) = %d, want 30", lm)
	}
	if ly := localDateLengthOfYear([]any{date}).(int64); ly != 365 {
		t.Errorf("lengthOfYear(2026) = %d, want 365", ly)
	}
	if ly := localDateLengthOfYear([]any{date2024}).(int64); ly != 366 {
		t.Errorf("lengthOfYear(2024) = %d, want 366", ly)
	}

	// ChronoField queries
	cfDOM := getChronoField("DAY_OF_MONTH")
	if sup := localDateIsSupportedField([]any{date, cfDOM}); sup != types.JavaBoolTrue && sup != true {
		t.Errorf("isSupported(DAY_OF_MONTH) = %v, want true", sup)
	}
	cfHour := getChronoField("HOUR_OF_DAY")
	if sup := localDateIsSupportedField([]any{date, cfHour}); sup != types.JavaBoolFalse && sup != false {
		t.Errorf("isSupported(HOUR_OF_DAY) = %v, want false", sup)
	}

	if val := localDateGet([]any{date, cfDOM}).(int64); val != 8 {
		t.Errorf("get(DAY_OF_MONTH) = %d, want 8", val)
	}
	if val := localDateGetLong([]any{date, cfDOM}).(int64); val != 8 {
		t.Errorf("getLong(DAY_OF_MONTH) = %d, want 8", val)
	}
}

func TestLocalDateAdjustmentsAndArithmetic(t *testing.T) {
	localDateTestInit()

	date := CreateLocalDate(2026, 11, 8)

	// withYear, withMonth, withDayOfMonth, withDayOfYear
	wy := localDateWithYear([]any{date, int64(2030)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(wy); y != 2030 || m != 11 || d != 8 {
		t.Errorf("withYear(2030) = %d-%d-%d, want 2030-11-8", y, m, d)
	}
	wm := localDateWithMonth([]any{date, int64(12)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(wm); y != 2026 || m != 12 || d != 8 {
		t.Errorf("withMonth(12) = %d-%d-%d, want 2026-12-8", y, m, d)
	}
	wd := localDateWithDayOfMonth([]any{date, int64(25)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(wd); y != 2026 || m != 11 || d != 25 {
		t.Errorf("withDayOfMonth(25) = %d-%d-%d, want 2026-11-25", y, m, d)
	}
	wdoy := localDateWithDayOfYear([]any{date, int64(1)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(wdoy); y != 2026 || m != 1 || d != 1 {
		t.Errorf("withDayOfYear(1) = %d-%d-%d, want 2026-1-1", y, m, d)
	}

	// plusYears, plusMonths, plusWeeks, plusDays
	py := localDatePlusYears([]any{date, int64(2)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(py); y != 2028 || m != 11 || d != 8 {
		t.Errorf("plusYears(2) = %d-%d-%d, want 2028-11-8", y, m, d)
	}
	pm := localDatePlusMonths([]any{date, int64(3)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(pm); y != 2027 || m != 2 || d != 8 {
		t.Errorf("plusMonths(3) = %d-%d-%d, want 2027-2-8", y, m, d)
	}
	pw := localDatePlusWeeks([]any{date, int64(2)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(pw); y != 2026 || m != 11 || d != 22 {
		t.Errorf("plusWeeks(2) = %d-%d-%d, want 2026-11-22", y, m, d)
	}
	pd := localDatePlusDays([]any{date, int64(10)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(pd); y != 2026 || m != 11 || d != 18 {
		t.Errorf("plusDays(10) = %d-%d-%d, want 2026-11-18", y, m, d)
	}

	// minusYears, minusMonths, minusWeeks, minusDays
	my := localDateMinusYears([]any{date, int64(1)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(my); y != 2025 || m != 11 || d != 8 {
		t.Errorf("minusYears(1) = %d-%d-%d, want 2025-11-8", y, m, d)
	}
	mm := localDateMinusMonths([]any{date, int64(2)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(mm); y != 2026 || m != 9 || d != 8 {
		t.Errorf("minusMonths(2) = %d-%d-%d, want 2026-9-8", y, m, d)
	}
	mw := localDateMinusWeeks([]any{date, int64(1)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(mw); y != 2026 || m != 11 || d != 1 {
		t.Errorf("minusWeeks(1) = %d-%d-%d, want 2026-11-1", y, m, d)
	}
	md := localDateMinusDays([]any{date, int64(7)}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(md); y != 2026 || m != 11 || d != 1 {
		t.Errorf("minusDays(7) = %d-%d-%d, want 2026-11-1", y, m, d)
	}
}

func TestLocalDateConversionsAndComparisons(t *testing.T) {
	localDateTestInit()

	d1 := CreateLocalDate(2026, 11, 8)
	d2 := CreateLocalDate(2026, 11, 9)

	// compareTo, isAfter, isBefore, isEqual, equals, hashCode, toString
	if cmp := localDateCompareTo([]any{d1, d2}).(int64); cmp >= 0 {
		t.Errorf("d1.compareTo(d2) = %d, want < 0", cmp)
	}
	if cmp := localDateCompareTo([]any{d1, d1}).(int64); cmp != 0 {
		t.Errorf("d1.compareTo(d1) = %d, want 0", cmp)
	}
	if isA := localDateIsAfter([]any{d2, d1}); isA != types.JavaBoolTrue && isA != true {
		t.Errorf("d2.isAfter(d1) = %v, want true", isA)
	}
	if isB := localDateIsBefore([]any{d1, d2}); isB != types.JavaBoolTrue && isB != true {
		t.Errorf("d1.isBefore(d2) = %v, want true", isB)
	}
	if isE := localDateIsEqual([]any{d1, CreateLocalDate(2026, 11, 8)}); isE != types.JavaBoolTrue && isE != true {
		t.Errorf("d1.isEqual(same) = %v, want true", isE)
	}
	if eq := localDateEquals([]any{d1, CreateLocalDate(2026, 11, 8)}); eq != types.JavaBoolTrue && eq != true {
		t.Errorf("d1.equals(same) = %v, want true", eq)
	}
	if eq := localDateEquals([]any{d1, d2}); eq != types.JavaBoolFalse && eq != false {
		t.Errorf("d1.equals(d2) = %v, want false", eq)
	}

	strObj := localDateToString([]any{d1}).(*object.Object)
	if str := object.GoStringFromStringObject(strObj); str != "2026-11-08" {
		t.Errorf("d1.toString() = %q, want '2026-11-08'", str)
	}

	// atTime(LocalTime)
	timeObj := CreateLocalTime(14, 30, 0, 0)
	dtObj := localDateAtTimeLocalTime([]any{d1, timeObj}).(*object.Object)
	y, m, d, h, min, s, n := getLocalDateTimeYMDHMSN(dtObj)
	if y != 2026 || m != 11 || d != 8 || h != 14 || min != 30 || s != 0 || n != 0 {
		t.Errorf("d1.atTime(14:30) = %d-%d-%d %d:%d:%d.%d", y, m, d, h, min, s, n)
	}

	// atStartOfDay()
	sodObj := localDateAtStartOfDay([]any{d1}).(*object.Object)
	y, m, d, h, min, s, n = getLocalDateTimeYMDHMSN(sodObj)
	if y != 2026 || m != 11 || d != 8 || h != 0 || min != 0 || s != 0 || n != 0 {
		t.Errorf("d1.atStartOfDay() = %d-%d-%d %d:%d:%d.%d", y, m, d, h, min, s, n)
	}

	// until(LocalDate, TemporalUnit)
	diffDays := localDateUntilUnit([]any{d1, d2, ghelpers.CreateTemporalUnit("DAYS")}).(int64)
	if diffDays != 1 {
		t.Errorf("d1.until(d2, DAYS) = %d, want 1", diffDays)
	}

	// until(LocalDate) -> Period
	pObj := localDateUntilPeriod([]any{d1, d2}).(*object.Object)
	if py, pm, pd, _ := getPeriodYMD(pObj); py != 0 || pm != 0 || pd != 1 {
		t.Errorf("d1.until(d2) = %d-%d-%d, want 0-0-1", py, pm, pd)
	}
}
