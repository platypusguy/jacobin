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
	"jacobin/types"
	"os"
	"path/filepath"
)

func Load_Io_File() map[string]GMeth {

	MethodSignatures["java/io/File.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
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

	return MethodSignatures
}

// "java/io/File.<init>(Ljava/lang/String;)V"
// File file = new File(path);
func fileInit(params []interface{}) interface{} {

	// Initialise the status as "invalid".
	fld := object.Field{Ftype: types.Int, Fvalue: int64(0)}
	params[0].(*object.Object).FieldTable[FileStatus] = fld

	// Get the argument path string.
	argPathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	if argPathStr == "" {
		errMsg := "fileInit: String argument for path is null"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Create an absolute path string.
	absPathStr, err := filepath.Abs(argPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileInit: filepath.Abs(%s) returned: %s", argPathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Fill in File attributes that might get accessed by OpenJDK library member functions.

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(absPathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	fld = object.Field{Ftype: types.Int, Fvalue: os.PathSeparator}
	params[0].(*object.Object).FieldTable["separatorChar"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte{os.PathSeparator}}
	params[0].(*object.Object).FieldTable["separator"] = fld

	fld = object.Field{Ftype: types.Int, Fvalue: os.PathListSeparator}
	params[0].(*object.Object).FieldTable["pathSeparatorChar"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte{os.PathListSeparator}}
	params[0].(*object.Object).FieldTable["pathSeparator"] = fld

	// Set status to "checked" (=1).
	fld = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	params[0].(*object.Object).FieldTable[FileStatus] = fld

	return nil
}

// "java/io/File.getPath()Ljava/lang/String;"
func fileGetPath(params []interface{}) interface{} {
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "fileGetPath: File object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	bytes := fld.Fvalue.([]byte)
	return object.StringObjectFromByteArray(bytes)
}

// "java/io/File.isInvalid()Z"
func fileIsInvalid(params []interface{}) interface{} {
	status, ok := params[0].(*object.Object).FieldTable[FileStatus].Fvalue.(int64)
	if !ok {
		errMsg := "fileIsInvalid: File object lacks a FileStatus field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	if status == 0 {
		return int64(1)
	} else {
		return int64(0)
	}
}

// "java/io/File.delete()Ljava/lang/String;"
func fileDelete(params []interface{}) interface{} {
	bytes, ok := params[0].(*object.Object).FieldTable[FilePath].Fvalue.([]byte)
	if !ok {
		errMsg := "fileDelete: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	path := string(bytes)
	err := os.Remove(path)
	if err != nil {
		return int64(0)
	}
	return int64(1)
}

// "java/io/File.createNewFile()Ljava/lang/String;"
func fileCreate(params []interface{}) interface{} {
	bytes, ok := params[0].(*object.Object).FieldTable[FilePath].Fvalue.([]byte)
	if !ok {
		errMsg := "fileCreate: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	path := string(bytes)
	_, err := os.Create(path)
	if err != nil {
		return int64(0)
	}
	return int64(1)
}
