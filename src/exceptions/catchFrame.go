/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"container/list"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/util"
)

// the three likely superclasses of a given exception
var javaLangError = "java/lang/Error"
var javaLangException = "java/lang/Exception"
var javaLangThrowable = "java/lang/Throwable"

// This routine looks for a handler for the given exception (excName) in the
// current frame stack working its way up the frame stack (fs). If one is found,
// it returns a pointer to that frame, otherwise it returns nil. Param pc is the
// program counter in the current frame where the execption was thrown.
func FindCatchFrame(fs *list.List, exceptName string, pc int) (*frames.Frame, int) {
	excName := util.ConvertClassFilenameToInternalFormat(exceptName)

	var excFrame *frames.Frame // the catch frame
	var excPC int              // the program counter for the catch logic in the catch frame

	for fr := fs.Front(); fr != nil; {
		var f = fr.Value.(*frames.Frame)
		var searchPC int
		if f.ExceptionPC == -1 {
			searchPC = f.PC
		} else {
			searchPC = f.ExceptionPC
		}

		excFrame, excPC = locateExceptionFrame(f, excName, searchPC)
		if excFrame != nil {
			break
		} else { // if the exception was not found in this frame, we delete the current frame
			// unless it's the last frame
			if fr.Next() == nil {
				return nil, -1
			}

			fr = fr.Next()
		}
	}
	return excFrame, excPC
}

// locateExceptionFrame (private to package exceptions) is a helper function for FindCatchFrame
func locateExceptionFrame(f *frames.Frame, excName string, pc int) (*frames.Frame, int) {
	// get the method and check for an exception catch table
	// get the full method nameclassloader.MTable = {map[string]classloader.MTentry}
	fullMethName := f.ClName + "." + f.MethName + f.MethType
	methEntry, found := classloader.MTable[fullMethName]
	if !found {
		errMsg := fmt.Sprintf("locateExceptionFrame: Method %s not found in MTable", fullMethName)
		minimalAbort(excNames.InternalException, errMsg)
	}

	if methEntry.MType != 'J' {
		errMsg := fmt.Sprintf("locateExceptionFrame: Method %s is a native method", fullMethName)
		minimalAbort(excNames.InternalException, errMsg)
	}

	method := methEntry.Meth.(classloader.JmEntry)
	if method.Exceptions == nil {
		if globals.TraceVerbose {
			infoMsg := fmt.Sprintf("locateExceptionFrame: Method %s has no exception table", fullMethName)
			trace.Trace(infoMsg)
		}
		return nil, -1 // no exception handler was found
	}

	// if we got this far, the method has an exception table
	for i := 0; i < len(method.Exceptions); i++ {
		entry := method.Exceptions[i]
		// per https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.7.3
		// the StartPC value is inclusive, the EndPC value is exclusive
		if pc >= entry.StartPc && pc < entry.EndPc {
			// found a handler, now check that it's for the right exception
			CP := f.CP.(*classloader.CPool)
			catchName :=
				classloader.GetClassNameFromCPclassref(CP, uint16(entry.CatchType))

			// TODO: add support for checking for subclasses
			// In the meantime, check for a direct match or one of the typical
			// superclasses.
			if catchName == excName ||
				catchName == "java/lang/Throwable" ||
				catchName == "java/lang/Exception" ||
				catchName == "java/lang/Error" {
				return f, entry.HandlerPc
			} else {
				catchClass := classloader.MethAreaFetch(catchName)
				if catchClass == nil { // if the class isn't found, skip it
					continue // in theory, this should be impossible
				}
				if catchClass.Data.SuperclassIndex == stringPool.GetStringIndex(&javaLangThrowable) ||
					catchClass.Data.SuperclassIndex == stringPool.GetStringIndex(&javaLangException) ||
					catchClass.Data.SuperclassIndex == stringPool.GetStringIndex(&javaLangError) {
					return f, entry.HandlerPc
				}
			}
		}
	}
	// if we got this far, no exception handler was found
	return nil, -1
}
