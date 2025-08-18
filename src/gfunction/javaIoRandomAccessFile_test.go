package gfunction

import (
	"io"
	"jacobin/src/globals"
	"os"
	"testing"

	"jacobin/src/object"
	"jacobin/src/types"
)

// helper to create a new RandomAccessFile object with initialized FieldTable
func newRAFObject() *object.Object {
	return &object.Object{FieldTable: make(map[string]object.Field)}
}

func TestClinitGeneric(t *testing.T) {
	ret := clinitGeneric(nil)
	if ret != nil {
		t.Errorf("clinitGeneric should return nil, got %v", ret)
	}
}

func TestJustReturn(t *testing.T) {
	ret := justReturn(nil)
	if ret != nil {
		t.Errorf("justReturn should return nil, got %v", ret)
	}
}

func TestRafInitStringAndGetFilePointer(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rafObj := newRAFObject()

	pathStrObj := object.StringObjectFromGoString(tmpFile.Name())
	modeStrObj := object.StringObjectFromGoString("r")

	params := []interface{}{rafObj, pathStrObj, modeStrObj}
	ret := rafInitString(params)
	if ret != nil {
		t.Fatalf("rafInitString returned error: %v", ret)
	}

	fld, ok := rafObj.FieldTable[FilePath]
	if !ok {
		t.Fatal("FilePath field not set")
	}
	gotPath := string(fld.Fvalue.([]byte))
	if gotPath != tmpFile.Name() {
		t.Fatalf("FilePath mismatch, want %s, got %s", tmpFile.Name(), gotPath)
	}

	fld, ok = rafObj.FieldTable[FileHandle]
	if !ok {
		t.Fatal("FileHandle field not set")
	}
	fh, ok := fld.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("FileHandle field has wrong type %T", fld.Fvalue)
	}

	getPointerParams := []interface{}{rafObj}
	pos := rafGetFilePointer(getPointerParams)
	offset, ok := pos.(int64)
	if !ok {
		t.Fatalf("rafGetFilePointer returned wrong type %T", pos)
	}
	if offset != 0 {
		t.Errorf("Initial file pointer expected 0, got %d", offset)
	}

	fh.Close()
}

func TestRafInitFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	fileObj.FieldTable[FilePath] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: object.JavaByteArrayFromGoString(tmpFile.Name()),
	}

	rafObj := newRAFObject()

	modeStrObj := object.StringObjectFromGoString("r")

	params := []interface{}{rafObj, fileObj, modeStrObj}
	ret := rafInitFile(params)
	if ret != nil {
		t.Fatalf("rafInitFile returned error: %v", ret)
	}

	fld, ok := rafObj.FieldTable[FilePath]
	if !ok {
		t.Fatal("FilePath field not set")
	}
	gotPath := string(fld.Fvalue.([]byte))
	if gotPath != tmpFile.Name() {
		t.Fatalf("FilePath mismatch, want %s, got %s", tmpFile.Name(), gotPath)
	}

	fld, ok = rafObj.FieldTable[FileHandle]
	if !ok {
		t.Fatal("FileHandle field not set")
	}
	fh, ok := fld.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("FileHandle field has wrong type %T", fld.Fvalue)
	}

	fh.Close()
}

func TestFisClose(t *testing.T) {
	rafObj := newRAFObject()

	// Set FileHandle with a pipe writer to avoid closing os.Stdout accidentally
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	rafObj.FieldTable[FileHandle] = object.Field{Ftype: types.FileHandle, Fvalue: w}

	ret := fisClose([]interface{}{rafObj})
	if ret != nil {
		t.Errorf("fisClose returned error: %v", ret)
	}

	// Writing after close should fail
	_, err = w.Write([]byte("test"))
	if err == nil {
		t.Errorf("Write succeeded after close, expected failure")
	}
}

func TestFisReadOne(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "raf_read_one")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte{0x42}
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[FileHandle] = object.Field{Ftype: types.FileHandle, Fvalue: tmpFile}

	ret := fisReadOne([]interface{}{rafObj})

	intRet, ok := ret.(int64)
	if !ok {
		t.Fatalf("fisReadOne returned wrong type %T", ret)
	}
	if intRet != int64(content[0]) {
		t.Errorf("fisReadOne expected %d, got %d", content[0], intRet)
	}
}

func TestFisReadByteArray(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_ba")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[FileHandle] = object.Field{Ftype: types.FileHandle, Fvalue: tmpFile}

	byteArray := make([]types.JavaByte, len(content))

	javaByteArrayObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, byteArray)
	params := []interface{}{rafObj, javaByteArrayObj}
	ret := fisReadByteArray(params)

	numRead, ok := ret.(int64)
	if !ok {
		t.Fatalf("fisReadByteArray returned wrong type %T", ret)
	}
	if numRead != int64(len(content)) {
		t.Errorf("fisReadByteArray expected read %d bytes, got %d", len(content), numRead)
	}
}

func TestFisReadByteArrayOffset(t *testing.T) {
	globals.InitStringPool()
	tmpFile, err := os.CreateTemp("", "raf_read_ba_offset")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("hello world")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	tmpFile.Sync()
	tmpFile.Seek(0, io.SeekStart)

	rafObj := newRAFObject()
	rafObj.FieldTable[FileHandle] = object.Field{Ftype: types.FileHandle, Fvalue: tmpFile}

	byteArray := make([]types.JavaByte, len(content))
	javaByteArrayObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, byteArray)

	offset := int64(2)
	length := int64(5)

	params := []interface{}{rafObj, javaByteArrayObj, offset, length}
	ret := fisReadByteArrayOffset(params)

	numRead, ok := ret.(int64)
	if !ok {
		t.Fatalf("fisReadByteArrayOffset returned wrong type %T", ret)
	}
	if numRead != length {
		t.Errorf("fisReadByteArrayOffset expected read %d bytes, got %d", length, numRead)
	}
}
