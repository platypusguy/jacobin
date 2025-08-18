/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

// Implementation of some of the functions in Java/util/Locale.
// Strategy: Locale = jacobin Object wrapping a Go string.

func Load_Util_Locale() {

	MethodSignatures["java/util/Locale.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/util/Locale.getDefault()Ljava/util/Locale;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getDefaultLocale,
		}

	MethodSignatures["java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  getDefaultLocale, // ignore input
		}

	MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  getDefaultLocale, // ignore input
		}

	MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		GMeth{
			ParamSlots: 5,
			GFunction:  getDefaultLocale, // ignore input
		}

	MethodSignatures["java/util/Locale.getInstance(Lsun/util/locale/BaseLocale;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  getDefaultLocale, // ignore input
		}

}

// "java/util/Locale.getDefault()Ljava/util/Locale;"
// "java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"
// "java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"
func getDefaultLocale([]interface{}) interface{} {
	// Ignore parameters.
	langStr := os.Getenv("LANGUAGE")
	classStr := "java/lang/Locale"
	obj := object.MakeEmptyObjectWithClassName(&classStr)
	fld := object.Field{Ftype: types.ByteArray, Fvalue: []byte(langStr)}
	obj.FieldTable["value"] = fld
	return obj
}
