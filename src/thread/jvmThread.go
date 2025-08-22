/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package thread

import (
	"container/list"
	"jacobin/src/globals"
)

// Creates a JVM program execution thread. These threads are extremely limited.
// They basically hold a Stack of frames. They push and popFrame frames as required.
// They begin execution; they exit when execution ends.

type ExecThread struct {
	ID    int        // the thread ID
	Stack *list.List // the JVM Stack (frame stack, that is) for this thread
	Trace bool       // do we trace instructions?
}

// CreateThread creates an execution thread and initializes it with default values
// All Jacobin execution threads *must* use this function to create a thread
func CreateThread() ExecThread {
	t := ExecThread{}
	t.ID = IncrementThreadNumber()
	t.Stack = nil
	t.Trace = false
	return t
}

// Adds a thread to the global thread table using the ID as the key,
// and a pointer to the ExecThread as the value
func (t *ExecThread) AddThreadToTable(glob *globals.Globals) {
	glob.ThreadLock.Lock()
	glob.Threads[t.ID] = t
	glob.ThreadLock.Unlock()
}

// threads are assigned a monotonically incrementing integer ID. This function
// increments the counter and returns its value as the integer ID to use
func IncrementThreadNumber() int {
	glob := globals.GetGlobalRef()

	glob.ThreadLock.Lock()
	forCaller := glob.ThreadNumber + 1 // ensure that caller sees this one
	glob.ThreadNumber = forCaller
	glob.ThreadLock.Unlock() // I don't care if glob.ThreadNumber races ahead
	return forCaller
}

// ======= Items for Java threads ======
// Thread state constants matching Java's Thread.State enum
type State int

const (
	NEW           State = 0
	RUNNABLE            = 1
	BLOCKED             = 2
	WAITING             = 3
	TIMED_WAITING       = 4
	TERMINATED          = 5
)
