/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Util_Objects() {

	MethodSignatures["java/util/Objects.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Objects.checkFromIndexSize(III)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  objectsCheckFromIndexSize,
		}

	MethodSignatures["java/util/Objects.checkFromIndexSize(JJJ)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  objectsCheckFromIndexSize,
		}

	MethodSignatures["java/util/Objects.checkFromToIndex(III)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  objectsCheckFromToIndex,
		}

	MethodSignatures["java/util/Objects.checkFromToIndex(JJJ)J"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  objectsCheckFromToIndex,
		}

	MethodSignatures["java/util/Objects.checkIndex(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  objectsCheckIndex,
		}

	MethodSignatures["java/util/Objects.checkIndex(JJ)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  objectsCheckIndex,
		}

	// TODO: Not even trapped: static <T> int compare(T a, T b, Comparator<? super T> c)

	MethodSignatures["java/util/Objects.deepEquals(Ljava/lang/Object;Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Objects.equals(Ljava/lang/Object;Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  objectsEquals,
		}

	MethodSignatures["java/util/Objects.hash([Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Objects.hashCode(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Objects.isNull(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  objectsIsNull,
		}

	MethodSignatures["java/util/Objects.nonNull(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  objectsNonNull,
		}

	// TODO: Not even trapped: requireNonNull
	// TODO: Not even trapped: requireNonNullElse
	// TODO: Not even trapped: requireNonNullElseGet
	// TODO: Not even trapped: static String toIdentityString(Object o)

	MethodSignatures["java/util/Objects.toIdentityString(Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Objects.toString(Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Objects.toString(Ljava/lang/Object;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

}

func objectsCheckFromIndexSize(params []interface{}) interface{} {
	fromIndex := params[0].(int64)
	size := params[1].(int64)
	length := params[1].(int64)
	if fromIndex < 0 || size < 0 || fromIndex+size > length || length < 0 {
		errMsg := fmt.Sprintf("objectsFromIndexSize: Invalid parameters: fromIndex=%d, size=%d, length=%d",
			fromIndex, size, length)
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}
	return fromIndex
}

func objectsCheckFromToIndex(params []interface{}) interface{} {
	fromIndex := params[0].(int64)
	toIndex := params[1].(int64)
	length := params[1].(int64)
	if fromIndex < 0 || fromIndex > toIndex || toIndex > length || length < 0 {
		errMsg := fmt.Sprintf("objectsCheckFromToIndex: Invalid parameters: fromIndex=%d, toIndex=%d, length=%d",
			fromIndex, toIndex, length)
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}
	return fromIndex
}

func objectsCheckIndex(params []interface{}) interface{} {
	index := params[0].(int64)
	length := params[1].(int64)
	if length < 0 || index < 0 || index >= length {
		errMsg := fmt.Sprintf("objectsCheckIndex: Invalid parameters: index=%d, length=%d", index, length)
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}
	return index
}

// java/util/Objects.equals(Object a, Object b) -> boolean
// Minimal implementation:
// - returns true if both are null
// - returns false if exactly one is null
// - returns true if both are the same reference
// - otherwise returns false (does not invoke a.equals(b))
func objectsEquals(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "objectsEquals: too few arguments")
	}

	if params[0] == nil && params[1] == nil {
		return types.JavaBoolTrue
	}
	if params[0] == nil || params[1] == nil {
		return types.JavaBoolFalse
	}

	a, okA := params[0].(*object.Object)
	b, okB := params[1].(*object.Object)
	if !okA || !okB {
		return types.JavaBoolFalse
	}

	if a == b {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func objectsIsNull(params []interface{}) interface{} {
	_, ok := params[0].(*object.Object)
	if ok {
		return types.JavaBoolFalse
	}
	return types.JavaBoolTrue
}

func objectsNonNull(params []interface{}) interface{} {
	_, ok := params[0].(*object.Object)
	return ok
}
