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

type ArrType interface {
	int64 | byte | float64
}

type ArrayHolder[T ArrType] struct {
	Type ArrayType
	Arr  []T
}
