/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package frames

import "testing"

func TestFrameStack(t *testing.T) {
	fs := CreateFrameStack()
	if fs.Len() != 0 {
		t.Errorf("Newly allocated framestack. Expected size to be 0, got: %d", fs.Len())
	}
}

func TestFrameStackPushAndPop(t *testing.T) {
	fs := CreateFrameStack()
	_ = PushFrame(fs, &Frame{})
	if fs.Len() != 1 {
		t.Errorf("Pushed frame on to empty stack. Expected size of stack to be 1, got: %d", fs.Len())
	}

	_ = PopFrame(fs)
	if fs.Len() != 0 {
		t.Errorf("Poppped only frame. Expected size of stack to be 0, got: %d", fs.Len())
	}

	// the test stack is empty at this point, so PopFrame() should return an error
	if PopFrame(fs) == nil {
		t.Error("PopFrame() on an empty frame stack did not generate an error.")
	}
}
