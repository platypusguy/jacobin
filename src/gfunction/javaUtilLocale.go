/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"os"
)

// Implementation of some of the functions in Java/util/Locale.
// Strategy: Locale = jacobin Object wrapping a Go string.

func Load_Util_Locale() map[string]GMeth {

	MethodSignatures["java/util/Locale.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  localeFromLanguage,
		}

	MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  localeFromLanguageCountry,
		}

	MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  localeFromLanguageCountryVariant,
		}

	MethodSignatures["java/util/Locale.getDefault()Ljava/util/Locale;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getDefaultLocale,
		}

	return MethodSignatures
}

func localeFromLanguage(params []interface{}) interface{} {
	// params[0]: input string
	propObj := params[0].(*object.Object) // string
	bytes := propObj.FieldTable["value"].Fvalue.([]byte)
	str := string(bytes)
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func localeFromLanguageCountry(params []interface{}) interface{} {
	// params[0]: input string
	propObj := params[0].(*object.Object) // string
	bytes := propObj.FieldTable["value"].Fvalue.([]byte)
	str1 := string(bytes)

	propObj = params[1].(*object.Object) // string
	bytes = propObj.FieldTable["value"].Fvalue.([]byte)
	str2 := string(bytes)

	str := str1 + "_" + str2
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func localeFromLanguageCountryVariant(params []interface{}) interface{} {
	// params[0]: input string
	propObj := params[0].(*object.Object)
	bytes := propObj.FieldTable["value"].Fvalue.([]byte)
	str1 := string(bytes)

	propObj = params[1].(*object.Object)
	bytes = propObj.FieldTable["value"].Fvalue.([]byte)
	str2 := string(bytes)

	propObj = params[2].(*object.Object)
	bytes = propObj.FieldTable["value"].Fvalue.([]byte)
	str3 := string(bytes)

	str := str1 + "_" + str2 + "_" + str3
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func getDefaultLocale(params []interface{}) interface{} {
	str := os.Getenv("LANGUAGE")
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}
