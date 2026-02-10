package javaxCrypto

import (
	"crypto/ecdsa"

	"fmt"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaSecurity"
	"jacobin/src/object"
	"jacobin/src/types"

	"golang.org/x/crypto/curve25519"
)

func Load_Crypto_KeyAgreement() {
	// <clinit>
	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.doPhase(Ljava/security/Key;Z)Ljava/security/Key;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  keyagreementDoPhase,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.generateSecret()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyagreementGenerateSecret,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.generateSecret([BI)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.generateSecret(Ljava/lang/String;)Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyagreementGetAlgorithm,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.getInstance(Ljava/lang/String;)Ljavax/crypto/KeyAgreement;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  keyagreementGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/crypto/KeyAgreement;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/crypto/KeyAgreement;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  javaSecurity.SecurityGetProvider,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.init(Ljava/security/Key;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  keyagreementInit,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.init(Ljava/security/Key;Ljava/security/spec/AlgorithmParameterSpec;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KeyAgreement.init(Ljava/security/Key;Ljava/security/spec/AlgorithmParameterSpec;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

}

// keyagreementGetInstance creates a new KeyAgreement object
func keyagreementGetInstance(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("keyagreementGetInstance: expected 1 parameter, got %d", len(params)),
		)
	}

	algorithmObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keyagreementGetInstance: algorithm is not a String object",
		)
	}

	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Validate algorithm
	if !isSupportedKeyAgreementAlgorithm(algorithm) {
		return ghelpers.GetGErrBlk(
			excNames.NoSuchAlgorithmException,
			fmt.Sprintf("keyagreementGetInstance: unsupported key agreement algorithm: %s", algorithm),
		)
	}

	// Create KeyAgreement object
	kaClassName := "javax/crypto/KeyAgreement"
	kaObj := object.MakeEmptyObjectWithClassName(&kaClassName)

	// Store algorithm
	kaObj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algorithm),
	}

	// State: 0=uninitialized, 1=initialized, 2=phase complete
	kaObj.FieldTable["state"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(0),
	}

	return kaObj
}

// keyagreementGetAlgorithm returns the algorithm name
func keyagreementGetAlgorithm(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keyagreementGetAlgorithm: expected 0 parameters",
		)
	}

	kaObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keyagreementGetAlgorithm: this is not an Object",
		)
	}

	return kaObj.FieldTable["algorithm"].Fvalue
}

// keyagreementInit initializes the KeyAgreement with a private key
func keyagreementInit(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keyagreementInit: expected 1 parameter",
		)
	}

	kaObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keyagreementInit: `this` is not an Object",
		)
	}

	privateKeyObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keyagreementInit: key is not an Object",
		)
	}

	// Get algorithm
	algoObj := kaObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algoObj)

	// Validate key type matches algorithm
	keyClassName := object.GoStringFromStringPoolIndex(privateKeyObj.KlassName)

	if algorithm == "ECDH" || algorithm == "EC" {
		if keyClassName != types.ClassNameECPrivateKey {
			return ghelpers.GetGErrBlk(
				excNames.InvalidKeyException,
				fmt.Sprintf("ECDH requires EC private key, got %s", keyClassName),
			)
		}
	} else if algorithm == "DH" {
		if keyClassName != types.ClassNameDHPrivateKey {
			return ghelpers.GetGErrBlk(
				excNames.InvalidKeyException,
				fmt.Sprintf("DH requires DH private key, got %s", keyClassName),
			)
		}
	} else if algorithm == "XDH" || algorithm == "X25519" {
		if keyClassName != types.ClassNameX25519PrivateKey {
			return ghelpers.GetGErrBlk(
				excNames.InvalidKeyException,
				fmt.Sprintf("X25519 requires X25519 private key, got %s", keyClassName),
			)
		}
	}

	// Store private key
	kaObj.FieldTable["privateKey"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: privateKeyObj,
	}

	// Set state to initialized
	kaObj.FieldTable["state"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(1),
	}

	return nil
}

// keyagreementDoPhase performs a phase of the key agreement
func keyagreementDoPhase(params []any) any {
	if len(params) != 3 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyAgreement.doPhase: expected 2 parameters",
		)
	}

	kaObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyAgreement.doPhase: this is not an Object",
		)
	}

	publicKeyObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyAgreement.doPhase: key is not an Object",
		)
	}

	lastPhase, ok := params[2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyAgreement.doPhase: lastPhase is not a boolean",
		)
	}

	// Check state
	state := kaObj.FieldTable["state"].Fvalue.(int64)
	if state != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"KeyAgreement.doPhase: not initialized",
		)
	}

	// Get algorithm
	algoObj := kaObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algoObj)

	// Validate key type
	keyClassName := object.GoStringFromStringPoolIndex(publicKeyObj.KlassName)

	if algorithm == "ECDH" || algorithm == "EC" {
		if keyClassName != types.ClassNameECPublicKey {
			return ghelpers.GetGErrBlk(
				excNames.InvalidKeyException,
				fmt.Sprintf("ECDH requires EC public key, got %s", keyClassName),
			)
		}
	} else if algorithm == "X25519" || algorithm == "XDH" {
		if keyClassName != types.ClassNameX25519PublicKey {
			return ghelpers.GetGErrBlk(
				excNames.InvalidKeyException,
				fmt.Sprintf("X25519 requires X25519 public key, got %s", keyClassName),
			)
		}
	}

	// Store public key
	kaObj.FieldTable["publicKey"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: publicKeyObj,
	}

	// Set state to phase complete
	kaObj.FieldTable["state"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(2),
	}

	// For single-phase protocols (ECDH, X25519), return null
	if lastPhase != 0 { // Java boolean true
		return object.Null
	}

	// For multi-phase, would return intermediate key (not implemented)
	return object.Null
}

// keyagreementGenerateSecret generates the shared secret
func keyagreementGenerateSecret(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyAgreement.generateSecret: expected 0 parameters",
		)
	}

	kaObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"KeyAgreement.generateSecret: this is not an Object",
		)
	}

	// Check state
	state := kaObj.FieldTable["state"].Fvalue.(int64)
	if state != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalStateException,
			"KeyAgreement.generateSecret: doPhase not completed",
		)
	}

	// Get algorithm
	algoObj := kaObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algoObj)

	// Get private and public keys
	privateKeyObj := kaObj.FieldTable["privateKey"].Fvalue.(*object.Object)
	publicKeyObj := kaObj.FieldTable["publicKey"].Fvalue.(*object.Object)

	var secretBytes []byte
	var err error

	switch algorithm {
	case "ECDH", "EC":
		secretBytes, err = performECDH(privateKeyObj, publicKeyObj)
	case "X25519", "XDH":
		secretBytes, err = performX25519(privateKeyObj, publicKeyObj)
	default:
		return ghelpers.GetGErrBlk(
			excNames.NoSuchAlgorithmException,
			fmt.Sprintf("unsupported algorithm: %s", algorithm),
		)
	}

	if err != nil {
		return ghelpers.GetGErrBlk(excNames.InvalidKeyException, err.Error())
	}

	// Return as byte array
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(secretBytes))
}

// Helper functions

func isSupportedKeyAgreementAlgorithm(algo string) bool {
	supported := []string{"ECDH", "EC", "X25519", "XDH", "DH"}
	for _, s := range supported {
		if s == algo {
			return true
		}
	}
	return false
}

// performECDH performs Elliptic Curve Diffie-Hellman
func performECDH(privateKeyObj, publicKeyObj *object.Object) ([]byte, error) {
	// Extract private key
	privKeyValue := privateKeyObj.FieldTable["value"].Fvalue.(*ecdsa.PrivateKey)

	// Extract public key
	pubKeyValue := publicKeyObj.FieldTable["value"].Fvalue.(*ecdsa.PublicKey)

	// Validate curves match
	if privKeyValue.Curve.Params().Name != pubKeyValue.Curve.Params().Name {
		return nil, fmt.Errorf("curve mismatch")
	}

	// Perform scalar multiplication: shared_secret = private_key * public_key_point
	x, _ := pubKeyValue.Curve.ScalarMult(pubKeyValue.X, pubKeyValue.Y, privKeyValue.D.Bytes())

	// The shared secret is the x-coordinate
	return x.Bytes(), nil
}

// performX25519 performs X25519 key agreement
func performX25519(privateKeyObj, publicKeyObj *object.Object) ([]byte, error) {
	// Extract private key (32 bytes)
	privKey := privateKeyObj.FieldTable["value"].Fvalue.([]byte)

	// Extract public key (32 bytes)
	pubKey := publicKeyObj.FieldTable["value"].Fvalue.([]byte)

	// Perform X25519 scalar multiplication
	secret, err := curve25519.X25519(privKey, pubKey)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
