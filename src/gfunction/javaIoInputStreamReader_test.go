package gfunction

import (
    "jacobin/excNames"
    "jacobin/globals"
    "jacobin/object"
    "jacobin/types"
    "os"
    "path/filepath"
    "testing"
)

// helper: make an InputStream-like object that carries FilePath and FileHandle for reading
func makeInputStreamObjForFile(t *testing.T, filePath string) *object.Object {
    t.Helper()
    fh, err := os.Open(filePath)
    if err != nil {
        t.Fatalf("failed to open file for reading: %v", err)
    }
    return &object.Object{FieldTable: map[string]object.Field{
        FilePath:  {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(filePath)},
        FileHandle:{Ftype: types.FileHandle, Fvalue: fh},
    }}
}

func TestInputStreamReader_Init_And_ReadOne_And_EOF(t *testing.T) {
    globals.InitStringPool()

    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "isr_test1.txt")
    content := []byte("ABC")
    if err := os.WriteFile(filePath, content, 0o644); err != nil {
        t.Fatalf("failed to write test file: %v", err)
    }
    defer os.Remove(filePath)

    inStreamObj := makeInputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()

    if res := inputStreamReaderInit([]interface{}{target, inStreamObj}); res != nil {
        t.Fatalf("inputStreamReaderInit returned error: %v", res)
    }

    // Read 'A'
    if v := isrReadOneChar([]interface{}{target}); v == nil {
        t.Fatalf("isrReadOneChar returned nil")
    } else if v.(int64) != int64('A') {
        t.Fatalf("got %d want %d", v.(int64), int64('A'))
    }
    // Read 'B'
    if v := isrReadOneChar([]interface{}{target}); v.(int64) != int64('B') {
        t.Fatalf("unexpected second char: %d", v.(int64))
    }
    // Read 'C'
    if v := isrReadOneChar([]interface{}{target}); v.(int64) != int64('C') {
        t.Fatalf("unexpected third char: %d", v.(int64))
    }
    // Next read -> EOF (-1)
    if v := isrReadOneChar([]interface{}{target}); v.(int64) != int64(-1) {
        t.Fatalf("expected -1 at EOF, got %d", v.(int64))
    }
}

func TestInputStreamReader_ReadCharBufferSubset(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "isr_test2.txt")
    content := []byte("HELLO")
    if err := os.WriteFile(filePath, content, 0o644); err != nil {
        t.Fatalf("failed to write test file: %v", err)
    }
    defer os.Remove(filePath)

    inStreamObj := makeInputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := inputStreamReaderInit([]interface{}{target, inStreamObj}); res != nil {
        t.Fatalf("inputStreamReaderInit returned error: %v", res)
    }

    // dest buffer of size 8
    dest := make([]int64, 8)
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: dest},
    }}

    // read 3 bytes into offset 2
    v := isrReadCharBufferSubset([]interface{}{target, bufObj, int64(2), int64(3)})
    if v == nil {
        t.Fatalf("expected count, got nil")
    }
    if n := v.(int64); n != 3 {
        t.Fatalf("expected 3 bytes read, got %d", n)
    }
    // verify positions 2..4 are 'H','E','L'
    got := bufObj.FieldTable["value"].Fvalue.([]int64)
    if got[2] != int64('H') || got[3] != int64('E') || got[4] != int64('L') {
        t.Fatalf("buffer contents incorrect: %+v", got)
    }
}

func TestInputStreamReader_ReadCharBufferSubset_ParamError(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "isr_test3.txt")
    if err := os.WriteFile(filePath, []byte("12345"), 0o644); err != nil {
        t.Fatalf("failed to write test file: %v", err)
    }
    defer os.Remove(filePath)

    inStreamObj := makeInputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := inputStreamReaderInit([]interface{}{target, inStreamObj}); res != nil {
        t.Fatalf("inputStreamReaderInit returned error: %v", res)
    }

    dest := make([]int64, 5)
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: dest},
    }}

    res := isrReadCharBufferSubset([]interface{}{target, bufObj, int64(4), int64(3)})
    if res == nil {
        t.Fatalf("expected error for OOB params")
    }
    if geb, ok := res.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    } else {
        t.Fatalf("expected *GErrBlk, got %T", res)
    }
}

func TestInputStreamReader_Ready_And_Close(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "isr_test4.txt")
    if err := os.WriteFile(filePath, []byte("X"), 0o644); err != nil {
        t.Fatalf("failed to write test file: %v", err)
    }
    defer os.Remove(filePath)

    inStreamObj := makeInputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := inputStreamReaderInit([]interface{}{target, inStreamObj}); res != nil {
        t.Fatalf("inputStreamReaderInit returned error: %v", res)
    }

    // ready([reader, inputStream]) should return 1 for open handle
    if v := isrReady([]interface{}{target, inStreamObj}); v.(int64) != 1 {
        t.Fatalf("expected ready==1, got %d", v.(int64))
    }

    // close underlying handle to provoke not-ready
    fh := inStreamObj.FieldTable[FileHandle].Fvalue.(*os.File)
    _ = fh.Close()
    if v := isrReady([]interface{}{target, inStreamObj}); v.(int64) != 0 {
        t.Fatalf("expected ready==0 after closing handle, got %d", v.(int64))
    }

    // also test isrClose closes the reader's own handle (copied during init)
    // re-open for close test
    inStreamObj = makeInputStreamObjForFile(t, filePath)
    target = object.MakeEmptyObject()
    if res := inputStreamReaderInit([]interface{}{target, inStreamObj}); res != nil {
        t.Fatalf("inputStreamReaderInit returned error: %v", res)
    }
    if res := isrClose([]interface{}{target}); res != nil {
        t.Fatalf("isrClose returned error: %v", res)
    }
    // second close on same os.File should fail at Go level; verify by trying to close again
    fh2 := target.FieldTable[FileHandle].Fvalue.(*os.File)
    if err := fh2.Close(); err == nil {
        t.Fatalf("expected error on closing already closed file, got nil")
    }
}
