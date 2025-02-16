/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"jacobin/types"
	"unicode"
)

func Load_Lang_Character() {

	MethodSignatures["java/lang/Character.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Character.isDigit(C)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charIsDigit,
		}

	MethodSignatures["java/lang/Character.isLetter(C)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charIsLetter,
		}

	MethodSignatures["java/lang/Character.charValue()C"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  charValue,
		}

	MethodSignatures["java/lang/Character.toLowerCase(C)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charToLowerCase,
		}

	MethodSignatures["java/lang/Character.toUpperCase(C)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  charToUpperCase,
		}

	MethodSignatures["java/lang/Character.valueOf(C)Ljava/lang/Character;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  characterValueOf,
		}

}

// "java/lang/Character.isDigit(C)Z"
func charIsDigit(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsDigit(rune(ii)) {
		return int64(1)
	}
	return int64(0)
}

// "java/lang/Character.isLetter(C)Z"
func charIsLetter(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsLetter(rune(ii)) {
		return int64(1)
	}
	return int64(0)
}

// "java/lang/Character.toLowerCase(C)C"
func charToLowerCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	rr := unicode.ToLower(rune(ii))
	return int64(rr)
}

// "java/lang/Character.toUpperCase(C)C"
func charToUpperCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	rr := unicode.ToUpper(rune(ii))
	return int64(rr)
}

// "java/lang/Character.valueOf(C)Ljava/lang/Character;"
func characterValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Character", types.Char, int64Value)
}

// "java/lang/Character.charValue()C"
func charValue(params []interface{}) interface{} {
	var ch int64
	parmObj := params[0].(*object.Object)
	ch = parmObj.FieldTable["value"].Fvalue.(int64)
	return ch
}
