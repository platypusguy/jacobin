/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"errors"
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

	ghelpers.MethodSignatures["java/io/BufferedReader.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

}

// "java/io/BufferedReader.<init>(Ljava/io/Reader;])V"
func bufferedReaderInit(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil || object.IsNull(params[0]) || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bufferedReaderInit: null parameter")
	}
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, "bufferedReaderInit: self is not an object")
	}
	readerObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, "bufferedReaderInit: reader parameter is not an object")
	}

	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}

	if fldHandle, ok := readerObj.FieldTable[ghelpers.FileHandle]; ok && fldHandle.Fvalue != nil && !object.IsNull(fldHandle.Fvalue) {
		self.FieldTable[ghelpers.FileHandle] = fldHandle
		if fldPath, ok := readerObj.FieldTable[ghelpers.FilePath]; ok {
			self.FieldTable[ghelpers.FilePath] = fldPath
		}
		return nil
	}

	fld1, ok := readerObj.FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "Reader object lacks a ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, errMsg)
	}
	var inPathStr string
	switch v := fld1.Fvalue.(type) {
	case []types.JavaByte:
		inPathStr = object.GoStringFromJavaByteArray(v)
	case []byte:
		inPathStr = string(v)
	case string:
		inPathStr = v
	default:
		errMsg := "Reader object has invalid ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, errMsg)
	}

	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("os.Open(%s) failed, reason: %s", inPathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	self.FieldTable[ghelpers.FilePath] = fld1

	// Field ghelpers.FileHandle = Golang *os.File from os.Open
	self.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: types.Ref, Fvalue: osFile}

	return nil
}

// "java/io/BufferedReader.markSupported()Z"
func bufferedReaderMarkSupported([]any) any {
	return types.JavaBoolFalse // false
}

// "java/io/BufferedReader.readLine()Ljava/lang/String;"
func bufferedReaderReadLine(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bufferedReaderReadLine: self is null")
	}
	// Get BufferedReader object.
	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.InvalidTypeException, "bufferedReaderReadLine: self is not an object")
	}

	// Already at EOF?
	if ghelpers.EofGet(obj) {
		return object.Null
	}

	// Get file handle.
	fldHandle, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok || fldHandle.Fvalue == nil || object.IsNull(fldHandle.Fvalue) {
		errMsg := "Reader object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	reader, ok := fldHandle.Fvalue.(io.Reader)
	if !ok {
		errMsg := "Reader object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Need a one-byte buffer.
	byteBuf := make([]byte, 1)
	var buffer []byte
	var err error
	for {
		_, err = reader.Read(byteBuf)
		if errors.Is(err, io.EOF) {
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
