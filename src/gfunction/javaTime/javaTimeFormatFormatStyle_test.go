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

func TestFormatStyle(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Format_FormatStyle()

	// 1. Test registration of methods in MethodSignatures
	expectedSignatures := []string{
		"java/time/format/FormatStyle.<clinit>()V",
		"java/time/format/FormatStyle.<init>(Ljava/lang/String;I)V",
		"java/time/format/FormatStyle.compareTo(Ljava/lang/Enum;)I",
		"java/time/format/FormatStyle.compareTo(Ljava/lang/Object;)I",
		"java/time/format/FormatStyle.compareTo(Ljava/time/format/FormatStyle;)I",
		"java/time/format/FormatStyle.equals(Ljava/lang/Object;)Z",
		"java/time/format/FormatStyle.hashCode()I",
		"java/time/format/FormatStyle.name()Ljava/lang/String;",
		"java/time/format/FormatStyle.ordinal()I",
		"java/time/format/FormatStyle.toString()Ljava/lang/String;",
		"java/time/format/FormatStyle.valueOf(Ljava/lang/String;)Ljava/time/format/FormatStyle;",
		"java/time/format/FormatStyle.values()[Ljava/time/format/FormatStyle;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}

	// 2. Run clinit and check statics
	formatStyleClinit(nil)

	enumNames := []string{
		"FULL",
		"LONG",
		"MEDIUM",
		"SHORT",
	}

	for _, name := range enumNames {
		s, ok := statics.QueryStatic("java/time/format/FormatStyle", name)
		if !ok || s.Value == nil {
			t.Fatalf("Static java/time/format/FormatStyle.%s not found", name)
		}
	}

	// 3. Test values()
	valuesResult := formatStyleValues(nil)
	valuesArray, okArray := valuesResult.(*object.Object)
	if !okArray || object.IsNull(valuesArray) {
		t.Fatalf("Expected non-null array from values(), got %v", valuesResult)
	}
	elements := valuesArray.FieldTable["value"].Fvalue.([]*object.Object)
	if len(elements) != 4 {
		t.Fatalf("Expected 4 elements in values(), got %d", len(elements))
	}

	for i, expectedName := range enumNames {
		inst := elements[i]
		nameObj := formatStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected element %d to be %s, got %s", i, expectedName, object.GoStringFromStringObject(nameObj))
		}
		ordObj := formatStyleOrdinal([]any{inst}).(int64)
		if ordObj != int64(i) {
			t.Errorf("Expected element %d ordinal to be %d, got %d", i, i, ordObj)
		}
	}

	// 4. Test valueOf(String)
	for _, expectedName := range enumNames {
		nameStr := object.StringObjectFromGoString(expectedName)
		res := formatStyleValueOf([]any{nameStr})
		inst, ok := res.(*object.Object)
		if !ok || object.IsNull(inst) {
			t.Fatalf("valueOf(%s) failed, got %v", expectedName, res)
		}
		nameObj := formatStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected %s, got %s", expectedName, object.GoStringFromStringObject(nameObj))
		}
	}

	// Test invalid valueOf
	invalidStr := object.StringObjectFromGoString("INVALID_NAME")
	errRes := formatStyleValueOf([]any{invalidStr})
	if !isError(errRes) {
		t.Errorf("Expected error for invalid valueOf name, got %v", errRes)
	}

	nullErr := formatStyleValueOf([]any{(*object.Object)(nil)})
	if !isError(nullErr) {
		t.Errorf("Expected error for null valueOf name, got %v", nullErr)
	}

	// 5. Test toString(), compareTo(), equals(), hashCode()
	fullInst := elements[0]
	longInst := elements[1]
	mediumInst := elements[2]

	toStr := formatStyleToString([]any{fullInst})
	if object.GoStringFromStringObject(toStr.(*object.Object)) != "FULL" {
		t.Errorf("toString() expected 'FULL', got %s", object.GoStringFromStringObject(toStr.(*object.Object)))
	}

	// compareTo
	cmpSelf := formatStyleCompareTo([]any{fullInst, fullInst})
	if cmpSelf.(int64) != 0 {
		t.Errorf("compareTo self expected 0, got %v", cmpSelf)
	}

	cmpLess := formatStyleCompareTo([]any{fullInst, longInst})
	if cmpLess.(int64) >= 0 {
		t.Errorf("compareTo FULL vs LONG expected < 0, got %v", cmpLess)
	}

	cmpGreater := formatStyleCompareTo([]any{mediumInst, fullInst})
	if cmpGreater.(int64) <= 0 {
		t.Errorf("compareTo MEDIUM vs FULL expected > 0, got %v", cmpGreater)
	}

	// equals
	eqSelf := formatStyleEquals([]any{fullInst, fullInst})
	if eqSelf != types.JavaBoolTrue {
		t.Errorf("equals self expected true, got %v", eqSelf)
	}

	eqDiff := formatStyleEquals([]any{fullInst, longInst})
	if eqDiff != types.JavaBoolFalse {
		t.Errorf("equals different expected false, got %v", eqDiff)
	}

	// hashCode
	hc := formatStyleHashCode([]any{longInst})
	if hc.(int64) != 1 {
		t.Errorf("hashCode expected 1, got %v", hc)
	}

	// Null self checks
	if !isError(formatStyleName([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for name with null self")
	}
	if !isError(formatStyleOrdinal([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for ordinal with null self")
	}
	if !isError(formatStyleCompareTo([]any{(*object.Object)(nil), fullInst})) {
		t.Errorf("Expected error for compareTo with null self")
	}
	if !isError(formatStyleCompareTo([]any{fullInst, (*object.Object)(nil)})) {
		t.Errorf("Expected error for compareTo with null other")
	}
	if !isError(formatStyleHashCode([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for hashCode with null self")
	}

	// Test init
	freshObj := object.MakeEmptyObjectWithClassName(&classNameFormatStyle)
	formatStyleInit([]any{
		freshObj,
		object.StringObjectFromGoString("CUSTOM"),
		int64(10),
	})
	if nameVal, ok := freshObj.FieldTable["name"]; !ok || object.GoStringFromStringObject(nameVal.Fvalue.(*object.Object)) != "CUSTOM" {
		t.Errorf("init did not set name properly")
	}
	if ordVal, ok := freshObj.FieldTable["ordinal"]; !ok || ordVal.Fvalue.(int64) != 10 {
		t.Errorf("init did not set ordinal properly")
	}
}
