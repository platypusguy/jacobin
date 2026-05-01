/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

func Load_Util_Zip_CheckedInputStream() {
	ghelpers.MethodSignatures["java/util/zip/CheckedInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/zip/CheckedInputStream.<init>(Ljava/io/InputStream;Ljava/util/zip/Checksum;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  CheckedInputStreamInit,
		}

	ghelpers.MethodSignatures["java/util/zip/CheckedInputStream.getChecksum()Ljava/util/zip/Checksum;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  CheckedInputStreamGetChecksum,
		}

	ghelpers.MethodSignatures["java/util/zip/CheckedInputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  CheckedInputStreamRead,
		}

	ghelpers.MethodSignatures["java/util/zip/CheckedInputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  CheckedInputStreamReadArray,
		}

	ghelpers.MethodSignatures["java/util/zip/CheckedInputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  CheckedInputStreamSkip,
		}
}

func CheckedInputStreamInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	inputStream := params[1].(*object.Object)
	cksum := params[2].(*object.Object)

	self.FieldTable["in"] = object.Field{Ftype: "Ljava/io/InputStream;", Fvalue: inputStream}
	self.FieldTable["cksum"] = object.Field{Ftype: "Ljava/util/zip/Checksum;", Fvalue: cksum}

	return nil
}

func CheckedInputStreamGetChecksum(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	cksum, ok := self.FieldTable["cksum"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "CheckedInputStream.cksum is null")
	}
	return cksum
}

func CheckedInputStreamRead(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	in, ok := self.FieldTable["in"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "CheckedInputStream.in is null")
	}

	// Read one byte from 'in'
	inClassName := stringPool.GetStringPointer(in.KlassName)
	readMethod := fmt.Sprintf("%s.read()I", *inClassName)
	res := ghelpers.Invoke(readMethod, []interface{}{in})

	if err, ok := res.(*ghelpers.GErrBlk); ok {
		return err
	}

	val := res.(int64)
	if val != -1 {
		cksum := self.FieldTable["cksum"].Fvalue.(*object.Object)
		// Update checksum with the byte
		cksumClassName := stringPool.GetStringPointer(cksum.KlassName)
		updateMethod := fmt.Sprintf("%s.update(I)V", *cksumClassName)
		ghelpers.Invoke(updateMethod, []interface{}{cksum, val})
	}

	return val
}

func CheckedInputStreamReadArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	buf := params[1].(*object.Object)
	off := params[2].(int64)
	lenVal := params[3].(int64)

	in, ok := self.FieldTable["in"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "CheckedInputStream.in is null")
	}

	// Call in.read(buf, off, len)
	inClassName := stringPool.GetStringPointer(in.KlassName)
	readMethod := fmt.Sprintf("%s.read([BII)I", *inClassName)
	res := ghelpers.Invoke(readMethod, []interface{}{in, buf, off, lenVal})

	if err, ok := res.(*ghelpers.GErrBlk); ok {
		return err
	}

	n := res.(int64)
	if n > 0 {
		cksum := self.FieldTable["cksum"].Fvalue.(*object.Object)
		// Update checksum with the bytes read
		cksumClassName := stringPool.GetStringPointer(cksum.KlassName)
		updateMethod := fmt.Sprintf("%s.update([BII)V", *cksumClassName)
		ghelpers.Invoke(updateMethod, []interface{}{cksum, buf, off, n})
	}

	return n
}

func CheckedInputStreamSkip(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	n := params[1].(int64)

	if n <= 0 {
		return int64(0)
	}

	// CheckedInputStream.skip updates checksum as well.
	// It reads into a temporary buffer and updates.
	// We'll use a 512 byte buffer.
	bufSize := int64(512)
	if n < bufSize {
		bufSize = n
	}

	// Create a byte array object for reading
	javaBytes := make([]types.JavaByte, bufSize)
	bufObj := object.MakePrimitiveObject("[B", types.ByteArray, javaBytes) // Using "[B" for byte array class name

	totalSkipped := int64(0)
	for totalSkipped < n {
		toRead := n - totalSkipped
		if toRead > bufSize {
			toRead = bufSize
		}

		res := CheckedInputStreamReadArray([]interface{}{self, bufObj, int64(0), toRead})
		if err, ok := res.(*ghelpers.GErrBlk); ok {
			return err
		}

		nr := res.(int64)
		if nr < 0 {
			break
		}
		totalSkipped += nr
	}

	return totalSkipped
}
