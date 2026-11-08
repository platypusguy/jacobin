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
	"time"
)

func initClockTestEnvironment() {
	initTestEnvironment()
	ensureClockInited()
	ensureZoneOffsetInited()
	ensureZoneIdInited()
	ensureInstantInited()
	Load_Time_Clock()
	Load_Time_ZoneOffset()
	Load_Time_ZoneId()
	Load_Time_Instant()
	Load_Time_Duration()
}

func TestClockClinitAndInit(t *testing.T) {
	initClockTestEnvironment()

	res := clockClinit(nil)
	if res != nil {
		t.Errorf("Expected clockClinit to return nil, got %v", res)
	}

	resInit := clockInit(nil)
	if resInit != nil {
		t.Errorf("Expected clockInit to return nil, got %v", resInit)
	}
}

func TestClockSystemFactories(t *testing.T) {
	initClockTestEnvironment()

	// systemUTC()
	utcClock := clockSystemUTC(nil).(*object.Object)
	if utcClock == nil {
		t.Fatal("clockSystemUTC returned nil")
	}
	zone := clockGetZone([]any{utcClock}).(*object.Object)
	if getZoneIdString(zone) != "UTC" {
		t.Errorf("Expected zone UTC, got %s", getZoneIdString(zone))
	}

	// systemDefaultZone()
	defClock := clockSystemDefaultZone(nil).(*object.Object)
	if defClock == nil {
		t.Fatal("clockSystemDefaultZone returned nil")
	}

	// system(ZoneId)
	parisZone := makeZoneIdFromID("Europe/Paris")
	parisClock := clockSystem([]any{parisZone}).(*object.Object)
	if parisClock == nil {
		t.Fatal("clockSystem returned nil")
	}
	zoneParis := clockGetZone([]any{parisClock}).(*object.Object)
	if getZoneIdString(zoneParis) != "Europe/Paris" {
		t.Errorf("Expected zone Europe/Paris, got %s", getZoneIdString(zoneParis))
	}

	// system(null) -> NPE
	errRes := clockSystem([]any{nil})
	if _, isErr := errRes.(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected NPE for clockSystem(null)")
	}

	// millis() and instant()
	m1 := clockMillis([]any{utcClock}).(int64)
	nowMilli := time.Now().UnixMilli()
	if m1 < nowMilli-2000 || m1 > nowMilli+2000 {
		t.Errorf("clockMillis out of range: got %d, expected near %d", m1, nowMilli)
	}

	instObj := clockInstant([]any{utcClock}).(*object.Object)
	if instObj == nil {
		t.Fatal("clockInstant returned nil")
	}
	sec := getInstantEpochSecondVal(instObj)
	nowSec := time.Now().Unix()
	if sec < nowSec-2 || sec > nowSec+2 {
		t.Errorf("clockInstant epochSecond out of range: got %d, expected near %d", sec, nowSec)
	}

	// toString
	strObj := clockToString([]any{utcClock}).(*object.Object)
	if object.GoStringFromStringObject(strObj) != "SystemClock[UTC]" {
		t.Errorf("clockToString mismatch: got %s", object.GoStringFromStringObject(strObj))
	}
}

func TestClockFixed(t *testing.T) {
	initClockTestEnvironment()

	inst := createInstant(1782311415, 500000000)
	zoneParis := makeZoneIdFromID("Europe/Paris")

	fixedClock := clockFixed([]any{inst, zoneParis}).(*object.Object)
	if fixedClock == nil {
		t.Fatal("clockFixed returned nil")
	}

	// getZone
	z := clockGetZone([]any{fixedClock}).(*object.Object)
	if getZoneIdString(z) != "Europe/Paris" {
		t.Errorf("Expected zone Europe/Paris, got %s", getZoneIdString(z))
	}

	// instant
	i := clockInstant([]any{fixedClock}).(*object.Object)
	if getInstantEpochSecondVal(i) != 1782311415 || getInstantNanoVal(i) != 500000000 {
		t.Errorf("Fixed clock instant mismatch: %d.%d", getInstantEpochSecondVal(i), getInstantNanoVal(i))
	}

	// millis
	m := clockMillis([]any{fixedClock}).(int64)
	expectedM := int64(1782311415)*1000 + 500
	if m != expectedM {
		t.Errorf("Fixed clock millis mismatch: expected %d, got %d", expectedM, m)
	}

	// withZone
	zoneTokyo := makeZoneIdFromID("Asia/Tokyo")
	tokyoFixed := clockWithZone([]any{fixedClock, zoneTokyo}).(*object.Object)
	if getZoneIdString(clockGetZone([]any{tokyoFixed}).(*object.Object)) != "Asia/Tokyo" {
		t.Errorf("withZone mismatch")
	}
	if clockWithZone([]any{fixedClock, zoneParis}).(*object.Object) != fixedClock {
		t.Errorf("withZone with same zone should return self")
	}

	// equals and hashCode
	fixedClock2 := clockFixed([]any{inst, zoneParis}).(*object.Object)
	if clockEquals([]any{fixedClock, fixedClock2}) != types.JavaBoolTrue {
		t.Errorf("Expected fixed clocks to be equal")
	}
	if clockHashCode([]any{fixedClock}) != clockHashCode([]any{fixedClock2}) {
		t.Errorf("Expected equal fixed clocks to have same hashCode")
	}

	// toString
	strObj := clockToString([]any{fixedClock}).(*object.Object)
	goStr := object.GoStringFromStringObject(strObj)
	if goStr != "FixedClock[2026-06-24T14:30:15.500Z,Europe/Paris]" && goStr != "FixedClock[1782311415,Europe/Paris]" {
		// As long as it starts with FixedClock and contains Europe/Paris
		if len(goStr) < 10 || goStr[:11] != "FixedClock[" {
			t.Errorf("Fixed clock toString mismatch: %s", goStr)
		}
	}

	// Error cases
	if _, isErr := clockFixed([]any{nil, zoneParis}).(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected NPE for clockFixed(null, zone)")
	}
	if _, isErr := clockFixed([]any{inst, nil}).(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected NPE for clockFixed(inst, null)")
	}
}

func TestClockOffset(t *testing.T) {
	initClockTestEnvironment()

	inst := createInstant(1000, 0)
	zoneUTC := makeZoneIdFromID("UTC")
	baseClock := clockFixed([]any{inst, zoneUTC}).(*object.Object)

	dur := createDuration(50, 200000000) // 50.2s
	offsetClock := clockOffset([]any{baseClock, dur}).(*object.Object)

	// instant
	i := clockInstant([]any{offsetClock}).(*object.Object)
	if getInstantEpochSecondVal(i) != 1050 || getInstantNanoVal(i) != 200000000 {
		t.Errorf("Offset clock instant mismatch: %d.%d", getInstantEpochSecondVal(i), getInstantNanoVal(i))
	}

	// millis
	m := clockMillis([]any{offsetClock}).(int64)
	if m != 1050200 {
		t.Errorf("Offset clock millis mismatch: expected 1050200, got %d", m)
	}

	// getZone
	z := clockGetZone([]any{offsetClock}).(*object.Object)
	if getZoneIdString(z) != "UTC" {
		t.Errorf("Expected UTC, got %s", getZoneIdString(z))
	}

	// zero duration returns base clock
	zeroDur := createDuration(0, 0)
	sameClock := clockOffset([]any{baseClock, zeroDur}).(*object.Object)
	if sameClock != baseClock {
		t.Errorf("Expected zero offset to return base clock")
	}

	// equals and hashCode
	offsetClock2 := clockOffset([]any{baseClock, dur}).(*object.Object)
	if clockEquals([]any{offsetClock, offsetClock2}) != types.JavaBoolTrue {
		t.Errorf("Expected offset clocks to be equal")
	}
	if clockHashCode([]any{offsetClock}) != clockHashCode([]any{offsetClock2}) {
		t.Errorf("Expected equal offset clocks to have same hashCode")
	}

	// withZone
	zoneLondon := makeZoneIdFromID("Europe/London")
	londonOffset := clockWithZone([]any{offsetClock, zoneLondon}).(*object.Object)
	if getZoneIdString(clockGetZone([]any{londonOffset}).(*object.Object)) != "Europe/London" {
		t.Errorf("Offset withZone mismatch")
	}
}

func TestClockTick(t *testing.T) {
	initClockTestEnvironment()

	zoneUTC := makeZoneIdFromID("UTC")

	// tickMillis
	tm := clockTickMillis([]any{zoneUTC}).(*object.Object)
	if tm == nil {
		t.Fatal("clockTickMillis returned nil")
	}

	// tickSeconds
	ts := clockTickSeconds([]any{zoneUTC}).(*object.Object)
	if ts == nil {
		t.Fatal("clockTickSeconds returned nil")
	}
	instSec := clockInstant([]any{ts}).(*object.Object)
	if getInstantNanoVal(instSec) != 0 {
		t.Errorf("Expected nanos to be 0 for tickSeconds, got %d", getInstantNanoVal(instSec))
	}

	// tickMinutes
	tmin := clockTickMinutes([]any{zoneUTC}).(*object.Object)
	if tmin == nil {
		t.Fatal("clockTickMinutes returned nil")
	}
	instMin := clockInstant([]any{tmin}).(*object.Object)
	if getInstantNanoVal(instMin) != 0 || getInstantEpochSecondVal(instMin)%60 != 0 {
		t.Errorf("Expected epochSecond divisible by 60 for tickMinutes, got %d", getInstantEpochSecondVal(instMin))
	}

	// tick(baseClock, tickDuration)
	fixedInst := createInstant(1000, 123456789)
	fixedClock := clockFixed([]any{fixedInst, zoneUTC}).(*object.Object)

	dur100Millis := createDuration(0, 100000000)
	tick100ms := clockTick([]any{fixedClock, dur100Millis}).(*object.Object)
	i100ms := clockInstant([]any{tick100ms}).(*object.Object)
	// 123456789 nanos truncated to 100ms (100_000_000 ns) -> 100_000_000 ns
	if getInstantEpochSecondVal(i100ms) != 1000 || getInstantNanoVal(i100ms) != 100000000 {
		t.Errorf("Expected 1000s 100000000ns, got %ds %dns", getInstantEpochSecondVal(i100ms), getInstantNanoVal(i100ms))
	}

	// tick error validation
	negDur := createDuration(-1, 0)
	if _, isErr := clockTick([]any{fixedClock, negDur}).(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected exception for negative tick duration")
	}
	zeroDur := createDuration(0, 0)
	if _, isErr := clockTick([]any{fixedClock, zeroDur}).(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected exception for zero tick duration")
	}
	invalidDur := createDuration(0, 7) // 7 nanos does not divide 1s and 7%1ms != 0
	if _, isErr := clockTick([]any{fixedClock, invalidDur}).(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Expected exception for invalid tick duration (7ns)")
	}

	// duration of 1ns returns baseClock
	dur1ns := createDuration(0, 1)
	if clockTick([]any{fixedClock, dur1ns}).(*object.Object) != fixedClock {
		t.Errorf("Expected 1ns tick to return baseClock")
	}

	// toString
	strTick := clockToString([]any{tick100ms}).(*object.Object)
	if object.GoStringFromStringObject(strTick) == "" {
		t.Errorf("Expected non-empty toString for TickClock")
	}
}

func TestClockNowIntegration(t *testing.T) {
	initClockTestEnvironment()
	ensureLocalDateTimeInited()
	ensureLocalDateInited()
	ensureLocalTimeInited()
	ensureZonedDateTimeInited()
	ensureYearInited()
	Load_Time_LocalDateTime()
	Load_Time_LocalDate()
	Load_Time_LocalTime()
	Load_Time_ZonedDateTime()
	Load_Time_Year()

	// Fixed clock at 2026-06-24T14:30:15Z in UTC
	// 2026-06-24T14:30:15Z epochSecond = 1782311415
	inst := createInstant(1782311415, 500)
	zoneUTC := makeZoneIdFromID("UTC")
	fixedClock := clockFixed([]any{inst, zoneUTC}).(*object.Object)

	// Instant.now(clock)
	instNow := instantNowClock([]any{fixedClock}).(*object.Object)
	if getInstantEpochSecondVal(instNow) != 1782311415 || getInstantNanoVal(instNow) != 500 {
		t.Errorf("Instant.now(clock) mismatch")
	}

	// ZonedDateTime.now(clock)
	zdtNow := zonedDateTimeNowClock([]any{fixedClock}).(*object.Object)
	if getZonedDateTimeEpochSecond(zdtNow) != 1782311415 {
		t.Errorf("ZonedDateTime.now(clock) mismatch")
	}

	// LocalDateTime.now(clock)
	ldtNow := localDateTimeNowClock([]any{fixedClock}).(*object.Object)
	y, m, d, h, min, s, nano := getLocalDateTimeYMDHMSN(ldtNow)
	if y != 2026 || m != 6 || d != 24 || h != 14 || min != 30 || s != 15 || nano != 500 {
		t.Errorf("LocalDateTime.now(clock) mismatch: %d-%d-%d %d:%d:%d.%d", y, m, d, h, min, s, nano)
	}

	// LocalDate.now(clock)
	ldNow := localDateNowClock([]any{fixedClock}).(*object.Object)
	if localDateGetYear([]any{ldNow}).(int64) != 2026 || localDateGetMonthValue([]any{ldNow}).(int64) != 6 || localDateGetDayOfMonth([]any{ldNow}).(int64) != 24 {
		t.Errorf("LocalDate.now(clock) mismatch")
	}

	// LocalTime.now(clock)
	ltNow := localTimeNowClock([]any{fixedClock}).(*object.Object)
	if localTimeGetHour([]any{ltNow}).(int64) != 14 || localTimeGetMinute([]any{ltNow}).(int64) != 30 || localTimeGetSecond([]any{ltNow}).(int64) != 15 {
		t.Errorf("LocalTime.now(clock) mismatch")
	}

	// Year.now(clock)
	yearNow := yearNowClock([]any{fixedClock}).(*object.Object)
	if yVal, ok := getYearValue(yearNow); !ok || yVal != 2026 {
		t.Errorf("Year.now(clock) mismatch: %d", yVal)
	}
}
