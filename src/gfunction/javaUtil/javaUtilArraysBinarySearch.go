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

// Arrays.binarySearch implementations for all primitive and reference arrays
func utilArraysBinarySearch(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysBinarySearch: too few arguments")
	}
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysBinarySearch: null array argument")
	}
	arrObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysBinarySearch: first arg not an array object")
	}
	field, ok := arrObj.FieldTable["value"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysBinarySearch: missing array value field")
	}

	var fromIndex, toIndex int
	var key interface{}
	// var comparator *object.Object // unused for now but present in some signatures

	if len(params) == 2 {
		fromIndex = 0
		key = params[1]
		// Determine toIndex from array length later
	} else if len(params) == 3 {
		// Could be (Object[], Object, Comparator)
		fromIndex = 0
		key = params[1]
		// params[2] would be comparator
	} else if len(params) == 4 {
		// (Primitive[], int from, int to, Primitive key)
		fromIndex = int(params[1].(int64))
		toIndex = int(params[2].(int64))
		key = params[3]
	} else if len(params) == 5 {
		// (Object[], int from, int to, Object key, Comparator c)
		fromIndex = int(params[1].(int64))
		toIndex = int(params[2].(int64))
		key = params[3]
		// params[4] is comparator
	}

	switch a := field.Fvalue.(type) {
	case []types.JavaByte:
		if len(params) == 2 {
			toIndex = len(a)
		}
		k := types.JavaByte(key.(int64))
		low := fromIndex
		high := toIndex - 1
		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]
			if midVal < k {
				low = mid + 1
			} else if midVal > k {
				high = mid - 1
			} else {
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	case []int64:
		if len(params) == 2 {
			toIndex = len(a)
		}
		k := key.(int64)
		low := fromIndex
		high := toIndex - 1
		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]
			if midVal < k {
				low = mid + 1
			} else if midVal > k {
				high = mid - 1
			} else {
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	case []float64:
		if len(params) == 2 {
			toIndex = len(a)
		}
		k := key.(float64)
		low := fromIndex
		high := toIndex - 1
		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]
			if midVal < k {
				low = mid + 1
			} else if midVal > k {
				high = mid - 1
			} else {
				// bitwise equals for NaN and -0.0/+0.0 would go here in full impl
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	case []float32:
		if len(params) == 2 {
			toIndex = len(a)
		}
		k := float32(key.(float64))
		low := fromIndex
		high := toIndex - 1
		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]
			if midVal < k {
				low = mid + 1
			} else if midVal > k {
				high = mid - 1
			} else {
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	case []int32:
		if len(params) == 2 {
			toIndex = len(a)
		}
		k := int32(key.(int64))
		low := fromIndex
		high := toIndex - 1
		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]
			if midVal < k {
				low = mid + 1
			} else if midVal > k {
				high = mid - 1
			} else {
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	case []int16:
		if len(params) == 2 {
			toIndex = len(a)
		}
		k := int16(key.(int64))
		low := fromIndex
		high := toIndex - 1
		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]
			if midVal < k {
				low = mid + 1
			} else if midVal > k {
				high = mid - 1
			} else {
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	case []*object.Object:
		if len(params) == 2 || len(params) == 3 {
			toIndex = len(a)
		}
		// Minimal: only handles Comparable if key is string or similar, or uses Comparator if provided.
		// For now, let's just do a basic implementation for String objects.
		low := fromIndex
		high := toIndex - 1
		kObj, ok := key.(*object.Object)
		if !ok || kObj == nil {
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysBinarySearch: null key")
		}

		for low <= high {
			mid := (low + high) >> 1
			midVal := a[mid]

			// Use CompareTo if available, or basic logic.
			// This is complex because we might need to invoke a Java method.
			// For this task, we'll implement a basic version for Strings as a placeholder.
			var cmp int
			if object.IsStringObject(kObj) && object.IsStringObject(midVal) {
				s1 := object.GoStringFromStringObject(midVal)
				s2 := object.GoStringFromStringObject(kObj)
				if s1 < s2 {
					cmp = -1
				} else if s1 > s2 {
					cmp = 1
				} else {
					cmp = 0
				}
			} else {
				// Fallback or trap if not supported
				return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "utilArraysBinarySearch: Object comparison not fully implemented")
			}

			if cmp < 0 {
				low = mid + 1
			} else if cmp > 0 {
				high = mid - 1
			} else {
				return int64(mid)
			}
		}
		return int64(-(low + 1))

	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysBinarySearch: unsupported array type")
	}
}
