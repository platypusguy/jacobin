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
)

func Load_Io_FileReader() map[string]GMeth {

	MethodSignatures["java/io/FileReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/FileReader.<init>(Ljava/io/File;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileReader,
		}

	MethodSignatures["java/io/FileReader.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileReaderString,
		}

	// -----------------------------------------
	// Traps that do nothing but return an error
	// -----------------------------------------

	MethodSignatures["java/io/FileReader.<init>(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFileDescriptor,
		}

	MethodSignatures["java/io/FileReader.<init>(Ljava/io/File;Ljava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapCharset,
		}

	MethodSignatures["java/io/FileReader.<init>(Ljava/lang/String;Ljava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapCharset,
		}

	return MethodSignatures
}

// "java/io/FileReader.<init>(Ljava/io/File;])V"
func initFileReader(params []interface{}) interface{} {
	fld1, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initFileReader: File argument lacks a FilePath field"
		return getGErrBlk(exceptions.InvalidTypeException, errMsg)
	}
	inPathStr := string(fld1.Fvalue.([]byte))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileReader: os.Open(%s) returned: %s", inPathStr, err.Error())
		return getGErrBlk(exceptions.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := fld1
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Field FileHandle = Golang *os.File from os.Open
	fld = object.Field{Ftype: types.Ref, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileReader.<init>(Ljava/lang/String;])V"
func initFileReaderString(params []interface{}) interface{} {
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileReaderString: os.Open(%s) returned: %s", pathStr, err.Error())
		return getGErrBlk(exceptions.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Field FileHandle = Golang *os.File
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}
