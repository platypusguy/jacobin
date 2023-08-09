/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import "jacobin/types"

/*  This file contains some data structures and some functions
 	for array handling in Jacobin

    An array is implemented as a struct with two fields:
	a value indicating the type of elements in the array and
    a pointer to the array itself

	We use a pointer to the array b/c in Go, if you pass an
	array to a function, the entire array is copied over. We
	don't want that!

    For our purposes, there are four possible array types:
    int64 (all integral types ), float64 (all FP types), bytes
	(for bytes and boolean/bits), and references (i.e. pointers)

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
	REF   = 4 // arrays of object references
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

const T_REF = 12 // used only in Jacobin

// converts one the of the JDK values indicating the primitive
// used in the elements of an array into one of the values used
// by Jacobin in array creation. Returns zero on error.
func JdkArrayTypeToJacobinType(jdkType int) int {
	switch jdkType {
	case T_BOOLEAN, T_BYTE:
		return BYTE
	case T_CHAR, T_SHORT, T_INT, T_LONG:
		return INT
	case T_FLOAT, T_DOUBLE:
		return FLOAT
	case T_REF:
		return REF // technically not one of the JDK categories
		// but needed for our purposes.
	default: // this would indicate an error
		return 0
	}
}

// Make2DimArray creates a the last two dimensions of a multi-
// dimensional array. (All the dimensions > 2 are simply arrays
// of pointers to arrays.)
func Make2DimArray(ptrArrSize, leafArrSize int64, arrType uint8) (*Object, error) {
	ptrArr := MakeEmptyObject()                      // ptrArr is the pointer to the array of pointers to the leaf arrays
	value := make([]*Object, ptrArrSize, ptrArrSize) // the actual ptr-level array
	ptrArr.Fields = append(ptrArr.Fields, Field{Fvalue: &(value)})
	for i := 0; i < len(value); i++ { // for each entry in the ptr array
		value[i] = Make1DimArray(arrType, leafArrSize)
	}

	// the type of the pointer array will be the type of the leaf
	// array with a [ pre pended.
	ptrArrType := "[" + value[0].Fields[0].Ftype
	ptrArr.Fields[0].Ftype = ptrArrType

	ptrArr.Klass = &value[0].Fields[0].Ftype

	return ptrArr, nil
}

// Make1DimArray creates and 1-diminensional Jacobin-style array
// of the specified type (passed as a byte) and size.
func Make1DimArray(arrType uint8, size int64) *Object {
	o := MakeEmptyObject()
	var of Field

	switch arrType {
	// case 'B': // byte arrays
	case BYTE:
		// barArr := make([]types.JavaByte, size) // changed with JACOBIN-282
		barArr := make([]byte, size)
		of = Field{Ftype: types.ByteArray, Fvalue: &barArr}
		o.Fields = append(o.Fields, of)
	// case 'F', 'D': // float arrays
	case FLOAT:
		farArr := make([]float64, size)
		of := Field{Ftype: types.FloatArray, Fvalue: &farArr}
		o.Fields = append(o.Fields, of)
	case REF: // reference/pointer arrays
		rarArr := make([]*Object, size)
		of := Field{Ftype: types.RefArray, Fvalue: &rarArr}
		o.Fields = append(o.Fields, of)
	default: // all the integer types
		iarArr := make([]int64, size)
		of := Field{Ftype: types.IntArray, Fvalue: &iarArr}
		o.Fields = append(o.Fields, of)
	}
	o.Klass = &o.Fields[0].Ftype // in arrays, Klass field is a pointer to the array type string
	return o
}

// MakeArrayFromRawArray accepts a raw array (such as []byte) and
// converts it into an array *object*.
func MakeArrayFromRawArray(rawArray interface{}) *Object {
	switch rawArray.(type) {
	case *Object: // if it's a ref to an array object, just return it
		arr := rawArray.(*Object)
		return arr
	case *[]uint8: // an array of bytes
		raw := rawArray.(*[]uint8)
		o := MakeEmptyObject()
		o.Klass = nil
		of := Field{Ftype: types.ByteArray, Fvalue: raw}
		o.Fields = append(o.Fields, of)
		return o
	}
	return nil
}
