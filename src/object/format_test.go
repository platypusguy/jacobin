/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"path/filepath"
	"testing"
)

func TestDumpObjectFieldTable(t *testing.T) {
	t.Log("Test Object.FieldTable DumpObject processing")
	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.Klass = &klassType

	myFloatField := Field{
		Ftype:  "F",
		Fvalue: 1.0,
	}
	obj.FieldTable["myFloat"] = myFloatField

	myDoubleField := Field{
		Ftype:  "D",
		Fvalue: 2.0,
	}
	obj.FieldTable["myDouble"] = myDoubleField

	myIntField := Field{
		Ftype:  "I",
		Fvalue: 42,
	}
	obj.FieldTable["myInt"] = myIntField

	myLongField := Field{
		Ftype:  "J",
		Fvalue: 42,
	}
	obj.FieldTable["myLong"] = myLongField

	myShortField := Field{
		Ftype:  "S",
		Fvalue: 42,
	}
	obj.FieldTable["myShort"] = myShortField

	myByteField := Field{
		Ftype:  "B",
		Fvalue: 0x61,
	}
	obj.FieldTable["myByte"] = myByteField

	myStaticTrueField := Field{
		Ftype:  "XZ",
		Fvalue: true,
	}
	obj.FieldTable["myStaticTrue"] = myStaticTrueField

	myFalseField := Field{
		Ftype:  "Z",
		Fvalue: false,
	}
	obj.FieldTable["myFalse"] = myFalseField

	myCharField := Field{
		Ftype:  "C",
		Fvalue: 'C',
	}
	obj.FieldTable["myChar"] = myCharField

	myStringField := Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: "Hello, Unka Andoo !",
	}
	obj.FieldTable["myString"] = myStringField

	obj.DumpObject(klassType, 3)
}

// Test field slice toString processing
func TestDumpObjectFieldSlice(t *testing.T) {
	t.Log("Test Object.Fields slice DumpObject processing")
	literal := "This is a compact string from a Go string"
	csObj := CreateCompactStringFromGoString(&literal)
	csObj.DumpObject(literal, 0)

	// Create a custom object.
	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.Klass = &klassType

	// Now, dump the same string as a byte array.
	csObj.Klass = &klassType
	csObj.DumpObject("[B string]", 0)

	myFloatField := Field{
		Ftype:  "F",
		Fvalue: 1.0,
	}
	obj.Fields = append(obj.Fields, myFloatField)
	obj.DumpObject("F", 0)

	myDoubleField := Field{
		Ftype:  "D",
		Fvalue: 2.0,
	}
	obj.Fields[0] = myDoubleField
	obj.DumpObject("D", 0)

	myIntField := Field{
		Ftype:  "I",
		Fvalue: 42,
	}
	obj.Fields[0] = myIntField
	obj.DumpObject("I", 0)

	myLongField := Field{
		Ftype:  "J",
		Fvalue: 42,
	}
	obj.Fields[0] = myLongField
	obj.DumpObject("J", 0)

	myShortField := Field{
		Ftype:  "S",
		Fvalue: 42,
	}
	obj.Fields[0] = myShortField
	obj.DumpObject("S", 0)

	myByteField := Field{
		Ftype:  "B",
		Fvalue: 0x61,
	}
	obj.Fields[0] = myByteField
	obj.DumpObject("B", 0)

	myFalseField := Field{
		Ftype:  "Z",
		Fvalue: false,
	}
	obj.Fields[0] = myFalseField
	obj.DumpObject("false Z", 0)

	myStaticTrueField := Field{
		Ftype:  "XZ",
		Fvalue: true,
	}
	obj.Fields[0] = myStaticTrueField
	obj.DumpObject("true XZ", 0)

	myCharField := Field{
		Ftype:  "C",
		Fvalue: 'C',
	}
	obj.Fields[0] = myCharField
	obj.DumpObject("C", 0)

	myCharField = Field{
		Ftype:  "XC",
		Fvalue: 'D',
	}
	obj.Fields[0] = myCharField
	obj.DumpObject("XC", 0)

}

func TestFormatField(t *testing.T) {
	t.Log("Test field slice DumpObject processing")

	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.Klass = &klassType

	myFloatField := Field{
		Ftype:  "F",
		Fvalue: 1.0,
	}
	obj.FieldTable["myFloat"] = myFloatField

	myDoubleField := Field{
		Ftype:  "D",
		Fvalue: 2.0,
	}
	obj.FieldTable["myDouble"] = myDoubleField

	myIntField := Field{
		Ftype:  "I",
		Fvalue: 42,
	}
	obj.FieldTable["myInt"] = myIntField

	myLongField := Field{
		Ftype:  "J",
		Fvalue: 42,
	}
	obj.FieldTable["myLong"] = myLongField

	myShortField := Field{
		Ftype:  "S",
		Fvalue: 42,
	}
	obj.FieldTable["myShort"] = myShortField

	myByteField := Field{
		Ftype:  "B",
		Fvalue: 0x61,
	}
	obj.FieldTable["myByte"] = myByteField

	myStaticTrueField := Field{
		Ftype:  "XZ",
		Fvalue: true,
	}
	obj.FieldTable["myStaticTrue"] = myStaticTrueField

	myFalseField := Field{
		Ftype:  "Z",
		Fvalue: false,
	}
	obj.FieldTable["myFalse"] = myFalseField

	myCharField := Field{
		Ftype:  "C",
		Fvalue: 'C',
	}
	obj.FieldTable["myChar"] = myCharField

	myStringField1 := Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: "Hello, Unka Andoo !",
	}
	obj.FieldTable["myString"] = myStringField1

	t.Log("NOTE: Key \"Fred\" will be diagnosed as missing:")
	str := obj.FormatField("Fred")
	t.Log(str)

	t.Log("NOTE: Will add a key \"value\" field.")
	myStringField2 := Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: "Hello, Unka Andoo !",
	}
	obj.FieldTable["Fred"] = myStringField2

	t.Log("Will try FormatField again.")
	str = obj.FormatField("Fred")
	t.Log(str)

}
