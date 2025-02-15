/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"crypto/rand"
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"runtime"
)

func Load_Security_SecureRandom() {

	MethodSignatures["java/security/SecureRandom.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/security/SecureRandom.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.<init>(Ljava/security/SecureRandomSpi;Ljava/security/Provider;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.generateSeed(I)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomGenerateSeed,
		}

	MethodSignatures["java/security/SecureRandom.getAlgorithm()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomGetAlgorithm,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;Ljava/security/Provider;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.getInstanceStrong()Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandom_Init,
		}

	MethodSignatures["java/security/SecureRandom.getParameters()Ljava/security/SecureRandomParameters;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getProvider()Ljava/security/Provider;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getSeed(I)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomNextBytes,
		}

	MethodSignatures["java/security/SecureRandom.next(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomNext,
		}

	MethodSignatures["java/security/SecureRandom.nextBytes([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomNextBytes,
		}

	MethodSignatures["java/security/SecureRandom.nextBoolean()Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomNextBoolean,
		}

	MethodSignatures["java/security/SecureRandom.nextBytes([BLjava/security/SecureRandomParameters;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.nextDouble()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomNextFloat,
		}

	MethodSignatures["java/security/SecureRandom.nextFloat()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomNextFloat,
		}

	MethodSignatures["java/security/SecureRandom.nextGaussian()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.nextInt()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomNextInt,
		}

	MethodSignatures["java/security/SecureRandom.nextLong()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomNextInt,
		}

	MethodSignatures["java/security/SecureRandom.reseed()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/security/SecureRandom.reseed(Ljava/security/SecureRandomParameters;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/security/SecureRandom.setSeed([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/security/SecureRandom.setSeed(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/security/SecureRandom.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomToString,
		}

}

var secureRandomClassName = "java/security/SecureRandom"

// SecureRandom.<init> or SecureRandom.getInstance() with a dummy seed
// "java/security/SecureRandom.<init>()V"
func SecureRandom_Init(params []interface{}) interface{} {

	// Create dummy seed.
	dummySeed := []byte{0}

	// Create SecureRandom object with dummy seed value.
	obj := object.MakeEmptyObjectWithClassName(&secureRandomClassName)
	obj.FieldTable["seed"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: dummySeed,
	}

	return obj

}

// SecureRandomNextBytes generates a specified number of random bytes
func SecureRandomNextBytes(params []interface{}) interface{} {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("SecureRandomNextBytes: Expected 2 parameters (SecureRandom object, int64 size), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	secureRandomObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextBytes: First parameter must be a SecureRandom object")
	}

	if secureRandomObj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextBytes: SecureRandom object cannot be nil")
	}

	size, ok := params[1].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextBytes: Second parameter must be an int64")
	}

	// Generate random bytes
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextBytes: Failed to generate random bytes: %v", err))
	}

	result := object.JavaByteArrayFromGoByteArray(bytes)
	return result
}

// SecureRandomNext generates an integer containing the user-specified number of pseudo-random bits
// (right justified, with leading zeros).
func SecureRandomNext(params []interface{}) interface{} {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("SecureRandomNext: Expected 2 parameters (SecureRandom object, bit count), observed %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	secureRandomObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNext: Parameter must be a SecureRandom object")
	}

	if secureRandomObj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNext: SecureRandom object cannot be nil")
	}

	intArg := params[1].(int64)
	if intArg < 0 || intArg > 32 {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNext: bit count must be >= 0 and <=32 ")
	}

	// Generate random int64
	var result int64
	randBytes := make([]byte, 8) // int64 is 8 bytes
	_, err := rand.Read(randBytes)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNext: Failed to generate random int64: %v", err))
	}

	// Convert bytes to int64
	for i := 0; i < 8; i++ {
		result = (result << 8) | int64(randBytes[i])
	}

	// Mask in only the bits requested.
	mask := 2 ^ intArg
	result &= mask

	return result
}

// SecureRandomNextInt generates a random int64
func SecureRandomNextInt(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("SecureRandomNextInt: Expected 1 parameter (SecureRandom object), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	secureRandomObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextInt: Parameter must be a SecureRandom object")
	}

	if secureRandomObj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextInt: SecureRandom object cannot be nil")
	}

	// Generate random int64
	var result int64
	randBytes := make([]byte, 8) // int64 is 8 bytes
	_, err := rand.Read(randBytes)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextInt: Failed to generate random int64: %v", err))
	}

	// Convert bytes to int64
	for i := 0; i < 8; i++ {
		result = (result << 8) | int64(randBytes[i])
	}

	return result
}

// SecureRandomNextFloat generates a random float64
func SecureRandomNextFloat(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("SecureRandomNextFloat: Expected 1 parameter (SecureRandom object), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	secureRandomObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextFloat: Parameter must be a SecureRandom object")
	}

	if secureRandomObj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextFloat: SecureRandom object cannot be nil")
	}

	// Generate random float64 in the range [0, 1)
	randBytes := make([]byte, 8) // float64 is 8 bytes
	_, err := rand.Read(randBytes)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextFloat: Failed to generate random float64: %v", err))
	}

	// Convert bytes to a value in [0, 1)
	var result float64
	for i := 0; i < 8; i++ {
		result = result*256 + float64(randBytes[i])
	}
	result /= 1 << 64

	return result
}

// SecureRandomGenerateSeed generates a new seed as a slice of JavaByte
func SecureRandomGenerateSeed(params []interface{}) interface{} {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("SecureRandomGenerateSeed: Expected 2 parameters (SecureRandom object, int64 numBytes), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	secureRandomObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGenerateSeed: First parameter must be a SecureRandom object")
	}

	if secureRandomObj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGenerateSeed: SecureRandom object cannot be nil")
	}

	numBytes, ok := params[1].(int64)
	if !ok || numBytes <= 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGenerateSeed: Second parameter must be a positive int64")
	}

	bytes := make([]byte, numBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomGenerateSeed: Failed to generate seed: %v", err))
	}

	result := object.JavaByteArrayFromGoByteArray(bytes)
	return result
}

// Get PRNG algorithm. This is simply whatever the O/S provides. Reference: Go package crypto/rand.
func SecureRandomGetAlgorithm(params []interface{}) interface{} {
	return object.StringObjectFromGoString(runtime.GOOS)
}

// toString - return the class name string.
func SecureRandomToString(params []interface{}) interface{} {
	return object.StringObjectFromGoString(secureRandomClassName)
}

func SecureRandomNextBoolean(params []interface{}) interface{} {
	randByte := make([]byte, 1) // float64 is 8 bytes
	_, err := rand.Read(randByte)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextBoolean: Failed to generate random byte: %v", err))
	}
	if randByte[0]&0x01 == 0x01 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}
