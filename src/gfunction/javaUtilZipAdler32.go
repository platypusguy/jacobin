/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"jacobin/types"
)

func Load_Util_Zip_Adler32() {

	MethodSignatures["java/util/zip/Adler32.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/zip/Adler32.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  adlerInit,
		}

	MethodSignatures["java/util/zip/Adler32.getValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  adlerGetValue,
		}

	MethodSignatures["java/util/zip/Adler32.reset()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  adlerReset,
		}

	MethodSignatures["java/util/zip/Adler32.update([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  adlerUpdateFromArray,
		}

	MethodSignatures["java/util/zip/Adler32.update(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  adlerUpdateFromInt,
		}

	MethodSignatures["java/util/zip/Adler32.update(Ljava/nio/ByteBuffer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

}

// Instantiate an Adler32 object.
func adlerInit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fld := object.Field{Ftype: types.Long, Fvalue: int64(1)}
	obj.FieldTable["value"] = fld
	obj.FieldTable["resetValue"] = fld
	return nil
}

// Get the current Adler32 value.
func adlerGetValue(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	value := obj.FieldTable["value"].Fvalue.(int64)
	return value
}

// Set the current Adler32 value to the initial (reset) value.
func adlerReset(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	resetValue := obj.FieldTable["resetValue"].Fvalue.(int64)
	fld := obj.FieldTable["value"]
	fld.Fvalue = resetValue
	obj.FieldTable["value"] = fld

	return nil
}

// Update the current Adler32 value from an array of bytes.
func adlerUpdateFromArray(params []interface{}) interface{} {
	// Collect parameters.
	obj := params[0].(*object.Object)
	objBB := params[1].(*object.Object)
	bbWhole := objBB.FieldTable["value"].Fvalue.([]byte)
	offset := params[2].(int64)
	length := params[3].(int64)
	bbSubset := bbWhole[offset:length]

	// Get current Adler32 value.
	fld := obj.FieldTable["value"]
	value := fld.Fvalue.(int64)
	initialChecksum := uint32(value)

	// Compute new checksum and store it back.
	fld.Fvalue = int64(updateAdler32(initialChecksum, bbSubset))
	obj.FieldTable["value"] = fld

	return nil
}

// Update the current Adler32 value from a single byte.
func adlerUpdateFromInt(params []interface{}) interface{} {
	// Collect parameters.
	obj := params[0].(*object.Object)
	bb := make([]byte, 1)
	bb[0] = byte(params[1].(int64))

	// Get current Adler32 value.
	fld := obj.FieldTable["value"]
	value := fld.Fvalue.(int64)
	initialChecksum := uint32(value)

	// Compute new checksum and store it back.
	fld.Fvalue = int64(updateAdler32(initialChecksum, bb))
	obj.FieldTable["value"] = fld

	return nil
}

// UpdateAdler32 updates an existing Adler32 checksum with new input data.
func updateAdler32(existingChecksum uint32, input []byte) uint32 {
	// Adler-32 is split into two parts: A and B
	// Existing checksum = (B << 16) | A
	A := existingChecksum & 0xffff
	B := (existingChecksum >> 16) & 0xffff

	// Modulo used in Adler-32 checksum
	modulo := uint32(65521)

	// Process each byte of the new data
	for _, nextByte := range input {
		A = (A + uint32(nextByte)) % modulo
		B = (B + A) % modulo
	}

	// Combine A and B to get the new checksum
	return (B << 16) | A
}
