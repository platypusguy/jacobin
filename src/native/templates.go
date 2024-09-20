/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

/*

During initialization,
* The NfLibXrefTable is built by either a POSIX loader or a Windows loader. Note that both the library path and handle are populated.
* The nfToTmplTable remains nil.

At run-time, RunNativeFunction will do the following in order to get (1) a native function handle
and (2) the corresponding template function address:
* Look up the methodName in the nfToTmplTable.
* If not found,
     - Look up methodName in nfToLibTable. Not found ---> error.
     - Derive the template function to use for this methodName based on the methodType.
     - Store the template address in nfToTmplTable.
* Call the template function (by address) with arguments: library handle and the function name.

*/

package native

import (
	"github.com/ebitengine/purego"
)

/*
Map a method type to a template function handle.
Input:

	Java language method type string E.g. (II)I

Return:
  - Reference to a template function if successful; else nil
  - Ok-boolean = true if successful; else false

Methodology:

	Brute-force switch (gasp!)
*/
func mapToTemplateHandle(methodType string) (typeTemplateFunction, bool) {
	switch methodType {
	// to add a new function, add the template here as a new case statement
	case "(II)I":
		var templateFunction = template_II_I
		return templateFunction, true
	}
	return nil, false
}

/*
Template function for method type (II)I
E.g. Java_java_util_zip_CRC32_update
*/
func template_II_I(libHandle uintptr, nativeFunctionName string, params []interface{}, tracing bool) interface{} {
	// Register the native function prototype
	// env = JNI environment
	// class = a reference to the object whose method is being called
	// arg1, arg2 are the II that we're passing in
	var fn func(env, class uintptr, arg1, arg2 NFint) NFint
	purego.RegisterLibFunc(&fn, libHandle, nativeFunctionName)

	// Get arguments.
	arg1 := NFint(params[0].(int64))
	arg2 := NFint(params[1].(int64))

	// Compute result and return it.
	// out := fn(HandleENV, 0, arg1, arg2)
	out := fn(0, 0, arg1, arg2)
	return int64(out)
}
