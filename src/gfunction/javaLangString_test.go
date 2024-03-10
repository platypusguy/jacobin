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
	retval := stringClinit(nil)
	if retval == nil {
		t.Error("TestStringClinit: Was expecting an error message, but got none.")
	}
	switch retval.(type) {
	case *GErrBlk:
		gErr := retval.(*GErrBlk)
		if !strings.Contains(gErr.ErrMsg, "TestStringClinit: Could not find java/lang/String") {
			t.Errorf("TestStringClinit: Unexpected error message. got %s", gErr.ErrMsg)
		}
		if gErr.ExceptionType != exceptions.ClassNotLoadedException {
			t.Errorf("TestStringClinit: Unexpected exception type. got %d", gErr.ExceptionType)
		}
	default:
		t.Errorf("TestStringClinit: Did not get expected error message, got %v", retval)
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
