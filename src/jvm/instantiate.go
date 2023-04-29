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
	"os"
	"unsafe"
	"sync"
    "time"
    "errors"
)

// instantiating a class is a two-part process:
// 1) the class needs to be loaded, so that its details and its methods are knowable
// 2) the class fields (if static) and instance fields (if non-static) are allocated. Details
//    for this second step appear in front of the initializeFields() method.

// Mutex for protecting the Log function during multithreading.
var mutex = sync.Mutex{}

func instantiateClass(classname string) (*object.Object, error) {
    _ = log.Log("instantiateClass: Instantiating class: "+classname, log.FINE)
    countDown := 4 // 2 seconds maximum time to load a class
    for true {
        mutex.Lock()
        k, present := classloader.Classes[classname]
        mutex.Unlock()
        if k.Status == 'I' {                         // the class is being loaded
            if countDown < 1 { // I've waited too long!
                msg := "instantiateClass: Status is still 'I' waiting for class: "+classname+". Overdue!"
                err := errors.New(msg)
                _ = log.Log(msg, log.SEVERE)
                return nil, err
            }
            countDown -= 1
            time.Sleep(500 * time.Millisecond)
            continue // recheck the status until it changes (i.e., until the class is loaded)
        }
        
        if present { break }
        
        // Not present - try to load from name
        if classloader.LoadClassFromNameOnly(classname) != nil {
            msg := "instantiateClass: LoadClassFromNameOnly("+classname+") failed. Exiting."
            err := errors.New(msg)
            _ = log.Log(msg, log.SEVERE)
            return nil, err
        }
        break // loaded by name
    }

	// at this point the class has been loaded into the method area (Classes).
	k, _ := classloader.Classes[classname]

	obj := object.Object{
		Klass: &k,
	}

	// the object's mark field contains the lower 32-bits of the object's
	// address, which serves as the hash code for the object
	uintp := uintptr(unsafe.Pointer(&obj))
	obj.Mark.Hash = uint32(uintp)

	if len(k.Data.Fields) > 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			f := k.Data.Fields[i]
			desc := k.Data.CP.CpIndex[f.Desc]
			if desc.Type != classloader.NameAndType {
				_ = log.Log("error creating field in: "+classname, log.SEVERE)
				return nil, classloader.CFE("invalid field type")
			}
			nameAndType := k.Data.CP.NameAndTypes[desc.Slot]
			ftype := classloader.FetchUTF8stringFromCPEntryNumber(
				&k.Data.CP, nameAndType.DescIndex)

			fieldToAdd := new(object.Field)
			fieldToAdd.Ftype = ftype
			switch string(fieldToAdd.Ftype[0]) {
			case "L", "[": // it's a reference
				fieldToAdd.Fvalue = nil
			case "B", "C", "I", "J", "S", "Z":
				fieldToAdd.Fvalue = 0
			case "D", "F":
				fieldToAdd.Fvalue = 0.0
			default:
				_ = log.Log("error creating field in: "+classname+
					" Invalid type: "+string(fieldToAdd.Ftype), log.SEVERE)
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
