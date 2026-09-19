/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunMisc

import "jacobin/src/gfunction/ghelpers"

func Load_Traps_Sun_Misc() {

	ghelpers.MethodSignatures["sun/nio/ch/SocketDispatcher.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

}
