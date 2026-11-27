/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
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
	"runtime/debug"
)

// 0xB9 INVOKEINTERFACE
func doInvokeInterface(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	count := fr.Meth[fr.PC+3]
	zeroByte := fr.Meth[fr.PC+4]

	CP := fr.CP.(*classloader.CPool)
	CP.Mutex.RLock()
	CPentry := CP.CpIndex[CPslot]
	CP.Mutex.RUnlock()
	if CPentry.Type != classloader.Interface || zeroByte != 0 { // remove the zeroByte test later
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("INVOKEINTERFACE: CP entry type (%d) did not point to an interface method type (%d)",
			CPentry.Type, classloader.Interface)
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, fr) // this is the error thrown by JDK
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	method := CP.InterfaceRefs[CPentry.Slot]

	// get the class entry from this method
	interfaceRef := method.ClassIndex
	CP.Mutex.RLock()
	interfaceNameIndex := CP.ClassRefs[CP.CpIndex[interfaceRef].Slot]
	CP.Mutex.RUnlock()
	interfaceNamePtr := stringPool.GetStringPointer(interfaceNameIndex)
	interfaceName := *interfaceNamePtr

	// get the method name for this method
	nAndTindex := method.NameAndType
	CP.Mutex.RLock()
	nAndTentry := CP.CpIndex[nAndTindex]
	CP.Mutex.RUnlock()
	nAndTslot := nAndTentry.Slot
	nAndT := CP.NameAndTypes[nAndTslot]
	interfaceMethodNameIndex := nAndT.NameIndex
	interfaceMethodName := classloader.FetchUTF8stringFromCPEntryNumber(CP, interfaceMethodNameIndex)

	// get the signature for this method
	interfaceMethodSigIndex := nAndT.DescIndex
	interfaceMethodType := classloader.FetchUTF8stringFromCPEntryNumber(
		CP, interfaceMethodSigIndex)

	// now get the objRef pointing to the class containing the call to the method
	// described just previously. It is located on the f.OpStack below the args to
	// be passed to the method.
	// The objRef object has previously been instantiated and its constructor called.
	objRef := fr.OpStack[fr.TOS-int(count)+1]
	if objRef == nil {
		errMsg := fmt.Sprintf("INVOKEINTERFACE: object whose method, %s, is invoked is null",
			interfaceName+interfaceMethodName+interfaceMethodType)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Make sure that the reference is truly an object.
	switch objRef.(type) {
	case *object.Object:
	default:
		errMsg := fmt.Sprintf("INVOKEINTERFACE: objRef is not an object, observed type: %T", objRef)
		status := exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// get the name of the objectRef's class, and make sure it's loaded
	objRefClassName := *(stringPool.GetStringPointer(objRef.(*object.Object).KlassName))
	if err := classloader.LoadClassFromNameOnly(objRefClassName); err != nil {
		// in this case, LoadClassFromNameOnly() will have already thrown the exception
		if globals.JacobinHome() == "test" {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	class := classloader.MethAreaFetch(objRefClassName)
	if class == nil {
		// in theory, this can't happen due to immediately previous loading, but making sure
		errMsg := fmt.Sprintf("INVOKEINTERFACE: class %s not found", objRefClassName)
		status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	var mtEntry classloader.MTentry
	var err error
	mtEntry, err = locateInterfaceMeth(class, fr, objRefClassName, interfaceName,
		interfaceMethodName, interfaceMethodType)
	if err != nil { // any error will already have been handled
		return ERROR_OCCURRED
	}

	if mtEntry.MType == 'J' {
		var className string
		entry := mtEntry.Meth.(classloader.JmEntry)
		if entry.AccessFlags&classloader.ACC_PRIVATE > 0 {
			className = interfaceName
		} else {
			className = (*class.Data).Name
		}
		// the args and the objRef are popped off the stack by following call
		fram, err := createAndInitNewFrame(
			className, interfaceMethodName, interfaceMethodType, &entry, true, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEINTERFACE: Error creating frame in: " + className + "." +
				interfaceMethodName + interfaceMethodType
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
		}
		fr.PC += 5                    // 2 for CP slot, 1 for count, 1 for zero byte, 1 for next bytecode
		fr.FrameStack.PushFront(fram) // push the new frame
		return 0                      // forcing execution of the new frame
	} else if mtEntry.MType == 'G' { // it's a gfunction (i.e., a native function implemented in golang)
		gmethData := mtEntry.Meth.(ghelpers.GMeth)
		paramCount := gmethData.ParamSlots
		var params []any
		for i := 0; i < paramCount; i++ {
			params = append(params, pop(fr))
		}

		// now push the objectRef (the object whose method we're invoking)
		params = append(params, objRef)
		_ = pop(fr) // the objRef is still on the stack, so pop it off now

		if globals.TraceInst {
			infoMsg := fmt.Sprintf("G-function: interface=%s, meth=%s%s", interfaceName, interfaceName, interfaceMethodType)
			trace.Trace(infoMsg)
		}
		ret := gfunction.RunGfunction(
			mtEntry, fr.FrameStack, &params, true,
			globals.TraceVerbose)
		if ret != nil {
			switch ret.(type) {
			case error:
				if globals.GetGlobalRef().JacobinName == "test" {
					return ERROR_OCCURRED
				} else if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
					return ERROR_OCCURRED
				}
			default: // if it's not an error, then it's a legitimate return value, which we simply push
				push(fr, ret)
				return 5 // 2 for CP slot + 1 for count, 1 for zero byte, and 1 for next bytecode
			}
		}
		// Normal nil return.
		return 5 // 2 for CP slot + 1 for count, 1 for zero byte, and 1 for next bytecode
	}

	// Impossible return point unless there is a preceding software error.
	globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
	errMsg := "INVOKEINTERFACE: Logic error creating frame in: " + (*class.Data).Name + "." +
		interfaceMethodName + interfaceMethodType
	status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
	if status != exceptions.Caught {
		return ERROR_OCCURRED // applies only if in test
	}
	return notImplemented(fr, 0) // in theory, unreachable code
}

// searchForDefaultInterfaceFunction searches for a default method in all interfaces
// implemented by className (and their superinterfaces), following JVM rules.
//
// JVM rules for default methods:
//  1. Start with all interfaces directly implemented by the class.
//  2. Recursively check superinterfaces.
//  3. If multiple default methods are found with the same signature, throw
//     IncompatibleClassChangeError (ambiguity).
//  4. Otherwise, return the most-specific default method (closest to the class).
//
// Returns the interface name and mtEntry if a default method is found, or nil if not.
//
// Reference: https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.4.3.4
//
// Note: this function is similar in many aspects to locateInterfaceMeth() in run.go.
// The two functions might eventually be integrated into one.
func searchForDefaultInterfaceFunction(
	fr *frames.Frame, className, methodName, methodType string,
) (any, classloader.MTentry) {

	var mtEntry classloader.MTentry

	// Fetch runtime class metadata
	klass := classloader.MethAreaFetch(className)
	if klass == nil || len(klass.Data.Interfaces) == 0 {
		// Class does not implement any interfaces
		return nil, mtEntry
	}

	// Recursive helper to search an interface and its superinterfaces.
	// visited is passed in fresh for every top-level search so that
	// interfaces shared by different branches (diamond inheritance of
	// default methods) are not incorrectly skipped.
	var searchInterface func(intfName string, visited map[string]bool) (string, classloader.MTentry)
	searchInterface = func(intfName string, visited map[string]bool) (string, classloader.MTentry) {
		if visited[intfName] {
			return "", classloader.MTentry{}
		}
		visited[intfName] = true

		intfClass := classloader.MethAreaFetch(intfName)
		if intfClass == nil {
			return "", classloader.MTentry{}
		}

		// Step 3a: Build fully qualified method name and check method table
		fqn := intfName + "." + methodName + methodType
		mtEntry = classloader.GetMtableEntry(fqn)
		if mtEntry.Meth != nil {
			m := mtEntry.Meth.(classloader.JmEntry)
			if len(m.Code) > 0 {
				// Found a non-abstract method → candidate default method
				return intfName, mtEntry
			}
		}

		// Step 3b: Recurse into superinterfaces
		for _, idx := range intfClass.Data.Interfaces {
			superIntfName := *stringPool.GetStringPointer(uint32(idx))
			if resName, resMT := searchInterface(superIntfName, visited); resName != "" {
				return resName, resMT
			}
		}

		return "", classloader.MTentry{}
	}

	// Track all candidates found
	candidates := make([]struct {
		intf string
		mt   classloader.MTentry
	}, 0)

	// Search all interfaces directly implemented by this class
	for _, idx := range klass.Data.Interfaces {
		intfName := *stringPool.GetStringPointer(uint32(idx))
		if resName, resMT := searchInterface(intfName, make(map[string]bool)); resName != "" {
			candidates = append(candidates, struct {
				intf string
				mt   classloader.MTentry
			}{resName, resMT})
		}
	}

	// De-duplicate candidates by the interface that actually declares the
	// default method. In diamond-inheritance cases (e.g. a class implements
	// two interfaces that both extend a common super-interface and neither
	// overrides the default method), the search above will find the *same*
	// declaring interface via more than one directly-implemented interface.
	// That is not a genuine conflict per JVMS 5.4.3.3 -- only distinct
	// declaring interfaces providing competing default methods are ambiguous.
	if len(candidates) > 1 {
		seen := make(map[string]bool, len(candidates))
		deduped := candidates[:0:0]
		for _, c := range candidates {
			if !seen[c.intf] {
				seen[c.intf] = true
				deduped = append(deduped, c)
			}
		}
		candidates = deduped
	}

	// JVM rule: handle multiple candidates
	if len(candidates) > 1 {
		// Conflict detected: multiple default methods with same signature
		errMsg := "INVOKEVIRTUAL: Ambiguous default method for " + className + "." + methodName + methodType
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, fr)
		// The exception has already been thrown here (and either caught, in which
		// case execution should resume in the catch block, or it was uncaught, in
		// which case the interpreter loop must stop). Either way, the caller must
		// not attempt any further resolution or throw a second exception for the
		// same bytecode, so we signal "already handled" via defaultMethodException.
		if status == exceptions.Caught {
			return defaultMethodException(RESUME_HERE), classloader.MTentry{}
		}
		return defaultMethodException(ERROR_OCCURRED), classloader.MTentry{}
	}

	// Return the default method (unambiguous single candidate)
	if len(candidates) == 1 {
		return candidates[0].intf, candidates[0].mt
	}

	// No default method found
	return nil, mtEntry
}

// defaultMethodException wraps an interpreter return code (RESUME_HERE or
// ERROR_OCCURRED) so that searchForDefaultInterfaceFunction's caller can
// distinguish "an exception was already thrown, just return this code" from
// "no default method was found, go ahead and throw NoSuchMethodError".
type defaultMethodException int
