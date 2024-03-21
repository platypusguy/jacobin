/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
)

// Implementation of some of the functions in in Java/lang/Class.

func Load_Util_HashMap() map[string]GMeth {

	MethodSignatures["java/util/HashMap.hash(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hashMapHash,
		}

	return MethodSignatures
}

// hashMapHash accepts a pointer to an object and returns
// a uint64 MD5 hash value of the pointed-to thing
func hashMapHash(params []interface{}) interface{} {
	var hashValue uint64 = 0
	var bytes []byte
	switch params[0].(type) {
	case *object.Object:
		obj := params[0].(*object.Object) // force golang to treat it as the object we know it to be
		fld := obj.FieldTable["value"]
		switch fld.Ftype {
		case types.ByteArray:
			bytes = obj.FieldTable["value"].Fvalue.([]byte)
		case types.Bool, types.Byte, types.Char, types.Int, types.Long, types.Short:
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, uint64(fld.Fvalue.(int64)))
		case types.Double, types.Float:
			bytes := make([]byte, 8)
			binary.BigEndian.PutUint64(bytes, uint64(fld.Fvalue.(float64)))
		default:
			str := fmt.Sprintf("hashMapHash: unrecognized object field type: %T", fld.Ftype)
			return getGErrBlk(exceptions.VirtualMachineError, str)
		}
		roughHash := md5.Sum(bytes)            // md5.sum returns an array of bytes
		hash := roughHash[:]                   // convert the array to a slice
		uHash := binary.BigEndian.Uint64(hash) // convert slice to a uint64
		return int64(uHash)                    // return an int64
	default:
		str := fmt.Sprintf("hashMapHash: unrecognized parameter type: %T", params[0])
		return getGErrBlk(exceptions.VirtualMachineError, str)
	}
	return hashValue
}
