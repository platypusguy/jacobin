/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package shutdown

import (
	"fmt"
	"jacobin/src/config"
	"jacobin/src/globals"
	"jacobin/src/statics"
	"jacobin/src/trace"
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

// This is the exit-to-O/S function.
// TODO: Check a list of JVM Shutdown hooks before closing down in order to have an orderly exit.
func Exit(errorCondition ExitStatus) int {
	globals.LoaderWg.Wait()
	g := globals.GetGlobalRef()
	if g.JacobinName == "test" || g.JacobinName == "testWithoutShutdown" {
		if errorCondition == OK {
			errorCondition = TEST_OK
		} else {
			errorCondition = TEST_ERR
		}
	}

	if globals.TraceVerbose {
		msg := fmt.Sprintf("shutdown.Exit(%d) requested", errorCondition)
		trace.Trace(msg)
	}

	if errorCondition == TEST_OK {
		return 0
	} else if errorCondition == TEST_ERR {
		return 1
	}

	if errorCondition != OK {
		statics.DumpStatics("exit.Exit", statics.SelectUser, "")
		config.DumpConfig(os.Stderr)
	}
	
	os.Stderr.Sync() // ensure all output is written before exiting
	os.Exit(errorCondition)

	return 0 // required by go
}
