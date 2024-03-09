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

func localeFromLanguage(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	inObj := params[1].(*object.Object)
	outObj := params[0].(*object.Object)
	outObj.FieldTable["value"] = inObj.FieldTable["value"]
	return nil
}

func localeFromLanguageCountry(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	// params[2]: input country string
	langObj := params[1].(*object.Object) // string
	langStr := object.GetGoStringFromObject(langObj)

	countryObj := params[2].(*object.Object) // string
	countryStr := object.GetGoStringFromObject(countryObj)

	str := langStr + "_" + countryStr
	fld := params[0].(*object.Object).FieldTable["value"]
	fld.Ftype = types.StringIndex
	fld.Fvalue = stringPool.GetStringIndex(&str)
	params[0].(*object.Object).FieldTable["value"] = fld

	return nil
}

func localeFromLanguageCountryVariant(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	// params[2]: input country string
	// params[3]: input variant string
	langObj := params[1].(*object.Object)
	langStr := object.GetGoStringFromObject(langObj)

	countryObj := params[2].(*object.Object)
	countryStr := object.GetGoStringFromObject(countryObj)

	variantObj := params[3].(*object.Object)
	variantStr := object.GetGoStringFromObject(variantObj)

	str := langStr + "_" + countryStr + "_" + variantStr
	fld := params[0].(*object.Object).FieldTable["value"]
	fld.Ftype = types.StringIndex
	fld.Fvalue = stringPool.GetStringIndex(&str)
	params[0].(*object.Object).FieldTable["value"] = fld

	return nil
}

func getDefaultLocale([]interface{}) interface{} {
	langStr := os.Getenv("LANGUAGE")
	classStr := "java/lang/Locale"
	index := stringPool.GetStringIndex(&langStr)
	obj := object.MakeEmptyObjectWithClassName(&classStr)
	fld := object.Field{Ftype: types.StringIndex, Fvalue: index}
	obj.FieldTable["value"] = fld
	return obj
}
