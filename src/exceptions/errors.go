/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package exceptions

import (
	"container/list"
	"fmt"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/thread"
	"runtime/debug"
	"strings"
)

// routines for formatting error data when an error occurs inside the JVM

// Prints out the frame stack
func ShowFrameStack(source interface{}) {
	if globals.GetGlobalRef().JvmFrameStackShown == false {
		var entries *[]string
		switch source.(type) {
		case *thread.ExecThread:
			t := source.(*thread.ExecThread)
			entries = GrabFrameStack(t.Stack)
		case *list.List:
			entries = GrabFrameStack(source.(*list.List))
		}

		if len(*entries) == 0 {
			_ = log.Log("no further data available", log.SEVERE)
			return
		}

		// step through the list-based stack of called methods and print contents
		literals := *entries
		for i := 0; i < len(literals); i++ {
			_ = log.Log(literals[i], log.SEVERE)
		}
		globals.GetGlobalRef().JvmFrameStackShown = true
	}
}

// gets the JVM frame stack data and returns it as a slice of strings
func GrabFrameStack(fs *list.List) *[]string {
	var stackListing []string

	if fs == nil {
		// return an empty stack listing
		return &stackListing
	}
	frameStack := fs.Front()
	if frameStack == nil {
		// return an empty stack listing
		return &stackListing
	}

	// step through the list-based stack of called methods and print contents
	for e := frameStack; e != nil; e = e.Next() {
		val := e.Value.(*frames.Frame)
		methName := fmt.Sprintf("%s.%s", val.ClName, val.MethName)
		entry := fmt.Sprintf("Method: %-40s PC: %03d", methName, val.PC)
		stackListing = append(stackListing, entry)
	}
	return &stackListing
}

// takes the panic cause (as returned by the golang runtime) and prints the
// cause as determined by the runtime. Not sure it could ever be nil, but
// covering our bases nonetheless.
func ShowPanicCause(reason any) {
	// don't show the cause a second time
	if globals.GetGlobalRef().PanicCauseShown {
		return
	}

	// show the event that caused the panic
	if reason != nil {
		cause := fmt.Sprintf("%v", reason)
		_ = log.Log("\nerror: go panic because of: "+cause+"", log.SEVERE)
	} else {
		_ = log.Log("\nerror: go panic -- cause unknown", log.SEVERE)
	}
	globals.GetGlobalRef().PanicCauseShown = true
}

// ShowGoStackTrace explains that a panic has occurred and
// (to a limited extent why) and then prints the golang stack trace.
// Param stackInfo is the error returned when the panic occurred.
//
// Important: if you change the name of this function, you need to update
// exceptions.ShowGoStackTrace()--this function--which explicitly tests
// for this function name.
func ShowGoStackTrace(stackInfo any) {
	var stack string

	global := globals.GetGlobalRef()
	if global.GoStackShown {
		return
	}

	if stackInfo != nil && global.PanicCauseShown == false {
		ShowPanicCause(stackInfo)
	}

	// get the golang stack either b/c it was saved or fetch it new here
	if global.ErrorGoStack != "" {
		stack = global.ErrorGoStack
	} else {
		stack = string(debug.Stack())
	}
	entries := strings.Split(stack, "\n")

	_ = log.Log(" ", log.SEVERE) // print a blank line

	// print the remaining strings in the golang stack trace
	var i = 0
	for {
		if i < len(entries) {
			if strings.HasPrefix(entries[i], "runtime") ||
				strings.HasPrefix(entries[i], "jacobin/exceptions.ShowGoStackTrace") ||
				strings.HasPrefix(entries[i], "jacobin/exceptions.ThrowEx") {
				i += 2 // skip over runtime traces, we just want app data
				continue
			}
			_ = log.Log(entries[i], log.SEVERE)
			i += 1
		} else {
			break
		}
	}
	global.GoStackShown = true
}

// GetExceptionNameFromClassName extracts the name of the exception from the name of the exception class
func GetExceptionNameFromClassName(className string) string {
	var excName = ""

	// if it's not an excepted exception or error class name, return an empty string
	if !strings.HasSuffix(className, "xception") && !strings.HasSuffix(className, "rror") {
		return ""
	}

	lastSlash := strings.LastIndex(className, "/")
	excName = className[lastSlash+1:]

	return excName
}
