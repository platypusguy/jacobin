/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

// Loader
func Load_Security_Spec_NamedParameterSpec() {

	ghelpers.MethodSignatures["java/security/spec/NamedParameterSpec.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  namedParameterSpecClinit,
		}

	ghelpers.MethodSignatures["java/security/spec/NamedParameterSpec.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  namedParameterSpecInit,
		}

	ghelpers.MethodSignatures["java/security/spec/NamedParameterSpec.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  namedParameterSpecGetName,
		}
}

func namedParameterSpecClinit([]interface{}) interface{} {

	className := "java/security/spec/NamedParameterSpec"

	makeSpec := func(name string) *object.Object {
		specObj := object.MakeEmptyObjectWithClassName(&className)
		nameObj := object.StringObjectFromGoString(name)

		// call the constructor <init>(String)
		ret := namedParameterSpecInit([]interface{}{specObj, nameObj})
		if ret != nil {
			// fail early if constructor fails
			errBlk := *ret.(*ghelpers.GErrBlk)
			errMsg := fmt.Sprintf("namedParameterSpecClinit/makeSpec(name=%s): %s", name, errBlk.ErrMsg)
			exceptions.MinimalAbort(errBlk.ExceptionType, errMsg)
		}
		return specObj
	}

	// initialize static fields
	_ = statics.AddStatic(className+".X25519", statics.Static{Type: types.Ref, Value: makeSpec("X25519")})
	_ = statics.AddStatic(className+".X448", statics.Static{Type: types.Ref, Value: makeSpec("X448")})
	_ = statics.AddStatic(className+".ED25519", statics.Static{Type: types.Ref, Value: makeSpec("Ed25519")})
	_ = statics.AddStatic(className+".ED448", statics.Static{Type: types.Ref, Value: makeSpec("Ed448")})

	return nil
}

// NamedParameterSpec.<init>(String)
func namedParameterSpecInit(params []interface{}) interface{} {

	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "namedParameterSpecInit: expected 1 argument")
	}

	selfObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "namedParameterSpecInit: invalid self object")
	}

	nameParam, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "namedParameterSpecInit: invalid name object")
	}
	if !object.IsStringObject(nameParam) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "namedParameterSpecInit: expected name to be a String")
	}

	// store name in FieldTable
	selfObj.FieldTable["name"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: nameParam,
	}

	return nil // constructor returns void
}

// NamedParameterSpec.getName()Ljava/lang/String;
func namedParameterSpecGetName(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "namedParameterSpecGetName: expected 0 arguments")
	}

	selfObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "namedParameterSpecGetName: invalid self object")
	}

	fieldEntry, exists := selfObj.FieldTable["name"]
	if !exists {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"namedParameterSpecGetName: name field not set",
		)
	}

	nameObj, ok := fieldEntry.Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"namedParameterSpecGetName: name field has invalid type",
		)
	}

	return nameObj
}
