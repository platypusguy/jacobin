package javaSecurity

import "jacobin/src/gfunction/ghelpers"

// Loader
func Load_Security_Spec_AlgorithmParameterSpec() {

	ghelpers.MethodSignatures["java/security/spec/AlgorithmParameterSpec.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}
}
