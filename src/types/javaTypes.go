/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package types

import "strings"

// These are the types that are used in an object's field.Ftype and elsewhere.

const Bool = "Z"
const Byte = "B"
const Char = "C"
const Double = "D"
const Float = "F"
const Int = "I" // can be either 32- or 64-bit int
const Long = "J"
const Ref = "L"
const Rune = "R"
const Short = "S"

const Array = "["
const BoolArray = "[Z"
const ByteArray = "[B"
const CharArray = "[R" // char array is = rune array
const DoubleArray = "[D"
const IntArray = "[I"
const FloatArray = "[F"
const LongArray = "[J"
const RefArray = "[L"
const RuneArray = CharArray // used only in strings that are not compact
const ShortArray = "[S"

// Jacobin-specific types
const Interface = "i"
const NonArrayObject = "n"

const StringIndex = "T"
const GolangString = "G"
const FileHandle = "FH" // The related Fvalue is a Golang *os.File
const BigInteger = "BI" // The related Fvalue is a Golang *big.Int

const Static = "X"
const StaticDouble = "XD"
const StaticLong = "XJ"

const GoMeth = "G" // a go mehod

const Error = "0"  // if an error occurred in getting a type
const Struct = "9" // used primarily in returning items from the CP

func IsIntegral(t string) bool {
	if t == Byte || t == Char || t == Int ||
		t == Long || t == Short || t == Bool {
		return true
	}
	return false
}

func IsFloatingPoint(t string) bool {
	if t == Float || t == Double {
		return true
	}
	return false
}

func IsAddress(t string) bool {
	if strings.HasPrefix(t, Ref) || strings.HasPrefix(t, Array) {
		return true
	}
	return false
}

func IsStatic(t string) bool {
	if strings.HasPrefix(t, Static) {
		return true
	}
	return false
}

func IsError(t string) bool {
	if t == Error {
		return true
	}
	return false
}

// bytes in Go are uint8, whereas in Java they are int8. Hence this type alias.
type JavaByte = int8
type JavaByteSlice = []int8

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
