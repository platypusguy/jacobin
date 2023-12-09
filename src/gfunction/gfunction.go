/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import "jacobin/classloader"

// GMeth is the entry in the MTable for Go functions. See MTable comments for details.
// Fu is a go function. All go functions accept a possibly empty slice of interface{} and
// return a possibly nil interface{}
type GMeth struct {
	ParamSlots   int
	GFunction    func([]interface{}) interface{}
	NeedsContext bool // does this method need a pointer to the frame stack? Defaults to false.
}

// MTableLoadNatives loads the Go methods from files that contain them. It does this
// by calling the Load_* function in each of those files to load whatever Go functions
// they make available.
func MTableLoadNatives(MTable *classloader.MT) {

	loadlib(MTable, Load_Io_PrintStream())         // load the java.io.prinstream golang functions
	loadlib(MTable, Load_Lang_Class())             // load the java.lang.Class golang functions
	loadlib(MTable, Load_Lang_Math())              // load the java.lang.Math golang functions
	loadlib(MTable, Load_Lang_Object())            // load the java.lang.Class golang functions
	loadlib(MTable, Load_Misc_Unsafe())            // load the jdk.internal/misc/Unsafe functions
	loadlib(MTable, Load_Lang_String())            // load the java.lang.String golang functions
	loadlib(MTable, Load_Lang_System())            // load the java.lang.System golang functions
	loadlib(MTable, Load_Lang_StackTraceELement()) //  java.lang.StackTraceElement golang functions
	loadlib(MTable, Load_Lang_Thread())            // load the java.lang.Thread golang functions
	loadlib(MTable, Load_Lang_Throwable())         // load the java.lang.Throwable golang functions (errors & exceptions)
	loadlib(MTable, Load_Lang_UTF16())             // load the java.lang.UTF16 golang functions
	loadlib(MTable, Load_Util_HashMap())           // load the java.util.HashMap golang functions
	loadlib(MTable, Load_Util_Locale())            // load the java.util.Locale golang functions
	loadlib(MTable, Load_Primitives())             // load the Java primitives golang functions
}

func loadlib(tbl *classloader.MT, libMeths map[string]GMeth) {
	for key, val := range libMeths {
		gme := GMeth{}
		gme.ParamSlots = val.ParamSlots
		gme.GFunction = val.GFunction
		gme.NeedsContext = val.NeedsContext

		tableEntry := classloader.MTentry{
			MType: 'G',
			Meth:  gme,
		}

		classloader.AddEntry(tbl, key, tableEntry)
	}
}

// do-nothing Go function shared by several source files
func justReturn([]interface{}) interface{} {
	return nil
}
