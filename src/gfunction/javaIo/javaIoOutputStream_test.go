/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"os"
	"testing"
)

func TestOutputStreamMethodRegistration(t *testing.T) {
	Load_Io_OutputStream()

	expectedMethods := []struct {
		sig   string
		slots int
	}{
		{"java/io/OutputStream.<clinit>()V", 0},
		{"java/io/OutputStream.<init>()V", 0},
		{"java/io/OutputStream.close()V", 0},
		{"java/io/OutputStream.flush()V", 0},
		{"java/io/OutputStream.nullOutputStream()Ljava/io/OutputStream;", 0},
		{"java/io/OutputStream.write(I)V", 1},
		{"java/io/OutputStream.write([B)V", 1},
		{"java/io/OutputStream.write([BII)V", 3},
	}

	for _, expected := range expectedMethods {
		gm, ok := ghelpers.MethodSignatures[expected.sig]
		if !ok {
			t.Errorf("Expected method signature %s was not registered", expected.sig)
			continue
		}
		if gm.ParamSlots != expected.slots {
			t.Errorf("Expected ParamSlots %d for %s, got %d", expected.slots, expected.sig, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Errorf("Expected non-nil GFunction for %s", expected.sig)
		}
	}
}

func TestOutputStreamInitAndNullOutputStream(t *testing.T) {
	Load_Io_OutputStream()

	// <init>()V returns nil
	initFunc := ghelpers.MethodSignatures["java/io/OutputStream.<init>()V"].GFunction
	res := initFunc([]any{&object.Object{}})
	if res != nil {
		t.Errorf("Expected init to return nil, got %v", res)
	}

	// nullOutputStream() returns an OutputStream object
	nullStreamFunc := ghelpers.MethodSignatures["java/io/OutputStream.nullOutputStream()Ljava/io/OutputStream;"].GFunction
	resNull := nullStreamFunc([]any{})
	obj, ok := resNull.(*object.Object)
	if !ok || obj == nil {
		t.Fatalf("Expected nullOutputStream to return *object.Object, got %T", resNull)
	}
	klassName := *(stringPool.GetStringPointer(obj.KlassName))
	if klassName != "java/io/OutputStream" {
		t.Errorf("Expected class java/io/OutputStream, got %s", klassName)
	}
}

func TestOutputStreamWriteInt(t *testing.T) {
	Load_Io_OutputStream()
	writeIntFunc := ghelpers.MethodSignatures["java/io/OutputStream.write(I)V"].GFunction

	// Missing or invalid params
	if res := writeIntFunc([]any{}); res != nil {
		t.Errorf("Expected nil for empty params, got %v", res)
	}
	if res := writeIntFunc([]any{&object.Object{}}); res == nil {
		t.Errorf("Expected error block for missing int parameter, got nil")
	}

	// Plain output stream without file handle (no-op)
	outObj := &object.Object{
		FieldTable: make(map[string]object.Field),
	}
	res := writeIntFunc([]any{outObj, int64(65)})
	if res != nil {
		t.Errorf("Expected write(int) without FileHandle to return nil, got %v", res)
	}

	// With temporary file handle
	tmpFile, err := os.CreateTemp("", "test_outputstream_write_int")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	outObj.FieldTable[ghelpers.FileHandle] = object.Field{
		Ftype:  "Ljava/io/FileDescriptor;",
		Fvalue: tmpFile,
	}

	res = writeIntFunc([]any{outObj, int64(65)})
	if res != nil {
		t.Errorf("Expected write(int) with FileHandle to return nil, got %v", res)
	}

	// Flush and close
	flushFunc := ghelpers.MethodSignatures["java/io/OutputStream.flush()V"].GFunction
	closeFunc := ghelpers.MethodSignatures["java/io/OutputStream.close()V"].GFunction

	if resFlush := flushFunc([]any{outObj}); resFlush != nil {
		t.Errorf("Expected flush to return nil, got %v", resFlush)
	}
	if resClose := closeFunc([]any{outObj}); resClose != nil {
		t.Errorf("Expected close to return nil, got %v", resClose)
	}

	// Verify written file content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "A" {
		t.Errorf("Expected content 'A', got '%s'", string(content))
	}
}

func TestOutputStreamWriteBytes(t *testing.T) {
	Load_Io_OutputStream()
	writeBytesFunc := ghelpers.MethodSignatures["java/io/OutputStream.write([B)V"].GFunction

	outObj := &object.Object{
		FieldTable: make(map[string]object.Field),
	}

	// Null byte array parameter
	errRes := writeBytesFunc([]any{outObj, nil})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for nil byte array, got %v", errRes)
	}

	// Byte array missing value field
	badArr := &object.Object{
		FieldTable: make(map[string]object.Field),
	}
	errRes = writeBytesFunc([]any{outObj, badArr})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for invalid byte array, got %v", errRes)
	}

	// Valid byte array with FileHandle
	tmpFile, err := os.CreateTemp("", "test_outputstream_write_bytes")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	outObj.FieldTable[ghelpers.FileHandle] = object.Field{
		Ftype:  "Ljava/io/FileDescriptor;",
		Fvalue: tmpFile,
	}

	data := []byte("Hello, World!")
	jbytes := object.JavaByteArrayFromGoByteArray(data)
	arrObj := object.MakeArrayFromRawArray(jbytes)

	res := writeBytesFunc([]any{outObj, arrObj})
	if res != nil {
		t.Errorf("Expected write([B) to return nil, got %v", res)
	}

	// Read content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", string(content))
	}
}

func TestOutputStreamWriteBytesRange(t *testing.T) {
	Load_Io_OutputStream()
	writeRangeFunc := ghelpers.MethodSignatures["java/io/OutputStream.write([BII)V"].GFunction

	outObj := &object.Object{
		FieldTable: make(map[string]object.Field),
	}

	data := []byte("Hello, World!")
	jbytes := object.JavaByteArrayFromGoByteArray(data)
	arrObj := object.MakeArrayFromRawArray(jbytes)

	// Out of bounds tests
	cases := []struct {
		off int64
		len int64
	}{
		{-1, 5},
		{0, -1},
		{10, 5},
		{0, 14},
	}

	for _, c := range cases {
		errRes := writeRangeFunc([]any{outObj, arrObj, c.off, c.len})
		if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Errorf("Expected IndexOutOfBoundsException for off=%d, len=%d, got %v", c.off, c.len, errRes)
		}
	}

	// Length 0 (no-op)
	if res := writeRangeFunc([]any{outObj, arrObj, int64(0), int64(0)}); res != nil {
		t.Errorf("Expected nil for len=0, got %v", res)
	}

	// Subarray write to file
	tmpFile, err := os.CreateTemp("", "test_outputstream_write_range")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	outObj.FieldTable[ghelpers.FileHandle] = object.Field{
		Ftype:  "Ljava/io/FileDescriptor;",
		Fvalue: tmpFile,
	}

	res := writeRangeFunc([]any{outObj, arrObj, int64(7), int64(5)})
	if res != nil {
		t.Errorf("Expected write([BII) to return nil, got %v", res)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "World" {
		t.Errorf("Expected 'World', got '%s'", string(content))
	}
}

func TestOutputStreamFlushCloseNil(t *testing.T) {
	// Call flush and close on nil or empty params
	if res := outputStreamFlush(nil); res != nil {
		t.Errorf("Expected nil, got %v", res)
	}
	if res := outputStreamClose(nil); res != nil {
		t.Errorf("Expected nil, got %v", res)
	}
	if res := outputStreamFlush([]any{nil}); res != nil {
		t.Errorf("Expected nil, got %v", res)
	}
	if res := outputStreamClose([]any{nil}); res != nil {
		t.Errorf("Expected nil, got %v", res)
	}
}
