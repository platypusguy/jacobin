/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021-2 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"os"
)

var Global globals.Globals

// where everything begins
func main() {
	Global = globals.InitGlobals(os.Args[0])
	log.Init()

	// during development, let's use the most verbose logging level
	// log.Level = log.FINEST  // no longer needed
	_ = log.Log("running program: "+Global.JacobinName, log.FINE)

	// handle the command-line interface (cli) -- i.e., process the args
	LoadOptionsTable(Global)
	err := HandleCli(os.Args, &Global)
	if err != nil {
		shutdown(true)
	}
	// some CLI options, like -version, show data and immediately exit. This tests for that.
	if Global.ExitNow == true {
		shutdown(false)
	}

	if Global.StartingClass == "" {
		_ = log.Log("Error: No executable program specified. Exiting.", log.INFO)
		showUsage(os.Stdout)
		shutdown(true)
	}

	// load the starting class, classes it references, and some base classes
	_ = classloader.Init()
	classloader.LoadBaseClasses(&Global)
	mainClass, err := classloader.LoadClassFromFile(classloader.BootstrapCL, Global.StartingClass)
	if err != nil { // the error message will already have been shown to user
		shutdown(true)
	}
	classloader.LoadReferencedClasses(mainClass)

	// begin execution
	_ = log.Log("Starting execution with: "+Global.StartingClass, log.INFO)
	if StartExec(mainClass, &Global) != nil {
		shutdown(true)
	}

	shutdown(false)
}

// the exit function. Later on, this will check a list of JVM shutdown hooks
// before closing down in order to have an orderly exit
func shutdown(errorCondition bool) int {
	globals.LoaderWg.Wait()
	g := globals.GetGlobalRef()

	err := errorCondition
	if log.Log("shutdown", log.INFO) != nil {
		err = true
	}

	if err {
		if g.JacobinName == "test" {
			return 1
		} else {
			os.Exit(1)
		}
	}

	if g.JacobinName == "test" {
		return 0
	} else {
		os.Exit(0)
	}
	return 0 // required by go
}
