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

func durationTestInit() {
	globals.InitGlobals("test")
	durationMutex.Lock()
	durationInitialized = false
	durationMutex.Unlock()
	Load_Time_Duration()
	Load_Time_Instant()
	Load_Time_LocalTime()
	Load_Time_LocalDate()
	Load_Time_LocalDateTime()
	Load_Time_ZonedDateTime()
	Load_Time_OffsetDateTime()
	Load_Time_Temporal_ChronoUnit()
	Load_Time_Temporal_TemporalUnit()
	Load_Time_Temporal_TemporalAmount()
}

func TestDurationClinitAndFactories(t *testing.T) {
	durationTestInit()

	if res := durationClinit(nil); res != nil {
		t.Errorf("durationClinit returned %v", res)
	}

	st, ok := statics.QueryStatic(classNameDuration, "ZERO")
	if !ok || st.Value == nil {
		t.Fatalf("Duration.ZERO not initialized")
	}
	zeroObj := st.Value.(*object.Object)
	s, n, okSN := getDurationSecondsNanos(zeroObj)
	if !okSN || s != 0 || n != 0 {
		t.Errorf("Duration.ZERO = %ds %dns, want 0s 0ns", s, n)
	}

	// ofSeconds, ofSecondsNanos
	d1 := durationOfSeconds([]any{int64(10)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(d1); s != 10 || n != 0 {
		t.Errorf("ofSeconds(10) = %ds %dns, want 10s 0ns", s, n)
	}

	d2 := durationOfSecondsNanos([]any{int64(10), int64(500000000)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(d2); s != 10 || n != 500000000 {
		t.Errorf("ofSecondsNanos(10, 500000000) = %ds %dns, want 10s 500000000ns", s, n)
	}

	// ofDays, ofHours, ofMinutes, ofMillis, ofNanos
	dDays := durationOfDays([]any{int64(2)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dDays); s != 172800 {
		t.Errorf("ofDays(2) = %ds, want 172800s", s)
	}

	dHours := durationOfHours([]any{int64(3)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dHours); s != 10800 {
		t.Errorf("ofHours(3) = %ds, want 10800s", s)
	}

	dMins := durationOfMinutes([]any{int64(5)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dMins); s != 300 {
		t.Errorf("ofMinutes(5) = %ds, want 300s", s)
	}

	dMillis := durationOfMillis([]any{int64(1500)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dMillis); s != 1 || n != 500000000 {
		t.Errorf("ofMillis(1500) = %ds %dns, want 1s 500000000ns", s, n)
	}

	dNanos := durationOfNanos([]any{int64(2500000000)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dNanos); s != 2 || n != 500000000 {
		t.Errorf("ofNanos(2500000000) = %ds %dns, want 2s 500000000ns", s, n)
	}

	// of(amount, TemporalUnit)
	dOfSec := durationOf([]any{int64(45), ghelpers.CreateTemporalUnit("SECONDS")}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dOfSec); s != 45 {
		t.Errorf("of(45, SECONDS) = %ds, want 45s", s)
	}

	dOfMicros := durationOf([]any{int64(1500000), ghelpers.CreateTemporalUnit("MICROS")}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dOfMicros); s != 1 || n != 500000000 {
		t.Errorf("of(1500000, MICROS) = %ds %dns, want 1s 500000000ns", s, n)
	}

	dOfHalfDays := durationOf([]any{int64(2), ghelpers.CreateTemporalUnit("HALF_DAYS")}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dOfHalfDays); s != 86400 {
		t.Errorf("of(2, HALF_DAYS) = %ds, want 86400s", s)
	}

	errUnit := durationOf([]any{int64(1), ghelpers.CreateTemporalUnit("YEARS")})
	if gerr, ok := errUnit.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeException {
		t.Errorf("of(1, YEARS) expected DateTimeException, got %v", errUnit)
	}
}

func TestDurationBetweenAndFrom(t *testing.T) {
	durationTestInit()

	// between(Instant, Instant)
	i1 := CreateInstant(100, 200)
	i2 := CreateInstant(110, 500)
	dInst := durationBetween([]any{i1, i2}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dInst); s != 10 || n != 300 {
		t.Errorf("between(Instant) = %ds %dns, want 10s 300ns", s, n)
	}

	// between(LocalTime, LocalTime)
	t1 := CreateLocalTime(10, 0, 0, 0)
	t2 := CreateLocalTime(11, 30, 0, 0)
	dTime := durationBetween([]any{t1, t2}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dTime); s != 5400 || n != 0 {
		t.Errorf("between(LocalTime) = %ds %dns, want 5400s 0ns", s, n)
	}

	// between(LocalDateTime, LocalDateTime)
	dt1 := CreateLocalDateTime(CreateLocalDate(2026, 1, 1), CreateLocalTime(0, 0, 0, 0))
	dt2 := CreateLocalDateTime(CreateLocalDate(2026, 1, 2), CreateLocalTime(1, 0, 0, 0))
	dDT := durationBetween([]any{dt1, dt2}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dDT); s != 90000 { // 86400 + 3600
		t.Errorf("between(LocalDateTime) = %ds, want 90000s", s)
	}

	// from(Duration)
	dOrig := CreateDuration(123, 456)
	dFrom := durationFrom([]any{dOrig}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dFrom); s != 123 || n != 456 {
		t.Errorf("from(Duration) = %ds %dns, want 123s 456ns", s, n)
	}

	// from(null)
	errNull := durationFrom([]any{nil})
	if gerr, ok := errNull.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("from(null) expected NullPointerException, got %v", errNull)
	}
}

func TestDurationGettersAndQueries(t *testing.T) {
	durationTestInit()

	d := CreateDuration(3665, 123456789)

	if s := durationGetSeconds([]any{d}).(int64); s != 3665 {
		t.Errorf("getSeconds() = %d, want 3665", s)
	}
	if n := durationGetNano([]any{d}).(int64); n != 123456789 {
		t.Errorf("getNano() = %d, want 123456789", n)
	}

	// get(TemporalUnit)
	if s := durationGet([]any{d, ghelpers.CreateTemporalUnit("SECONDS")}).(int64); s != 3665 {
		t.Errorf("get(SECONDS) = %d, want 3665", s)
	}
	if n := durationGet([]any{d, ghelpers.CreateTemporalUnit("NANOS")}).(int64); n != 123456789 {
		t.Errorf("get(NANOS) = %d, want 123456789", n)
	}
	errUnit := durationGet([]any{d, ghelpers.CreateTemporalUnit("DAYS")})
	if gerr, ok := errUnit.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeException {
		t.Errorf("get(DAYS) expected DateTimeException, got %v", errUnit)
	}

	// getUnits
	unitsList := durationGetUnits([]any{d}).(*object.Object)
	if valFld, ok := unitsList.FieldTable["value"]; !ok || len(valFld.Fvalue.([]any)) != 2 {
		t.Errorf("getUnits() expected list of 2 units, got %v", unitsList)
	}

	// isZero, isNegative, isPositive
	dZero := CreateDuration(0, 0)
	dNeg := CreateDuration(-5, 0)
	dPos := CreateDuration(5, 0)

	if durationIsZero([]any{dZero}).(int64) != types.JavaBoolTrue {
		t.Errorf("dZero.isZero() = false, want true")
	}
	if durationIsZero([]any{dPos}).(int64) != types.JavaBoolFalse {
		t.Errorf("dPos.isZero() = true, want false")
	}
	if durationIsNegative([]any{dNeg}).(int64) != types.JavaBoolTrue {
		t.Errorf("dNeg.isNegative() = false, want true")
	}
	if durationIsPositive([]any{dPos}).(int64) != types.JavaBoolTrue {
		t.Errorf("dPos.isPositive() = false, want true")
	}
}

func TestDurationArithmetic(t *testing.T) {
	durationTestInit()

	d1 := CreateDuration(10, 500000000)
	d2 := CreateDuration(5, 600000000)

	// plus
	dPlus := durationPlus([]any{d1, d2}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dPlus); s != 16 || n != 100000000 {
		t.Errorf("plus = %ds %dns, want 16s 100000000ns", s, n)
	}

	// plus(amount, unit)
	dPlusUnit := durationPlusUnit([]any{d1, int64(10), ghelpers.CreateTemporalUnit("SECONDS")}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dPlusUnit); s != 20 || n != 500000000 {
		t.Errorf("plus(10, SECONDS) = %ds %dns, want 20s 500000000ns", s, n)
	}

	// plusDays, plusHours, plusMinutes, plusSeconds, plusMillis, plusNanos
	dpDays := durationPlusDays([]any{d1, int64(1)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dpDays); s != 86410 {
		t.Errorf("plusDays(1) = %ds, want 86410s", s)
	}
	dpHours := durationPlusHours([]any{d1, int64(1)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dpHours); s != 3610 {
		t.Errorf("plusHours(1) = %ds, want 3610s", s)
	}
	dpMins := durationPlusMinutes([]any{d1, int64(1)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dpMins); s != 70 {
		t.Errorf("plusMinutes(1) = %ds, want 70s", s)
	}
	dpSec := durationPlusSeconds([]any{d1, int64(5)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dpSec); s != 15 {
		t.Errorf("plusSeconds(5) = %ds, want 15s", s)
	}
	dpMillis := durationPlusMillis([]any{d1, int64(600)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dpMillis); s != 11 || n != 100000000 {
		t.Errorf("plusMillis(600) = %ds %dns, want 11s 100000000ns", s, n)
	}
	dpNanos := durationPlusNanos([]any{d1, int64(600000000)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dpNanos); s != 11 || n != 100000000 {
		t.Errorf("plusNanos(600000000) = %ds %dns, want 11s 100000000ns", s, n)
	}

	// minus
	dMinus := durationMinus([]any{dPlus, d2}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dMinus); s != 10 || n != 500000000 {
		t.Errorf("minus = %ds %dns, want 10s 500000000ns", s, n)
	}

	// minus(amount, unit)
	dMinusUnit := durationMinusUnit([]any{d1, int64(5), ghelpers.CreateTemporalUnit("SECONDS")}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dMinusUnit); s != 5 || n != 500000000 {
		t.Errorf("minus(5, SECONDS) = %ds %dns, want 5s 500000000ns", s, n)
	}

	// minusDays, minusHours, minusMinutes, minusSeconds, minusMillis, minusNanos
	dmDays := durationMinusDays([]any{dpDays, int64(1)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dmDays); s != 10 {
		t.Errorf("minusDays(1) = %ds, want 10s", s)
	}
	dmHours := durationMinusHours([]any{dpHours, int64(1)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dmHours); s != 10 {
		t.Errorf("minusHours(1) = %ds, want 10s", s)
	}
	dmMins := durationMinusMinutes([]any{dpMins, int64(1)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dmMins); s != 10 {
		t.Errorf("minusMinutes(1) = %ds, want 10s", s)
	}
	dmSec := durationMinusSeconds([]any{dpSec, int64(5)}).(*object.Object)
	if s, _, _ := getDurationSecondsNanos(dmSec); s != 10 {
		t.Errorf("minusSeconds(5) = %ds, want 10s", s)
	}
	dmMillis := durationMinusMillis([]any{dpMillis, int64(600)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dmMillis); s != 10 || n != 500000000 {
		t.Errorf("minusMillis(600) = %ds %dns, want 10s 500000000ns", s, n)
	}
	dmNanos := durationMinusNanos([]any{dpNanos, int64(600000000)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dmNanos); s != 10 || n != 500000000 {
		t.Errorf("minusNanos(600000000) = %ds %dns, want 10s 500000000ns", s, n)
	}
}

func TestDurationModificationsAndTruncation(t *testing.T) {
	durationTestInit()

	d := CreateDuration(10, 500000000)

	// multipliedBy
	dMul := durationMultipliedBy([]any{d, int64(2)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dMul); s != 21 || n != 0 {
		t.Errorf("multipliedBy(2) = %ds %dns, want 21s 0ns", s, n)
	}

	// dividedBy(long)
	dDivL := durationDividedByLong([]any{dMul, int64(2)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dDivL); s != 10 || n != 500000000 {
		t.Errorf("dividedBy(2) = %ds %dns, want 10s 500000000ns", s, n)
	}

	// dividedBy(Duration)
	dDivD := durationDividedByDuration([]any{dMul, d}).(int64)
	if dDivD != 2 {
		t.Errorf("dividedBy(d) = %d, want 2", dDivD)
	}

	// negated & abs
	dNeg := durationNegated([]any{d}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dNeg); s != -11 || n != 500000000 {
		t.Errorf("negated = %ds %dns, want -11s 500000000ns", s, n)
	}
	dAbs := durationAbs([]any{dNeg}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dAbs); s != 10 || n != 500000000 {
		t.Errorf("abs = %ds %dns, want 10s 500000000ns", s, n)
	}

	// withSeconds, withNanos
	dWS := durationWithSeconds([]any{d, int64(20)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dWS); s != 20 || n != 500000000 {
		t.Errorf("withSeconds(20) = %ds %dns, want 20s 500000000ns", s, n)
	}
	dWN := durationWithNanos([]any{d, int64(100)}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dWN); s != 10 || n != 100 {
		t.Errorf("withNanos(100) = %ds %dns, want 10s 100ns", s, n)
	}

	// truncatedTo
	dTruncSec := durationTruncatedTo([]any{d, ghelpers.CreateTemporalUnit("SECONDS")}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dTruncSec); s != 10 || n != 0 {
		t.Errorf("truncatedTo(SECONDS) = %ds %dns, want 10s 0ns", s, n)
	}

	dHour := CreateDuration(3665, 500000000)
	dTruncHour := durationTruncatedTo([]any{dHour, ghelpers.CreateTemporalUnit("HOURS")}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(dTruncHour); s != 3600 || n != 0 {
		t.Errorf("truncatedTo(HOURS) = %ds %dns, want 3600s 0ns", s, n)
	}
}

func TestDurationConversionsAndParts(t *testing.T) {
	durationTestInit()

	d := CreateDuration(90065, 123456789) // 1 day, 1 hour, 1 min, 5 sec + nanos

	if days := durationToDays([]any{d}).(int64); days != 1 {
		t.Errorf("toDays() = %d, want 1", days)
	}
	if hours := durationToHours([]any{d}).(int64); hours != 25 {
		t.Errorf("toHours() = %d, want 25", hours)
	}
	if mins := durationToMinutes([]any{d}).(int64); mins != 1501 {
		t.Errorf("toMinutes() = %d, want 1501", mins)
	}
	if secs := durationToSeconds([]any{d}).(int64); secs != 90065 {
		t.Errorf("toSeconds() = %d, want 90065", secs)
	}
	if millis := durationToMillis([]any{d}).(int64); millis != 90065123 {
		t.Errorf("toMillis() = %d, want 90065123", millis)
	}

	// Parts
	if dp := durationToDaysPart([]any{d}).(int64); dp != 1 {
		t.Errorf("toDaysPart() = %d, want 1", dp)
	}
	if hp := durationToHoursPart([]any{d}).(int64); hp != 1 {
		t.Errorf("toHoursPart() = %d, want 1", hp)
	}
	if mp := durationToMinutesPart([]any{d}).(int64); mp != 1 {
		t.Errorf("toMinutesPart() = %d, want 1", mp)
	}
	if sp := durationToSecondsPart([]any{d}).(int64); sp != 5 {
		t.Errorf("toSecondsPart() = %d, want 5", sp)
	}
	if mlp := durationToMillisPart([]any{d}).(int64); mlp != 123 {
		t.Errorf("toMillisPart() = %d, want 123", mlp)
	}
	if np := durationToNanosPart([]any{d}).(int64); np != 123456789 {
		t.Errorf("toNanosPart() = %d, want 123456789", np)
	}
}

func TestDurationAddToSubtractFrom(t *testing.T) {
	durationTestInit()

	d := CreateDuration(3600, 0) // 1 hour
	inst := CreateInstant(1000, 0)

	added := durationAddTo([]any{d, inst}).(*object.Object)
	if s, _, _ := getInstantSecNano(added); s != 4600 {
		t.Errorf("addTo(Instant) = %ds, want 4600s", s)
	}

	subbed := durationSubtractFrom([]any{d, added}).(*object.Object)
	if s, _, _ := getInstantSecNano(subbed); s != 1000 {
		t.Errorf("subtractFrom(Instant) = %ds, want 1000s", s)
	}
}

func TestDurationComparisonsAndFormatting(t *testing.T) {
	durationTestInit()

	d1 := CreateDuration(10, 0)
	d2 := CreateDuration(20, 0)
	d3 := CreateDuration(10, 0)

	if cmp := durationCompareTo([]any{d1, d2}).(int64); cmp >= 0 {
		t.Errorf("d1.compareTo(d2) = %d, want <0", cmp)
	}
	if cmp := durationCompareTo([]any{d2, d1}).(int64); cmp <= 0 {
		t.Errorf("d2.compareTo(d1) = %d, want >0", cmp)
	}
	if cmp := durationCompareTo([]any{d1, d3}).(int64); cmp != 0 {
		t.Errorf("d1.compareTo(d3) = %d, want 0", cmp)
	}

	if eq := durationEquals([]any{d1, d3}).(int64); eq != types.JavaBoolTrue {
		t.Errorf("d1.equals(d3) = false, want true")
	}
	if eq := durationEquals([]any{d1, d2}).(int64); eq != types.JavaBoolFalse {
		t.Errorf("d1.equals(d2) = true, want false")
	}

	if h1, h3 := durationHashCode([]any{d1}).(int64), durationHashCode([]any{d3}).(int64); h1 != h3 {
		t.Errorf("hashCode mismatch: %d != %d", h1, h3)
	}

	// toString & parse
	dStr := CreateDuration(3665, 123000000)
	strObj := durationToString([]any{dStr}).(*object.Object)
	str := object.GoStringFromStringObject(strObj)
	if str != "PT1H1M5.123S" {
		t.Errorf("toString = %q, want 'PT1H1M5.123S'", str)
	}

	parsed := durationParse([]any{object.StringObjectFromGoString("PT1H1M5.123S")}).(*object.Object)
	if s, n, _ := getDurationSecondsNanos(parsed); s != 3665 || n != 123000000 {
		t.Errorf("parse('PT1H1M5.123S') = %ds %dns, want 3665s 123000000ns", s, n)
	}
}
