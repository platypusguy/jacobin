/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"encoding/binary"
	"fmt"
	"jacobin/frames"
	"jacobin/log"
	"jacobin/thread"
	"runtime/debug"
	"strings"
)

// routines for formatting error data when an error occurs inside the JVM

// Stack overflow error (e.g., pushing a value when the stack is full, etc.)
func formatStackOverflowError(f *frames.Frame) {
	// Change the bytecode to be IMPDEP2 and give info in four bytes:
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

// Stack underflow error (e.g., trying to pop when the stack is empty, etc.)
func formatStackUnderflowError(f *frames.Frame) {
	// Change the bytecode to be IMPDEP2 and give info in four bytes:
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

// Prints out the frame stack
func showFrameStack(t *thread.ExecThread) {
	frameStack := t.Stack.Front()
	if frameStack == nil {
		_ = log.Log("No further data available", log.SEVERE)
		return
	}

	// step through the list-based stack of called methods and print contents
	for e := frameStack; e != nil; e = e.Next() {
		val := e.Value.(*frames.Frame)
		methName := fmt.Sprintf("%s.%s", val.ClName, val.MethName)
		data := fmt.Sprintf("Method: %-40s PC: %03d", methName, val.PC)
		_ = log.Log(data, log.SEVERE)
	}
	return
}

func showPanicCause(reason any) {
	// show the event that caused the panic
	if reason != nil {
		cause := fmt.Sprintf("%v", reason)
		_ = log.Log("\nerror: go panic because of "+cause+"\n", log.SEVERE)
	}
}

// in the event of a panic, this routine explains that a panic occurred and
// (to a limited extent why) and then prints the Jacobin frame stack and then
// the golang stack trace. r is the error returned when the panic occurs
func showGoStackTrace(reason any) {
	//
	// // show the Jaocbin frame stack
	// showFrameStack(&MainThread)
	// _ = log.Log("\n", log.SEVERE)

	// capture the golang function stack and convert it to
	// a slice of strings
	stack := string(debug.Stack())
	entries := strings.Split(stack, "\n")

	// remove the strings showing the internals of golang's panic stack trace
	var i int
	for i = 0; i < len(entries); i++ {
		if strings.HasPrefix(entries[i], "panic") {
			i += 2 //
			break
		}
	}

	// print the remaining strings in the golang stack trace
	for {
		if i < len(entries) {
			_ = log.Log(entries[i], log.SEVERE)
			i += 1
		} else {
			break
		}
	}
}
