/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"jacobin/classloader"
)

// With regard to the layout of a created object in Jacobin, note that
// on some architectures, but not Jacobin, there is an additional field
// that insures that the fields that follow the oops (the mark word and
// the class pointer) are aligned in memory for maximal performance.
type Object struct {
	Mark   MarkWord
	Klass  *classloader.Klass // pointer to the loaded class
	Fields []Field            // slice containing the fields
}

// These mark word contains values for different purposes. Here,
// we use the first four bytes for a hash value, which is taken
// from the address of the object. The 'misc' field will eventually
// contain other values, such as locking and monitoring items.
type MarkWord struct {
	Hash uint32 // contains hash code which is the lower 32 bits of the address
	Misc uint32 // at present unused
}

// We need to know the type of the field only to tell whether
// it occupies one or two slots on the stack when getfield and
// putfield bytecodes are executed.
type Field struct {
	Ftype  string // what type of value is stored in the field
	Fvalue any    // the actual value
}
