/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"bytes"
	"jacobin/object"
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
	fld := inObj.FieldTable["value"]
	bites := bytes.ToLower(fld.Fvalue.([]byte))
	fld.Fvalue = bites
	params[0].(*object.Object).FieldTable["value"] = fld
	return nil
}

func localeFromLanguageCountry(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	// params[2]: input country string
	inObj := params[1].(*object.Object) // string
	bites := inObj.FieldTable["value"].Fvalue.([]byte)

	inObj = params[2].(*object.Object) // string
	bytesCountry := inObj.FieldTable["value"].Fvalue.([]byte)
	bites = append(bites, '_')
	bites = append(bites, bytesCountry...)

	bites = bytes.ToLower(bites)
	fld := object.Field{Ftype: types.ByteArray, Fvalue: bites}
	params[0].(*object.Object).FieldTable["value"] = fld
	return nil
}

func localeFromLanguageCountryVariant(params []interface{}) interface{} {
	// params[0]: Locale object to update
	// params[1]: input language string
	// params[2]: input country string
	// params[3]: input variant string
	inObj := params[1].(*object.Object) // string
	bites := inObj.FieldTable["value"].Fvalue.([]byte)

	inObj = params[2].(*object.Object) // string
	bytesCountry := inObj.FieldTable["value"].Fvalue.([]byte)
	bites = append(bites, '_')
	bites = append(bites, bytesCountry...)

	inObj = params[3].(*object.Object) // string
	bytesVariant := inObj.FieldTable["value"].Fvalue.([]byte)
	bites = append(bites, '_')
	bites = append(bites, bytesVariant...)

	bites = bytes.ToLower(bites)
	fld := object.Field{Ftype: types.ByteArray, Fvalue: bites}
	params[0].(*object.Object).FieldTable["value"] = fld
	return nil
}

func getDefaultLocale(params []interface{}) interface{} {
	str := os.Getenv("LANGUAGE")
	obj := object.MakePrimitiveObject("java/util/Locale", types.ByteArray, []byte(str))
	return obj
}
