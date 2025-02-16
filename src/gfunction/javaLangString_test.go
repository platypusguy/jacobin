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
		case *object.Object:
			str := object.GoStringFromStringObject(retval.(*object.Object))
			if str != "stringClinit" {
				t.Errorf("TestStringClinit: Expected \"stringClinit\", observed %v", str)
			}
			return
		default:
			t.Errorf("TestStringClinit: Did not get expected error message, got %v", retval)
		}
	}
	t.Error("TestStringClinit: stringClinit() returned nil")
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
					"value": {Fvalue: object.JavaByteArrayFromGoString("Hello")},
				},
			},
			obj2: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: object.JavaByteArrayFromGoString("Hello")},
				},
			},
			want: types.JavaBoolTrue,
		},

		{
			name: "Test not equal strings",
			obj1: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: object.JavaByteArrayFromGoString("Hello")},
				},
			},
			obj2: object.Object{
				FieldTable: map[string]object.Field{
					"value": {Fvalue: object.JavaByteArrayFromGoString("World")},
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
func TestLastIndexOfCharacter(t *testing.T) {
	tests := []struct {
		name       string
		base       string
		searchChar int64
		start      int64
		want       int64
	}{
		{
			"Find character, from start",
			"Hello World!",
			108, // ASCII for 'l'
			11,
			9,
		},
		{
			"Find non-existent character, from start",
			"Hello World!",
			102, // ASCII for 'f'
			5,
			-1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseObject := &object.Object{
				KlassName: types.StringPoolStringIndex,
				FieldTable: map[string]object.Field{
					"value": {Fvalue: object.JavaByteArrayFromGoString(tt.base)},
				},
			}

			if got := lastIndexOfCharacter([]interface{}{baseObject, tt.searchChar, tt.start}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lastIndexOfCharacter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastIndexOfString(t *testing.T) {
	str := "Hello, World!"
	searchString := "World"
	start := int64(12)
	want := int64(7)

	strObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(str)},
		},
	}

	searchObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(searchString)},
		},
	}

	if got := lastIndexOfString([]interface{}{strObject, searchObject, start}); got != want {
		t.Errorf("lastIndexOfString() = %v, want %v", got, want)
	}
}

func TestStringRegionMatchesWithoutIgnoreCase(t *testing.T) {
	baseStr := "Hello, World!"
	baseOffset := int64(0)

	compareStr := "Hello, Go!"
	compareOffset := int64(0)

	regionLength := int64(5) // Compare the first 5 characters of both strings

	baseStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseStr)},
		},
	}

	compareStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareStr)},
		},
	}

	want := types.JavaBoolTrue // because the first 5 characters of both strings are "Hello"

	if got := stringRegionMatches([]interface{}{baseStringObject, baseOffset, compareStringObject, compareOffset, regionLength}); got != want {
		t.Errorf("stringRegionMatches() = %v, want %v", got, want)
	}
}

func TestStringRegionMatchesWithIgnoreCase(t *testing.T) {
	baseStr := "HELLO, WORLD!"
	baseOffset := int64(0)

	compareStr := "hello, go!"
	compareOffset := int64(0)

	regionLength := int64(5) // Compare the first 5 characters of both strings

	baseStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseStr)},
		},
	}

	compareStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareStr)},
		},
	}

	want := types.JavaBoolTrue // because the first 5 characters of both strings are "hello" ignoring cases

	got := stringRegionMatches(
		[]interface{}{baseStringObject, types.JavaBoolTrue, baseOffset,
			compareStringObject, compareOffset, regionLength})
	if got != want {
		t.Errorf("stringRegionMatches() = %v, want %v", got, want)
	}
}

// test the variants of String.equals()
func TestStringEquals(t *testing.T) {
	baseStr := "Hello"
	baseStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseStr)},
		},
	}

	compareStr := "Hello"
	compareStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareStr)},
		},
	}

	if result := stringEquals([]interface{}{
		baseStringObject, compareStringObject}); result != types.JavaBoolTrue {
		t.Fatalf("The strings do not match when they should")
	}

	compareStr = "World"
	compareStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareStr)},
		},
	}

	if result := stringEquals([]interface{}{
		baseStringObject, compareStringObject}); result != types.JavaBoolFalse {
		t.Fatalf("The strings match when they shouldn't")
	}

	referenceStr := ""
	baseStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(referenceStr)},
		},
	}

	compareToStr := ""
	compareStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareToStr)},
		},
	}

	if result := stringEquals([]interface{}{
		baseStringObject, compareStringObject}); result != types.JavaBoolTrue {
		t.Fatalf("The empty strings do not match when they should")
	}
}

func TestStringEqualsIgnoreCase(t *testing.T) {
	// test same string, different case
	referenceStr := "hello"
	baseStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(referenceStr)},
		},
	}

	compareToStr := "HELLO"
	compareStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareToStr)},
		},
	}

	want := types.JavaBoolTrue
	if got := stringEqualsIgnoreCase([]interface{}{
		baseStringObject, compareStringObject}); got != want {
		t.Errorf("stringEqualsIgnoreCase() = %v, want %v", got, want)
	}

	// test same string, same case
	referenceStr = "hello"
	baseStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(referenceStr)},
		},
	}

	compareToStr = "hello"
	compareStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareToStr)},
		},
	}

	want = types.JavaBoolTrue
	if got := stringEqualsIgnoreCase([]interface{}{
		baseStringObject, compareStringObject}); got != want {
		t.Errorf("stringEqualsIgnoreCase() = %v, want %v", got, want)
	}

	// test unequal strings
	referenceStr = "hello"
	baseStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(referenceStr)},
		},
	}

	compareToStr = "world"
	compareStringObject = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(compareToStr)},
		},
	}

	want = types.JavaBoolFalse
	if got := stringEqualsIgnoreCase([]interface{}{
		baseStringObject, compareStringObject}); got != want {
		t.Errorf("stringEqualsIgnoreCase() = %v, want %v", got, want)
	}
}

// test the variants of index of string
func TestStringIndexOfString(t *testing.T) {
	baseStr := "hello world"
	baseStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseStr)},
		},
	}

	argStr := "world"
	argStringObject := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(argStr)},
		},
	}

	// Test "int indexOf(String str)"
	res := stringIndexOfString([]interface{}{
		baseStringObject, argStringObject})
	if res != int64(6) {
		t.Errorf("Expected 6 but received %d", res)
	}

	// Test "int indexOf(String str, int fromIndex)"
	res = stringIndexOfString([]interface{}{
		baseStringObject, argStringObject, int64(0)})
	if res != int64(6) {
		t.Errorf("Expected 6 but received %d", res)
	}

	res = stringIndexOfString([]interface{}{
		baseStringObject, argStringObject, int64(7)})
	if res != int64(-1) {
		t.Errorf("Expected -1 but received %d", res)
	}

	// Test "int indexOf(String str, int beginIndex, int endIndex)"
	res = stringIndexOfString([]interface{}{
		baseStringObject, argStringObject, int64(0), int64(11)})
	if res != int64(6) {
		t.Errorf("Expected 6 but received %d", res)
	}

	res = stringIndexOfString([]interface{}{
		baseStringObject, argStringObject, int64(7), int64(11)})
	if res != int64(-1) {
		t.Errorf("Expected -1 but received %d", res)
	}
}

func TestStringStartsWith(t *testing.T) {
	baseLiteral := "helloWorld"
	baseStr := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseLiteral)},
		},
	}

	prefixLiteral := "hello"
	prefix := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(prefixLiteral)},
		},
	}

	// single param startswith positive test
	result := stringStartsWith([]interface{}{baseStr, prefix})
	if result != types.JavaBoolTrue {
		t.Errorf("In TestStringStartsWith(), expected true but got false")
	}

	prefixLiteral = "World"
	prefix = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(prefixLiteral)},
		},
	}
	// single param startswith Negative test
	result = stringStartsWith([]interface{}{baseStr, prefix})
	if result != types.JavaBoolFalse {
		t.Errorf("In TestStringStartsWith() expected false but got true")
	}
}

func TestStringStartsWithUsingOffset(t *testing.T) {
	baseLiteral := "helloWorld"
	baseStr := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseLiteral)},
		},
	}

	offset := int64(3)
	prefixLiteral := "loW"
	prefix := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(prefixLiteral)},
		},
	}

	// with offset startswith positive test
	result := stringStartsWith([]interface{}{baseStr, prefix, offset})
	if result != types.JavaBoolTrue {
		t.Errorf(
			"In TestStringStartsWithUsingOffset() expected true but got false")
	}

	offset = int64(5)
	prefixLiteral = "Hello"
	prefix = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(prefixLiteral)},
		},
	}
	// with offset startswith negative test
	result = stringStartsWith([]interface{}{baseStr, prefix, offset})
	if result != types.JavaBoolFalse {
		t.Errorf(
			"In TestStringStartsWithUsingOffset() expected false but got true")
	}
}

func TestStringStrip(t *testing.T) {
	// First, we would create inputs that we would want to test
	testInputString := "    Hello, world!    "
	testObject := object.StringObjectFromGoString(testInputString)
	testInput := []interface{}{testObject}

	// Expected output for the input is trimmed string
	expectedOutput := "Hello, world!"

	output := object.GoStringFromStringObject(
		stringStrip(testInput).(*object.Object))

	// Check if the returned output matches the expected output
	if output != expectedOutput {
		t.Errorf("stringStrip(%v) = %v; wanted %v", testInputString, output, expectedOutput)
	}

	// Test with an empty string
	emptyInput := "                  "
	emptyObject := object.StringObjectFromGoString(emptyInput)
	emptyTestInput := []interface{}{emptyObject}
	expectedEmptyOutput := ""

	emptyOutput := object.GoStringFromStringObject(
		stringStrip(emptyTestInput).(*object.Object))

	if emptyOutput != expectedEmptyOutput {
		t.Errorf("stringStrip(%s) = %s; wanted %s", emptyInput, emptyOutput, expectedEmptyOutput)
	}

	// Test with no leading or trailing spaces
	noSpaceInput := "NoSpace"
	noSpaceObject := object.StringObjectFromGoString(noSpaceInput)
	noSpaceTestInput := []interface{}{noSpaceObject}
	expectedNoSpaceOutput := "NoSpace"

	noSpaceOutput := object.GoStringFromStringObject(
		stringStrip(noSpaceTestInput).(*object.Object))

	if noSpaceOutput != expectedNoSpaceOutput {
		t.Errorf("stringStrip(%s) = %s; wanted %s", noSpaceInput, noSpaceOutput, expectedNoSpaceOutput)
	}
}

func TestStringStripLeading(t *testing.T) {
	baseLiteral := "    Hello, World!"
	baseStr := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseLiteral)},
		},
	}

	expected := "Hello, World!"
	outputObj := stringStripLeading([]interface{}{baseStr}).(*object.Object)
	output :=
		object.GoStringFromJavaByteArray(outputObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if output != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, output)
	}

	// ---

	baseLiteral = "Hello, World!   "
	baseStr = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseLiteral)},
		},
	}

	expected = "Hello, World!   "
	outputObj = stringStripLeading([]interface{}{baseStr}).(*object.Object)
	output =
		object.GoStringFromJavaByteArray(outputObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if output != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, output)
	}

	// ---
	baseLiteral = "Hello, World!"
	baseStr = &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(baseLiteral)},
		},
	}

	expected = "Hello, World!"
	outputObj = stringStripLeading([]interface{}{baseStr}).(*object.Object)
	output =
		object.GoStringFromJavaByteArray(outputObj.FieldTable["value"].Fvalue.([]types.JavaByte))
	if output != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, output)
	}
}

func TestStringStripTrailing(t *testing.T) {
	inputLiteral := "Hello, World!   "
	input := &object.Object{
		KlassName: types.StringPoolStringIndex,
		FieldTable: map[string]object.Field{
			"value": {Fvalue: object.JavaByteArrayFromGoString(inputLiteral)},
		},
	}

	expected := "Hello, World!"

	outputRaw := stringStripTrailing([]interface{}{input})

	output := *outputRaw.(*object.Object)
	strippedString :=
		object.GoStringFromJavaByteArray(output.FieldTable["value"].Fvalue.([]types.JavaByte))
	// strippedString := string(output.FieldTable["value"].Fvalue.([]byte))
	if strippedString != expected {
		t.Errorf("Expected '%s' but got '%s'",
			expected, strippedString)
	}
}

// stringRegionMatches() tests
func TestStringRegionMatches(t *testing.T) {

	// Test case: baseOffset < 0
	params := []any{createStringObject("hello"), int64(-1), createStringObject("hello"), int64(0), int64(5)}
	result := stringRegionMatches(params)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false for baseOffset < 0, got %v", result)
	}

	// Test case: compareOffset < 0
	params = []any{createStringObject("hello"), int64(0), createStringObject("hello"), int64(-1), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false for compareOffset < 0, got %v", result)
	}

	// Test case: baseOffset + regionLength > len(baseByteArray)
	params = []any{createStringObject("hello"), int64(1), createStringObject("hello"), int64(0), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false for baseOffset + regionLength > len(baseByteArray), got %v", result)
	}

	// Test case: compareOffset + regionLength > len(compareByteArray)
	params = []any{createStringObject("hello"), int64(0), createStringObject("hello"), int64(1), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false for compareOffset + regionLength > len(compareByteArray), got %v", result)
	}

	// Test case: sections are equal, case-sensitive
	params = []any{createStringObject("hello"), int64(0), createStringObject("hello"), int64(0), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolTrue {
		t.Errorf("Expected true for equal sections, case-sensitive, got %v", result)
	}

	// Test case: sections are not equal, case-sensitive
	params = []any{createStringObject("hello"), int64(0), createStringObject("world"), int64(0), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false for unequal sections, case-sensitive, got %v", result)
	}

	// Test case: sections are equal, ignoring case
	params = []any{createStringObject("hello"), types.JavaBoolTrue, int64(0), createStringObject("HELLO"), int64(0), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolTrue {
		t.Errorf("Expected true for equal sections, ignoring case, got %v", result)
	}

	// Test case: sections are not equal, ignoring case
	params = []any{createStringObject("hello"), types.JavaBoolTrue, int64(0), createStringObject("WORLD"), int64(0), int64(5)}
	result = stringRegionMatches(params)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false for unequal sections, ignoring case, got %v", result)
	}
}

// stringRepeat() tests
func TestStringRepeat(t *testing.T) {

	// Test case: repeat string 0 times
	params := []interface{}{createStringObject("hello"), int64(0)}
	result := stringRepeat(params).(*object.Object)
	expected := ""
	if object.GoStringFromStringObject(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, object.GoStringFromStringObject(result))
	}

	// Test case: repeat string 1 time
	params = []interface{}{createStringObject("hello"), int64(1)}
	result = stringRepeat(params).(*object.Object)
	expected = "hello"
	if object.GoStringFromStringObject(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, object.GoStringFromStringObject(result))
	}

	// Test case: repeat string multiple times
	params = []interface{}{createStringObject("hello"), int64(3)}
	result = stringRepeat(params).(*object.Object)
	expected = "hellohellohello"
	if object.GoStringFromStringObject(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, object.GoStringFromStringObject(result))
	}

	// Test case: repeat empty string multiple times
	params = []interface{}{createStringObject(""), int64(3)}
	result = stringRepeat(params).(*object.Object)
	expected = ""
	if object.GoStringFromStringObject(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, object.GoStringFromStringObject(result))
	}
}

// stringSplit() tests
func TestStringSplitLimit(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	// Test case: valid input with no limit
	params := []interface{}{createStringObject("a,b,c"), createStringObject(","), int64(0)}
	result := stringSplitLimit(params).(*object.Object)
	expected := []string{"a", "b", "c"}
	for i, obj := range result.FieldTable["value"].Fvalue.([]*object.Object) {
		if object.GoStringFromStringObject(obj) != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], object.GoStringFromStringObject(obj))
		}
	}

	// Test case: valid input with limit
	params = []interface{}{createStringObject("a,b,c"), createStringObject(","), int64(2)}
	result = stringSplitLimit(params).(*object.Object)
	expected = []string{"a", "b,c"}
	for i, obj := range result.FieldTable["value"].Fvalue.([]*object.Object) {
		if object.GoStringFromStringObject(obj) != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], object.GoStringFromStringObject(obj))
		}
	}

	// Test case: invalid regex pattern
	params = []interface{}{createStringObject("a,b,c"), createStringObject("[a-"), int64(0)}
	errBlk := stringSplitLimit(params).(*GErrBlk)
	if errBlk == nil {
		t.Errorf("Expected PatternSyntaxException, got %v", result)
	}

	if errBlk.ExceptionType != excNames.PatternSyntaxException {
		t.Errorf("Expected PatternSyntaxException, got %v",
			excNames.JVMexceptionNames[errBlk.ExceptionType])
	}

	// Test case: empty input string
	params = []interface{}{createStringObject(""), createStringObject(","), int64(0)}
	result = stringSplitLimit(params).(*object.Object)
	expected = []string{""}
	for i, obj := range result.FieldTable["value"].Fvalue.([]*object.Object) {
		if object.GoStringFromStringObject(obj) != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], object.GoStringFromStringObject(obj))
		}
	}
}

// --- utility functions for tests above ---
func createStringObject(s string) *object.Object {
	return object.StringObjectFromGoString(s)
}
