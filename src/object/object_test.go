/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"fmt"
	"testing"
)

func TestObjectToString1(t *testing.T) {
	obj := MakeEmptyObject()
	klassType := "java\\lang\\madeUpClass"
	obj.Klass = &klassType
	myIntField := Field{
		Ftype:  "I",
		Fvalue: 42,
	}
	obj.FieldTable["myInt"] = &myIntField

	myByteField := Field{
		Ftype:  "B",
		Fvalue: 0x61,
	}
	obj.FieldTable["myByte"] = &myByteField

	myStringField := Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: "Hello, Richard",
	}
	obj.FieldTable["myString"] = &myStringField

	str := obj.ToString()
	if len(str) == 0 {
		t.Errorf("empty string for object.ToString()")
	} else {
		fmt.Println(str)
	}
}

func TestObjectToString2(t *testing.T) {
	literal := "Hello, Jacobin!"
	str := CreateCompactStringFromGoString(&literal)

	retStr := str.ToString()
	if len(retStr) == 0 {
		t.Errorf("empty string for object.ToString()")
	} else {
		fmt.Println(retStr)
	}
}
