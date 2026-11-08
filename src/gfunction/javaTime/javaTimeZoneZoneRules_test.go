/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func zoneRulesTestInit() {
	globals.InitGlobals("test")
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
}

func TestZoneRulesFixedOffset(t *testing.T) {
	zoneRulesTestInit()

	offset := createZoneOffset(3600) // +01:00
	rulesRes := zoneRulesOfOffset([]any{offset})
	rules := rulesRes.(*object.Object)

	// isFixedOffset
	if zoneRulesIsFixedOffset([]any{rules}) != types.JavaBoolTrue {
		t.Fatal("isFixedOffset() = false, want true")
	}

	// instant & localDateTime
	instantObj := createInstantInternal(1700000000, 0)
	dateObj := CreateLocalDate(2026, 6, 24)
	timeObj := createLocalTime(14, 30, 0, 0)
	ldtObj := CreateLocalDateTime(dateObj, timeObj)

	// getOffset(Instant)
	offInst := zoneRulesGetOffsetInstant([]any{rules, instantObj}).(*object.Object)
	secInst, _ := getZoneOffsetTotalSecondsVal(offInst)
	if secInst != 3600 {
		t.Fatalf("getOffset(Instant) = %d, want 3600", secInst)
	}

	// getOffset(LocalDateTime)
	offLDT := zoneRulesGetOffsetLocalDateTime([]any{rules, ldtObj}).(*object.Object)
	secLDT, _ := getZoneOffsetTotalSecondsVal(offLDT)
	if secLDT != 3600 {
		t.Fatalf("getOffset(LocalDateTime) = %d, want 3600", secLDT)
	}

	// getStandardOffset(Instant)
	offStd := zoneRulesGetStandardOffset([]any{rules, instantObj}).(*object.Object)
	secStd, _ := getZoneOffsetTotalSecondsVal(offStd)
	if secStd != 3600 {
		t.Fatalf("getStandardOffset(Instant) = %d, want 3600", secStd)
	}

	// getDaylightSavings(Instant)
	durObj := zoneRulesGetDaylightSavings([]any{rules, instantObj}).(*object.Object)
	durSec, _ := getInt64Val(durObj.FieldTable["seconds"].Fvalue)
	if durSec != 0 {
		t.Fatalf("getDaylightSavings(Instant) = %d, want 0", durSec)
	}

	// isDaylightSavings(Instant)
	if zoneRulesIsDaylightSavings([]any{rules, instantObj}) != types.JavaBoolFalse {
		t.Fatal("isDaylightSavings(Instant) = true, want false")
	}

	// getValidOffsets(LocalDateTime)
	validList := zoneRulesGetValidOffsets([]any{rules, ldtObj}).(*object.Object)
	elems := validList.FieldTable["elements"].Fvalue.([]any)
	if len(elems) != 1 || elems[0].(*object.Object) != offset {
		t.Fatalf("getValidOffsets() = %v, want [%v]", elems, offset)
	}

	// isValidOffset
	if zoneRulesIsValidOffset([]any{rules, ldtObj, offset}) != types.JavaBoolTrue {
		t.Fatal("isValidOffset(offset) = false, want true")
	}
	otherOffset := createZoneOffset(7200)
	if zoneRulesIsValidOffset([]any{rules, ldtObj, otherOffset}) != types.JavaBoolFalse {
		t.Fatal("isValidOffset(otherOffset) = true, want false")
	}

	// getTransitions / getTransitionRules
	trList := zoneRulesGetTransitions([]any{rules}).(*object.Object)
	if len(trList.FieldTable["elements"].Fvalue.([]any)) != 0 {
		t.Fatal("getTransitions() expected empty list")
	}
	trRulesList := zoneRulesGetTransitionRules([]any{rules}).(*object.Object)
	if len(trRulesList.FieldTable["elements"].Fvalue.([]any)) != 0 {
		t.Fatal("getTransitionRules() expected empty list")
	}

	// equals, hashCode, toString
	rules2 := zoneRulesOfOffset([]any{offset}).(*object.Object)
	if zoneRulesEquals([]any{rules, rules2}) != types.JavaBoolTrue {
		t.Fatal("equals() = false, want true")
	}
	rulesOther := zoneRulesOfOffset([]any{otherOffset}).(*object.Object)
	if zoneRulesEquals([]any{rules, rulesOther}) != types.JavaBoolFalse {
		t.Fatal("equals() = true, want false")
	}

	strObj := zoneRulesToString([]any{rules}).(*object.Object)
	str := object.GoStringFromStringObject(strObj)
	if str != "ZoneRules[currentStandardOffset=+01:00]" {
		t.Fatalf("toString() = %q, want 'ZoneRules[currentStandardOffset=+01:00]'", str)
	}
}

func TestZoneRulesNamedZone(t *testing.T) {
	zoneRulesTestInit()

	nyRules := createZoneRulesForZoneId("America/New_York")
	if zoneRulesIsFixedOffset([]any{nyRules}) != types.JavaBoolFalse {
		t.Fatal("America/New_York isFixedOffset() = true, want false")
	}

	// Summer instant (June 2026): EDT (-04:00)
	summerInstant := createInstantInternal(1782259200, 0) // approx June 2026
	offSummer := zoneRulesGetOffsetInstant([]any{nyRules, summerInstant}).(*object.Object)
	if offSummer == nil {
		t.Fatal("getOffset(summerInstant) returned nil")
	}
}
