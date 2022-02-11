/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package frames

import (
	"container/list"
	"fmt"
	"jacobin/classloader"
	"jacobin/log"
)

// The data structures and functions related to JVM frames

// Frame is the fundamental execution environment for a single function/method call.
// Note that the operand stack (opStack) is made up of int64 items, rather than the JVM-
// prescribed 32-bit entries. The rationale is that longs and doubles can be stored
// without manipulation at this width. (However, there will still be need for the dummy
// second stack entry for these data items.
type Frame struct {
	Thread   int
	MethName string             // method name
	ClName   string             // class name
	Meth     []byte             // bytecode of method
	CP       *classloader.CPool // constant pool of class
	Locals   []int64            // local variables
	OpStack  []int64            // operand stack
	TOS      int                // top of the operand stack
	PC       int                // program counter (index into the bytecode of the method)
	Ftype    byte               // type of method in frame: 'J' = java, 'G' = Golang, 'N' = native
}

// CreateFrameStack creates a stack of frames. Implemented as a list in which
// the current running frame is always the frame at the head
func CreateFrameStack() *list.List {
	l := list.New()
	return l
}

// CreateFrame creates a raw frame and allocates an opStack of the passed-in size.
func CreateFrame(opStackSize int) *Frame {
	fram := Frame{}

	// allocate the operand stack
	for j := 0; j < opStackSize; j++ {
		fram.OpStack = append(fram.OpStack, int64(0))
	}

	// set top of stack to an empty stack
	fram.TOS = -1
	fram.PC = 0
	return &fram
}

// PushFrame pushes a frame. This simply adds a frame to the head of the list.
func PushFrame(fs *list.List, f *Frame) error {
	fs.PushFront(f)
	// TODO: move this to instrumentation system
	if log.Level == log.FINEST {
		var s string
		for e := fs.Front(); e != nil; e = e.Next() {
			fr := e.Value.(*Frame)
			s = s + "\n" + "> " + fr.MethName
		}
		_ = log.Log("Present stack frame:"+s, log.FINEST)
	}
	return nil
}

// PopFrame deletes the frame at the head of the list.
func PopFrame(fs *list.List) error {
	if fs.Len() == 0 {
		return fmt.Errorf("invalid PopFrame of empty JVM frame stack")
	}

	fs.Remove(fs.Front())
	return nil
}

// PeekFrame peeks at a given frame without popping or deleting it.
func PeekFrame(fs *list.List, which int) *Frame {
	var e *list.Element
	i := 0
	for e = fs.Front(); e != nil; e = e.Next() {
		if i == which {
			break
		} else {
			i += 1
		}
	}
	return e.Value.(*Frame)
}
