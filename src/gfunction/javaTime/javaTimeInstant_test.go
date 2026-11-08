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

func initInstantTestEnv() {
	globals.InitGlobals("test")
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
	Load_Time_Temporal_TemporalQueries()
	Load_Time_Temporal_TemporalAccessor()
	Load_Time_Temporal_Temporal()
	Load_Time_Temporal_ValueRange()
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
	Load_Time_ZoneId()
	Load_Time_Instant()
}

func TestInstantClinitAndStaticFields(t *testing.T) {
	initInstantTestEnv()

	epochStat := statics.GetStaticValue(classNameInstantVal, "EPOCH")
	if epochStat == nil {
		t.Fatalf("Instant.EPOCH static field not found")
	}
	epochObj := epochStat.(*object.Object)
	s, n := getInstantSecondsAndNanos(epochObj)
	if s != 0 || n != 0 {
		t.Errorf("Instant.EPOCH expected (0, 0), got (%d, %d)", s, n)
	}

	minStat := statics.GetStaticValue(classNameInstantVal, "MIN")
	if minStat == nil {
		t.Fatalf("Instant.MIN static field not found")
	}
	minObj := minStat.(*object.Object)
	sMin, nMin := getInstantSecondsAndNanos(minObj)
	if sMin != instantMinSecond || nMin != 0 {
		t.Errorf("Instant.MIN expected (%d, 0), got (%d, %d)", instantMinSecond, sMin, nMin)
	}

	maxStat := statics.GetStaticValue(classNameInstantVal, "MAX")
	if maxStat == nil {
		t.Fatalf("Instant.MAX static field not found")
	}
	maxObj := maxStat.(*object.Object)
	sMax, nMax := getInstantSecondsAndNanos(maxObj)
	if sMax != instantMaxSecond || nMax != 999999999 {
		t.Errorf("Instant.MAX expected (%d, 999999999), got (%d, %d)", instantMaxSecond, sMax, nMax)
	}
}

func TestInstantNowAndFactories(t *testing.T) {
	initInstantTestEnv()

	// now()
	nowRes := instantNow(nil)
	nowObj, ok := nowRes.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from instantNow, got %T", nowRes)
	}
	sNow, _ := getInstantSecondsAndNanos(nowObj)
	if sNow < 1700000000 { // Reasonable sanity check for modern epoch
		t.Errorf("Unexpected instantNow epoch second: %d", sNow)
	}

	// ofEpochSecond(long)
	inst1Res := instantOfEpochSecond1([]any{int64(1234567890)})
	inst1, ok := inst1Res.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object from ofEpochSecond, got %T", inst1Res)
	}
	s, n := getInstantSecondsAndNanos(inst1)
	if s != 1234567890 || n != 0 {
		t.Errorf("ofEpochSecond(1234567890) expected (1234567890, 0), got (%d, %d)", s, n)
	}

	// ofEpochSecond(long, long) with nano normalization
	inst2Res := instantOfEpochSecond2([]any{int64(100), int64(1500000000)})
	inst2 := inst2Res.(*object.Object)
	s2, n2 := getInstantSecondsAndNanos(inst2)
	if s2 != 101 || n2 != 500000000 {
		t.Errorf("ofEpochSecond(100, 1500000000) expected (101, 500000000), got (%d, %d)", s2, n2)
	}

	// ofEpochSecond with negative nano adjustment
	inst3Res := instantOfEpochSecond2([]any{int64(100), int64(-500000000)})
	inst3 := inst3Res.(*object.Object)
	s3, n3 := getInstantSecondsAndNanos(inst3)
	if s3 != 99 || n3 != 500000000 {
		t.Errorf("ofEpochSecond(100, -500000000) expected (99, 500000000), got (%d, %d)", s3, n3)
	}

	// ofEpochMilli
	instMilliRes := instantOfEpochMilli([]any{int64(1234567890123)})
	instMilli := instMilliRes.(*object.Object)
	sM, nM := getInstantSecondsAndNanos(instMilli)
	if sM != 1234567890 || nM != 123000000 {
		t.Errorf("ofEpochMilli(1234567890123) expected (1234567890, 123000000), got (%d, %d)", sM, nM)
	}

	// ofEpochMilli negative
	instMilliNegRes := instantOfEpochMilli([]any{int64(-1500)})
	instMilliNeg := instMilliNegRes.(*object.Object)
	sMNeg, nMNeg := getInstantSecondsAndNanos(instMilliNeg)
	if sMNeg != -2 || nMNeg != 500000000 {
		t.Errorf("ofEpochMilli(-1500) expected (-2, 500000000), got (%d, %d)", sMNeg, nMNeg)
	}
}

func TestInstantParseAndToString(t *testing.T) {
	initInstantTestEnv()

	testCases := []struct {
		input       string
		expectedSec int64
		expectedNos int64
		expectedStr string
	}{
		{"1970-01-01T00:00:00Z", 0, 0, "1970-01-01T00:00:00Z"},
		{"2007-12-03T10:15:30.000Z", 1196676930, 0, "2007-12-03T10:15:30Z"},
		{"2007-12-03T10:15:30.123Z", 1196676930, 123000000, "2007-12-03T10:15:30.123Z"},
		{"2007-12-03T10:15:30.123456Z", 1196676930, 123456000, "2007-12-03T10:15:30.123456Z"},
		{"2007-12-03T10:15:30.123456789Z", 1196676930, 123456789, "2007-12-03T10:15:30.123456789Z"},
		{"2007-12-03T12:15:30+02:00", 1196676930, 0, "2007-12-03T10:15:30Z"},
		{"2007-12-03T05:15:30-05:00", 1196676930, 0, "2007-12-03T10:15:30Z"},
	}

	for _, tc := range testCases {
		strObj := object.StringObjectFromGoString(tc.input)
		parseRes := instantParse([]any{strObj})
		inst, ok := parseRes.(*object.Object)
		if !ok {
			t.Fatalf("Failed to parse %s: got %v", tc.input, parseRes)
		}
		s, n := getInstantSecondsAndNanos(inst)
		if s != tc.expectedSec || n != tc.expectedNos {
			t.Errorf("parse(%s) expected (%d, %d), got (%d, %d)", tc.input, tc.expectedSec, tc.expectedNos, s, n)
		}

		toStrRes := instantToString([]any{inst})
		strRes := object.GoStringFromStringObject(toStrRes.(*object.Object))
		if strRes != tc.expectedStr {
			t.Errorf("toString for %s expected %s, got %s", tc.input, tc.expectedStr, strRes)
		}
	}
}

func TestInstantFieldsAndRange(t *testing.T) {
	initInstantTestEnv()

	inst := createInstant(1234567890, 987654321)

	cfInstantSec := getChronoFieldForInstant("INSTANT_SECONDS")
	cfNanoOfSec := getChronoFieldForInstant("NANO_OF_SECOND")
	cfMicroOfSec := getChronoFieldForInstant("MICRO_OF_SECOND")
	cfMilliOfSec := getChronoFieldForInstant("MILLI_OF_SECOND")
	cfDayOfMonth := getChronoFieldForInstant("DAY_OF_MONTH")

	if instantIsSupportedField([]any{inst, cfInstantSec}) != types.JavaBoolTrue {
		t.Errorf("Expected INSTANT_SECONDS to be supported")
	}
	if instantIsSupportedField([]any{inst, cfNanoOfSec}) != types.JavaBoolTrue {
		t.Errorf("Expected NANO_OF_SECOND to be supported")
	}
	if instantIsSupportedField([]any{inst, cfDayOfMonth}) != types.JavaBoolFalse {
		t.Errorf("Expected DAY_OF_MONTH to not be supported")
	}

	// getLong
	if instantGetLong([]any{inst, cfInstantSec}) != int64(1234567890) {
		t.Errorf("getLong(INSTANT_SECONDS) mismatch")
	}
	if instantGetLong([]any{inst, cfNanoOfSec}) != int64(987654321) {
		t.Errorf("getLong(NANO_OF_SECOND) mismatch")
	}
	if instantGetLong([]any{inst, cfMicroOfSec}) != int64(987654) {
		t.Errorf("getLong(MICRO_OF_SECOND) mismatch")
	}
	if instantGetLong([]any{inst, cfMilliOfSec}) != int64(987) {
		t.Errorf("getLong(MILLI_OF_SECOND) mismatch")
	}

	// get (int)
	if instantGet([]any{inst, cfNanoOfSec}) != int64(987654321) {
		t.Errorf("get(NANO_OF_SECOND) mismatch")
	}
	if instantGet([]any{inst, cfMicroOfSec}) != int64(987654) {
		t.Errorf("get(MICRO_OF_SECOND) mismatch")
	}
	if instantGet([]any{inst, cfMilliOfSec}) != int64(987) {
		t.Errorf("get(MILLI_OF_SECOND) mismatch")
	}

	// get(INSTANT_SECONDS) should fail
	errGet := instantGet([]any{inst, cfInstantSec})
	if _, isErr := errGet.(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected error when calling get(INSTANT_SECONDS)")
	}

	// getEpochSecond and getNano
	if instantGetEpochSecond([]any{inst}) != int64(1234567890) {
		t.Errorf("getEpochSecond mismatch")
	}
	if instantGetNano([]any{inst}) != int64(987654321) {
		t.Errorf("getNano mismatch")
	}

	// with field
	withNano := instantWithField([]any{inst, cfNanoOfSec, int64(111222333)}).(*object.Object)
	_, wN := getInstantSecondsAndNanos(withNano)
	if wN != 111222333 {
		t.Errorf("with(NANO_OF_SECOND) mismatch: got %d", wN)
	}

	withSec := instantWithField([]any{inst, cfInstantSec, int64(99999)}).(*object.Object)
	wS, _ := getInstantSecondsAndNanos(withSec)
	if wS != 99999 {
		t.Errorf("with(INSTANT_SECONDS) mismatch: got %d", wS)
	}
}

func TestInstantTruncatedTo(t *testing.T) {
	initInstantTestEnv()

	// 2007-12-03T10:15:30.123456789Z (epoch sec = 1196676930)
	inst := createInstant(1196676930, 123456789)

	uNanos := ghelpers.CreateTemporalUnit("NANOS")
	uMicros := ghelpers.CreateTemporalUnit("MICROS")
	uMillis := ghelpers.CreateTemporalUnit("MILLIS")
	uSeconds := ghelpers.CreateTemporalUnit("SECONDS")
	uMinutes := ghelpers.CreateTemporalUnit("MINUTES")
	uHours := ghelpers.CreateTemporalUnit("HOURS")
	uDays := ghelpers.CreateTemporalUnit("DAYS")

	// Truncate NANOS
	tNanos := instantTruncatedTo([]any{inst, uNanos}).(*object.Object)
	s, n := getInstantSecondsAndNanos(tNanos)
	if s != 1196676930 || n != 123456789 {
		t.Errorf("truncatedTo(NANOS) mismatch: (%d, %d)", s, n)
	}

	// Truncate MICROS
	tMicros := instantTruncatedTo([]any{inst, uMicros}).(*object.Object)
	s, n = getInstantSecondsAndNanos(tMicros)
	if s != 1196676930 || n != 123456000 {
		t.Errorf("truncatedTo(MICROS) mismatch: (%d, %d)", s, n)
	}

	// Truncate MILLIS
	tMillis := instantTruncatedTo([]any{inst, uMillis}).(*object.Object)
	s, n = getInstantSecondsAndNanos(tMillis)
	if s != 1196676930 || n != 123000000 {
		t.Errorf("truncatedTo(MILLIS) mismatch: (%d, %d)", s, n)
	}

	// Truncate SECONDS
	tSec := instantTruncatedTo([]any{inst, uSeconds}).(*object.Object)
	s, n = getInstantSecondsAndNanos(tSec)
	if s != 1196676930 || n != 0 {
		t.Errorf("truncatedTo(SECONDS) mismatch: (%d, %d)", s, n)
	}

	// Truncate MINUTES (1196676930 % 60 = 30 -> 1196676900)
	tMin := instantTruncatedTo([]any{inst, uMinutes}).(*object.Object)
	s, n = getInstantSecondsAndNanos(tMin)
	if s != 1196676900 || n != 0 {
		t.Errorf("truncatedTo(MINUTES) mismatch: (%d, %d)", s, n)
	}

	// Truncate HOURS (10:15:30 -> 10:00:00 -> sec = 1196676000)
	tHour := instantTruncatedTo([]any{inst, uHours}).(*object.Object)
	s, n = getInstantSecondsAndNanos(tHour)
	if s != 1196676000 || n != 0 {
		t.Errorf("truncatedTo(HOURS) mismatch: (%d, %d)", s, n)
	}

	// Truncate DAYS (2007-12-03T00:00:00Z -> sec = 1196640000)
	tDay := instantTruncatedTo([]any{inst, uDays}).(*object.Object)
	s, n = getInstantSecondsAndNanos(tDay)
	if s != 1196640000 || n != 0 {
		t.Errorf("truncatedTo(DAYS) mismatch: (%d, %d)", s, n)
	}
}

func TestInstantArithmetic(t *testing.T) {
	initInstantTestEnv()

	inst := createInstant(100, 500000000)

	// plusSeconds
	pSec := instantPlusSeconds([]any{inst, int64(10)}).(*object.Object)
	s, n := getInstantSecondsAndNanos(pSec)
	if s != 110 || n != 500000000 {
		t.Errorf("plusSeconds(10) expected (110, 500000000), got (%d, %d)", s, n)
	}

	// plusMillis
	pMil := instantPlusMillis([]any{inst, int64(600)}).(*object.Object)
	s, n = getInstantSecondsAndNanos(pMil)
	if s != 101 || n != 100000000 {
		t.Errorf("plusMillis(600) expected (101, 100000000), got (%d, %d)", s, n)
	}

	// plusNanos
	pNano := instantPlusNanos([]any{inst, int64(600000000)}).(*object.Object)
	s, n = getInstantSecondsAndNanos(pNano)
	if s != 101 || n != 100000000 {
		t.Errorf("plusNanos(600000000) expected (101, 100000000), got (%d, %d)", s, n)
	}

	// minusSeconds
	mSec := instantMinusSeconds([]any{inst, int64(10)}).(*object.Object)
	s, n = getInstantSecondsAndNanos(mSec)
	if s != 90 || n != 500000000 {
		t.Errorf("minusSeconds(10) expected (90, 500000000), got (%d, %d)", s, n)
	}

	// minusMillis
	mMil := instantMinusMillis([]any{inst, int64(600)}).(*object.Object)
	s, n = getInstantSecondsAndNanos(mMil)
	if s != 99 || n != 900000000 {
		t.Errorf("minusMillis(600) expected (99, 900000000), got (%d, %d)", s, n)
	}

	// minusNanos
	mNano := instantMinusNanos([]any{inst, int64(600000000)}).(*object.Object)
	s, n = getInstantSecondsAndNanos(mNano)
	if s != 99 || n != 900000000 {
		t.Errorf("minusNanos(600000000) expected (99, 900000000), got (%d, %d)", s, n)
	}

	// plus unit HOURS
	uHours := ghelpers.CreateTemporalUnit("HOURS")
	pHour := instantPlusUnit([]any{inst, int64(2), uHours}).(*object.Object)
	s, _ = getInstantSecondsAndNanos(pHour)
	if s != 100+7200 {
		t.Errorf("plus(2, HOURS) expected %d, got %d", 100+7200, s)
	}

	// minus unit DAYS
	uDays := ghelpers.CreateTemporalUnit("DAYS")
	mDay := instantMinusUnit([]any{inst, int64(1), uDays}).(*object.Object)
	s, _ = getInstantSecondsAndNanos(mDay)
	if s != 100-86400 {
		t.Errorf("minus(1, DAYS) expected %d, got %d", 100-86400, s)
	}
}

func TestInstantUntil(t *testing.T) {
	initInstantTestEnv()

	inst1 := createInstant(100, 200000000)
	inst2 := createInstant(105, 700000000)

	uNanos := ghelpers.CreateTemporalUnit("NANOS")
	uMillis := ghelpers.CreateTemporalUnit("MILLIS")
	uSeconds := ghelpers.CreateTemporalUnit("SECONDS")

	// until NANOS: (5 * 1B) + (700M - 200M) = 5500000000
	nanosDiff := instantUntil([]any{inst1, inst2, uNanos})
	if nanosDiff != int64(5500000000) {
		t.Errorf("until(NANOS) expected 5500000000, got %v", nanosDiff)
	}

	// until MILLIS: 5500 ms
	millisDiff := instantUntil([]any{inst1, inst2, uMillis})
	if millisDiff != int64(5500) {
		t.Errorf("until(MILLIS) expected 5500, got %v", millisDiff)
	}

	// until SECONDS: 5 s
	secsDiff := instantUntil([]any{inst1, inst2, uSeconds})
	if secsDiff != int64(5) {
		t.Errorf("until(SECONDS) expected 5, got %v", secsDiff)
	}
}

func TestInstantAtOffsetAndAtZone(t *testing.T) {
	initInstantTestEnv()

	// 2007-12-03T10:15:30Z
	inst := createInstant(1196676930, 0)
	offsetObj := createZoneOffset(3600) // +01:00

	odt := instantAtOffset([]any{inst, offsetObj}).(*object.Object)
	if odt == nil {
		t.Fatalf("atOffset returned nil")
	}
	ldt := odt.FieldTable["dateTime"].Fvalue.(*object.Object)
	y, m, d, h, min, s, _ := getLocalDateTimeYMDHMSN(ldt)
	if y != 2007 || m != 12 || d != 3 || h != 11 || min != 15 || s != 30 {
		t.Errorf("atOffset(+01:00) expected 2007-12-03T11:15:30, got %04d-%02d-%02dT%02d:%02d:%02d", y, m, d, h, min, s)
	}

	// atZone("UTC")
	zoneUTC := makeZoneIdFromID("UTC")
	zdt := instantAtZone([]any{inst, zoneUTC}).(*object.Object)
	if zdt == nil {
		t.Fatalf("atZone returned nil")
	}
	ldtUTC := zdt.FieldTable["dateTime"].Fvalue.(*object.Object)
	y, m, d, h, min, s, _ = getLocalDateTimeYMDHMSN(ldtUTC)
	if y != 2007 || m != 12 || d != 3 || h != 10 || min != 15 || s != 30 {
		t.Errorf("atZone(UTC) expected 2007-12-03T10:15:30, got %04d-%02d-%02dT%02d:%02d:%02d", y, m, d, h, min, s)
	}
}

func TestInstantCompareEqualsHashCodeAndMilli(t *testing.T) {
	initInstantTestEnv()

	inst1 := createInstant(100, 500)
	inst2 := createInstant(100, 500)
	inst3 := createInstant(100, 600)
	inst4 := createInstant(101, 100)

	if instantEquals([]any{inst1, inst2}) != types.JavaBoolTrue {
		t.Errorf("Expected inst1.equals(inst2) to be true")
	}
	if instantEquals([]any{inst1, inst3}) != types.JavaBoolFalse {
		t.Errorf("Expected inst1.equals(inst3) to be false")
	}

	if instantCompareTo([]any{inst1, inst2}) != int64(0) {
		t.Errorf("Expected compareTo(inst1, inst2) == 0")
	}
	if instantCompareTo([]any{inst1, inst3}) != int64(-1) {
		t.Errorf("Expected compareTo(inst1, inst3) == -1")
	}
	if instantCompareTo([]any{inst4, inst1}) != int64(1) {
		t.Errorf("Expected compareTo(inst4, inst1) == 1")
	}

	if instantIsAfter([]any{inst4, inst1}) != types.JavaBoolTrue {
		t.Errorf("Expected isAfter(inst4, inst1) to be true")
	}
	if instantIsBefore([]any{inst1, inst4}) != types.JavaBoolTrue {
		t.Errorf("Expected isBefore(inst1, inst4) to be true")
	}

	if instantHashCode([]any{inst1}) != instantHashCode([]any{inst2}) {
		t.Errorf("HashCodes for equal instants must match")
	}

	// toEpochMilli
	instM := createInstant(1234567890, 500000000)
	if instantToEpochMilli([]any{instM}) != int64(1234567890500) {
		t.Errorf("toEpochMilli expected 1234567890500, got %v", instantToEpochMilli([]any{instM}))
	}
}
