/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

/*
This file represents the sole Security Provider Service permitted by Jacobin: The Go Run-time
*/

// Load_Security_Provider_Service initializes java/security/Provider$Service methods
func Load_Security_Provider_Service() {

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".<init>(Ljava/security/Provider;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 5, // provider, type, algorithm, className, aliases
			GFunction:  securityProvSvcInit,
		}

	// ---------- Member Functions ----------
	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProvSvcGetAlgorithm}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".getAliases()Ljava/util/List;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProvSvcGetAliases}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".getAttribute(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: securityProvSvcGetAttribute}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".getClassName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProvSvcGetClassName}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProvSvcGetProvider}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".getType()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProvSvcGetType}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".newInstance(Ljava/lang/Object[])Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures[types.ClassNameSecurityProviderService+".toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProvSvcToString}
}

// ----------------------- Constructor -----------------------
func securityProvSvcInit(params []any) any {
	// params: provider, type, algorithm, className, aliases
	this := params[0].(*object.Object)

	if params[1] == nil || params[2] == nil || params[3] == nil || params[4] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "securityProvSvcInit: one or more arguments are null")
	}

	provider := params[1].(*object.Object)
	typeStr := strings.TrimSpace(object.GoStringFromStringObject(params[2].(*object.Object)))
	algorithmStr := strings.TrimSpace(object.GoStringFromStringObject(params[3].(*object.Object)))
	classNameStr := strings.TrimSpace(object.GoStringFromStringObject(params[4].(*object.Object)))
	aliasesArray := []string{}
	if len(params) > 5 {
		for _, a := range params[5].([]*object.Object) {
			if a != nil {
				aliasesArray = append(aliasesArray, strings.TrimSpace(object.GoStringFromStringObject(a)))
			}
		}
	}
	aliasObjArray := object.StringObjectArrayFromGoStringArray(aliasesArray)

	this.FieldTable["provider"] = object.Field{Ftype: types.Ref, Fvalue: provider}
	this.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(typeStr)}
	this.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(algorithmStr)}
	this.FieldTable["className"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(classNameStr)}
	this.FieldTable["aliases"] = object.Field{Ftype: types.StringArrayClassName, Fvalue: aliasObjArray}
	var attributes = map[string]*object.Object{}
	attributes["ImplementedIn"] = object.StringObjectFromGoString("Software")
	attributes["blockSize"] = object.StringObjectFromGoString("null")
	this.FieldTable["attributes"] = object.Field{Ftype: types.Map, Fvalue: attributes}

	return nil
}

// ----------------------- Member Functions -----------------------
func securityProvSvcGetProvider(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["provider"].Fvalue.(*object.Object)
}

func securityProvSvcGetType(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["type"].Fvalue.(*object.Object)
}

func securityProvSvcGetAlgorithm(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["algorithm"].Fvalue.(*object.Object)
}

func securityProvSvcGetClassName(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["className"].Fvalue.(*object.Object)
}

func securityProvSvcGetAliases(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["aliases"].Fvalue
}

func securityProvSvcGetAttribute(params []any) any {
	this, ok := params[0].(*object.Object)
	if !ok {
		return object.Null
	}
	nameObj, ok := params[1].(*object.Object)
	if !ok {
		return object.Null
	}
	name := object.GoStringFromStringObject(nameObj)

	attrs := this.FieldTable["attributes"].Fvalue.(map[string]*object.Object)
	if val, ok := attrs[name]; ok {
		return val
	}
	return object.Null // match Hotspot: missing attribute returns null
}

func securityProvSvcToString(params []any) any {
	this := params[0].(*object.Object)
	typeStr := object.GoStringFromStringObject(this.FieldTable["type"].Fvalue.(*object.Object))
	algStr := object.GoStringFromStringObject(this.FieldTable["algorithm"].Fvalue.(*object.Object))
	return object.StringObjectFromGoString(typeStr + "." + algStr)
}
