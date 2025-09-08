/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
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

func Load_Io_FilterOutputStream() {

	MethodSignatures["java/io/FilterOutputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  filteroutputstreamInit,
		}

	MethodSignatures["java/io/FilterOutputStream.close()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  filteroutputstreamClose,
		}

	MethodSignatures["java/io/FilterOutputStream.flush()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  filteroutputstreamFlush,
		}

	MethodSignatures["java/io/FilterOutputStream.write(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  filteroutputstreamWrite,
		}

	MethodSignatures["java/io/FilterOutputStream.write([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  filteroutputstreamWriteBytes,
		}

	MethodSignatures["java/io/FilterOutputStream.write([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  filteroutputstreamWriteBytesRange,
		}
}

// "java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V"
func filteroutputstreamInit(params []interface{}) interface{} {
	// Expect params[1] to be an OutputStream-like object that carries FilePath and FileHandle
	if len(params) < 2 {
		return getGErrBlk(excNames.IOException, "filteroutputstreamInit: missing OutputStream parameter")
	}
	underlying, _ := params[1].(*object.Object)
	if underlying == nil {
		return getGErrBlk(excNames.NullPointerException, "filteroutputstreamInit: underlying OutputStream is null")
	}
	fldPath, ok := underlying.FieldTable[FilePath]
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamInit: underlying OutputStream lacks FilePath field")
	}
	fldHandle, ok := underlying.FieldTable[FileHandle]
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamInit: underlying OutputStream lacks FileHandle field")
	}
	osFile, ok := fldHandle.Fvalue.(*os.File)
	if !ok || osFile == nil {
		return getGErrBlk(excNames.IOException, "filteroutputstreamInit: FileHandle is not *os.File")
	}
	// Validate the handle points to a file by calling Stat
	if _, err := osFile.Stat(); err != nil {
		pathStr := string(fldPath.Fvalue.([]byte))
		return getGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamInit: os.Stat(%s) failed: %v", pathStr, err))
	}
	// Copy fields onto this FilterOutputStream instance
	self, _ := params[0].(*object.Object)
	if self == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "filteroutputstreamInit: self is not object")
	}
	self.FieldTable[FilePath] = fldPath
	self.FieldTable[FileHandle] = fldHandle
	return nil
}

func filteroutputstreamClose(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamClose: missing FileHandle field")
	}
	if err := osFile.Close(); err != nil {
		return getGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamClose: Close failed: %v", err))
	}
	return nil
}

func filteroutputstreamFlush(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamFlush: missing FileHandle field")
	}
	if err := osFile.Sync(); err != nil {
		return getGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamFlush: Sync failed: %v", err))
	}
	return nil
}

func filteroutputstreamWrite(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamWrite: missing FileHandle field")
	}
	wint, ok := params[1].(int64)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamWrite: missing int argument")
	}
	buf := []byte{byte(wint % 256)}
	if _, err := osFile.Write(buf); err != nil {
		return getGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamWrite: Write failed: %v", err))
	}
	return nil
}

func filteroutputstreamWriteBytes(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamWriteBytes: missing FileHandle field")
	}
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamWriteBytes: byte[] lacks value field")
	}
	buf := object.GoByteArrayFromJavaByteArray(javaBytes)
	if _, err := osFile.Write(buf); err != nil {
		return getGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamWriteBytes: Write failed: %v", err))
	}
	return nil
}

func filteroutputstreamWriteBytesRange(params []interface{}) interface{} {
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamWriteBytesRange: missing FileHandle field")
	}
	javaBytes, ok := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return getGErrBlk(excNames.IOException, "filteroutputstreamWriteBytesRange: byte[] lacks value field")
	}
	buf := object.GoByteArrayFromJavaByteArray(javaBytes)
	offset := params[2].(int64)
	length := params[3].(int64)
	if length == 0 {
		return int64(0)
	}
	if length < 0 || offset < 0 || length > int64(len(buf))-offset {
		return getGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("filteroutputstreamWriteBytesRange: bad params offset=%d length=%d len=%d", offset, length, len(buf)))
	}
	if _, err := osFile.Write(buf[offset : offset+length]); err != nil {
		return getGErrBlk(excNames.IOException, fmt.Sprintf("filteroutputstreamWriteBytesRange: Write failed: %v", err))
	}
	return nil
}
