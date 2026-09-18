/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/classloader"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/statics"
	"jacobin/src/trace"
	"jacobin/src/types"
)

func Load_Lang_Void() {

	ghelpers.MethodSignatures["java/lang/Void.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  voidClinit,
		}
}

var classNameVoid = "java/lang/Void"

// voidClinit initializes the static fields of java.lang.Void.
// Specifically, it sets the TYPE field to the primitive class for "void".
func voidClinit(_ []interface{}) interface{} {
	// Fetch the dummy "void" class from the Method Area
	k := classloader.MethAreaFetch("void")
	if k == nil || k.Data.ClassObject == nil {
		// Fatal error: boot sequence failed
		trace.Error("voidClinit: primitive 'void' class not found in MethArea")
		return nil
	}

	// Set the static field Void.TYPE to this object
	statics.AddStatic("java/lang/Void.TYPE", statics.Static{
		Type:  types.Ref,
		Value: k.Data.ClassObject,
	})

	return nil

}
