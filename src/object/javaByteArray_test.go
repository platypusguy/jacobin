/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"jacobin/types"
	"testing"
)

func TestGoStringFromJavaByteArray(t *testing.T) {
	// Test case: non-empty array
	jbarr := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	expected := "Hello"
	result := GoStringFromJavaByteArray(jbarr)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test case: empty array
	jbarr = []types.JavaByte{}
	expected = ""
	result = GoStringFromJavaByteArray(jbarr)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// test that string creation from a JavaByte array works correctly
func TestGoStringFromJavaByteArrayAddl(t *testing.T) {
	javaByteArray :=
		[]types.JavaByte{0x4d, 0x61, 0x72, 0x79, 0x20, 0x68, 0x61, 0x64,
			0x20, 0x61, 0x20, 0x6c, 0x69, 0x74, 0x74, 0x6c, 0x65, 0x20,
			0x6c, 0x61, 0x6d, 0x62, 0x20, 0x77, 0x68, 0x6f, 0x73, 0x65,
			0x20, 0x66, 0x6c, 0x65, 0x65, 0x63, 0x65, 0x20, 0x77, 0x61,
			0x73, 0x20, 0x77, 0x68, 0x69, 0x74, 0x65, 0x20, 0x61, 0x73,
			0x20, 0x73, 0x6e, 0x6f, 0x77, 0x2e}
	goString := GoStringFromJavaByteArray(javaByteArray)
	if goString != "Mary had a little lamb whose fleece was white as snow." {
		t.Errorf("Expected 'Mary had a little lamb whose fleece was white as snow.', got '%s'", goString)
	}
}

func TestJavaByteArrayFromGoString(t *testing.T) {
	// Test case: non-empty string
	str := "Hello"
	expected := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	result := JavaByteArrayFromGoString(str)
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}

	// Test case: empty string
	str = ""
	expected = []types.JavaByte{}
	result = JavaByteArrayFromGoString(str)
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// test that a JavaByte array is correctly created from a Go string
func TestJavaByteArrayFromGoStringAddl(t *testing.T) {
	constStr := "Mary"
	jba := JavaByteArrayFromGoString(constStr)
	if len(jba) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(jba))
	}
	if jba[0] != 0x4d || jba[1] != 0x61 || jba[2] != 0x72 || jba[3] != 0x79 {
		t.Errorf("Expected [0x4d, 0x61, 0x72, 0x79], got %v", jba)
	}
}

func TestJavaByteArrayFromGoByteArray(t *testing.T) {
	// Test case: non-empty byte array
	gbarr := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	expected := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	result := JavaByteArrayFromGoByteArray(gbarr)
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}

	// Test case: empty byte array
	gbarr = []byte{}
	expected = []types.JavaByte{}
	result = JavaByteArrayFromGoByteArray(gbarr)
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGoByteArrayFromJavaByteArray(t *testing.T) {
	// Test case: non-empty Java byte array
	jbarr := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	expected := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	result := GoByteArrayFromJavaByteArray(jbarr)
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}

	// Test case: empty Java byte array
	jbarr = []types.JavaByte{}
	expected = []byte{}
	result = GoByteArrayFromJavaByteArray(jbarr)
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJavaByteArrayFromStringObject(t *testing.T) {
	// Test case: valid string object
	strObj := StringObjectFromGoString("Hello")
	expected := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	result := JavaByteArrayFromStringObject(strObj)
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}

	// Test case: nil object
	result = JavaByteArrayFromStringObject(nil)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestStringObjectFromJavaByteArray(t *testing.T) {
	// Test case: valid Java byte array
	jbarr := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	expected := "Hello"
	result := StringObjectFromJavaByteArray(jbarr)
	if GoStringFromStringObject(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, GoStringFromStringObject(result))
	}

	// Test case: empty Java byte array
	jbarr = []types.JavaByte{}
	expected = ""
	result = StringObjectFromJavaByteArray(jbarr)
	if GoStringFromStringObject(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, GoStringFromStringObject(result))
	}
}

func TestJavaByteArrayFromStringPoolIndex(t *testing.T) {
	// Test case: valid index
	index := uint32(0)
	expected := []types.JavaByte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // Assuming "Hello" is at index 0
	result := JavaByteArrayFromStringPoolIndex(index)
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}

	// Test case: invalid index
	index = uint32(999999)
	result = JavaByteArrayFromStringPoolIndex(index)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestJavaByteArrayEquals(t *testing.T) {
	// Test case: both arrays are nil
	if !JavaByteArrayEquals(nil, nil) {
		t.Errorf("Expected true, got false")
	}

	// Test case: one array is nil
	if JavaByteArrayEquals(nil, []types.JavaByte{0x48}) {
		t.Errorf("Expected false, got true")
	}

	// Test case: arrays have different lengths
	if JavaByteArrayEquals([]types.JavaByte{0x48}, []types.JavaByte{0x48, 0x65}) {
		t.Errorf("Expected false, got true")
	}

	// Test case: arrays are equal
	if !JavaByteArrayEquals([]types.JavaByte{0x48, 0x65}, []types.JavaByte{0x48, 0x65}) {
		t.Errorf("Expected true, got false")
	}

	// Test case: arrays are not equal
	if JavaByteArrayEquals([]types.JavaByte{0x48, 0x65}, []types.JavaByte{0x48, 0x66}) {
		t.Errorf("Expected false, got true")
	}
}

func TestJavaByteArrayEqualsIgnoreCase(t *testing.T) {
	// Test case: both arrays are nil
	if !JavaByteArrayEqualsIgnoreCase(nil, nil) {
		t.Errorf("Expected true, got false")
	}

	// Test case: one array is nil
	if JavaByteArrayEqualsIgnoreCase(nil, []types.JavaByte{0x48}) {
		t.Errorf("Expected false, got true")
	}

	// Test case: arrays have different lengths
	if JavaByteArrayEqualsIgnoreCase([]types.JavaByte{0x48}, []types.JavaByte{0x48, 0x65}) {
		t.Errorf("Expected false, got true")
	}

	// Test case: arrays are equal ignoring case
	if !JavaByteArrayEqualsIgnoreCase([]types.JavaByte{0x48, 0x65}, []types.JavaByte{0x48, 0x45}) {
		t.Errorf("Expected true, got false")
	}

	// Test case: arrays are not equal ignoring case
	if JavaByteArrayEqualsIgnoreCase([]types.JavaByte{0x48, 0x65}, []types.JavaByte{0x48, 0x66}) {
		t.Errorf("Expected false, got true")
	}
}
