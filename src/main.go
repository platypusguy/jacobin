/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"os"
)

var Global *Globals

// where everything begins
func main() {
	Global = initGlobals(os.Args[0])

	// during development, let's use the most verbose logging level
	Global.logLevel = FINEST
	Log("running program: "+Global.jacobinName, FINE)

	// handle the command-line interface (cli) -- i.e., process the args
	LoadOptionsTable(Global)
	err := HandleCli(os.Args)
	if err != nil {
		shutdown(true)
	}
	// some CLI options, like -version, show data and immediately exit. This tests for that.
	if Global.exitNow == true {
		shutdown(false)
	}

	if Global.startingClass == "" {
		Log("Error: No executable program specified. Exiting.", INFO)
		showUsage(os.Stdout)
		shutdown(true)
	} else {
		Log("Starting execution with: "+Global.startingClass, INFO)
	}

	shutdown(false)
}

// the exit function. Later on, this will check a list of JVM shutdown hooks
// before closing down in order to have an orderly exit
func shutdown(errorCondition bool) {
	Log("shutdown", FINE)
	if errorCondition {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
