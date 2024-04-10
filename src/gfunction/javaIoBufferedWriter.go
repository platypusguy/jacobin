/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"os"
	"runtime"
)

func Load_Io_BufferedWriter() map[string]GMeth {

	MethodSignatures["java/io/BufferedWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/BufferedWriter.<init>(Ljava/io/OutputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  initOutputStreamWriter,
		}

	MethodSignatures["java/io/BufferedWriter.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  oswClose,
		}

	MethodSignatures["java/io/BufferedWriter.flush()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  oswFlush,
		}

	MethodSignatures["java/io/BufferedWriter.newLine()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bwNewLine,
		}

	MethodSignatures["java/io/BufferedWriter.write(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  oswWriteOneChar,
		}

	MethodSignatures["java/io/BufferedWriter.write([CII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  oswWriteCharBuffer,
		}

	MethodSignatures["java/io/BufferedWriter.write(Ljava/lang/String;II)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  oswWriteStringBuffer,
		}

	return MethodSignatures
}

func bwNewLine(params []interface{}) interface{} {

	// Get file handle.
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		errMsg := "bwNewLine: BufferedWriter object lacks a FileHandle field"
		return getGErrBlk(exceptions.IOException, errMsg)
	}

	// Derive the newline byte array. the file.
	var newline string
	if runtime.GOOS == "windows" {
		newline = "\\r\\n"
	} else {
		newline = "\\n"
	}
	bytes := []byte(newline)

	// Write newline byte array to the file.
	_, err := osFile.Write(bytes)
	if err != nil {
		errMsg := fmt.Sprintf("bwNewLine: osFile.Close() failed, reason: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	return nil
}
