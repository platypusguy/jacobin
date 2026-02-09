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

func Load_Security_KeyPair() {

	ghelpers.MethodSignatures["java/security/KeyPair.<init>(Ljava/security/PublicKey;Ljava/security/PrivateKey;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  keypairInit,
		}

	ghelpers.MethodSignatures["java/security/KeyPair.getPublic()Ljava/security/PublicKey;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keypairGetPublic,
		}

	ghelpers.MethodSignatures["java/security/KeyPair.getPrivate()Ljava/security/PrivateKey;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keypairGetPrivate,
		}

}

func keypairInit(params []any) any {
	if len(params) != 3 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("keypairInit: expected 3 params, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairInit: param[0] is not an Object")
	}

	pub, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairInit: param[1] is not an Object")
	}

	priv, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairInit: param[2] is not an Object")
	}

	thisObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}
	thisObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

	return nil
}

func keypairGetPublic(params []any) any {
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairGetPublic: param[0] is not an Object")
	}
	return thisObj.FieldTable["public"].Fvalue
}

func keypairGetPrivate(params []any) any {
	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairGetPrivate: param[0] is not an Object")
	}
	return thisObj.FieldTable["private"].Fvalue
}
