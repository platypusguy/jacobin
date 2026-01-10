/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"testing"
)

func TestCopyOfObjectPointers_TooFewArguments(t *testing.T) {
	result := *(utilArraysCopyOf([]interface{}{}).(*ghelpers.GErrBlk))
	if result.ExceptionType != excNames.IllegalArgumentException || result.ErrMsg != "utilArraysCopyOf: too few arguments" {
		t.Errorf("Expected IllegalArgumentException for too few arguments")
	}
}

func TestCopyOfObjectPointers_NullArray(t *testing.T) {
	result := *(utilArraysCopyOf([]interface{}{nil, int64(5)}).(*ghelpers.GErrBlk))
	if result.ExceptionType != excNames.NullPointerException || result.ErrMsg != "utilArraysCopyOf: null array argument" {
		t.Errorf("Expected NullPointerException for null array argument")
	}
}

func TestCopyOfObjectPointers_NegativeLength(t *testing.T) {
	obj := &object.Object{}
	result := *(utilArraysCopyOf([]interface{}{obj, int64(-1)}).(*ghelpers.GErrBlk))
	if result.ExceptionType != excNames.NegativeArraySizeException || result.ErrMsg != "utilArraysCopyOf: negative array length" {
		// if result != ghelpers.GetGErrBlk(excNames.NegativeArraySizeException, "copyOf: negative array length") {
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
	result := utilArraysCopyOf([]interface{}{oldArray, int64(4)})
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

func TestArraysAsList(t *testing.T) {
	globals.InitGlobals("test")
	stringPool.PreloadArrayClassesToStringPool()

	// 1. Create a mock array object [Ljava/lang/Object;
	arrayObj := object.Make1DimRefArray("Ljava/lang/Object;", 3)
	rawArray := arrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")
	rawArray[0] = s1
	rawArray[1] = s2
	rawArray[2] = nil // Arrays.asList allows nulls

	// 2. Call utilArraysAsList
	res := utilArraysAsList([]interface{}{arrayObj})
	listObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}

	// 3. Verify it's an ArrayList
	if className := object.GoStringFromStringPoolIndex(listObj.KlassName); className != "java/util/ArrayList" {
		t.Errorf("expected className java/util/ArrayList, got %s", className)
	}

	// 4. Verify elements
	al, err := GetArrayListFromObject(listObj)
	if err != nil {
		t.Fatalf("getArrayListFromObject failed: %v", err)
	}

	if len(al) != 3 {
		t.Errorf("expected size 3, got %d", len(al))
	}

	if al[0] != s1 || al[1] != s2 || al[2] != object.Null {
		t.Errorf("elements mismatch")
	}
}
