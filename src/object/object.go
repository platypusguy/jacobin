/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-23 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"jacobin/types"
	"path/filepath"
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

func toStringHelper(klassString string, field Field) string {
	if klassString == filepath.FromSlash("java/lang/String") {
		return fmt.Sprintf("%s", *field.Fvalue.(*[]byte))
	}
	switch field.Ftype {
	case types.Double, types.Float:
		return fmt.Sprintf("%f", field.Fvalue)
	case types.Int, types.Long, types.Short:
		return fmt.Sprintf("%d", field.Fvalue)
	case types.Byte:
		return fmt.Sprintf("%02x", field.Fvalue)
	case types.Bool:
		return fmt.Sprintf("%t", field.Fvalue)
	case types.Char:
		return fmt.Sprintf("%q", field.Fvalue)
	case "Ljava/lang/String;":
		return field.Fvalue.(string)
	case types.ByteArray:
		return fmt.Sprintf("% x", *field.Fvalue.(*[]byte))
	}

	return fmt.Sprintf("%v", field.Fvalue)
}

// ToString dumps the contents of an object to a formatted multi-line string
func (objPtr *Object) ToString() string {
	var str string
	obj := *objPtr
	klassString := *obj.Klass
	if obj.Klass != nil {
		str = klassString + "\n"
	} else {
		str = "class type: n/a \n"
	}

	if len(obj.FieldTable) > 0 {
		for key := range obj.FieldTable {
			str += fmt.Sprintf("\tFld: %s: (%s) %s\n", key, obj.FieldTable[key].Ftype, toStringHelper(klassString, *obj.FieldTable[key]))
		}
	} else {
		//for i, field := range obj.Fields {
		//	str += fmt.Sprintf("\tFld: %02d: (%s) %s\n", i, field.Ftype, toStringHelper(field))
		//}
		str += fmt.Sprintf("\tFld:(%s) %s", obj.Fields[0].Ftype, toStringHelper(klassString, obj.Fields[0]))
	}

	return str
}
