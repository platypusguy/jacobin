/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
)

func Load_Time_DateTimeException() {

	ghelpers.MethodSignatures["java/time/DateTimeException.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    javaLang.ThrowableInitString,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/time/DateTimeException.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots:   2,
			GFunction:    javaLang.ThrowableInitString,
			NeedsContext: true,
		}

}
