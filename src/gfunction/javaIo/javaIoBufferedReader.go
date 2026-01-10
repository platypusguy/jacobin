/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"fmt"
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

func Load_Io_BufferedReader() {

	ghelpers.MethodSignatures["java/io/BufferedReader.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bufferedReaderInit,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  isrClose,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.lines()Ljava/util/stream/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderMarkSupported,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  isrReadOneChar,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.read([CII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  isrReadCharBufferSubset,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.readLine()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderReadLine,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.ready()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  isrReady,
		}

	ghelpers.MethodSignatures["java/io/BufferedReader.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

}

// "java/io/BufferedReader.<init>(Ljava/io/Reader;])V"
func bufferedReaderInit(params []interface{}) interface{} {
	fld1, ok := params[1].(*object.Object).FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "Reader object lacks a ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, errMsg)
	}
	inPathStr := object.GoStringFromJavaByteArray(fld1.Fvalue.([]types.JavaByte))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("os.Open(%s) failed, reason: %s", inPathStr, err.Error())
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

// "java/io/BufferedReader.markSupported()Z"
func bufferedReaderMarkSupported([]interface{}) interface{} {
	return types.JavaBoolFalse // false
}

// "java/io/BufferedReader.readLine()Ljava/lang/String;"
func bufferedReaderReadLine(params []interface{}) interface{} {
	// Get BufferedReader object.
	obj := params[0].(*object.Object)

	// Already at EOF?
	if ghelpers.EofGet(obj) {
		return object.Null
	}

	// Get file handle.
	osFile, ok := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "Reader object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Need a one-byte buffer.
	byteBuf := make([]byte, 1)
	var buffer []byte
	var err error
	for {
		_, err = osFile.Read(byteBuf)
		if err == io.EOF {
			ghelpers.EofSet(obj, true)
			if len(buffer) > 0 {
				break
			}
			return object.Null
		}
		if err != nil {
			errMsg := fmt.Sprintf("osFile.Read failed, reason: %s", err.Error())
			return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
		}
		if byteBuf[0] == '\r' {
			continue
		}
		if byteBuf[0] == '\n' {
			break
		}
		buffer = append(buffer, byteBuf[0])
	}

	// Return the string.
	return object.StringObjectFromByteArray(buffer)
}
