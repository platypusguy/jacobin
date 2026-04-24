/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaSecurity"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

func Load_Crypto_Cipher() {
	ghelpers.MethodSignatures["javax/crypto/Cipher.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherClinit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal(Ljava/nio/ByteBuffer;Ljava/nio/ByteBuffer;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal([BII[BI)I"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal([BII[B)I"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal([BII)[B"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal([BI)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.doFinal([B)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherDoFinal,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherGetAlgorithm,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getBlockSize()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherGetBlockSize,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getExemptionMechanism()Ljavax/crypto/ExemptionMechanism;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getIV()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherGetIV,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getInstance(Ljava/lang/String;)Ljavax/crypto/Cipher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/crypto/Cipher;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  cipherGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/crypto/Cipher;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  cipherGetInstance,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getMaxAllowedKeyLength(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getMaxAllowedParameterSpec(Ljava/lang/String;)Ljava/security/spec/AlgorithmParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherGetMaxAllowedParameterSpec,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getOutputSize(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherGetOutputSize,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getParameters()Ljava/security/AlgorithmParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherGetParameters,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  cipherGetProvider,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/Key;Ljava/security/AlgorithmParameters;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/Key;Ljava/security/AlgorithmParameters;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/Key;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/Key;Ljava/security/spec/AlgorithmParameterSpec;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/Key;Ljava/security/spec/AlgorithmParameterSpec;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/Key;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/cert/Certificate;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.init(ILjava/security/cert/Certificate;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  cipherInit,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.unwrap([BLjava/lang/String;I)Ljava/security/Key;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.update(Ljava/nio/ByteBuffer;Ljava/nio/ByteBuffer;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.update([BII[BI)I"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  cipherUpdate,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.update([BII[B)I"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  cipherUpdate,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.update([BII)[B"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  cipherUpdate,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.update([B)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherUpdate,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.updateAAD(Ljava/nio/ByteBuffer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.updateAAD([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.updateAAD([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.wrap(Ljava/security/Key;)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}

func cipherClinit(params []any) any {
	_ = statics.AddStatic(types.ClassNameCipher+".DECRYPT_MODE", statics.Static{Type: types.Int, Value: int64(2)})
	_ = statics.AddStatic(types.ClassNameCipher+".ENCRYPT_MODE", statics.Static{Type: types.Int, Value: int64(1)})
	_ = statics.AddStatic(types.ClassNameCipher+".PRIVATE_KEY", statics.Static{Type: types.Int, Value: int64(2)})
	_ = statics.AddStatic(types.ClassNameCipher+".PUBLIC_KEY", statics.Static{Type: types.Int, Value: int64(1)})
	_ = statics.AddStatic(types.ClassNameCipher+".SECRET_KEY", statics.Static{Type: types.Int, Value: int64(3)})
	_ = statics.AddStatic(types.ClassNameCipher+".UNWRAP_MODE", statics.Static{Type: types.Int, Value: int64(4)})
	_ = statics.AddStatic(types.ClassNameCipher+".WRAP_MODE", statics.Static{Type: types.Int, Value: int64(3)})
	return nil
}

func cipherDoFinal(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	var input []byte
	if len(params) > 1 {
		inputObj, ok := params[1].(*object.Object)
		if ok && !object.IsNull(inputObj) {
			input = object.GoByteArrayFromJavaByteArray(inputObj.FieldTable["value"].Fvalue.([]types.JavaByte))
			if len(params) > 3 {
				offset := params[2].(int64)
				length := params[3].(int64)
				input = input[offset : offset+length]
			}
		}
	}

	buffer, _ := self.FieldTable["buffer"].Fvalue.([]byte)
	fullInput := append(buffer, input...)
	self.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: []byte{}}

	config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "cipherDoFinal: config missing")
	}

	opmodeField, ok := self.FieldTable["opmode"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "cipherDoFinal: cipher not initialized")
	}
	opmode := opmodeField.Fvalue.(int64)

	keyField, ok := self.FieldTable["key"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "cipherDoFinal: key missing")
	}
	keyObj := keyField.Fvalue.(*object.Object)
	var keyBytes []byte

	if kb, ok := keyObj.FieldTable["key"]; ok {
		keyBytes, _ = kb.Fvalue.([]byte)
		config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
		if ok && strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC") {
			// Legacy PBE keys generated via SecretKeyFactory may contain both Key and IV.
			// Extract just the key part for the underlying block cipher.
			switch {
			case strings.Contains(config.Name, "DESede") || strings.Contains(config.Name, "TripleDES"):
				if len(keyBytes) >= 24 {
					keyBytes = keyBytes[:24]
				}
			case strings.Contains(config.Name, "DES"):
				if len(keyBytes) >= 8 {
					keyBytes = keyBytes[:8]
				}
			case strings.Contains(config.Name, "RC2"), strings.Contains(config.Name, "RC4"):
				bits := 128
				if strings.Contains(config.Name, "40") {
					bits = 40
				}
				if len(keyBytes) >= bits/8 {
					keyBytes = keyBytes[:bits/8]
				}
			}
		}
	} else if vb, ok := keyObj.FieldTable["value"]; ok {
		// Handle cases where the key is just a byte array object (common in tests)
		if jBytes, ok := vb.Fvalue.([]types.JavaByte); ok {
			keyBytes = object.GoByteArrayFromJavaByteArray(jBytes)
		} else if bBytes, ok := vb.Fvalue.([]byte); ok {
			keyBytes = bBytes
		}
		config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
		if ok && strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC") {
			// If we fell back to the "value" field, and it's a legacy PBE,
			// this is likely the raw password which is NOT the derived key.
			// We MUST NOT use it as the DES/RC2 key.
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException, fmt.Sprintf("cipherDoFinal: derived key missing for PBE (found raw password of length %d instead)", len(keyBytes)))
		}
	}

	if keyBytes == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "cipherDoFinal: invalid key or key bytes missing")
	}

	var iv []byte
	ivField, ok := self.FieldTable["iv"]
	if ok && !object.IsNull(ivField.Fvalue) {
		switch v := ivField.Fvalue.(type) {
		case []types.JavaByte:
			iv = object.GoByteArrayFromJavaByteArray(v)
		case []byte:
			iv = v
		case *object.Object:
			if v.KlassName == object.StringPoolIndexFromGoString(types.ByteArray) {
				iv = object.GoByteArrayFromJavaByteArray(v.FieldTable["value"].Fvalue.([]types.JavaByte))
			}
		}
	}

	result, err := performCipher(config, opmode, keyBytes, iv, fullInput)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.GeneralSecurityException, "cipherDoFinal: "+err.Error())
	}

	// Handle methods that return I (length written to provided [B)
	// javax/crypto/Cipher.doFinal([BII[BI)I
	// javax/crypto/Cipher.doFinal([BII[B)I
	// javax/crypto/Cipher.doFinal([BI)I
	// javax/crypto/Cipher.doFinal([B)I

	// The G-function signature determines what we should return.
	// However, we can use a heuristic or just implement a more robust dispatcher.
	// For now, let's check if there's an output buffer provided in the params.

	var outputBuf *object.Object
	var outputOffset int64

	// Heuristic to find output buffer and offset based on param count and types
	// [this, input, inOff, inLen, outBuf, outOff] -> 6 params
	// [this, input, inOff, inLen, outBuf] -> 5 params
	// [this, outBuf, outOff] -> 3 params
	// [this, outBuf] -> 2 params

	if len(params) >= 2 {
		if len(params) == 6 { // doFinal([BII[BI)I
			outputBuf = params[4].(*object.Object)
			outputOffset = params[5].(int64)
		} else if len(params) == 5 { // doFinal([BII[B)I
			outputBuf = params[4].(*object.Object)
			outputOffset = 0
		} else if len(params) == 3 { // doFinal([BI)I
			outputBuf = params[1].(*object.Object)
			outputOffset = params[2].(int64)
		} else if len(params) == 2 {
			// Could be doFinal([B)I or doFinal([B)[B
			// If the second param is an object and it's a byte array, it might be the output buffer.
			// But it also might be the input buffer.
			// This is tricky without knowing the exact method signature being called.
			// However, in Jacobin G-functions, we can usually tell from the ParamSlots.
		}
	}

	if outputBuf != nil && !object.IsNull(outputBuf) {
		dest := outputBuf.FieldTable["value"].Fvalue.([]types.JavaByte)
		if int(outputOffset)+len(result) > len(dest) {
			return ghelpers.GetGErrBlk(excNames.ShortBufferException, "cipherDoFinal: output buffer too short")
		}
		jResult := object.JavaByteArrayFromGoByteArray(result)
		copy(dest[outputOffset:], jResult)
		outputBuf.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: dest}
		return int64(len(result))
	}

	jBytes := object.JavaByteArrayFromGoByteArray(result)
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}

func cipherGetAlgorithm(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	transformationField, ok := self.FieldTable["transformation"]
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	return transformationField.Fvalue
}

func cipherGetBlockSize(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return int64(0)
	}

	config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
	if !ok {
		return int64(0)
	}

	// Determine block size based on algorithm
	switch config.Algorithm {
	case "AES":
		return int64(16)
	case "DES", "DESede", "Blowfish":
		return int64(8)
	default:
		return int64(0)
	}
}

func cipherGetInstance(params []any) any {
	transformationObj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(transformationObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "cipherGetInstance: transformation cannot be null")
	}
	transformation := object.GoStringFromStringObject(transformationObj)

	config, ok := ValidateCipherTransformation(transformation)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, fmt.Sprintf("cipherGetInstance: %s not found", transformation))
	}

	if !config.Enabled {
		return ghelpers.GetGErrBlk(excNames.NoSuchAlgorithmException, fmt.Sprintf("cipherGetInstance: %s is disabled", transformation))
	}

	// Check provider parameter if provided
	if len(params) > 1 {
		p := params[1]
		if pObj, ok := p.(*object.Object); ok && pObj != nil {
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
				return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException, fmt.Sprintf("cipherGetInstance: provider %s not found", pName))
			}
		}
	}

	cipher := object.MakeEmptyObjectWithClassName(new(types.ClassNameCipher))

	cipher.FieldTable["transformation"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: transformationObj,
	}

	providerObj := ghelpers.GetDefaultSecurityProvider()
	cipher.FieldTable["provider"] = object.Field{
		Ftype:  types.ClassNameSecurityProvider,
		Fvalue: providerObj,
	}

	// Internal state to hold transformation info
	cipher.FieldTable["config"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: config,
	}

	return cipher
}

func cipherGetIV(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	ivField, ok := self.FieldTable["iv"]
	if !ok || object.IsNull(ivField.Fvalue) {
		return ghelpers.ReturnNull(nil)
	}

	// ivField.Fvalue can be []types.JavaByte (from IvParameterSpec) or []byte.
	// Ensure it's returned as a Java byte array object.
	switch v := ivField.Fvalue.(type) {
	case []types.JavaByte:
		return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, v)
	case []byte:
		jBytes := object.JavaByteArrayFromGoByteArray(v)
		return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
	default:
		// Fallback for cases where it might already be an object (unlikely here but safe)
		return object.MakeArrayFromRawArray(v)
	}
}

func cipherGetMaxAllowedParameterSpec(params []any) any {
	// For now, return null, which in many cases means unlimited or no specific restrictions
	return ghelpers.ReturnNull(params)
}

func cipherGetOutputSize(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return int64(0)
	}

	inputLen := params[1].(int64)
	buffer, _ := self.FieldTable["buffer"].Fvalue.([]byte)
	totalLen := inputLen + int64(len(buffer))

	blockSize := cipherGetBlockSize(params).(int64)
	if blockSize == 0 {
		return totalLen
	}

	// Simplified: assume padding might add up to one block
	return totalLen + blockSize
}

func cipherGetParameters(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	// Check if we already have an AlgorithmParameters object
	if paramsField, ok := self.FieldTable["parameters"]; ok && !object.IsNull(paramsField.Fvalue) {
		return paramsField.Fvalue
	}

	// If not, but we have an IV (and it's a suitable algorithm), we can create one
	ivField, ok := self.FieldTable["iv"]
	if !ok || object.IsNull(ivField.Fvalue) {
		return ghelpers.ReturnNull(params)
	}

	config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	// Create AlgorithmParameters for the algorithm
	paramsObj := javaSecurity.AlgparamsGetInstance([]any{nil, object.StringObjectFromGoString(config.Algorithm)})
	if errBlk, ok := paramsObj.(*ghelpers.GErrBlk); ok {
		// If we can't get an instance for this algorithm, just return null
		_ = errBlk
		return ghelpers.ReturnNull(params)
	}

	pObj := paramsObj.(*object.Object)

	// Initialize it with the IV
	var ivObj *object.Object
	switch v := ivField.Fvalue.(type) {
	case *object.Object:
		ivObj = v
	case []types.JavaByte:
		ivObj = object.MakePrimitiveObject(types.ByteArray, types.ByteArray, v)
	case []byte:
		ivObj = object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(v))
	}

	if ivObj != nil {
		javaSecurity.AlgparamsInit([]any{pObj, ivObj})
		self.FieldTable["parameters"] = object.Field{Ftype: types.ClassNameAlgorithmParameters, Fvalue: pObj}
		return pObj
	}

	return ghelpers.ReturnNull(params)
}

func cipherGetProvider(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	providerField, ok := self.FieldTable["provider"]
	if !ok {
		return ghelpers.ReturnNull(params)
	}

	return providerField.Fvalue
}

func cipherInit(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "cipherInit: 'this' cannot be null")
	}

	opmode := params[1].(int64)
	self.FieldTable["opmode"] = object.Field{Ftype: types.Int, Fvalue: opmode}

	// Key is params[2]
	if len(params) > 2 {
		key := params[2].(*object.Object)
		config, hasConfig := self.FieldTable["config"].Fvalue.(CipherTransformation)
		isLegacyPBE := hasConfig && strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC")
		isHmacPBE := hasConfig && strings.Contains(config.Name, "PBEWith") && (strings.Contains(config.Name, "Hmac") || strings.Contains(config.Name, "HMAC"))

		if isLegacyPBE || isHmacPBE {
			if _, ok := key.FieldTable["key"]; !ok {
				// Key is NOT derived. Check if we have an old key that WAS derived for the same password.
				if oldKeyField, ok := self.FieldTable["key"]; ok {
					if oldKeyObj, ok := oldKeyField.Fvalue.(*object.Object); ok {
						if _, ok := oldKeyObj.FieldTable["key"]; ok {
							if oldValField, ok := self.FieldTable["pbe_password"]; ok {
								oldVal := oldValField.Fvalue
								if newValField, ok := key.FieldTable["value"]; ok {
									newVal := newValField.Fvalue
									var oldGoBytes []byte
									switch v := oldVal.(type) {
									case []types.JavaByte:
										oldGoBytes = object.GoByteArrayFromJavaByteArray(v)
									case []byte:
										oldGoBytes = v
									}
									var newGoBytes []byte
									switch v := newVal.(type) {
									case []types.JavaByte:
										newGoBytes = object.GoByteArrayFromJavaByteArray(v)
									case []byte:
										newGoBytes = v
									}
									if oldGoBytes != nil && newGoBytes != nil && bytes.Equal(oldGoBytes, newGoBytes) {
										// Match! Re-use the already derived key.
										key = oldKeyObj
									}
								}
							}
						}
					}
				}
			}
		}
		self.FieldTable["key"] = object.Field{Ftype: "java/security/Key", Fvalue: key}
	}

	ivProvided := false
	// Handle parameters if provided via AlgorithmParameters or AlgorithmParameterSpec
	if len(params) > 3 {
		param := params[3]
		if spec, ok := param.(*object.Object); ok && !object.IsNull(spec) {
			if spec.KlassName == object.StringPoolIndexFromGoString(types.ClassNameAlgorithmParameters) {
				// java.security.AlgorithmParameters
				self.FieldTable["parameters"] = object.Field{Ftype: types.ClassNameAlgorithmParameters, Fvalue: spec}
				// Extract spec from AlgorithmParameters
				innerSpec := javaSecurity.AlgparamsGetParameterSpec([]any{spec})
				if specObj, ok := innerSpec.(*object.Object); ok && !object.IsNull(specObj) {
					// Recurse or handle directly
					handleInitSpec(self, specObj, &ivProvided)
				}
			} else {
				// java.security.spec.AlgorithmParameterSpec
				handleInitSpec(self, spec, &ivProvided)
			}
		} else if param == nil || object.IsNull(param) {
			config, hasConfig := self.FieldTable["config"].Fvalue.(CipherTransformation)
			isLegacyPBE := hasConfig && strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC")
			isHmacPBE := hasConfig && strings.Contains(config.Name, "PBEWith") && (strings.Contains(config.Name, "Hmac") || strings.Contains(config.Name, "HMAC"))
			if isLegacyPBE || isHmacPBE {
				// For legacy PBE, a null spec might mean "use defaults/extract from key"
				// So we DON'T set ivProvided = true yet.
				self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
				ivProvided = false
			} else {
				// If spec is explicitly provided as null, we treat it as "provided" but empty.
				// This prevents auto-generation.
				self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
				ivProvided = true
			}
		}
	}

	// After handling spec (which might have derived the key), update 'old_derived_key' if we now have a derived key.
	if keyField, ok := self.FieldTable["key"]; ok {
		keyObj := keyField.Fvalue.(*object.Object)
		if _, ok := keyObj.FieldTable["key"]; ok {
			self.FieldTable["old_derived_key"] = keyField
		}
	}

	if len(params) <= 3 {
		// If no spec parameter provided at all, it's NOT provided, so it's a candidate for auto-generation.
		// However, we should clear any existing IV from a previous init, UNLESS it's a legacy PBE
		// where we might want to keep the derived IV.
		config, hasConfig := self.FieldTable["config"].Fvalue.(CipherTransformation)
		isLegacyPBE := hasConfig && strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC")
		isHmacPBE := hasConfig && strings.Contains(config.Name, "PBEWith") && (strings.Contains(config.Name, "Hmac") || strings.Contains(config.Name, "HMAC"))

		if isLegacyPBE || isHmacPBE {
			// Try to handle PBE key derivation if no spec is provided but key might need it.
			keyField, ok := self.FieldTable["key"]
			if ok {
				keyObj := keyField.Fvalue.(*object.Object)
				if _, ok := keyObj.FieldTable["key"]; !ok {
					// Key is NOT derived. We already handled potential reuse above,
					// so if it's still not derived, it will fail in cipherDoFinal.
				}
			}
		}

		ivObj, ok := self.FieldTable["iv"]
		if !ok || object.IsNull(ivObj.Fvalue) {
			self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
			ivProvided = false
		} else if isLegacyPBE {
			// Keep existing derived IV for legacy PBE
			ivProvided = true
		} else {
			// For non-legacy PBE, we typically want a fresh IV if none provided
			self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
			ivProvided = false
		}
	}

	// Automatically generate IV if needed and not provided
	if !ivProvided {
		config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
		if ok && config.NeedsIV {
			ivLen := config.IVLength
			if ivLen == 0 {
				// Default block size for AES/DES if IVLength not specified in table
				ivLen = int(cipherGetBlockSize(params).(int64))
			}

			if ivLen > 0 {
				// Special case for legacy PBE: check if key already contains derived IV
				if strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC") {
					keyField, ok := self.FieldTable["key"]
					if ok {
						keyObj := keyField.Fvalue.(*object.Object)
						if kb, ok := keyObj.FieldTable["key"]; ok {
							keyBytes := kb.Fvalue.([]byte)
							// Legacy PBE: Extract Key AND IV if not provided
							var derivedIV []byte
							if strings.Contains(config.Name, "DESede") || strings.Contains(config.Name, "TripleDES") {
								if len(keyBytes) >= 32 {
									derivedIV = keyBytes[24:32]
								}
							} else if strings.Contains(config.Name, "DES") {
								if len(keyBytes) >= 16 {
									derivedIV = keyBytes[8:16]
								}
							} else if strings.Contains(config.Name, "RC2") {
								bits := 128
								if strings.Contains(config.Name, "40") {
									bits = 40
								}
								if len(keyBytes) >= (bits/8)+8 {
									derivedIV = keyBytes[bits/8 : (bits/8)+8]
								}
							}

							if len(derivedIV) > 0 {
								jBytes := object.JavaByteArrayFromGoByteArray(derivedIV)
								ivObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
								self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: ivObj}
								ivProvided = true
							}
						}
					}
				}

				if !ivProvided {
					if opmode == 1 { // ENCRYPT_MODE is 1
						// Generate random IV only for ENCRYPT_MODE
						iv := make([]byte, ivLen)
						if _, err := rand.Read(iv); err == nil {
							jBytes := object.JavaByteArrayFromGoByteArray(iv)
							ivObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
							self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: ivObj}
						}
					} else {
						// For decryption and other modes, if IV is still not provided and it's not a legacy PBE,
						// it might be an error or it might default to zero IV.
						// Legacy Java PBE algorithms (when not deriving IV) often default to a zero IV.
						if strings.Contains(config.Name, "PBEWith") {
							iv := make([]byte, ivLen)
							jBytes := object.JavaByteArrayFromGoByteArray(iv)
							ivObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
							self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: ivObj}
						}
					}
				}
			}
		}
	}

	self.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: []byte{}}

	return nil
}

func handleInitSpec(self *object.Object, spec *object.Object, ivProvided *bool) {
	if spec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/IvParameterSpec") {
		ivField, ok := spec.FieldTable["iv"]
		if ok {
			self.FieldTable["iv"] = ivField
			*ivProvided = true
		}
	} else if spec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/GCMParameterSpec") {
		ivField, ok := spec.FieldTable["iv"]
		if ok {
			self.FieldTable["iv"] = ivField
			*ivProvided = true
		}
		// GCMParameterSpec also has tLen (tag length in bits)
		if tLenField, ok := spec.FieldTable["tLen"]; ok {
			self.FieldTable["tLen"] = tLenField
		}
	} else if spec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/PBEParameterSpec") {
		saltField, ok := spec.FieldTable["salt"]
		if ok {
			self.FieldTable["iv"] = saltField // Default: use salt as IV. May be overwritten by derived IV later.
			*ivProvided = true
		}

		// Perform PBE key derivation if the key is not already derived
		keyField, ok := self.FieldTable["key"]
		if ok {
			keyObj := keyField.Fvalue.(*object.Object)
			config := self.FieldTable["config"].Fvalue.(CipherTransformation)
			if config.KeyDerivation {
				// We need to re-derive if the current key looks like a password (un-derived)
				// Or more simply, if we have a PBEParameterSpec, we should use it to derive.
				// For PBKDF2-based PBE, we need salt and iterations.
				iterations := spec.FieldTable["iterationCount"].Fvalue.(int64)

				// Re-derive the key only if it's not already derived.
				// A derived key will have a "key" field in its FieldTable.
				if _, ok := keyObj.FieldTable["key"]; !ok {
					// Create a PBEKeySpec from the existing password and the new salt/iterations
					passwordBytes := keyObj.FieldTable["value"].Fvalue.([]byte)
					passwordChars := make([]int64, len(passwordBytes))
					for i, b := range passwordBytes {
						passwordChars[i] = int64(b)
					}
					passwordCharsObj := object.MakePrimitiveObject(types.CharArray, types.CharArray, passwordChars)

					pbeKeySpecClassName := "javax/crypto/spec/PBEKeySpec"
					pbeKeySpec := object.MakeEmptyObjectWithClassName(&pbeKeySpecClassName)
					pbeKeySpec.FieldTable["password"] = object.Field{Ftype: "[C", Fvalue: passwordCharsObj}

					var jSalt []types.JavaByte
					switch v := saltField.Fvalue.(type) {
					case []types.JavaByte:
						jSalt = v
					case []byte:
						jSalt = object.JavaByteArrayFromGoByteArray(v)
					}
					pbeKeySpec.FieldTable["salt"] = object.Field{Ftype: "[B", Fvalue: object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jSalt)}
					pbeKeySpec.FieldTable["iterationCount"] = object.Field{Ftype: types.Int, Fvalue: iterations}

					keyLength := int64(0)
					if f, ok := keyObj.FieldTable["inferred_key_length"]; ok {
						keyLength = f.Fvalue.(int64)
					}
					pbeKeySpec.FieldTable["keyLength"] = object.Field{Ftype: types.Int, Fvalue: keyLength}

					// Reuse secretKeyFactoryGenerateSecret to derive the key
					// We need a SKF object for this.
					skfAlgoObj := self.FieldTable["transformation"].Fvalue.(*object.Object)
					skf := secretKeyFactoryGetInstance([]any{skfAlgoObj}).(*object.Object)

					derivedKey := secretKeyFactoryGenerateSecret([]any{skf, pbeKeySpec})
					if derivedKeyObj, ok := derivedKey.(*object.Object); ok {
						self.FieldTable["key"] = object.Field{Ftype: "java/security/Key", Fvalue: derivedKeyObj}
						// Store the password in the cipher for potential reuse later (in case of another init with underived key)
						if pwVal, ok := keyObj.FieldTable["value"]; ok {
							self.FieldTable["pbe_password"] = pwVal
						}

						// If this is a legacy PBE algorithm that should have a derived IV,
						// we signal that the IV is NOT provided so that cipherInit can extract it from the newly derived key.
						// We also clear any existing salt-based IV to ensure extraction happens.
						if strings.Contains(config.Name, "PBEWith") && !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC") {
							self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
							*ivProvided = false
						} else if strings.Contains(config.Name, "PBEWith") {
							// For HMAC-based PBE, we use salt as IV by default (if 16 bytes).
							// But if we already set *ivProvided = true above with saltField, it's fine.
							// The issue is that for legacy PBE we want to extract IV from key material.
						}
					}
				} else {
				}
			}
		}
	} else {
		// Other spec types: we don't handle them as IV providers for now,
		// but they are "provided" so we don't auto-generate.
		*ivProvided = true
	}
}

func cipherUpdate(params []any) any {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return ghelpers.ReturnNull(params)
	}

	var input []byte
	inputObj, ok := params[1].(*object.Object)
	if ok && !object.IsNull(inputObj) {
		input = object.GoByteArrayFromJavaByteArray(inputObj.FieldTable["value"].Fvalue.([]types.JavaByte))
		if len(params) > 3 {
			offset := params[2].(int64)
			length := params[3].(int64)
			input = input[offset : offset+length]
		}
	}

	buffer, _ := self.FieldTable["buffer"].Fvalue.([]byte)
	self.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: append(buffer, input...)}

	// Return null or empty byte array for now, as we buffer everything
	jBytes := object.JavaByteArrayFromGoByteArray([]byte{})
	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
}
