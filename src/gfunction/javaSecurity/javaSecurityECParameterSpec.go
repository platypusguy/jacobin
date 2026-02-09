/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/elliptic"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
)

func Load_ECParameterSpec() {

	// Constructor from curve name string
	ghelpers.MethodSignatures["java/security/spec/ECGenParameterSpec.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ecparameterGenSpecInitString,
		}

	// Constructor: ECParameterSpec(EllipticCurve curve, ECPoint g, BigInteger n, int h)
	ghelpers.MethodSignatures["java/security/spec/ECParameterSpec.<init>(Ljava/security/spec/EllipticCurve;Ljava/security/spec/ECPoint;Ljava/math/BigInteger;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ecparameterSpecInit,
		}

	// ---------------------------------------------------------
	// Public API getters
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/spec/ECParameterSpec.getCurve()Ljava/security/spec/EllipticCurve;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ecparameterSpecGetCurve,
		}

	ghelpers.MethodSignatures["java/security/spec/ECParameterSpec.getGenerator()Ljava/security/spec/ECPoint;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ecparameterSpecGetGenerator,
		}

	ghelpers.MethodSignatures["java/security/spec/ECParameterSpec.getOrder()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ecparameterSpecGetOrder,
		}

	ghelpers.MethodSignatures["java/security/spec/ECParameterSpec.getCofactor()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ecparameterSpecGetCofactor,
		}
}

// Constructor
// ecparameterGenSpecInitString constructs an ECParameterSpec from a curve name string
func ecparameterGenSpecInitString(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecparameterSpecInitString: expected 1 parameter, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInitString: this is not an Object",
		)
	}

	curveNameObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInitString: curveName is not a String object",
		)
	}

	curveName := object.GoStringFromStringObject(curveNameObj)

	// Map curve name to Go's elliptic curve
	var curve elliptic.Curve
	switch curveName {
	case "secp224r1", "P-224":
		curve = elliptic.P224()
	case "secp256r1", "P-256":
		curve = elliptic.P256()
	case "secp384r1", "P-384":
		curve = elliptic.P384()
	case "secp521r1", "P-521":
		curve = elliptic.P521()
	default:
		return ghelpers.GetGErrBlk(
			excNames.InvalidAlgorithmParameterException,
			fmt.Sprintf("ecparameterSpecInitString: unsupported curve name: %s", curveName),
		)
	}

	curveParams := curve.Params()

	// --- Construct EllipticCurve object ---
	curveObj := NewGoRuntimeService("EllipticCurve", "EC", types.ClassNameEllipticCurve)

	// Wrap P (prime field)
	curveObj.FieldTable["p"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, curveParams.P),
	}

	// Wrap A (coefficient, -3 for NIST curves)
	curveObj.FieldTable["a"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(-3)),
	}

	// Wrap B (coefficient)
	curveObj.FieldTable["b"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, curveParams.B),
	}

	// --- Create generator point (G) ---
	generatorObj := NewGoRuntimeService("ECPoint", "", types.ClassNameECPoint)

	// Wrap Gx
	generatorObj.FieldTable["x"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, curveParams.Gx),
	}

	// Wrap Gy
	generatorObj.FieldTable["y"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, curveParams.Gy),
	}

	curveObj.FieldTable["generator"] = object.Field{Ftype: types.ECPoint, Fvalue: generatorObj}

	// --- Populate ECParameterSpec fields ---
	thisObj.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: curveNameObj}
	thisObj.FieldTable["curve"] = object.Field{Ftype: types.Ref, Fvalue: curveObj}
	thisObj.FieldTable["g"] = object.Field{Ftype: types.Ref, Fvalue: generatorObj}

	// Wrap N (order)
	thisObj.FieldTable["n"] = object.Field{
		Ftype:  types.BigInteger,
		Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, curveParams.N),
	}

	thisObj.FieldTable["h"] = object.Field{Ftype: types.Int, Fvalue: int64(1)} // cofactor

	return nil
}

// Constructor
// ECParameterSpec(EllipticCurve curve, ECPoint g, BigInteger n, int h)
func ecparameterSpecInit(params []any) any {
	if len(params) != 5 { // this + curve + g + n + h
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecparameterSpecInit: expected 4 parameters, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInit: this is not an Object",
		)
	}

	curveObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInit: param curve is not an Object",
		)
	}

	gObj, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInit: param g is not an Object",
		)
	}

	nObj, ok := params[3].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInit: param n is not an Object",
		)
	}

	hVal, ok := params[4].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecInit: param h is not int",
		)
	}

	// Populate thisObj fields
	thisObj.FieldTable = map[string]object.Field{
		"curve": {Ftype: types.Ref, Fvalue: curveObj},
		"g":     {Ftype: types.Ref, Fvalue: gObj},
		"n":     {Ftype: types.Ref, Fvalue: nObj},
		"h":     {Ftype: types.Int, Fvalue: hVal},
	}

	return nil // <init> returns void
}

// ---------------------------------------------------------
// Getter: getCurve
// ---------------------------------------------------------
func ecparameterSpecGetCurve(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecparameterSpecGetCurve: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecGetCurve: this is not an Object",
		)
	}
	curveObj, ok := thisObj.FieldTable["curve"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecparameterSpecGetCurve: curve field is missing or invalid",
		)
	}
	return curveObj
}

// ---------------------------------------------------------
// Getter: getGenerator
// ---------------------------------------------------------
func ecparameterSpecGetGenerator(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecparameterSpecGetGenerator: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecGetGenerator: this is not an Object",
		)
	}
	genObj, ok := thisObj.FieldTable["g"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecparameterSpecGetGenerator: generator field is missing or invalid",
		)
	}
	return genObj
}

// ---------------------------------------------------------
// Getter: getOrder
// ---------------------------------------------------------
func ecparameterSpecGetOrder(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecparameterSpecGetOrder: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecGetOrder: this is not an Object",
		)
	}
	orderObj, ok := thisObj.FieldTable["n"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecparameterSpecGetOrder: order field is missing or invalid",
		)
	}
	return orderObj
}

// ---------------------------------------------------------
// Getter: getCofactor
// ---------------------------------------------------------
func ecparameterSpecGetCofactor(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecparameterSpecGetCofactor: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecparameterSpecGetCofactor: this is not an Object",
		)
	}
	cofactor, ok := thisObj.FieldTable["h"].Fvalue.(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecparameterSpecGetCofactor: cofactor field is missing or invalid",
		)
	}
	return cofactor
}
