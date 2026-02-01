package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaSecurity"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

func Load_Sun_Security_Jca_ProviderList() {

	// <clinit> — load statics
	ghelpers.MethodSignatures["sun/security/jca/ProviderList.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clinitProviderList,
		}

	// <init>
	ghelpers.MethodSignatures["sun/security/jca/ProviderList.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	// ---- member functions ----

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.add(Lsun/security/jca/ProviderConfig;)Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.fromSecurityProperties()Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: providerListFromSecurityProperties}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.getDefault()Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: providerListGetDefault}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.getJarList(Ljava/lang/String;)Ljava/util/List;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.getProvider(Ljava/lang/String;)Ljava/security/Provider;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: providerListGetProvider}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.getProviderConfig(Ljava/lang/String;)Lsun/security/jca/ProviderConfig;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.getProviderConfigs()Ljava/util/List;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.getService(Ljava/lang/String;Ljava/lang/String;)Ljava/security/Provider$Service;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.insertAt(Lsun/security/jca/ProviderConfig;I)Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.isEmpty()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ReturnFalse}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.loadAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.JustReturn}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.newList(Lsun/security/jca/ProviderConfig;)Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.providers()Ljava/util/List;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: providerListProviders}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.remove(Ljava/lang/String;)Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.removeInvalid()Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.size()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: providerListSize}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.toArray()[Ljava/security/Provider;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["sun/security/jca/ProviderList$ServiceList.tryGet(I)Ljava/security/Provider$Service;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: providerListTryGet}

}

func clinitProviderList(params []interface{}) interface{} {

	thisClassName := "sun/security/jca/ProviderList"

	// debug stub → null
	_ = statics.AddStatic(thisClassName+".debug",
		statics.Static{Type: types.Ref, Value: object.Null})

	// PC0 → ProviderConfig[] with one empty ProviderConfig
	providerConfigClassName := "sun/security/jca/ProviderConfig"
	providerConfig := object.MakeEmptyObjectWithClassName(&providerConfigClassName)
	pcArray := object.MakePrimitiveObject(
		"[Lsun/security/jca/ProviderConfig;", // proper array class
		types.RefArray,
		[]*object.Object{providerConfig},
	)
	_ = statics.AddStatic("L"+thisClassName+".PC0",
		statics.Static{Type: types.Ref, Value: pcArray})

	// P0 → Provider[] with the runtime provider
	provider := javaSecurity.NewGoRuntimeProvider()
	providerArray := object.MakePrimitiveObject(
		"[Ljava/security/Provider;", // proper array class
		types.RefArray,
		[]*object.Object{provider},
	)
	_ = statics.AddStatic("L"+thisClassName+".P0",
		statics.Static{Type: types.Ref, Value: providerArray})

	// EMPTY → ProviderList object with zero providers
	emptyProviderList := object.MakeEmptyObjectWithClassName(&thisClassName)
	emptyProviderList.FieldTable["providers"] = object.Field{
		Ftype: types.Ref,
		Fvalue: object.MakePrimitiveObject(
			"[Ljava/security/Provider;",
			types.RefArray,
			[]*object.Object{}, // zero providers
		),
	}
	_ = statics.AddStatic(thisClassName+".EMPTY",
		statics.Static{Type: types.Ref, Value: emptyProviderList})

	// preferredPropList → null
	_ = statics.AddStatic(thisClassName+".preferredPropList",
		statics.Static{Type: types.Ref, Value: object.Null})

	return nil
}

func providerListFromSecurityProperties(params []interface{}) interface{} {
	thisClassName := "sun/security/jca/ProviderList"

	// fetch P0 array from statics
	providerArray := statics.GetStaticValue("L"+thisClassName, "P0")
	if providerArray == nil {
		return object.Null
	}

	// create ProviderList object containing exactly the runtime provider
	providerList := object.MakeEmptyObjectWithClassName(&thisClassName)
	providerList.FieldTable["providers"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: providerArray.(*object.Object),
	}

	return providerList
}

func providerListGetDefault(params []interface{}) interface{} {
	thisClassName := "sun/security/jca/ProviderList"

	// fetch P0 array from statics
	providerArray := statics.GetStaticValue("L"+thisClassName, "P0")
	if providerArray == nil {
		return object.Null
	}

	// create ProviderList object
	providerList := object.MakeEmptyObjectWithClassName(&thisClassName)
	providerList.FieldTable["providers"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: providerArray.(*object.Object),
	}

	return providerList
}

func providerListGetProvider([]interface{}) interface{} {
	return javaSecurity.NewGoRuntimeProvider()
}

func providerListProviders(params []interface{}) interface{} {
	// params[0] is `this` object
	thisObj := params[0].(*object.Object)
	thisClassName := object.GoStringFromStringPoolIndex(thisObj.KlassName)

	providersField, ok := thisObj.FieldTable["providers"]
	if ok {
		return providersField.Fvalue
	}

	// fallback to EMPTY if missing
	emptyPL := statics.GetStaticValue(thisClassName, "EMPTY")
	if emptyPL == nil {
		return object.Null
	}

	return emptyPL.(*object.Object).FieldTable["providers"].Fvalue
}

func providerListSize([]interface{}) interface{} {
	return int64(1)
}

func providerListTryGet(params []interface{}) interface{} {
	// params[0] = this object (ProviderList)
	// params[1] = integer key / index
	if len(params) < 2 {
		return object.Null
	}

	// The integer is ignored because we have no services
	// Always return null safely
	return object.Null
}
