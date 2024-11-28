/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import (
	"fmt"
	"jacobin/globals"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"strings"
)

var DEBUGGING = false
var SLICING = "_SLICING_"

// This file contains functions for the display of objects and their contents.
// These functions are used primarily in logging and debugging.

// fmtHelper - Return a string representing the field type and value.
//
// Input fields:
// * field: structure of field type and field value (if not static).
// * className: statics fields abd debugging
// * fieldName: Key to the jacobin statics table.
func fmtHelper(field Field, className string, fieldName string) string {
	ftype := field.Ftype
	fvalue := field.Fvalue
	if DEBUGGING {
		fmt.Printf("DEBUG fmtHelper ftype=[%s], fvalue=[%v], className=[%s], fieldName=[%s]\n", ftype, fvalue, className, fieldName)
	}

	// Static?
	flagStatic := strings.HasPrefix(ftype, types.Static)

	// Process Java String class reference.
	if ftype == types.StringClassRef {
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

	// Process the other types.
	switch ftype {
	case types.StringIndex:
		return *stringPool.GetStringPointer(fvalue.(uint32))
	case types.Bool:
		// Special handling for boolean.
		if flagStatic {
			return fmt.Sprintf("%v [static]", statics.GetStaticValue(className, fieldName))
		} else {
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
		var bytes []byte
		if flagStatic {
			return fmt.Sprintf("% x [static]", statics.GetStaticValue(className, fieldName))
		} else {
			if field.Fvalue == nil {
				return "<ERROR nil Fvalue>"
			}
			switch field.Fvalue.(type) {
			case *Object:
				return "*** embedded object ***"
			}
			switch field.Fvalue.(type) {
			case *[]byte:
				bptr := field.Fvalue.(*[]byte)
				bytes = *bptr
			case []byte:
				bytes = field.Fvalue.([]byte)
			default:
				errMsg := fmt.Sprintf("<type is byte array but value is of type %T>", field.Fvalue)
				return errMsg
			}
			if len(bytes) < 1 {
				return "<byte array of zero length>"
			}
			return fmt.Sprintf("% x", bytes)
		}
	}

	// Default action for anything else.
	if flagStatic {
		return fmt.Sprintf("%v [static]", statics.GetStaticValue(className, fieldName))
	} else {
		return fmt.Sprintf("%v", field.Fvalue)
	}
}

// FormatField creates a string that represents a single field of an Object.
func (objPtr *Object) FormatField(fieldName string) string {
	var output string
	var klassString string // string class name
	obj := *objPtr         // whole object

	if obj.KlassName != types.InvalidStringIndex {
		klassString = *stringPool.GetStringPointer(obj.KlassName)
	} else {
		klassString = "<ERROR nil class pointer>" // Why is there no class name pointer for this object?
		obj.DumpObject(klassString, 0)
		return klassString
	}

	// Use the FieldTable map with key fieldName?
	if len(fieldName) > 0 && len(obj.FieldTable) > 0 {
		// Using key="value" in the FieldTable
		ptr, ok := obj.FieldTable[fieldName]
		if !ok {
			str := fmt.Sprintf("<ERROR FieldTable[\"%s\"] not found>", fieldName)
			obj.DumpObject(str, 0)
			return str
		}
		field := ptr
		str := fmtHelper(field, klassString, fieldName)
		if strings.HasPrefix(str, "<ERROR") {
			obj.DumpObject(str, 0)
		}
		output = fmt.Sprintf("%s: (%s) %s", fieldName, field.Ftype, str)
		return output
	}

	// Empty FieldTable. fieldName supplied?
	if len(fieldName) > 0 && DEBUGGING {
		// fieldName supplied but FieldTable is empty.
		title := fmt.Sprintf("DEBUG FormatField: fieldName=%s but FieldTable is empty", fieldName)
		obj.DumpObject(title, 0)
	}

	// fieldName was not supplied. FieldTable populated?
	if len(obj.FieldTable) > 0 && DEBUGGING {
		title := "DEBUG FormatField: FieldTable nonempty but fieldName is a nil string"
		obj.DumpObject(title, 0)
	}

	// Field table is empty.
	if DEBUGGING {
		output = "<Field table is empty>"
		obj.DumpObject(output, 0)
	}
	return klassString
}

// DumpObject displays every attribute of an Object, formatted as multi-line printed output.
// 3 sections:
// * Class name
// * Field table
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
	if obj.KlassName != types.InvalidStringIndex {
		klassString = "\tClass: " + *(stringPool.GetStringPointer(obj.KlassName))
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
		// Create a sorted slice of keys.
		keys := make([]string, 0, len(obj.FieldTable))
		for key := range obj.FieldTable {
			keys = append(keys, key)
		}
		globals.SortCaseInsensitive(&keys)

		for _, fieldName := range keys {
			if indent > 0 {
				output += strings.Repeat(" ", indent)
			}
			_, ok := obj.FieldTable[fieldName]
			if !ok {
				output += fmt.Sprintf("\t\t<ERROR nil FieldTable[%s] ptr>\n", fieldName)
			} else {
				str := fmtHelper(obj.FieldTable[fieldName], klassString, fieldName)
				output += fmt.Sprintf("\t\tFld %s: (%s) %s\n", fieldName, obj.FieldTable[fieldName].Ftype, str)
			}
		}
	} else {
		output += fmt.Sprintf("\tField Table is <empty>\n")
	}

	// Emit END
	if indent > 0 {
		output += strings.Repeat(" ", indent)
	}
	output += "}\n"

	// Print output all at once.
	fmt.Print(output)
}
