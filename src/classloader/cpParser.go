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

// each entry in the CP table consists of an identifying integer (see enums above)
// and a structure that has varying fields and fulfills the posUdater interface.
// all entries are defined at the end of this file
type cpEntry struct {
	cpType int
	posUdater
}

var cp = make([]posUdater, 10)

func parseConstantPool(rawBytes []byte, klass *parsedClass) (int, error) {
	pos := 10 // position of the first byte of the constant pool

	// cpE := cpEntry{-1, dumbbell.updatePos(0) }
	// cpDummyEntry := dummyEntry{-1,interface{cpEntry.update}}
	cp[0] = cpEntry{cpType: -1}
	cp[1] = dumbbell{i: -1}
	println(cp)

	return pos, nil
}

type dumbbell struct {
	i int
}

func (dumbbell) update(i int) int {
	return 0
}

// d := dumbbell {i: -1, d.update: udpatePos, }

// func (dumbbell) update( i int) int { return 0 }

//
// var dumdum posUdater
func update(i int) int {
	return 0
}

// type dummyEntry posUdater {
//
// }
// func (dumbbell) update( i int ) int {
// 	return 0
// }
