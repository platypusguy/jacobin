/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"jacobin/log"
)

type Klass struct {
}

// reads in a class file, parses it, and puts the values into the fields of the
// class that will be loaded into the classloader. Some verification performed
// receives the rawBytes of the class that were previously read in
//
// ClassFormatError - if the parser finds anything unexpected
func Parse(rawBytes []byte) (Klass, error) {
	var parsedClass Klass

	if len(rawBytes) < 10 {
		log.Log("Classfile format error.", log.SEVERE)
		return Klass{}, errors.New("Class format error")
	}

	return parsedClass, nil
}
