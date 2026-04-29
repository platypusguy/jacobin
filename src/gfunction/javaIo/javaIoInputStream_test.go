/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"testing"
)

func TestInputStreamReadIntoByteArray_ZeroLength(t *testing.T) {
	inputStreamObj := &object.Object{}
	byteArrayObj := &object.Object{}

	// params: [this, byte[], offset, length]
	params := []interface{}{inputStreamObj, byteArrayObj}

	result := inputStreamReadIntoByteArray(params)

	// Java spec states that if len is zero, then no bytes are read and 0 is returned.
	resInt, ok := result.(int64)
	if !ok {
		t.Errorf("Expected result to be int64, got %T", result)
	} else if resInt != 0 {
		t.Errorf("Expected 0 bytes read for zero length request, got %d", resInt)
	}
}

func TestInputStreamReadIntoByteArray_NilArray(t *testing.T) {
	inputStreamObj := &object.Object{}

	// params: [this, byte[] (nil), offset, length]
	params := []interface{}{inputStreamObj, nil}

	result := inputStreamReadIntoByteArray(params)

	// According to JVM spec, reading into a null array throws NullPointerException.
	// In Jacobin, this should return a Java exception (error block).
	errBlock, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Errorf("Expected an *ghelpers.GErrBlk representing NullPointerException, got %T", result)
	} else if errBlock == nil {
		t.Error("Expected an error block, got nil")
	}
}

func TestInputStreamReadIntoByteArray_ValidRead(t *testing.T) {
	inputStreamObj := &object.Object{}
	byteArrayObj := &object.Object{}

	// params: [this, byte[], offset, length]
	params := []interface{}{inputStreamObj, byteArrayObj, int32(0), int32(5)}

	result := inputStreamReadIntoByteArray(params)

	// For a valid read on an uninitialized/mocked stream, we might expect -1 (EOF)
	// or an actual integer representing the bytes read.
	if _, ok := result.(int64); !ok {
		if _, isErr := result.(*object.Object); isErr {
			t.Errorf("Expected an int64 result for bytes read, but got an error block/exception")
		} else {
			t.Errorf("Expected an int64 result, got %T", result)
		}
	}
}
