/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/object"
	"strings"
	"testing"
)

func TestStringClinit(t *testing.T) {
	classloader.InitMethodArea()
	retval := stringClinit(nil)
	if retval == nil {
		t.Error("TestStringClinit: Was expecting an error message, but got none.")
	}
	switch retval.(type) {
	case *GErrBlk:
		gErr := retval.(*GErrBlk)
		if !strings.Contains(gErr.ErrMsg, "TestStringClinit: Could not find java/lang/String") {
			t.Errorf("TestStringClinit: Unexpected error message. got %s", gErr.ErrMsg)
		}
		if gErr.ExceptionType != exceptions.ClassNotLoadedException {
			t.Errorf("TestStringClinit: Unexpected exception type. got %d", gErr.ExceptionType)
		}
	default:
		t.Errorf("TestStringClinit: Did not get expected error message, got %v", retval)
	}
}
func TestStringToUpperCase(t *testing.T) {
	originalString := "hello"
	originalObj := object.NewPoolStringFromGoString(originalString)
	params := []interface{}{originalObj}
	ucObj := toUpperCase(params)
	strUpper := object.GetGoStringFromObject(ucObj.(*object.Object))
	expValue := "HELLO"
	if string(strUpper) != expValue {
		t.Errorf("TestStringToUpperCase failed, expected: %s, observed: %s", expValue, strUpper)
	}
}
