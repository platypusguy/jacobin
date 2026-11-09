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

func zoneOffsetTestInit() {
	globals.InitGlobals("test")
	Load_Time_ZoneOffset()
}

func TestZoneOffsetClinitAndConstants(t *testing.T) {
	zoneOffsetTestInit()
	ensureZoneOffsetInited()

	utc := statics.GetStaticValue(classNameZoneOffsetVal, "UTC")
	if utc == nil {
		t.Fatal("ZoneOffset.UTC static is nil")
	}
	utcObj := utc.(*object.Object)
	sec, ok := getZoneOffsetTotalSecondsVal(utcObj)
	if !ok || sec != 0 {
		t.Fatalf("ZoneOffset.UTC totalSeconds = %d, want 0", sec)
	}
	id := getZoneOffsetID(utcObj)
	if id != "Z" {
		t.Fatalf("ZoneOffset.UTC id = %q, want 'Z'", id)
	}

	minVal := statics.GetStaticValue(classNameZoneOffsetVal, "MIN")
	if minVal == nil {
		t.Fatal("ZoneOffset.MIN static is nil")
	}
	minObj := minVal.(*object.Object)
	sec, ok = getZoneOffsetTotalSecondsVal(minObj)
	if !ok || sec != -64800 {
		t.Fatalf("ZoneOffset.MIN totalSeconds = %d, want -64800", sec)
	}

	maxVal := statics.GetStaticValue(classNameZoneOffsetVal, "MAX")
	if maxVal == nil {
		t.Fatal("ZoneOffset.MAX static is nil")
	}
	maxObj := maxVal.(*object.Object)
	sec, ok = getZoneOffsetTotalSecondsVal(maxObj)
	if !ok || sec != 64800 {
		t.Fatalf("ZoneOffset.MAX totalSeconds = %d, want 64800", sec)
	}
}

func TestZoneOffsetOfFactories(t *testing.T) {
	zoneOffsetTestInit()

	// of(String)
	tests := []struct {
		input   string
		wantSec int64
		wantID  string
		wantErr bool
	}{
		{"Z", 0, "Z", false},
		{"+0", 0, "Z", false},
		{"+5", 18000, "+05:00", false},
		{"-5", -18000, "-05:00", false},
		{"+05", 18000, "+05:00", false},
		{"-05", -18000, "-05:00", false},
		{"+0530", 19800, "+05:30", false},
		{"-0530", -19800, "-05:30", false},
		{"+05:30", 19800, "+05:30", false},
		{"-05:30", -19800, "-05:30", false},
		{"+053045", 19845, "+05:30:45", false},
		{"-053045", -19845, "-05:30:45", false},
		{"+05:30:45", 19845, "+05:30:45", false},
		{"-05:30:45", -19845, "-05:30:45", false},
		{"+18:00", 64800, "+18:00", false},
		{"-18:00", -64800, "-18:00", false},
		{"+18:01", 0, "", true},
		{"-18:01", 0, "", true},
		{"+19:00", 0, "", true},
		{"invalid", 0, "", true},
		{"", 0, "", true},
		{"+1", 3600, "+01:00", false},
	}

	for _, tt := range tests {
		res := zoneOffsetOf([]any{object.StringObjectFromGoString(tt.input)})
		if tt.wantErr {
			if _, isErr := res.(*ghelpers.GErrBlk); !isErr {
				t.Errorf("zoneOffsetOf(%q) expected error, got %v", tt.input, res)
			}
		} else {
			if err, isErr := res.(*ghelpers.GErrBlk); isErr {
				t.Errorf("zoneOffsetOf(%q) unexpected error: %v", tt.input, err.ErrMsg)
				continue
			}
			obj := res.(*object.Object)
			sec, _ := getZoneOffsetTotalSecondsVal(obj)
			if sec != tt.wantSec {
				t.Errorf("zoneOffsetOf(%q) seconds = %d, want %d", tt.input, sec, tt.wantSec)
			}
			id := getZoneOffsetID(obj)
			if id != tt.wantID {
				t.Errorf("zoneOffsetOf(%q) id = %q, want %q", tt.input, id, tt.wantID)
			}
		}
	}

	// ofHours
	resH := zoneOffsetOfHours([]any{int64(5)})
	if sec, _ := getZoneOffsetTotalSecondsVal(resH.(*object.Object)); sec != 18000 {
		t.Fatalf("zoneOffsetOfHours(5) = %d, want 18000", sec)
	}
	resHErr := zoneOffsetOfHours([]any{int64(19)})
	if _, isErr := resHErr.(*ghelpers.GErrBlk); !isErr {
		t.Fatal("zoneOffsetOfHours(19) expected error")
	}

	// ofHoursMinutes
	resHM := zoneOffsetOfHoursMinutes([]any{int64(-5), int64(-30)})
	if sec, _ := getZoneOffsetTotalSecondsVal(resHM.(*object.Object)); sec != -19800 {
		t.Fatalf("zoneOffsetOfHoursMinutes(-5, -30) = %d, want -19800", sec)
	}
	resHMErr := zoneOffsetOfHoursMinutes([]any{int64(-5), int64(30)})
	if _, isErr := resHMErr.(*ghelpers.GErrBlk); !isErr {
		t.Fatal("zoneOffsetOfHoursMinutes(-5, 30) expected error for sign mismatch")
	}

	// ofHoursMinutesSeconds
	resHMS := zoneOffsetOfHoursMinutesSeconds([]any{int64(2), int64(15), int64(30)})
	if sec, _ := getZoneOffsetTotalSecondsVal(resHMS.(*object.Object)); sec != 8130 {
		t.Fatalf("zoneOffsetOfHoursMinutesSeconds(2, 15, 30) = %d, want 8130", sec)
	}

	// ofTotalSeconds
	resTS := zoneOffsetOfTotalSeconds([]any{int64(3600)})
	if sec, _ := getZoneOffsetTotalSecondsVal(resTS.(*object.Object)); sec != 3600 {
		t.Fatalf("zoneOffsetOfTotalSeconds(3600) = %d, want 3600", sec)
	}
	resTSErr := zoneOffsetOfTotalSeconds([]any{int64(64801)})
	if _, isErr := resTSErr.(*ghelpers.GErrBlk); !isErr {
		t.Fatal("zoneOffsetOfTotalSeconds(64801) expected error")
	}
}

func TestZoneOffsetGettersAndMethods(t *testing.T) {
	zoneOffsetTestInit()

	offset := createZoneOffset(7200) // +02:00
	if offset == nil {
		t.Fatal("createZoneOffset(7200) returned nil")
	}

	// getTotalSeconds
	totSec := zoneOffsetGetTotalSeconds([]any{offset})
	if totSec != int64(7200) {
		t.Fatalf("getTotalSeconds() = %v, want 7200", totSec)
	}

	// getId
	idRes := zoneOffsetGetId([]any{offset})
	idStr := object.GoStringFromStringObject(idRes.(*object.Object))
	if idStr != "+02:00" {
		t.Fatalf("getId() = %q, want '+02:00'", idStr)
	}

	// toString
	toStrRes := zoneOffsetToString([]any{offset})
	toStr := object.GoStringFromStringObject(toStrRes.(*object.Object))
	if toStr != "+02:00" {
		t.Fatalf("toString() = %q, want '+02:00'", toStr)
	}

	// normalized
	normRes := zoneOffsetNormalized([]any{offset})
	if normRes != offset {
		t.Fatalf("normalized() = %v, want self", normRes)
	}

	// getRules
	rulesRes := zoneOffsetGetRules([]any{offset})
	rulesObj, ok := rulesRes.(*object.Object)
	if !ok || rulesObj == nil {
		t.Fatal("getRules() returned nil or non-object")
	}
	if !isZoneRulesFixed(rulesObj) {
		t.Fatal("rules.isFixedOffset() expected true")
	}
}

func TestZoneOffsetTemporalAccessor(t *testing.T) {
	zoneOffsetTestInit()
	Load_Time_Temporal_ChronoField()
	Load_Time_Temporal_TemporalQueries()

	offset := createZoneOffset(-18000) // -05:00
	cfOffsetSeconds := getChronoField("OFFSET_SECONDS")
	cfNanoOfSecond := getChronoField("NANO_OF_SECOND")

	// isSupported
	if zoneOffsetIsSupported([]any{offset, cfOffsetSeconds}) != types.JavaBoolTrue {
		t.Fatal("isSupported(OFFSET_SECONDS) = false, want true")
	}
	if zoneOffsetIsSupported([]any{offset, cfNanoOfSecond}) != types.JavaBoolFalse {
		t.Fatal("isSupported(NANO_OF_SECOND) = true, want false")
	}

	// range
	r := zoneOffsetRange([]any{offset, cfOffsetSeconds}).(*object.Object)
	minVal, _ := getInt64Val(r.FieldTable["minSmallest"].Fvalue)
	maxVal, _ := getInt64Val(r.FieldTable["maxLargest"].Fvalue)
	if minVal != -64800 || maxVal != 64800 {
		t.Fatalf("range(OFFSET_SECONDS) = %d..%d, want -64800..64800", minVal, maxVal)
	}

	// get / getLong
	getVal := zoneOffsetGet([]any{offset, cfOffsetSeconds})
	if getVal != int64(-18000) {
		t.Fatalf("get(OFFSET_SECONDS) = %v, want -18000", getVal)
	}
	getLongVal := zoneOffsetGetLong([]any{offset, cfOffsetSeconds})
	if getLongVal != int64(-18000) {
		t.Fatalf("getLong(OFFSET_SECONDS) = %v, want -18000", getLongVal)
	}

	// query
	qOffset := temporalQueriesOffset(nil).(*object.Object)
	qZone := temporalQueriesZone(nil).(*object.Object)
	qZoneId := temporalQueriesZoneId(nil).(*object.Object)
	qPrec := temporalQueriesPrecision(nil).(*object.Object)

	if zoneOffsetQuery([]any{offset, qOffset}) != offset {
		t.Fatal("query(offset) did not return offset")
	}
	if zoneOffsetQuery([]any{offset, qZone}) != offset {
		t.Fatal("query(zone) did not return offset")
	}
	if zoneOffsetQuery([]any{offset, qZoneId}) != nil {
		t.Fatal("query(zoneId) did not return nil")
	}
	if zoneOffsetQuery([]any{offset, qPrec}) != nil {
		t.Fatal("query(precision) did not return nil")
	}

	// from
	fromRes := zoneOffsetFrom([]any{offset})
	if fromRes != offset {
		t.Fatal("from(ZoneOffset) did not return self")
	}
}

func TestZoneOffsetCompareAndEquals(t *testing.T) {
	zoneOffsetTestInit()

	off1 := createZoneOffset(7200)  // +02:00
	off2 := createZoneOffset(3600)  // +01:00
	off3 := createZoneOffset(7200)  // +02:00
	offNeg := createZoneOffset(-3600) // -01:00

	// In java.time.ZoneOffset, +02:00 is earlier on timeline than +01:00,
	// so compareTo returns negative when comparing +02:00 with +01:00.
	cmp := zoneOffsetCompareTo([]any{off1, off2}).(int64)
	if cmp >= 0 {
		t.Fatalf("compareTo(+02:00, +01:00) = %d, want < 0", cmp)
	}

	cmpNeg := zoneOffsetCompareTo([]any{off2, off1}).(int64)
	if cmpNeg <= 0 {
		t.Fatalf("compareTo(+01:00, +02:00) = %d, want > 0", cmpNeg)
	}

	cmpEqual := zoneOffsetCompareTo([]any{off1, off3}).(int64)
	if cmpEqual != 0 {
		t.Fatalf("compareTo(+02:00, +02:00) = %d, want 0", cmpEqual)
	}

	// equals
	if zoneOffsetEquals([]any{off1, off3}) != types.JavaBoolTrue {
		t.Fatal("equals(+02:00, +02:00) = false, want true")
	}
	if zoneOffsetEquals([]any{off1, off2}) != types.JavaBoolFalse {
		t.Fatal("equals(+02:00, +01:00) = true, want false")
	}
	if zoneOffsetEquals([]any{off1, offNeg}) != types.JavaBoolFalse {
		t.Fatal("equals(+02:00, -01:00) = true, want false")
	}

	// hashCode
	h := zoneOffsetHashCode([]any{off1}).(int64)
	if h != 7200 {
		t.Fatalf("hashCode(+02:00) = %d, want 7200", h)
	}
}
