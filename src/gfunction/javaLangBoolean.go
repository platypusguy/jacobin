/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"jacobin/types"
)

func Load_Lang_Boolean() map[string]GMeth {

	MethodSignatures["java/lang/Boolean.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  booleanValueOf,
		}

	return MethodSignatures
}

// java/lang/Boolean.valueOf(Z)Ljava/lang/Boolean;
func booleanValueOf(params []interface{}) interface{} {
	zz := params[0].(int64)
	objPtr := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, zz)
	return objPtr
}
