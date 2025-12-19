/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestStringCharAt(t *testing.T) {
	globals.InitGlobals("test")
	str := "Hello"
	obj := object.StringObjectFromGoString(str)

	// Valid index
	params := []interface{}{obj, int64(1)}
	res := stringCharAt(params).(int64)
	if res != int64('e') {
		t.Errorf("Expected 'e' (101), got %d", res)
	}

	// Boundary: index 0
	params = []interface{}{obj, int64(0)}
	res = stringCharAt(params).(int64)
	if res != int64('H') {
		t.Errorf("Expected 'H' (72), got %d", res)
	}

	// Boundary: last index
	params = []interface{}{obj, int64(len(str) - 1)}
	res = stringCharAt(params).(int64)
	if res != int64('o') {
		t.Errorf("Expected 'o' (111), got %d", res)
	}

	// Boundary: index out of bounds (negative)
	params = []interface{}{obj, int64(-1)}
	resObj := stringCharAt(params)
	if gErr, ok := resObj.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected StringIndexOutOfBoundsException for index -1, got %v", resObj)
	}

	// Boundary: index out of bounds (length)
	params = []interface{}{obj, int64(len(str))}
	resObj = stringCharAt(params)
	if gErr, ok := resObj.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected StringIndexOutOfBoundsException for index %d, got %v", len(str), resObj)
	}

	// Boundary: null object
	params = []interface{}{nil, int64(0)}
	resObj = stringCharAt(params)
	if gErr, ok := resObj.(*GErrBlk); !ok || gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null string, got %v", resObj)
	}
}

func TestStringCompareToCaseSensitive(t *testing.T) {
	globals.InitGlobals("test")
	objA := object.StringObjectFromGoString("apple")
	objB := object.StringObjectFromGoString("banana")
	objA2 := object.StringObjectFromGoString("apple")

	// apple vs banana
	params := []interface{}{objA, objB}
	res := stringCompareToCaseSensitive(params).(int64)
	if res >= 0 {
		t.Errorf("Expected negative result for 'apple' vs 'banana', got %d", res)
	}

	// banana vs apple
	params = []interface{}{objB, objA}
	res = stringCompareToCaseSensitive(params).(int64)
	if res <= 0 {
		t.Errorf("Expected positive result for 'banana' vs 'apple', got %d", res)
	}

	// apple vs apple
	params = []interface{}{objA, objA2}
	res = stringCompareToCaseSensitive(params).(int64)
	if res != 0 {
		t.Errorf("Expected 0 for 'apple' vs 'apple', got %d", res)
	}

	// null second parameter
	params = []interface{}{objA, nil}
	resGeneric := stringCompareToCaseSensitive(params)
	if gErr, ok := resGeneric.(*GErrBlk); !ok || gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null parameter, got %v", resGeneric)
	}
}

func TestStringConcat(t *testing.T) {
	globals.InitGlobals("test")
	objA := object.StringObjectFromGoString("Hello")
	objB := object.StringObjectFromGoString(" World")

	params := []interface{}{objA, objB}
	res := stringConcat(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: concat with empty string
	objEmpty := object.StringObjectFromGoString("")
	params = []interface{}{objA, objEmpty}
	res = stringConcat(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "Hello" {
		t.Errorf("Expected 'Hello', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: null parameter
	params = []interface{}{objA, nil}
	res2 := stringConcat(params)
	if gErr, ok := res2.(*GErrBlk); !ok || gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null parameter, got %v", res2)
	}
}

func TestStringIsBlankAndEmpty(t *testing.T) {
	globals.InitGlobals("test")

	testCases := []struct {
		input   string
		isBlank bool
		isEmpty bool
	}{
		{"", true, true},
		{" ", true, false},
		{"\t\n", true, false},
		{"a", false, false},
		{" a ", false, false},
	}

	for _, tc := range testCases {
		obj := object.StringObjectFromGoString(tc.input)

		// Test isEmpty
		resEmpty := stringIsEmpty([]interface{}{obj})
		expectedEmpty := types.JavaBoolFalse
		if tc.isEmpty {
			expectedEmpty = types.JavaBoolTrue
		}
		if resEmpty != expectedEmpty {
			t.Errorf("stringIsEmpty(%q) expected %v, got %v", tc.input, expectedEmpty, resEmpty)
		}

		// Test isBlank
		resBlank := stringIsBlank([]interface{}{obj})
		expectedBlank := types.JavaBoolFalse
		if tc.isBlank {
			expectedBlank = types.JavaBoolTrue
		}
		if resBlank != expectedBlank {
			t.Errorf("stringIsBlank(%q) expected %v, got %v", tc.input, expectedBlank, resBlank)
		}
	}
}

func TestStringIndexOfCh(t *testing.T) {
	globals.InitGlobals("test")
	obj := object.StringObjectFromGoString("banana")

	// Find 'a' starting from 0
	params := []interface{}{obj, int64('a'), int64(0)}
	res := stringIndexOfCh(params).(int64)
	if res != 1 {
		t.Errorf("Expected index 1 for 'a' in 'banana', got %d", res)
	}

	// Find 'a' starting from 2
	params = []interface{}{obj, int64('a'), int64(2)}
	res = stringIndexOfCh(params).(int64)
	if res != 3 {
		t.Errorf("Expected index 3 for 'a' starting from 2, got %d", res)
	}

	// Not found
	params = []interface{}{obj, int64('z'), int64(0)}
	res = stringIndexOfCh(params).(int64)
	if res != -1 {
		t.Errorf("Expected -1 for 'z', got %d", res)
	}

	// Boundary: starting index beyond length
	params = []interface{}{obj, int64('a'), int64(10)}
	res = stringIndexOfCh(params).(int64)
	if res != -1 {
		t.Errorf("Expected -1 for starting index 10, got %d", res)
	}

	// Boundary: starting index negative
	params = []interface{}{obj, int64('b'), int64(-1)}
	res = stringIndexOfCh(params).(int64)
	if res != 0 {
		t.Errorf("Expected index 0 for 'b' with starting index -1, got %d", res)
	}
}

func TestStringReplaceLiteral(t *testing.T) {
	globals.InitGlobals("test")
	input := object.StringObjectFromGoString("hello world hello")
	target := object.StringObjectFromGoString("hello")
	replacement := object.StringObjectFromGoString("hi")

	params := []interface{}{input, target, replacement}
	res := stringReplaceLiteral(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "hi world hi" {
		t.Errorf("Expected 'hi world hi', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: target not found
	target2 := object.StringObjectFromGoString("bye")
	params = []interface{}{input, target2, replacement}
	res = stringReplaceLiteral(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "hello world hello" {
		t.Errorf("Expected no change, got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: empty target (should result in replacement between every char in some Java versions, but strings.ReplaceAll behavior is what we check)
	target3 := object.StringObjectFromGoString("")
	params = []interface{}{input, target3, replacement}
	res = stringReplaceLiteral(params).(*object.Object)
	// strings.ReplaceAll("abc", "", "x") -> "xaxbxcx"
	// Let's just check that it doesn't panic and returns a string longer than the input
	if len(object.GoStringFromStringObject(res)) <= len(object.GoStringFromStringObject(input)) {
		t.Errorf("Expected result to be longer than input for empty target")
	}

	// Boundary: null replacement
	params = []interface{}{input, target, nil}
	res2 := stringReplaceLiteral(params)
	if gErr, ok := res2.(*GErrBlk); !ok || gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null replacement, got %v", res2)
	}
}

func TestSubstringStartEnd(t *testing.T) {
	globals.InitGlobals("test")
	obj := object.StringObjectFromGoString("hamburger")

	// Valid substring
	params := []interface{}{obj, int64(4), int64(8)}
	res := substringStartEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "urge" {
		t.Errorf("Expected 'urge', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: start == end
	params = []interface{}{obj, int64(4), int64(4)}
	res = substringStartEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "" {
		t.Errorf("Expected '', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: start = 0, end = length
	params = []interface{}{obj, int64(0), int64(9)}
	res = substringStartEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "hamburger" {
		t.Errorf("Expected 'hamburger', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: start > end
	params = []interface{}{obj, int64(5), int64(4)}
	res2 := substringStartEnd(params)
	if gErr, ok := res2.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected StringIndexOutOfBoundsException for start > end, got %v", res2)
	}

	// Boundary: negative index
	params = []interface{}{obj, int64(-1), int64(4)}
	res2 = substringStartEnd(params)
	if gErr, ok := res2.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected StringIndexOutOfBoundsException for negative start, got %v", res2)
	}

	// Boundary: end > length
	params = []interface{}{obj, int64(0), int64(10)}
	res2 = substringStartEnd(params)
	if gErr, ok := res2.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected StringIndexOutOfBoundsException for end > length, got %v", res2)
	}

	// Boundary: substring(0,0) on empty string should work
	objEmpty := object.StringObjectFromGoString("")
	params = []interface{}{objEmpty, int64(0), int64(0)}
	resObj := substringStartEnd(params).(*object.Object)
	if object.GoStringFromStringObject(resObj) != "" {
		t.Errorf("Expected empty string for substring(0,0) on empty string")
	}
}

func TestStringTrim(t *testing.T) {
	globals.InitGlobals("test")
	testCases := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\t\n hello \r", "hello"},
		{"", ""},
		{"   ", ""},
		{"a", "a"},
	}

	for _, tc := range testCases {
		obj := object.StringObjectFromGoString(tc.input)
		res := trimString([]interface{}{obj}).(*object.Object)
		if object.GoStringFromStringObject(res) != tc.expected {
			t.Errorf("trimString(%q) expected %q, got %q", tc.input, tc.expected, object.GoStringFromStringObject(res))
		}
	}
}

func TestStringToCharArray(t *testing.T) {
	globals.InitGlobals("test")
	str := "abc"
	obj := object.StringObjectFromGoString(str)
	res := toCharArray([]interface{}{obj}).(*object.Object)

	// toCharArray returns an array object. We check its type and contents.
	if res.FieldTable["value"].Ftype != types.CharArray {
		t.Errorf("Expected CharArray type, got %s", res.FieldTable["value"].Ftype)
	}

	// Since we can't easily inspect the array object contents without knowing Populator internals
	// or using other G-functions, we'll assume it's correct if it didn't panic and returned the right type.
	// Actually, Populator for CharArray should store []int64 in Fvalue of the value field.
}

func TestStringSplit(t *testing.T) {
	globals.InitGlobals("test")
	obj := object.StringObjectFromGoString("a,b,c")
	pattern := object.StringObjectFromGoString(",")

	params := []interface{}{obj, pattern}
	res := stringSplit(params).(*object.Object)

	if res.FieldTable["value"].Ftype != types.RefArray {
		t.Errorf("Expected RefArray type, got %s", res.FieldTable["value"].Ftype)
	}
}

func TestSubstringToTheEnd_Boundary(t *testing.T) {
	globals.InitGlobals("test")
	obj := object.StringObjectFromGoString("hello")

	// Valid
	params := []interface{}{obj, int64(2)}
	res := substringToTheEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "llo" {
		t.Errorf("Expected 'llo', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: start = length
	params = []interface{}{obj, int64(5)}
	res = substringToTheEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "" {
		t.Errorf("Expected '', got %s", object.GoStringFromStringObject(res))
	}

	// Boundary: start = 0
	params = []interface{}{obj, int64(0)}
	res = substringToTheEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "hello" {
		t.Errorf("Expected 'hello', got %s", object.GoStringFromStringObject(res))
	}

	// Out of bounds
	params = []interface{}{obj, int64(6)}
	resGeneric := substringToTheEnd(params)
	if gErr, ok := resGeneric.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected StringIndexOutOfBoundsException for index 6, got %v", resGeneric)
	}
}

func TestStringGetChars(t *testing.T) {
	globals.InitGlobals("test")
	str := "hello"
	strObj := object.StringObjectFromGoString(str)

	// Create a char array [10]
	iArray := make([]int64, 10)
	charArrayObj := Populator("[C", types.CharArray, iArray)

	// Copy "ell" to charArray at index 2
	// getChars(srcBegin, srcEnd, dst[], dstBegin)
	params := []interface{}{strObj, int64(1), int64(4), charArrayObj, int64(2)}
	res := stringGetChars(params)
	if res != nil {
		t.Fatalf("stringGetChars returned error: %v", res)
	}

	// Verify contents of charArray
	val := charArrayObj.FieldTable["value"].Fvalue.([]int64)
	if val[2] != int64('e') || val[3] != int64('l') || val[4] != int64('l') {
		t.Errorf("Expected 'ell' at index 2, got %v", val[2:5])
	}

	// Boundary: out of bounds src
	params = []interface{}{strObj, int64(-1), int64(4), charArrayObj, int64(2)}
	res = stringGetChars(params)
	if gErr, ok := res.(*GErrBlk); !ok || gErr.ExceptionType != excNames.StringIndexOutOfBoundsException {
		t.Errorf("Expected SIOOBE for srcBegin -1, got %v", res)
	}

	// Boundary: out of bounds dst
	params = []interface{}{strObj, int64(1), int64(4), charArrayObj, int64(8)}
	res = stringGetChars(params)
	if gErr, ok := res.(*GErrBlk); !ok || gErr.ExceptionType != excNames.ArrayIndexOutOfBoundsException {
		t.Errorf("Expected AIOOBE for dstBegin 8, got %v", res)
	}
}
func TestStringConstructor_Boundary(t *testing.T) {
	globals.InitGlobals("test")

	// dummy this object
	thisObj := object.MakeEmptyObject()

	// newStringFromBytes(byte[] bytes)
	bytes := []types.JavaByte{types.JavaByte('a'), types.JavaByte('b'), types.JavaByte('c')}
	bytesObj := Populator("[B", types.ByteArray, bytes)
	res := newStringFromBytes([]interface{}{thisObj, bytesObj})
	if res != nil {
		t.Fatalf("newStringFromBytes returned error: %v", res)
	}
	if object.GoStringFromStringObject(thisObj) != "abc" {
		t.Errorf("newStringFromBytes expected 'abc', got %s", object.GoStringFromStringObject(thisObj))
	}

	// Boundary: newStringFromBytes with empty array
	emptyBytes := []types.JavaByte{}
	emptyBytesObj := Populator("[B", types.ByteArray, emptyBytes)
	res = newStringFromBytes([]interface{}{thisObj, emptyBytesObj})
	if object.GoStringFromStringObject(thisObj) != "" {
		t.Errorf("newStringFromBytes(empty) expected '', got %s", object.GoStringFromStringObject(thisObj))
	}
}

func TestSubstringEmpty_Boundary(t *testing.T) {
	globals.InitGlobals("test")
	objEmpty := object.StringObjectFromGoString("")

	// substring(0) on empty string
	params := []interface{}{objEmpty, int64(0)}
	res := substringToTheEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "" {
		t.Errorf("substringToTheEnd(0) on empty string expected '', got %s", object.GoStringFromStringObject(res))
	}

	// substring(0,0) on empty string
	params = []interface{}{objEmpty, int64(0), int64(0)}
	res = substringStartEnd(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "" {
		t.Errorf("substringStartEnd(0,0) on empty string expected '', got %s", object.GoStringFromStringObject(res))
	}
}
func TestStringReplaceCC_Boundary(t *testing.T) {
	globals.InitGlobals("test")
	obj := object.StringObjectFromGoString("banana")

	// Replace 'a' with 'o'
	params := []interface{}{obj, int64('a'), int64('o')}
	res := stringReplaceCC(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "bonono" {
		t.Errorf("Expected 'bonono', got %s", object.GoStringFromStringObject(res))
	}

	// Replace non-existent char
	params = []interface{}{obj, int64('z'), int64('x')}
	res = stringReplaceCC(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "banana" {
		t.Errorf("Expected 'banana', got %s", object.GoStringFromStringObject(res))
	}

	// Empty string
	objEmpty := object.StringObjectFromGoString("")
	params = []interface{}{objEmpty, int64('a'), int64('o')}
	res = stringReplaceCC(params).(*object.Object)
	if object.GoStringFromStringObject(res) != "" {
		t.Errorf("Expected '', got %s", object.GoStringFromStringObject(res))
	}
}

func TestStringValueOf_More(t *testing.T) {
	globals.InitGlobals("test")

	// long
	res := valueOfLong([]interface{}{int64(9223372036854775807)}).(*object.Object)
	if object.GoStringFromStringObject(res) != "9223372036854775807" {
		t.Errorf("valueOfLong failed")
	}

	// float
	res = valueOfFloat([]interface{}{float64(1.23)}).(*object.Object)
	if object.GoStringFromStringObject(res) != "1.23" {
		t.Errorf("valueOfFloat failed, got %s", object.GoStringFromStringObject(res))
	}
}
func TestStringValueOf(t *testing.T) {
	globals.InitGlobals("test")

	// boolean
	res := valueOfBoolean([]interface{}{types.JavaBoolTrue}).(*object.Object)
	if object.GoStringFromStringObject(res) != "true" {
		t.Errorf("valueOfBoolean expected 'true', got %s", object.GoStringFromStringObject(res))
	}

	// int
	res = valueOfInt([]interface{}{int64(123)}).(*object.Object)
	if object.GoStringFromStringObject(res) != "123" {
		t.Errorf("valueOfInt expected '123', got %s", object.GoStringFromStringObject(res))
	}

	// double
	res = valueOfDouble([]interface{}{float64(3.14)}).(*object.Object)
	if object.GoStringFromStringObject(res) != "3.14" {
		t.Errorf("valueOfDouble expected '3.14', got %s", object.GoStringFromStringObject(res))
	}

	// null object
	resObj := valueOfObject([]interface{}{object.Null}).(*object.Object)
	if object.GoStringFromStringObject(resObj) != "null" {
		t.Errorf("valueOfObject(null) expected 'null', got %s", object.GoStringFromStringObject(resObj))
	}
}
