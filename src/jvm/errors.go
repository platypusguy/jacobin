/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"encoding/binary"
	"jacobin/frames"
)

// routines for formatting error data when an error occurs inside the JVM

func formatStackOverflowError(f *frames.Frame) {
	// Stack overflow error. Change the bytecode to be IMPDEP2 and give info
	// in four bytes:
	// IMDEP2 (0xFF), 0x01 code for stack underflow, bytes 2 and 3:
	// the present PC written as an int16 value. First check that there
	// are enough bytes in the method that we can overwrite the first four bytes.
	currPC := int16(f.PC)
	if len(f.Meth) < 5 { // the present bytecode + 4 bytes for error info
		f.Meth = make([]byte, 5)
	}

	f.Meth[0] = 0x00 // dummy for the current bytecode
	f.Meth[1] = IMPDEP2
	f.Meth[2] = 0x01

	// now convert the PC at time of error into a two-byte value
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(currPC))
	f.Meth[3] = bytes[0]
	f.Meth[4] = bytes[1]
	f.PC = 0 // reset the current PC to point to the zeroth byte of our error data
}

func formatStackUnderflowError(f *frames.Frame) {
	// Stack underflow error. Change the bytecode to be IMPDEP2 and give info
	// in four bytes:
	// IMDEP2 (0xFF), 0x02 code for stack underflow, bytes 2 and 3:
	// the present PC written as an int16 value. First check that there
	// are enough bytes in the method that we can overwrite the first four bytes.
	currPC := int16(f.PC)
	if len(f.Meth) < 5 { // the present bytecode + 4 bytes for error info
		f.Meth = make([]byte, 5)
	}

	f.Meth[0] = 0x00 // dummy for the current bytecode
	f.Meth[1] = IMPDEP2
	f.Meth[2] = 0x02

	// now convert the PC at time of error into a two-byte value
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(currPC))
	f.Meth[3] = bytes[0]
	f.Meth[4] = bytes[1]
	f.PC = 0 // reset the current PC to point to the zeroth byte of our error data
}
