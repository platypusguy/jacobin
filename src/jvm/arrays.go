/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/shutdown"
	"unsafe"
)

/* This file contains the principal operations that Jacobin
   performs on arrays.

    An array is implemented as a struct with two fields:
    a pointer to the array and a second field, which holds
	a value indicating the type of elements in the array.

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
    even though it uses more RAM for the benefits it delivers.

    The code here was implemented by @alb as part of JACOBIN-203,
    based on code ginned up by @suresk.
*/

type ArrayType int

const (
	ERROR = 0
	FLOAT = 1
	INT   = 2
	BYTE  = 3
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

// type ArrType interface {
// 	int64 | byte | float64
// }

type Array struct {
	Type   ArrayType
	ArrPtr uintptr // will point to one of three kinds of arrays
}

// the fundamental way that an array is represented in Jacobin
// type ArrayHolder[T ArrType] struct {
// 	Type ArrayType
// 	Arr  []T
// }

type ByteArray []byte
type IntArray []int64
type FloatArray []float64

/*
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

*/

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
// (in the form of an unsafe pointer cast as a unintptr)
// and an int which identifies the type of array (where 0 = error)
//
// Note that once the array pointer is transmorgrified into an
// unsafe pointer in a uintptr field, go no longer realizes that
// there is a reference to the array. As a result, it could be
// garbage-collected at any moment. To avoid this, the address
// of the array is carefully saved in a linked list in the globals
// package, so that there will remain an active pointer to the array.
func createArray(arrayType int, size int64) (int, uintptr) {
	if size < 0 {
		exceptions.Throw(
			exceptions.NegativeArraySizeException,
			"Invalid size for array")
		shutdown.Exit(shutdown.APP_EXCEPTION)
	}

	g := globals.GetGlobalRef()

	aType := jdkArrayTypeToJacobinType(arrayType)
	if aType == BYTE {
		a := make([]byte, size)
		g.ArrayAddressList.PushFront(&a) // add address to list so array is not GC'd
		up := unsafe.Pointer(&a)
		ba := Array{Type: BYTE, ArrPtr: uintptr(up)}
		return BYTE, uintptr(unsafe.Pointer(&ba))

	} else if aType == INT {
		a := make([]int64, size)
		g.ArrayAddressList.PushFront(&a)
		up := unsafe.Pointer(&a)
		ia := Array{Type: INT, ArrPtr: uintptr(up)}
		return INT, uintptr(unsafe.Pointer(&ia))

	} else if aType == FLOAT {
		a := make([]float64, size)
		g.ArrayAddressList.PushFront(&a)
		up := unsafe.Pointer(&a)
		fa := Array{Type: FLOAT, ArrPtr: uintptr(up)}
		return FLOAT, uintptr(unsafe.Pointer(&fa))
	} else {
		return ERROR, 0
	}
}
