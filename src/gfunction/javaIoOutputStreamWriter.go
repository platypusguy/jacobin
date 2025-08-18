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

func Load_Io_OutputStreamWriter() {

	MethodSignatures["java/io/OutputStreamWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/OutputStreamWriter.<init>(Ljava/io/OutputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initOutputStreamWriter,
		}

	MethodSignatures["java/io/OutputStreamWriter.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  oswClose,
		}

	MethodSignatures["java/io/OutputStreamWriter.flush()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  oswFlush,
		}

	MethodSignatures["java/io/OutputStreamWriter.write(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  oswWriteOneChar,
		}

	MethodSignatures["java/io/OutputStreamWriter.write([CII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  oswWriteCharBuffer,
		}

	MethodSignatures["java/io/OutputStreamWriter.write(Ljava/lang/String;II)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  oswWriteStringBuffer,
		}

	// -----------------------------------------
	// Traps that do nothing but return an error
	// -----------------------------------------

	MethodSignatures["java/io/OutputStreamWriter.<init>(Ljava/io/OutputStream;Ljava/lang.String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/OutputStreamWriter.<init>(Ljava/io/OutputStream;Ljava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/OutputStreamWriter.<init>(Ljava/io/OutputStream;Ljava/nio/charset/CharsetDecoder;)Ljava/lang.String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/OutputStreamWriter.getEncoding()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

}

// "java/io/OutputStreamWriter.<init>(Ljava/io/OutputStream;)V"
func initOutputStreamWriter(params []interface{}) interface{} {

	// Get file path field.
	fldPath, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initOutputStreamWriter: OutputStream object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get file handle field.
	fldHandle, ok := params[1].(*object.Object).FieldTable[FileHandle]
	if !ok {
		errMsg := "initOutputStreamWriter: OutputStream object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	osFile := fldHandle.Fvalue.(*os.File)

	// Get file statistics.
	_, err := osFile.Stat()
	if err != nil {
		pathStr := string(fldPath.Fvalue.([]byte))
		errMsg := fmt.Sprintf("initOutputStreamWriter: os.Stat(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy file path into the OutputStreamWriter object.
	params[0].(*object.Object).FieldTable[FilePath] = fldPath

	// Copy file handle into the OutputStreamWriter object.
	params[0].(*object.Object).FieldTable[FileHandle] = fldHandle

	return nil
}

func oswClose(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "oswClose: OutputStreamWriter object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Close the file.
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("oswClose: osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

func oswFlush(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "oswFlush: OutputStreamWriter object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Flush the file's buffers.
	err := osFile.Sync()
	if err != nil {
		errMsg := fmt.Sprintf("oswFlush: osFile.Sync() failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/io/OutputStreamWriter.write(I)"
func oswWriteOneChar(params []interface{}) interface{} {

	// Get OutputStream object.
	obj := params[0].(*object.Object)

	// Get file handle.
	osFile, ok := obj.FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "oswWriteOneChar: OutputStreamWriter object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the integer argument.
	wint, ok := params[1].(int64)
	if !ok {
		errMsg := "oswWriteOneChar: Error in integer argument"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Create a one-byte buffer.
	buffer := make([]byte, 1)
	buffer[0] = byte(wint % 256)

	// Write one byte.
	_, err := osFile.Write(buffer)
	if err != nil {
		errMsg := fmt.Sprintf("oswWriteOneChar: osFile.Write failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/OutputStreamWriter.write([CII)I"
func oswWriteCharBuffer(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "oswWriteCharBuffer: OutputStreamWriter object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the parameter buffer, offset, and length.
	intArray, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]int64)
	if !ok {
		errMsg := "oswWriteCharBuffer: Trouble with value field ([]int64)"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	offset := params[2].(int64)
	length := params[3].(int64)

	// Check parameters.
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(intArray))-offset) {
		errMsg := fmt.Sprintf("oswWriteCharBuffer: Error in parameters: offset=%d, length=%d, char.array.length=%d",
			offset, length, len(intArray))
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Create and fill a byte buffer.
	outBytes := make([]byte, length)
	for ii := int64(0); ii < length; ii++ {
		outBytes[ii] = byte(intArray[offset+ii])
	}

	// Write the byte buffer.
	_, err := osFile.Write(outBytes)
	if err != nil {
		errMsg := fmt.Sprintf("oswWriteCharBuffer: osFile.Write failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/OutputStreamWriter.write(Ljava/lang/String;II)I"
func oswWriteStringBuffer(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "oswWriteStringBuffer: OutputStreamWriter object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the parameter string byte array, offset, and length.
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := "oswWriteStringBuffer: Trouble with value field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	paramBytes := object.GoByteArrayFromJavaByteArray(javaBytes)
	offset := params[2].(int64)
	length := params[3].(int64)

	// Check parameters.
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(paramBytes))-offset) {
		errMsg := fmt.Sprintf("oswWriteStringBuffer: Error in parameters: offset=%d, length=%d, char.array.length=%d",
			offset, length, len(paramBytes))
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Create and fill a byte buffer.
	outBytes := make([]byte, length)
	for ii := int64(0); ii < length; ii++ {
		outBytes[ii] = paramBytes[offset+ii]
	}

	// Write the byte buffer.
	_, err := osFile.Write(outBytes)
	if err != nil {
		errMsg := fmt.Sprintf("oswWriteStringBuffer: osFile.Write failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}
