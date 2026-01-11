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

// helper to create a java/io/File-like object carrying ghelpers.FilePath bytes
func makeJavaFileObj(path string) *object.Object {
	return &object.Object{FieldTable: map[string]object.Field{
		ghelpers.FilePath: {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(path)},
	}}
}

// helper to make an empty target object (FileReader target)
func makeTargetObj() *object.Object { return object.MakeEmptyObject() }

func TestFileReader_Init_WithFile_Success_And_Read(t *testing.T) {
	globals.InitStringPool()

	// Prepare a temp file with known content
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "fr_test1.txt")
	content := []byte("Hello, FileReader!")
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(filePath)

	fileObj := makeJavaFileObj(filePath)
	target := makeTargetObj()

	// Call constructor that takes a File object
	res := initFileReader([]interface{}{target, fileObj})
	if res != nil {
		t.Fatalf("initFileReader returned error: %v", res)
	}

	// Verify ghelpers.FilePath copied
	gotPathBytes := target.FieldTable[ghelpers.FilePath].Fvalue.([]types.JavaByte)
	if string(object.GoByteArrayFromJavaByteArray(gotPathBytes)) != filePath {
		t.Fatalf("ghelpers.FilePath not copied correctly")
	}

	// Verify we can read via the stored file handle
	fh, ok := target.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok || fh == nil {
		t.Fatalf("ghelpers.FileHandle is not *os.File: %T", target.FieldTable[ghelpers.FileHandle].Fvalue)
	}
	// Read a few bytes and compare to content
	buf := make([]byte, len(content))
	n, err := fh.ReadAt(buf, 0)
	if err != nil && err.Error() != "EOF" { // ReadAt may return EOF exactly at end
		if n == 0 {
			t.Fatalf("reading from FileReader handle failed: %v", err)
		}
	}
	if string(buf[:n]) != string(content[:n]) {
		t.Fatalf("read content mismatch: got %q want prefix of %q", string(buf[:n]), string(content))
	}
	// Close the handle opened by FileReader
	_ = fh.Close()
}

func TestFileReader_Init_WithString_Success_And_Read(t *testing.T) {
	globals.InitStringPool()

	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "fr_test2.txt")
	content := []byte("Another read!")
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(filePath)

	target := makeTargetObj()
	pathObj := object.StringObjectFromGoString(filePath)

	res := initFileReaderString([]interface{}{target, pathObj})
	if res != nil {
		t.Fatalf("initFileReaderString returned error: %v", res)
	}

	// Verify ghelpers.FilePath set on target
	if target.FieldTable[ghelpers.FilePath].Ftype != types.ByteArray {
		t.Fatalf("ghelpers.FilePath field type unexpected: %v", target.FieldTable[ghelpers.FilePath].Ftype)
	}

	// Verify handle is usable
	fh, ok := target.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok || fh == nil {
		t.Fatalf("ghelpers.FileHandle is not *os.File: %T", target.FieldTable[ghelpers.FileHandle].Fvalue)
	}
	b := make([]byte, len(content))
	n, err := fh.ReadAt(b, 0)
	if err != nil && err.Error() != "EOF" {
		if n == 0 {
			t.Fatalf("reading from handle failed: %v", err)
		}
	}
	if string(b[:n]) != string(content[:n]) {
		t.Fatalf("content mismatch: got %q want prefix of %q", string(b[:n]), string(content))
	}
	_ = fh.Close()
}

func TestFileReader_Error_Cases(t *testing.T) {
	globals.InitStringPool()

	// Missing ghelpers.FilePath on File object -> InvalidTypeException
	badFileObj := object.MakeEmptyObject()
	target := makeTargetObj()
	if res := initFileReader([]interface{}{target, badFileObj}); res == nil {
		t.Fatalf("expected error for missing ghelpers.FilePath field")
	} else if geb, ok := res.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.InvalidTypeException {
			t.Fatalf("expected InvalidTypeException, got %d", geb.ExceptionType)
		}
	}

	// Nonexistent file via File parameter -> FileNotFoundException
	target2 := makeTargetObj()
	nfPath := filepath.Join(os.TempDir(), "fr_no_such_file.txt")
	nfFileObj := makeJavaFileObj(nfPath)
	if res := initFileReader([]interface{}{target2, nfFileObj}); res == nil {
		t.Fatalf("expected FileNotFoundException for nonexistent file (File)")
	} else if geb, ok := res.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.FileNotFoundException {
			t.Fatalf("expected FileNotFoundException, got %d", geb.ExceptionType)
		}
	}

	// Nonexistent file via String parameter -> FileNotFoundException
	target3 := makeTargetObj()
	if res := initFileReaderString([]interface{}{target3, object.StringObjectFromGoString(nfPath)}); res == nil {
		t.Fatalf("expected FileNotFoundException for nonexistent file (String)")
	} else if geb, ok := res.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.FileNotFoundException {
			t.Fatalf("expected FileNotFoundException, got %d", geb.ExceptionType)
		}
	}
}
