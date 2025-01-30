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

// Initialize the trace frame.
func Init() {
	StartTime = time.Now()
}

// Trace is the principal tracing function. Note that it currently
// writes to stderr. At some future point, this might become an option.
func Trace(argMsg string) {

	var err error

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
		abruptEnd(excNames.IOError, errMsg)
	}
}

// An error message is a prefix-decorated message that has no time-stamp.
func Error(argMsg string) {
	var err error
	errMsg := "ERROR: " + argMsg
	mutex.Lock()
	_, err = fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	mutex.Unlock()
	if err != nil {
		errMsg := fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		abruptEnd(excNames.IOError, errMsg)
	}
}

// Similar to Error, except it's a warning, not an error.
func Warning(argMsg string) {
	errMsg := "WARNING: " + argMsg
	mutex.Lock()
	_, err := fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	mutex.Unlock()
	if err != nil {
		errMsg := fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		abruptEnd(excNames.IOError, errMsg)
	}
}

// Duplicated from minimalAbort in the exceptions package exceptions.go due to a Go-diagnosed cycle.
func abruptEnd(whichException int, msg string) {
	globals.GetGlobalRef().FuncMinimalAbort(whichException, msg)
	/*
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
		// exceptions.ShowGoStackTrace(nil) // <------------ causes Go-diagnosed cycle: classloader > exceptions > classloader
		// _ = shutdown.Exit(shutdown.APP_EXCEPTION) <--- causes Go-diagnosed cycle: shutdown > statics > trace > shutdown
		os.Exit(UNKNOWN_ERROR)
	*/
}
