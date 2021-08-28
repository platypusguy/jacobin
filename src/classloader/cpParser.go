/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"os"
)

// this file contains the parser for the constant pool and the verifier.
// Refer to: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4-140

// the various types of entries in the constant pool
const (
	Invalid       = -1 // used for initialization and for dummy entries (viz. for longs, doubles)
	UTF8          = 1
	IntConst      = 3
	FloatConst    = 4
	LongConst     = 5
	DoubleConst   = 6
	ClassRef      = 7
	String        = 8
	Field         = 9
	Method        = 10
	Interface     = 11
	NameAndType   = 12
	MethodHandle  = 15
	MethodType    = 16
	Dynamic       = 17
	InvokeDynamic = 18
	Module        = 19
	Package       = 20
)

// every CP entry must update the position in the classfile
type cpTyper interface {
	getType() int
}

// the constant pool, which is an array of different record types. Each entry in the table
// consists of an identifying integer (see enums above)
// // and a structure that has varying fields and fulfills the cpTyper interface.
// // all entries are defined at the end of this file
var cp []cpTyper

func parseConstantPool(rawBytes []byte, klass *parsedClass) (int, error) {
	cp = make([]cpTyper, klass.cpCount)
	var i int
	pos := 9 // position of the last byte before the constant pool

	cp[0] = dummyEntry{Invalid}
	for i = 1; i <= klass.cpCount-1; {
		pos += 1
		entryType := int(rawBytes[pos])
		switch entryType {
		case Method:
			{
				classIndex, _ := intFrom2Bytes(rawBytes, pos+1)
				nameAndTypeIndex, _ := intFrom2Bytes(rawBytes, pos+3)
				mre := methodRefEntry{Method, classIndex, nameAndTypeIndex}
				cp[i] = mre
				pos += 4
				i += 1
			}
		default:
			klass.cpCount = i // just to get it over with for the moment
		}
	}

	for j := 0; j < i; j += 1 {
		fmt.Fprintf(os.Stderr, "CP entry: %d, type %d\n", j, cp[j].getType())
		//TODO: see if I can retrieve an entry from cp and read its individual fields
	}
	return pos, nil
}

// ==== the various entry types in the constant pool (listed in order of the enums above) ====
type dummyEntry struct { // type -1 (invalid or dummy entry)
	cpeType int
}

func (dummyEntry) getType() int {
	return -1
}

type methodRefEntry struct { // type: 10 (method reference)
	cpeType           int
	nameIndex         int
	classAndTypeIndex int
}

func (methodRefEntry) getType() int {
	return 10
}
