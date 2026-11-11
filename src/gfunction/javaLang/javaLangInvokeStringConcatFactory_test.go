/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"testing"
)

func TestStringConcatFactory_MakeConcatWithConstants_PtypesKlass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()

	// Prepare method type with parameter types (e.g. String, int)
	integerClinit(nil)
	mtObjRaw := MethodTypeFromMethodDescriptorString([]interface{}{
		object.StringObjectFromGoString("(Ljava/lang/String;I)Ljava/lang/String;"),
		nil,
	})
	mtObj, ok := mtObjRaw.(*object.Object)
	if !ok || mtObj == nil {
		t.Fatalf("Failed to create MethodType object: %v", mtObjRaw)
	}

	lookupObj := object.StringObjectFromGoString("caller/MyClass")
	nameObj := object.StringObjectFromGoString("makeConcatWithConstants")
	recipeObj := object.StringObjectFromGoString("Hello \x01, score \x01")
	constants := []*object.Object{}

	params := []any{
		lookupObj,
		nameObj,
		mtObj,
		recipeObj,
		constants,
	}

	result := stringConcatFactory_makeConcatWithConstants(params)
	if result != nil {
		t.Errorf("Expected nil return, got %v", result)
	}
}

func TestStringConcatFactory_MakeConcatWithConstants_IncorrectParamCount(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []any{
		object.StringObjectFromGoString("caller/MyClass"),
		object.StringObjectFromGoString("makeConcatWithConstants"),
	}

	result := stringConcatFactory_makeConcatWithConstants(params)
	if result == nil {
		t.Errorf("Expected error block for incorrect parameter count, got nil")
	}
}

func TestConcat1param_Success(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()

	callerClass := "caller/MyClass"
	classloader.MethAreaInsert(callerClass, &classloader.Klass{Status: 'N', Data: &classloader.ClData{Name: callerClass}})

	cp := &classloader.CPool{}
	cp.CpIndex = make([]classloader.CpEntry, 3)
	cp.CpIndex[0] = classloader.CpEntry{Type: classloader.Dummy, Slot: 0}
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.InvokeDynamic, Slot: 0}
	cp.CpIndex[2] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}

	cp.InvokeDynamics = []classloader.InvokeDynamicEntry{
		{BootstrapIndex: 0, NameAndType: 0},
	}
	cp.Utf8Refs = []string{"Hello \x01"}
	cp.Bootstraps = []classloader.BootstrapMethod{
		{
			MethodRef: 1,
			Args:      []uint16{2}, // points to CP index 2 (UTF8 "Hello \x01")
		},
	}

	fr := frames.CreateFrame(2)
	fr.ClName = callerClass
	fr.CP = cp
	// INVOKEDYNAMIC opcode (0xba), followed by CP slot 1 (0x00, 0x01), and 2 dummy bytes (0x00, 0x00)
	fr.Meth = []byte{0xBA, 0x00, 0x01, 0x00, 0x00}
	fr.PC = 0

	fs := list.New()
	fs.PushFront(fr)

	// params: [arg0 (e.g. type or item), params[1] = frame stack]
	params := []any{
		object.StringObjectFromGoString("dummyParam"),
		fs,
	}

	res := ConcatStringElements(params)
	if res != nil {
		t.Fatalf("Expected nil return on success, got %v", res)
	}
}

func TestConcat1param_InvalidParams(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Case 1: Incorrect param count
	res1 := ConcatStringElements([]any{list.New()})
	if res1 == nil {
		t.Errorf("Expected error block for incorrect parameter count, got nil")
	}

	// Case 2: params[1] is not a list
	res2 := ConcatStringElements([]any{"arg0", "not-a-list"})
	if _, ok := res2.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for invalid frame stack, got %T", res2)
	}

	// Case 3: empty frame stack
	emptyFs := list.New()
	res3 := ConcatStringElements([]any{"arg0", emptyFs})
	if _, ok := res3.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected *ghelpers.GErrBlk for empty frame stack, got %T", res3)
	}
}
