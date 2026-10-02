package jvm

import (
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/util"
	"math"
	"runtime/debug"
	"sync"
)

// dispatchKey identifies one resolved virtual call: the runtime class of the
// receiver plus the name and descriptor of the invoked method. All three are
// string-pool indices, so building a key hashes no strings.
//
// The class named in the constant-pool methodref is deliberately not part of the
// key: resolution starts from the receiver's class, so the result depends only
// on these three values.
type dispatchKey struct{ recv, name, desc uint32 }

// dispatchVal is the cached result of resolving a dispatchKey.
//   - mt:        the resolved method table entry.
//   - className: the owner of the resolved method: the receiver's class, or the
//     interface when the method is an interface default method.
//   - jm:        the resolved Java method, shared by every caller, so treat it as
//     read-only. It is nil for a G (Go) method, which has no JmEntry.
type dispatchVal struct {
	mt        classloader.MTentry
	className string
	jm        *classloader.JmEntry
}

var (
	dispatchCache sync.Map // dispatchKey -> *dispatchVal; entries are never evicted
	argSlotsCache sync.Map // method descriptor string -> int
)

// 0xB6 INVOKEVIRTUAL
//
// Two resolution steps happen, and both are cached:
//  1. The CP methodref is resolved to the declared method. This is cached per call
//     site in CP.CachedMethods.
//  2. The declared method is resolved against the receiver's runtime class, which
//     may override it. This is cached per (receiver class, method) in dispatchCache.
func doInvokeVirtual(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)                                // codeCheck.go ensures that CPslot is a valid index to a methodRef

	// A single read-lock section covers both the CpIndex read and the CachedMethods
	// read, so the slot cannot go stale between the two.
	var mtEntry classloader.MTentry
	cached := false
	CP.Mutex.RLock()
	entry := CP.CpIndex[CPslot]
	if globals.CacheMeths && entry.Type == classloader.CachedMeth {
		mtEntry = CP.CachedMethods[entry.Slot]
		cached = true
	}
	CP.Mutex.RUnlock()

	// Step 1: the declared method.
	var className, methodName, methodType string
	if cached {
		className = *stringPool.GetStringPointer(mtEntry.MethClass)
		methodName = *stringPool.GetStringPointer(mtEntry.MethName)
		methodType = *stringPool.GetStringPointer(mtEntry.MethType)
	} else {
		var rc int
		var ok bool
		mtEntry, className, methodName, methodType, rc, ok = resolveDeclaredVirtual(fr, CP, CPslot)
		if !ok {
			return rc
		}
		mtEntry.MethClass = stringPool.GetStringIndex(&className)
		mtEntry.MethName = stringPool.GetStringIndex(&methodName)
		mtEntry.MethType = stringPool.GetStringIndex(&methodType)

		if globals.CacheMeths {
			CP.Mutex.Lock()
			// Re-check under the write lock: another goroutine may have cached this
			// call site first. Also never overflow the uint16 slot.
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

	// A G function is implemented in Go rather than Java. Per the JVM spec, the
	// objectRef and the parameters are passed to the function as parameters. Consult:
	// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.invokevirtual
	if mtEntry.MType == 'G' {
		return invokeVirtualGfunction(fr, mtEntry, className, methodName, methodType)
	}
	if mtEntry.MType != 'J' {
		return ERROR_OCCURRED // neither G nor J: in theory, unreachable
	}

	// Step 2: the receiver. It sits on the operand stack below the method arguments.
	nslots := argSlots(methodType)
	refObj, ok := fr.OpStack[fr.TOS-nslots].(*object.Object)
	if !ok || refObj == nil { // an untyped nil fails the assertion; a typed nil is caught by == nil
		return throwInvokeVirtual(fr, excNames.NullPointerException, "INVOKEVIRTUAL: Stack reference object is nil")
	}

	// Look up the method that this receiver class actually runs, resolving it on the first call.
	key := dispatchKey{refObj.KlassName, mtEntry.MethName, mtEntry.MethType}
	var dv *dispatchVal
	if v, hit := dispatchCache.Load(key); hit {
		dv = v.(*dispatchVal)
	} else {
		var rc int
		dv, rc = resolveReceiverMethod(fr, mtEntry, refObj, methodName, methodType)
		if dv == nil {
			return rc // resolution failed and has already thrown
		}
		dispatchCache.Store(key, dv) // a concurrent duplicate store is harmless: both values are equivalent
	}

	// The receiver's class may override a Java method with a Go (G) one.
	if dv.mt.MType == 'G' {
		return invokeVirtualGfunction(fr, dv.mt, dv.className, methodName, methodType)
	}

	nextFrame, err := createAndInitNewFrame(dv.className, methodName, methodType, dv.jm, true, fr)
	if err != nil {
		return throwInvokeVirtual(fr, excNames.InvalidStackFrameException,
			"INVOKEVIRTUAL: Error creating frame in: "+dv.className+"."+methodName+methodType)
	}

	fr.PC += 3                         // 2 for the CP slot operand, plus 1 to move past the opcode
	fr.FrameStack.PushFront(nextFrame) // push the new frame; the next interpreter loop runs it
	return 0
}

// ResetDispatchCaches drops every memoized dispatch result. It must be called whenever
// the string pool or the method area is re-initialized (e.g., in test setup).
func ResetDispatchCaches() {
	dispatchCache.Clear()
	argSlotsCache.Clear()
}

// argSlots returns the number of parameters in a method descriptor, receiver
// excluded: one per entry returned by util.ParseIncomingParamsFromMethTypeString.
// The result is memoized so the hot path does not allocate a slice just to take len().
func argSlots(methodType string) int {
	if v, ok := argSlotsCache.Load(methodType); ok {
		return v.(int)
	}
	n := len(util.ParseIncomingParamsFromMethTypeString(methodType))
	argSlotsCache.Store(methodType, n)
	return n
}

// throwInvokeVirtual records the Go stack, throws excType with msg in frame fr,
// and returns the interpreter return code: RESUME_HERE if a Java handler caught
// the exception, otherwise ERROR_OCCURRED (the uncaught case applies only in tests).
// NOTE: if the excNames constants are not plain ints, change the type of excType.
func throwInvokeVirtual(fr *frames.Frame, excType int, msg string) int {
	globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
	if exceptions.ThrowEx(excType, msg, fr) != exceptions.Caught {
		return ERROR_OCCURRED // applies only if in test
	}
	return RESUME_HERE // caught
}

// resolveDeclaredVirtual resolves the method named by the CP methodref. It looks in
// the method table, then the class and its superclasses, then the interfaces
// (a class can inherit default methods without overriding them).
// On success ok is true. On failure the NoSuchMethodError has already been thrown,
// rc holds the interpreter return code, and ok is false.
func resolveDeclaredVirtual(fr *frames.Frame, CP *classloader.CPool, CPslot int) (
	mt classloader.MTentry, className, methodName, methodType string, rc int, ok bool) {

	var fqn string
	className, methodName, methodType, fqn = classloader.GetMethInfoFromCPmethref(CP, CPslot)

	mt = classloader.GetMtableEntry(className + "." + methodName + methodType)
	if mt.Meth == nil { // not in the method table: search the class and its superclasses
		var err error
		mt, err = classloader.FetchMethodAndCP(className, methodName, methodType)
		if err != nil {
			mt = classloader.MTentry{} // an error means not found, whatever was returned
		}
	}

	if mt.Meth == nil { // not in the superclasses: check the interfaces
		klass := classloader.MethAreaFetch(className)
		if klass != nil {
			for _, ifIdx := range klass.Data.Interfaces {
				interfaceName := *stringPool.GetStringPointer(uint32(ifIdx))
				found, err := locateInterfaceMeth(klass, fr, className, interfaceName, methodName, methodType)
				if err == nil && found.Meth != nil {
					mt = found
					break
				}
			}
		}
	}

	// One not-found check for every path. Previously it ran only when the class had
	// interfaces, so a class without any returned ERROR_OCCURRED with no exception.
	if mt.Meth == nil {
		// JVMS 5.4.3.3: a failed method resolution is a NoSuchMethodError.
		rc = throwInvokeVirtual(fr, excNames.NoSuchMethodError, "INVOKEVIRTUAL: Class method not found: "+fqn)
		return mt, className, methodName, methodType, rc, false
	}
	return mt, className, methodName, methodType, 0, true
}

// resolveReceiverMethod is the slow path of step 2. It runs once per
// (receiver class, method) pair and its result is cached in dispatchCache.
// It checks the declared method, searches the receiver's class chain, then falls
// back to interface default methods (JVMS 5.4.3.4).
// On success it returns the cached value. On failure the exception has already
// been thrown and it returns (nil, rc).
func resolveReceiverMethod(fr *frames.Frame, declared classloader.MTentry, refObj *object.Object,
	methodName, methodType string) (*dispatchVal, int) {

	declaredClass := *stringPool.GetStringPointer(declared.MethClass)
	declaredFqn := declaredClass + "." + methodName + methodType

	// Jacobin cannot run native (C/C++) methods. Go (G) methods are handled elsewhere.
	if dm, isJ := declared.Meth.(classloader.JmEntry); isJ && dm.AccessFlags&classloader.ACC_NATIVE > 0 {
		return nil, throwInvokeVirtual(fr, excNames.UnsupportedOperationException,
			"INVOKEVIRTUAL: Native method requested: "+declaredFqn)
	}

	// First the receiver's class and its superclasses, then interface default methods.
	className := *stringPool.GetStringPointer(refObj.KlassName)
	mt, err := classloader.FetchMethodAndCP(className, methodName, methodType)
	if err != nil || mt.Meth == nil {
		ret, imt := searchForDefaultInterfaceFunction(className, methodName, methodType)
		if ret == nil {
			return nil, throwInvokeVirtual(fr, excNames.NoSuchMethodError,
				"INVOKEVIRTUAL: Concreted class method not found: "+declaredFqn)
		}
		mt = imt
		className = ret.(string) // from here on, className is the interface that supplied the default method
	}

	if mt.MType == 'G' {
		return &dispatchVal{mt: mt, className: className}, 0
	}

	jm, isJ := mt.Meth.(classloader.JmEntry)
	if !isJ {
		return nil, throwInvokeVirtual(fr, excNames.NoSuchMethodError,
			"INVOKEVIRTUAL: Resolved entry is not a method: "+className+"."+methodName+methodType)
	}
	fqn := className + "." + methodName + methodType

	// The method that was resolved (not just the declared one) could be native.
	if jm.AccessFlags&classloader.ACC_NATIVE > 0 {
		return nil, throwInvokeVirtual(fr, excNames.UnsupportedOperationException,
			"INVOKEVIRTUAL: Native method requested: "+fqn)
	}
	// An empty code segment means the method is abstract or an interface stub.
	if len(jm.Code) == 0 {
		return nil, throwInvokeVirtual(fr, excNames.AbstractMethodError,
			"INVOKEVIRTUAL: J class method code is empty: "+fqn)
	}
	return &dispatchVal{mt: mt, className: className, jm: &jm}, 0
}
