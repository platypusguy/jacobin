/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Io_File_Native_Fakes() map[string]GMeth {

	MethodSignatures["java/io/UnixFileSystem.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/WinNTFileSystem.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	return MethodSignatures
}
