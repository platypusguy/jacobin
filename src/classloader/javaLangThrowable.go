/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/log"
)

func Load_Lang_Throwable() map[string]GMeth {

	MethodSignatures["java/lang/Throwable.fillInStackTrace()Ljava/lang/Throwable;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fillInStackTrace,
		}
	return MethodSignatures
}

func fillInStackTrace([]interface{}) interface{} {
	glob := globals.GetGlobalRef()
	if glob.JVMframeStack == nil { // if we haven't captured the JVM stack before now, we're hosed.
		_ = log.Log("No stack data available for this error. Incomplete data will be shown.", log.SEVERE)
		return nil
	}
	// CURR: next steps
	// instantiate an []StackTraceElements
	// fill it in with 1 item per frame.
	// This should probably be done at the time the error/exception is detected.
	// This might require that we add the logic to the class parse showing the Java code source line number.
	// JACOBIN-224 refers to this.
	return nil
}
