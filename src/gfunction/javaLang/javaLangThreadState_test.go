/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestMain(m *testing.M) {
	// Initialize global string pool used by object and arrays utilities
	globals.InitStringPool()
	m.Run()
}

// --- threadStateToString() tests ---

func TestThreadStateValueOf_HappyPath(t *testing.T) {
	nameObj := object.StringObjectFromGoString("WAITING")
	res := threadStateValueOf([]interface{}{nameObj})
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected state object, got %T", res)
	}
	val := obj.FieldTable["value"].Fvalue.(int)
	if val != WAITING {
		t.Errorf("valueOf(\"WAITING\") = %d; want %d", val, WAITING)
	}
}

func TestThreadStateValueOf_MissingArg(t *testing.T) {
	res := threadStateValueOf([]interface{}{})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_WrongTypeArg(t *testing.T) {
	res := threadStateValueOf([]interface{}{123})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_NullName(t *testing.T) {
	res := threadStateValueOf([]interface{}{object.Null})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_NotStringObject(t *testing.T) {
	// Create a non-string object
	obj := object.MakeEmptyObject()
	res := threadStateValueOf([]interface{}{obj})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_NoMatch(t *testing.T) {
	nameObj := object.StringObjectFromGoString("BOGUS")
	res := threadStateValueOf([]interface{}{nameObj})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

// --- threadStateValues() tests ---

func TestThreadStateValues(t *testing.T) {
	res := threadStateValues(nil)
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected state object, got %T", res)
	}
	arr, ok := obj.FieldTable["value"].Fvalue.([]*object.Object)
	if !ok {
		t.Fatalf("expected []*object.Object, got %T", obj.FieldTable["value"].Fvalue)
	}
	if len(arr) != 6 {
		t.Errorf("expected 6 states, got %d", len(arr))
	}
	// Verify they are in order
	expectedNames := []string{"NEW", "RUNNABLE", "BLOCKED", "WAITING", "TIMED_WAITING", "TERMINATED"}
	for i, stateObj := range arr {
		val := stateObj.FieldTable["value"].Fvalue.(int)
		if val != i {
			t.Errorf("state at index %d has value %d; want %d", i, val, i)
		}
		// Verify toString matches
		strObj := threadStateToString([]interface{}{stateObj})
		name := object.GoStringFromStringObject(strObj.(*object.Object))
		if name != expectedNames[i] {
			t.Errorf("state at index %d has name %s; want %s", i, name, expectedNames[i])
		}
	}
}

// --- threadStateToString() tests ---

func TestThreadStateToString_HappyPath(t *testing.T) {
	stateObj := object.MakePrimitiveObject("java/lang/Thread$State", types.Int, RUNNABLE)
	res := threadStateToString([]interface{}{stateObj})
	strObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	name := object.GoStringFromStringObject(strObj)
	if name != "RUNNABLE" {
		t.Errorf("toString(RUNNABLE) = %s; want \"RUNNABLE\"", name)
	}
}

func TestThreadStateToString_MissingArg(t *testing.T) {
	res := threadStateToString([]interface{}{})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateToString_NotAnObject(t *testing.T) {
	res := threadStateToString([]interface{}{123})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateToString_MissingValueField(t *testing.T) {
	stateObj := object.MakeEmptyObject()
	res := threadStateToString([]interface{}{stateObj})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateToString_InvalidState(t *testing.T) {
	stateObj := object.MakePrimitiveObject("java/lang/Thread$State", types.Int, 99)
	res := threadStateToString([]interface{}{stateObj})
	gerr, ok := res.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}
