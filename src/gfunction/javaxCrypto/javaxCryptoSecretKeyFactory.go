/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Crypto_SecretKeyFactory() {
	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.getInstance(Ljava/lang/String;)Ljavax/crypto/SecretKeyFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  secretKeyFactoryGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/crypto/SecretKeyFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  secretKeyFactoryGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/crypto/SecretKeyFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  secretKeyFactoryGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.generateSecret(Ljava/security/spec/KeySpec;)Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  secretKeyFactoryGenerateSecret,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.getKeySpec(Ljavax/crypto/SecretKey;Ljava/lang/Class;)Ljava/security/spec/KeySpec;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.translateKey(Ljavax/crypto/SecretKey;)Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}

func secretKeyFactoryGetInstance(params []any) any {
	algorithmObj, ok := params[0].(*object.Object)
	if !ok || algorithmObj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "secretKeyFactoryGetInstance: algorithm cannot be null")
	}
	algorithm := object.GoStringFromStringObject(algorithmObj)

	if _, ok := CipherConfigTable[algorithm]; !ok {
		return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, fmt.Sprintf("secretKeyFactoryGetInstance: %s not found", algorithm))
	}

	// Check provider if provided
	if len(params) > 1 {
		p := params[1]
		if pObj, ok := p.(*object.Object); ok && pObj != nil {
			// provider can be String or Provider object
			var pName string
			if pObj.KlassName == object.StringPoolIndexFromGoString(types.StringClassName) {
				pName = object.GoStringFromStringObject(pObj)
			} else {
				// Assume it's a Provider object
				nameField, ok := pObj.FieldTable["name"]
				if ok {
					pName = object.GoStringFromStringObject(nameField.Fvalue.(*object.Object))
				}
			}

			if pName != "" && pName != types.SecurityProviderName {
				return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException, fmt.Sprintf("secretKeyFactoryGetInstance: provider %s not found", pName))
			}
		}
	}

	skf := object.MakeEmptyObjectWithClassName(new("javax/crypto/SecretKeyFactory"))

	skf.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: algorithmObj,
	}

	providerObj := ghelpers.GetDefaultSecurityProvider()
	skf.FieldTable["provider"] = object.Field{
		Ftype:  types.ClassNameSecurityProvider,
		Fvalue: providerObj,
	}

	// We don't strictly need a separate SPI object for now as we trap the methods on the factory itself
	return skf
}

func secretKeyFactoryGenerateSecret(params []any) any {
	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "secretKeyFactoryGenerateSecret: invalid 'this'")
	}

	keySpec, ok := params[1].(*object.Object)
	if !ok || keySpec == nil {
		return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: keySpec cannot be null")
	}

	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Implementation depends on what KeySpecs we support.
	// For now, let's look for SecretKeySpec which is common.
	if keySpec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/SecretKeySpec") {
		// SecretKeySpec is already a SecretKey, so we can just return it if algorithms match.
		specAlgoObj := keySpec.FieldTable["algorithm"].Fvalue.(*object.Object)
		specAlgo := object.GoStringFromStringObject(specAlgoObj)
		if specAlgo != algorithm {
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException, fmt.Sprintf("secretKeyFactoryGenerateSecret: algorithm mismatch: %s vs %s", specAlgo, algorithm))
		}
		return keySpec
	}

	// javax/crypto/spec/PBEKeySpec
	if keySpec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/PBEKeySpec") {
		isPBKDF2 := strings.HasPrefix(algorithm, "PBKDF2WithHmac")
		isPBEAES := strings.HasPrefix(algorithm, "PBEWithHmac") && strings.Contains(algorithm, "AndAES_")
		isLegacyPBE := strings.HasPrefix(algorithm, "PBEWith") && !isPBEAES

		if !isPBKDF2 && !isPBEAES && !isLegacyPBE {
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: PBEKeySpec only supported for PBKDF2 or PBEWithHmac*AndAES algorithms")
		}

		passwordVal := keySpec.FieldTable["password"].Fvalue
		if object.IsNull(passwordVal) {
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: password cannot be null")
		}
		passwordObj := passwordVal.(*object.Object)
		passwordChars := passwordObj.FieldTable["value"].Fvalue.([]int64)
		password := object.GoStringFromJavaCharArray(passwordChars)

		saltVal := keySpec.FieldTable["salt"].Fvalue
		if object.IsNull(saltVal) && (isPBKDF2 || isPBEAES) {
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: salt cannot be null")
		}
		var salt []byte
		if !object.IsNull(saltVal) {
			saltObj := saltVal.(*object.Object)
			saltJavaBytes := saltObj.FieldTable["value"].Fvalue.([]types.JavaByte)
			salt = object.GoByteArrayFromJavaByteArray(saltJavaBytes)
		}

		iterations := keySpec.FieldTable["iterationCount"].Fvalue.(int64)
		keyLength := keySpec.FieldTable["keyLength"].Fvalue.(int64)

		if keyLength == 0 {
			if isPBEAES {
				if strings.HasSuffix(algorithm, "_128") {
					keyLength = 128
				} else if strings.HasSuffix(algorithm, "_256") {
					keyLength = 256
				}
			} else if isLegacyPBE {
				// Legacy PBE usually has fixed key lengths based on the algorithm
				switch {
				case strings.Contains(algorithm, "DESede"), strings.Contains(algorithm, "TripleDES"):
					keyLength = 192
				case strings.Contains(algorithm, "DES"):
					keyLength = 64
				case strings.Contains(algorithm, "RC2_40"), strings.Contains(algorithm, "RC4_40"):
					keyLength = 40
				case strings.Contains(algorithm, "RC2_128"), strings.Contains(algorithm, "RC4_128"):
					keyLength = 128
				}
			}
		}

		if keyLength == 0 && (isPBKDF2 || isPBEAES) {
			// Some default or based on algorithm? Java docs say it should be specified for PBKDF2
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: keyLength must be specified if not implied by algorithm")
		}

		var derivedKey []byte
		if isLegacyPBE && !isPBEAES {
			// PKCS#5 v1.5 (PBEWithMD5AndDES, etc.)
			// For legacy PBE, we use the password and salt to derive the key using a simple MD5/SHA1 hash loop.
			// Ref: PKCS#5 v1.5, Section 6.
			var h hash.Hash
			switch {
			case strings.Contains(algorithm, "MD5"):
				h = md5.New()
			case strings.Contains(algorithm, "SHA1"):
				h = sha1.New()
			default:
				return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, "secretKeyFactoryGenerateSecret: unsupported digest for legacy PBE")
			}

			// DK = Hash(Password || Salt)
			// For iterations > 1: Hash(DK_{i-1})
			h.Write([]byte(password))
			h.Write(salt)
			derivedKey = h.Sum(nil)
			for i := 1; i < int(iterations); i++ {
				h.Reset()
				h.Write(derivedKey)
				derivedKey = h.Sum(nil)
			}

			// Trim to keyLength
			if keyLength > 0 && len(derivedKey) > int(keyLength/8) {
				derivedKey = derivedKey[:keyLength/8]
			}
		} else {
			// PBKDF2
			var h func() hash.Hash
			switch {
			case strings.Contains(algorithm, "SHA1"):
				h = sha1.New
			case strings.Contains(algorithm, "SHA224"):
				h = sha256.New224
			case strings.Contains(algorithm, "SHA256"):
				h = sha256.New
			case strings.Contains(algorithm, "SHA384"):
				h = sha512.New384
			case strings.Contains(algorithm, "SHA512"):
				h = sha512.New
			default:
				return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, "secretKeyFactoryGenerateSecret: unsupported PRF for PBKDF2/PBE")
			}

			derivedKey = pbkdf2.Key([]byte(password), salt, int(iterations), int(keyLength/8), h)
		}

		// Create the SecretKey object
		sk := object.MakePrimitiveObject(types.ClassNameSecretKey, types.ByteArray, derivedKey)
		sk.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}
		return sk
	}

	return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: unsupported KeySpec")
}
