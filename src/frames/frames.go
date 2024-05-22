/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package frames

import (
	"container/list"
	"fmt"
	"jacobin/log"
	"unsafe"
)

// The data structures and functions related to JVM frames
type StackValue interface {
	int64 | float64 | unsafe.Pointer
}

var debugging bool = false

type Number interface {
	int64 | float64
}

func ftag(f *Frame) string {
	pp := fmt.Sprintf("%p\n", f)
	jj := len(pp) - 5 // show last 4 hex digits
	return pp[jj:]
}

// Frame is the fundamental execution environment for a single function/method call.
// Note that the operand stack (opStack) is made up of int64 items, rather than the JVM-
// prescribed 32-bit entries. The rationale is that longs and doubles can be stored
// without manipulation at this width. (However, there will still be need for the dummy
// second stack entry for these data items.
type Frame struct {
	Thread      int
	MethName    string        // method name
	MethType    string        // method type (signature)
	ClName      string        // class name
	Meth        []byte        // bytecode of method
	CP          interface{}   // will hold a *classloader.CPool (constant pool ptr) but due to circularity must be done this way
	Locals      []interface{} // local variables
	OpStack     []interface{} // operand stack
	TOS         int           // top of the operand stack
	PC          int           // program counter (index into the bytecode of the method)
	Ftype       byte          // type of method in frame: 'J' = java, 'G' = Golang, 'N' = native
	ExceptionPC int           // program counter at the moment the PC threw an exception
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
	// fram.OpStack = *new(opStack)

	if opStackSize < 0 { // TODO: Check if this is possible. If so, decide what to do. Class is clearly malformed.
		opStackSize = 0
	}

	// allocate the operand stack
	for j := 0; j < opStackSize; j++ {
		fram.OpStack = append(fram.OpStack, 0)
	}

	// set top of stack to an empty stack
	fram.TOS = -1

	fram.PC = 0
	fram.ExceptionPC = -1
	return &fram
}

// PushFrame pushes a frame. This simply adds a frame to the head of the list.
func PushFrame(fs *list.List, f *Frame) error {
	if debugging {
		fmt.Printf("DEBUG PushFrame %s ClName=%s, MethName=%s TOS=%d, PC=%d\n", ftag(f), f.ClName, f.MethName, f.TOS, f.PC)
	}
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

	if debugging {
		f := PeekFrame(fs, 0)
		fmt.Printf("DEBUG PopFrame %s ClName=%s, MethName=%s TOS=%d, PC=%d\n", ftag(f), f.ClName, f.MethName, f.TOS, f.PC)
	}

	fs.Remove(fs.Front())
	return nil
}

// PeekFrame peeks at a given frame without popping or deleting it.
// The current frame (so, top of stack) is 0, the one below it is 1, etc.
// Pass that value in and you receive back a pointer to the frame.
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
