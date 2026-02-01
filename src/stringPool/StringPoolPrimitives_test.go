/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package stringPool

import (
	"jacobin/src/globals"
	"math/rand"
	"testing"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	stringLength = 100
)

func randomString(maxlength int) string {
	halflength := maxlength / 2
	length := rand.Intn(maxlength)
	if length < halflength {
		length += halflength
	}
	bb := make([]byte, length)
	for i := range bb {
		bb[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bb)
}

func TestStringIndexPrimitives_1(t *testing.T) {
	index := uint32(0)
	var str string
	str1 := "Mary had a little lamb"
	str2 := "Whose fleece was white as snow"

	globals.InitGlobals("test") // Start with pool size to 2.

	DumpStringPool("TestStringIndexPrimitives_1: should be empty")
	sz := GetStringPoolSize()
	if sz != 5 {
		t.Errorf("Expected string repo size 5 but observed: %d", sz)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 5 {
		t.Errorf("Expected string str1 index=5 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string str1 value=%s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 6 {
		t.Errorf("Expected string str2 index=6 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string str2 value=%s but observed: %s", str2, str)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 5 {
		t.Errorf("Expected string str1 index=5 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string str1 value=%s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 6 {
		t.Errorf("Expected string str2 index=6 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string str2 value=%s but observed: %s", str2, str)
	}

	// Add 16 random strings to pool giving a total of 21.
	for ix := 0; ix < 16; ix++ {
		str = randomString(stringLength)
		index = GetStringIndex(&str)
	}

	// Check resultant pool sizer.
	sz = GetStringPoolSize()
	if sz != 23 {
		t.Errorf("Expected string repo size 23 but observed: %d", sz)
	}

	// Dump the pool--if needed in case of test failure
	// DumpStringPool("TestStringIndexPrimitives_1: final repo")
}

func TestStringIndexPrimitives_2(t *testing.T) {
	// NOTE that TestStringIndexPrimitives_2 is dependent on globals.InitStringPool!
	globals.InitGlobals("test")
	postInitSize := GetStringPoolSize()
	if postInitSize != 5 {
		t.Errorf("Expected string repo size=4 but observed: %d", postInitSize)
	}

	var LIMIT uint32 = 1_000
	var LIMITp2 uint32 = LIMIT + postInitSize
	t.Logf("string slice size to be filled up: %d\n", LIMITp2)
	finalIndex := LIMITp2 - 1
	// 	t.Logf("final index value: %d\n", finalIndex)
	midIndex := finalIndex / 2
	t.Logf("mid index value: %d\n", midIndex)
	midString := "Mary had a little lamb"
	var str string
	var ix uint32
	var i uint32
	// Add LIMIT more strings.
	for ix = 2; ix < LIMITp2; ix++ {
		if ix == midIndex {
			str = midString
			i = GetStringIndex(&str)
		} else {
			str = randomString(stringLength)
			_ = GetStringIndex(&str)
		}
	}

	str = *GetStringPointer(i)
	if str != midString {
		t.Errorf("Expected mid-string value: %s. Observed: %s", midString, str)
	}
}

func TestVariousInitStringPool(t *testing.T) {
	globals.InitGlobals("testInit")
	s := GetStringPointer(uint32(0))
	if *s != "" {
		t.Errorf("Expected null string, got %s", *s)
	}

	if GetStringPoolSize() < 3 {
		t.Errorf("Expected initiailized string pool size >= 3, got %d", GetStringPoolSize())
	}
}

func TestEmptyStringPool(t *testing.T) {
	globals.InitGlobals("testInit")
	globals.InitStringPool()
	initialSize := GetStringPoolSize()

	str1 := "test1"
	GetStringIndex(&str1)

	str2 := "test2"
	GetStringIndex(&str2)

	newSize := GetStringPoolSize()
	if newSize != initialSize+2 {
		t.Errorf("Expected string pool size = init size + 2 but observed: %d", newSize)
	}

	EmptyStringPool()
	emptySize := GetStringPoolSize()
	if emptySize != initialSize {
		t.Errorf("Expected string pool size = init size, observed: %d", emptySize)
	}

}
