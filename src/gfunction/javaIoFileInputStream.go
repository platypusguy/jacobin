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

	MethodSignatures["java/io/FileInputStream.<init>(Ljava/io/File;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initFileInputStream,
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

	// -------------------
	// <clinit> justReturn
	// -------------------

	MethodSignatures["java/io/FileInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// ------------------
	// initIDs justReturn
	// ------------------

	MethodSignatures["java/io/FileInputStream.initIDs()V"] = // private function that calls C native functions
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

	MethodSignatures["java/nio/channels/FileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileChannel,
		}

	MethodSignatures["java/io/FileDescriptor.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileDescriptor,
		}

	return MethodSignatures
}

// "java/io/FileInputStream.<init>(Ljava/io/File;])V"
func initFileInputStream(params []interface{}) interface{} {
	fld1 := params[1].(*object.Object).FieldTable["filePath"]
	inPathStr := string(fld1.Fvalue.([]byte))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileInputStream: os.Open(%s) returned: %s", inPathStr, err.Error())
		return getGErrBlk(exceptions.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := fld1
	params[0].(*object.Object).FieldTable["filePath"] = fld

	// Field "osfile" = Golang *os.File from os.Open
	fld = object.Field{Ftype: types.Ref, Fvalue: osFile}
	params[0].(*object.Object).FieldTable["osfile"] = fld

	return nil
}

// "java/io/FileInputStream.<init>(Ljava/lang/String;])V"
func initFileInputStreamString(params []interface{}) interface{} {
	inPathStr := object.GoStringFromStringObject(params[1].(*object.Object))
	osFile, err := os.Open(inPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("initFileInputStreamString: os.Open(%s) returned: %s", inPathStr, err.Error())
		return getGErrBlk(exceptions.FileNotFoundException, errMsg)
	}

	// Copy java/io/File path
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(inPathStr)}
	params[0].(*object.Object).FieldTable["filePath"] = fld

	// Field "osfile" = Golang *os.File from os.Open
	fld = object.Field{Ftype: types.Ref, Fvalue: osFile}
	params[0].(*object.Object).FieldTable["osfile"] = fld

	return nil
}

// "java/io/FileInputStream.available()I"
func fisAvailable(params []interface{}) interface{} {
	osFile := params[0].(*object.Object).FieldTable["osfile"].Fvalue.(*os.File)
	fileInfo, err := osFile.Stat()
	if err != nil {
		path := string(params[0].(*object.Object).FieldTable["path"].Fvalue.([]byte))
		errMsg := fmt.Sprintf("fisAvailable: osFile.Stat(%s) returned: %s", path, err.Error())
		return getGErrBlk(exceptions.FileNotFoundException, errMsg)
	}
	return fileInfo.Size()
}

// "java/io/FileInputStream.read()I"
func fisReadOne(params []interface{}) interface{} {
	osFile := params[0].(*object.Object).FieldTable["osfile"].Fvalue.(*os.File)
	buffer := make([]byte, 1)

	// Try read.
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
	osFile := params[0].(*object.Object).FieldTable["osfile"].Fvalue.(*os.File)
	buffer := params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)

	// Try read.
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
	osFile := params[0].(*object.Object).FieldTable["osfile"].Fvalue.(*os.File)
	buf1 := params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
	offset := params[2].(int64)
	length := params[3].(int64)

	// Check parameters.
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(buf1))-offset) {
		errMsg := fmt.Sprintf("fisReadByteArrayOffset: error in parameters offset=%d length=%d bytes.length=%d",
			offset, length, len(buf1))
		return getGErrBlk(exceptions.IndexOutOfBoundsException, errMsg)
	}

	// Try read.
	buf2 := make([]byte, length)
	nbytes, err := osFile.Read(buf2)
	if err == io.EOF {
		return int64(-1) // return -1 on EOF
	}
	if err != nil {
		errMsg := fmt.Sprintf("fisReadByteArrayOffset osFile.Read failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// All is well - Copy the bytes read into the user buffer, beginning at the offset.
	copy(buf1[offset:], buf2)

	// Update the supplied buffer.
	fld := object.Field{Ftype: types.ByteArray, Fvalue: buf1}
	params[1].(*object.Object).FieldTable["value"] = fld

	// Return the number of bytes.
	return int64(nbytes)
}

// "java/io/FileInputStream.skip(J)J"
func fisSkip(params []interface{}) interface{} {
	osFile := params[0].(*object.Object).FieldTable["osfile"].Fvalue.(*os.File)
	count := params[1].(int64)
	_, err := osFile.Seek(count, 1)
	if err != nil {
		errMsg := fmt.Sprintf("fisSkip osFile.Seek(%d) failed, reason: %s", count, err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	return count
}

// "java/io/FileInputStream.close()V"
func fisClose(params []interface{}) interface{} {
	osFile := params[0].(*object.Object).FieldTable["osfile"].Fvalue.(*os.File)
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("fisSkip osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	return nil
}

// -------------------- Traps ----------------------------------

func trapFileDescriptor([]interface{}) interface{} {
	errMsg := "FileDescriptor class is not yet supported !!"
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

func trapFileChannel([]interface{}) interface{} {
	errMsg := "FileChannel class is not yet supported !!"
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}
