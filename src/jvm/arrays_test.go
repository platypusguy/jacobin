/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/frames"
	"jacobin/globals"
	"testing"
)

func TestJdkArrayTypeToJacobinType(t *testing.T) {

	a := jdkArrayTypeToJacobinType(T_BOOLEAN)
	if a != BYTE {
		t.Errorf("Expected Jacobin type of %d, got: %d", BYTE, a)
	}

	b := jdkArrayTypeToJacobinType(T_CHAR)
	if b != INT {
		t.Errorf("Expected Jacobin type of %d, got: %d", INT, b)
	}

	c := jdkArrayTypeToJacobinType(T_DOUBLE)
	if c != FLOAT {
		t.Errorf("Expected Jacobin type of %d, got: %d", FLOAT, c)
	}

	d := jdkArrayTypeToJacobinType(999)
	if d != ERROR {
		t.Errorf("Expected Jacobin type of %d, got: %d", ERROR, d)
	}
}

// NEWARRAY: creationg of array for primitive values
func TestNewrray(t *testing.T) {
	f := newFrame(NEWARRAY)
	push(&f, int64(13))             // make the array 13 elements big
	f.Meth = append(f.Meth, T_LONG) // make it an array of longs

	globals.InitGlobals("test")
	// globals.InitArrayAddressList()  // init address list in globals

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)
	if f.TOS != 0 {
		t.Errorf("Top of stack, expected 0, got: %d", f.TOS)
	}

	// did we capture the address of the new array in globals?
	g := globals.GetGlobalRef()
	if g.ArrayAddressList.Len() != 1 {
		t.Errorf("Expecting array address list to have length 1, got %d",
			g.ArrayAddressList.Len())
	}

	// now, test the length of the array, which should be 13
	element := g.ArrayAddressList.Front()
	val := element.Value.([]int64)
	if len(val) != 13 {
		t.Errorf("Expecting array length of 13, got %d", len(val))
	}
}
