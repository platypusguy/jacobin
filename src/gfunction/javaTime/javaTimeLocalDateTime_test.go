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

func initTestEnvironment() {
	globals.InitGlobals("test")
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_LocalDateTime()
	Load_Time_Month()
	Load_Time_DayOfWeek()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
	Load_Time_Temporal_ValueRange()
	Load_Time_Chrono_IsoChronology()

	chronoUnitClinit(nil)
	chronoFieldClinit(nil)
	monthClinit(nil)
	dayOfWeekClinit(nil)
	localDateClinit(nil)
	localTimeClinit(nil)
	localDateTimeClinit(nil)
}

func isTrue(v any) bool {
	switch b := v.(type) {
	case bool:
		return b
	case int64:
		return b != 0
	case int32:
		return b != 0
	case int:
		return b != 0
	default:
		return false
	}
}

func getLocalDateTimeYMDHMSN(dt *object.Object) (year, month, day, hour, minute, second, nano int64) {
	dObj := getLocalDateTimeDate(dt)
	tObj := getLocalDateTimeTime(dt)
	year, month, day, _ = getLocalDateYMD(dObj)
	hour, minute, second, nano = getLocalTimeHMSN(tObj)
	return
}

func TestLocalDateTimeRegistration(t *testing.T) {
	initTestEnvironment()

	expectedMethods := []string{
		"java/time/LocalDateTime.<clinit>()V",
		"java/time/LocalDateTime.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.atOffset(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;",
		"java/time/LocalDateTime.atZone(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;",
		"java/time/LocalDateTime.atZone(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;",
		"java/time/LocalDateTime.compareTo(Ljava/time/chrono/ChronoLocalDateTime;)I",
		"java/time/LocalDateTime.compareTo(Ljava/lang/Object;)I",
		"java/time/LocalDateTime.equals(Ljava/lang/Object;)Z",
		"java/time/LocalDateTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;",
		"java/time/LocalDateTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.get(Ljava/time/temporal/TemporalField;)I",
		"java/time/LocalDateTime.getDayOfMonth()I",
		"java/time/LocalDateTime.getDayOfWeek()Ljava/time/DayOfWeek;",
		"java/time/LocalDateTime.getDayOfYear()I",
		"java/time/LocalDateTime.getHour()I",
		"java/time/LocalDateTime.getLong(Ljava/time/temporal/TemporalField;)J",
		"java/time/LocalDateTime.getMinute()I",
		"java/time/LocalDateTime.getMonth()Ljava/time/Month;",
		"java/time/LocalDateTime.getMonthValue()I",
		"java/time/LocalDateTime.getNano()I",
		"java/time/LocalDateTime.getSecond()I",
		"java/time/LocalDateTime.getYear()I",
		"java/time/LocalDateTime.hashCode()I",
		"java/time/LocalDateTime.isAfter(Ljava/time/chrono/ChronoLocalDateTime;)Z",
		"java/time/LocalDateTime.isBefore(Ljava/time/chrono/ChronoLocalDateTime;)Z",
		"java/time/LocalDateTime.isEqual(Ljava/time/chrono/ChronoLocalDateTime;)Z",
		"java/time/LocalDateTime.isSupported(Ljava/time/temporal/TemporalField;)Z",
		"java/time/LocalDateTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z",
		"java/time/LocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDateTime;",
		"java/time/LocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDateTime;",
		"java/time/LocalDateTime.minusDays(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusHours(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusMinutes(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusMonths(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusNanos(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusSeconds(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusWeeks(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.minusYears(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.now()Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.now(Ljava/time/Clock;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.now(Ljava/time/ZoneId;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(IIIII)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(IIIIII)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(IIIIIII)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(ILjava/time/Month;III)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(ILjava/time/Month;IIII)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(ILjava/time/Month;IIIII)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.of(Ljava/time/LocalDate;Ljava/time/LocalTime;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.ofEpochSecond(JILjava/time/ZoneOffset;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.parse(Ljava/lang/CharSequence;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/chrono/ChronoLocalDateTime;",
		"java/time/LocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoLocalDateTime;",
		"java/time/LocalDateTime.plusDays(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusHours(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusMinutes(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusMonths(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusNanos(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusSeconds(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusWeeks(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.plusYears(J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;",
		"java/time/LocalDateTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;",
		"java/time/LocalDateTime.toEpochSecond(Ljava/time/ZoneOffset;)J",
		"java/time/LocalDateTime.toInstant(Ljava/time/ZoneOffset;)Ljava/time/Instant;",
		"java/time/LocalDateTime.toLocalDate()Ljava/time/LocalDate;",
		"java/time/LocalDateTime.toLocalDate()Ljava/time/chrono/ChronoLocalDate;",
		"java/time/LocalDateTime.toLocalTime()Ljava/time/LocalTime;",
		"java/time/LocalDateTime.toString()Ljava/lang/String;",
		"java/time/LocalDateTime.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J",
		"java/time/LocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/chrono/ChronoLocalDateTime;",
		"java/time/LocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;",
		"java/time/LocalDateTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/chrono/ChronoLocalDateTime;",
		"java/time/LocalDateTime.withDayOfMonth(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withDayOfYear(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withHour(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withMinute(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withMonth(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withNano(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withSecond(I)Ljava/time/LocalDateTime;",
		"java/time/LocalDateTime.withYear(I)Ljava/time/LocalDateTime;",
	}

	for _, sig := range expectedMethods {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected method signature %s was not registered", sig)
		}
	}
}

func TestLocalDateTimeClinitAndStatics(t *testing.T) {
	initTestEnvironment()

	for _, name := range []string{"MIN", "MAX"} {
		s, ok := statics.QueryStatic("java/time/LocalDateTime", name)
		if !ok || s.Value == nil {
			t.Fatalf("LocalDateTime.%s not found in statics", name)
		}
	}

	sMin, _ := statics.QueryStatic("java/time/LocalDateTime", "MIN")
	minObj := sMin.Value.(*object.Object)
	y, m, d, h, minute, sec, n := getLocalDateTimeYMDHMSN(minObj)
	if y != -999999999 || m != 1 || d != 1 || h != 0 || minute != 0 || sec != 0 || n != 0 {
		t.Errorf("LocalDateTime.MIN expected -999999999-01-01T00:00:00, got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}

	sMax, _ := statics.QueryStatic("java/time/LocalDateTime", "MAX")
	maxObj := sMax.Value.(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(maxObj)
	if y != 999999999 || m != 12 || d != 31 || h != 23 || minute != 59 || sec != 59 || n != 999999999 {
		t.Errorf("LocalDateTime.MAX expected 999999999-12-31T23:59:59.999999999, got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}
}

// TestLocalDateTimeOfWithMonthEnum verifies the Go equivalent of LocalDateTime.of(2026, Month.JUNE, 24, 14, 30, 15)
func TestLocalDateTimeOfWithMonthEnum(t *testing.T) {
	initTestEnvironment()

	// Month.JUNE
	sJune, ok := statics.QueryStatic("java/time/Month", "JUNE")
	if !ok || sJune.Value == nil {
		t.Fatal("Month.JUNE not found in statics")
	}
	monthJune := sJune.Value.(*object.Object)

	// LocalDateTime.of(2026, Month.JUNE, 24, 14, 30, 15)
	res := localDateTimeOfYMonthHMS([]any{int64(2026), monthJune, int64(24), int64(14), int64(30), int64(15)})
	dt, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from of(2026, Month.JUNE, 24, 14, 30, 15), got %T", res)
	}

	y, m, d, h, minute, sec, n := getLocalDateTimeYMDHMSN(dt)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 0 {
		t.Errorf("Expected 2026-06-24T14:30:15, got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}

	// toString should produce 2026-06-24T14:30:15
	var gotStr string
	strRes := localDateTimeToString([]any{dt})
	if strObj, okStr := strRes.(*object.Object); okStr {
		if gotStr = object.GoStringFromStringObject(strObj); gotStr != "2026-06-24T14:30:15" {
			t.Errorf("Expected '2026-06-24T14:30:15', got '%s'", gotStr)
		}
		t.Logf("LocalDateTimeToString produced '%s' as expected", gotStr)
	} else {
		t.Errorf("Expected string object from toString, got %T", strRes)
	}
}

func TestLocalDateTimeOfAllOverloads(t *testing.T) {
	initTestEnvironment()

	sJune, _ := statics.QueryStatic("java/time/Month", "JUNE")
	monthJune := sJune.Value.(*object.Object)

	// of(year, monthInt, day, hour, min)
	r1 := localDateTimeOfHM([]any{int64(2026), int64(6), int64(24), int64(14), int64(30)}).(*object.Object)
	y, m, d, h, minute, sec, n := getLocalDateTimeYMDHMSN(r1)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 0 || n != 0 {
		t.Errorf("of(IIIII) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// of(year, monthInt, day, hour, min, sec)
	r2 := localDateTimeOfHMS([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15)}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(r2)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 0 {
		t.Errorf("of(IIIIII) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// of(year, monthInt, day, hour, min, sec, nano)
	r3 := localDateTimeOfHMSN([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500)}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(r3)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 500 {
		t.Errorf("of(IIIIIII) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// of(year, Month, day, hour, min)
	r4 := localDateTimeOfYMonthHM([]any{int64(2026), monthJune, int64(24), int64(14), int64(30)}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(r4)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 0 || n != 0 {
		t.Errorf("of(I Month III) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// of(year, Month, day, hour, min, sec) - via MethodSignatures table to verify exact descriptor binding
	gmethYMonthHMS, okSig := ghelpers.MethodSignatures["java/time/LocalDateTime.of(ILjava/time/Month;IIII)Ljava/time/LocalDateTime;"]
	if !okSig {
		t.Fatal("Method java/time/LocalDateTime.of(ILjava/time/Month;IIII)Ljava/time/LocalDateTime; not found")
	}
	if gmethYMonthHMS.ParamSlots != 6 {
		t.Errorf("Expected ParamSlots 6 for of(ILjava/time/Month;IIII), got %d", gmethYMonthHMS.ParamSlots)
	}
	r5a := gmethYMonthHMS.GFunction([]any{int64(2026), monthJune, int64(24), int64(14), int64(30), int64(15)}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(r5a)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 0 {
		t.Errorf("of(I Month IIII) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// of(year, Month, day, hour, min, sec, nano)
	gmethYMonthHMSN, okSigN := ghelpers.MethodSignatures["java/time/LocalDateTime.of(ILjava/time/Month;IIIII)Ljava/time/LocalDateTime;"]
	if !okSigN {
		t.Fatal("Method java/time/LocalDateTime.of(ILjava/time/Month;IIIII)Ljava/time/LocalDateTime; not found")
	}
	if gmethYMonthHMSN.ParamSlots != 7 {
		t.Errorf("Expected ParamSlots 7 for of(ILjava/time/Month;IIIII), got %d", gmethYMonthHMSN.ParamSlots)
	}
	r5 := gmethYMonthHMSN.GFunction([]any{int64(2026), monthJune, int64(24), int64(14), int64(30), int64(15), int64(999)}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(r5)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 999 {
		t.Errorf("of(I Month IIIII) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// of(LocalDate, LocalTime)
	dObj := CreateLocalDate(2026, 6, 24)
	tObj := createLocalTime(14, 30, 15, 123)
	r6 := localDateTimeOfDateTime([]any{dObj, tObj}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(r6)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 123 {
		t.Errorf("of(LocalDate, LocalTime) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// Error cases for of()
	if _, isErr := localDateTimeOfHM([]any{int64(2026), int64(13), int64(1), int64(0), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for month = 13")
	}
	if _, isErr := localDateTimeOfHM([]any{int64(2026), int64(2), int64(30), int64(0), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for Feb 30")
	}
	if _, isErr := localDateTimeOfHMS([]any{int64(2026), int64(6), int64(24), int64(25), int64(0), int64(0)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for hour = 25")
	}
	if _, isErr := localDateTimeOfYMonthHM([]any{int64(2026), nil, int64(24), int64(14), int64(30)}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for null Month")
	}
	if _, isErr := localDateTimeOfDateTime([]any{nil, tObj}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for null LocalDate")
	}
	if _, isErr := localDateTimeOfDateTime([]any{dObj, nil}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error for null LocalTime")
	}
}

func TestLocalDateTimeOfEpochSecondAndInstant(t *testing.T) {
	initTestEnvironment()

	zoneOffsetObj := object.MakeEmptyObjectWithClassName(&classNameZoneOffset)
	zoneOffsetObj.FieldTable["totalSeconds"] = object.Field{Ftype: types.Int, Fvalue: int64(3600)}

	// 1700000000 + 3600 sec offset
	res := localDateTimeOfEpochSecond([]any{int64(1700000000), int64(500), zoneOffsetObj})
	dt, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from ofEpochSecond, got %T", res)
	}
	y, m, d, h, minute, sec, n := getLocalDateTimeYMDHMSN(dt)
	if y != 2023 || m != 11 || d != 14 || h != 23 || minute != 13 || sec != 20 || n != 500 {
		t.Errorf("ofEpochSecond mismatch: got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}

	// ofInstant
	instantObj := object.MakeEmptyObjectWithClassName(&classNameInstant)
	instantObj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: int64(1700000000)}
	instantObj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: int64(500)}

	resInst := localDateTimeOfInstant([]any{instantObj, zoneOffsetObj})
	dtInst, ok := resInst.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from ofInstant, got %T", resInst)
	}
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(dtInst)
	if y != 2023 || m != 11 || d != 14 || h != 23 || minute != 13 || sec != 20 || n != 500 {
		t.Errorf("ofInstant mismatch: got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}
}

func TestLocalDateTimeGetters(t *testing.T) {
	initTestEnvironment()

	dObj := CreateLocalDate(2026, 6, 24)
	tObj := createLocalTime(14, 30, 15, 123456789)
	dt := CreateLocalDateTime(dObj, tObj)

	if got := localDateTimeGetYear([]any{dt}).(int64); got != 2026 {
		t.Errorf("getYear: expected 2026, got %d", got)
	}
	if got := localDateTimeGetMonthValue([]any{dt}).(int64); got != 6 {
		t.Errorf("getMonthValue: expected 6, got %d", got)
	}
	if got := localDateTimeGetDayOfMonth([]any{dt}).(int64); got != 24 {
		t.Errorf("getDayOfMonth: expected 24, got %d", got)
	}
	if got := localDateTimeGetHour([]any{dt}).(int64); got != 14 {
		t.Errorf("getHour: expected 14, got %d", got)
	}
	if got := localDateTimeGetMinute([]any{dt}).(int64); got != 30 {
		t.Errorf("getMinute: expected 30, got %d", got)
	}
	if got := localDateTimeGetSecond([]any{dt}).(int64); got != 15 {
		t.Errorf("getSecond: expected 15, got %d", got)
	}
	if got := localDateTimeGetNano([]any{dt}).(int64); got != 123456789 {
		t.Errorf("getNano: expected 123456789, got %d", got)
	}

	monthObj := localDateTimeGetMonth([]any{dt}).(*object.Object)
	ord, _ := getMonthOrdinal(monthObj)
	if ord != 5 { // JUNE ordinal is 5 (1-based value is 6)
		t.Errorf("getMonth: expected JUNE (ord 5), got ordinal %d", ord)
	}

	dowObj := localDateTimeGetDayOfWeek([]any{dt}).(*object.Object)
	dowOrd, _ := getDayOfWeekOrdinal(dowObj)
	// 2026-06-24 is Wednesday (WEDNESDAY ordinal is 2)
	if dowOrd != 2 {
		t.Errorf("getDayOfWeek: expected WEDNESDAY (ord 2), got ordinal %d", dowOrd)
	}

	doy := localDateTimeGetDayOfYear([]any{dt}).(int64)
	// Jan (31) + Feb (28) + Mar (31) + Apr (30) + May (31) + 24 = 175
	if doy != 175 {
		t.Errorf("getDayOfYear: expected 175, got %d", doy)
	}

	// get / getLong with ChronoField
	cfYear := chronoFieldValueOf([]any{object.StringObjectFromGoString("YEAR")}).(*object.Object)
	if got := localDateTimeGet([]any{dt, cfYear}).(int64); got != 2026 {
		t.Errorf("get(YEAR): expected 2026, got %d", got)
	}
	if got := localDateTimeGetLong([]any{dt, cfYear}).(int64); got != 2026 {
		t.Errorf("getLong(YEAR): expected 2026, got %d", got)
	}

	cfHour := chronoFieldValueOf([]any{object.StringObjectFromGoString("HOUR_OF_DAY")}).(*object.Object)
	if got := localDateTimeGet([]any{dt, cfHour}).(int64); got != 14 {
		t.Errorf("get(HOUR_OF_DAY): expected 14, got %d", got)
	}
}

func TestLocalDateTimeWithModifiers(t *testing.T) {
	initTestEnvironment()

	dObj := CreateLocalDate(2026, 6, 24)
	tObj := createLocalTime(14, 30, 15, 100)
	dt := CreateLocalDateTime(dObj, tObj)

	// withYear
	w1 := localDateTimeWithYear([]any{dt, int64(2030)}).(*object.Object)
	y, m, d, h, minute, sec, n := getLocalDateTimeYMDHMSN(w1)
	if y != 2030 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 100 {
		t.Errorf("withYear: got %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}

	// withMonth
	w2 := localDateTimeWithMonth([]any{dt, int64(12)}).(*object.Object)
	y, m, d, _, _, _, _ = getLocalDateTimeYMDHMSN(w2)
	if y != 2026 || m != 12 || d != 24 {
		t.Errorf("withMonth: got %d-%d-%d", y, m, d)
	}

	// withDayOfMonth
	w3 := localDateTimeWithDayOfMonth([]any{dt, int64(10)}).(*object.Object)
	y, m, d, _, _, _, _ = getLocalDateTimeYMDHMSN(w3)
	if y != 2026 || m != 6 || d != 10 {
		t.Errorf("withDayOfMonth: got %d-%d-%d", y, m, d)
	}

	// withDayOfYear
	w4 := localDateTimeWithDayOfYear([]any{dt, int64(1)}).(*object.Object)
	y, m, d, _, _, _, _ = getLocalDateTimeYMDHMSN(w4)
	if y != 2026 || m != 1 || d != 1 {
		t.Errorf("withDayOfYear: got %d-%d-%d", y, m, d)
	}

	// withHour
	w5 := localDateTimeWithHour([]any{dt, int64(23)}).(*object.Object)
	_, _, _, h, minute, sec, n = getLocalDateTimeYMDHMSN(w5)
	if h != 23 || minute != 30 || sec != 15 || n != 100 {
		t.Errorf("withHour: got %d:%d:%d.%d", h, minute, sec, n)
	}

	// withMinute
	w6 := localDateTimeWithMinute([]any{dt, int64(59)}).(*object.Object)
	_, _, _, h, minute, sec, n = getLocalDateTimeYMDHMSN(w6)
	if h != 14 || minute != 59 || sec != 15 || n != 100 {
		t.Errorf("withMinute: got %d:%d:%d.%d", h, minute, sec, n)
	}

	// withSecond
	w7 := localDateTimeWithSecond([]any{dt, int64(45)}).(*object.Object)
	_, _, _, h, minute, sec, n = getLocalDateTimeYMDHMSN(w7)
	if h != 14 || minute != 30 || sec != 45 || n != 100 {
		t.Errorf("withSecond: got %d:%d:%d.%d", h, minute, sec, n)
	}

	// withNano
	w8 := localDateTimeWithNano([]any{dt, int64(999)}).(*object.Object)
	_, _, _, h, minute, sec, n = getLocalDateTimeYMDHMSN(w8)
	if h != 14 || minute != 30 || sec != 15 || n != 999 {
		t.Errorf("withNano: got %d:%d:%d.%d", h, minute, sec, n)
	}

	// with(TemporalField, long)
	cfMinute := chronoFieldValueOf([]any{object.StringObjectFromGoString("MINUTE_OF_HOUR")}).(*object.Object)
	w9 := localDateTimeWithField([]any{dt, cfMinute, int64(12)}).(*object.Object)
	_, _, _, _, minute, _, _ = getLocalDateTimeYMDHMSN(w9)
	if minute != 12 {
		t.Errorf("withField(MINUTE_OF_HOUR): got minute %d", minute)
	}

	// with(TemporalAdjuster) with another LocalDate
	newDate := CreateLocalDate(2028, 1, 15)
	w10 := localDateTimeWithAdjuster([]any{dt, newDate}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(w10)
	if y != 2028 || m != 1 || d != 15 || h != 14 || minute != 30 || sec != 15 || n != 100 {
		t.Errorf("withAdjuster(LocalDate): got %d-%d-%d %d:%d:%d.%d", y, m, d, h, minute, sec, n)
	}
}

func TestLocalDateTimePlusMinusArithmetic(t *testing.T) {
	initTestEnvironment()

	dObj := CreateLocalDate(2026, 6, 24)
	tObj := createLocalTime(14, 30, 15, 0)
	dt := CreateLocalDateTime(dObj, tObj)

	// plusYears
	pY := localDateTimePlusYears([]any{dt, int64(2)}).(*object.Object)
	y, m, d, _, _, _, _ := getLocalDateTimeYMDHMSN(pY)
	if y != 2028 || m != 6 || d != 24 {
		t.Errorf("plusYears: expected 2028-06-24, got %d-%d-%d", y, m, d)
	}

	// plusMonths
	pM := localDateTimePlusMonths([]any{dt, int64(7)}).(*object.Object)
	y, m, d, _, _, _, _ = getLocalDateTimeYMDHMSN(pM)
	if y != 2027 || m != 1 || d != 24 {
		t.Errorf("plusMonths: expected 2027-01-24, got %d-%d-%d", y, m, d)
	}

	// plusWeeks
	pW := localDateTimePlusWeeks([]any{dt, int64(1)}).(*object.Object)
	y, m, d, _, _, _, _ = getLocalDateTimeYMDHMSN(pW)
	if y != 2026 || m != 7 || d != 1 {
		t.Errorf("plusWeeks: expected 2026-07-01, got %d-%d-%d", y, m, d)
	}

	// plusDays
	pD := localDateTimePlusDaysMethod([]any{dt, int64(10)}).(*object.Object)
	y, m, d, _, _, _, _ = getLocalDateTimeYMDHMSN(pD)
	if y != 2026 || m != 7 || d != 4 {
		t.Errorf("plusDays: expected 2026-07-04, got %d-%d-%d", y, m, d)
	}

	// plusHours crossing midnight
	pH := localDateTimePlusHours([]any{dt, int64(12)}).(*object.Object)
	y, m, d, h, minute, sec, _ := getLocalDateTimeYMDHMSN(pH)
	if y != 2026 || m != 6 || d != 25 || h != 2 || minute != 30 || sec != 15 {
		t.Errorf("plusHours: expected 2026-06-25T02:30:15, got %04d-%02d-%02dT%02d:%02d:%02d", y, m, d, h, minute, sec)
	}

	// plusMinutes
	pMin := localDateTimePlusMinutes([]any{dt, int64(45)}).(*object.Object)
	_, _, _, h, minute, _, _ = getLocalDateTimeYMDHMSN(pMin)
	if h != 15 || minute != 15 {
		t.Errorf("plusMinutes: expected 15:15, got %02d:%02d", h, minute)
	}

	// plusSeconds
	pSec := localDateTimePlusSeconds([]any{dt, int64(50)}).(*object.Object)
	_, _, _, _, minute, sec, _ = getLocalDateTimeYMDHMSN(pSec)
	if minute != 31 || sec != 5 {
		t.Errorf("plusSeconds: expected 31:05, got %02d:%02d", minute, sec)
	}

	// plusNanos
	pN := localDateTimePlusNanos([]any{dt, int64(500000000)}).(*object.Object)
	_, _, _, _, _, sec, n := getLocalDateTimeYMDHMSN(pN)
	if sec != 15 || n != 500000000 {
		t.Errorf("plusNanos: expected sec 15, nano 500000000, got sec %d, nano %d", sec, n)
	}

	// minusYears, minusMonths, minusWeeks, minusDays, minusHours, minusMinutes, minusSeconds, minusNanos
	mY := localDateTimeMinusYears([]any{dt, int64(1)}).(*object.Object)
	y, _, _, _, _, _, _ = getLocalDateTimeYMDHMSN(mY)
	if y != 2025 {
		t.Errorf("minusYears: expected 2025, got %d", y)
	}

	mM := localDateTimeMinusMonths([]any{dt, int64(1)}).(*object.Object)
	_, m, _, _, _, _, _ = getLocalDateTimeYMDHMSN(mM)
	if m != 5 {
		t.Errorf("minusMonths: expected 5, got %d", m)
	}

	mW := localDateTimeMinusWeeks([]any{dt, int64(1)}).(*object.Object)
	_, _, d, _, _, _, _ = getLocalDateTimeYMDHMSN(mW)
	if d != 17 {
		t.Errorf("minusWeeks: expected 17, got %d", d)
	}

	mD := localDateTimeMinusDays([]any{dt, int64(4)}).(*object.Object)
	_, _, d, _, _, _, _ = getLocalDateTimeYMDHMSN(mD)
	if d != 20 {
		t.Errorf("minusDays: expected 20, got %d", d)
	}

	mH := localDateTimeMinusHours([]any{dt, int64(15)}).(*object.Object)
	y, m, d, h, minute, sec, _ = getLocalDateTimeYMDHMSN(mH)
	if y != 2026 || m != 6 || d != 23 || h != 23 || minute != 30 || sec != 15 {
		t.Errorf("minusHours: expected 2026-06-23T23:30:15, got %04d-%02d-%02dT%02d:%02d:%02d", y, m, d, h, minute, sec)
	}

	mMin := localDateTimeMinusMinutes([]any{dt, int64(35)}).(*object.Object)
	_, _, _, h, minute, _, _ = getLocalDateTimeYMDHMSN(mMin)
	if h != 13 || minute != 55 {
		t.Errorf("minusMinutes: expected 13:55, got %02d:%02d", h, minute)
	}

	mSec := localDateTimeMinusSeconds([]any{dt, int64(20)}).(*object.Object)
	_, _, _, _, minute, sec, _ = getLocalDateTimeYMDHMSN(mSec)
	if minute != 29 || sec != 55 {
		t.Errorf("minusSeconds: expected 29:55, got %02d:%02d", minute, sec)
	}

	mN := localDateTimeMinusNanos([]any{dt, int64(1)}).(*object.Object)
	_, _, _, _, minute, sec, n = getLocalDateTimeYMDHMSN(mN)
	if minute != 30 || sec != 14 || n != 999999999 {
		t.Errorf("minusNanos: expected sec 14, nano 999999999, got sec %d, nano %d", sec, n)
	}

	// plusUnit and minusUnit
	daysUnit := chronoUnitValueOf([]any{object.StringObjectFromGoString("DAYS")}).(*object.Object)
	pU := localDateTimePlusUnit([]any{dt, int64(5), daysUnit}).(*object.Object)
	_, _, d, _, _, _, _ = getLocalDateTimeYMDHMSN(pU)
	if d != 29 {
		t.Errorf("plusUnit(5 DAYS): expected day 29, got %d", d)
	}

	mU := localDateTimeMinusUnit([]any{dt, int64(5), daysUnit}).(*object.Object)
	_, _, d, _, _, _, _ = getLocalDateTimeYMDHMSN(mU)
	if d != 19 {
		t.Errorf("minusUnit(5 DAYS): expected day 19, got %d", d)
	}
}

func TestLocalDateTimeCompareEqualsHashCode(t *testing.T) {
	initTestEnvironment()

	dt1 := CreateLocalDateTime(CreateLocalDate(2026, 6, 24), createLocalTime(14, 30, 15, 0))
	dt2 := CreateLocalDateTime(CreateLocalDate(2026, 6, 24), createLocalTime(14, 30, 15, 0))
	dt3 := CreateLocalDateTime(CreateLocalDate(2026, 6, 24), createLocalTime(15, 0, 0, 0))
	dt0 := CreateLocalDateTime(CreateLocalDate(2026, 6, 20), createLocalTime(14, 30, 15, 0))

	// compareTo
	if cmp := localDateTimeCompareTo([]any{dt1, dt2}).(int64); cmp != 0 {
		t.Errorf("compareTo equal: expected 0, got %d", cmp)
	}
	if cmp := localDateTimeCompareTo([]any{dt1, dt3}).(int64); cmp >= 0 {
		t.Errorf("compareTo earlier vs later: expected < 0, got %d", cmp)
	}
	if cmp := localDateTimeCompareTo([]any{dt1, dt0}).(int64); cmp <= 0 {
		t.Errorf("compareTo later vs earlier: expected > 0, got %d", cmp)
	}

	// isBefore, isAfter, isEqual
	if !isTrue(localDateTimeIsBefore([]any{dt1, dt3})) {
		t.Error("isBefore: expected dt1 < dt3")
	}
	if isTrue(localDateTimeIsBefore([]any{dt3, dt1})) {
		t.Error("isBefore: expected dt3 not < dt1")
	}

	if !isTrue(localDateTimeIsAfter([]any{dt3, dt1})) {
		t.Error("isAfter: expected dt3 > dt1")
	}
	if isTrue(localDateTimeIsAfter([]any{dt1, dt3})) {
		t.Error("isAfter: expected dt1 not > dt3")
	}

	if !isTrue(localDateTimeIsEqual([]any{dt1, dt2})) {
		t.Error("isEqual: expected dt1 == dt2")
	}
	if isTrue(localDateTimeIsEqual([]any{dt1, dt3})) {
		t.Error("isEqual: expected dt1 != dt3")
	}

	// equals
	if !isTrue(localDateTimeEquals([]any{dt1, dt2})) {
		t.Error("equals: expected true for same values")
	}
	if isTrue(localDateTimeEquals([]any{dt1, dt3})) {
		t.Error("equals: expected false for different time")
	}
	if isTrue(localDateTimeEquals([]any{dt1, nil})) {
		t.Error("equals: expected false for nil")
	}

	// hashCode
	h1 := localDateTimeHashCode([]any{dt1}).(int64)
	h2 := localDateTimeHashCode([]any{dt2}).(int64)
	if h1 != h2 {
		t.Errorf("hashCode: expected %d == %d", h1, h2)
	}
}

func TestLocalDateTimeConversionsAndAdjustments(t *testing.T) {
	initTestEnvironment()

	dOrig := CreateLocalDate(2026, 6, 24)
	tOrig := createLocalTime(14, 30, 15, 500)
	dt := CreateLocalDateTime(dOrig, tOrig)

	// toLocalDate
	dRes := localDateTimeToLocalDate([]any{dt}).(*object.Object)
	y, m, d, _ := getLocalDateYMD(dRes)
	if y != 2026 || m != 6 || d != 24 {
		t.Errorf("toLocalDate: got %d-%d-%d", y, m, d)
	}

	// toLocalTime
	tRes := localDateTimeToLocalTime([]any{dt}).(*object.Object)
	h, minute, sec, n := getLocalTimeHMSN(tRes)
	if h != 14 || minute != 30 || sec != 15 || n != 500 {
		t.Errorf("toLocalTime: got %d:%d:%d.%d", h, minute, sec, n)
	}

	// toEpochSecond
	zoneOffsetObj := object.MakeEmptyObjectWithClassName(&classNameZoneOffset)
	zoneOffsetObj.FieldTable["totalSeconds"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}

	epochSec := localDateTimeToEpochSecond([]any{dt, zoneOffsetObj}).(int64)
	// 2026-06-24T14:30:15 UTC
	// days between 1970-01-01 and 2026-06-24:
	expectedSec := ghelpers.GetChronoLocalDateEpochDay(dOrig)*86400 + (14*3600 + 30*60 + 15)
	if epochSec != expectedSec {
		t.Errorf("toEpochSecond: expected %d, got %d", expectedSec, epochSec)
	}

	// toInstant
	instRes := localDateTimeToInstant([]any{dt, zoneOffsetObj}).(*object.Object)
	sFld, _ := instRes.FieldTable["seconds"]
	nFld, _ := instRes.FieldTable["nanos"]
	sVal, _ := getInt64Val(sFld.Fvalue)
	nVal, _ := getInt64Val(nFld.Fvalue)
	if sVal != expectedSec || nVal != 500 {
		t.Errorf("toInstant: expected %d sec / 500 nano, got %d sec / %d nano", expectedSec, sVal, nVal)
	}

	// adjustInto
	targetDT := CreateLocalDateTime(CreateLocalDate(2000, 1, 1), createLocalTime(0, 0, 0, 0))
	adjRes := localDateTimeAdjustInto([]any{dt, targetDT}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(adjRes)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 500 {
		t.Errorf("adjustInto: got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}

	// truncatedTo
	hoursUnit := chronoUnitValueOf([]any{object.StringObjectFromGoString("HOURS")}).(*object.Object)
	truncRes := localDateTimeTruncatedTo([]any{dt, hoursUnit}).(*object.Object)
	_, _, _, h, minute, sec, n = getLocalDateTimeYMDHMSN(truncRes)
	if h != 14 || minute != 0 || sec != 0 || n != 0 {
		t.Errorf("truncatedTo(HOURS): expected 14:00:00.0, got %02d:%02d:%02d.%d", h, minute, sec, n)
	}

	// until
	dtLater := CreateLocalDateTime(CreateLocalDate(2026, 6, 26), createLocalTime(14, 30, 15, 500))
	daysUnit := chronoUnitValueOf([]any{object.StringObjectFromGoString("DAYS")}).(*object.Object)
	daysUntil := localDateTimeUntil([]any{dt, dtLater, daysUnit}).(int64)
	if daysUntil != 2 {
		t.Errorf("until(DAYS): expected 2, got %d", daysUntil)
	}
}

func TestLocalDateTimeParseAndFormat(t *testing.T) {
	initTestEnvironment()

	// parse("2026-06-24T14:30:15")
	strInput := object.StringObjectFromGoString("2026-06-24T14:30:15")
	parsedRes := localDateTimeParse([]any{strInput}).(*object.Object)
	y, m, d, h, minute, sec, n := getLocalDateTimeYMDHMSN(parsedRes)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 0 {
		t.Errorf("parse: expected 2026-06-24T14:30:15, got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}

	// parse invalid string
	badStr := object.StringObjectFromGoString("invalid-date-time")
	if _, isErr := localDateTimeParse([]any{badStr}).(*ghelpers.GErrBlk); !isErr {
		t.Error("Expected error when parsing invalid string")
	}

	// from(TemporalAccessor)
	fromRes := localDateTimeFrom([]any{parsedRes}).(*object.Object)
	y, m, d, h, minute, sec, n = getLocalDateTimeYMDHMSN(fromRes)
	if y != 2026 || m != 6 || d != 24 || h != 14 || minute != 30 || sec != 15 || n != 0 {
		t.Errorf("from: expected 2026-06-24T14:30:15, got %04d-%02d-%02dT%02d:%02d:%02d.%d", y, m, d, h, minute, sec, n)
	}
}

func TestLocalDateTimeSupportedAndRange(t *testing.T) {
	initTestEnvironment()

	dt := CreateLocalDateTime(CreateLocalDate(2026, 6, 24), createLocalTime(14, 30, 15, 0))

	// isSupported(TemporalField)
	cfYear := chronoFieldValueOf([]any{object.StringObjectFromGoString("YEAR")}).(*object.Object)
	if !isTrue(localDateTimeIsSupportedField([]any{dt, cfYear})) {
		t.Error("isSupported(YEAR): expected true")
	}

	// isSupported(TemporalUnit)
	cuDays := chronoUnitValueOf([]any{object.StringObjectFromGoString("DAYS")}).(*object.Object)
	if !isTrue(localDateTimeIsSupportedUnit([]any{dt, cuDays})) {
		t.Error("isSupported(DAYS): expected true")
	}

	// range(TemporalField)
	cfDOM := chronoFieldValueOf([]any{object.StringObjectFromGoString("DAY_OF_MONTH")}).(*object.Object)
	rangeRes := localDateTimeRange([]any{dt, cfDOM}).(*object.Object)
	if rangeRes == nil || object.IsNull(rangeRes) {
		t.Error("range(DAY_OF_MONTH): expected valid ValueRange object")
	}
}

func TestLocalTemporalFormatUsesFormatter(t *testing.T) {
	initTestEnvironment()

	tests := []struct {
		name      string
		format    func([]any) any
		temporal  *object.Object
		formatter *object.Object
		expected  string
	}{
		{
			name:      "LocalDate",
			format:    localDateFormat,
			temporal:  CreateLocalDate(2026, 6, 24),
			formatter: createDateTimeFormatterObject("yyyy/MM/dd", nil, nil, nil),
			expected:  "2026/06/24",
		},
		{
			name:      "LocalTime",
			format:    localTimeFormat,
			temporal:  createLocalTime(14, 30, 15, 0),
			formatter: createDateTimeFormatterObject("HH:mm", nil, nil, nil),
			expected:  "14:30",
		},
		{
			name:      "LocalDateTime",
			format:    localDateTimeFormat,
			temporal:  CreateLocalDateTime(CreateLocalDate(2026, 6, 24), createLocalTime(14, 30, 15, 0)),
			formatter: createDateTimeFormatterObject("yyyy/MM/dd HH:mm", nil, nil, nil),
			expected:  "2026/06/24 14:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := tt.format([]any{tt.temporal, tt.formatter}).(*object.Object)
			if !ok {
				t.Fatalf("format returned unexpected value: %v", result)
			}
			if got := object.GoStringFromStringObject(result); got != tt.expected {
				t.Errorf("format returned %q, want %q", got, tt.expected)
			}
		})
	}
}
