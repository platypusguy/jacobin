/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"os"
	"path/filepath"
)

func Load_Io_File() map[string]GMeth {

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

	MethodSignatures["java/io/File.isInvalid()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileIsInvalid,
		}

	// -------------------
	// <clinit> justReturn
	// -------------------

	MethodSignatures["java/io/File.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// -----------------------------------------
	// Traps that do nothing but return an error
	// -----------------------------------------

	MethodSignatures["java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapGetDefaultFileSystem,
		}

	return MethodSignatures
}

// "java/io/File.<init>(Ljava/lang/String;)V"
// File file = new File(path);
func fileInit(params []interface{}) interface{} {
	inPathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	if inPathStr == "" {
		errMsg := "fileInit: String argument for path is null"
		return getGErrBlk(exceptions.NullPointerException, errMsg)
	}
	absPathStr, err := filepath.Abs(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileInit: filepath.Abs(%s) returned: %s", inPathStr, err.Error())
		return getGErrBlk(exceptions.FileSystemNotFoundException, errMsg)
	}

	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(absPathStr)}
	params[0].(*object.Object).FieldTable["filePath"] = fld

	fld.Ftype = types.Int // char
	fld.Fvalue = os.PathSeparator
	params[0].(*object.Object).FieldTable["separatorChar"] = fld

	fld.Ftype = types.ByteArray // string version of separatorChar
	fld.Fvalue = []byte{os.PathSeparator}
	params[0].(*object.Object).FieldTable["separator"] = fld

	fld.Ftype = types.Int // char
	fld.Fvalue = os.PathListSeparator
	params[0].(*object.Object).FieldTable["pathSeparatorChar"] = fld

	fld.Ftype = types.ByteArray // string version of separatorChar
	fld.Fvalue = []byte{os.PathListSeparator}
	params[0].(*object.Object).FieldTable["pathSeparator"] = fld

	fld.Ftype = types.Int // status: "checked" (1) as opposed to "invalid" (0)
	fld.Fvalue = int64(1)
	params[0].(*object.Object).FieldTable["status"] = fld

	return nil
}

// "java/io/File.getPath()Ljava/lang/String;"
func fileGetPath(params []interface{}) interface{} {
	fld := params[0].(*object.Object).FieldTable["filePath"]
	bytes := fld.Fvalue.([]byte)
	return object.StringObjectFromByteArray(bytes)
}

// "java/io/File.isInvalid()Ljava/lang/String;"
func fileIsInvalid(params []interface{}) interface{} {
	status := params[0].(*object.Object).FieldTable["status"].Fvalue.(int64)
	result := status == 0
	return result
}

// -------------------- Traps ----------------------------------

// "java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;"
func trapGetDefaultFileSystem([]interface{}) interface{} {
	errMsg := "DefaultFileSystem.getFileSystem() is not yet supported"
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}
