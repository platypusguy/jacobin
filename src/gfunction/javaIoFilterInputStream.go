/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"os"
)

func Load_Io_FilterInputStream() {

	MethodSignatures["java/io/FilterInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/FilterInputStream.<init>(Ljava.io.InputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFilterInputStreamFile,
		}

	MethodSignatures["java/io/FilterInputStream.<init>(Ljava.lang.String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFilterInputStreamString,
		}

	MethodSignatures["java/io/FilterInputStream.available()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisAvailable,
		}

	MethodSignatures["java/io/FilterInputStream.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	MethodSignatures["java/io/FilterInputStream.mark(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/FilterInputStream.markSupported()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderMarkSupported,
		}

	MethodSignatures["java/io/FilterInputStream.read()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	MethodSignatures["java/io/FilterInputStream.read([B)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	MethodSignatures["java/io/FilterInputStream.read([BII)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	MethodSignatures["java/io/FilterInputStream.reset()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/FilterInputStream.skip(J)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fisSkip,
		}

}

// "java/io/FilterInputStream.<init>(Ljava/io/File;])V"
func initFilterInputStreamFile(params []interface{}) interface{} {

	// Get file path field from the File argument.
	fld, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initFilterInputStreamFile: File object argument lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the file path.
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Open the file for read-only, yielding a file handle.
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FilterInputStream object.
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FilterInputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FilterInputStream.<init>(Ljava/lang/String;])V"
func initFilterInputStreamString(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFilterInputStreamString: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FilterInputStream object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FilterInputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}
