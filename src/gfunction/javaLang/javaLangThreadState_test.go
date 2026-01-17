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
	val := obj.FieldTable["value"].Fvalue.(int64)
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

func TestThreadStateValues(t *testing.T) {
	//expectedNames := []string{"NEW", "RUNNABLE", "BLOCKED", "WAITING", "TIMED_WAITING", "TERMINATED"}
	strValuesObj := threadStateValues([]interface{}{nil}).(*object.Object)
	strValues := strValuesObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(strValues) != 6 {
		t.Errorf("expected 6 states, got %d", len(strValues))
	}
	for ix, obj := range strValues {
		state := obj.FieldTable["value"].Fvalue.(int64)
		if state != int64(ix) {
			t.Errorf("expected state = %d but observed %d", ix, state)
		}
	}
}

func TestThreadStateToString_HappyPath(t *testing.T) {
	obj := object.MakePrimitiveObject("java/lang/Thread$State", types.Int, RUNNABLE)
	res := ThreadStateToString([]interface{}{obj})
	strObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected string object, got %T", res)
	}
	got := object.GoStringFromStringObject(strObj)
	if got != "RUNNABLE" {
		t.Errorf("ThreadStateToString(RUNNABLE) = %q; want \"RUNNABLE\"", got)
	}
}

func TestThreadStateToString_Errors(t *testing.T) {
	// Missing object
	res := ThreadStateToString([]interface{}{})
	if gerr := res.(*ghelpers.GErrBlk); gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for missing object")
	}

	// Wrong type
	res = ThreadStateToString([]interface{}{123})
	if gerr := res.(*ghelpers.GErrBlk); gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for wrong type")
	}

	// Missing value field
	obj := object.MakeEmptyObject()
	res = ThreadStateToString([]interface{}{obj})
	if gerr := res.(*ghelpers.GErrBlk); gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for missing value field")
	}

	// Unknown state value
	obj = object.MakePrimitiveObject("java/lang/Thread$State", types.Int, int64(999))
	res = ThreadStateToString([]interface{}{obj})
	if gerr := res.(*ghelpers.GErrBlk); gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for unknown state value")
	}
}
