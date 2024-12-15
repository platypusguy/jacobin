/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
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

	coder := str.FieldTable["coder"].Fvalue.(byte)
	if coder != 0 && coder != 1 {
		t.Errorf("coder field should be 0 or 1, observed: %d", coder)
	}

	hash := str.FieldTable["hash"].Fvalue.(uint32)
	if hash != uint32(0) {
		t.Errorf("hash field should be 0, observed: %d", hash)
	}

	hashIsZero := str.FieldTable["hashIsZero"].Fvalue.(byte)
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

func TestByteArrayFromStringObject(t *testing.T) {
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

	if ret != "null" {
		t.Errorf("Expected null, got %s", ret)
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

	// to inspect usage message, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	obj := StringObjectFromGoString("allo!")
	obj.FieldTable["testField"] = Field{
		Ftype:  "..",
		Fvalue: nil,
	}
	ret := ObjectFieldToString(obj, "testField")

	if !strings.Contains(ret, "java/lang/String") {
		t.Errorf("Expected different return string, got %s", ret)
	}

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "not yet supported") {
		t.Errorf("Expected different error message, got %s", msg)
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
