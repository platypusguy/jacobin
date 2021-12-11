/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import (
	"container/list"
	"fmt"
	"jacobin/log"
)

// The data structures and functions related to JVM frames

// Note that the operand stack (opStack) is made up of int64 items, rather than the JVM-
// prescribed 32-bit entries. The rationale is that longs and doubles can be stored
// without manipulation at this width. (However, there will still be need for the dummy
// second stack entry for these data items.
type frame struct {
	thread   int
	methName string  // method name
	clName   string  // class name
	meth     []byte  // bytecode of method
	cp       *CPool  // constant pool of class
	locals   []int64 // local variables
	opStack  []int64 // operand stack
	tos      int     // top of the operand stack
	pc       int     // program counter (index into the bytecode of the method)
	ftype    byte    // type of method in frame: 'J' = java, 'G' = Golang, 'N' = native
}

// a stack of frames. Implemented as a list in which the current running
// frame is always the frame at the head
func createFrameStack() *list.List {
	l := list.New()
	return l
}

// creates a raw frame and allocates an opStack of the passed-in size.
func createFrame(opStackSize int) *frame {
	fram := frame{}

	// allocate the operand stack
	for j := 0; j < opStackSize; j++ {
		fram.opStack = append(fram.opStack, int64(0))
	}

	// set top of stack to an empty stack
	fram.tos = -1
	fram.pc = 0
	return &fram
}

// push a frame. This simply adds a frame to the head of the list.
func pushFrame(fs *list.List, f *frame) error {
	fs.PushFront(f)
	// TODO: move this to instrumentation system
	if log.Level == log.FINEST {
		var s string
		for e := fs.Front(); e != nil; e = e.Next() {
			fr := e.Value.(*frame)
			s = s + "\n" + "> " + fr.methName
		}
		_ = log.Log("Present stack frame:"+s, log.FINEST)
	}
	return nil
}

// deletes the frame at the head of the list.
func popFrame(fs *list.List) error {
	if fs.Len() == 0 {
		return fmt.Errorf("invalid popFrame of empty JVM frame stack")
	}

	fs.Remove(fs.Front())
	return nil
}
