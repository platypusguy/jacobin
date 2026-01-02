/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"reflect"
)

/*  This file contains some data structures and some primitive
    functions for array handling in Jacobin

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
    eight a time into a single byte. Like the HotSpot JVM,
    we opted for the former option due to performance and simplicity,
    even though it uses more RAM.
*/

//const ( // the ArrayTypes
//	ERROR = 0
//	FLOAT = 1
//	INT   = 2
//	BYTE  = 3
//	REF   = 4 // arrays of object references
//)

// the primitive types as specified in the
// JVM instructions for arrays
const (
	T_ERROR   = 86
	T_BOOLEAN = 4
	T_CHAR    = 5
	T_FLOAT   = 6
	T_DOUBLE  = 7
	T_BYTE    = 8
	T_SHORT   = 9
	T_INT     = 10
	T_LONG    = 11
	T_REF     = 12
)

// converts one the of the JDK values indicating the primitive
// used in the elements of an array into one of the values used
// by Jacobin in array creation. Returns zero on error.
//func JdkArrayTypeToJacobinType(jdkType int) int {
//	switch jdkType {
//	case T_BOOLEAN, T_BYTE:
//		return BYTE
//	case T_CHAR, T_SHORT, T_INT, T_LONG:
//		return INT
//	case T_FLOAT, T_DOUBLE:
//		return FLOAT
//	case T_REF:
//		return REF // technically not one of the JDK categories
//		// but needed for our purposes.
//	default: // this would indicate an error
//		return 0
//	}
//}

// identifies the type of entry that the array is made up of
// i.e., primitives or specific references. Note if the array is a
// reference array, the trailing ; is *not* removed.
func GetArrayType(arrayType string) string {
	typeChars := []byte(arrayType)
	for index, char := range typeChars {
		if char != '[' {
			return string(typeChars[index:])
		}
	}
	return arrayType
}

// Make2DimArray creates the last two dimensions of a multi// dimensional array.
// (All the dimensions > 2 are simply arrays  of pointers to arrays.)
func Make2DimArray(ptrArrSize, leafArrSize int64, arrType uint8) (*Object, error) {
	ptrArr := MakeEmptyObject()          // ptrArr is the pointer to the array of pointers to the leaf arrays
	value := make([]*Object, ptrArrSize) // the actual ptr-level array

	// make the first array in the ptr array and get its type converted to a string
	firstArray := Make1DimArray(arrType, leafArrSize)
	aType := "[" + *(stringPool.GetStringPointer(firstArray.KlassName)) // the type of the first array
	ptrArr.FieldTable["value"] = Field{
		Ftype:  aType,
		Fvalue: value,
	}

	value[0] = firstArray
	for i := 1; i < len(value); i++ { // for each entry in the ptr array
		value[i] = Make1DimArray(arrType, leafArrSize)
	}

	ptrArr.KlassName = firstArray.KlassName
	return ptrArr, nil
}

// Make1DimArray creates and 1-dimensional Jacobin-style array
// of the specified type (passed as a byte) and size.
func Make1DimArray(arrType uint8, size int64) *Object {
	o := MakeEmptyObject()
	var of Field

	// JACOBIN-457: Converted to exclusive use of o.FieldTable and o.Fields
	// contain the actual value rather than a pointer to the value. 2024-02
	switch arrType {
	case T_BOOLEAN:
		barArr := make([]types.JavaBool, size)
		of = Field{Ftype: types.BoolArray, Fvalue: barArr}
	case T_BYTE:
		barArr := make([]types.JavaByte, size)
		of = Field{Ftype: types.ByteArray, Fvalue: barArr}
	case T_CHAR: // integer arrays
		farArr := make([]int64, size)
		of = Field{Ftype: types.CharArray, Fvalue: farArr}
	case T_DOUBLE: // double arrays
		farArr := make([]float64, size)
		of = Field{Ftype: types.DoubleArray, Fvalue: farArr}
	case T_FLOAT: // float arrays
		farArr := make([]float64, size)
		of = Field{Ftype: types.FloatArray, Fvalue: farArr}
	case T_INT: // integer arrays
		farArr := make([]int64, size)
		of = Field{Ftype: types.IntArray, Fvalue: farArr}
	case T_LONG: // long arrays
		farArr := make([]int64, size)
		of = Field{Ftype: types.LongArray, Fvalue: farArr}
	case T_SHORT: // short arrays
		farArr := make([]int64, size)
		of = Field{Ftype: types.ShortArray, Fvalue: farArr}
	case T_REF: // reference/pointer arrays
		rarArr := make([]*Object, size)
		of = Field{Ftype: types.RefArray, Fvalue: rarArr}
	default:
		errMsg := fmt.Sprintf("object.Make1DimArray() was passed an unsupported array type: %d", arrType)
		globals.GetGlobalRef().FuncThrowException(excNames.IllegalArgumentException, errMsg)
		return nil
	}
	o.FieldTable["value"] = of
	value := o.FieldTable["value"]
	o.KlassName = stringPool.GetStringIndex(&value.Ftype) // in arrays, Klass field is a pointer to the array type string
	return o
}

// Make1DimRefArray makes a 1-dimensional reference array. Its logic is nearly identical to
// Make1DimArray, except that it is passed a string identifying the type of object in
// the array and it inserts that value into the field and object type fields.
func Make1DimRefArray(objType string, size int64) *Object {
	o := MakeEmptyObject()
	rarArr := make([]*Object, size)
	arrayType := types.RefArray + objType
	of := Field{Ftype: arrayType, Fvalue: rarArr}
	o.FieldTable["value"] = of
	o.KlassName = stringPool.GetStringIndex(&of.Ftype)
	return o
}

// MakeArrayFromRawArray accepts a raw array (such as []byte) and
// converts it into an array *object*.
func MakeArrayFromRawArray(rawArray interface{}) *Object {
	if rawArray == nil {
		errMsg := fmt.Sprintf("object.MakeArrayFromRawArray() was passed a nil parameter")
		globals.GetGlobalRef().FuncThrowException(excNames.IllegalArgumentException, errMsg)
		// trace.Warning(errMsg)
		return nil
	}

	switch rawArray.(type) {
	case *Object: // if it's a ref to an array object, just return it
		arr := rawArray.(*Object)
		return arr
	case []*Object: // if it's a ref to an array of objects, just return it
		obj := MakePrimitiveObject("java/lang/Object", "[Ljava/lang/Object;", rawArray.([]*Object))
		return obj
	case *[]types.JavaByte: // an array of bytes
		objPtr :=
			MakePrimitiveObject(types.ByteArray, types.ByteArray, *rawArray.(*[]types.JavaByte))
		return objPtr
	}

	arrType := reflect.TypeOf(rawArray)
	if arrType.Kind() == reflect.Slice {
		switch arrType.Elem().Kind() {
		case reflect.Int8:
			return MakePrimitiveObject(types.ByteArray, types.ByteArray, rawArray.([]int8))
		case reflect.Uint8:
			// convert array of bytes into JavaBytes
			jba := JavaByteArrayFromGoByteArray(rawArray.([]uint8))
			return MakePrimitiveObject(types.ByteArray, types.ByteArray, jba)
		case reflect.Int64:
			return MakePrimitiveObject(types.IntArray, types.IntArray, rawArray.([]int64))
		}
	}

	// This code basically turns an array into a slice.
	if arrType.Kind() == reflect.Array {
		slice := make([]interface{}, arrType.Len())
		arrValue := reflect.ValueOf(rawArray)
		for i := 0; i < arrValue.Len(); i++ {
			item := arrValue.Index(i).Interface()
			slice[i] = item
		}

		switch arrType.Elem().Kind() {
		case reflect.Int8:
			return MakePrimitiveObject(types.ByteArray, types.ByteArray, slice)
		case reflect.Uint8:
			return MakePrimitiveObject(types.ByteArray, types.ByteArray, slice)
		case reflect.Int64:
			return MakePrimitiveObject(types.IntArray, types.IntArray, slice)
		}
	}

	errMsg := fmt.Sprintf("object.MakeArrayFromRawArray() was passed an unsupported type: %T", rawArray)
	globals.GetGlobalRef().FuncThrowException(excNames.IllegalArgumentException, errMsg)
	return nil
}

// ArrayLength returns the length of an array object, when passed a pointer to it
func ArrayLength(arrayRef *Object) int64 {
	var size int64
	o := arrayRef.FieldTable["value"]
	arrayType := o.Ftype
	switch arrayType {
	case types.ByteArray:
		array := o.Fvalue.([]types.JavaByte)
		size = int64(len(array))
	case types.BoolArray:
		array := o.Fvalue.([]types.JavaBool)
		size = int64(len(array))
	case types.RefArray:
		array := o.Fvalue.([]*Object)
		size = int64(len(array))
	case types.FloatArray, types.DoubleArray:
		array := o.Fvalue.([]float64)
		size = int64(len(array))
	case types.IntArray, types.LongArray, types.ShortArray, types.CharArray:
		array := o.Fvalue.([]int64)
		size = int64(len(array))
	default:
		array := o.Fvalue.([]*Object)
		size = int64(len(array))
	}
	return size
}
