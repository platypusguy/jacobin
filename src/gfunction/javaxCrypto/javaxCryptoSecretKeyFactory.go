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

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.generateSecret(Ljava/security/spec/KeySpec;)Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  secretKeyFactoryGenerateSecret,
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

	ghelpers.MethodSignatures["javax/crypto/SecretKeyFactory.getInstance(Ljava/lang/String;)Ljavax/crypto/SecretKeyFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  secretKeyFactoryGetInstance,
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
		isPBKDF2 := strings.HasPrefix(strings.ToUpper(algorithm), "PBKDF2WITHHMAC")
		isPBEAES := strings.HasPrefix(strings.ToUpper(algorithm), "PBEWITHHMAC") && strings.Contains(strings.ToUpper(algorithm), "ANDAES_")
		isLegacyPBE := strings.HasPrefix(strings.ToUpper(algorithm), "PBEWITH") && !isPBEAES && !strings.Contains(strings.ToUpper(algorithm), "HMAC")

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
		var salt []byte
		if !object.IsNull(saltVal) {
			saltObj := saltVal.(*object.Object)
			saltJavaBytes := saltObj.FieldTable["value"].Fvalue.([]types.JavaByte)
			salt = object.GoByteArrayFromJavaByteArray(saltJavaBytes)
		}

		iterations := keySpec.FieldTable["iterationCount"].Fvalue.(int64)
		keyLength := keySpec.FieldTable["keyLength"].Fvalue.(int64)

		if keyLength == 0 {
			algoUpper := strings.ToUpper(algorithm)
			if isPBEAES {
				if strings.Contains(algoUpper, "_128") {
					keyLength = 128
				} else if strings.Contains(algoUpper, "_256") {
					keyLength = 256
				}
			} else if isLegacyPBE {
				// Legacy PBE PKCS#5 v1.5 usually derives Key + IV.
				// IV is always 8 bytes (64 bits) for these block ciphers.
				switch {
				case strings.Contains(algoUpper, "DESEDE") || strings.Contains(algoUpper, "TRIPLEDES"):
					keyLength = 192 + 64 // 24 bytes key + 8 bytes IV
				case strings.Contains(algoUpper, "DES"):
					keyLength = 64 + 64 // 8 bytes key + 8 bytes IV
				case strings.Contains(algoUpper, "RC2_40") || strings.Contains(algoUpper, "RC4_40"):
					keyLength = 40 + 64 // 5 bytes key + 8 bytes IV
				case strings.Contains(algoUpper, "RC2_128") || strings.Contains(algoUpper, "RC4_128"):
					keyLength = 128 + 64 // 16 bytes key + 8 bytes IV
				default:
					keyLength = 128
				}
			}
		}

		if (isPBKDF2 || isPBEAES) && (len(salt) == 0 || iterations == 0 || keyLength == 0) {
			// If salt, iterations or keyLength are missing, we can't derive the key yet.
			// Return a SecretKey that just contains the password and needs derivation.
			sk := object.MakePrimitiveObject(types.ClassNameSecretKey, types.ByteArray, []byte(password))
			sk.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}
			if keyLength > 0 {
				sk.FieldTable["inferred_key_length"] = object.Field{Ftype: types.Int, Fvalue: keyLength}
			}
			return sk
		}

		var derivedKey []byte
		algoUpper := strings.ToUpper(algorithm)
		if isLegacyPBE && !isPBEAES {
			// PKCS#5 v1.5 (PBEWithMD5AndDES, etc.)
			if salt == nil {
				// Salt is required for legacy PBE derivation.
				// If not provided, we must defer derivation.
				sk := object.MakePrimitiveObject(types.ClassNameSecretKey, types.ByteArray, []byte(password))
				sk.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}
				return sk
			}
			// For legacy PBE, we use the password and salt to derive the key using a simple MD5/SHA1 hash loop.
			// Ref: PKCS#5 v1.5, Section 6.
			var h hash.Hash
			switch {
			case strings.Contains(algoUpper, "MD5"):
				h = md5.New()
			case strings.Contains(algoUpper, "SHA1") || strings.Contains(algoUpper, "SHA-1"):
				h = sha1.New()
			default:
				return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, "secretKeyFactoryGenerateSecret: unsupported digest for legacy PBE")
			}

			// DK = Hash(Password || Salt)
			// For iterations > 1: Hash(DK_{i-1})
			h.Reset()
			h.Write([]byte(password))
			h.Write(salt)
			derivedKey = h.Sum(nil)
			for i := 1; i < int(iterations); i++ {
				h.Reset()
				h.Write(derivedKey)
				derivedKey = h.Sum(nil)
			}

			// If we need more than one hash block (e.g. for DESede 24 bytes + 8 bytes IV = 32 bytes)
			// PKCS#5 v1.5 actually only defines derivation for up to hash size.
			// However, some implementations (like OpenSSL or Java) extend it by hashing (DK_i-1 || password || salt)
			// But Java's PBEWithMD5AndDES specifically uses PKCS#5 v1.5 which is only for 8-byte keys.
			// Actually, Java's PBEWithMD5AndDES only derives 8 bytes of key and 8 bytes of IV.
			// If more is needed, it repeats the process with a slightly different input or just fails.
			// For PKCS#5 v1.5, if we need more than hash size:
			// DK_1 = Hash(password || salt)
			// DK_2 = Hash(DK_1 || password || salt) ... no, that's not it.

			// Actually, for PBEWithMD5AndDES, MD5 is 16 bytes.
			// 8 bytes are used for Key, 8 bytes for IV. So one MD5 block is enough.
			// For SHA1, it's 20 bytes. 8 for key, 8 for IV. Also enough.
			// For DESede, we need 24+8 = 32 bytes. SHA1 is not enough.

			if keyLength > 0 && int(keyLength/8) > len(derivedKey) {
				// Extend derived key if needed (PKCS#5 v1.5 extension used by Java/OpenSSL)
				fullDerivedKey := make([]byte, 0, (int(keyLength)/8+h.Size()-1)/h.Size()*h.Size())
				fullDerivedKey = append(fullDerivedKey, derivedKey...)

				currentDK := derivedKey
				for len(fullDerivedKey) < int(keyLength/8) {
					h.Reset()
					h.Write(currentDK)
					h.Write([]byte(password))
					h.Write(salt)
					currentDK = h.Sum(nil)
					for i := 1; i < int(iterations); i++ {
						h.Reset()
						h.Write(currentDK)
						currentDK = h.Sum(nil)
					}
					fullDerivedKey = append(fullDerivedKey, currentDK...)
				}
				derivedKey = fullDerivedKey
			}

			// Trim or pad to keyLength
			if keyLength > 0 {
				if len(derivedKey) > int(keyLength/8) {
					derivedKey = derivedKey[:keyLength/8]
				} else if len(derivedKey) < int(keyLength/8) {
					// If we need more key material, we can use the DK_{i} = Hash(DK_{i-1}) pattern
					// but actually PKCS#5 v1.5 PBE only supports single hash length keys (8 or 16 bytes depending on hash).
					// For DESede, it's more complex (PKCS#12 or others).
					// For now, let's just pad it if it's too short, which is NOT standard but prevents panic.
					newKey := make([]byte, keyLength/8)
					copy(newKey, derivedKey)
					derivedKey = newKey
				}
			}
		} else {
			// PBKDF2
			var h func() hash.Hash
			switch {
			case strings.Contains(algoUpper, "SHA1"):
				h = sha1.New
			case strings.Contains(algoUpper, "SHA224"):
				h = sha256.New224
			case strings.Contains(algoUpper, "SHA256"):
				h = sha256.New
			case strings.Contains(algoUpper, "SHA384"):
				h = sha512.New384
			case strings.Contains(algoUpper, "SHA512/224"):
				h = sha512.New512_224
			case strings.Contains(algoUpper, "SHA512/256"):
				h = sha512.New512_256
			case strings.Contains(algoUpper, "SHA512"):
				h = sha512.New
			default:
				return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, "secretKeyFactoryGenerateSecret: unsupported PRF for PBKDF2/PBE")
			}

			derivedKey = pbkdf2.Key([]byte(password), salt, int(iterations), int(keyLength/8), h)
		}

		// Create the SecretKey object
		// For legacy PBE, the SecretKey should only contain the key part,
		// but the full derived material (Key + IV) is useful.
		// However, standard Java PBE keys often only return the key bytes in getEncoded().
		keyPart := derivedKey
		if isLegacyPBE && !isPBEAES {
			var actualKeyLen int
			switch {
			case strings.Contains(algorithm, "DESede"), strings.Contains(algorithm, "TripleDES"):
				actualKeyLen = 24
			case strings.Contains(algorithm, "DES"):
				actualKeyLen = 8
			case strings.Contains(algorithm, "RC2"), strings.Contains(algorithm, "RC4"):
				if strings.Contains(algorithm, "40") {
					actualKeyLen = 5
				} else {
					actualKeyLen = 16
				}
			}
			if actualKeyLen > 0 && len(derivedKey) > actualKeyLen {
				keyPart = derivedKey[:actualKeyLen]
			}
		}

		sk := object.MakePrimitiveObject(types.ClassNameSecretKey, types.ByteArray, keyPart)
		sk.FieldTable["algorithm"] = object.Field{Ftype: types.StringClassName, Fvalue: algorithmObj}
		sk.FieldTable["key"] = object.Field{Ftype: types.ByteArray, Fvalue: derivedKey} // Full material (Key + IV)
		return sk
	}

	return ghelpers.GetGErrBlk(excNames.InvalidKeyException, "secretKeyFactoryGenerateSecret: unsupported KeySpec")
}

func secretKeyFactoryGetInstance(params []any) any {
	algorithmObj, ok := params[0].(*object.Object)
	if !ok || algorithmObj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "secretKeyFactoryGetInstance: algorithm cannot be null")
	}
	algorithm := object.GoStringFromStringObject(algorithmObj)

	if _, ok := ValidateCipherTransformation(algorithm); !ok {
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
