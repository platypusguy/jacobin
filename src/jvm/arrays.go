/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import "unsafe"

/*  This file contains some data structures and some functions
 	for array handling in Jacobin

    An array is implemented as a struct with two fields:
	a value indicating the type of elements in the array and
    a pointer to the array itself
	.

	We use a pointer to the array b/c in Go, if you pass an
	array to a function, the entire array is copied over. We
	don't want that!

    For our purposes, there are three possible array types:
    int64 (all integral types and addresses), float64 (all
    FP types), and bytes (for bytes and boolean/bits)

    The official JVM docs suggest that bit arrays (so booleans)
    can be implemented as individual byte elements or aggregated
    eight a time into a single byte. Like the Oracle JVM,
    we opted for the former option due to performance and simplicity,
    even though it uses more RAM.
*/

const ( // the ArrayTypes
	ERROR = 0
	FLOAT = 1
	INT   = 2
	BYTE  = 3
	REF   = 4
	ARR   = 5  // points to arrays, used in multidimensional arrays
	ARRF  = 6  // points to arrays of floats--for multidimensional arrays
	ARRI  = 7  // points to arrays of ints--for multidimensional arrays
	ARRB  = 8  // points to arrays of bytes--for multidimensional arrays
	ARRR  = 9  // points to arrays of references--for multidimensional arrays
	ARRG  = 10 // generic array (of unsafe.Pointers)
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

// bytes in Go are uint8, whereas in Java they are int8. Hence this type alias.
type JavaByte = int8

type ArrayType int

type JacobinByteArray struct {
	Type ArrayType
	Arr  *[]JavaByte
}

type JacobinIntArray struct {
	Type ArrayType
	Arr  *[]int64
}

type JacobinFloatArray struct {
	Type ArrayType
	Arr  *[]float64
}

type JacobinRefArray struct {
	Type ArrayType
	Arr  *[]unsafe.Pointer
}

// === The following types are used only in multidimensional arrays
// Array that points to other arrays.
type JacobinArrArray struct {
	Type ArrayType
	Arr  *[]JacobinArrArray
}

type JacobinArrFloatArray struct {
	Type ArrayType
	Arr  *[]JacobinFloatArray
}

type JacobinArrIntArray struct {
	Type ArrayType
	Arr  *[]JacobinIntArray
}

type JacobinArrByteArray struct {
	Type ArrayType
	Arr  *[]JacobinByteArray
}

type JacobinArrRefArray struct {
	Type ArrayType
	Arr  *[]JacobinRefArray
}

type JacobinArrGenArray struct {
	Type ArrayType
	Arr  *[]unsafe.Pointer
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
