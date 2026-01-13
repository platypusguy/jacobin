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

// utilArraysCompare implements java.util.Arrays.compare and compareUnsigned for primitive types.
// It excludes Comparable and Comparator based overloads.
func utilArraysCompare(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCompare: too few arguments")
	}

	// Handle nulls per Java semantics
	if params[0] == nil || params[0] == object.Null {
		if params[1] == nil || params[1] == object.Null {
			return int64(0)
		}
		return int64(-1)
	}
	if params[1] == nil || params[1] == object.Null {
		return int64(1)
	}

	arrObj1, ok1 := params[0].(*object.Object)
	var arrObj2 *object.Object
	var ok2 bool
	if len(params) >= 6 {
		arrObj2, ok2 = params[3].(*object.Object)
	} else {
		arrObj2, ok2 = params[1].(*object.Object)
	}

	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCompare: arguments must be array objects")
	}

	field1, ok1 := arrObj1.FieldTable["value"]
	field2, ok2 := arrObj2.FieldTable["value"]
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCompare: missing array value field")
	}

	// Range parameters
	var fromIndex1, toIndex1, fromIndex2, toIndex2 int
	isRange := false

	if len(params) >= 6 {
		isRange = true
		fromIndex1 = int(params[1].(int64))
		toIndex1 = int(params[2].(int64))
		fromIndex2 = int(params[4].(int64)) // params[3] is the second array
		toIndex2 = int(params[5].(int64))
	}

	checkBounds := func(length, from, to int) bool {
		return from >= 0 && to >= from && to <= length
	}

	switch a := field1.Fvalue.(type) {
	case []int64: // handles long, int, short, char, boolean
		b, ok := field2.Fvalue.([]int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysCompare: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysCompare: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}

		// Determine if this is a boolean array to use boolean comparison logic
		isBool := false
		if field1.Ftype == types.BoolArray {
			isBool = true
		}

		for i := 0; i < n; i++ {
			v1 := a[fromIndex1+i]
			v2 := b[fromIndex2+i]
			if v1 != v2 {
				if isBool {
					// Java boolean compare: false < true
					if v1 == types.JavaBoolFalse {
						return int64(-1)
					}
					return int64(1)
				}
				if v1 < v2 {
					return int64(-1)
				}
				return int64(1)
			}
		}
		return int64(len1 - len2)

	case []types.JavaByte:
		b, ok := field2.Fvalue.([]types.JavaByte)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysCompare: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysCompare: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			v1 := a[fromIndex1+i]
			v2 := b[fromIndex2+i]
			if v1 != v2 {
				if v1 < v2 {
					return int64(-1)
				}
				return int64(1)
			}
		}
		return int64(len1 - len2)

	case []int32: // handles int or char
		b, ok := field2.Fvalue.([]int32)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysCompare: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysCompare: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			v1 := a[fromIndex1+i]
			v2 := b[fromIndex2+i]
			if v1 != v2 {
				if v1 < v2 {
					return int64(-1)
				}
				return int64(1)
			}
		}
		return int64(len1 - len2)

	case []int16: // handles short
		b, ok := field2.Fvalue.([]int16)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysCompare: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysCompare: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			v1 := a[fromIndex1+i]
			v2 := b[fromIndex2+i]
			if v1 != v2 {
				if v1 < v2 {
					return int64(-1)
				}
				return int64(1)
			}
		}
		return int64(len1 - len2)

	case []float32:
		b, ok := field2.Fvalue.([]float32)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysCompare: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysCompare: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			v1 := a[fromIndex1+i]
			v2 := b[fromIndex2+i]
			if v1 != v2 {
				// float comparison in Java: -0.0 < 0.0, NaN == NaN and NaN > everything else
				// For now, simple comparison.
				if v1 < v2 {
					return int64(-1)
				}
				return int64(1)
			}
		}
		return int64(len1 - len2)

	case []float64:
		b, ok := field2.Fvalue.([]float64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysCompare: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysCompare: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			v1 := a[fromIndex1+i]
			v2 := b[fromIndex2+i]
			if v1 != v2 {
				if v1 < v2 {
					return int64(-1)
				}
				return int64(1)
			}
		}
		return int64(len1 - len2)

	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysCompare: unsupported array type")
	}
}
