/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"os"
)

func Load_Io_RandomAccessFile() {

	ghelpers.MethodSignatures["java/io/RandomAccessFile.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafInitString,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.<init>(Ljava/io/File;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafInitFile,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisClose,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.getChannel()Ljava/nio/channels/FileChannel;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.getFD()Ljava/io/FileDescriptor;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.getFilePointer()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafGetFilePointer,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.length()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafLength,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.length0()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafLength,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.open0(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafOpen0,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read0()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fisReadOne,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fisReadByteArray,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readBytes([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  fisReadByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readFully([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafReadFully,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readFully([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  rafReadFullyOffset,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readBoolean()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadBoolean,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readByte()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadByte,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readChar()C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadChar,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readDouble()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadDouble,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readFloat()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadFloat,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readInt()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadInt,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readLine()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadLine,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readLong()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadLong,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readShort()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadShort,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readUnsignedByte()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadUnsignedByte,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readUnsignedShort()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadUnsignedShort,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.readUTF()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rafReadUTF,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.seek(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafSeek,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.seek0(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafSeek,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.setLength(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafSetLength,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.setLength0(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafSetLength,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWrite,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.write0(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWrite,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeBoolean(Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteBoolean,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeByte(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteByte,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeShort(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteShort,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeChar(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteChar,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeInt(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteInt,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeLong(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafWriteLong,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeFloat(F)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteFloat,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeDouble(D)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  rafWriteDouble,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeBytes(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteBytes,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeChars(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteChars,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeUTF(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteUTF,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.write([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rafWriteByteArray,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  rafWriteByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/RandomAccessFile.writeBytes([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  rafWriteByteArrayOffset,
		}

	// ----------------------------------------------------------
	// initIDs - ghelpers.JustReturn
	// This is a private function that call C native functions.
	// ----------------------------------------------------------

	ghelpers.MethodSignatures["java/io/RandomAccessFile.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

}

// "java/io/RandomAccessFile.open0(Ljava/lang/String;I)V"
func rafOpen0(params []interface{}) interface{} {
	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))

	// Mode flags from java.io.RandomAccessFile
	// O_RDONLY = 1
	// O_RDWR   = 2
	// O_SYNC   = 4
	// O_DSYNC  = 8
	mode := params[2].(int64)

	var modeInt int
	if (mode & 2) != 0 {
		modeInt = os.O_RDWR | os.O_CREATE
	} else {
		modeInt = os.O_RDONLY
	}

	if (mode & 4) != 0 {
		modeInt |= os.O_SYNC
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, ghelpers.CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafOpen0: os.OpenFile(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld := object.Field{Ftype: types.JavaByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil
}

// "java/io/RandomAccessFile.<init>(Ljava/lang/String;Ljava/lang/String;)V"
// RandomAccessFile raf = new RandomAccessFile(Stringname, Stringmode);
func rafInitString(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	pathStr := object.GoStringFromStringObject(params[1].(*object.Object))

	// Mode.
	var modeInt int
	modeStr := object.GoStringFromStringObject(params[2].(*object.Object))
	switch modeStr {
	case "r":
		modeInt = os.O_RDONLY
	case "rw", "rws", "rwd":
		modeInt = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		errMsg := fmt.Sprintf("rafInitString: mode string (%s) invalid", modeStr)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, ghelpers.CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafInitString: os.OpenFile(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld := object.Field{Ftype: types.JavaByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil

}

// "java/io/RandomAccessFile.<init>(Ljava/io/File;Ljava/lang/String;)V"
// RandomAccessFile raf = new RandomAccessFile(Fileobject, Stringmode);
func rafInitFile(params []interface{}) interface{} {

	// Using the argument path string, open the file for read-only.
	obj := params[1].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FilePath]
	if !ok {
		errMsg := "rafInitFile: java/io/File object is missing the ghelpers.FilePath field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Mode.
	var modeInt int
	modeStr := object.GoStringFromStringObject(params[2].(*object.Object))
	switch modeStr {
	case "r":
		modeInt = os.O_RDONLY
	case "rw", "rws", "rwd":
		modeInt = os.O_RDWR | os.O_CREATE | os.O_APPEND
	default:
		errMsg := fmt.Sprintf("rafInitFile: mode string (%s) invalid", modeStr)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Open the file in the specified mode.
	osFile, err := os.OpenFile(pathStr, modeInt, ghelpers.CreateFilePermissions)
	if err != nil {
		errMsg := fmt.Sprintf("rafInitFile: os.Open(%s) failed, reason: %s", pathStr, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy the file path field into the RandomAccessFile object.
	fld = object.Field{Ftype: types.JavaByteArray, Fvalue: []byte(pathStr)}
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fld

	// Copy the file handle into the RandomAccessFile object.
	fld = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fld

	return nil

}

// "java/io/RandomAccessFile.getFilePointer()J"
// Get current file position (offset from the beginning of file).
func rafGetFilePointer(params []interface{}) interface{} {

	// Get the open file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafGetFilePointer: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the current position relative to the beginning of file.
	posn, err := osFile.Seek(0, 1)
	if err != nil {
		errMsg := fmt.Sprintf("rafGetFilePointer: osFile.Seek(0, 1) failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return posn

}

// "java/io/RandomAccessFile.readFully([B)V"
func rafReadFully(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafReadFully: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the byte array value.
	arrObj := params[1].(*object.Object)
	fld, ok = arrObj.FieldTable["value"]
	if !ok {
		errMsg := "rafReadFully: byte array lacks a \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	javaBytes := fld.Fvalue.([]types.JavaByte)
	buffer := object.GoByteArrayFromJavaByteArray(javaBytes)

	// Read until the buffer is full.
	_, err := io.ReadFull(osFile, buffer)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			// In Java, readFully throws EOFException if EOF is reached.
			// Since EOFException is not in excNames, we use IOException.
			return ghelpers.GetGErrBlk(excNames.IOException, "rafReadFully: EOF reached before reading all bytes")
		}
		errMsg := fmt.Sprintf("rafReadFully: osFile.Read failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Update the byte array object with the read bytes.
	javaBytes = object.JavaByteArrayFromGoByteArray(buffer)
	arrObj.FieldTable["value"] = object.Field{Ftype: types.JavaByteArray, Fvalue: javaBytes}

	return nil
}

// "java/io/RandomAccessFile.setLength(J)V"
func rafSetLength(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafSetLength: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the new length.
	newLength := params[1].(int64)

	// Truncate the file to the new length.
	err := osFile.Truncate(newLength)
	if err != nil {
		errMsg := fmt.Sprintf("rafSetLength: osFile.Truncate(%d) failed, reason: %s", newLength, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/RandomAccessFile.length()J"
func rafLength(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafLength: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	fi, err := osFile.Stat()
	if err != nil {
		errMsg := fmt.Sprintf("rafLength: osFile.Stat failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return fi.Size()
}

// "java/io/RandomAccessFile.seek(J)V"
func rafSeek(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafSeek: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the new position.
	pos := params[1].(int64)

	_, err := osFile.Seek(pos, io.SeekStart)
	if err != nil {
		errMsg := fmt.Sprintf("rafSeek: osFile.Seek(%d) failed, reason: %s", pos, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/RandomAccessFile.readFully([BII)V"
func rafReadFullyOffset(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafReadFullyOffset: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the byte array value.
	arrObj := params[1].(*object.Object)
	fld, ok = arrObj.FieldTable["value"]
	if !ok {
		errMsg := "rafReadFullyOffset: byte array lacks a \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	javaBytes := fld.Fvalue.([]types.JavaByte)

	offset := int(params[2].(int64))
	length := int(params[3].(int64))

	if offset < 0 || length < 0 || offset+length > len(javaBytes) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "rafReadFullyOffset: index out of bounds")
	}

	buffer := make([]byte, length)

	// Read until the buffer is full.
	_, err := io.ReadFull(osFile, buffer)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return ghelpers.GetGErrBlk(excNames.IOException, "rafReadFullyOffset: EOF reached before reading all bytes")
		}
		errMsg := fmt.Sprintf("rafReadFullyOffset: osFile.Read failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Update the byte array object with the read bytes.
	for i := range length {
		javaBytes[offset+i] = types.JavaByte(buffer[i])
	}

	return nil
}

// "java/io/RandomAccessFile.write(I)V"
func rafWrite(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafWrite: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the byte to write.
	b := byte(params[1].(int64))

	_, err := osFile.Write([]byte{b})
	if err != nil {
		errMsg := fmt.Sprintf("rafWrite: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/RandomAccessFile.write([B)V"
func rafWriteByteArray(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafWriteByteArray: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the byte array value.
	arrObj := params[1].(*object.Object)
	fld, ok = arrObj.FieldTable["value"]
	if !ok {
		errMsg := "rafWriteByteArray: byte array lacks a \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	javaBytes := fld.Fvalue.([]types.JavaByte)
	buffer := object.GoByteArrayFromJavaByteArray(javaBytes)

	_, err := osFile.Write(buffer)
	if err != nil {
		errMsg := fmt.Sprintf("rafWriteByteArray: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/RandomAccessFile.write([BII)V"
func rafWriteByteArrayOffset(params []interface{}) interface{} {

	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "rafWriteByteArrayOffset: java/io/RandomAccessFile object is missing the ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	var osFile *os.File = fld.Fvalue.(*os.File)

	// Get the byte array value.
	arrObj := params[1].(*object.Object)
	fld, ok = arrObj.FieldTable["value"]
	if !ok {
		errMsg := "rafWriteByteArrayOffset: byte array lacks a \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	javaBytes := fld.Fvalue.([]types.JavaByte)

	offset := int(params[2].(int64))
	length := int(params[3].(int64))

	if offset < 0 || length < 0 || offset+length > len(javaBytes) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "rafWriteByteArrayOffset: index out of bounds")
	}

	buffer := make([]byte, length)
	for i := range length {
		buffer[i] = byte(javaBytes[offset+i])
	}

	_, err := osFile.Write(buffer)
	if err != nil {
		errMsg := fmt.Sprintf("rafWriteByteArrayOffset: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	return nil
}

// "java/io/RandomAccessFile.readBoolean()Z"
func rafReadBoolean(params []interface{}) interface{} {
	val := rafReadUnsignedByte(params)
	if err, ok := val.(*object.Object); ok {
		return err
	}
	if val.(int64) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// "java/io/RandomAccessFile.readByte()B"
func rafReadByte(params []interface{}) interface{} {
	val := rafReadUnsignedByte(params)
	if err, ok := val.(*object.Object); ok {
		return err
	}
	return int64(int8(val.(int64)))
}

// "java/io/RandomAccessFile.readChar()C"
func rafReadChar(params []interface{}) interface{} {
	return rafReadUnsignedShort(params)
}

// "java/io/RandomAccessFile.readDouble()D"
func rafReadDouble(params []interface{}) interface{} {
	val := rafReadLong(params)
	if err, ok := val.(*object.Object); ok {
		return err
	}
	return math.Float64frombits(uint64(val.(int64)))
}

// "java/io/RandomAccessFile.readFloat()F"
func rafReadFloat(params []interface{}) interface{} {
	val := rafReadInt(params)
	if err, ok := val.(*object.Object); ok {
		return err
	}
	return float64(math.Float32frombits(uint32(val.(int64))))
}

// "java/io/RandomAccessFile.readInt()I"
func rafReadInt(params []interface{}) interface{} {
	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadInt: missing FileHandle")
	}
	osFile := fld.Fvalue.(*os.File)

	var b [4]byte
	_, err := io.ReadFull(osFile, b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadInt: EOF reached")
	}
	return int64(int32(binary.BigEndian.Uint32(b[:])))
}

// "java/io/RandomAccessFile.readLong()J"
func rafReadLong(params []interface{}) interface{} {
	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadLong: missing FileHandle")
	}
	osFile := fld.Fvalue.(*os.File)

	var b [8]byte
	_, err := io.ReadFull(osFile, b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadLong: EOF reached")
	}
	return int64(binary.BigEndian.Uint64(b[:]))
}

// "java/io/RandomAccessFile.readShort()S"
func rafReadShort(params []interface{}) interface{} {
	val := rafReadUnsignedShort(params)
	if err, ok := val.(*object.Object); ok {
		return err
	}
	return int64(int16(val.(int64)))
}

// "java/io/RandomAccessFile.readUnsignedByte()I"
func rafReadUnsignedByte(params []interface{}) interface{} {
	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadUnsignedByte: missing FileHandle")
	}
	osFile := fld.Fvalue.(*os.File)

	var b [1]byte
	_, err := osFile.Read(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadUnsignedByte: EOF reached")
	}
	return int64(b[0])
}

// "java/io/RandomAccessFile.readUnsignedShort()I"
func rafReadUnsignedShort(params []interface{}) interface{} {
	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadUnsignedShort: missing FileHandle")
	}
	osFile := fld.Fvalue.(*os.File)

	var b [2]byte
	_, err := io.ReadFull(osFile, b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadUnsignedShort: EOF reached")
	}
	return int64(binary.BigEndian.Uint16(b[:]))
}

// "java/io/RandomAccessFile.readLine()Ljava/lang/String;"
func rafReadLine(params []interface{}) interface{} {
	// Get file handle.
	obj := params[0].(*object.Object)
	fld, ok := obj.FieldTable[ghelpers.FileHandle]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadLine: missing FileHandle")
	}
	osFile := fld.Fvalue.(*os.File)

	var line []byte
	var b [1]byte
	for {
		_, err := osFile.Read(b[:])
		if err != nil {
			if len(line) == 0 {
				return object.IsNull(nil)
			}
			break
		}
		if b[0] == '\n' {
			break
		}
		if b[0] == '\r' {
			// Check for next \n
			pos, _ := osFile.Seek(0, io.SeekCurrent)
			_, err = osFile.Read(b[:])
			if err == nil && b[0] != '\n' {
				osFile.Seek(pos, io.SeekStart)
			}
			break
		}
		line = append(line, b[0])
	}

	// In DataInputStream.readLine, it's NOT UTF-8, it's just ISO-8859-1 (one byte per char)
	// RandomAccessFile.readLine also follows this.
	return object.StringObjectFromGoString(string(line))
}

// "java/io/RandomAccessFile.readUTF()Ljava/lang/String;"
func rafReadUTF(params []interface{}) interface{} {
	lenVal := rafReadUnsignedShort(params)
	if err, ok := lenVal.(*object.Object); ok {
		return err
	}
	utflen := lenVal.(int64)

	// Get file handle.
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)

	bytearr := make([]byte, utflen)
	_, err := io.ReadFull(osFile, bytearr)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafReadUTF: EOF reached")
	}

	// Simple Modified UTF-8 to String conversion.
	var chararr = make([]rune, 0, utflen)
	var count = 0
	for count < int(utflen) {
		c := int(bytearr[count] & 0xff)
		if c > 127 {
			break
		}
		count++
		chararr = append(chararr, rune(c))
	}

	for count < int(utflen) {
		c := int(bytearr[count] & 0xff)
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			/* 0xxxxxxx*/
			count++
			chararr = append(chararr, rune(c))
		case 12, 13:
			/* 110x xxxx   10xx xxxx*/
			count += 2
			if count > int(utflen) {
				return ghelpers.GetGErrBlk(excNames.IOException, "malformed input: partial character at end")
			}
			char2 := int(bytearr[count-1])
			if (char2 & 0xC0) != 0x80 {
				return ghelpers.GetGErrBlk(excNames.IOException, "malformed input around byte "+fmt.Sprint(count))
			}
			chararr = append(chararr, rune(((c&0x1F)<<6)|(char2&0x3F)))
		case 14:
			/* 1110 xxxx  10xx xxxx  10xx xxxx */
			count += 3
			if count > int(utflen) {
				return ghelpers.GetGErrBlk(excNames.IOException, "malformed input: partial character at end")
			}
			char2 := int(bytearr[count-2])
			char3 := int(bytearr[count-1])
			if ((char2 & 0xC0) != 0x80) || ((char3 & 0xC0) != 0x80) {
				return ghelpers.GetGErrBlk(excNames.IOException, "malformed input around byte "+fmt.Sprint(count-1))
			}
			chararr = append(chararr, rune(((c&0x0F)<<12)|((char2&0x3F)<<6)|((char3&0x3F)<<0)))
		default:
			/* 10xx xxxx,  1111 xxxx */
			return ghelpers.GetGErrBlk(excNames.IOException, "malformed input around byte "+fmt.Sprint(count))
		}
	}

	return object.StringObjectFromGoString(string(chararr))
}

// "java/io/RandomAccessFile.writeBoolean(Z)V"
func rafWriteBoolean(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := params[1].(int64)
	var b byte
	if v != 0 {
		b = 1
	}
	_, err := osFile.Write([]byte{b})
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteBoolean: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeByte(I)V"
func rafWriteByte(params []interface{}) interface{} {
	return rafWrite(params)
}

// "java/io/RandomAccessFile.writeShort(I)V"
func rafWriteShort(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := int16(params[1].(int64))
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(v))
	_, err := osFile.Write(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteShort: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeChar(I)V"
func rafWriteChar(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := uint16(params[1].(int64))
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], v)
	_, err := osFile.Write(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteChar: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeInt(I)V"
func rafWriteInt(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := int32(params[1].(int64))
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	_, err := osFile.Write(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteInt: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeLong(J)V"
func rafWriteLong(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := params[1].(int64)
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(v))
	_, err := osFile.Write(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteLong: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeFloat(F)V"
func rafWriteFloat(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := float32(params[1].(float64))
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], math.Float32bits(v))
	_, err := osFile.Write(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteFloat: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeDouble(D)V"
func rafWriteDouble(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	v := params[1].(float64)
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], math.Float64bits(v))
	_, err := osFile.Write(b[:])
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteDouble: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeBytes(Ljava/lang/String;)V"
func rafWriteBytes(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	s := object.GoStringFromStringObject(params[1].(*object.Object))
	b := make([]byte, len(s))
	for i := range len(s) {
		b[i] = byte(s[i])
	}
	_, err := osFile.Write(b)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteBytes: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeChars(Ljava/lang/String;)V"
func rafWriteChars(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	s := object.GoStringFromStringObject(params[1].(*object.Object))
	b := make([]byte, len(s)*2)
	for i, r := range s {
		binary.BigEndian.PutUint16(b[i*2:], uint16(r))
	}
	_, err := osFile.Write(b)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteChars: "+err.Error())
	}
	return nil
}

// "java/io/RandomAccessFile.writeUTF(Ljava/lang/String;)V"
func rafWriteUTF(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	str := object.GoStringFromStringObject(params[1].(*object.Object))

	var utflen int
	for _, r := range str {
		if r >= 0x0001 && r <= 0x007F {
			utflen++
		} else if r > 0x07FF {
			utflen += 3
		} else {
			utflen += 2
		}
	}

	if utflen > 65535 {
		return ghelpers.GetGErrBlk(excNames.UTFDataFormatException, fmt.Sprintf("encoded string too long: %d bytes", utflen))
	}

	bytearr := make([]byte, utflen+2)
	binary.BigEndian.PutUint16(bytearr[0:], uint16(utflen))

	count := 2
	for _, r := range str {
		if r >= 0x0001 && r <= 0x007F {
			bytearr[count] = byte(r)
			count++
		} else if r > 0x07FF {
			bytearr[count] = byte(0xE0 | ((r >> 12) & 0x0F))
			bytearr[count+1] = byte(0x80 | ((r >> 6) & 0x3F))
			bytearr[count+2] = byte(0x80 | ((r >> 0) & 0x3F))
			count += 3
		} else {
			bytearr[count] = byte(0xC0 | ((r >> 6) & 0x1F))
			bytearr[count+1] = byte(0x80 | ((r >> 0) & 0x3F))
			count += 2
		}
	}

	_, err := osFile.Write(bytearr)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "rafWriteUTF: "+err.Error())
	}
	return nil
}
