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

func monthTestInit() {
	globals.InitGlobals("test")
	monthMutex.Lock()
	monthInitialized = false
	monthMutex.Unlock()
	Load_Time_Month()
	Load_Time_LocalDate()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
}

func TestMonthClinitAndConstants(t *testing.T) {
	monthTestInit()

	if res := monthClinit(nil); res != nil {
		t.Errorf("monthClinit returned %v", res)
	}

	expectedMonths := []string{
		"JANUARY", "FEBRUARY", "MARCH", "APRIL", "MAY", "JUNE",
		"JULY", "AUGUST", "SEPTEMBER", "OCTOBER", "NOVEMBER", "DECEMBER",
	}

	for i, mName := range expectedMonths {
		st, ok := statics.QueryStatic(classNameMonthVal, mName)
		if !ok || st.Value == nil {
			t.Fatalf("static Month.%s not found", mName)
		}
		mObj := st.Value.(*object.Object)
		val := monthGetValue([]any{mObj}).(int64)
		if val != int64(i+1) {
			t.Errorf("Month.%s.getValue() = %d, want %d", mName, val, i+1)
		}
		ord := monthOrdinal([]any{mObj}).(int64)
		if ord != int64(i) {
			t.Errorf("Month.%s.ordinal() = %d, want %d", mName, ord, i)
		}
		nameObj := monthName([]any{mObj}).(*object.Object)
		if nameStr := object.GoStringFromStringObject(nameObj); nameStr != mName {
			t.Errorf("Month.%s.name() = %q, want %q", mName, nameStr, mName)
		}
	}
}

func TestMonthFactories(t *testing.T) {
	monthTestInit()

	// of(1..12)
	for i := int64(1); i <= 12; i++ {
		res := monthOf([]any{i})
		mObj, ok := res.(*object.Object)
		if !ok || mObj == nil {
			t.Fatalf("Month.of(%d) failed: %v", i, res)
		}
		if val := monthGetValue([]any{mObj}).(int64); val != i {
			t.Errorf("Month.of(%d).getValue() = %d", i, val)
		}
	}

	// of invalid
	for _, badVal := range []int64{0, 13, -1, 100} {
		err := monthOf([]any{badVal})
		if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeException {
			t.Errorf("Month.of(%d) expected DateTimeException, got %v", badVal, err)
		}
	}

	// valueOf
	janObj := object.StringObjectFromGoString("JANUARY")
	res := monthValueOf([]any{janObj})
	if mObj, ok := res.(*object.Object); !ok || monthGetValue([]any{mObj}).(int64) != int64(1) {
		t.Errorf("Month.valueOf(JANUARY) failed: %v", res)
	}

	badStr := object.StringObjectFromGoString("INVALID_MONTH")
	err := monthValueOf([]any{badStr})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Month.valueOf(INVALID) expected IllegalArgumentException, got %v", err)
	}

	// values
	valsRes := monthValues(nil)
	arrObj, ok := valsRes.(*object.Object)
	if !ok || arrObj == nil {
		t.Fatalf("Month.values() returned non-array: %v", valsRes)
	}
	elems := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(elems) != 12 {
		t.Errorf("Month.values() length = %d, want 12", len(elems))
	}

	// from
	dateNov := CreateLocalDate(2026, 11, 8)
	fromRes := monthFrom([]any{dateNov})
	if mObj, ok := fromRes.(*object.Object); !ok || monthGetValue([]any{mObj}).(int64) != int64(11) {
		t.Errorf("Month.from(2026-11-08) = %v, want NOVEMBER(11)", fromRes)
	}
}

func TestMonthCalculationsAndOperations(t *testing.T) {
	monthTestInit()

	feb := monthOf([]any{int64(2)}).(*object.Object) // FEBRUARY
	nov := monthOf([]any{int64(11)}).(*object.Object) // NOVEMBER

	// length (leap vs non-leap)
	if l := monthLength([]any{feb, int64(0)}).(int64); l != 28 {
		t.Errorf("FEBRUARY.length(false) = %d, want 28", l)
	}
	if l := monthLength([]any{feb, int64(1)}).(int64); l != 29 {
		t.Errorf("FEBRUARY.length(true) = %d, want 29", l)
	}
	if ml := monthMinLength([]any{feb}).(int64); ml != 28 {
		t.Errorf("FEBRUARY.minLength() = %d, want 28", ml)
	}
	if ml := monthMaxLength([]any{feb}).(int64); ml != 29 {
		t.Errorf("FEBRUARY.maxLength() = %d, want 29", ml)
	}

	// firstDayOfYear
	if f := monthFirstDayOfYear([]any{feb, int64(0)}).(int64); f != 32 {
		t.Errorf("FEBRUARY.firstDayOfYear(false) = %d, want 32", f)
	}
	if f := monthFirstDayOfYear([]any{feb, int64(1)}).(int64); f != 32 {
		t.Errorf("FEBRUARY.firstDayOfYear(true) = %d, want 32", f)
	}

	mar := monthOf([]any{int64(3)}).(*object.Object)
	if f := monthFirstDayOfYear([]any{mar, int64(0)}).(int64); f != 60 {
		t.Errorf("MARCH.firstDayOfYear(false) = %d, want 60", f)
	}
	if f := monthFirstDayOfYear([]any{mar, int64(1)}).(int64); f != 61 {
		t.Errorf("MARCH.firstDayOfYear(true) = %d, want 61", f)
	}

	// firstMonthOfQuarter
	// NOV is in Q4 -> OCT (10)
	qMonth := monthFirstMonthOfQuarter([]any{nov}).(*object.Object)
	if val := monthGetValue([]any{qMonth}).(int64); val != 10 {
		t.Errorf("NOVEMBER.firstMonthOfQuarter() = %d, want 10 (OCTOBER)", val)
	}

	// plus & minus
	plus2 := monthPlus([]any{nov, int64(2)}).(*object.Object) // JAN (1)
	if val := monthGetValue([]any{plus2}).(int64); val != 1 {
		t.Errorf("NOVEMBER.plus(2) = %d, want 1 (JANUARY)", val)
	}
	minus1 := monthMinus([]any{nov, int64(1)}).(*object.Object) // OCT (10)
	if val := monthGetValue([]any{minus1}).(int64); val != 10 {
		t.Errorf("NOVEMBER.minus(1) = %d, want 10 (OCTOBER)", val)
	}

	// compareTo & equals
	if cmp := monthCompareTo([]any{feb, nov}).(int64); cmp >= 0 {
		t.Errorf("FEBRUARY.compareTo(NOVEMBER) = %d, want < 0", cmp)
	}
	if eq := monthEquals([]any{feb, feb}).(int64); eq != 1 {
		t.Errorf("FEBRUARY.equals(FEBRUARY) = %d, want 1", eq)
	}
	if eq := monthEquals([]any{feb, nov}).(int64); eq != 0 {
		t.Errorf("FEBRUARY.equals(NOVEMBER) = %d, want 0", eq)
	}

	// toString
	strObj := monthToString([]any{nov}).(*object.Object)
	if str := object.GoStringFromStringObject(strObj); str != "NOVEMBER" {
		t.Errorf("NOVEMBER.toString() = %q, want 'NOVEMBER'", str)
	}

	// isSupported & get & getLong
	cfMOY := getChronoField("MONTH_OF_YEAR")
	if sup := monthIsSupported([]any{nov, cfMOY}).(int64); sup != 1 {
		t.Errorf("NOVEMBER.isSupported(MONTH_OF_YEAR) = %d, want 1", sup)
	}
	if val := monthGet([]any{nov, cfMOY}).(int64); val != 11 {
		t.Errorf("NOVEMBER.get(MONTH_OF_YEAR) = %d, want 11", val)
	}
	if val := monthGetLong([]any{nov, cfMOY}).(int64); val != 11 {
		t.Errorf("NOVEMBER.getLong(MONTH_OF_YEAR) = %d, want 11", val)
	}
}
