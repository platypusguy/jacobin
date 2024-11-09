/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/types"
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

func Load_Lang_UTF16() {

	MethodSignatures["java/lang/StringUTF16.isBigEndian()Z"] = // is the current platform big-endian?
		GMeth{
			ParamSlots: 0,
			GFunction:  isBigEndian,
		}

}

// "java/lang/StringUTF16.isBigEndian()Z"
// Return whether the current platform is big-endian.
func isBigEndian([]interface{}) interface{} {
	return types.JavaBoolFalse // Jacobin runs on no big-endian platform at present
}
