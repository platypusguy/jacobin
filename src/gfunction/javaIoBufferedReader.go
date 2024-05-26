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
			ParamSlots: 1,
			GFunction:  bufferedReaderInitSz,
		}

	MethodSignatures["java/io/BufferedReader.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrClose,
		}

	MethodSignatures["java/io/BufferedReader.lines()Ljava/util/stream/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufferedReaderLines,
		}

	MethodSignatures["java/io/BufferedReader.mark(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bufferedReaderMarkAndReset,
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
			GFunction:  bufferedReaderMarkAndReset,
		}

}

// "java/io/BufferedReader.<init>(Ljava/io/Reader;])V"
func bufferedReaderInit(params []interface{}) interface{} {
	fld1, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "bufferedReaderInit: File argument lacks a FilePath field"
		return getGErrBlk(excNames.InvalidTypeException, errMsg)
	}
	inPathStr := string(fld1.Fvalue.([]byte))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("bufferedReaderInit: os.Open(%s) returned: %s", inPathStr, err.Error())
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

func bufferedReaderInitSz(params []interface{}) interface{} {
	errMsg := "Instantiating BufferedReader with a size is not yet supported by jacobin"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func bufferedReaderLines(params []interface{}) interface{} {
	errMsg := "Instantiating BufferedReader with a size is not yet supported by jacobin"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func bufferedReaderMarkAndReset(params []interface{}) interface{} {
	errMsg := "BufferedReader mark() & reset() are not yet supported by jacobin"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func bufferedReaderMarkSupported(params []interface{}) interface{} {
	return int64(0) // false
}

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
		errMsg := "bufferedReaderReadLine: BufferedReader object lacks a FileHandle field"
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
			errMsg := fmt.Sprintf("bufferedReaderReadLine: osFile.Read failed, reason: %s", err.Error())
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
