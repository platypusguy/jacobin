/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/exceptions"
)

// Implementation of some of the functions in Java/nio/charset/Charset.

func Load_Nio_Charset_Charset() map[string]GMeth {

	// Class initialisation for Console.
	MethodSignatures["java/nio/charset/Charset.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// Get the default character set.
	MethodSignatures["java/nio/charset/Charset.defaultCharset()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDefaultCharset,
		}

	return MethodSignatures
}

// Trap "java/nio/charset/Charset.defaultCharset()Ljava/nio/charset/Charset;"
func trapDefaultCharset([]interface{}) interface{} {
	errMsg := "java/nio/charset/Charset.defaultCharset TRAP: not yet supported !!"
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}
