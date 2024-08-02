/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Security_AccessController() {

	MethodSignatures["java/security/AccessController.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/security/AccessController.doPrivileged(Ljava/security/PrivilegedAction;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  returnNullObject,
		}

}
