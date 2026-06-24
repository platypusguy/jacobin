/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

const (
	defaultOutputBufferSize = 8192
)

func Load_Io_BufferedOutputStream() {
	ghelpers.MethodSignatures["java/io/BufferedOutputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/BufferedOutputStream.<init>(Ljava/io/OutputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  BufferedOutputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/BufferedOutputStream.<init>(Ljava/io/OutputStream;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  BufferedOutputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/BufferedOutputStream.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  BufferedOutputStreamWriteInt,
		}

	ghelpers.MethodSignatures["java/io/BufferedOutputStream.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  BufferedOutputStreamWriteRange,
		}

	ghelpers.MethodSignatures["java/io/BufferedOutputStream.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  BufferedOutputStreamFlush,
		}
}

func BufferedOutputStreamInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	out := params[1].(*object.Object)

	size := int64(defaultOutputBufferSize)
	if len(params) > 2 {
		size = params[2].(int64)
		if size <= 0 {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Buffer size <= 0")
		}
	}

	self.FieldTable["out"] = object.Field{Ftype: "Ljava/io/OutputStream;", Fvalue: out}
	self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: make([]types.JavaByte, size)}
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: int64(0)}

	return nil
}

func flushBuffer(self *object.Object) interface{} {
	count := self.FieldTable["count"].Fvalue.(int64)
	if count > 0 {
		out := self.FieldTable["out"].Fvalue.(*object.Object)
		bufObj := &object.Object{
			FieldTable: map[string]object.Field{
				"value": self.FieldTable["buf"],
			},
		}

		outClassName := stringPool.GetStringPointer(out.KlassName)
		method := fmt.Sprintf("%s.write([BII)V", *outClassName)
		res := ghelpers.Invoke(method, []interface{}{out, bufObj, int64(0), count})
		if res != nil {
			return res
		}
		self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: int64(0)}
	}
	return nil
}

func BufferedOutputStreamWriteInt(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	b := params[1].(int64)

	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)
	count := self.FieldTable["count"].Fvalue.(int64)

	if count >= int64(len(buf)) {
		if res := flushBuffer(self); res != nil {
			return res
		}
		count = 0
	}

	buf[count] = types.JavaByte(b)
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: count + 1}

	return nil
}

func BufferedOutputStreamWriteRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bObj := params[1].(*object.Object)
	off := params[2].(int64)
	lenVal := params[3].(int64)

	if lenVal == 0 {
		return nil
	}

	javaBytes := bObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if off < 0 || lenVal < 0 || off+lenVal > int64(len(javaBytes)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "")
	}

	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)
	count := self.FieldTable["count"].Fvalue.(int64)

	if lenVal >= int64(len(buf)) {
		/* If the request length exceeds the size of the output buffer,
		   flush the output buffer and then write the data directly.
		   In this way buffered streams will cascade harmlessly. */
		if res := flushBuffer(self); res != nil {
			return res
		}
		out := self.FieldTable["out"].Fvalue.(*object.Object)
		outClassName := stringPool.GetStringPointer(out.KlassName)
		method := fmt.Sprintf("%s.write([BII)V", *outClassName)
		return ghelpers.Invoke(method, []interface{}{out, bObj, off, lenVal})
	}

	if lenVal > int64(len(buf))-count {
		if res := flushBuffer(self); res != nil {
			return res
		}
		count = 0
	}

	copy(buf[count:], javaBytes[off:off+lenVal])
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: count + lenVal}

	return nil
}

func BufferedOutputStreamFlush(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	if res := flushBuffer(self); res != nil {
		return res
	}
	out := self.FieldTable["out"].Fvalue.(*object.Object)
	outClassName := stringPool.GetStringPointer(out.KlassName)
	method := fmt.Sprintf("%s.flush()V", *outClassName)

	if _, ok := ghelpers.MethodSignatures[method]; !ok {
		// Fallback to FilterOutputStream.flush() which is often inherited
		method = "java/io/FilterOutputStream.flush()V"
	}

	return ghelpers.Invoke(method, []interface{}{out})
}
