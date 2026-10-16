/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
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
			ParamSlots:   1,
			GFunction:    outputStreamWriteBytes,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/io/OutputStream.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots:   3,
			GFunction:    outputStreamWriteBytesRange,
			NeedsContext: true,
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

func outputStreamNullOutputStream([]any) any {
	return object.MakeEmptyObjectWithClassName(&types.ClassNameOutputStream)
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

// java/io/OutputStream.write([B)V
func outputStreamWriteBytes(params []any) any {
	// params[0]: frame stack
	// params[1]: java/io/OutputStream object
	// params[2]: Java byte array object
	if len(params) < 3 || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytes: params count wrong or byte array is null")
	}

	// Get Java byte array object.
	arrObj, ok := params[2].(*object.Object)
	if !ok || arrObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytes: byte array is null")
	}

	// Get Java byte array from the byte array object.
	javaBytes, ok := arrObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytes: byte array lacks value field")
	}

	// Let outputStreamWriteBytesRange carry on from here.
	return outputStreamWriteBytesRange([]any{params[0], params[1], arrObj, int64(0), int64(len(javaBytes))})
}

// java/io/OutputStream.write([BII)V
func outputStreamWriteBytesRange(params []any) any {
	// params[0]: frame stack
	// params[1]: java/io/OutputStream object
	// params[2]: Java byte array object
	// params[3]: int64 offset
	// params[4]: int64 length
	if len(params) < 5 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "outputStreamWriteBytesRange: expected 5 params")
	}

	// Get frame stack.
	fs, ok := params[0].(*list.List)
	if !ok || fs == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytesRange: frame stack is wrong type or null")
	}

	outputStreamObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(outputStreamObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytesRange: OutputStream object is null")
	}

	// Get Java byte array object.
	arrObj, ok := params[2].(*object.Object)
	if !ok || arrObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytesRange: byte array object is null")
	}

	// Get Java byte array from the byte array object.
	javaBytes, ok := arrObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "outputStreamWriteBytesRange: byte array lacks value field")
	}

	// Get offset.
	offset, ok := params[3].(int64)
	if !ok || offset < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException,
			fmt.Sprintf("outputStreamWriteBytesRange: offset(%v) is wrong type or < 0", offset))
	}

	// Get length.
	length, ok := params[4].(int64)
	if !ok || length < 0 || offset > int64(len(javaBytes))-length {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException,
			fmt.Sprintf("outputStreamWriteBytesRange: length(%v), offset=%d, len(byte array)=%d", length, offset, len(javaBytes)))
	}

	// Fast path: subclass backed by an OS file handle
	if fld, ok := outputStreamObj.FieldTable[ghelpers.FileHandle]; ok {
		if osFile, ok := fld.Fvalue.(*os.File); ok && osFile != nil {
			buf := object.GoByteArrayFromJavaByteArray(javaBytes[offset : offset+length])
			if _, err := osFile.Write(buf); err != nil {
				return ghelpers.GetGErrBlk(excNames.IOException, "outputStreamWriteBytesRange: "+err.Error())
			}
			return nil
		}
	}

	// Slower path: delegate to the subclass "write([BII)V" method.
	subclassName := *stringPool.GetStringPointer(outputStreamObj.KlassName)
	if subclassName == types.ClassNameOutputStream {
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError,
			fmt.Sprintf("outputStreamWriteBytesRange: impossible situation: OutputStream subclassed itself?!"))
	}

	// Run the FQN and wait for the method to complete.
	globals.GetGlobalRef().FuncRunJavaFromG(fs, subclassName, "write", "([BII)V", outputStreamObj, arrObj, offset, length)

	return nil
}
