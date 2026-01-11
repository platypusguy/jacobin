/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

// Implementation of some of the functions in Java/util/Locale.
// Strategy: Locale = jacobin Object wrapping a Go string.

func Load_Util_Locale() {

	ghelpers.MethodSignatures["java/util/Locale.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDefault()Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  getDefaultLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  getDefaultLocale, // ignore input
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  getDefaultLocale, // ignore input
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  getDefaultLocale, // ignore input
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Lsun/util/locale/BaseLocale;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
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
