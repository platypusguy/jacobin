/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-23 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"jacobin/statics"
	"jacobin/types"
	"strings"
	"unsafe"
)

var tracing = false

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

// Make an object for a Java primitive field (byte, int, etc.), given the class and field type.
func MakePrimitiveObject(classString string, ftype string, arg any) *Object {
	objPtr := MakeEmptyObject()
	(*objPtr).Klass = &classString
	var field Field
	field.Ftype = ftype
	field.Fvalue = arg
	(*objPtr).Fields = append((*objPtr).Fields, field)
	return objPtr
}

// determines whether a value is null or not
func IsNull(value any) bool {
	return value == nil || value == Null
}

// Return a string representing the field type and value.
// Input for all fields:
//
//	field: structure of field type and field value if not static.
//
// Input for static fields:
//
//	fieldName: Key to the jacobin statics table.
//	frameStack: Required for getStaticValue() if the field is static.
func fmtHelper(field Field, className string, fieldName string) string {
	ftype := field.Ftype
	fvalue := field.Fvalue
	if tracing {
		fmt.Printf("DEBUG fmtHelper ftype=[%s], fvalue=[%v], className=[%s], fieldName=[%s]\n", ftype, fvalue, className, fieldName)
	}
	flagStatic := false
	if strings.HasPrefix(ftype, types.Static) {
		flagStatic = true
		ftype = ftype[1:] // get rid of leading 'X'
	}
	if len(fieldName) < 1 && flagStatic {
		return "<ERROR field name nil but field is static>"
	}

	if ftype == StringClassRef {
		// Special handling for String.
		if flagStatic {
			return fmt.Sprintf("%v", statics.GetStaticValue(className, fieldName))
		} else {
			if fvalue != nil {
				switch fvalue.(type) {
				case *[]byte:
					bytes := *fvalue.(*[]byte)
					last := len(bytes) - 1
					if last < 0 {
						return "\"\""
					}
					if bytes[last] == '\n' {
						bytes = bytes[0:last]
					}
					return fmt.Sprintf("\"%s\"", string(bytes))
				case string:
					return fvalue.(string)
				}
			} else {
				return "<nil>"
			}
		}
	}

	switch ftype {

	case types.Bool:
		// Special handling for boolean.
		if flagStatic {
			return fmt.Sprintf("%v", statics.GetStaticValue(className, fieldName))
		} else {
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
				return fmt.Sprintf("<ERROR Ftype=bool but unexpected Fvalue variable type: %T >", field.Fvalue)
			}
		}

	case types.ByteArray:
		// Special handling for non-String byte array.
		if flagStatic {
			return fmt.Sprintf("% x", statics.GetStaticValue(className, fieldName))
		} else {
			if field.Fvalue == nil {
				return "<ERROR nil Fvalue>"
			}
			switch field.Fvalue.(type) {
			case *Object:
				return "*** embedded object ***"
			}
			bytesPtr := field.Fvalue.(*[]byte)
			if bytesPtr == nil {
				return "<ERROR nil byte array ptr>"
			}
			if len(*bytesPtr) < 1 {
				return "<byte array of zero length>"
			}
			return fmt.Sprintf("% x", *bytesPtr)
		}
	}

	// Default action for anything else.
	if flagStatic {
		return fmt.Sprintf("%v", statics.GetStaticValue(className, fieldName))
	} else {
		return fmt.Sprintf("%v", field.Fvalue)
	}
}

// FormatField creates a string that represents a single field of an Object.
func (objPtr *Object) FormatField(fieldName string) string {
	var output string
	var klassString string // string class name
	obj := *objPtr         // whole object

	if obj.Klass != nil {
		klassString = *obj.Klass
	} else {
		klassString = "<ERROR nil class pointer>" // Why is there no class name pointer for this object?
		obj.DumpObject(klassString, 0)
		return klassString
	}

	// Use the FieldTable map with key fieldName?
	if len(fieldName) > 0 && len(obj.FieldTable) > 0 {
		// Using key="value" in the FieldTable
		ptr := obj.FieldTable[fieldName]
		if ptr == nil {
			str := fmt.Sprintf("<ERROR FieldTable[\"%s\"] not found>", fieldName)
			obj.DumpObject(str, 0)
			return str
		}
		field := *ptr
		str := fmtHelper(field, klassString, fieldName)
		if strings.HasPrefix(str, "<ERROR") {
			obj.DumpObject(str, 0)
		}
		output = fmt.Sprintf("%s: (%s) %s", fieldName, field.Ftype, str)
		return output
	}

	// Empty FieldTable. fieldName supplied?
	if len(fieldName) > 0 && tracing {
		// fieldName supplied but FieldTable is empty.
		title := fmt.Sprintf("DEBUG FormatField: fieldName=%s but FieldTable is empty", fieldName)
		obj.DumpObject(title, 0)
	}

	// fieldName was not supplied. FieldTable populated?
	if len(obj.FieldTable) > 0 && tracing {
		title := "DEBUG FormatField: FieldTable nonempty but fieldName is a nil string"
		obj.DumpObject(title, 0)
	}

	// Check use of the Fields slice.
	if len(obj.Fields) > 0 {
		// Using [0] in the Fields slice
		field := obj.Fields[0]
		str := fmtHelper(field, klassString, "")
		if strings.HasPrefix(str, "<ERROR") {
			obj.DumpObject(str, 0)
		}
		output = fmt.Sprintf("(%s) %s", obj.Fields[0].Ftype, str)
		return output
	}

	// Field table and field slice are both empty.
	if tracing {
		output = "<Field table and field slice are both empty>"
		obj.DumpObject(output, 0)
	}
	return klassString
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
		klassString = "\t<class MISSING>"
	}
	output += klassString + "\n"

	// Emit FieldTable.
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	nflds := len(obj.FieldTable)
	if nflds > 0 {
		output += fmt.Sprintf("\tField Table (%d):\n", nflds)
		for fieldName := range obj.FieldTable {
			if indent > 0 {
				output += strings.Repeat(" ", indent)
			}
			ptr := obj.FieldTable[fieldName]
			if ptr == nil {
				output += fmt.Sprintf("\t\t<ERROR nil FieldTable[%s] ptr>\n", fieldName)
			} else {
				str := fmtHelper(*obj.FieldTable[fieldName], klassString, fieldName)
				output += fmt.Sprintf("\t\tFld %s: (%s) %s\n", fieldName, obj.FieldTable[fieldName].Ftype, str)
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
			str := fmtHelper(fld, klassString, "")
			output += fmt.Sprintf("\t\tFld (%s) %s\n", fld.Ftype, str)
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
