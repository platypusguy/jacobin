/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-23 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"unsafe"
)

// With regard to the layout of a created object in Jacobin, note that
// on some architectures, but not Jacobin, there is an additional field
// that insures that the fields that follow the oops (the mark word and
// the class pointer) are aligned in memory for maximal performance.
type Object struct {
	Mark       MarkWord
	Klass      *string // the class name in the method area
	Fields     []Field // slice containing the fields
	FieldTable map[string]*Field
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

// Null is the Jacobin implementation of Java's null
var Null *Object = nil

// MakeEmptyObject() creates an empty basis Object. It is expected that other
// code will fill in the fields and the Klass field.
func MakeEmptyObject() *Object {
	o := Object{}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	o.Klass = &EmptyString // s/be filled in later, when class is filled in.

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]*Field)
	return &o
}

// determines whether a value is null or not
func IsNull(value any) bool {
	return value == nil || value == Null
}

// ToString dumps the contents of an object to a formatted multi-line string
func (objPtr *Object) ToString() string {
	var str string
	obj := *objPtr
	if obj.Klass != nil {
		str = *obj.Klass + "\n"
	} else {
		str = "class type: n/a \n"
	}

	if len(obj.FieldTable) > 0 {
		for key := range obj.FieldTable {
			str += fmt.Sprintf("\tFld: %s: (%s) %v\n", key, obj.FieldTable[key].Ftype, obj.FieldTable[key].Fvalue)
		}
	} else {
		for i, f := range obj.Fields {
			str += fmt.Sprintf("\tFld: %02d: (%s) %v\n", i, f.Ftype, f.Fvalue)
		}
	}

	return str
}
