/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
)

func Load_Crypto_Spec_DHParameterSpec() {

	// <init>(BigInteger p, BigInteger g)
	ghelpers.MethodSignatures["javax/crypto/spec/DHParameterSpec.<init>(Ljava/math/BigInteger;Ljava/math/BigInteger;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  dhparameterspecInit,
		}

	// <init>(BigInteger p, BigInteger g, int l)
	ghelpers.MethodSignatures["javax/crypto/spec/DHParameterSpec.<init>(Ljava/math/BigInteger;Ljava/math/BigInteger;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  dhparameterspecInit,
		}

	// getP()
	ghelpers.MethodSignatures["javax/crypto/spec/DHParameterSpec.getP()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dhparameterspecGetP,
		}

	// getG()
	ghelpers.MethodSignatures["javax/crypto/spec/DHParameterSpec.getG()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dhparameterspecGetG,
		}

	// getL()
	ghelpers.MethodSignatures["javax/crypto/spec/DHParameterSpec.getL()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dhparameterspecGetL,
		}
}

func dhparameterspecInit(args []interface{}) interface{} {

	const funcName = "dhparameterspecInit"

	if len(args) != 3 && len(args) != 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": wrong number of arguments")
	}

	obj, ok := args[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": receiver is not an object")
	}

	pObj, ok1 := args[1].(*object.Object)
	if !ok1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": p must be a BigInteger object")
	}
	p, ok1 := pObj.FieldTable["value"].Fvalue.(*big.Int)

	gObj, ok2 := args[2].(*object.Object)
	if !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": g must be a BigInteger object")
	}
	g, ok2 := gObj.FieldTable["value"].Fvalue.(*big.Int)

	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": p and g must be valid BigInteger objects")
	}

	var lVal int64 = 0

	if len(args) == 4 {
		var ok3 bool
		lVal, ok3 = args[3].(int64)
		if !ok3 {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				funcName+": l must be int64")
		}
	}

	obj.FieldTable["p"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: p,
	}

	obj.FieldTable["g"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: g,
	}

	obj.FieldTable["l"] = object.Field{
		Ftype:  types.Int,
		Fvalue: lVal,
	}

	return nil
}

func dhparameterspecGetP(args []interface{}) interface{} {

	const funcName = "dhparameterspecGetP"

	if len(args) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": wrong number of arguments")
	}

	obj := args[0].(*object.Object)

	p, ok := obj.FieldTable["p"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException,
			funcName+": p not initialized")
	}

	return object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, p)
}

func dhparameterspecGetG(args []interface{}) interface{} {

	const funcName = "dhparameterspecGetG"

	if len(args) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": wrong number of arguments")
	}

	obj := args[0].(*object.Object)

	g, ok := obj.FieldTable["g"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException,
			funcName+": g not initialized")
	}

	return object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, g)
}

func dhparameterspecGetL(args []interface{}) interface{} {

	const funcName = "dhparameterspecGetL"

	if len(args) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			funcName+": wrong number of arguments")
	}

	obj := args[0].(*object.Object)

	field, ok := obj.FieldTable["l"].Fvalue.(int64)
	if !ok {
		return int64(0) // matches JDK behavior
	}

	return field
}
