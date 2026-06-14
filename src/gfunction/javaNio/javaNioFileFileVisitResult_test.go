/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func Test_FileVisitResult_Enum(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_FileVisitResult()

	// Test clinit
	fvResultClinit(nil)

	// Test values()
	res := fvResultValues(nil)
	arr, ok := res.(*object.Object)
	if !ok || arr == nil {
		t.Fatalf("values() should return array object")
	}
	vals := arr.FieldTable["value"].Fvalue.([]*object.Object)
	if len(vals) != 4 {
		t.Fatalf("expected 4 FileVisitResult values, got %d", len(vals))
	}

	expectedNames := []string{"CONTINUE", "TERMINATE", "SKIP_SUBTREE", "SKIP_SIBLINGS"}
	for i, name := range expectedNames {
		nmObj := vals[i].FieldTable["name"].Fvalue.(*object.Object)
		if object.GoStringFromStringObject(nmObj) != name {
			t.Errorf("expected %s at index %d, got %s", name, i, object.GoStringFromStringObject(nmObj))
		}
		ord := vals[i].FieldTable["ordinal"].Fvalue.(int64)
		if ord != int64(i) {
			t.Errorf("expected ordinal %d for %s, got %d", i, name, ord)
		}
	}

	// Test valueOf(String) success
	for _, name := range expectedNames {
		sObj := object.StringObjectFromGoString(name)
		v := fvResultValueOfString([]interface{}{sObj})
		if v == nil || object.IsNull(v) {
			t.Fatalf("valueOf(%s) returned null", name)
		}
		resObj := v.(*object.Object)
		nmObj := resObj.FieldTable["name"].Fvalue.(*object.Object)
		if object.GoStringFromStringObject(nmObj) != name {
			t.Errorf("valueOf(%s) returned wrong constant: %s", name, object.GoStringFromStringObject(nmObj))
		}
	}

	// Test valueOf(String) error cases
	// 1. Missing argument
	err1 := fvResultValueOfString([]interface{}{})
	if geb, ok := err1.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for missing arg, got %T", err1)
	}

	// 2. Null argument
	err2 := fvResultValueOfString([]interface{}{object.Null})
	if geb, ok := err2.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for null arg, got %T", err2)
	}

	// 3. Not a String
	err3 := fvResultValueOfString([]interface{}{object.MakeEmptyObjectWithClassName(new(string))})
	if geb, ok := err3.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for non-string arg, got %T", err3)
	}

	// 4. Invalid name
	err4 := fvResultValueOfString([]interface{}{object.StringObjectFromGoString("INVALID")})
	if geb, ok := err4.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for invalid name, got %T", err4)
	}
}
