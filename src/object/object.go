/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import "unsafe"

type Object struct {
	mark   MarkWord
	class  unsafe.Pointer // pointer to the loaded class
	fields []any          // slice containing the fields
}

type MarkWord struct {
	hash uint32 // contains hash code which is the lower 32 bits of the address
	misc uint32 // at present unused
}
