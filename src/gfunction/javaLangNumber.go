/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import "jacobin/src/excNames"

func Load_Lang_Number() {

	// Class initializer
	MethodSignatures["java/lang/Number.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	// Constructor: Number()
	MethodSignatures["java/lang/Number.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// abstract int intValue()
	MethodSignatures["java/lang/Number.intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  numberIntValue,
		}

	// abstract long longValue()
	MethodSignatures["java/lang/Number.longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  numberLongValue,
		}

	// abstract float floatValue()
	MethodSignatures["java/lang/Number.floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  numberFloatValue,
		}

	// abstract double doubleValue()
	MethodSignatures["java/lang/Number.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  numberDoubleValue,
		}

	// abstract byte byteValue()
	MethodSignatures["java/lang/Number.byteValue()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  numberByteValue,
		}

	// abstract short shortValue()
	MethodSignatures["java/lang/Number.shortValue()S"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  numberShortValue,
		}
}

// ========================
// Skeleton function bodies
// ========================

// int intValue()
func numberIntValue([]interface{}) interface{} {
	return getGErrBlk(excNames.AbstractMethodError, "numberIntValue is abstract")
}

// long longValue()
func numberLongValue([]interface{}) interface{} {
	return getGErrBlk(excNames.AbstractMethodError, "numberLongValue is abstract")
}

// float floatValue()
func numberFloatValue([]interface{}) interface{} {
	return getGErrBlk(excNames.AbstractMethodError, "numberFloatValue is abstract")
}

// double doubleValue()
func numberDoubleValue([]interface{}) interface{} {
	return getGErrBlk(excNames.AbstractMethodError, "numberDoubleValue is abstract")
}

// byte byteValue()
func numberByteValue([]interface{}) interface{} {
	return getGErrBlk(excNames.AbstractMethodError, "numberByteValue is abstract")
}

// short shortValue()
func numberShortValue([]interface{}) interface{} {
	return getGErrBlk(excNames.AbstractMethodError, "numberShortValue is abstract")
}
