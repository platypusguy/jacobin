/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/globals"
	"jacobin/types"
	"math/big"
	"strconv"
	"testing"
)

var fieldName = "fred"
var fieldStrValue = "42"
var fieldIntValue = int64(42)
var fieldFltValue = float64(42.0)
var theClassName = "jacobin.unit.test.object"

func makeGenericObject() *Object {
	globals.InitGlobals("test")
	className := theClassName
	return MakeEmptyObjectWithClassName(&className)
}

func setField(obj *Object, fieldName string, fieldType string, fieldValue interface{}) {
	fld := Field{Ftype: fieldType, Fvalue: fieldValue}
	obj.FieldTable[fieldName] = fld
}

func TestObjectFieldToStringPos(t *testing.T) {
	var observed string
	var obsFloat float64
	var err error
	obj := makeGenericObject()

	// BigInteger
	var bi = new(big.Int)
	_, ok := bi.SetString(fieldStrValue, 10)
	if !ok {
		t.Errorf("(big.Int).SetString() failed, skipping BigInteger test\n")
	} else {
		setField(obj, fieldName, types.BigInteger, bi)
		str := ObjectFieldToString(obj, fieldName)
		if str != fieldStrValue {
			t.Errorf("BigInteger, expected: %s, observed: %s\n", fieldStrValue, str)
		}
	}

	// Boolean scalar
	setField(obj, fieldName, types.Bool, types.JavaBoolFalse)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != "false" {
		t.Errorf("Boolean scalar, expected: false, observed: %s\n", observed)
	}
	setField(obj, fieldName, types.Bool, types.JavaBoolTrue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != "true" {
		t.Errorf("Boolean scalar, expected: true, observed: %s\n", observed)
	}

	// Boolean array
	var boolAry = []int64{types.JavaBoolTrue, types.JavaBoolFalse, types.JavaBoolTrue, types.JavaBoolFalse, types.JavaBoolTrue}
	setField(obj, fieldName, types.BoolArray, boolAry)
	observed = ObjectFieldToString(obj, fieldName)
	expected := "true false true false true"
	if observed != expected {
		t.Errorf("Boolean array, expected: %s, observed: %s\n", expected, observed)
	}

	// Scalars: byte, character, float, double, int, long, rune, short

	setField(obj, fieldName, types.Byte, fieldIntValue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != fieldStrValue {
		t.Errorf("Byte scalar, expected: %s, observed: %s\n", fieldStrValue, observed)
	}

	setField(obj, fieldName, types.Char, fieldIntValue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != fieldStrValue {
		t.Errorf("Character scalar, expected: %s, observed: %s\n", fieldStrValue, observed)
	}

	setField(obj, fieldName, types.Double, fieldFltValue)
	observed = ObjectFieldToString(obj, fieldName)
	obsFloat, err = strconv.ParseFloat(observed, 64)
	if err != nil {
		t.Errorf("Double strconv.ParseFloat(%s) failed, skipping types.Double test\n", observed)
	} else {
		if obsFloat != fieldFltValue {
			t.Errorf("Double scalar, expected: %f, observed: %f\n", fieldFltValue, obsFloat)
		}
	}

	setField(obj, fieldName, types.Float, fieldFltValue)
	observed = ObjectFieldToString(obj, fieldName)
	obsFloat, err = strconv.ParseFloat(observed, 64)
	if err != nil {
		t.Errorf("Float strconv.ParseFloat(%s) failed, skipping types.Double test\n", observed)
	} else {
		if obsFloat != fieldFltValue {
			t.Errorf("Float scalar, expected: %f, observed: %f\n", fieldFltValue, obsFloat)
		}
	}

	setField(obj, fieldName, types.Int, fieldIntValue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != fieldStrValue {
		t.Errorf("Integer scalar, expected: %s, observed: %s\n", fieldStrValue, observed)
	}

	setField(obj, fieldName, types.Long, fieldIntValue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != fieldStrValue {
		t.Errorf("Long scalar, expected: %s, observed: %s\n", fieldStrValue, observed)
	}

	setField(obj, fieldName, types.Rune, fieldIntValue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != fieldStrValue {
		t.Errorf("Rune scalar, expected: %s, observed: %s\n", fieldStrValue, observed)
	}

	setField(obj, fieldName, types.Short, fieldIntValue)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != fieldStrValue {
		t.Errorf("Short scalar, expected: %s, observed: %s\n", fieldStrValue, observed)
	}

	// Byte array
	var byteAry = []byte{65, 66, 67}
	setField(obj, fieldName, types.ByteArray, byteAry)
	observed = ObjectFieldToString(obj, fieldName)
	expected = "ABC"
	if observed != expected {
		t.Errorf("Byte array, expected: %s, observed: %s\n", expected, observed)
	}

	// Double array
	var doubleAry = []float64{1.1, 2.2, 3.3}
	setField(obj, fieldName, types.DoubleArray, doubleAry)
	observed = ObjectFieldToString(obj, fieldName)
	expected = "1.1 2.2 3.3"
	if observed != expected {
		t.Errorf("Double array, expected: %s, observed: %s\n", expected, observed)
	}

	// Reference scalar
	str := "ABCDEF"
	strObj := StringObjectFromGoString(str)
	setField(obj, fieldName, types.Ref, strObj)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != theClassName {
		t.Errorf("Reference scalar, expected: %s, observed: %s\n", theClassName, observed)
	}

	// Reference array
	str = "ABCDEF"
	strObj = StringObjectFromGoString(str)
	strObjAry := []*Object{strObj, strObj, strObj, strObj, strObj}
	setField(obj, fieldName, types.RefArray, strObjAry)
	observed = ObjectFieldToString(obj, fieldName)
	if observed != theClassName {
		t.Errorf("Reference array, expected: %s, observed: %s\n", theClassName, observed)
	}

	// nil
	observed = ObjectFieldToString(nil, fieldName)
	if observed != "null" {
		t.Errorf("nil, expected: %s, observed: %s\n", "null", observed)
	}

	// Null
	observed = ObjectFieldToString(Null, fieldName)
	if observed != "null" {
		t.Errorf("object.Null, expected: %s, observed: %s\n", "null", observed)
	}
}
