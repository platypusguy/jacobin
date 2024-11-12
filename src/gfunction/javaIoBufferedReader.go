/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"io"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"os"
)

func Load_Io_BufferedReader() {

	MethodSignatures["java/io/BufferedReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bufferedReaderInit,
		}

	MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/BufferedReader.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrClose,
		}

	MethodSignatures["java/io/BufferedReader.lines()Ljava/util/stream/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/BufferedReader.mark(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/BufferedReader.markSupported()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderMarkSupported,
		}

	MethodSignatures["java/io/BufferedReader.read()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrReadOneChar,
		}

	MethodSignatures["java/io/BufferedReader.read([CII)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  isrReadCharBufferSubset,
		}

	MethodSignatures["java/io/BufferedReader.readLine()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderReadLine,
		}

	MethodSignatures["java/io/BufferedReader.ready()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrReady,
		}

	MethodSignatures["java/io/BufferedReader.reset()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

}

// "java/io/BufferedReader.<init>(Ljava/io/Reader;])V"
func bufferedReaderInit(params []interface{}) interface{} {
	fld1, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "Reader object lacks a FilePath field"
		return getGErrBlk(excNames.InvalidTypeException, errMsg)
	}
	inPathStr := string(fld1.Fvalue.([]byte))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("os.Open(%s) failed, reason: %s", inPathStr, err.Error())
		return getGErrBlk(excNames.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := fld1
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Field FileHandle = Golang *os.File from os.Open
	fld = object.Field{Ftype: types.Ref, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/BufferedReader.markSupported()Z"
func bufferedReaderMarkSupported(params []interface{}) interface{} {
	return int64(0) // false
}

// "java/io/BufferedReader.readLine()Ljava/lang/String;"
func bufferedReaderReadLine(params []interface{}) interface{} {
	// Get BufferedReader object.
	obj := params[0].(*object.Object)

	// Already at EOF?
	if eofGet(obj) {
		return object.Null
	}

	// Get file handle.
	osFile, ok := obj.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "Reader object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Need a one-byte buffer.
	byteBuf := make([]byte, 1)
	var buffer []byte
	var err error
	for {
		_, err = osFile.Read(byteBuf)
		if err == io.EOF {
			eofSet(obj, true)
			if len(buffer) > 0 {
				break
			}
			return object.Null
		}
		if err != nil {
			errMsg := fmt.Sprintf("osFile.Read failed, reason: %s", err.Error())
			return getGErrBlk(excNames.IOException, errMsg)
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
