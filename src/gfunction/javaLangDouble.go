/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"strconv"
)

func Load_Lang_Double() map[string]GMeth {

	MethodSignatures["java/lang/Double.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Double.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleByteValue,
		}

	MethodSignatures["java/lang/Double.compare(DD)I"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  doubleCompare,
		}

	MethodSignatures["java/lang/Double.compareTo(Ljava/lang/Double;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleCompareTo,
		}

	MethodSignatures["java/lang/Double.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleDoubleValue,
		}

	MethodSignatures["java/lang/Double.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleEquals,
		}

	MethodSignatures["java/lang/Double.parseDouble(Ljava/lang/String;)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  doubleParseDouble,
		}

	MethodSignatures["java/lang/Double.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  doubleToString,
		}

	return MethodSignatures
}

func doubleByteValue(params []interface{}) interface{} {
	var dd float64
	parmObj := params[0].(*object.Object)
	dd = parmObj.FieldTable["value"].Fvalue.(float64)
	return int64(byte(dd))
}

func doubleCompare(params []interface{}) interface{} {
	dd1 := params[0].(float64)
	dd2 := params[1].(float64)
	if dd1 == dd2 {
		return int64(0)
	}
	if dd1 < dd2 {
		return int64(-1)
	}
	return int(1)
}

func doubleCompareTo(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	dd1 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Get the actual Java Double parameter
	parmObj = params[1].(*object.Object)
	dd2 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Now, its just like doubleCompare.
	if dd1 == dd2 {
		return int64(0)
	}
	if dd1 < dd2 {
		return int64(-1)
	}
	return int64(1)
}

func doubleDoubleValue(params []interface{}) interface{} {
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		return parmObj.FieldTable["value"].Fvalue.(float64)
	}
	return parmObj.FieldTable["value"].Fvalue.(float64)
}

func doubleEquals(params []interface{}) interface{} {
	var dd1, dd2 float64

	// Get the Double object reference
	parmObj := params[0].(*object.Object)
	dd1 = parmObj.FieldTable["value"].Fvalue.(float64)

	// Get the actual Java Object parameter
	parmObj = params[1].(*object.Object)
	if parmObj.Klass == nil {
		return int64(0)
	}
	// fmt.Printf("DEBUG doubleEquals Klass --> %s\n", *parmObj.Klass)
	if *parmObj.Klass != "java/lang/Double" {
		return int64(0)
	}
	dd2 = parmObj.FieldTable["value"].Fvalue.(float64)

	// If equal, return true; else return false.
	// fmt.Printf("DEBUG doubleEquals dd1=%f, dd2=%f\n", dd1, dd2)
	if dd1 == dd2 {
		return int64(1)
	}
	return int64(0)
}

func doubleParseDouble(params []interface{}) interface{} {
	// Extract and validate the string argument.
	parmObj := params[0].(*object.Object)
	strArg := object.GetGoStringFromObject(parmObj)
	if len(strArg) < 1 {
		return getGErrBlk(exceptions.NumberFormatException, "doubleParseDouble: string length < 1")
	}

	// Compute output.
	output, err := strconv.ParseFloat(strArg, 64)
	if err != nil {
		return getGErrBlk(exceptions.NumberFormatException, "doubleParseDouble Error(): "+err.Error())
	}
	return output

}

func doubleToString(params []interface{}) interface{} {
	var dd float64
	parmObj := params[0].(*object.Object)
	dd = parmObj.FieldTable["value"].Fvalue.(float64)
	str := fmt.Sprintf("%f", dd)
	objPtr := object.NewPoolStringFromGoString(str)
	return objPtr
}
