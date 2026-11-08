/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func initZonedDateTimeTestEnvironment() {
	initTestEnvironment()
	ensureZonedDateTimeInited()
	ensureZoneOffsetInited()
	ensureZoneIdInited()
	ensureInstantInited()
	ensureLocalDateTimeInited()
	ensureLocalDateInited()
	ensureLocalTimeInited()
	ensureChronoFieldInited()
	ensureChronoUnitInited()
	ensureTemporalQueriesInited()
	Load_Time_ZonedDateTime()
	Load_Time_ZoneOffset()
	Load_Time_ZoneId()
	Load_Time_Instant()
	Load_Time_LocalDateTime()
	Load_Time_LocalDate()
	Load_Time_LocalTime()
}

func TestZonedDateTimeClinitAndInit(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	res := zonedDateTimeClinit(nil)
	if res != nil {
		t.Errorf("Expected zonedDateTimeClinit to return nil, got %v", res)
	}

	resInit := zonedDateTimeInit(nil)
	if resInit != nil {
		t.Errorf("Expected zonedDateTimeInit to return nil, got %v", resInit)
	}
}

func TestZonedDateTimeFactories(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	zoneUTC := makeZoneIdFromID("UTC")
	zoneParis := makeZoneIdFromID("Europe/Paris")

	// now()
	nowObj := zonedDateTimeNow(nil).(*object.Object)
	if nowObj == nil {
		t.Fatal("zonedDateTimeNow returned nil")
	}

	// now(zone)
	nowParis := zonedDateTimeNowZone([]any{zoneParis}).(*object.Object)
	if nowParis == nil {
		t.Fatal("zonedDateTimeNowZone returned nil")
	}

	// now(clock)
	nowClock := zonedDateTimeNowClock([]any{zoneUTC}).(*object.Object)
	if nowClock == nil {
		t.Fatal("zonedDateTimeNowClock returned nil")
	}

	// of(LocalDate, LocalTime, ZoneId)
	dateObj := CreateLocalDate(2026, 6, 24)
	timeObj := createLocalTime(14, 30, 15, 500)
	zdtDate := zonedDateTimeOfDate([]any{dateObj, timeObj, zoneUTC}).(*object.Object)
	y, m, d, h, min, s, nano := getZonedDateTimeYMDHMSN(zdtDate)
	if y != 2026 || m != 6 || d != 24 || h != 14 || min != 30 || s != 15 || nano != 500 {
		t.Errorf("zonedDateTimeOfDate mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, min, s, nano)
	}

	// of(LocalDateTime, ZoneId)
	ldt := createLocalDateTime(dateObj, timeObj)
	zdtLDT := zonedDateTimeOfLDT([]any{ldt, zoneUTC}).(*object.Object)
	if getZonedDateTimeEpochSecond(zdtLDT) != getZonedDateTimeEpochSecond(zdtDate) {
		t.Errorf("zonedDateTimeOfLDT mismatch")
	}

	// of(int, int, int, int, int, int, int, ZoneId)
	zdtInts := zonedDateTimeOfInts([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500), zoneUTC}).(*object.Object)
	if getZonedDateTimeEpochSecond(zdtInts) != getZonedDateTimeEpochSecond(zdtDate) {
		t.Errorf("zonedDateTimeOfInts mismatch")
	}

	// ofLocal
	offsetUTC := createZoneOffset(0)
	zdtLocal := zonedDateTimeOfLocalMethod([]any{ldt, zoneUTC, offsetUTC}).(*object.Object)
	if zdtLocal == nil {
		t.Fatal("zonedDateTimeOfLocal returned nil")
	}

	// ofInstant(Instant, ZoneId)
	inst := createInstant(1782311415, 500)
	zdtInst := zonedDateTimeOfInstantMethod([]any{inst, zoneUTC}).(*object.Object)
	if getZonedDateTimeEpochSecond(zdtInst) != 1782311415 {
		t.Errorf("Expected epoch second 1782311415, got %d", getZonedDateTimeEpochSecond(zdtInst))
	}

	// ofInstant(LocalDateTime, ZoneOffset, ZoneId)
	zdtInst2 := zonedDateTimeOfInstantLdtOffsetZoneMethod([]any{ldt, offsetUTC, zoneUTC}).(*object.Object)
	if zdtInst2 == nil {
		t.Fatal("zonedDateTimeOfInstantLdtOffsetZoneMethod returned nil")
	}

	// ofStrict
	zdtStrict := zonedDateTimeOfStrictMethod([]any{ldt, offsetUTC, zoneUTC}).(*object.Object)
	if zdtStrict == nil {
		t.Fatal("zonedDateTimeOfStrictMethod returned nil")
	}

	// from(temporal)
	zdtFrom := zonedDateTimeFrom([]any{zdtStrict}).(*object.Object)
	if zdtFrom != zdtStrict {
		t.Errorf("Expected zonedDateTimeFrom to return same object for ZonedDateTime")
	}

	// from(Instant with ZoneId)
	instWithZone := object.MakeEmptyObjectWithClassName(&classNameInstantVal)
	instWithZone.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: int64(1782311415)}
	instWithZone.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: int64(500)}
	instWithZone.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zoneUTC}
	zdtFromInst := zonedDateTimeFrom([]any{instWithZone}).(*object.Object)
	if zdtFromInst == nil || getZonedDateTimeEpochSecond(zdtFromInst) != 1782311415 {
		t.Errorf("Expected zonedDateTimeFrom instWithZone to succeed")
	}
}

func TestZonedDateTimeParseAndToString(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	// parse "2007-12-03T10:15:30+01:00[Europe/Paris]"
	strObj := object.StringObjectFromGoString("2007-12-03T10:15:30+01:00[Europe/Paris]")
	parsedRes := zonedDateTimeParse([]any{strObj})
	if err, isErr := parsedRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Failed to parse: %v", err)
	}
	parsed := parsedRes.(*object.Object)
	y, m, d, h, min, s, _ := getZonedDateTimeYMDHMSN(parsed)
	if y != 2007 || m != 12 || d != 3 || h != 10 || min != 15 || s != 30 {
		t.Errorf("Parsed mismatch: %d-%d-%d %d:%d:%d", y, m, d, h, min, s)
	}

	// toString
	strVal := zonedDateTimeToString([]any{parsed}).(*object.Object)
	goStr := object.GoStringFromStringObject(strVal)
	if goStr != "2007-12-03T10:15:30+01:00[Europe/Paris]" {
		t.Errorf("toString mismatch: expected '2007-12-03T10:15:30+01:00[Europe/Paris]', got '%s'", goStr)
	}

	// parse "2007-12-03T10:15:30Z"
	strUTC := object.StringObjectFromGoString("2007-12-03T10:15:30Z")
	parsedUTC := zonedDateTimeParse([]any{strUTC}).(*object.Object)
	strUTCVal := zonedDateTimeToString([]any{parsedUTC}).(*object.Object)
	if object.GoStringFromStringObject(strUTCVal) != "2007-12-03T10:15:30Z" {
		t.Errorf("toString UTC mismatch: got %s", object.GoStringFromStringObject(strUTCVal))
	}

	// parseFormatter
	parsedFmt := zonedDateTimeParseFormatter([]any{strUTC, strUTC}).(*object.Object)
	if parsedFmt == nil {
		t.Errorf("zonedDateTimeParseFormatter failed")
	}
}

func TestZonedDateTimeGettersAndQueries(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	zoneParis := makeZoneIdFromID("Europe/Paris")
	zdt := zonedDateTimeOfInts([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500), zoneParis}).(*object.Object)

	if zonedDateTimeGetYear([]any{zdt}).(int64) != 2026 {
		t.Errorf("getYear mismatch")
	}
	if zonedDateTimeGetMonthValue([]any{zdt}).(int64) != 6 {
		t.Errorf("getMonthValue mismatch")
	}
	if zonedDateTimeGetDayOfMonth([]any{zdt}).(int64) != 24 {
		t.Errorf("getDayOfMonth mismatch")
	}
	if zonedDateTimeGetHour([]any{zdt}).(int64) != 14 {
		t.Errorf("getHour mismatch")
	}
	if zonedDateTimeGetMinute([]any{zdt}).(int64) != 30 {
		t.Errorf("getMinute mismatch")
	}
	if zonedDateTimeGetSecond([]any{zdt}).(int64) != 15 {
		t.Errorf("getSecond mismatch")
	}
	if zonedDateTimeGetNano([]any{zdt}).(int64) != 500 {
		t.Errorf("getNano mismatch")
	}
	if zonedDateTimeGetMonth([]any{zdt}) == nil {
		t.Errorf("getMonth returned nil")
	}
	if zonedDateTimeGetDayOfWeek([]any{zdt}) == nil {
		t.Errorf("getDayOfWeek returned nil")
	}
	if zonedDateTimeGetDayOfYear([]any{zdt}).(int64) <= 0 {
		t.Errorf("getDayOfYear invalid")
	}

	// toLocalDateTime, toLocalDate, toLocalTime
	if zonedDateTimeToLocalDateTime([]any{zdt}) == nil {
		t.Errorf("toLocalDateTime returned nil")
	}
	if zonedDateTimeToLocalDate([]any{zdt}) == nil {
		t.Errorf("toLocalDate returned nil")
	}
	if zonedDateTimeToLocalTime([]any{zdt}) == nil {
		t.Errorf("toLocalTime returned nil")
	}
	if zonedDateTimeToOffsetDateTime([]any{zdt}) == nil {
		t.Errorf("toOffsetDateTime returned nil")
	}
	if zonedDateTimeToInstant([]any{zdt}) == nil {
		t.Errorf("toInstant returned nil")
	}
	if zonedDateTimeToEpochSecond([]any{zdt}).(int64) == 0 {
		t.Errorf("toEpochSecond invalid")
	}
	if zonedDateTimeGetChronology([]any{zdt}) == nil {
		t.Errorf("getChronology returned nil")
	}
	if zonedDateTimeGetZone([]any{zdt}) == nil {
		t.Errorf("getZone returned nil")
	}
	if zonedDateTimeGetOffset([]any{zdt}) == nil {
		t.Errorf("getOffset returned nil")
	}

	// isSupported
	cfInstSec := getChronoField("INSTANT_SECONDS")
	if zonedDateTimeIsSupportedField([]any{zdt, cfInstSec}) != types.JavaBoolTrue {
		t.Errorf("Expected INSTANT_SECONDS supported")
	}
	cuDays := ghelpers.CreateTemporalUnit("DAYS")
	if zonedDateTimeIsSupportedUnit([]any{zdt, cuDays}) != types.JavaBoolTrue {
		t.Errorf("Expected DAYS supported")
	}

	// range
	vr := zonedDateTimeRange([]any{zdt, cfInstSec}).(*object.Object)
	if vr == nil {
		t.Errorf("range returned nil")
	}

	// get & getLong
	cfHour := getChronoField("HOUR_OF_DAY")
	if zonedDateTimeGet([]any{zdt, cfHour}).(int64) != 14 {
		t.Errorf("get(HOUR_OF_DAY) mismatch")
	}
	if zonedDateTimeGetLong([]any{zdt, cfInstSec}).(int64) == 0 {
		t.Errorf("getLong(INSTANT_SECONDS) mismatch")
	}

	// query
	qZone := temporalQueriesZone(nil).(*object.Object)
	if zonedDateTimeQuery([]any{zdt, qZone}) == nil {
		t.Errorf("query(zone) returned nil")
	}
	qOffset := temporalQueriesOffset(nil).(*object.Object)
	if zonedDateTimeQuery([]any{zdt, qOffset}) == nil {
		t.Errorf("query(offset) returned nil")
	}
	qLocalDate := temporalQueriesLocalDate(nil).(*object.Object)
	if zonedDateTimeQuery([]any{zdt, qLocalDate}) == nil {
		t.Errorf("query(localDate) returned nil")
	}
	qLocalTime := temporalQueriesLocalTime(nil).(*object.Object)
	if zonedDateTimeQuery([]any{zdt, qLocalTime}) == nil {
		t.Errorf("query(localTime) returned nil")
	}
	qChrono := temporalQueriesChronology(nil).(*object.Object)
	if zonedDateTimeQuery([]any{zdt, qChrono}) == nil {
		t.Errorf("query(chronology) returned nil")
	}
	qPrecision := temporalQueriesPrecision(nil).(*object.Object)
	if zonedDateTimeQuery([]any{zdt, qPrecision}) == nil {
		t.Errorf("query(precision) returned nil")
	}
}

func TestZonedDateTimeWithModifications(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	zoneUTC := makeZoneIdFromID("UTC")
	zdt := zonedDateTimeOfInts([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500), zoneUTC}).(*object.Object)

	// withYear
	wY := zonedDateTimeWithYear([]any{zdt, int64(2030)}).(*object.Object)
	if zonedDateTimeGetYear([]any{wY}).(int64) != 2030 {
		t.Errorf("withYear mismatch: %v", zonedDateTimeGetYear([]any{wY}))
	}

	// withMonth
	wM := zonedDateTimeWithMonth([]any{zdt, int64(11)}).(*object.Object)
	if zonedDateTimeGetMonthValue([]any{wM}).(int64) != 11 {
		t.Errorf("withMonth mismatch")
	}

	// withDayOfMonth
	wD := zonedDateTimeWithDayOfMonth([]any{zdt, int64(5)}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{wD}).(int64) != 5 {
		t.Errorf("withDayOfMonth mismatch")
	}

	// withHour
	wH := zonedDateTimeWithHour([]any{zdt, int64(8)}).(*object.Object)
	if zonedDateTimeGetHour([]any{wH}).(int64) != 8 {
		t.Errorf("withHour mismatch")
	}

	// withMinute
	wMin := zonedDateTimeWithMinute([]any{zdt, int64(45)}).(*object.Object)
	if zonedDateTimeGetMinute([]any{wMin}).(int64) != 45 {
		t.Errorf("withMinute mismatch")
	}

	// withSecond
	wS := zonedDateTimeWithSecond([]any{zdt, int64(59)}).(*object.Object)
	if zonedDateTimeGetSecond([]any{wS}).(int64) != 59 {
		t.Errorf("withSecond mismatch")
	}

	// withNano
	wN := zonedDateTimeWithNano([]any{zdt, int64(12345)}).(*object.Object)
	if zonedDateTimeGetNano([]any{wN}).(int64) != 12345 {
		t.Errorf("withNano mismatch")
	}

	// withDayOfYear
	wDoy := zonedDateTimeWithDayOfYear([]any{zdt, int64(100)}).(*object.Object)
	if zonedDateTimeGetDayOfYear([]any{wDoy}).(int64) != 100 {
		t.Errorf("withDayOfYear mismatch")
	}

	// withAdjuster (LocalDate)
	newDate := CreateLocalDate(2028, 1, 15)
	wAdjDate := zonedDateTimeWithAdjuster([]any{zdt, newDate}).(*object.Object)
	if zonedDateTimeGetYear([]any{wAdjDate}).(int64) != 2028 {
		t.Errorf("withAdjuster(LocalDate) mismatch")
	}

	// withField (ChronoField)
	cfYear := getChronoField("YEAR")
	wField := zonedDateTimeWithField([]any{zdt, cfYear, int64(2040)}).(*object.Object)
	if zonedDateTimeGetYear([]any{wField}).(int64) != 2040 {
		t.Errorf("withField(YEAR) mismatch")
	}

	// truncatedTo
	cuHours := ghelpers.CreateTemporalUnit("HOURS")
	trunc := zonedDateTimeTruncatedTo([]any{zdt, cuHours}).(*object.Object)
	if zonedDateTimeGetMinute([]any{trunc}).(int64) != 0 || zonedDateTimeGetSecond([]any{trunc}).(int64) != 0 {
		t.Errorf("truncatedTo(HOURS) mismatch")
	}

	// withEarlierOffsetAtOverlap / withLaterOffsetAtOverlap / withFixedOffsetZone
	earlier := zonedDateTimeWithEarlierOffsetAtOverlap([]any{zdt}).(*object.Object)
	if earlier == nil {
		t.Errorf("withEarlierOffsetAtOverlap returned nil")
	}
	later := zonedDateTimeWithLaterOffsetAtOverlap([]any{zdt}).(*object.Object)
	if later == nil {
		t.Errorf("withLaterOffsetAtOverlap returned nil")
	}
	fixed := zonedDateTimeWithFixedOffsetZone([]any{zdt}).(*object.Object)
	if fixed == nil {
		t.Errorf("withFixedOffsetZone returned nil")
	}

	// withZoneSameLocal / withZoneSameInstant
	zoneTokyo := makeZoneIdFromID("Asia/Tokyo")
	sameLocal := zonedDateTimeWithZoneSameLocal([]any{zdt, zoneTokyo}).(*object.Object)
	if zonedDateTimeGetHour([]any{sameLocal}).(int64) != 14 {
		t.Errorf("withZoneSameLocal hour changed")
	}
	sameInstant := zonedDateTimeWithZoneSameInstant([]any{zdt, zoneTokyo}).(*object.Object)
	if getZonedDateTimeEpochSecond(sameInstant) != getZonedDateTimeEpochSecond(zdt) {
		t.Errorf("withZoneSameInstant epochSecond changed")
	}
}

func TestZonedDateTimeArithmetic(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	zoneUTC := makeZoneIdFromID("UTC")
	zdt := zonedDateTimeOfInts([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500), zoneUTC}).(*object.Object)

	// plusYears / minusYears
	pY := zonedDateTimePlusYears([]any{zdt, int64(2)}).(*object.Object)
	if zonedDateTimeGetYear([]any{pY}).(int64) != 2028 {
		t.Errorf("plusYears mismatch")
	}
	mY := zonedDateTimeMinusYears([]any{zdt, int64(2)}).(*object.Object)
	if zonedDateTimeGetYear([]any{mY}).(int64) != 2024 {
		t.Errorf("minusYears mismatch")
	}

	// plusMonths / minusMonths
	pM := zonedDateTimePlusMonths([]any{zdt, int64(3)}).(*object.Object)
	if zonedDateTimeGetMonthValue([]any{pM}).(int64) != 9 {
		t.Errorf("plusMonths mismatch")
	}
	mM := zonedDateTimeMinusMonths([]any{zdt, int64(3)}).(*object.Object)
	if zonedDateTimeGetMonthValue([]any{mM}).(int64) != 3 {
		t.Errorf("minusMonths mismatch")
	}

	// plusWeeks / minusWeeks
	pW := zonedDateTimePlusWeeks([]any{zdt, int64(1)}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{pW}).(int64) != 1 || zonedDateTimeGetMonthValue([]any{pW}).(int64) != 7 {
		t.Errorf("plusWeeks mismatch")
	}
	mW := zonedDateTimeMinusWeeks([]any{zdt, int64(1)}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{mW}).(int64) != 17 {
		t.Errorf("minusWeeks mismatch")
	}

	// plusDays / minusDays
	pD := zonedDateTimePlusDays([]any{zdt, int64(10)}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{pD}).(int64) != 4 || zonedDateTimeGetMonthValue([]any{pD}).(int64) != 7 {
		t.Errorf("plusDays mismatch")
	}
	mD := zonedDateTimeMinusDays([]any{zdt, int64(10)}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{mD}).(int64) != 14 {
		t.Errorf("minusDays mismatch")
	}

	// plusHours / minusHours
	pH := zonedDateTimePlusHours([]any{zdt, int64(5)}).(*object.Object)
	if zonedDateTimeGetHour([]any{pH}).(int64) != 19 {
		t.Errorf("plusHours mismatch")
	}
	mH := zonedDateTimeMinusHours([]any{zdt, int64(5)}).(*object.Object)
	if zonedDateTimeGetHour([]any{mH}).(int64) != 9 {
		t.Errorf("minusHours mismatch")
	}

	// plusMinutes / minusMinutes
	pMin := zonedDateTimePlusMinutes([]any{zdt, int64(20)}).(*object.Object)
	if zonedDateTimeGetMinute([]any{pMin}).(int64) != 50 {
		t.Errorf("plusMinutes mismatch")
	}
	mMin := zonedDateTimeMinusMinutes([]any{zdt, int64(20)}).(*object.Object)
	if zonedDateTimeGetMinute([]any{mMin}).(int64) != 10 {
		t.Errorf("minusMinutes mismatch")
	}

	// plusSeconds / minusSeconds
	pS := zonedDateTimePlusSeconds([]any{zdt, int64(30)}).(*object.Object)
	if zonedDateTimeGetSecond([]any{pS}).(int64) != 45 {
		t.Errorf("plusSeconds mismatch")
	}
	mS := zonedDateTimeMinusSeconds([]any{zdt, int64(10)}).(*object.Object)
	if zonedDateTimeGetSecond([]any{mS}).(int64) != 5 {
		t.Errorf("minusSeconds mismatch")
	}

	// plusNanos / minusNanos
	pN := zonedDateTimePlusNanos([]any{zdt, int64(500)}).(*object.Object)
	if zonedDateTimeGetNano([]any{pN}).(int64) != 1000 {
		t.Errorf("plusNanos mismatch")
	}
	mN := zonedDateTimeMinusNanos([]any{zdt, int64(200)}).(*object.Object)
	if zonedDateTimeGetNano([]any{mN}).(int64) != 300 {
		t.Errorf("minusNanos mismatch")
	}

	// plusUnit / minusUnit
	cuDays := ghelpers.CreateTemporalUnit("DAYS")
	pU := zonedDateTimePlusUnit([]any{zdt, int64(3), cuDays}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{pU}).(int64) != 27 {
		t.Errorf("plusUnit(DAYS) mismatch")
	}
	mU := zonedDateTimeMinusUnit([]any{zdt, int64(3), cuDays}).(*object.Object)
	if zonedDateTimeGetDayOfMonth([]any{mU}).(int64) != 21 {
		t.Errorf("minusUnit(DAYS) mismatch")
	}

	// plusAmount / minusAmount
	periodObj := createPeriod(0, 1, 5)
	pAmt := zonedDateTimePlusAmount([]any{zdt, periodObj}).(*object.Object)
	if zonedDateTimeGetMonthValue([]any{pAmt}).(int64) != 7 || zonedDateTimeGetDayOfMonth([]any{pAmt}).(int64) != 29 {
		t.Errorf("plusAmount(Period) mismatch")
	}
	mAmt := zonedDateTimeMinusAmount([]any{zdt, periodObj}).(*object.Object)
	if zonedDateTimeGetMonthValue([]any{mAmt}).(int64) != 5 || zonedDateTimeGetDayOfMonth([]any{mAmt}).(int64) != 19 {
		t.Errorf("minusAmount(Period) mismatch")
	}

	// until
	zdtLater := zonedDateTimePlusDays([]any{zdt, int64(5)}).(*object.Object)
	daysDiff := zonedDateTimeUntil([]any{zdt, zdtLater, cuDays}).(int64)
	if daysDiff != 5 {
		t.Errorf("until(DAYS) expected 5, got %d", daysDiff)
	}
}

func TestZonedDateTimeCompareEqualsAndHash(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	zoneUTC := makeZoneIdFromID("UTC")
	zdt1 := zonedDateTimeOfInts([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500), zoneUTC}).(*object.Object)
	zdt2 := zonedDateTimeOfInts([]any{int64(2026), int64(6), int64(24), int64(14), int64(30), int64(15), int64(500), zoneUTC}).(*object.Object)
	zdtLater := zonedDateTimePlusDays([]any{zdt1, int64(1)}).(*object.Object)

	// equals
	if zonedDateTimeEquals([]any{zdt1, zdt2}) != types.JavaBoolTrue {
		t.Errorf("Expected zdt1 equals zdt2")
	}
	if zonedDateTimeEquals([]any{zdt1, zdtLater}) != types.JavaBoolFalse {
		t.Errorf("Expected zdt1 not equals zdtLater")
	}

	// hashCode
	h1 := zonedDateTimeHashCode([]any{zdt1}).(int64)
	h2 := zonedDateTimeHashCode([]any{zdt2}).(int64)
	if h1 != h2 {
		t.Errorf("Expected equal hashcodes, got %d and %d", h1, h2)
	}

	// compareTo
	if zonedDateTimeCompareTo([]any{zdt1, zdt2}).(int64) != 0 {
		t.Errorf("compareTo expected 0")
	}
	if zonedDateTimeCompareTo([]any{zdt1, zdtLater}).(int64) >= 0 {
		t.Errorf("compareTo expected < 0")
	}
	if zonedDateTimeCompareTo([]any{zdtLater, zdt1}).(int64) <= 0 {
		t.Errorf("compareTo expected > 0")
	}

	// isAfter / isBefore / isEqual
	if zonedDateTimeIsAfter([]any{zdtLater, zdt1}) != types.JavaBoolTrue {
		t.Errorf("isAfter expected true")
	}
	if zonedDateTimeIsBefore([]any{zdt1, zdtLater}) != types.JavaBoolTrue {
		t.Errorf("isBefore expected true")
	}
	if zonedDateTimeIsEqual([]any{zdt1, zdt2}) != types.JavaBoolTrue {
		t.Errorf("isEqual expected true")
	}
}

func TestZonedDateTimeDSTTransitions(t *testing.T) {
	initZonedDateTimeTestEnvironment()

	nyZone := makeZoneIdFromID("America/New_York")

	// 1. Spring-forward gap: 2026-03-08 02:30 does not exist in America/New_York;
	// it should resolve forward by 1 hour to 03:30 with offset -04:00.
	nonExistentLDT := createLocalDateTime(CreateLocalDate(2026, 3, 8), createLocalTime(2, 30, 0, 0))
	resolved := zonedDateTimeOfLocal(nonExistentLDT, nyZone, nil).(*object.Object)
	resolvedStr := object.GoStringFromStringObject(zonedDateTimeToString([]any{resolved}).(*object.Object))
	expectedGapStr := "2026-03-08T03:30-04:00[America/New_York]"
	if resolvedStr != expectedGapStr {
		t.Errorf("Spring-forward gap resolution expected %s, got %s", expectedGapStr, resolvedStr)
	}

	// 2. Fall-back overlap: 2026-11-01 00:30-04:00 plus 2 hours should be 01:30-05:00
	beforeFallBack := zonedDateTimeOfInts([]any{int64(2026), int64(11), int64(1), int64(0), int64(30), int64(0), int64(0), nyZone}).(*object.Object)
	plusTwoHours := zonedDateTimePlusHours([]any{beforeFallBack, int64(2)}).(*object.Object)
	plusTwoStr := object.GoStringFromStringObject(zonedDateTimeToString([]any{plusTwoHours}).(*object.Object))
	expectedFallBackStr := "2026-11-01T01:30-05:00[America/New_York]"
	if plusTwoStr != expectedFallBackStr {
		t.Errorf("Fall-back plusHours(2) expected %s, got %s", expectedFallBackStr, plusTwoStr)
	}

	// 3. now() uses systemDefault zone
	nowObj := zonedDateTimeNow(nil).(*object.Object)
	nowZone := zonedDateTimeGetZone([]any{nowObj}).(*object.Object)
	defZone := zoneIdSystemDefault(nil).(*object.Object)
	nowZoneStr := object.GoStringFromStringObject(zoneIdToString([]any{nowZone}).(*object.Object))
	defZoneStr := object.GoStringFromStringObject(zoneIdToString([]any{defZone}).(*object.Object))
	if nowZoneStr != defZoneStr {
		t.Errorf("now() zone expected %s, got %s", defZoneStr, nowZoneStr)
	}
}
