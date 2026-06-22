/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Util_SequencedSet() {

	ghelpers.MethodSignatures["java/util/SequencedSet.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/SequencedSet.reversed()Ljava/util/SequencedSet;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
