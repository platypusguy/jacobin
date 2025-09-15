/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"io"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

func Load_Io_InputStreamReader() {

	MethodSignatures["java/io/InputStreamReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/InputStreamReader.<init>(Ljava/io/InputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamReaderInit,
		}

	MethodSignatures["java/io/InputStreamReader.<init>(Ljava/io/InputStream;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/InputStreamReader.<init>(Ljava/io/InputStream;Ljava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/InputStreamReader.<init>(Ljava/io/InputStream;Ljava/nio/charset/CharsetDecoder;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/InputStreamReader.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrClose,
		}

	MethodSignatures["java/io/InputStreamReader.getEncoding()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnCharsetName,
		}

	MethodSignatures["java/io/InputStreamReader.read()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrReadOneChar,
		}

	MethodSignatures["java/io/InputStreamReader.read([CII)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  isrReadCharBufferSubset,
		}

	MethodSignatures["java/io/InputStreamReader.ready()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  isrReady,
		}

}

// "java/io/InputStreamReader.<init>(Ljava/io/InputStream;)V"
func inputStreamReaderInit(params []interface{}) interface{} {

	// Get file path field.
	fldPath, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "inputStreamReaderInit: InputStream object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get file handle field.
	fldHandle, ok := params[1].(*object.Object).FieldTable[FileHandle]
	if !ok {
		errMsg := "inputStreamReaderInit: InputStream object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	osFile := fldHandle.Fvalue.(*os.File)

	// Get file statistics.
	_, err := osFile.Stat()
	if err != nil {
		pathStr := string(fldPath.Fvalue.([]byte))
		errMsg := fmt.Sprintf("inputStreamReaderInit: os.Stat(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy file path into the InputStreamReader object.
	params[0].(*object.Object).FieldTable[FilePath] = fldPath

	// Copy file handle into the InputStreamReader object.
	params[0].(*object.Object).FieldTable[FileHandle] = fldHandle

	return nil
}

// "java/io/InputStreamReader.close()V"
func isrClose(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "isrClose: InputStreamReader object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Close the file.
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("isrClose: osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// Almost a duplicate of fisReadOne in fileInputStream.go
// "java/io/InputStreamReader.read()I"
func isrReadOneChar(params []interface{}) interface{} {

	// Get InputStream object.
	obj := params[0].(*object.Object)

	// Get file handle.
	osFile, ok := obj.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "isrReadOneChar: InputStreamReader object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Need a one-byte buffer.
	buffer := make([]byte, 1)

	// Read one byte.
	_, err := osFile.Read(buffer)
	if err == io.EOF {
		eofSet(obj, true)
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("isrReadOneChar: osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Return the byte as an integer.
	return int64(buffer[0])
}

// "java/io/InputStreamReader.read([CII)I"
func isrReadCharBufferSubset(params []interface{}) interface{} {

	// Get InputStream object.
	obj := params[0].(*object.Object)

	// Get file handle.
	osFile, ok := obj.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "isrReadCharBufferSubset: InputStreamReader object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the parameter buffer, offset, and length.
	intArray, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]int64)
	if !ok {
		errMsg := "isrReadCharBufferSubset: InputStreamReader trouble with character array buffer"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	offset := params[2].(int64)
	length := params[3].(int64)

	// Check parameters.
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(intArray))-offset) {
		errMsg := fmt.Sprintf("isrReadCharBufferSubset: Error in parameters: offset=%d, length=%d, char.array.length=%d",
			offset, length, len(intArray))
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Fill the replacement byte buffer.
	inBytes := make([]byte, length)
	nbytes, err := osFile.Read(inBytes)
	if err == io.EOF {
		eofSet(obj, true)
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("isrReadCharBufferSubset: osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Update the parameter buffer, beginning at the offset.
	for ii := int64(0); ii < int64(nbytes); ii++ {
		intArray[offset+ii] = int64(inBytes[ii])
	}

	// Update the parameter buffer.
	fld := object.Field{Ftype: types.IntArray, Fvalue: intArray}
	params[1].(*object.Object).FieldTable["value"] = fld

	// Return the number of bytes.
	return int64(nbytes)

}

// "java/io/InputStreamReader.ready()Z"
func isrReady(params []interface{}) interface{} {

	// Get file path.
	fldPath, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "isrReady: InputStreamReader object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get file handle.
	fldHandle, ok := params[1].(*object.Object).FieldTable[FileHandle]
	if !ok {
		errMsg := "isrReady: InputStreamReader object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy java/io/File path field into the InputStreamReader object.
	fld := fldPath
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Get file handle and get file statistics.
	osFile := fldHandle.Fvalue.(*os.File)
	_, err := osFile.Stat()
	if err != nil {
		return types.JavaBoolFalse // Ready: false
	}

	return types.JavaBoolTrue
}
