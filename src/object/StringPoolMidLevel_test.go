/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/globals"
	"jacobin/stringPool"
	"testing"
)

// Test MakeEmptyStringObject
// This test does not use the String Pool
func TestStringMidLevel_1(t *testing.T) {
	_ = globals.InitGlobals("test") // Initialize the String Pool

	objPtr := MakeEmptyStringObject()
	if *(stringPool.GetStringPointer(objPtr.KlassName)) != StringClassName {
		t.Errorf("Expected Klass to be %s but observed: %s",
			StringClassName, *(stringPool.GetStringPointer(objPtr.KlassName)))
	}
	sz := len(objPtr.FieldTable)
	if sz > 0 {
		t.Errorf("Expected FieldTable size 0 but observed: %d", sz)
	}
}

// Test NewPoolStringFromGoString and GetGoStringFromObject
func TestStringMidLevel_2(t *testing.T) {
	_ = globals.InitGlobals("test") // Initialize the String Pool

	str1 := "The rain in Spain falls mainly on the plain"
	objPtr := NewPoolStringFromGoString(str1)
	objPtr.DumpObject("TestStringMidLevel_2", 0)
	str2 := GetGoStringFromObject(objPtr)
	if str1 != str2 {
		t.Errorf("Expected GetGoStringFromObject to return %s but observed: %s", str1, str2)
	}
}
