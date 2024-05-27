/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/log"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
)

func Load_Math_Big_Integer() {

	MethodSignatures["java/math/BigInteger.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerClinit,
		}

	MethodSignatures["java/math/BigInteger.valueOf(J)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerValueOf,
		}

}

var bigIntegerClassName = "java/math/BigInteger"

// initBigIntegerFields: Initiaalise the object fields.
func initBigInteger(obj *object.Object, argValue int64) {
	fld := object.Field{Ftype: types.IntArray, Fvalue: []int64{argValue}}
	obj.FieldTable["mag"] = fld
	fld = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	obj.FieldTable["signum"] = fld
	fld = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	obj.FieldTable["bitCountPlusOne"] = fld
	obj.FieldTable["bitLengthPlusOne"] = fld
	obj.FieldTable["lowestSetBitPlusTwo"] = fld
	obj.FieldTable["firstNonzeroIntNumPlusTwo"] = fld
}

// addStaticBigInteger: Form a BigInteger object based on the parameter value.
func addStaticBigInteger(argName string, argValue int64) {
	name := fmt.Sprintf("java/math/BigInteger.%s", argName)
	obj := object.MakeEmptyObjectWithClassName(&bigIntegerClassName)
	initBigInteger(obj, argValue)
	_ = statics.AddStatic(name, statics.Static{Type: "Ljava/math/BigInteger;", Value: obj})
}

// "java/math/BigInteger.<clinit>()V"
func bigIntegerClinit(params []interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/math/BigInteger")
	if klass == nil {
		errMsg := "bigIntegerClinit: Expected java/math/BigInteger to be in the MethodArea, but it was not"
		_ = log.Log(errMsg, log.SEVERE)
		return getGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}
	if klass.Data.ClInit != types.ClInitRun {
		addStaticBigInteger("ONE", int64(1))
		addStaticBigInteger("TEN", int64(10))
		addStaticBigInteger("TWO", int64(2))
		addStaticBigInteger("ZERO", int64(0))
		klass.Data.ClInit = types.ClInitRun
	}
	return nil
}

// "java/math/BigInteger.valueOf(J)Ljava/math/BigInteger;"
func bigIntegerValueOf(params []interface{}) interface{} {
	// params[0] holds the class object
	argValue := params[1].(int64)

	obj := object.MakeEmptyObjectWithClassName(&bigIntegerClassName)
	initBigInteger(obj, argValue)

	return obj
}
