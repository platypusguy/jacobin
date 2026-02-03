package javaSecurity

import (
	"fmt"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

var SecurityProviderServices = map[string]map[string]func() *object.Object{
	"Runtime": {
		"Security": func() *object.Object {
			return NewGoRuntimeService("Runtime", "Security")
		},
	},
	"MessageDigest": {
		"MD5": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "MD5")
		},
		"SHA-1": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-1")
		},
		"SHA-224": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-224")
		},
		"SHA-256": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-256")
		},
		"SHA-384": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-384")
		},
		"SHA-512": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-512")
		},
		"SHA-512/224": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-512/224")
		},
		"SHA-512/256": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-512/256")
		},
	},
}

// NewGoRuntimeService creates a basic Provider$Service object for a Go runtime algorithm
func NewGoRuntimeService(typ, algo string) *object.Object {
	className := "java/security/Provider$Service"
	svc := object.MakeEmptyObjectWithClassName(&className)

	svc.FieldTable["type"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(typ),
	}
	svc.FieldTable["algorithm"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(algo),
	}
	svc.FieldTable["className"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(typ),
	}
	svc.FieldTable["aliases"] = object.Field{
		Ftype:  types.StringArrayClassName,
		Fvalue: []*object.Object{},
	}

	// Attributes map
	attributes := map[string]*object.Object{}
	blockSize := getBlockSizeForAlgorithm(algo)
	if blockSize > 0 {
		attributes["blockSize"] = object.StringObjectFromGoString(fmt.Sprintf("%d", blockSize))
	}

	svc.FieldTable["attributes"] = object.Field{
		Ftype:  types.Map,
		Fvalue: attributes,
	}

	return svc
}

func getBlockSizeForAlgorithm(algorithm string) int {
	switch strings.ToUpper(algorithm) {
	case "MD5", "SHA-1", "SHA-224", "SHA-256":
		return 64
	case "SHA-384", "SHA-512", "SHA-512/224", "SHA-512/256":
		return 128
	default:
		return 0
	}
}
