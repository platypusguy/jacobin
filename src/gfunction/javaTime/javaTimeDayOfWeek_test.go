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

func dayOfWeekTestInit() {
	globals.InitGlobals("test")
	dayOfWeekMutex.Lock()
	dayOfWeekInitialized = false
	dayOfWeekMutex.Unlock()
	Load_Time_DayOfWeek()
	Load_Time_LocalDate()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_ChronoUnit()
}

func TestDayOfWeekClinitAndConstants(t *testing.T) {
	dayOfWeekTestInit()

	if res := dayOfWeekClinit(nil); res != nil {
		t.Errorf("dayOfWeekClinit returned %v", res)
	}

	expectedDays := []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}
	for i, dayName := range expectedDays {
		st, ok := statics.QueryStatic(classNameDayOfWeekVal, dayName)
		if !ok || st.Value == nil {
			t.Fatalf("static DayOfWeek.%s not found", dayName)
		}
		dayObj := st.Value.(*object.Object)
		val := dayOfWeekGetValue([]any{dayObj}).(int64)
		if val != int64(i+1) {
			t.Errorf("DayOfWeek.%s.getValue() = %d, want %d", dayName, val, i+1)
		}
		ord := dayOfWeekOrdinal([]any{dayObj}).(int64)
		if ord != int64(i) {
			t.Errorf("DayOfWeek.%s.ordinal() = %d, want %d", dayName, ord, i)
		}
		nameObj := dayOfWeekName([]any{dayObj}).(*object.Object)
		if nameStr := object.GoStringFromStringObject(nameObj); nameStr != dayName {
			t.Errorf("DayOfWeek.%s.name() = %q, want %q", dayName, nameStr, dayName)
		}
	}
}

func TestDayOfWeekFactories(t *testing.T) {
	dayOfWeekTestInit()

	// of(1..7)
	for i := int64(1); i <= 7; i++ {
		res := dayOfWeekOf([]any{i})
		dayObj, ok := res.(*object.Object)
		if !ok || dayObj == nil {
			t.Fatalf("DayOfWeek.of(%d) failed: %v", i, res)
		}
		val := dayOfWeekGetValue([]any{dayObj}).(int64)
		if val != i {
			t.Errorf("DayOfWeek.of(%d).getValue() = %d", i, val)
		}
	}

	// of invalid
	for _, badVal := range []int64{0, 8, -1, 100} {
		err := dayOfWeekOf([]any{badVal})
		if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.DateTimeException {
			t.Errorf("DayOfWeek.of(%d) expected DateTimeException, got %v", badVal, err)
		}
	}

	// valueOf
	monObj := object.StringObjectFromGoString("MONDAY")
	res := dayOfWeekValueOf([]any{monObj})
	dayObj, ok := res.(*object.Object)
	if !ok || dayOfWeekGetValue([]any{dayObj}).(int64) != int64(1) {
		t.Errorf("DayOfWeek.valueOf(MONDAY) failed: %v", res)
	}

	badStr := object.StringObjectFromGoString("INVALID_DAY")
	err := dayOfWeekValueOf([]any{badStr})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("DayOfWeek.valueOf(INVALID) expected IllegalArgumentException, got %v", err)
	}

	// values
	valsRes := dayOfWeekValues(nil)
	arrObj, ok := valsRes.(*object.Object)
	if !ok || arrObj == nil {
		t.Fatalf("DayOfWeek.values() returned non-array: %v", valsRes)
	}
	elems := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(elems) != 7 {
		t.Errorf("DayOfWeek.values() length = %d, want 7", len(elems))
	}

	// from
	// 2026-11-08 is Sunday (7)
	dateSun := CreateLocalDate(2026, 11, 8)
	fromRes := dayOfWeekFrom([]any{dateSun})
	dayObj, ok = fromRes.(*object.Object)
	if !ok || dayOfWeekGetValue([]any{dayObj}).(int64) != int64(7) {
		t.Errorf("DayOfWeek.from(2026-11-08) = %v, want SUNDAY(7)", fromRes)
	}
}

func TestDayOfWeekOperations(t *testing.T) {
	dayOfWeekTestInit()

	monday := dayOfWeekOf([]any{int64(1)}).(*object.Object)
	sunday := dayOfWeekOf([]any{int64(7)}).(*object.Object)

	// plus
	plus3 := dayOfWeekPlus([]any{monday, int64(3)}).(*object.Object)
	if val := dayOfWeekGetValue([]any{plus3}).(int64); val != 4 { // Thursday
		t.Errorf("MONDAY.plus(3) = %d, want 4 (THURSDAY)", val)
	}
	plus7 := dayOfWeekPlus([]any{monday, int64(7)}).(*object.Object)
	if val := dayOfWeekGetValue([]any{plus7}).(int64); val != 1 {
		t.Errorf("MONDAY.plus(7) = %d, want 1 (MONDAY)", val)
	}

	// minus
	minus1 := dayOfWeekMinus([]any{monday, int64(1)}).(*object.Object)
	if val := dayOfWeekGetValue([]any{minus1}).(int64); val != 7 { // Sunday
		t.Errorf("MONDAY.minus(1) = %d, want 7 (SUNDAY)", val)
	}

	// compareTo & equals
	if cmp := dayOfWeekCompareTo([]any{monday, sunday}).(int64); cmp >= 0 {
		t.Errorf("MONDAY.compareTo(SUNDAY) = %d, want < 0", cmp)
	}
	if cmp := dayOfWeekCompareTo([]any{monday, monday}).(int64); cmp != 0 {
		t.Errorf("MONDAY.compareTo(MONDAY) = %d, want 0", cmp)
	}
	if eq := dayOfWeekEquals([]any{monday, monday}).(int64); eq != 1 {
		t.Errorf("MONDAY.equals(MONDAY) = %d, want 1", eq)
	}
	if eq := dayOfWeekEquals([]any{monday, sunday}).(int64); eq != 0 {
		t.Errorf("MONDAY.equals(SUNDAY) = %d, want 0", eq)
	}

	// toString
	strObj := dayOfWeekToString([]any{monday}).(*object.Object)
	if str := object.GoStringFromStringObject(strObj); str != "MONDAY" {
		t.Errorf("MONDAY.toString() = %q, want 'MONDAY'", str)
	}

	// isSupported, range, get, getLong
	cfDOW := getChronoField("DAY_OF_WEEK")
	if sup := dayOfWeekIsSupported([]any{monday, cfDOW}).(int64); sup != 1 {
		t.Errorf("MONDAY.isSupported(DAY_OF_WEEK) = %d, want 1", sup)
	}
	cfDOM := getChronoField("DAY_OF_MONTH")
	if sup := dayOfWeekIsSupported([]any{monday, cfDOM}).(int64); sup != 0 {
		t.Errorf("MONDAY.isSupported(DAY_OF_MONTH) = %d, want 0", sup)
	}

	val := dayOfWeekGet([]any{monday, cfDOW}).(int64)
	if val != 1 {
		t.Errorf("MONDAY.get(DAY_OF_WEEK) = %d, want 1", val)
	}
	longVal := dayOfWeekGetLong([]any{monday, cfDOW}).(int64)
	if longVal != 1 {
		t.Errorf("MONDAY.getLong(DAY_OF_WEEK) = %d, want 1", longVal)
	}

	rng := dayOfWeekRange([]any{monday, cfDOW})
	if rngObj, ok := rng.(*object.Object); !ok || rngObj == nil {
		t.Errorf("MONDAY.range(DAY_OF_WEEK) failed: %v", rng)
	}
}
