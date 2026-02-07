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
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"keypairgeneratorGenerateKeyPair: keySize not set",
		)
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
			privObj := object.MakePrimitiveObject(types.ClassNamePrivateKey, types.PrivateKey, privRsa)
			pubObj := object.MakePrimitiveObject(types.ClassNamePublicKey, types.PublicKey, pubRsa)
			keyPairObj = NewGoRuntimeService("KeyPair", "RSA", types.ClassNameKeyPair)
			keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: privObj}
			keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: pubObj}
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
			// Generate ECDSA key
			priv, err2 := ecdsa.GenerateKey(curve, rand.Reader)
			if err2 != nil {
				err = err2
			} else {
				// Jacobin KeyPair object
				keyPairObj = NewGoRuntimeService("KeyPair", "EC", types.ClassNameKeyPair)
				keyPairObj.FieldTable["private"] = object.Field{Ftype: types.PrivateKey, Fvalue: priv}
				keyPairObj.FieldTable["public"] = object.Field{Ftype: types.PublicKey, Fvalue: &priv.PublicKey}

				// --- Construct EllipticCurve object ---
				curveObj := NewGoRuntimeService("EllipticCurve", "EC", types.ClassNameEllipticCurve)
				curveObj.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: &object.Object{ /*wrap curve.Params().P*/ }}
				curveObj.FieldTable["a"] = object.Field{Ftype: types.BigInteger, Fvalue: &object.Object{ /* wrap -3 */ }}
				curveObj.FieldTable["b"] = object.Field{Ftype: types.BigInteger, Fvalue: &object.Object{ /* wrap curve.Params().B */ }}
				curveObj.FieldTable["generator"] = object.Field{Ftype: types.ECPoint, Fvalue: NewGoRuntimeService("ECPoint", "", types.ClassNameECPoint)}
				curveObj.FieldTable["generator"].Fvalue.(*object.Object).FieldTable["x"] = object.Field{Ftype: types.BigInteger, Fvalue: &object.Object{ /* wrap curve.Params().Gx */ }}
				curveObj.FieldTable["generator"].Fvalue.(*object.Object).FieldTable["y"] = object.Field{Ftype: types.BigInteger, Fvalue: &object.Object{ /* wrap curve.Params().Gy */ }}

				// --- Construct ECParameterSpec object ---
				ecSpecObj := NewGoRuntimeService("ECParameterSpec", "EC", types.ClassNameECParameterSpec)
				ecSpecObj.FieldTable["curve"] = object.Field{Ftype: types.Ref, Fvalue: curveObj}
				ecSpecObj.FieldTable["g"] = object.Field{Ftype: types.Ref, Fvalue: curveObj.FieldTable["generator"].Fvalue}
				ecSpecObj.FieldTable["n"] = object.Field{Ftype: types.BigInteger, Fvalue: &object.Object{ /* wrap curve.Params().N */ }}
				ecSpecObj.FieldTable["h"] = object.Field{Ftype: types.Int, Fvalue: 1} // cofactor is 1

				// Attach to public key
				keyPairObj.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: ecSpecObj}
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
