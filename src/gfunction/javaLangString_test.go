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
		t.Error("Was expecting an error message, but got none.")
	}
	switch retval.(type) {
	case *GErrBlk:
		gErr := retval.(*GErrBlk)
		if !strings.Contains(gErr.ErrMsg, "stringClinit: Could not find java/lang/String") {
			t.Errorf("Unexpected error message. got %s", gErr.ErrMsg)
		}
		if gErr.ExceptionType != exceptions.VirtualMachineError {
			t.Errorf("Unexpected exception type. got %d", gErr.ExceptionType)
		}
	default:
		t.Errorf("Did not get expected error message, got %v", retval)
	}
}
func TestStringToUpperCase(t *testing.T) {
	originalString := "hello"
	s := object.CreateCompactStringFromGoString(&originalString)
	params := []interface{}{s}
	s2 := toUpperCase(params)
	sUpper := object.GetGoStringFromJavaStringPtr(s2.(*object.Object))
	expValue := "HELLO"
	if string(sUpper) != expValue {
		t.Errorf("string toUpperCase failed, expected: %s, observed: %s", expValue, sUpper)
	}
}
