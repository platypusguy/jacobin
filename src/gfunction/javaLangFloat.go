/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"math"
	"unsafe"
)

func Load_Lang_Float() {

	MethodSignatures["java/lang/Float.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	// Native functions or caller to native functions

	MethodSignatures["java/lang/Float.floatToIntBits(F)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatToIntBits,
		}

	MethodSignatures["java/lang/Float.floatToRawIntBits(F)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  floatToRawIntBits,
		}

	MethodSignatures["java/lang/Float.intBitsToFloat(I)F"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  intBitsToFloat,
		}

}

// Simulating intBitsToFloat in Go
// "java/lang/Float.intBitsToFloat(I)F"
func intBitsToFloat(params []interface{}) interface{} {
	bits := params[0].(int64)
	return math.Float64frombits(uint64(bits))
}

// Simulating floatToRawIntBits in Go
// "java/lang/Float.floatToRawIntBits(F)I"
func floatToRawIntBits(params []interface{}) interface{} {
	value := params[0].(float64)
	return *(*int64)(unsafe.Pointer(&value))
}

// Simulating floatToIntBits in Go
// "java/lang/Float.floatToIntBits(F)I"
func floatToIntBits(params []interface{}) interface{} {
	value := params[0].(float64)
	if !math.IsNaN(float64(value)) {
		return *(*int64)(unsafe.Pointer(&value))
	}
	return 0x7fc00000 // equivalent to Java's 0x7fc00000
}
