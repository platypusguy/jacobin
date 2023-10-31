/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin Authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package libs

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
)

// These are java/lang/String native functions that need to be segregated in libs
// due to golang circularity issues. Normally, they'd be part of javaLangString.go
// in classloader (which, it should be noted, is only in classloader package due to
// circularity issues, alas).

// get the bytes of a string. To find the string involved, we go to the TOS of the calling
// stack which has pushed a pointer to the string prior to this call. Returns a pointer to
// a raw slice of bytes.

// NOTE:
// * params[0] = extra parameter, the String object
func GetBytesVoid(params []interface{}) interface{} {
	switch params[0].(type) {
	case *object.Object:
		parmObj := params[0].(*object.Object)
		bytes := parmObj.Fields[0].Fvalue.(*[]byte)
		return bytes
	default:
		errMsg := fmt.Sprintf("In libs.GetBytesVoid, unexpected params[0] type=%T, value=%v", params[0], params[0])
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
		return nil
	}
}
