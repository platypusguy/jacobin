/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package stringPool

import (
	"jacobin/globals"
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
	if sz != 2 {
		t.Errorf("Expected string repo size 2 but observed: %d", sz)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 2 {
		t.Errorf("Expected string str1 index=2 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string str1 value=%s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 3 {
		t.Errorf("Expected string str2 index=3 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string str2 value=%s but observed: %s", str2, str)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 2 {
		t.Errorf("Expected string str1 index=2 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string str1 value=%s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 3 {
		t.Errorf("Expected string str2 index=3 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string str2 value=%s but observed: %s", str2, str)
	}

	// Add 16 random strings to pool giving a total of 20.
	for ix := 0; ix < 16; ix++ {
		str = randomString(stringLength)
		index = GetStringIndex(&str)
	}

	// Check resultant pool sizer.
	sz = GetStringPoolSize()
	if sz != 20 {
		t.Errorf("Expected string repo size 20 but observed: %d", sz)
	}

	// Dump the pool.
	DumpStringPool("TestStringIndexPrimitives_1: final repo")
}

func TestStringIndexPrimitives_2(t *testing.T) {
	// NOTE that TestStringIndexPrimitives_2 is dependent on globals::InitStringPool!

	var LIMIT uint32 = 1000000
	var LIMITp2 uint32 = LIMIT + 2
	t.Logf("string slice size to be filled up: 2 + %d\n", LIMIT)
	finalIndex := LIMITp2 - 1
	t.Logf("final index value: %d\n", finalIndex)
	midIndex := LIMITp2 / 2
	t.Logf("mid index value: %d\n", midIndex)
	midString := "Mary had a little lamb"
	var str string
	var ix uint32

	globals.InitGlobals("test") // Start with pool size to 2.

	DumpStringPool("TestStringIndexPrimitives_2: should only have 2 entries")
	sz := GetStringPoolSize()
	if sz != 2 {
		t.Errorf("Expected string repo size=2 but observed: %d", sz)
	}

	// Add LIMIT more strings.
	for ix = 2; ix < LIMITp2; ix++ {
		if ix == midIndex {
			str = midString
		} else {
			str = randomString(stringLength)
		}
		_ = GetStringIndex(&str)
		// t.Logf("DEBUG %d) string %d: %s\n", ix, index, str)
	}

	// Report
	str = *GetStringPointer(2)
	t.Logf("First index (2): %s\n", str)
	str = *GetStringPointer(midIndex)
	t.Logf("Mid index (%d): %s\n", midIndex, str)
	if str != midString {
		t.Errorf("Expected mid-string value: %s. Observed: %s", midString, str)
	}
	str = *GetStringPointer(finalIndex)
	t.Logf("Last index (%d): %s\n", finalIndex, str)
	sz = GetStringPoolSize()
	if sz != (LIMIT + 2) {
		t.Errorf("Expected string repo size=%d but observed: %d", LIMIT+2, sz)
	}
	if sz < 100 {
		DumpStringPool("TestStringIndexPrimitives_2: final repo")
	}

}
