/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"fmt"
	"hash"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

// Load_Security_MessageDigest initializes java/security/MessageDigest methods
func Load_Security_MessageDigest() {

	ghelpers.MethodSignatures["java/security/MessageDigest.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapProtected,
		}

	// ---------- Member Functions ----------
	ghelpers.MethodSignatures["java/security/MessageDigest.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigClone}

	ghelpers.MethodSignatures["java/security/MessageDigest.digest()[B"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigDigest}

	ghelpers.MethodSignatures["java/security/MessageDigest.digest([B)[B"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: msgdigDigestBytes}

	ghelpers.MethodSignatures["java/security/MessageDigest.digest([BII)I"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: msgdigDigestBytesII}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineDigest()[B"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapProtected}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineDigest([BII)I"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapProtected}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineGetDigestLength()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapProtected}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineReset()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapProtected}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineUpdate(B)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapProtected}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineUpdate([BII)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapProtected}

	ghelpers.MethodSignatures["java/security/MessageDigest.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigGetAlgorithm}

	ghelpers.MethodSignatures["java/security/MessageDigest.getDigestLength()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigGetDigestLength}

	ghelpers.MethodSignatures["java/security/MessageDigest.getInstance(Ljava/lang/String;)Ljava/security/MessageDigest;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: msgdigGetInstance}

	ghelpers.MethodSignatures["java/security/MessageDigest.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/MessageDigest;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: msgdigGetInstanceProvider}

	ghelpers.MethodSignatures["java/security/MessageDigest.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljava/security/MessageDigest;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: msgdigGetInstanceProviderObj}

	ghelpers.MethodSignatures["java/security/MessageDigest.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigGetProvider}

	ghelpers.MethodSignatures["java/security/MessageDigest.isEqual([B[B)Z"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: msgdigIsEqual}

	ghelpers.MethodSignatures["java/security/MessageDigest.reset()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigReset}

	ghelpers.MethodSignatures["java/security/MessageDigest.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: msgdigToString}

	ghelpers.MethodSignatures["java/security/MessageDigest.update(B)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: msgdigUpdateByte}

	ghelpers.MethodSignatures["java/security/MessageDigest.update([B)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: msgdigUpdateBytes}

	ghelpers.MethodSignatures["java/security/MessageDigest.update([BII)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: msgdigUpdateBytesII}

	ghelpers.MethodSignatures["java/security/MessageDigest.update(Ljava/nio/ByteBuffer;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
}

// ===================== Helper Functions =====================

func getHashForAlgorithm(algorithm string) (hash.Hash, error) {
	alg := strings.ToUpper(algorithm)
	switch alg {
	case "MD5":
		return md5.New(), nil
	case "SHA-1", "SHA1":
		return sha1.New(), nil
	case "SHA-224", "SHA224":
		return sha256.New224(), nil
	case "SHA-256", "SHA256":
		return sha256.New(), nil
	case "SHA-384", "SHA384":
		return sha512.New384(), nil
	case "SHA-512", "SHA512":
		return sha512.New(), nil
	case "SHA-512/224", "SHA512/224":
		return sha512.New512_224(), nil
	case "SHA-512/256", "SHA512/256":
		return sha512.New512_256(), nil
	case "SHA3-224", "SHA3_224":
		return sha3.New224(), nil
	case "SHA3-256", "SHA3_256":
		return sha3.New256(), nil
	case "SHA3-384", "SHA3_384":
		return sha3.New384(), nil
	case "SHA3-512", "SHA3_512":
		return sha3.New512(), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func getDigestLengthForAlgorithm(algorithm string) int {
	alg := strings.ToUpper(algorithm)
	switch alg {
	case "MD5":
		return 16
	case "SHA-1", "SHA1":
		return 20
	case "SHA-224", "SHA224", "SHA-512/224", "SHA512/224":
		return 28
	case "SHA-256", "SHA256", "SHA-512/256", "SHA512/256":
		return 32
	case "SHA-384", "SHA384":
		return 48
	case "SHA-512", "SHA512":
		return 64
	default:
		return 0
	}
}

func appendToBuffer(this *object.Object, bytes []types.JavaByte) {
	current := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	current = append(current, bytes...)
	this.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: current,
	}
}

func extractJavaBytes(bytesObj *object.Object) ([]types.JavaByte, error) {
	field, ok := bytesObj.FieldTable["value"]
	if !ok {
		return nil, fmt.Errorf("missing 'value' field")
	}

	switch v := field.Fvalue.(type) {
	case []types.JavaByte:
		return v, nil
	case []byte:
		return object.JavaByteArrayFromGoByteArray(v), nil
	default:
		return nil, fmt.Errorf("invalid byte array type")
	}
}

// ===================== MessageDigest Methods =====================

func msgdigGetInstance(params []any) any {
	algorithmObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "msgdigGetInstance: Algorithm cannot be null")
	}
	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Get the default (only) security provider.
	providerObj := ghelpers.GetDefaultSecurityProvider() // single Go runtime provider

	// Try to get a service from the provider
	svcObj := securityProviderGetService([]interface{}{providerObj, object.StringObjectFromGoString("MessageDigest"), algorithmObj})
	if errBlk, ok := svcObj.(*ghelpers.GErrBlk); ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "msgdigGetInstance: "+errBlk.ErrMsg)
	}

	// Create MessageDigest object
	md := object.MakeEmptyObjectWithClassName(&types.ClassNameMessageDigest)

	// Store algorithm name
	md.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algorithm),
	}

	// Store provider
	md.FieldTable["provider"] = object.Field{
		Ftype:  types.ClassNameSecurityProvider,
		Fvalue: providerObj,
	}

	// Optionally store reference to the service (can be useful for future extensions)
	md.FieldTable["service"] = object.Field{
		Ftype:  "java/security/Provider$Service",
		Fvalue: svcObj.(*object.Object),
	}

	// Initialize empty buffer for accumulating data
	md.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	return md
}

func msgdigGetInstanceProvider(params []any) any {
	const fname = "msgdigGetInstanceProvider"
	algorithmObj := params[0].(*object.Object)
	providerObj := params[1].(*object.Object)
	providerName := object.GoStringFromStringObject(providerObj)

	if providerName != types.SecurityProviderName {
		return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException,
			fmt.Sprintf("%s: Provider %s not supported. Only %s is supported.", fname, providerName, types.SecurityProviderName))
	}

	return msgdigGetInstance([]any{algorithmObj})
}

func msgdigGetInstanceProviderObj(params []any) any {
	const fname = "msgdigGetInstanceProviderObj"
	algorithmObj := params[0].(*object.Object)
	providerObj := params[1].(*object.Object)

	if providerObj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("%s: Provider cannot be null", fname))
	}

	providerNameField, ok := providerObj.FieldTable["name"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("%s: Invalid provider object", fname))
	}

	providerNameObj := providerNameField.Fvalue.(*object.Object)
	providerName := object.GoStringFromStringObject(providerNameObj)

	if providerName != types.SecurityProviderName {
		return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException,
			fmt.Sprintf("%s: Provider %s not supported. Only %s is supported.", fname, providerName, types.SecurityProviderName))
	}

	return msgdigGetInstance([]any{algorithmObj})
}

func msgdigGetAlgorithm(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["algorithm"].Fvalue.(*object.Object)
}

func msgdigGetProvider(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["provider"].Fvalue.(*object.Object)
}

func msgdigGetDigestLength(params []any) any {
	this := params[0].(*object.Object)
	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)
	return int64(getDigestLengthForAlgorithm(algorithm))
}

func msgdigUpdateByte(params []any) any {
	const fname = "msgdigUpdateByte"
	this := params[0].(*object.Object)
	b := params[1].(int64)
	appendToBuffer(this, []types.JavaByte{types.JavaByte(b)})
	return nil
}

func msgdigUpdateBytes(params []any) any {
	const fname = "msgdigUpdateBytes"
	this := params[0].(*object.Object)
	bytesObj := params[1].(*object.Object)

	bytes, err := extractJavaBytes(bytesObj)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("%s: %s", fname, err.Error()))
	}

	appendToBuffer(this, bytes)
	return nil
}

func msgdigUpdateBytesII(params []any) any {
	const fname = "msgdigUpdateBytesII"
	this := params[0].(*object.Object)
	bytesObj := params[1].(*object.Object)
	offset := params[2].(int64)
	length := params[3].(int64)

	bytes, err := extractJavaBytes(bytesObj)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("%s: %s", fname, err.Error()))
	}

	if offset < 0 || length < 0 || offset+length > int64(len(bytes)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("%s: Invalid offset or length", fname))
	}

	appendToBuffer(this, bytes[offset:offset+length])
	return nil
}

func msgdigDigest(params []any) any {
	const fname = "msgdigDigest"
	this := params[0].(*object.Object)
	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	h, err := getHashForAlgorithm(algorithm)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("%s: %s", fname, err.Error()))
	}

	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	h.Write(object.GoByteArrayFromJavaByteArray(buffer))
	digest := h.Sum(nil)
	javaDigest := object.JavaByteArrayFromGoByteArray(digest)

	this.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{}}

	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, javaDigest)
}

func msgdigDigestBytes(params []any) any {
	const fname = "msgdigDigestBytes"
	this := params[0].(*object.Object)
	bytesObj := params[1].(*object.Object)

	if res := msgdigUpdateBytes([]any{this, bytesObj}); res != nil {
		if _, ok := res.(*ghelpers.GErrBlk); ok {
			return res
		}
	}

	return msgdigDigest([]any{this})
}

func msgdigDigestBytesII(params []any) any {
	const fname = "msgdigDigestBytesII"
	this := params[0].(*object.Object)
	bufObj := params[1].(*object.Object)
	offset := params[2].(int64)
	length := params[3].(int64)

	bufferField := bufObj.FieldTable["value"]
	var buf []types.JavaByte
	switch v := bufferField.Fvalue.(type) {
	case []types.JavaByte:
		buf = v
	case []byte:
		buf = object.JavaByteArrayFromGoByteArray(v)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("%s: Invalid byte array", fname))
	}

	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)
	digestLen := getDigestLengthForAlgorithm(algorithm)

	if offset < 0 || length < int64(digestLen) || offset+int64(digestLen) > int64(len(buf)) {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, fmt.Sprintf("%s: Buffer too small or invalid offset", fname))
	}

	digestObj := msgdigDigest([]any{this})
	if errBlk, ok := digestObj.(*ghelpers.GErrBlk); ok {
		return errBlk
	}

	digestBytes := digestObj.(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	copy(buf[offset:], digestBytes)
	bufObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: buf}

	return int64(len(digestBytes))
}

func msgdigReset(params []any) any {
	const fname = "msgdigReset"
	this := params[0].(*object.Object)
	this.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{}}
	return nil
}

func msgdigIsEqual(params []any) any {
	const fname = "msgdigIsEqual"
	a := params[0].(*object.Object)
	b := params[1].(*object.Object)

	bytesA := a.FieldTable["value"].Fvalue.([]types.JavaByte)
	bytesB := b.FieldTable["value"].Fvalue.([]types.JavaByte)

	if object.JavaByteArrayEquals(bytesA, bytesB) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func msgdigToString(params []any) any {
	this := params[0].(*object.Object)
	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	return object.StringObjectFromGoString(fmt.Sprintf("MessageDigest[%s]", algorithm))
}

func msgdigClone(params []any) any {
	this := params[0].(*object.Object)
	clone := object.MakeEmptyObjectWithClassName(&types.ClassNameMessageDigest)

	clone.FieldTable["algorithm"] = this.FieldTable["algorithm"]
	clone.FieldTable["provider"] = this.FieldTable["provider"]

	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	bufferCopy := make([]types.JavaByte, len(buffer))
	copy(bufferCopy, buffer)
	clone.FieldTable["buffer"] = object.Field{Ftype: types.ByteArray, Fvalue: bufferCopy}

	return clone
}
