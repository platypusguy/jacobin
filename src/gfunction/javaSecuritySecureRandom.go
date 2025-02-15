/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"math/rand"
	"time"
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
			GFunction:  SecureRandomInit,
		}

	MethodSignatures["java/security/SecureRandom.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomInit,
		}

	MethodSignatures["java/security/SecureRandom.<init>(Ljava/security/SecureRandomSpi;Ljava/security/Provider;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
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
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;Ljava/lang/String;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getInstance(Ljava/lang/String;Ljava/security/SecureRandomParameters;Ljava/security/Provider;)Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.getInstanceStrong()Ljava/security/SecureRandom;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
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
			GFunction:  SecureRandomGetSeed,
		}

	MethodSignatures["java/security/SecureRandom.next(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
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
			GFunction:  SecureRandomReseed,
		}

	MethodSignatures["java/security/SecureRandom.reseed(Ljava/security/SecureRandomParameters;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/security/SecureRandom.setSeed([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomSetSeed,
		}

	MethodSignatures["java/security/SecureRandom.setSeed(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  SecureRandomSetSeed,
		}

	MethodSignatures["java/security/SecureRandom.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  SecureRandomToString,
		}

}

var secureRandomClassName = "java.security.SecureRandom"

// Return a byte array holding a generated seed of the specified byte size (count).
func genSeed(count int64) []byte {
	seed := time.Now().UnixNano()
	byteArray := types.Int64ToBytesBE(seed)
	if count < 0 {
		count = 0
	}
	return byteArray[:count]
}

// Re-seed and update rng of a SecureRandom object with the specified new seed expressed as an int64.
func reSeedObject(obj *object.Object, newSeed int64) {

	// New random number generator based on the specified seed.
	rng := rand.New(rand.NewSource(newSeed))

	// Update the seed field.
	obj.FieldTable["seed"] = object.Field{
		Ftype:  types.Int,
		Fvalue: newSeed,
	}

	// Update the rng field.
	obj.FieldTable["rng"] = object.Field{
		Ftype:  types.RNG,
		Fvalue: rng,
	}

}

// SecureRandomInit - instantiate a SecureRandom object.
func SecureRandomInit(params []interface{}) interface{} {
	var seed int64
	switch len(params) {
	case 1: // no parameters furnished - <init>()
		// Create default seed.
		seed = time.Now().UnixNano()
	case 2: // seed byte array was furnished - <init>([B)
		fld, ok := params[1].(*object.Object).FieldTable["value"]
		if !ok {
			errMsg := "SecureRandomInit: in byte array field \"value\" missing"
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		var byteArray []byte
		switch fld.Fvalue.(type) {
		case []byte:
			byteArray = fld.Fvalue.([]byte)
		case []types.JavaByte:
			byteArray = object.GoByteArrayFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		default:
			errMsg := fmt.Sprintf("SecureRandomInit: unrecognized type for field \"value\", observed: %T", fld.Fvalue)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		seed = types.BytesToInt64BE(byteArray)
	default:
		errMsg := fmt.Sprintf("SecureRandomInit: Number of parameters (%d) is not correct", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Create SecureRandom object with default seed value.
	obj := object.MakeEmptyObjectWithClassName(&secureRandomClassName)
	reSeedObject(obj, seed)

	// Return SecureRandom object to caller.
	return obj

}

// Re-seed this SecureRandom object.
func SecureRandomReseed(params []interface{}) interface{} {

	// Validate object.
	_, ok := params[0].(*object.Object).FieldTable["rng"].Fvalue.([]byte)
	if !ok {
		errMsg := "SecureRandomInit: field rng missing"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Re-seed object.
	reSeedObject(params[0].(*object.Object), time.Now().UnixNano())
	return nil
}

// Set the specified seed in this SecureRandom object.
func SecureRandomSetSeed(params []interface{}) interface{} {

	// Validate object.
	_, ok := params[0].(*object.Object).FieldTable["rng"].Fvalue.(*rand.Rand)
	if !ok {
		errMsg := "SecureRandomInit: field rng missing"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate seed parameter.
	switch params[1].(type) {
	case int64:
		reSeedObject(params[0].(*object.Object), params[1].(int64))
		return nil
	case *object.Object:
	default:
		errMsg := fmt.Sprintf("SecureRandomSetSeed: seed parameter must be int64 or an object, observed: %T", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get seed field.
	fld, ok := params[1].(*object.Object).FieldTable["value"]
	if !ok {
		errMsg := "SecureRandomSetSeed: parameter field \"value\" missing"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Use seed to re-seed the object.
	switch fld.Fvalue.(type) {
	case []byte:
		ii := types.BytesToInt64BE(fld.Fvalue.([]byte))
		reSeedObject(params[0].(*object.Object), ii)
	case []types.JavaByte:
		bb := object.GoByteArrayFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		ii := types.BytesToInt64BE(bb)
		reSeedObject(params[0].(*object.Object), ii)
	case int64:
		reSeedObject(params[0].(*object.Object), fld.Fvalue.(int64))
	default:
		errMsg := fmt.Sprintf("SecureRandomSetSeed: unrecognized type for field \"value\", observed: %T", fld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return nil
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
	rng, ok := secureRandomObj.FieldTable["rng"].Fvalue.(*rand.Rand)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextBytes: SecureRandom object missing \"rng\" field")
	}

	baObject, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextBytes: Second parameter must be a byte array")
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
		errMsg := fmt.Sprintf("SecureRandomNextBytes: unrecognized type for field \"value\", observed: %T", fld.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	_, err := rng.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextBytes: rng.Read(byteArray), err: %v", err))
	}

	result := object.JavaByteArrayFromGoByteArray(byteArray)
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

	rng, ok := secureRandomObj.FieldTable["rng"].Fvalue.(*rand.Rand)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextNextInt: SecureRandom object missing \"rng\" field")
	}

	// Generate random int64
	var result int64
	byteArray := make([]byte, 8) // int64 is 8 bytes
	_, err := rng.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextInt: Failed to generate random int64: %v", err))
	}

	// Convert bytes to int64
	for i := 0; i < 8; i++ {
		result = (result << 8) | int64(byteArray[i])
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

	rng, ok := secureRandomObj.FieldTable["rng"].Fvalue.(*rand.Rand)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomNextFloat: SecureRandom object missing \"rng\" field")
	}

	// Generate random float64 in the range [0, 1)
	byteArray := make([]byte, 8) // float64 is 8 bytes
	_, err := rng.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomNextFloat: rng.Read(byteArray) failed, err: %v", err))
	}

	// Convert bytes to a value in [0, 1)
	var result float64
	for i := 0; i < 8; i++ {
		result = result*256 + float64(byteArray[i])
	}
	result /= 1 << 64

	return result
}

// SecureRandomGenerateSeed generates a new seed as a slice of JavaByte
func SecureRandomGenerateSeed(params []interface{}) interface{} {

	// Validate parameters and set up rng.
	if len(params) != 2 {
		errMsg := fmt.Sprintf("SecureRandomGenerateSeed: Expected 2 parameters (SecureRandom object, int64 numBytes), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	secureRandomObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGenerateSeed: First parameter must be a SecureRandom object")
	}
	rng, ok := secureRandomObj.FieldTable["rng"].Fvalue.(*rand.Rand)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGenerateSeed: SecureRandom object missing \"rng\" field")
	}

	// Get seed byte array size.
	numBytes, ok := params[1].(int64)
	if !ok || numBytes <= 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGenerateSeed: Second parameter must be a positive int64")
	}

	// Generate output byte array.
	byteArray := genSeed(numBytes)
	_, err := rng.Read(byteArray)
	if err != nil {
		return getGErrBlk(excNames.RuntimeException, fmt.Sprintf("SecureRandomGenerateSeed: rng.Read(byteArray) failed, err: %v", err))
	}

	classStr := "[B"
	result := object.MakePrimitiveObject(classStr, types.ByteArray, byteArray)
	return result
}

// Get PRNG algorithm. This is simply whatever the O/S provides. Reference: Go package crypto/rand.
func SecureRandomGetAlgorithm(params []interface{}) interface{} {
	return object.StringObjectFromGoString("NativePRNG")
}

// toString - return the class name string.
func SecureRandomToString(params []interface{}) interface{} {
	return object.StringObjectFromGoString("NativePRNG")
}

func SecureRandomNextBoolean(params []interface{}) interface{} {
	rng := params[0].(*object.Object).FieldTable["rng"].Fvalue.(*rand.Rand)
	ii := rng.Int()
	if ii&0x01 == 0x01 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// SecureRandomGetSeed returns the current seed as a byte array of the specified size
func SecureRandomGetSeed(params []interface{}) interface{} {

	// Validate parameter count.
	if len(params) != 1 {
		errMsg := fmt.Sprintf("SecureRandomGetSeed: Expected 1 parameter (int64 size), got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	size, ok := params[0].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "SecureRandomGetSeed: Size parameter must be an int64")
	}

	// Form byte array for return to caller.
	byteArray := genSeed(size)
	classStr := "[B"
	result := object.MakePrimitiveObject(classStr, types.ByteArray, byteArray[:size])
	return result

}
