/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"testing"
)

func TestInitFileOutputStreamFile(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			"FilePath": {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(filePath)},
		},
	}
	emptyObj := object.MakeEmptyObject()
	emptyObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: 0} // dummy placeholder value
	params := []interface{}{emptyObj, fileObj}

	result := initFileOutputStreamFile(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it does not exist")
	} else {
		os.Remove(filePath)
	}
}

func TestInitFileOutputStreamFileBoolean(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			"FilePath": {Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(filePath)},
		},
	}
	emptyObj := object.MakeEmptyObject()
	emptyObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: 0} // dummy placeholder value
	params := []interface{}{emptyObj, fileObj, int64(1)}

	result := initFileOutputStreamFileBoolean(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it does not exist")
	} else {
		os.Remove(filePath)
	}

}

func TestInitFileOutputStreamString(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	strObj := object.StringObjectFromGoString(filePath)
	emptyObj := object.MakeEmptyObject()
	emptyObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: 0} // dummy placeholder value
	params := []interface{}{emptyObj, strObj, int64(1)}

	result := initFileOutputStreamString(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it does not exist")
	} else {
		os.Remove(filePath)
	}
}

func TestInitFileOutputStreamStringBoolean(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	strObj := object.StringObjectFromGoString(filePath)
	emptyObj := object.MakeEmptyObject()
	emptyObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: 0} // dummy placeholder value
	params := []interface{}{emptyObj, strObj, int64(1)}

	result := initFileOutputStreamStringBoolean(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it does not exist")
	} else {
		os.Remove(filePath)
	}
}

func TestFosWriteOne(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	osFile, _ := os.Create(filePath)
	defer osFile.Close()

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			"FileHandle": {Ftype: ghelpers.FileHandle, Fvalue: osFile},
		},
	}
	params := []interface{}{fileObj, int64(65)}

	result := fosWriteOne(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "A" {
		t.Errorf("Expected file content to be 'A', got %s", string(content))
	} else {
		os.Remove(filePath)
	}
}

func TestFosWriteByteArray(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	osFile, _ := os.Create(filePath)
	defer osFile.Close()

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			"FileHandle": {Ftype: ghelpers.FileHandle, Fvalue: osFile},
		},
	}
	byteArray := []types.JavaByte{65, 66, 67}
	byteArrayObj := &object.Object{
		FieldTable: map[string]object.Field{
			"value": {Ftype: types.ByteArray, Fvalue: byteArray},
		},
	}
	params := []interface{}{fileObj, byteArrayObj}

	result := fosWriteByteArray(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "ABC" {
		t.Errorf("Expected file content to be 'ABC', got %s", string(content))
	} else {
		os.Remove(filePath)
	}
}

func TestFosWriteByteArrayOffset(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	osFile, _ := os.Create(filePath)
	defer osFile.Close()

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			"FileHandle": {Ftype: ghelpers.FileHandle, Fvalue: osFile},
		},
	}
	byteArray := []types.JavaByte{65, 66, 67, 68, 69}
	byteArrayObj := &object.Object{
		FieldTable: map[string]object.Field{
			"value": {Ftype: types.ByteArray, Fvalue: byteArray},
		},
	}
	params := []interface{}{fileObj, byteArrayObj, int64(1), int64(3)}

	result := fosWriteByteArrayOffset(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	content, _ := os.ReadFile(filePath)
	if string(content) != "BCD" {
		t.Errorf("Expected file content to be 'BCD', got %s", string(content))
	} else {
		os.Remove(filePath)
	}
}

func TestFosClose(t *testing.T) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")
	osFile, _ := os.Create(filePath)

	fileObj := &object.Object{
		FieldTable: map[string]object.Field{
			"FileHandle": {Ftype: ghelpers.FileHandle, Fvalue: osFile},
		},
	}
	params := []interface{}{fileObj}

	result := fosClose(params)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	err := osFile.Close()
	if err == nil {
		t.Errorf("Expected error on closing already closed file, got nil")
	} else {
		os.Remove(filePath)
	}
}
