package gfunction

import (
	"jacobin/src/object"
)

// Load_Security initializes java/security/Security methods
func Load_Security() {

	MethodSignatures["java/security/Security.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	// ---------- Static Methods ----------

	// Security.getProvider(String)
	MethodSignatures["java/security/Security.getProvider(Ljava/lang/String;)Ljava/security/Provider;"] =
		GMeth{
			ParamSlots: 1, // provider name
			GFunction:  securityGetProvider,
		}

	// Security.getProviders()
	MethodSignatures["java/security/Security.getProviders()[Ljava/security/Provider;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  securityGetProviders,
		}

	// Security.addProvider(Provider) - trapFunction, since we only allow DefaultSecurityProvider
	MethodSignatures["java/security/Security.addProvider(Ljava/security/Provider;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	// Security.insertProviderAt(Provider, int) - trapFunction
	MethodSignatures["java/security/Security.insertProviderAt(Ljava/security/Provider;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Security.removeProvider(String) - trapFunction
	MethodSignatures["java/security/Security.removeProvider(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}
}

// ----------------------- Member Functions -----------------------

func securityGetProvider([]any) any {
	return GetDefaultSecurityProvider()
}

// getProviders() -> Provider[]
func securityGetProviders(params []any) any {
	provider := GetDefaultSecurityProvider()
	if provider == nil {
		return []*object.Object{}
	}
	return []*object.Object{provider}
}
