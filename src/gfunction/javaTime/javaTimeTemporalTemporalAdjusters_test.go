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
	"testing"
)

func temporalAdjustersTestInit() {
	globals.InitGlobals("test")
	temporalAdjustersClinit(nil)
	Load_Time_DayOfWeek()
	Load_Time_LocalDate()
	Load_Time_Temporal_TemporalAdjusters()
	Load_Time_Temporal_TemporalAdjuster()
}

func TestTemporalAdjustersClinitAndFactories(t *testing.T) {
	temporalAdjustersTestInit()

	expectedSigs := []string{
		"java/time/temporal/TemporalAdjusters.<clinit>()V",
		"java/time/temporal/TemporalAdjusters.<init>()V",
		"java/time/temporal/TemporalAdjusters.firstDayOfMonth()Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.lastDayOfMonth()Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.firstDayOfNextMonth()Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.firstDayOfYear()Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.lastDayOfYear()Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.firstDayOfNextYear()Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.firstInMonth(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.lastInMonth(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.dayOfWeekInMonth(ILjava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.next(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.nextOrSame(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.previous(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
		"java/time/temporal/TemporalAdjusters.previousOrSame(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;",
	}

	for _, sig := range expectedSigs {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("expected signature %s to be registered", sig)
		}
	}

	// Test <clinit> and <init>
	if res := temporalAdjustersClinit(nil); res != nil {
		t.Errorf("expected nil from clinit, got %v", res)
	}
	if res := temporalAdjustersInit(nil); res != nil {
		t.Errorf("expected nil from init, got %v", res)
	}

	// Test singleton factories
	factories := []struct {
		name string
		fn   func([]any) any
		want string
	}{
		{"firstDayOfMonth", temporalAdjustersFirstDayOfMonth, "firstDayOfMonth"},
		{"lastDayOfMonth", temporalAdjustersLastDayOfMonth, "lastDayOfMonth"},
		{"firstDayOfNextMonth", temporalAdjustersFirstDayOfNextMonth, "firstDayOfNextMonth"},
		{"firstDayOfYear", temporalAdjustersFirstDayOfYear, "firstDayOfYear"},
		{"lastDayOfYear", temporalAdjustersLastDayOfYear, "lastDayOfYear"},
		{"firstDayOfNextYear", temporalAdjustersFirstDayOfNextYear, "firstDayOfNextYear"},
	}

	for _, f := range factories {
		res := f.fn(nil)
		obj, ok := res.(*object.Object)
		if !ok || obj == nil {
			t.Fatalf("%s returned non-object %v", f.name, res)
		}
		adjType := getAdjusterType(obj)
		if adjType != f.want {
			t.Errorf("%s got adjusterType %q, want %q", f.name, adjType, f.want)
		}
	}

	// DayOfWeek parameterized factories
	dowMonday := dayOfWeekOf([]any{int64(1)}).(*object.Object)
	if dowMonday == nil {
		t.Fatalf("could not get DayOfWeek.MONDAY")
	}

	dowFactories := []struct {
		name string
		fn   func([]any) any
		want string
	}{
		{"firstInMonth", temporalAdjustersFirstInMonth, "dayOfWeekInMonth"},
		{"lastInMonth", temporalAdjustersLastInMonth, "dayOfWeekInMonth"},
		{"next", temporalAdjustersNext, "next"},
		{"nextOrSame", temporalAdjustersNextOrSame, "nextOrSame"},
		{"previous", temporalAdjustersPrevious, "previous"},
		{"previousOrSame", temporalAdjustersPreviousOrSame, "previousOrSame"},
	}

	for _, df := range dowFactories {
		// Null arg test
		err := df.fn([]any{nil})
		if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
			t.Errorf("%s with nil expected NullPointerException, got %v", df.name, err)
		}

		// Valid arg test
		res := df.fn([]any{dowMonday})
		obj, ok := res.(*object.Object)
		if !ok || obj == nil {
			t.Fatalf("%s returned non-object %v", df.name, res)
		}
		if adjType := getAdjusterType(obj); adjType != df.want {
			t.Errorf("%s got adjusterType %q, want %q", df.name, adjType, df.want)
		}
	}

	// dayOfWeekInMonth
	err := temporalAdjustersDayOfWeekInMonth([]any{int64(2), nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("dayOfWeekInMonth with nil DOW expected NullPointerException, got %v", err)
	}

	res := temporalAdjustersDayOfWeekInMonth([]any{int64(2), dowMonday})
	obj, ok := res.(*object.Object)
	if !ok || obj == nil {
		t.Fatalf("dayOfWeekInMonth returned %v", res)
	}
	if adjType := getAdjusterType(obj); adjType != "dayOfWeekInMonth" {
		t.Errorf("dayOfWeekInMonth got adjusterType %q, want 'dayOfWeekInMonth'", adjType)
	}

	// ofDateAdjuster
	res = temporalAdjustersOfDateAdjuster([]any{obj})
	dateAdjObj, ok := res.(*object.Object)
	if !ok || dateAdjObj == nil {
		t.Fatalf("ofDateAdjuster returned %v", res)
	}
	if adjType := getAdjusterType(dateAdjObj); adjType != "ofDateAdjuster" {
		t.Errorf("ofDateAdjuster got adjusterType %q, want 'ofDateAdjuster'", adjType)
	}
	err = temporalAdjustersOfDateAdjuster([]any{nil})
	if gerr, ok := err.(*ghelpers.GErrBlk); !ok || gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("ofDateAdjuster with nil expected NullPointerException, got %v", err)
	}
}

func TestTemporalAdjustersExecutionWithLocalDate(t *testing.T) {
	temporalAdjustersTestInit()

	// 2026-11-08 is a Sunday (DayOfWeek = 7)
	date := CreateLocalDate(2026, 11, 8)

	// firstDayOfMonth -> 2026-11-01
	adjFirstDOM := temporalAdjustersFirstDayOfMonth(nil).(*object.Object)
	res := temporalAdjusterAdjustInto([]any{adjFirstDOM, date})
	y, m, d, ok := getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 1 {
		t.Errorf("firstDayOfMonth on 2026-11-08 = %d-%d-%d, want 2026-11-1", y, m, d)
	}

	// lastDayOfMonth -> 2026-11-30
	adjLastDOM := temporalAdjustersLastDayOfMonth(nil).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjLastDOM, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 30 {
		t.Errorf("lastDayOfMonth on 2026-11-08 = %d-%d-%d, want 2026-11-30", y, m, d)
	}

	// firstDayOfNextMonth -> 2026-12-01
	adjFirstNextMonth := temporalAdjustersFirstDayOfNextMonth(nil).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjFirstNextMonth, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 12 || d != 1 {
		t.Errorf("firstDayOfNextMonth on 2026-11-08 = %d-%d-%d, want 2026-12-1", y, m, d)
	}

	// firstDayOfYear -> 2026-01-01
	adjFirstDOY := temporalAdjustersFirstDayOfYear(nil).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjFirstDOY, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 1 || d != 1 {
		t.Errorf("firstDayOfYear on 2026-11-08 = %d-%d-%d, want 2026-1-1", y, m, d)
	}

	// lastDayOfYear -> 2026-12-31
	adjLastDOY := temporalAdjustersLastDayOfYear(nil).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjLastDOY, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 12 || d != 31 {
		t.Errorf("lastDayOfYear on 2026-11-08 = %d-%d-%d, want 2026-12-31", y, m, d)
	}

	// firstDayOfNextYear -> 2027-01-01
	adjFirstNextYear := temporalAdjustersFirstDayOfNextYear(nil).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjFirstNextYear, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2027 || m != 1 || d != 1 {
		t.Errorf("firstDayOfNextYear on 2026-11-08 = %d-%d-%d, want 2027-1-1", y, m, d)
	}

	// next MONDAY after 2026-11-08 (Sun) -> 2026-11-09
	dowMonday := dayOfWeekOf([]any{int64(1)}).(*object.Object)
	adjNextMon := temporalAdjustersNext([]any{dowMonday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjNextMon, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 9 {
		t.Errorf("next MONDAY on 2026-11-08 = %d-%d-%d, want 2026-11-09", y, m, d)
	}

	// nextOrSame SUNDAY on 2026-11-08 (Sun) -> 2026-11-08
	dowSunday := dayOfWeekOf([]any{int64(7)}).(*object.Object)
	adjNextOrSameSun := temporalAdjustersNextOrSame([]any{dowSunday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjNextOrSameSun, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 8 {
		t.Errorf("nextOrSame SUNDAY on 2026-11-08 = %d-%d-%d, want 2026-11-08", y, m, d)
	}

	// next SUNDAY after 2026-11-08 (Sun) -> 2026-11-15
	adjNextSun := temporalAdjustersNext([]any{dowSunday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjNextSun, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 15 {
		t.Errorf("next SUNDAY on 2026-11-08 = %d-%d-%d, want 2026-11-15", y, m, d)
	}

	// previous MONDAY before 2026-11-08 (Sun) -> 2026-11-02
	adjPrevMon := temporalAdjustersPrevious([]any{dowMonday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjPrevMon, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 2 {
		t.Errorf("previous MONDAY on 2026-11-08 = %d-%d-%d, want 2026-11-02", y, m, d)
	}

	// previousOrSame SUNDAY on 2026-11-08 (Sun) -> 2026-11-08
	adjPrevOrSameSun := temporalAdjustersPreviousOrSame([]any{dowSunday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjPrevOrSameSun, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 8 {
		t.Errorf("previousOrSame SUNDAY on 2026-11-08 = %d-%d-%d, want 2026-11-08", y, m, d)
	}

	// firstInMonth MONDAY in Nov 2026 -> 2026-11-02
	adjFirstInMonthMon := temporalAdjustersFirstInMonth([]any{dowMonday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjFirstInMonthMon, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 2 {
		t.Errorf("firstInMonth MONDAY in Nov 2026 = %d-%d-%d, want 2026-11-02", y, m, d)
	}

	// lastInMonth MONDAY in Nov 2026 -> 2026-11-30
	adjLastInMonthMon := temporalAdjustersLastInMonth([]any{dowMonday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjLastInMonthMon, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 30 {
		t.Errorf("lastInMonth MONDAY in Nov 2026 = %d-%d-%d, want 2026-11-30", y, m, d)
	}

	// dayOfWeekInMonth (2nd Monday in Nov 2026) -> 2026-11-09
	adj2ndMon := temporalAdjustersDayOfWeekInMonth([]any{int64(2), dowMonday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adj2ndMon, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 9 {
		t.Errorf("2nd Monday in Nov 2026 = %d-%d-%d, want 2026-11-09", y, m, d)
	}

	// dayOfWeekInMonth (-1st / last Monday in Nov 2026) -> 2026-11-30
	adjLastMon := temporalAdjustersDayOfWeekInMonth([]any{int64(-1), dowMonday}).(*object.Object)
	res = temporalAdjusterAdjustInto([]any{adjLastMon, date})
	y, m, d, ok = getLocalDateYMD(res.(*object.Object))
	if !ok || y != 2026 || m != 11 || d != 30 {
		t.Errorf("-1 Monday in Nov 2026 = %d-%d-%d, want 2026-11-30", y, m, d)
	}
}
