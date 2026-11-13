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
	"jacobin/src/trace"
	"testing"
)

// setupHelper initializes globals and base class definitions for tests
func setupConcatTest(t *testing.T) {
	t.Helper()
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()
}

// helper to build a minimal valid CPool and Frame for ConcatStringElements
func buildValidConcatFrame(callerClass, invdynType, recipe string, opStackArgs []any) (*frames.Frame, *list.List) {
	cp := &classloader.CPool{}
	// CpIndex layout:
	// 0: Dummy
	// 1: InvokeDynamic (slot 0)
	// 2: NAT type UTF8 reference (slot 0 in Utf8Refs, type 0)
	// 3: Recipe UTF8 (slot 1 in Utf8Refs)
	cp.CpIndex = make([]classloader.CpEntry, 4)
	cp.CpIndex[0] = classloader.CpEntry{Type: classloader.Dummy, Slot: 0}
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.InvokeDynamic, Slot: 0}
	// invdynTypeIndex.Type is used to index cp.Utf8Refs
	cp.CpIndex[2] = classloader.CpEntry{Type: 0, Slot: 0}
	cp.CpIndex[3] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}

	cp.InvokeDynamics = []classloader.InvokeDynamicEntry{
		{BootstrapIndex: 0, NameAndType: 2},
	}
	cp.Utf8Refs = []string{invdynType, recipe}
	cp.Bootstraps = []classloader.BootstrapMethod{
		{
			MethodRef: 1,
			Args:      []uint16{3}, // points to CP index 3 (recipe)
		},
	}

	fr := frames.CreateFrame(10)
	fr.ClName = callerClass
	fr.CP = cp
	// INVOKEDYNAMIC opcode (0xba), followed by CP slot 1 (0x00, 0x01), and 2 trailing bytes
	fr.Meth = []byte{0xBA, 0x00, 0x01, 0x00, 0x00}
	fr.PC = 0

	// Push operands in stack order
	for i, arg := range opStackArgs {
		fr.OpStack[i] = arg
	}
	if len(opStackArgs) > 0 {
		fr.TOS = len(opStackArgs) - 1
	} else {
		fr.TOS = -1
	}

	fs := list.New()
	fs.PushFront(fr)
	return fr, fs
}

func TestLoad_Lang_Invoke_ConcatStringElements(t *testing.T) {
	setupConcatTest(t)
	Load_Lang_Invoke_ConcatStringElements()

	sig := "java/lang/invoke.ConcatStringElements()Ljava/lang/String;"
	gMeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gMeth.ParamSlots != 0 {
		t.Errorf("Expected ParamSlots 0, got %d", gMeth.ParamSlots)
	}
	if !gMeth.NeedsContext {
		t.Errorf("Expected NeedsContext true, got false")
	}
	if gMeth.GFunction == nil {
		t.Errorf("Expected non-nil GFunction")
	}
}
