/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package ghelpers

import "testing"

func TestConvertArgsToParams_ZeroArgs(t *testing.T) {
	result := ConvertArgsToParams()
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

func TestConvertArgsToParams_OneArg(t *testing.T) {
	result := ConvertArgsToParams(42)
	if len(result) != 1 {
		t.Errorf("Expected length 1, got %d", len(result))
	}
	if result[0] != 42 {
		t.Errorf("Expected 42, got %v", result[0])
	}
}

func TestConvertArgsToParams_OneArgNil(t *testing.T) {
	result := ConvertArgsToParams(nil)
	if len(result) != 1 {
		t.Errorf("Expected length 1, got %d", len(result))
	}
	if result[0] != nil {
		t.Errorf("Expected nil, got %v", result[0])
	}
}

func TestConvertArgsToParams_TwoArgs(t *testing.T) {
	result := ConvertArgsToParams("hello", 100)
	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
	if result[0] != "hello" {
		t.Errorf("Expected 'hello', got %v", result[0])
	}
	if result[1] != 100 {
		t.Errorf("Expected 100, got %v", result[1])
	}
}

func TestConvertArgsToParams_TwoArgsWithNil(t *testing.T) {
	result := ConvertArgsToParams(nil, "world")
	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
	if result[0] != nil {
		t.Errorf("Expected nil at index 0, got %v", result[0])
	}
	if result[1] != "world" {
		t.Errorf("Expected 'world', got %v", result[1])
	}
}

func TestConvertArgsToParams_ThreeArgs(t *testing.T) {
	result := ConvertArgsToParams(1, 2.5, true)
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
	if result[0] != 1 {
		t.Errorf("Expected 1, got %v", result[0])
	}
	if result[1] != 2.5 {
		t.Errorf("Expected 2.5, got %v", result[1])
	}
	if result[2] != true {
		t.Errorf("Expected true, got %v", result[2])
	}
}

func TestConvertArgsToParams_ThreeArgsWithNil(t *testing.T) {
	result := ConvertArgsToParams("first", nil, "third")
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}
	if result[0] != "first" {
		t.Errorf("Expected 'first', got %v", result[0])
	}
	if result[1] != nil {
		t.Errorf("Expected nil at index 1, got %v", result[1])
	}
	if result[2] != "third" {
		t.Errorf("Expected 'third', got %v", result[2])
	}
}

func TestConvertArgsToParams_FiveArgs(t *testing.T) {
	result := ConvertArgsToParams(10, 20, 30, 40, 50)
	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}
	for i := 0; i < 5; i++ {
		expected := (i + 1) * 10
		if result[i] != expected {
			t.Errorf("Expected %d at index %d, got %v", expected, i, result[i])
		}
	}
}

func TestConvertArgsToParams_FiveArgsWithNil(t *testing.T) {
	result := ConvertArgsToParams(nil, "two", nil, 4, nil)
	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}
	if result[0] != nil {
		t.Errorf("Expected nil at index 0, got %v", result[0])
	}
	if result[1] != "two" {
		t.Errorf("Expected 'two', got %v", result[1])
	}
	if result[2] != nil {
		t.Errorf("Expected nil at index 2, got %v", result[2])
	}
	if result[3] != 4 {
		t.Errorf("Expected 4, got %v", result[3])
	}
	if result[4] != nil {
		t.Errorf("Expected nil at index 4, got %v", result[4])
	}
}
