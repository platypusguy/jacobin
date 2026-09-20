/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
)

// Load_Net_SocketAddress registers trap entries for java.net.SocketAddress.
// None of this class is currently implemented, so every method routes to
// TrapFunction.
func Load_Net_SocketAddress() {

	ghelpers.MethodSignatures["java/net/SocketAddress.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/net/SocketAddress.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
