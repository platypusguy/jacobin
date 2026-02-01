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
			GFunction:  ghelpers.ClinitGeneric,
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

}

// ----------------------- Constructor -----------------------
func securityProviderInit(params []any) any {
	var version float64
	var err error

	this := params[0].(*object.Object)

	// name
	nameObj := params[1].(*object.Object)
	nameStr := strings.TrimSpace(object.GoStringFromStringObject(nameObj))
	this.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(nameStr)}

	// version
	versionAny := params[2]
	switch v := versionAny.(type) {
	case *object.Object:
		versionStr := strings.TrimSpace(object.GoStringFromStringObject(v))
		version, err = strconv.ParseFloat(versionStr, 64)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.VirtualMachineError, fmt.Sprintf("securityProviderInit: failed parsing version '%s'", versionStr))
		}
	case float64:
		version = v
	default:
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, fmt.Sprintf("securityProviderInit: invalid version type %T", versionAny))
	}
	this.FieldTable["version"] = object.Field{Ftype: types.Double, Fvalue: version}

	// info
	infoObj := params[3].(*object.Object)
	infoStr := strings.TrimSpace(object.GoStringFromStringObject(infoObj))
	this.FieldTable["info"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString(infoStr)}

	// initialize services map
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
	this := params[0].(*object.Object)
	name := object.GoStringFromStringObject(this.FieldTable["name"].Fvalue.(*object.Object))
	version := this.FieldTable["version"].Fvalue.(float64)
	info := object.GoStringFromStringObject(this.FieldTable["info"].Fvalue.(*object.Object))
	return object.StringObjectFromGoString(fmt.Sprintf("%s %.1f\n%s", name, version, info))
}

// ----------------------- Services -----------------------
func securityProviderGetService(params []any) any {
	this := params[0].(*object.Object)

	if params[1] == nil || params[2] == nil {
		return nil
	}
	typeStr := object.GoStringFromStringObject(params[1].(*object.Object))
	algStr := object.GoStringFromStringObject(params[2].(*object.Object))

	services := this.FieldTable["services"].Fvalue.(map[string]*object.Object)
	key := typeStr + "/" + algStr
	if svc, ok := services[key]; ok {
		return svc
	}
	return nil
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

// ----------------------- Helper: Default Go Runtime Provider -----------------------
// ----------------------- Used at Jacobin startup -----------------------------------
func NewGoRuntimeProvider() *object.Object {
	// Create the Provider object
	provider := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)

	// Initialize the provider with name, version, info
	nameObj := object.StringObjectFromGoString(types.SecurityProviderName)
	infoObj := object.StringObjectFromGoString(types.SecurityProviderInfo)
	params := []any{provider, nameObj, 1.0, infoObj} // version=1.0
	securityProviderInit(params)

	// Create the default Provider$Service
	className := "java/security/Provider$Service"
	service := object.MakeEmptyObjectWithClassName(&className)

	// Service fields: provider, type, algorithm, className, aliases
	serviceType := object.StringObjectFromGoString("Runtime")
	serviceAlgorithm := object.StringObjectFromGoString("Security")
	serviceClassName := object.StringObjectFromGoString(types.SecurityProviderName)
	aliases := []*object.Object{} // empty aliases
	serviceParams := []any{service, provider, serviceType, serviceAlgorithm, serviceClassName, aliases}
	securityProvSvcInit(serviceParams)

	// Register the service with the provider
	securityProviderPutService([]any{provider, service})

	return provider
}
