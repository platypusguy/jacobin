/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/object"
	"os"
	"testing"
)

func TestPathsGet(t *testing.T) {
	testSep := string(os.PathSeparator)

	// Test 1: single string
	first := object.StringObjectFromGoString("a")
	res := pathsGet([]interface{}{first, nil}).(*object.Object)
	val := res.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val) != "a" {
		t.Errorf("Expected 'a', got %s", object.GoStringFromStringObject(val))
	}

	// Test 2: multiple strings
	moreArr := []*object.Object{
		object.StringObjectFromGoString("b"),
		object.StringObjectFromGoString("c"),
	}
	moreObj := object.MakeArrayFromRawArray(moreArr)
	res2 := pathsGet([]interface{}{first, moreObj}).(*object.Object)
	val2 := res2.FieldTable["value"].Fvalue.(*object.Object)
	expected2 := fmt.Sprintf("a%sb%sc", testSep, testSep)
	if object.GoStringFromStringObject(val2) != expected2 {
		t.Errorf("Expected '%s', got %s", expected2, object.GoStringFromStringObject(val2))
	}

	// Test 3: first with separator
	firstWithSep := object.StringObjectFromGoString("a" + testSep)
	res3 := pathsGet([]interface{}{firstWithSep, moreObj}).(*object.Object)
	val3 := res3.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val3) != expected2 {
		t.Errorf("Expected '%s', got %s", expected2, object.GoStringFromStringObject(val3))
	}

	// Test 4: empty more
	emptyMoreObj := object.MakeArrayFromRawArray([]*object.Object{})
	res4 := pathsGet([]interface{}{first, emptyMoreObj}).(*object.Object)
	val4 := res4.FieldTable["value"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(val4) != "a" {
		t.Errorf("Expected 'a', got %s", object.GoStringFromStringObject(val4))
	}
}
