/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

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
type posUdater interface {
	update(int) int
}

// the constant pool, which is an array of different record types. Each entry in the table
// consists of an identifying integer (see enums above)
// // and a structure that has varying fields and fulfills the posUdater interface.
// // all entries are defined at the end of this file
var cp []posUdater

func parseConstantPool(rawBytes []byte, klass *parsedClass) (int, error) {
	cp = make([]posUdater, klass.cpCount)
	pos := 10 // position of the first byte of the constant pool

	cp[0] = dummyEntry{Invalid}
	println(cp)

	return pos, nil
}

type dummyEntry struct {
	i int
}

func (dummyEntry) update(i int) int {
	return 0
}
