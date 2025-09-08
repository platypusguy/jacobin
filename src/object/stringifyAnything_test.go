/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"container/list"
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"math/big"
	"strings"
	"testing"
)

var byteName = "java/lang/Byte"
var boolName = "java/lang/Boolean"
var charName = "java/lang/Character"
var doubleName = "java/lang/Double"
var floatName = "java/lang/Float"
var intName = "java/lang/Integer"
var longName = "java/lang/Long"
var shortName = "java/lang/Short"
var stringName = "java/lang/String"

// array literals
var ba = "[B"
var da = "[D"
var fa = "[F"
var ia = "[I"
var la = "[J"
var sa = "[S"
var za = "[Z" // boolean array

func TestStringifyAnythingGo_NilArgument(t *testing.T) {
	result := StringifyAnythingGo(nil)
	if result != types.NullString {
		t.Errorf("Expected %s for nil argument, got %s", types.NullString, result)
	}
}

func TestStringifyAnythingGo_NullObject(t *testing.T) {
	globals.InitGlobals("test")
	result := StringifyAnythingGo(Null)
	if result != types.NullString {
		t.Errorf("Expected %s for null object, got %s", types.NullString, result)
	}
}

func TestStringifyAnythingGo_StringObject(t *testing.T) {
	globals.InitGlobals("test")
	// Create a String object with value field
	obj := MakeEmptyObjectWithClassName(&stringName)
	javaBytes := JavaByteArrayFromGoString("hello")
	obj.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: javaBytes}

	result := StringifyAnythingGo(obj)
	if result != "hello" {
		t.Errorf("Expected 'hello', got %s", result)
	}
}

func TestStringifyAnythingGo_StringObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&stringName)
	// Don't add value field

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted String object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_BooleanObject_True(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&boolName)
	obj.FieldTable["value"] = Field{Ftype: types.Bool, Fvalue: types.JavaBoolTrue}

	result := StringifyAnythingGo(obj)
	if result != "true" {
		t.Errorf("Expected 'true', got %s", result)
	}
}

func TestStringifyAnythingGo_BooleanObject_False(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&boolName)
	obj.FieldTable["value"] = Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}

	result := StringifyAnythingGo(obj)
	if result != "false" {
		t.Errorf("Expected 'false', got %s", result)
	}
}

func TestStringifyAnythingGo_BooleanObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&boolName)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted Boolean object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_ByteObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&byteName)
	obj.FieldTable["value"] = Field{Ftype: types.Byte, Fvalue: int64(0x42)}

	result := StringifyAnythingGo(obj)
	if result != "0x42" {
		t.Errorf("Expected '0x42', got %s", result)
	}
}

func TestStringifyAnythingGo_ByteObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&byteName)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted Byte object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_CharacterObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&charName)
	obj.FieldTable["value"] = Field{Ftype: types.Char, Fvalue: int64(65)}

	result := StringifyAnythingGo(obj)
	if result != "65" {
		t.Errorf("Expected '65', got %s", result)
	}
}

func TestStringifyAnythingGo_CharacterObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&charName)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted Character object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_DoubleObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&doubleName)
	obj.FieldTable["value"] = Field{Ftype: types.Double, Fvalue: 42.5}

	result := StringifyAnythingGo(obj)
	if result != "42.5" {
		t.Errorf("Expected '42.5', got %s", result)
	}
}

func TestStringifyAnythingGo_DoubleObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")

	obj := MakeEmptyObjectWithClassName(&doubleName)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted Double object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_FloatObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&floatName)
	obj.FieldTable["value"] = Field{Ftype: types.Float, Fvalue: 42.5}

	result := StringifyAnythingGo(obj)
	if result != "42.5" {
		t.Errorf("Expected '42.5', got %s", result)
	}
}

func TestStringifyAnythingGo_FloatObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&floatName)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted Float object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_IntegerObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&intName)
	obj.FieldTable["value"] = Field{Ftype: types.Int, Fvalue: int64(42)}

	result := StringifyAnythingGo(obj)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_LongObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&longName)
	obj.FieldTable["value"] = Field{Ftype: types.Long, Fvalue: int64(123456789)}

	result := StringifyAnythingGo(obj)
	if result != "123456789" {
		t.Errorf("Expected '123456789', got %s", result)
	}
}

func TestStringifyAnythingGo_ShortObject(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&shortName)
	obj.FieldTable["value"] = Field{Ftype: types.Short, Fvalue: int64(42)}

	result := StringifyAnythingGo(obj)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_IntegerObject_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&intName)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted Integer/Long/Short object, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_ByteArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&ba)
	bytes := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello"
	obj.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: bytes}

	result := StringifyAnythingGo(obj)
	if result != "0x48656c6c6f" {
		t.Errorf("Expected '0x48656c6c6f', got %s", result)
	}
}

func TestStringifyAnythingGo_ByteArray_JavaBytes(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&ba)
	javaBytes := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	obj.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: javaBytes}

	result := StringifyAnythingGo(obj)
	if result != "0x48656c6c6f" {
		t.Errorf("Expected '0x48656c6c6f', got %s", result)
	}
}

func TestStringifyAnythingGo_ByteArray_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&ba)

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted byte array, missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_ByteArray_CorruptedValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&ba)
	obj.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: "invalid"} // Wrong type

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted byte array \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_BoolArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&za)
	boolValues := []int64{1, 0, 1, 0}
	obj.FieldTable["value"] = Field{Ftype: types.BoolArray, Fvalue: boolValues}

	result := StringifyAnythingGo(obj)
	if result != "[true, false, true, false]" {
		t.Errorf("Expected '[true, false, true, false]', got %s", result)
	}
}

func TestStringifyAnythingGo_BoolArray_MissingValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&za)
	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: boolean object missing \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_BoolArray_CorruptedValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&za)
	obj.FieldTable["value"] = Field{Ftype: types.BoolArray, Fvalue: "invalid"}

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: corrupted boolean array \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_DoubleArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&da)
	doubleValues := []float64{1.5, 2.7, 3.14}
	obj.FieldTable["value"] = Field{Ftype: types.DoubleArray, Fvalue: doubleValues}

	result := StringifyAnythingGo(obj)
	if result != "[1.5, 2.7, 3.14]" {
		t.Errorf("Expected '[1.5, 2.7, 3.14]', got %s", result)
	}
}

func TestStringifyAnythingGo_DoubleArray_CorruptedValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&da)
	obj.FieldTable["value"] = Field{Ftype: types.DoubleArray, Fvalue: "invalid"}

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: double array missing \"value\" field or array value corrupted"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_FloatArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&fa)
	floatValues := []float64{1.5, 2.7}
	obj.FieldTable["value"] = Field{Ftype: types.FloatArray, Fvalue: floatValues}

	result := StringifyAnythingGo(obj)
	if result != "[1.5, 2.7]" {
		t.Errorf("Expected '[1.5, 2.7]', got %s", result)
	}
}

func TestStringifyAnythingGo_FloatArray_CorruptedValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&fa)
	obj.FieldTable["value"] = Field{Ftype: types.FloatArray, Fvalue: "invalid"}

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: float array missing \"value\" field or array value corrupted"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_IntArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&ia)
	intValues := []int64{1, 2, 3, 42}
	obj.FieldTable["value"] = Field{Ftype: types.IntArray, Fvalue: intValues}

	result := StringifyAnythingGo(obj)
	if result != "[1, 2, 3, 42]" {
		t.Errorf("Expected '[1, 2, 3, 42]', got %s", result)
	}
}

func TestStringifyAnythingGo_LongArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&la)
	longValues := []int64{100, 200, 300}
	obj.FieldTable["value"] = Field{Ftype: types.LongArray, Fvalue: longValues}

	result := StringifyAnythingGo(obj)
	if result != "[100, 200, 300]" {
		t.Errorf("Expected '[100, 200, 300]', got %s", result)
	}
}

func TestStringifyAnythingGo_ShortArray(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&sa)
	shortValues := []int64{10, 20, 30}
	obj.FieldTable["value"] = Field{Ftype: types.ShortArray, Fvalue: shortValues}

	result := StringifyAnythingGo(obj)
	if result != "[10, 20, 30]" {
		t.Errorf("Expected '[10, 20, 30]', got %s", result)
	}
}

func TestStringifyAnythingGo_IntArray_CorruptedValue(t *testing.T) {
	globals.InitGlobals("test")
	obj := MakeEmptyObjectWithClassName(&ia)
	obj.FieldTable["value"] = Field{Ftype: types.IntArray, Fvalue: "invalid"}

	result := StringifyAnythingGo(obj)
	expected := "StringifyAnythingGo: int/long/short array missing \"value\" field or array value corrupted"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_DefaultObjectCase(t *testing.T) {
	globals.InitGlobals("test")
	arrList := "java/util/ArrayList"
	obj := MakeEmptyObjectWithClassName(&arrList)
	obj.FieldTable["size"] = Field{Ftype: types.Int, Fvalue: int64(5)}
	obj.FieldTable["capacity"] = Field{Ftype: types.Int, Fvalue: int64(10)}

	result := StringifyAnythingGo(obj)

	// Should contain class name and field values
	if !strings.Contains(result, "ArrayList{") || !strings.Contains(result, "size=5") || !strings.Contains(result, "capacity=10") {
		t.Errorf("Expected format with class name and fields, got %s", result)
	}
}

func TestStringifyAnythingGo_Field_StringClassRef(t *testing.T) {
	javaBytes := JavaByteArrayFromGoString("test")
	field := Field{Ftype: types.StringClassRef, Fvalue: javaBytes}

	result := StringifyAnythingGo(field)
	if result != "test" {
		t.Errorf("Expected 'test', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Byte_Int64(t *testing.T) {
	field := Field{Ftype: types.Byte, Fvalue: int64(0x42)}

	result := StringifyAnythingGo(field)
	if result != "0x42" {
		t.Errorf("Expected '0x42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Byte_Byte(t *testing.T) {
	field := Field{Ftype: types.Byte, Fvalue: byte(0x42)}

	result := StringifyAnythingGo(field)
	if result != "0x42" {
		t.Errorf("Expected '0x42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Byte_CorruptedValue(t *testing.T) {
	field := Field{Ftype: types.Byte, Fvalue: "invalid"}

	result := StringifyAnythingGo(field)
	expected := "StringifyAnythingGo Field types.Byte: corrupted byte \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_Field_ByteArray_JavaBytes(t *testing.T) {
	javaBytes := []types.JavaByte{0x48, 0x65}
	field := Field{Ftype: types.ByteArray, Fvalue: javaBytes}

	result := StringifyAnythingGo(field)
	if result != "0x4865" {
		t.Errorf("Expected '0x4865', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_ByteArray_GoBytes(t *testing.T) {
	bytes := []byte{0x48, 0x65}
	field := Field{Ftype: types.ByteArray, Fvalue: bytes}

	result := StringifyAnythingGo(field)
	if result != "0x4865" {
		t.Errorf("Expected '0x4865', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_ByteArray_CorruptedValue(t *testing.T) {
	field := Field{Ftype: types.ByteArray, Fvalue: "invalid"}

	result := StringifyAnythingGo(field)
	expected := "StringifyAnythingGo Field types.ByteArray: corrupted byte array \"value\" field"
	if result != expected {
		t.Errorf("Expected error message, got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Bool_True(t *testing.T) {
	field := Field{Ftype: types.Bool, Fvalue: types.JavaBoolTrue}

	result := StringifyAnythingGo(field)
	if result != "true" {
		t.Errorf("Expected 'true', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Bool_False(t *testing.T) {
	field := Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}

	result := StringifyAnythingGo(field)
	if result != "false" {
		t.Errorf("Expected 'false', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Int_Int64(t *testing.T) {
	field := Field{Ftype: types.Int, Fvalue: int64(42)}

	result := StringifyAnythingGo(field)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Int_Uint64(t *testing.T) {
	field := Field{Ftype: types.Int, Fvalue: uint64(42)}

	result := StringifyAnythingGo(field)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Int_Int32(t *testing.T) {
	field := Field{Ftype: types.Int, Fvalue: int32(42)}

	result := StringifyAnythingGo(field)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Int_Uint32(t *testing.T) {
	field := Field{Ftype: types.Int, Fvalue: uint32(42)}

	result := StringifyAnythingGo(field)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Short_Int16(t *testing.T) {
	field := Field{Ftype: types.Short, Fvalue: int16(42)}

	result := StringifyAnythingGo(field)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Short_Uint16(t *testing.T) {
	field := Field{Ftype: types.Short, Fvalue: uint16(42)}

	result := StringifyAnythingGo(field)
	if result != "42" {
		t.Errorf("Expected '42', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Long(t *testing.T) {
	field := Field{Ftype: types.Long, Fvalue: int64(123456789)}

	result := StringifyAnythingGo(field)
	if result != "123456789" {
		t.Errorf("Expected '123456789', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Int_UnrecognizedType(t *testing.T) {
	field := Field{Ftype: types.Int, Fvalue: "invalid"}

	result := StringifyAnythingGo(field)
	if !strings.Contains(result, "StringifyAnythingGo  Field types.Int: unrecognized field value type") {
		t.Errorf("Expected error message about unrecognized type, got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Double(t *testing.T) {
	field := Field{Ftype: types.Double, Fvalue: 42.5}

	result := StringifyAnythingGo(field)
	if result != "42.5" {
		t.Errorf("Expected '42.5', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_Float(t *testing.T) {
	field := Field{Ftype: types.Float, Fvalue: 42.5}

	result := StringifyAnythingGo(field)
	if result != "42.5" {
		t.Errorf("Expected '42.5', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_BigInteger(t *testing.T) {
	bigInt := big.NewInt(123456789)
	field := Field{Ftype: types.BigInteger, Fvalue: bigInt}

	result := StringifyAnythingGo(field)
	if result != "123456789" {
		t.Errorf("Expected '123456789', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_LinkedList(t *testing.T) {
	globals.InitGlobals("test")
	ll := list.New()
	ll.PushBack(StringObjectFromGoString("first"))
	ll.PushBack(StringObjectFromGoString("second"))
	field := Field{Ftype: types.LinkedList, Fvalue: ll}

	result := StringifyAnythingGo(field)
	if result != "[first, second]" {
		t.Errorf("Expected '[first, second]', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_LinkedList_Empty(t *testing.T) {
	ll := list.New()
	field := Field{Ftype: types.LinkedList, Fvalue: ll}

	result := StringifyAnythingGo(field)
	if result != "[]" { // Empty list
		t.Errorf("Expected '[]', got %s", result)
	}
}

func TestStringifyAnythingGo_Field_DefaultCase(t *testing.T) {
	field := Field{Ftype: "UNKNOWN", Fvalue: "test"}

	result := StringifyAnythingGo(field)
	if !strings.Contains(result, "StringifyAnythingGo Field default: unrecognized argument type") {
		t.Errorf("Expected error message about unrecognized argument type, got %s", result)
	}
}

func TestStringifyAnythingGo_UnrecognizedArgumentType(t *testing.T) {
	result := StringifyAnythingGo("invalid string argument")

	if !strings.Contains(result, "StringifyAnythingGo: neither *Object nor Field") {
		t.Errorf("Expected error message about unrecognized argument, got %s", result)
	}
}

func TestStringifyAnythingJava(t *testing.T) {
	globals.InitGlobals("test")
	result := StringifyAnythingJava("test input")

	// Should return a String object
	if result == nil {
		t.Errorf("Expected String object, got nil")
		return
	}

	// Verify it's a String object
	klassName := stringPool.GetStringPointer(result.KlassName)
	if *klassName != "java/lang/String" {
		t.Errorf("Expected java/lang/String, got %s", *klassName)
	}

	// Verify the content by stringifying it back
	content := StringifyAnythingGo(result)
	expectedContent := StringifyAnythingGo("test input")
	if content != expectedContent {
		t.Errorf("Expected %s, got %s", expectedContent, content)
	}
}

func TestStringifyAnythingJava_NilInput(t *testing.T) {
	globals.InitGlobals("test")
	result := StringifyAnythingJava(nil)

	// Should return a String object containing types.NullString
	if result == nil {
		t.Errorf("Expected String object, got nil")
		return
	}

	content := StringifyAnythingGo(result)
	if content != types.NullString {
		t.Errorf("Expected %s, got %s", types.NullString, content)
	}
}
