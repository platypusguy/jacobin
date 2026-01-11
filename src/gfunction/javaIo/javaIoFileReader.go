/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
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

func Load_Io_FileReader() {

	ghelpers.MethodSignatures["java/io/FileReader.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/FileReader.<init>(Ljava/io/File;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFileReader,
		}

	ghelpers.MethodSignatures["java/io/FileReader.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFileReaderString,
		}

	// -----------------------------------------
	// traps that do nothing but return an error
	// -----------------------------------------

	ghelpers.MethodSignatures["java/io/FileReader.<init>(Ljava/io/FileDescriptor;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.<init>(Ljava/io/File;Ljava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.<init>(Ljava/lang/String;Ljava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.getEncoding()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.read([C)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.read([CII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.read(Ljava/nio/CharBuffer;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.ready()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileReader.transferTo(Ljava/io/Writer;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

}

// "java/io/FileReader.<init>(Ljava/io/File;)V"
func initFileReader(params []interface{}) interface{} {
	fld1, ok := params[1].(*object.Object).FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "initFileReader: File object lacks a ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, errMsg)
	}
	inPathStr := object.GoStringFromJavaByteArray(fld1.Fvalue.([]types.JavaByte))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileReader: os.Open(%s) failed, reason: %s", inPathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := fld1
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Field ghelpers.FileHandle = Golang *os.File from os.Open
	fld = object.Field{Ftype: types.Ref, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}

// "java/io/FileReader.<init>(Ljava/lang/String;)V"
func initFileReaderString(params []interface{}) interface{} {
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileReaderString: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Field ghelpers.FileHandle = Golang *os.File
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}
