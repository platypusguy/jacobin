package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
)

// Load_PublicPrivateKeys registers minimal interface classes
func Load_PublicAndPrivateKeys() {
	// ---------------------------------------------------------
	// PublicKey interface placeholder
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/PublicKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	// ---------------------------------------------------------
	// PrivateKey interface placeholder
	// ---------------------------------------------------------
	ghelpers.MethodSignatures["java/security/PrivateKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}
}
