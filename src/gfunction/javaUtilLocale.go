/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"jacobin/stringPool"
	"jacobin/types"
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

// "java/util/Locale.<init>(Ljava/lang/String;)V"
func localeFromLanguage(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	inObj := params[1].(*object.Object)
	outObj := params[0].(*object.Object)
	outObj.FieldTable["value"] = inObj.FieldTable["value"]
	return nil
}

// "java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;)V"
func localeFromLanguageCountry(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	// params[2]: input country string
	langObj := params[1].(*object.Object) // string
	langStr := object.GoStringFromStringObject(langObj)

	countryObj := params[2].(*object.Object) // string
	countryStr := object.GoStringFromStringObject(countryObj)

	bytes := []byte(langStr + "_" + countryStr)
	object.UpdateStringObjectFromBytes(params[0].(*object.Object), bytes)

	return nil
}

// "java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"
func localeFromLanguageCountryVariant(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	// params[2]: input country string
	// params[3]: input variant string
	langObj := params[1].(*object.Object)
	langStr := object.GoStringFromStringObject(langObj)

	countryObj := params[2].(*object.Object)
	countryStr := object.GoStringFromStringObject(countryObj)

	variantObj := params[3].(*object.Object)
	variantStr := object.GoStringFromStringObject(variantObj)

	bytes := []byte(langStr + "_" + countryStr + "_" + variantStr)
	object.UpdateStringObjectFromBytes(params[0].(*object.Object), bytes)

	return nil
}

// "java/util/Locale.getDefault()Ljava/util/Locale;"
func getDefaultLocale([]interface{}) interface{} {
	langStr := os.Getenv("LANGUAGE")
	classStr := "java/lang/Locale"
	index := stringPool.GetStringIndex(&langStr)
	obj := object.MakeEmptyObjectWithClassName(&classStr)
	fld := object.Field{Ftype: types.ByteArray, Fvalue: index}
	obj.FieldTable["value"] = fld
	return obj
}
