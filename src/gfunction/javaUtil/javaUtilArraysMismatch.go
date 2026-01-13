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

// utilArraysMismatch implements all java.util.Arrays.mismatch overloads.
func utilArraysMismatch(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysMismatch: too few arguments")
	}

	if params[0] == nil || params[0] == object.Null || params[1] == nil || params[1] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "utilArraysMismatch: null array argument")
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysMismatch: arguments must be array objects")
	}

	field1, ok1 := arrObj1.FieldTable["value"]
	field2, ok2 := arrObj2.FieldTable["value"]
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysMismatch: missing array value field")
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

	// Helper function for bounds checking
	checkBounds := func(length, from, to int) bool {
		return from >= 0 && to >= from && to <= length
	}

	switch a := field1.Fvalue.(type) {
	case []types.JavaByte:
		b, ok := field2.Fvalue.([]types.JavaByte)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			if a[fromIndex1+i] != b[fromIndex2+i] {
				return int64(i)
			}
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	case []int64:
		b, ok := field2.Fvalue.([]int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			if a[fromIndex1+i] != b[fromIndex2+i] {
				return int64(i)
			}
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	case []int32:
		b, ok := field2.Fvalue.([]int32)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			if a[fromIndex1+i] != b[fromIndex2+i] {
				return int64(i)
			}
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	case []int16:
		b, ok := field2.Fvalue.([]int16)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			if a[fromIndex1+i] != b[fromIndex2+i] {
				return int64(i)
			}
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	case []float32:
		b, ok := field2.Fvalue.([]float32)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			// In Java mismatch, float/double comparison handles NaN specially
			// but here we do simple comparison for now.
			if a[fromIndex1+i] != b[fromIndex2+i] {
				return int64(i)
			}
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	case []float64:
		b, ok := field2.Fvalue.([]float64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			if a[fromIndex1+i] != b[fromIndex2+i] {
				return int64(i)
			}
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	case []*object.Object:
		b, ok := field2.Fvalue.([]*object.Object)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ClassCastException, "utilArraysMismatch: incompatible array types")
		}
		if !isRange {
			fromIndex1, toIndex1 = 0, len(a)
			fromIndex2, toIndex2 = 0, len(b)
		} else {
			if !checkBounds(len(a), fromIndex1, toIndex1) || !checkBounds(len(b), fromIndex2, toIndex2) {
				return ghelpers.GetGErrBlk(excNames.ArrayIndexOutOfBoundsException, "utilArraysMismatch: index out of bounds")
			}
		}
		len1 := toIndex1 - fromIndex1
		len2 := toIndex2 - fromIndex2
		n := len1
		if len2 < n {
			n = len2
		}
		for i := 0; i < n; i++ {
			obj1 := a[fromIndex1+i]
			obj2 := b[fromIndex2+i]
			if obj1 == obj2 {
				continue
			}
			if obj1 == nil || obj2 == nil {
				return int64(i)
			}
			// For Object[], mismatch uses equals()
			// For now, minimal support: reference equality was already checked above.
			// If we want to support Object.equals(), we'd need to call into the VM.
			// But since we are in a GFunction, let's stick to reference equality for now
			// as it's common in this codebase's minimal implementations.
			// However, if they are Strings, we can compare them.
			if object.IsStringObject(obj1) && object.IsStringObject(obj2) {
				if object.GoStringFromStringObject(obj1) == object.GoStringFromStringObject(obj2) {
					continue
				}
			}
			return int64(i)
		}
		if len1 != len2 {
			return int64(n)
		}
		return int64(-1)

	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "utilArraysMismatch: unsupported array type")
	}
}
