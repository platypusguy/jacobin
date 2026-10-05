/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/util"
	"runtime/debug"
)

/* Implementation of the INVOKEVIRTUAL instruction (opcode 0xB6)
 * Based on the Java SE 21 JVM Specification:
 *   - §6.5.invokevirtual
 *   - §5.4.3.3  Method Resolution
 *   - §5.4.5    Method Overriding
 *   - §5.4.6    Method Selection
 */

// 0xB6 INVOKEVIRTUAL
func doInvokeVirtual(fr *frames.Frame, _ int64) int {
	var className, methodName, methodType, fqn string
	var mtEntry classloader.MTentry
	var shouldCacheMeth bool
	var err error

	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)                                // codeCheck.go ensures that CPslot is a valid index to a methodRef
	CP.Mutex.RLock()
	entry := CP.CpIndex[CPslot]
	CP.Mutex.RUnlock()

	shouldCacheMeth = false
	if globals.CacheMeths { // this is the optimized and default path
		if entry.Type == classloader.CachedMeth {
			className, methodName, methodType, fqn = classloader.GetMethInfoFromCPmethref(CP, CPslot)
			CP.Mutex.RLock()
			mtEntry = CP.CachedMethods[entry.Slot]
			CP.Mutex.RUnlock()
			goto processMTentry
		} else { // it's our first time running this method, mark the method for caching
			shouldCacheMeth = true // and proceed with standard method lookup
		}
	}

	/* Determining the method to invoke is a two-step process: method resolution and
	 * method selection.
	 *
	 * Resolution is the process of locating the method to invoke as specified in
	 * the method reference (and the CP).
	 *
	 * Selection is the process of selecting  the method to invoke, by searching for
	 * any methods in the objectRef's class hierarchy that override the method in the
	 * method reference. Interfaces are searched in this selection process. If no
	 * overriding method is found, the method located in the resolution process is invoked.
	 */

	// ==== Resolution step. See JVM spec §5.4.3.3 =====
	// Get the method table entry for the FQN indicated in CP.
	className, methodName, methodType, fqn = classloader.GetMethInfoFromCPmethref(CP, CPslot)
	mtEntry = classloader.GetMtableEntry(fqn)
	if mtEntry.Meth == nil { // if the method is not in the method table, search classes or superclasses
		mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
	}

	// Not found after a class-superclass search. Check the interfaces.
	if err != nil || mtEntry.Meth == nil { // the method is not in the superclasses, so check interfaces.
		// When a class implements an interface and inherits default methods (or doesn't override them),
		// the compiler generates INVOKEVIRTUAL
		klass := classloader.MethAreaFetch(className)
		if klass != nil && len(klass.Data.Interfaces) > 0 {
			for i := 0; i < len(klass.Data.Interfaces); i++ {
				index := uint32(klass.Data.Interfaces[i])
				interfaceName := *stringPool.GetStringPointer(index)
				mtEntry, err = locateInterfaceMeth(klass, fr, className, interfaceName, methodName, methodType)
				if mtEntry.Meth != nil {
					// found a match.
					break
				}
			} // end of search of interfaces if method has any

			// Any matches in the interfaces?
			if err != nil || mtEntry.Meth == nil {
				// method was not found in interfaces, so throw an exception
				return throwInvokevirtualException(fr, excNames.NoSuchMethodError, "Class method not found: "+fqn)
			}
		}
	}

	// if we got here, we have a method to call in mtEntry.Meth

	// ==== Selection step See §5.4.5, §5.4.6 ====
processMTentry:
	// if this is the first time calling this method and we're using cached methods,
	// then cache this mtEntry
	if (globals.CacheMeths && shouldCacheMeth) || !globals.CacheMeths {
		mtEntry.MethClass = stringPool.GetStringIndex(&className)
		mtEntry.MethName = stringPool.GetStringIndex(&methodName)
		mtEntry.MethType = stringPool.GetStringIndex(&methodType)
		if globals.CacheMeths && shouldCacheMeth {
			CP.Mutex.Lock()
			if CP.CpIndex[CPslot].Type != classloader.CachedMeth { // someone else may have won the race
				CP.CachedMethods = append(CP.CachedMethods, mtEntry)
				CP.CpIndex[CPslot] = classloader.CpEntry{
					Type: classloader.CachedMeth,
					Slot: uint16(len(CP.CachedMethods) - 1),
				}
			}
			CP.Mutex.Unlock()
			shouldCacheMeth = false
		}
	}

	if globals.CacheMeths {
		className = *stringPool.GetStringPointer(mtEntry.MethClass)
		methodName = *stringPool.GetStringPointer(mtEntry.MethName)
		methodType = *stringPool.GetStringPointer(mtEntry.MethType)
	}

	// if we have a gFunction (that is, one implemented in golang, rather than Java),
	// then follow the JVM spec and push the objectRef and the parameters to the method
	// as parameters and execute that method.  Consult:
	// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.invokevirtual

	if mtEntry.MType == 'G' { // so we have a golang function
		return invokevirtualGfunction(fr, mtEntry, className, methodName, methodType)
	}

	// 	To resolve a J method (i.e., a Java method) for invokevirtual:
	//  - If it's a native Java function (written in C/C++), Jacobin does not support it.
	//  - Get the reference object from the stack.
	// 	- S the reference object class and its superclasses for any override of the method.
	// 	- If an override is not found, try the reference object class's interface hierarchy.
	//  - If still not found, use the present mtEntry.
	if mtEntry.MType == 'J' { // it's a Java function
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_NATIVE > 0 { // it's native code
			return throwInvokevirtualException(fr, excNames.UnsatisfiedLinkError,
				"Native method requested: "+fqn)
		}

		// The run-time class object is on the stack, below the method arguments.
		// To locate it, get the number of arguments for the method.
		nslots := len(util.ParseIncomingParamsFromMethTypeString(methodType))

		// Extract the reference object from the stack.
		refObjSlot := fr.TOS - nslots
		if refObjSlot < 0 || refObjSlot >= len(fr.OpStack) {
			return throwInvokevirtualException(fr, excNames.NullPointerException,
				"Operand stack underflow locating reference object")
		}

		refObj, ok := fr.OpStack[refObjSlot].(*object.Object)
		if !ok {
			return throwInvokevirtualException(fr, excNames.ClassCastException,
				"Stack reference object is nil")
		}

		// Get the reference object class name.
		clNameIdx := refObj.KlassName
		className = *(stringPool.GetStringPointer(clNameIdx))

		// === Method resolution ===
		// First, try superclass resolution.
		mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
		if err != nil || mtEntry.Meth == nil {
			// That did not succeed. So, try for an interface default method.
			var ret any
			ret, mtEntry = searchForDefaultInterfaceFunction(fr, className, methodName, methodType)
			if code, alreadyHandled := ret.(defaultMethodException); alreadyHandled {
				// searchForDefaultInterfaceFunction already threw (and possibly
				// caught) an exception, e.g. IncompatibleClassChangeError for an
				// ambiguous default method. Do not throw a second exception.
				return int(code)
			}
			if ret == nil {
				return throwInvokevirtualException(fr, excNames.NoSuchMethodError,
					"Class method not found: "+fqn)
			}

			// Found an interface default method.
			className = ret.(string)
		}

		// Resolve to a G function?
		if mtEntry.MType == 'G' {
			return invokevirtualGfunction(fr, mtEntry, className, methodName, methodType)
		}

		// It's a J function. Get its JmEntry.
		m = mtEntry.Meth.(classloader.JmEntry)
		fqn = className + "." + methodName + methodType

		// If an empty code segment, that's an error. It's probably abstract or an interface.
		// In this case, flag it as an AbstractMethodError.
		if len(m.Code) == 0 {
			return throwInvokevirtualException(fr, excNames.AbstractMethodError,
				"Empty code segment: "+fqn)
		}

		// Create the next frame to execute.
		nextFrame, err := createAndInitNewFrame(
			className, methodName, methodType, &m, true, fr)
		if err != nil {
			return throwInvokevirtualException(fr, excNames.InvalidStackFrameException,
				"Error creating frame in: "+fqn)
		}

		fr.PC += 3                         // 2 for PC slot, move to next bytecode before exiting
		fr.FrameStack.PushFront(nextFrame) // push the new frame, it'll be run by the next interpreter loop
		return 0
	}
	return ERROR_OCCURRED // in theory, unreachable
}

// invoke a G function.
func invokevirtualGfunction(fr *frames.Frame,
	mtEntry classloader.MTentry,
	className, methodName, methodType string) int {

	// Parameter array for G function.
	var params []any

	// Append the parameters/args off the stack to params.
	gmethData := mtEntry.Meth.(ghelpers.GMeth)
	paramCount := gmethData.ParamSlots
	for i := 0; i < paramCount; i++ {
		params = append(params, pop(fr))
	}

	// now get the objectRef (the object whose method we're invoking)
	popped := pop(fr)
	params = append(params, popped)

	// DYNAMIC DISPATCH for G-functions:
	// Check if the object's actual class has an override for this G-function.
	if objRef, ok := popped.(*object.Object); ok {
		objClassName := *(stringPool.GetStringPointer(objRef.KlassName))
		if objClassName != className {
			// Try to find a more specific G-function registration.
			specificFQN := objClassName + "." + methodName + methodType
			if specificGmeth, ok := ghelpers.MethodSignatures[specificFQN]; ok {
				// We found a more specific G-function. Use it instead.
				mtEntry = classloader.MTentry{
					Meth:  specificGmeth,
					MType: 'G',
				}
				className = objClassName
			}
		}
	}

	if globals.TraceInst {
		infoMsg := fmt.Sprintf("G-function: class=%s, meth=%s%s", className, methodName, methodType)
		trace.Trace(infoMsg)
	}

	// Execute the G function.
	ret := gfunction.RunGfunction(
		mtEntry, fr.FrameStack, &params, true, globals.TraceInst)
	if ret != nil {
		switch ret.(type) {
		case error: // only occurs in testing
			if globals.GetGlobalRef().JacobinName == "test" {
				return ERROR_OCCURRED
			}
			if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
				return RESUME_HERE // caught
			}
		default: // if it's not an error, then it's a legitimate return value, which we simply push
			push(fr, ret)
		}
		// any exception will already have been handled.
	}
	return 3 // 2 for CP slot + 1 for next bytecode
}

// throwInvokevirtualException records the Go stack, throws excType with msg in frame fr,
// and returns the interpreter return code: RESUME_HERE if a Java handler caught
// the exception, otherwise ERROR_OCCURRED (the uncaught case applies only in tests).
func throwInvokevirtualException(fr *frames.Frame, excType int, msg string) int {
	globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
	if exceptions.ThrowEx(excType, "INVOKEVIRTUAL: "+msg, fr) != exceptions.Caught {
		return ERROR_OCCURRED // applies only if in test
	}
	return RESUME_HERE // caught
}
