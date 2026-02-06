package javaSecurity

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Load_ECKeys registers EC interfaces and concrete classes
func Load_EC_Keys() {
	// ---------------------------------------------------------
	// Interfaces
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/interfaces/PublicKey.<init>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.JustReturn}

	ghelpers.MethodSignatures["java/security/interfaces/PrivateKey.<init>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.JustReturn}

	ghelpers.MethodSignatures["java/security/interfaces/ECKey.<init>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.JustReturn}

	// ---------------------------------------------------------
	// Concrete EC public key
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/interfaces/ECPublicKey.<init>(Ljava/security/spec/ECParameterSpec;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ECPublicKeyInit}

	ghelpers.MethodSignatures["java/security/interfaces/ECPublicKey.getParams()Ljava/security/spec/ECParameterSpec;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ecPublicKeyGetParams}

	// ---------------------------------------------------------
	// Concrete EC private key
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/interfaces/ECPrivateKey.<init>(Ljava/security/spec/ECParameterSpec;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ECPrivateKeyInit}

	ghelpers.MethodSignatures["java/security/interfaces/ECPrivateKey.getParams()Ljava/security/spec/ECParameterSpec;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ecPrivateKeyGetParams}
}

// ---------------------------------------------------------
// ECPublicKey concrete class
// ---------------------------------------------------------
func ECPublicKeyInit(params []any) any {
	if len(params) != 2 { // this + ECParameterSpec
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPublicKeyInit: expected 1 parameter (params), got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyInit: this is not an Object",
		)
	}

	specObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyInit: param is not ECParameterSpec",
		)
	}

	thisObj.FieldTable = map[string]object.Field{
		"params": {Ftype: types.Ref, Fvalue: specObj},
	}

	return nil
}

func ecPublicKeyGetParams(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPublicKeyGetParams: expected 0 parameters, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyGetParams: this is not an Object",
		)
	}

	specObj, exists := thisObj.FieldTable["params"]
	if !exists {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPublicKeyGetParams: params field missing",
		)
	}

	return specObj.Fvalue
}

// ---------------------------------------------------------
// ECPrivateKey concrete class
// ---------------------------------------------------------
func ECPrivateKeyInit(params []any) any {
	if len(params) != 2 { // this + ECParameterSpec
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPrivateKeyInit: expected 1 parameter (params), got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyInit: this is not an Object",
		)
	}

	specObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyInit: param is not ECParameterSpec",
		)
	}

	thisObj.FieldTable = map[string]object.Field{
		"params": {Ftype: types.Ref, Fvalue: specObj},
	}

	return nil
}

func ecPrivateKeyGetParams(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("ecPrivateKeyGetParams: expected 0 parameters, got %d", len(params)-1),
		)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyGetParams: this is not an Object",
		)
	}

	specObj, exists := thisObj.FieldTable["params"]
	if !exists {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"ecPrivateKeyGetParams: params field missing",
		)
	}

	return specObj.Fvalue
}
