/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
)

func Load_Io_OutputStream() {
	ghelpers.MethodSignatures["java/io/OutputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  outputStreamClose,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  outputStreamFlush,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.nullOutputStream()Ljava/io/OutputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  outputStreamNullOutputStream,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  outputStreamWriteInt,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.write([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  outputStreamWriteBytes,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  outputStreamWriteBytesRange,
		}
}

func outputStreamClose(params []any) any {
	if len(params) == 0 || object.IsNull(params[0]) {
		return nil
	}
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return nil
	}
	if fld, ok := self.FieldTable[ghelpers.FileHandle]; ok {
		if osFile, ok := fld.Fvalue.(*os.File); ok && osFile != nil {
			if err := osFile.Close(); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("outputStreamClose: %v", err))
			}
		}
	}
	return nil
}

func outputStreamFlush(params []any) any {
	if len(params) == 0 || object.IsNull(params[0]) {
		return nil
	}
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return nil
	}
	if fld, ok := self.FieldTable[ghelpers.FileHandle]; ok {
		if osFile, ok := fld.Fvalue.(*os.File); ok && osFile != nil {
			if err := osFile.Sync(); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("outputStreamFlush: flush: %v", err))
			}
		}
	}
	return nil
}

func outputStreamNullOutputStream(params []any) any {
	className := "java/io/OutputStream"
	return object.MakeEmptyObjectWithClassName(&className)
}

func outputStreamWriteInt(params []any) any {
	if len(params) == 0 || object.IsNull(params[0]) {
		return nil
	}
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return nil
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IOException, "outputStreamWriteInt: missing int argument")
	}
	wint, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "outputStreamWriteInt: missing int argument")
	}

	if fld, ok := self.FieldTable[ghelpers.FileHandle]; ok {
		if osFile, ok := fld.Fvalue.(*os.File); ok && osFile != nil {
			buf := []byte{byte(wint & 0xFF)}
			if _, err := osFile.Write(buf); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("outputStreamWriteInt: %v", err))
			}
		}
	}
	return nil
}

func outputStreamWriteBytes(params []any) any {
	if len(params) < 2 || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array is null")
	}

	arrObj, ok := params[1].(*object.Object)
	if !ok || arrObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array is null")
	}

	javaBytes, ok := arrObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array lacks value field")
	}

	return outputStreamWriteBytesRange([]any{params[0], params[1], int64(0), int64(len(javaBytes))})
}

func outputStreamWriteBytesRange(params []any) any {
	if len(params) < 2 || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array is null")
	}

	arrObj, ok := params[1].(*object.Object)
	if !ok || arrObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array is null")
	}

	javaBytes, ok := arrObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array lacks value field")
	}

	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "write: missing offset or length")
	}

	offset, okOff := params[2].(int64)
	length, okLen := params[3].(int64)
	if !okOff || !okLen {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "write: invalid offset or length")
	}

	bufLen := int64(len(javaBytes))
	if offset < 0 || length < 0 || offset > bufLen-length {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("write: bad params offset=%d length=%d len=%d", offset, length, bufLen))
	}

	if length == 0 {
		return nil
	}

	if object.IsNull(params[0]) {
		return nil
	}
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return nil
	}

	if fld, ok := self.FieldTable[ghelpers.FileHandle]; ok {
		if osFile, ok := fld.Fvalue.(*os.File); ok && osFile != nil {
			buf := object.GoByteArrayFromJavaByteArray(javaBytes[offset : offset+length])
			if _, err := osFile.Write(buf); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("write: %v", err))
			}
			return nil
		}
	}

	if self.KlassName != 0 {
		selfClassName := *(stringPool.GetStringPointer(self.KlassName))
		if selfClassName != "java/io/OutputStream" {
			method := selfClassName + ".write(I)V"
			for i := range length {
				res := ghelpers.Invoke(method, []any{self, int64(javaBytes[offset+i])})
				if res != nil {
					return res
				}
			}
		}
	}

	return nil
}
