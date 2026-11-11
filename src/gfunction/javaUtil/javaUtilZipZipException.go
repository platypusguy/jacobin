/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
)

func Load_Util_Zip_ZipException() {

	ghelpers.MethodSignatures["java/util/zip/ZipException.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    javaLang.ThrowableInitNull,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/util/zip/ZipException.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    javaLang.ThrowableInitString,
			NeedsContext: true,
		}

}
