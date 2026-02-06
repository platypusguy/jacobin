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
// Loader for EllipticCurve
// ---------------------------------------------------------
func Load_EllipticCurve() {
	// Constructor
	ghelpers.MethodSignatures["java/security/spec/EllipticCurve.<init>(Ljava/security/spec/ECField;Ljava/math/BigInteger;Ljava/math/BigInteger;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ellipticCurveInit,
		}

	// Getter for field
	ghelpers.MethodSignatures["java/security/spec/EllipticCurve.getField()Ljava/security/spec/ECField;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ellipticCurveGetField,
		}

	// Getter for a
	ghelpers.MethodSignatures["java/security/spec/EllipticCurve.getA()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ellipticCurveGetA,
		}

	// Getter for b
	ghelpers.MethodSignatures["java/security/spec/EllipticCurve.getB()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ellipticCurveGetB,
		}
}

// ---------------------------------------------------------
// EllipticCurve G functions
// ---------------------------------------------------------

// Constructor: EllipticCurve(ECField field, BigInteger a, BigInteger b)
func ellipticCurveInit(params []any) any {
	if len(params) != 4 { // this + field + a + b
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ellipticCurveInit: expected 3 parameters (field, a, b), got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveInit: this is not an Object",
		)
	}

	fieldObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveInit: param field is not an Object",
		)
	}

	aObj, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveInit: param a is not an Object",
		)
	}

	bObj, ok := params[3].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveInit: param b is not an Object",
		)
	}

	// Populate thisObj FieldTable
	thisObj.FieldTable = map[string]object.Field{
		"field": {Ftype: types.Ref, Fvalue: fieldObj},
		"a":     {Ftype: types.Ref, Fvalue: aObj},
		"b":     {Ftype: types.Ref, Fvalue: bObj},
	}

	return nil // <init> always returns void
}

// Getter: getField()
func ellipticCurveGetField(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ellipticCurveGetField: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveGetField: this is not an Object",
		)
	}
	fieldObj, ok := thisObj.FieldTable["field"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ellipticCurveGetField: field is missing or invalid",
		)
	}
	return fieldObj
}

// Getter: getA()
func ellipticCurveGetA(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ellipticCurveGetA: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveGetA: this is not an Object",
		)
	}
	aObj, ok := thisObj.FieldTable["a"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ellipticCurveGetA: a field is missing or invalid",
		)
	}
	return aObj
}

// Getter: getB()
func ellipticCurveGetB(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ellipticCurveGetB: expected 1 parameter (this), got %d", len(params)),
		)
	}
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ellipticCurveGetB: this is not an Object",
		)
	}
	bObj, ok := thisObj.FieldTable["b"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.NullPointerException,
			"ellipticCurveGetB: b field is missing or invalid",
		)
	}
	return bObj
}
