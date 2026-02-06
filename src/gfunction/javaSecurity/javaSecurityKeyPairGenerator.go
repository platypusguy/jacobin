package javaSecurity

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
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
