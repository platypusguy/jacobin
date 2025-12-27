/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

// Bytecodes tested in alphabetical order. Non-bytecode tests at end of file.
// Note: array bytecodes are in interpreter_arrayBytecodes_test.go

// LLOAD: test load of long in locals[index] on to stack
func TestLload(t *testing.T) {
	f := newFrame(opcodes.LLOAD)
	f.Meth = append(f.Meth, 0x04) // use local var #4
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x1234562)) // put value in locals[4]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	x := pop(&f).(int64)
	if x != 0x1234562 {
		t.Errorf("LLOAD: Expecting 0x1234562 on stack, got: 0x%x", x)
	}
	if f.TOS != -1 {
		t.Errorf("LLOAD: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
	if f.PC != 2 {
		t.Errorf("LLOAD: Expected pc to be pointing at byte 2, got: %d", f.PC)
	}
}

// LLOAD_0: Load long from locals[0]
func TestLload0(t *testing.T) {
	f := newFrame(opcodes.LLOAD_0)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[0]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	x := pop(&f).(int64)
	if x != 0x12345678 {
		t.Errorf("LLOAD_0: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_0: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LLOAD_1: Load long from locals[1]
func TestLload1(t *testing.T) {
	f := newFrame(opcodes.LLOAD_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[1]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	x := pop(&f).(int64)

	if x != 0x12345678 {
		t.Errorf("LLOAD_1: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_1: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LLOAD_2: Load long from locals[2]
func TestLload2(t *testing.T) {
	f := newFrame(opcodes.LLOAD_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[2]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	x := pop(&f).(int64)

	if x != 0x12345678 {
		t.Errorf("LLOAD_2: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.Locals[2] != x {
		t.Errorf("LLOAD_2: Local variable[3] holds invalid value: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_1: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LLOAD_3: Load long from locals[3]
func TestLload3(t *testing.T) {
	f := newFrame(opcodes.LLOAD_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, int64(0x12345678)) // put value in locals[3]

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	x := pop(&f).(int64)
	if x != 0x12345678 {
		t.Errorf("LLOAD_3: Expecting 0x12345678 on stack, got: 0x%x", x)
	}

	if f.TOS != -1 {
		t.Errorf("LLOAD_3: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LMUL: pop 2 longs, multiply them, push result
func TestLmul(t *testing.T) {
	f := newFrame(opcodes.LMUL)
	push(&f, int64(10))
	push(&f, int64(7))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("LMUL, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != 70 {
		t.Errorf("LMUL: Expected popped value to be 70, got: %d", value)
	}
}

// LNEG: negate a long
func TestLneg(t *testing.T) {
	f := newFrame(opcodes.LNEG)
	push(&f, int64(10))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("LNEG, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != -10 {
		t.Errorf("LNEG: Expected popped value to be -10, got: %d", value)
	}
}

// LOR: Logical OR of two longs
func TestLor(t *testing.T) {
	f := newFrame(opcodes.LOR)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 23 { // 21 | 22 = 23
		t.Errorf("LOR: expected a result of 23, but got: %d", value)
	}

	if f.TOS != -1 {
		t.Errorf("LOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LREM: remainder of long division (the % operator)
func TestLrem(t *testing.T) {
	f := newFrame(opcodes.LREM)
	push(&f, int64(74))
	push(&f, int64(6))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("LREM, Top of stack, expected 0, got: %d", f.TOS)
	}

	value := pop(&f).(int64)
	if value != 2 {
		t.Errorf("LREM: Expected popped value to be 2, got: %d", value)
	}
}

// LREM: long modulo -- divide by zero
func TestLremDivideByZero(t *testing.T) {
	globals.InitGlobals("test")

	// hide the error message to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.LREM)
	push(&f, int64(6))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	// reset stderr to its normal stream
	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "division by zero") {
		t.Errorf("LREM: Expected division by zero error msg, got: %s", errMsg)
	}
}

// LRETURN: Return a long from a function
func TestLreturn(t *testing.T) {
	f0 := newFrame(0)
	push(&f0, int64(20))
	fs := frames.CreateFrameStack()
	fs.PushFront(&f0)
	f1 := newFrame(opcodes.LRETURN)
	push(&f1, int64(21))

	fs.PushFront(&f1)
	interpret(fs) // LRETURN pops the topmost frame it's returning from

	f3 := fs.Front().Value.(*frames.Frame)
	newVal := pop(f3).(int64)
	if newVal != 21 {
		t.Errorf("After LRETURN, expected a value of 21 in previous frame, got: %d", newVal)
	}
}

// LSHL: Left shift of long
func TestLshl(t *testing.T) {
	f := newFrame(opcodes.LSHL)
	push(&f, int64(22))
	push(&f, int64(3)) // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 176 { // 22 << 3 = 176
		t.Errorf("LSHL: expected a result of 176, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LSHL: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LSHR: Right shift of long
func TestLshr(t *testing.T) {
	f := newFrame(opcodes.LSHR)
	push(&f, int64(200))

	push(&f, int64(3)) // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("LSHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LSHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LSTORE: Store long from stack into local specified by following byte, and the local var after it.
func TestLstore(t *testing.T) {
	f := newFrame(opcodes.LSTORE)
	f.Meth = append(f.Meth, 0x02) // use local var #2
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x22223))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[2] != int64(0x22223) {
		t.Errorf("LSTORE: Expecting 0x22223 in locals[2], got: 0x%x", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE: Expecting an empty stack, but tos points to item: %d", f.TOS)
	}
}

// LSTORE_0: Store long from stack in localVar[0] and again in localVar[1]
func TestLstore0(t *testing.T) {
	f := newFrame(opcodes.LSTORE_0)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[0] != int64(0x12345678) {
		t.Errorf("LSTORE_0: expected locals[0] to be 0x12345678, got: %d", f.Locals[0])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_0: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSTORE_1: Store long from stack in localVar[1]
func TestLstore1(t *testing.T) {
	f := newFrame(opcodes.LSTORE_1)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[1] != int64(0x12345678) {
		t.Errorf("LSTORE_1: expected locals[1] to be 0x12345678, got: %d", f.Locals[1])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_1: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSTORE_2: Store long from stack in localVar[2] and again in localVar[3]
func TestLstore2(t *testing.T) {
	f := newFrame(opcodes.LSTORE_2)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[2] != int64(0x12345678) {
		t.Errorf("LSTORE_2: expected locals[2] to be 0x12345678, got: %d", f.Locals[2])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_2: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSTORE_3: Store long from stack in localVar[3] and again in localVar[]
func TestLstore3(t *testing.T) {
	f := newFrame(opcodes.LSTORE_3)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero)
	f.Locals = append(f.Locals, zero) // LSTORE instructions fill two local variables (with the same value)
	push(&f, int64(0x12345678))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[3] != int64(0x12345678) {
		t.Errorf("LSTORE_3: expected locals[3] to be 0x12345678, got: %d", f.Locals[3])
	}

	if f.TOS != -1 {
		t.Errorf("LSTORE_3: Expected op stack to be empty, got tos: %d", f.TOS)
	}
}

// LSUB: Subtract two longs
func TestLsub(t *testing.T) {
	f := newFrame(opcodes.LSUB)
	push(&f, int64(10))
	push(&f, int64(7))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64)

	if f.TOS != -1 {
		t.Errorf("LSUB, Top of stack, expected -1, got: %d", f.TOS)
	}

	if value != 3 {
		t.Errorf("LSUB: Expected popped value to be 3, got: %d", value)
	}
}

// LUSHR: Right unsigned shift of long
func TestLushr(t *testing.T) {
	f := newFrame(opcodes.LUSHR)
	push(&f, int64(200))
	push(&f, int64(3)) // shift left 3 bits

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 25 { // 200 >> 3 = 25
		t.Errorf("LUSHR: expected a result of 25, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LUSHR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// LXOR: Logical XOR of two longs
func TestLxor(t *testing.T) {
	f := newFrame(opcodes.LXOR)
	push(&f, int64(21))
	push(&f, int64(22))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	value := pop(&f).(int64) // longs require two slots, so popped twice

	if value != 3 { // 21 ^ 22 = 3
		t.Errorf("LXOR: expected a result of 3, but got: %d", value)
	}
	if f.TOS != -1 {
		t.Errorf("LXOR: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// MONITORENTER: The JDK JVM does not implement this, nor do we. So just pop the ref off stack
func TestMonitorEnter(t *testing.T) {
	f := newFrame(opcodes.MONITORENTER)
	push(&f, &f) // push any value and make sure it gets popped off

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != -1 {
		t.Errorf("MONITORENTER: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// MONITOREXIT: The JDK JVM does not implement this, nor do we. So just pop the ref off stack
func TestMonitorExit(t *testing.T) {
	f := newFrame(opcodes.MONITOREXIT)
	push(&f, &f) // push any value and make sure it gets popped off

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != -1 {
		t.Errorf("MONITOREXIT: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// NEW: Instantiate object -- here with an error
func TestNewWithError(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.NEW)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be class or interface
	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("NEW: Expected error message, but got none")
	}

	if !strings.Contains(errMsg, "Invalid type for new object") {
		t.Errorf("NEW: got unexpected error message: %s", errMsg)
	}
}

// PEEK: test peek, stack underflow with Jacobin error message
func TestPeekWithStackUnderflow(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("testWithoutShutdown")
	gl := globals.GetGlobalRef()

	gl.FuncInstantiateClass = InstantiateClass
	gl.FuncThrowException = exceptions.ThrowExNil
	gl.FuncFillInStackTrace = gfunction.FillInStackTrace
	gl.StrictJDK = false

	stringPool.PreloadArrayClassesToStringPool()
	trace.Init()

	err := classloader.Init()
	if err != nil {
		t.Fail()
	}
	classloader.LoadBaseClasses()

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	classloader.LoadBaseClasses()
	_ = classloader.LoadClassFromNameOnly("java/lang/Object")

	// Create a Java-level Thread object (no use of jvmThread.go ExecThread)
	InitGlobalFunctionPointers()
	gfunction.InitializeGlobalThreadGroups()

	thObj := gfunction.ThreadCreateNoarg(nil).(*object.Object)
	main := object.StringObjectFromGoString("main")
	params := []any{thObj, main}
	gfunction.ThreadInitWithName(params)
	gfunction.RegisterThread(thObj) // put into globals.Threads map
	thID := int(thObj.FieldTable["ID"].Fvalue.(int64))

	f := frames.CreateFrame(1)
	f.ClName = "java/lang/Double" // Not a G-function so catchFrame won't vomit.
	f.MethName = "hashCode"       // -------------------------------------------
	f.MethType = "()I"            // -------------------------------------------
	_, err = classloader.FetchMethodAndCP(f.ClName, f.MethName, f.MethType)
	for i := 0; i < 4; i++ {
		f.OpStack = append(f.OpStack, int64(0))
	}
	f.TOS = -1
	f.Thread = thID

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	// Attach the JVM frame stack to the Java thread object
	thObj.FieldTable["framestack"] = object.Field{Ftype: types.LinkedList, Fvalue: fs}

	_ = peek(f)

	_ = w.Close()
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	// txt, _ := io.ReadAll(rout)

	os.Stderr = normalStderr
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "stack underflow") ||
		!strings.Contains(msg, "org.jacobin.InternalException") { // use the Jacobin error message
		t.Errorf("got unexpected error message: %s", msg)
	}
}

// PEEK: test peek, stack underflow with Jacobin error message
func TestPeekWithStackUnderflowStrictJDK(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("testWithoutShutdown")
	gl := globals.GetGlobalRef()

	gl.FuncInstantiateClass = InstantiateClass
	gl.FuncThrowException = exceptions.ThrowExNil
	gl.FuncFillInStackTrace = gfunction.FillInStackTrace
	gl.StrictJDK = true

	stringPool.PreloadArrayClassesToStringPool()
	trace.Init()

	err := classloader.Init()
	if err != nil {
		t.Fail()
	}
	classloader.LoadBaseClasses()

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	classloader.LoadBaseClasses()
	_ = classloader.LoadClassFromNameOnly("java/lang/Object")

	// Create a Java-level Thread object (no use of jvmThread.go ExecThread)
	InitGlobalFunctionPointers()
	gfunction.InitializeGlobalThreadGroups()
	thObj := gfunction.ThreadCreateNoarg(nil).(*object.Object)
	main := object.StringObjectFromGoString("main")
	params := []any{thObj, main}
	gfunction.ThreadInitWithName(params)
	gfunction.RegisterThread(thObj) // put into globals.Threads map
	thID := int(thObj.FieldTable["ID"].Fvalue.(int64))

	f := frames.CreateFrame(1)
	f.ClName = "java/lang/Double" // Not a G-function so catchFrame won't vomit.
	f.MethName = "hashCode"       // -------------------------------------------
	f.MethType = "()I"            // -------------------------------------------
	_, err = classloader.FetchMethodAndCP(f.ClName, f.MethName, f.MethType)
	for i := 0; i < 4; i++ {
		f.OpStack = append(f.OpStack, int64(0))
	}
	f.TOS = -1
	f.Thread = thID

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	// Attach the JVM frame stack to the Java thread object
	thObj.FieldTable["framestack"] = object.Field{Ftype: types.LinkedList, Fvalue: fs}

	_ = peek(f)

	_ = w.Close()
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	// txt, _ := io.ReadAll(rout)

	os.Stderr = normalStderr
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "stack underflow") ||
		!strings.Contains(msg, "com.sun.jdi.InternalException") { // use the HotSpot error message
		t.Errorf("got unexpected error message: %s", msg)
	}
}

// POP: pop item off stack and discard it
func TestPop(t *testing.T) {
	f := newFrame(opcodes.POP)
	push(&f, int64(34)) // push three different values
	push(&f, int64(21))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 1 {
		t.Errorf("POP: Expected stack with 2 items, but got a tos of: %d", f.TOS)
	}

	top := pop(&f).(int64)

	if top != 21 {
		t.Errorf("POP: expected top's value to be 21, but got: %d", top)
	}
}

// POP with tracing enabled
func TestPopWithTracing(t *testing.T) {
	f := newFrame(opcodes.POP)
	push(&f, int64(34)) // push three different values
	push(&f, int64(21))
	push(&f, int64(0))

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)         // push the new frame
	globals.TraceInst = true // turn on tracing
	interpret(fs)

	if f.TOS != 1 {
		t.Errorf("POP: Expected stack with 2 items, but got a tos of: %d", f.TOS)
	}

	top := pop(&f).(int64)

	if top != 21 {
		t.Errorf("POP: expected top's value to be 21, but got: %d", top)
	}

	if globals.TraceInst != true {
		t.Errorf("POP: globals.TraceInst was not re-enabled after the POP execution")
	}
}

// POP with stack underflow error
func TestPopWithStackUnderflow(t *testing.T) {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("testWithoutShutdown")
	gfunction.InitializeGlobalThreadGroups()
	gl := globals.GetGlobalRef()

	gl.FuncInstantiateClass = InstantiateClass
	gl.FuncThrowException = exceptions.ThrowExNil
	gl.FuncFillInStackTrace = gfunction.FillInStackTrace
	gl.FuncInvokeGFunction = gfunction.Invoke
	gl.FuncMinimalAbort = exceptions.MinimalAbort
	gl.FuncRunThread = RunJavaThread

	stringPool.PreloadArrayClassesToStringPool()
	trace.Init()

	err := classloader.Init()
	if err != nil {
		t.Fail()
	}
	classloader.LoadBaseClasses()

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	classloader.LoadBaseClasses()
	_ = classloader.LoadClassFromNameOnly("java/lang/Object")
	classloader.FetchMethodAndCP("java/lang/Object", "wait", "(JI)V")

	var f *frames.Frame
	globals.InitGlobals("test")
	gfunction.InitializeGlobalThreadGroups()
	thObj := gfunction.ThreadCreateNoarg(nil).(*object.Object)
	main := object.StringObjectFromGoString("main")
	params := []any{thObj, main}
	gfunction.ThreadInitWithName(params)
	gfunction.RegisterThread(thObj) // put into globals.Threads map
	thID := int(thObj.FieldTable["ID"].Fvalue.(int64))

	f = frames.CreateFrame(1)
	f.ClName = "java/lang/Object"
	f.MethName = "wait"
	f.MethType = "(JI)V"
	_, err = classloader.FetchMethodAndCP(f.ClName, f.MethName, f.MethType)
	for i := 0; i < 4; i++ {
		f.OpStack = append(f.OpStack, int64(0))
	}
	f.TOS = -1
	f.Thread = thID

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	// Attach the JVM frame stack to the Java thread object
	thObj.FieldTable["framestack"] = object.Field{Ftype: types.LinkedList, Fvalue: fs}

	_ = pop(f)

	_ = w.Close()
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	// txt, _ := io.ReadAll(rout)

	os.Stderr = normalStderr
	os.Stdout = normalStdout

	msg := string(out[:])

	if !strings.Contains(msg, "stack underflow") {
		t.Errorf("got unexpected error message: %s", msg)
	}
}

// The previous tests for pop test it as an action performed by Jacobin in the course
// of handling other bytecodes. Here we test the POP bytecode. We know from previous
// tests it works correctly. So, here we test only that it handles errors correctly.
func TestPopBytecodrUnderflow(t *testing.T) {
	globals.InitGlobals("test")

	// hide the error message to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.POP)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "stack underflow in POP") {
		t.Errorf("Did not get expected error from invalide POP, got: %s", errMsg)
	}
}

// POP2: pop two items
func TestPop2(t *testing.T) {

	f := newFrame(opcodes.POP2)
	// push three different values; 34 at bottom
	push(&f, int64(34)) // fload_0 : Load float from local variable
	push(&f, int64(21)) // iload : Load int from local variable
	push(&f, int64(10)) // lconst_1 : Push long constant

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.TOS != 1 {
		t.Errorf("POP2: Expected stack with 1 item, but got a tos of: %d", f.TOS)
	}

	top := pop(&f).(int64)

	if top != 21 {
		t.Errorf("POP2: expected top's value to be 21, but got: %d", top)
	}
}

// POP2: pop two items off stack -- make sure tracing doesn't affect the output
func TestPop2WithTrace(t *testing.T) {
	f := newFrame(opcodes.POP2)
	// push three different values; 34 at bottom
	push(&f, int64(34)) // fload_0 : Load float from local variable
	push(&f, int64(21)) // iload : Load int from local variable
	push(&f, int64(10)) // lconst_1 : Push long constant

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)         // push the new frame
	globals.TraceInst = true // turn on tracing
	interpret(fs)

	if f.TOS != 1 {
		t.Errorf("POP2: Expected stack with 1 item, but got a tos of: %d", f.TOS)
	}

	top := pop(&f).(int64)

	if top != 21 {
		t.Errorf("POP2: expected top's value to be 21, but got: %d", top)
	}

	if globals.TraceInst != true {
		t.Errorf("POP2: globals.TraceInst was not re-enabled after the POP2 execution")
	}
}

// POP2: Test underflow error
func TestPop2Underflow(t *testing.T) {
	globals.InitGlobals("test")

	// hide the error message to stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.POP2)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "stack underflow in POP") {
		t.Errorf("Did not get expected error from invalid POP, got: %s", errMsg)
	}
}

// PUSH: Push a value on the op stack
func TestPushWithStackOverflow(t *testing.T) {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("testWithoutShutdown")
	gfunction.InitializeGlobalThreadGroups()

	gl := globals.GetGlobalRef()

	gl.FuncInstantiateClass = InstantiateClass
	gl.FuncThrowException = exceptions.ThrowExNil
	gl.FuncFillInStackTrace = gfunction.FillInStackTrace
	gl.FuncInvokeGFunction = gfunction.Invoke
	gl.FuncMinimalAbort = exceptions.MinimalAbort
	gl.FuncRunThread = RunJavaThread

	stringPool.PreloadArrayClassesToStringPool()
	trace.Init()

	err := classloader.Init()
	if err != nil {
		t.Fail()
	}
	classloader.LoadBaseClasses()

	// initialize the MTable (table caching methods)
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	classloader.LoadBaseClasses()
	_ = classloader.LoadClassFromNameOnly("java/lang/Object")
	classloader.FetchMethodAndCP("java/lang/Object", "wait", "(JI)V")

	var f *frames.Frame
	thObj := gfunction.ThreadCreateNoarg(nil).(*object.Object)
	main := object.StringObjectFromGoString("main")
	params := []any{thObj, main}
	gfunction.ThreadInitWithName(params)
	gfunction.RegisterThread(thObj) // put into globals.Threads map
	thID := int(thObj.FieldTable["ID"].Fvalue.(int64))
	f = frames.CreateFrame(1)
	f.ClName = "java/lang/Object"
	f.MethName = "wait"
	f.MethType = "(JI)V"
	_, err = classloader.FetchMethodAndCP(f.ClName, f.MethName, f.MethType)
	for i := 0; i < 4; i++ {
		f.OpStack = append(f.OpStack, int64(0))
	}
	f.TOS = 4
	f.Thread = thID
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	// Attach the JVM frame stack to the Java thread object
	thObj.FieldTable["framestack"] = object.Field{Ftype: types.LinkedList, Fvalue: fs}

	push(f, int64(34))

	_ = w.Close()
	out, _ := io.ReadAll(r)

	_ = wout.Close()
	// txt, _ := io.ReadAll(rout)

	os.Stderr = normalStderr
	os.Stdout = normalStdout

	msg := string(out[:])
	// woutMsg := string(txt[:])

	if !strings.Contains(msg, "exceeded op stack size of 5") {
		t.Errorf("got unexpected error message: %s", msg)
	}
}

// PUTFIELD: Update a non-static field
func TestPutFieldSimpleInt(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    false,
		IsFinal:     false,
		ClName:      "",
		FldName:     "value",
		FldType:     types.Int,
	}

	f.CP = &CP

	// now create the object we're updating, with one int field
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(42),
	}
	push(&f, obj)

	push(&f, int64(26)) // update the field to 26

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("PUTFIELD: Got unexpected error msg: %s", errMsg)
	}

	res := obj.FieldTable["value"].Fvalue.(int64)
	if res != 26 {
		t.Errorf("PUTFIELD: Expected a new value of 26, got: %d", res)
	}
}

// PUTFIELD for a double
func TestPutFieldDouble(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    false,
		IsFinal:     false,
		ClName:      "",
		FldName:     "value",
		FldType:     types.Double,
	}

	f.CP = &CP

	// now create the object we're updating, with one int field
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.Double,
		Fvalue: float64(42.0),
	}
	push(&f, obj)

	push(&f, float64(26.8)) // update the field to 26.8

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Errorf("PUTFIELD: Got unexpected error msg: %s", errMsg)
	}

	res := obj.FieldTable["value"].Fvalue.(float64)
	if res != 26.8 {
		t.Errorf("PUTFIELD: Expected a new value of 26.8, got: %f", res)
	}
}

// PUTFIELD: Update a field in an object -- error doesn't point to a field
func TestPutFieldNonFieldCPentry(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: 8, Slot: 0} // point to non-fieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "PUTFIELD: Expected a field ref, but") {
		t.Errorf("PUTFIELD: Did not get expected error msg: %s", msg)
	}
}

// PUTFIELD: Error: attempt to update a static field (which should be done by PUTSTATIC, not PUTFIELD)
func TestPutFieldErrorUpdatingStatic(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "",
		FldName:     "value",
		FldType:     types.Static + types.Int,
	}

	f.CP = &CP

	// now create the object we're updating, with one int field
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.Static + types.Int,
		Fvalue: int64(42),
	}
	push(&f, obj)

	push(&f, int64(26)) // update the field to 26

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("PUTFIELD: Expected error message but got none")
	}

	if !strings.Contains(errMsg, "invalid attempt to update a static variable") {
		t.Errorf("PUTFIELD: Did not get expected error message, got %s", errMsg)
	}
}

// PUTSTATIC: Update a static field, an int, successfully
func TestPutStaticInt(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = false

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, int64(420))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "I",
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "I",
		Value: 42,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64)
	if val != 420 {
		t.Errorf("PUTSTATIC: Expected static value to be 420, got: %d", val)
	}
}

// PUTSTATIC: Update a static field, an int, successfully (same as previous test, with tracing on)
func TestPutStaticIntWithTrace(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = true

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, int64(420))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "I",
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "I",
		Value: 42,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "PUTSTATIC") || !strings.Contains(errMsg, "field1") {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64)
	if val != 420 {
		t.Errorf("PUTSTATIC: Expected static value to be 420, got: %d", val)
	}
}

// PUTSTATIC: Update a static field, a boolean, successfully
func TestPutStaticBool(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = false

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, types.JavaBoolTrue)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Bool,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "Z",
		Value: types.JavaBoolFalse,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64)
	if val != types.JavaBoolTrue {
		t.Errorf("PUTSTATIC: Expected static value to be true (1), got: %d", val)
	}
}

// PUTSTATIC: Update a static field, a byte, successfully
func TestPutStaticByte(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = false

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, byte('A'))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Byte,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  types.Byte,
		Value: byte('B'),
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64) // GeStaticValue converts bytes to int64s
	if rune(val) != 'A' {
		t.Errorf("PUTSTATIC: Expected static value to be 'A', got: %c", rune(val))
	}
}

// PUTSTATIC: Update a static field, a byte, successfully
func TestPutStaticJavaByte(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = false

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, types.JavaByte('A'))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Byte,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  types.Byte,
		Value: types.JavaByte('B'),
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64) // GeStaticValue converts bytes to int64s
	if rune(val) != 'A' {
		t.Errorf("PUTSTATIC: Expected static value to be 'A', got: %c", rune(val))
	}
}

// PUTSTATIC: Update a static field, a byte, successfully. This byte value is passed in as int64
func TestPutStaticByteAsInt64(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = false

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, int64('A'))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Byte,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  types.Byte,
		Value: int64('B'),
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64) // GeStaticValue converts bytes to int64s
	if rune(val) != 'A' {
		t.Errorf("PUTSTATIC: Expected static value to be 'A', got: %c", rune(val))
	}
}

// PUTSTATIC: Update a static field, an float/double, successfully
func TestPutStaticFloat(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = false

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, float64(420.1))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "F",
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "F",
		Value: 42.0,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(float64)
	if val != 420.1 {
		t.Errorf("PUTSTATIC: Expected static value to be 420.9, got: %f", val)
	}
}

// PUTSTATIC: this should bonk because the class of the static cannot be found/loaded
func TestPutStaticInvalidNoSuchClass(t *testing.T) {
	globals.InitGlobals("test")
	globals.TraceInst = true

	classloader.InitMethodArea()
	statics.Statics = make(map[string]statics.Static)

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(&f, float64(420.1))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "F",
	}
	f.CP = &CP

	ret := doPutStatic(&f, 0)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	_, _ = io.ReadAll(r)
	os.Stderr = normalStderr

	if ret != exceptions.ERROR_OCCURRED {
		t.Errorf("TestPutStaticInvalidNoSuchClass: Expected ret=exceptions.ERROR_OCCURRED, observed: %d", ret)
		t.Log(string(msg))
	}
}

// PUTSTATIC: Update a static field -- invalid b/c does not point to a field ref in the CP
func TestPutStaticInvalid(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{}
	f.CP = &CP

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("TestPutStaticInvalidNoSuchClass: Expected error message but errMsg is \"\".")
	} else {
		expected := "Expected a field ref, but got"
		if !strings.Contains(errMsg, expected) {
			t.Errorf("TestPutStaticInvalidNoSuchClass: expected: %s, observed: %s", expected, errMsg)
		}
	}
}

// RET: the complement to JSR. The wide version of RET is tested farther below with
// the other WIDE bytecodes
func TestRET(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.RET)
	f.Meth = append(f.Meth, 0x02) // index pointing to local variable 2
	f.PC = 0
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(0), int64(0), int64(456))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.PC-1 != 455 { // -1 because PC++ after processing RET
		t.Errorf("WIDE,RET: expected frame PC value to be 455, got: %d", f.PC)
	}
}

// RETURN: Does a function return correctly?
func TestReturn(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.RETURN)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != -1 {
		t.Errorf("Top of stack, expected -1, got: %d", f.TOS)
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg != "" {
		t.Error("RETURN: Expected popped value to be nil, got: " + errMsg)
	}
}

// SIPUSH: create int from next two bytes and push the int
func TestSipush(t *testing.T) {
	f := newFrame(opcodes.SIPUSH)
	f.Meth = append(f.Meth, 0x01)
	f.Meth = append(f.Meth, 0x02)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value != 258 {
		t.Errorf("SIPUSH: Expected popped value to be 258, got: %d", value)
	}

	if f.PC != 3 {
		t.Errorf("SIPUSH: Expected PC to be 3, got: %d", f.PC)
	}
}

// SIPUSH: create a negative int from next two bytes and push the int
func TestSipushNegative(t *testing.T) {
	f := newFrame(opcodes.SIPUSH)
	val := -1
	f.Meth = append(f.Meth, byte(val))
	f.Meth = append(f.Meth, 0x02)
	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	if f.TOS != 0 {
		t.Errorf("BIPUSH: Top of stack, expected 0, got: %d", f.TOS)
	}
	value := pop(&f).(int64)
	if value >= 0 {
		t.Errorf("SIPUSH: Expected popped value to be negative, got: %d", value)
	}

	if f.PC != 3 {
		t.Errorf("SIPUSH: Expected PC to be 3, got: %d", f.PC)
	}
}

// SWAP: Swap top two items on stack
func TestSwap(t *testing.T) {
	// set the logger to low granularity, so that logging messages are not also captured in this test

	f := newFrame(opcodes.SWAP)
	push(&f, int64(34)) // push two different values
	push(&f, int64(21)) // TOS now = 21

	fs := frames.CreateFrameStack()
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	top := pop(&f).(int64)
	next := pop(&f).(int64)

	if top != 34 {
		t.Errorf("SWAP: expected top's value to be 34, but got: %d", top)
	}

	if next != 21 {
		t.Errorf("SWAP: expected next's value to be 21, but got: %d", next)
	}

	if f.TOS != -1 {
		t.Errorf("SWAP: Expected an empty stack, but got a tos of: %d", f.TOS)
	}
}

// TABLESWITCH: Test with index matching high value
func TestTableswitchMatchHigh(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.TABLESWITCH)

	push(&f, int64(7)) // matches high value

	// Padding bytes
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00)

	// Default jump offset
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x64) // default: 100

	// Low value: 5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x05)

	// High value: 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x07)

	// Jump offsets for values 5, 6, 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x14) // case 5: jump 20
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x28) // case 6: jump 40
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x3C) // case 7: jump 60

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	if f.PC != 60 {
		t.Errorf("TABLESWITCH: Expected jump offset 60 for index 7, got: %d", f.PC)
	}

	if f.TOS != -1 {
		t.Errorf("TABLESWITCH: Expected empty stack, got TOS: %d", f.TOS)
	}
}

// TABLESWITCH: Test with index matching middle value
func TestTableswitchMatchMiddle(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.TABLESWITCH)

	push(&f, int64(6)) // matches middle value

	// Padding bytes
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00)

	// Default jump offset
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x64) // default: 100

	// Low value: 5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x05)

	// High value: 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x07)

	// Jump offsets for values 5, 6, 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x14) // case 5: jump 20
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x28) // case 6: jump 40
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x3C) // case 7: jump 60

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	if f.PC != 40 {
		t.Errorf("TABLESWITCH: Expected jump offset 40 for index 6, got: %d", f.PC)
	}

	if f.TOS != -1 {
		t.Errorf("TABLESWITCH: Expected empty stack, got TOS: %d", f.TOS)
	}
}

// TABLESWITCH: Test with index below low value (default case)
func TestTableswitchDefaultBelowLow(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.TABLESWITCH)

	push(&f, int64(3)) // below low value, should use default

	// Padding bytes
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00)

	// Default jump offset
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x64) // default: 100

	// Low value: 5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x05)

	// High value: 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x07)

	// Jump offsets for values 5, 6, 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x14) // case 5: jump 20
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x28) // case 6: jump 40
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x3C) // case 7: jump 60

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	if f.PC != 100 {
		t.Errorf("TABLESWITCH: Expected default jump offset 100 for index 3, got: %d", f.PC)
	}

	if f.TOS != -1 {
		t.Errorf("TABLESWITCH: Expected empty stack, got TOS: %d", f.TOS)
	}
}

// TABLESWITCH: Test with index above high value (default case)
func TestTableswitchDefaultAboveHigh(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.TABLESWITCH)

	push(&f, int64(10)) // above high value, should use default

	// Padding bytes
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00)

	// Default jump offset
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x64) // default: 100

	// Low value: 5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x05)

	// High value: 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x07)

	// Jump offsets for values 5, 6, 7
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x14) // case 5: jump 20
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x28) // case 6: jump 40
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x3C) // case 7: jump 60

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	if f.PC != 100 {
		t.Errorf("TABLESWITCH: Expected default jump offset 100 for index 10, got: %d", f.PC)
	}

	if f.TOS != -1 {
		t.Errorf("TABLESWITCH: Expected empty stack, got TOS: %d", f.TOS)
	}
}

// TABLESWITCH: Test with single case (low == high)
func TestTableswitchSingleCase(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.TABLESWITCH)

	push(&f, int64(5)) // only case

	// Padding bytes
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00)

	// Default jump offset
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x64) // default: 100

	// Low value: 5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x05)

	// High value: 5 (same as low)
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x05)

	// Jump offset for value 5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x2A) // case 5: jump 42

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	if f.PC != 42 {
		t.Errorf("TABLESWITCH: Expected jump offset 42 for single case, got: %d", f.PC)
	}

	if f.TOS != -1 {
		t.Errorf("TABLESWITCH: Expected empty stack, got TOS: %d", f.TOS)
	}
}

// TABLESWITCH: Test with negative jump offset (backward jump)
func TestTableswitchNegativeJump(t *testing.T) {
	globals.InitGlobals("test")
	f := newFrame(opcodes.TABLESWITCH)

	push(&f, int64(1)) // match case 1

	// Padding bytes
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00)

	// Default jump offset
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x0A) // default: 10

	// Low value: 1
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x01)

	// High value: 2
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x02)

	// Jump offsets - negative value for backward jump
	// -5 in two's complement 32-bit: 0xFFFFFFFB
	f.Meth = append(f.Meth, 0xFF, 0xFF, 0xFF, 0xFB) // case 1: jump -5
	f.Meth = append(f.Meth, 0x00, 0x00, 0x00, 0x0F) // case 2: jump 15

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)
	interpret(fs)

	if f.PC != -5 {
		t.Errorf("TABLESWITCH: Expected jump offset -5 for negative jump, got: %d", f.PC)
	}

	if f.TOS != -1 {
		t.Errorf("TABLESWITCH: Expected empty stack, got TOS: %d", f.TOS)
	}
}

// WIDE version of DLOAD
func TestWideDLOAD(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.DLOAD)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 1
	f.Meth = append(f.Meth, 0x01)
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, float64(0), float64(33.3), float64(0))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	ret := pop(&f).(float64)
	if ret != 33.3 {
		t.Errorf("WIDE,DLOAD: expected return of 33.3, got: %f", ret)
	}

}

// WIDE version of DSTORE
func TestWideDSTORE(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.DSTORE)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 1
	f.Meth = append(f.Meth, 0x01)
	f.TOS = 1 // top of stack = 1 b/c two values are pushed for longs
	f.OpStack[1] = float64(26.2)
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, float64(0), float64(0), float64(0))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	ret := f.Locals[1]
	if ret != float64(26.2) {
		t.Errorf("WIDE,ILOAD: expected locals[1] value to be 26.2, got: %f", ret)
	}
}

// WIDE version of IINC
func TestWideIINC(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.IINC)
	f.Meth = append(f.Meth, 0x00) // index = 2, i.e. locals[2]
	f.Meth = append(f.Meth, 0x02)

	f.Meth = append(f.Meth, 0x00) // amount of increment, 0x24 = 36
	f.Meth = append(f.Meth, 0x24)
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(10), int64(20), int64(30))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.Locals[2] != int64(66) {
		t.Errorf("WIDE,IINC: expected result of 66, got: %d", f.Locals[2])
	}
}

// WIDE version of ILOAD (covers FLOAD AND ALOAD as well b/c they use the same logic)
func TestWideILOAD(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.ILOAD)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 2
	f.Meth = append(f.Meth, 0x02)
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(10), int64(20), int64(30))
	fs.PushFront(&f) // push the new frame
	interpret(fs)
	ret := pop(&f).(int64)
	if ret != int64(30) {
		t.Errorf("WIDE,ILOAD: expected return of 30, got: %d", ret)
	}
}

// WIDE version of ISTORE
func TestWideISTORE(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.ISTORE)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 2
	f.Meth = append(f.Meth, 0x02)
	f.TOS = 0
	f.OpStack[0] = int64(25)
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(0), int64(0), int64(0))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	ret := f.Locals[2]
	if ret != int64(25) {
		t.Errorf("WIDE,ILOAD: expected locals[2] value to be 25, got: %d", ret)
	}
}

// WIDE version of LLOAD
func TestWideLLOAD(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.LLOAD)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 1
	f.Meth = append(f.Meth, 0x01)

	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(33), int64(33), int64(0))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	ret := pop(&f).(int64)
	if ret != 33 {
		t.Errorf("WIDE,DLOAD: expected return of 33, got: %d", ret)
	}
}

// WIDE version of LSTORE
func TestWideLSTORE(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.LSTORE)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 1
	f.Meth = append(f.Meth, 0x01)
	f.TOS = 0
	f.OpStack[0] = int64(25)
	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(0), int64(0), int64(0))

	fs.PushFront(&f) // push the new frame
	interpret(fs)

	ret := f.Locals[1]
	if ret != int64(25) {
		t.Errorf("WIDE,ILOAD: expected locals[1] value to be 25, got: %d", ret)
	}
}

// WIDE version of RET
func TestWideRET(t *testing.T) {
	globals.InitGlobals("test")

	f := newFrame(opcodes.WIDE)
	f.Meth = append(f.Meth, opcodes.RET)
	f.Meth = append(f.Meth, 0x00) // index pointing to local variable 2
	f.Meth = append(f.Meth, 0x02)
	f.PC = 0

	fs := frames.CreateFrameStack()
	f.Locals = append(f.Locals, int64(0), int64(0), int64(123456))
	fs.PushFront(&f) // push the new frame
	interpret(fs)

	if f.PC != 123457 { // 123456 + 1 for the WIDE bytecode, which takes up 1 byte
		t.Errorf("WIDE,RET: expected frame PC value to be 123457, got: %d", f.PC)
	}
}

func TestInvalidInstruction(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(252) // an invalid bytecode

	fs := frames.CreateFrameStack()
	fs.PushFront(&f)          // push the new frame
	globals.TraceInst = false // turn off tracing
	interpret(fs)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if !strings.Contains(msg, "Invalid bytecode") {
		t.Errorf("Error message for invalid bytecode not as expected, got: %s", msg)
	}
}
