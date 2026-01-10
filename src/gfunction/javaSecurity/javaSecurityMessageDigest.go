package javaSecurity

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

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

	ghelpers.MethodSignatures["java/security/MessageDigest.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigClone,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.digest()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigDigest,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.digest([B)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  msgdigDigestBytes,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.digest([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  msgdigDigestBytesII,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineDigest()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapProtected,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineDigest([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapProtected,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineGetDigestLength()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapProtected,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineReset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapProtected,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineUpdate(B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapProtected,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.engineUpdate([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapProtected,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigGetAlgorithm,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.getDigestLength()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigGetDigestLength,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.getInstance(Ljava/lang/String;)Ljava/security/MessageDigest;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  msgdigGetInstance,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/MessageDigest;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  msgdigGetInstanceProvider,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljava/security/MessageDigest;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  msgdigGetInstanceProviderObj,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigGetProvider,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.isEqual([B[B)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  msgdigIsEqual,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigReset,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  msgdigToString,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.update(B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  msgdigUpdateByte,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.update(Ljava/nio/ByteBuffer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.update([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  msgdigUpdateBytes,
		}

	ghelpers.MethodSignatures["java/security/MessageDigest.update([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  msgdigUpdateBytesII,
		}
}

// ===================== Helper Functions =====================

// getHashForAlgorithm returns a hash.Hash for the given algorithm name
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
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// getDigestLengthForAlgorithm returns the digest length in bytes for the given algorithm
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

// ===================== MessageDigest Functions =====================

// getInstance(Ljava/lang/String;)Ljava/security/MessageDigest;
func msgdigGetInstance(params []any) any {
	algorithmObj := params[0].(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Validate algorithm
	_, err := getHashForAlgorithm(algorithm)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, err.Error())
	}

	// Create MessageDigest object
	className := "java/security/MessageDigest"
	md := object.MakeEmptyObjectWithClassName(&className)

	// Store algorithm name
	md.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algorithm),
	}

	// Store provider
	md.FieldTable["provider"] = object.Field{
		Ftype:  "Ljava/security/Provider;",
		Fvalue: ghelpers.GetDefaultSecurityProvider(),
	}

	// Initialize empty buffer for accumulating data
	md.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	return md
}

// getInstance(Ljava/lang/String;Ljava/lang/String;)Ljava/security/MessageDigest;
func msgdigGetInstanceProvider(params []any) any {
	algorithmObj := params[0].(*object.Object)
	providerObj := params[1].(*object.Object)

	providerName := object.GoStringFromStringObject(providerObj)

	// Only accept our security provider
	if providerName != types.SecurityProviderName {
		return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException,
			fmt.Sprintf("Provider %s not supported. Only %s is supported.", providerName, types.SecurityProviderName))
	}

	// Delegate to single-parameter getInstance
	return msgdigGetInstance([]any{algorithmObj})
}

// getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljava/security/MessageDigest;
func msgdigGetInstanceProviderObj(params []any) any {
	algorithmObj := params[0].(*object.Object)
	providerObj := params[1].(*object.Object)

	// Check if provider is our default provider
	if providerObj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Provider cannot be null")
	}

	// Get provider name from the provider object
	providerNameField, ok := providerObj.FieldTable["name"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid provider object")
	}

	providerNameObj := providerNameField.Fvalue.(*object.Object)
	providerName := object.GoStringFromStringObject(providerNameObj)

	// Only accept our security provider
	if providerName != types.SecurityProviderName {
		return ghelpers.GetGErrBlk(excNames.ProviderNotFoundException,
			fmt.Sprintf("Provider %s not supported. Only %s is supported.", providerName, types.SecurityProviderName))
	}

	// Delegate to single-parameter getInstance
	return msgdigGetInstance([]any{algorithmObj})
}

// getAlgorithm()Ljava/lang/String;
func msgdigGetAlgorithm(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["algorithm"].Fvalue.(*object.Object)
}

// getProvider()Ljava/security/Provider;
func msgdigGetProvider(params []any) any {
	this := params[0].(*object.Object)
	return this.FieldTable["provider"].Fvalue.(*object.Object)
}

// getDigestLength()I
func msgdigGetDigestLength(params []any) any {
	this := params[0].(*object.Object)
	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)
	return int64(getDigestLengthForAlgorithm(algorithm))
}

// update(B)V
func msgdigUpdateByte(params []any) any {
	this := params[0].(*object.Object)
	b := params[1].(int64)

	// Get current buffer
	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)

	// Append byte
	buffer = append(buffer, types.JavaByte(b))

	// Update buffer
	this.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: buffer,
	}

	return nil
}

// update([B)V
func msgdigUpdateBytes(params []any) any {
	this := params[0].(*object.Object)
	bytesObj := params[1].(*object.Object)

	// Get bytes from object
	bytesField := bytesObj.FieldTable["value"]
	var bytes []types.JavaByte
	switch v := bytesField.Fvalue.(type) {
	case []types.JavaByte:
		bytes = v
	case []byte:
		bytes = object.JavaByteArrayFromGoByteArray(v)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid byte array")
	}

	// Get current buffer
	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)

	// Append bytes
	buffer = append(buffer, bytes...)

	// Update buffer
	this.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: buffer,
	}

	return nil
}

// update([BII)V
func msgdigUpdateBytesII(params []any) any {
	this := params[0].(*object.Object)
	bytesObj := params[1].(*object.Object)
	offset := params[2].(int64)
	length := params[3].(int64)

	// Get bytes from object
	bytesField := bytesObj.FieldTable["value"]
	var bytes []types.JavaByte
	switch v := bytesField.Fvalue.(type) {
	case []types.JavaByte:
		bytes = v
	case []byte:
		bytes = object.JavaByteArrayFromGoByteArray(v)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid byte array")
	}

	// Validate offset and length
	if offset < 0 || length < 0 || offset+length > int64(len(bytes)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "Invalid offset or length")
	}

	// Get current buffer
	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)

	// Append bytes from offset to offset+length
	buffer = append(buffer, bytes[offset:offset+length]...)

	// Update buffer
	this.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: buffer,
	}

	return nil
}

// digest()[B
func msgdigDigest(params []any) any {
	this := params[0].(*object.Object)

	// Get algorithm
	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Get hash function
	h, err := getHashForAlgorithm(algorithm)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, err.Error())
	}

	// Get buffer
	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)

	// Convert to Go bytes and compute hash
	goBytes := object.GoByteArrayFromJavaByteArray(buffer)
	h.Write(goBytes)
	digest := h.Sum(nil)

	// Reset buffer
	this.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	// Convert digest to Java bytes and return as byte array object
	javaDigest := object.JavaByteArrayFromGoByteArray(digest)
	return object.StringObjectFromJavaByteArray(javaDigest)
}

// digest([B)[B
func msgdigDigestBytes(params []any) any {
	this := params[0].(*object.Object)
	bytesObj := params[1].(*object.Object)

	// Update with the provided bytes
	result := msgdigUpdateBytes([]any{this, bytesObj})
	if result != nil {
		return result // Error occurred
	}

	// Compute and return digest
	return msgdigDigest([]any{this})
}

// digest([BII)I
func msgdigDigestBytesII(params []any) any {
	this := params[0].(*object.Object)
	bufObj := params[1].(*object.Object)
	offset := params[2].(int64)
	length := params[3].(int64)

	// Get algorithm
	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	// Get expected digest length
	digestLen := getDigestLengthForAlgorithm(algorithm)

	// Validate that buffer has enough space
	bufField := bufObj.FieldTable["value"]
	var buf []types.JavaByte
	switch v := bufField.Fvalue.(type) {
	case []types.JavaByte:
		buf = v
	case []byte:
		buf = object.JavaByteArrayFromGoByteArray(v)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid byte array")
	}

	if offset < 0 || length < int64(digestLen) || offset+int64(digestLen) > int64(len(buf)) {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "Buffer too small or invalid offset")
	}

	// Compute digest
	digestObj := msgdigDigest([]any{this})
	if errBlk, ok := digestObj.(*ghelpers.GErrBlk); ok {
		return errBlk
	}

	// Get digest bytes
	digestBytes := digestObj.(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)

	// Copy digest into buffer at offset
	copy(buf[offset:], digestBytes)

	// Update the buffer object
	bufObj.FieldTable["value"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: buf,
	}

	// Return number of bytes written
	return int64(len(digestBytes))
}

// reset()V
func msgdigReset(params []any) any {
	this := params[0].(*object.Object)

	// Reset buffer
	this.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: []types.JavaByte{},
	}

	return nil
}

// isEqual([B[B)Z - static method
func msgdigIsEqual(params []any) any {
	digesta := params[0].(*object.Object)
	digestb := params[1].(*object.Object)

	// Get bytes from both objects
	bytesA := digesta.FieldTable["value"].Fvalue.([]types.JavaByte)
	bytesB := digestb.FieldTable["value"].Fvalue.([]types.JavaByte)

	// Compare using constant-time comparison
	if object.JavaByteArrayEquals(bytesA, bytesB) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// toString()Ljava/lang/String;
func msgdigToString(params []any) any {
	this := params[0].(*object.Object)

	algorithmObj := this.FieldTable["algorithm"].Fvalue.(*object.Object)
	algorithm := object.GoStringFromStringObject(algorithmObj)

	str := fmt.Sprintf("MessageDigest[%s]", algorithm)
	return object.StringObjectFromGoString(str)
}

// clone()Ljava/lang/Object;
func msgdigClone(params []any) any {
	this := params[0].(*object.Object)

	// Create new MessageDigest object
	className := "java/security/MessageDigest"
	clone := object.MakeEmptyObjectWithClassName(&className)

	// Copy algorithm
	clone.FieldTable["algorithm"] = this.FieldTable["algorithm"]

	// Copy provider
	clone.FieldTable["provider"] = this.FieldTable["provider"]

	// Copy buffer (make a new slice)
	buffer := this.FieldTable["buffer"].Fvalue.([]types.JavaByte)
	bufferCopy := make([]types.JavaByte, len(buffer))
	copy(bufferCopy, buffer)
	clone.FieldTable["buffer"] = object.Field{
		Ftype:  types.ByteArray,
		Fvalue: bufferCopy,
	}

	return clone
}
