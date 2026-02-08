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
	"math/big"

	"golang.org/x/crypto/curve25519"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

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
		keySize = int64(-1)
	}

	var keyPairObj *object.Object

	switch algorithm {
	case "RSA":
		var privRsa *rsa.PrivateKey
		privRsa, err = rsa.GenerateKey(rand.Reader, int(keySize))
		if err == nil {
			publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privRsa.PublicKey)
			if err != nil {
				return ghelpers.GetGErrBlk(
					excNames.GeneralSecurityException,
					"keypairgeneratorGenerateKeyPair: RSA x509.MarshalPKIXPublicKey failed: "+err.Error(),
				)
			}
			pubKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBytes)
			if err != nil {
				return ghelpers.GetGErrBlk(
					excNames.GeneralSecurityException,
					"keypairgeneratorGenerateKeyPair: RSA x509.ParsePKIXPublicKey failed: "+err.Error(),
				)
			}
			pubRsa, ok := pubKeyInterface.(*rsa.PublicKey)
			if !ok {
				return ghelpers.GetGErrBlk(
					excNames.GeneralSecurityException,
					"keypairgeneratorGenerateKeyPair: RSA pubKeyInterface.(*rsa.PublicKey) failed",
				)
			}

			// Finalize the private key.
			privateKeyObj := NewGoRuntimeService("RSAPrivateKey", "RSA", types.ClassNameRSAPrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: privRsa}

			// Finalize the public key.
			publicKeyObj := NewGoRuntimeService("RSAPublicKey", "RSA", types.ClassNameRSAPublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubRsa}

			// Put both private and public into the key pair object.
			keyPairObj = NewGoRuntimeService("KeyPair", "RSA", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.ClassNameRSAPrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.ClassNameRSAPublicKey, Fvalue: publicKeyObj}
		}

	case "DH":
		// Mickey Mouse (constant) simple but valid DH big.Int values
		// p := big.NewInt(23)   // prime modulus
		// g := big.NewInt(5)    // generator
		prv := big.NewInt(6) // private key
		// g**prv mod p (public key)
		pub := big.NewInt(8)

		// Finalize the private key.
		privateKeyObj := NewGoRuntimeService("DHPrivateKey", "DH", types.ClassNameDHPrivateKey)
		privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: prv}

		// Finalize the public key.
		publicKeyObj := NewGoRuntimeService("DHPublicKey", "DH", types.ClassNameDHPublicKey)
		publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pub}

		// Put both private and public into the key pair object.
		keyPairObj = NewGoRuntimeService("KeyPair", "DH", types.ClassNameKeyPair)
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
				// --- Construct DSAParameterSpec object ---
				dsaParamsObj := NewGoRuntimeService("DSAParameterSpec", "DSA", types.ClassNameDSAParameterSpec)

				// Wrap P (prime modulus)
				dsaParamsObj.FieldTable["p"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.P),
				}

				// Wrap Q (subprime)
				dsaParamsObj.FieldTable["q"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.Q),
				}

				// Wrap G (generator)
				dsaParamsObj.FieldTable["g"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.G),
				}

				// --- Create SEPARATE public key with params ---
				pubKey := &dsa.PublicKey{
					Parameters: *params,
					Y:          new(big.Int).Set(priv.PublicKey.Y), // Copy Y value
				}

				publicKeyObj := NewGoRuntimeService("DSAPublicKey", "DSA", types.ClassNameDSAPublicKey)
				publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKey}
				publicKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: dsaParamsObj}

				// --- Create private key with params ---
				privateKeyObj := NewGoRuntimeService("DSAPrivateKey", "DSA", types.ClassNameDSAPrivateKey)
				privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				privateKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: dsaParamsObj}

				// --- Create KeyPair object ---
				keyPairObj = NewGoRuntimeService("KeyPair", "DSA", types.ClassNameKeyPair)
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
			// Generate ECDSA key
			var priv *ecdsa.PrivateKey
			priv, err = ecdsa.GenerateKey(curve, rand.Reader)
			if err == nil {
				// Get curve parameters
				params := curve.Params()

				// --- Construct EllipticCurve object ---
				curveObj := NewGoRuntimeService("EllipticCurve", "EC", types.ClassNameEllipticCurve)

				// Wrap P (prime field)
				curveObj.FieldTable["p"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.P),
				}

				// Wrap A (coefficient, -3 for NIST curves)
				curveObj.FieldTable["a"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(-3)),
				}

				// Wrap B (coefficient)
				curveObj.FieldTable["b"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.B),
				}

				// --- Create generator point (G) ---
				generatorObj := NewGoRuntimeService("ECPoint", "", types.ClassNameECPoint)

				// Wrap Gx
				generatorObj.FieldTable["x"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.Gx),
				}

				// Wrap Gy
				generatorObj.FieldTable["y"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.Gy),
				}

				curveObj.FieldTable["generator"] = object.Field{Ftype: types.ECPoint, Fvalue: generatorObj}

				// --- Construct ECParameterSpec object ---
				ecSpecObj := NewGoRuntimeService("ECParameterSpec", "EC", types.ClassNameECParameterSpec)
				ecSpecObj.FieldTable["curve"] = object.Field{Ftype: types.Ref, Fvalue: curveObj}
				ecSpecObj.FieldTable["g"] = object.Field{Ftype: types.Ref, Fvalue: generatorObj}

				// Wrap N (order)
				ecSpecObj.FieldTable["n"] = object.Field{
					Ftype:  types.BigInteger,
					Fvalue: object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, params.N),
				}

				ecSpecObj.FieldTable["h"] = object.Field{Ftype: types.Int, Fvalue: int64(1)} // cofactor

				// --- Create SEPARATE public key with params ---
				pubKey := &ecdsa.PublicKey{
					Curve: priv.PublicKey.Curve,
					X:     new(big.Int).Set(priv.PublicKey.X), // Copy X coordinate
					Y:     new(big.Int).Set(priv.PublicKey.Y), // Copy Y coordinate
				}

				// --- Create public key with params ---
				publicKeyObj := NewGoRuntimeService("ECPublicKey", "EC", types.ClassNameECPublicKey)
				publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKey}
				publicKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: ecSpecObj}

				// --- Create private key with params ---
				privateKeyObj := NewGoRuntimeService("ECPrivateKey", "EC", types.ClassNameECPrivateKey)
				privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				privateKeyObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: ecSpecObj}

				// Put both private and public key objects into the key pair object.
				keyPairObj = NewGoRuntimeService("KeyPair", "EC", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.ClassNameECPrivateKey, Fvalue: privateKeyObj}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.ClassNameECPublicKey, Fvalue: publicKeyObj}
			}
		}

	case "Ed25519":
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err == nil {
			// Ed25519 has no parameters like EC or DSA - it's a fixed curve
			// The keys are just byte slices, not structs with embedded public keys

			// --- Create separate public key copy ---
			pubKeyCopy := make(ed25519.PublicKey, len(pub))
			copy(pubKeyCopy, pub)

			publicKeyObj := NewGoRuntimeService("Ed25519PublicKey", "Ed25519", types.ClassNameEd25519PublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKeyCopy}

			// --- Create private key ---
			privateKeyObj := NewGoRuntimeService("Ed25519PrivateKey", "Ed25519", types.ClassNameEd25519PrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

			// --- Create KeyPair object ---
			keyPairObj = NewGoRuntimeService("KeyPair", "Ed25519", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}
		}

	case "XDH", "X25519":
		priv := make([]byte, 32)
		_, err = rand.Read(priv)
		if err == nil {
			pub, err := curve25519.X25519(priv, curve25519.Basepoint)
			if err != nil {
				break
			}

			// --- Create separate public key copy ---
			pubKeyCopy := make([]byte, len(pub))
			copy(pubKeyCopy, pub)

			publicKeyObj := NewGoRuntimeService("X25519PublicKey", "XDH", types.ClassNameX25519PublicKey)
			publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKeyCopy}

			// --- Create private key ---
			privateKeyObj := NewGoRuntimeService("X25519PrivateKey", "XDH", types.ClassNameX25519PrivateKey)
			privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

			// --- Create KeyPair object ---
			keyPairObj = NewGoRuntimeService("KeyPair", algorithm, types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}
		}

	case "Ed448":
		// Ed448 not in standard Go; use placeholder for now
		priv := make([]byte, 57) // Ed448 private key length
		_, err = rand.Read(priv)
		if err != nil {
			break
		}
		pub := make([]byte, 57)
		copy(pub, priv) // placeholder: real Ed448 requires proper library

		// --- Create separate public key copy ---
		pubKeyCopy := make([]byte, len(pub))
		copy(pubKeyCopy, pub)

		publicKeyObj := NewGoRuntimeService("Ed448PublicKey", "Ed448", types.ClassNameEd448PublicKey)
		publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKeyCopy}

		// --- Create private key ---
		privateKeyObj := NewGoRuntimeService("Ed448PrivateKey", "Ed448", types.ClassNameEd448PrivateKey)
		privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

		// --- Create KeyPair object ---
		keyPairObj = NewGoRuntimeService("KeyPair", "Ed448", types.ClassNameKeyPair)
		keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
		keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}

	case "X448":
		// X448 not in standard Go; use placeholder for now
		priv := make([]byte, 56) // X448 private key length
		_, err = rand.Read(priv)
		if err != nil {
			break
		}
		pub := make([]byte, 56) // X448 public key length
		copy(pub, priv)         // placeholder: real X448 requires proper library

		// --- Create separate public key copy ---
		pubKeyCopy := make([]byte, len(pub))
		copy(pubKeyCopy, pub)

		publicKeyObj := NewGoRuntimeService("X448PublicKey", "X448", types.ClassNameX448PublicKey)
		publicKeyObj.FieldTable["value"] = object.Field{Ftype: types.PublicKey, Fvalue: pubKeyCopy}

		// --- Create private key ---
		privateKeyObj := NewGoRuntimeService("X448PrivateKey", "X448", types.ClassNameX448PrivateKey)
		privateKeyObj.FieldTable["value"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}

		// --- Create KeyPair object ---
		keyPairObj = NewGoRuntimeService("KeyPair", "X448", types.ClassNameKeyPair)
		keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privateKeyObj}
		keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: publicKeyObj}

	default:
		return ghelpers.GetGErrBlk(
			excNames.GeneralSecurityException,
			"keypairgeneratorGenerateKeyPair: Unrecognized key generation algorithm: "+algorithm,
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
