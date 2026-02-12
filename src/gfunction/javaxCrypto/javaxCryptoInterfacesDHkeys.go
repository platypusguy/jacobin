/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Private License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"fmt"
	"math/big"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Placeholder: use simple DH big.Int values

func Load_Crypto_Interfaces_DH_Keys() {

	// =====DHKey =====

	ghelpers.MethodSignatures["java/security/interfaces/DHKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHKey.getParams()Ljavax/crypto/spec/DHParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnNull,
		}

	// =====DHPrivateKey =====

	ghelpers.MethodSignatures["java/security/interfaces/DHPrivateKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHPrivateKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHPrivateKey.getX()()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dhPrivateGetX,
		}

	// =====DHPublicKey =====

	ghelpers.MethodSignatures["java/security/interfaces/DHPublicKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHPublicKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHPublicKey.getY()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dhPublicKeyGetY,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DHPublicKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnNull,
		}

}

// ---------------------------------------------------------
// G functions
// ---------------------------------------------------------

func dhPrivateGetX(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dhPrivateGetExponent: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dhPrivateGetExponent: `this` is not an Object",
		)
	}

	dhprvkey, ok := thisObj.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dhPrivateGetX: DH private key extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dhprvkey)

	return bigint
}

func dhPublicKeyGetY(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dhPublicKeyGetY: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dhPublicKeyGetY: `this` is not an Object",
		)
	}

	dhpubkey, ok := thisObj.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"dhPublicKeyGetY: DH public key extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dhpubkey)

	return bigint
}
