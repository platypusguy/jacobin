/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Io_ByteArrayOutputStream() {

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayOutputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayOutputStreamInit,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayOutputStreamReset,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayOutputStreamSize,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.toByteArray()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayOutputStreamToByteArray,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ByteArrayOutputStreamToString,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.toString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.toString(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayOutputStreamToStringCharsetName,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.toString(Ljava/nio/charset/Charset;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayOutputStreamToStringCharset,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayOutputStreamWriteInt,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.write([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ByteArrayOutputStreamWriteBytes,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.writeBytes([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayOutputStreamWriteBytesAll,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayOutputStream.writeTo(Ljava/io/OutputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ByteArrayOutputStreamWriteTo,
		}
}

// ByteArrayOutputStreamInit implements <init>()V and <init>(I)V
func ByteArrayOutputStreamInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)

	// Ensure the target object's FieldTable map is initialized before assignment.
	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}

	size := int64(32) // Default size
	if len(params) == 2 {
		size = params[1].(int64)
	}

	buf := make([]types.JavaByte, size)
	self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: buf}
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: int64(0)}

	return nil
}

func ensureCapacity(self *object.Object, minCapacity int64) {
	bufField := self.FieldTable["buf"]
	buf := bufField.Fvalue.([]types.JavaByte)
	if minCapacity > int64(len(buf)) {
		newCapacity := int64(len(buf)) * 2
		if newCapacity < minCapacity {
			newCapacity = minCapacity
		}
		newBuf := make([]types.JavaByte, newCapacity)
		copy(newBuf, buf)
		self.FieldTable["buf"] = object.Field{Ftype: "[B", Fvalue: newBuf}
	}
}

func ByteArrayOutputStreamWriteInt(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	b := params[1].(int64)

	count := self.FieldTable["count"].Fvalue.(int64)
	ensureCapacity(self, count+1)

	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)
	buf[count] = types.JavaByte(b)
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: count + 1}

	return nil
}

func ByteArrayOutputStreamWriteBytes(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bObj, _ := params[1].(*object.Object)
	off := params[2].(int64)
	lenVal := params[3].(int64)

	if bObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "write: byte array is null")
	}

	b := bObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if off < 0 || lenVal < 0 || off+lenVal > int64(len(b)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "write: index out of bounds")
	}

	count := self.FieldTable["count"].Fvalue.(int64)
	ensureCapacity(self, count+lenVal)

	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)
	copy(buf[count:], b[off:off+lenVal])
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: count + lenVal}

	return nil
}

func ByteArrayOutputStreamWriteBytesAll(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bObj, _ := params[1].(*object.Object)

	if bObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "writeBytes: byte array is null")
	}

	b := bObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	lenVal := int64(len(b))

	count := self.FieldTable["count"].Fvalue.(int64)
	ensureCapacity(self, count+lenVal)

	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)
	copy(buf[count:], b)
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: count + lenVal}

	return nil
}

func ByteArrayOutputStreamReset(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	self.FieldTable["count"] = object.Field{Ftype: "I", Fvalue: int64(0)}
	return nil
}

func ByteArrayOutputStreamSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["count"].Fvalue.(int64)
}

func ByteArrayOutputStreamToByteArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	count := self.FieldTable["count"].Fvalue.(int64)
	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)

	newBuf := make([]types.JavaByte, count)
	copy(newBuf, buf[:count])

	return object.MakePrimitiveObject("[B", types.ByteArray, newBuf)
}

func ByteArrayOutputStreamToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	count := self.FieldTable["count"].Fvalue.(int64)
	buf := self.FieldTable["buf"].Fvalue.([]types.JavaByte)

	return object.StringObjectFromJavaByteArray(buf[:count])
}

func ByteArrayOutputStreamToStringCharsetName(params []interface{}) interface{} {
	// For now, we don't have full charset support in gfunctions easily accessible
	// without more complex logic. Let's trap it as a reminder or implement if simple.
	// Actually, most Jacobin strings are UTF-8.
	return ghelpers.TrapFunction(params)
}

func ByteArrayOutputStreamToStringCharset(params []interface{}) interface{} {
	return ghelpers.TrapFunction(params)
}

func ByteArrayOutputStreamWriteTo(params []interface{}) interface{} {
	outObj, _ := params[1].(*object.Object)

	if outObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "writeTo: output stream is null")
	}

	// We need to call out.write(buf, 0, count)
	// Since out is a Java OutputStream, we need to invoke it.
	// However, if it's a GFunction-based stream, we might be able to call it directly.
	// For now, let's trap it until we have a clear way to invoke Java methods from GFunctions.
	// Actually, some streams might have a Go-side handle.

	return ghelpers.TrapFunction(params)
}
