/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"jacobin/classloader"
	"jacobin/log"
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
			initializeField(f)
		}
	}
	return nil, nil
}

func initializeField(f classloader.Field) {
	//
}
