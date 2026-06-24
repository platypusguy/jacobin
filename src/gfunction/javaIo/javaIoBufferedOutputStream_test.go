/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"testing"
)

func TestBufferedOutputStream_Basic(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Io_FileOutputStream()
	Load_Io_FilterOutputStream()
	Load_Io_BufferedOutputStream()

	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "bostream_test.bin")
	defer os.Remove(filePath)

	under := makeOutputStreamObjForFileBuffered(t, filePath)
	self := object.MakeEmptyObject()

	// Init with small buffer for testing flushes
	if res := BufferedOutputStreamInit([]interface{}{self, under, int64(4)}); res != nil {
		t.Fatalf("Init failed: %v", res)
	}

	// Write 3 bytes (less than buffer size 4)
	for i := range 3 {
		if res := BufferedOutputStreamWriteInt([]interface{}{self, int64('a' + i)}); res != nil {
			t.Fatalf("WriteInt failed at %d: %v", i, res)
		}
	}

	// Verify nothing written yet
	data, _ := os.ReadFile(filePath)
	if len(data) != 0 {
		t.Fatalf("Expected 0 bytes, got %d", len(data))
	}

	// Write 4th byte, should flush
	if res := BufferedOutputStreamWriteInt([]interface{}{self, int64('d')}); res != nil {
		t.Fatalf("WriteInt failed: %v", res)
	}

	// Still nothing? Wait, BufferedOutputStream.write(int) flushes when it's FULL?
	// In Java's BufferedOutputStream.write(int):
	// if (count >= buf.length) { flushBuffer(); } buf[count++] = (byte)b;
	// So if size is 4, when writing 5th byte it flushes.

	if res := BufferedOutputStreamWriteInt([]interface{}{self, int64('e')}); res != nil {
		t.Fatalf("WriteInt failed: %v", res)
	}

	data, _ = os.ReadFile(filePath)
	if len(data) != 4 {
		t.Fatalf("Expected 4 bytes, got %d: %q", len(data), string(data))
	}

	// Flush the rest
	if res := BufferedOutputStreamFlush([]interface{}{self}); res != nil {
		t.Fatalf("Flush failed: %v", res)
	}

	data, _ = os.ReadFile(filePath)
	if string(data) != "abcde" {
		t.Fatalf("Expected 'abcde', got %q", string(data))
	}
}

func TestBufferedOutputStream_WriteRange(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Io_FileOutputStream()
	Load_Io_FilterOutputStream()
	Load_Io_BufferedOutputStream()

	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "bostream_range_test.bin")
	defer os.Remove(filePath)

	under := makeOutputStreamObjForFileBuffered(t, filePath)
	self := object.MakeEmptyObject()

	// Init with buffer size 10
	if res := BufferedOutputStreamInit([]interface{}{self, under, int64(10)}); res != nil {
		t.Fatalf("Init failed: %v", res)
	}

	// Write range that fits in buffer
	buf1 := []types.JavaByte{types.JavaByte('1'), types.JavaByte('2'), types.JavaByte('3')}
	bObj1 := &object.Object{FieldTable: map[string]object.Field{"value": {Ftype: "[B", Fvalue: buf1}}}
	if res := BufferedOutputStreamWriteRange([]interface{}{self, bObj1, int64(0), int64(3)}); res != nil {
		t.Fatalf("WriteRange 1 failed: %v", res)
	}

	// Write range that exceeds buffer (should write directly)
	buf2 := []types.JavaByte{types.JavaByte('A'), types.JavaByte('B'), types.JavaByte('C'), types.JavaByte('D'), types.JavaByte('E'), types.JavaByte('F'), types.JavaByte('G'), types.JavaByte('H'), types.JavaByte('I'), types.JavaByte('J'), types.JavaByte('K')}
	bObj2 := &object.Object{FieldTable: map[string]object.Field{"value": {Ftype: "[B", Fvalue: buf2}}}
	// len(buf2) = 11 > 10
	if res := BufferedOutputStreamWriteRange([]interface{}{self, bObj2, int64(0), int64(11)}); res != nil {
		t.Fatalf("WriteRange 2 failed: %v", res)
	}

	// After large write, the previous 3 bytes should also be flushed
	data, _ := os.ReadFile(filePath)
	if string(data) != "123ABCDEFGHIJK" {
		t.Fatalf("Expected '123ABCDEFGHIJK', got %q", string(data))
	}
}

func makeOutputStreamObjForFileBuffered(t *testing.T, filePath string) *object.Object {
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("os.Create failed: %v", err)
	}
	under := object.MakeEmptyObject()
	s := "java/io/FileOutputStream"
	under.KlassName = stringPool.GetStringIndex(&s)
	under.FieldTable[ghelpers.FilePath] = object.Field{Ftype: "[B", Fvalue: []byte(filePath)}
	under.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: "Ljava/io/FileDescriptor;", Fvalue: f}
	return under
}
