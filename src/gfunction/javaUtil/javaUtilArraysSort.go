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
	"sort"
)

// utilArraysSort implements java.util.Arrays.sort and java.util.Arrays.parallelSort overloads.
func utilArraysSort(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysSort: too few arguments")
	}

	if params[0] == nil || params[0] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysSort: null array argument")
	}

	arrObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysSort: first arg not an array object")
	}

	field, ok := arrObj.FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysSort: missing array value field")
	}

	var fromIndex, toIndex int
	isRange := false

	if len(params) == 3 || len(params) == 4 {
		// sort(type[] a, int fromIndex, int toIndex)
		// or sort(Object[] a, int fromIndex, int toIndex, Comparator c)
		if _, ok := params[1].(int64); ok {
			if _, ok := params[2].(int64); ok {
				isRange = true
				fromIndex = int(params[1].(int64))
				toIndex = int(params[2].(int64))
			}
		}
	}

	switch a := field.Fvalue.(type) {
	case []types.JavaByte:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]
		sort.Slice(sub, func(i, j int) bool {
			return sub[i] < sub[j]
		})

	case []int64:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]
		sort.Slice(sub, func(i, j int) bool {
			return sub[i] < sub[j]
		})

	case []int32:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]
		sort.Slice(sub, func(i, j int) bool {
			return sub[i] < sub[j]
		})

	case []int16:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]
		sort.Slice(sub, func(i, j int) bool {
			return sub[i] < sub[j]
		})

	case []float32:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]
		sort.Slice(sub, func(i, j int) bool {
			// Basic float sort; Java has special handling for NaN and -0.0
			return sub[i] < sub[j]
		})

	case []float64:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]
		sort.Slice(sub, func(i, j int) bool {
			return sub[i] < sub[j]
		})

	case []*object.Object:
		if !isRange {
			fromIndex, toIndex = 0, len(a)
		}
		if fromIndex < 0 || toIndex > len(a) || fromIndex > toIndex {
			return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysSort: index out of bounds")
		}
		sub := a[fromIndex:toIndex]

		var comparator *object.Object
		if len(params) == 2 {
			// sort(Object[] a, Comparator c)
			if p1, ok := params[1].(*object.Object); ok {
				comparator = p1
			}
		} else if len(params) == 4 {
			// sort(Object[] a, int fromIndex, int toIndex, Comparator c)
			if p3, ok := params[3].(*object.Object); ok {
				comparator = p3
			}
		}

		if comparator != nil && !object.IsNull(comparator) {
			// TODO: Implement sorting with Comparator. Needs calling into JVM to execute compare().
			return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "utilArraysSort: sorting with Comparator not yet supported")
		}

		// Natural order sorting for Objects (must be Comparable)
		sort.Slice(sub, func(i, j int) bool {
			objI := sub[i]
			objJ := sub[j]
			if objI == nil || object.IsNull(objI) {
				return true // nulls first? Java sort(Object[]) throws NPE if it finds nulls usually
			}
			if objJ == nil || object.IsNull(objJ) {
				return false
			}

			// For now, only support String objects (Comparable)
			if object.IsStringObject(objI) && object.IsStringObject(objJ) {
				sI := object.GoStringFromStringObject(objI)
				sJ := object.GoStringFromStringObject(objJ)
				return sI < sJ
			}
			// Fallback: compare by hash or just say i < j is false
			return false
		})

	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysSort: unsupported array type")
	}

	return nil
}
