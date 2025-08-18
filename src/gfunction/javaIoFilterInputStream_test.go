/*
 * Jacobin VM - A Java virtual machine
 * Tests for javaIoFilterInputStream.go
 */
package gfunction

import (
    "os"
    "path/filepath"
    "testing"

    "jacobin/src/excNames"
    "jacobin/src/globals"
    "jacobin/src/object"
    "jacobin/src/types"
)

func makeTempFileFIS(t *testing.T, content []byte) (string, func()) {
    t.Helper()
    tmpDir := t.TempDir()
    tmpFile := filepath.Join(tmpDir, "filter_in_test.txt")
    if err := os.WriteFile(tmpFile, content, 0o644); err != nil {
        t.Fatalf("failed to create temp file: %v", err)
    }
    return tmpFile, func() { _ = os.Remove(tmpFile) }
}

func newJavaFileObjPath(path string) *object.Object {
    return &object.Object{FieldTable: map[string]object.Field{
        FilePath: {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(path)},
    }}
}

func newFilterInputStreamObj() *object.Object {
    return &object.Object{FieldTable: make(map[string]object.Field)}
}

func newJavaByteArrayObj(size int) *object.Object {
    jb := make([]types.JavaByte, size)
    return &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.ByteArray, Fvalue: jb},
    }}
}

func mustOpenFileFIS(t *testing.T, path string) *os.File {
    t.Helper()
    f, err := os.Open(path)
    if err != nil { t.Fatalf("open %s: %v", path, err) }
    return f
}

func TestFilterInputStream_Init_WithFile_Success(t *testing.T) {
    globals.InitStringPool()

    content := []byte("hello filter")
    path, _ := makeTempFileFIS(t, content)

    fileObj := newJavaFileObjPath(path)
    fis := newFilterInputStreamObj()

    if res := initFilterInputStreamFile([]interface{}{fis, fileObj}); res != nil {
        t.Fatalf("initFilterInputStreamFile error: %v", res)
    }

    // FilePath copied
    fld, ok := fis.FieldTable[FilePath]
    if !ok { t.Fatalf("FilePath not set") }
    got := string(object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte)))
    if got != path { t.Fatalf("FilePath mismatch: got %q want %q", got, path) }

    // Handle present
    if _, ok := fis.FieldTable[FileHandle].Fvalue.(*os.File); !ok {
        t.Fatalf("FileHandle not *os.File")
    }
    _ = fis.FieldTable[FileHandle].Fvalue.(*os.File).Close()
}

func TestFilterInputStream_Init_WithString_Success_And_NotFound(t *testing.T) {
    globals.InitStringPool()

    content := []byte("abc")
    path, _ := makeTempFileFIS(t, content)

    // Success path
    fis := newFilterInputStreamObj()
    if res := initFilterInputStreamString([]interface{}{fis, object.StringObjectFromGoString(path)}); res != nil {
        t.Fatalf("initFilterInputStreamString error: %v", res)
    }
    if _, ok := fis.FieldTable[FileHandle].Fvalue.(*os.File); !ok {
        t.Fatalf("FileHandle not set")
    }
    _ = fis.FieldTable[FileHandle].Fvalue.(*os.File).Close()

    // Not found
    bad := filepath.Join(t.TempDir(), "no_such_file.txt")
    fis2 := newFilterInputStreamObj()
    res := initFilterInputStreamString([]interface{}{fis2, object.StringObjectFromGoString(bad)})
    if _, ok := res.(*GErrBlk); !ok {
        t.Fatalf("expected error for nonexistent file, got %T", res)
    }
}

func TestFilterInputStream_Read_Available_Skip_Close(t *testing.T) {
    globals.InitStringPool()

    content := []byte("ABCDEFGHIJ")
    path, _ := makeTempFileFIS(t, content)

    obj := newFilterInputStreamObj()
    obj.FieldTable[FileHandle] = object.Field{Ftype: types.FileHandle, Fvalue: mustOpenFileFIS(t, path)}

    // available should be > 0
    if v := fisAvailable([]interface{}{obj}); v == nil {
        t.Fatalf("available returned nil")
    } else if n, ok := v.(int64); !ok {
        t.Fatalf("available type %T", v)
    } else if n <= 0 {
        t.Fatalf("available <= 0: %d", n)
    }

    // read() one byte
    r := fisReadOne([]interface{}{obj})
    if _, ok := r.(int64); !ok {
        t.Fatalf("read() did not return int64, got %T", r)
    }

    // read([B) into buffer
    bufObj := newJavaByteArrayObj(4)
    if v := fisReadByteArray([]interface{}{obj, bufObj}); v == nil {
        t.Fatalf("read([B) returned nil")
    } else if n, ok := v.(int64); !ok || n <= 0 {
        t.Fatalf("read([B) invalid result: %v", v)
    }

    // read([B,off,len) with bounds
    big := newJavaByteArrayObj(10)
    if v := fisReadByteArrayOffset([]interface{}{obj, big, int64(3), int64(4)}); v == nil {
        t.Fatalf("read([BII) returned nil")
    } else if n, ok := v.(int64); !ok || n <= 0 {
        t.Fatalf("read([BII) invalid result: %v", v)
    }

    // invalid bounds -> IndexOutOfBoundsException
    inv := fisReadByteArrayOffset([]interface{}{obj, big, int64(20), int64(5)})
    if geb, ok := inv.(*GErrBlk); !ok || geb.ExceptionType != excNames.IndexOutOfBoundsException {
        t.Fatalf("expected IndexOutOfBoundsException, got %T", inv)
    }

    // skip some bytes
    if v := fisSkip([]interface{}{obj, int64(2)}); v == nil {
        t.Fatalf("skip returned nil")
    } else if n, ok := v.(int64); !ok || n != 2 {
        t.Fatalf("skip expected 2, got %v", v)
    }

    // close
    if res := fisClose([]interface{}{obj}); res != nil {
        t.Fatalf("close error: %v", res)
    }
}

func TestFilterInputStream_MarkSupported_False(t *testing.T) {
    globals.InitStringPool()
    // markSupported is shared with BufferedReader impl returning false
    if v := bufferedReaderMarkSupported([]interface{}{}); v.(int64) != 0 {
        t.Fatalf("markSupported expected false (0), got %v", v)
    }
}

func TestFilterInputStream_Init_FilePathMissing_Error(t *testing.T) {
    globals.InitStringPool()

    fis := newFilterInputStreamObj()
    badFile := &object.Object{FieldTable: map[string]object.Field{}}
    res := initFilterInputStreamFile([]interface{}{fis, badFile})
    if geb, ok := res.(*GErrBlk); !ok || geb.ExceptionType != excNames.IOException {
        t.Fatalf("expected IOException for missing FilePath, got %T", res)
    }
}
