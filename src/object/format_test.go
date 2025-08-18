/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/types"
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

// === the following tests were geenerated by JetBrains Junie to cover gaps in testing

// Test DEBUGGING flag when enabled
func TestFmtHelperWithDebuggingEnabled(t *testing.T) {
	globals.InitGlobals("test")

	// Save original DEBUGGING state and restore after test
	originalDebugging := DEBUGGING
	defer func() { DEBUGGING = originalDebugging }()

	// Enable debugging
	DEBUGGING = true

	field := Field{Ftype: types.Int, Fvalue: int64(42)}
	result := fmtHelper(field, "TestClass", "testField")

	// The function should still work correctly with debugging enabled
	if result != "42" {
		t.Errorf("Expected '42', got '%s'", result)
	}
	// Note: The debug output goes to stdout, which is hard to capture in tests
}

// Test String handling with trailing newline characters
func TestFmtHelperStringWithNewline(t *testing.T) {
	globals.InitGlobals("test")

	// Test with JavaByte array ending with newline
	javaBytes := []types.JavaByte{'H', 'e', 'l', 'l', 'o', '\n'}
	field := Field{Ftype: types.StringClassRef, Fvalue: javaBytes}
	result := fmtHelper(field, "TestClass", "testField")
	expected := "\"Hello\""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with []byte ending with newline
	bytes := []byte{'W', 'o', 'r', 'l', 'd', '\n'}
	field = Field{Ftype: types.StringClassRef, Fvalue: bytes}
	result = fmtHelper(field, "TestClass", "testField")
	expected = "\"World\""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with *[]byte ending with newline
	bytesPtr := &[]byte{'T', 'e', 's', 't', '\n'}
	field = Field{Ftype: types.StringClassRef, Fvalue: bytesPtr}
	result = fmtHelper(field, "TestClass", "testField")
	expected = "\"Test\""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test String handling with empty arrays
func TestFmtHelperStringWithEmptyArrays(t *testing.T) {
	globals.InitGlobals("test")

	// Test with empty JavaByte array
	javaBytes := []types.JavaByte{}
	field := Field{Ftype: types.StringClassRef, Fvalue: javaBytes}
	result := fmtHelper(field, "TestClass", "testField")
	expected := "\"\""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with empty []byte
	bytes := []byte{}
	field = Field{Ftype: types.StringClassRef, Fvalue: bytes}
	result = fmtHelper(field, "TestClass", "testField")
	expected = "\"\""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test with empty *[]byte
	emptyBytes := []byte{}
	field = Field{Ftype: types.StringClassRef, Fvalue: &emptyBytes}
	result = fmtHelper(field, "TestClass", "testField")
	expected = "\"\""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test String handling with nil fvalue
func TestFmtHelperStringWithNilValue(t *testing.T) {
	globals.InitGlobals("test")

	field := Field{Ftype: types.StringClassRef, Fvalue: nil}
	result := fmtHelper(field, "TestClass", "testField")
	expected := "<nil>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test static boolean handling
func TestFmtHelperStaticBoolean(t *testing.T) {
	globals.InitGlobals("test")
	globals.GetGlobalRef().FuncThrowException = _formatCycleKiller

	// Test static boolean field
	field := Field{Ftype: types.Static + types.Bool, Fvalue: true}
	result := fmtHelper(field, "TestClass", "testField")

	// Should contain "static" in the result
	if !strings.Contains(result, "static") {
		t.Errorf("Expected result to contain 'static', got '%s'", result)
	}
}

// Test boolean with int64 values
func TestFmtHelperBooleanInt64Values(t *testing.T) {
	globals.InitGlobals("test")

	// Test boolean with int64 value of 0 (should be false)
	field := Field{Ftype: types.Bool, Fvalue: int64(0)}
	result := fmtHelper(field, "TestClass", "testField")
	expected := "false"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test boolean with int64 value non-zero (should be true)
	field = Field{Ftype: types.Bool, Fvalue: int64(42)}
	result = fmtHelper(field, "TestClass", "testField")
	expected = "true"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test boolean with unexpected type
func TestFmtHelperBooleanUnexpectedType(t *testing.T) {
	globals.InitGlobals("test")

	// Test boolean with unexpected type (string)
	field := Field{Ftype: types.Bool, Fvalue: "not a boolean"}
	result := fmtHelper(field, "TestClass", "testField")

	if !strings.Contains(result, "ERROR") || !strings.Contains(result, "unexpected Fvalue variable type") {
		t.Errorf("Expected error message about unexpected type, got '%s'", result)
	}
}

// Test static byte array handling
func TestFmtHelperStaticByteArray(t *testing.T) {
	globals.InitGlobals("test")
	globals.GetGlobalRef().FuncThrowException = _formatCycleKiller

	field := Field{Ftype: types.Static + types.ByteArray, Fvalue: []byte{1, 2, 3}}
	result := fmtHelper(field, "TestClass", "testField")

	// Should contain "static" in the result
	if !strings.Contains(result, "static") {
		t.Errorf("Expected result to contain 'static', got '%s'", result)
	}
}

// Test byte array with nil fvalue
func TestFmtHelperByteArrayNilValue(t *testing.T) {
	globals.InitGlobals("test")

	field := Field{Ftype: types.ByteArray, Fvalue: nil}
	result := fmtHelper(field, "TestClass", "testField")
	expected := "<ERROR nil Fvalue>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test byte array with embedded *Object
func TestFmtHelperByteArrayWithEmbeddedObject(t *testing.T) {
	globals.InitGlobals("test")

	embeddedObj := MakeEmptyObject()
	field := Field{Ftype: types.ByteArray, Fvalue: embeddedObj}
	result := fmtHelper(field, "TestClass", "testField")
	expected := "*** embedded object ***"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test FormatField with invalid KlassName
func TestFormatFieldInvalidKlassName(t *testing.T) {
	globals.InitGlobals("test")

	obj := MakeEmptyObject()
	obj.KlassName = types.InvalidStringIndex
	obj.FieldTable["testField"] = Field{Ftype: types.Int, Fvalue: 42}

	result := obj.FormatField("testField")
	expected := "<ERROR nil class pointer>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test FormatField with empty fieldName and DEBUGGING enabled
func TestFormatFieldEmptyFieldNameWithDebugging(t *testing.T) {
	globals.InitGlobals("test")

	// Save original DEBUGGING state and restore after test
	originalDebugging := DEBUGGING
	defer func() { DEBUGGING = originalDebugging }()

	// Enable debugging
	DEBUGGING = true

	obj := MakeEmptyObject()
	klassType := "TestClass"
	obj.KlassName = stringPool.GetStringIndex(&klassType)

	// Call with empty fieldName - should trigger debugging path
	result := obj.FormatField("")

	// Should return the class name
	if result != klassType {
		t.Errorf("Expected '%s', got '%s'", klassType, result)
	}
}

// Test FormatField with non-empty FieldTable but empty fieldName and DEBUGGING
func TestFormatFieldNonEmptyTableEmptyFieldNameWithDebugging(t *testing.T) {
	globals.InitGlobals("test")

	// Save original DEBUGGING state and restore after test
	originalDebugging := DEBUGGING
	defer func() { DEBUGGING = originalDebugging }()

	// Enable debugging
	DEBUGGING = true

	obj := MakeEmptyObject()
	klassType := "TestClass"
	obj.KlassName = stringPool.GetStringIndex(&klassType)
	obj.FieldTable["existingField"] = Field{Ftype: types.Int, Fvalue: 42}

	// Call with empty fieldName but non-empty FieldTable
	result := obj.FormatField("")

	// Should return the class name
	if result != klassType {
		t.Errorf("Expected '%s', got '%s'", klassType, result)
	}
}

// Test FormatField with empty field table and DEBUGGING enabled
func TestFormatFieldEmptyTableWithDebugging(t *testing.T) {
	globals.InitGlobals("test")

	// Save original DEBUGGING state and restore after test
	originalDebugging := DEBUGGING
	defer func() { DEBUGGING = originalDebugging }()

	// Enable debugging
	DEBUGGING = true

	obj := MakeEmptyObject()
	klassType := "TestClass"
	obj.KlassName = stringPool.GetStringIndex(&klassType)

	// Call with empty field table
	result := obj.FormatField("")

	// Should return the class name
	if result != klassType {
		t.Errorf("Expected '%s', got '%s'", klassType, result)
	}
}

// Test DumpObject with missing KlassName
func TestDumpObjectMissingKlassName(t *testing.T) {
	globals.InitGlobals("test")

	obj := MakeEmptyObject()
	obj.KlassName = types.InvalidStringIndex

	// Capture output by redirecting stdout (or just test that it doesn't panic)
	obj.DumpObject("Test with missing class name", 0)
	// If we reach here without panicking, the test passes
}

// Test DumpObject with field table lookup error simulation
func TestDumpObjectFieldTableLookupError(t *testing.T) {
	globals.InitGlobals("test")

	obj := MakeEmptyObject()
	klassType := "TestClass"
	obj.KlassName = stringPool.GetStringIndex(&klassType)

	// Add a field
	obj.FieldTable["testField"] = Field{Ftype: types.Int, Fvalue: 42}

	// This tests the normal case, but the error case (lines 252-253) is hard to trigger
	// since Go's map implementation is consistent. The error case would only occur
	// in race conditions or if the map is modified during iteration.
	obj.DumpObject("Test field table", 2)
	// If we reach here without panicking, the test passes
}

// Test DumpObject with various indent levels
func TestDumpObjectWithIndents(t *testing.T) {
	globals.InitGlobals("test")

	obj := MakeEmptyObject()
	klassType := "TestClass"
	obj.KlassName = stringPool.GetStringIndex(&klassType)
	obj.FieldTable["testField"] = Field{Ftype: types.Int, Fvalue: 42}

	// Test with different indent levels
	obj.DumpObject("Test indent 0", 0)
	obj.DumpObject("Test indent 4", 4)
	obj.DumpObject("Test indent 8", 8)
	// If we reach here without panicking, the test passes
}

// === end of tests generated by JetBrains Junie
