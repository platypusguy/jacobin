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

func Load_Security_Interfaces_RSA_Keys() {

	// =====RSAKey =====

	ghelpers.MethodSignatures["java/security/interfaces/RSAKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
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

	// =====RSAPrivateKey =====

	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateKey.getPrivateExponent()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rsaprivateGetExponent,
		}

	// =====RSAPublicKey =====

	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.getModulus()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rsaKeyGetModulus,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;"] =
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

	// The RSAKey super-interface is shared by both RSAPublicKey and RSAPrivateKey.
	// Support extracting modulus from either kind of underlying Go key.
	if rsapubkey, ok := thisObj.FieldTable["value"].Fvalue.(*rsa.PublicKey); ok {
		bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, rsapubkey.N)
		return bigint
	}
	if rsaprivkey, ok := thisObj.FieldTable["value"].Fvalue.(*rsa.PrivateKey); ok {
		bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, rsaprivkey.N)
		return bigint
	}

	return ghelpers.GetGErrBlk(
		excNames.VirtualMachineError,
		"rsaKeyGetModulus: RSA key field is the wrong type",
	)
}

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
			excNames.VirtualMachineError,
			"rsaprivateGetExponent: RSA private key field is the wrong type",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, rsaprvkey.D)

	return bigint
}
