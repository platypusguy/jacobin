/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

// Implementation of some of the functions in Java/util/HexFormat.

var JavaUtilHexFormat string = "java/util/HexFormat"

func Load_Util_HexFormat() map[string]GMeth {

	MethodSignatures["java/util/HexFormat.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	return MethodSignatures
}
