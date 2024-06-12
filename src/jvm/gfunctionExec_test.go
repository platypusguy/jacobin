/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"io"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/opcodes"
	"jacobin/stringPool"
	"os"
	"strings"
	"testing"
)

func TestGfunctionExecValid(t *testing.T) {

	globals.InitGlobals("test")

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.CpIndex[6] = classloader.CpEntry{Type: classloader.StringConst, Slot: 7}
	CP.CpIndex[7] = classloader.CpEntry{Type: classloader.UTF8, Slot: 2} // point to UTF8[2]

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	methRef := classloader.MethodRefEntry{
		ClassIndex:  2, // these are CP entries
		NameAndType: 3,
	}
	CP.MethodRefs[0] = methRef

	printlnClassName := "java/io/PrintStream"
	CP.ClassRefs = append(CP.ClassRefs, stringPool.GetStringIndex(&printlnClassName))

	CP.Utf8Refs = append(CP.Utf8Refs, "println")
	CP.Utf8Refs = append(CP.Utf8Refs, "(Ljava/lang/String;)V")
	CP.Utf8Refs = append(CP.Utf8Refs, "Hello from test of gfunctionExec")
	nAndT := classloader.NameAndTypeEntry{
		NameIndex: uint16(4),
		DescIndex: uint16(5),
	}
	CP.NameAndTypes = append(CP.NameAndTypes, nAndT)

	f := newFrame(opcodes.LDC)
	f.Meth = append(f.Meth, uint8(6)) // point to the string constant referred to by CPindex[6]
	f.Meth = append(f.Meth, opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to method referred to in 0x0001 of the CP

	f.CP = &CP

	// create the opStack
	for j := 0; j < 10; j++ {
		f.OpStack = append(f.OpStack, 0)
	}
	f.OpStack[0] = os.Stdout
	f.TOS = 0

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	rawMsg, _ := io.ReadAll(rout)
	os.Stdout = normalStdout

	msg := string(rawMsg[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Hello from test of gfunctionExec") {
		t.Errorf("gfunctionExec: Did not get expected msg, got: %s", msg)
	}
}

// TestGfunctionExecTemplate is a way to test gfunctions that accept a string
// parameter via calls from the INVOKEVIRTUAL bytecode. It sets up the frame,
// the CP, and the environment so that the INVOKEVIRTUAL call is correctly simulated.
// Copy this template, insert your test values as explained in the comments,
// and rename and run your test. Also check that the test at the end is what you want.

func TestGfuncINVOKEVIRTUALwith1stringArgtemplate(t *testing.T) {
	// INVOKEVIRTUAL calls a method of an object. Below, enter the name of the object
	// class as a string (i.e., "java/io/PrintStream"), the method name, and the signature

	var objClassName string // e.g. "java/io/PrintStream"
	var methName string     // e.g. "println"
	var methType string     // e.g. "(Ljava/lang/String;)V"
	var stringParam string  // e.g. "test string"

	// ---------------------------------
	objClassName = "java/io/PrintStream"
	methName = "println"
	methType = "(Ljava/lang/String;)V"
	stringParam = "Hello from test of gfunctionExec"
	// ---------------------------------

	globals.InitGlobals("test")

	log.Init()
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	rerr, werr, _ := os.Pipe()
	os.Stderr = werr

	normalStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.CpIndex[6] = classloader.CpEntry{Type: classloader.StringConst, Slot: 7}
	CP.CpIndex[7] = classloader.CpEntry{Type: classloader.UTF8, Slot: 2} // point to UTF8[2]

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	methRef := classloader.MethodRefEntry{
		ClassIndex:  2, // these are CP entries
		NameAndType: 3,
	}
	CP.MethodRefs[0] = methRef

	CP.ClassRefs = append(CP.ClassRefs, stringPool.GetStringIndex(&objClassName))

	CP.Utf8Refs = append(CP.Utf8Refs, methName)
	CP.Utf8Refs = append(CP.Utf8Refs, methType)
	CP.Utf8Refs = append(CP.Utf8Refs, stringParam)
	nAndT := classloader.NameAndTypeEntry{
		NameIndex: uint16(4),
		DescIndex: uint16(5),
	}
	CP.NameAndTypes = append(CP.NameAndTypes, nAndT)

	f := newFrame(opcodes.LDC)
	f.Meth = append(f.Meth, uint8(6)) // point to the string constant parameter indexed by CPindex[6]
	f.Meth = append(f.Meth, opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to method referred to in 0x0001 of the CP

	f.CP = &CP

	// create the opStack
	for j := 0; j < 10; j++ {
		f.OpStack = append(f.OpStack, 0)
	}
	f.OpStack[0] = os.Stdout
	f.TOS = 0

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	_ = runFrame(fs)

	// get contents written by stderr and stdout, then
	// restore stderr and stdout to what they were before the test
	_ = werr.Close()
	rawStderrMsg, _ := io.ReadAll(rerr)
	os.Stderr = normalStderr

	_ = wout.Close()
	rawStdoutMsg, _ := io.ReadAll(rout)
	os.Stdout = normalStdout

	// convert the contents written to stderr and stdout into strings
	// and run tests on those strings

	errMsg := string(rawStderrMsg[:])
	if len(errMsg) != 0 {
		t.Errorf("gfunctionExec: Got unexpected error message: %s", errMsg)
	}

	outMsg := string(rawStdoutMsg[:])
	if !strings.Contains(outMsg, stringParam) {
		t.Errorf("gfunctionExec: Did not get expected msg, got: %s", outMsg)
	}
}
