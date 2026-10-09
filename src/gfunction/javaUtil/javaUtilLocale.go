/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"strings"
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
			GFunction:  localeGetDefaultLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDefaultLocale, // ignore input
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localeGetDefaultLocale, // ignore input
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  localeGetDefaultLocale, // ignore input
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Lsun/util/locale/BaseLocale;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localeGetDefaultLocale, // ignore input
		}

	// traps for not-yet-implemented Locale instance/static methods.

	ghelpers.MethodSignatures["java/util/Locale.setDefault(Ljava/util/Locale;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.setDefault(Ljava/util/Locale$Category;Ljava/util/Locale;)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getAvailableLocales()[Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISOCountries()[Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISOCountries(Ljava/util/Locale$IsoCountryCode;)Ljava/util/Set;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.getISOLanguages()[Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.forLanguageTag(Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: localeOf}
	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: localeOf}
	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] = ghelpers.GMeth{ParamSlots: 3, GFunction: localeOf}

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
	ghelpers.MethodSignatures["java/util/Locale.toString()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: LocaleToString}

}

func _getLocaleFromEnv() string {
	for _, name := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := os.Getenv(name); value != "" {
			// strip ".encoding" and "@modifier", e.g. de_DE.UTF-8@euro -> de_DE
			if i := strings.IndexAny(value, ".@"); i >= 0 {
				value = value[:i]
			}
			if value != "" && value != "C" && value != "POSIX" {
				return value
			}
		}
	}
	return "en_US"
}

// "java/util/Locale.getDefault()Ljava/util/Locale;"
// "java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"
// "java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"
func localeGetDefaultLocale([]interface{}) interface{} {
	// Ignore parameters.
	var langStr string
	if globals.OnWindows {
		langStr = _getDefaultLanguageFromWindows()
	} else {
		langStr = _getLocaleFromEnv()
	}
	obj := object.MakePrimitiveObject(types.ClassNameLocale, types.StringClassName, object.StringObjectFromGoString(langStr))
	return obj
}

// of(String language)
// of(String language, String country)
// of(String language, String country, String variant)
func localeOf(params []interface{}) interface{} {
	if len(params) < 1 || len(params) > 3 {
		// Cannot happen if the MethodSignatures table is right; trap it as your
		// other gfunctions do (IllegalArgumentException / internal error).
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("localeOf: unexpected parameter count %d", len(params)))
	}

	// Locale.of throws NullPointerException for a null argument.
	parts := make([]string, 3) // language, country, variant
	for i, p := range params {
		strObj, ok := p.(*object.Object)
		if !ok || object.IsNull(strObj) {
			return ghelpers.GetGErrBlk(excNames.NullPointerException,
				fmt.Sprintf("localeOf: argument %d is null", i+1))
		}
		parts[i] = object.GoStringFromStringObject(strObj)
	}

	language := strings.ToLower(parts[0])
	country := strings.ToUpper(parts[1])
	variant := parts[2] // Java leaves the variant's case alone

	// Same text Locale.toString() produces:
	//   en | en_US | en__POSIX | en_US_POSIX | _US
	localeData := language
	if country != "" || variant != "" {
		localeData += "_" + country
	}
	if variant != "" {
		localeData += "_" + variant
	}
	valueObj := object.StringObjectFromGoString(localeData)
	localeObj := object.MakeEmptyObjectWithClassName(&types.ClassNameLocale)
	fld := object.Field{Ftype: types.Ref, Fvalue: valueObj}
	localeObj.FieldTable["value"] = fld

	return localeObj
}

func LocaleToString(params []interface{}) interface{} {
	objectRef, ok := params[0].(*object.Object)
	if !ok || object.IsNull(objectRef) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException,
			"localeToString: objectRef is null")
	}

	valueObj, ok := objectRef.FieldTable["value"].Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			"localeToString: Locale object has no value field or value field is not a string")
	}

	return valueObj
}
