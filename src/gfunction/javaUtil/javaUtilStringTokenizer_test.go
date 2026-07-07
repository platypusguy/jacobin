/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestStringTokenizer_Basic(t *testing.T) {
	globals.InitStringPool()

	classNameST := "java/util/StringTokenizer"
	stObj := object.MakeEmptyObjectWithClassName(&classNameST)

	strObj := object.StringObjectFromGoString("Hello World Jacobin")
	delimObj := object.StringObjectFromGoString(" ")

	// <init>(String, String)
	stringTokenizerInit2([]interface{}{stObj, strObj, delimObj})

	// countTokens()
	count := stringTokenizerCountTokens([]interface{}{stObj})
	if count.(int64) != 3 {
		t.Errorf("countTokens expected 3, got %v", count)
	}

	// hasMoreTokens()
	if stringTokenizerHasMoreTokens([]interface{}{stObj}) != types.JavaBoolTrue {
		t.Errorf("hasMoreTokens expected true")
	}

	// nextToken()
	token1 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token1.(*object.Object)) != "Hello" {
		t.Errorf("token1 expected 'Hello', got %v", object.GoStringFromStringObject(token1.(*object.Object)))
	}

	token2 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token2.(*object.Object)) != "World" {
		t.Errorf("token2 expected 'World', got %v", object.GoStringFromStringObject(token2.(*object.Object)))
	}

	token3 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token3.(*object.Object)) != "Jacobin" {
		t.Errorf("token3 expected 'Jacobin', got %v", object.GoStringFromStringObject(token3.(*object.Object)))
	}

	// hasMoreTokens() false
	if stringTokenizerHasMoreTokens([]interface{}{stObj}) != types.JavaBoolFalse {
		t.Errorf("hasMoreTokens expected false")
	}
}

func TestStringTokenizer_Delims(t *testing.T) {
	globals.InitStringPool()

	classNameST := "java/util/StringTokenizer"
	stObj := object.MakeEmptyObjectWithClassName(&classNameST)

	strObj := object.StringObjectFromGoString("a:b;c")
	delimObj := object.StringObjectFromGoString(":")

	// <init>(String, String)
	stringTokenizerInit2([]interface{}{stObj, strObj, delimObj})

	token1 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token1.(*object.Object)) != "a" {
		t.Errorf("token1 expected 'a', got %v", object.GoStringFromStringObject(token1.(*object.Object)))
	}

	// nextToken(String)
	delimObj2 := object.StringObjectFromGoString(";")
	token2 := stringTokenizerNextTokenWithDelims([]interface{}{stObj, delimObj2})
	if object.GoStringFromStringObject(token2.(*object.Object)) != ":b" {
		t.Errorf("token2 expected ':b', got %v", object.GoStringFromStringObject(token2.(*object.Object)))
	}

	token3 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token3.(*object.Object)) != "c" {
		t.Errorf("token3 expected 'c', got %v", object.GoStringFromStringObject(token3.(*object.Object)))
	}
}

func TestStringTokenizer_ReturnDelims(t *testing.T) {
	globals.InitStringPool()

	classNameST := "java/util/StringTokenizer"
	stObj := object.MakeEmptyObjectWithClassName(&classNameST)

	strObj := object.StringObjectFromGoString("a:b")
	delimObj := object.StringObjectFromGoString(":")

	// <init>(String, String, boolean)
	stringTokenizerInit([]interface{}{stObj, strObj, delimObj, int64(1)})

	// countTokens()
	count := stringTokenizerCountTokens([]interface{}{stObj})
	if count.(int64) != 3 {
		t.Errorf("countTokens expected 3, got %v", count)
	}

	token1 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token1.(*object.Object)) != "a" {
		t.Errorf("token1 expected 'a', got %v", object.GoStringFromStringObject(token1.(*object.Object)))
	}

	token2 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token2.(*object.Object)) != ":" {
		t.Errorf("token2 expected ':', got %v", object.GoStringFromStringObject(token2.(*object.Object)))
	}

	token3 := stringTokenizerNextToken([]interface{}{stObj})
	if object.GoStringFromStringObject(token3.(*object.Object)) != "b" {
		t.Errorf("token3 expected 'b', got %v", object.GoStringFromStringObject(token3.(*object.Object)))
	}
}
