/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/exceptions"
	"jacobin/shutdown"
	"unsafe"
)

/* This file contains the principal operations that Jacobin
   performs on arrays.

    An array is implemented as a struct with two fields:
    the array itself and a second field, which holds a value
    that indicates what type of elements are in the array.

    For our purposes, there are three possible array types:
    int64 (all integral types and addresses), float64 (all
    FP types), and bytes (for bytes and boolean/bits)

    The official JVM docs suggest that bit arrays (so booleans)
    can be implemented as individual byte elements or aggregated
    eight a time into a single byte. Like the Oracle JVM,
    we opted for the former option due to performance and simplicity,
    even though it uses more RAM for the benefits it delivers.

    The code here was implemented by @alb as part of JACOBIN-203,
    based on code ginned up by @suresk.
*/

type ArrayType int

const (
	ERROR = 0
	FLOAT = iota
	INT   = iota
	BYTE  = iota
)

// the primitive types as specified in the
// JVM instructions for arrays
const (
	T_BOOLEAN = 4
	T_CHAR    = 5
	T_FLOAT   = 6
	T_DOUBLE  = 7
	T_BYTE    = 8
	T_SHORT   = 9
	T_INT     = 10
	T_LONG    = 11
)

type ArrType interface {
	int64 | byte | float64
}

// the fundamental way that an array is represented in Jacobin
// type ArrayHolder[T ArrType] struct {
// 	Type ArrayType
// 	Arr  []T
// }

type ByteArray struct {
	Type  ArrayType
	Array []byte
}

type IntArray struct {
	Type  ArrayType
	Array []int64
}

type FloatArray struct {
	Type  ArrayType
	Array []float64
}

// converts one the of the JDK values indicating the primitive
// used in the elements of an array into one of the values used
// by Jacobin in array creation. Returns zero on error.
func jdkArrayTypeToJacobinType(jdkType int) int {
	switch jdkType {
	case T_BOOLEAN, T_BYTE:
		return BYTE
	case T_CHAR, T_SHORT, T_INT, T_LONG:
		return INT
	case T_FLOAT, T_DOUBLE:
		return FLOAT
	default: // this would indicate an error
		return 0
	}
}

// creates an array struct and returns a pointer to it
// (in the form of an unsafe pointer cast as an int)
// and an int which identifies the type of array (where 0 = error)
func createArray(arrayType int, size int64) (int, uintptr) {
	if size < 0 {
		exceptions.Throw(
			exceptions.NegativeArraySizeException,
			"Invalid size for array")
		shutdown.Exit(shutdown.APP_EXCEPTION)
	}

	aType := jdkArrayTypeToJacobinType(arrayType)
	if aType == BYTE {
		a := make([]byte, size)
		ba := ByteArray{Type: BYTE, Array: a}
		return BYTE, uintptr(unsafe.Pointer(&ba))
	} else if aType == INT {
		a := make([]int64, size)
		ia := IntArray{Type: INT, Array: a}
		return INT, uintptr(unsafe.Pointer(&ia))
	} else if aType == FLOAT {
		a := make([]float64, size)
		fa := FloatArray{Type: FLOAT, Array: a}
		return FLOAT, uintptr(unsafe.Pointer(&fa))
	} else {
		return ERROR, 0
	}
}
