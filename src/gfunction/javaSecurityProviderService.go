package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

/*
This file represents the sole Security Provider Service permitted by Jacobin: The Go Run-time
*/

// Load_Security_Provider_Service initializes java/security/Provider$Service methods
func Load_Security_Provider_Service() {

	MethodSignatures["java/security/Provider$Service.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/security/Provider$Service.<init>(Ljava/security/Provider;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 5, // provider, type, algorithm, className, aliases
			GFunction:  securityProvSvcInit,
		}

	// ---------- Member Functions ----------
	MethodSignatures["java/security/Provider$Service.getAlgorithm()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: securityProvSvcGetAlgorithm}

	MethodSignatures["java/security/Provider$Service.getAliases()Ljava/util/List;"] =
		GMeth{ParamSlots: 0, GFunction: securityProvSvcGetAliases}

	MethodSignatures["java/security/Provider$Service.getAttribute(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{ParamSlots: 1, GFunction: securityProvSvcGetAttribute}

	MethodSignatures["java/security/Provider$Service.getClassName()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: securityProvSvcGetClassName}

	MethodSignatures["java/security/Provider$Service.getProvider()Ljava/security/Provider;"] =
		GMeth{ParamSlots: 0, GFunction: securityProvSvcGetProvider}

	MethodSignatures["java/security/Provider$Service.getType()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: securityProvSvcGetType}

	MethodSignatures["java/security/Provider$Service.newInstance(Ljava/lang/Object[])Ljava/lang/Object;"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/security/Provider$Service.toString()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: securityProvSvcToString}
}

// ----------------------- Constructor -----------------------
func securityProvSvcInit(params []any) any {
	this := params[0].(*object.Object)

	if params[1] == nil || params[2] == nil || params[3] == nil || params[4] == nil {
		return getGErrBlk(excNames.NullPointerException, "securityProvSvcInit: one or more arguments are null")
	}

	provider := params[1].(*object.Object)
	typeStr := strings.TrimSpace(object.GoStringFromStringObject(params[2].(*object.Object)))
	algorithmStr := strings.TrimSpace(object.GoStringFromStringObject(params[3].(*object.Object)))
	classNameStr := strings.TrimSpace(object.GoStringFromStringObject(params[4].(*object.Object)))
	aliasGoArray := []string{}
	if len(params) > 5 {
		for _, a := range params[5].([]*object.Object) {
			if a != nil {
				aliasGoArray = append(aliasGoArray, strings.TrimSpace(object.GoStringFromStringObject(a)))
			}
		}
	}
	aliasArray := object.StringObjectArrayFromGoStringArray(aliasGoArray)

	this.FieldTable["provider"] = object.Field{Ftype: types.Ref, Fvalue: provider}
	this.FieldTable["type"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(typeStr)}
	this.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(algorithmStr)}
	this.FieldTable["className"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(classNameStr)}
	this.FieldTable["aliases"] = object.Field{Ftype: types.StringClassNameArray, Fvalue: aliasArray}
	this.FieldTable["attributes"] = object.Field{Ftype: types.Map, Fvalue: map[string]string{}}

	return nil
}

// ----------------------- Member Functions -----------------------
func securityProvSvcGetProvider(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["provider"].Fvalue
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
	this := params[0].(*object.Object)
	if params[1] == nil {
		return nil
	}
	keyObj, ok := params[1].(*object.Object)
	if !ok {
		return nil
	}
	key := object.GoStringFromStringObject(keyObj)
	attributes := this.FieldTable["attributes"].Fvalue.(map[string]string)
	if val, ok := attributes[key]; ok {
		return object.StringObjectFromGoString(val)
	}
	return nil
}

func securityProvSvcToString(params []any) any {
	this := params[0].(*object.Object)
	typeStr := object.GoStringFromStringObject(this.FieldTable["type"].Fvalue.(*object.Object))
	algStr := object.GoStringFromStringObject(this.FieldTable["algorithm"].Fvalue.(*object.Object))
	return object.StringObjectFromGoString(typeStr + "/" + algStr)
}
