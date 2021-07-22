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
}

func showCopyright() {
	fmt.Println("Jacobin VM, v. 0.1.0, © 2021 by Andrew Binstock. All rights reserved. MPL 2.0 License.")
}
