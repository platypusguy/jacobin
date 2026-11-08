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

func TestResolverStyle(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Format_ResolverStyle()

	// 1. Test registration of methods in MethodSignatures
	expectedSignatures := []string{
		"java/time/format/ResolverStyle.<clinit>()V",
		"java/time/format/ResolverStyle.<init>(Ljava/lang/String;I)V",
		"java/time/format/ResolverStyle.compareTo(Ljava/lang/Enum;)I",
		"java/time/format/ResolverStyle.compareTo(Ljava/lang/Object;)I",
		"java/time/format/ResolverStyle.compareTo(Ljava/time/format/ResolverStyle;)I",
		"java/time/format/ResolverStyle.equals(Ljava/lang/Object;)Z",
		"java/time/format/ResolverStyle.hashCode()I",
		"java/time/format/ResolverStyle.name()Ljava/lang/String;",
		"java/time/format/ResolverStyle.ordinal()I",
		"java/time/format/ResolverStyle.toString()Ljava/lang/String;",
		"java/time/format/ResolverStyle.valueOf(Ljava/lang/String;)Ljava/time/format/ResolverStyle;",
		"java/time/format/ResolverStyle.values()[Ljava/time/format/ResolverStyle;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}

	// 2. Run clinit and check statics
	resolverStyleClinit(nil)

	enumNames := []string{
		"STRICT",
		"SMART",
		"LENIENT",
	}

	for _, name := range enumNames {
		s, ok := statics.QueryStatic("java/time/format/ResolverStyle", name)
		if !ok || s.Value == nil {
			t.Fatalf("Static java/time/format/ResolverStyle.%s not found", name)
		}
	}

	// 3. Test values()
	valuesResult := resolverStyleValues(nil)
	valuesArray, okArray := valuesResult.(*object.Object)
	if !okArray || object.IsNull(valuesArray) {
		t.Fatalf("Expected non-null array from values(), got %v", valuesResult)
	}
	elements := valuesArray.FieldTable["value"].Fvalue.([]*object.Object)
	if len(elements) != 3 {
		t.Fatalf("Expected 3 elements in values(), got %d", len(elements))
	}

	for i, expectedName := range enumNames {
		inst := elements[i]
		nameObj := resolverStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected element %d to be %s, got %s", i, expectedName, object.GoStringFromStringObject(nameObj))
		}
		ordObj := resolverStyleOrdinal([]any{inst}).(int64)
		if ordObj != int64(i) {
			t.Errorf("Expected element %d ordinal to be %d, got %d", i, i, ordObj)
		}
	}

	// 4. Test valueOf(String)
	for _, expectedName := range enumNames {
		nameStr := object.StringObjectFromGoString(expectedName)
		res := resolverStyleValueOf([]any{nameStr})
		inst, ok := res.(*object.Object)
		if !ok || object.IsNull(inst) {
			t.Fatalf("valueOf(%s) failed, got %v", expectedName, res)
		}
		nameObj := resolverStyleName([]any{inst}).(*object.Object)
		if object.GoStringFromStringObject(nameObj) != expectedName {
			t.Errorf("Expected %s, got %s", expectedName, object.GoStringFromStringObject(nameObj))
		}
	}

	// Test invalid valueOf
	invalidStr := object.StringObjectFromGoString("INVALID_NAME")
	errRes := resolverStyleValueOf([]any{invalidStr})
	if !isError(errRes) {
		t.Errorf("Expected error for invalid valueOf name, got %v", errRes)
	}

	nullErr := resolverStyleValueOf([]any{(*object.Object)(nil)})
	if !isError(nullErr) {
		t.Errorf("Expected error for null valueOf name, got %v", nullErr)
	}

	// 5. Test toString(), compareTo(), equals(), hashCode()
	strictInst := elements[0]
	smartInst := elements[1]
	lenientInst := elements[2]

	toStr := resolverStyleToString([]any{strictInst})
	if object.GoStringFromStringObject(toStr.(*object.Object)) != "STRICT" {
		t.Errorf("toString() expected 'STRICT', got %s", object.GoStringFromStringObject(toStr.(*object.Object)))
	}

	// compareTo
	cmpSelf := resolverStyleCompareTo([]any{strictInst, strictInst})
	if cmpSelf.(int64) != 0 {
		t.Errorf("compareTo self expected 0, got %v", cmpSelf)
	}

	cmpLess := resolverStyleCompareTo([]any{strictInst, smartInst})
	if cmpLess.(int64) >= 0 {
		t.Errorf("compareTo STRICT vs SMART expected < 0, got %v", cmpLess)
	}

	cmpGreater := resolverStyleCompareTo([]any{lenientInst, strictInst})
	if cmpGreater.(int64) <= 0 {
		t.Errorf("compareTo LENIENT vs STRICT expected > 0, got %v", cmpGreater)
	}

	// equals
	eqSelf := resolverStyleEquals([]any{strictInst, strictInst})
	if eqSelf != types.JavaBoolTrue {
		t.Errorf("equals self expected true, got %v", eqSelf)
	}

	eqDiff := resolverStyleEquals([]any{strictInst, smartInst})
	if eqDiff != types.JavaBoolFalse {
		t.Errorf("equals different expected false, got %v", eqDiff)
	}

	// hashCode
	hc := resolverStyleHashCode([]any{smartInst})
	if hc.(int64) != 1 {
		t.Errorf("hashCode expected 1, got %v", hc)
	}

	// Null self checks
	if !isError(resolverStyleName([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for name with null self")
	}
	if !isError(resolverStyleOrdinal([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for ordinal with null self")
	}
	if !isError(resolverStyleCompareTo([]any{(*object.Object)(nil), strictInst})) {
		t.Errorf("Expected error for compareTo with null self")
	}
	if !isError(resolverStyleCompareTo([]any{strictInst, (*object.Object)(nil)})) {
		t.Errorf("Expected error for compareTo with null other")
	}
	if !isError(resolverStyleHashCode([]any{(*object.Object)(nil)})) {
		t.Errorf("Expected error for hashCode with null self")
	}

	// Test init
	freshObj := object.MakeEmptyObjectWithClassName(&classNameResolverStyle)
	initRes := resolverStyleInit([]any{freshObj, object.StringObjectFromGoString("STRICT"), int64(0)})
	if initRes != nil {
		t.Errorf("Expected nil from init, got %v", initRes)
	}
	if freshObj.FieldTable["ordinal"].Fvalue.(int64) != 0 {
		t.Errorf("Expected ordinal 0, got %v", freshObj.FieldTable["ordinal"].Fvalue)
	}
}
