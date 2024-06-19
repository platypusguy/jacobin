/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Security_SecureRandom() {

	MethodSignatures["java/security/SecureRandom.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/security/SecureRandom.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secureRandomInit,
		}

	MethodSignatures["java/security/SecureRandom.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  secureRandomInit,
		}

	MethodSignatures["java/security/SecureRandom.<init>(Ljava.security.SecureRandomSpi;Ljava.security.Provider;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

}

// "java/security/SecureRandom.<init>()V"
// "java/security/SecureRandom.<init>([B)V"
func secureRandomInit(params []interface{}) interface{} {
	return nil
}
