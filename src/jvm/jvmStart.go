/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/thread"
	"jacobin/types"
	"os"
)

var globPtr *globals.Globals

// JVMrun is where everything begins
// The call to shutdown.Exit() exits the program (after some clean-up and logging); the reason
// it is here returned is because in testing mode, the actual exit() call is side-stepped and
// instead an int is returned (because calling exit() during testing exits the testing run as well).
func JVMrun() int {

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
	// globals and log have been set in the testing function. So, don't reset them here.
	if globals.GetGlobalRef().JacobinName != "test" {
		// Not a test!
		_ = globals.InitGlobals(os.Args[0])
		stringPool.PreloadArrayClassToStringPool()
		log.Init()
	}
	globPtr = globals.GetGlobalRef()

	// Enable functions call InstantiateClass through a global function variable. (This avoids circularity issues.)
	globPtr.FuncInstantiateClass = InstantiateClass

	_ = log.Log("running program: "+globPtr.JacobinName, log.FINE)

	var status error

	// load static variables. Needs to be here b/c CLI might modify their values
	statics.PreloadStatics()

	// handle the command-line interface (cli) -- i.e., process the args
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
	if globPtr.StartingJar != "" {
		manifestClass, err := classloader.GetMainClassFromJar(classloader.BootstrapCL, globPtr.StartingJar)

		if err != nil {
			_ = log.Log(err.Error(), log.INFO)
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}

		if manifestClass == "" {
			_ = log.Log(fmt.Sprintf("no main manifest attribute, in %s", globPtr.StartingJar), log.INFO)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		mainClassNameIndex, err = classloader.LoadClassFromJar(classloader.BootstrapCL, manifestClass, globPtr.StartingJar)
		if err != nil { // the exceptions message will already have been shown to user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	} else if globPtr.StartingClass != "" {
		mainClassNameIndex, err = classloader.LoadClassFromFile(classloader.BootstrapCL, globPtr.StartingClass)
		if err != nil { // the exceptions message will already have been shown to user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	} else {
		_ = log.Log("Error: No executable program specified. Exiting.", log.INFO)
		ShowUsage(os.Stdout)
		return shutdown.Exit(shutdown.APP_EXCEPTION)
	}

	// if assertions were enable on the command line for the program, then make sure
	// that it's set in the Statics table w/ an entry corresponding to the main class
	// Otherwise, it was previously set to disabled
	if globPtr.Options["-ea"].Set {
		_ = statics.AddStatic("main.$assertionsDisabled",
			statics.Static{Type: types.Int, Value: types.JavaBoolFalse})
	}

	// the following was commented out per JACOBIN-327.
	// Likely to be reinstated at some later point
	// classloader.LoadReferencedClasses(mainClass)

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadNatives(&classloader.MTable)

	// create the main thread
	MainThread = thread.CreateThread()
	MainThread.AddThreadToTable(globPtr)

	// begin execution
	mainClass := stringPool.GetStringPointer(mainClassNameIndex)
	_ = log.Log("Starting execution with: "+*mainClass, log.INFO)
	status = StartExec(*mainClass, &MainThread, globPtr)

	if status != nil {
		return shutdown.Exit(shutdown.APP_EXCEPTION)
	}
	return shutdown.Exit(shutdown.OK)
}
