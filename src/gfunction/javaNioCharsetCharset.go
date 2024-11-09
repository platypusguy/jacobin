/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

// Implementation of some of the functions in Java/nio/charset/Charset.

func Load_Nio_Charset_Charset() {

	// Class initialisation for Console.
	MethodSignatures["java/nio/charset/Charset.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Get the default character set.
	MethodSignatures["java/nio/charset/Charset.defaultCharset()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

}
