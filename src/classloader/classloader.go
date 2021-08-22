/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"jacobin/log"
	"os"
)

// Classloaders hold the parsed bytecode in classes, where they can be retrieved
// and moved to an execution role. Most of the comments and code presuppose some
// familiarity with the role of classloaders. More information can be found at:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-5.html#jvms-5.3
type Classloader struct {
	Name    string
	Parent  string
	Classes map[string]Klass
}

var AppCL Classloader
var BootstrapCL Classloader
var ExtensionCL Classloader

// first canonicalizes the filename, checks whether the class is already loaded,
// and if not, then parses the class and loads it.
/*
// 1 TODO: canonicalize class name
// 2 TODO: search through classloaders for this class
// 3 TODO: determine which classloader should load the class, then
// 4 TODO: have *it* parse and load the class.
*/
func (cl Classloader) LoadClassFromFile(filename string) (Klass, error) {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Log("Could not read file: "+filename+". Exiting.", log.SEVERE)
		return Klass{}, fmt.Errorf("file I/O error")
	} else {
		log.Log(filename+" read", log.FINE)
	}

	parsedClass, err := Parse(rawBytes)
	if err != nil {
		log.Log("error parsing "+filename+". Exiting.", log.SEVERE)
		return Klass{}, fmt.Errorf("parsing error")
	} else {
		return parsedClass, nil
	}
}

func Init() error {
	BootstrapCL.Name = "bootstrap"
	BootstrapCL.Parent = ""
	BootstrapCL.Classes = make(map[string]Klass)

	ExtensionCL.Name = "extension"
	ExtensionCL.Parent = "bootstrap"
	ExtensionCL.Classes = make(map[string]Klass)

	AppCL.Name = "app"
	AppCL.Parent = "system"
	AppCL.Classes = make(map[string]Klass)
	return nil
}
