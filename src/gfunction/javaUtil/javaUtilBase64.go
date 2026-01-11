/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

// Implementation of some of the functions in Java/util/Base64.
// Strategy: Locale = jacobin Object wrapping a Go string.

func Load_Util_Base64() {

	ghelpers.MethodSignatures["java/util/Base64.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Base64$Decoder.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Base64.getDecoder()Ljava/util/Base64$Decoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64GetStdDecoder,
		}

	ghelpers.MethodSignatures["java/util/Base64.getEncoder()Ljava/util/Base64$Encoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64GetStdEncoder,
		}

	ghelpers.MethodSignatures["java/util/Base64.getMimeDecoder()Ljava/util/Base64$Decoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64GetMimeDecoder,
		}

	ghelpers.MethodSignatures["java/util/Base64.getMimeEncoder()Ljava/util/Base64$Encoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64GetMimeEncoder,
		}

	ghelpers.MethodSignatures["java/util/Base64.getMimeEncoder(I[B)Ljava/util/Base64$Encoder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Base64.getUrlDecoder()Ljava/util/Base64$Decoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64GetUrlDecoder,
		}

	ghelpers.MethodSignatures["java/util/Base64.getUrlEncoder()Ljava/util/Base64$Encoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64GetUrlEncoder,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.encode([B)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  base64EncodeBsrc,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.encode([B[B)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  base64EncodeBsrcBdst,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.encode(Ljava/nio/ByteBuffer;)Ljava/nio/ByteBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.encodeToString([B)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  base64EncodeBsrcToString,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.withoutPadding()Ljava/util/Base64$Encoder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  base64WithoutPadding,
		}

	ghelpers.MethodSignatures["java/util/Base64$Encoder.wrap(Ljava/io/OutputStream;)Ljava/io/OutputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Base64$Decoder.decode([B)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  base64Decode,
		}

	ghelpers.MethodSignatures["java/util/Base64$Decoder.decode([B[B)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  base64DecodeBsrcBdst,
		}

	ghelpers.MethodSignatures["java/util/Base64$Decoder.decode(Ljava/lang/String;)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  base64Decode,
		}

	ghelpers.MethodSignatures["java/util/Base64$Decoder.decode(Ljava/nio/ByteBuffer;)Ljava/nio/ByteBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Base64$Decoder.wrap(Ljava/io/InputStream;)Ljava/io/InputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

}

// Classes:
var classNameBase64Decoder = "java/util/Base64/Decoder"
var classNameBase64Encoder = "java/util/Base64/Encoder"

// Schemes:
var base64SchemeStd = int64(1)
var base64SchemeStdRaw = int64(101) // no padding
var base64SchemeUrl = int64(2)
var base64SchemeUrlRaw = int64(102) // no padding
var base64SchemeMime = int64(3)     // essentially, the same as base64SchemeStd

// Fields:
var fieldNameScheme = "scheme"
var fieldNameValue = "value"

// GetDecoder returns a standard Base64 decoder instance.
func base64GetStdDecoder([]interface{}) interface{} {
	return object.MakeOneFieldObject(classNameBase64Decoder, fieldNameScheme, types.Int, base64SchemeStd)
}

// GetEncoder returns a standard Base64 encoder instance.
func base64GetStdEncoder([]interface{}) interface{} {
	return object.MakeOneFieldObject(classNameBase64Encoder, fieldNameScheme, types.Int, base64SchemeStd)
}

// GetDecoder returns a Mime Base64 decoder instance.
func base64GetMimeDecoder([]interface{}) interface{} {
	return object.MakeOneFieldObject(classNameBase64Decoder, fieldNameScheme, types.Int, base64SchemeMime)
}

// GetEncoder returns a Mime Base64 encoder instance.
func base64GetMimeEncoder([]interface{}) interface{} {
	return object.MakeOneFieldObject(classNameBase64Encoder, fieldNameScheme, types.Int, base64SchemeMime)
}

// GetDecoder returns a URL Base64 decoder instance.
func base64GetUrlDecoder([]interface{}) interface{} {
	return object.MakeOneFieldObject(classNameBase64Decoder, fieldNameScheme, types.Int, base64SchemeUrl)
}

// GetEncoder returns a URL Base64 encoder instance.
func base64GetUrlEncoder([]interface{}) interface{} {
	return object.MakeOneFieldObject(classNameBase64Encoder, fieldNameScheme, types.Int, base64SchemeUrl)
}

// Returns an encoder instance that encodes equivalently to this one,
// but without adding any padding character at the end of the encoded byte data.
// The encoding scheme of this encoder instance is unaffected by this invocation.
// The returned encoder instance should be used for non-padding encoding operation.
func base64WithoutPadding(params []interface{}) interface{} {

	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "base64WithoutPadding: base-object parameter is not an object")
	}

	observedClassName := *stringPool.GetStringPointer(obj.KlassName)
	if observedClassName != classNameBase64Encoder {
		errMsg := fmt.Sprintf("base64WithoutPadding: Expected class: %s, observed class: %s",
			classNameBase64Encoder, observedClassName)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	schemeField, ok := obj.FieldTable[fieldNameScheme]
	if !ok {
		errMsg := fmt.Sprintf("base64WithoutPadding: Scheme field not found: %s", fieldNameScheme)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	scheme, ok := schemeField.Fvalue.(int64)
	if !ok {
		errMsg := fmt.Sprintf("base64WithoutPadding: Corrupt %s field: %v, %T",
			fieldNameScheme, schemeField.Ftype, schemeField.Fvalue)
		return errMsg
	}

	// Compute the Base 64 encoding new scheme.
	newScheme := scheme
	switch scheme {
	case base64SchemeStd:
		newScheme = base64SchemeStdRaw
	case base64SchemeUrl:
		newScheme = base64SchemeUrlRaw
	}

	// Return the new Base 64 encoding object with the new scheme.
	return object.MakeOneFieldObject(classNameBase64Encoder, fieldNameScheme, types.Int, newScheme)
}

// Base 64 encode all bytes from the specified byte array into a newly-allocated byte array using the Base64 encoding scheme.
// Return the byte array in an object to caller.
func base64EncodeBsrc(params []interface{}) interface{} {
	// Do the Base 64 encoding.
	jba, ok := _encodeBase64(params)
	if !ok {
		return jba
	}

	// Return the result array in an object.
	return object.MakePrimitiveObject("[B", types.ByteArray, jba)
}

// Base 64 encode all bytes from the specified byte array using the Base64 encoding scheme specified in params[0].
// Write the resulting bytes to the output byte array specified in params[2].
func base64EncodeBsrcBdst(params []interface{}) interface{} {
	// Do the Base 64 encoding.
	jba, ok := _encodeBase64(params)
	if !ok {
		return jba
	}

	// Get the object holding the destination array.
	dstObject, ok := params[2].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("base64EncodeBsrcBdst: Destination argument should be an object, observed: %T", params[2])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get field holding the destination array.
	fld, ok := dstObject.FieldTable[fieldNameValue]
	if !ok {
		errMsg := fmt.Sprintf("base64EncodeBsrcBdst: Field missing in destination object: %s", fieldNameValue)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Put the encoded array in the destination object.
	fld.Ftype = types.ByteArray
	fld.Fvalue = jba.([]types.JavaByte)
	dstObject.FieldTable[fieldNameValue] = fld

	// Return array length to caller.
	return int64(len(jba.([]types.JavaByte)))
}

// Base 64 encode all bytes from the specified byte array using the Base64 encoding scheme specified in params[0].
// Write the resulting bytes to a newly-allocated String.
// Return the String to caller.
func base64EncodeBsrcToString(params []interface{}) interface{} {
	// Do the Base 64 encoding.
	jba, ok := _encodeBase64(params)
	if !ok {
		return jba
	}

	// Return the result array in a String object.
	return object.StringObjectFromJavaByteArray(jba.([]types.JavaByte))
}

// Base 64 decode all bytes from the specified byte array into a newly-allocated byte array using the Base64 encoding scheme.
// Return the byte array in an object to caller.
func base64Decode(params []interface{}) interface{} {
	// Do the Base 64 decoding.
	jba, ok := _decodeBase64(params)
	if !ok {
		return jba
	}

	// Return the result array in an object.
	return object.MakePrimitiveObject("[B", types.ByteArray, jba)
}

// Base 64 decode all bytes from the specified byte array using the Base64 encoding scheme specified in params[0].
// Write the resulting bytes to the output byte array specified in params[2].
func base64DecodeBsrcBdst(params []interface{}) interface{} {
	// Do the Base 64 decoding.
	jba, ok := _decodeBase64(params)
	if !ok {
		return jba
	}

	// Get the object holding the destination array.
	dstObject, ok := params[2].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("base64DecodeBsrcBdst: Destination argument should be an object, observed: %T", params[2])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get field holding the destination array.
	fld, ok := dstObject.FieldTable[fieldNameValue]
	if !ok {
		errMsg := fmt.Sprintf("base64DecodeBsrcBdst: Field missing in destination object: %s", fieldNameValue)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Put the encoded array in the destination object.
	fld.Ftype = types.ByteArray
	fld.Fvalue = jba.([]types.JavaByte)
	dstObject.FieldTable[fieldNameValue] = fld

	// Return array length to caller.
	return int64(len(jba.([]types.JavaByte)))
}

/*
_encodeBase64 encodes data in params[1] based on the specified encoding scheme field in params[0].
Return #1

	normal - JavaByte array, Base64-encoded according to the scheme.
	error - ghelpers.GetGErrBlk

Return #2

	true : success
	false : an error occurred
*/
func _encodeBase64(params []interface{}) (interface{}, bool) {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "_encodeBase64: expects at least 2 parameters"), false
	}

	// this: Base64 encoder object.
	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "_encodeBase64: base-object parameter is not an object"), false
	}

	// Validate class name.
	observedClassName := *stringPool.GetStringPointer(this.KlassName)
	if observedClassName != classNameBase64Encoder {
		errMsg := fmt.Sprintf("_encodeBase64: Expected class: %s, observed class: %s",
			classNameBase64Encoder, observedClassName)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
	}

	// Get scheme field.
	schemeField, ok := this.FieldTable[fieldNameScheme]
	if !ok {
		errMsg := fmt.Sprintf("_encodeBase64: Expected scheme: %s", fieldNameScheme)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg), false
	}
	scheme, ok := schemeField.Fvalue.(int64)
	if !ok {
		errMsg := fmt.Sprintf("_encodeBase64: Corrupt %s field: %v, %T",
			fieldNameScheme, schemeField.Ftype, schemeField.Fvalue)
		return errMsg, false
	}

	// Get source JavaByte array converted to []byte.
	srcObject, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "_encodeBase64: source parameter is not an object"), false
	}
	srcJB, ok := srcObject.FieldTable[fieldNameValue].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := fmt.Sprintf("_encodeBase64: source parameter has no field %s or its the wrong type", fieldNameValue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
	}
	src := object.GoByteArrayFromJavaByteArray(srcJB)

	// Base 64 encode, depending on scheme.
	var dst []byte
	switch scheme {
	case base64SchemeStd:
		stdEnc := base64.StdEncoding
		dst = make([]byte, stdEnc.EncodedLen(len(src)))
		stdEnc.Encode(dst, src)
	case base64SchemeStdRaw:
		stdEnc := base64.StdEncoding.WithPadding(base64.NoPadding)
		dst = make([]byte, stdEnc.EncodedLen(len(src)))
		stdEnc.Encode(dst, src)
	case base64SchemeUrl:
		urlEnc := base64.URLEncoding
		dst = make([]byte, urlEnc.EncodedLen(len(src)))
		urlEnc.Encode(dst, src)
	case base64SchemeUrlRaw:
		urlEnc := base64.URLEncoding.WithPadding(base64.NoPadding)
		dst = make([]byte, urlEnc.EncodedLen(len(src)))
		urlEnc.Encode(dst, src)
	case base64SchemeMime:
		var encodedBuffer bytes.Buffer
		mimeEncoder := base64.NewEncoder(base64.StdEncoding, &encodedBuffer)
		_, err := mimeEncoder.Write(src)
		if err != nil {
			errMsg := fmt.Sprintf("_encodeBase64: mimeEncoder.Write failed, err: %v", err)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
		}
		_ = mimeEncoder.Close() // Ensure encoding is flushed
		dst = encodedBuffer.Bytes()
	default:
		errMsg := fmt.Sprintf("_encodeBase64: Impossible %s field: %d", fieldNameScheme, scheme)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg), false
	}

	// Return JavaByte array result to caller.
	dstJB := object.JavaByteArrayFromGoByteArray(dst)
	return dstJB, true
}

/*
_decodeBase64 decodes data in params[1] based on the specified encoding scheme field in params[0].
Return #1

	normal - JavaByte array, Base64-encoded according to the scheme.
	error - ghelpers.GetGErrBlk

Return #2

	true : success
	false : an error occurred
*/
func _decodeBase64(params []interface{}) (interface{}, bool) {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "_decodeBase64: expects at least 2 parameters"), false
	}

	// this: Base64 encoder object.
	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "_decodeBase64: base-object parameter is not an object"), false
	}

	// Validate class name.
	observedClassName := *stringPool.GetStringPointer(this.KlassName)
	if observedClassName != classNameBase64Decoder {
		errMsg := fmt.Sprintf("_decodeBase64: Expected class: %s, observed class: %s",
			classNameBase64Decoder, observedClassName)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
	}

	// Get scheme field.
	schemeField, ok := this.FieldTable[fieldNameScheme]
	if !ok {
		errMsg := fmt.Sprintf("_decodeBase64: Expected scheme: %s", fieldNameScheme)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg), false
	}
	scheme, ok := schemeField.Fvalue.(int64)
	if !ok {
		errMsg := fmt.Sprintf("_decodeBase64: Corrupt %s field: %v, %T",
			fieldNameScheme, schemeField.Ftype, schemeField.Fvalue)
		return errMsg, false
	}

	// Get source JavaByte array converted to []byte.
	srcObject, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "_decodeBase64: source parameter is not an object"), false
	}
	srcJB, ok := srcObject.FieldTable[fieldNameValue].Fvalue.([]types.JavaByte)
	if !ok {
		errMsg := fmt.Sprintf("_decodeBase64: source parameter has no field %s or its the wrong type", fieldNameValue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
	}
	encoded := object.GoByteArrayFromJavaByteArray(srcJB)

	// Base 64 encode, depending on scheme.
	var decoded []byte
	switch scheme {
	case base64SchemeStd, base64SchemeStdRaw:
		decoder := base64.StdEncoding
		decoded = make([]byte, decoder.DecodedLen(len(encoded)))
		num, err := decoder.Decode(decoded, encoded)
		if err != nil {
			decoder = base64.StdEncoding.WithPadding(base64.NoPadding)
			decoded = make([]byte, decoder.DecodedLen(len(encoded)))
			num, err = decoder.Decode(decoded, encoded)
			if err != nil {
				errMsg := fmt.Sprintf("_decodeBase64: Decoding/standard error: %v", err)
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
			}
		}
		decoded = decoded[:num]
	case base64SchemeUrl, base64SchemeUrlRaw:
		decoder := base64.URLEncoding
		decoded = make([]byte, decoder.DecodedLen(len(encoded)))
		num, err := decoder.Decode(decoded, encoded)
		if err != nil {
			decoder = base64.URLEncoding.WithPadding(base64.NoPadding)
			decoded = make([]byte, decoder.DecodedLen(len(encoded)))
			num, err = decoder.Decode(decoded, encoded)
			if err != nil {
				errMsg := fmt.Sprintf("_decodeBase64: Decoding/URL error: %v", err)
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
			}
		}
		decoded = decoded[:num]
	case base64SchemeMime:
		decoded = make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
		num, err := base64.StdEncoding.Decode(decoded, encoded)
		if err != nil {
			errMsg := fmt.Sprintf("_decodeBase64: Decoding/MIME error: %v", err)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
		}
		decoded = decoded[:num]
	default:
		errMsg := fmt.Sprintf("_decodeBase64: Impossible %s field: %d", fieldNameScheme, scheme)
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg), false
	}

	// Return JavaByte array result to caller.
	jbo := object.JavaByteArrayFromGoByteArray(decoded)
	return jbo, true
}
