/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

func Load_Io_FilterInputStream() {

	ghelpers.MethodSignatures["java/io/FilterInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.<init>(Ljava.io.InputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFilterInputStreamFile,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.<init>(Ljava.lang.String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFilterInputStreamString,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisAvailable,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderMarkSupported,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fisSkip,
		}

}

// "java/io/FilterInputStream.<init>(Ljava/io/File;])V"
func initFilterInputStreamFile(params []interface{}) interface{} {

	// Get file path field from the File argument.
	fld, ok := params[1].(*object.Object).FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "initFilterInputStreamFile: File object argument lacks a ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Get the file path.
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Open the file for read-only, yielding a file handle.
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FilterInputStream object.
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the FilterInputStream object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}

// "java/io/FilterInputStream.<init>(Ljava/lang/String;])V"
func initFilterInputStreamString(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFilterInputStreamString: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FilterInputStream object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the FilterInputStream object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}
