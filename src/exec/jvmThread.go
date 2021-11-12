/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import "fmt"

// Creates a JVM program execution thread. These threads are extremely limited.
// They basically hold a stack of frames. They push and pop frames as required.
// They begin execution; they exit when execution ends; and they emit diagnostic
// and performance data.

type execThread struct {
	id    int
	stack []frame
}

type frame struct {
}

func CreateThread() execThread {
	t := execThread{}
	return t
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
func (frameStack) init() frameStack {
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
func (frameStack) push(fs *frameStack, f frame) error {
	fs.top += 1
	if fs.top >= fs.size {
		fs.frames = append(fs.frames, f)
		fs.size += 1
	} else {
		fs.frames[fs.top] = f
	}
	return nil
}

// unlike most stacks, pop() here does not return an item. It simply
// decrements to the top of stack variable. Nothing is erased. Popping
// from an empty stack returns an error.
func (frameStack) pop(fs *frameStack) error {
	if fs.top == 0 {
		return fmt.Errorf("invalid pop of empty JVM frame stack")
	}
	fs.top -= 1
	return nil
}
