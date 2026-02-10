package javaSecurity

import (
	"crypto/dsa"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"math/big"
	"unsafe"
)

func Load_Security_Key() {

	ghelpers.MethodSignatures["java/security/Key.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/security/Key.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyGetAlgorithm,
		}

	ghelpers.MethodSignatures["java/security/Key.getEncoded()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyGetEncoded,
		}

	ghelpers.MethodSignatures["java/security/Key.getFormat()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyGetFormat,
		}

	ghelpers.MethodSignatures["java/security/Key.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyHashCode,
		}
}

func keyGetAlgorithm(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetAlgorithm: missing 'this'")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetAlgorithm: 'this' is not an object")
	}

	alg, ok := obj.FieldTable["algorithm"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetAlgorithm: algorithm not found")
	}
	return alg
}

func keyGetFormat(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetFormat: missing 'this'")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetFormat: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(obj.KlassName)
	var format string

	switch className {
	case types.ClassNameRSAPublicKey, types.ClassNameECPublicKey, types.ClassNameDSAPublicKey,
		types.ClassNameEd25519PublicKey, types.ClassNameEd448PublicKey,
		types.ClassNameX25519PublicKey, types.ClassNameX448PublicKey,
		types.ClassNameDHPublicKey:
		format = "X.509"
	case types.ClassNameRSAPrivateKey, types.ClassNameECPrivateKey, types.ClassNameDSAPrivateKey,
		types.ClassNameEd25519PrivateKey, types.ClassNameEd448PrivateKey,
		types.ClassNameX25519PrivateKey, types.ClassNameX448PrivateKey,
		types.ClassNameDHPrivateKey:
		format = "PKCS#8"
	default:
		format = "RAW"
	}

	return object.StringObjectFromGoString(format)
}

func keyGetEncoded(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetEncoded: missing 'this'")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyGetEncoded: 'this' is not an object")
	}

	val := obj.FieldTable["value"].Fvalue
	var encoded []byte
	var err error = nil

	switch k := val.(type) {
	case *rsa.PublicKey:
		encoded, err = x509.MarshalPKIXPublicKey(k)
	case *rsa.PrivateKey:
		encoded, err = x509.MarshalPKCS8PrivateKey(k)
	case *ecdsa.PublicKey:
		encoded, err = x509.MarshalPKIXPublicKey(k)
	case *ecdsa.PrivateKey:
		encoded, err = x509.MarshalPKCS8PrivateKey(k)
	case *dsa.PublicKey:
		encoded, err = x509.MarshalPKIXPublicKey(k)
	case *dsa.PrivateKey:
		encoded, err = x509.MarshalPKCS8PrivateKey(k)
	case ed25519.PublicKey:
		encoded, err = x509.MarshalPKIXPublicKey(k)
	case ed25519.PrivateKey:
		encoded, err = x509.MarshalPKCS8PrivateKey(k)
	case *ecdh.PublicKey:
		// ecdh keys can be converted to pkix/pkcs8 via x509 in Go 1.20+
		encoded, err = x509.MarshalPKIXPublicKey(k)
	case *ecdh.PrivateKey:
		encoded, err = x509.MarshalPKCS8PrivateKey(k)
	case []byte:
		// Some keys (X25519, X448) might be stored as raw bytes
		encoded = k
	case *big.Int:
		// DH keys often stored as *big.Int
		encoded = k.Bytes()
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("keyGetEncoded: unsupported key type %T", val))
	}

	if err != nil {
		return ghelpers.GetGErrBlk(excNames.GeneralSecurityException, "keyGetEncoded: encoding failed: "+err.Error())
	}

	return object.MakePrimitiveObject(types.ByteArray, types.ByteArray, object.JavaByteArrayFromGoByteArray(encoded))
}

func keyHashCode(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyHashCode: missing 'this'")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keyHashCode: 'this' is not an object")
	}

	ptr := uintptr(unsafe.Pointer(obj))
	hashCode := int64(ptr ^ (ptr >> 32))
	return hashCode
}
