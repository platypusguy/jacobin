/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"

	"github.com/cespare/xxhash/v2"
)

// Covers: case int64
func TestHashAnything_Int64(t *testing.T) {
	var input int64 = -1234567890123
	h, err := HashAnything(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Compute expected using the same little-endian layout
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(input))
	expected := xxhash.Sum64(b)
	if h != expected {
		t.Fatalf("hash mismatch for int64: got %v, want %v", h, expected)
	}
}

// Covers: case float64
func TestHashAnything_Float64(t *testing.T) {
	input := 3.141592653589793
	h, err := HashAnything(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Compute expected using Float64bits then little-endian
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(input))
	expected := xxhash.Sum64(b)
	if h != expected {
		t.Fatalf("hash mismatch for float64: got %v, want %v", h, expected)
	}
}

// Covers: case []byte
func TestHashAnything_Bytes(t *testing.T) {
	input := []byte{0x00, 0xFF, 0x10, 0x20, 0x7F}
	h, err := HashAnything(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := xxhash.Sum64(input)
	if h != expected {
		t.Fatalf("hash mismatch for []byte: got %v, want %v", h, expected)
	}
}

// Covers: default case (successful JSON marshalling)
func TestHashAnything_Default_Struct(t *testing.T) {
	type sample struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	input := sample{A: 42, B: "answer"}

	h, err := HashAnything(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("unexpected marshal error in test: %v", err)
	}
	expected := xxhash.Sum64(jsonBytes)
	if h != expected {
		t.Fatalf("hash mismatch for default(struct): got %v, want %v", h, expected)
	}
}

// Covers: default case when JSON marshalling fails (error path)
func TestHashAnything_Default_JSONMarshalError(t *testing.T) {
	// json.Marshal on channels (and funcs) returns an error: unsupported type
	ch := make(chan int)
	h, err := HashAnything(ch)
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
	if h != 0 {
		t.Fatalf("expected returned hash to be 0 on error, got %v", h)
	}
}

// Sanity: deterministic hashing for identical inputs (not table-driven)
func TestHashAnything_Deterministic_Int64(t *testing.T) {
	var input int64 = 987654321
	h1, err1 := HashAnything(input)
	h2, err2 := HashAnything(input)
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v, %v", err1, err2)
	}
	if h1 != h2 {
		t.Fatalf("expected deterministic result, got %v and %v", h1, h2)
	}
}
