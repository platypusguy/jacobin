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

}

// addStaticBigInteger: Form a BigInteger object based on the parameter value.
func addStaticBigInteger(argName string, argValue int) {
	name := fmt.Sprintf("java/math/BigInteger.%s", argName)
	obj := object.MakePrimitiveObject("java/math/BigInteger", types.Int, int64(argValue))
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
		addStaticBigInteger("ONE", 1)
		addStaticBigInteger("TEN", 10)
		addStaticBigInteger("TWO", 2)
		addStaticBigInteger("ZERO", 0)
		klass.Data.ClInit = types.ClInitRun
	}
	return nil
}
