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
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"runtime/debug"
	"strings"
)

// 0xB8 INVOKESTATIC
func doInvokeStatic(fr *frames.Frame, _ int64) int {
	var className, methodName, methodType, fqn string
	var mtEntry classloader.MTentry
	var shouldCacheMeth bool
	var k *classloader.Klass
	var clInitRun byte
	var err error

	CPslot := (int(fr.Meth[fr.PC+1]) << 8) | int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	// we don't verify the validity of the CP slot b/c that's done in codeCheck.
	CP := fr.CP.(*classloader.CPool)
	CP.Mutex.RLock()
	entry := CP.CpIndex[CPslot]
	CP.Mutex.RUnlock()

	shouldCacheMeth = false
	if globals.CacheMeths { // this is the optimized and default path
		if entry.Type == classloader.CachedMeth {
			CP.Mutex.RLock()
			mtEntry = CP.CachedMethods[entry.Slot]
			CP.Mutex.RUnlock()
			goto processMTentry // don't check if ClInit has been run b/c this must be the 2nd (or later) run of this method
		} else { // it's our first time running this method, mark the method for caching
			shouldCacheMeth = true // and proceed with standard method lookup
		}
	}

	if entry.Type == classloader.Interface {
		className, methodName, methodType =
			classloader.GetMethInfoFromCPinterfaceRef(CP, CPslot)
		fqn = className + "." + methodName + methodType
	} else {
		className, methodName, methodType, fqn = // fqn is the fully qualified name of the method
			classloader.GetMethInfoFromCPmethref(CP, CPslot)
	}

	mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
	if err != nil || mtEntry.Meth == nil {
		// TODO: search the classpath and retry  <---  still a valid comment?
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "INVOKESTATIC: Class method not found: " + fqn
		status := exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	} else {
		if mtEntry.MethClass == 0 { // true in the case of a Gfunction
			mtEntry.MethClass = stringPool.GetStringIndex(&className)
			mtEntry.MethName = stringPool.GetStringIndex(&methodName)
			mtEntry.MethType = stringPool.GetStringIndex(&methodType)
		}
	}

	// before we can run the method, we need to either instantiate the class and/or
	// make sure that its static intializer block (if any) has been run. At this point,
	// all we know is that the class exists and has been loaded.
	k = classloader.MethAreaFetch(className)
	if k != nil {
		k.Data.CP.Mutex.Lock()
		clInitRun = k.Data.ClInit
		k.Data.CP.Mutex.Unlock()

		if clInitRun == types.ClInitNotRun {
			err = runInitializationBlock(k, nil, fr.FrameStack)
			if err != nil {
				globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("INVOKESTATIC: error running initializer block in %s", fqn)
				status := exceptions.ThrowEx(excNames.ExceptionInInitializerError, errMsg, fr)
				if status != exceptions.Caught {
					return ERROR_OCCURRED // applies only if in test
				}
				return RESUME_HERE // caught
			}
		}
	}

processMTentry: // at this point, we have the mtEntry

	if mtEntry.MethClass != 0 {
		className = *stringPool.GetStringPointer(mtEntry.MethClass)
		methodName = *stringPool.GetStringPointer(mtEntry.MethName)
		methodType = *stringPool.GetStringPointer(mtEntry.MethType)
	}

	// if this is the first time calling this method and we're using cached methods,
	// then cache this mtEntry
	if globals.CacheMeths && shouldCacheMeth {
		CP.Mutex.Lock() // update the CP with the cached method
		if CP.CpIndex[CPslot].Type != classloader.CachedMeth {
			CP.CachedMethods = append(CP.CachedMethods, mtEntry)
			CP.CpIndex[CPslot] = classloader.CpEntry{
				Type: classloader.CachedMeth,
				Slot: uint16(len(CP.CachedMethods) - 1)}
		}
		CP.Mutex.Unlock()
		shouldCacheMeth = false
	}

	if mtEntry.MType == 'G' {
		gmethData := mtEntry.Meth.(ghelpers.GMeth)
		paramCount := gmethData.ParamSlots
		var params []any
		for range paramCount {
			params = append(params, pop(fr))
		}

		if globals.TraceInst {
			var cachedStatus = ""
			if entry.Type == classloader.CachedMeth {
				cachedStatus = "(cached)"
			}
			infoMsg := fmt.Sprintf("G-function: %s.%s%s %s",
				className,
				methodName,
				methodType, cachedStatus)
			trace.Trace(infoMsg)
		}

		ret := gfunction.RunGfunction(mtEntry, fr.FrameStack, &params, false, globals.TraceInst)
		if ret != nil {
			switch ret.(type) {
			case error:
				if globals.GetGlobalRef().JacobinName == "test" {
					return ERROR_OCCURRED
				} else if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
					return RESUME_HERE // resume at the present PC, which points to the exception code
				}
			default:
				if !strings.HasSuffix(methodType, "V") { // if it's not an error,
					// then it's a legitimate return value, which we push provided
					// the method does not return void
					push(fr, ret)
				}
			}
		}
		return 3
		// any exception will already have been handled.
	} else if mtEntry.MType == 'J' {
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_STATIC == 0 {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESTATIC: Method is not static: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		if m.AccessFlags&classloader.ACC_ABSTRACT > 0 {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESTATIC: Abstract method requested: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.AbstractMethodError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		if m.AccessFlags&classloader.ACC_NATIVE > 0 {
			// Native code
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESTATIC: Native method requested: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.UnsatisfiedLinkError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		fram, err := createAndInitNewFrame(
			className, methodName, methodType, &m, false, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESTATIC: Error creating frame in: " +
				className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		fr.PC += 3                    // 2 == initial PC advance in this bytecode + 1 for next bytecode
		fr.FrameStack.PushFront(fram) // push the new frame
		return 0
	}
	return ERROR_OCCURRED // in theory, unreachable code
}
