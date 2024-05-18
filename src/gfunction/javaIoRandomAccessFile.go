/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Io_RandomAccessFile() {

	MethodSignatures["java/io/RandomAccessFile.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// ----------------------------------------------------------
	// initIDs justReturn
	// These are private functiona that calls C native functions.
	// ----------------------------------------------------------

	MethodSignatures["java/io/RandomAccessFile.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

}
