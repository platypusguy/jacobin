package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Load_Security initializes java/security/Security methods
func Load_Security() {

	ghelpers.MethodSignatures["java/security/Security.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// ---------- Static Methods ----------

	// Security.getProvider(String)
	ghelpers.MethodSignatures["java/security/Security.getProvider(Ljava/lang/String;)Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 1, // provider name
			GFunction:  securityGetProvider,
		}

	// Security.getProviders()
	ghelpers.MethodSignatures["java/security/Security.getProviders()[Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  securityGetProviders,
		}

	// Security.addProvider(Provider) - ghelpers.TrapFunction, since we only allow DefaultSecurityProvider
	ghelpers.MethodSignatures["java/security/Security.addProvider(Ljava/security/Provider;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// Security.insertProviderAt(Provider, int) - ghelpers.TrapFunction
	ghelpers.MethodSignatures["java/security/Security.insertProviderAt(Ljava/security/Provider;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Security.removeProvider(String) - ghelpers.TrapFunction
	ghelpers.MethodSignatures["java/security/Security.removeProvider(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

}

// ----------------------- Member Functions -----------------------

func securityGetProvider([]any) any {
	return ghelpers.GetDefaultSecurityProvider()
}

// getProviders() -> Provider[]
func securityGetProviders(params []any) any {
	provider := ghelpers.GetDefaultSecurityProvider()
	if provider == nil {
		return []*object.Object{}
	}
	return object.MakeOneFieldObject(types.ObjectClassName, "value", types.RefArray, []*object.Object{provider})
}
