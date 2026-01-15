/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
    "jacobin/src/excNames"
    "jacobin/src/gfunction/ghelpers"
    "jacobin/src/globals"
    "jacobin/src/object"
    "jacobin/src/types"
    "os"
    "path/filepath"
    "testing"
)

// helper: make a Writer-like object that carries ghelpers.FilePath and ghelpers.FileHandle
func makeWriterObjForFile(t *testing.T, filePath string) *object.Object {
    t.Helper()
    fh, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
    if err != nil {
        t.Fatalf("failed to open test file: %v", err)
    }
    return &object.Object{FieldTable: map[string]object.Field{
        ghelpers.FilePath:   {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(filePath)},
        ghelpers.FileHandle: {Ftype: ghelpers.FileHandle, Fvalue: fh},
    }}
}

func TestBufferedWriter_Init_And_WriteOne(t *testing.T) {
    globals.InitStringPool()

    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test1.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()

    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    if res := bwWriteOneChar([]interface{}{target, int64('Z')}); res != nil {
        t.Fatalf("bwWriteOneChar error: %v", res)
    }

    if res := bwFlush([]interface{}{target}); res != nil {
        t.Fatalf("bwFlush error: %v", res)
    }

    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "Z" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "Z")
    }

    if res := bwClose([]interface{}{target}); res != nil {
        t.Fatalf("bwClose error: %v", res)
    }
}

func TestBufferedWriter_NewLine(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_newline.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    if res := bufferedWriterNewLine([]interface{}{target}); res != nil {
        t.Fatalf("bufferedWriterNewLine error: %v", res)
    }
    _ = bwFlush([]interface{}{target})

    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "\n" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "\n")
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteCharBuffer(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_chars.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    charVals := []int64{66, 67, 68} // B C D
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: charVals},
    }}

    if res := bwWriteCharBuffer([]interface{}{target, bufObj, int64(0), int64(3)}); res != nil {
        t.Fatalf("bwWriteCharBuffer error: %v", res)
    }

    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "BCD" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "BCD")
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteCharBuffer_ZeroLength(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_chars_zero.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    charVals := []int64{65, 66, 67}
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: charVals},
    }}

    res := bwWriteCharBuffer([]interface{}{target, bufObj, int64(1), int64(0)})
    if res == nil {
        t.Fatalf("expected int64(0) for zero-length write, got nil")
    }
    if v, ok := res.(int64); !ok || v != 0 {
        t.Fatalf("expected int64(0), got %T %v", res, res)
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteCharBuffer_ParamError(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_chars_paramerr.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    charVals := []int64{1, 2, 3}
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: charVals},
    }}

    res := bwWriteCharBuffer([]interface{}{target, bufObj, int64(2), int64(4)})
    if res == nil {
        t.Fatalf("expected error for out-of-bounds params")
    }
    if geb, ok := res.(*ghelpers.GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    } else {
        t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteCharBuffer_ValueTypeError(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_chars_typeerr.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    // "value" is wrong type (not []int64)
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.ByteArray, Fvalue: []byte{1, 2, 3}},
    }}

    res := bwWriteCharBuffer([]interface{}{target, bufObj, int64(0), int64(1)})
    if res == nil {
        t.Fatalf("expected IOException for wrong value type")
    }
    if geb, ok := res.(*ghelpers.GErrBlk); ok {
        if geb.ExceptionType != excNames.IOException {
            t.Fatalf("expected IOException, got %d", geb.ExceptionType)
        }
    } else {
        t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteStringBuffer(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_str.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    strObj := object.StringObjectFromGoString("Hello")
    if res := bwWriteStringBuffer([]interface{}{target, strObj, int64(1), int64(3)}); res != nil {
        t.Fatalf("bwWriteStringBuffer error: %v", res)
    }

    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "ell" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "ell")
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteStringBuffer_ZeroLength(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_str_zero.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    strObj := object.StringObjectFromGoString("abc")
    res := bwWriteStringBuffer([]interface{}{target, strObj, int64(0), int64(0)})
    if res == nil {
        t.Fatalf("expected int64(0) for zero-length write, got nil")
    }
    if v, ok := res.(int64); !ok || v != 0 {
        t.Fatalf("expected int64(0), got %T %v", res, res)
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteStringBuffer_ParamError(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_str_paramerr.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    strObj := object.StringObjectFromGoString("xyz")
    res := bwWriteStringBuffer([]interface{}{target, strObj, int64(2), int64(4)})
    if res == nil {
        t.Fatalf("expected error for out-of-bounds params")
    }
    if geb, ok := res.(*ghelpers.GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    } else {
        t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_WriteStringBuffer_ValueTypeError(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_str_typeerr.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    // wrong type in String.value
    bogus := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: []int64{1, 2, 3}},
    }}
    res := bwWriteStringBuffer([]interface{}{target, bogus, int64(0), int64(1)})
    if res == nil {
        t.Fatalf("expected IOException for wrong value type")
    }
    if geb, ok := res.(*ghelpers.GErrBlk); ok {
        if geb.ExceptionType != excNames.IOException {
            t.Fatalf("expected IOException, got %d", geb.ExceptionType)
        }
    } else {
        t.Fatalf("expected *ghelpers.GErrBlk, got %T", res)
    }
    _ = bwClose([]interface{}{target})
}

func TestBufferedWriter_Flush_And_AfterClose(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_flush.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    if res := bwFlush([]interface{}{target}); res != nil {
        t.Fatalf("bwFlush error: %v", res)
    }

    // close, then flush should error (underlying file is closed)
    _ = bwClose([]interface{}{target})
    res := bwFlush([]interface{}{target})
    if res == nil {
        t.Fatalf("expected error on flush after close")
    }
}

func TestBufferedWriter_WriteAfterClose_Errors(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_afterclose.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    // close first
    _ = bwClose([]interface{}{target})

    // write one char should error
    if res := bwWriteOneChar([]interface{}{target, int64('A')}); res == nil {
        t.Fatalf("expected error writing after close")
    }
    // newline should error
    if res := bufferedWriterNewLine([]interface{}{target}); res == nil {
        t.Fatalf("expected error newline after close")
    }
    // char buffer should error
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: []int64{65}},
    }}
    if res := bwWriteCharBuffer([]interface{}{target, bufObj, int64(0), int64(1)}); res == nil {
        t.Fatalf("expected error char buffer after close")
    }
    // string buffer should error
    strObj := object.StringObjectFromGoString("X")
    if res := bwWriteStringBuffer([]interface{}{target, strObj, int64(0), int64(1)}); res == nil {
        t.Fatalf("expected error string buffer after close")
    }
}

func TestBufferedWriter_Init_ErrorPaths(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_init_errs.txt")
    defer os.Remove(filePath)

    // base writer with both fields
    writerObj := makeWriterObjForFile(t, filePath)

    // target
    target := object.MakeEmptyObject()

    // missing FilePath
    writerMissingPath := &object.Object{FieldTable: map[string]object.Field{
        ghelpers.FileHandle: writerObj.FieldTable[ghelpers.FileHandle],
    }}
    res1 := bufferedWriterInit([]interface{}{target, writerMissingPath})
    if res1 == nil {
        t.Fatalf("expected error for missing FilePath")
    }

    // missing FileHandle
    writerMissingHandle := &object.Object{FieldTable: map[string]object.Field{
        ghelpers.FilePath: writerObj.FieldTable[ghelpers.FilePath],
    }}
    res2 := bufferedWriterInit([]interface{}{target, writerMissingHandle})
    if res2 == nil {
        t.Fatalf("expected error for missing FileHandle")
    }

    // wrong type for FileHandle
    writerBadHandleType := &object.Object{FieldTable: map[string]object.Field{
        ghelpers.FilePath:   writerObj.FieldTable[ghelpers.FilePath],
        ghelpers.FileHandle: {Ftype: ghelpers.FileHandle, Fvalue: "not-a-file"},
    }}
    res3 := bufferedWriterInit([]interface{}{target, writerBadHandleType})
    if res3 == nil {
        t.Fatalf("expected error for non-*os.File handle")
    }
}

func TestBufferedWriter_WriteOne_WrongArgType(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "bw_test_wrongarg.txt")
    defer os.Remove(filePath)

    writerObj := makeWriterObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := bufferedWriterInit([]interface{}{target, writerObj}); res != nil {
        t.Fatalf("bufferedWriterInit returned error: %v", res)
    }

    res := bwWriteOneChar([]interface{}{target, "bad"})
    if res == nil {
        t.Fatalf("expected error for wrong integer argument type")
    }
    _ = bwClose([]interface{}{target})
}
