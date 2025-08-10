/*
 * Jacobin VM - A Java virtual machine
 * Tests for javaIoFileInputStream.go
 * Generated according to user rules.
 */

package gfunction

import (
	"os"
	"path/filepath"
	"testing"

	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
)

func makeTempFile(t *testing.T, content []byte) (string, func()) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "testfile.txt")
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile, func() { os.Remove(tmpFile) }
}

func newFileObjectWithPath(t *testing.T, path string) *object.Object {
	fld := object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(path)}
	obj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: fld,
		},
	}
	return obj
}

func newFileInputStreamObject() *object.Object {
	return &object.Object{
		FieldTable: make(map[string]object.Field),
	}
}

func newJavaByteArrayObject(size int) *object.Object {
	ba := make([]types.JavaByte, size)
	return &object.Object{
		FieldTable: map[string]object.Field{
			"value": {Ftype: types.ByteArray, Fvalue: ba},
		},
	}
}

func TestInitFileInputStreamFile_Success(t *testing.T) {
	content := []byte("hello world")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fileObj := newFileObjectWithPath(t, path)
	fisObj := newFileInputStreamObject()

	params := []interface{}{fisObj, fileObj}

	res := initFileInputStreamFile(params)
	if res != nil {
		t.Fatalf("Expected nil, got error: %v", res)
	}

	// Check FilePath copied
	fld, ok := fisObj.FieldTable[FilePath]
	if !ok {
		t.Fatalf("FilePath not set in FileInputStream object")
	}
	gotPath := string(object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte)))
	if gotPath != path {
		t.Errorf("FilePath mismatch: got %q want %q", gotPath, path)
	}

	// Check FileHandle is set and valid
	fh, ok := fisObj.FieldTable[FileHandle]
	if !ok {
		t.Fatalf("FileHandle not set in FileInputStream object")
	}
	if fh.Ftype != types.FileHandle {
		t.Errorf("FileHandle field type mismatch: got %v want %v", fh.Ftype, types.FileHandle)
	}
	fileHandle, ok := fh.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("FileHandle Fvalue is not *os.File")
	}
	fileHandle.Close()
}

func TestInitFileInputStreamFile_FilePathMissing(t *testing.T) {
	fisObj := newFileInputStreamObject()
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)} // No FilePath

	params := []interface{}{fisObj, fileObj}
	res := initFileInputStreamFile(params)
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestInitFileInputStreamString_Success(t *testing.T) {
	content := []byte("file input stream test")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	strObj := object.StringObjectFromGoString(path)
	fisObj := newFileInputStreamObject()
	params := []interface{}{fisObj, strObj}

	res := initFileInputStreamString(params)
	if res != nil {
		t.Fatalf("Expected nil, got error: %v", res)
	}

	// Check FilePath set
	fld, ok := fisObj.FieldTable[FilePath]
	if !ok {
		t.Fatalf("FilePath not set in FileInputStream object")
	}
	gotPath := string(object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte)))
	if gotPath != path {
		t.Errorf("FilePath mismatch: got %q want %q", gotPath, path)
	}

	// Check FileHandle is set and valid
	fh, ok := fisObj.FieldTable[FileHandle]
	if !ok {
		t.Fatalf("FileHandle not set in FileInputStream object")
	}
	fileHandle, ok := fh.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("FileHandle Fvalue is not *os.File")
	}
	fileHandle.Close()
}

func TestInitFileInputStreamString_FileNotFound(t *testing.T) {
	fisObj := newFileInputStreamObject()
	badPath := "/nonexistent/file/hopefully/not/there.txt"
	strObj := object.StringObjectFromGoString(badPath)
	params := []interface{}{fisObj, strObj}

	res := initFileInputStreamString(params)
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFisAvailable_Success(t *testing.T) {
	content := []byte("hello available test")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fisObj := newFileInputStreamObject()
	fisObj.FieldTable[FileHandle] = object.Field{
		Ftype:  types.FileHandle,
		Fvalue: mustOpenFile(t, path),
	}
	defer fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()

	res := fisAvailable([]interface{}{fisObj})
	avail, ok := res.(int64)
	if !ok {
		t.Fatalf("Expected int64, got %T", res)
	}
	if avail <= 0 {
		t.Errorf("Expected positive available bytes, got %d", avail)
	}
}

func mustOpenFile(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", path, err)
	}
	return f
}

func TestFisAvailable_NoFileHandle(t *testing.T) {
	fisObj := newFileInputStreamObject()
	res := fisAvailable([]interface{}{fisObj})
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFisReadOne_SuccessAndEOF(t *testing.T) {
	content := []byte("a")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fisObj := newFileInputStreamObject()
	fisObj.FieldTable[FileHandle] = object.Field{
		Ftype:  types.FileHandle,
		Fvalue: mustOpenFile(t, path),
	}
	defer fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()

	res := fisReadOne([]interface{}{fisObj})
	b, ok := res.(int64)
	if !ok {
		t.Fatalf("Expected int64, got %T", res)
	}
	if b != int64(content[0]) {
		t.Errorf("Expected byte %d, got %d", content[0], b)
	}

	// Read until EOF (skip the single byte read previously)
	// Close and reopen file to reset
	fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()
	fld := fisObj.FieldTable[FileHandle]
	fld.Fvalue = mustOpenFile(t, path)
	fisObj.FieldTable[FileHandle] = fld

	// Read all bytes
	for i := 0; i < len(content); i++ {
		fisReadOne([]interface{}{fisObj})
	}

	// Now read again, should return -1 at EOF
	res = fisReadOne([]interface{}{fisObj})
	val, ok := res.(int64)
	if !ok {
		t.Fatalf("Expected int64, got %T", res)
	}
	if val != int64(-1) {
		t.Errorf("Expected -1 at EOF, got %d", val)
	}

	// Work-around to prevent Windows from getting lost in TempDir RemoveAll cleanup
	err := fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()
	if err != nil {
		t.Fatalf("Failed to close file handle: %v", err)
	}
	err = os.Remove(path)
	if err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}

}

func TestFisReadOne_NoFileHandle(t *testing.T) {
	fisObj := newFileInputStreamObject()
	res := fisReadOne([]interface{}{fisObj})
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFisReadByteArray_SuccessAndEOF(t *testing.T) {
	content := []byte("abcdefg")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fisObj := newFileInputStreamObject()
	fisObj.FieldTable[FileHandle] = object.Field{
		Ftype:  types.FileHandle,
		Fvalue: mustOpenFile(t, path),
	}
	defer fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()

	javaByteArrayObj := newJavaByteArrayObject(10)
	params := []interface{}{fisObj, javaByteArrayObj}

	res := fisReadByteArray(params)
	n, ok := res.(int64)
	if !ok {
		t.Fatalf("Expected int64, got %T", res)
	}
	if n <= 0 {
		t.Errorf("Expected positive bytes read, got %d", n)
	}
}

func TestFisReadByteArrayOffset_Success(t *testing.T) {
	content := []byte("1234567890")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fisObj := newFileInputStreamObject()
	fisObj.FieldTable[FileHandle] = object.Field{
		Ftype:  types.FileHandle,
		Fvalue: mustOpenFile(t, path),
	}
	defer fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()

	javaByteArrayObj := newJavaByteArrayObject(20)
	offset := int64(5)
	length := int64(4)

	params := []interface{}{fisObj, javaByteArrayObj, offset, length}

	res := fisReadByteArrayOffset(params)
	n, ok := res.(int64)
	if !ok {
		t.Fatalf("Expected int64, got %T", res)
	}
	if n <= 0 {
		t.Errorf("Expected positive bytes read, got %d", n)
	}

	// Test invalid offset and length (too large)
	paramsInvalid := []interface{}{fisObj, javaByteArrayObj, int64(1000), int64(10)}
	resInvalid := fisReadByteArrayOffset(paramsInvalid)
	errObj, ok := resInvalid.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", resInvalid)
	}
	if errObj.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Errorf("Expected IndexOutOfBoundsException, got %v", errObj.ExceptionType)
	}
}

func TestFisSkip_Success(t *testing.T) {
	content := []byte("1234567890")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fisObj := newFileInputStreamObject()
	fisObj.FieldTable[FileHandle] = object.Field{
		Ftype:  types.FileHandle,
		Fvalue: mustOpenFile(t, path),
	}
	defer fisObj.FieldTable[FileHandle].Fvalue.(*os.File).Close()

	skipCount := int64(5)
	params := []interface{}{fisObj, skipCount}

	res := fisSkip(params)
	n, ok := res.(int64)
	if !ok {
		t.Fatalf("Expected int64, got %T", res)
	}
	if n != skipCount {
		t.Errorf("Expected skip count %d, got %d", skipCount, n)
	}
}

func TestFisSkip_NoFileHandle(t *testing.T) {
	fisObj := newFileInputStreamObject()
	params := []interface{}{fisObj, int64(5)}
	res := fisSkip(params)
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFisClose_Success(t *testing.T) {
	content := []byte("close test")
	path, cleanup := makeTempFile(t, content)
	defer cleanup()

	fisObj := newFileInputStreamObject()
	fisObj.FieldTable[FileHandle] = object.Field{
		Ftype:  types.FileHandle,
		Fvalue: mustOpenFile(t, path),
	}

	res := fisClose([]interface{}{fisObj})
	if res != nil {
		t.Fatalf("Expected nil, got error %v", res)
	}
}

func TestFisClose_NoFileHandle(t *testing.T) {
	fisObj := newFileInputStreamObject()
	res := fisClose([]interface{}{fisObj})
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk error, got %T", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}
