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
	// g := globals.GetGlobalRef()
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
	CP.CpIndex[6] = classloader.CpEntry{Type: classloader.StringConst, Slot: 2} // point to UTF8[2]

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

	f.TOS = 0 // opStack[0] in theory contains a pointer to stdout, however, here we just use a zero value

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
