/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
    "fmt"
    "jacobin/src/classloader"
    "jacobin/src/exceptions"
    "jacobin/src/gfunction"
    "jacobin/src/globals"
    "jacobin/src/object"
    "jacobin/src/shutdown"
    "jacobin/src/statics"
    "jacobin/src/stringPool"
    "jacobin/src/thread"
    "jacobin/src/trace"
    "jacobin/src/types"
    "os"
    "strings"
)

var globPtr *globals.Globals

// JVMrun is where everything begins
// The call to shutdown.Exit() exits the program (after some cleanup and logging); the reason
// it is here returned is because: in testing mode, the actual exit() call is side-stepped and
// instead an int is returned. This is necessary because calling exit() during testing exits
// the testing run as well.
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

	// initialize the globals package.
	// if globals.JacobinName == "test", then we're in test mode, which means
	// globals and log have been set in the testing function, likely to specific
	// values, so, don't reset them here.
	if globals.GetGlobalRef().JacobinName != "test" {
		// Not a test!
		_ = globals.InitGlobals(os.Args[0])
		stringPool.PreloadArrayClassesToStringPool()
	}
	globPtr = globals.GetGlobalRef()

	// Enable select functions via a global function variable. (This avoids circularity issues.)
	InitGlobalFunctionPointers()

	if globals.TraceInit {
		trace.Trace("running program: " + globPtr.JacobinName)
	}

	// load select static variables. Needs to be here b/c CLI might modify their values
	statics.PreloadStatics()

	// check for environmental variables that set JVM options
	globPtr.ClasspathRaw = os.Getenv("CLASSPATH")
	expandClasspth(globPtr)

	// handle the command-line interface (CLI) -- i.e., process the args
	LoadOptionsTable(*globPtr)
	err := HandleCli(os.Args, globPtr)
	if err != nil {
		return shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// some CLI options, like -version, show information and immediately exit.
	// This tests for that.
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
	if globPtr.StartingJar != "" {

		// A jar file was specified. Get the main class name from it.
		manifestClass, archive, err := classloader.GetMainClassFromJar(classloader.BootstrapCL, globPtr.StartingJar)
		if err != nil {
			errMsg := fmt.Sprintf("JVMrun: GetMainClassFromJar(%s) failed, err: %s", globPtr.StartingJar, err.Error())
			trace.Error(errMsg)
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}

		if manifestClass == "" {
			errMsg := fmt.Sprintf("JVMrun: no main manifest attribute in %s", globPtr.StartingJar)
			trace.Error(errMsg)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}

		// Set the main class name in the globals package.
		// This is used by the classloader to load the main class.
		globPtr.StartingClass = manifestClass
		mainClassNameIndex, _, err = classloader.LoadClassFromArchive(classloader.BootstrapCL, manifestClass, globPtr.StartingJar)
		if err != nil { // the exceptions message will already have been shown to the user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}

		// globals.InitGlobals has already run, setting up a provisional classpath by invoking globals.InitClasspath().
		// We will now override it with the classpath from the manifest or simply the jar path itself if there is no Class-Path manifest attribute.
		// * Update the archive classpath info.
		// * Update the globals classpath info.
		archive.UpdateArchiveWithClassPath()
		globPtr.ClasspathRaw = archive.ClasspathRaw
		globPtr.Classpath = archive.Classpath

 } else if globPtr.StartingClass != "" { // if a class file or class name was specified
        // Determine whether StartingClass is a filesystem path to a .class file
        // or a class name intended to be resolved via the classpath (including jars).
        starting := globPtr.StartingClass
        // Fast path: if the exact path exists on disk, load from file (preserves old behavior)
        if _, statErr := os.Stat(starting); statErr == nil {
            mainClassNameIndex, _, err = classloader.LoadClassFromFile(classloader.BootstrapCL, starting)
            if err != nil { // the exceptions message will already have been shown to the user
                return shutdown.Exit(shutdown.JVM_EXCEPTION)
            }
        } else {
            // Treat it as a class name. Normalize:
            // 1) strip trailing .class if present
            // 2) convert dots/backslashes to forward slashes for internal name
            // 3) trim any leading slashes
            name := strings.TrimSuffix(starting, ".class")
            nameSlash := strings.ReplaceAll(name, ".", "/")
            nameSlash = strings.ReplaceAll(nameSlash, "\\", "/")
            nameSlash = strings.TrimLeft(nameSlash, "/")

            // Load by name using the classpath search (dirs and jars)
            if err = classloader.LoadClassFromNameOnly(nameSlash); err != nil {
                // LoadClassFromNameOnly already emitted diagnostics
                return shutdown.Exit(shutdown.JVM_EXCEPTION)
            }
            // Record the main class string in the StringPool for later retrieval
            nameDot := strings.ReplaceAll(nameSlash, "/", ".")
            mainClassNameIndex = stringPool.GetStringIndex(&nameDot)
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
	// and in addition execute some initialization gfunctions (e.g., in javaLangThreadGroup.go)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)

	// Initialize the initial global thread groups
	gfunction.InitializeGlobalThreadGroups()

	// create the main thread
	if globPtr.UseOldThread { //
		MainThread = thread.CreateThread()
		MainThread.AddThreadToTable(globPtr)

		mainClass := stringPool.GetStringPointer(mainClassNameIndex)
		if globals.TraceInit {
			trace.Trace("Starting execution with: " + *mainClass)
		}

		// StartExec() runs the main thread. It does not return an error because all errors
		// will be handled one of three ways: 1) trapped in an exception, which shuts down the
		// JVM after processing the error; 2) a deferred catch of a go panic, which also shuts
		// down after processing the error; 3) an undeferred go panic, which should never occur.
		// Consequently, if StartExec() finishes, no errors were encountered.
		//
		// To test for errors, trap stderr, as many of the unit tests do.
		StartExec(*mainClass, &MainThread, globPtr)

	} else { // JACOBIN-732
		mainClass := stringPool.GetStringPointer(mainClassNameIndex)
		runnable := gfunction.NewRunnable(
			object.JavaByteArrayFromGoString(*mainClass),
			object.JavaByteArrayFromGoString("main"),
			object.JavaByteArrayFromGoString("([Ljava/lang/String;)V"))
		params := []interface{}{runnable, object.StringObjectFromGoString("main")}
		t := globals.GetGlobalRef().FuncInvokeGFunction(
			"java/lang/Thread.<init>(Ljava/lang/Runnable;Ljava/lang/String;)V", params)
		if globals.TraceInit {
			trace.Trace("Starting execution with: " + *mainClass)
		}
		// the thread is registered in thread.Run()
		thread.Run(t.(*object.Object))
	}

	return shutdown.Exit(shutdown.OK)
}

// InitGlobalFunctionPointers initializes the global function pointers in the globals package.
// These circumvent circular dependencies. A JVM is a textbook example of circular dependencies:
// For example, the interpreter necessarily needs to deal with objects and to call exceptions;
// exceptions are objects in the JVM, and certain functions that the object package contains
// can throw exceptions. A typical golang solution is to stuff objects and exceptions into the
// same package, but we prefer to keep them separate for ease of comprehension and navigability,
// so we use global function pointers.
func InitGlobalFunctionPointers() {
	globalPtr := globals.GetGlobalRef()
	globalPtr.FuncInstantiateClass = InstantiateClass
	globalPtr.FuncInvokeGFunction = gfunction.Invoke
	globalPtr.FuncMinimalAbort = exceptions.MinimalAbort
	globalPtr.FuncRunThread = RunJavaThread
	globalPtr.FuncThrowException = exceptions.ThrowExNil
	globalPtr.FuncFillInStackTrace = gfunction.FillInStackTrace
}
