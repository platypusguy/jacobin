/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
)

// Load_Util_Function_Supplier loads the method signatures for the
// java/util/function/Supplier interface. Supplier is a functional
// interface whose sole abstract method is get(), which supplies a
// result with no input.
func Load_Util_Function_Supplier() {

	ghelpers.MethodSignatures["java/util/function/Supplier.get()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
