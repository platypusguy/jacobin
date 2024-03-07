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

func TestStringIndexPrimitives_1(t *testing.T) {
	var index uint32
	var str string
	str1 := "abc"
	str2 := "def"

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 0 {
		t.Errorf("Expected string index 0 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string value %s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 1 {
		t.Errorf("Expected string index 1 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string value %s but observed: %s", str2, str)
	}

	index = GetStringIndex(&str1)
	str = *GetStringPointer(index)
	t.Logf("str1 index %d: %s\n", index, str)
	if index != 0 {
		t.Errorf("Expected string index 0 but observed: %d", index)
	}
	if str != str1 {
		t.Errorf("Expected string value %s but observed: %s", str1, str)
	}

	index = GetStringIndex(&str2)
	str = *GetStringPointer(index)
	t.Logf("str2 index %d: %s\n", index, str)
	if index != 1 {
		t.Errorf("Expected string index 1 but observed: %d", index)
	}
	if str != str2 {
		t.Errorf("Expected string value %s but observed: %s", str2, str)
	}
}

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	stringLength = 100
)

func randomString(length int) string {
	bb := make([]byte, length)
	for i := range bb {
		bb[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bb)
}

func TestStringIndexPrimitives_2(t *testing.T) {
	var LIMIT uint32 = 10000
	t.Logf("string slice size to be filled up: %d\n", LIMIT)
	finalIndex := LIMIT - 1
	midIndex := LIMIT / 2
	var str string
	var index uint32
	var ix uint32
	for ix = 0; ix < LIMIT; ix++ {
		str = randomString(stringLength)
		index = GetStringIndex(&str)
	}
	t.Logf("last index value: %d\n", index)
	str = *GetStringPointer(0)
	t.Logf("str1 index 0: %s\n", str)
	str = *GetStringPointer(midIndex)
	t.Logf("str1 index %d: %s\n", midIndex, str)
	str = *GetStringPointer(finalIndex)
	t.Logf("str1 index %d: %s\n", finalIndex, str)

}
