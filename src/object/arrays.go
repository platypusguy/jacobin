/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
    "jacobin/javaTypes"
    "unsafe"
)

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
    REF   = 4 // arrays of object references
    STR   = 5 // arrays of strings (these are references too, but for speed)
    ARR   = 6 // points to arrays, used in multidimensional arrays
    // ARRF  = 6  // points to arrays of floats--for multidimensional arrays
    // ARRI  = 7  // points to arrays of ints--for multidimensional arrays
    // ARRB  = 8  // points to arrays of bytes--for multidimensional arrays
    ARRR = 9  // points to arrays of references--for multidimensional arrays
    ARRG = 10 // generic array (of unsafe.Pointers)
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

type ArrayType int

type Ilength interface {
    Length() int64
}

// type JacobinByteArray struct {
// 	Type ArrayType
// 	Arr  *[]javaTypes.JavaByte
// }

// func (jba JacobinByteArray) Length() int64 {
// 	i := len(*(jba.Arr))
// 	return int64(i)
// }

type JacobinIntArray struct {
    Type ArrayType
    Arr  *[]int64
}

func (jba JacobinIntArray) Length() int64 {
    i := len(*(jba.Arr))
    return int64(i)
}

type JacobinFloatArray struct {
    Type ArrayType
    Arr  *[]float64
}

func (jba JacobinFloatArray) Length() int64 {
    i := len(*(jba.Arr))
    return int64(i)
}

type JacobinRefArray struct {
    Type ArrayType
    Arr  *[]*Object
}

func (jba JacobinRefArray) Length() int64 {
    i := len(*(jba.Arr))
    return int64(i)
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

// type JacobinArrByteArray struct {
// 	Type ArrayType
// 	Arr  *[]JacobinByteArray
// }

type JacobinArrRefArray struct {
    Type ArrayType
    Arr  *[]JacobinArrRefArray
}

type JacobinArrGenArray struct {
    Type ArrayType
    Arr  *[]unsafe.Pointer
}

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
    default: // this would indicate an error
        return 0
    }
}

// Make2DimArray creates a the last two dimensions of a multi-
// dimensional array. (All the dimensions > 2 are simply arrays
// of pointers to arrays.)
func Make2DimArray(ptrArrSize, leafArrSize int64,
    arrType uint8) (*Object, error) {

    ptrArr := MakeObject()                           // ptrArr is the pointer to the array of pointers to the leaf arrays
    value := make([]*Object, ptrArrSize, ptrArrSize) // the actual ptr-level array
    ptrArr.Fields = append(ptrArr.Fields, Field{Fvalue: &(value)})
    for i := 0; i < len(value); i++ { // for each entry in the ptr array
        value[i] = Make1DimArray(int(arrType), leafArrSize)
    }

    // the type of the pointer array will be the type of the leaf
    // array with a [ pre pended.
    ptrArrType := "[" + value[0].Fields[0].Ftype
    ptrArr.Fields[0].Ftype = ptrArrType
    // switch arrType {
    // case 'B': // byte arrays
    // 	ptrArr.Fields[0].Ftype = "[[B"
    // 	value[i].Fields[0].Ftype = "[B"
    // 	value[i].Fields[0].Fvalue =

    // value[i].Fields[0].Ftype = "[B"
    // actualArray := make([]byte, leafArrSize, leafArrSize)
    // value[i].Fields[0].Fvalue = &actualArray
    // 	}
    // }
    // TODO: ** COMMENTED OUT DUE to partial implementation of JACOBIN-261
    // var i int64
    // for i = 0; i < ptrArrSize; i++ {
    // 	switch arrType {
    // 	case 'B': // byte arrays
    // 		barArr := make([]javaTypes.JavaByte, leafArrSize)
    // 		ba := JacobinByteArray{
    // 			Type: BYTE,
    // 			Arr:  &barArr,
    // 		}
    // 		ptrArr[i] = unsafe.Pointer(&ba)
    // 	case 'F', 'D': // float arrays
    // 		farArr := make([]float64, leafArrSize)
    // 		fa := JacobinFloatArray{
    // 			Type: FLOAT,
    // 			Arr:  &farArr,
    // 		}
    // 		ptrArr[i] = unsafe.Pointer(&fa)
    // 	case 'L': // reference/pointer arrays
    // 		rarArr := make([]*object.Object, leafArrSize)
    // 		ra := JacobinRefArray{
    // 			Type: REF,
    // 			Arr:  &rarArr,
    // 		}
    // 		ptrArr[i] = unsafe.Pointer(&ra)
    // 	default: // all the integer types
    // 		iarArr := make([]int64, leafArrSize)
    // 		ia := JacobinIntArray{
    // 			Type: INT,
    // 			Arr:  &iarArr,
    // 		}
    // 		ptrArr[i] = unsafe.Pointer(&ia)
    // 	}
    // }
    //
    // multiArr := JacobinRefArray{
    // 	Type: ARRG,
    // 	Arr:  &ptrArr,
    // }
    return ptrArr, nil
}

func Make1DimArray(arrType int, size int64) *Object {
    o := MakeObject()
    o.Klass = nil // arrays don't have a pointer to a parsed class

    switch arrType {
    // case 'B': // byte arrays
    case BYTE:
        barArr := make([]javaTypes.JavaByte, size)
        of := Field{Ftype: "[B", Fvalue: &barArr}
        o.Fields = append(o.Fields, of)
        return o
    // case 'F', 'D': // float arrays
    case FLOAT:
        farArr := make([]float64, size)
        of := Field{Ftype: "[F", Fvalue: &farArr}
        o.Fields = append(o.Fields, of)
        return o
    case REF: // reference/pointer arrays
        rarArr := make([]*Object, size)
        of := Field{Ftype: "[L", Fvalue: &rarArr}
        o.Fields = append(o.Fields, of)
        return o
    default: // all the integer types
        iarArr := make([]int64, size)
        of := Field{Ftype: "[I", Fvalue: &iarArr}
        o.Fields = append(o.Fields, of)
        return o
    }
}

// MakeArrRefArray makes an array of pointers to other
// arrays of pointers. Each of these represents the elements
// of the dimensions > 2.
func MakeArrRefArray(size int64) *JacobinArrRefArray {
    rarArr := make([]JacobinArrRefArray, size)
    ra := JacobinArrRefArray{
        Type: ARRR,
        Arr:  &rarArr,
    }
    return &ra
}
