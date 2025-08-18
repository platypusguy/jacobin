package gfunction

import (
	"os"
	"path/filepath"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestFileInit_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "file.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	pathObj := object.StringObjectFromGoString(testFile)
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}

	params := []interface{}{fileObj, pathObj}

	res := fileInit(params)
	if res != nil {
		t.Fatalf("Expected nil error, got %#v", res)
	}

	fld, ok := fileObj.FieldTable[FilePath]
	if !ok {
		t.Errorf("FilePath field missing after fileInit")
	} else {
		bytes, ok := fld.Fvalue.([]types.JavaByte)
		if !ok {
			t.Errorf("FilePath field value is not []types.JavaByte")
		} else {
			goStr := object.GoStringFromStringObject(object.StringObjectFromJavaByteArray(bytes))
			abs, err := filepath.Abs(testFile)
			if err != nil {
				t.Fatalf("filepath.Abs error in test: %v", err)
			}
			if goStr != abs {
				t.Errorf("FilePath mismatch, want %q got %q", abs, goStr)
			}
		}
	}

	statusFld, ok := fileObj.FieldTable[FileStatus]
	if !ok {
		t.Errorf("FileStatus field missing after fileInit")
	} else if statusFld.Fvalue.(int64) != 1 {
		t.Errorf("FileStatus expected 1, got %v", statusFld.Fvalue)
	}
}

func TestFileInit_NullPath(t *testing.T) {
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	params := []interface{}{fileObj, object.Null}

	res := fileInit(params)
	if res == nil {
		t.Fatal("Expected NullPointerException error, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk, got %#v", res)
	}
	if errObj.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", errObj.ExceptionType)
	}
}

func TestFileInit_EmptyPath(t *testing.T) {
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	emptyStrObj := object.StringObjectFromGoString("")
	params := []interface{}{fileObj, emptyStrObj}

	res := fileInit(params)
	if res == nil {
		t.Fatal("Expected NullPointerException error, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk, got %#v", res)
	}
	if errObj.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", errObj.ExceptionType)
	}
}

func TestFileGetPath_Success(t *testing.T) {
	pathStr := t.TempDir()
	byteArr := object.JavaByteArrayFromGoString(pathStr)
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: {Ftype: types.Array, Fvalue: byteArr},
		},
	}
	params := []interface{}{fileObj}

	res := fileGetPath(params)
	if !object.IsStringObject(res) {
		t.Fatalf("Expected StringObject, got %#v", res)
	}

	goStr := object.GoStringFromStringObject(res.(*object.Object))
	if goStr != pathStr {
		t.Errorf("Expected %q, got %q", pathStr, goStr)
	}
}

func TestFileGetPath_MissingFilePath(t *testing.T) {
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	params := []interface{}{fileObj}

	res := fileGetPath(params)
	if res == nil {
		t.Fatal("Expected IOException error, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk, got %#v", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFileIsInvalid_ZeroStatus(t *testing.T) {
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FileStatus: {Ftype: types.Int, Fvalue: int64(0)},
		},
	}
	params := []interface{}{fileObj}

	res := fileIsInvalid(params)
	if val, ok := res.(int64); !ok || val != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue (1) for invalid file status, got %#v", res)
	}
}

func TestFileIsInvalid_NonZeroStatus(t *testing.T) {
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FileStatus: {Ftype: types.Int, Fvalue: int64(1)},
		},
	}
	params := []interface{}{fileObj}

	res := fileIsInvalid(params)
	if val, ok := res.(int64); !ok || val != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse (0) for valid file status, got %#v", res)
	}
}

func TestFileIsInvalid_MissingField(t *testing.T) {
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	params := []interface{}{fileObj}

	res := fileIsInvalid(params)
	if res == nil {
		t.Fatal("Expected IOException error, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk, got %#v", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFileDelete_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "delete.txt")
	err := os.WriteFile(testFile, []byte("to delete"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString(testFile)},
		},
	}
	params := []interface{}{fileObj}

	res := fileDelete(params)
	if val, ok := res.(int64); !ok || val != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue (1) on successful delete, got %#v", res)
	}

	_, err = os.Stat(testFile)
	if !os.IsNotExist(err) {
		t.Errorf("File still exists after delete")
	}
}

func TestFileDelete_CloseFileHandle(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "delete2.txt")
	err := os.WriteFile(testFile, []byte("to delete"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	f, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath:   {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString(testFile)},
			FileHandle: {Ftype: types.Ref, Fvalue: f},
		},
	}
	params := []interface{}{fileObj}

	res := fileDelete(params)
	if val, ok := res.(int64); !ok || val != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue (1) on successful delete, got %#v", res)
	}

	_, err = os.Stat(testFile)
	if !os.IsNotExist(err) {
		t.Errorf("File still exists after delete")
	}
}

func TestFileDelete_MissingFilePath(t *testing.T) {
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	params := []interface{}{fileObj}

	res := fileDelete(params)
	if res == nil {
		t.Fatal("Expected IOException error, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk, got %#v", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFileCreate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	newFile := filepath.Join(tmpDir, "newfile.txt")

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString(newFile)},
		},
	}
	params := []interface{}{fileObj}

	res := fileCreate(params)
	if val, ok := res.(int64); !ok || val != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue (1) on successful create, got %#v", res)
	}

	fh, ok := fileObj.FieldTable[FileHandle]
	if !ok {
		t.Errorf("FileHandle field missing after fileCreate")
	}
	if _, ok := fh.Fvalue.(*os.File); !ok {
		t.Errorf("FileHandle field value is not *os.File")
	}

	// Work-around to prevent Windows from getting lost in TempDir RemoveAll cleanup
	err := fh.Fvalue.(*os.File).Close()
	if err != nil {
		t.Fatalf("Failed to close file handle: %v", err)
	}
	err = os.Remove(newFile)
	if err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}
}

func TestFileCreate_MissingFilePath(t *testing.T) {
	fileObj := &object.Object{FieldTable: make(map[string]object.Field)}
	params := []interface{}{fileObj}

	res := fileCreate(params)
	if res == nil {
		t.Fatal("Expected IOException error, got nil")
	}
	errObj, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("Expected *GErrBlk, got %#v", res)
	}
	if errObj.ExceptionType != excNames.IOException {
		t.Errorf("Expected IOException, got %v", errObj.ExceptionType)
	}
}

func TestFileCreate_Failure(t *testing.T) {
	// Attempt to create file in directory without permissions to simulate failure
	// Note: This test might require root permissions or will be skipped.
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			FilePath: {Ftype: types.Array, Fvalue: object.JavaByteArrayFromGoString("/root/forbiddenfile")},
		},
	}
	params := []interface{}{fileObj}

	res := fileCreate(params)
	if val, ok := res.(int64); !ok || val != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse (0) on failure to create file, got %#v", res)
	}
}
