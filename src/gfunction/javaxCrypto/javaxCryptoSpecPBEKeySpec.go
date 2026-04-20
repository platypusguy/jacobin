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
)

func Load_Crypto_Spec_PBEKeySpec() {
	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.<init>([C)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  pbeKeySpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.<init>([C[BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  pbeKeySpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.clearPassword()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeKeySpecClearPassword,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.getIterationCount()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeKeySpecGetIterationCount,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.getKeyLength()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeKeySpecGetKeyLength,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.getPassword()[C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeKeySpecGetPassword,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/PBEKeySpec.getSalt()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  pbeKeySpecGetSalt,
		}

}

func pbeKeySpecClearPassword(params []any) any {
	this := params[0].(*object.Object)
	passwordVal := this.FieldTable["password"].Fvalue
	if !object.IsNull(passwordVal) {
		passwordObj := passwordVal.(*object.Object)
		passwordChars := passwordObj.FieldTable["value"].Fvalue.([]int64)
		for i := range passwordChars {
			passwordChars[i] = 0
		}
	}
	this.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: object.Null}
	return nil
}

func pbeKeySpecGetIterationCount(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["iterationCount"].Fvalue
}

func pbeKeySpecGetKeyLength(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["keyLength"].Fvalue
}

func pbeKeySpecGetPassword(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["password"].Fvalue
}

func pbeKeySpecGetSalt(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["salt"].Fvalue
}

func pbeKeySpecInit(params []any) any {
	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "pbeKeySpecInit: invalid 'this'")
	}

	if len(params) == 2 {
		// PBEKeySpec(char[] password)
		password := params[1]
		this.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: password}
		this.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: object.Null}
		this.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		this.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	} else if len(params) == 5 {
		// PBEKeySpec(char[] password, byte[] salt, int iterationCount, int keyLength)
		password := params[1]
		salt := params[2]
		iterationCount := params[3].(int64)
		keyLength := params[4].(int64)

		if iterationCount < 0 {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "pbeKeySpecInit: iterationCount must be >= 0")
		}
		if keyLength < 0 {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "pbeKeySpecInit: keyLength must be >= 0")
		}

		this.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: password}
		this.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: salt}
		this.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterationCount}
		this.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: keyLength}
	}

	return nil
}
