/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Private License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"fmt"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
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

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.getG()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetG,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.getP()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetP,
		}

	ghelpers.MethodSignatures["java/security/interfaces/DSAParams.getQ()Ljava/math/BigInteger;"] =
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

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.getG()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetG,
		}

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.getP()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  dsaParamsGetP,
		}

	ghelpers.MethodSignatures["java/security/spec/DSAParameterSpec.getQ()Ljava/math/BigInteger;"] =
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

	ghelpers.MethodSignatures["java/security/interfaces/DSAPrivateKey.getX()Ljava/math/BigInteger;"] =
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
			fmt.Sprintf("dsaPrivateGetX: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"dsaPrivateGetX: `this` is not an Object",
		)
	}

	dsaprvkeyObj, ok := thisObj.FieldTable["value"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaPrivateGetX: DSA private key extraction failed (not an Object)",
		)
	}

	return dsaprvkeyObj
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

	dsaParamsGObj, ok := thisObj.FieldTable["g"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaParamsGetG: DSA field G extraction failed (not an Object)",
		)
	}

	return dsaParamsGObj
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

	dsaParamsPObj, ok := thisObj.FieldTable["p"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaParamsGetP: DSA field P extraction failed (not an Object)",
		)
	}

	return dsaParamsPObj
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

	dsaParamsQObj, ok := thisObj.FieldTable["q"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.VirtualMachineError,
			"dsaParamsGetQ: DSA field Q extraction failed (not an Object)",
		)
	}

	return dsaParamsQObj
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

	dsapubkeyObj, ok := thisObj.FieldTable["value"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"dsaPublicKeyGetY: DSA public key extraction failed (not an Object)",
		)
	}

	return dsapubkeyObj
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
			"dsaKeysGetParams: DSA parameters extraction failed",
		)
	}

	return dsaParams
}
