/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/log"
	"jacobin/object"
	"jacobin/shutdown"
	"os"
	"unsafe"
)

// instantiating a class is a two-part process:
// 1) the class needs to be loaded, so that its details and its methods are knowable
// 2) the class fields (if static) and instance fields (if non-static) are allocated. Details
//    for this second step appear in front of the initializeFields() method.

// Mutex for protecting the Log function during multithreading.
// var mutex = sync.Mutex{}

func instantiateClass(classname string) (*object.Object, error) {
	_ = log.Log("instantiateClass: Instantiating class: "+classname, log.FINE)
	k := classloader.MethAreaFetch(classname)
	if k == nil || k.Data == nil {

		// Not present - try to load from name
		if classloader.LoadClassFromNameOnly(classname) != nil {
			msg := "instantiateClass: Failed to load class " + classname
			_ = log.Log(msg, log.SEVERE)
			shutdown.Exit(shutdown.APP_EXCEPTION)
		}

		// Success in loaded by name

	}

	// at this point the class has been loaded into the method area (MethArea).
	k = classloader.MethAreaFetch(classname)

	obj := object.Object{
		Klass: k,
	}

	// the object's mark field contains the lower 32-bits of the object's
	// address, which serves as the hash code for the object
	uintp := uintptr(unsafe.Pointer(&obj))
	obj.Mark.Hash = uint32(uintp)

	if len(k.Data.Fields) > 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			f := k.Data.Fields[i]
			if (f.AccessFlags & 0b00001000) != 0 {
				fmt.Fprintf(os.Stdout, "%s is static",
					classloader.FetchUTF8stringFromCPEntryNumber(&k.Data.CP, f.Name))
			}
			// name := k.Data.CP.CpIndex[f.Name]
			desc := k.Data.CP.Utf8Refs[f.Desc]
			// if desc.Type != classloader.NameAndType {
			// 	_ = log.Log("error creating field in: "+classname, log.SEVERE)
			// 	return nil, classloader.CFE("invalid field type")
			// }
			// nameAndType := k.Data.CP.NameAndTypes[desc.Slot]
			// ftype := classloader.FetchUTF8stringFromCPEntryNumber(
			// 	&k.Data.CP, nameAndType.DescIndex)

			fieldToAdd := new(object.Field)
			fieldToAdd.Ftype = desc
			switch string(fieldToAdd.Ftype[0]) {
			case "L", "[": // it's a reference
				fieldToAdd.Fvalue = nil
			case "B", "C", "I", "J", "S", "Z":
				fieldToAdd.Fvalue = 0
			case "D", "F":
				fieldToAdd.Fvalue = 0.0
			default:
				_ = log.Log("error creating field in: "+classname+
					" Invalid type: "+fieldToAdd.Ftype, log.SEVERE)
				return nil, classloader.CFE("invalid field type")
			}
			obj.Fields = append(obj.Fields, *fieldToAdd)
			// CURR: resume here
		}
	}
	return &obj, nil
}

// the only fields allocated during class instantiation are instance fields--
// method-local fields are created on the stack during method execution.
// The allocated fields are in a structure that is defined in object.go
func initializeField(f classloader.Field, cp *classloader.CPool, cn string, obj *object.Object) {
	name := cp.Utf8Refs[int(f.Name)]
	desc := cp.Utf8Refs[int(f.Desc)]
	var attr = ""
	// var value int64
	if len(f.Attributes) > 0 {
		for i := 0; i < len(f.Attributes); i++ {
			attr = cp.Utf8Refs[int(f.Attributes[i].AttrName)]
			if attr == "ConstantValue" {
				// valueIndex := int(f.Attributes[i].AttrContent[0])*256 +
				//     int(f.Attributes[i].AttrContent[1])
				// // valueType := cp.CpIndex[valueIndex].Type
				// // valueSlot := cp.CpIndex[valueIndex].Slot
			}
			// } else {
			// 	value = 0
			// }

		}
	}
	// append field to the object.fields slice TODO: check that this is a good solution.
	// obj.fields = append(obj.fields, Field{
	// 	metadata: classloader.Field{},
	// 	value:    value,
	// })
	// CURR: Resume here, entering the new field into obj.
	_, _ = fmt.Fprintf(os.Stdout, "Class: %s, Field to initialize: %s, type: %s\n", cn, name, desc)
	if attr != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Attribute name: %s\n", attr)
	}
}
