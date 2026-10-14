/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Java_Lang_Exception() {

	ghelpers.MethodSignatures["java/lang/Exception.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    ThrowableInitNull,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/Exception.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    throwableInitString,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/Exception.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Exception.<init>(Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Exception.<init>(Ljava/lang/String;Ljava/lang/Throwable;ZZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

}
