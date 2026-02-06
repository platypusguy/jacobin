/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

func Load_Traps_Java_Security() {

	MethodSignatures["java/security/interfaces/RSAMultiPrimePrivateCrtKey.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/security/interfaces/RSAPrivateCrtKey.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

}
