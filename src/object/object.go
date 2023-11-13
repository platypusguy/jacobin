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
	"strings"
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
	Fvalue any    // the actual value or a pointer to the value (ftype="[something)
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

func fmtHelper(klassString string, field Field) string {
	if klassString == filepath.FromSlash(StringClassName) {
		bytes := *field.Fvalue.(*[]byte)
		//fmt.Printf("DEBUG fmtHelper bytes: %d % x\n", len(bytes), bytes)
		last := len(bytes) - 1
		if last < 0 {
			return "\"\""
		}
		if bytes[last] == '\n' {
			bytes = bytes[0:last]
		}
		return fmt.Sprintf("\"%s\"", string(bytes))
	}
	switch field.Ftype {
	case types.Double, types.Float, types.Static + types.Double, types.Static + types.Float:
		return fmt.Sprintf("%f", field.Fvalue)
	case types.Int, types.Long, types.Short, types.Static + types.Int, types.Static + types.Long, types.Static + types.Short:
		return fmt.Sprintf("%d", field.Fvalue)
	case types.Byte, types.Static + types.Byte:
		return fmt.Sprintf("%02x", field.Fvalue)
	case types.Bool, types.Static + types.Bool:
		// TODO: Why does FieldTable[key] pass an int64 YET Fields[index] passes a bool???
		switch field.Fvalue.(type) {
		case bool:
			if field.Fvalue.(bool) {
				return "true"
			} else {
				return "false"
			}
		case int64:
			if field.Fvalue.(int64) != 0 {
				return "true"
			} else {
				return "false"
			}
		default:
			return fmt.Sprintf("<ERROR Ftype=bool but unexpected Fvalue variable type: %T !>", field.Fvalue)
		}
	case types.Char, types.Static + types.Char:
		return fmt.Sprintf("%q", field.Fvalue)
	case types.ByteArray, types.Static + types.ByteArray:
		fvalue := field.Fvalue
		if fvalue == nil {
			return "<ERROR nil Fvalue!>"
		}
		bytesPtr := fvalue.(*[]byte)
		if bytesPtr == nil {
			return "<ERROR nil byte array ptr!>"
		}
		if len(*bytesPtr) < 1 {
			return "<nil byte array>"
		}
		return fmt.Sprintf("% x", *bytesPtr)
	}

	// Default action:
	return fmt.Sprintf("%v", field.Fvalue)
}

// FormatField creates a string that represents a single field of an Object.
func (objPtr *Object) FormatField() string {
	var output string
	var klassString string // string class name
	obj := *objPtr         // whole object
	key := "value"         // key to the FieldTable map

	if obj.Klass != nil {
		klassString = *obj.Klass
	} else {
		klassString = "<ERROR nil class pointer!>" // Why is there no class name pointer for this object?
		obj.DumpObject(klassString, 0)
		return klassString
	}

	if len(obj.FieldTable) > 0 {
		// Using key="value" in the FieldTable
		ptr := obj.FieldTable[key]
		if ptr == nil {
			str := fmt.Sprintf("<ERROR FieldTable[\"%s\"] not found!>", key)
			obj.DumpObject(str, 0)
			return str
		}
		field := *ptr
		str := fmtHelper(klassString, field)
		output = fmt.Sprintf("%s: (%s) %s", key, obj.FieldTable[key].Ftype, str)
		if strings.HasPrefix(str, "<ERROR") {
			obj.DumpObject(str, 0)
		}
	} else {
		// Using [0] in the Fields slice
		if len(obj.Fields) > 0 {
			field := obj.Fields[0]
			str := fmtHelper(klassString, field)
			output = fmt.Sprintf("(%s) %s", obj.Fields[0].Ftype, str)
			if strings.HasPrefix(str, "<ERROR") {
				obj.DumpObject(str, 0)
			}
		} else {
			// Field table and field slice are both empty.
			output = "<ERROR field empty!>"
			obj.DumpObject(output, 0)
		}
	}

	return output
}

// DumpObject displays every attribute of an Object, formatted as multi-line printed output.
// 3 sections:
// * Class name
// * Field table
// * Field slice
func (objPtr *Object) DumpObject(title string, indent int) {
	obj := *objPtr
	output := ""
	var klassString string

	// Emit BEGIN
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	output += "DumpObject " + title + " {\n"

	// Emit klass line
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	if obj.Klass != nil {
		klassString = "\tClass: " + *obj.Klass
	} else {
		klassString = "\t<class MISSING!>"
	}
	output += klassString + "\n"

	// Emit FieldTable.
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	nflds := len(obj.FieldTable)
	if nflds > 0 {
		output += fmt.Sprintf("\tField Table (%d):\n", nflds)
		for key := range obj.FieldTable {
			if indent > 0 {
				output += strings.Repeat(" ", indent)
			}
			ptr := obj.FieldTable[key]
			if ptr == nil {
				output += fmt.Sprintf("\t\t<ERROR nil FieldTable[%s] ptr!>\n", key)
			} else {
				output += fmt.Sprintf("\t\tFld %s: (%s) %s\n", key, obj.FieldTable[key].Ftype, fmtHelper(klassString, *obj.FieldTable[key]))
			}
		}
	} else {
		output += fmt.Sprintf("\tField Table is <empty>\n")
	}

	// Emit Fields slice.
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	nflds = len(obj.Fields)
	if nflds > 0 {
		output += fmt.Sprintf("\tField Slice (%d):\n", nflds)
		for _, fld := range obj.Fields {
			output += fmt.Sprintf("\t\tFld (%s) %s\n", fld.Ftype, fmtHelper(klassString, fld))
		}
	} else {
		output += "\tField Slice is <empty>\n"
	}

	// Emit END
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	output += "}\n"

	// Print output all at once.
	fmt.Print(output)
}
