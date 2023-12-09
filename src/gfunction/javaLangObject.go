/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_Object() map[string]GMeth {

	MethodSignatures["java/lang/Object.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}
	return MethodSignatures
}
