/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
)

func Load_Traps() {

	MethodSignatures["java/rmi/RMISecurityManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/rmi/RMISecurityManager.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/sql/Driver.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/sql/DriverAction.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/sql/DriverPropertyInfo.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/sql/DriverPropertyInfo.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapClass,
		}

	MethodSignatures["java/sql/DriverManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/util/Iterator.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

}

// Generic trap for classes
func trapClass([]interface{}) interface{} {
	errMsg := "TRAP: The requested class is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for deprecated classes and functions
func trapDeprecated([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is deprecated and, therefore, not supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for deprecated classes and functions
func trapUndocumented([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is undocumented and, therefore, not supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for functions
func trapFunction([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for functions
func trapProtected([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is protected"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}
