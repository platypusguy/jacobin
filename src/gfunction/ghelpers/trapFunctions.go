/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"jacobin/src/excNames"
)

// TrapClass is a generic Trap for classes
func TrapClass([]interface{}) interface{} {
	errMsg := "TRAP: The requested class is not yet supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapDeprecated is a generic Trap for deprecated classes and functions
func TrapDeprecated([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is deprecated and, therefore, not supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for functions
func TrapFunction([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is not yet supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func TrapKeyPairGeneration([]interface{}) interface{} {
	errMsg := "TRAP: Use KeyPairGenerator to create public and private keys"
	return GetGErrBlk(excNames.SecurityException, errMsg)
}

// TrapProtected is a generic Trap for functions
func TrapProtected([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is protected"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapUndocumented is a generic Trap for deprecated classes and functions
func TrapUndocumented([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is undocumented and, therefore, not supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func TrapUnicode([]interface{}) interface{} {
	return GetGErrBlk(
		excNames.UnsupportedOperationException,
		"Character Unicode method not yet implemented",
	)
}
