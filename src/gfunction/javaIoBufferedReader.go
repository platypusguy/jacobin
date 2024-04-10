/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"io"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"os"
)

func Load_Io_BufferedReader() map[string]GMeth {

	MethodSignatures["java/io/BufferedReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initBufferedReader,
		}

	MethodSignatures["java/io/BufferedReader.<init>(Ljava/io/Reader;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  initBufferedReader,
		}

	MethodSignatures["java/io/BufferedReader.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufrdrClose,
		}

	MethodSignatures["java/io/BufferedReader.mark(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bufrdrMark,
		}

	MethodSignatures["java/io/BufferedReader.markSupported()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufrdrMarkSupported,
		}

	MethodSignatures["java/io/BufferedReader.read()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufrdrReadOneChar,
		}

	MethodSignatures["java/io/BufferedReader.read([CII)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  ReadCharBuffer,
		}

	MethodSignatures["java/io/BufferedReader.readLine()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufrdrReadLine,
		}

	MethodSignatures["java/io/BufferedReader.ready()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufrdrReady,
		}

	MethodSignatures["java/io/BufferedReader.reset()"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bufrdrReset,
		}

	// -----------------------------------------
	// Traps that do nothing but return an error
	// -----------------------------------------

	MethodSignatures["java/io/BufferedReader.lines()Ljava.util.stream.Stream<String>;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapBufrdrStream,
		}

	return MethodSignatures
}

// "java/io/BufferedReader.<init>(Ljava/io/Reader;)V"
func initBufferedReader(params []interface{}) interface{} {

	// Get file path.
	fldPath, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initBufferedReader: Reader argument lacks a FilePath field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Copy java/io/File path field into the BufferedReader object.
	fld := fldPath
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Get O/S file handle.
	fldHandle, ok := params[1].(*object.Object).FieldTable[FileHandle]
	if !ok {
		errMsg := "initBufferedReader: Reader argument lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	osFile := fldHandle.Fvalue.(*os.File)
	_, err := osFile.Stat()
	if err != nil {
		pathStr := string(fldPath.Fvalue.([]byte))
		errMsg := fmt.Sprintf("initBufferedReader: os.Stat(%s) returned: %s", pathStr, err.Error())
		return getGErrBlk(exceptions.FileNotFoundException, errMsg)
	}

	// Copy Reader FilePath field to the BufferedReader object.
	params[0].(*object.Object).FieldTable[FilePath] = fldPath

	// Field FileHandle = Golang *os.File from os.Open.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

func bufrdrClose(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bufrdrClose: BufferedReader object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("bufrdrClose osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	return nil
}

// "java/io/BufferedReader.mark(I)V"
func bufrdrMark(params []interface{}) interface{} {

	// Get handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bufrdrMark: BufferedReader object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Get current file offset.
	posn, err := osFile.Seek(0, io.SeekCurrent)
	if err != nil {
		errMsg := fmt.Sprintf("bufrdrMark osFile.Seek() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Record file offset in field FileMark.
	params[0].(*object.Object).FieldTable[FileMark] = object.Field{Ftype: types.Int, Fvalue: posn}
	return nil

}

// "java/io/BufferedReader.markSupported()Z"
func bufrdrMarkSupported([]interface{}) interface{} {
	return int64(1) // true
}

// Almost a duplicate of fisReadOne in fileInputStream.go
// "java/io/BufferedReader.read()I"
func bufrdrReadOneChar(params []interface{}) interface{} {

	// Get Reader object.
	obj := params[0].(*object.Object)
	if eofGet(obj) {
		return int64(-1) // return -1 on EOF
	}

	// Get file handle.
	osFile, ok := obj.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bufrdrReadOneChar: BufferedReader object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Need a one-byte buffer.
	buffer := make([]byte, 1)

	// Try read.
	_, err := osFile.Read(buffer)
	if err == io.EOF {
		eofSet(obj, true)
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("bufrdrReadOneChar: osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Return the byte as an integer.
	return int64(buffer[0])
}

// "java/io/BufferedReader.readLine()Ljava/lang/String;"
func bufrdrReadLine(params []interface{}) interface{} {

	// Get Reader object.
	objRdr := params[0].(*object.Object)
	if eofGet(objRdr) {
		return int64(-1) // return -1 on EOF
	}

	// Get file handle.
	osFile, ok := objRdr.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bufrdrReadLine: BufferedReader object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Read one byte at a time, assembling a byte array.
	// Stop at EOF or \n.
	// Discard \r.
	var bytes []byte
	buffer := make([]byte, 1)
	for {
		// Read one byte.
		_, err := osFile.Read(buffer)
		if err == io.EOF {
			eofSet(objRdr, true)
			objStr := object.Null // jvm/run.go IFNULL will catch this.
			return objStr
		}
		if err != nil {
			errMsg := fmt.Sprintf("bufrdrReadLine: osFile.Read failed, reason: %s", err.Error())
			return getGErrBlk(exceptions.IOException, errMsg)
		}

		// If \r, discard it.
		if buffer[0] == '\r' {
			continue
		}

		// If \n, that is the end of line signal.
		if buffer[0] == '\n' {
			break
		}

		// It's a text byte. Append the bytes slice.
		bytes = append(bytes, buffer[0])
	}

	// Make a String object from the assembled bytes and return it.
	objStr := object.StringObjectFromByteArray(bytes)
	return objStr
}

// "java/io/BufferedReader.ready()Z"
func bufrdrReady(params []interface{}) interface{} {

	obj := params[0].(*object.Object)

	// Get file path.
	_, ok := obj.FieldTable[FilePath]
	if !ok {
		errMsg := "bufrdrReady: BufferedReader object lacks a FilePath field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Copy java/io/File path field into the BufferedReader object.
	if eofGet(obj) {
		return int64(0) // return -1 on EOF
	}

	// This file is not at EOF.
	return int64(1)
}

// "java/io/BufferedReader.reset()V"
func bufrdrReset(params []interface{}) interface{} {

	// Get Reader object.
	obj := params[0].(*object.Object)

	// Get file handle.
	osFile, ok := obj.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bufrdrReset: BufferedReader object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Get marked position.
	posn, ok := params[0].(*object.Object).FieldTable[FileMark].Fvalue.(int64)
	if !ok {
		errMsg := "bufrdrReset: BufferedReader object lacks a FileMark field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Set current file offset.
	_, err := osFile.Seek(posn, io.SeekStart)
	if err != nil {
		errMsg := fmt.Sprintf("bufrdrReset osFile.Seek() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Turn off at EOF.
	eofSet(obj, false)

	return nil
}

// Trap for Stream<String> reference
func trapBufrdrStream([]interface{}) interface{} {
	errMsg := "Method java/io/BufferedReader.lines().Stream<String> is not yet supported"
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}
