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
	"jacobin/src/stringPool"
	"testing"
)

func TestCopyOfObjectPointers_TooFewArguments(t *testing.T) {
	result := *(copyOfObjectPointers([]interface{}{}).(*GErrBlk))
	if result.ExceptionType != excNames.IllegalArgumentException || result.ErrMsg != "copyOfObjectPointers: too few arguments" {
		t.Errorf("Expected IllegalArgumentException for too few arguments")
	}
}

func TestCopyOfObjectPointers_NullArray(t *testing.T) {
	result := *(copyOfObjectPointers([]interface{}{nil, int64(5)}).(*GErrBlk))
	if result.ExceptionType != excNames.NullPointerException || result.ErrMsg != "copyOfObjectPointers: null array argument" {
		t.Errorf("Expected NullPointerException for null array argument")
	}
}

func TestCopyOfObjectPointers_NegativeLength(t *testing.T) {
	obj := &object.Object{}
	result := *(copyOfObjectPointers([]interface{}{obj, int64(-1)}).(*GErrBlk))
	if result.ExceptionType != excNames.NegativeArraySizeException || result.ErrMsg != "copyOfObjectPointers: negative array length" {
		// if result != getGErrBlk(excNames.NegativeArraySizeException, "copyOf: negative array length") {
		t.Errorf("Expected NegativeArraySizeException for negative array length")
	}
}

func TestCopyOfObjectPointers_CopyArray(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool() // needed for creation of array

	// Create a mock array object
	oldArray := object.Make1DimRefArray("Ljava/lang/Object;", 2)
	rawOldArray := oldArray.FieldTable["value"].Fvalue.([]*object.Object)
	rawOldArray[0] = object.StringObjectFromGoString("foo")
	rawOldArray[1] = object.StringObjectFromGoString("bar")

	// Test copying to a larger array
	result := copyOfObjectPointers([]interface{}{oldArray, int64(4)})
	newArray := result.(*object.Object).FieldTable["value"].Fvalue.([]*object.Object)
	if len(newArray) != 4 {
		t.Errorf("Expected new array length of 4, got %d", len(newArray))
	}

	if len(newArray) != 4 {
		t.Errorf("Expected new array length of 4, got %d", len(newArray))
	}

	if object.GoStringFromStringObject(newArray[0]) != "foo" || object.GoStringFromStringObject(newArray[1]) != "bar" {
		t.Errorf("Array elements not copied correctly")
	}
}
