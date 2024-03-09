/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
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

	EmptyStringPool()
	DumpStringPool("TestStringIndexPrimitives_1: should be empty")
	sz := GetStringPoolSize()
	if sz != 0 {
		t.Errorf("Expected string repo size 0 but observed: %d", sz)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 0 {
		t.Errorf("Expected string str1 index 0 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string str1 value %s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 1 {
		t.Errorf("Expected string str2 index 1 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string str2 value %s but observed: %s", str2, str)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 0 {
		t.Errorf("Expected string str1 index 0 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string str1 value %s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 1 {
		t.Errorf("Expected string str2 index 1 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string str2 value %s but observed: %s", str2, str)
	}

	for ix := 0; ix < 18; ix++ {
		str = randomString(stringLength)
		index = GetStringIndex(&str)
	}

	sz = GetStringPoolSize()
	if sz != 20 {
		t.Errorf("Expected string repo size 20 but observed: %d", sz)
	}

	DumpStringPool("TestStringIndexPrimitives_1: final repo")
}

func TestStringIndexPrimitives_2(t *testing.T) {
	var LIMIT uint32 = 1000000
	t.Logf("string slice size to be filled up: %d\n", LIMIT)
	finalIndex := LIMIT - 1
	t.Logf("final index value: %d\n", finalIndex)
	midIndex := LIMIT / 2
	t.Logf("mid index value: %d\n", midIndex)
	midString := "Mary had a little lamb"
	var index uint32 = 0
	var str string
	var ix uint32

	EmptyStringPool()
	DumpStringPool("TestStringIndexPrimitives_2: should be empty")
	sz := GetStringPoolSize()
	if sz != 0 {
		t.Errorf("Expected string repo size 0 but observed: %d", sz)
	}

	for ix = 0; ix < LIMIT; ix++ {
		if ix == midIndex {
			str = midString
		} else {
			str = randomString(stringLength)
		}
		index = GetStringIndex(&str)
		//t.Logf("DEBUG %d) string %d %s\n", ix, index, str)
	}
	t.Logf("last index value: %d\n", index)
	str = *GetStringPointer(0)
	t.Logf("str1 index 0: %s\n", str)
	str = *GetStringPointer(midIndex)
	t.Logf("str1 index %d: %s\n", midIndex, str)
	if str != midString {
		t.Errorf("Expected mid-string value %s but observed: %s", midString, str)
	}
	str = *GetStringPointer(finalIndex)
	t.Logf("str1 index %d: %s\n", finalIndex, str)

	sz = GetStringPoolSize()
	if sz != LIMIT {
		t.Errorf("Expected string repo size %d but observed: %d", LIMIT, sz)
	}
	if sz < 100 {
		DumpStringPool("TestStringIndexPrimitives_2: final repo")
	}

}
