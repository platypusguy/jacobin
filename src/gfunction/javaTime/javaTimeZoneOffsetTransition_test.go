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
	"jacobin/src/types"
	"testing"
)

func zoneOffsetTransitionTestInit() {
	globals.InitGlobals("test")
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
	Load_Time_Zone_ZoneOffsetTransition()
}

func TestZoneOffsetTransitionRegistration(t *testing.T) {
	zoneOffsetTransitionTestInit()

	expectedSignatures := []string{
		"java/time/zone/ZoneOffsetTransition.<clinit>()V",
		"java/time/zone/ZoneOffsetTransition.<init>(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;)V",
		"java/time/zone/ZoneOffsetTransition.<init>(JLjava/time/ZoneOffset;Ljava/time/ZoneOffset;)V",
		"java/time/zone/ZoneOffsetTransition.of(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;)Ljava/time/zone/ZoneOffsetTransition;",
		"java/time/zone/ZoneOffsetTransition.getInstant()Ljava/time/Instant;",
		"java/time/zone/ZoneOffsetTransition.toEpochSecond()J",
		"java/time/zone/ZoneOffsetTransition.getDateTimeBefore()Ljava/time/LocalDateTime;",
		"java/time/zone/ZoneOffsetTransition.getDateTimeAfter()Ljava/time/LocalDateTime;",
		"java/time/zone/ZoneOffsetTransition.getOffsetBefore()Ljava/time/ZoneOffset;",
		"java/time/zone/ZoneOffsetTransition.getOffsetAfter()Ljava/time/ZoneOffset;",
		"java/time/zone/ZoneOffsetTransition.getDuration()Ljava/time/Duration;",
		"java/time/zone/ZoneOffsetTransition.isGap()Z",
		"java/time/zone/ZoneOffsetTransition.isOverlap()Z",
		"java/time/zone/ZoneOffsetTransition.isValidOffset(Ljava/time/ZoneOffset;)Z",
		"java/time/zone/ZoneOffsetTransition.compareTo(Ljava/time/zone/ZoneOffsetTransition;)I",
		"java/time/zone/ZoneOffsetTransition.compareTo(Ljava/lang/Object;)I",
		"java/time/zone/ZoneOffsetTransition.equals(Ljava/lang/Object;)Z",
		"java/time/zone/ZoneOffsetTransition.hashCode()I",
		"java/time/zone/ZoneOffsetTransition.toString()Ljava/lang/String;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}
}

func TestZoneOffsetTransitionOfAndValidation(t *testing.T) {
	zoneOffsetTransitionTestInit()

	d := CreateLocalDate(2020, 3, 8)
	tm := createLocalTime(2, 0, 0, 0)
	ldt := createLocalDateTime(d, tm)

	offBefore := createZoneOffset(-5 * 3600) // -05:00
	offAfter := createZoneOffset(-4 * 3600)  // -04:00

	// Valid creation
	res := zoneOffsetTransitionOf([]any{ldt, offBefore, offAfter})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Unexpected error: %v", err)
	}
	transObj, ok := res.(*object.Object)
	if !ok || transObj == nil {
		t.Fatalf("Expected object, got %T", res)
	}

	// Null transition
	resNullTrans := zoneOffsetTransitionOf([]any{object.Null, offBefore, offAfter})
	if err, isErr := resNullTrans.(*ghelpers.GErrBlk); !isErr || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE for null transition, got %v", resNullTrans)
	}

	// Null offsetBefore
	resNullBefore := zoneOffsetTransitionOf([]any{ldt, object.Null, offAfter})
	if err, isErr := resNullBefore.(*ghelpers.GErrBlk); !isErr || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE for null offsetBefore, got %v", resNullBefore)
	}

	// Null offsetAfter
	resNullAfter := zoneOffsetTransitionOf([]any{ldt, offBefore, object.Null})
	if err, isErr := resNullAfter.(*ghelpers.GErrBlk); !isErr || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE for null offsetAfter, got %v", resNullAfter)
	}

	// Equal offsets -> IllegalArgumentException
	resEqual := zoneOffsetTransitionOf([]any{ldt, offBefore, offBefore})
	if err, isErr := resEqual.(*ghelpers.GErrBlk); !isErr || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for equal offsets, got %v", resEqual)
	}

	// Nano of second not zero -> IllegalArgumentException
	tmWithNano := createLocalTime(2, 0, 0, 500)
	ldtWithNano := createLocalDateTime(d, tmWithNano)
	resNano := zoneOffsetTransitionOf([]any{ldtWithNano, offBefore, offAfter})
	if err, isErr := resNano.(*ghelpers.GErrBlk); !isErr || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for non-zero nano, got %v", resNano)
	}
}

func TestZoneOffsetTransitionGapOperations(t *testing.T) {
	zoneOffsetTransitionTestInit()

	// Spring forward: 2020-03-08 02:00 -05:00 -> -04:00 (Gap of 1 hour)
	d := CreateLocalDate(2020, 3, 8)
	tm := createLocalTime(2, 0, 0, 0)
	ldt := createLocalDateTime(d, tm)

	offBefore := createZoneOffset(-5 * 3600) // -05:00
	offAfter := createZoneOffset(-4 * 3600)  // -04:00

	trans := zoneOffsetTransitionOf([]any{ldt, offBefore, offAfter}).(*object.Object)

	// isGap and isOverlap
	isGap := zoneOffsetTransitionIsGap([]any{trans})
	if isGap != types.JavaBoolTrue {
		t.Errorf("Expected isGap=true, got %v", isGap)
	}
	isOverlap := zoneOffsetTransitionIsOverlap([]any{trans})
	if isOverlap != types.JavaBoolFalse {
		t.Errorf("Expected isOverlap=false, got %v", isOverlap)
	}

	// getDateTimeBefore
	dtBefore := zoneOffsetTransitionGetDateTimeBefore([]any{trans}).(*object.Object)
	if localDateTimeEquals([]any{dtBefore, ldt}) != types.JavaBoolTrue {
		t.Errorf("getDateTimeBefore mismatch")
	}

	// getDateTimeAfter (2020-03-08T03:00:00)
	dtAfter := zoneOffsetTransitionGetDateTimeAfter([]any{trans}).(*object.Object)
	expectedAfterLdt := createLocalDateTime(d, createLocalTime(3, 0, 0, 0))
	if localDateTimeEquals([]any{dtAfter, expectedAfterLdt}) != types.JavaBoolTrue {
		t.Errorf("getDateTimeAfter mismatch: expected 03:00")
	}

	// getOffsetBefore and getOffsetAfter
	b := zoneOffsetTransitionGetOffsetBefore([]any{trans}).(*object.Object)
	if zoneOffsetEquals([]any{b, offBefore}) != types.JavaBoolTrue {
		t.Errorf("getOffsetBefore mismatch")
	}
	a := zoneOffsetTransitionGetOffsetAfter([]any{trans}).(*object.Object)
	if zoneOffsetEquals([]any{a, offAfter}) != types.JavaBoolTrue {
		t.Errorf("getOffsetAfter mismatch")
	}

	// getDuration (3600 seconds)
	dur := zoneOffsetTransitionGetDuration([]any{trans}).(*object.Object)
	durSec := dur.FieldTable["seconds"].Fvalue.(int64)
	if durSec != 3600 {
		t.Errorf("Expected duration 3600s, got %d", durSec)
	}

	// isValidOffset: gap has no valid offsets
	if zoneOffsetTransitionIsValidOffset([]any{trans, offBefore}) != types.JavaBoolFalse {
		t.Errorf("Gap should not validate offsetBefore")
	}
	if zoneOffsetTransitionIsValidOffset([]any{trans, offAfter}) != types.JavaBoolFalse {
		t.Errorf("Gap should not validate offsetAfter")
	}

	// toString
	strObj := zoneOffsetTransitionToString([]any{trans}).(*object.Object)
	str := object.GoStringFromStringObject(strObj)
	expectedStr := "Transition[Gap at 2020-03-08T02:00-05:00 to -04:00]"
	if str != expectedStr {
		t.Errorf("toString mismatch:\ngot:  %s\nwant: %s", str, expectedStr)
	}
}

func TestZoneOffsetTransitionOverlapOperations(t *testing.T) {
	zoneOffsetTransitionTestInit()

	// Fall back: 2020-11-01 01:00 -04:00 -> -05:00 (Overlap of 1 hour)
	d := CreateLocalDate(2020, 11, 1)
	tm := createLocalTime(1, 0, 0, 0)
	ldt := createLocalDateTime(d, tm)

	offBefore := createZoneOffset(-4 * 3600) // -04:00
	offAfter := createZoneOffset(-5 * 3600)  // -05:00

	trans := zoneOffsetTransitionOf([]any{ldt, offBefore, offAfter}).(*object.Object)

	// isGap and isOverlap
	if zoneOffsetTransitionIsGap([]any{trans}) != types.JavaBoolFalse {
		t.Errorf("Expected isGap=false")
	}
	if zoneOffsetTransitionIsOverlap([]any{trans}) != types.JavaBoolTrue {
		t.Errorf("Expected isOverlap=true")
	}

	// getDateTimeAfter (2020-11-01T00:00:00)
	dtAfter := zoneOffsetTransitionGetDateTimeAfter([]any{trans}).(*object.Object)
	expectedAfterLdt := createLocalDateTime(d, createLocalTime(0, 0, 0, 0))
	if localDateTimeEquals([]any{dtAfter, expectedAfterLdt}) != types.JavaBoolTrue {
		t.Errorf("getDateTimeAfter mismatch for overlap: expected 00:00")
	}

	// getDuration (-3600 seconds)
	dur := zoneOffsetTransitionGetDuration([]any{trans}).(*object.Object)
	durSec := dur.FieldTable["seconds"].Fvalue.(int64)
	if durSec != -3600 {
		t.Errorf("Expected duration -3600s, got %d", durSec)
	}

	// isValidOffset: overlap validates both before and after
	if zoneOffsetTransitionIsValidOffset([]any{trans, offBefore}) != types.JavaBoolTrue {
		t.Errorf("Overlap should validate offsetBefore")
	}
	if zoneOffsetTransitionIsValidOffset([]any{trans, offAfter}) != types.JavaBoolTrue {
		t.Errorf("Overlap should validate offsetAfter")
	}
	otherOffset := createZoneOffset(0)
	if zoneOffsetTransitionIsValidOffset([]any{trans, otherOffset}) != types.JavaBoolFalse {
		t.Errorf("Overlap should not validate unrelated offset")
	}

	// toString
	strObj := zoneOffsetTransitionToString([]any{trans}).(*object.Object)
	str := object.GoStringFromStringObject(strObj)
	expectedStr := "Transition[Overlap at 2020-11-01T01:00-04:00 to -05:00]"
	if str != expectedStr {
		t.Errorf("toString mismatch:\ngot:  %s\nwant: %s", str, expectedStr)
	}
}

func TestZoneOffsetTransitionEpochAndInstant(t *testing.T) {
	zoneOffsetTransitionTestInit()

	d := CreateLocalDate(2020, 3, 8)
	tm := createLocalTime(2, 0, 0, 0)
	ldt := createLocalDateTime(d, tm)

	offBefore := createZoneOffset(-5 * 3600)
	offAfter := createZoneOffset(-4 * 3600)

	trans := zoneOffsetTransitionOf([]any{ldt, offBefore, offAfter}).(*object.Object)

	epochSec := zoneOffsetTransitionToEpochSecond([]any{trans}).(int64)
	// 2020-03-08 02:00:00 -05:00 is 2020-03-08 07:00:00 UTC
	instObj := zoneOffsetTransitionGetInstant([]any{trans}).(*object.Object)
	instSec, ok := getInstantEpochSecond(instObj)
	instNano := instantGetNano([]any{instObj}).(int64)
	if !ok || instSec != epochSec || instNano != 0 {
		t.Errorf("Instant mismatch: epochSec=%d, instSec=%d, nano=%d", epochSec, instSec, instNano)
	}

	// Test init from epoch
	transFromEpoch := object.MakeEmptyObjectWithClassName(&classNameZoneOffsetTransition)
	zoneOffsetTransitionInitEpoch([]any{transFromEpoch, epochSec, offBefore, offAfter})
	if zoneOffsetTransitionEquals([]any{trans, transFromEpoch}) != types.JavaBoolTrue {
		t.Errorf("Transition created from epoch mismatch")
	}
}

func TestZoneOffsetTransitionCompareAndEquals(t *testing.T) {
	zoneOffsetTransitionTestInit()

	d1 := CreateLocalDate(2020, 3, 8)
	tm1 := createLocalTime(2, 0, 0, 0)
	ldt1 := createLocalDateTime(d1, tm1)

	d2 := CreateLocalDate(2020, 11, 1)
	tm2 := createLocalTime(1, 0, 0, 0)
	ldt2 := createLocalDateTime(d2, tm2)

	offMinus5 := createZoneOffset(-5 * 3600)
	offMinus4 := createZoneOffset(-4 * 3600)

	trans1 := zoneOffsetTransitionOf([]any{ldt1, offMinus5, offMinus4}).(*object.Object)
	trans1Copy := zoneOffsetTransitionOf([]any{ldt1, offMinus5, offMinus4}).(*object.Object)
	trans2 := zoneOffsetTransitionOf([]any{ldt2, offMinus4, offMinus5}).(*object.Object)

	// equals & hashCode
	if zoneOffsetTransitionEquals([]any{trans1, trans1Copy}) != types.JavaBoolTrue {
		t.Errorf("trans1 should equal trans1Copy")
	}
	if zoneOffsetTransitionEquals([]any{trans1, trans2}) != types.JavaBoolFalse {
		t.Errorf("trans1 should not equal trans2")
	}
	h1 := zoneOffsetTransitionHashCode([]any{trans1}).(int64)
	h1Copy := zoneOffsetTransitionHashCode([]any{trans1Copy}).(int64)
	if h1 != h1Copy {
		t.Errorf("Hash code mismatch for equal objects: %d vs %d", h1, h1Copy)
	}

	// compareTo
	cmp1_2 := zoneOffsetTransitionCompareTo([]any{trans1, trans2}).(int64)
	if cmp1_2 >= 0 {
		t.Errorf("trans1 should be < trans2, got %d", cmp1_2)
	}
	cmp2_1 := zoneOffsetTransitionCompareTo([]any{trans2, trans1}).(int64)
	if cmp2_1 <= 0 {
		t.Errorf("trans2 should be > trans1, got %d", cmp2_1)
	}
	cmp1_1 := zoneOffsetTransitionCompareTo([]any{trans1, trans1Copy}).(int64)
	if cmp1_1 != 0 {
		t.Errorf("trans1 compareTo trans1Copy should be 0, got %d", cmp1_1)
	}
}
