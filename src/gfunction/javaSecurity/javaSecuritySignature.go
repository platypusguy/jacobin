/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
	"strings"
)

func Load_Security_Signature() {

	// <clinit>
	ghelpers.MethodSignatures["java/security/Signature.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// Constructors (protected)
	ghelpers.MethodSignatures["java/security/Signature.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// Static methods
	ghelpers.MethodSignatures["java/security/Signature.getInstance(Ljava/lang/String;)Ljava/security/Signature;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signatureGetInstance,
		}

	ghelpers.MethodSignatures["java/security/Signature.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/Signature;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/security/Signature.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljava/security/Signature;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Instance methods
	ghelpers.MethodSignatures["java/security/Signature.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  signatureGetAlgorithm,
		}

	ghelpers.MethodSignatures["java/security/Signature.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  SecurityGetProvider,
		}

	ghelpers.MethodSignatures["java/security/Signature.initSign(Ljava/security/PrivateKey;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signatureInitSign,
		}

	ghelpers.MethodSignatures["java/security/Signature.initSign(Ljava/security/PrivateKey;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  signatureInitSignRandom,
		}

	ghelpers.MethodSignatures["java/security/Signature.initVerify(Ljava/security/PublicKey;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signatureInitVerify,
		}

	ghelpers.MethodSignatures["java/security/Signature.sign()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  signatureSign,
		}

	ghelpers.MethodSignatures["java/security/Signature.sign([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  signatureSignToBuffer,
		}

	ghelpers.MethodSignatures["java/security/Signature.update(B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signatureUpdateByte,
		}

	ghelpers.MethodSignatures["java/security/Signature.update([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signatureUpdateBytes,
		}

	ghelpers.MethodSignatures["java/security/Signature.update([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  signatureUpdateBytesRange,
		}

	ghelpers.MethodSignatures["java/security/Signature.verify([B)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signatureVerify,
		}

	ghelpers.MethodSignatures["java/security/Signature.verify([BII)Z"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  signatureVerifyRange,
		}
}

// signatureGetInstance creates a new Signature object for the given algorithm
func signatureGetInstance(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			fmt.Sprintf("Signature.getInstance: expected 1 parameter, got %d", len(params)),
		)
	}

	algorithm, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.getInstance: algorithm is not a String object",
		)
	}

	algoStr := object.GoStringFromStringObject(algorithm)

	// Validate algorithm
	if !isSupportedSignatureAlgorithm(algoStr) {
		return ghelpers.GetGErrBlk(
			excNames.NoSuchAlgorithmException,
			fmt.Sprintf("unsupported signature algorithm: %s", algoStr),
		)
	}

	// Create Signature object
	sigObj := object.MakeEmptyObjectWithClassName(&types.ClassNameSignature)

	// Store algorithm
	sigObj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algoStr),
	}

	// Initialize buffer for update() calls
	sigObj.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	// State: 0=uninitialized, 1=sign, 2=verify
	sigObj.FieldTable["state"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(0),
	}

	return sigObj
}

// signatureGetAlgorithm returns the algorithm name
func signatureGetAlgorithm(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.getAlgorithm: expected 0 parameters",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.getAlgorithm: this is not an Object",
		)
	}

	return sigObj.FieldTable["algorithm"].Fvalue
}

// signatureInitSign initializes the signature for signing
func signatureInitSign(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initSign: expected 1 parameter",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initSign: this is not an Object",
		)
	}

	privateKeyObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initSign: privateKey is not an Object",
		)
	}

	// Get algorithm
	algoStr := object.GoStringFromStringObject(sigObj.FieldTable["algorithm"].Fvalue.(*object.Object))

	// Validate key type matches algorithm
	requiredKeyType := getRequiredPrivateKeyType(algoStr)
	actualKeyType := object.GoStringFromStringPoolIndex(privateKeyObj.KlassName)

	if requiredKeyType != actualKeyType {
		return ghelpers.GetGErrBlk(
			excNames.InvalidKeyException,
			fmt.Sprintf("%s requires %s, got %s", algoStr, requiredKeyType, actualKeyType),
		)
	}

	// Store private key
	sigObj.FieldTable["privateKey"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: privateKeyObj,
	}

	// Set state to signing
	sigObj.FieldTable["state"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(1),
	}

	// Reset buffer
	sigObj.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	return nil
}

// signatureInitSignRandom initializes the signature for signing with custom random
func signatureInitSignRandom(params []any) any {
	if len(params) != 3 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initSign: expected 2 parameters",
		)
	}

	// For now, ignore the SecureRandom parameter and call regular initSign
	return signatureInitSign([]any{params[0], params[1]})
}

// signatureInitVerify initializes the signature for verification
func signatureInitVerify(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initVerify: expected 1 parameter",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initVerify: this is not an Object",
		)
	}

	publicKeyObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.initVerify: publicKey is not an Object",
		)
	}

	// Get algorithm
	algoStr := object.GoStringFromStringObject(sigObj.FieldTable["algorithm"].Fvalue.(*object.Object))

	// Validate key type matches algorithm
	requiredKeyType := getRequiredPublicKeyType(algoStr)
	actualKeyType := object.GoStringFromStringPoolIndex(publicKeyObj.KlassName)

	if requiredKeyType != actualKeyType {
		return ghelpers.GetGErrBlk(
			excNames.InvalidKeyException,
			fmt.Sprintf("%s requires %s, got %s", algoStr, requiredKeyType, actualKeyType),
		)
	}

	// Store public key
	sigObj.FieldTable["publicKey"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: publicKeyObj,
	}

	// Set state to verifying
	sigObj.FieldTable["state"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(2),
	}

	// Reset buffer
	sigObj.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	return nil
}

// signatureUpdateByte updates the signature with a single byte
func signatureUpdateByte(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: expected 1 parameter",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: this is not an Object",
		)
	}

	b, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: byte parameter is not a byte",
		)
	}

	// Get current buffer
	bufferJbytes := sigObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	buffer := object.GoByteArrayFromJavaByteArray(bufferJbytes)

	// Append byte
	buffer = append(buffer, byte(b))

	// Store updated buffer
	sigObj.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: object.JavaByteArrayFromGoByteArray(buffer),
	}

	return nil
}

// signatureUpdateBytes updates the signature with a byte array
func signatureUpdateBytes(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: expected 1 parameter",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: this is not an Object",
		)
	}

	dataObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: data is not an Object",
		)
	}

	dataJbytes := dataObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	data := object.GoByteArrayFromJavaByteArray(dataJbytes)

	// Get current buffer
	bufferJbytes := sigObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	buffer := object.GoByteArrayFromJavaByteArray(bufferJbytes)

	// Append data
	buffer = append(buffer, data...)

	// Store updated buffer
	sigObj.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: object.JavaByteArrayFromGoByteArray(buffer),
	}

	return nil
}

// signatureUpdateBytesRange updates the signature with a range of bytes
func signatureUpdateBytesRange(params []any) any {
	if len(params) != 4 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: expected 3 parameters",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: this is not an Object",
		)
	}

	dataObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: data is not an Object",
		)
	}

	offset, ok := params[2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: offset is not an int",
		)
	}

	length, ok := params[3].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.update: length is not an int",
		)
	}

	dataJbytes := dataObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	data := object.GoByteArrayFromJavaByteArray(dataJbytes)

	// Validate bounds
	if offset < 0 || length < 0 || int(offset+length) > len(data) {
		return ghelpers.GetGErrBlk(
			excNames.ArrayIndexOutOfBoundsException,
			"Signature.update: invalid offset or length",
		)
	}

	// Get current buffer
	bufferJbytes := sigObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	buffer := object.GoByteArrayFromJavaByteArray(bufferJbytes)

	// Append data range
	buffer = append(buffer, data[offset:offset+length]...)

	// Store updated buffer
	sigObj.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: object.JavaByteArrayFromGoByteArray(buffer),
	}

	return nil
}

// signatureSign generates the signature
func signatureSign(params []any) any {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.sign: expected 0 parameters",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.sign: this is not an Object",
		)
	}

	// Check state
	state := sigObj.FieldTable["state"].Fvalue.(int64)
	if state != 1 {
		return ghelpers.GetGErrBlk(
			excNames.SignatureException,
			"Signature.sign: signature not initialized for signing",
		)
	}

	// Get algorithm
	algoStr := object.GoStringFromStringObject(sigObj.FieldTable["algorithm"].Fvalue.(*object.Object))

	// Get private key
	privateKeyObj := sigObj.FieldTable["privateKey"].Fvalue.(*object.Object)
	privateKeyValue := privateKeyObj.FieldTable["value"].Fvalue

	// Get data buffer
	dataJBytes := sigObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	data := object.GoByteArrayFromJavaByteArray(dataJBytes)

	// Sign based on algorithm
	var signature []byte
	var err error

	switch {
	case strings.Contains(algoStr, "RSA"):
		signature, err = signRSA(privateKeyValue.(*rsa.PrivateKey), data, algoStr)
	case strings.Contains(algoStr, "ECDSA"):
		signature, err = signECDSA(privateKeyValue.(*ecdsa.PrivateKey), data, algoStr)
	case strings.Contains(algoStr, "DSA"):
		signature, err = signDSA(privateKeyValue.(*dsa.PrivateKey), data, algoStr)
	case algoStr == "Ed25519":
		signature, err = signEd25519(privateKeyValue.(ed25519.PrivateKey), data)
	default:
		return ghelpers.GetGErrBlk(
			excNames.SignatureException,
			fmt.Sprintf("unsupported signature algorithm: %s", algoStr),
		)
	}

	if err != nil {
		return ghelpers.GetGErrBlk(excNames.SignatureException, err.Error())
	}

	// Return signature as byte array object
	sigJbytes := object.JavaByteArrayFromGoByteArray(signature)
	sigArrayObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, sigJbytes)
	return sigArrayObj
}

// signatureSignToBuffer signs and stores result in provided buffer
func signatureSignToBuffer(params []any) any {
	if len(params) != 4 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.sign: expected 3 parameters",
		)
	}

	// Get signature bytes
	sigResult := signatureSign([]any{params[0]})

	// Check for error
	if _, isError := sigResult.(*ghelpers.GErrBlk); isError {
		return sigResult
	}

	sigArrayObj := sigResult.(*object.Object)
	sigJbytes := sigArrayObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	signature := object.GoByteArrayFromJavaByteArray(sigJbytes)

	// Get output buffer
	outbufObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.sign: outbuf is not an Object",
		)
	}

	offset, ok := params[2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.sign: offset is not an int",
		)
	}

	length, ok := params[3].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.sign: length is not an int",
		)
	}

	outbuf := outbufObj.FieldTable["value"].Fvalue.([]byte)

	// Validate bounds
	if offset < 0 || int(offset+int64(len(signature))) > len(outbuf) {
		return ghelpers.GetGErrBlk(
			excNames.SignatureException,
			"Signature.sign: buffer too small",
		)
	}

	// Enforce caller-provided length
	if int64(len(signature)) > length {
		return ghelpers.GetGErrBlk(
			excNames.SignatureException,
			"Signature.sign: buffer too small for signature",
		)
	}

	// Copy output buffer to signature.
	copy(outbuf[offset:], signature)
	sigJbytes = object.JavaByteArrayFromGoByteArray(signature)
	sigArrayObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: sigJbytes}

	// Return length of signature
	return int64(len(signature))
}

// signatureVerify verifies a signature
func signatureVerify(params []any) any {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: expected 1 parameter",
		)
	}

	sigObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: this is not an Object",
		)
	}

	sigArrayObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: signature is not an Object",
		)
	}

	// Check state
	state := sigObj.FieldTable["state"].Fvalue.(int64)
	if state != 2 {
		return ghelpers.GetGErrBlk(
			excNames.SignatureException,
			"Signature.verify: signature not initialized for verification",
		)
	}

	// Get algorithm
	algoStr := object.GoStringFromStringObject(sigObj.FieldTable["algorithm"].Fvalue.(*object.Object))

	// Get public key
	publicKeyObj := sigObj.FieldTable["publicKey"].Fvalue.(*object.Object)
	publicKeyValue := publicKeyObj.FieldTable["value"].Fvalue

	// Get data buffer
	dataJbytes := sigObj.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	data := object.GoByteArrayFromJavaByteArray(dataJbytes)

	// Get signature bytes
	sigJbytes := sigArrayObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	signature := object.GoByteArrayFromJavaByteArray(sigJbytes)

	// Verify based on algorithm
	var valid bool

	switch {
	case strings.Contains(algoStr, "RSA"):
		valid = verifyRSA(publicKeyValue.(*rsa.PublicKey), data, signature, algoStr)
	case strings.Contains(algoStr, "ECDSA"):
		valid = verifyECDSA(publicKeyValue.(*ecdsa.PublicKey), data, signature, algoStr)
	case strings.Contains(algoStr, "DSA"):
		valid = verifyDSA(publicKeyValue.(*dsa.PublicKey), data, signature, algoStr)
	case algoStr == "Ed25519":
		valid = verifyEd25519(publicKeyValue.(ed25519.PublicKey), data, signature)
	default:
		return ghelpers.GetGErrBlk(
			excNames.SignatureException,
			fmt.Sprintf("unsupported signature algorithm: %s", algoStr),
		)
	}

	return valid
}

// signatureVerifyRange verifies a signature from a range of bytes
func signatureVerifyRange(params []any) any {
	if len(params) != 4 {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: expected 3 parameters",
		)
	}

	sigArrayObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: signature is not an Object",
		)
	}

	offset, ok := params[2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: offset is not an int",
		)
	}

	length, ok := params[3].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(
			excNames.IllegalArgumentException,
			"Signature.verify: length is not an int",
		)
	}

	sigJbytes := sigArrayObj.FieldTable["value"].Fvalue.([]types.JavaByte)
	signature := object.GoByteArrayFromJavaByteArray(sigJbytes)

	// Validate bounds
	if offset < 0 || length < 0 || int(offset+length) > len(signature) {
		return ghelpers.GetGErrBlk(
			excNames.ArrayIndexOutOfBoundsException,
			"Signature.verify: invalid offset or length",
		)
	}

	// Extract signature range
	sigRange := signature[offset : offset+length]
	sigJbytes = object.JavaByteArrayFromGoByteArray(sigRange)

	// Create new signature array object
	sigRangeObj := object.MakePrimitiveObject(types.ByteArray, types.ByteArray, sigJbytes)

	// Call regular verify with extracted range
	return signatureVerify([]any{params[0], sigRangeObj})
}

// Helper functions

func isSupportedSignatureAlgorithm(algo string) bool {
	supported := []string{
		"SHA256withRSA", "SHA384withRSA", "SHA512withRSA",
		"SHA256withECDSA", "SHA384withECDSA", "SHA512withECDSA",
		"SHA256withDSA",
		"Ed25519",
	}
	for _, s := range supported {
		if s == algo {
			return true
		}
	}
	return false
}

func getRequiredPrivateKeyType(algo string) string {
	if strings.Contains(algo, "RSA") {
		return types.ClassNameRSAPrivateKey
	} else if strings.Contains(algo, "ECDSA") {
		return types.ClassNameECPrivateKey
	} else if strings.Contains(algo, "DSA") {
		return types.ClassNameDSAPrivateKey
	} else if algo == "Ed25519" {
		return types.ClassNameEdECPrivateKey
	}
	return ""
}

func getRequiredPublicKeyType(algo string) string {
	if strings.Contains(algo, "RSA") {
		return types.ClassNameRSAPublicKey
	} else if strings.Contains(algo, "ECDSA") {
		return types.ClassNameECPublicKey
	} else if strings.Contains(algo, "DSA") {
		return types.ClassNameDSAPublicKey
	} else if algo == "Ed25519" {
		return types.ClassNameEdECPublicKey
	}
	return ""
}

func getSigHashForAlgorithm(algo string) crypto.Hash {
	if strings.Contains(algo, "SHA256") {
		return crypto.SHA256
	} else if strings.Contains(algo, "SHA384") {
		return crypto.SHA384
	} else if strings.Contains(algo, "SHA512") {
		return crypto.SHA512
	}
	return crypto.SHA256 // default
}

func hashData(data []byte, hashType crypto.Hash) []byte {
	switch hashType {
	case crypto.SHA256:
		h := sha256.Sum256(data)
		return h[:]
	case crypto.SHA384:
		h := sha512.Sum384(data)
		return h[:]
	case crypto.SHA512:
		h := sha512.Sum512(data)
		return h[:]
	default:
		h := sha256.Sum256(data)
		return h[:]
	}
}

// RSA signing/verification
func signRSA(privateKey *rsa.PrivateKey, data []byte, algo string) ([]byte, error) {
	hashType := getSigHashForAlgorithm(algo)
	hashed := hashData(data, hashType)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, hashType, hashed)
}

func verifyRSA(publicKey *rsa.PublicKey, data []byte, signature []byte, algo string) bool {
	hashType := getSigHashForAlgorithm(algo)
	hashed := hashData(data, hashType)
	err := rsa.VerifyPKCS1v15(publicKey, hashType, hashed, signature)
	return err == nil
}

// ECDSA signing/verification
func signECDSA(privateKey *ecdsa.PrivateKey, data []byte, algo string) ([]byte, error) {
	hashType := getSigHashForAlgorithm(algo)
	hashed := hashData(data, hashType)
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func verifyECDSA(publicKey *ecdsa.PublicKey, data []byte, signature []byte, algo string) bool {
	hashType := getSigHashForAlgorithm(algo)
	hashed := hashData(data, hashType)
	return ecdsa.VerifyASN1(publicKey, hashed, signature)
}

// DSA signing/verification
func signDSA(privateKey *dsa.PrivateKey, data []byte, algo string) ([]byte, error) {
	hashType := getSigHashForAlgorithm(algo)
	hashed := hashData(data, hashType)
	r, s, err := dsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return nil, err
	}

	// Concatenate r and s
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func verifyDSA(publicKey *dsa.PublicKey, data []byte, signature []byte, algo string) bool {
	hashType := getSigHashForAlgorithm(algo)
	hashed := hashData(data, hashType)

	// Split signature into r and s
	sigLen := len(signature) / 2
	r := new(big.Int).SetBytes(signature[:sigLen])
	s := new(big.Int).SetBytes(signature[sigLen:])

	return dsa.Verify(publicKey, hashed, r, s)
}

// Ed25519 signing/verification
func signEd25519(privateKey ed25519.PrivateKey, data []byte) ([]byte, error) {
	// Ed25519 doesn't need pre-hashing
	signature := ed25519.Sign(privateKey, data)
	return signature, nil
}

func verifyEd25519(publicKey ed25519.PublicKey, data []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}
