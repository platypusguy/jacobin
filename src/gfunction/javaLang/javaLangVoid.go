/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/classloader"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/statics"
	"jacobin/src/trace"
	"jacobin/src/types"
	"slices"
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
	// Create the primitive class object for "void"
	primJlc := classloader.MakeJlcEntry("void", true)

	// Register it in the JLCmap so it can be found by name "void"
	classloader.JlcMapLock.Lock()
	classloader.JLCmap["void"] = primJlc
	classloader.JlcMapLock.Unlock()

	// Set the static field Void.TYPE to this object
	_ = statics.AddStatic("java/lang/Void.TYPE", statics.Static{
		Type:  types.Jlc,
		Value: primJlc,
	})

	// Also update the Jlc entry for Void to include this static field in its Statics list
	classloader.JlcMapLock.RLock()
	voidJlc, ok := classloader.JLCmap[classNameVoid]
	classloader.JlcMapLock.RUnlock()

	if ok {
		fieldName := "TYPE"
		fieldDesc := types.Jlc
		entry := fieldName + fieldDesc

		found := false
		voidJlc.Lock.Lock()
		if slices.Contains(voidJlc.Statics, entry) {
			found = true
		}
		if !found {
			voidJlc.Statics = append(voidJlc.Statics, entry)
		}
		voidJlc.Lock.Unlock()
	} else {
		// This should not happen if LoadBaseClasses ran and loaded Integer
		if globals.TraceClass {
			trace.Warning("voidClinit: java/lang/Void not found in JLCmap")
		}
	}
	return nil
}
