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
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([BB)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([BIIB)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([CC)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([CIIC)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([DD)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([DIID)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([FF)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([FIIF)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([II)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([IIII)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([JJ)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([JIIJ)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([SS)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([SIIS)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Ljava/lang/Object;Ljava/lang/Object;)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Ljava/lang/Object;IILjava/lang/Object;)I"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Ljava/lang/Object;Ljava/lang/Object;Ljava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysBinarySearch}
	ghelpers.MethodSignatures["java/util/Arrays.binarySearch([Ljava/lang/Object;IILjava/lang/Object;Ljava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 5, GFunction: utilArraysBinarySearch}

	// compare
	ghelpers.MethodSignatures["java/util/Arrays.compare([Z[Z)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([ZII[ZII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([B[B)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([BII[BII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([C[C)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([CII[CII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([D[D)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([DII[DII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([F[F)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([FII[FII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([I[I)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([III[III)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([J[J)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([JII[JII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([S[S)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compare([SII[SII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}

	// compare (Object/Comparator)
	ghelpers.MethodSignatures["java/util/Arrays.compare([Ljava/lang/Object;[Ljava/lang/Object;)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.compare([Ljava/lang/Object;II[Ljava/lang/Object;II)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.compare([Ljava/lang/Object;[Ljava/lang/Object;Ljava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.compare([Ljava/lang/Object;II[Ljava/lang/Object;IILjava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 7, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([Ljava/lang/Object;[Ljava/lang/Object;Ljava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([Ljava/lang/Object;II[Ljava/lang/Object;IILjava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 7, GFunction: ghelpers.TrapFunction}

	// compareUnsigned
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([B[B)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([BII[BII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([I[I)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([III[III)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([J[J)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([JII[JII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([S[S)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCompare}
	ghelpers.MethodSignatures["java/util/Arrays.compareUnsigned([SII[SII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysCompare}

	// copyOf
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([BI)[B"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([CI)[C"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([DI)[D"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([FI)[F"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([II)[I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([JI)[J"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([SI)[S"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([ZI)[Z"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysCopyOfPrimitive}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;I)[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  utilArraysCopyOf,
		}
	ghelpers.MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;ILjava/lang/Class;)[Ljava/lang/Object;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysCopyOfObjectWithClass}

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
	ghelpers.MethodSignatures["java/util/Arrays.fill([BB)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([BBII)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([CC)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([CIIC)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([DD)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([DIID)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([FF)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([FIIF)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([II)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([IIII)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([JJ)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([JIIJ)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([SS)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([SIIS)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([ZZ)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([ZIIZ)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([Ljava/lang/Object;Ljava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysFill}
	ghelpers.MethodSignatures["java/util/Arrays.fill([Ljava/lang/Object;IILjava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysFill}

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
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([B[B)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([BII[BII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([C[C)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([CII[CII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([D[D)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([DII[DII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([F[F)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([FII[FII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([I[I)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([III[III)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([J[J)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([JII[JII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([S[S)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([SII[SII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([Z[Z)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([ZII[ZII)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([Ljava/lang/Object;[Ljava/lang/Object;)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([Ljava/lang/Object;II[Ljava/lang/Object;II)I"] = ghelpers.GMeth{ParamSlots: 6, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([Ljava/lang/Object;II[Ljava/lang/Object;IILjava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 7, GFunction: utilArraysMismatch}
	ghelpers.MethodSignatures["java/util/Arrays.mismatch([Ljava/lang/Object;[Ljava/lang/Object;Ljava/util/Comparator;)I"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysMismatch}

	// parallelPrefix
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([DLjava/util/function/DoubleBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([DIILjava/util/function/DoubleBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([FLjava/util/function/FloatBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([FIILjava/util/function/FloatBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([ILjava/util/function/IntBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([IIILjava/util/function/IntBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([JLjava/util/function/LongBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([JIILjava/util/function/LongBinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelPrefix([Ljava/lang/Object;Ljava/util/function/BinaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// parallelSetAll
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([DLjava/util/function/IntToDoubleFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([FLjava/util/function/IntToFloatFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([ILjava/util/function/IntUnaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([JLjava/util/function/IntToLongFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSetAll([Ljava/lang/Object;Ljava/util/function/IntFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// parallelSort
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([D)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([DII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([F)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([FII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([III)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([J)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([JII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;II)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;IILjava/util/Comparator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.parallelSort([Ljava/lang/Object;Ljava/util/function/Comparator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// setAll
	ghelpers.MethodSignatures["java/util/Arrays.setAll([DLjava/util/function/IntToDoubleFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([FLjava/util/function/IntToFloatFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([ILjava/util/function/IntUnaryOperator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([JLjava/util/function/IntToLongFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Arrays.setAll([Ljava/lang/Object;Ljava/util/function/IntFunction;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	// sort
	ghelpers.MethodSignatures["java/util/Arrays.sort([B)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([BII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([C)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([CII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([D)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([DII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([F)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([FII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([III)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([J)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([JII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([S)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([SII)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;II)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;Ljava/util/Comparator;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: utilArraysSort}
	ghelpers.MethodSignatures["java/util/Arrays.sort([Ljava/lang/Object;IILjava/util/Comparator;)V"] = ghelpers.GMeth{ParamSlots: 4, GFunction: utilArraysSort}

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

	// Get the array length and type info.
	arrField := parmObj.FieldTable["value"]

	switch old := arrField.Fvalue.(type) {
	case []*object.Object:
		oldLen := len(old)
		// Preserve component type if available: arrField.Ftype is like "[L<class>;"
		objType := object.GetArrayType(arrField.Ftype) // e.g., Ljava/lang/String;
		newArrayObj := object.Make1DimRefArray(objType, int64(newLen))
		rawArrayNew := newArrayObj.FieldTable["value"].Fvalue.([]*object.Object)
		for i := 0; i < oldLen && i < newLen; i++ {
			rawArrayNew[i] = old[i]
		}
		// remaining entries already nil
		return newArrayObj
	case []types.JavaByte:
		oldLen := len(old)
		newArr := make([]types.JavaByte, newLen)
		for i := 0; i < oldLen && i < newLen; i++ {
			newArr[i] = old[i]
		}
		return object.MakePrimitiveObject(arrField.Ftype, arrField.Ftype, newArr)
	case []int64:
		oldLen := len(old)
		newArr := make([]int64, newLen)
		for i := 0; i < oldLen && i < newLen; i++ {
			newArr[i] = old[i]
		}
		return object.MakePrimitiveObject(arrField.Ftype, arrField.Ftype, newArr)
	case []float64:
		oldLen := len(old)
		newArr := make([]float64, newLen)
		for i := 0; i < oldLen && i < newLen; i++ {
			newArr[i] = old[i]
		}
		return object.MakePrimitiveObject(arrField.Ftype, arrField.Ftype, newArr)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCopyOf: unsupported array type")
	}
}

// utilArraysCopyOfPrimitive handles primitive array copyOf variants uniformly
func utilArraysCopyOfPrimitive(params []interface{}) interface{} {
	return utilArraysCopyOf(params)
}

// utilArraysCopyOfObjectWithClass: copyOf(Object[] a, int newLen, Class newType)
// Minimal implementation: ignore newType and behave like copyOf(Object[], int)
func utilArraysCopyOfObjectWithClass(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCopyOfObjectWithClass: too few arguments")
	}
	// Reuse utilArraysCopyOf for copying behavior
	return utilArraysCopyOf([]interface{}{params[0], params[1]})
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
