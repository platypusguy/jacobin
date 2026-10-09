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

	// traps for not-yet-implemented Locale instance/static methods.

	ghelpers.MethodSignatures["java/util/Locale.setDefault(Ljava/util/Locale;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.setDefault(Ljava/util/Locale$Category;Ljava/util/Locale;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getAvailableLocales()[Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISOCountries()[Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISOCountries(Ljava/util/Locale$IsoCountryCode;)Ljava/util/Set;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISOLanguages()[Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.forLanguageTag(Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/Locale.getLanguage()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getCountry()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getVariant()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getScript()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getExtensionKeys()Ljava/util/Set;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getUnicodeLocaleAttributes()Ljava/util/Set;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getUnicodeLocaleKeys()Ljava/util/Set;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getExtension(C)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getUnicodeLocaleType(Ljava/lang/String;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayLanguage()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayLanguage(Ljava/util/Locale;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayScript()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayScript(Ljava/util/Locale;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayCountry()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayCountry(Ljava/util/Locale;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayVariant()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayVariant(Ljava/util/Locale;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayName()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getDisplayName(Ljava/util/Locale;)Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISO3Language()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISO3Country()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.toLanguageTag()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.toString()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.equals(Ljava/lang/Object;)Z"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.hashCode()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.clone()Ljava/lang/Object;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.hasExtensions()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.stripExtensions()Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

}

// "java/util/Locale.getDefault()Ljava/util/Locale;"
// "java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"
// "java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"
func getDefaultLocale([]interface{}) interface{} {
	// Ignore parameters.
	langStr := os.Getenv("LANGUAGE")
	classStr := "java/lang/Locale"
	obj := object.MakeEmptyObjectWithClassName(&classStr)
	fld := object.Field{Ftype: types.JavaByteArray, Fvalue: []byte(langStr)}
	obj.FieldTable["value"] = fld
	return obj
}
