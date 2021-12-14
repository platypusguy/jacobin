/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import "testing"

func TestFrameStack(t *testing.T) {
	fs := createFrameStack()
	if fs.Len() != 0 {
		t.Errorf("Newly allocated framestack. Expected size to be 0, got: %d", fs.Len())
	}
}

func TestFrameStackPushAndPop(t *testing.T) {
	fs := createFrameStack()
	_ = pushFrame(fs, &frame{})
	if fs.Len() != 1 {
		t.Errorf("Pushed frame on to empty stack. Expected size of stack to be 1, got: %d", fs.Len())
	}

	_ = popFrame(fs)
	if fs.Len() != 0 {
		t.Errorf("Poppped only frame. Expected size of stack to be 0, got: %d", fs.Len())
	}

	// the test stack is empty at this point, so popFrame() should return an error
	if popFrame(fs) == nil {
		t.Error("popFrame() on an empty frame stack did not generate an error.")
	}
}
