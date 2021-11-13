/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import "testing"

func TestFrameStack(t *testing.T) {
	fs := createFrameStack()
	if fs.top != 0 {
		t.Errorf("Newly allocated framestack. Expected top to be 0, got: %d", fs.top)
	}

	if fs.size != 10 {
		t.Errorf("Newly allocated framestack. Expected size to be 10, got: %d", fs.size)
	}
}

func TestFrameStackPushAndPop(t *testing.T) {
	fs := createFrameStack()
	pushFrame(&fs, frame{})
	if fs.top != 1 {
		t.Errorf("Pushed frame on to empty stack. Expected top to be 1, got: %d", fs.top)
	}

	popFrame(&fs)
	if fs.top != 0 {
		t.Errorf("Poppped only frame. Expected top to be 0, got: %d", fs.top)
	}

	if fs.size != 10 {
		t.Errorf("Expected framestack size to be 10, got: %d", fs.size)
	}

	// the test stack is empty at this point, so popFrame() should return an error
	if popFrame(&fs) == nil {
		t.Error("popFrame() on an empty frame stack did not generate an error.")
	}
}

//TODO: Push 11 items and make sure that the stack grows correctly.
