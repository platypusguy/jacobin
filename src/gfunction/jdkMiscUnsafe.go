/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/exceptions"
	"jacobin/object"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

func Load_Misc_Unsafe() map[string]GMeth {

	MethodSignatures["jdk/internal/misc/Unsafe.arrayBaseOffset(Ljava/lang/Class;)I"] = // offset to start of first item in an array
		GMeth{
			ParamSlots: 1,
			GFunction:  arrayBaseOffset,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.<clinit>()V"] = // offset to start of first item in an array
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn, // Unsafe <clinit>
		}

	return MethodSignatures
}

// Return the number of bytes between the beginning of the object and the first element.
// This is used in computing the pointer to a given element
func arrayBaseOffset(param []interface{}) interface{} {
	p := param[0]
	if p == nil || p == object.Null {
		errMsg := "jdk.internal.misc.Unsafe::arrayBaseOffset() was passed a null pointer"
		return getGErrBlk(exceptions.NullPointerException, errMsg)
	}
	return int64(0) // this should work...
}
