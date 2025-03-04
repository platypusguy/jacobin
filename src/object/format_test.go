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
	"strings"
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

// used instead of throwing an exception, which creates a circularity problem
func _formatCycleKiller(_ int, _ string) bool {
	return true
}

func TestFmtHelper(t *testing.T) {
	globals.InitGlobals("test")
	globals.GetGlobalRef().FuncThrowException = _formatCycleKiller
	className := filepath.FromSlash("java/lang/madeUpClass")
	fieldName := "Fred"

	// String.
	javaBytes := []types.JavaByte{'A', 'B', 'C'}
	fld := Field{types.StringClassRef, javaBytes}
	str := fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "ABC") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"ABC\", got \"%v\"", str)
	}

	// Static String.
	fld = Field{types.Static + types.StringClassRef, javaBytes}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "ABC") && !strings.Contains(str, "static") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"ABC\", got \"%v\"", str)
	}

	// String, unexpectedly using []byte.
	bites := []byte{'A', 'B', 'C'}
	fld = Field{types.StringClassRef, bites}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "ABC") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"ABC\", got \"%v\"", str)
	}

	// String, unexpectedly using *[]byte.
	fld = Field{types.StringClassRef, &bites}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "ABC") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"ABC\", got \"%v\"", str)
	}

	// String, unexpectedly using string.
	fld = Field{types.StringClassRef, "ABC"}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "ABC") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"ABC\", got \"%v\"", str)
	}

	// String, unexpectedly using string.
	// debug: stringPool.DumpStringPool("format_test.go-TestFmtHelper")
	fld = Field{types.StringIndex, types.StringPoolStringIndex}
	str = fmtHelper(fld, className, fieldName)
	if str != types.StringClassName {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"%s\", got \"%v\"", types.StringClassName, str)
	}

	// Primitive integer.
	fld = Field{types.Int, int64(42)}
	str = fmtHelper(fld, className, fieldName)
	if str != "42" {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"42\", got \"%v\"", str)
	}

	// Primitive static long.
	fld = Field{types.StaticLong, int64(42)}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "42") && !strings.Contains(str, "static") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return STATIC \"42\", got \"%v\"", str)
	}

	// Primitive boolean.
	fld = Field{types.Bool, types.JavaBoolTrue}
	str = fmtHelper(fld, className, fieldName)
	if str != "true" {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"true\", got \"%v\"", str)
	}

	// Primitive boolean, unexpected Go bool.
	fld = Field{types.Bool, true}
	str = fmtHelper(fld, className, fieldName)
	if str != "true" {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"true\", got \"%v\"", str)
	}

	// Primitive byte array.
	bites = []byte{'A', 'B', 'C'}
	fld = Field{types.ByteArray, bites}
	str = fmtHelper(fld, className, fieldName)
	if str != "41 42 43" {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"41 42 43\", got \"%v\"", str)
	}

	// Primitive byte array pointer.
	bites = []byte{'A', 'B', 'C'}
	fld = Field{types.ByteArray, &bites}
	str = fmtHelper(fld, className, fieldName)
	if str != "41 42 43" {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"41 42 43\", got \"%v\"", str)
	}

	// Primitive byte array, wrong format.
	jb := []types.JavaByte{1, 2, 3}
	fld = Field{types.ByteArray, jb}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "but value is of type") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"but value is of type\", got \"%v\"", str)
	}

	// Primitive byte array, length=0.
	bites = []byte{}
	fld = Field{types.ByteArray, bites}
	str = fmtHelper(fld, className, fieldName)
	if !strings.Contains(str, "array of zero") {
		t.Errorf("TestFmtHelper: expected fmtHelper to return \"array of zero\", got \"%v\"", str)
	}

}
