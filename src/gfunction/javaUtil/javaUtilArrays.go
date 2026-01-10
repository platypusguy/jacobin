/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

// A partial implementation of the java/util/Arrays class.

func Load_Util_Arrays() {

	ghelpers.MethodSignatures["java/util/Arrays.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ClinitGeneric}

	// asList
	ghelpers.MethodSignatures["java/util/Arrays.asList([Ljava/lang/Object;)Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  utilArraysAsList,
		}

	// binarySearch
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([B)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([B[B)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([B[BI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([BII)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([C)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([C[C)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([C[CI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([CI)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([D)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([D[D)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([D[DI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([DI)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([F)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([F[F)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([F[FI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([FI)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([I)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([I[I)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([I[II)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([II)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([J)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([J[J)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([J[JI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([JI)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([S)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([S[S)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([S[SI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([SI)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Z)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Z[Z)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Z[ZI)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([ZI)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// copyOf
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;I)[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  utilArraysCopyOf,
		}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;ILjava/lang/Class;)[Ljava/lang/Object;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	// equals
	ghelpers.MethodSignatures["java/util/Arrays.equals([B[B)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([C[C)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([D[D)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([F[F)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([I[I)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([J[J)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([S[S)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}
	ghelpers.MethodSignatures["java/util/Arrays.equals([Z[Z)Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysEquals}

	// fill
	ghelpers.MethodSignatures["java/util/Arrays.fill([BB)V"] = ghelpers.GMeth{
		ParamSlots: 2,
		GFunction:  utilArraysFillBytes,
	}
	ghelpers.MethodSignatures["java/util/Arrays.fill([BBII)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([CC)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([DD)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([FF)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([II)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([JJ)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([SS)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([ZZ)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([Ljava/lang/Object;Ljava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.fill([Ljava/lang/Object;IILjava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}

	// hashCode
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([B)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([C)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([D)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([F)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([I)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([J)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([Ljava/lang/Object;)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([S)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.hashCode([Z)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	// mismatch
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([B[B)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([C[C)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([D[D)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([F[F)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([I[I)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([J[J)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([S[S)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([Z[Z)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// parallelPrefix
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([DLjava/util/function/DoubleBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([DIILjava/util/function/DoubleBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([FLjava/util/function/FloatBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([FIILjava/util/function/FloatBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([ILjava/util/function/IntBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([IIILjava/util/function/IntBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([JLjava/util/function/LongBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([JIILjava/util/function/LongBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}

	// parallelSetAll
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([DLjava/util/function/IntToDoubleFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([FLjava/util/function/IntToFloatFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([ILjava/util/function/IntUnaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([JLjava/util/function/IntToLongFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([Ljava/lang/Object;Ljava/util/function/IntFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// parallelSort
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([D)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([DII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([F)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([FII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([III)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([J)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([JII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;II)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;IILjava/util/Comparator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}

	// setAll
	ghelpers.MethodSignatures["java/util/Arrays.setAll([DLjava/util/function/IntToDoubleFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([FLjava/util/function/IntToFloatFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([ILjava/util/function/IntUnaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([JLjava/util/function/IntToLongFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([Ljava/lang/Object;Ljava/util/function/IntFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// sort
	ghelpers.MethodSignatures["java/util/Arrays.sort([D)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([DII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([F)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([FII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([III)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([J)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([JII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;II)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;IILjava/util/Comparator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}

	// spliterator
	ghelpers.MethodSignatures["java/util/Arrays.spliterator([D)Ljava/util/Spliterator$OfDouble;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.spliterator([F)Ljava/util/Spliterator$OfFloat;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.spliterator([I)Ljava/util/Spliterator$OfInt;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.spliterator([J)Ljava/util/Spliterator$OfLong;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.spliterator([Ljava/lang/Object;)Ljava/util/Spliterator;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	// stream
	ghelpers.MethodSignatures["java/util/Arrays.stream([D)Ljava/util/DoubleStream;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([DII)Ljava/util/DoubleStream;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([F)Ljava/util/FloatStream;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([FII)Ljava/util/FloatStream;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([I)Ljava/util/IntStream;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([III)Ljava/util/IntStream;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([J)Ljava/util/LongStream;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([JII)Ljava/util/LongStream;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([Ljava/lang/Object;)Ljava/util/Stream;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.stream([Ljava/lang/Object;II)Ljava/util/Stream;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	// toString
	ghelpers.MethodSignatures["java/util/Arrays.toString([B)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([C)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([D)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([F)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([I)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([J)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([Ljava/lang/Object;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([S)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}
	ghelpers.MethodSignatures["java/util/Arrays.toString([Z)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysToString}

}

// Arrays.equals(Object[] a, Object[] b) -> boolean
// Minimal implementation that compares reference arrays element-wise by reference equality.
func utilArraysEquals(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysEquals: too few arguments")
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCopyOf: too few arguments")
	}

	// Check for a null array.
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysCopyOf: null array argument")
	}

	// Extract the array and the new length.
	parmObj := params[0].(*object.Object)
	newLen := int(params[1].(int64))

	// Check for a negative length.
	if newLen < 0 {
		return ghelpers.GetGErrBlk(excNames.NegativeArraySizeException, "utilArraysCopyOf: negative array length")
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: too few arguments")
	}

	if params[0] == nil || params[0] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysAsList: array is null")
	}

	arrayObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: first arg not an object")
	}

	field, ok := arrayObj.FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: missing array value field")
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysAsList: unsupported array type")
	}

	// Return an ArrayList object
	listObj := object.MakePrimitiveObject("java/util/ArrayList", types.ArrayList, ifaceElements)
	return listObj
}

// Arrays.fill(byte[] a, byte val) -> void (also used for boolean[])
func utilArraysFillBytes(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: too few arguments")
	}
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysFillBytes: null array argument")
	}
	arrObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: first arg not an array object")
	}
	field, ok := arrObj.FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: missing array value field")
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFillBytes: unsupported array type")
	}
}

// Arrays.toString for primitive and reference arrays
func utilArraysToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysToString: too few arguments")
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
