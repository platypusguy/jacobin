/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
)

// Load_Util_Logging_SocketHandler loads the method signatures for the
// java/util/logging/SocketHandler class. None of SocketHandler's methods
// are deprecated, so all entries use ghelpers.TrapFunction.
func Load_Util_Logging_SocketHandler() {

	ghelpers.MethodSignatures["java/util/logging/SocketHandler.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/logging/SocketHandler.publish(Ljava/util/logging/LogRecord;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}
