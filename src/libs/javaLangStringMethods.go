/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin Authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package libs

import (
	"jacobin/frames"
	"jacobin/object"
	jvmThread "jacobin/thread"
)

// These are java/lang/String native functions that need to be segregated in libs
// due to golang circularity issues. Normally, they'd be part of javaLangString.go
// in classloader (which, I should add, is only in classloader package due to
// circularity issues, alas.

// get the bytes of a string. To find the string involved, we go to the TOS of the calling
// stack which has pushed a pointer to the string prior to this call.
func GetBytesVoid(params []interface{}) interface{} {
	threadPtr := params[0].(*jvmThread.ExecThread)
	frameStack := threadPtr.Stack
	prevFrame := frameStack.Front().Next().Value.(*frames.Frame)
	str := prevFrame.OpStack[prevFrame.TOS].(*object.Object)
	bytes := str.Fields[0].Fvalue.([]byte)
	return bytes
}
