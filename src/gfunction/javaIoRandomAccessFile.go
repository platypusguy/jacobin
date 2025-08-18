/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

func Load_Io_RandomAccessFile() {

	MethodSignatures["java/io/RandomAccessFile.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/RandomAccessFile.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  rafInitString,
		}

	MethodSignatures["java/io/RandomAccessFile.<init>(Ljava/io/File;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  rafInitFile,
		}

	MethodSignatures["java/io/RandomAccessFile.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	MethodSignatures["java/io/RandomAccessFile.getFilePointer()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  rafGetFilePointer,
		}

	MethodSignatures["java/io/RandomAccessFile.read()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	MethodSignatures["java/io/RandomAccessFile.read([B)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	MethodSignatures["java/io/RandomAccessFile.read([BII)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	// ----------------------------------------------------------
	// initIDs - justReturn
	// This is a private function that call C native functions.
	// ----------------------------------------------------------

	MethodSignatures["java/io/RandomAccessFile.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafInitString: os.OpenFile(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil

}

// "java/io/RandomAccessFile.<init>(Ljava/io/File;Ljava/lang/String;)V"
// RandomAccessFile raf = new RandomAccessFile(Fileobject, Stringmode);
func rafInitFile(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	obj := params[1].(*object.Object)
	fld, ok := obj.FieldTable[FilePath]
	if !ok {
		errMsg := "rafInitFile: java/io/File object is missing the FilePath field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafInitFile: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil

}

// "java/io/RandomAccessFile.getFilePointer()J"
// Get current file position (offset from the beginning of file).
func rafGetFilePointer(params []interface{}) interface{} {

	// Get the open file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[FileHandle]
	if !ok {
		errMsg := "rafGetFilePointer: java/io/RandomAccessFile object is missing the FileHandle field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the current position relative to the beginning of file.
	posn, err := osFile.Seek(0, 1)
	if err != nil {
		errMsg := fmt.Sprintf("rafGetFilePointer: osFile.Seek(0, 1) failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return posn

}
