/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import "jacobin/src/gfunction/ghelpers"

func Load_Traps_Java_Security() {

	ghelpers.MethodSignatures["java/security/AccessController.doPrivileged(Ljava/security/PrivilegedAction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAMultiPrimePrivateCrtKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/security/interfaces/RSAPrivateCrtKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/security/AccessController.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

}
