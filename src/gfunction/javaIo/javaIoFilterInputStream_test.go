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
		ghelpers.FilePath: {Ftype: types.JavaByteArray, Fvalue: object.JavaByteArrayFromGoString(path)},
	}}
}

func newFilterInputStreamObj() *object.Object {
	return &object.Object{FieldTable: make(map[string]object.Field)}
}

func newJavaByteArrayObj(size int) *object.Object {
	jb := make([]types.JavaByte, size)
	return &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: types.JavaByteArray, Fvalue: jb},
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

func newBAISObj() *object.Object {
	className := "java/io/ByteArrayInputStream"
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

func TestFilterInputStream_MarkResetMarkSupported(t *testing.T) {
	globals.InitStringPool()
	Load_Io_ByteArrayInputStream()

	content := []byte("ABCDEFGHIJ")
	bais := newBAISObj()
	bufObj := &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: types.JavaByteArray, Fvalue: object.JavaByteArrayFromGoString(string(content))},
	}}
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	filter := newFilterInputStreamObj()
	initFilterInputStream([]interface{}{filter, bais})

	// markSupported() should delegate to true
	if ms := filterInputStreamMarkSupported([]interface{}{filter}); ms != types.JavaBoolTrue {
		t.Fatalf("markSupported() expected true, got %v", ms)
	}

	// mark(5) should delegate without error
	if res := filterInputStreamMark([]interface{}{filter, int64(5)}); res != nil {
		t.Fatalf("mark(5) error: %v", res)
	}

	// consume a couple of bytes, then reset()
	filterInputStreamRead([]interface{}{filter})
	filterInputStreamRead([]interface{}{filter})

	if res := filterInputStreamReset([]interface{}{filter}); res != nil {
		t.Fatalf("reset() error: %v", res)
	}
}

func TestFilterInputStream_ReadByteArrayOffset(t *testing.T) {
	globals.InitStringPool()
	Load_Io_ByteArrayInputStream()

	content := []byte("ABCDEFGHIJ")
	bais := newBAISObj()
	bufObj := &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: types.JavaByteArray, Fvalue: object.JavaByteArrayFromGoString(string(content))},
	}}
	ByteArrayInputStreamInit([]interface{}{bais, bufObj})

	filter := newFilterInputStreamObj()
	initFilterInputStream([]interface{}{filter, bais})

	dst := newJavaByteArrayObj(4)
	r := filterInputStreamReadByteArrayOffset([]interface{}{filter, dst, int64(0), int64(4)})
	if n, ok := r.(int64); !ok || n != 4 {
		t.Fatalf("read([BII)I expected 4, got %v", r)
	}
}

func TestFilterInputStream_ReadByteArray(t *testing.T) {
	globals.InitStringPool()
	Load_Io_FileInputStream()

	content := []byte("ABCDEFGHIJ")
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test2.txt")
	os.WriteFile(path, content, 0644)

	fis := newFISObj()
	pathObj := object.StringObjectFromGoString(path)
	InitFileInputStreamString([]interface{}{fis, pathObj})

	filter := newFilterInputStreamObj()
	initFilterInputStream([]interface{}{filter, fis})

	dst := newJavaByteArrayObj(4)
	r := filterInputStreamReadByteArray([]interface{}{filter, dst})
	if n, ok := r.(int64); !ok || n != 4 {
		t.Fatalf("read([B)I expected 4, got %v", r)
	}
}

func TestFilterInputStream_LoadRegistersMethods(t *testing.T) {
	globals.InitStringPool()
	Load_Io_FilterInputStream()

	expected := []string{
		"java/io/FilterInputStream.<clinit>()V",
		"java/io/FilterInputStream.<init>(Ljava/io/InputStream;)V",
		"java/io/FilterInputStream.available()I",
		"java/io/FilterInputStream.close()V",
		"java/io/FilterInputStream.mark(I)V",
		"java/io/FilterInputStream.markSupported()Z",
		"java/io/FilterInputStream.read()I",
		"java/io/FilterInputStream.read([B)I",
		"java/io/FilterInputStream.read([BII)I",
		"java/io/FilterInputStream.reset()V",
		"java/io/FilterInputStream.skip(J)J",
	}
	for _, sig := range expected {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Fatalf("expected method signature %q to be registered", sig)
		}
	}
}
