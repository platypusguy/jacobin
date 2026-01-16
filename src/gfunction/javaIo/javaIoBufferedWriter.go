/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

func Load_Io_BufferedWriter() {

	ghelpers.MethodSignatures["java/io/BufferedWriter.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// Constructors
	ghelpers.MethodSignatures["java/io/BufferedWriter.<init>(Ljava/io/Writer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bufferedWriterInit,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.<init>(Ljava/io/Writer;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Core operations
	ghelpers.MethodSignatures["java/io/BufferedWriter.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bwClose,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bwFlush,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.newLine()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bufferedWriterNewLine,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bwWriteOneChar,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.write([CII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  bwWriteCharBuffer,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.write(Ljava/lang/String;II)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  bwWriteStringBuffer,
		}

	// Append variants — currently trapped
	ghelpers.MethodSignatures["java/io/BufferedWriter.append(C)Ljava/io/Writer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.append(Ljava/lang/CharSequence;)Ljava/io/Writer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/BufferedWriter.append(Ljava/lang/CharSequence;II)Ljava/io/Writer;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}
}

// "java/io/BufferedWriter.<init>(Ljava/io/Writer;)V"
func bufferedWriterInit(params []interface{}) interface{} {
    // Copy FilePath and FileHandle from the provided Writer to this BufferedWriter.
    // Ensure the target object's FieldTable map is initialized before assignment.
    if params[0].(*object.Object).FieldTable == nil {
        params[0].(*object.Object).FieldTable = make(map[string]object.Field)
    }
    fldPath, ok := params[1].(*object.Object).FieldTable[ghelpers.FilePath]
    if !ok {
        errMsg := "bufferedWriterInit: Writer object lacks a ghelpers.FilePath field"
        return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
    }

	fldHandle, ok := params[1].(*object.Object).FieldTable[ghelpers.FileHandle]
	if !ok {
		errMsg := "bufferedWriterInit: Writer object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Validate the handle by stat'ing it if possible.
	if osFile, ok2 := fldHandle.Fvalue.(*os.File); ok2 {
		if _, err := osFile.Stat(); err != nil {
			pathStr := string(fldPath.Fvalue.([]byte))
			errMsg := fmt.Sprintf("bufferedWriterInit: os.Stat(%s) failed, reason: %s", pathStr, err.Error())
			return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
		}
	} else {
		errMsg := "bufferedWriterInit: Writer object's ghelpers.FileHandle is not an *os.File"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	// Copy path and handle into this object.
	params[0].(*object.Object).FieldTable[ghelpers.FilePath] = fldPath
	params[0].(*object.Object).FieldTable[ghelpers.FileHandle] = fldHandle

	return nil
}

func bwClose(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bwClose: BufferedWriter object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	if err := osFile.Close(); err != nil {
		errMsg := fmt.Sprintf("bwClose: osFile.Close() failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

func bwFlush(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bwFlush: BufferedWriter object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	if err := osFile.Sync(); err != nil {
		errMsg := fmt.Sprintf("bwFlush: osFile.Sync() failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/io/BufferedWriter.newLine()V"
func bufferedWriterNewLine(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bufferedWriterNewLine: BufferedWriter object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	// Java uses platform-independent newline via writer; here we use \n
	if _, err := osFile.Write([]byte{'\n'}); err != nil {
		errMsg := fmt.Sprintf("bufferedWriterNewLine: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/io/BufferedWriter.write(I)V"
func bwWriteOneChar(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	osFile, ok := obj.FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bwWriteOneChar: BufferedWriter object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	wint, ok := params[1].(int64)
	if !ok {
		errMsg := "bwWriteOneChar: Error in integer argument"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	buffer := []byte{byte(wint % 256)}
	if _, err := osFile.Write(buffer); err != nil {
		errMsg := fmt.Sprintf("bwWriteOneChar: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/io/BufferedWriter.write([CII)V"
func bwWriteCharBuffer(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bwWriteCharBuffer: BufferedWriter object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	intArray, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]int64)
	if !ok {
		errMsg := "bwWriteCharBuffer: Trouble with value field ([]int64)"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	offset := params[2].(int64)
	length := params[3].(int64)

	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(intArray))-offset) {
		errMsg := fmt.Sprintf("bwWriteCharBuffer: Error in parameters: offset=%d, length=%d, char.array.length=%d",
			offset, length, len(intArray))
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	outBytes := make([]byte, length)
	for ii := int64(0); ii < length; ii++ {
		outBytes[ii] = byte(intArray[offset+ii])
	}

	if _, err := osFile.Write(outBytes); err != nil {
		errMsg := fmt.Sprintf("bwWriteCharBuffer: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/io/BufferedWriter.write(Ljava/lang/String;II)V"
func bwWriteStringBuffer(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bwWriteStringBuffer: BufferedWriter object lacks a ghelpers.FileHandle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := "bwWriteStringBuffer: Trouble with value field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	paramBytes := object.GoByteArrayFromJavaByteArray(javaBytes)
	offset := params[2].(int64)
	length := params[3].(int64)

	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > (int64(len(paramBytes))-offset) {
		errMsg := fmt.Sprintf("bwWriteStringBuffer: Error in parameters: offset=%d, length=%d, char.array.length=%d",
			offset, length, len(paramBytes))
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	outBytes := make([]byte, length)
	for ii := int64(0); ii < length; ii++ {
		outBytes[ii] = paramBytes[offset+ii]
	}

	if _, err := osFile.Write(outBytes); err != nil {
		errMsg := fmt.Sprintf("bwWriteStringBuffer: osFile.Write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}
