/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"fmt"
	"os"
)

var Global *Globals

// where everything begins
func main() {
	showCopyright()
	Global = initGlobals(os.Args[0])

	// during development, let's use the most verbose logging level
	Global.logLevel = FINEST
	Log("running program: "+Global.jacobinName, FINE)

	// handle the command-line interface (cli) -- that is, process the args
	err := HandleCli()
	if err != nil {
		closedown(true)
	}

	closedown(false)
}

// the exit function. Later on, this will check a list of JVM shutdown hooks
// before closing down in order to have an orderly exit
func closedown(errorCondition bool) {
	if errorCondition {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func showCopyright() {
	fmt.Println("Jacobin VM, v. 0.1.0, © 2021 by Andrew Binstock. All rights reserved. MPL 2.0 License.")
}
