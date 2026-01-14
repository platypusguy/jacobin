/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// utilArraysFill implements all java.util.Arrays.fill overloads.
// Overloads:
// fill(type[] a, type val)
// fill(type[] a, int fromIndex, int toIndex, type val)
func utilArraysFill(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFill: too few arguments")
	}

	if params[0] == nil || params[0] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysFill: null array argument")
	}

	arrObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFill: first arg not an array object")
	}

	field, ok := arrObj.FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFill: missing array value field")
	}

	var fromIndex, toIndex int
	var val interface{}
	isRange := false

	if len(params) == 4 {
		isRange = true
		fromIndex = int(params[1].(int64))
		toIndex = int(params[2].(int64))
		val = params[3]
	} else {
		val = params[1]
	}

	switch a := field.Fvalue.(type) {
	case []types.JavaByte:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		fillVal := types.JavaByte(val.(int64))
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	case []int64:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		fillVal := val.(int64)
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	case []int32:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		fillVal := int32(val.(int64))
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	case []int16:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		fillVal := int16(val.(int64))
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	case []float32:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		fillVal := val.(float32)
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	case []float64:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		fillVal := val.(float64)
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	case []*object.Object:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysFill: index out of bounds")
		}
		var fillVal *object.Object
		if val != nil {
			fillVal = val.(*object.Object)
		}
		for i := fromIndex; i < toIndex; i++ {
			a[i] = fillVal
		}
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysFill: unsupported array type")
	}

	return nil
}
