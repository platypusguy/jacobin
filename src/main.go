/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"jacobin/src/jvm"
	"jacobin/src/prof"
	"os"
)

func main() {
	path := os.Getenv("JACOBIN_CPUPROFILE")
	if path != "" {
		prof.StartProfiling(path)
	}
	jvm.JVMrun()
	prof.StopProfiling()
}
