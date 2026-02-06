/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Private License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/rsa"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// ---------------------------------------------------------
// Loader for RSAPrivateKey
// ---------------------------------------------------------
func Load_Security_Interfaces_RSAPrivateKey() {

	// Interface constructor placeholder
	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// Interface constructor placeholder
	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateKey.getPrivateExponent()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rsaprivateGetExponent,
		}
}

// ---------------------------------------------------------
// G functions
// ---------------------------------------------------------

func rsaprivateGetExponent(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("rsaprivateGetExponent: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"rsaprivateGetExponent: `this` is not an Object",
		)
	}

	rsaprvkey, ok := thisObj.FieldTable["value"].Fvalue.(*rsa.PrivateKey)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"rsaprivateGetExponent: RSA public key extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, rsaprvkey.D)

	return bigint
}
