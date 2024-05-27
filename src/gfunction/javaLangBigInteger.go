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

func Load_Lang_Big_Integer() {

	MethodSignatures["java/lang/BigInteger.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerClinit,
		}

}

// addStaticBigInteger: Form a BigInteger object based on the parameter value.
func addStaticBigInteger(argName string, argValue int) {
	name := fmt.Sprintf("java/lang/BigInteger.%s", argName)
	obj := object.MakePrimitiveObject("java/lang/BigInteger", types.Int, int64(argValue))
	_ = statics.AddStatic(name, statics.Static{Type: "Ljava/lang/BigInteger;", Value: obj})
}

// "java/lang/BigInteger.<clinit>()V"
func bigIntegerClinit(params []interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/lang/BigInteger")
	if klass == nil {
		errMsg := "bigIntegerClinit: Expected java/lang/BigInteger to be in the MethodArea, but it was not"
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
