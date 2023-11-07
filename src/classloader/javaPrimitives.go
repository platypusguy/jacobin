/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/object"
	"jacobin/types"
)

// Implementation of some of the functions in in Java/lang/Class.

func Load_Primitives() map[string]GMeth {

	MethodSignatures["java/lang/Byte.valueOf(B)Ljava/lang/Byte;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  byteValueOf,
		}

	MethodSignatures["java/lang/Character.valueOf(C)Ljava/lang/Character;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  characterValueOf,
		}

	MethodSignatures["java/lang/Integer.valueOf(I)Ljava/lang/Integer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  integerValueOf,
		}

	MethodSignatures["java/lang/Long.valueOf(J)Ljava/lang/Long;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longValueOf,
		}

	MethodSignatures["java/lang/Short.valueOf(S)Ljava/lang/Short;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  shortValueOf,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
		}

	MethodSignatures["java/lang/Boolean.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  booleanJustReturn,
		}

	return MethodSignatures
}

func makePrimitiveObject(klass string, ftype string, arg any) *object.Object {
	objPtr := object.MakeEmptyObject()
	(*objPtr).Klass = &klass
	(*objPtr).Fields = append((*objPtr).Fields, object.Field{Ftype: ftype, Fvalue: arg})
	return objPtr
}

func byteValueOf(params []interface{}) interface{} {
	bb := params[0].(int64)
	objPtr := makePrimitiveObject("java/lang/Byte", types.Byte, bb)
	return objPtr
}

func characterValueOf(params []interface{}) interface{} {
	cc := params[0].(int64)
	objPtr := makePrimitiveObject("java/lang/Character", types.Char, cc)
	return objPtr
}

func integerValueOf(params []interface{}) interface{} {
	ii := params[0].(int64)
	objPtr := makePrimitiveObject("java/lang/Integer", types.Int, ii)
	return objPtr
}

func longValueOf(params []interface{}) interface{} {
	jj := params[0].(int64)
	objPtr := makePrimitiveObject("java/lang/Long", types.Long, jj)
	return objPtr
}

func shortValueOf(params []interface{}) interface{} {
	ss := params[0].(int64)
	objPtr := makePrimitiveObject("java/lang/Short", types.Short, ss)
	return objPtr
}

func booleanValueOf(params []interface{}) interface{} {
	zz := params[0].(int64)
	objPtr := makePrimitiveObject("java/lang/Boolean", types.Bool, zz)
	return objPtr
}

func booleanJustReturn(params []interface{}) interface{} {
	return nil
}
