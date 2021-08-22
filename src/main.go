/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"sync"
)

// var Global *globals.Globals
// Global := globals.Globals{}

// where everything begins
func main() {
	Global := globals.InitGlobals(os.Args[0])
	log.Init()

	// during development, let's use the most verbose logging level
	log.LogLevel = log.FINEST
	log.Log("running program: "+Global.JacobinName, log.FINE)

	// handle the command-line interface (cli) -- i.e., process the args
	LoadOptionsTable(Global)
	err := HandleCli(os.Args, Global)
	if err != nil {
		shutdown(true)
	}
	// some CLI options, like -version, show data and immediately exit. This tests for that.
	if Global.ExitNow == true {
		shutdown(false)
	}

	if Global.StartingClass == "" {
		log.Log("Error: No executable program specified. Exiting.", log.INFO)
		showUsage(os.Stdout)
		shutdown(true)
	} else {
		log.Log("Starting execution with: "+Global.StartingClass, log.INFO)
		classloader.AppCL.LoadClassFromFile(Global.StartingClass)
	}

	shutdown(false)
}

// the exit function. Later on, this will check a list of JVM shutdown hooks
// before closing down in order to have an orderly exit
func shutdown(errorCondition bool) {

	var mutex = sync.Mutex{}
	mutex.Lock()
	log.Log("shutdown", log.FINE)
	mutex.Unlock()

	if errorCondition {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
