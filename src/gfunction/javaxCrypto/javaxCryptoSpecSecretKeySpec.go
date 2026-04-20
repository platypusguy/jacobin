package javaxCrypto

import (
	"bytes"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"slices"
)

func Load_Crypto_Spec_SecretKeySpec() {
	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.<init>([BIILjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  secretKeySpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.<init>([BLjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  secretKeySpecInit,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  secretKeySpecEquals,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  secretKeySpecGetAlgorithm,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.getEncoded()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  secretKeySpecGetEncoded,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.getFormat()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  secretKeySpecGetFormat,
		}

	ghelpers.MethodSignatures["javax/crypto/spec/SecretKeySpec.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  secretKeySpecHashCode,
		}
}

func secretKeySpecEquals(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecEquals: missing parameters")
	}

	// params[0] is 'this'
	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecEquals: 'this' is not an object")
	}

	// params[1] is the object to compare
	otherObj, ok := params[1].(*object.Object)
	if !ok || otherObj == nil {
		return int64(0) // false - not an object or null
	}

	// Fast path: same instance
	if thisObj == otherObj {
		return int64(1) // true
	}

	// Check if same class
	if thisObj.KlassName != otherObj.KlassName {
		return int64(0) // false
	}

	// Compare algorithms
	thisAlgoObj, ok1 := thisObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	otherAlgoObj, ok2 := otherObj.FieldTable["algorithm"].Fvalue.(*object.Object)
	if !ok1 || !ok2 {
		return int64(0) // false
	}
	thisAlgo := object.GoStringFromStringObject(thisAlgoObj)
	otherAlgo := object.GoStringFromStringObject(otherAlgoObj)
	if thisAlgo != otherAlgo {
		return int64(0) // false
	}

	// Compare key bytes
	thisKey, ok1 := thisObj.FieldTable["key"].Fvalue.([]byte)
	otherKey, ok2 := otherObj.FieldTable["key"].Fvalue.([]byte)
	if !ok1 || !ok2 {
		return int64(0) // false
	}

	if bytes.Equal(thisKey, otherKey) {
		return int64(1) // true
	}

	return int64(0) // false
}

func secretKeySpecGetAlgorithm(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecGetAlgorithm: missing 'this'")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecGetAlgorithm: 'this' is not an object")
	}

	algorithmObj, ok := obj.FieldTable["algorithm"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException,
			"secretKeySpecGetAlgorithm: algorithm field not found or invalid")
	}

	return algorithmObj
}

func secretKeySpecGetEncoded(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecGetEncoded: missing 'this'")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecGetEncoded: 'this' is not an object")
	}

	keyBytes, ok := obj.FieldTable["key"].Fvalue.([]byte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException,
			"secretKeySpecGetEncoded: key field not found or invalid")
	}

	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray,
		object.JavaByteArrayFromGoByteArray(slices.Clone(keyBytes)))
}

func secretKeySpecGetFormat(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecGetFormat: missing 'this'")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecGetFormat: 'this' is not an object")
	}

	// SecretKeySpec always uses RAW format
	return object.StringObjectFromGoString("RAW")
}

func secretKeySpecHashCode(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecHashCode: missing 'this'")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecHashCode: 'this' is not an object")
	}

	algorithmObj, ok1 := obj.FieldTable["algorithm"].Fvalue.(*object.Object)
	keyBytes, ok2 := obj.FieldTable["key"].Fvalue.([]byte)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException,
			"secretKeySpecHashCode: fields not found or invalid")
	}

	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Compute hash based on algorithm and key bytes
	// Using Java's algorithm: hash = algorithm.hashCode() ^ Arrays.hashCode(key)
	var hash int64 = 0

	// Hash the algorithm string (Java's String.hashCode algorithm)
	for i := 0; i < len(algorithm); i++ {
		hash = 31*hash + int64(algorithm[i])
	}

	// XOR with key bytes hash (Java's Arrays.hashCode algorithm)
	var keyHash int64 = 1
	for _, b := range keyBytes {
		keyHash = 31*keyHash + int64(b)
	}

	hash ^= keyHash

	return hash
}

func secretKeySpecInit(params []any) any {
	var algoObj *object.Object

	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecInit: insufficient parameters")
	}

	// params[0] is 'this'
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecInit: 'this' is not an object")
	}

	var keyBytes []byte
	var algorithm string

	if len(params) == 3 {
		// Constructor: SecretKeySpec(byte[] key, String algorithm)
		keyObj, ok := params[1].(*object.Object)
		if !ok || keyObj == nil {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"secretKeySpecInit: key is not a valid byte array")
		}
		keyBytes = object.GoByteArrayFromJavaByteArray(keyObj.FieldTable["value"].Fvalue.([]types.JavaByte))

		algoObj, ok = params[2].(*object.Object)
		if !ok || algoObj == nil {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"secretKeySpecInit: algorithm is not a valid string")
		}
		algorithm = object.GoStringFromStringObject(algoObj)

	} else if len(params) == 5 {
		// Constructor: SecretKeySpec(byte[] key, int offset, int len, String algorithm)
		keyObj, ok := params[1].(*object.Object)
		if !ok || keyObj == nil {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"secretKeySpecInit: key is not a valid byte array")
		}
		fullKeyBytes := object.GoByteArrayFromJavaByteArray(keyObj.FieldTable["value"].Fvalue.([]types.JavaByte))

		offset, ok := params[2].(int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"secretKeySpecInit: offset is not an integer")
		}

		length, ok := params[3].(int64)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"secretKeySpecInit: length is not an integer")
		}

		algoObj, ok = params[4].(*object.Object)
		if !ok || algoObj == nil {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
				"secretKeySpecInit: algorithm is not a valid string")
		}
		algorithm = object.GoStringFromStringObject(algoObj)

		// Validate offset and length
		if offset < 0 || length < 0 || int(offset+length) > len(fullKeyBytes) {
			return ghelpers.GetGErrBlk(excNames.InvalidKeyException,
				"secretKeySpecInit: invalid offset or length")
		}

		// Extract the subset of bytes
		keyBytes = fullKeyBytes[offset : offset+length]

	} else {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecInit: invalid number of parameters")
	}

	// Validate algorithm is not empty
	if algorithm == "" {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"secretKeySpecInit: algorithm cannot be empty")
	}

	// Get configuration for this algorithm.
	_, enabled := ValidateSecretKeySpecAlgorithm(algorithm)
	if !enabled {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("secretKeySpecInit: unknown or invalid SecretKeySpec algorithm: %s", algorithm))
	}

	// Store the key and algorithm in the object's fields
	obj.FieldTable["key"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: slices.Clone(keyBytes),
	}
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: slices.Clone(keyBytes),
	}
	obj.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: algoObj,
	}

	return nil
}
