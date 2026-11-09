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
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestLocalTimeClinit(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()
	localTimeClinit(nil)

	for _, name := range []string{"MIN", "MAX", "MIDNIGHT", "NOON"} {
		s, ok := statics.QueryStatic("java/time/LocalTime", name)
		if !ok || s.Value == nil {
			t.Fatalf("LocalTime.%s not found in statics", name)
		}
	}

	sMin, ok := statics.QueryStatic("java/time/LocalTime", "MIN")
	if !ok || sMin.Value == nil {
		t.Fatal("MIN not found in statics")
	}
	minObj := sMin.Value.(*object.Object)
	h, m, s, n := getLocalTimeHMSN(minObj)
	if h != 0 || m != 0 || s != 0 || n != 0 {
		t.Errorf("MIN expected 00:00:00.0, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	sMax, ok := statics.QueryStatic("java/time/LocalTime", "MAX")
	if !ok || sMax.Value == nil {
		t.Fatal("MAX not found in statics")
	}
	maxObj := sMax.Value.(*object.Object)
	h, m, s, n = getLocalTimeHMSN(maxObj)
	if h != 23 || m != 59 || s != 59 || n != 999999999 {
		t.Errorf("MAX expected 23:59:59.999999999, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	sNoon, ok := statics.QueryStatic("java/time/LocalTime", "NOON")
	if !ok || sNoon.Value == nil {
		t.Fatal("NOON not found in statics")
	}
	noonObj := sNoon.Value.(*object.Object)
	h, m, s, n = getLocalTimeHMSN(noonObj)
	if h != 12 || m != 0 || s != 0 || n != 0 {
		t.Errorf("NOON expected 12:00:00.0, got %02d:%02d:%02d.%d", h, m, s, n)
	}
}

func TestLocalTimeOfFactories(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	// of(hour, minute)
	resHM := localTimeOfHM([]any{int64(14), int64(30)})
	ltHM, ok := resHM.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from of(14, 30), got %T", resHM)
	}
	h, m, s, n := getLocalTimeHMSN(ltHM)
	if h != 14 || m != 30 || s != 0 || n != 0 {
		t.Errorf("of(14, 30): got %d:%d:%d.%d", h, m, s, n)
	}

	// of(hour, minute, second)
	resHMS := localTimeOfHMS([]any{int64(14), int64(30), int64(45)})
	ltHMS, ok := resHMS.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from of(14, 30, 45), got %T", resHMS)
	}
	h, m, s, n = getLocalTimeHMSN(ltHMS)
	if h != 14 || m != 30 || s != 45 || n != 0 {
		t.Errorf("of(14, 30, 45): got %d:%d:%d.%d", h, m, s, n)
	}

	// of(hour, minute, second, nano)
	resHMSN := localTimeOfHMSN([]any{int64(14), int64(30), int64(45), int64(123456789)})
	ltHMSN, ok := resHMSN.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from of(14, 30, 45, 123456789), got %T", resHMSN)
	}
	h, m, s, n = getLocalTimeHMSN(ltHMSN)
	if h != 14 || m != 30 || s != 45 || n != 123456789 {
		t.Errorf("of(14, 30, 45, 123456789): got %d:%d:%d.%d", h, m, s, n)
	}

	// Invalid hour/minute/second/nano tests
	if _, isErr := localTimeOfHM([]any{int64(24), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for hour = 24")
	}
	if _, isErr := localTimeOfHM([]any{int64(-1), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for hour = -1")
	}
	if _, isErr := localTimeOfHMS([]any{int64(10), int64(60), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for minute = 60")
	}
	if _, isErr := localTimeOfHMSN([]any{int64(10), int64(0), int64(60), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for second = 60")
	}
	if _, isErr := localTimeOfHMSN([]any{int64(10), int64(0), int64(0), int64(1000000000)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for nano = 1000000000")
	}
}

func TestLocalTimeOfSecondAndNanoOfDay(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	// ofSecondOfDay(3665) -> 01:01:05
	resSec := localTimeOfSecondOfDay([]any{int64(3665)})
	ltSec := resSec.(*object.Object)
	h, m, s, n := getLocalTimeHMSN(ltSec)
	if h != 1 || m != 1 || s != 5 || n != 0 {
		t.Errorf("ofSecondOfDay(3665) expected 01:01:05, got %02d:%02d:%02d", h, m, s)
	}

	// ofNanoOfDay(3665000000123) -> 01:01:05.000000123
	resNano := localTimeOfNanoOfDay([]any{int64(3665000000123)})
	ltNano := resNano.(*object.Object)
	h, m, s, n = getLocalTimeHMSN(ltNano)
	if h != 1 || m != 1 || s != 5 || n != 123 {
		t.Errorf("ofNanoOfDay expected 01:01:05.000000123, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	// Boundary checks
	if _, isErr := localTimeOfSecondOfDay([]any{int64(86400)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for secondOfDay = 86400")
	}
	if _, isErr := localTimeOfSecondOfDay([]any{int64(-1)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for secondOfDay = -1")
	}
	if _, isErr := localTimeOfNanoOfDay([]any{int64(86400000000000)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for nanoOfDay = 86400000000000")
	}
}

func TestLocalTimeGetters(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	lt := createLocalTime(15, 42, 19, 500000000)

	if got := localTimeGetHour([]any{lt}).(int64); got != 15 {
		t.Errorf("getHour: expected 15, got %d", got)
	}
	if got := localTimeGetMinute([]any{lt}).(int64); got != 42 {
		t.Errorf("getMinute: expected 42, got %d", got)
	}
	if got := localTimeGetSecond([]any{lt}).(int64); got != 19 {
		t.Errorf("getSecond: expected 19, got %d", got)
	}
	if got := localTimeGetNano([]any{lt}).(int64); got != 500000000 {
		t.Errorf("getNano: expected 500000000, got %d", got)
	}

	secOfDay := localTimeToSecondOfDay([]any{lt}).(int64)
	expectedSec := int64(15*3600 + 42*60 + 19)
	if secOfDay != expectedSec {
		t.Errorf("toSecondOfDay: expected %d, got %d", expectedSec, secOfDay)
	}

	nanoOfDay := localTimeToNanoOfDay([]any{lt}).(int64)
	expectedNano := expectedSec*1000000000 + 500000000
	if nanoOfDay != expectedNano {
		t.Errorf("toNanoOfDay: expected %d, got %d", expectedNano, nanoOfDay)
	}
}

func TestLocalTimeWith(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	lt := createLocalTime(10, 20, 30, 400)

	// withHour
	lt2 := localTimeWithHour([]any{lt, int64(15)}).(*object.Object)
	h, m, s, n := getLocalTimeHMSN(lt2)
	if h != 15 || m != 20 || s != 30 || n != 400 {
		t.Errorf("withHour: got %d:%d:%d.%d", h, m, s, n)
	}

	// withMinute
	lt3 := localTimeWithMinute([]any{lt, int64(45)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(lt3)
	if h != 10 || m != 45 || s != 30 || n != 400 {
		t.Errorf("withMinute: got %d:%d:%d.%d", h, m, s, n)
	}

	// withSecond
	lt4 := localTimeWithSecond([]any{lt, int64(55)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(lt4)
	if h != 10 || m != 20 || s != 55 || n != 400 {
		t.Errorf("withSecond: got %d:%d:%d.%d", h, m, s, n)
	}

	// withNano
	lt5 := localTimeWithNano([]any{lt, int64(999999)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(lt5)
	if h != 10 || m != 20 || s != 30 || n != 999999 {
		t.Errorf("withNano: got %d:%d:%d.%d", h, m, s, n)
	}
}

func TestLocalTimeArithmetic(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	lt := createLocalTime(23, 30, 45, 500000000)

	// plusHours(2) -> 01:30:45.5
	p1 := localTimePlusHours([]any{lt, int64(2)}).(*object.Object)
	h, m, s, n := getLocalTimeHMSN(p1)
	if h != 1 || m != 30 || s != 45 || n != 500000000 {
		t.Errorf("plusHours(2): expected 01:30:45.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	// plusMinutes(40) -> 00:10:45.5
	p2 := localTimePlusMinutes([]any{lt, int64(40)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(p2)
	if h != 0 || m != 10 || s != 45 || n != 500000000 {
		t.Errorf("plusMinutes(40): expected 00:10:45.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	// plusSeconds(30) -> 23:31:15.5
	p3 := localTimePlusSeconds([]any{lt, int64(30)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(p3)
	if h != 23 || m != 31 || s != 15 || n != 500000000 {
		t.Errorf("plusSeconds(30): expected 23:31:15.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	// minusHours(25) -> 22:30:45.5
	m1 := localTimeMinusHours([]any{lt, int64(25)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(m1)
	if h != 22 || m != 30 || s != 45 || n != 500000000 {
		t.Errorf("minusHours(25): expected 22:30:45.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	// minusMinutes(45) -> 22:45:45.5
	m2 := localTimeMinusMinutes([]any{lt, int64(45)}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(m2)
	if h != 22 || m != 45 || s != 45 || n != 500000000 {
		t.Errorf("minusMinutes(45): expected 22:45:45.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}
}

func TestLocalTimeComparisons(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	t1 := createLocalTime(10, 30, 0, 0)
	t2 := createLocalTime(10, 30, 0, 0)
	t3 := createLocalTime(11, 0, 0, 0)

	// equals
	if localTimeEquals([]any{t1, t2}) != types.JavaBoolTrue {
		t.Error("t1 should equal t2")
	}
	if localTimeEquals([]any{t1, t3}) != types.JavaBoolFalse {
		t.Error("t1 should not equal t3")
	}

	// compareTo
	if cmp := localTimeCompareTo([]any{t1, t2}).(int64); cmp != 0 {
		t.Errorf("compareTo equal: expected 0, got %d", cmp)
	}
	if cmp := localTimeCompareTo([]any{t1, t3}).(int64); cmp >= 0 {
		t.Errorf("compareTo earlier: expected < 0, got %d", cmp)
	}
	if cmp := localTimeCompareTo([]any{t3, t1}).(int64); cmp <= 0 {
		t.Errorf("compareTo later: expected > 0, got %d", cmp)
	}

	// isBefore / isAfter
	if localTimeIsBefore([]any{t1, t3}) != types.JavaBoolTrue {
		t.Error("t1 isBefore t3 should be true")
	}
	if localTimeIsAfter([]any{t3, t1}) != types.JavaBoolTrue {
		t.Error("t3 isAfter t1 should be true")
	}

	// hashCode
	h1 := localTimeHashCode([]any{t1}).(int64)
	h2 := localTimeHashCode([]any{t2}).(int64)
	if h1 != h2 {
		t.Errorf("hashCodes should match: %d vs %d", h1, h2)
	}
}

func TestLocalTimeParseAndFormat(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	tests := []struct {
		input      string
		expected   string
		h, m, s, n int64
	}{
		{"12:30", "12:30", 12, 30, 0, 0},
		{"12:30:45", "12:30:45", 12, 30, 45, 0},
		{"12:30:45.123", "12:30:45.123", 12, 30, 45, 123000000},
		{"12:30:45.123456", "12:30:45.123456", 12, 30, 45, 123456000},
		{"12:30:45.123456789", "12:30:45.123456789", 12, 30, 45, 123456789},
		{"T08:05:00", "08:05", 8, 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			strObj := object.StringObjectFromGoString(tt.input)
			res := localTimeParse([]any{strObj})
			lt, ok := res.(*object.Object)
			if !ok {
				t.Fatalf("parse failed for %s: %v", tt.input, res)
			}
			h, m, s, n := getLocalTimeHMSN(lt)
			if h != tt.h || m != tt.m || s != tt.s || n != tt.n {
				t.Errorf("parse(%s): got (%d, %d, %d, %d), expected (%d, %d, %d, %d)",
					tt.input, h, m, s, n, tt.h, tt.m, tt.s, tt.n)
			}

			// toString
			strRes := localTimeToString([]any{lt}).(*object.Object)
			formatted := object.GoStringFromStringObject(strRes)
			if formatted != tt.expected {
				t.Errorf("toString: expected %s, got %s", tt.expected, formatted)
			}
		})
	}
}

func TestLocalTimeTruncatedToAndUntil(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	lt1 := createLocalTime(10, 30, 45, 500000000)
	lt2 := createLocalTime(12, 45, 55, 600000000)

	// truncatedTo SECONDS
	unitSec := ghelpers.CreateTemporalUnit("Seconds")
	trunc := localTimeTruncatedTo([]any{lt1, unitSec}).(*object.Object)
	h, m, s, n := getLocalTimeHMSN(trunc)
	if h != 10 || m != 30 || s != 45 || n != 0 {
		t.Errorf("truncatedTo(Seconds): expected 10:30:45.0, got %d:%d:%d.%d", h, m, s, n)
	}

	// until HOURS
	unitHours := ghelpers.CreateTemporalUnit("Hours")
	untilH := localTimeUntilUnit([]any{lt1, lt2, unitHours}).(int64)
	if untilH != 2 {
		t.Errorf("until(Hours): expected 2, got %d", untilH)
	}

	// until MINUTES
	unitMin := ghelpers.CreateTemporalUnit("Minutes")
	untilM := localTimeUntilUnit([]any{lt1, lt2, unitMin}).(int64)
	if untilM != 135 {
		t.Errorf("until(Minutes): expected 135, got %d", untilM)
	}
}

func TestLocalTimeFieldsAndRanges(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	lt := createLocalTime(13, 45, 30, 250000000)

	cfHourOfDay := getChronoField("HOUR_OF_DAY")
	cfMinuteOfHour := getChronoField("MINUTE_OF_HOUR")
	cfSecondOfMinute := getChronoField("SECOND_OF_MINUTE")
	cfNanoOfSecond := getChronoField("NANO_OF_SECOND")
	cfNanoOfDay := getChronoField("NANO_OF_DAY")
	cfHourOfAmPm := getChronoField("HOUR_OF_AMPM")
	cfClockHourOfAmPm := getChronoField("CLOCK_HOUR_OF_AMPM")
	cfClockHourOfDay := getChronoField("CLOCK_HOUR_OF_DAY")
	cfAmPmOfDay := getChronoField("AMPM_OF_DAY")

	// isSupportedField
	if localTimeIsSupportedField([]any{lt, cfHourOfDay}) != types.JavaBoolTrue {
		t.Error("HOUR_OF_DAY should be supported")
	}
	cfYear := getChronoField("YEAR")
	if localTimeIsSupportedField([]any{lt, cfYear}) != types.JavaBoolFalse {
		t.Error("YEAR should not be supported")
	}

	// get / getLong
	if got := localTimeGet([]any{lt, cfHourOfDay}).(int64); got != 13 {
		t.Errorf("get(HOUR_OF_DAY): expected 13, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfMinuteOfHour}).(int64); got != 45 {
		t.Errorf("get(MINUTE_OF_HOUR): expected 45, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfSecondOfMinute}).(int64); got != 30 {
		t.Errorf("get(SECOND_OF_MINUTE): expected 30, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfNanoOfSecond}).(int64); got != 250000000 {
		t.Errorf("get(NANO_OF_SECOND): expected 250000000, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfHourOfAmPm}).(int64); got != 1 {
		t.Errorf("get(HOUR_OF_AMPM): expected 1, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfClockHourOfAmPm}).(int64); got != 1 {
		t.Errorf("get(CLOCK_HOUR_OF_AMPM): expected 1, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfClockHourOfDay}).(int64); got != 13 {
		t.Errorf("get(CLOCK_HOUR_OF_DAY): expected 13, got %d", got)
	}
	if got := localTimeGet([]any{lt, cfAmPmOfDay}).(int64); got != 1 {
		t.Errorf("get(AMPM_OF_DAY): expected 1, got %d", got)
	}
	if got := localTimeGetLong([]any{lt, cfNanoOfDay}).(int64); got != getLocalTimeNanoOfDay(lt) {
		t.Errorf("getLong(NANO_OF_DAY): got %d, expected %d", got, getLocalTimeNanoOfDay(lt))
	}

	// range
	r := localTimeRange([]any{lt, cfHourOfDay}).(*object.Object)
	if r.FieldTable["maxLargest"].Fvalue.(int64) != 23 {
		t.Errorf("range(HOUR_OF_DAY) maxLargest expected 23, got %v", r.FieldTable["maxLargest"].Fvalue)
	}

	// with(field, newValue)
	wHour := localTimeWithField([]any{lt, cfHourOfDay, int64(8)}).(*object.Object)
	h, _, _, _ := getLocalTimeHMSN(wHour)
	if h != 8 {
		t.Errorf("with(HOUR_OF_DAY, 8): expected hour 8, got %d", h)
	}
}

func TestLocalTimeAdjustIntoAndQuery(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()

	lt := createLocalTime(9, 15, 0, 0)
	target := createLocalTime(0, 0, 0, 0)

	adjusted := localTimeAdjustInto([]any{lt, target}).(*object.Object)
	h, m, s, n := getLocalTimeHMSN(adjusted)
	if h != 9 || m != 15 || s != 0 || n != 0 {
		t.Errorf("adjustInto: expected 09:15:00, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	// query
	qLocalTime := makeTemporalQueryObject("localTime")
	resQ := localTimeQuery([]any{lt, qLocalTime}).(*object.Object)
	if resQ != lt {
		t.Error("query(localTime) should return self")
	}

	qPrecision := makeTemporalQueryObject("precision")
	resPrec := localTimeQuery([]any{lt, qPrecision}).(*object.Object)
	if getTemporalUnitName(resPrec) != "Nanos" {
		t.Errorf("query(precision): expected Nanos, got %v", getTemporalUnitName(resPrec))
	}

	qLocalDate := makeTemporalQueryObject("localDate")
	if resNull := localTimeQuery([]any{lt, qLocalDate}); resNull != nil {
		t.Errorf("query(localDate) expected nil, got %v", resNull)
	}
}

func TestLocalTimeAtDateAndOffset(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()
	Load_Time_LocalDate()

	lt := createLocalTime(14, 30, 0, 0)
	ld := CreateLocalDate(2026, 11, 5)

	// atDate
	ldt := localTimeAtDate([]any{lt, ld}).(*object.Object)
	dObj := getChronoLocalDateTimeDate(ldt)
	tObj := getChronoLocalDateTimeTime(ldt)
	y, m, d, _ := getLocalDateYMD(dObj)
	th, tm, _, _ := getLocalTimeHMSN(tObj)
	if y != 2026 || m != 11 || d != 5 || th != 14 || tm != 30 {
		t.Errorf("atDate: got %04d-%02d-%02d %02d:%02d", y, m, d, th, tm)
	}

	// atOffset
	zo := createZoneOffset(0)
	ot := localTimeAtOffset([]any{lt, zo}).(*object.Object)
	if ot == nil {
		t.Fatal("atOffset returned nil")
	}

	// toEpochSecond
	epochSec := localTimeToEpochSecond([]any{lt, ld, zo}).(int64)
	expectedSec := ymdToEpochDay(2026, 11, 5)*86400 + 14*3600 + 30*60
	if epochSec != expectedSec {
		t.Errorf("toEpochSecond: expected %d, got %d", expectedSec, epochSec)
	}
}

func TestLocalTimePlusMinusAmount(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_LocalTime()
	Load_Time_Duration()

	lt := createLocalTime(12, 0, 0, 0)
	dur := createDuration(3665, 500000000) // 1 hour, 1 min, 5.5 sec

	plusDur := localTimePlusAmount([]any{lt, dur}).(*object.Object)
	h, m, s, n := getLocalTimeHMSN(plusDur)
	if h != 13 || m != 1 || s != 5 || n != 500000000 {
		t.Errorf("plus(Duration): expected 13:01:05.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}

	minusDur := localTimeMinusAmount([]any{lt, dur}).(*object.Object)
	h, m, s, n = getLocalTimeHMSN(minusDur)
	if h != 10 || m != 58 || s != 54 || n != 500000000 {
		t.Errorf("minus(Duration): expected 10:58:54.5, got %02d:%02d:%02d.%d", h, m, s, n)
	}
}
