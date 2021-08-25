/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/globals"
	"strconv"
)

// reads in a class file, parses it, and puts the values into the fields of the
// class that will be loaded into the classloader. Some verification performed
// receives the rawBytes of the class that were previously read in
//
// ClassFormatError - if the parser finds anything unexpected
func parse(rawBytes []byte) (parsedClass, error) {

	// the parsed class as we'll give it to the classloader
	var pClass = parsedClass{}

	err := parseMagicNumber(rawBytes)
	if err != nil {
		return pClass, err
	}

	err = parseJavaVersionNumber(rawBytes, &pClass)
	if err != nil {
		return pClass, err
	}

	err = getConstantPoolCount(rawBytes, &pClass)
	return pClass, nil
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

// get the Java version number used in creating this class file. If it's higher than the
// version Jacobin presently supports, report an error.
func parseJavaVersionNumber(bytes []byte, klass *parsedClass) error {
	version, err := intFrom2Bytes(bytes, 6)
	if err != nil {
		return err
	}

	if version > globals.GetInstance().MaxJavaVersionRaw {
		errMsg := "Jacobin supports only Java versions through Java " +
			strconv.Itoa(globals.GetInstance().MaxJavaVersion)
		return cfe(errMsg)
	}

	klass.javaVersion = version
	println("Java version: " + strconv.Itoa(version))
	return nil
}

// get the number of entries in the constant pool. This number will
// be used later on to verify that the number of entries we fetch is
// correct. Note that this number is technically 1 greater than the
// number of actual entries, because the first entry in the constant
// pool is an empty placeholder, rather than an actual entry.
func getConstantPoolCount(bytes []byte, klass *parsedClass) error {
	cpEntryCount, err := intFrom2Bytes(bytes, 8)
	if err != nil || cpEntryCount <= 2 {
		return cfe("Invalid number of entries in constant pool: " +
			strconv.Itoa(globals.GetInstance().MaxJavaVersion))
	} else {
		klass.cpCount = cpEntryCount
		println("Number of CP entries: " + strconv.Itoa(cpEntryCount))
		return nil
	}
}
