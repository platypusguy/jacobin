/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Io_ByteArrayInputStream() {
	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.<init>([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayInputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.<init>([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ByteArrayInputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayInputStreamAvailable,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayInputStreamMark,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnTrue,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayInputStreamRead,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ByteArrayInputStreamReadRange,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.readAllBytes()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayInputStreamReadAllBytes,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.readNBytes([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ByteArrayInputStreamReadNBytes,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayInputStreamReset,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayInputStreamSkip,
		}
}

func ByteArrayInputStreamInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}

	bufObj := params[1].(*object.Object)
	buf := bufObj.FieldTable["value"].Fvalue.([]types.JavaByte)

	var offset, length int64
	if len(params) == 2 {
		offset = 0
		length = int64(len(buf))
	} else {
		offset = params[2].(int64)
		length = params[3].(int64)
	}

	self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: buf}
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: offset}
	self.FieldTable["mark"] = object.Field{Ftype: "I", Fvalue: offset}

	count := offset + length
	if count > int64(len(buf)) {
		count = int64(len(buf))
	}
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: count}

	return nil
}

func ByteArrayInputStreamRead(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	pos := self.FieldTable["pos"].Fvalue.(int64)
	count := self.FieldTable["count"].Fvalue.(int64)
	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)

	if pos < count {
		res := int64(buf[pos]) & 0xff
		self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: pos + 1}
		return res
	}
	return int64(-1)
}

func ByteArrayInputStreamReadRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	if object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ByteArrayInputStream.read: b is null")
	}

	destObj := params[1].(*object.Object)
	off := params[2].(int64)
	lenReq := params[3].(int64)

	destBuf := destObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if off < 0 || lenReq < 0 || lenReq > int64(len(destBuf))-off {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "ByteArrayInputStream.read: bounds check failed")
	}

	pos := self.FieldTable["pos"].Fvalue.(int64)
	count := self.FieldTable["count"].Fvalue.(int64)
	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)

	if pos >= count {
		return int64(-1)
	}

	avail := count - pos
	if lenReq > avail {
		lenReq = avail
	}
	if lenReq <= 0 {
		return int64(0)
	}

	copy(destBuf[off:], buf[pos:pos+lenReq])
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: pos + lenReq}

	return lenReq
}

func ByteArrayInputStreamAvailable(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	pos := self.FieldTable["pos"].Fvalue.(int64)
	count := self.FieldTable["count"].Fvalue.(int64)

	return count - pos
}

func ByteArrayInputStreamSkip(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	n := params[1].(int64)

	pos := self.FieldTable["pos"].Fvalue.(int64)
	count := self.FieldTable["count"].Fvalue.(int64)

	k := count - pos
	if n < k {
		if n < 0 {
			k = 0
		} else {
			k = n
		}
	}

	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: pos + k}
	return k
}

func ByteArrayInputStreamMark(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	// readAheadLimit is ignored by ByteArrayInputStream
	self.FieldTable["mark"] = object.Field{Ftype: "I", Fvalue: self.FieldTable["pos"].Fvalue}
	return nil
}

func ByteArrayInputStreamReset(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: self.FieldTable["mark"].Fvalue}
	return nil
}

func ByteArrayInputStreamReadAllBytes(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	pos := self.FieldTable["pos"].Fvalue.(int64)
	count := self.FieldTable["count"].Fvalue.(int64)
	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)

	avail := count - pos
	if avail <= 0 {
		// Return empty byte array
		res, _ := globals.GetGlobalRef().FuncInstantiateClass("[B", nil)
		resObj := res.(*object.Object)
		resObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: []types.JavaByte{}}
		return resObj
	}

	data := make([]types.JavaByte, avail)
	copy(data, buf[pos:count])
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: count}

	res, _ := globals.GetGlobalRef().FuncInstantiateClass("[B", nil)
	resObj := res.(*object.Object)
	resObj.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: data}
	return resObj
}

func ByteArrayInputStreamReadNBytes(params []interface{}) interface{} {
	// readNBytes([BII)I in ByteArrayInputStream just calls read([BII)I
	return ByteArrayInputStreamReadRange(params)
}
