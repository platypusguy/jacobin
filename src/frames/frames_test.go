/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package frames

import "testing"

func TestNewFrame(t *testing.T) {
	f := CreateFrame(6)
	if len(f.OpStack) != 6 || f.TOS != -1 || f.PC != 0 {
		t.Error("Created frame is invalid")
	}
}

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

func TestFramePeek(t *testing.T) {
	fs := CreateFrameStack()
	f1 := CreateFrame(1)
	f2 := CreateFrame(2)
	_ = PushFrame(fs, f1)
	_ = PushFrame(fs, f2)

	peek := PeekFrame(fs, 1)
	if len(peek.OpStack) != 1 {
		t.Errorf("Peeked at prior frame. Expected size of opstack to be 1, got: %d", len(peek.OpStack))
	}
}
