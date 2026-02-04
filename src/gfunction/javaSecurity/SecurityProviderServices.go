package javaSecurity

import (
	"fmt"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

// InitDefaultSecurityProvider: initialize the one and only security provider.
func InitDefaultSecurityProvider() *object.Object {

	// Set up the security provider object.
	provider := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProvider)

	// Initialize the security provider with default values.
	securityProviderInit([]any{
		provider,
		object.StringObjectFromGoString(types.SecurityProviderName),
		types.SecurityProviderVersion,
		object.StringObjectFromGoString(types.SecurityProviderInfo),
	})

	// For each distinct type-algorithm entry in SecurityProviderServices,
	// create a Provider$Service object and add it to the provider.
	// Note that the provider to service is a one to many relationship
	// and each service has a provider field as a convenience.
	// -------------------------------------------------------------------
	// For each SecurityProviderServices type entry (outer map), get a list of algorithm maps.
	for _, algos := range SecurityProviderServices {
		// For each algorithm map, extract its service initialization function.
		for _, serviceInit := range algos {
			// Use the serviceInit function to create the service.
			svc := serviceInit()
			// Add the provider to the service.
			svc.FieldTable["provider"] = object.Field{
				Ftype:  types.ClassNameSecurityProvider,
				Fvalue: provider,
			}
			// Add the service to the provider.
			securityProviderPutService([]any{provider, svc})
		}
	}

	return provider
}

/*
SecurityProviderServices is a map of maps of functions that return Provider$Service objects.
The outer map is keyed by the service type (KeyPairGenerator, MessageDigest, etc.).
The inner map is keyed by the algorithm name.
*/
var SecurityProviderServices = map[string]map[string]func() *object.Object{
	"KeyPairGenerator": {
		"DiffieHellman": func() *object.Object {
			return NewGoRuntimeService("KeyPairGenerator", "DiffieHellman", types.ClassNameKeyPairGenerator)
		},
		"DSA": func() *object.Object {
			return NewGoRuntimeService("KeyPairGenerator", "DSA", types.ClassNameKeyPairGenerator)
		},
		"RSA": func() *object.Object {
			return NewGoRuntimeService("KeyPairGenerator", "RSA", types.ClassNameKeyPairGenerator)
		},
		"RSASSA-PSS": func() *object.Object {
			return NewGoRuntimeService("KeyPairGenerator", "RSASSA-PSS", types.ClassNameKeyPairGenerator)
		},
		"EC": func() *object.Object {
			return NewGoRuntimeService("KeyPairGenerator", "EC", types.ClassNameKeyPairGenerator)
		},
	},
	"MessageDigest": {
		"MD5": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "MD5", types.ClassNameMessageDigest)
		},
		"SHA-1": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-1", types.ClassNameMessageDigest)
		},
		"SHA-224": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-224", types.ClassNameMessageDigest)
		},
		"SHA-256": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-256", types.ClassNameMessageDigest)
		},
		"SHA-384": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-384", types.ClassNameMessageDigest)
		},
		"SHA-512": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-512", types.ClassNameMessageDigest)
		},
		"SHA-512/224": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-512/224", types.ClassNameMessageDigest)
		},
		"SHA-512/256": func() *object.Object {
			return NewGoRuntimeService("MessageDigest", "SHA-512/256", types.ClassNameMessageDigest)
		},
	},
}

// NewGoRuntimeService creates a basic Provider$Service object for a security runtime algorithm.
func NewGoRuntimeService(typ, algo, className string) *object.Object {
	svc := object.MakeEmptyObjectWithClassName(&types.ClassNameSecurityProviderService)

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
		Fvalue: object.StringObjectFromGoString(className),
	}
	svc.FieldTable["aliases"] = object.Field{
		Ftype:  types.StringArrayClassName,
		Fvalue: []*object.Object{},
	}

	// Initialize the service attribute ttributes map.
	attributes := map[string]*object.Object{}

	// Block size.
	blockSize := getBlockSizeForAlgorithm(algo)
	if blockSize > 0 {
		attributes["blockSize"] = object.StringObjectFromGoString(fmt.Sprintf("%d", blockSize))
	}

	// Software or hardware?
	attributes["ImplementedIn"] = object.StringObjectFromGoString("Software")

	// Add the attributes map to the services object field table.
	svc.FieldTable["attributes"] = object.Field{
		Ftype:  types.Map,
		Fvalue: attributes,
	}

	return svc
}

// getBlockSizeForAlgorithm returns the block size for a given algorithm, or 0 if not applicable.
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
