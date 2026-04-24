/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"slices"
)

func Load_Crypto_Spec_PBEParameterSpec() {
	ghelpers.MethodSignatures["javax/crypto/spec/PBEParameterSpec.<init>([BI)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  pbeParameterSpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEParameterSpec.<init>([BILjava/security/spec/AlgorithmParameterSpec;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  pbeParameterSpecInitWithSpec,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEParameterSpec.getIterationCount()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeParameterSpecGetIterationCount,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEParameterSpec.getSalt()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeParameterSpecGetSalt,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEParameterSpec.getParameterSpec()Ljava/security/spec/AlgorithmParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeParameterSpecGetParameterSpec,
		}
}

func pbeParameterSpecGetIterationCount(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecGetIterationCount: missing 'this'")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecGetIterationCount: 'this' is not an object")
	}

	iterationCount, ok := self.FieldTable["iterationCount"].Fvalue.(int64)
	if !ok {
		return int64(0)
	}

	return iterationCount
}

func pbeParameterSpecGetSalt(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecGetSalt: missing 'this'")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecGetSalt: 'this' is not an object")
	}

	salt, ok := self.FieldTable["salt"].Fvalue.([]byte)
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	// Returns a copy of the salt
	jBytes := object.JavaByteArrayFromGoByteArray(slices.Clone(salt))
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}

func pbeParameterSpecGetParameterSpec(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecGetParameterSpec: missing 'this'")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecGetParameterSpec: 'this' is not an object")
	}

	paramSpec, ok := self.FieldTable["paramSpec"].Fvalue.(*object.Object)
	if !ok || object.IsNull(paramSpec) {
		return ghelpers.ReturnNull(params)
	}

	return paramSpec
}

func pbeParameterSpecInit(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecInit: insufficient parameters")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecInit: 'this' is not an object")
	}

	saltObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(saltObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException,
			"pbeParameterSpecInit: salt cannot be null")
	}

	iterationCount := params[2].(int64)
	if iterationCount < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecInit: iterationCount must be >= 0")
	}

	salt := object.GoByteArrayFromJavaByteArray(saltObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	self.FieldTable["salt"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: slices.Clone(salt),
	}
	self.FieldTable["iterationCount"] = object.Field{
		Ftype:  types.Int,
		Fvalue: iterationCount,
	}

	return nil
}

func pbeParameterSpecInitWithSpec(params []any) any {
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecInitWithSpec: insufficient parameters")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecInitWithSpec: 'this' is not an object")
	}

	saltObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(saltObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException,
			"pbeParameterSpecInitWithSpec: salt cannot be null")
	}

	iterationCount := params[2].(int64)
	if iterationCount < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"pbeParameterSpecInitWithSpec: iterationCount must be >= 0")
	}

	paramSpec, ok := params[3].(*object.Object)
	// paramSpec can be null according to Java docs

	salt := object.GoByteArrayFromJavaByteArray(saltObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	self.FieldTable["salt"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: slices.Clone(salt),
	}
	self.FieldTable["iterationCount"] = object.Field{
		Ftype:  types.Int,
		Fvalue: iterationCount,
	}
	self.FieldTable["paramSpec"] = object.Field{
		Ftype:  "Ljava/security/spec/AlgorithmParameterSpec;",
		Fvalue: paramSpec,
	}

	return nil
}
