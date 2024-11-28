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
)

// jj (Jacobin JVM) functions are functions that can be inserted inside Java programs
// for diagnostic purposes. They simply return when run in the JDK, but do what they're
// supposed to do when run under Jacobin.
//
// Note this is a rough first design that will surely be refined. (JACOBIN-624)

func Load_jj() {

	MethodSignatures["jj._dumpStatics(Ljava/lang/String;ILjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpStatics,
		}

	MethodSignatures["jj._dumpObject(Ljava/lang/Object;Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpObject,
		}
}

func jjDumpStatics(params []interface{}) interface{} {
	fromObj := params[0].(*object.Object)
	if fromObj == nil || fromObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("Invalid object in objectGetClass(): %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	from := object.ObjectFieldToString(fromObj, "value")
	selection := params[1].(int64)
	classNameObj := params[2].(*object.Object)
	className := object.ObjectFieldToString(classNameObj, "value")

	statics.DumpStatics(from, selection, className)
	return nil
}

func jjDumpObject(params []interface{}) interface{} {
	this := params[0].(*object.Object)
	objTitle := params[1].(*object.Object)
	title := object.ObjectFieldToString(objTitle, "value")
	indent := params[2].(int64)
	this.DumpObject(title, int(indent))
	return nil
}
