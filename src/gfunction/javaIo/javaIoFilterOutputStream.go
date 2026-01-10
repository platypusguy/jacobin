/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
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

func Load_Io_FilterOutputStream() {

	ghelpers.MethodSignatures["java/io/FilterOutputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  filteroutputstreamInit,
		}

	ghelpers.MethodSignatures["java/io/FilterOutputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filteroutputstreamClose,
		}

	ghelpers.MethodSignatures["java/io/FilterOutputStream.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filteroutputstreamFlush,
		}

	ghelpers.MethodSignatures["java/io/FilterOutputStream.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  filteroutputstreamWrite,
		}

	ghelpers.MethodSignatures["java/io/FilterOutputStream.write([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  filteroutputstreamWriteBytes,
		}

	ghelpers.MethodSignatures["java/io/FilterOutputStream.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  filteroutputstreamWriteBytesRange,
		}
}

// "java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V"
func filteroutputstreamInit(params []interface{}) interface{} {
	// Expect params[1] to be an OutputStream-like object that carries ghelpers.FilePath and ghelpers.FileHandle
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamInit: missing OutputStream parameter")
	}
	underlying, _ := params[1].(*object.Object)
	if underlying == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "filteroutputstreamInit: underlying OutputStream is null")
	}
	fldPath, ok := underlying.FieldTable[ghelpers.FilePath]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamInit: underlying OutputStream lacks ghelpers.FilePath field")
	}
	fldHandle, ok := underlying.FieldTable[ghelpers.FileHandle]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamInit: underlying OutputStream lacks ghelpers.FileHandle field")
	}
	osFile, ok := fldHandle.Fvalue.(*os.File)
	if !ok || osFile == nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamInit: ghelpers.FileHandle is not *os.File")
	}
	// Validate the handle points to a file by calling Stat
	if _, err := osFile.Stat(); err != nil {
		pathStr := string(fldPath.Fvalue.([]byte))
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamInit: os.Stat(%s) failed: %v", pathStr, err))
	}
	// Copy fields onto this FilterOutputStream instance
	self, _ := params[0].(*object.Object)
	if self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "filteroutputstreamInit: self is not object")
	}
	self.FieldTable[ghelpers.FilePath] = fldPath
	self.FieldTable[ghelpers.FileHandle] = fldHandle
	return nil
}

func filteroutputstreamClose(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamClose: missing ghelpers.FileHandle field")
	}
	if err := osFile.Close(); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamClose: Close failed: %v", err))
	}
	return nil
}

func filteroutputstreamFlush(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamFlush: missing ghelpers.FileHandle field")
	}
	if err := osFile.Sync(); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamFlush: Sync failed: %v", err))
	}
	return nil
}

func filteroutputstreamWrite(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamWrite: missing ghelpers.FileHandle field")
	}
	wint, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamWrite: missing int argument")
	}
	buf := []byte{byte(wint % 256)}
	if _, err := osFile.Write(buf); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamWrite: Write failed: %v", err))
	}
	return nil
}

func filteroutputstreamWriteBytes(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamWriteBytes: missing ghelpers.FileHandle field")
	}
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamWriteBytes: byte[] lacks value field")
	}
	buf := object.GoByteArrayFromJavaByteArray(javaBytes)
	if _, err := osFile.Write(buf); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamWriteBytes: Write failed: %v", err))
	}
	return nil
}

func filteroutputstreamWriteBytesRange(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[ghelpers.FileHandle].Fvalue.(*os.File)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamWriteBytesRange: missing ghelpers.FileHandle field")
	}
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "filteroutputstreamWriteBytesRange: byte[] lacks value field")
	}
	buf := object.GoByteArrayFromJavaByteArray(javaBytes)
	offset := params[2].(int64)
	length := params[3].(int64)
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > int64(len(buf))-offset {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("filteroutputstreamWriteBytesRange: bad params offset=%d length=%d len=%d", offset, length, len(buf)))
	}
	if _, err := osFile.Write(buf[offset : offset+length]); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamWriteBytesRange: Write failed: %v", err))
	}
	return nil
}
