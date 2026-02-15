/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

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
