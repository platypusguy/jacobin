/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"jacobin/thread"
	"os"
	"runtime/debug"
	"strings"
)

var Global globals.Globals

// JVMrun is where everything begins
// The call to shutdown.Exit() exits the program (after some clean-up and logging); the reason
// it is here returned is because in testing mode, the actual exit() call is side-stepped and
// instead an int is returned (because calling exit() during testing exits the testing run as well).
func JVMrun() int {

	// if globals.JacobinName == "test", then we're in test mode and globals and log have been set
	// in the testing function. So, don't reset them here.
	if globals.GetGlobalRef().JacobinName != "test" {
		Global = globals.InitGlobals(os.Args[0])
		log.Init()
	} else {
		Global = *globals.GetGlobalRef()
	}

	_ = log.Log("running program: "+Global.JacobinName, log.FINE)

	var status error

	// handle the command-line interface (cli) -- i.e., process the args
	LoadOptionsTable(Global)
	err := HandleCli(os.Args, &Global)
	if err != nil {
		return shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	// some CLI options, like -version, show data and immediately exit. This tests for that.
	if Global.ExitNow == true {
		return shutdown.Exit(shutdown.OK)
	}

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		return shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	classloader.LoadBaseClasses() // must follow classloader.Init
	classloader.StaticsPreload()

	var mainClass string

	if Global.StartingJar != "" {
		manifestClass, err := classloader.GetMainClassFromJar(classloader.BootstrapCL, Global.StartingJar)

		if err != nil {
			_ = log.Log(err.Error(), log.INFO)
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}

		if manifestClass == "" {
			_ = log.Log(fmt.Sprintf("no main manifest attribute, in %s", Global.StartingJar), log.INFO)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		mainClass, err = classloader.LoadClassFromJar(classloader.BootstrapCL, manifestClass, Global.StartingJar)
		if err != nil { // the exceptions message will already have been shown to user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	} else if Global.StartingClass != "" {
		mainClass, err = classloader.LoadClassFromFile(classloader.BootstrapCL, Global.StartingClass)
		if err != nil { // the exceptions message will already have been shown to user
			return shutdown.Exit(shutdown.JVM_EXCEPTION)
		}
	} else {
		_ = log.Log("Error: No executable program specified. Exiting.", log.INFO)
		ShowUsage(os.Stdout)
		return shutdown.Exit(shutdown.APP_EXCEPTION)
	}

	// the following was commented out per JACOBIN-327.
	// Likely to be reinstated at some later point
	// classloader.LoadReferencedClasses(mainClass)

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	classloader.MTableLoadNatives()

	// create the main thread
	MainThread = thread.CreateThread()

	// capture any panics and print diagnostic data
	defer func() int {
		if r := recover(); r != nil {
			showGoStackTrace(r)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	// begin execution
	_ = log.Log("Starting execution with: "+mainClass, log.INFO)
	status = StartExec(mainClass, &MainThread, &Global)

	if status != nil {
		return shutdown.Exit(shutdown.APP_EXCEPTION)
	}
	return shutdown.Exit(shutdown.OK)
}

// in the event of a panic, this routine explains that a panic occurred and
// (to a limited extent why) and then prints the Jacobin frame stack and then
// the golang stack trace. r is the error returned when the panic occurs
func showGoStackTrace(r any) {
	// show the event that caused the panic
	cause := fmt.Sprintf("%v", r)
	_ = log.Log("\nerror: go panic because of "+cause+"\n", log.SEVERE)

	// show the Jaocbin frame stack
	showFrameStack(&MainThread)
	_ = log.Log("\n", log.SEVERE)

	// capture the golang function stack and convert it to
	// a slice of strings
	stack := string(debug.Stack())
	entries := strings.Split(stack, "\n")

	// remove the strings showing the internals of golang's panic stack trace
	var i int
	for i = 0; i < len(entries); i++ {
		if strings.HasPrefix(entries[i], "panic") {
			i += 2 //
			break
		}
	}

	// print the remaining strings in the golang stack trace
	for {
		if i < len(entries) {
			_ = log.Log(entries[i], log.SEVERE)
			i += 1
		} else {
			break
		}
	}
}
