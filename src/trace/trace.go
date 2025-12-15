/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package trace

// The principal logging function. Note it currently logs to stderr.
// At some future point, might allow the user to specify where logging should go.
import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"os"
	"sync"
	"time"
)

// Mutex for protecting the Log function during multithreading.
var mutex = sync.Mutex{}

// StartTime is the start time of this instance of the Jacoby VM.
var StartTime time.Time

// Identical to shutdown.UNKNOWN_ERROR (avoiding a cycle)
const UNKNOWN_ERROR = 5

var disabled = false

// Initialize the trace frame.
func Init() {
	StartTime = time.Now()
	disabled = false
}

// Disable the trace function. This is useful primarily in testing.
func Disable() {
	disabled = true
}

// Trace is the principal tracing function. Note that it currently
// writes to stderr. At some future point, this might become an option.
func Trace(argMsg string) {
	var err error

	if disabled {
		return
	}
	// if the message is more low-level than a WARNING,
	// prefix it with the elapsed time in millisecs.
	// check duration accuracy: time.Sleep(100 * time.Millisecond)
	duration := time.Since(StartTime)
	var millis = duration.Milliseconds()

	// Lock access to the logging stream to prevent inter-thread overwrite issues
	mutex.Lock()
	_, err = fmt.Fprintf(os.Stderr, "[%3d.%03ds] %s\n", millis/1000, millis%1000, argMsg)
	mutex.Unlock()
	if err != nil {
		errMsg := fmt.Sprintf("Trace: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// An error message is a prefix-decorated message that has no time-stamp.
func Error(argMsg string) {
	if disabled {
		return
	}

	var err error
	errMsg := "ERROR: " + argMsg
	mutex.Lock()
	_, err = fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	mutex.Unlock()
	if err != nil {
		errMsg = fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// Trace as-is.
// Useful for tracing stack traceback lines.
func AsIs(argMsg string) {
	if disabled {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()

	_, err := fmt.Fprintf(os.Stderr, "%s\n", argMsg)
	if err != nil {
		errMsg := fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// Similar to Error, except it's a warning, not an error.
func Warning(argMsg string) {
	if disabled {
		return
	}

	errMsg := "WARNING: " + argMsg
	mutex.Lock()
	_, err := fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	mutex.Unlock()
	if err != nil {
		errMsg = fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// Perform a minimal abort, which is a direct call to the global minimal abort function.
// Clearly, if trace is not working, then something is grievously wrong and the abort
// must be immediate.
func rawAbort(whichException int, msg string) {
	globals.GetGlobalRef().FuncMinimalAbort(whichException, msg)
}
