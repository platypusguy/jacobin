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

func Load_Io_FileOutputStream() {

	MethodSignatures["java/io/FileOutputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/FileOutputStream.<init>(Ljava/io/File;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileOutputStreamFile,
		}

	MethodSignatures["java/io/FileOutputStream.<init>(Ljava/io/File;Z)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  initFileOutputStreamFileBoolean,
		}

	MethodSignatures["java/io/FileOutputStream.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileOutputStreamString,
		}

	MethodSignatures["java/io/FileOutputStream.<init>(Ljava/lang/String;Z)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  initFileOutputStreamStringBoolean,
		}

	MethodSignatures["java/io/FileOutputStream.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fosClose,
		}

	MethodSignatures["java/io/FileOutputStream.write(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fosWriteOne,
		}

	MethodSignatures["java/io/FileOutputStream.write([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fosWriteByteArray,
		}

	MethodSignatures["java/io/FileOutputStream.write([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  fosWriteByteArrayOffset,
		}

	MethodSignatures["java/io/FileOutputStream.<init>(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/FileOutputStream.getFD()Ljava/io/FileDescriptor;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

}

// "java/io/FileOutputStream.<init>(Ljava/io/File;])V"
func initFileOutputStreamFile(params []interface{}) interface{} {

	// Get the file path.
	fld, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initFileOutputStreamFile: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Open the file for write-only, yielding a file handle.
	osFile, err := os.Create(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileOutputStreamFile: os.Create(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FileOutputStream object.
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FileOutputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileOutputStream.<init>(Ljava/io/File;Z])V"
func initFileOutputStreamFileBoolean(params []interface{}) interface{} {

	// Get file path field from the File argument.
	fld, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initFileOutputStreamFileBoolean: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the file path.
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Get the boolean argument.
	boolarg, ok := params[2].(int64)
	if !ok {
		errMsg := "initFileOutputStreamFileBoolean: Missing append-boolean argument"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Open the file for write-only, yielding a file handle.
	var osFile *os.File
	var err error
	if boolarg != 0 { // append: true
		osFile, err = os.OpenFile(pathStr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, CreateFilePermissions)
	} else {
		osFile, err = os.Create(pathStr)
	}
	if err != nil {
		errMsg := fmt.Sprintf("initFileOutputStreamFileBoolean: os.OpenFile|os.Create(%s) failed, reason: %s",
			pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FileOutputStream object.
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FileOutputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileOutputStream.<init>(Ljava/lang/String;])V"
func initFileOutputStreamString(params []interface{}) interface{} {

	// Using the argument path string, open the file for write-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))

	// Open the file for write-only, yielding a file handle.
	osFile, err := os.Create(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileOutputStreamString: os.Create(%s) failed, reason: %s", pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FileOutputStream object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FileOutputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileOutputStream.<init>(Ljava/lang/String;])V"
func initFileOutputStreamStringBoolean(params []interface{}) interface{} {

	// Using the argument path string, open the file for write-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))

	// Get the boolean argument.
	boolarg, ok := params[2].(int64)
	if !ok {
		errMsg := "initFileOutputStreamStringBoolean: Missing append-boolean argument"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Open the file for write-only, yielding a file handle.
	var osFile *os.File
	var err error
	if boolarg != 0 { // append: true
		osFile, err = os.OpenFile(pathStr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		osFile, err = os.Create(pathStr)
	}
	if err != nil {
		errMsg := fmt.Sprintf("initFileOutputStreamStringBoolean: os.OpenFile|os.Create(%s) failed, reason: %s",
			pathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FileOutputStream object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FileOutputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileOutputStream.write(I)"
func fosWriteOne(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fosWriteOne: FileOutputStream object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Get the integer argument.
	wint, ok := params[1].(int64)
	if !ok {
		errMsg := "fosWriteOne: Missing integer argument"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Create a one-byte buffer.
	buffer := make([]byte, 1)
	buffer[0] = byte(wint % 256)

	// Write one byte.
	_, err := osFile.Write(buffer)
	if err != nil {
		errMsg := fmt.Sprintf("fosWriteOne: osFile.Write failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/FileOutputStream.write([B)I"
func fosWriteByteArray(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fosWriteByteArray: FileOutputStream object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Set buffer to the byte array parameter.
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := "fosWriteByteArray: Byte array parameter lacks a \"value\" field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	buffer := object.GoByteArrayFromJavaByteArray(javaBytes)

	// Write the buffer.
	_, err := osFile.Write(buffer)
	if err != nil {
		errMsg := fmt.Sprintf("fosWriteByteArray: osFile.Write failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/FileOutputStream.write([BII)I"
func fosWriteByteArrayOffset(params []interface{}) interface{} {

	// Get the file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fosWriteByteArrayOffset: FileOutputStream object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Set buffer (buf1) to the byte array parameter.
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := "fosWriteByteArrayOffset: Byte array parameter lacks a \"value\" field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	buf1 := object.GoByteArrayFromJavaByteArray(javaBytes)

	// Collect the offset and length parameter values.
	offset := params[2].(int64)
	length := params[3].(int64)

	// Check the parameters.
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(buf1))-offset) {
		errMsg := fmt.Sprintf("fosWriteByteArrayOffset: Error in parameters offset=%d length=%d bytes.length=%d",
			offset, length, len(buf1))
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Write the byte buffer.
	_, err := osFile.Write(buf1[offset : offset+length])
	if err != nil {
		errMsg := fmt.Sprintf("fosWriteByteArrayOffset: osFile.Write failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/FileOutputStream.close()V"
func fosClose(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fosClose: FileOutputStream object lacks a FileHandle field"
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Close the file.
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("fosClose: osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}
