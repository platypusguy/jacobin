/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/log"
	"os"
)

func instantiateClass(classname string) (interface{}, error) {
	log.Log("Instantiating class: "+classname, log.FINEST)
recheck:
	k, present := classloader.Classes[classname]
	if k.Status == 'I' { // the class is being loaded
		goto recheck // recheck the status until it changes (i.e., the class is loaded)
	} else if !present { // the class has not yet been loaded
		if classloader.LoadClassFromNameOnly(classname) != nil {
			log.Log("Error loading class: "+classname+". Exiting.", log.SEVERE)
		}
	}

	// at this point the class has been loaded
	k, _ = classloader.Classes[classname]
	if len(k.Data.Fields) > 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			f := k.Data.Fields[i]
			initializeField(f, &k.Data.CP)
		}
	}
	return nil, nil
}

func initializeField(f classloader.Field, cp *classloader.CPool) {
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
	fmt.Fprintf(os.Stdout, "Field to initialize: %s, type: %s\n", name, desc)
	if attr != "" {
		fmt.Fprintf(os.Stdout, "Attribute name: %s\n", attr)
	}
}
