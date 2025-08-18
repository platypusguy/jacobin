package gfunction

import (
    "jacobin/src/excNames"
    "jacobin/src/globals"
    "jacobin/src/object"
    "jacobin/src/types"
    "os"
    "path/filepath"
    "testing"
)

// helper: make an OutputStream-like object that carries FilePath and FileHandle
func makeOutputStreamObjForFile(t *testing.T, filePath string) *object.Object {
    t.Helper()
    // open for write/create/truncate
    fh, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
    if err != nil {
        t.Fatalf("failed to open test file: %v", err)
    }
    return &object.Object{FieldTable: map[string]object.Field{
        FilePath:  {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(filePath)},
        FileHandle:{Ftype: types.FileHandle, Fvalue: fh},
    }}
}

func TestOutputStreamWriter_Init_And_WriteOne(t *testing.T) {
    globals.InitStringPool()

    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "osw_test1.txt")
    defer os.Remove(filePath)

    // underlying OutputStream object
    outStreamObj := makeOutputStreamObjForFile(t, filePath)
    // target OutputStreamWriter object
    target := object.MakeEmptyObject()

    // init(OutputStream)
    if res := initOutputStreamWriter([]interface{}{target, outStreamObj}); res != nil {
        t.Fatalf("initOutputStreamWriter returned error: %v", res)
    }

    // write(int)
    if res := oswWriteOneChar([]interface{}{target, int64('A')}); res != nil {
        t.Fatalf("oswWriteOneChar error: %v", res)
    }

    // flush to ensure persisted
    if res := oswFlush([]interface{}{target}); res != nil {
        t.Fatalf("oswFlush error: %v", res)
    }

    // verify content
    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "A" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "A")
    }

    // close
    if res := oswClose([]interface{}{target}); res != nil {
        t.Fatalf("oswClose error: %v", res)
    }
}

func TestOutputStreamWriter_WriteCharBuffer(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "osw_test2.txt")
    defer os.Remove(filePath)

    outStreamObj := makeOutputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()

    if res := initOutputStreamWriter([]interface{}{target, outStreamObj}); res != nil {
        t.Fatalf("initOutputStreamWriter returned error: %v", res)
    }

    // prepare char[] as []int64 in an object field "value"
    charVals := []int64{66, 67, 68} // B C D
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: charVals},
    }}

    if res := oswWriteCharBuffer([]interface{}{target, bufObj, int64(0), int64(3)}); res != nil {
        t.Fatalf("oswWriteCharBuffer error: %v", res)
    }

    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "BCD" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "BCD")
    }

    _ = oswClose([]interface{}{target})
}

func TestOutputStreamWriter_WriteCharBuffer_ParamError(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "osw_test3.txt")
    defer os.Remove(filePath)

    outStreamObj := makeOutputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := initOutputStreamWriter([]interface{}{target, outStreamObj}); res != nil {
        t.Fatalf("initOutputStreamWriter returned error: %v", res)
    }

    // length goes past end of buffer -> expect IndexOutOfBoundsException
    charVals := []int64{1, 2, 3}
    bufObj := &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.IntArray, Fvalue: charVals},
    }}

    res := oswWriteCharBuffer([]interface{}{target, bufObj, int64(2), int64(4)})
    if res == nil {
        t.Fatalf("expected error for out-of-bounds params")
    }
    if geb, ok := res.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    } else {
        t.Fatalf("expected *GErrBlk, got %T", res)
    }

    _ = oswClose([]interface{}{target})
}

func TestOutputStreamWriter_WriteStringBuffer(t *testing.T) {
    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "osw_test4.txt")
    defer os.Remove(filePath)

    outStreamObj := makeOutputStreamObjForFile(t, filePath)
    target := object.MakeEmptyObject()
    if res := initOutputStreamWriter([]interface{}{target, outStreamObj}); res != nil {
        t.Fatalf("initOutputStreamWriter returned error: %v", res)
    }

    strObj := object.StringObjectFromGoString("Hello")
    // write subset: "ell"
    if res := oswWriteStringBuffer([]interface{}{target, strObj, int64(1), int64(3)}); res != nil {
        t.Fatalf("oswWriteStringBuffer error: %v", res)
    }

    bytes, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    if string(bytes) != "ell" {
        t.Fatalf("content mismatch: got %q want %q", string(bytes), "ell")
    }

    _ = oswClose([]interface{}{target})
}
