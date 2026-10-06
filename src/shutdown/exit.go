/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package shutdown

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/config"
	"jacobin/src/globals"
	"jacobin/src/prof"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"os"
	"runtime"
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

// Exit exits the JVM and returns the status code to the OS
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
		showStats()
		statics.DumpStatics("exit.Exit", statics.SelectUser, "")
	}

	if globals.TraceStats && !globals.TraceVerbose {
		showStats()
	}

	if errorCondition == TEST_OK {
		return 0
	} else if errorCondition == TEST_ERR {
		return 1
	}

	if errorCondition != OK {
		if !g.StrictJDK { // dump statics on error, unless in strict JDK mode
			statics.DumpStatics("exit.Exit", statics.SelectUser, "")
			_ = config.DumpConfig(os.Stderr)
		}
	}

	os.Stderr.Sync() // ensure all output is written before exiting
	prof.ExitToOS(errorCondition)

	return 0 // required by go
}

// showStats prints statistics about the JVM to stderr. We don't use trace.Trace() because
// it prints execution time, which is not useful here and disrupts formatting.
func showStats() {
	_, _ = fmt.Fprintf(os.Stderr, "\nString pool:   %6d entries\n", stringPool.GetStringPoolSize())
	_, _ = fmt.Fprintf(os.Stderr, "Classloader:   %6d classes loaded\n", classloader.MethAreaSize())
	_, _ = fmt.Fprintf(os.Stderr, "Statics Table: %6d statics accessed\n\n", len(statics.Statics))

	// Read full memory statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Memory currently in use by live data on the heap
	_, _ = fmt.Fprintf(os.Stderr, "Heap memory presently in use: %4d MB\n", m.Alloc/1024/1024)

	// Total memory requested and reserved from the OS
	_, _ = fmt.Fprintf(os.Stderr, "Total allocated memory:       %4d MB\n", m.Sys/1024/1024)
}
