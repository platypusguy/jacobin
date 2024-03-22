/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"bytes"
	"jacobin/globals"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"testing"
)

func TestNewStringObject(t *testing.T) {
	globals.InitGlobals("test")

	str := *NewStringObject()
	klassStr := *(stringPool.GetStringPointer(str.KlassName))
	if klassStr != "java/lang/String" {
		t.Errorf("Klass should be java/lang/String, observed: %s", klassStr)
	}

	value := str.FieldTable["value"].Fvalue.([]byte)
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

func TestByteArrayFromStringObject(t *testing.T) {
	globals.InitGlobals("test")
	statics.LoadStaticsString()

	constStr := "Mary had a little lamb whose fleece was white as snow."
	constBytes := []byte(constStr)

	strObj := StringObjectFromGoString(constStr)
	bb := ByteArrayFromStringObject(strObj)
	if !bytes.Equal(bb, constBytes) {
		t.Errorf("expected string value to be '%s', observed: '%s'", constStr, string(bb))
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
	if strValue == EmptyString { // if ""
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

	bb := ByteArrayFromStringPoolIndex(index)
	if bb == nil {
		t.Errorf("byte array from string pool index %d is nil", index)
	}
}

func TestUpdateStringObjectFromBytes(t *testing.T) {
	constStr := "Mary had a little lamb whose fleece was white as snow."
	constBytes := []byte(constStr)
	strObj := StringObjectFromGoString("To be updated")
	if !IsStringObject(strObj) {
		t.Errorf("expected IsStringObject(valid string object) to be true, observed false")
	}
	UpdateStringObjectFromBytes(strObj, constBytes)
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
