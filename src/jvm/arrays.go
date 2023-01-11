/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

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
    eight a time into a single byte. For simplicity (and very likely
    for performance benefits) we opted for the former option, even
    though it uses more RAM for the benefits it delivers.

    The code here was implemented by @alb as part of JACOBIN-203,
    based on code ginned up by @suresk.
*/

type ArrayType int

const (
	FLOAT = iota
	INT   = iota
	BYTE  = iota
)

// the primitive types as specified in the
// JVM instructions // for arrays
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

type ArrayHolder[T ArrType] struct {
	Type ArrayType
	Arr  []T
}

// converts one the of the JDK values indicating the primitive
// used in the elements of an array into one of the values used
// by Jacobin in array creation. Returns zero on error.
func jdkArrayTypeToJacobinType(jdkType byte) int {
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
