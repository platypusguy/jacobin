/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"crypto/rand"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
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
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getOutputSize(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  cipherGetOutputSize,
		}

	ghelpers.MethodSignatures["javax/crypto/Cipher.getParameters()Ljava/security/AlgorithmParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
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
		keyBytes, ok = kb.Fvalue.([]byte)
	} else if vb, ok := keyObj.FieldTable["value"]; ok {
		// Handle cases where the key is just a byte array object (common in tests)
		if jBytes, ok := vb.Fvalue.([]types.JavaByte); ok {
			keyBytes = object.GoByteArrayFromJavaByteArray(jBytes)
		} else if bBytes, ok := vb.Fvalue.([]byte); ok {
			keyBytes = bBytes
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
		self.FieldTable["key"] = object.Field{Ftype: "java/security/Key", Fvalue: key}
	}

	ivProvided := false
	// Handle IV if provided via AlgorithmParameterSpec (IvParameterSpec)
	if len(params) > 3 {
		spec := params[3].(*object.Object)
		if !object.IsNull(spec) {
			// Check if it's IvParameterSpec
			if spec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/IvParameterSpec") {
				ivField, ok := spec.FieldTable["iv"]
				if ok {
					self.FieldTable["iv"] = ivField
					ivProvided = true
				}
			} else if spec.KlassName == object.StringPoolIndexFromGoString("javax/crypto/spec/GCMParameterSpec") {
				ivField, ok := spec.FieldTable["iv"]
				if ok {
					self.FieldTable["iv"] = ivField
					ivProvided = true
				}
				// GCMParameterSpec also has tLen (tag length in bits)
				if tLenField, ok := spec.FieldTable["tLen"]; ok {
					self.FieldTable["tLen"] = tLenField
				}
			} else {
				// Other spec types: we don't handle them as IV providers for now,
				// but they are "provided" so we don't auto-generate.
				ivProvided = true
			}
		} else {
			// If spec is explicitly provided as null, we treat it as "provided" but empty.
			// This prevents auto-generation.
			self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
			ivProvided = true
		}
	} else {
		// If no spec parameter provided at all, it's NOT provided, so it's a candidate for auto-generation.
		// However, we should clear any existing IV from a previous init.
		self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: object.Null}
	}

	// Automatically generate IV if needed and not provided (only for ENCRYPT_MODE)
	if !ivProvided && opmode == 1 { // ENCRYPT_MODE is 1
		config, ok := self.FieldTable["config"].Fvalue.(CipherTransformation)
		if ok && config.NeedsIV {
			ivLen := config.IVLength
			if ivLen == 0 {
				// Default block size for AES/DES if IVLength not specified in table
				ivLen = int(cipherGetBlockSize(params).(int64))
			}

			if ivLen > 0 {
				iv := make([]byte, ivLen)
				if _, err := rand.Read(iv); err == nil {
					jBytes := object.JavaByteArrayFromGoByteArray(iv)
					ivObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, jBytes)
					self.FieldTable["iv"] = object.Field{Ftype: types.ByteArray, Fvalue: ivObj}
				}
			}
		}
	}

	self.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: []byte{}}

	return nil
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
