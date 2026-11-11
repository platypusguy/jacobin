/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"container/list"
	"fmt"
	"math"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

func asGErrBlk(val any) (*ghelpers.GErrBlk, bool) {
	if val == nil {
		return nil, false
	}
	switch val.(type) {
	case error:
		return ghelpers.GetGErrBlk(excNames.IOException, val.(error).Error()), true
	case *ghelpers.GErrBlk:
		return val.(*ghelpers.GErrBlk), true
	}
	return nil, false
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func getInt64(val any) int64 {
	v, ok := ghelpers.GetInt64Val(val)
	if !ok {
		switch n := val.(type) {
		case int:
			return int64(n)
		case int32:
			return int64(n)
		case uint32:
			return int64(n)
		}
	}
	return v
}

func getObjectClassName(obj *object.Object) string {
	if obj == nil || object.IsNull(obj) {
		return ""
	}
	ptr := stringPool.GetStringPointer(obj.KlassName)
	if ptr == nil {
		return ""
	}
	return *ptr
}

func Load_Io_InputStream() {
	ghelpers.MethodSignatures["java/io/InputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/InputStream.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/InputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  inputStreamAvailable,
		}

	ghelpers.MethodSignatures["java/io/InputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/InputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/InputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["java/io/InputStream.nullInputStream()Ljava/io/InputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  inputStreamNullInputStream,
		}

	ghelpers.MethodSignatures["java/io/InputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  inputStreamReadByte,
		}

	ghelpers.MethodSignatures["java/io/InputStream.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamReadIntoByteArray,
		}

	ghelpers.MethodSignatures["java/io/InputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  inputStreamReadRange,
		}

	ghelpers.MethodSignatures["java/io/InputStream.readAllBytes()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  inputStreamReadAllBytes,
		}

	ghelpers.MethodSignatures["java/io/InputStream.readNBytes(I)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamReadNBytesLen,
		}

	ghelpers.MethodSignatures["java/io/InputStream.readNBytes([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  inputStreamReadNBytesBuf,
		}

	ghelpers.MethodSignatures["java/io/InputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  inputStreamReset,
		}

	ghelpers.MethodSignatures["java/io/InputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamSkip,
		}

	ghelpers.MethodSignatures["java/io/InputStream.skipNBytes(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamSkipNBytes,
		}

	ghelpers.MethodSignatures["java/io/InputStream.transferTo(Ljava/io/OutputStream;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamTransferTo,
		}

	// Standard Object traps
	ghelpers.MethodSignatures["java/io/InputStream.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.hashCode()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/io/InputStream.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
}

func inputStreamAvailable(params []any) any {
	return int64(0)
}

func inputStreamNullInputStream([]any) any {
	className := "java/io/InputStream"
	return object.MakeEmptyObjectWithClassName(&className)
}

func inputStreamReset(params []any) any {
	return ghelpers.GetGErrBlk(excNames.IOException, "mark/reset not supported")
}

func inputStreamReadByte(params []any) any {
	return int64(-1)
}

func inputStreamInvokeReadRange(self *object.Object, bufObj *object.Object, off int64, lenVal int64) any {
	if lenVal == 0 {
		return int64(0)
	}

	className := getObjectClassName(self)
	if className != "" && className != "java/io/InputStream" {
		method3 := fmt.Sprintf("%s.read([BII)I", className)
		if gmeth, ok := ghelpers.MethodSignatures[method3]; ok && gmeth.GFunction != nil {
			res := ghelpers.Invoke(method3, []any{self, bufObj, off, lenVal})
			if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
				return res
			}
		}

		method1 := fmt.Sprintf("%s.read()I", className)
		if gmeth, ok := ghelpers.MethodSignatures[method1]; ok && gmeth.GFunction != nil {
			res0 := ghelpers.Invoke(method1, []any{self})
			if errBlk, ok := asGErrBlk(res0); ok {
				return errBlk
			}
			b0, ok := res0.(int64)
			if !ok || b0 < 0 {
				return int64(-1)
			}
			buf := bufObj.FieldTable["value"].Fvalue.([]types.JavaByte)
			buf[off] = types.JavaByte(b0)
			var i int64 = 1
			for ; i < lenVal; i++ {
				resI := ghelpers.Invoke(method1, []any{self})
				if _, ok := asGErrBlk(resI); ok {
					break
				}
				bI, ok := resI.(int64)
				if !ok || bI < 0 {
					break
				}
				buf[off+i] = types.JavaByte(bI)
			}
			return i
		}
	}

	return inputStreamReadRange([]any{self, bufObj, off, lenVal})
}

func inputStreamInvokeReadByte(self *object.Object) any {
	className := getObjectClassName(self)
	if className != "" && className != "java/io/InputStream" {
		method1 := fmt.Sprintf("%s.read()I", className)
		if gmeth, ok := ghelpers.MethodSignatures[method1]; ok && gmeth.GFunction != nil {
			res := ghelpers.Invoke(method1, []any{self})
			if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
				return res
			}
		}

		oneByte := make([]types.JavaByte, 1)
		oneObj := object.MakePrimitiveObject("[B", types.JavaByteArray, oneByte)
		res := inputStreamInvokeReadRange(self, oneObj, 0, 1)
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		n, ok := res.(int64)
		if !ok || n <= 0 {
			return int64(-1)
		}
		actual := oneObj.FieldTable["value"].Fvalue.([]types.JavaByte)
		return int64(uint8(actual[0]))
	}

	return int64(-1)
}

func inputStreamInvokeSkip(self *object.Object, n int64) any {
	if n <= 0 {
		return int64(0)
	}

	className := getObjectClassName(self)
	if className != "" && className != "java/io/InputStream" {
		method := fmt.Sprintf("%s.skip(J)J", className)
		if gmeth, ok := ghelpers.MethodSignatures[method]; ok && gmeth.GFunction != nil {
			res := ghelpers.Invoke(method, []any{self, n})
			if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
				return res
			}
		}
	}

	return inputStreamSkip([]any{self, n})
}

func outputStreamInvokeWriteRange(out *object.Object, bufObj *object.Object, off int64, lenVal int64) any {
	if lenVal == 0 {
		return nil
	}

	className := getObjectClassName(out)
	if className != "" {
		method3 := fmt.Sprintf("%s.write([BII)V", className)
		if gmeth, ok := ghelpers.MethodSignatures[method3]; ok && gmeth.GFunction != nil {
			var res any
			if gmeth.NeedsContext {
				res = gmeth.GFunction([]any{list.New(), out, bufObj, off, lenVal})
			} else {
				res = ghelpers.Invoke(method3, []any{out, bufObj, off, lenVal})
			}
			if errBlk, ok := asGErrBlk(res); !ok || errBlk.ExceptionType != excNames.NoSuchMethodException {
				return res
			}
		}

		method1 := fmt.Sprintf("%s.write(I)V", className)
		if gmeth, ok := ghelpers.MethodSignatures[method1]; ok && gmeth.GFunction != nil {
			buf := bufObj.FieldTable["value"].Fvalue.([]types.JavaByte)
			for i := int64(0); i < lenVal; i++ {
				b := int64(uint8(buf[off+i]))
				res := ghelpers.Invoke(method1, []any{out, b})
				if res != nil {
					return res
				}
			}
			return nil
		}
	}

	if gmeth, ok := ghelpers.MethodSignatures["java/io/OutputStream.write([BII)V"]; ok && gmeth.GFunction != nil {
		return gmeth.GFunction([]any{list.New(), out, bufObj, off, lenVal})
	}

	return ghelpers.GetGErrBlk(excNames.IOException, "transferTo: No suitable write method found for OutputStream")
}

// java/io/InputStream.read([BII)I
func inputStreamReadRange(params []any) any {
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "read: null receiver")
	}
	self := params[0].(*object.Object)

	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "read: null buffer")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "read: null buffer")
	}

	var buf []types.JavaByte
	if bufObj.FieldTable != nil {
		if fld, ok := bufObj.FieldTable["value"]; ok && fld.Fvalue != nil {
			if b, ok := fld.Fvalue.([]types.JavaByte); ok {
				buf = b
			} else if bRaw, ok := fld.Fvalue.([]byte); ok {
				buf = object.JavaByteArrayFromGoByteArray(bRaw)
			}
		}
	}

	off := getInt64(params[2])
	lenVal := getInt64(params[3])

	if off < 0 || lenVal < 0 || lenVal > int64(len(buf))-off {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException,
			fmt.Sprintf("read: offset=%d, length=%d, array length=%d", off, lenVal, len(buf)))
	}

	if lenVal == 0 {
		return int64(0)
	}

	res0 := inputStreamInvokeReadByte(self)
	if errBlk, ok := asGErrBlk(res0); ok {
		return errBlk
	}
	b0, ok := res0.(int64)
	if !ok || b0 < 0 {
		return int64(-1)
	}
	buf[off] = types.JavaByte(b0)

	var i int64 = 1
	for ; i < lenVal; i++ {
		resI := inputStreamInvokeReadByte(self)
		if _, ok := asGErrBlk(resI); ok {
			break
		}
		bI, ok := resI.(int64)
		if !ok || bI < 0 {
			break
		}
		buf[off+i] = types.JavaByte(bI)
	}

	return i
}

// java/io/InputStream.read([B)I
func inputStreamReadIntoByteArray(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "read: null receiver")
	}
	self := params[0].(*object.Object)

	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "read: null buffer")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "read: null buffer")
	}

	var buf []types.JavaByte
	if bufObj.FieldTable != nil {
		if fld, ok := bufObj.FieldTable["value"]; ok && fld.Fvalue != nil {
			if b, ok := fld.Fvalue.([]types.JavaByte); ok {
				buf = b
			} else if bRaw, ok := fld.Fvalue.([]byte); ok {
				buf = object.JavaByteArrayFromGoByteArray(bRaw)
			}
		}
	}

	if len(buf) == 0 {
		return int64(0)
	}

	return inputStreamInvokeReadRange(self, bufObj, 0, int64(len(buf)))
}

// java/io/InputStream.skip(J)J
func inputStreamSkip(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "skip: null receiver")
	}
	self := params[0].(*object.Object)

	n := getInt64(params[1])
	if n <= 0 {
		return int64(0)
	}

	var skipped int64 = 0
	bufSize := minInt64(n, 2048)
	tempBuf := make([]types.JavaByte, bufSize)
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, tempBuf)

	for skipped < n {
		toRead := minInt64(n-skipped, bufSize)
		res := inputStreamInvokeReadRange(self, bufObj, 0, toRead)
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		nr, ok := res.(int64)
		if !ok || nr <= 0 {
			break
		}
		skipped += nr
	}

	return skipped
}

// java/io/InputStream.readAllBytes()[B
func inputStreamReadAllBytes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "readAllBytes: null receiver")
	}
	return inputStreamReadNBytesLen([]any{params[0], int64(math.MaxInt32)})
}

// java/io/InputStream.readNBytes([BII)I
func inputStreamReadNBytesBuf(params []any) any {
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "readNBytes: null receiver")
	}
	self := params[0].(*object.Object)

	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "readNBytes: null buffer")
	}
	bufObj, ok := params[1].(*object.Object)
	if !ok || bufObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "readNBytes: null buffer")
	}

	var buf []types.JavaByte
	if bufObj.FieldTable != nil {
		if fld, ok := bufObj.FieldTable["value"]; ok && fld.Fvalue != nil {
			if b, ok := fld.Fvalue.([]types.JavaByte); ok {
				buf = b
			} else if bRaw, ok := fld.Fvalue.([]byte); ok {
				buf = object.JavaByteArrayFromGoByteArray(bRaw)
			}
		}
	}

	off := getInt64(params[2])
	lenVal := getInt64(params[3])

	if off < 0 || lenVal < 0 || lenVal > int64(len(buf))-off {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException,
			fmt.Sprintf("readNBytes: offset=%d, length=%d, array length=%d", off, lenVal, len(buf)))
	}

	if lenVal == 0 {
		return int64(0)
	}

	var n int64 = 0
	for n < lenVal {
		res := inputStreamInvokeReadRange(self, bufObj, off+n, lenVal-n)
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		count, ok := res.(int64)
		if !ok || count < 0 {
			break
		}
		if count == 0 {
			continue
		}
		n += count
	}

	return n
}

// java/io/InputStream.readNBytes(I)[B
func inputStreamReadNBytesLen(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "readNBytes: null receiver")
	}
	self := params[0].(*object.Object)

	lenVal := getInt64(params[1])
	if lenVal < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "len < 0")
	}
	if lenVal == 0 {
		return object.MakePrimitiveObject("[B", types.JavaByteArray, []types.JavaByte{})
	}
	if lenVal > 2147483639 && lenVal != math.MaxInt32 {
		return ghelpers.GetGErrBlk(excNames.OutOfMemoryError, "Required array size too large")
	}

	var totalRead int64 = 0
	capacity := minInt64(lenVal, 8192)
	result := make([]types.JavaByte, 0, capacity)

	for totalRead < lenVal {
		toRead := minInt64(8192, lenVal-totalRead)
		chunk := make([]types.JavaByte, toRead)
		chunkObj := object.MakePrimitiveObject("[B", types.JavaByteArray, chunk)

		res := inputStreamInvokeReadRange(self, chunkObj, 0, toRead)
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		nRead, ok := res.(int64)
		if !ok || nRead < 0 {
			break
		}
		if nRead == 0 {
			continue
		}

		actualBytes := chunkObj.FieldTable["value"].Fvalue.([]types.JavaByte)
		result = append(result, actualBytes[:nRead]...)
		totalRead += nRead
	}

	return object.MakePrimitiveObject("[B", types.JavaByteArray, result)
}

// java/io/InputStream.skipNBytes(J)V
func inputStreamSkipNBytes(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "skipNBytes: null receiver")
	}
	self := params[0].(*object.Object)

	nVal := getInt64(params[1])
	if nVal <= 0 {
		return nil
	}

	for nVal > 0 {
		res := inputStreamInvokeSkip(self, nVal)
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		ns, ok := res.(int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IOException, "Unable to skip")
		}
		if ns > 0 && ns <= nVal {
			nVal -= ns
		} else if ns == 0 {
			resRead := inputStreamInvokeReadByte(self)
			if errBlk, ok := asGErrBlk(resRead); ok {
				return errBlk
			}
			b, ok := resRead.(int64)
			if !ok || b < 0 {
				return ghelpers.GetGErrBlk(excNames.IOException, "skipNBytes: reached end of stream")
			}
			nVal--
		} else if ns > nVal {
			nVal = 0
		} else if ns < 0 {
			return ghelpers.GetGErrBlk(excNames.IOException, "Unable to skip exactly")
		}
	}

	return nil
}

// java/io/InputStream.transferTo(Ljava/io/OutputStream;)J
func inputStreamTransferTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "transferTo: null receiver")
	}
	self := params[0].(*object.Object)

	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "transferTo: out is null")
	}
	out, ok := params[1].(*object.Object)
	if !ok || out == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "transferTo: out is null")
	}

	var transferred int64 = 0
	buf := make([]types.JavaByte, 8192)
	bufObj := object.MakePrimitiveObject("[B", types.JavaByteArray, buf)

	for {
		res := inputStreamInvokeReadRange(self, bufObj, 0, int64(len(buf)))
		if errBlk, ok := asGErrBlk(res); ok {
			return errBlk
		}
		nRead, ok := res.(int64)
		if !ok || nRead < 0 {
			break
		}
		if nRead == 0 {
			continue
		}

		resWrite := outputStreamInvokeWriteRange(out, bufObj, 0, nRead)
		if errBlk, ok := asGErrBlk(resWrite); ok {
			return errBlk
		}
		transferred += nRead
	}

	return transferred
}
