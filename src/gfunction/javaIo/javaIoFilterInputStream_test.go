/*
 * Jacobin VM - A Java virtual machine
 * Tests for javaIoFilterInputStream.go
 */
package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"testing"
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
		ghelpers.FilePath: {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(path)},
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
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	return f
}

func newFISObj() *object.Object {
	className := "java/io/FileInputStream"
	return object.MakeEmptyObjectWithClassName(&className)
}

func TestFilterInputStream_Init_Success(t *testing.T) {
	globals.InitStringPool()

	in := newFISObj()
	fis := newFilterInputStreamObj()

	if res := initFilterInputStream([]interface{}{fis, in}); res != nil {
		t.Fatalf("initFilterInputStream error: %v", res)
	}

	// field "in" should be set
	fld, ok := fis.FieldTable["in"]
	if !ok {
		t.Fatalf("field 'in' not set")
	}
	if fld.Fvalue != in {
		t.Fatalf("field 'in' mismatch")
	}
}

func TestFilterInputStream_Delegation(t *testing.T) {
	globals.InitStringPool()
	Load_Io_FileInputStream()

	content := []byte("ABCDEFGHIJ")
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(path, content, 0644)

	// 1. Setup FileInputStream
	fis := newFISObj()
	pathObj := object.StringObjectFromGoString(path)
	InitFileInputStreamString([]interface{}{fis, pathObj})

	// 2. Setup FilterInputStream wrapping the FIS
	filter := newFilterInputStreamObj()
	initFilterInputStream([]interface{}{filter, fis})

	// 3. Test delegation for read()
	r := filterInputStreamRead([]interface{}{filter})
	if v, ok := r.(int64); !ok || v != int64('A') {
		t.Fatalf("read() expected 'A' (65), got %v", r)
	}

	// 4. Test delegation for available()
	v := filterInputStreamAvailable([]interface{}{filter})
	if n, ok := v.(int64); !ok || n <= 0 {
		t.Fatalf("available() invalid result: %v", v)
	}

	// 5. Test delegation for skip()
	s := filterInputStreamSkip([]interface{}{filter, int64(2)})
	if n, ok := s.(int64); !ok || n != 2 {
		t.Fatalf("skip(2) expected 2, got %v", s)
	}

	// 6. Test delegation for close()
	if res := filterInputStreamClose([]interface{}{filter}); res != nil {
		t.Fatalf("close() error: %v", res)
	}
}
