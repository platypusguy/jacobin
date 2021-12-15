/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package log

// The principal logging function. Note it currently logs to stderr.
// At some future point, might allow the user to specify where logging should go.
import (
	"errors"
	"fmt"
	"jacobin/globals"
	"os"
	"sync"
	"time"
)

// the various logging levels (Note that higher numbers means more granular)
const (
	SEVERE = iota + 1
	WARNING
	CLASS
	INFO
	FINE
	FINEST
	TRACE_INST
)

// Level is the level the logger currently supports. See the enums above.
var Level int

// Mutex for protecting the Log function during multithreading.
var mutex = sync.Mutex{}

// StartTime is the start time of this instance of the Jacoby VM.
var StartTime time.Time

// Init initialize the logger, which by default is set to WARNING. Note: that it cannot be
// set any coarser. At all times, SEVERE and WARNING messages must be visible to the user.
func Init() {
	Level = WARNING
	StartTime = time.Now()
}

// Log is the principal logging function. Note that it currently
// logs to stderr. At some future point, this might become an option.
func Log(msg string, level int) (err error) {
	if len(msg) == 0 {
		return errors.New("empty logging message")
	}

	if level < SEVERE || level > TRACE_INST {
		return errors.New("invalid logging level")
	}

	// if the message is a trace and we're not tracing, then return.
	if level == TRACE_INST && globals.GetInstance().Options["-trace"].Set != true {
		return
	}

	// if the message is for a finer logging level than currently being logged, return.
	if level <= FINEST && level > Level {
		return
	}

	// if the message is more low-level than a WARNING,
	// prefix it with the elapsed time in millisecs.
	duration := time.Since(StartTime)
	var millis = duration.Milliseconds()

	// lock the write to the logging stream to prevent overwrite issues
	// if some other operation is also writing to the stream
	mutex.Lock()
	if level > WARNING { // show elapsed time only if messages are finer than warning
		_, _ = fmt.Fprintf(os.Stderr, "[%3d.%03ds] ", millis/1000, millis%1000)
	}
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", msg)
	mutex.Unlock()
	return
}

// SetLogLevel seta the level of granularity.
func SetLogLevel(level int) (err error) {
	// SEVERE is here just to fill the hierarchy. You cannot actually set the logging
	// level coarser than WARNING. In other words, all warnings must be shown.
	if level <= SEVERE || level > TRACE_INST {
		return errors.New("invalid logging level")
	}

	Level = level
	return nil
}
