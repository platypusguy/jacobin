package javaSecurity

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"math/big"

	"golang.org/x/crypto/curve25519"

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

// keypairgeneratorGenerateKeyPair generates a KeyPair for supported algorithms.
func keypairgeneratorGenerateKeyPair(params []any) any {
	var err error
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGenerateKeyPair: missing KeyPairGenerator object",
		)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGenerateKeyPair: first parameter must be KeyPairGenerator object",
		)
	}

	algObj, ok := obj.FieldTable["algorithm"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGenerateKeyPair: algorithm not set",
		)
	}
	algorithm := object.GoStringFromStringObject(algObj)

	keySize, ok := obj.FieldTable["keySize"].Fvalue.(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGenerateKeyPair: keySize not set",
		)
	}

	var keyPairObj *object.Object

	switch algorithm {
	case "RSA":
		var priv *rsa.PrivateKey
		priv, err = rsa.GenerateKey(rand.Reader, int(keySize))
		if err == nil {
			keyPairObj = NewGoRuntimeService("KeyPair", "RSA", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: &priv.PublicKey}
		}

	case "DSA":
		params := new(dsa.Parameters)
		if err = dsa.GenerateParameters(params, rand.Reader, dsa.L2048N256); err == nil {
			priv := new(dsa.PrivateKey)
			priv.Parameters = *params
			if err = dsa.GenerateKey(priv, rand.Reader); err == nil {
				keyPairObj = NewGoRuntimeService("KeyPair", "DSA", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: &priv.PublicKey}
			}
		}

	case "DH":
		// Placeholder: generate simple DH big.Int values
		p, g := big.NewInt(0), big.NewInt(0)
		priv := big.NewInt(0)
		pub := big.NewInt(0)
		keyPairObj = NewGoRuntimeService("KeyPair", "DH", types.ClassNameKeyPair)
		keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
		keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}
		keyPairObj.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: p}
		keyPairObj.FieldTable["g"] = object.Field{Ftype: types.BigInteger, Fvalue: g}

	case "EC":
		var curve elliptic.Curve
		switch keySize {
		case 224:
			curve = elliptic.P224()
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		case 521:
			curve = elliptic.P521()
		default:
			err = errors.New("unsupported EC key size")
		}
		if err == nil {
			priv, err2 := ecdsa.GenerateKey(curve, rand.Reader)
			if err2 != nil {
				err = err2
			} else {
				keyPairObj = NewGoRuntimeService("KeyPair", "EC", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: &priv.PublicKey}
			}
		}

	case "Ed25519":
		pub, priv, err2 := ed25519.GenerateKey(rand.Reader)
		if err2 != nil {
			err = err2
		} else {
			keyPairObj = NewGoRuntimeService("KeyPair", "Ed25519", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}
		}

	case "XDH", "X25519":
		priv := make([]byte, 32)
		_, err2 := rand.Read(priv)
		if err2 != nil {
			err = err2
			break
		}
		pub, err2 := curve25519.X25519(priv, curve25519.Basepoint)
		if err2 != nil {
			err = err2
			break
		}
		keyPairObj = NewGoRuntimeService("KeyPair", algorithm, types.ClassNameKeyPair)
		keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
		keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}

	case "Ed448":
		// Ed448 not in standard Go; use placeholder for now
		priv := make([]byte, 57) // Ed448 private key length
		_, err2 := rand.Read(priv)
		if err2 != nil {
			err = err2
			break
		}
		pub := make([]byte, 57)
		copy(pub, priv) // placeholder: real Ed448 requires proper library
		keyPairObj = NewGoRuntimeService("KeyPair", "Ed448", types.ClassNameKeyPair)
		keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
		keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}

	default:
		keyPairObj = NewGoRuntimeService("KeyPair", algorithm, types.ClassNameKeyPair)
	}

	if err != nil {
		return ghelpers.GetGErrBlk(
			excNames.GeneralSecurityException,
			"keypairgeneratorGenerateKeyPair: "+algorithm+" key generation failed: "+err.Error(),
		)
	}

	return keyPairObj
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
