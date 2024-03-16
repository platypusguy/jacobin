/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/object"
	"strings"
	"testing"
)

func TestStringClinit(t *testing.T) {
	globals.InitGlobals("test")
	_ = classloader.Init()
	classloader.LoadBaseClasses()
	retval := stringClinit(nil)
	if retval != nil {
		switch retval.(type) {
		case *GErrBlk:
			gErr := retval.(*GErrBlk)
			if !strings.Contains(gErr.ErrMsg, "TestStringClinit: Could not find java/lang/String") {
				classloader.MethAreaDump()
				t.Errorf("TestStringClinit: Unexpected error message. got %s", gErr.ErrMsg)
			}
			if gErr.ExceptionType != exceptions.ClassNotLoadedException {
				t.Errorf("TestStringClinit: Unexpected exception type. got %d", gErr.ExceptionType)
			}
		default:
			t.Errorf("TestStringClinit: Did not get expected error message, got %v", retval)
		}
	}
}

func TestStringToUpperCase(t *testing.T) {
	globals.InitGlobals("test")
	originalString := "He did the Monster Mash!"
	originalObj := object.NewPoolStringFromGoString(originalString)
	params := []interface{}{originalObj}
	ucObj := toUpperCase(params).(*object.Object)
	strUpper := object.GetGoStringFromObject(ucObj)
	expValue := strings.ToUpper(originalString)
	if string(strUpper) != expValue {
		t.Errorf("TestStringToUpperCase failed, expected: %s, observed: %s", expValue, strUpper)
	}
}

func TestStringToLowerCase(t *testing.T) {
	globals.InitGlobals("test")
	originalString := "It was a graveyard smash!"
	originalObj := object.NewPoolStringFromGoString(originalString)
	params := []interface{}{originalObj}
	ucObj := toLowerCase(params).(*object.Object)
	strUpper := object.GetGoStringFromObject(ucObj)
	expValue := strings.ToLower(originalString)
	if string(strUpper) != expValue {
		t.Errorf("TestStringToLowerCase failed, expected: %s, observed: %s", expValue, strUpper)
	}
}

func TestCompareToIgnoreCaseOk(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	bString := "It waS a graveYARD sMash!"
	aObj := object.NewPoolStringFromGoString(aString)
	bObj := object.NewPoolStringFromGoString(bString)
	params := []interface{}{aObj, bObj}
	result := compareToIgnoreCase(params).(int64)
	if result != 0 {
		t.Errorf("TestCompareToIgnoreCaseOk: expected: 0, observed: %d", result)
	}
}

func TestCompareToIgnoreCaseNotOk_1(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	bString := "It waS a graveYARE sMash!"
	aObj := object.NewPoolStringFromGoString(aString)
	bObj := object.NewPoolStringFromGoString(bString)
	params := []interface{}{aObj, bObj}
	result := compareToIgnoreCase(params).(int64)
	if result >= 0 {
		t.Errorf("TestCompareToIgnoreCaseOk_1: expected: <0, observed: %d", result)
	}
}

func TestCompareToIgnoreCaseNotOk_2(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	bString := "It waS a graveYARc sMash!"
	aObj := object.NewPoolStringFromGoString(aString)
	bObj := object.NewPoolStringFromGoString(bString)
	params := []interface{}{aObj, bObj}
	result := compareToIgnoreCase(params).(int64)
	if result <= 0 {
		t.Errorf("TestCompareToIgnoreCaseOk_2: expected: >0, observed: %d", result)
	}
}

func TestStringLength_1(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	aObj := object.NewPoolStringFromGoString(aString)
	params := []interface{}{aObj}
	result := stringLength(params).(int64)
	if result != 25 {
		t.Errorf("TestStringLength_1: expected: 25, observed: %d", result)
	}
}

func TestStringLength_2(t *testing.T) {
	globals.InitGlobals("test")
	aString := ""
	aObj := object.NewPoolStringFromGoString(aString)
	params := []interface{}{aObj}
	result := stringLength(params).(int64)
	if result != 0 {
		t.Errorf("TestStringLength_2: expected: 0, observed: %d", result)
	}
}

func TestSprintf_1(t *testing.T) {
	globals.InitGlobals("test")
	aString := "Mary had a %s little lamb"
	aObj := object.NewPoolStringFromGoString(aString)
	params := []interface{}{aObj}
	resultObj := (sprintf(params)).(*object.Object)
	str := object.GetGoStringFromObject(resultObj)
	if str != aString {
		t.Errorf("TestSprintf_1: expected: %s, observed: %s", aString, str)
	}
}

func TestSprintf_2(t *testing.T) {
	globals.InitGlobals("test")
	aString := "Mary had a %s lamb"
	bString := "little"
	cString := "Mary had a little lamb"
	aObj := object.NewPoolStringFromGoString(aString)
	aObj.DumpObject("TestSprintf_2 aObj", 0)
	bObj := object.NewPoolStringFromGoString(bString)
	bObj.DumpObject("TestSprintf_2 bObj", 0)
	var bArray []*object.Object
	bArray = append(bArray, bObj)
	classStr := "[Ljava/lang/Object"
	lsObj := object.MakeEmptyObjectWithClassName(&classStr)
	lsObj.FieldTable["value"] = object.Field{Ftype: classStr, Fvalue: bArray}
	params := []interface{}{aObj, lsObj}
	result := sprintf(params)

	switch result.(type) {
	case *GErrBlk:
		geptr := *(result.(*GErrBlk))
		errMsg := geptr.ErrMsg
		t.Errorf("TestSprintf_2: %s\n", errMsg)
	case *object.Object:
		obj := result.(*object.Object)
		obj.DumpObject("TestSprintf_2 result", 0)
		str := object.GetGoStringFromObject(obj)
		if str != cString {
			t.Errorf("TestSprintf_2: expected: %s, observed: %s", cString, str)
		}
	default:
		t.Errorf("TestSprintf_2: result type %T makes no sense", result)
	}
}
