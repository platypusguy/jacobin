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

func Load_Io_FileInputStream() {

	ghelpers.MethodSignatures["java/io/FileInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.<init>(Ljava/io/File;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFileInputStreamFile,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFileInputStreamString,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisAvailable,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.readNBytes([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fisSkip,
		}

	// ----------------------------------------------------------
	// These are private functiona that calls C native functions.
	// ----------------------------------------------------------

	ghelpers.MethodSignatures["java/io/FileInputStream.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/UnixFileSystem.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/WinNTFileSystem.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	// -----------------------------------------
	// traps that do nothing but return an error
	// -----------------------------------------

	ghelpers.MethodSignatures["java/io/FileInputStream.<init>(Ljava/io/FileDescriptor;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.getChannel()Ljava/nio/channels/FileChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileInputStream.getFD()Ljava/io/FileDescriptor;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

}

// "java/io/FileInputStream.<init>(Ljava/io/File;])V"
func initFileInputStreamFile(params []interface{}) interface{} {

	// Get file path field from the File argument.
	fld, ok := params[1].(*object.Object).FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "initFileInputStreamFile: File object argument lacks a ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Get the file path.
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Open the file for read-only, yielding a file handle.
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileInputStreamFile: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FileInputStream object.
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the FileInputStream object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}

// "java/io/FileInputStream.<init>(Ljava/lang/String;])V"
func initFileInputStreamString(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileInputStreamString: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the FileInputStream object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the FileInputStream object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}

// "java/io/FileInputStream.available()I"
func fisAvailable(params []interface{}) interface{} {

	// Get the file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisAvailable: FileInputStream object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Compute total file size.
	fileInfo, err := osFile.Stat()
	if err != nil {
		path := object.GoStringFromJavaByteArray(params[0].(*object.Object).FieldTable["path"].Fvalue.([]types.JavaByte))
		errMsg := fmt.Sprintf("fisAvailable: osFile.Stat(%s) failed, reason: %s", path, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	fsize := fileInfo.Size()

	// Get current file offset.
	posn, err := osFile.Seek(0, io.SeekCurrent)
	if err != nil {
		errMsg := fmt.Sprintf("fisAvailable: osFile.Seek() failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Compute and return the number of bytes remaining.
	return fsize - posn
}

// "java/io/FileInputStream.read()I"
func fisReadOne(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisReadOne: FileInputStream object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Create a one-byte buffer.
	buffer := make([]byte, 1)

	// Read one byte.
	_, err := osFile.Read(buffer)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadOne: osFile.Read failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Return the read byte as an integer.
	return int64(buffer[0])
}

// "java/io/FileInputStream.read([B)I"
func fisReadByteArray(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisReadByteArray: FileInputStream object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Set buffer to the byte array parameter.
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := "fisReadByteArray: Byte array parameter lacks a \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	buffer := object.GoByteArrayFromJavaByteArray(javaBytes)

	// Fill the buffer.
	nbytes, err := osFile.Read(buffer)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadByteArray: osFile.Read failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// All is well - update the supplied buffer.
	javaBytes = object.JavaByteArrayFromGoByteArray(buffer[:nbytes])
	fld := object.Field{Ftype: types.ByteArray, Fvalue: javaBytes}
	params[1].(*object.Object).FieldTable["value"] = fld

	// Return the number of bytes.
	return int64(nbytes)
}

// "java/io/FileInputStream.read([BII)I"
func fisReadByteArrayOffset(params []interface{}) interface{} {

	// Get the file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisReadByteArrayOffset: FileInputStream object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Set buffer (buf1) to the byte array parameter.
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := "fisReadByteArrayOffset: Byte array parameter lacks a \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
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
		errMsg := fmt.Sprintf("fisReadByteArrayOffset: Error in parameters offset=%d length=%d bytes.length=%d",
			offset, length, len(buf1))
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Try read with a second buffer.
	buf2 := make([]byte, length)
	nbytes, err := osFile.Read(buf2)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadByteArrayOffset: osFile.Read failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// All is well - Copy the bytes read into the original buffer, beginning at the offset.
	copy(buf1[offset:], buf2)

	// Update the parameter buffer.
	javaBytes = object.JavaByteArrayFromGoByteArray(buf1)
	fld := object.Field{Ftype: types.ByteArray, Fvalue: javaBytes}
	params[1].(*object.Object).FieldTable["value"] = fld

	// Return the number of bytes.
	return int64(nbytes)
}

// "java/io/FileInputStream.skip(J)J"
func fisSkip(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisSkip: FileInputStream object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Get skip count.
	count := params[1].(int64)

	// Skip.
	_, err := osFile.Seek(count, 1)
	if err != nil {
		errMsg := fmt.Sprintf("fisSkip: osFile.Seek(%d) failed, reason: %s", count, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Return skip count.
	return count
}

// "java/io/FileInputStream.close()V"
func fisClose(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisClose: FileInputStream object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Close the file.
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("fisClose: osFile.Close() failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}
