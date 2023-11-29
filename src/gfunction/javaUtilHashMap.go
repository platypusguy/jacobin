/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"crypto/md5"
	"encoding/binary"
	"jacobin/object"
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

	obj := params[0]
	switch obj.(type) {
	case *object.Object:
		o := obj.(*object.Object)              // force golang to treat it as the object we know it to be
		f := o.Fields[0].Fvalue.(*[]byte)      // get the first field
		roughHash := md5.Sum(*f)               // md5.sum returns an array of bytes
		hash := roughHash[:]                   // convert the array to a slice so we can convert to int
		uHash := binary.BigEndian.Uint64(hash) // convert slice of bytes to Uint (int is not available)
		return int64(uHash)                    // convert uint64 to int64
	default:
		panic("unrecognized type to hash in hashMapHash()")
	}
	return hashValue
}
