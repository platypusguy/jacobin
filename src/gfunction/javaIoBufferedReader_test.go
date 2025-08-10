package gfunction

import (
	"os"
	"path/filepath"
	"testing"

	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
)

func TestBufferedReaderInit_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello\nWorld\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString(filePath)},
		},
	}

	brObj := &object.Object{FieldTable: map[string]object.Field{}}

	params := []interface{}{brObj, fileObj}
	res := bufferedReaderInit(params)
	if res != nil {
		t.Fatalf("Expected success, got error: %v", res)
	}

	if _, ok := brObj.FieldTable[FilePath]; !ok {
		t.Errorf("FilePath field not set on BufferedReader object")
	}
	if fh, ok := brObj.FieldTable[FileHandle]; !ok {
		t.Errorf("FileHandle field not set on BufferedReader object")
	} else {
		if f, ok := fh.Fvalue.(*os.File); ok {
			_ = f.Close()
		}
	}
}

func TestBufferedReaderInit_FileNotFound(t *testing.T) {
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString("/no/such/file")},
		},
	}
	brObj := &object.Object{FieldTable: map[string]object.Field{}}

	params := []interface{}{brObj, fileObj}
	res := bufferedReaderInit(params)
	if res == nil {
		t.Errorf("Expected FileNotFoundException error, got nil")
		return
	}
	errObj, ok := res.(*GErrBlk)
	if !ok || errObj.ExceptionType != excNames.FileNotFoundException {
		t.Errorf("Expected FileNotFoundException, got %#v", res)
	}
}

func TestBufferedReaderInit_MissingFilePathField(t *testing.T) {
	fileObj := &object.Object{FieldTable: map[string]object.Field{}}
	brObj := &object.Object{FieldTable: map[string]object.Field{}}

	params := []interface{}{brObj, fileObj}
	res := bufferedReaderInit(params)
	if res == nil {
		t.Errorf("Expected InvalidTypeException error, got nil")
		return
	}
	errObj, ok := res.(*GErrBlk)
	if !ok || errObj.ExceptionType != excNames.InvalidTypeException {
		t.Errorf("Expected InvalidTypeException, got %#v", res)
	}
}

func TestBufferedReaderMarkSupported_False(t *testing.T) {
	params := []interface{}{}
	res := bufferedReaderMarkSupported(params)
	if res != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse (0), got %v", res)
	}
}

func TestBufferedReaderReadLine_FirstLine(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("line1\nline2\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	brObj := &object.Object{
		FieldTable: map[string]object.Field{
			FileHandle: {Ftype: types.Ref, Fvalue: f},
		},
	}

	params := []interface{}{brObj}
	res := bufferedReaderReadLine(params)
	if !object.IsStringObject(res) {
		t.Fatalf("Expected Java String object, got %#v", res)
	}
	goStr := object.GoStringFromStringObject(res.(*object.Object))
	if goStr != "line1" {
		t.Errorf("Expected 'line1', got %q", goStr)
	}
}

func TestBufferedReaderReadLine_AtEOF(t *testing.T) {
	brObj := &object.Object{
		FieldTable: make(map[string]object.Field),
	}
	eofSet(brObj, true)

	params := []interface{}{brObj}
	res := bufferedReaderReadLine(params)
	if res != object.Null {
		t.Errorf("Expected Null at EOF, got %#v", res)
	}
}

func TestBufferedReaderReadLine_MissingFileHandle(t *testing.T) {
	brObj := &object.Object{FieldTable: map[string]object.Field{}}
	params := []interface{}{brObj}
	res := bufferedReaderReadLine(params)
	if res == nil {
		t.Errorf("Expected IOException error, got nil")
		return
	}
	errObj, ok := res.(*GErrBlk)
	if !ok || errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %#v", res)
	}
}

func TestBufferedReaderReadLine_MultiLineSequential(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("first\nsecond\nthird")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	brObj := &object.Object{
		FieldTable: map[string]object.Field{
			FileHandle: {Ftype: types.Ref, Fvalue: f},
		},
	}
	params := []interface{}{brObj}

	expected := []string{"first", "second", "third"}
	for i, want := range expected {
		res := bufferedReaderReadLine(params)
		if res == object.Null {
			t.Fatalf("Unexpected EOF on line %d", i+1)
		}
		if !object.IsStringObject(res) {
			t.Fatalf("Expected Java String object, got %#v", res)
		}
		got := object.GoStringFromStringObject(res.(*object.Object))
		if got != want {
			t.Errorf("Line %d: expected %q, got %q", i+1, want, got)
		}
	}
}
