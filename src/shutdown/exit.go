/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
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
type ExitStatus = int

const (
	OK ExitStatus = iota
	JVM_EXCEPTION
	APP_EXCEPTION
	TEST_OK
	TEST_ERR
	UNKNOWN_ERROR
)

// Shutdown is the exit function. Later on, this will check a list of JVM Shutdown hooks
// before closing down in order to have an orderly exit
func Exit(errorCondition ExitStatus) int {
	globals.LoaderWg.Wait()
	g := globals.GetGlobalRef()
	if g.JacobinName == "test" {
		if errorCondition == OK {
			errorCondition = TEST_OK
		} else {
			errorCondition = TEST_ERR
		}
	}

	if log.Log("shutdown", log.INFO) != nil {
		errorCondition = UNKNOWN_ERROR
	}

	if errorCondition == TEST_OK {
		return 0
	} else if errorCondition == TEST_ERR {
		return 1
	}

	if errorCondition == OK {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

	return 0 // required by go
}
