/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"

	"github.com/unix-world/smartgo/crypto/eddsa/ed448"
	"golang.org/x/crypto/curve25519"
)

// keypairgeneratorGenerateKeyPair generates a KeyPair for supported algorithms.
// Parameters: KeyPairGenerator object
// FieldTable
// 0: Algorithm string
// 1: Key size (optional, default 2048)
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
		keySize = int64(2048)
	}

	var keyPairObj *object.Object

	switch algorithm {
	case "RSA":
		// --- RSA Key Generation ---
		var privRsa *rsa.PrivateKey
		privRsa, err = rsa.GenerateKey(rand.Reader, int(keySize))
		if err == nil {
			pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privRsa.PublicKey)
			if err != nil {
				return ghelpers.GetGErrBlk(excNames.GeneralSecurityException,
					"RSA x509.MarshalPKIXPublicKey failed: "+err.Error())
			}
			pubKeyInterface, err := x509.ParsePKIXPublicKey(pubKeyBytes)
			if err != nil {
				return ghelpers.GetGErrBlk(excNames.GeneralSecurityException,
					"RSA x509.ParsePKIXPublicKey failed: "+err.Error())
			}
			pubRsa, ok := pubKeyInterface.(*rsa.PublicKey)
			if !ok {
				return ghelpers.GetGErrBlk(excNames.GeneralSecurityException,
					"RSA pubKeyInterface.(*rsa.PublicKey) failed")
			}

			privateKeyObj := NewGoRuntimeService("RSA", "RSA", types.ClassNameRSAPrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privRsa}

			publicKeyObj := NewGoRuntimeService("RSA", "RSA", types.ClassNameRSAPublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubRsa}

			keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "RSA", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.ClassNameRSAPrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.ClassNameRSAPublicKey, Fvalue: publicKeyObj}
		}

	case "DH":
		// --- Simple DH example ---
		priv := big.NewInt(6)
		pub := big.NewInt(8)

		privateKeyObj := NewGoRuntimeService("DH", "DH", types.ClassNameDHPrivateKey)
		privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

		publicKeyObj := NewGoRuntimeService("DH", "DH", types.ClassNameDHPublicKey)
		publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}

		keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "DH", types.ClassNameKeyPair)
		keyPairObj.FieldTable["private"] = object.Field{Ftype: types.ClassNameDHPrivateKey, Fvalue: privateKeyObj}
		keyPairObj.FieldTable["public"] = object.Field{Ftype: types.ClassNameDHPublicKey, Fvalue: publicKeyObj}

	case "DSA":
		params := new(dsa.Parameters)
		err = dsa.GenerateParameters(params, rand.Reader, dsa.L2048N256)
		if err == nil {
			priv := new(dsa.PrivateKey)
			priv.Parameters = *params
			err = dsa.GenerateKey(priv, rand.Reader)
			if err == nil {
				dsaParamsObj := NewGoRuntimeService("DSA", "DSA", types.ClassNameDSAParameterSpec)
				dsaParamsObj.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.P)}
				dsaParamsObj.FieldTable["q"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.Q)}
				dsaParamsObj.FieldTable["g"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.G)}

				pubKey := &dsa.PublicKey{Parameters: *params, Y: new(big.Int).Set(priv.PublicKey.Y)}
				publicKeyObj := NewGoRuntimeService("DSA", "DSA", types.ClassNameDSAPublicKey)
				publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKey}
				publicKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: dsaParamsObj}

				privateKeyObj := NewGoRuntimeService("DSA", "DSA", types.ClassNameDSAPrivateKey)
				privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				privateKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: dsaParamsObj}

				keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "DSA", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}
			}
		}

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
			priv, err := ecdsa.GenerateKey(curve, rand.Reader)
			if err == nil {
				params := curve.Params()
				curveObj := NewGoRuntimeService("EC", "EC", types.ClassNameEllipticCurve)
				curveObj.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.P)}
				curveObj.FieldTable["a"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(-3))}
				curveObj.FieldTable["b"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.B)}

				generatorObj := NewGoRuntimeService("EC", "EC", types.ClassNameECPoint)
				generatorObj.FieldTable["x"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.Gx)}
				generatorObj.FieldTable["y"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.Gy)}

				curveObj.FieldTable["generator"] = object.Field{Ftype: types.ECPoint, Fvalue: generatorObj}

				ecSpecObj := NewGoRuntimeService("EC", "EC", types.ClassNameECParameterSpec)
				ecSpecObj.FieldTable["curve"] = object.Field{Ftype: types.Ref, Fvalue: curveObj}
				ecSpecObj.FieldTable["g"] = object.Field{Ftype: types.Ref, Fvalue: generatorObj}
				ecSpecObj.FieldTable["n"] = object.Field{Ftype: types.BigInteger, Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.N)}
				ecSpecObj.FieldTable["h"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}

				pubKey := &ecdsa.PublicKey{Curve: priv.PublicKey.Curve, X: new(big.Int).Set(priv.PublicKey.X), Y: new(big.Int).Set(priv.PublicKey.Y)}
				publicKeyObj := NewGoRuntimeService("EC", "EC", types.ClassNameECPublicKey)
				publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKey}
				publicKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: ecSpecObj}

				privateKeyObj := NewGoRuntimeService("EC", "EC", types.ClassNameECPrivateKey)
				privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				privateKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: ecSpecObj}

				keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "EC", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.ClassNameECPrivateKey, Fvalue: privateKeyObj}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.ClassNameECPublicKey, Fvalue: publicKeyObj}
			}
		}

	case "EdDSA", "Ed25519", "Ed448":
		// --- Handle EdDSA / Ed25519 / Ed448 ---
		var curveName string
		if algorithm == "EdDSA" {
			paramSpecObj, ok := obj.FieldTable["paramSpec"].Fvalue.(*object.Object)
			if !ok {
				return ghelpers.GetGErrBlk(
					excNames.InvalidAlgorithmParameterException,
					"EdDSA requires NamedParameterSpec",
				)
			}
			nameObj, exists := paramSpecObj.FieldTable["name"].Fvalue.(*object.Object)
			if !exists {
				return ghelpers.GetGErrBlk(
					excNames.InvalidAlgorithmParameterException,
					"EdDSA NamedParameterSpec missing name field",
				)
			}
			curveName = object.GoStringFromStringObject(nameObj)
		} else {
			// algorithm == "Ed25519" or "Ed448"
			curveName = algorithm
		}

		// For Ed25519, Ed448:
		// Curve type is the outer switch algorithm.
		// Curve name is the inner switch curveName.
		curveType := algorithm

		switch curveName {
		case "Ed25519":
			pub, priv, err := ed25519.GenerateKey(rand.Reader)
			if err == nil {
				publicKeyObj := NewGoRuntimeService(curveType, curveName, types.ClassNameEdECPublicKey)
				pubCopy := make(ed25519.PublicKey, len(pub))
				copy(pubCopy, pub)
				publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubCopy}

				privateKeyObj := NewGoRuntimeService(curveType, curveName, types.ClassNameEdECPrivateKey)
				privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

				keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "EdDEc", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}
			}
		case "Ed448":
			pub, priv, err := ed448.GenerateKey(rand.Reader)
			if err != nil {
				return ghelpers.GetGErrBlk(
					excNames.GeneralSecurityException,
					"Ed448 key generation failed: "+err.Error(),
				)
			}

			pubKeyCopy := make([]byte, len(pub))
			copy(pubKeyCopy, pub)
			publicKeyObj := NewGoRuntimeService(curveType, curveName, types.ClassNameEdECPublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKeyCopy}

			privKeyCopy := make([]byte, len(priv))
			copy(privKeyCopy, priv)
			privateKeyObj := NewGoRuntimeService(curveType, curveName, types.ClassNameEdECPrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privKeyCopy}

			keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "EdEc", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}

		default:
			return ghelpers.GetGErrBlk(
				excNames.InvalidAlgorithmParameterException,
				"unsupported EdDSA curve: "+curveName,
			)
		}

	case "XDH", "X25519":
		priv := make([]byte, 32)
		_, err = rand.Read(priv)
		if err == nil {
			pub, err := curve25519.X25519(priv, curve25519.Basepoint)
			if err != nil {
				break
			}

			pubKeyCopy := make([]byte, len(pub))
			copy(pubKeyCopy, pub)
			publicKeyObj := NewGoRuntimeService(algorithm, algorithm, types.ClassNameEdECPublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKeyCopy}

			privateKeyObj := NewGoRuntimeService(algorithm, algorithm, types.ClassNameEdECPrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

			keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "EdEc", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}
		}

	case "X448":
		priv := make([]byte, 56)
		_, err = rand.Read(priv)
		if err == nil {
			pub := make([]byte, 56)
			copy(pub, priv)

			publicKeyObj := NewGoRuntimeService(algorithm, algorithm, types.ClassNameEdECPublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}

			privateKeyObj := NewGoRuntimeService(algorithm, algorithm, types.ClassNameEdECPrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

			keyPairObj = NewGoRuntimeService(types.SecurityServiceKeyPairGenerator, "EdEc", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}
		}

	default:
		return ghelpers.GetGErrBlk(
			excNames.GeneralSecurityException,
			"Unrecognized key generation algorithm: "+algorithm,
		)
	}

	if err != nil {
		return ghelpers.GetGErrBlk(
			excNames.GeneralSecurityException,
			"keypairgeneratorGenerateKeyPair: "+algorithm+" key generation failed: "+err.Error(),
		)
	}

	return keyPairObj
}
