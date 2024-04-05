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

func Load_Io_FileInputStream() map[string]GMeth {

	MethodSignatures["java/io/FileInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/FileInputStream.<init>(Ljava/io/File;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileInputStreamFile,
		}

	MethodSignatures["java/io/FileInputStream.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileInputStreamString,
		}

	MethodSignatures["java/io/FileInputStream.available()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisAvailable,
		}

	MethodSignatures["java/io/FileInputStream.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	MethodSignatures["java/io/FileInputStream.read()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	MethodSignatures["java/io/FileInputStream.read([B)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	MethodSignatures["java/io/FileInputStream.read([BII)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	MethodSignatures["java/io/FileInputStream.skip(J)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fisSkip,
		}

	// ----------------------------------------------------------
	// initIDs justReturn
	// These are private functiona that calls C native functions.
	// ----------------------------------------------------------

	MethodSignatures["java/io/FileInputStream.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/UnixFileSystem.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/WinNTFileSystem.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// -----------------------------------------
	// Traps that do nothing but return an error
	// -----------------------------------------

	MethodSignatures["java/io/FileInputStream.<init>(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFileDescriptor,
		}

	MethodSignatures["java/io/FileInputStream.getChannel()Ljava/nio/channels/FileChannel;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileChannel,
		}

	MethodSignatures["java/io/FileInputStream.getFD()Ljava/io/FileDescriptor;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileDescriptor,
		}

	return MethodSignatures
}

// "java/io/FileInputStream.<init>(Ljava/io/File;])V"
func initFileInputStreamFile(params []interface{}) interface{} {

	// Get file path field from the File argument.
	fld, ok := params[1].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "initFileInputStreamFile: File argument lacks a FilePath field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Get the file path.
	pathStr := string(fld.Fvalue.([]byte))

	// Open the file for read-only, yielding a file handle.
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileInputStreamFile: os.Open(%s) returned: %s", pathStr, err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Copy the file path field into the FileInputStream object.
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FileInputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileInputStream.<init>(Ljava/lang/String;])V"
func initFileInputStreamString(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileInputStreamString: os.Open(%s) returned: %s", pathStr, err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Copy the file path field into the FileInputStream object.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[FilePath] = fld

	// Copy the file handle into the FileInputStream object.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return nil
}

// "java/io/FileInputStream.available()I"
func fisAvailable(params []interface{}) interface{} {

	// Get the file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisAvailable: FileInputStream object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Compute total file size.
	fileInfo, err := osFile.Stat()
	if err != nil {
		path := string(params[0].(*object.Object).FieldTable["path"].Fvalue.([]byte))
		errMsg := fmt.Sprintf("fisAvailable: osFile.Stat(%s) returned: %s", path, err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	fsize := fileInfo.Size()

	// Get current file offset.
	posn, err := osFile.Seek(0, io.SeekCurrent)
	if err != nil {
		errMsg := fmt.Sprintf("fisAvailable: osFile.Seek() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Compute and return the number of bytes remaining.
	return fsize - posn
}

// "java/io/FileInputStream.read()I"
func fisReadOne(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisReadOne: FileInputStream object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Create a one-byte buffer.
	buffer := make([]byte, 1)

	// Read one byte.
	_, err := osFile.Read(buffer)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadOne osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Return the read byte as an integer.
	return int64(buffer[0])
}

// "java/io/FileInputStream.read([B)I"
func fisReadByteArray(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisReadByteArray: FileInputStream object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Set buffer to the byte array parameter.
	buffer, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
	if !ok {
		errMsg := "fisReadByteArray: Byte array parameter lacks a \"value\" field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Fill the buffer.
	nbytes, err := osFile.Read(buffer)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadByteArray osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// All is well - update the supplied buffer.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: buffer}
	params[1].(*object.Object).FieldTable["value"] = fld

	// Return the number of bytes.
	return int64(nbytes)
}

// "java/io/FileInputStream.read([BII)I"
func fisReadByteArrayOffset(params []interface{}) interface{} {

	// Get the file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisReadByteArrayOffset: FileInputStream object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Set buffer (buf1) to the byte array parameter.
	buf1, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
	if !ok {
		errMsg := "fisReadByteArrayOffset: Byte array parameter lacks a \"value\" field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Collect the offset and length parameter values.
	offset := params[2].(int64)
	length := params[3].(int64)

	// Check the parameters.
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(buf1))-offset) {
		errMsg := fmt.Sprintf("fisReadByteArrayOffset: error in parameters offset=%d length=%d bytes.length=%d",
			offset, length, len(buf1))
		return getGErrBlk(exceptions.IndexOutOfBoundsException, errMsg)
	}

	// Try read with a second buffer.
	buf2 := make([]byte, length)
	nbytes, err := osFile.Read(buf2)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadByteArrayOffset osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// All is well - Copy the bytes read into the original buffer, beginning at the offset.
	copy(buf1[offset:], buf2)

	// Update the parameter buffer.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: buf1}
	params[1].(*object.Object).FieldTable["value"] = fld

	// Return the number of bytes.
	return int64(nbytes)
}

// "java/io/FileInputStream.skip(J)J"
func fisSkip(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisSkip: FileInputStream object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Get skip count.
	count := params[1].(int64)

	// Skip.
	_, err := osFile.Seek(count, 1)
	if err != nil {
		errMsg := fmt.Sprintf("fisSkip osFile.Seek(%d) failed, reason: %s", count, err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Return skip count.
	return count
}

// "java/io/FileInputStream.close()V"
func fisClose(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "fisClose: FileInputStream object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Close the file.
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("fisSkip osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	return nil
}
