/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

/*
We don't run String's static initializer block because the initialization
is already handled in new String creation
*/

func Load_Lang_String() map[string]GMeth {
	// need to replace eventually by enbling the Java intializer to run
	MethodSignatures["java/lang/String.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	return MethodSignatures
}
