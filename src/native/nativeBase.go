/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import "errors"

var CaughtNativeFunctionException = errors.New("caught native function exception")

// Native function to library handle table
// Input: Native function name.
// Output: Library file handle.
var nfToLibTable = map[string]uintptr{}

// Native function to template function handle table
// Input: Native function name.
// Output: Template function handle.
var nfToTmplTable = map[string]typeTemplateFunction{}

// Native function error block.
type NativeErrBlk struct {
	ExceptionType int
	ErrMsg        string
}

// Type definition for all the template functions
type typeTemplateFunction func(libHandle uintptr, functionName string, params []interface{}) interface{}

// Argument types for template functions.
type NFboolean uint8
type NFbyte uint8
type NFchar uint16
type NFshort int16
type NFint int32
type NFlong int64
type NFfloat float32
type NFdouble float64
