/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
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
	Classes map[string]parsedClass
}

var AppCL Classloader
var BootstrapCL Classloader
var ExtensionCL Classloader

// the parsed class
type parsedClass struct {
}

// cfe = class format error, which is the error thrown by the parser for most
// of the errors arising from malformed bytecode
func cfe(msg string) error {
	errMsg := "Class Format Error: " + msg
	log.Log(errMsg+". Exiting.", log.SEVERE)
	return errors.New(errMsg)
}

// first canonicalizes the filename, checks whether the class is already loaded,
// and if not, then parses the class and loads it.
/*
// 1 TODO: canonicalize class name
// 2 TODO: search through classloaders for this class
// 3 TODO: determine which classloader should load the class, then
// 4 TODO: have *it* parse and load the class.
*/
func (cl Classloader) LoadClassFromFile(filename string) error {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Log("Could not read file: "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("file I/O error")
	} else {
		log.Log(filename+" read", log.FINE)
	}

	fullyParsedClass, err := parse(rawBytes)
	if err != nil {
		log.Log("error parsing "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("parsing error")
	} else {
		return insert(fullyParsedClass)
	}
}

func Init() error {
	BootstrapCL.Name = "bootstrap"
	BootstrapCL.Parent = ""
	BootstrapCL.Classes = make(map[string]parsedClass)

	ExtensionCL.Name = "extension"
	ExtensionCL.Parent = "bootstrap"
	ExtensionCL.Classes = make(map[string]parsedClass)

	AppCL.Name = "app"
	AppCL.Parent = "system"
	AppCL.Classes = make(map[string]parsedClass)
	return nil
}

// insert the fully parsed class into the classloader
func insert(class parsedClass) error {
	return nil //TODO: fill out after finishing parser
}

// the main parsing operation on the classfile. Drives the numberous functions that follow.
func parse(rawBytes []byte) (parsedClass, error) {
	var klass = parsedClass{}
	err := parseMagicNumber(rawBytes)
	if err != nil {
		return klass, err
	}

	return klass, nil
}

// all bytecode files start with 0xCAFEBABE ( it was the 90s!)
// this checks for that.
func parseMagicNumber(bytes []byte) error {
	if len(bytes) < 4 {
		return cfe("invalid magic number")
	} else if (bytes[0] != 0xCA) || (bytes[1] != 0xFE) || (bytes[2] != 0xBA) || (bytes[3] != 0xBE) {
		return cfe("invalid magic number")
	} else {
		return nil
	}
}
