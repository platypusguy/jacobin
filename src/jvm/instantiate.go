/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/log"
	"os"
)

// Object is the layout of the data fields of an object. It's explained in more detail
// in the comments for initializeField()
type Object struct {
	klass  classloader.Klass
	mark   *MarkWord
	fields []Field
}

type MarkWord struct {
	m interface{}
}

type Field struct {
	metadata classloader.Field
	value    interface{}
}

// instantiating a class is a two-part process:
// 1) the class needs to be loaded, so that its details and its methods are knowable
// 2) the class fields (if static) and instance fields (if non-static) are allocated. Details
//    for this second step appear in front of the initializeFields() method.

func instantiateClass(classname string) (*Object, error) {
	_ = log.Log("Instantiating class: "+classname, log.FINE)
recheck:
	k, present := classloader.Classes[classname] // TODO: Put a mutex around this the same one used for writing.
	if k.Status == 'I' {                         // the class is being loaded
		goto recheck // recheck the status until it changes (i.e., until the class is loaded)
	} else if !present { // the class has not yet been loaded
		if classloader.LoadClassFromNameOnly(classname) != nil {
			_ = log.Log("Error loading class: "+classname+". Exiting.", log.SEVERE)
		}
	}

	// at this point the class has been loaded into the method area (Classes).
	k, _ = classloader.Classes[classname]

	obj := Object{
		klass:  k,
		mark:   &MarkWord{m: nil},
		fields: nil,
	}

	if len(k.Data.Fields) > 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			f := k.Data.Fields[i]
			initializeField(f, &k.Data.CP, classname, &obj)
		}
	}
	return &obj, nil
}

// the only fields allocated during class instantiation are class fields and instance fields--
// method-local fields are created on the stack during method execution.
// The allocated fields are in a structure that starts with a header area containing fields
// collectively referred to as oops: ordinary object pointers. These include two fields:
// * the mark word, which points to a struct with data about locking, a hashcode, and GC metadata
// * the klass word, which is a pointer back to the class definition as loaded by the classloader
// On some architectures, but not Jacobin, there is an additional field that insures that the
// fields that follow the oops are properly aligned for maximal performance.
func initializeField(f classloader.Field, cp *classloader.CPool, cn string, obj *Object) {
	name := cp.Utf8Refs[int(f.Name)]
	desc := cp.Utf8Refs[int(f.Desc)]
	var attr string = ""
	if len(f.Attributes) > 0 {
		for i := 0; i < len(f.Attributes); i++ {
			attr = cp.Utf8Refs[int(f.Attributes[i].AttrName)]
			if attr == "ConstantValue" {
				// valueIndex := int(f.Attributes[i].AttrContent[0])*256 +
				//     int(f.Attributes[i].AttrContent[1])
				// // valueType := cp.CpIndex[valueIndex].Type
				// // valueSlot := cp.CpIndex[valueIndex].Slot

			}
		}
	}
	// CURR: Resume here, entering the new field into obj.
	_, _ = fmt.Fprintf(os.Stdout, "Class: %s, Field to initialize: %s, type: %s\n", cn, name, desc)
	if attr != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Attribute name: %s\n", attr)
	}
}
