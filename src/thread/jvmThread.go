/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
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

func CreateThread() ExecThread {
	t := ExecThread{}
	t.ID = 0
	t.PC = 0
	t.Stack = nil
	t.Trace = false
	return t
}

func AddThreadToTable(t *ExecThread, tbl *globals.ThreadList) int {
	lock := *tbl.ThreadsMutex
	lock.Lock()

	tbl.ThreadsList.PushBack(t)
	t.ID = tbl.ThreadsList.Len() - 1
	lock.Unlock()

	return t.ID
}
