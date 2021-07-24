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

	// handle the command-line interface (cli) -- that is, process the args
	err := HandleCli(os.Args)
	if err != nil {
		if err.Error() == "end of processing" { // this is not an error but an end of processing
			shutdown(false)
		}
		shutdown(true)
	}

	Log("shutdown", FINE) // eventually move this to the shutdown func
	// shutdown(false)
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
