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

func Load_Security_AlgorithmParameters() {
	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  AlgparamsGetAlgorithm,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getEncoded()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  AlgparamsGetEncoded,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getEncoded(Ljava/lang/String;)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  AlgparamsGetEncoded,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getInstance(Ljava/lang/String;)Ljava/security/AlgorithmParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  AlgparamsGetInstance,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/AlgorithmParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  AlgparamsGetInstance,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljava/security/AlgorithmParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  AlgparamsGetInstance,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getParameterSpec(Ljava/lang/Class;)Ljava/security/spec/AlgorithmParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  AlgparamsGetParameterSpec,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  AlgparamsGetProvider,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.init(Ljava/security/spec/AlgorithmParameterSpec;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  AlgparamsInit,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.init([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  AlgparamsInit,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.init([BLjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  AlgparamsInit,
		}

	ghelpers.MethodSignatures["java/security/AlgorithmParameters.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  AlgparamsToString,
		}
}

func AlgparamsGetAlgorithm(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	algo, ok := self.FieldTable["algorithm"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	return algo
}

func AlgparamsGetEncoded(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "algparamsGetEncoded: this is null")
	}

	initialized, ok := self.FieldTable["initialized"].Fvalue.(bool)
	if !ok || !initialized {
		return ghelpers.GetGErrBlk(excNames.IOException, "algparamsGetEncoded: not initialized")
	}

	// For now, we only support basic encoding if the parameter was initialized with byte[]
	// In a real implementation, we would use the spec to encode.
	var encoded []byte
	if encField, ok := self.FieldTable["encoded"]; ok {
		encoded, _ = encField.Fvalue.([]byte)
	}

	if encoded == nil {
		// If we were initialized with a spec, we might need to derive encoded bytes
		if specField, ok := self.FieldTable["spec"]; ok {
			specObj, ok := specField.Fvalue.(*object.Object)
			if ok && specObj != nil {
				// Try to get IV if it's an IvParameterSpec
				if ivField, ok := specObj.FieldTable["iv"]; ok {
					if iv, ok := ivField.Fvalue.([]byte); ok {
						encoded = iv
					}
				}
			}
		}
	}

	if encoded == nil {
		return ghelpers.GetGErrBlk(excNames.IOException, "algparamsGetEncoded: no encoded form available")
	}

	jBytes := object.JavaByteArrayFromGoByteArray(encoded)
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}

func AlgparamsGetInstance(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "algparamsGetInstance: missing parameters")
	}

	algorithmObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(algorithmObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "algparamsGetInstance: algorithm cannot be null")
	}
	algorithm := object.GoStringFromStringObject(algorithmObj)

	providerObj := ghelpers.GetDefaultSecurityProvider()
	if object.IsNull(providerObj) {
		providerObj = InitDefaultSecurityProvider()
	}
	// Check if the service is supported
	svcObj := securityProviderGetService([]any{providerObj, object.StringObjectFromGoString(types.SecurityServiceAlgorithmParameters), algorithmObj})
	if _, ok := svcObj.(*ghelpers.GErrBlk); ok {
		return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, fmt.Sprintf("algparamsGetInstance: %s not found", algorithm))
	}

	// Handle provider parameter if provided
	if len(params) > 2 {
		p := params[2]
		if pObj, ok := p.(*object.Object); ok && !object.IsNull(pObj) {
			var pName string
			if pObj.KlassName == object.StringPoolIndexFromGoString(types.StringClassName) {
				pName = object.GoStringFromStringObject(pObj)
			} else {
				nameField, ok := pObj.FieldTable["name"]
				if ok {
					pName = object.GoStringFromStringObject(nameField.Fvalue.(*object.Object))
				}
			}

			if pName != "" && pName != types.SecurityProviderName {
				return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException, fmt.Sprintf("algparamsGetInstance: provider %s not found", pName))
			}
		}
	}

	algParams := object.MakeEmptyObjectWithClassName(&types.ClassNameAlgorithmParameters)
	algParams.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}
	algParams.FieldTable["provider"] = object.Field{Ftype: types.ClassNameSecurityProvider, Fvalue: providerObj}
	algParams.FieldTable["initialized"] = object.Field{Ftype: types.Bool, Fvalue: false}

	return algParams
}

func AlgparamsGetParameterSpec(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "algparamsGetParameterSpec: this is null")
	}

	initialized, ok := self.FieldTable["initialized"].Fvalue.(bool)
	if !ok || !initialized {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "algparamsGetParameterSpec: not initialized")
	}

	spec, ok := self.FieldTable["spec"].Fvalue.(*object.Object)
	if !ok || object.IsNull(spec) {
		return ghelpers.ReturnNull(params)
	}

	return spec
}

func AlgparamsGetProvider(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	provider, ok := self.FieldTable["provider"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	return provider
}

func AlgparamsInit(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "algparamsInit: this is null")
	}

	initialized, ok := self.FieldTable["initialized"].Fvalue.(bool)
	if ok && initialized {
		return ghelpers.GetGErrBlk(excNames.IOException, "algparamsInit: already initialized")
	}

	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "algparamsInit: missing parameters")
	}

	param := params[1]
	if paramObj, ok := param.(*object.Object); ok && !object.IsNull(paramObj) {
		if paramObj.KlassName == object.StringPoolIndexFromGoString(types.ByteArray) {
			// init([B) or init([B, String)
			encoded := object.GoByteArrayFromJavaByteArray(paramObj.FieldTable["value"].Fvalue.([]types.JavaByte))
			self.FieldTable["encoded"] = object.Field{Ftype: types.ByteArray, Fvalue: encoded}
			self.FieldTable["initialized"] = object.Field{Ftype: types.Bool, Fvalue: true}

			// If we know the algorithm, we might want to create a spec from the bytes
			algoObj := self.FieldTable["algorithm"].Fvalue.(*object.Object)
			algo := object.GoStringFromStringObject(algoObj)
			if algo == "AES" || algo == "DES" || algo == "DESede" {
				// Create IvParameterSpec
				ivSpecClass := "javax/crypto/spec/IvParameterSpec"
				ivSpec := object.MakeEmptyObjectWithClassName(&ivSpecClass)
				ivSpec.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: encoded}
				self.FieldTable["spec"] = object.Field{Ftype: "Ljavax/crypto/spec/IvParameterSpec;", Fvalue: ivSpec}
			}
		} else {
			// init(AlgorithmParameterSpec)
			self.FieldTable["spec"] = object.Field{Ftype: "Ljava/security/spec/AlgorithmParameterSpec;", Fvalue: paramObj}
			self.FieldTable["initialized"] = object.Field{Ftype: types.Bool, Fvalue: true}
		}
	} else {
		return ghelpers.GetGErrBlk(excNames.InvalidAlgorithmParameterException, "algparamsInit: invalid parameter")
	}

	return nil
}

func AlgparamsToString(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return object.StringObjectFromGoString("null")
	}

	algoObj := self.FieldTable["algorithm"].Fvalue.(*object.Object)
	algo := object.GoStringFromStringObject(algoObj)

	return object.StringObjectFromGoString(fmt.Sprintf("AlgorithmParameters[%s]", algo))
}
