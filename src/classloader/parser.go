/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

type Klass struct {
}

// reads in a class file, parses it, and puts the values into the fields of the
// class that will be loaded into the classloader. Some verification performed
// receives the rawBytes of the class that were previously read in
//
// ClassFormatError - if the parser finds anything unexpected
func parse(rawBytes []byte) (parsedClass, error) {
	var klass = parsedClass{}

	err := parseMagicNumber(rawBytes)
	if err != nil {
		// log.Log("Classfile format error: "+err.Error(), log.SEVERE);
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
