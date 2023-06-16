/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package javaTypes

// bytes in Go are uint8, whereas in Java they are int8. Hence this type alias.
type JavaByte = int8

// booleans in Java are defined as integer values of 0 and 1
// in arrays, they're stored as bytes, everywhere else as 32-bit ints.
// Jacobin, however, uses 64-bit ints.

const JavaBoolTrue int64 = 1
const JavaBoolFalse int64 = 0

var JavaBool int64

// ConvertGoBoolToJavaBool takes a go boolean which is not a numeric
// value (and can't be cast to one) and converts into into an integral
// type using the constraints defined in section 2.3.4 of the JVM spec,
// with the notable difference that we're using an int64, rather than
// Java's 32-bit int.
func ConvertGoBoolToJavaBool(goBool bool) int64 {
	if goBool {
		return JavaBoolTrue
	} else {
		return JavaBoolFalse
	}
}
