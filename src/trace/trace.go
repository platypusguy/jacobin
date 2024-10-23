/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package trace

// The principal logging function. Note it currently logs to stderr.
// At some future point, might allow the user to specify where logging should go.
import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Should never see an indication of an empty trace message!
const EmptyMsg = "*** EMPTY LOGGING MESSAGE !!!"

// Mutex for protecting the Log function during multithreading.
var mutex = sync.Mutex{}

// StartTime is the start time of this instance of the Jacoby VM.
var StartTime time.Time
var okStderr bool

// Initialize the trace frame.
func Init() {
	StartTime = time.Now()
	okStderr = true
}

// Trace is the principal tracing function. Note that it currently
// writes to stderr. At some future point, this might become an option.
func Trace(msg string) {

	var err error

	if len(msg) == 0 {
		msg = EmptyMsg
	}

	// if the message is more low-level than a WARNING,
	// prefix it with the elapsed time in millisecs.
	duration := time.Since(StartTime)
	var millis = duration.Milliseconds()

	// Lock access to the logging stream to prevent inter-thread overwrite issues
	if okStderr {
		mutex.Lock()
		_, err = fmt.Fprintf(os.Stderr, "[%3d.%03ds] %s\n", millis/1000, millis%1000, msg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[%3d.%03ds] *** stderr failed, err: %v\n", millis/1000, millis%1000, err)
			okStderr = false
		}
		mutex.Unlock()
	}

}
