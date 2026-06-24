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
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

const defaultBufferSize = 8192

func Load_Io_BufferedInputStream() {
	ghelpers.MethodSignatures["java/io/BufferedInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.<init>(Ljava/io/InputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  BufferedInputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.<init>(Ljava/io/InputStream;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  BufferedInputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  BufferedInputStreamAvailable,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  BufferedInputStreamClose,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  BufferedInputStreamMark,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnTrue,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  BufferedInputStreamRead,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  BufferedInputStreamReadRange,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  BufferedInputStreamReset,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  BufferedInputStreamSkip,
		}
}

func BufferedInputStreamInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := params[1].(*object.Object)

	var size int64 = defaultBufferSize
	if len(params) == 3 {
		size = params[2].(int64)
		if size <= 0 {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Buffer size <= 0")
		}
	}

	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}

	self.FieldTable["in"] = object.Field{Ftype: "Ljava/io/InputStream;", Fvalue: in}

	buf := make([]types.JavaByte, size)
	self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: buf}
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: int64(0)}
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: int64(0)}
	self.FieldTable["markpos"] = object.Field{Ftype: "I", Fvalue: int64(-1)}
	self.FieldTable["marklimit"] = object.Field{Ftype: "I", Fvalue: int64(0)}

	return nil
}

func getInIfOpen(self *object.Object) (*object.Object, interface{}) {
	inField, ok := self.FieldTable["in"]
	if !ok || object.IsNull(inField.Fvalue) {
		return nil, ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}
	return inField.Fvalue.(*object.Object), nil
}

func getBufIfOpen(self *object.Object) ([]types.JavaByte, interface{}) {
	bufField, ok := self.FieldTable["buf"]
	if !ok || object.IsNull(bufField.Fvalue) {
		return nil, ghelpers.GetGErrBlk(excNames.IOException, "Stream closed")
	}
	return bufField.Fvalue.([]types.JavaByte), nil
}

func fill(self *object.Object) interface{} {
	buf, err := getBufIfOpen(self)
	if err != nil {
		return err
	}

	markpos := self.FieldTable["markpos"].Fvalue.(int64)
	pos := self.FieldTable["pos"].Fvalue.(int64)
	marklimit := self.FieldTable["marklimit"].Fvalue.(int64)

	if markpos < 0 {
		pos = 0
	} else if pos >= int64(len(buf)) {
		if markpos > 0 {
			sz := pos - markpos
			copy(buf, buf[markpos:pos])
			pos = sz
			markpos = 0
		} else if int64(len(buf)) >= marklimit {
			markpos = -1
			pos = 0
		} else {
			// Grow buffer
			ns := pos * 2
			if ns > marklimit {
				ns = marklimit
			}
			nbuf := make([]types.JavaByte, ns)
			copy(nbuf, buf)
			buf = nbuf
			self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: buf}
		}
	}

	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: pos}
	self.FieldTable["markpos"] = object.Field{Ftype: "I", Fvalue: markpos}

	in, err := getInIfOpen(self)
	if err != nil {
		return err
	}

	// Read into buf from pos
	inClassName := stringPool.GetStringPointer(in.KlassName)
	// We need to call read([BII)I on 'in'
	method := fmt.Sprintf("%s.read([BII)I", *inClassName)

	// Prepare [B object for buf
	// Since we are operating on the internal byte slice, we might need a temporary object wrapper
	// if ghelpers.Invoke expects it.
	// Actually, most read implementations in Jacobin expect an *object.Object representing [B.

	// Create a temporary [B object wrapper
	tempBufObj, _ := globals.GetGlobalRef().FuncInstantiateClass("[B", nil)
	bufWrapper := tempBufObj.(*object.Object)
	bufWrapper.FieldTable["value"] = object.Field{Ftype: "[B", Fvalue: buf}

	n := ghelpers.Invoke(method, []interface{}{in, bufWrapper, pos, int64(len(buf)) - pos})

	if errBlk, ok := n.(ghelpers.GErrBlk); ok {
		return errBlk
	}

	nRead := n.(int64)
	if nRead > 0 {
		self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: pos + nRead}
	} else {
		self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: pos}
	}

	return nil
}

func BufferedInputStreamRead(params []interface{}) interface{} {
	self := params[0].(*object.Object)

	pos := self.FieldTable["pos"].Fvalue.(int64)
	count := self.FieldTable["count"].Fvalue.(int64)

	if pos >= count {
		res := fill(self)
		if err, ok := res.(ghelpers.GErrBlk); ok {
			return err
		}
		pos = self.FieldTable["pos"].Fvalue.(int64)
		count = self.FieldTable["count"].Fvalue.(int64)
		if pos >= count {
			return int64(-1)
		}
	}

	buf, err := getBufIfOpen(self)
	if err != nil {
		return err
	}

	b := int64(buf[pos]) & 0xff
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: pos + 1}
	return b
}

func BufferedInputStreamReadRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	if object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "b is null")
	}
	bObj := params[1].(*object.Object)
	off := params[2].(int64)
	length := params[3].(int64)

	if off < 0 || length < 0 || off+length > int64(len(bObj.FieldTable["value"].Fvalue.([]types.JavaByte))) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "Read bounds check failed")
	}

	if length == 0 {
		return int64(0)
	}

	totalRead := int64(0)
	for {
		avail := self.FieldTable["count"].Fvalue.(int64) - self.FieldTable["pos"].Fvalue.(int64)
		if avail <= 0 {
			res := fill(self)
			if err, ok := res.(ghelpers.GErrBlk); ok {
				return err
			}
			avail = self.FieldTable["count"].Fvalue.(int64) - self.FieldTable["pos"].Fvalue.(int64)
			if avail <= 0 {
				if totalRead == 0 {
					return int64(-1)
				}
				return totalRead
			}
		}

		cnt := avail
		if cnt > length {
			cnt = length
		}

		buf, err := getBufIfOpen(self)
		if err != nil {
			return err
		}
		pos := self.FieldTable["pos"].Fvalue.(int64)
		destBuf := bObj.FieldTable["value"].Fvalue.([]types.JavaByte)
		copy(destBuf[off:], buf[pos:pos+cnt])

		self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: pos + cnt}
		totalRead += cnt
		length -= cnt
		off += cnt

		if length <= 0 {
			return totalRead
		}

		// If input stream has no more data immediately available, return what we have
		in, err := getInIfOpen(self)
		if err != nil {
			return err
		}
		inClassName := stringPool.GetStringPointer(in.KlassName)
		method := fmt.Sprintf("%s.available()I", *inClassName)
		inAvail := ghelpers.Invoke(method, []interface{}{in}).(int64)
		if inAvail <= 0 {
			return totalRead
		}
	}
}

func BufferedInputStreamAvailable(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in, err := getInIfOpen(self)
	if err != nil {
		return err
	}

	avail := self.FieldTable["count"].Fvalue.(int64) - self.FieldTable["pos"].Fvalue.(int64)

	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.available()I", *inClassName)
	inAvail := ghelpers.Invoke(method, []interface{}{in}).(int64)

	return avail + inAvail
}

func BufferedInputStreamSkip(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	n := params[1].(int64)
	if n <= 0 {
		return int64(0)
	}

	avail := self.FieldTable["count"].Fvalue.(int64) - self.FieldTable["pos"].Fvalue.(int64)
	if avail <= 0 {
		in, err := getInIfOpen(self)
		if err != nil {
			return err
		}

		markpos := self.FieldTable["markpos"].Fvalue.(int64)
		if markpos < 0 {
			inClassName := stringPool.GetStringPointer(in.KlassName)
			method := fmt.Sprintf("%s.skip(J)J", *inClassName)
			return ghelpers.Invoke(method, []interface{}{in, n})
		}

		res := fill(self)
		if err, ok := res.(ghelpers.GErrBlk); ok {
			return err
		}
		avail = self.FieldTable["count"].Fvalue.(int64) - self.FieldTable["pos"].Fvalue.(int64)
		if avail <= 0 {
			return int64(0)
		}
	}

	skipped := avail
	if skipped > n {
		skipped = n
	}
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: self.FieldTable["pos"].Fvalue.(int64) + skipped}
	return skipped
}

func BufferedInputStreamMark(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	readlimit := params[1].(int64)
	self.FieldTable["marklimit"] = object.Field{Ftype: "I", Fvalue: readlimit}
	self.FieldTable["markpos"] = object.Field{Ftype: "I", Fvalue: self.FieldTable["pos"].Fvalue}
	return nil
}

func BufferedInputStreamReset(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	_, err := getBufIfOpen(self)
	if err != nil {
		return err
	}
	markpos := self.FieldTable["markpos"].Fvalue.(int64)
	if markpos < 0 {
		return ghelpers.GetGErrBlk(excNames.IOException, "Resetting to invalid mark")
	}
	self.FieldTable["pos"] = object.Field{Ftype: "I", Fvalue: markpos}
	return nil
}

func BufferedInputStreamClose(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	buf, _ := self.FieldTable["buf"]
	if !object.IsNull(buf.Fvalue) {
		self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: nil}
		inField, ok := self.FieldTable["in"]
		if ok && !object.IsNull(inField.Fvalue) {
			in := inField.Fvalue.(*object.Object)
			inClassName := stringPool.GetStringPointer(in.KlassName)
			method := fmt.Sprintf("%s.close()V", *inClassName)
			ghelpers.Invoke(method, []interface{}{in})
			self.FieldTable["in"] = object.Field{Ftype: "Ljava/io/InputStream;", Fvalue: nil}
		}
	}
	return nil
}
