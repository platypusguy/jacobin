/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"jacobin/globals"
	"jacobin/stringPool"
	"jacobin/types"
	"path/filepath"
	"testing"
)

func TestDumpObjectFieldTable(t *testing.T) {
	t.Log("Test Object.FieldTable DumpObject processing")

	globals.InitGlobals("test")
	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.KlassName = stringPool.GetStringIndex(&klassType)

	myFloatField := Field{
		Ftype:  types.Float,
		Fvalue: 1.0,
	}
	obj.FieldTable["myFloat"] = myFloatField

	myDoubleField := Field{
		Ftype:  types.Double,
		Fvalue: 2.0,
	}
	obj.FieldTable["myDouble"] = myDoubleField

	myIntField := Field{
		Ftype:  types.Int,
		Fvalue: 42,
	}
	obj.FieldTable["myInt"] = myIntField

	myLongField := Field{
		Ftype:  types.Long,
		Fvalue: 42,
	}
	obj.FieldTable["myLong"] = myLongField

	myShortField := Field{
		Ftype:  types.Short,
		Fvalue: 42,
	}
	obj.FieldTable["myShort"] = myShortField

	myByteField := Field{
		Ftype:  types.Byte,
		Fvalue: 0x61,
	}
	obj.FieldTable["myByte"] = myByteField

	myStaticTrueField := Field{
		Ftype:  types.Static + types.Bool,
		Fvalue: true,
	}
	obj.FieldTable["myStaticTrue"] = myStaticTrueField

	myFalseField := Field{
		Ftype:  types.Bool,
		Fvalue: false,
	}
	obj.FieldTable["myFalse"] = myFalseField

	myCharField := Field{
		Ftype:  types.Char,
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

func TestFormatField(t *testing.T) {
	t.Log("Test field slice DumpObject processing")

	globals.InitGlobals("test")
	obj := MakeEmptyObject()
	klassType := filepath.FromSlash("java/lang/madeUpClass")
	obj.KlassName = stringPool.GetStringIndex(&klassType)

	myFloatField := Field{
		Ftype:  types.Float,
		Fvalue: 1.0,
	}
	obj.FieldTable["myFloat"] = myFloatField

	myDoubleField := Field{
		Ftype:  types.Double,
		Fvalue: 2.0,
	}
	obj.FieldTable["myDouble"] = myDoubleField

	myIntField := Field{
		Ftype:  types.Int,
		Fvalue: 42,
	}
	obj.FieldTable["myInt"] = myIntField

	myLongField := Field{
		Ftype:  types.Long,
		Fvalue: 42,
	}
	obj.FieldTable["myLong"] = myLongField

	myShortField := Field{
		Ftype:  types.Short,
		Fvalue: 42,
	}
	obj.FieldTable["myShort"] = myShortField

	myByteField := Field{
		Ftype:  types.Byte,
		Fvalue: 0x61,
	}
	obj.FieldTable["myByte"] = myByteField

	myStaticTrueField := Field{
		Ftype:  types.Static + types.Bool,
		Fvalue: true,
	}
	obj.FieldTable["myStaticTrue"] = myStaticTrueField

	myFalseField := Field{
		Ftype:  types.Bool,
		Fvalue: false,
	}
	obj.FieldTable["myFalse"] = myFalseField

	myCharField := Field{
		Ftype:  types.Char,
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
