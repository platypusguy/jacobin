/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Java_Lang_NullPointerException() {

	ghelpers.MethodSignatures["java/lang/NullPointerException.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    ThrowableInitNull,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/NullPointerException.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    throwableInitString,
			NeedsContext: true,
		}

}
