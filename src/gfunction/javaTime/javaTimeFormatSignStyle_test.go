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

func TestSignStyle(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Format_SignStyle()

	// 1. Test registration of methods in MethodSignatures
	expectedSignatures := []string{
		"java/time/format/SignStyle.<clinit>()V",
		"java/time/format/SignStyle.<init>(Ljava/lang/String;I)V",
		"java/time/format/SignStyle.compareTo(Ljava/lang/Enum;)I",
		"java/time/format/SignStyle.compareTo(Ljava/lang/Object;)I",
		"java/time/format/SignStyle.compareTo(Ljava/time/format/SignStyle;)I",
		"java/time/format/SignStyle.equals(Ljava/lang/Object;)Z",
		"java/time/format/SignStyle.hashCode()I",
		"java/time/format/SignStyle.name()Ljava/lang/String;",
		"java/time/format/SignStyle.ordinal()I",
		"java/time/format/SignStyle.toString()Ljava/lang/String;",
		"java/time/format/SignStyle.valueOf(Ljava/lang/String;)Ljava/time/format/SignStyle;",
		"java/time/format/SignStyle.values()[Ljava/time/format/SignStyle;",
		"java/time/format/SignStyle.parse(ZZZ)Z",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}

	// 2. Run clinit and check statics
	signStyleClinit(nil)

	enumNames := []string{
		"NORMAL",
		"ALWAYS",
		"NEVER",
		"NOT_NEGATIVE",
		"EXCEEDS_PAD",
	}

	for _, name := range enumNames {
		s, ok := statics.QueryStatic("java/time/format/SignStyle", name)
		if !ok || s.Value == nil {
			t.Fatalf("Static java/time/format/SignStyle.%s not found", name)
		}
	}

	// 3. Test values()
	valuesResult := signStyleValues(nil)
	valuesArray, okArray := valuesResult.(*object.Object)
	if !okArray || object.IsNull(valuesArray) {
		t.Fatalf("Expected non-null array from values(), got %v", valuesResult)
	}
	elements := valuesArray.FieldTable["value"].Fvalue.([]*object.Object)
	if len(elements) != 5 {
		t.Fatalf("Expected 5 elements in values(), got %d", len(elements))
	}

	for i, expectedName := range enumNames {
		inst := elements[i]
		nameObj := signStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected element %d to be %s, got %s", i, expectedName, object.GoStringFromStringObject(nameObj))
		}
		ordObj := signStyleOrdinal([]any{inst}).(int64)
		if ordObj != int64(i) {
			t.Errorf("Expected element %d ordinal to be %d, got %d", i, i, ordObj)
		}
	}

	// 4. Test valueOf(String)
	for _, expectedName := range enumNames {
		nameStr := object.StringObjectFromGoString(expectedName)
		res := signStyleValueOf([]any{nameStr})
		inst, ok := res.(*object.Object)
		if !ok || object.IsNull(inst) {
			t.Fatalf("valueOf(%s) failed, got %v", expectedName, res)
		}
		nameObj := signStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected %s, got %s", expectedName, object.GoStringFromStringObject(nameObj))
		}
	}

	// Test invalid valueOf
	invalidStr := object.StringObjectFromGoString("INVALID_NAME")
	errRes := signStyleValueOf([]any{invalidStr})
	if !isError(errRes) {
		t.Errorf("Expected error for invalid valueOf name, got %v", errRes)
	}

	nullErr := signStyleValueOf([]any{(*object.Object)(nil)})
	if !isError(nullErr) {
		t.Errorf("Expected error for null valueOf name, got %v", nullErr)
	}

	// 5. Test toString(), compareTo(), equals(), hashCode()
	normalInst := elements[0]
	alwaysInst := elements[1]
	neverInst := elements[2]
	notNegInst := elements[3]
	exceedsPadInst := elements[4]

	toStr := signStyleToString([]any{normalInst})
	if object.GoStringFromStringObject(toStr.(*object.Object)) != "NORMAL" {
		t.Errorf("toString() expected 'NORMAL', got %s", object.GoStringFromStringObject(toStr.(*object.Object)))
	}

	// compareTo
	cmpSelf := signStyleCompareTo([]any{normalInst, normalInst})
	if cmpSelf.(int64) != 0 {
		t.Errorf("compareTo self expected 0, got %v", cmpSelf)
	}

	cmpLess := signStyleCompareTo([]any{normalInst, alwaysInst})
	if cmpLess.(int64) >= 0 {
		t.Errorf("compareTo NORMAL vs ALWAYS expected < 0, got %v", cmpLess)
	}

	cmpGreater := signStyleCompareTo([]any{exceedsPadInst, normalInst})
	if cmpGreater.(int64) <= 0 {
		t.Errorf("compareTo EXCEEDS_PAD vs NORMAL expected > 0, got %v", cmpGreater)
	}

	// equals
	eqSelf := signStyleEquals([]any{normalInst, normalInst})
	if eqSelf != types.JavaBoolTrue {
		t.Errorf("equals self expected true, got %v", eqSelf)
	}

	eqDiff := signStyleEquals([]any{normalInst, alwaysInst})
	if eqDiff != types.JavaBoolFalse {
		t.Errorf("equals different expected false, got %v", eqDiff)
	}

	// hashCode
	hc := signStyleHashCode([]any{alwaysInst})
	if hc.(int64) != 1 {
		t.Errorf("hashCode expected 1, got %v", hc)
	}

	// 6. Test parse(positive, strict, fixedWidth)
	// NORMAL (0): !positive || !strict
	if res := signStyleParse([]any{normalInst, types.JavaBoolTrue, types.JavaBoolTrue, types.JavaBoolFalse}); res != types.JavaBoolFalse {
		t.Errorf("NORMAL parse(true, true, false) expected false, got %v", res)
	}
	if res := signStyleParse([]any{normalInst, types.JavaBoolFalse, types.JavaBoolTrue, types.JavaBoolFalse}); res != types.JavaBoolTrue {
		t.Errorf("NORMAL parse(false, true, false) expected true, got %v", res)
	}
	if res := signStyleParse([]any{normalInst, types.JavaBoolTrue, types.JavaBoolFalse, types.JavaBoolFalse}); res != types.JavaBoolTrue {
		t.Errorf("NORMAL parse(true, false, false) expected true, got %v", res)
	}

	// ALWAYS (1): true
	if res := signStyleParse([]any{alwaysInst, types.JavaBoolTrue, types.JavaBoolTrue, types.JavaBoolTrue}); res != types.JavaBoolTrue {
		t.Errorf("ALWAYS parse expected true, got %v", res)
	}

	// EXCEEDS_PAD (4): true
	if res := signStyleParse([]any{exceedsPadInst, types.JavaBoolFalse, types.JavaBoolTrue, types.JavaBoolTrue}); res != types.JavaBoolTrue {
		t.Errorf("EXCEEDS_PAD parse expected true, got %v", res)
	}

	// NEVER (2) & NOT_NEGATIVE (3): !strict && !fixedWidth
	if res := signStyleParse([]any{neverInst, types.JavaBoolTrue, types.JavaBoolFalse, types.JavaBoolFalse}); res != types.JavaBoolTrue {
		t.Errorf("NEVER parse(true, false, false) expected true, got %v", res)
	}
	if res := signStyleParse([]any{neverInst, types.JavaBoolTrue, types.JavaBoolTrue, types.JavaBoolFalse}); res != types.JavaBoolFalse {
		t.Errorf("NEVER parse(true, true, false) expected false, got %v", res)
	}
	if res := signStyleParse([]any{notNegInst, types.JavaBoolTrue, types.JavaBoolFalse, types.JavaBoolTrue}); res != types.JavaBoolFalse {
		t.Errorf("NOT_NEGATIVE parse(true, false, true) expected false, got %v", res)
	}

	// Null self checks
	if !isError(signStyleName([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for name with null self")
	}
	if !isError(signStyleOrdinal([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for ordinal with null self")
	}
	if !isError(signStyleCompareTo([]any{(*object.Object)(nil), normalInst})) {
		t.Errorf("Expected error for compareTo with null self")
	}
	if !isError(signStyleCompareTo([]any{normalInst, (*object.Object)(nil)})) {
		t.Errorf("Expected error for compareTo with null other")
	}
	if !isError(signStyleHashCode([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for hashCode with null self")
	}
	if !isError(signStyleParse([]any{(*object.Object)(nil), types.JavaBoolTrue, types.JavaBoolTrue, types.JavaBoolFalse})) {
		t.Errorf("Expected error for parse with null self")
	}

	// Test init
	freshObj := object.MakeEmptyObjectWithClassName(&classNameSignStyle)
	initRes := signStyleInit([]any{freshObj, object.StringObjectFromGoString("NORMAL"), int64(0)})
	if initRes != nil {
		t.Errorf("Expected nil from init, got %v", initRes)
	}
}
