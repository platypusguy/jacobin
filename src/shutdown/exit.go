/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package shutdown

import (
	"jacobin/globals"
	"jacobin/log"
	"os"
)

// The various flags that can be passed to the exit() function, reflecting
// the various reasons a shutdown is requested. (OK = normal end of program)
const (
	OK = iota
	JVM_EXCEPTION
	APP_EXCEPTION
	TEST
	UNKNOWN_ERROR
)

// Shutdown is the exit function. Later on, this will check a list of JVM Shutdown hooks
// before closing down in order to have an orderly exit
func Exit(errorCondition bool) int {
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
