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
	"jacobin/src/types"
	"testing"
)

// runInstanceof runs doInstanceof on obj against targetName and returns a result pushed on the stack.
func runInstanceof(t *testing.T, obj interface{}, targetName string) int64 {
	fr := newFrame(opcodes.INSTANCEOF)
	setCPClassRefFor(&fr, targetName)
	push(&fr, obj)
	fs := frames.CreateFrameStack()
	fs.PushFront(&fr)
	interpret(fs)
	val, _ := pop(&fr).(int64)
	return val
}

func TestInstanceofMisc(t *testing.T) {
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
	// place element (not strictly needed for instanceof of array itself)
	arrVals := strArray.FieldTable["value"].Fvalue.([]*object.Object)
	arrVals[0] = strSimple

	strMyClass := "MyClass"
	objSimple := object.MakeEmptyObjectWithClassName(&strMyClass)
	objArray := object.Make1DimRefArray("MyClass", 1)
	objArrVals := objArray.FieldTable["value"].Fvalue.([]*object.Object)
	objArrVals[0] = objSimple

	if got := runInstanceof(t, strSimple, "java/lang/String"); got != types.JavaBoolTrue {
		t.Fatalf("strSimple instanceof String: want 1 got %d", got)
	}
	if got := runInstanceof(t, strSimple, "java/lang/Object"); got != types.JavaBoolTrue {
		t.Fatalf("strSimple instanceof Object: CURRENT want 0 got %d", got)
	}

	if got := runInstanceof(t, strArray, "[Ljava/lang/String"); got != types.JavaBoolTrue {
		t.Fatalf("strArray instanceof String[]: want 1 got %d", got)
	}
	if got := runInstanceof(t, strArray, "[Ljava/lang/Object"); got != types.JavaBoolTrue {
		t.Fatalf("strArray instanceof Object[]: CURRENT want 0 got %d", got)
	}

	if got := runInstanceof(t, objSimple, strMyClass); got != types.JavaBoolTrue {
		t.Fatalf("objSimple instanceof MyClass: want 1 got %d", got)
	}
	if got := runInstanceof(t, objSimple, "java/lang/Object"); got != types.JavaBoolTrue {
		t.Fatalf("objSimple instanceof Object: CURRENT want 0 got %d", got)
	}

	if got := runInstanceof(t, objArray, "[L"+strMyClass); got != types.JavaBoolTrue {
		t.Fatalf("objArray instanceof MyClass[]: want 1 got %d", got)
	}
	if got := runInstanceof(t, objArray, "[Ljava/lang/Object"); got != types.JavaBoolTrue {
		t.Fatalf("objArray instanceof Object[]: CURRENT want 0 got %d", got)
	}
}
