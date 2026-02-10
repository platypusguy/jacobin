/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/ecdsa"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Load_ECKeys registers EC interfaces and concrete classes
func Load_Security_Interfaces_EC_Keys() {
	// ---------------------------------------------------------
	// <clinit> and <init> with no arguments
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/interfaces/ECKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/ECKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/ECPrivateKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/ECPrivateKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/ECPublicKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/interfaces/ECPublicKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapKeyPairGeneration,
		}

	ghelpers.MethodSignatures["java/security/interfaces/ECPrivateKey.getParams()Ljava/security/spec/ECParameterSpec;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ecPrivateKeyGetParams}

	ghelpers.MethodSignatures["java/security/interfaces/ECPrivateKey.getS()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ecPrivateKeyGetS}

	ghelpers.MethodSignatures["java/security/interfaces/ECPublicKey.getParams()Ljava/security/spec/ECParameterSpec;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ecPublicKeyGetParams}

	ghelpers.MethodSignatures["java/security/interfaces/ECPublicKey.getW()Ljava/security/spec/ECPoint;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ecPublicKeyGetW}

}

// === ECPrivateKey ===

func ecPrivateKeyGetParams(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPrivateKeyGetParams: expected 0 parameters, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyGetParams: this is not an Object",
		)
	}

	specObj, exists := thisObj.FieldTable["params"].Fvalue.(*object.Object)
	if !exists {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyGetParams: params field missing",
		)
	}

	return specObj

}

func ecPrivateKeyGetS(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPrivateKeyGetS: expected `this` object only, got %d params", len(params)),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyGetS: `this` is not an Object",
		)
	}

	ecprivkey, ok := thisObj.FieldTable["value"].Fvalue.(*ecdsa.PrivateKey)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"ecPrivateKeyGetS: EC private key extraction failed",
		)
	}

	bigintObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, ecprivkey.D)

	return bigintObj
}

// === ECPublicKey ===

func ecPublicKeyGetParams(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPublicKeyGetParams: expected 0 parameters, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyGetParams: this is not an Object",
		)
	}

	specObj, exists := thisObj.FieldTable["params"].Fvalue.(*object.Object)
	if !exists {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyGetParams: params field missing",
		)
	}

	return specObj
}

func ecPublicKeyGetW(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPublicKeyGetW: expected 0 parameters, got %d", len(params)-1),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyGetW: this is not an Object",
		)
	}

	// Try "w" field first (manually populated)
	if pointObj, exists := thisObj.FieldTable["w"].Fvalue.(*object.Object); exists {
		return pointObj
	}

	// Otherwise extract from "value" (*ecdsa.PublicKey)
	ecpubkey, ok := thisObj.FieldTable["value"].Fvalue.(*ecdsa.PublicKey)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"ecPublicKeyGetW: EC public key extraction failed",
		)
	}

	// Create ECPoint object
	pointObj := NewGoRuntimeService("ECPoint", "", types.ClassNameECPoint)
	pointObj.FieldTable["x"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, ecpubkey.X),
	}
	pointObj.FieldTable["y"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, ecpubkey.Y),
	}

	return pointObj
}
