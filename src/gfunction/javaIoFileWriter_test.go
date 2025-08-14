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

// helper to make a generic target object; FileWriter/OutputStreamWriter impls only rely on FilePath/FileHandle fields
func makeEmptyObj() *object.Object {
    return object.MakeEmptyObject()
}

// helper to build a char array object expected by oswWriteCharBuffer ([C as []int64)
func makeCharArray(vals []int64) *object.Object {
    return &object.Object{FieldTable: map[string]object.Field{
        "value": {Ftype: types.CharArray, Fvalue: vals},
    }}
}

func TestFileWriter_Write_OneChar_CharBuffer_StringBuffer(t *testing.T) {
    globals.InitStringPool()

    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "fw_test1.txt")
    defer os.Remove(filePath)

    // Create FileWriter via String-path constructor (reusing FileOutputStream initializer)
    fw := makeEmptyObj()
    pathStr := object.StringObjectFromGoString(filePath)
    if ret := initFileOutputStreamString([]interface{}{fw, pathStr}); ret != nil {
        t.Fatalf("initFileOutputStreamString error: %v", ret)
    }

    // write(I)V -> 'A'
    if ret := oswWriteOneChar([]interface{}{fw, int64('A')}); ret != nil {
        t.Fatalf("write(I) error: %v", ret)
    }

    // write([CII)V -> write 'B','C'
    chars := makeCharArray([]int64{'B', 'C', 'Z'})
    if ret := oswWriteCharBuffer([]interface{}{fw, chars, int64(0), int64(2)}); ret != nil {
        t.Fatalf("write([CII) error: %v", ret)
    }

    // write(String,II)V -> write "DEF"
    sObj := object.StringObjectFromGoString("DEF")
    if ret := oswWriteStringBuffer([]interface{}{fw, sObj, int64(0), int64(3)}); ret != nil {
        t.Fatalf("write(String,II) error: %v", ret)
    }

    // flush()V and close()V
    if ret := oswFlush([]interface{}{fw}); ret != nil {
        t.Fatalf("flush error: %v", ret)
    }
    if ret := oswClose([]interface{}{fw}); ret != nil {
        t.Fatalf("close error: %v", ret)
    }

    // Verify content
    content, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("reading file failed: %v", err)
    }
    if string(content) != "ABCDEF" {
        t.Fatalf("file content mismatch: got %q", string(content))
    }
}

func TestFileWriter_Append_Mode(t *testing.T) {
    globals.InitStringPool()

    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "fw_test2.txt")
    defer os.Remove(filePath)

    // First, create file and write "X"
    fw1 := makeEmptyObj()
    if ret := initFileOutputStreamString([]interface{}{fw1, object.StringObjectFromGoString(filePath)}); ret != nil {
        t.Fatalf("init (create) error: %v", ret)
    }
    _ = oswWriteOneChar([]interface{}{fw1, int64('X')})
    _ = oswClose([]interface{}{fw1})

    // Now reopen with append=true and write "YZ"
    fw2 := makeEmptyObj()
    if ret := initFileOutputStreamStringBoolean([]interface{}{fw2, object.StringObjectFromGoString(filePath), int64(1)}); ret != nil {
        t.Fatalf("init (append) error: %v", ret)
    }
    _ = oswWriteStringBuffer([]interface{}{fw2, object.StringObjectFromGoString("YZ"), int64(0), int64(2)})
    _ = oswClose([]interface{}{fw2})

    // Verify content is "XYZ"
    content, err := os.ReadFile(filePath)
    if err != nil { t.Fatalf("read failed: %v", err) }
    if string(content) != "XYZ" {
        t.Fatalf("append content mismatch: got %q", string(content))
    }
}

func TestFileWriter_Write_ParamErrors(t *testing.T) {
    globals.InitStringPool()

    tmpDir := os.TempDir()
    filePath := filepath.Join(tmpDir, "fw_test3.txt")
    defer os.Remove(filePath)

    fw := makeEmptyObj()
    if ret := initFileOutputStreamString([]interface{}{fw, object.StringObjectFromGoString(filePath)}); ret != nil {
        t.Fatalf("init error: %v", ret)
    }

    // oswWriteCharBuffer with invalid (offset,length) -> IndexOutOfBoundsException
    chars := makeCharArray([]int64{'A', 'B'})
    if res := oswWriteCharBuffer([]interface{}{fw, chars, int64(1), int64(5)}); res == nil {
        t.Fatalf("expected error for char buffer bounds")
    } else if geb, ok := res.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    }

    // oswWriteStringBuffer with invalid (offset,length)
    sObj := object.StringObjectFromGoString("HI")
    if res := oswWriteStringBuffer([]interface{}{fw, sObj, int64(0), int64(5)}); res == nil {
        t.Fatalf("expected error for string buffer bounds")
    } else if geb, ok := res.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IndexOutOfBoundsException {
            t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
        }
    }

    _ = oswClose([]interface{}{fw})
}
