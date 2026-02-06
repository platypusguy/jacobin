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
	"strconv"
	"strings"
)

func Load_Security_Provider() {

	ghelpers.MethodSignatures["java/security/Provider.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  securityProviderInit,
		}

	// ---------- Constructors ----------
	ghelpers.MethodSignatures["java/security/Provider.<init>(Ljava/lang/String;DLjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/security/Provider.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  securityProviderInit,
		}

	// ---------- Member Functions ----------

	ghelpers.MethodSignatures["java/security/Provider.clear()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/security/Provider.getInfo()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProviderGetInfo}

	ghelpers.MethodSignatures["java/security/Provider.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProviderGetName}

	ghelpers.MethodSignatures["java/security/Provider.getProperty(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/security/Provider.getService(Ljava/lang/String;Ljava/lang/String;)Ljava/security/Provider$Service;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: securityProviderGetService}

	ghelpers.MethodSignatures["java/security/Provider.getVersion()D"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/security/Provider.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/security/Provider.putAll(Ljava/util/Map;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/security/Provider.putService(Ljava/security/Provider$Service;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: securityProviderPutService}

	ghelpers.MethodSignatures["java/security/Provider.remove(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/security/Provider.removeService(Ljava/security/Provider$Service;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/security/Provider.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: securityProviderToString}

	// Set up a vector so that other functions can find the one and only Security Provider.
	ghelpers.DefaultSecurityProvider = InitDefaultSecurityProvider()

}

// ----------------------- Constructor -----------------------
func securityProviderInit(params []any) any {
	var version float64
	var err error

	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderInit: invalid `this` object, saw: %T", params[0]))
	}

	// name
	nameObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderInit: invalid name object, saw: %T", params[1]))
	}
	nameStr := strings.TrimSpace(object.GoStringFromStringObject(nameObj))
	this.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(nameStr)}

	// version
	versionAny := params[2]
	switch v := versionAny.(type) {
	case *object.Object:
		if object.IsNull(versionAny) {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				fmt.Sprintf("securityProviderInit: invalid version object, saw: %T", params[2]))
		}
		versionStr := strings.TrimSpace(object.GoStringFromStringObject(v))
		version, err = strconv.ParseFloat(versionStr, 64)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.VirtualMachineError,
				fmt.Sprintf("securityProviderInit: failed parsing version: '%s'", versionStr))
		}
	case float64:
		version = v
	default:
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError,
			fmt.Sprintf("securityProviderInit: invalid version type, saw: %T", params[2]))
	}
	this.FieldTable["version"] = object.Field{Ftype: types.Double, Fvalue: version}

	// info
	infoObj, ok := params[3].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderInit: invalid info object, saw: %T", infoObj))
	}
	infoStr := strings.TrimSpace(object.GoStringFromStringObject(infoObj))
	this.FieldTable["info"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(infoStr)}

	// Initialize services map in an empty state.
	this.FieldTable["services"] = object.Field{Ftype: types.Map, Fvalue: map[string]*object.Object{}}

	return nil
}

// ----------------------- Getters -----------------------

func securityProviderGetName([]any) any {
	return object.StringObjectFromGoString(types.SecurityProviderName)
}

func securityProviderGetInfo([]any) any {
	return object.StringObjectFromGoString(types.SecurityProviderInfo)
}

func securityProviderToString(params []any) any {
	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderToString: invalid `this` object, saw: %T", params[0]))
	}
	nameObj, ok := this.FieldTable["name"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderToString: invalid name object, saw: %T", nameObj))
	}
	name := object.GoStringFromStringObject(nameObj)
	version, ok := this.FieldTable["version"].Fvalue.(float64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderToString: invalid version, should be float64, saw: %T", version))
	}
	infoObj, ok := this.FieldTable["info"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderToString: invalid info object, saw: %T", infoObj))
	}
	info := object.GoStringFromStringObject(infoObj)

	return object.StringObjectFromGoString(fmt.Sprintf("%s %.1f\n%s", name, version, info))
}

// ----------------------- Get/Put Services -----------------------

func securityProviderGetService(params []any) any {
	if len(params) != 3 {
		return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException,
			fmt.Sprintf("securityProviderGetService: wrong number of args, expected 3, saw: %d", len(params)))
	}

	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderGetService: invalid `this` object, saw: %T", params[0]))
	}

	if object.IsNull(params[1]) {
		return nil
	}
	typeObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderGetService: invalid type object, saw: %T", params[1]))
	}
	typeStr := object.GoStringFromStringObject(typeObj)

	if object.IsNull(params[2]) {
		return nil
	}
	algObj, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("securityProviderGetService: invalid alg object, saw: %T", params[2]))
	}
	algStr := object.GoStringFromStringObject(algObj)

	services := this.FieldTable["services"].Fvalue.(map[string]*object.Object)
	key := typeStr + "/" + algStr
	if svc, ok := services[key]; ok {
		return svc
	}

	//if secSvcTypeMap, ok := SecurityProviderServices[typeStr]; ok {
	//	if svcInit, ok2 := secSvcTypeMap[algStr]; ok2 {
	//		return svcInit()
	//	}
	//}

	return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException,
		fmt.Sprintf("securityProviderGetService: unsupported type/algorithm %s/%s", typeStr, algStr))
}

func securityProviderPutService(params []any) any {
	this := params[0].(*object.Object)
	svc := params[1].(*object.Object)

	svcType := object.GoStringFromStringObject(svc.FieldTable["type"].Fvalue.(*object.Object))
	svcAlgo := object.GoStringFromStringObject(svc.FieldTable["algorithm"].Fvalue.(*object.Object))
	key := svcType + "/" + svcAlgo

	services := this.FieldTable["services"].Fvalue.(map[string]*object.Object)
	services[key] = svc
	return nil
}
