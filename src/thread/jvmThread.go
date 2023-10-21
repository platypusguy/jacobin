/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package thread

import (
	"container/list"
	"jacobin/globals"
)

// Creates a JVM program execution thread. These threads are extremely limited.
// They basically hold a Stack of frames. They push and popFrame frames as required.
// They begin execution; they exit when execution ends; and they emit diagnostic
// and performance data.

type ExecThread struct {
	ID    int        // the thread ID
	Stack *list.List // the JVM Stack (frame stack, that is) for this thread
	PC    int        // the program counter (the index to the instruction being executed)
	Trace bool       // do we Trace instructions?
}

// CreateThread creates an execution thread and initializes it with default values
// All Jacobin execution threads *must* use this function to create a thread
func CreateThread() ExecThread {
	t := ExecThread{}
	t.ID = incrementThreadNumber()
	t.Stack = nil
	t.PC = 0
	t.Trace = false
	return t
}

// Adds a thread to the global thread table using the ID as the key,
// and a pointer to the ExecThread as the value
func (t *ExecThread) AddThreadToTable() {
	glob := globals.GetGlobalRef()

	glob.ThreadLock.Lock()
	glob.Threads[t.ID] = t
	glob.ThreadLock.Unlock()
}

// func AddThreadToTable(t *ExecThread, tbl *globals.ThreadList) int {
// 	tbl.ThreadsMutex.Lock()
//
// 	tbl.ThreadsList.PushBack(t)
// 	t.ID = tbl.ThreadsList.Len() - 1
// 	tbl.ThreadsMutex.Unlock()
//
// 	return t.ID
// }

// threads are assigned a monotonically incrementing integer ID. This function
// increments the counter and returns its value as the integer ID to use
func incrementThreadNumber() int {
	glob := globals.GetGlobalRef()
	glob.ThreadLock.Lock()
	glob.ThreadNumber += 1
	glob.ThreadLock.Unlock()
	return glob.ThreadNumber
}
