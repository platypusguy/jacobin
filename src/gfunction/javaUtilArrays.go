/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
)

// A partial implementation of the java/util/Arrays class.

func Load_Util_Arrays() {
	MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;I)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  copyOfObjectPointers,
		}

	MethodSignatures["java/util/Arrays.copyOf([Ljava/lang/Object;ILjava/lang/Class;)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}
}

// Copy the specified array of pointers, truncating or padding with nulls so the copy has the specified length.
func copyOfObjectPointers(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "copyOfObjectPointers: too few arguments")
	}

	// Check for a null array.
	if params[0] == nil {
		return getGErrBlk(excNames.NullPointerException, "copyOfObjectPointers: null array argument")
	}

	// Extract the array and the new length.
	parmObj := params[0].(*object.Object)
	newLen := int(params[1].(int64))

	// Check for a negative length.
	if newLen < 0 {
		return getGErrBlk(excNames.NegativeArraySizeException, "copyOfObjectPointers: negative array length")
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
