/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

// Creates a JVM program execution thread. These threads are extremely limited.
// They basically hold a stack of frames. They push and popFrame frames as required.
// They begin execution; they exit when execution ends; and they emit diagnostic
// and performance data.

type execThread struct {
	id    int        // the thread ID
	stack frameStack // the JVM stack for this thread
	pc    int        // the program counter (the index to the instruction being executed)
}

func CreateThread(threadNum int) execThread {
	t := execThread{}
	t.id = threadNum
	t.pc = 0
	t.stack = createFrameStack()
	return t
}
