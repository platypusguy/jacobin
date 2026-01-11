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

func Load_Io_RandomAccessFile() {

	ghelpers.MethodSignatures["java/io/RandomAccessFile.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafInitString,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.<init>(Ljava/io/File;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafInitFile,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.getFilePointer()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafGetFilePointer,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	// ----------------------------------------------------------
	// initIDs - ghelpers.JustReturn
	// This is a private function that call C native functions.
	// ----------------------------------------------------------

	ghelpers.MethodSignatures["java/io/RandomAccessFile.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

}

// "java/io/RandomAccessFile.<init>(Ljava/lang/String;Ljava/lang/String;)V"
// RandomAccessFile raf = new RandomAccessFile(Stringname, Stringmode);
func rafInitString(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))

	// Mode.
	var modeInt int
	modeStr := object.GoStringFromStringObject(params[2].(*object.Object))
	switch modeStr {
	case "r":
		modeInt = os.O_RDONLY
	case "rw", "rws", "rwd":
		modeInt = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		errMsg := fmt.Sprintf("rafInitString: mode string (%s) invalid", modeStr)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, ghelpers.CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafInitString: os.OpenFile(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil

}

// "java/io/RandomAccessFile.<init>(Ljava/io/File;Ljava/lang/String;)V"
// RandomAccessFile raf = new RandomAccessFile(Fileobject, Stringmode);
func rafInitFile(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	obj := params[1].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "rafInitFile: java/io/File object is missing the ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Mode.
	var modeInt int
	modeStr := object.GoStringFromStringObject(params[2].(*object.Object))
	switch modeStr {
	case "r":
		modeInt = os.O_RDONLY
	case "rw", "rws", "rwd":
		modeInt = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		errMsg := fmt.Sprintf("rafInitFile: mode string (%s) invalid", modeStr)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, ghelpers.CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafInitFile: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil

}

// "java/io/RandomAccessFile.getFilePointer()J"
// Get current file position (offset from the beginning of file).
func rafGetFilePointer(params []interface{}) interface{} {

	// Get the open file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafGetFilePointer: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the current position relative to the beginning of file.
	posn, err := osFile.Seek(0, 1)
	if err != nil {
		errMsg := fmt.Sprintf("rafGetFilePointer: osFile.Seek(0, 1) failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return posn

}
