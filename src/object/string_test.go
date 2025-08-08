/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"fmt"
	"io"
	"jacobin/globals"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"math"
	"os"
	"strings"
	"testing"
)

func TestNewStringObject(t *testing.T) {
	globals.InitGlobals("test")

	str := *NewStringObject()
	klassStr := *(stringPool.GetStringPointer(str.KlassName))
	if klassStr != types.StringClassName {
		t.Errorf("Klass should be java/lang/String, observed: %s", klassStr)
	}

	value := str.FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(value) != 0 {
		t.Errorf("value field should be an empty byte, observed length of %d", len(value))
	}

	coder := str.FieldTable["coder"].Fvalue.(types.JavaByte)
	if coder != 0 && coder != 1 {
		t.Errorf("coder field should be 0 or 1, observed: %d", coder)
	}

	hash := str.FieldTable["hash"].Fvalue.(uint32)
	if hash != uint32(0) {
		t.Errorf("hash field should be 0, observed: %d", hash)
	}

	hashIsZero := str.FieldTable["hashIsZero"].Fvalue.(types.JavaByte)
	if hashIsZero != 0 {
		t.Errorf("hashIsZero field should be false(0), observed: %d", hashIsZero)
	}
}

func TestStringObjectFromGoString(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."

	strObj := StringObjectFromGoString(constStr)
	strValue := GoStringFromStringObject(strObj)
	if strValue != constStr {
		t.Errorf("expected string value to be '%s', observed: '%s'", constStr, strValue)
	}
}

func TestGoStringFromInvalidStringObject(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()
	className := "Oregano"
	obj := MakeEmptyObjectWithClassName(&className)
	strValue := GoStringFromStringObject(obj)
	if strValue != "" {
		t.Errorf("expected empty string , observed: '%s'", strValue)
	}
}

func TestGoStringFromStringNullObject(t *testing.T) {
	str := GoStringFromStringObject(nil)
	if str != "" {
		t.Errorf("expected empty string, observed: '%s'", str)
	}
}

func TestGoStringFromInvalidObject(t *testing.T) {
	obj := MakeEmptyObject()
	obj.FieldTable["value"] = Field{Ftype: types.Int, Fvalue: 42}
	str := GoStringFromStringObject(obj)
	if str != "" {
		t.Errorf("expected empty string, observed: '%s'", str)
	}
}

func TestGoStringFromGoString(t *testing.T) {
	obj := MakeEmptyObject()
	obj.FieldTable["value"] = Field{Ftype: types.GolangString, Fvalue: "Hello"}
	str := GoStringFromStringObject(obj)
	if str != "Hello" {
		t.Errorf("expected string 'Hello', observed: '%s'", str)
	}
}

func TestJavaByteArrayFromStringObjectValie(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."
	constBytes := JavaByteArrayFromGoString(constStr)

	strObj := StringObjectFromGoString(constStr)
	bb := JavaByteArrayFromStringObject(strObj)
	if !JavaByteArrayEquals(bb, constBytes) {
		t.Errorf("expected string value to be '%s', observed: '%s'", constStr, GoStringFromJavaByteArray(bb))
	}
}

func TestJavaByteArrayFromStringObjectInvalid(t *testing.T) {
	str := StringObjectFromGoString("ABC")
	ba := ByteArrayFromStringObject(str)
	if ba[0] != types.JavaByte('A') && ba[1] != types.JavaByte('B') && ba[2] != types.JavaByte('C') {
		t.Errorf("expected 'ABC', observed: %s", string(GoStringFromJavaByteArray(ba)))
	}
}

func TestStringPoolIndexFromStringObjectWithByteArray(t *testing.T) {
	stringPool.EmptyStringPool()
	stringPool.PreloadArrayClassesToStringPool()
	ba := []byte("Hello")
	obj := StringObjectFromByteArray(ba)
	index := StringPoolIndexFromStringObject(obj)
	if *(stringPool.GetStringPointer(index)) != "Hello" {
		t.Errorf("expected 'Hello', observed: %s", *(stringPool.GetStringPointer(index)))
	}
}

func TestByteArrayFromStringObjectInvalid(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."

	strObj := StringObjectFromGoString(constStr)
	strObj.KlassName = uint32(200) // the invalid part; KlassName is not java/lang/String
	bb := ByteArrayFromStringObject(strObj)
	if bb != nil {
		t.Errorf("expected nil return b/c of error, observed: %v", bb)
	}
}

func TestStringObjectFromByteArray(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."
	constBytes := []byte(constStr)

	strObj := StringObjectFromByteArray(constBytes)
	strValue := GoStringFromStringObject(strObj)
	if strValue != constStr {
		t.Errorf("expected string value to be '%s', observed: '%s'", constStr, strValue)
	}
}

func TestStringPoolStringOperations(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."
	strObj := StringObjectFromGoString(constStr)
	index := StringPoolIndexFromStringObject(strObj)
	if index == types.InvalidStringIndex {
		t.Errorf("string pool index is types.InvalidStringIndex")
		return
	}

	strValue := GoStringFromStringPoolIndex(index)
	if strValue == types.EmptyString { // if ""
		t.Errorf("strValue from string pool index %d is an empty string", index)
	}

	strObj = StringObjectFromPoolIndex(index)
	if strObj == nil {
		t.Errorf("strObj from string pool index %d is nil", index)
	}

	index2 := StringPoolIndexFromStringObject(strObj)
	if index2 != index {
		t.Errorf("string pool index=%d but index2=%d (expected equality)", index, index2)
		return
	}

	bb := JavaByteArrayFromStringPoolIndex(index)
	if bb == nil {
		t.Errorf("byte array from string pool index %d is nil", index)
	}
}

func TestStringPoolStringIndexFromStringObjectInvalid(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."
	strObj := StringObjectFromGoString(constStr)
	strObj.KlassName = uint32(200) // the invalid part: KlassName is not java/lang/String
	index := StringPoolIndexFromStringObject(strObj)
	if index != types.InvalidStringIndex {
		t.Errorf("Expected types.InvalidStringIndex, got %d", index)
		return
	}
}

func TestByteArrayFromStringPoolIndexInvalid(t *testing.T) {
	index := math.MaxInt32 // use a string pool index that will always be too big
	byteArray := JavaByteArrayFromStringPoolIndex(uint32(index))
	if byteArray != nil {
		t.Errorf("expected nil due to error, got %v", byteArray)
	}
}

func TestUpdateStringObjectFromBytes(t *testing.T) {
	constStr := "Mary had a little lamb whose fleece was white as snow."
	constBytes := JavaByteArrayFromGoString(constStr)
	strObj := StringObjectFromGoString("To be updated")
	if !IsStringObject(strObj) {
		t.Errorf("expected IsStringObject(valid string object) to be true, observed false")
	}
	UpdateValueFieldFromJavaBytes(strObj, constBytes)
	strValue := GoStringFromStringObject(strObj)
	if strValue != constStr {
		t.Errorf("strValue from updated string object has wrong value: %s", strValue)
	}
}

func TestIsStringObjectValid(t *testing.T) {
	constStr := "Mary had a little lamb whose fleece was white as snow."
	strObj := StringObjectFromGoString(constStr)
	if !IsStringObject(strObj) {
		t.Errorf("expected IsStringObject(valid string object) to be true, observed false")
	}
}

func TestIsStringObjectWithNil(t *testing.T) {
	if IsStringObject(nil) {
		t.Errorf("expected IsStringObject(nil) to be false, observed true")
	}
}

func TestIsStringObjectWithGoString(t *testing.T) {
	if IsStringObject("go string") {
		t.Errorf("expected IsStringObject(go string) to be false, observed true")
	}
}

func TestGetStringFromStringPoolIndex(t *testing.T) {
	globals.InitGlobals("test")
	goStr := GoStringFromStringPoolIndex(types.StringPoolStringIndex)
	if goStr != types.StringClassName {
		t.Errorf("Got unexpected string value: %s", goStr)
	}

	goStr = GoStringFromStringPoolIndex(200_000)
	if goStr != "" {
		t.Errorf("Expected empty string, got %s", goStr)
	}
}

func TestGetStringObjectFromStringPoolIndex(t *testing.T) {
	globals.InitGlobals("test")
	stObj := StringObjectFromPoolIndex(types.StringPoolStringIndex)
	goStr := GoStringFromStringObject(stObj)
	if goStr != types.StringClassName {
		t.Errorf("Got unexpected string value: %s", goStr)
	}

	stObj = StringObjectFromPoolIndex(200_000)
	if stObj != nil {
		t.Errorf("Expected nil, got %s", GoStringFromStringObject(stObj))
	}
}

func TestIsStringObject(t *testing.T) {
	globals.InitGlobals("test")

	strObj := StringObjectFromGoString("test object")
	if !IsStringObject(strObj) {
		t.Errorf("expected IsStringObject(valid string object) to be true, got false")
	}

	emptyObj := Make1DimArray(BYTE, 10)
	if IsStringObject(emptyObj) {
		t.Errorf("expected IsStringObject(emptyObj) to be false, got true")
	}
}

func TestObjectFieldToStringInvalidFieldName(t *testing.T) {
	globals.InitGlobals("test")

	obj := StringObjectFromGoString("a little lamb")
	ret := ObjectFieldToString(obj, "non-existentField")

	if ret != types.NullString {
		t.Errorf("Expected types.NullString, got %s", ret)
	}

}

func TestObjectFieldToStringForIntArray(t *testing.T) {
	globals.InitGlobals("test")

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	obj := Make1DimArray(INT, 10)
	ret := ObjectFieldToString(obj, "value")

	if !strings.Contains(ret, " 0") {
		t.Errorf("Expected different return string, got %s", ret)
	}

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if msg != "" {
		t.Errorf("Expected no error message, got %s", msg)
	}
}

func TestObjectFieldToStringForUnknownType(t *testing.T) {
	globals.InitGlobals("test")
	globals.GetGlobalRef().FuncThrowException = func(i int, s string) bool {
		fmt.Fprintf(os.Stderr, "Exception thrown in TestObjectFieldToStringForUnknownType: %s\n", s)
		return false
	}

	obj := StringObjectFromGoString("allo!")
	obj.FieldTable["testFieldName"] = Field{
		Ftype:  "testFieldType",
		Fvalue: nil,
	}
	ret := ObjectFieldToString(obj, "testFieldName")

	if !strings.Contains(ret, "java/lang/String") {
		t.Errorf("Expected different return string, got %s", ret)
	}

	if !strings.Contains(ret, "testFieldName") { // just the field name
		t.Errorf("Looking for field name, got %s", ret)
	}
	if !strings.Contains(ret, "testFieldType") { // just the field name
		t.Errorf("Looking for field type, got %s", ret)
	}
}

func TestObjectFieldToStringForFileHandle(t *testing.T) {
	globals.InitGlobals("test")

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	obj := StringObjectFromGoString("allo!")
	obj.FieldTable["testField"] = Field{
		Ftype:  types.FileHandle,
		Fvalue: nil,
	}
	ret := ObjectFieldToString(obj, "testField")

	if !strings.Contains(ret, "FileHandle") {
		t.Errorf("Expected different return string, got %s", ret)
	}

	// restore stderr to what it was before
	_ = w.Close()
	os.Stderr = normalStderr
}

func TestObjectFieldToStringForStaticBool(t *testing.T) {
	globals.InitGlobals("test")

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	obj := StringObjectFromGoString("allo!")
	obj.FieldTable["testField"] = Field{
		Ftype:  "XZ",
		Fvalue: int64(1),
	}
	ret := ObjectFieldToString(obj, "testField")

	if !strings.Contains(ret, "true") {
		t.Errorf("Expected different return string, got %s", ret)
	}

	// restore stderr to what it was before
	_ = w.Close()
	os.Stderr = normalStderr
}

// === The following tests were generated by JetBrains Junie to cover testing gaps

// Test StringObjectArrayFromGoStringArray function (completely untested)
func TestStringObjectArrayFromGoStringArray(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	// Test empty array
	emptyArray := []string{}
	result := StringObjectArrayFromGoStringArray(emptyArray)
	if len(result) != 0 {
		t.Errorf("Expected empty array, got length %d", len(result))
	}

	// Test single element array
	singleArray := []string{"hello"}
	result = StringObjectArrayFromGoStringArray(singleArray)
	if len(result) != 1 {
		t.Errorf("Expected array length 1, got %d", len(result))
	}
	if !IsStringObject(result[0]) {
		t.Errorf("Expected element to be string object, got non-string object")
	}
	str := GoStringFromStringObject(result[0])
	if str != "hello" {
		t.Errorf("Expected 'hello', got '%s'", str)
	}

	// Test multiple elements array
	multiArray := []string{"hello", "world", "test"}
	result = StringObjectArrayFromGoStringArray(multiArray)
	if len(result) != 3 {
		t.Errorf("Expected array length 3, got %d", len(result))
	}
	for i, expected := range multiArray {
		if !IsStringObject(result[i]) {
			t.Errorf("Expected element %d to be string object", i)
		}
		actual := GoStringFromStringObject(result[i])
		if actual != expected {
			t.Errorf("Expected element %d to be '%s', got '%s'", i, expected, actual)
		}
	}

	// Test array with empty strings
	mixedArray := []string{"", "non-empty", ""}
	result = StringObjectArrayFromGoStringArray(mixedArray)
	if len(result) != 3 {
		t.Errorf("Expected array length 3, got %d", len(result))
	}
	for i, expected := range mixedArray {
		actual := GoStringFromStringObject(result[i])
		if actual != expected {
			t.Errorf("Expected element %d to be '%s', got '%s'", i, expected, actual)
		}
	}
}

// Test GoStringArrayFromStringObjectArray function (completely untested)
func TestGoStringArrayFromStringObjectArray(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	// Test empty array
	emptyArray := []*Object{}
	result := GoStringArrayFromStringObjectArray(emptyArray)
	if len(result) != 0 {
		t.Errorf("Expected empty array, got length %d", len(result))
	}

	// Test single element array
	singleObj := StringObjectFromGoString("hello")
	singleArray := []*Object{singleObj}
	result = GoStringArrayFromStringObjectArray(singleArray)
	if len(result) != 1 {
		t.Errorf("Expected array length 1, got %d", len(result))
	}
	if result[0] != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result[0])
	}

	// Test multiple elements array
	obj1 := StringObjectFromGoString("hello")
	obj2 := StringObjectFromGoString("world")
	obj3 := StringObjectFromGoString("test")
	multiArray := []*Object{obj1, obj2, obj3}
	result = GoStringArrayFromStringObjectArray(multiArray)
	expected := []string{"hello", "world", "test"}
	if len(result) != 3 {
		t.Errorf("Expected array length 3, got %d", len(result))
	}
	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Expected element %d to be '%s', got '%s'", i, exp, result[i])
		}
	}

	// Test array with nil object
	nilArray := []*Object{nil}
	result = GoStringArrayFromStringObjectArray(nilArray)
	if len(result) != 1 {
		t.Errorf("Expected array length 1, got %d", len(result))
	}
	if result[0] != "" {
		t.Errorf("Expected empty string for nil object, got '%s'", result[0])
	}

	// Test mixed array (valid objects and nil)
	mixedArray := []*Object{StringObjectFromGoString("valid"), nil, StringObjectFromGoString("")}
	result = GoStringArrayFromStringObjectArray(mixedArray)
	expectedMixed := []string{"valid", "", ""}
	if len(result) != 3 {
		t.Errorf("Expected array length 3, got %d", len(result))
	}
	for i, exp := range expectedMixed {
		if result[i] != exp {
			t.Errorf("Expected element %d to be '%s', got '%s'", i, exp, result[i])
		}
	}
}

// Test StringPoolIndexFromGoString function (completely untested)
func TestStringPoolIndexFromGoString(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	// Test normal string
	testStr := "test string"
	index := StringPoolIndexFromGoString(testStr)
	retrievedStr := GoStringFromStringPoolIndex(index)
	if retrievedStr != testStr {
		t.Errorf("Expected '%s', got '%s'", testStr, retrievedStr)
	}

	// Test empty string
	emptyStr := ""
	index = StringPoolIndexFromGoString(emptyStr)
	retrievedStr = GoStringFromStringPoolIndex(index)
	if retrievedStr != emptyStr {
		t.Errorf("Expected empty string, got '%s'", retrievedStr)
	}

	// Test string with special characters
	specialStr := "Hello\nWorld\t!"
	index = StringPoolIndexFromGoString(specialStr)
	retrievedStr = GoStringFromStringPoolIndex(index)
	if retrievedStr != specialStr {
		t.Errorf("Expected '%s', got '%s'", specialStr, retrievedStr)
	}

	/*
		// Test Unicode string -- not sure we are ready for unicode yet
		unicodeStr := "你好世界"
		index = StringPoolIndexFromGoString(unicodeStr)
		retrievedStr = GoStringFromStringPoolIndex(index)
		if retrievedStr != unicodeStr {
			t.Errorf("Expected '%s', got '%s'", unicodeStr, retrievedStr)
		}
	*/
}

// Test GoStringFromJavaCharArray function (completely untested)
func TestGoStringFromJavaCharArray(t *testing.T) {
	// Test empty array
	emptyArray := []int64{}
	result := GoStringFromJavaCharArray(emptyArray)
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}

	// Test single character
	singleChar := []int64{65} // ASCII 'A'
	result = GoStringFromJavaCharArray(singleChar)
	if result != "A" {
		t.Errorf("Expected 'A', got '%s'", result)
	}

	// Test multiple characters
	multiChars := []int64{72, 101, 108, 108, 111} // "Hello"
	result = GoStringFromJavaCharArray(multiChars)
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}

	// Test Unicode characters
	unicodeChars := []int64{20320, 22909} // 你好
	result = GoStringFromJavaCharArray(unicodeChars)
	expected := "你好"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test special characters
	specialChars := []int64{9, 10, 13} // tab, newline, carriage return
	result = GoStringFromJavaCharArray(specialChars)
	expected = "\t\n\r"
	if result != expected {
		t.Errorf("Expected tab+newline+CR, got %q", result)
	}

	// Test zero values
	zeroChars := []int64{0, 65, 0}
	result = GoStringFromJavaCharArray(zeroChars)
	expected = "\x00A\x00"
	if result != expected {
		t.Errorf("Expected null+A+null, got %q", result)
	}
}

// Test EqualStringObjects function (completely untested)
func TestEqualStringObjects(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	// Test equal strings
	str1 := StringObjectFromGoString("hello")
	str2 := StringObjectFromGoString("hello")
	if !EqualStringObjects(str1, str2) {
		t.Errorf("Expected equal strings to return true")
	}

	// Test different strings
	str3 := StringObjectFromGoString("world")
	if EqualStringObjects(str1, str3) {
		t.Errorf("Expected different strings to return false")
	}

	// Test empty strings
	empty1 := StringObjectFromGoString("")
	empty2 := StringObjectFromGoString("")
	if !EqualStringObjects(empty1, empty2) {
		t.Errorf("Expected equal empty strings to return true")
	}

	// Test first argument is not string object
	nonStringObj := MakeEmptyObject()
	if EqualStringObjects(nonStringObj, str1) {
		t.Errorf("Expected non-string object vs string object to return false")
	}

	// Test second argument is not string object
	if EqualStringObjects(str1, nonStringObj) {
		t.Errorf("Expected string object vs non-string object to return false")
	}

	// Test both arguments are not string objects
	nonStringObj2 := MakeEmptyObject()
	if EqualStringObjects(nonStringObj, nonStringObj2) {
		t.Errorf("Expected non-string object vs non-string object to return false")
	}

	// Test nil arguments
	if EqualStringObjects(nil, str1) {
		t.Errorf("Expected nil vs string object to return false")
	}

	if EqualStringObjects(str1, nil) {
		t.Errorf("Expected string object vs nil to return false")
	}

	if EqualStringObjects(nil, nil) {
		t.Errorf("Expected nil vs nil to return false")
	}

	// Test same object reference
	if !EqualStringObjects(str1, str1) {
		t.Errorf("Expected same object reference to return true")
	}
}

// Test ObjectFieldToString additional cases (many untested branches)
func TestObjectFieldToStringAdditionalCases(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()
	globals.InitStringPool()

	obj := MakeEmptyObject()

	// Test StringClassRef type
	obj.FieldTable["stringField"] = Field{
		Ftype:  types.StringClassRef,
		Fvalue: JavaByteArrayFromGoString("test string"),
	}
	result := ObjectFieldToString(obj, "stringField")
	if result != "test string" {
		t.Errorf("Expected 'test string', got '%s'", result)
	}

	// Test BigInteger type
	obj.FieldTable["bigIntField"] = Field{
		Ftype:  types.BigInteger,
		Fvalue: "12345678901234567890",
	}
	result = ObjectFieldToString(obj, "bigIntField")
	if result != "12345678901234567890" {
		t.Errorf("Expected big integer string, got '%s'", result)
	}

	// Test Bool type - false case
	obj.FieldTable["boolFalseField"] = Field{
		Ftype:  types.Bool,
		Fvalue: int64(0),
	}
	result = ObjectFieldToString(obj, "boolFalseField")
	if result != "false" {
		t.Errorf("Expected 'false', got '%s'", result)
	}

	// Test BoolArray type
	obj.FieldTable["boolArrayField"] = Field{
		Ftype:  types.BoolArray,
		Fvalue: []int64{1, 0, 1},
	}
	result = ObjectFieldToString(obj, "boolArrayField")
	if result != "true false true" {
		t.Errorf("Expected 'true false true', got '%s'", result)
	}

	// Test BoolArray with single element
	obj.FieldTable["singleBoolArray"] = Field{
		Ftype:  types.BoolArray,
		Fvalue: []int64{1},
	}
	result = ObjectFieldToString(obj, "singleBoolArray")
	if result != "true" {
		t.Errorf("Expected 'true', got '%s'", result)
	}

	// Test BoolArray empty
	obj.FieldTable["emptyBoolArray"] = Field{
		Ftype:  types.BoolArray,
		Fvalue: []int64{},
	}
	result = ObjectFieldToString(obj, "emptyBoolArray")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}

	// Test ByteArray with []byte
	obj.FieldTable["byteArrayField"] = Field{
		Ftype:  types.ByteArray,
		Fvalue: []byte{0x41, 0x42, 0x43}, // ABC
	}
	result = ObjectFieldToString(obj, "byteArrayField")
	if result != "414243" {
		t.Errorf("Expected '414243', got '%s'", result)
	}

	// Test "Ljava/lang/String;" with []byte
	obj.FieldTable["javaStringField"] = Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}, // Hello
	}
	result = ObjectFieldToString(obj, "javaStringField")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}

	// Test "Ljava/lang/String;" with []JavaByte
	obj.FieldTable["javaStringJBField"] = Field{
		Ftype:  "Ljava/lang/String;",
		Fvalue: JavaByteArrayFromGoString("World"),
	}
	result = ObjectFieldToString(obj, "javaStringJBField")
	if result != "World" {
		t.Errorf("Expected 'World', got '%s'", result)
	}

	// Test CharArray type
	obj.FieldTable["charArrayField"] = Field{
		Ftype:  types.CharArray,
		Fvalue: []int64{72, 101, 108, 108, 111}, // Hello
	}
	result = ObjectFieldToString(obj, "charArrayField")
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}

	// Test IntArray type
	obj.FieldTable["intArrayField"] = Field{
		Ftype:  types.IntArray,
		Fvalue: []int64{1, 2, 3},
	}
	result = ObjectFieldToString(obj, "intArrayField")
	if result != "1 2 3" {
		t.Errorf("Expected '1 2 3', got '%s'", result)
	}

	// Test LongArray type
	obj.FieldTable["longArrayField"] = Field{
		Ftype:  types.LongArray,
		Fvalue: []int64{100, 200},
	}
	result = ObjectFieldToString(obj, "longArrayField")
	if result != "100 200" {
		t.Errorf("Expected '100 200', got '%s'", result)
	}

	// Test ShortArray type
	obj.FieldTable["shortArrayField"] = Field{
		Ftype:  types.ShortArray,
		Fvalue: []int64{5, 10, 15},
	}
	result = ObjectFieldToString(obj, "shortArrayField")
	if result != "5 10 15" {
		t.Errorf("Expected '5 10 15', got '%s'", result)
	}

	// Test empty integer arrays
	obj.FieldTable["emptyIntArray"] = Field{
		Ftype:  types.IntArray,
		Fvalue: []int64{},
	}
	result = ObjectFieldToString(obj, "emptyIntArray")
	if result != "" {
		t.Errorf("Expected empty string for empty array, got '%s'", result)
	}

	// Test Double type
	obj.FieldTable["doubleField"] = Field{
		Ftype:  types.Double,
		Fvalue: float64(3.14159),
	}
	result = ObjectFieldToString(obj, "doubleField")
	if result != "3.14159" {
		t.Errorf("Expected '3.14159', got '%s'", result)
	}

	// Test Float type
	obj.FieldTable["floatField"] = Field{
		Ftype:  types.Float,
		Fvalue: float64(2.71828),
	}
	result = ObjectFieldToString(obj, "floatField")
	if result != "2.71828" {
		t.Errorf("Expected '2.71828', got '%s'", result)
	}

	// Test DoubleArray type
	obj.FieldTable["doubleArrayField"] = Field{
		Ftype:  types.DoubleArray,
		Fvalue: []float64{1.1, 2.2, 3.3},
	}
	result = ObjectFieldToString(obj, "doubleArrayField")
	if result != "1.1 2.2 3.3" {
		t.Errorf("Expected '1.1 2.2 3.3', got '%s'", result)
	}

	// Test FloatArray type
	obj.FieldTable["floatArrayField"] = Field{
		Ftype:  types.FloatArray,
		Fvalue: []float64{0.5, 1.5},
	}
	result = ObjectFieldToString(obj, "floatArrayField")
	if result != "0.5 1.5" {
		t.Errorf("Expected '0.5 1.5', got '%s'", result)
	}

	// Test empty float arrays
	obj.FieldTable["emptyFloatArray"] = Field{
		Ftype:  types.FloatArray,
		Fvalue: []float64{},
	}
	result = ObjectFieldToString(obj, "emptyFloatArray")
	if result != "" {
		t.Errorf("Expected empty string for empty float array, got '%s'", result)
	}

	// Test Ref type
	className := "[B"
	obj.KlassName = uint32(stringPool.GetStringIndex(&className))
	obj.FieldTable["refField"] = Field{
		Ftype:  types.Ref,
		Fvalue: "some reference",
	}
	result = ObjectFieldToString(obj, "refField")
	// This should return the class name from string pool
	if result != "[B" {
		t.Errorf("Expected '[B' for Ref type, got %s", result)
	}

	// Test RefArray type
	obj.FieldTable["refArrayField"] = Field{
		Ftype:  types.RefArray,
		Fvalue: "some ref array",
	}
	result = ObjectFieldToString(obj, "refArrayField")
	if result == "" {
		t.Errorf("Expected non-empty string for RefArray type, got empty string")
	}

	// Test "[Ljava/lang/Object;" type
	obj.FieldTable["objectArrayField"] = Field{
		Ftype:  "[Ljava/lang/Object;",
		Fvalue: "some object array",
	}
	result = ObjectFieldToString(obj, "objectArrayField")
	if result == "" {
		t.Errorf("Expected non-empty string for object array type, got empty string")
	}
}

// Test ObjectFieldToString with null object
func TestObjectFieldToStringWithNullObject(t *testing.T) {
	result := ObjectFieldToString(nil, "anyField")
	if result != types.NullString {
		t.Errorf("Expected types.NullString, got '%s'", result)
	}
}

// Test ByteArrayFromStringObject with nil object
func TestByteArrayFromStringObjectWithNil(t *testing.T) {
	result := ByteArrayFromStringObject(nil)
	if result != nil {
		t.Errorf("Expected nil for nil object, got %v", result)
	}
}

// Test ByteArrayFromStringObject with wrong KlassName
func TestByteArrayFromStringObjectWrongKlass(t *testing.T) {
	obj := MakeEmptyObject()
	obj.KlassName = uint32(999) // wrong class name
	result := ByteArrayFromStringObject(obj)
	if result != nil {
		t.Errorf("Expected nil for wrong class name, got %v", result)
	}
}

// Test StringPoolIndexFromStringObject with nil object
func TestStringPoolIndexFromStringObjectWithNil(t *testing.T) {
	result := StringPoolIndexFromStringObject(nil)
	if result != types.InvalidStringIndex {
		t.Errorf("Expected InvalidStringIndex for nil object, got %d", result)
	}
}

// Test StringObjectFromPoolIndex with invalid index
func TestStringObjectFromPoolIndexInvalid(t *testing.T) {
	globals.InitGlobals("test")

	// Test with very large index
	result := StringObjectFromPoolIndex(uint32(999999))
	if result != nil {
		t.Errorf("Expected nil for invalid index, got %v", result)
	}
}

// Test GoStringFromStringPoolIndex with invalid index
func TestGoStringFromStringPoolIndexInvalid(t *testing.T) {
	globals.InitGlobals("test")

	// Test with very large index
	result := GoStringFromStringPoolIndex(uint32(999999))
	if result != "" {
		t.Errorf("Expected empty string for invalid index, got '%s'", result)
	}
}

// Test additional type cases for individual primitive types in ObjectFieldToString
func TestObjectFieldToStringPrimitiveTypes(t *testing.T) {
	globals.InitGlobals("test")

	obj := MakeEmptyObject()

	// Test Byte type
	obj.FieldTable["byteField"] = Field{
		Ftype:  types.Byte,
		Fvalue: int64(127),
	}
	result := ObjectFieldToString(obj, "byteField")
	if result != "127" {
		t.Errorf("Expected '127', got '%s'", result)
	}

	// Test Char type
	obj.FieldTable["charField"] = Field{
		Ftype:  types.Char,
		Fvalue: int64(65), // 'A'
	}
	result = ObjectFieldToString(obj, "charField")
	if result != "65" {
		t.Errorf("Expected '65', got '%s'", result)
	}

	// Test Int type
	obj.FieldTable["intField"] = Field{
		Ftype:  types.Int,
		Fvalue: int64(42),
	}
	result = ObjectFieldToString(obj, "intField")
	if result != "42" {
		t.Errorf("Expected '42', got '%s'", result)
	}

	// Test Long type
	obj.FieldTable["longField"] = Field{
		Ftype:  types.Long,
		Fvalue: int64(9223372036854775807),
	}
	result = ObjectFieldToString(obj, "longField")
	if result != "9223372036854775807" {
		t.Errorf("Expected max long value, got '%s'", result)
	}

	// Test Rune type
	obj.FieldTable["runeField"] = Field{
		Ftype:  types.Rune,
		Fvalue: int64(8364), // € symbol
	}
	result = ObjectFieldToString(obj, "runeField")
	if result != "8364" {
		t.Errorf("Expected '8364', got '%s'", result)
	}

	// Test Short type
	obj.FieldTable["shortField"] = Field{
		Ftype:  types.Short,
		Fvalue: int64(-32768),
	}
	result = ObjectFieldToString(obj, "shortField")
	if result != "-32768" {
		t.Errorf("Expected '-32768', got '%s'", result)
	}
}

// end of Junie-generated tests
