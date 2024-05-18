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

func Load_Lang_Long() {

	MethodSignatures["java/lang/Long.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Long.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  longDoubleValue,
		}

	MethodSignatures["java/lang/Long.valueOf(J)Ljava/lang/Long;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  longValueOf,
		}

}

// "java/lang/Long.doubleValue()D"
func longDoubleValue(params []interface{}) interface{} {
	var jj int64
	parmObj := params[0].(*object.Object)
	jj = parmObj.FieldTable["value"].Fvalue.(int64)
	return float64(jj)
}

// "java/lang/Long.valueOf(J)Ljava/lang/Long;"
func longValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return populator("java/lang/Long", types.Long, int64Value)
}
