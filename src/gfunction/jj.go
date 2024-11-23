/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
	"os"
)

// jj (Jacobin JVM) functions are functions that can be inserted inside Java programs
// for diagnostic purposes. They simply return when run in the JDK, but do what they're
// supposed to do when run under Jacobin.
//
// Note this is a rough first design that will surely be refined. (JACOBIN-624)

func Load_jj() {

	MethodSignatures["jj._dumpStatics(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  jjDumpStatics,
		}
}

func jjDumpStatics(params []interface{}) interface{} {
	objPtr := params[0].(*object.Object)
	if objPtr == nil || objPtr.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("Invalid object in objectGetClass(): %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	str := object.ObjectFieldToString(objPtr, "value")
	fmt.Fprintf(os.Stderr, "%s: non-JDK statics:\n", str)
	statics.DumpStatics()
	return nil
}
