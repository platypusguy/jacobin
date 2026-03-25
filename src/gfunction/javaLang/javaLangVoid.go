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

	/*
		// Create the primitive class object for "void"
		primJlc := classloader.MakeJlcEntry("void", true)

		// Register it in the JLCmap so it can be found by name "void"
		classloader.JlcMapLock.Lock()
		classloader.JLCmap["void"] = primJlc
		classloader.JlcMapLock.Unlock()

		// Set the static field Void.TYPE to this object
		_ = statics.AddStatic("java/lang/Void.TYPE", statics.Static{
			Type:  types.Ref,
			Value: object.MakePrimitiveObjectFromJlcInstance("void"),
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

	*/
}
