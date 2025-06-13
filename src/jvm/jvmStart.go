/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/thread"
	"jacobin/trace"
	"jacobin/types"
	"os"
)

var globPtr *globals.Globals

// JVMrun is where everything begins
// The call to shutdown.Exit() exits the program (after some clean-up and logging); the reason
// it is here returned is because: in testing mode, the actual exit() call is side-stepped and
// instead an int is returned. This is necessary because calling exit() during testing exits
// the testing run as well).
func JVMrun() int {

	trace.Init()

	// capture any panics and print diagnostic data
	defer func() int {
		if r := recover(); r != nil {
			// we get here only on errors that are not intercepted at
			// the thread level. Essentially, very unexpected JVM errors
			rglobPtr := globals.GetGlobalRef()
			if rglobPtr.ErrorGoStack != "" {
				// if the ErrorGoStack is not empty, we earlier intercepted
				// the error, so print the stack captured at that point
				exceptions.ShowGoStackTrace(nil)
			} else {
				// otherwise show the stack as it is now
				exceptions.ShowGoStackTrace(r)
			}
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	// if globals.JacobinName == "test", then we're in test mode, which means
	// globals and log have been set in the testing function, likely to specific values,
	// so, don't reset them here.
	if globals.GetGlobalRef().JacobinName != "test" {
		// Not a test!
		_ = globals.InitGlobals(os.Args[0])
		stringPool.PreloadArrayClassesToStringPool()
	}
	globPtr = globals.GetGlobalRef()

	// Enable select functions via a global function variable. (This avoids circularity issues.)
	InitGlobalFunctionPointers(globPtr)

	if globals.TraceInit {
		trace.Trace("running program: " + globPtr.JacobinName)
	}

	// load static variables. Needs to be here b/c CLI might modify their values
	statics.PreloadStatics()

	// handle the command-line interface (CLI) -- i.e., process the args
	LoadOptionsTable(*globPtr)
	err := HandleCli(os.Args, globPtr)
	if err != nil {
		return shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// some CLI options, like -version, show data and immediately exit. This tests for that.
	if globPtr.ExitNow == true {
		return shutdown.Exit(shutdown.OK)
	}

	// Initialize classloaders and method area
	err = classloader.Init()
	if err != nil {
		return shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	classloader.LoadBaseClasses() // must follow classloader.Init()

	var mainClassNameIndex uint32
	if globPtr.StartingJar != "" { // if a jar file was specified, then load the main class from it
		manifestClass, err := classloader.GetMainClassFromJar(classloader.BootstrapCL, globPtr.StartingJar)

		if err != nil {
			errMsg := fmt.Sprintf("JVMrun: GetMainClassFromJar(%s) failed, err: %v", globPtr.StartingJar, err)
			trace.Error(errMsg)
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}

		if manifestClass == "" {
			errMsg := fmt.Sprintf("JVMrun: no main manifest attribute in %s", globPtr.StartingJar)
			trace.Error(errMsg)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		mainClassNameIndex, _, err = classloader.LoadClassFromJar(classloader.BootstrapCL, manifestClass, globPtr.StartingJar)
		if err != nil { // the exceptions message will already have been shown to user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	} else if globPtr.StartingClass != "" { // if a class file was specified, then load the main class from it
		mainClassNameIndex, _, err = classloader.LoadClassFromFile(classloader.BootstrapCL, globPtr.StartingClass)
		if err != nil { // the exceptions message will already have been shown to user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	} else {
		trace.Error("JVMrun: No starting class from a class file nor a jar")
		ShowUsage(os.Stdout)
		return shutdown.Exit(shutdown.APP_EXCEPTION)
	}

	// if assertions were enabled on the command line for the program, then make sure
	// that the assertion status is set in the Statics table w/ an entry corresponding
	// to the main class. 	// Otherwise, it was previously initialized to "disabled".
	if globPtr.Options["-ea"].Set {
		_ = statics.AddStatic("main.$assertionsDisabled",
			statics.Static{Type: types.Int, Value: types.JavaBoolFalse})
	}

	// the following was commented out per JACOBIN-327. Likely to be reinstated at some later point.
	// Preload the main class and its dependencies.
	// classloader.LoadReferencedClasses(mainClass)

	// initialize the MTable (table caching methods) and load the gfunctions
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)

	// create the main thread
	MainThread = thread.CreateThread()
	MainThread.AddThreadToTable(globPtr)

	mainClass := stringPool.GetStringPointer(mainClassNameIndex)
	if globals.TraceInit {
		trace.Trace("Starting execution with: " + *mainClass)
	}

	// StartExec() runs the main thread. It does not return an error because all errors
	// will be handled one of three ways: 1) trapped in an exception, which shuts down the
	// JVM after processing the error; 2) a deferred catch of a go panic, which also shuts
	// down after processing the error; 3) a undeferred go panic, which should never occur.
	// Consequently, if StartExec() finishes, no errors were encountered.
	//
	// To test for errors, trap stderr, as many of the unit tests do.
	StartExec(*mainClass, &MainThread, globPtr)

	return shutdown.Exit(shutdown.OK)
}

// InitGlobalFunctionPointers initializes the global function pointers in the globals package.
// These circumvent circular dependencies. A JVM is a textbook example of circular dependencies:
// For example, the intepreter necessarily needs to deal with objects and to call exceptions;
// exceptions are objects in the JVM, and certain functions that the object package contains
// can throw exceptions. A typical golang solution is to stuff objects and exceptions into the
// same package, but we prefer to keep them separate for ease of comprehension and navigability,
// so we use global function pointers.
func InitGlobalFunctionPointers(globalPtr *globals.Globals) {
	globalPtr.FuncInstantiateClass = InstantiateClass
	globalPtr.FuncMinimalAbort = exceptions.MinimalAbort
	globalPtr.FuncThrowException = exceptions.ThrowExNil
	globalPtr.FuncFillInStackTrace = gfunction.FillInStackTrace
}
