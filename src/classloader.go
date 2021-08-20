/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"fmt"
	"os"
)

// Classloaders hold the parsed bytecode in classes, where they can be retrieved
// and moved to an execution role. Most of the comments and code presuppose some
// familiarity with the role of classloaders. More information can be found at:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-5.html#jvms-5.3
type classloader struct {
	name    string
	parent  string
	classes map[string]loadedClass
}

type loadedClass struct {
}

// first canonicalizes the filename, checks whether the class is already loaded,
// and if not, then parses the class and loads it.
/*
// 1 TODO: canonicalize class name
// 2 TODO: search through classloaders for this class
// 3 TODO: determine which classloader should load the class, then
// 4 TODO: have *it* parse and load the class.
*/
func (cl classloader) loadClassFromFile(filename string) (loadedClass, error) {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		Log("Could not read file: "+filename+". Exiting.", SEVERE)
		return loadedClass{}, fmt.Errorf("file I/O error")
	}

	loadedKlass, err := cl.parseClass(rawBytes)
	if err != nil {
		Log("error parsing "+filename+". Exiting.", SEVERE)
		return loadedClass{}, fmt.Errorf("parsing error")
	} else {
		return loadedKlass, nil
	}
}

func (cl classloader) parseClass(rawBytes []byte) (loadedClass, error) {
	return loadedClass{}, nil
}
