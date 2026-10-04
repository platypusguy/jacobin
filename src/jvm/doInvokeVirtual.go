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
				globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
				errMsg := "INVOKEVIRTUAL: Class method not found: " + fqn
				status := exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, fr)
				if status != exceptions.Caught {
					return ERROR_OCCURRED // applies only if in test
				}
				return RESUME_HERE // caught
			}
		}
	}

	// if we got here, we have a method to call in mtEntry.Meth

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
	// then follow the JVM spec and push the objectRef and the parameters to the function
	// as parameters. Consult:
	// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.invokevirtual

	if mtEntry.MType == 'G' { // so we have a golang function
		return invokeVirtualGfunction(fr, mtEntry, className, methodName, methodType)
	}

	// 	To resolve a J method (i.e., a Java method) for invokevirtual:
	//  - If it's a native Java function (written in C/C++), Jacobin does not support it.
	//  - Get the reference object from the stack.
	// 	- Try searching the reference object class and its superclass chain.
	// 	- If the method is not found, try the reference object class's interface hierarchy (JVM spec 5.4.3.4).
	if mtEntry.MType == 'J' { // it's a Java function
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_NATIVE > 0 {
			// Native code
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Native method requested: " + fqn
			status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// The run-time class object is on the stack, below the method arguments.
		// To locate it, get the number of arguments for the method.
		nslots := len(util.ParseIncomingParamsFromMethTypeString(methodType))

		// Extract the reference object from the stack.
		refObjSlot := fr.TOS - nslots
		if refObjSlot < 0 || refObjSlot >= len(fr.OpStack) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Operand stack underflow locating reference object"
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		refObj, ok := fr.OpStack[refObjSlot].(*object.Object)
		if !ok {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Stack reference object is nil"
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
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
				globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
				errMsg := "INVOKEVIRTUAL: Concreted class method not found: " + fqn
				status := exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, fr)
				if status != exceptions.Caught {
					return ERROR_OCCURRED // applies only if in test
				}
				return RESUME_HERE // caught
			}

			// Found an interface default method.
			className = ret.(string)
		}

		// Resolve to a G function?
		if mtEntry.MType == 'G' {
			return invokeVirtualGfunction(fr, mtEntry, className, methodName, methodType)
		}

		// It's a J function. Get its JmEntry.
		m = mtEntry.Meth.(classloader.JmEntry)
		fqn = className + "." + methodName + methodType

		// If an empty code segment, that's an error. It's probably abstract or an interface.
		// In this case, flag it as an AbstractMethodError.
		if len(m.Code) == 0 {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: J class method code is empty: " + fqn
			status := exceptions.ThrowEx(excNames.AbstractMethodError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// Create the next frame to execute.
		nextFrame, err := createAndInitNewFrame(
			className, methodName, methodType, &m, true, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Error creating frame in: " + fqn
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		fr.PC += 3                         // 2 for PC slot, move to next bytecode before exiting
		fr.FrameStack.PushFront(nextFrame) // push the new frame, it'll be run by the next interpreter loop
		return 0
	}
	return ERROR_OCCURRED // in theory, unreachable
}

// Execute an INVOKEVIRTUAL G function.
func invokeVirtualGfunction(fr *frames.Frame,
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
