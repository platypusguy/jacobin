/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import (
	"fmt"
)

// The data structures and functions related to JVM frames

type frame struct {
	thread   int
	methName string  // method name
	clName   string  // class name
	meth     []byte  // bytecode of method
	cp       *CPool  // constant pool of class
	locals   []int32 // local variables
	opStack  []int32 // operand stack
	tos      int     // top of the operand stack
	ftype    byte    // type of method in frame: 'J' = java, 'G' = Golang, 'N' = native
}

// a stack of frames. Top points to the present top of the stack.
// When top == 0, the stack is empty. (In other words, there is no
// zero entry in this stack. size is the total number of allocated
// slots in the stack. When a push() forces top to go past size, a
// new slot is allocated.
type frameStack struct {
	frames []frame
	top    int
	size   int
}

// we preallocate 10 frames for this stack. If more are needed,
// they'll be appended in push()
func createFrameStack() frameStack {
	fs := frameStack{}
	for i := 0; i < 10; i++ {
		fs.frames = append(fs.frames, frame{})
	}
	fs.top = 0
	fs.size = 10
	return fs
}

// push a frame. If more frames need to be allocated to the frame stack,
// then append one for the new frame. Always returns nil at present, but see:
// TODO: test for out of memory error.
func pushFrame(fs *frameStack, f frame) error {
	fs.top += 1
	if fs.top >= fs.size {
		fs.frames = append(fs.frames, f)
		fs.size += 1
	} else {
		fs.frames[fs.top] = f
	}
	return nil
}

// unlike most stacks, popFrame() here does not return an item. It simply
// decrements to the top of stack variable. Nothing is erased. Popping
// from an empty stack returns an error.
func popFrame(fs *frameStack) error {
	if fs.top == 0 {
		return fmt.Errorf("invalid popFrame of empty JVM frame stack")
	}
	fs.top -= 1
	return nil
}
