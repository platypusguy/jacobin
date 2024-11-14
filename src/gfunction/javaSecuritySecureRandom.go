/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Security_SecureRandom() {

	MethodSignatures["java/security/SecureRandom.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	/**
	MethodSignatures["java/security/SecureRandom.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/security/SecureRandom.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/security/SecureRandom.<init>(Ljava.security.SecureRandomSpi;Ljava.security.Provider;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}
	**/

}
