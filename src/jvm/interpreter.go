/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

/* ================================================
   THIS IS AN EXPERIMENTAL ALTERNATIVE TO run.go
   The chages it makes:
   * Uses an array of functions rather than a switch for each bytecode
   * Does only one push and pull for 64-bit values (longs and doubles)
*/

package jvm

import (
	"jacobin/frames"
	"jacobin/object"
)

// set up a DispatchTable with 256 slots that correspond to the bytecodes
// each slot being a pointer to a function that accepts a pointer to the
// current frame and an int parameter. It returns an int that indicates
// how much to increase that frame's PC (program counter) by.
type BytecodeFunc func(*frames.Frame, int64) int

var DispatchTable = [256]BytecodeFunc{
	doNop,        // NOP         0x00
	doAconstNull, // ACONST_NULL 0x01
	doIconstM1,   // ICONST_M1   0x02
	doIconst0,    // ICONST_0    0x03
	doIconst1,    // ICONST_1    0x04
	doIconst2,    // ICONST_2    0x05
	doIconst3,    // ICONST_3    0x06
	doIconst4,    // ICONST_4    0x07
	doIconst5,    // ICONST_5    0x08
	doLconst0,    // LCONST_0    0x09
	doLconst1,    // LCONST_1    0x0A
	doFconst0,    // FCONST_0    0x0B
	doFconst1,    // FCONST_1    0x0C
	doFconst2,    // FCONST_2    0x0D
	doDconst0,    // DCONST_0    0x0E
	doDconst1,    // DCONST_1    0x0F
}

// the functions, listed here in numerical order of the bytecode
func doNop(_ *frames.Frame, _ int64) int { return 1 }
func doAconstNull(f *frames.Frame, _ int64) int {
	push(f, object.Null)
	return 1
}

func doIconstM1(f *frames.Frame, _ int64) int { return pushInt(f, int64(-1)) }
func doIconst0(f *frames.Frame, _ int64) int  { return pushInt(f, int64(0)) }
func doIconst1(f *frames.Frame, _ int64) int  { return pushInt(f, int64(1)) }
func doIconst2(f *frames.Frame, _ int64) int  { return pushInt(f, int64(2)) }
func doIconst3(f *frames.Frame, _ int64) int  { return pushInt(f, int64(3)) }
func doIconst4(f *frames.Frame, _ int64) int  { return pushInt(f, int64(4)) }
func doIconst5(f *frames.Frame, _ int64) int  { return pushInt(f, int64(5)) }
func doLconst0(f *frames.Frame, _ int64) int  { return pushInt(f, int64(0)) }
func doLconst1(f *frames.Frame, _ int64) int  { return pushInt(f, int64(1)) }
func doFconst0(f *frames.Frame, _ int64) int  { return pushFloat(f, int64(0)) }
func doFconst1(f *frames.Frame, _ int64) int  { return pushFloat(f, int64(1)) }
func doFconst2(f *frames.Frame, _ int64) int  { return pushFloat(f, int64(2)) }
func doDconst0(f *frames.Frame, _ int64) int  { return pushFloat(f, int64(0)) }
func doDconst1(f *frames.Frame, _ int64) int  { return pushFloat(f, int64(1)) }

// the functions call by the dispatched functions
func pushInt(f *frames.Frame, intToPush int64) int {
	push(f, intToPush)
	return 1
}

func pushFloat(f *frames.Frame, intToPush int64) int {
	push(f, float64(intToPush))
	return 1
}

func interpretBytecodes(bytecode int, f *frames.Frame) int {
	PC := DispatchTable[bytecode](f, 0)
	println("PC after call to DispatchTable[", bytecode, "] = ", PC)
	return PC
}
