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
	"time"
)

/***

Jacobin should employ the most secure SecureRandom implementation available.
The Go crypto/rand package is the best choice for Jacobin.

Algorithm, Provider: The Go crypto/rand package automatically uses the most secure cryptographic random source
available on the system. It does not allow selecting a specific pseudo-random generator or cryptographic provider
because it relies on O/S-backed randomness (e.g., /dev/urandom on Linux, CryptGenRandom on Windows).

Seeding: Cryptographic pseudo-random generators derive their randomness from a system entropy pool,
which is automatically initialized. Allowing user-provided seeding would weaken security by
making the output more predictable. Unlike math/rand, which uses a deterministic algorithm based on a seed,
crypto/rand is designed for secure applications like encryption keys and authentication tokens.

***/

func Load_Security_SecureRandom() {

	MethodSignatures["java/security/SecureRandom.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/security/SecureRandom.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomInit,
		}

	MethodSignatures["java/security/SecureRandom.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomInit,
		}

	MethodSignatures["java/security/SecureRandom.<init>(Ljava/security/SecureRandomSpi;Ljava/security/Provider;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.generateSeed(I)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomGenerateSeed,
		}

	MethodSignatures["java/security/SecureRandom.getAlgorithm()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomGetAlgorithm,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 1, // String algorithm
			GFunction:  secureRandomGetInstance,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2, // String algorithm, String provider
			GFunction:  secureRandomGetInstance,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2, // String algorithm, Provider provider
			GFunction:  secureRandomGetInstance,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2, // String algorithm, SecureRandomParameters params
			GFunction:  secureRandomGetInstance,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 3, // String algorithm, SecureRandomParameters params, String provider
			GFunction:  secureRandomGetInstance,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;Ljava/security/Provider;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 3, // String algorithm, SecureRandomParameters params, Provider provider
			GFunction:  secureRandomGetInstance,
		}

	MethodSignatures["java/security/SecureRandom.getInstanceStrong()Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomGetInstanceStrong,
		}

	MethodSignatures["java/security/SecureRandom.getParameters()Ljava/security/SecureRandomParameters;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomGetAlgorithm,
		}

	MethodSignatures["java/security/SecureRandom.getProvider()Ljava/security/Provider;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomGetAlgorithm,
		}

	MethodSignatures["java/security/SecureRandom.getSeed(I)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomGetSeed,
		}

	MethodSignatures["java/security/SecureRandom.next(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapProtected,
		}

	MethodSignatures["java/security/SecureRandom.nextBoolean()Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomNextBoolean,
		}

	MethodSignatures["java/security/SecureRandom.nextBytes([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomNextBytes,
		}

	MethodSignatures["java/security/SecureRandom.nextBytes([BLjava/security/SecureRandomParameters;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  secureRandomNextBytes,
		}

	MethodSignatures["java/security/SecureRandom.nextDouble()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomNextFloat,
		}

	MethodSignatures["java/security/SecureRandom.nextFloat()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomNextFloat,
		}

	MethodSignatures["java/security/SecureRandom.nextGaussian()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.nextInt()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomNextInt,
		}

	MethodSignatures["java/security/SecureRandom.nextLong()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomNextInt,
		}

	MethodSignatures["java/security/SecureRandom.reseed()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomReseed,
		}

	MethodSignatures["java/security/SecureRandom.reseed(Ljava/security/SecureRandomParameters;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomReseed,
		}

	MethodSignatures["java/security/SecureRandom.setSeed([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomSetSeed,
		}

	MethodSignatures["java/security/SecureRandom.setSeed(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomSetSeed,
		}

	MethodSignatures["java/security/SecureRandom.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomToString,
		}

}

var secureRandomClassName = "java/security/SecureRandom"

// Return a byte array holding a generated dummy seed of the specified byte size (count).
func _genSeed(count int64) []byte {
	seed := time.Now().UnixNano()
	byteArray := types.Int64ToBytesBE(seed)
	if count < 0 {
		count = 0
	}
	return byteArray[:count]
}

// Re-seed and update rng of a SecureRandom object with the specified new seed expressed as an int64.
func _reSeedObject(obj *object.Object, newSeed int64) {

	// Update the seed field.
	obj.FieldTable["seed"] = object.Field{
		Ftype:  types.Int,
		Fvalue: newSeed,
	}

}

// secureRandomInit - instantiate a SecureRandom object.
func secureRandomInit(params []interface{}) interface{} {

	// Create SecureRandom object with default seed value.
	obj := params[0].(*object.Object)
	seed := time.Now().UnixNano()
	_reSeedObject(obj, seed)

	return obj

}

// SecureRandomGetInstance - several variations of SecureRandom getInstance.
func secureRandomGetInstance(params []interface{}) interface{} {
	return secureRandomInit(params)
}

// secureRandomGetInstanceStrong.
func secureRandomGetInstanceStrong(params []interface{}) interface{} {

	// Create SecureRandom object with default seed value.
	obj := object.MakeEmptyObjectWithClassName(&secureRandomClassName)
	seed := time.Now().UnixNano()
	_reSeedObject(obj, seed)

	return obj

}

// Re-seed this SecureRandom object.
func secureRandomReseed(params []interface{}) interface{} {
	_reSeedObject(params[0].(*object.Object), time.Now().UnixNano())
	return nil
}

// Set the specified seed in this SecureRandom object.
func secureRandomSetSeed(params []interface{}) interface{} {

	// Validate seed parameter.
	switch params[1].(type) {
	case int64:
		_reSeedObject(params[0].(*object.Object), params[1].(int64))
		return nil
	case *object.Object:
	default:
		errMsg := fmt.Sprintf("secureRandomSetSeed: seed parameter must be int64 or an object, observed: %T", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get seed field.
	fld, ok := params[1].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := "secureRandomSetSeed: parameter field \"value\" missing"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Use seed to re-seed the object.
	switch fld.Fvalue.(type) {
	case []byte:
		ii := types.BytesToInt64BE(fld.Fvalue.([]byte))
		_reSeedObject(params[0].(*object.Object), ii)
	case []types.JavaByte:
		bb := object.GoByteArrayFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		ii := types.BytesToInt64BE(bb)
		_reSeedObject(params[0].(*object.Object), ii)
	case int64:
		_reSeedObject(params[0].(*object.Object), fld.Fvalue.(int64))
	default:
		errMsg := fmt.Sprintf("secureRandomSetSeed: unrecognized type for field \"value\", observed: %T", fld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return nil
}

// secureRandomNextBytes generates a specified number of random bytes
func secureRandomNextBytes(params []interface{}) interface{} {
	switch len(params) {
	case 2, 3:
	default:
		errMsg := fmt.Sprintf("secureRandomNextBytes: Wrong number of parameters, observed %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get byte array object.
	baObject, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "secureRandomNextBytes: Second parameter must be a byte array")
	}

	// Generate random bytes
	fld := baObject.FieldTable["value"]
	var byteArray []byte
	switch fld.Fvalue.(type) {
	case []byte:
		byteArray = fld.Fvalue.([]byte)
	case []types.JavaByte:
		byteArray = object.GoByteArrayFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
	default:
		errMsg := fmt.Sprintf("secureRandomNextBytes: unrecognized type for field \"value\", observed: %T", fld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	_, err := rand.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("secureRandomNextBytes: rng.Read(byteArray), err: %v", err))
	}

	// Return objectified Java byte array.
	fld.Fvalue = object.JavaByteArrayFromGoByteArray(byteArray)
	baObject.FieldTable["value"] = fld
	return nil
}

// secureRandomNextInt generates a random int64
func secureRandomNextInt(params []interface{}) interface{} {

	// Validate parameter count.
	if len(params) != 1 {
		errMsg := fmt.Sprintf("secureRandomNextInt: Expected 1 parameter (SecureRandom object), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Generate random int64.
	var result int64
	byteArray := make([]byte, 8) // int64 is 8 bytes
	_, err := rand.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("secureRandomNextInt: Failed to generate random int64: %v", err))
	}

	// Convert bytes to int64.
	for i := 0; i < 8; i++ {
		result = (result << 8) | int64(byteArray[i])
	}

	return result
}

// secureRandomNextFloat generates a random float64
func secureRandomNextFloat(params []interface{}) interface{} {

	// Validate parameter count.
	if len(params) != 1 {
		errMsg := fmt.Sprintf("secureRandomNextFloat: Expected 1 parameter (SecureRandom object), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Generate random float64 in the range [0, 1).
	byteArray := make([]byte, 8) // float64 is 8 bytes
	_, err := rand.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("secureRandomNextFloat: rng.Read(byteArray) failed, err: %v", err))
	}

	// Convert bytes to a value in [0, 1).
	var result float64
	for i := 0; i < 8; i++ {
		result = result*256 + float64(byteArray[i])
	}
	result /= 1 << 64

	return result
}

// secureRandomGenerateSeed generates a new seed as a slice of JavaByte
func secureRandomGenerateSeed(params []interface{}) interface{} {

	// Validate parameters and set up rng.
	if len(params) != 2 {
		errMsg := fmt.Sprintf("secureRandomGenerateSeed: Expected 2 parameters (SecureRandom object, int64 numBytes), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get seed byte array size.
	numBytes, ok := params[1].(int64)
	if !ok || numBytes <= 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "secureRandomGenerateSeed: Second parameter must be a positive int64")
	}

	// Generate output byte array.
	byteArray := _genSeed(numBytes)
	_, err := rand.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("secureRandomGenerateSeed: rng.Read(byteArray) failed, err: %v", err))
	}

	classStr := "[B"
	result := object.MakePrimitiveObject(classStr, types.ByteArray, byteArray)
	return result
}

// Get PRNG algorithm. This is simply whatever the O/S provides. Reference: Go package crypto/rand.
func secureRandomGetAlgorithm(params []interface{}) interface{} {
	return object.StringObjectFromGoString("go/crypto/rand") // matches OpenJDK JVM
}

// toString - return the class name string.
func secureRandomToString(params []interface{}) interface{} {
	return object.StringObjectFromGoString("go/crypto/rand")
}

func secureRandomNextBoolean(params []interface{}) interface{} {
	byteArray := make([]byte, 1) // int64 is 8 bytes
	_, err := rand.Read(byteArray)
	if err != nil {
		errMsg := fmt.Sprintf("secureRandomNextBoolean: rand.Read(byteArray) failed, err: %s", err.Error())
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if byteArray[0]&0x01 == 0x01 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// secureRandomGetSeed returns the current seed as a byte array of the specified size
func secureRandomGetSeed(params []interface{}) interface{} {

	// Validate parameter count.
	if len(params) != 1 {
		errMsg := fmt.Sprintf("secureRandomGetSeed: Expected 1 parameter (int64 size), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	size, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "secureRandomGetSeed: Size parameter must be an int64")
	}

	// Form byte array for return to caller.
	byteArray := _genSeed(size)
	classStr := "[B"
	result := object.MakePrimitiveObject(classStr, types.ByteArray, byteArray[:size])
	return result

}
