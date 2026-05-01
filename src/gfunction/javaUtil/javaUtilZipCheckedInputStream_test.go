/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil_test

import (
	"jacobin/src/gfunction/javaIo"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"testing"
)

func newCheckedInputStreamObj() *object.Object {
	className := "java/util/zip/CheckedInputStream"
	return object.MakeEmptyObjectWithClassName(&className)
}

func newFISObj() *object.Object {
	className := "java/io/FileInputStream"
	return object.MakeEmptyObjectWithClassName(&className)
}

func newCRC32Obj() *object.Object {
	className := "java/util/zip/CRC32"
	return object.MakeEmptyObjectWithClassName(&className)
}

func TestCheckedInputStream_Read(t *testing.T) {
	globals.InitStringPool()
	javaIo.Load_Io_FileInputStream()
	javaUtil.Load_Util_Zip_Crc32_Crc32c()
	javaUtil.Load_Util_Zip_CheckedInputStream()

	// Create a dummy file for testing
	tmpFile, err := os.CreateTemp("", "test_checked_input_stream")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("hello world")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 1. Setup Checksum (CRC32)
	crc := newCRC32Obj()
	javaUtil.Crc32InitIEEE([]interface{}{crc})

	// 2. Setup FileInputStream
	fis := newFISObj()
	pathObj := object.StringObjectFromGoString(tmpFile.Name())
	javaIo.InitFileInputStreamString([]interface{}{fis, pathObj})

	// 3. Setup CheckedInputStream
	cis := newCheckedInputStreamObj()
	javaUtil.CheckedInputStreamInit([]interface{}{cis, fis, crc})

	// Test read()
	val := javaUtil.CheckedInputStreamRead([]interface{}{cis})
	if v, ok := val.(int64); !ok || v != int64('h') {
		t.Errorf("Expected 'h' (104), got %v", val)
	}

	// Test read(byte[], off, len)
	buf := make([]types.JavaByte, 5)
	bufObj := object.MakePrimitiveObject("java/util/ArrayList", types.ByteArray, buf)
	n := javaUtil.CheckedInputStreamReadArray([]interface{}{cis, bufObj, int64(0), int64(5)})
	if nr, ok := n.(int64); !ok || nr != 5 {
		t.Errorf("Expected 5 bytes read, got %v", n)
	}

	// Check checksum
	cksum := javaUtil.CheckedInputStreamGetChecksum([]interface{}{cis})
	if cksum != crc {
		t.Errorf("Checksum object mismatch")
	}

	valC := javaUtil.Crc32GetValue([]interface{}{crc})
	if v, ok := valC.(int64); !ok || v == 0 {
		t.Errorf("Checksum should not be 0 after reading")
	}
}

func TestCheckedInputStream_Skip(t *testing.T) {
	globals.InitStringPool()
	javaIo.Load_Io_FileInputStream()
	javaUtil.Load_Util_Zip_Crc32_Crc32c()
	javaUtil.Load_Util_Zip_CheckedInputStream()

	// Create a dummy file for testing
	tmpFile, err := os.CreateTemp("", "test_checked_input_stream_skip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("hello world")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	crc := newCRC32Obj()
	javaUtil.Crc32InitIEEE([]interface{}{crc})

	fis := newFISObj()
	pathObj := object.StringObjectFromGoString(tmpFile.Name())
	javaIo.InitFileInputStreamString([]interface{}{fis, pathObj})

	cis := newCheckedInputStreamObj()
	javaUtil.CheckedInputStreamInit([]interface{}{cis, fis, crc})

	// Skip 6 bytes ("hello ")
	skipped := javaUtil.CheckedInputStreamSkip([]interface{}{cis, int64(6)})
	if s, ok := skipped.(int64); !ok || s != 6 {
		t.Errorf("Expected 6 bytes skipped, got %v", skipped)
	}

	// Read next byte, should be 'w'
	val := javaUtil.CheckedInputStreamRead([]interface{}{cis})
	if v, ok := val.(int64); !ok || v != int64('w') {
		t.Errorf("Expected 'w' (119), got %v", val)
	}
}
