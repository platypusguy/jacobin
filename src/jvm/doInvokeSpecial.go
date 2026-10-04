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
	"math"
	"runtime/debug"
)

// throwInvokeSpecial records the Go stack, throws excType with msg in frame fr,
// and returns the interpreter return code: RESUME_HERE if a Java handler caught
// the exception, otherwise ERROR_OCCURRED (the uncaught case applies only in tests).
// NOTE: if the excNames constants are not plain ints, change the type of excType.
func throwInvokeSpecial(fr *frames.Frame, excType int, msg string) int {
	globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
	if exceptions.ThrowEx(excType, msg, fr) != exceptions.Caught {
		return ERROR_OCCURRED // applies only if in test
	}
	return RESUME_HERE // caught
}

// 0xB7 INVOKESPECIAL
//
// Invokes a method WITHOUT dynamic dispatch: the method named by the CP entry is
// the one that runs. It is used for constructors (<init>), private methods, and
// super.method() calls. Compare doInvokeVirtual, which dispatches on the
// receiver's runtime class.
//
// The declared method is cached per call site in CP.CachedMethods, using the same
// CachedMeth protocol as doInvokeVirtual and doInvokeStatic, so repeated calls
// skip the CP decoding and the method lookup.
func doInvokeSpecial(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)

	// This CP entry can be rewritten as a CachedMeth by doInvokeVirtual or
	// doInvokeStatic running in other threads that share this class's constant pool,
	// so the read must be lock-protected to avoid a torn or stale CpEntry. The
	// CachedMethods read shares the same read-lock section, so the slot cannot go
	// stale between the two reads.
	var mtEntry classloader.MTentry
	cached := false
	CP.Mutex.RLock()
	entry := CP.CpIndex[CPslot]
	if globals.CacheMeths && entry.Type == classloader.CachedMeth {
		mtEntry = CP.CachedMethods[entry.Slot]
		cached = true
	}
	CP.Mutex.RUnlock()

	var className, methodName, methodType string
	if cached {
		// Fast path: the cached entry carries string-pool indices for the names.
		className = *stringPool.GetStringPointer(mtEntry.MethClass)
		methodName = *stringPool.GetStringPointer(mtEntry.MethName)
		methodType = *stringPool.GetStringPointer(mtEntry.MethType)
	} else {
		// Slow path: decode the CP entry and look the method up.
		if entry.Type == classloader.Interface { // e.g. a private or default method in an interface
			className, methodName, methodType = classloader.GetMethInfoFromCPinterfaceRef(CP, CPslot)
		} else {
			className, methodName, methodType, _ = classloader.GetMethInfoFromCPmethref(CP, CPslot)
		}

		var err error
		mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
		if err != nil || mtEntry.Meth == nil {
			// TODO: search the classpath and retry
			// JVMS 5.4.3.3: a failed method resolution is a NoSuchMethodError.
			// The fully qualified name is built here, not up front, because the
			// interface-ref decoder does not return one.
			return throwInvokeSpecial(fr, excNames.NoSuchMethodError,
				"INVOKESPECIAL: Class method not found: "+className+"."+methodName+methodType)
		}

		// Cache the declared method for this call site. Interface refs are not
		// cached: other code may depend on the slot keeping its Interface type.
		if globals.CacheMeths && entry.Type != classloader.Interface {
			mtEntry.MethClass = stringPool.GetStringIndex(&className)
			mtEntry.MethName = stringPool.GetStringIndex(&methodName)
			mtEntry.MethType = stringPool.GetStringIndex(&methodType)

			CP.Mutex.Lock()
			// Re-check under the write lock: another goroutine may have cached this
			// slot first. Also never overflow the uint16 slot.
			if CP.CpIndex[CPslot].Type != classloader.CachedMeth && len(CP.CachedMethods) <= math.MaxUint16 {
				CP.CachedMethods = append(CP.CachedMethods, mtEntry)
				CP.CpIndex[CPslot] = classloader.CpEntry{
					Type: classloader.CachedMeth,
					Slot: uint16(len(CP.CachedMethods) - 1),
				}
			}
			CP.Mutex.Unlock()
		}
	}

	if mtEntry.MType == 'G' { // it's a golang method
		return invokeSpecialGfunction(fr, mtEntry, className, methodName, methodType)
	}

	if mtEntry.MType == 'J' { // it's a Java method
		// The arguments are correctly handled in createAndInitNewFrame()
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_NATIVE > 0 {
			// Jacobin cannot run native (C/C++) methods.
			return throwInvokeSpecial(fr, excNames.UnsupportedOperationException,
				"INVOKESPECIAL: Native method requested: "+className+"."+methodName+methodType)
		}
		// An empty code segment means the method is abstract (JVMS 6.5.invokespecial).
		if len(m.Code) == 0 {
			return throwInvokeSpecial(fr, excNames.AbstractMethodError,
				"INVOKESPECIAL: J class method code is empty: "+className+"."+methodName+methodType)
		}

		fram, err := createAndInitNewFrame(className, methodName, methodType, &m, true, fr)
		if err != nil {
			return throwInvokeSpecial(fr, excNames.InvalidStackFrameException,
				"INVOKESPECIAL: Error creating frame in: "+className+"."+methodName+methodType)
		}

		fr.PC += 3                    // 2 for the CP slot operand, plus 1 to move past the opcode; used when we return from the invoked method
		fr.FrameStack.PushFront(fram) // push the new frame; the next interpreter loop runs it
		return 0
	}
	return ERROR_OCCURRED // neither G nor J: in theory, unreachable
}

// invokeSpecialGfunction runs a G function (one implemented in Go rather than Java)
// for INVOKESPECIAL. It pops the arguments and the receiver off the caller's
// operand stack, calls the function, and pushes any return value.
// It returns 3 (2 for the CP slot operand + 1 for the opcode) on normal completion,
// otherwise an interpreter return code (RESUME_HERE or ERROR_OCCURRED).
func invokeSpecialGfunction(fr *frames.Frame, mtEntry classloader.MTentry,
	className, methodName, methodType string) int {

	// ParamSlots counts slots, not items, so doubles and longs are listed as two
	// slots. Popping that many stack entries is therefore correct.
	paramCount := mtEntry.Meth.(ghelpers.GMeth).ParamSlots

	// Arguments come off the stack in pop order; the receiver follows them.
	// Preallocated: arguments plus the receiver.
	params := make([]any, 0, paramCount+1)
	for i := 0; i < paramCount; i++ {
		params = append(params, pop(fr))
	}

	// The objectRef (the object whose method we're invoking) sits below the arguments.
	// JVMS: invokespecial on a null reference throws NullPointerException. The checked
	// assertion also covers a typed nil *object.Object.
	objRef, ok := pop(fr).(*object.Object)
	if !ok || objRef == nil {
		return throwInvokeSpecial(fr, excNames.NullPointerException,
			"INVOKESPECIAL: Stack reference object is nil in "+className+"."+methodName+methodType)
	}
	params = append(params, objRef)

	if globals.TraceInst {
		infoMsg := fmt.Sprintf("G-function: class=%s, meth=%s%s", className, methodName, methodType)
		trace.Trace(infoMsg)
	}

	ret := gfunction.RunGfunction(mtEntry, fr.FrameStack, &params, true, globals.TraceInst)
	if ret != nil {
		if gerr, isErr := ret.(error); isErr {
			if globals.GetGlobalRef().JacobinName == "test" {
				return ERROR_OCCURRED
			}
			if errors.Is(gerr, gfunction.CaughtGfunctionException) {
				return RESUME_HERE // resume at the present PC, which points to the exception code
			}
			// An error that is not a caught exception must not be ignored: previously
			// execution silently continued as if the call had succeeded.
			return ERROR_OCCURRED
		}
		// Not an error, so it is a legitimate return value, which we simply push.
		push(fr, ret)
	}
	// Any exception thrown by the G function has already been handled above.
	return 3 // 2 for the CP slot operand + 1 for the next bytecode
}
