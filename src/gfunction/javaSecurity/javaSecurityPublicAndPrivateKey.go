/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

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
