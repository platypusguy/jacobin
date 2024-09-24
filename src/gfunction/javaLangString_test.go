/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/types"
	"reflect"
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
			if gErr.ExceptionType != excNames.ClassNotLoadedException {
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
	originalObj := object.StringObjectFromGoString(originalString)
	params := []interface{}{originalObj}
	ucObj := toUpperCase(params).(*object.Object)
	strUpper := object.GoStringFromStringObject(ucObj)
	expValue := strings.ToUpper(originalString)
	if string(strUpper) != expValue {
		t.Errorf("TestStringToUpperCase failed, expected: %s, observed: %s", expValue, strUpper)
	}
}

func TestStringToLowerCase(t *testing.T) {
	globals.InitGlobals("test")
	originalString := "It was a graveyard smash!"
	originalObj := object.StringObjectFromGoString(originalString)
	params := []interface{}{originalObj}
	ucObj := toLowerCase(params).(*object.Object)
	strUpper := object.GoStringFromStringObject(ucObj)
	expValue := strings.ToLower(originalString)
	if string(strUpper) != expValue {
		t.Errorf("TestStringToLowerCase failed, expected: %s, observed: %s", expValue, strUpper)
	}
}

func TestCompareToIgnoreCaseOk(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	bString := "It waS a graveYARD sMash!"
	aObj := object.StringObjectFromGoString(aString)
	bObj := object.StringObjectFromGoString(bString)
	params := []interface{}{aObj, bObj}
	result := stringCompareToIgnoreCase(params).(int64)
	if result != 0 {
		t.Errorf("TestCompareToIgnoreCaseOk: expected: 0, observed: %d", result)
	}
}

func TestCompareToIgnoreCaseNotOk_1(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	bString := "It waS a graveYARE sMash!"
	aObj := object.StringObjectFromGoString(aString)
	bObj := object.StringObjectFromGoString(bString)
	params := []interface{}{aObj, bObj}
	result := stringCompareToIgnoreCase(params).(int64)
	if result >= 0 {
		t.Errorf("TestCompareToIgnoreCaseOk_1: expected: <0, observed: %d", result)
	}
}

func TestCompareToIgnoreCaseNotOk_2(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	bString := "It waS a graveYARc sMash!"
	aObj := object.StringObjectFromGoString(aString)
	bObj := object.StringObjectFromGoString(bString)
	params := []interface{}{aObj, bObj}
	result := stringCompareToIgnoreCase(params).(int64)
	if result <= 0 {
		t.Errorf("TestCompareToIgnoreCaseOk_2: expected: >0, observed: %d", result)
	}
}

func TestStringLength_1(t *testing.T) {
	globals.InitGlobals("test")
	aString := "It was a graveyard smash!"
	aObj := object.StringObjectFromGoString(aString)
	params := []interface{}{aObj}
	result := stringLength(params).(int64)
	if result != 25 {
		t.Errorf("TestStringLength_1: expected: 25, observed: %d", result)
	}
}

func TestStringLength_2(t *testing.T) {
	globals.InitGlobals("test")
	aString := ""
	aObj := object.StringObjectFromGoString(aString)
	params := []interface{}{aObj}
	result := stringLength(params).(int64)
	if result != 0 {
		t.Errorf("TestStringLength_2: expected: 0, observed: %d", result)
	}
}

func TestSprintf_1(t *testing.T) {
	globals.InitGlobals("test")
	aString := "Mary had a %s little lamb"
	aObj := object.StringObjectFromGoString(aString)
	params := []interface{}{aObj}
	resultObj := (sprintf(params)).(*object.Object)
	str := object.GoStringFromStringObject(resultObj)
	if str != aString {
		t.Errorf("TestSprintf_1: expected: %s, observed: %s", aString, str)
	}
}

func TestSprintf_2(t *testing.T) {
	globals.InitGlobals("test")
	aString := "Mary had a %s lamb"
	bString := "little"
	cString := "Mary had a little lamb"

	aObj := object.StringObjectFromGoString(aString)
	aObj.DumpObject("TestSprintf_2 aObj", 0)
	bObj := object.StringObjectFromGoString(bString)
	bObj.DumpObject("TestSprintf_2 bObj", 0)

	var bArray []*object.Object
	bArray = append(bArray, bObj)
	classStr := "[Ljava/lang/Object"
	lsObj := object.MakeEmptyObjectWithClassName(&classStr)
	lsObj.FieldTable["value"] = object.Field{Ftype: classStr, Fvalue: bArray}
	lsObj.DumpObject("TestSprintf_2 lsObj", 0)

	params := []interface{}{aObj, lsObj}
	t.Logf("#params = %d\n", len(params))
	result := sprintf(params)

	switch result.(type) {
	case *GErrBlk:
		geptr := *(result.(*GErrBlk))
		errMsg := geptr.ErrMsg
		t.Errorf("TestSprintf_2: %s\n", errMsg)
	case *object.Object:
		obj := result.(*object.Object)
		obj.DumpObject("TestSprintf_2 result", 0)
		str := object.GoStringFromStringObject(obj)
		if str != cString {
			t.Errorf("TestSprintf_2: expected: %s, observed: %s", cString, str)
		}
	default:
		t.Errorf("TestSprintf_2: result type %T makes no sense", result)
	}
}

func TestSprintf_3(t *testing.T) {
	globals.InitGlobals("test")
	aString := "Mary had %d little lambs and her favorite number is %.4f"
	bCount := int64(3)
	bPi := float64(3.1416)
	cString := "Mary had 3 little lambs and her favorite number is 3.1416"

	aObj := object.StringObjectFromGoString(aString)
	aObj.DumpObject("TestSprintf_2 aObj", 0)

	classStr := "java/lang/Integer"
	bCountObj := object.MakeEmptyObjectWithClassName(&classStr)
	bCountObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: bCount}
	bCountObj.DumpObject("TestSprintf_2 bCountObj", 0)

	classStr = "java/lang/Double"
	bPiObj := object.MakeEmptyObjectWithClassName(&classStr)
	bPiObj.FieldTable["value"] = object.Field{Ftype: types.Double, Fvalue: bPi}
	bPiObj.DumpObject("TestSprintf_2 bCountObj", 0)

	var bArray []*object.Object
	bArray = append(bArray, bCountObj)
	bArray = append(bArray, bPiObj)
	classStr = "[Ljava/lang/Object"
	lsObj := object.MakeEmptyObjectWithClassName(&classStr)
	lsObj.FieldTable["value"] = object.Field{Ftype: classStr, Fvalue: bArray}
	lsObj.DumpObject("TestSprintf_2 lsObj", 0)

	params := []interface{}{aObj, lsObj}
	t.Logf("#params = %d\n", len(params))
	result := sprintf(params)

	switch result.(type) {
	case *GErrBlk:
		geptr := *(result.(*GErrBlk))
		errMsg := geptr.ErrMsg
		t.Errorf("TestSprintf_2: %s\n", errMsg)
	case *object.Object:
		obj := result.(*object.Object)
		obj.DumpObject("TestSprintf_2 result", 0)
		str := object.GoStringFromStringObject(obj)
		if str != cString {
			t.Errorf("TestSprintf_2: expected: %s, observed: %s", cString, str)
		}
	default:
		t.Errorf("TestSprintf_2: result type %T makes no sense", result)
	}
}

func TestContainsString(t *testing.T) {
	globals.InitGlobals("test")
	targetString := "I love seafood"
	targetStringObj := object.StringObjectFromGoString(targetString)
	searchString := "food"
	searchStringObj := object.StringObjectFromGoString(searchString)
	params := []interface{}{targetStringObj, searchStringObj}

	res := stringContains(params)
	if res != types.JavaBoolTrue {
		t.Errorf("TestContainsString failed, expected: %s to contain: %s", targetString, searchString)
	}
}

// === the following tests were generated by JetBrains AI Assistant

func TestJavaLangStringContentEquals(t *testing.T) {
	tests := []struct {
		name string
		obj1 object.Object
		obj2 object.Object
		want interface{}
	}{
		{
			name: "Test equal strings",
			obj1: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: []byte("Hello")},
				},
			},
			obj2: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: []byte("Hello")},
				},
			},
			want: types.JavaBoolTrue,
		},

		{
			name: "Test not equal strings",
			obj1: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: []byte("Hello")},
				},
			},
			obj2: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: []byte("World")},
				},
			},
			want: types.JavaBoolFalse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := javaLangStringContentEquals([]interface{}{&tt.obj1, &tt.obj2}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("javaLangStringContentEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}

// for the moment commented out due matters noted in JACOBIN-555

// func TestLastIndexOfCharacter(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		base       string
// 		searchChar int64
// 		start      int64
// 		want       int64
// 	}{
// 		{
// 			"Find character, from start",
// 			"Hello World!",
// 			108, // ASCII for 'l'
// 			0,
// 			9,
// 		},
// 		{
// 			"Find non-existent character, from start",
// 			"Hello World!",
// 			101, // ASCII for 'e'
// 			5,
// 			-1,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			baseObject := &object.Object{
// 				FieldTable: map[string]object.Field{
// 					"value": {Fvalue: []byte(tt.base)},
// 				},
// 			}
//
// 			if got := lastIndexOfCharacter([]interface{}{baseObject, tt.searchChar, tt.start}); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("lastIndexOfCharacter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestLastIndexOfString(t *testing.T) {
// 	str := "Hello, World!"
// 	searchString := "World"
// 	start := int64(0)
// 	want := int64(7)
//
// 	strObject := &object.Object{
// 		FieldTable: map[string]object.Field{
// 			"value": {Fvalue: []byte(str)},
// 		},
// 	}
//
// 	searchObject := &object.Object{
// 		FieldTable: map[string]object.Field{
// 			"value": {Fvalue: []byte(searchString)},
// 		},
// 	}
//
// 	if got := lastIndexOfString([]interface{}{strObject, searchObject, start}); got != want {
// 		t.Errorf("lastIndexOfString() = %v, want %v", got, want)
// 	}
// }
//
// func TestStringRegionMatchesWithoutIgnoreCase(t *testing.T) {
// 	baseStr := "Hello, World!"
// 	baseOffset := int64(0)
//
// 	compareStr := "Hello, Go!"
// 	compareOffset := int64(0)
//
// 	regionLength := int64(5) // Compare the first 5 characters of both strings
//
// 	baseStringObject := &object.Object{
// 		FieldTable: map[string]object.Field{
// 			"value": {Fvalue: []byte(baseStr)},
// 		},
// 	}
//
// 	compareStringObject := &object.Object{
// 		FieldTable: map[string]object.Field{
// 			"value": {Fvalue: []byte(compareStr)},
// 		},
// 	}
//
// 	want := types.JavaBoolTrue // because the first 5 characters of both strings are "Hello"
//
// 	if got := stringRegionMatches([]interface{}{baseStringObject, baseOffset, compareStringObject, compareOffset, regionLength}); got != want {
// 		t.Errorf("stringRegionMatches() = %v, want %v", got, want)
// 	}
// }
//
// func TestStringRegionMatchesWithIgnoreCase(t *testing.T) {
// 	baseStr := "HELLO, WORLD!"
// 	baseOffset := int64(0)
//
// 	compareStr := "hello, go!"
// 	compareOffset := int64(0)
//
// 	regionLength := int64(5) // Compare the first 5 characters of both strings
//
// 	baseStringObject := &object.Object{
// 		FieldTable: map[string]object.Field{
// 			"value": {Fvalue: []byte(baseStr)},
// 		},
// 	}
//
// 	compareStringObject := &object.Object{
// 		FieldTable: map[string]object.Field{
// 			"value": {Fvalue: []byte(compareStr)},
// 		},
// 	}
//
// 	want := types.JavaBoolTrue // because the first 5 characters of both strings are "hello" ignoring cases
//
// 	got := stringRegionMatches(
// 		[]interface{}{baseStringObject, types.JavaBoolTrue, baseOffset,
// 			compareStringObject, compareOffset, regionLength})
// 	if got != want {
// 		t.Errorf("stringRegionMatches() = %v, want %v", got, want)
// 	}
// }
