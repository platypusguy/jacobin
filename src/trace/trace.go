/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package trace

// The principal logging function. Note it currently logs to stderr.
// At some future point, might allow the user to specify where logging should go.
import (
	"fmt"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/shutdown"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Mutex for protecting the Log function during multithreading.
var mutex = sync.Mutex{}

// StartTime is the start time of this instance of the Jacoby VM.
var StartTime time.Time

// Initialize the trace frame.
func Init() {
	StartTime = time.Now()
}

// Trace is the principal tracing function. Note that it currently
// writes to stderr. At some future point, this might become an option.
func Trace(msg string) {

	var err error

	if len(msg) == 0 {
		errMsg := fmt.Sprintf("Zero-length trace argument")
		abruptEnd(excNames.IllegalArgumentException, errMsg)
		return
	}

	// if the message is more low-level than a WARNING,
	// prefix it with the elapsed time in millisecs.
	duration := time.Since(StartTime)
	var millis = duration.Milliseconds()

	// Lock access to the logging stream to prevent inter-thread overwrite issues
	mutex.Lock()
	_, err = fmt.Fprintf(os.Stderr, "[%3d.%03ds] %s\n", millis/1000, millis%1000, msg)
	if err != nil {
		errMsg := fmt.Sprintf("*** stderr failed, err: %v", err)
		abruptEnd(excNames.IOError, errMsg)
	}
	mutex.Unlock()

}

// Duplicated from minimalAbort in the exceptions package exceptions.go due to a Go-diagnosed cycle.
func abruptEnd(whichException int, msg string) {
	var stack string
	bytes := debug.Stack()
	if len(bytes) > 0 {
		stack = string(bytes)
	} else {
		stack = ""
	}
	glob := globals.GetGlobalRef()
	glob.ErrorGoStack = stack
	errMsg := fmt.Sprintf("%s: %s", excNames.JVMexceptionNames[whichException], msg)
	_, _ = fmt.Fprintln(os.Stderr, errMsg)
	// exceptions.ShowGoStackTrace(nil) <---------------- causes Go-diagnosed cycle: classloader >exceptions > classloader
	_ = shutdown.Exit(shutdown.APP_EXCEPTION)
	// os.Exit(1)
}
