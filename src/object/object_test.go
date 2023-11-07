/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"path/filepath"
	"testing"
)

func TestObjectToString1(t *testing.T) {
	t.Log("Test field table toString processing")
	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.Klass = &klassType

	myFloatField := Field{
		Ftype:  "F",
		Fvalue: 1.0,
	}
	obj.FieldTable["myFloat"] = &myFloatField

	myDoubleField := Field{
		Ftype:  "D",
		Fvalue: 2.0,
	}
	obj.FieldTable["myDouble"] = &myDoubleField

	myIntField := Field{
		Ftype:  "I",
		Fvalue: 42,
	}
	obj.FieldTable["myInt"] = &myIntField

	myLongField := Field{
		Ftype:  "J",
		Fvalue: 42,
	}
	obj.FieldTable["myLong"] = &myLongField

	myShortField := Field{
		Ftype:  "S",
		Fvalue: 42,
	}
	obj.FieldTable["myShort"] = &myShortField

	myByteField := Field{
		Ftype:  "B",
		Fvalue: 0x61,
	}
	obj.FieldTable["myByte"] = &myByteField

	myStaticTrueField := Field{
		Ftype:  "XZ",
		Fvalue: true,
	}
	obj.FieldTable["myStaticTrue"] = &myStaticTrueField

	myFalseField := Field{
		Ftype:  "Z",
		Fvalue: false,
	}
	obj.FieldTable["myFalse"] = &myFalseField

	myCharField := Field{
		Ftype:  "C",
		Fvalue: 'C',
	}
	obj.FieldTable["myChar"] = &myCharField

	myStringField := Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: "Hello, Unka Andoo !",
	}
	obj.FieldTable["myString"] = &myStringField

	str := obj.ToString(42)
	if len(str) == 0 {
		t.Errorf("empty string for object.ToString()")
	} else {
		t.Log(str)
	}
}

// Test field slice toString processing
func TestObjectToString2(t *testing.T) {
	t.Log("Test field slice toString processing")
	literal := "This is a compact string from a Go string"
	csObj := CreateCompactStringFromGoString(&literal)
	retStr := csObj.ToString(0)
	if len(retStr) == 0 {
		t.Errorf("empty string for object.ToString()")
	} else {
		t.Log(retStr)
	}

	// Create a custom object.
	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.Klass = &klassType

	// Now, dump the same string as a byte array.
	csObj.Klass = &klassType
	retStr = csObj.ToString(0)
	if len(retStr) == 0 {
		t.Errorf("empty string for object.ToString()")
	} else {
		t.Log(retStr)
	}

	myFloatField := Field{
		Ftype:  "F",
		Fvalue: 1.0,
	}
	obj.Fields = append(obj.Fields, myFloatField)
	t.Log(obj.ToString(0))

	myDoubleField := Field{
		Ftype:  "D",
		Fvalue: 2.0,
	}
	obj.Fields[0] = myDoubleField
	t.Log(obj.ToString(0))

	myIntField := Field{
		Ftype:  "I",
		Fvalue: 42,
	}
	obj.Fields[0] = myIntField
	t.Log(obj.ToString(0))

	myLongField := Field{
		Ftype:  "J",
		Fvalue: 42,
	}
	obj.Fields[0] = myLongField
	t.Log(obj.ToString(0))

	myShortField := Field{
		Ftype:  "S",
		Fvalue: 42,
	}
	obj.Fields[0] = myShortField
	t.Log(obj.ToString(0))

	myByteField := Field{
		Ftype:  "B",
		Fvalue: 0x61,
	}
	obj.Fields[0] = myByteField
	t.Log(obj.ToString(0))

	myStaticTrueField := Field{
		Ftype:  "XZ",
		Fvalue: true,
	}
	obj.Fields[0] = myStaticTrueField
	t.Log(obj.ToString(0))

	myFalseField := Field{
		Ftype:  "Z",
		Fvalue: false,
	}
	obj.Fields[0] = myFalseField
	t.Log(obj.ToString(0))

	myCharField := Field{
		Ftype:  "C",
		Fvalue: 'C',
	}
	obj.Fields[0] = myCharField
	t.Log(obj.ToString(0))

}
