/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"reflect"
	"testing"
)

func TestByteArrayOutputStreamInit(t *testing.T) {
	// Default constructor
	obj := object.MakeEmptyObject()
	params := []any{obj}
	ByteArrayOutputStreamInit(params)

	bufField, ok := obj.FieldTable["buf"]
	if !ok {
		t.Fatal("buf field not found")
	}
	buf := bufField.Fvalue.([]types.JavaByte)
	if len(buf) != 32 {
		t.Errorf("Expected default buffer size 32, got %d", len(buf))
	}

	countField, ok := obj.FieldTable["count"]
	if !ok {
		t.Fatal("count field not found")
	}
	if countField.Fvalue.(int64) != 0 {
		t.Errorf("Expected count 0, got %d", countField.Fvalue)
	}

	// Constructor with size
	obj2 := object.MakeEmptyObject()
	params2 := []any{obj2, int64(64)}
	ByteArrayOutputStreamInit(params2)

	buf2 := obj2.FieldTable["buf"].Fvalue.([]types.JavaByte)
	if len(buf2) != 64 {
		t.Errorf("Expected buffer size 64, got %d", len(buf2))
	}
}

func TestByteArrayOutputStreamWriteInt(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj})

	ByteArrayOutputStreamWriteInt([]any{obj, int64(65)}) // 'A'
	ByteArrayOutputStreamWriteInt([]any{obj, int64(66)}) // 'B'

	count := obj.FieldTable["count"].Fvalue.(int64)
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	buf := obj.FieldTable["buf"].Fvalue.([]types.JavaByte)
	if buf[0] != 65 || buf[1] != 66 {
		t.Errorf("Unexpected buffer content: %v", buf[:2])
	}
}

func TestByteArrayOutputStreamWriteBytes(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj, int64(2)}) // small initial size

	data := []types.JavaByte{1, 2, 3, 4, 5}
	bObj := object.MakePrimitiveObject("[B", types.ByteArray, data)

	// write(b, 1, 3) -> should write [2, 3, 4]
	ByteArrayOutputStreamWriteBytes([]any{obj, bObj, int64(1), int64(3)})

	count := obj.FieldTable["count"].Fvalue.(int64)
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	buf := obj.FieldTable["buf"].Fvalue.([]types.JavaByte)
	expected := []types.JavaByte{2, 3, 4}
	if !reflect.DeepEqual(buf[:count], expected) {
		t.Errorf("Expected %v, got %v", expected, buf[:count])
	}

	// Test NullPointerException
	res := ByteArrayOutputStreamWriteBytes([]any{obj, nil, int64(0), int64(0)})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", res)
	}

	// Test IndexOutOfBoundsException
	res = ByteArrayOutputStreamWriteBytes([]any{obj, bObj, int64(-1), int64(3)})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Errorf("Expected IndexOutOfBoundsException, got %v", res)
	}
}

func TestByteArrayOutputStreamWriteBytesAll(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj})

	data := []types.JavaByte{10, 20, 30}
	bObj := object.MakePrimitiveObject("[B", types.ByteArray, data)

	ByteArrayOutputStreamWriteBytesAll([]any{obj, bObj})

	count := obj.FieldTable["count"].Fvalue.(int64)
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	buf := obj.FieldTable["buf"].Fvalue.([]types.JavaByte)
	if !reflect.DeepEqual(buf[:count], data) {
		t.Errorf("Expected %v, got %v", data, buf[:count])
	}
}

func TestByteArrayOutputStreamReset(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj})
	ByteArrayOutputStreamWriteInt([]any{obj, int64(1)})

	ByteArrayOutputStreamReset([]any{obj})

	count := obj.FieldTable["count"].Fvalue.(int64)
	if count != 0 {
		t.Errorf("Expected count 0 after reset, got %d", count)
	}
}

func TestByteArrayOutputStreamSize(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj})
	ByteArrayOutputStreamWriteInt([]any{obj, int64(1)})
	ByteArrayOutputStreamWriteInt([]any{obj, int64(2)})

	size := ByteArrayOutputStreamSize([]any{obj})
	if size.(int64) != 2 {
		t.Errorf("Expected size 2, got %v", size)
	}
}

func TestByteArrayOutputStreamToByteArray(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj})
	data := []types.JavaByte{5, 10, 15}
	bObj := object.MakePrimitiveObject("[B", types.ByteArray, data)
	ByteArrayOutputStreamWriteBytesAll([]any{obj, bObj})

	res := ByteArrayOutputStreamToByteArray([]any{obj})
	resObj := res.(*object.Object)
	resData := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)

	if !reflect.DeepEqual(resData, data) {
		t.Errorf("Expected %v, got %v", data, resData)
	}
}

func TestByteArrayOutputStreamToString(t *testing.T) {
	obj := object.MakeEmptyObject()
	ByteArrayOutputStreamInit([]any{obj})
	// Write "Hello"
	for _, b := range "Hello" {
		ByteArrayOutputStreamWriteInt([]any{obj, int64(b)})
	}

	res := ByteArrayOutputStreamToString([]any{obj})
	resObj := res.(*object.Object)
	// Jacobin's StringObjectFromJavaByteArray creates a String object with "value" field
	resData := resObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if object.GoStringFromJavaByteArray(resData) != "Hello" {
		t.Errorf("Expected 'Hello', got %s", object.GoStringFromJavaByteArray(resData))
	}
}
