/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package classloader

import (
	"jacobin/exceptions"
	"jacobin/libs"
	"jacobin/log"
	"jacobin/types"
)

// IMPORTANT NOTE: Some String functions are placed in libs\javaLangStringMethods.go
// due to golang circularity concerns, alas.

/*
   We don't run String's static initializer block because the initialization
   is already handled in String creation
*/

func Load_Lang_String() map[string]GMeth {
	// need to replace eventually by enbling the Java intializer to run
	MethodSignatures["java/lang/String.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringClinit,
		}

	// get the bytes from a string
	MethodSignatures["java/lang/String.getBytes()[B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  libs.GetBytesVoid,
		}
	return MethodSignatures
}

func stringClinit([]interface{}) interface{} {
	klass := MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In <clinit>, expected java/lang/String to be in the MethodArea, but it was not"
		_ = log.Log(errMsg, log.SEVERE)
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}

// // get the bytes of a string. To find the string involved, we go to the TOS of the calling
// // stack which has pushed a pointer to the string prior to this call.
// func getBytesVoid(params []interface{}) interface{} {
// 	threadPtr := params[0].(*jvmThread.ExecThread)
// 	frameStack := threadPtr.Stack
// 	prevFrame := frameStack.Front().Next().Value.(*frames.Frame)
// 	str := prevFrame.OpStack[prevFrame.TOS].(*object.Object)
// 	bytes := str.Fields[0].Fvalue.([]byte)
// 	return bytes
// }
