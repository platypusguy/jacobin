/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"path/filepath"
)

func Load_Io_File() {

	MethodSignatures["java/io/File.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/File.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/File.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileInit,
		}

	MethodSignatures["java/io/File.getPath()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetPath,
		}

	MethodSignatures["java/io/File.delete()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileDelete,
		}

	MethodSignatures["java/io/File.createNewFile()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileCreate,
		}

	MethodSignatures["java/io/File.isInvalid()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileIsInvalid,
		}

}

// "java/io/File.<init>(Ljava/lang/String;)V"
// File file = new File(path);
func fileInit(params []interface{}) interface{} {

	// Get File object. Initialise the field map if required.
	objFile := params[0].(*object.Object)
	if objFile.FieldTable == nil {
		objFile.FieldTable = make(map[string]object.Field)
	}

	// Initialise the file status as "invalid".
	fld := object.Field{Ftype: types.Int, Fvalue: int64(0)}
	objFile.FieldTable[FileStatus] = fld

	// Get the argument path string object.
	objPath := params[1]
	if object.IsNull(objPath) {
		errMsg := "Path object is null"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}
	argPathStr := object.GoStringFromStringObject(objPath.(*object.Object))
	if argPathStr == "" {
		errMsg := "String argument for path is empty"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Create an absolute path string.
	absPathStr, err := filepath.Abs(argPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("filepath.Abs(%s) failed, reason: %s", argPathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Fill in File attributes that might get accessed by OpenJDK library member functions.

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(absPathStr)}
	objFile.FieldTable[FilePath] = fld

	fld = object.Field{Ftype: types.Int, Fvalue: os.PathSeparator}
	objFile.FieldTable["separatorChar"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte{os.PathSeparator}}
	objFile.FieldTable["separator"] = fld

	fld = object.Field{Ftype: types.Int, Fvalue: os.PathListSeparator}
	objFile.FieldTable["pathSeparatorChar"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte{os.PathListSeparator}}
	objFile.FieldTable["pathSeparator"] = fld

	// Set status to "checked" (=1).
	fld = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	objFile.FieldTable[FileStatus] = fld

	return nil
}

// "java/io/File.getPath()Ljava/lang/String;"
func fileGetPath(params []interface{}) interface{} {
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	bytes := fld.Fvalue.([]byte)
	return object.StringObjectFromByteArray(bytes)
}

// "java/io/File.isInvalid()Z"
func fileIsInvalid(params []interface{}) interface{} {
	status, ok := params[0].(*object.Object).FieldTable[FileStatus].Fvalue.(int64)
	if !ok {
		errMsg := "File object lacks a FileStatus field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	if status == 0 {
		return int64(1)
	} else {
		return int64(0)
	}
}

// "java/io/File.delete()Z"
func fileDelete(params []interface{}) interface{} {
	// Close the file if it is open (Windows).
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if ok {
		_ = osFile.Close()
	}

	// Get file path string.
	bytes, ok := params[0].(*object.Object).FieldTable[FilePath].Fvalue.([]byte)
	if !ok {
		errMsg := "File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	pathStr := string(bytes)

	err := os.Remove(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileDelete: Failed to remove file %s, reason: %s", pathStr, err.Error())
		trace.Error(errMsg)
		return int64(0)
	}
	return int64(1)
}

// "java/io/File.createNewFile()Z"
func fileCreate(params []interface{}) interface{} {
	// Get path string.
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	pathStr := string(fld.Fvalue.([]byte))

	// Create the file.
	osFile, err := os.Create(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileCreate: Failed to create file %s, reason: %s", pathStr, err.Error())
		trace.Error(errMsg)
		return int64(0)
	}

	// Copy the file handle into the FileOutputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return int64(1)
}
