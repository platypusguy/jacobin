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

func Load_Crypto_Spec_GCMParameterSpec() {
	ghelpers.MethodSignatures["javax/crypto/spec/GCMParameterSpec.<init>(I[B)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  gcmParameterSpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/GCMParameterSpec.<init>(I[BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  gcmParameterSpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/GCMParameterSpec.getIV()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  gcmParameterSpecGetIV,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/GCMParameterSpec.getTLen()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  gcmParameterSpecGetTLen,
		}
}

func gcmParameterSpecGetIV(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecGetIV: missing 'this'")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecGetIV: 'this' is not an object")
	}

	iv, ok := self.FieldTable["iv"].Fvalue.([]byte)
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	// Returns a copy of the IV
	jBytes := object.JavaByteArrayFromGoByteArray(slices.Clone(iv))
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}

func gcmParameterSpecGetTLen(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecGetTLen: missing 'this'")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecGetTLen: 'this' is not an object")
	}

	tLen, ok := self.FieldTable["tLen"].Fvalue.(int64)
	if !ok {
		return int64(0)
	}

	return tLen
}

func gcmParameterSpecInit(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecInit: insufficient parameters")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecInit: 'this' is not an object")
	}

	tLen := params[1].(int64)

	ivObj, ok := params[2].(*object.Object)
	if !ok || object.IsNull(ivObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException,
			"gcmParameterSpecInit: iv cannot be null")
	}

	fullIv := object.GoByteArrayFromJavaByteArray(ivObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	var iv []byte
	if len(params) == 3 {
		// GCMParameterSpec(int tLen, byte[] src)
		iv = slices.Clone(fullIv)
	} else if len(params) == 5 {
		// GCMParameterSpec(int tLen, byte[] src, int offset, int len)
		offset := params[3].(int64)
		length := params[4].(int64)

		if offset < 0 || length < 0 || int(offset+length) > len(fullIv) {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"gcmParameterSpecInit: invalid offset or length")
		}
		iv = slices.Clone(fullIv[offset : offset+length])
	} else {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecInit: wrong number of parameters")
	}

	if len(iv) == 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"gcmParameterSpecInit: iv cannot be empty")
	}

	self.FieldTable["iv"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: iv,
	}
	self.FieldTable["tLen"] = object.Field{
		Ftype:  types.Int,
		Fvalue: tLen,
	}

	return nil
}
