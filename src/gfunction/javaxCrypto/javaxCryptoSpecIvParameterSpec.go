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

func Load_Crypto_Spec_IvParameterSpec() {
	ghelpers.MethodSignatures["javax/crypto/spec/IvParameterSpec.<init>([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ivParameterSpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/IvParameterSpec.<init>([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ivParameterSpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/IvParameterSpec.getIV()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ivParameterSpecGetIV,
		}
}

func ivParameterSpecGetIV(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"ivParameterSpecGetIV: missing 'this'")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"ivParameterSpecGetIV: 'this' is not an object")
	}

	iv, ok := self.FieldTable["iv"].Fvalue.([]byte)
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	// Returns a copy of the IV
	jBytes := object.JavaByteArrayFromGoByteArray(slices.Clone(iv))
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}

func ivParameterSpecInit(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"ivParameterSpecInit: insufficient parameters")
	}

	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"ivParameterSpecInit: 'this' is not an object")
	}

	ivObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(ivObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException,
			"ivParameterSpecInit: iv cannot be null")
	}

	fullIv := object.GoByteArrayFromJavaByteArray(ivObj.FieldTable["value"].Fvalue.([]types.JavaByte))

	var iv []byte
	if len(params) == 2 {
		// IvParameterSpec(byte[] iv)
		iv = slices.Clone(fullIv)
	} else if len(params) == 4 {
		// IvParameterSpec(byte[] iv, int offset, int len)
		offset := params[2].(int64)
		length := params[3].(int64)

		if offset < 0 || length < 0 || int(offset+length) > len(fullIv) {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"ivParameterSpecInit: invalid offset or length")
		}
		iv = slices.Clone(fullIv[offset : offset+length])
	} else {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"ivParameterSpecInit: wrong number of parameters")
	}

	if len(iv) == 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"ivParameterSpecInit: iv cannot be empty")
	}

	self.FieldTable["iv"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: iv,
	}

	return nil
}
