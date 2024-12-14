/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

/*
 Two types of CRC32 calculation:
 CRC-32            crc32 using an IEEE polynomial
 CRC-32C           crc32 a Castagnoli polynomial
*/

package gfunction

import (
	"hash/crc32"
	"jacobin/object"
	"jacobin/types"
)

func Load_Util_Zip_Crc32_Crc32c() {

	MethodSignatures["java/util/zip/CRC32.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/zip/CRC32C.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/zip/CRC32.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  crc32InitIEEE,
		}

	MethodSignatures["java/util/zip/CRC32C.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  crc32InitCastagnoli,
		}

	MethodSignatures["java/util/zip/CRC32.getValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  crc32GetValue,
		}

	MethodSignatures["java/util/zip/CRC32C.getValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  crc32GetValue,
		}

	MethodSignatures["java/util/zip/CRC32.reset()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  crc32Reset,
		}

	MethodSignatures["java/util/zip/CRC32C.reset()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  crc32Reset,
		}

	MethodSignatures["java/util/zip/CRC32.update([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  crc32UpdateFromArray,
		}

	MethodSignatures["java/util/zip/CRC32C.update([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  crc32UpdateFromArray,
		}

	MethodSignatures["java/util/zip/CRC32.update(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  crc32UpdateFromInt,
		}

	MethodSignatures["java/util/zip/CRC32C.update(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  crc32UpdateFromInt,
		}

	MethodSignatures["java/util/zip/CRC32.update(Ljava/nio/ByteBuffer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/zip/CRC32C.update(Ljava/nio/ByteBuffer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

}

// Instantiate an CRC32 object using an IEEE polynomial.
func crc32InitIEEE(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fld := object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj.FieldTable["value"] = fld
	obj.FieldTable["resetValue"] = fld
	fld = object.Field{Ftype: types.Int, Fvalue: int64(crc32.IEEE)}
	obj.FieldTable["polynomial"] = fld
	return nil
}

// Instantiate an CRC32C object using a Castagnoli polynomial.
func crc32InitCastagnoli(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fld := object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj.FieldTable["value"] = fld
	obj.FieldTable["resetValue"] = fld
	fld = object.Field{Ftype: types.Int, Fvalue: int64(crc32.Castagnoli)}
	obj.FieldTable["polynomial"] = fld
	return nil
}

// Get the current CRC32 value.
func crc32GetValue(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	value := obj.FieldTable["value"].Fvalue.(int64)
	return value
}

// Set the current CRC32 value to the initial (reset) value.
func crc32Reset(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	resetValue := obj.FieldTable["resetValue"].Fvalue.(int64)
	fld := obj.FieldTable["value"]
	fld.Fvalue = resetValue
	obj.FieldTable["value"] = fld

	return nil
}

// Update the current CRC32 value from an array of bytes.
func crc32UpdateFromArray(params []interface{}) interface{} {
	// Collect parameters.
	obj := params[0].(*object.Object)
	objBB := params[1].(*object.Object)
	bbWhole := objBB.FieldTable["value"].Fvalue.([]types.JavaByte)
	offset := params[2].(int64)
	length := params[3].(int64)
	bbSubset := bbWhole[offset:length]

	// Get current CRC32 value.
	fldValue := obj.FieldTable["value"]
	value := fldValue.Fvalue.(int64)
	initialChecksum := uint32(value)

	// Get CRC32 polynomial.
	fldPoly := obj.FieldTable["polynomial"]
	valuePoly := uint32(fldPoly.Fvalue.(int64))

	// Compute new checksum and store it back.
	fldValue.Fvalue = int64(updateCRC32(initialChecksum,
		object.GoByteArrayFromJavaByteArray(bbSubset), valuePoly))
	obj.FieldTable["value"] = fldValue

	return nil
}

// Update the current CRC32 value from a single byte.
func crc32UpdateFromInt(params []interface{}) interface{} {
	// Collect parameters.
	obj := params[0].(*object.Object)
	bb := make([]byte, 1)
	bb[0] = byte(params[1].(int64))

	// Get current CRC32 value.
	fldValue := obj.FieldTable["value"]
	value := fldValue.Fvalue.(int64)
	initialChecksum := uint32(value)

	// Get CRC32 polynomial.
	fldPoly := obj.FieldTable["polynomial"]
	valuePoly := uint32(fldPoly.Fvalue.(int64))

	// Compute new checksum and store it back.
	fldValue.Fvalue = int64(updateCRC32(initialChecksum, bb, valuePoly))
	obj.FieldTable["value"] = fldValue

	return nil
}

// UpdateCRC32 updates an existing CRC32 checksum with new input data
func updateCRC32(existingChecksum uint32, input []byte, polynomial uint32) uint32 {
	// Create a CRC32 table with the given polynomial
	table := crc32.MakeTable(polynomial)

	// Calculate checksum and return it.
	checksum := existingChecksum
	checksum = crc32.Update(checksum, table, input)
	return checksum
}
