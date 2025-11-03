/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/excNames"
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

func TestThreadStateToString_HappyPath(t *testing.T) {
	res := threadStateToString([]interface{}{RUNNABLE})
	if !object.IsStringObject(res) {
		t.Fatalf("expected String object, got %T", res)
	}
	got := object.GoStringFromStringObject(res.(*object.Object))
	if got != "RUNNABLE" {
		t.Errorf("toString() returned %q; want %q", got, "RUNNABLE")
	}
}

func TestThreadStateToString_MissingParam(t *testing.T) {
	res := threadStateToString([]interface{}{})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateToString_WrongTypeParam(t *testing.T) {
	// Pass a non-int; code treats it as null type error
	res := threadStateToString([]interface{}{"not-an-int"})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateToString_InvalidState(t *testing.T) {
	res := threadStateToString([]interface{}{42})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

// --- threadStateValueOf() tests ---

func TestThreadStateValueOf_HappyPath(t *testing.T) {
	nameObj := object.StringObjectFromGoString("WAITING")
	res := threadStateValueOf([]interface{}{nameObj})
	val, ok := res.(int)
	if !ok {
		t.Fatalf("expected int state, got %T", res)
	}
	if val != WAITING {
		t.Errorf("valueOf(\"WAITING\") = %d; want %d", val, WAITING)
	}
}

func TestThreadStateValueOf_MissingArg(t *testing.T) {
	res := threadStateValueOf([]interface{}{})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_WrongTypeArg(t *testing.T) {
	res := threadStateValueOf([]interface{}{123})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_NullName(t *testing.T) {
	res := threadStateValueOf([]interface{}{object.Null})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_NotStringObject(t *testing.T) {
	// Create a non-string object
	obj := object.MakeEmptyObject()
	res := threadStateValueOf([]interface{}{obj})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadStateValueOf_NoMatch(t *testing.T) {
	nameObj := object.StringObjectFromGoString("BOGUS")
	res := threadStateValueOf([]interface{}{nameObj})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

// --- threadStateValues() tests ---

func TestThreadStateValues_ArrayContentAndOrder(t *testing.T) {
	// Ensure non-zero size so returned array object is meaningful
	threadStateInstances = make([]*object.Object, len(ThreadState))

	res := threadStateValues([]interface{}{})
	arrObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object array, got %T", res)
	}

	// Validate array field type and size
	valueField := arrObj.FieldTable["value"]
	if valueField.Ftype != types.RefArray+"Ljava/lang/Thread$State;" {
		t.Errorf("unexpected array Ftype: %q", valueField.Ftype)
	}

	// The implementation stores the sorted keys slice in Fvalue
	keys, ok := valueField.Fvalue.([]int)
	if !ok {
		t.Fatalf("expected []int keys in array value, got %T", valueField.Fvalue)
	}
	if len(keys) != len(ThreadState) {
		t.Fatalf("expected %d keys, got %d", len(ThreadState), len(keys))
	}
	// Verify ascending order 0..5 based on defined constants
	if !(keys[0] == NEW && keys[1] == RUNNABLE && keys[2] == BLOCKED && keys[3] == WAITING && keys[4] == TIMED_WAITING && keys[5] == TERMINATED) {
		t.Errorf("keys not in declaration order: %v", keys)
	}
}
