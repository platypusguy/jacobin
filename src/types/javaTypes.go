/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package types

import "strings"

const Bool = "Z"
const Byte = "B"
const Char = "C"
const Double = "D"
const Float = "F"
const Int = "I" // can be either 32- or 64-bit int
const Long = "J"
const Short = "S"

const Array = "["
const ByteArray = "[B"
const Ref = "Z"

// Jacobin-specific types
const String = "T"
const Static = "X"

const GoMeth = "G" // a go mehod

const Error = "0"  // if an error occurred in getting a type
const Struct = "9" // used primarily in returning items from the CP

func IsIntegral(t string) bool {
	if t == "B" || t == "C" || t == "I" ||
		t == "J" || t == "S" || t == "Z" {
		return true
	}
	return false
}

func IsFloatingPoint(t string) bool {
	if t == "F" || t == "D" {
		return true
	}
	return false
}

func IsAddress(t string) bool {
	if strings.HasPrefix(t, "L") || strings.HasPrefix(t, "[") || t == "T" {
		return true
	}
	return false
}

func IsStatic(t string) bool {
	if strings.HasPrefix(t, "X") {
		return true
	}
	return false
}

func IsError(t string) bool {
	if t == "0" {
		return true
	}
	return false
}

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
