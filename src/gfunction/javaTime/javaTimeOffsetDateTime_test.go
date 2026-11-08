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

func offsetDateTimeTestInit() {
	globals.InitGlobals("test")
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_LocalDateTime()
	Load_Time_ZoneOffset()
	Load_Time_ZoneId()
	Load_Time_Instant()
	Load_Time_OffsetDateTime()
	Load_Time_ZonedDateTime()
}

func TestOffsetDateTimeRegistration(t *testing.T) {
	offsetDateTimeTestInit()

	signatures := []string{
		"java/time/OffsetDateTime.<clinit>()V",
		"java/time/OffsetDateTime.<init>()V",
		"java/time/OffsetDateTime.<init>(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;)V",
		"java/time/OffsetDateTime.now()Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.now(Ljava/time/ZoneId;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.now(Ljava/time/Clock;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.of(Ljava/time/LocalDate;Ljava/time/LocalTime;Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.of(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.of(IIIIIIILjava/time/ZoneOffset;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.parse(Ljava/lang/CharSequence;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.timeLineOrder()Ljava/util/Comparator;",
		"java/time/OffsetDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z",
		"java/time/OffsetDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z",
		"java/time/OffsetDateTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;",
		"java/time/OffsetDateTime.get(Ljava/time/temporal/TemporalField;)I",
		"java/time/OffsetDateTime.getLong(Ljava/time/temporal/TemporalField;)J",
		"java/time/OffsetDateTime.getOffset()Ljava/time/ZoneOffset;",
		"java/time/OffsetDateTime.withOffsetSameLocal(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withOffsetSameInstant(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.toLocalDateTime()Ljava/time/LocalDateTime;",
		"java/time/OffsetDateTime.toLocalDate()Ljava/time/LocalDate;",
		"java/time/OffsetDateTime.getYear()I",
		"java/time/OffsetDateTime.getMonthValue()I",
		"java/time/OffsetDateTime.getMonth()Ljava/time/Month;",
		"java/time/OffsetDateTime.getDayOfMonth()I",
		"java/time/OffsetDateTime.getDayOfYear()I",
		"java/time/OffsetDateTime.getDayOfWeek()Ljava/time/DayOfWeek;",
		"java/time/OffsetDateTime.toLocalTime()Ljava/time/LocalTime;",
		"java/time/OffsetDateTime.getHour()I",
		"java/time/OffsetDateTime.getMinute()I",
		"java/time/OffsetDateTime.getSecond()I",
		"java/time/OffsetDateTime.getNano()I",
		"java/time/OffsetDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withYear(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withMonth(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withDayOfMonth(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withDayOfYear(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withHour(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withMinute(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withSecond(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.withNano(I)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusYears(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusMonths(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusWeeks(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusDays(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusHours(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusMinutes(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusSeconds(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.plusNanos(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusYears(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusMonths(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusWeeks(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusDays(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusHours(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusMinutes(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusSeconds(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.minusNanos(J)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;",
		"java/time/OffsetDateTime.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;",
		"java/time/OffsetDateTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J",
		"java/time/OffsetDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;",
		"java/time/OffsetDateTime.atZoneSameInstant(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;",
		"java/time/OffsetDateTime.atZoneSimilarLocal(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;",
		"java/time/OffsetDateTime.toOffsetTime()Ljava/time/OffsetTime;",
		"java/time/OffsetDateTime.toEpochSecond()J",
		"java/time/OffsetDateTime.toInstant()Ljava/time/Instant;",
		"java/time/OffsetDateTime.toZonedDateTime()Ljava/time/ZonedDateTime;",
		"java/time/OffsetDateTime.compareTo(Ljava/time/OffsetDateTime;)I",
		"java/time/OffsetDateTime.compareTo(Ljava/lang/Object;)I",
		"java/time/OffsetDateTime.isAfter(Ljava/time/OffsetDateTime;)Z",
		"java/time/OffsetDateTime.isBefore(Ljava/time/OffsetDateTime;)Z",
		"java/time/OffsetDateTime.isEqual(Ljava/time/OffsetDateTime;)Z",
		"java/time/OffsetDateTime.equals(Ljava/lang/Object;)Z",
		"java/time/OffsetDateTime.hashCode()I",
		"java/time/OffsetDateTime.toString()Ljava/lang/String;",
		"java/time/OffsetDateTime.writeExternal(Ljava/io/DataOutput;)V",
		"java/time/OffsetDateTime.readExternal(Ljava/io/DataInput;)Ljava/time/OffsetDateTime;",
		"java/time/OffsetDateTime.write(Ljava/io/DataOutput;)V",
		"java/time/OffsetDateTime$1.compare(Ljava/lang/Object;Ljava/lang/Object;)I",
		"java/time/OffsetDateTime$1.compare(Ljava/time/OffsetDateTime;Ljava/time/OffsetDateTime;)I",
	}

	for _, sig := range signatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Missing signature: %s", sig)
		}
	}
}

func TestOffsetDateTimeFactoriesAndGetters(t *testing.T) {
	offsetDateTimeTestInit()

	// Check MIN and MAX statics
	minField := statics.GetStaticValue(classNameOffsetDateTime, "MIN")
	if minField == nil {
		t.Fatalf("Static MIN field missing")
	}
	maxField := statics.GetStaticValue(classNameOffsetDateTime, "MAX")
	if maxField == nil {
		t.Fatalf("Static MAX field missing")
	}

	dateObj := CreateLocalDate(2023, 10, 15)
	timeObj := createLocalTime(14, 30, 45, 500)
	offsetObj := createZoneOffset(2 * 3600) // +02:00

	odtRes := offsetDateTimeOfDate([]any{dateObj, timeObj, offsetObj})
	if err, isErr := odtRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("offsetDateTimeOfDate error: %v", err)
	}
	odt := odtRes.(*object.Object)

	// Test Getters
	if yr := offsetDateTimeGetYear([]any{odt}); yr != int64(2023) {
		t.Errorf("Expected year 2023, got %v", yr)
	}
	if mo := offsetDateTimeGetMonthValue([]any{odt}); mo != int64(10) {
		t.Errorf("Expected month 10, got %v", mo)
	}
	if day := offsetDateTimeGetDayOfMonth([]any{odt}); day != int64(15) {
		t.Errorf("Expected day 15, got %v", day)
	}
	if hr := offsetDateTimeGetHour([]any{odt}); hr != int64(14) {
		t.Errorf("Expected hour 14, got %v", hr)
	}
	if min := offsetDateTimeGetMinute([]any{odt}); min != int64(30) {
		t.Errorf("Expected minute 30, got %v", min)
	}
	if sec := offsetDateTimeGetSecond([]any{odt}); sec != int64(45) {
		t.Errorf("Expected second 45, got %v", sec)
	}
	if nano := offsetDateTimeGetNano([]any{odt}); nano != int64(500) {
		t.Errorf("Expected nano 500, got %v", nano)
	}

	// toString
	strRes := offsetDateTimeToString([]any{odt})
	str := object.GoStringFromStringObject(strRes.(*object.Object))
	if str != "2023-10-15T14:30:45.000000500+02:00" {
		t.Errorf("Unexpected toString output: %s", str)
	}
}

func TestOffsetDateTimeParseAndConversions(t *testing.T) {
	offsetDateTimeTestInit()

	strObj := object.StringObjectFromGoString("2024-05-20T10:15:30+03:00")
	parsedRes := offsetDateTimeParse([]any{strObj})
	if err, isErr := parsedRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("offsetDateTimeParse error: %v", err)
	}
	odt := parsedRes.(*object.Object)

	if yr := offsetDateTimeGetYear([]any{odt}); yr != int64(2024) {
		t.Errorf("Expected year 2024, got %v", yr)
	}
	if hr := offsetDateTimeGetHour([]any{odt}); hr != int64(10) {
		t.Errorf("Expected hour 10, got %v", hr)
	}

	// Conversions
	ldt := offsetDateTimeToLocalDateTime([]any{odt}).(*object.Object)
	if ldt == nil {
		t.Errorf("toLocalDateTime returned nil")
	}
	localDate := offsetDateTimeToLocalDate([]any{odt}).(*object.Object)
	if localDate == nil {
		t.Errorf("toLocalDate returned nil")
	}
	localTime := offsetDateTimeToLocalTime([]any{odt}).(*object.Object)
	if localTime == nil {
		t.Errorf("toLocalTime returned nil")
	}
	offsetTime := offsetDateTimeToOffsetTime([]any{odt}).(*object.Object)
	if offsetTime == nil {
		t.Errorf("toOffsetTime returned nil")
	}
	inst := offsetDateTimeToInstant([]any{odt}).(*object.Object)
	if inst == nil {
		t.Errorf("toInstant returned nil")
	}
	zdt := offsetDateTimeToZonedDateTime([]any{odt}).(*object.Object)
	if zdt == nil {
		t.Errorf("toZonedDateTime returned nil")
	}

	epochSec := offsetDateTimeToEpochSecond([]any{odt}).(int64)
	if epochSec <= 0 {
		t.Errorf("toEpochSecond invalid: %v", epochSec)
	}

	// from TemporalAccessor
	fromRes := offsetDateTimeFrom([]any{odt})
	if fromRes != odt {
		t.Errorf("from(OffsetDateTime) should return self")
	}
	fromZdtRes := offsetDateTimeFrom([]any{zdt})
	if fromZdtRes == nil {
		t.Errorf("from(ZonedDateTime) returned nil")
	}
}

func TestOffsetDateTimeWithAndAdjustments(t *testing.T) {
	offsetDateTimeTestInit()

	odt := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T10:15:30+03:00")}).(*object.Object)

	// withOffsetSameLocal (+01:00)
	off1 := createZoneOffset(3600)
	odtSameLocal := offsetDateTimeWithOffsetSameLocal([]any{odt, off1}).(*object.Object)
	if hr := offsetDateTimeGetHour([]any{odtSameLocal}); hr != int64(10) {
		t.Errorf("withOffsetSameLocal hour should still be 10, got %v", hr)
	}
	offRes := offsetDateTimeGetOffset([]any{odtSameLocal}).(*object.Object)
	sec, _ := getZoneOffsetTotalSecondsVal(offRes)
	if sec != 3600 {
		t.Errorf("Expected offset 3600, got %v", sec)
	}

	// withOffsetSameInstant (+01:00) -> 10:15:30 +03:00 is 07:15:30 UTC -> 08:15:30 +01:00
	odtSameInstant := offsetDateTimeWithOffsetSameInstant([]any{odt, off1}).(*object.Object)
	if hr := offsetDateTimeGetHour([]any{odtSameInstant}); hr != int64(8) {
		t.Errorf("withOffsetSameInstant hour should be 8, got %v", hr)
	}

	// withYear, withMonth, withDayOfMonth, etc.
	wYear := offsetDateTimeWithYear([]any{odt, int64(2030)}).(*object.Object)
	if yr := offsetDateTimeGetYear([]any{wYear}); yr != int64(2030) {
		t.Errorf("Expected year 2030, got %v", yr)
	}

	wHour := offsetDateTimeWithHour([]any{odt, int64(22)}).(*object.Object)
	if hr := offsetDateTimeGetHour([]any{wHour}); hr != int64(22) {
		t.Errorf("Expected hour 22, got %v", hr)
	}
}

func TestOffsetDateTimeArithmetic(t *testing.T) {
	offsetDateTimeTestInit()

	odt := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T10:15:30+00:00")}).(*object.Object)

	// plusDays(5) -> 2024-05-25
	plus5 := offsetDateTimePlusDays([]any{odt, int64(5)}).(*object.Object)
	if day := offsetDateTimeGetDayOfMonth([]any{plus5}); day != int64(25) {
		t.Errorf("Expected day 25, got %v", day)
	}

	// minusDays(5)
	minus5 := offsetDateTimeMinusDays([]any{plus5, int64(5)}).(*object.Object)
	if day := offsetDateTimeGetDayOfMonth([]any{minus5}); day != int64(20) {
		t.Errorf("Expected day 20, got %v", day)
	}

	// plusHours(3) -> 13
	plus3h := offsetDateTimePlusHours([]any{odt, int64(3)}).(*object.Object)
	if hr := offsetDateTimeGetHour([]any{plus3h}); hr != int64(13) {
		t.Errorf("Expected hour 13, got %v", hr)
	}

	// minusHours(3)
	minus3h := offsetDateTimeMinusHours([]any{plus3h, int64(3)}).(*object.Object)
	if hr := offsetDateTimeGetHour([]any{minus3h}); hr != int64(10) {
		t.Errorf("Expected hour 10, got %v", hr)
	}
}

func TestOffsetDateTimeComparisonsAndEquality(t *testing.T) {
	offsetDateTimeTestInit()

	odt1 := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T10:15:30+02:00")}).(*object.Object)
	odt2 := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T08:15:30Z")}).(*object.Object) // Same instant
	odt3 := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T11:15:30+02:00")}).(*object.Object)

	// isEqual
	if eq := offsetDateTimeIsEqual([]any{odt1, odt2}); eq != types.JavaBoolTrue {
		t.Errorf("odt1 and odt2 should be isEqual (same instant)")
	}

	// equals (different offsets -> not equals)
	if eq := offsetDateTimeEquals([]any{odt1, odt2}); eq != types.JavaBoolFalse {
		t.Errorf("odt1 and odt2 have different offsets, equals should be false")
	}

	// isBefore / isAfter
	if b := offsetDateTimeIsBefore([]any{odt1, odt3}); b != types.JavaBoolTrue {
		t.Errorf("odt1 should be before odt3")
	}
	if a := offsetDateTimeIsAfter([]any{odt3, odt1}); a != types.JavaBoolTrue {
		t.Errorf("odt3 should be after odt1")
	}

	// timeline comparator
	cmpObj := offsetDateTimeTimelineOrder(nil).(*object.Object)
	resCmp := offsetDateTimeTimelineCompare([]any{cmpObj, odt1, odt2})
	if resCmp != int64(0) {
		t.Errorf("Timeline order should compare instant equally, got %v", resCmp)
	}
}

func TestOffsetDateTimeQueriesAndUntil(t *testing.T) {
	offsetDateTimeTestInit()

	odt1 := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T10:00:00Z")}).(*object.Object)
	odt2 := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-22T12:00:00Z")}).(*object.Object)

	// query offset
	ensureTemporalQueriesInited()
	offRes := offsetDateTimeQuery([]any{odt1, offsetQueryInstance}).(*object.Object)
	if sec, _ := getZoneOffsetTotalSecondsVal(offRes); sec != 0 {
		t.Errorf("Expected offset 0, got %d", sec)
	}

	// query zoneId (must be nil for OffsetDateTime)
	zoneRes := offsetDateTimeQuery([]any{odt1, zoneIdQueryInstance})
	if zoneRes != nil {
		t.Errorf("query zoneId on OffsetDateTime must be nil")
	}

	// until DAYS
	unitDays := ghelpers.CreateTemporalUnit("DAYS")
	daysRes := offsetDateTimeUntil([]any{odt1, odt2, unitDays})
	if daysRes != int64(2) {
		t.Errorf("Expected 2 days until, got %v", daysRes)
	}

	// until HOURS
	unitHours := ghelpers.CreateTemporalUnit("HOURS")
	hoursRes := offsetDateTimeUntil([]any{odt1, odt2, unitHours})
	if hoursRes != int64(50) {
		t.Errorf("Expected 50 hours until, got %v", hoursRes)
	}
}

func TestOffsetDateTimeZoneConversionsAndFormat(t *testing.T) {
	offsetDateTimeTestInit()

	odt := offsetDateTimeParse([]any{object.StringObjectFromGoString("2024-05-20T10:15:30+02:00")}).(*object.Object)
	parisZone := makeZoneIdFromID("Europe/Paris")

	zdtSameInstant := offsetDateTimeAtZoneSameInstant([]any{odt, parisZone}).(*object.Object)
	if zdtSameInstant == nil {
		t.Errorf("atZoneSameInstant returned nil")
	}

	zdtSimilarLocal := offsetDateTimeAtZoneSimilarLocal([]any{odt, parisZone}).(*object.Object)
	if zdtSimilarLocal == nil {
		t.Errorf("atZoneSimilarLocal returned nil")
	}

	// format
	formatter := createDateTimeFormatterObject("yyyy-MM-dd", nil, nil, nil)
	formattedStr := offsetDateTimeFormat([]any{odt, formatter})
	if strObj, ok := formattedStr.(*object.Object); ok {
		str := object.GoStringFromStringObject(strObj)
		if str != "2024-05-20" {
			t.Errorf("Expected formatted string '2024-05-20', got '%s'", str)
		}
	}
}
