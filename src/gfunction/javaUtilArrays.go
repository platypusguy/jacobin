/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

// A partial implementation of the java/util/Arrays class.

func Load_Util_Arrays() {

	MethodSignatures["java/util/Arrays.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: clinitGeneric}

	// asList
	MethodSignatures["java/util/Arrays.asList([Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  utilArraysAsList,
		}

	// binarySearch
	MethodSignatures["java/util/Arrays.binarySearch([B)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([B[B)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([B[BI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([BII)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([C)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([C[C)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([C[CI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([CI)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([D)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([D[D)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([D[DI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([DI)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([F)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([F[F)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([F[FI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([FI)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([I)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([I[I)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([I[II)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([II)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([J)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([J[J)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([J[JI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([JI)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([S)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([S[S)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([S[SI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([SI)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([Z)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([Z[Z)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([Z[ZI)I"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.binarySearch([ZI)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}

	// copyOf
	MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;I)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  utilArraysCopyOf,
		}
	MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;ILjava/lang/Class;)[Ljava/lang/Object;"] = GMeth{ParamSlots: 3, GFunction: trapFunction}

	// equals
	MethodSignatures["java/util/Arrays.equals([B[B)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([C[C)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([D[D)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([F[F)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([I[I)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([J[J)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([S[S)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	MethodSignatures["java/util/Arrays.equals([Z[Z)Z"] = GMeth{ParamSlots: 2, GFunction: utilArraysEquals}

	// fill
	MethodSignatures["java/util/Arrays.fill([BB)V"] = GMeth{
		ParamSlots: 2,
		GFunction:  utilArraysFillBytes,
	}
	MethodSignatures["java/util/Arrays.fill([BBII)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([CC)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([DD)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([FF)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([II)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([JJ)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([SS)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([ZZ)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([Ljava/lang/Object;Ljava/lang/Object;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.fill([Ljava/lang/Object;IILjava/lang/Object;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}

	// hashCode
	MethodSignatures["java/util/Arrays.hashCode([B)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([C)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([D)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([F)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([I)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([J)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([Ljava/lang/Object;)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([S)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.hashCode([Z)I"] = GMeth{ParamSlots: 1, GFunction: trapFunction}

	// mismatch
	MethodSignatures["java/util/Arrays.mismatch([B[B)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([C[C)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([D[D)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([F[F)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([I[I)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([J[J)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([S[S)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.mismatch([Z[Z)I"] = GMeth{ParamSlots: 2, GFunction: trapFunction}

	// parallelPrefix
	MethodSignatures["java/util/Arrays.parallelPrefix([DLjava/util/function/DoubleBinaryOperator;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([DIILjava/util/function/DoubleBinaryOperator;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([FLjava/util/function/FloatBinaryOperator;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([FIILjava/util/function/FloatBinaryOperator;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([ILjava/util/function/IntBinaryOperator;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([IIILjava/util/function/IntBinaryOperator;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([JLjava/util/function/LongBinaryOperator;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelPrefix([JIILjava/util/function/LongBinaryOperator;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}

	// parallelSetAll
	MethodSignatures["java/util/Arrays.parallelSetAll([DLjava/util/function/IntToDoubleFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSetAll([FLjava/util/function/IntToFloatFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSetAll([ILjava/util/function/IntUnaryOperator;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSetAll([JLjava/util/function/IntToLongFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSetAll([Ljava/lang/Object;Ljava/util/function/IntFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}

	// parallelSort
	MethodSignatures["java/util/Arrays.parallelSort([D)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([DII)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([F)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([FII)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([I)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([III)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([J)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([JII)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;II)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;IILjava/util/Comparator;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}

	// setAll
	MethodSignatures["java/util/Arrays.setAll([DLjava/util/function/IntToDoubleFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.setAll([FLjava/util/function/IntToFloatFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.setAll([ILjava/util/function/IntUnaryOperator;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.setAll([JLjava/util/function/IntToLongFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.setAll([Ljava/lang/Object;Ljava/util/function/IntFunction;)V"] = GMeth{ParamSlots: 2, GFunction: trapFunction}

	// sort
	MethodSignatures["java/util/Arrays.sort([D)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([DII)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([F)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([FII)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([I)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([III)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([J)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([JII)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;II)V"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;IILjava/util/Comparator;)V"] = GMeth{ParamSlots: 4, GFunction: trapFunction}

	// spliterator
	MethodSignatures["java/util/Arrays.spliterator([D)Ljava/util/Spliterator$OfDouble;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.spliterator([F)Ljava/util/Spliterator$OfFloat;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.spliterator([I)Ljava/util/Spliterator$OfInt;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.spliterator([J)Ljava/util/Spliterator$OfLong;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.spliterator([Ljava/lang/Object;)Ljava/util/Spliterator;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}

	// stream
	MethodSignatures["java/util/Arrays.stream([D)Ljava/util/DoubleStream;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([DII)Ljava/util/DoubleStream;"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([F)Ljava/util/FloatStream;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([FII)Ljava/util/FloatStream;"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([I)Ljava/util/IntStream;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([III)Ljava/util/IntStream;"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([J)Ljava/util/LongStream;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([JII)Ljava/util/LongStream;"] = GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([Ljava/lang/Object;)Ljava/util/Stream;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/Arrays.stream([Ljava/lang/Object;II)Ljava/util/Stream;"] = GMeth{ParamSlots: 3, GFunction: trapFunction}

	// toString
	MethodSignatures["java/util/Arrays.toString([B)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([C)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([D)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([F)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([I)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([J)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([Ljava/lang/Object;)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([S)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	MethodSignatures["java/util/Arrays.toString([Z)Ljava/lang/String;"] = GMeth{ParamSlots: 1, GFunction: utilArraysToString}

}

// Arrays.equals(Object[] a, Object[] b) -> boolean
// Minimal implementation that compares reference arrays element-wise by reference equality.
func utilArraysEquals(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysEquals: too few arguments")
	}

	// Handle nulls per Java semantics
	if params[0] == nil && params[1] == nil {
		return types.JavaBoolTrue
	}
	if params[0] == nil || params[1] == nil {
		return types.JavaBoolFalse
	}

	objA, okA := params[0].(*object.Object)
	objB, okB := params[1].(*object.Object)
	if !okA || !okB {
		return types.JavaBoolFalse
	}

	fieldA, ok := objA.FieldTable["value"]
	if !ok {
		return types.JavaBoolFalse
	}
	fieldB, ok := objB.FieldTable["value"]
	if !ok {
		return types.JavaBoolFalse
	}

	switch a := fieldA.Fvalue.(type) {
	case []*object.Object:
		b, ok := fieldB.Fvalue.([]*object.Object)
		if !ok {
			return types.JavaBoolFalse
		}
		if len(a) != len(b) {
			return types.JavaBoolFalse
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] { // reference equality (nil handled by !=)
				return types.JavaBoolFalse
			}
		}
		return types.JavaBoolTrue
	case []types.JavaByte:
		b, ok := fieldB.Fvalue.([]types.JavaByte)
		if !ok {
			return types.JavaBoolFalse
		}
		if len(a) != len(b) {
			return types.JavaBoolFalse
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				return types.JavaBoolFalse
			}
		}
		return types.JavaBoolTrue
	case []int64:
		b, ok := fieldB.Fvalue.([]int64)
		if !ok {
			return types.JavaBoolFalse
		}
		if len(a) != len(b) {
			return types.JavaBoolFalse
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				return types.JavaBoolFalse
			}
		}
		return types.JavaBoolTrue
	case []float64:
		b, ok := fieldB.Fvalue.([]float64)
		if !ok {
			return types.JavaBoolFalse
		}
		if len(a) != len(b) {
			return types.JavaBoolFalse
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] { // minimal; does not handle NaN bitwise cases
				return types.JavaBoolFalse
			}
		}
		return types.JavaBoolTrue
	default:
		// Unsupported array backing type
		return types.JavaBoolFalse
	}
}

// Copy the specified array of pointers, truncating or padding with nulls so the copy has the specified length.
func utilArraysCopyOf(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysCopyOf: too few arguments")
	}

	// Check for a null array.
	if params[0] == nil {
		return getGErrBlk(excNames.NullPointerException, "utilArraysCopyOf: null array argument")
	}

	// Extract the array and the new length.
	parmObj := params[0].(*object.Object)
	newLen := int(params[1].(int64))

	// Check for a negative length.
	if newLen < 0 {
		return getGErrBlk(excNames.NegativeArraySizeException, "utilArraysCopyOf: negative array length")
	}

	// Get the array length.
	parmObject := *parmObj
	arr := parmObject.FieldTable["value"]
	rawArrayOld := arr.Fvalue.([]*object.Object)
	oldLen := len(rawArrayOld)

	// Create a new array of the desired length.
	newArrayObj := object.Make1DimRefArray("java/lang/Object;", int64(newLen))
	rawArrayNew := newArrayObj.FieldTable["value"].Fvalue.([]*object.Object)

	// Copy the elements from the old array to the new array.
	for i := 0; i < oldLen && i < newLen; i++ {
		rawArrayNew[i] = rawArrayOld[i]
	}

	if newLen > oldLen {
		for i := oldLen; i < newLen; i++ {
			rawArrayNew[i] = nil
		}
	}

	return newArrayObj
}

// Arrays.asList([Ljava/lang/Object;)Ljava/util/List;
func utilArraysAsList(params []interface{}) interface{} {
	if len(params) < 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: too few arguments")
	}

	if params[0] == nil || params[0] == object.Null {
		return getGErrBlk(excNames.NullPointerException, "utilArraysAsList: array is null")
	}

	arrayObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: first arg not an object")
	}

	field, ok := arrayObj.FieldTable["value"]
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: missing array value field")
	}

	// Reference arrays store elements as []*object.Object
	var ifaceElements []interface{}
	switch elements := field.Fvalue.(type) {
	case []*object.Object:
		ifaceElements = make([]interface{}, len(elements))
		for i, e := range elements {
			ifaceElements[i] = e
		}
	case []interface{}:
		ifaceElements = make([]interface{}, len(elements))
		copy(ifaceElements, elements)
	default:
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: unsupported array type")
	}

	// Return an ArrayList object
	listObj := object.MakePrimitiveObject("java/util/ArrayList", types.ArrayList, ifaceElements)
	return listObj
}

// Arrays.fill(byte[] a, byte val) -> void (also used for boolean[])
func utilArraysFillBytes(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: too few arguments")
	}
	if params[0] == nil {
		return getGErrBlk(excNames.NullPointerException, "utilArraysFillBytes: null array argument")
	}
	arrObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: first arg not an array object")
	}
	field, ok := arrObj.FieldTable["value"]
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: missing array value field")
	}
	// Java byte value is passed as int64
	val := types.JavaByte(params[1].(int64))
	switch a := field.Fvalue.(type) {
	case []types.JavaByte:
		for i := range a {
			a[i] = val
		}
		// write back (not strictly necessary as slice is by reference)
		arrObj.FieldTable["value"] = field
		return nil
	default:
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: unsupported array type")
	}
}

// Arrays.toString for primitive and reference arrays
func utilArraysToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "utilArraysToString: too few arguments")
	}
	if params[0] == nil {
		return object.StringObjectFromGoString("null")
	}
	arrObj, ok := params[0].(*object.Object)
	if !ok {
		return object.StringObjectFromGoString("null")
	}
	field, ok := arrObj.FieldTable["value"]
	if !ok {
		return object.StringObjectFromGoString("null")
	}

	var b strings.Builder
	b.WriteByte('[')

	appendComma := func(idx, l int) {
		if idx+1 < l {
			b.WriteString(", ")
		}
	}

	switch v := field.Fvalue.(type) {
	case []types.JavaByte:
		for i, e := range v {
			b.WriteString(fmt.Sprintf("%d", int8(e)))
			appendComma(i, len(v))
		}
	case []int64:
		// Special-case char arrays by type
		if field.Ftype == types.CharArray {
			for i, e := range v {
				b.WriteString(string(rune(e)))
				appendComma(i, len(v))
			}
		} else {
			for i, e := range v {
				b.WriteString(fmt.Sprintf("%d", e))
				appendComma(i, len(v))
			}
		}
	case []float64:
		for i, e := range v {
			b.WriteString(fmt.Sprintf("%v", e))
			appendComma(i, len(v))
		}
	case []*object.Object:
		for i, e := range v {
			if e == nil || object.IsNull(e) {
				b.WriteString("null")
			} else if object.IsStringObject(e) {
				b.WriteString(object.GoStringFromStringObject(e))
			} else {
				// Fallback to class name
				b.WriteString(object.GoStringFromStringPoolIndex(e.KlassName))
			}
			appendComma(i, len(v))
		}
	default:
		// Unknown backing; return []
	}

	b.WriteByte(']')
	return object.StringObjectFromGoString(b.String())
}
