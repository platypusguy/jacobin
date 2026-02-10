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
	"math/big"
)

// Load_KeyPairGenerator registers KeyPairGenerator methods in MethodSignatures
func Load_KeyPairGenerator() {
	// ---------------------------------------------------------
	// Constructor: protected
	// Should not be called directly
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/KeyPairGenerator.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapProtected,
		}

	// ---------------------------------------------------------
	// SPI class: trap everything
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/KeyPairGeneratorSpi.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/security/KeyPairGeneratorSpi.generateKeyPair()Ljava/security/KeyPair;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/security/KeyPairGeneratorSpi.initialize(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/security/KeyPairGeneratorSpi.initialize(ILjava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// ---------------------------------------------------------
	// Public API: getInstance
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/KeyPairGenerator.getInstance(Ljava/lang/String;)Ljava/security/KeyPairGenerator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  keypairgeneratorGetInstance,
		}
	ghelpers.MethodSignatures["java/security/KeyPairGenerator.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/KeyPairGenerator;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// ---------------------------------------------------------
	// Public API: initialize variants
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/KeyPairGenerator.initialize(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  keypairgeneratorInitialize,
		}

	ghelpers.MethodSignatures["java/security/KeyPairGenerator.initialize(ILjava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  keypairgeneratorInitializeWithRandom,
		}

	ghelpers.MethodSignatures["java/security/KeyPairGenerator.initialize(Ljava/security/spec/AlgorithmParameterSpec;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  keypairgeneratorInitializeParmSpec,
		}

	// ---------------------------------------------------------
	// Public API: generateKeyPair & genKeyPair
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/KeyPairGenerator.generateKeyPair()Ljava/security/KeyPair;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keypairgeneratorGenerateKeyPair,
		}

	ghelpers.MethodSignatures["java/security/KeyPairGenerator.genKeyPair()Ljava/security/KeyPair;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keypairgeneratorGenerateKeyPair,
		}

	// ---------------------------------------------------------
	// Optional member functions
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/KeyPairGenerator.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: keypairgeneratorGetAlgorithm}

	ghelpers.MethodSignatures["java/security/KeyPairGenerator.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: keypairgeneratorGetProvider}

	ghelpers.MethodSignatures["java/security/KeyPairGenerator.getKeySize()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: keypairgeneratorGetKeySize}
}

func keypairgeneratorGetInstance(params []any) any {
	algorithmObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairgeneratorGetInstance: Algorithm cannot be null")
	}
	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Get the default (only) security provider.
	providerObj := ghelpers.GetDefaultSecurityProvider() // single Go runtime provider

	// Try to get a service from the provider
	svcObj := securityProviderGetService([]interface{}{providerObj, object.StringObjectFromGoString("KeyPairGenerator"), algorithmObj})
	if errBlk, ok := svcObj.(*ghelpers.GErrBlk); ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keypairgeneratorGetInstance: "+errBlk.ErrMsg)
	}

	// Create KeyPairGenerator object
	kpgObj := object.MakeEmptyObjectWithClassName(&types.ClassNameKeyPairGenerator)

	// Store algorithm name
	kpgObj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algorithm),
	}

	// Store provider
	kpgObj.FieldTable["provider"] = object.Field{
		Ftype:  types.ClassNameSecurityProvider,
		Fvalue: providerObj,
	}

	// Store reference to the service.
	kpgObj.FieldTable["service"] = object.Field{
		Ftype:  types.ClassNameSecurityProviderService,
		Fvalue: svcObj.(*object.Object),
	}

	return kpgObj
}

func keypairgeneratorInitialize(params []any) any {
	// params[0] = KeyPairGenerator object
	// params[1] = int keySize
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitialize: missing keySize parameter",
		)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitialize: first parameter must be KeyPairGenerator object",
		)
	}

	keySize, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitialize: keySize must be an int",
		)
	}

	// Store keySize in the object's FieldTable
	obj.FieldTable["keySize"] = object.Field{Ftype: types.Int, Fvalue: keySize}

	return nil
}

func keypairgeneratorInitializeWithRandom(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitializeWithRandom: missing parameters",
		)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitializeWithRandom: first parameter must be KeyPairGenerator object",
		)
	}

	keySize, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitializeWithRandom: keySize must be an int",
		)
	}

	randomObj, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorInitializeWithRandom: random must be SecureRandom object",
		)
	}

	obj.FieldTable["keySize"] = object.Field{Ftype: types.Int, Fvalue: keySize}
	obj.FieldTable["random"] = object.Field{Ftype: types.Ref, Fvalue: randomObj}
	return nil
}

// keypairgeneratorInitializeParmSpec initializes the key pair generator with algorithm parameter spec
func keypairgeneratorInitializeParmSpec(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("KeyPairGenerator.initialize: expected 1 parameter, got %d", len(params)-1),
		)
	}

	kpgObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyPairGenerator.initialize: `this` is not an Object",
		)
	}

	paramSpecObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyPairGenerator.initialize: paramSpec is not an Object",
		)
	}

	// Get the algorithm from the generator
	algoObj := kpgObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algoObj)

	// Store the parameter spec for use during key generation
	kpgObj.FieldTable["paramSpec"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: paramSpecObj,
	}

	// Validate parameter spec type matches algorithm and extract keySize
	paramSpecClassName := object.GoStringFromStringPoolIndex(paramSpecObj.KlassName)
	var keySize int64

	switch algorithm {
	case "RSA":
		if paramSpecClassName != types.ClassNameRSAKeyGenParameterSpec {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				fmt.Sprintf("RSA requires RSAKeyGenParameterSpec, got %s", paramSpecClassName),
			)
		}
		// Extract key size from RSAKeyGenParameterSpec
		keySizeField, exists := paramSpecObj.FieldTable["keysize"]
		if !exists {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				"RSAKeyGenParameterSpec missing keysize field",
			)
		}
		keySize = keySizeField.Fvalue.(int64)

	case "EC":
		if paramSpecClassName != types.ClassNameECGenParameterSpec &&
			paramSpecClassName != types.ClassNameECParameterSpec {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				fmt.Sprintf("EC requires ECGenParameterSpec or ECParameterSpec, got %s", paramSpecClassName),
			)
		}
		// For EC, derive keySize from the curve
		if paramSpecClassName == types.ClassNameECGenParameterSpec {
			// ECGenParameterSpec has a curve name (e.g., "secp256r1")
			curveNameField, exists := paramSpecObj.FieldTable["name"]
			if !exists {
				return ghelpers.GetGErrBlk(
					excNames.InvalidAlgorithmParameterException,
					"ECGenParameterSpec missing name field",
				)
			}
			curveNameObj := curveNameField.Fvalue.(*object.Object)
			curveName := object.GoStringFromStringObject(curveNameObj)
			keySize = getKeySizeFromCurveName(curveName)
			if keySize == 0 {
				return ghelpers.GetGErrBlk(
					excNames.InvalidAlgorithmParameterException,
					fmt.Sprintf("unsupported curve name: %s", curveName),
				)
			}
		} else {
			// ECParameterSpec - extract from the order (n)
			nField, exists := paramSpecObj.FieldTable["n"]
			if !exists {
				return ghelpers.GetGErrBlk(
					excNames.InvalidAlgorithmParameterException,
					"ECParameterSpec missing n field",
				)
			}
			nObj := nField.Fvalue.(*object.Object)
			n := nObj.FieldTable["value"].Fvalue.(*big.Int)
			keySize = int64(n.BitLen())
		}

	case "DSA":
		if paramSpecClassName != types.ClassNameDSAParameterSpec {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				fmt.Sprintf("DSA requires DSAParameterSpec, got %s", paramSpecClassName),
			)
		}
		// Extract key size from P parameter
		pField, exists := paramSpecObj.FieldTable["p"]
		if !exists {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				"DSAParameterSpec missing p field",
			)
		}
		pObj := pField.Fvalue.(*object.Object)
		p := pObj.FieldTable["value"].Fvalue.(*big.Int)
		keySize = int64(p.BitLen())

	case "DH":
		if paramSpecClassName != types.ClassNameDHParameterSpec {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				fmt.Sprintf("DH requires DHParameterSpec, got %s", paramSpecClassName),
			)
		}
		// Extract key size from P parameter
		pField, exists := paramSpecObj.FieldTable["p"]
		if !exists {
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				"DHParameterSpec missing p field",
			)
		}
		pObj := pField.Fvalue.(*object.Object)
		p := pObj.FieldTable["value"].Fvalue.(*big.Int)
		keySize = int64(p.BitLen())

	case "Ed25519", "Ed448", "X25519", "X448", "XDH":
		// These algorithms don't use parameter specs
		return ghelpers.GetGErrBlk(
			excNames.InvalidAlgorithmParameterException,
			fmt.Sprintf("%s does not support AlgorithmParameterSpec initialization", algorithm),
		)

	default:
		return ghelpers.GetGErrBlk(
			excNames.InvalidAlgorithmParameterException,
			fmt.Sprintf("unsupported algorithm: %s", algorithm),
		)
	}

	// Set the key size
	kpgObj.FieldTable["keySize"] = object.Field{
		Ftype:  types.Int,
		Fvalue: keySize,
	}

	// Mark as initialized
	kpgObj.FieldTable["initialized"] = object.Field{
		Ftype:  types.Bool,
		Fvalue: types.JavaBoolTrue,
	}

	return nil
}

func keypairgeneratorGetAlgorithm(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetAlgorithm: missing KeyPairGenerator object",
		)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetAlgorithm: first parameter must be KeyPairGenerator object",
		)
	}

	algObj, ok := obj.FieldTable["algorithm"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetAlgorithm: algorithm not set",
		)
	}

	return algObj
}

func keypairgeneratorGetProvider(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetProvider: missing KeyPairGenerator object",
		)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetProvider: first parameter must be KeyPairGenerator object",
		)
	}

	providerObj, ok := obj.FieldTable["provider"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetProvider: provider not set",
		)
	}

	return providerObj
}

func keypairgeneratorGetKeySize(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetKeySize: missing KeyPairGenerator object",
		)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetKeySize: first parameter must be KeyPairGenerator object",
		)
	}

	keySize, ok := obj.FieldTable["keySize"].Fvalue.(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGetKeySize: keySize not set",
		)
	}

	return keySize
}

// Helper function to map curve names to key sizes
func getKeySizeFromCurveName(curveName string) int64 {
	curveMap := map[string]int64{
		"secp224r1": 224,
		"secp256r1": 256,
		"secp384r1": 384,
		"secp521r1": 521,
		"P-224":     224,
		"P-256":     256,
		"P-384":     384,
		"P-521":     521,
	}
	return curveMap[curveName]
}
