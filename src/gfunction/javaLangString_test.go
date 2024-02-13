/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/object"
	"testing"
)

func TestStringToUpperCase(t *testing.T) {
	originalString := "hello"
	s := object.CreateCompactStringFromGoString(&originalString)
	params := []interface{}{s}
	s2 := toUpperCase(params)
	// sUpper := s2.(*object.Object).Fields[0].Fvalue.(*[]byte)
	sUpper := object.GetGoStringFromJavaStringPtr(s2.(*object.Object))
	if string(sUpper) != "HELLO" {
		t.Errorf("string ToUpperCase failed, got %s", sUpper)
	}
}
