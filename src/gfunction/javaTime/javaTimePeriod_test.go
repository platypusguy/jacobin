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
	"testing"
)

func periodTestInit() {
	globals.InitGlobals("test")
	periodMutex.Lock()
	periodInitialized = false
	periodMutex.Unlock()
	Load_Time_Period()
	Load_Time_LocalDate()
	Load_Time_Temporal_ChronoUnit()
}

func TestPeriodClinitAndFactories(t *testing.T) {
	periodTestInit()

	if res := periodClinit(nil); res != nil {
		t.Errorf("periodClinit returned %v", res)
	}

	st, ok := statics.QueryStatic(classNamePeriod, "ZERO")
	if !ok || st.Value == nil {
		t.Fatalf("Period.ZERO not initialized")
	}
	zeroObj := st.Value.(*object.Object)
	y, m, d, okYMD := getPeriodYMD(zeroObj)
	if !okYMD || y != 0 || m != 0 || d != 0 {
		t.Errorf("Period.ZERO = %d-%d-%d, want 0-0-0", y, m, d)
	}

	// of(Y, M, D)
	p := periodOf([]any{int64(1), int64(2), int64(3)}).(*object.Object)
	if py := periodGetYears([]any{p}).(int64); py != 1 {
		t.Errorf("Period.of(1,2,3).getYears() = %d, want 1", py)
	}
	if pm := periodGetMonths([]any{p}).(int64); pm != 2 {
		t.Errorf("Period.of(1,2,3).getMonths() = %d, want 2", pm)
	}
	if pd := periodGetDays([]any{p}).(int64); pd != 3 {
		t.Errorf("Period.of(1,2,3).getDays() = %d, want 3", pd)
	}

	// ofYears, ofMonths, ofWeeks, ofDays
	pyObj := periodOfYears([]any{int64(5)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(pyObj); y != 5 || m != 0 || d != 0 {
		t.Errorf("ofYears(5) = %d-%d-%d, want 5-0-0", y, m, d)
	}

	pmObj := periodOfMonths([]any{int64(7)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(pmObj); y != 0 || m != 7 || d != 0 {
		t.Errorf("ofMonths(7) = %d-%d-%d, want 0-7-0", y, m, d)
	}

	pwObj := periodOfWeeks([]any{int64(3)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(pwObj); y != 0 || m != 0 || d != 21 {
		t.Errorf("ofWeeks(3) = %d-%d-%d, want 0-0-21", y, m, d)
	}

	pdObj := periodOfDays([]any{int64(15)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(pdObj); y != 0 || m != 0 || d != 15 {
		t.Errorf("ofDays(15) = %d-%d-%d, want 0-0-15", y, m, d)
	}

	// between(d1, d2)
	d1 := CreateLocalDate(2020, 1, 15)
	d2 := CreateLocalDate(2022, 4, 25)
	pBetween := periodBetween([]any{d1, d2}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(pBetween); y != 2 || m != 3 || d != 10 {
		t.Errorf("between(2020-01-15, 2022-04-25) = %d-%d-%d, want 2-3-10", y, m, d)
	}

	// parse
	parseP := periodParse([]any{object.StringObjectFromGoString("P1Y2M3D")}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(parseP); y != 1 || m != 2 || d != 3 {
		t.Errorf("parse('P1Y2M3D') = %d-%d-%d, want 1-2-3", y, m, d)
	}

	parseNeg := periodParse([]any{object.StringObjectFromGoString("P-1Y-2M-3D")}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(parseNeg); y != -1 || m != -2 || d != -3 {
		t.Errorf("parse('P-1Y-2M-3D') = %d-%d-%d, want -1--2--3", y, m, d)
	}

	err := periodParse([]any{object.StringObjectFromGoString("INVALID")})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeParseException {
		t.Errorf("parse('INVALID') expected DateTimeParseException, got %v", err)
	}
}

func TestPeriodOperations(t *testing.T) {
	periodTestInit()

	p1 := CreatePeriod(1, 2, 3)
	p2 := CreatePeriod(2, 4, 6)

	// isZero, isNegative
	pZero := CreatePeriod(0, 0, 0)
	pNeg := CreatePeriod(-1, 0, 0)
	if z := periodIsZero([]any{pZero}).(int64); z != 1 {
		t.Errorf("ZERO.isZero() = %d, want 1", z)
	}
	if z := periodIsZero([]any{p1}).(int64); z != 0 {
		t.Errorf("P1Y2M3D.isZero() = %d, want 0", z)
	}
	if n := periodIsNegative([]any{pNeg}).(int64); n != 1 {
		t.Errorf("P-1Y.isNegative() = %d, want 1", n)
	}
	if n := periodIsNegative([]any{p1}).(int64); n != 0 {
		t.Errorf("P1Y2M3D.isNegative() = %d, want 0", n)
	}

	// withYears, withMonths, withDays
	wy := periodWithYears([]any{p1, int64(10)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(wy); y != 10 || m != 2 || d != 3 {
		t.Errorf("p1.withYears(10) = %d-%d-%d, want 10-2-3", y, m, d)
	}
	wm := periodWithMonths([]any{p1, int64(11)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(wm); y != 1 || m != 11 || d != 3 {
		t.Errorf("p1.withMonths(11) = %d-%d-%d, want 1-11-3", y, m, d)
	}
	wd := periodWithDays([]any{p1, int64(20)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(wd); y != 1 || m != 2 || d != 20 {
		t.Errorf("p1.withDays(20) = %d-%d-%d, want 1-2-20", y, m, d)
	}

	// plus & minus
	plusP := periodPlus([]any{p1, p2}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(plusP); y != 3 || m != 6 || d != 9 {
		t.Errorf("p1.plus(p2) = %d-%d-%d, want 3-6-9", y, m, d)
	}
	minusP := periodMinus([]any{p2, p1}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(minusP); y != 1 || m != 2 || d != 3 {
		t.Errorf("p2.minus(p1) = %d-%d-%d, want 1-2-3", y, m, d)
	}

	// multipliedBy, negated, normalized, toTotalMonths
	mulP := periodMultipliedBy([]any{p1, int64(3)}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(mulP); y != 3 || m != 6 || d != 9 {
		t.Errorf("p1.multipliedBy(3) = %d-%d-%d, want 3-6-9", y, m, d)
	}
	negP := periodNegated([]any{p1}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(negP); y != -1 || m != -2 || d != -3 {
		t.Errorf("p1.negated() = %d-%d-%d, want -1--2--3", y, m, d)
	}

	pToNorm := CreatePeriod(1, 15, 0)
	normP := periodNormalized([]any{pToNorm}).(*object.Object)
	if y, m, d, _ := getPeriodYMD(normP); y != 2 || m != 3 || d != 0 {
		t.Errorf("P1Y15M.normalized() = %d-%d-%d, want 2-3-0", y, m, d)
	}

	totalM := periodToTotalMonths([]any{p1}).(int64)
	if totalM != 14 { // 1*12 + 2
		t.Errorf("p1.toTotalMonths() = %d, want 14", totalM)
	}

	// addTo & subtractFrom on LocalDate
	baseDate := CreateLocalDate(2026, 1, 1)
	addedDate := periodAddTo([]any{p1, baseDate}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(addedDate); y != 2027 || m != 3 || d != 4 {
		t.Errorf("baseDate + p1 = %d-%d-%d, want 2027-3-4", y, m, d)
	}

	subDate := periodSubtractFrom([]any{p1, addedDate}).(*object.Object)
	if y, m, d, _ := getLocalDateYMD(subDate); y != 2026 || m != 1 || d != 1 {
		t.Errorf("addedDate - p1 = %d-%d-%d, want 2026-1-1", y, m, d)
	}

	// equals, hashCode, toString
	if eq := periodEquals([]any{p1, CreatePeriod(1, 2, 3)}).(int64); eq != 1 {
		t.Errorf("p1.equals(P1Y2M3D) = %d, want 1", eq)
	}
	if eq := periodEquals([]any{p1, p2}).(int64); eq != 0 {
		t.Errorf("p1.equals(p2) = %d, want 0", eq)
	}

	strObj := periodToString([]any{p1}).(*object.Object)
	if str := object.GoStringFromStringObject(strObj); str != "P1Y2M3D" {
		t.Errorf("p1.toString() = %q, want 'P1Y2M3D'", str)
	}
}
