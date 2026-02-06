/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// ---------------------------------------------------------
// Loader for RSAKey
// ---------------------------------------------------------
func Load_Security_Interfaces_RSAKey() {

	// Interface constructor placeholder
	ghelpers.MethodSignatures["java/security/interfaces/RSAKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// Interface constructor placeholder
	ghelpers.MethodSignatures["java/security/interfaces/RSAKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAKey.getModulus()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rsaKeyGetModulus,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnNull,
		}
}

// ---------------------------------------------------------
// G functions
// ---------------------------------------------------------

func rsaKeyGetModulus(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("rsaKeyGetModulus: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"rsaKeyGetModulus: `this` is not an Object",
		)
	}

	modField, ok := thisObj.FieldTable["modulus"]
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"rsaKeyGetModulus: modulus field not set",
		)
	}

	if modField.Ftype != types.BigInteger {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"rsaKeyGetModulus: modulus is not BigInteger",
		)
	}

	return modField.Fvalue
}
