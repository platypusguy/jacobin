/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Java_Lang_Invoke_LambdaConversionException() {

	ghelpers.MethodSignatures["java/lang/invoke/LambdaConversionException.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    ThrowableInitNull,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaConversionException.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    ThrowableInitString,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaConversionException.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots:   2,
			GFunction:    ghelpers.TrapFunction,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaConversionException.<init>(Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    ThrowableInitThrowable,
			NeedsContext: true,
		}
}
