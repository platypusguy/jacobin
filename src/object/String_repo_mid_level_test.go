/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"testing"
)

// Test MakeEmptyStringObject
func TestStringMidLevel_1(t *testing.T) {
	objPtr := MakeEmptyStringObject()
	if *(objPtr.Klass) != StringClassName {
		t.Errorf("Expected Klass to be %s but observed: %s", StringClassName, *(objPtr.Klass))
	}
	sz := len(objPtr.FieldTable)
	if sz > 0 {
		t.Errorf("Expected FieldTable size 0 but observed: %d", sz)
	}
}

// Test NewRepoStringFromGoString and GetGoStringFromObject
func TestStringMidLevel_2(t *testing.T) {
	str1 := "The rain in Spain falls mainly on the plain"
	objPtr := NewRepoStringFromGoString(str1)
	objPtr.DumpObject("TestNewRepoStringFromGoString", 0)
	str2 := GetGoStringFromObject(objPtr)
	if str1 != str2 {
		t.Errorf("Expected GetGoStringFromObject to return %s but observed: %s", str1, str2)
	}
}
