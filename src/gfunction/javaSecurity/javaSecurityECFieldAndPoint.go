/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
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
// Loader for ECField, ECFieldFp, and ECPoint
// ---------------------------------------------------------
func Load_ECFieldAndPoint() {
	// --------------------------
	// ECField
	// --------------------------
	ghelpers.MethodSignatures["java/security/spec/ECField.<init>(I)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ecFieldInit}
	ghelpers.MethodSignatures["java/security/spec/ECField.getFieldSize()I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ecFieldGetFieldSize}

	// --------------------------
	// ECFieldFp
	// --------------------------
	ghelpers.MethodSignatures["java/security/spec/ECFieldFp.<init>(Ljava/math/BigInteger;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ecFieldFpInit}
	ghelpers.MethodSignatures["java/security/spec/ECFieldFp.getP()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ecFieldFpGetP}

	// --------------------------
	// ECPoint
	// --------------------------
	ghelpers.MethodSignatures["java/security/spec/ECPoint.<init>(Ljava/math/BigInteger;Ljava/math/BigInteger;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ecPointInit}
	ghelpers.MethodSignatures["java/security/spec/ECPoint.getAffineX()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ecPointGetAffineX}
	ghelpers.MethodSignatures["java/security/spec/ECPoint.getAffineY()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ecPointGetAffineY}
}

// ---------------------------------------------------------
// ECField G functions
// --------------------------
func ecFieldInit(params []any) any {
	if len(params) != 1+1 { // this + fieldSize
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecFieldInit: expected 1 parameter (fieldSize), got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecFieldInit: this is not an Object",
		)
	}

	fieldSize, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecFieldInit: fieldSize is not int",
		)
	}

	thisObj.FieldTable = map[string]object.Field{
		"fieldSize": {Ftype: types.Int, Fvalue: fieldSize},
	}

	return nil // <init> always returns void
}

func ecFieldGetFieldSize(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecFieldGetFieldSize: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecFieldGetFieldSize: this is not an Object",
		)
	}
	size, ok := thisObj.FieldTable["fieldSize"].Fvalue.(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecFieldGetFieldSize: fieldSize missing or invalid",
		)
	}
	return size
}

// ---------------------------------------------------------
// ECFieldFp G functions
// --------------------------
func ecFieldFpInit(params []any) any {
	if len(params) != 2 { // this + p
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecFieldFpInit: expected 1 parameter (p), got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecFieldFpInit: this is not an Object",
		)
	}

	pObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecFieldFpInit: param is not BigInteger",
		)
	}

	thisObj.FieldTable = map[string]object.Field{
		"p": {Ftype: types.BigInteger, Fvalue: pObj},
	}

	return nil // <init> returns void
}

func ecFieldFpGetP(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecFieldFpGetP: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecFieldFpGetP: this is not an Object",
		)
	}
	pObj, ok := thisObj.FieldTable["p"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecFieldFpGetP: p field missing or invalid",
		)
	}
	return pObj
}

// ---------------------------------------------------------
// ECPoint G functions
// --------------------------
func ecPointInit(params []any) any {
	if len(params) != 3 { // this + x + y
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPointInit: expected 2 parameters (x, y), got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPointInit: this is not an Object",
		)
	}

	xObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPointInit: param 0 is not BigInteger",
		)
	}

	yObj, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPointInit: param 1 is not BigInteger",
		)
	}

	// Populate the thisObj FieldTable
	thisObj.FieldTable = map[string]object.Field{
		"x": {Ftype: types.BigInteger, Fvalue: xObj},
		"y": {Ftype: types.BigInteger, Fvalue: yObj},
	}

	return nil // <init> always returns void
}

func ecPointGetAffineX(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPointGetAffineX: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPointGetAffineX: this is not an Object",
		)
	}
	xObj, ok := thisObj.FieldTable["x"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecPointGetAffineX: x field missing or invalid",
		)
	}
	return xObj
}

func ecPointGetAffineY(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPointGetAffineY: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPointGetAffineY: this is not an Object",
		)
	}
	yObj, ok := thisObj.FieldTable["y"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ecPointGetAffineY: y field missing or invalid",
		)
	}
	return yObj
}
