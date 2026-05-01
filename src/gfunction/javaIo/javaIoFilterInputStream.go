/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"fmt"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
)

func Load_Io_FilterInputStream() {

	ghelpers.MethodSignatures["java/io/FilterInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.<init>(Ljava/io/InputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFilterInputStream,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filterInputStreamAvailable,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filterInputStreamClose,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  filterInputStreamMark,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filterInputStreamMarkSupported,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filterInputStreamRead,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  filterInputStreamReadByteArray,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  filterInputStreamReadByteArrayOffset,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  filterInputStreamReset,
		}

	ghelpers.MethodSignatures["java/io/FilterInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  filterInputStreamSkip,
		}

}

func initFilterInputStream(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := params[1].(*object.Object)
	self.FieldTable["in"] = object.Field{Ftype: "Ljava/io/InputStream;", Fvalue: in}
	return nil
}

func filterInputStreamAvailable(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.available()I", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in})
}

func filterInputStreamClose(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.close()V", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in})
}

func filterInputStreamMark(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	readlimit := params[1].(int64)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.mark(I)V", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in, readlimit})
}

func filterInputStreamMarkSupported(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.markSupported()Z", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in})
}

func filterInputStreamRead(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.read()I", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in})
}

func filterInputStreamReadByteArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	buf := params[1].(*object.Object)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.read([B)I", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in, buf})
}

func filterInputStreamReadByteArrayOffset(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	buf := params[1].(*object.Object)
	off := params[2].(int64)
	lenVal := params[3].(int64)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.read([BII)I", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in, buf, off, lenVal})
}

func filterInputStreamReset(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.reset()V", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in})
}

func filterInputStreamSkip(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	n := params[1].(int64)
	in := self.FieldTable["in"].Fvalue.(*object.Object)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	method := fmt.Sprintf("%s.skip(J)J", *inClassName)
	return ghelpers.Invoke(method, []interface{}{in, n})
}
