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

func Load_Lang_Short() map[string]GMeth {

	MethodSignatures["java/lang/Short.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Short.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  shortDoubleValue,
		}

	MethodSignatures["java/lang/Short.valueOf(S)Ljava/lang/Short;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  shortValueOf,
		}

	return MethodSignatures
}

// "java/lang/Short.doubleValue()D"
func shortDoubleValue(params []interface{}) interface{} {
	var ii int64
	parmObj := params[0].(*object.Object)
	ii = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(ii)
}

// "java/lang/Short.valueOf(S)Ljava/lang/Short;"
func shortValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Short", types.Short, int64Value)
}
