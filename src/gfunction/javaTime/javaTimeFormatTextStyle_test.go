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

func isError(val any) bool {
	if val == nil {
		return false
	}
	_, ok := val.(*ghelpers.GErrBlk)
	return ok
}

func TestTextStyle(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Format_TextStyle()

	// 1. Test registration of methods in MethodSignatures
	expectedSignatures := []string{
		"java/time/format/TextStyle.<clinit>()V",
		"java/time/format/TextStyle.<init>(Ljava/lang/String;I)V",
		"java/time/format/TextStyle.<init>(Ljava/lang/String;III)V",
		"java/time/format/TextStyle.asNormal()Ljava/time/format/TextStyle;",
		"java/time/format/TextStyle.asStandalone()Ljava/time/format/TextStyle;",
		"java/time/format/TextStyle.compareTo(Ljava/lang/Enum;)I",
		"java/time/format/TextStyle.compareTo(Ljava/lang/Object;)I",
		"java/time/format/TextStyle.compareTo(Ljava/time/format/TextStyle;)I",
		"java/time/format/TextStyle.equals(Ljava/lang/Object;)Z",
		"java/time/format/TextStyle.hashCode()I",
		"java/time/format/TextStyle.isStandalone()Z",
		"java/time/format/TextStyle.name()Ljava/lang/String;",
		"java/time/format/TextStyle.ordinal()I",
		"java/time/format/TextStyle.toCalendarStyle()I",
		"java/time/format/TextStyle.toString()Ljava/lang/String;",
		"java/time/format/TextStyle.valueOf(Ljava/lang/String;)Ljava/time/format/TextStyle;",
		"java/time/format/TextStyle.values()[Ljava/time/format/TextStyle;",
		"java/time/format/TextStyle.zoneNameStyleIndex()I",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}

	// 2. Run clinit and check statics
	textStyleClinit(nil)

	enumNames := []string{
		"FULL",
		"FULL_STANDALONE",
		"SHORT",
		"SHORT_STANDALONE",
		"NARROW",
		"NARROW_STANDALONE",
	}

	for _, name := range enumNames {
		s, ok := statics.QueryStatic("java/time/format/TextStyle", name)
		if !ok || s.Value == nil {
			t.Fatalf("Static java/time/format/TextStyle.%s not found", name)
		}
	}

	// 3. Test values()
	valuesResult := textStyleValues(nil)
	valuesArray, okArray := valuesResult.(*object.Object)
	if !okArray || object.IsNull(valuesArray) {
		t.Fatalf("Expected non-null array from values(), got %v", valuesResult)
	}
	elements := valuesArray.FieldTable["value"].Fvalue.([]*object.Object)
	if len(elements) != 6 {
		t.Fatalf("Expected 6 elements in values(), got %d", len(elements))
	}

	for i, expectedName := range enumNames {
		inst := elements[i]
		nameObj := textStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected element %d to be %s, got %s", i, expectedName, object.GoStringFromStringObject(nameObj))
		}
		ordObj := textStyleOrdinal([]any{inst}).(int64)
		if ordObj != int64(i) {
			t.Errorf("Expected element %d ordinal to be %d, got %d", i, i, ordObj)
		}
	}

	// 4. Test valueOf(String)
	for _, expectedName := range enumNames {
		nameStr := object.StringObjectFromGoString(expectedName)
		res := textStyleValueOf([]any{nameStr})
		inst, ok := res.(*object.Object)
		if !ok || object.IsNull(inst) {
			t.Fatalf("valueOf(%s) failed, got %v", expectedName, res)
		}
		nameObj := textStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected %s, got %s", expectedName, object.GoStringFromStringObject(nameObj))
		}
	}

	// Test invalid valueOf
	invalidStr := object.StringObjectFromGoString("INVALID_NAME")
	errRes := textStyleValueOf([]any{invalidStr})
	if !isError(errRes) {
		t.Errorf("Expected error for invalid valueOf name, got %v", errRes)
	}

	nullErr := textStyleValueOf([]any{(*object.Object)(nil)})
	if !isError(nullErr) {
		t.Errorf("Expected error for null valueOf name, got %v", nullErr)
	}

	// 5. Test isStandalone()
	expectedStandalone := []bool{false, true, false, true, false, true}
	for i, inst := range elements {
		isSt := textStyleIsStandalone([]any{inst})
		expected := types.JavaBoolFalse
		if expectedStandalone[i] {
			expected = types.JavaBoolTrue
		}
		if isSt != expected {
			t.Errorf("%s isStandalone expected %v, got %v", enumNames[i], expectedStandalone[i], isSt)
		}
	}

	// 6. Test asStandalone()
	// FULL (0) -> FULL_STANDALONE (1)
	// FULL_STANDALONE (1) -> FULL_STANDALONE (1)
	// SHORT (2) -> SHORT_STANDALONE (3)
	// SHORT_STANDALONE (3) -> SHORT_STANDALONE (3)
	// NARROW (4) -> NARROW_STANDALONE (5)
	// NARROW_STANDALONE (5) -> NARROW_STANDALONE (5)
	expectedStandaloneIdx := []int{1, 1, 3, 3, 5, 5}
	for i, inst := range elements {
		res := textStyleAsStandalone([]any{inst})
		resInst, ok := res.(*object.Object)
		if !ok || resInst != elements[expectedStandaloneIdx[i]] {
			t.Errorf("%s asStandalone expected %s, got %v", enumNames[i], enumNames[expectedStandaloneIdx[i]], res)
		}
	}

	// 7. Test asNormal()
	// FULL (0) -> FULL (0)
	// FULL_STANDALONE (1) -> FULL (0)
	// SHORT (2) -> SHORT (2)
	// SHORT_STANDALONE (3) -> SHORT (2)
	// NARROW (4) -> NARROW (4)
	// NARROW_STANDALONE (5) -> NARROW (4)
	expectedNormalIdx := []int{0, 0, 2, 2, 4, 4}
	for i, inst := range elements {
		res := textStyleAsNormal([]any{inst})
		resInst, ok := res.(*object.Object)
		if !ok || resInst != elements[expectedNormalIdx[i]] {
			t.Errorf("%s asNormal expected %s, got %v", enumNames[i], enumNames[expectedNormalIdx[i]], res)
		}
	}

	// 8. Test toCalendarStyle()
	expectedCalStyles := []int64{2, 32770, 1, 32769, 4, 32772}
	for i, inst := range elements {
		calStyle := textStyleToCalendarStyle([]any{inst})
		if calStyle.(int64) != expectedCalStyles[i] {
			t.Errorf("%s toCalendarStyle expected %d, got %v", enumNames[i], expectedCalStyles[i], calStyle)
		}
	}

	// 9. Test zoneNameStyleIndex()
	expectedZoneIdx := []int64{0, 0, 1, 1, 2, 2}
	for i, inst := range elements {
		zoneIdx := textStyleZoneNameStyleIndex([]any{inst})
		if zoneIdx.(int64) != expectedZoneIdx[i] {
			t.Errorf("%s zoneNameStyleIndex expected %d, got %v", enumNames[i], expectedZoneIdx[i], zoneIdx)
		}
	}

	// 10. Test toString(), compareTo(), equals(), hashCode()
	fullInst := elements[0]
	fullStandInst := elements[1]
	shortInst := elements[2]

	toStr := textStyleToString([]any{fullInst})
	if object.GoStringFromStringObject(toStr.(*object.Object)) != "FULL" {
		t.Errorf("toString() expected 'FULL', got %s", object.GoStringFromStringObject(toStr.(*object.Object)))
	}

	// compareTo
	cmpSelf := textStyleCompareTo([]any{fullInst, fullInst})
	if cmpSelf.(int64) != 0 {
		t.Errorf("compareTo self expected 0, got %v", cmpSelf)
	}

	cmpLess := textStyleCompareTo([]any{fullInst, shortInst})
	if cmpLess.(int64) >= 0 {
		t.Errorf("compareTo FULL vs SHORT expected < 0, got %v", cmpLess)
	}

	cmpGreater := textStyleCompareTo([]any{shortInst, fullInst})
	if cmpGreater.(int64) <= 0 {
		t.Errorf("compareTo SHORT vs FULL expected > 0, got %v", cmpGreater)
	}

	// equals
	eqSelf := textStyleEquals([]any{fullInst, fullInst})
	if eqSelf != types.JavaBoolTrue {
		t.Errorf("equals self expected true, got %v", eqSelf)
	}

	eqDiff := textStyleEquals([]any{fullInst, fullStandInst})
	if eqDiff != types.JavaBoolFalse {
		t.Errorf("equals different expected false, got %v", eqDiff)
	}

	// hashCode
	hc := textStyleHashCode([]any{fullStandInst})
	if hc.(int64) != 1 {
		t.Errorf("hashCode expected 1, got %v", hc)
	}

	// Null self checks
	if !isError(textStyleIsStandalone([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for isStandalone with null self")
	}
	if !isError(textStyleAsStandalone([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for asStandalone with null self")
	}
	if !isError(textStyleAsNormal([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for asNormal with null self")
	}
	if !isError(textStyleToCalendarStyle([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for toCalendarStyle with null self")
	}
	if !isError(textStyleZoneNameStyleIndex([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for zoneNameStyleIndex with null self")
	}
	if !isError(textStyleCompareTo([]any{(*object.Object)(nil), fullInst})) {
		t.Errorf("Expected error for compareTo with null self")
	}
	if !isError(textStyleCompareTo([]any{fullInst, (*object.Object)(nil)})) {
		t.Errorf("Expected error for compareTo with null other")
	}
	if !isError(textStyleHashCode([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for hashCode with null self")
	}

	// Test init
	freshObj := object.MakeEmptyObjectWithClassName(&classNameTextStyle)
	textStyleInit([]any{
		freshObj,
		object.StringObjectFromGoString("CUSTOM"),
		int64(10),
		int64(99),
		int64(42),
	})
	if nameVal, ok := freshObj.FieldTable["name"]; !ok || object.GoStringFromStringObject(nameVal.Fvalue.(*object.Object)) != "CUSTOM" {
		t.Errorf("init did not set name properly")
	}
	if ordVal, ok := freshObj.FieldTable["ordinal"]; !ok || ordVal.Fvalue.(int64) != 10 {
		t.Errorf("init did not set ordinal properly")
	}
	if calVal, ok := freshObj.FieldTable["calendarStyle"]; !ok || calVal.Fvalue.(int64) != 99 {
		t.Errorf("init did not set calendarStyle properly")
	}
	if zoneVal, ok := freshObj.FieldTable["zoneNameStyleIndex"]; !ok || zoneVal.Fvalue.(int64) != 42 {
		t.Errorf("init did not set zoneNameStyleIndex properly")
	}
}
