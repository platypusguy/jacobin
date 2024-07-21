/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_Object() {

	MethodSignatures["java/lang/Object.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Object.getClass()Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectGetClass,
		}

	MethodSignatures["java/lang/Object.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectToString,
		}

}

// "java/lang/Object.getClass()Ljava/lang/Class;"
func objectGetClass(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	wint := obj.KlassName
	name := object.GoStringFromStringPoolIndex(wint)
	return object.StringObjectFromGoString("class " + name)
}

// "java/lang/Object.toString()Ljava/lang/String;"
func objectToString(params []interface{}) interface{} {
	// params[0]: input Object
	var str string

	switch params[0].(type) {
	case *object.Object:
		inObj := params[0].(*object.Object)
		str = object.ObjectFieldToString(inObj, "value")
		return object.StringObjectFromGoString(str)
	}

	errMsg := fmt.Sprintf("Unsupported parameter type: %T", params[0])
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)
}
