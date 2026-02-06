/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
)

// ---------------------------------------------------------
// Loader for RSAPublicKey
// ---------------------------------------------------------
func Load_Security_Interfaces_RSAPublicKey() {

	// Interface constructor placeholder
	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// Interface constructor placeholder
	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.getModulus()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rsaKeyGetModulus,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPublicKey.getParams()Ljava/security/spec/AlgorithmParameterSpec;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnNull,
		}
}
