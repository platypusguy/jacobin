/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/object"
	"jacobin/src/types"
)

// In our implementation, Runnable is a set of three strings: class name,
// method name and signature--by which we can find the method.

func newRunnable(clName string, methName string, signature string) *object.Object {
	runnableClassName := "java/lang/Runnable"
	o := object.MakeEmptyObjectWithClassName(&runnableClassName)
	o.FieldTable["clName"] = object.Field{Ftype: types.GolangString, Fvalue: clName}
	o.FieldTable["methName"] = object.Field{Ftype: types.GolangString, Fvalue: methName}
	o.FieldTable["signature"] = object.Field{Ftype: types.GolangString, Fvalue: signature}
	return o
}
