/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Private License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"fmt"
	"math/big"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Placeholder: use simple DSA big.Int values

func Load_Security_Interfaces_DSA_Keys() {

	// ===== DSAKey =====

	ghelpers.MethodSignatures["java/security/interfaces/DSAKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAKey.getParams()Ljava/security/interfaces/DSAParams;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaKeysGetParams,
		}

	// ===== DSAParams =====

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.getG()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetG,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.getP()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetP,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.getG()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetQ,
		}

	// ===== DSAParameterSpec =====

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.getG()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetG,
		}

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.getP()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetP,
		}

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.getG()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetQ,
		}

	// ===== DSAPrivateKey =====

	ghelpers.MethodSignatures["java/security/interfaces/DSAPrivateKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAPrivateKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAPrivateKey.getX()()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaPrivateGetX,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAPrivateKey.getParams()Ljava/security/interfaces/DSAParams;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaKeysGetParams,
		}

	// ===== DSAPublicKey =====

	ghelpers.MethodSignatures["java/security/interfaces/DSAPublicKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAPublicKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAPublicKey.getY()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaPublicKeyGetY,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAPublicKey.getParams()Ljava/security/interfaces/DSAParams;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaKeysGetParams,
		}

}

// ---------------------------------------------------------
// Member functions
// ---------------------------------------------------------

func dsaPrivateGetX(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dsaPrivateGetExponent: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaPrivateGetExponent: `this` is not an Object",
		)
	}

	dsaprvkey, ok := thisObj.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaPrivateGetX: DSA private key extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dsaprvkey)

	return bigint
}

func dsaParamsGetG(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dsaParamsGetG: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaParamsGetG: `this` is not an Object",
		)
	}

	dsaParamsG, ok := thisObj.FieldTable["g"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaParamsGetG: DSA field G extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dsaParamsG)

	return bigint
}

func dsaParamsGetP(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dsaParamsGetP: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaParamsGetP: `this` is not an Object",
		)
	}

	dsaParamsP, ok := thisObj.FieldTable["p"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaParamsGetP: DSA field P extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dsaParamsP)

	return bigint
}
func dsaParamsGetQ(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dsaParamsGetQ: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaParamsGetQ: `this` is not an Object",
		)
	}

	dsaParamsQ, ok := thisObj.FieldTable["q"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaParamsGetQ: DSA field Q extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dsaParamsQ)

	return bigint
}

func dsaPublicKeyGetY(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dsaPublicKeyGetY: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaPublicKeyGetY: `this` is not an Object",
		)
	}

	dsapubkey, ok := thisObj.FieldTable["value"].Fvalue.(*big.Int)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"dsaPublicKeyGetY: DSA public key extraction failed",
		)
	}

	bigint := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, dsapubkey)

	return bigint
}

func dsaKeysGetParams(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("dsaKeysGetParams: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaKeysGetParams: `this` is not an Object",
		)
	}

	dsaParams, ok := thisObj.FieldTable["params"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"dsaKeysGetParams: DSA public key extraction failed",
		)
	}

	return dsaParams
}
