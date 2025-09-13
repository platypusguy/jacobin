/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"testing"
)

// helper to set CP slot 1 to a ClassRef for the given class name
func setCPClassRefFor(f *frames.Frame, className string) {
	// Place two bytes after opcode to reference CP slot 1
	f.Meth = append(f.Meth, 0)
	f.Meth = append(f.Meth, 1)
	// Build a minimal constant pool with [1] as ClassRef -> string pool index of className
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 2)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	idx := stringPool.GetStringIndex(&className)
	CP.ClassRefs = append(CP.ClassRefs, idx)
	f.CP = &CP
}

// helper to insert a class record in the Method Area if not present
func ensureClassInMethArea(name string) {
	if classloader.MethAreaFetch(name) == nil {
		classloader.MethAreaInsert(name, &classloader.Klass{Status: 'X', Loader: "bootstrap", Data: &classloader.ClData{}})
	}
}

// helper to insert a class with a given superclass linkage into the Method Area
func ensureClassWithSuper(name string, super string) {
	if classloader.MethAreaFetch(name) != nil {
		return
	}
	var superIdx uint32
	if super == "" {
		superIdx = types.InvalidStringIndex
	} else if super == "java/lang/Object" {
		superIdx = types.ObjectPoolStringIndex
	} else {
		superIdx = stringPool.GetStringIndex(&super)
	}
	classloader.MethAreaInsert(name, &classloader.Klass{Status: 'X', Loader: "bootstrap", Data: &classloader.ClData{SuperclassIndex: superIdx}})
}

// runCheckcast runs doCheckcast on obj against targetName and returns whether the cast succeeded (no exception)
// and the object that remains on the stack (should be unchanged on success).
func runCheckcast(t *testing.T, obj interface{}, targetName string) (ok bool, top interface{}, pc int) {
	fr := newFrame(opcodes.CHECKCAST)
	setCPClassRefFor(&fr, targetName)
	push(&fr, obj)
	fs := frames.CreateFrameStack()
	fs.PushFront(&fr)
	interpret(fs)
	pc = fr.PC
	// If no exception occurred, CHECKCAST leaves the object on the stack unchanged.
	if fr.TOS >= 0 {
		top = peek(&fr)
		ok = true
	} else {
		ok = false
	}
	return
}

func TestCheckcastMisc(t *testing.T) {
	// Initialize runtime and minimal classes
	globals.InitGlobals("test")
	_ = classloader.Init()

	// Prepare classes in the method area used as targets
	ensureClassWithSuper("java/lang/Object", "")
	ensureClassWithSuper("java/lang/String", "java/lang/Object")
	ensureClassInMethArea("[Ljava/lang/String")
	ensureClassInMethArea("[Ljava/lang/Object")
	ensureClassWithSuper("MyClass", "java/lang/Object")
	ensureClassInMethArea("[LMyClass")

	// Build objects per the Java program
	strSimple := object.StringObjectFromGoString("ABC")
	strArray := object.Make1DimRefArray("java/lang/String", 1)
	arrVals := strArray.FieldTable["value"].Fvalue.([]*object.Object)
	arrVals[0] = strSimple

	strMyClass := "MyClass"
	objSimple := object.MakeEmptyObjectWithClassName(&strMyClass)
	objArray := object.Make1DimRefArray("MyClass", 1)
	objArrVals := objArray.FieldTable["value"].Fvalue.([]*object.Object)
	objArrVals[0] = objSimple

	// Success cases: CHECKCAST should allow upcasts and leave object unchanged on stack
	ok, top, pc := runCheckcast(t, strSimple, "java/lang/String")
	if !ok || top != strSimple || pc != 3 {
		t.Fatalf("CHECKCAST strSimple->String failed: ok=%v top==obj? %v pc=%d", ok, top == strSimple, pc)
	}

	ok, top, pc = runCheckcast(t, strSimple, "java/lang/Object")
	if !ok || top != strSimple || pc != 3 {
		t.Fatalf("CHECKCAST strSimple->Object should succeed: ok=%v top==obj? %v pc=%d", ok, top == strSimple, pc)
	}

	ok, top, pc = runCheckcast(t, strArray, "[Ljava/lang/String")
	if !ok || top != strArray || pc != 3 {
		t.Fatalf("CHECKCAST strArray->String[] failed: ok=%v top==obj? %v pc=%d", ok, top == strArray, pc)
	}

	ok, top, pc = runCheckcast(t, strArray, "[Ljava/lang/Object")
	if !ok || top != strArray || pc != 3 {
		t.Fatalf("CHECKCAST strArray->Object[] should succeed: ok=%v top==obj? %v pc=%d", ok, top == strArray, pc)
	}

	ok, top, pc = runCheckcast(t, objSimple, strMyClass)
	if !ok || top != objSimple || pc != 3 {
		t.Fatalf("CHECKCAST objSimple->MyClass failed: ok=%v top==obj? %v pc=%d", ok, top == objSimple, pc)
	}

	ok, top, pc = runCheckcast(t, objSimple, "java/lang/Object")
	if !ok || top != objSimple || pc != 3 {
		t.Fatalf("CHECKCAST objSimple->Object should succeed: ok=%v top==obj? %v pc=%d", ok, top == objSimple, pc)
	}

	ok, top, pc = runCheckcast(t, objArray, "[L"+strMyClass)
	if !ok || top != objArray || pc != 3 {
		t.Fatalf("CHECKCAST objArray->MyClass[] failed: ok=%v top==obj? %v pc=%d", ok, top == objArray, pc)
	}

	ok, top, pc = runCheckcast(t, objArray, "[Ljava/lang/Object")
	if !ok || top != objArray || pc != 3 {
		t.Fatalf("CHECKCAST objArray->Object[] should succeed: ok=%v top==obj? %v pc=%d", ok, top == objArray, pc)
	}
}
