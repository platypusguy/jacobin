/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"strings"
	"testing"
)

// setupCpResolutionTest initializes globals (and the string pool) before
// each test in this file.
func setupCpResolutionTest(t *testing.T) {
	globals.InitGlobals("test")
}

// buildInterfaceRefCP builds a minimal CPool containing a single
// InterfaceRefEntry (referring to className/methodName/methodSig via
// intermediate ClassRefs/NameAndTypes/Utf8Refs), suitable for
// ResolveCPinterfaceRefs.
func buildInterfaceRefCP(className, methodName, methodSig string) *CPool {
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0}, // index 1: class ref
		{Type: 12, Slot: 0},       // index 2: name and type
		{Type: UTF8, Slot: 0},     // index 3: method name utf8
		{Type: UTF8, Slot: 1},     // index 4: method sig utf8
	}
	cp.ClassRefs = []uint32{stringPool.GetStringIndex(&className)}
	cp.Utf8Refs = []string{methodName, methodSig}
	cp.NameAndTypes = []NameAndTypeEntry{
		{NameIndex: 3, DescIndex: 4},
	}
	cp.InterfaceRefs = []InterfaceRefEntry{
		{ClassIndex: 1, NameAndType: 2},
	}
	cp.MethodRefs = []MethodRefEntry{{ClassIndex: 1, NameAndType: 2}} // needed only to pass the nil check
	return cp
}

// buildMethodRefCP builds a minimal CPool containing a single MethodRefEntry,
// suitable for ResolveCPmethRefs.
func buildMethodRefCP(className, methodName, methodSig string) *CPool {
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0}, // index 1: class ref
		{Type: 12, Slot: 0},       // index 2: name and type
		{Type: UTF8, Slot: 0},     // index 3: method name utf8
		{Type: UTF8, Slot: 1},     // index 4: method sig utf8
	}
	cp.ClassRefs = []uint32{stringPool.GetStringIndex(&className)}
	cp.Utf8Refs = []string{methodName, methodSig}
	cp.NameAndTypes = []NameAndTypeEntry{
		{NameIndex: 3, DescIndex: 4},
	}
	cp.MethodRefs = []MethodRefEntry{
		{ClassIndex: 1, NameAndType: 2},
	}
	return cp
}

// --- ResolveCPinterfaceRefs ---

func TestResolveCPinterfaceRefs_NilCpIndex(t *testing.T) {
	setupCpResolutionTest(t)
	cp := &CPool{CpIndex: nil, MethodRefs: []MethodRefEntry{{}}}

	err := ResolveCPinterfaceRefs(cp)
	if err == nil {
		t.Fatalf("Expected error for nil CpIndex, got nil")
	}
}

func TestResolveCPinterfaceRefs_NilMethodRefs(t *testing.T) {
	setupCpResolutionTest(t)
	cp := &CPool{CpIndex: []CpEntry{{}}, MethodRefs: nil}

	err := ResolveCPinterfaceRefs(cp)
	if err == nil {
		t.Fatalf("Expected error for nil MethodRefs, got nil")
	}
}

func TestResolveCPinterfaceRefs_EmptyInterfaceRefs(t *testing.T) {
	setupCpResolutionTest(t)
	cp := &CPool{
		CpIndex:       []CpEntry{{}},
		MethodRefs:    []MethodRefEntry{{}},
		InterfaceRefs: nil, // no interfaces to resolve
	}

	err := ResolveCPinterfaceRefs(cp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cp.ResolvedInterfaceRefs) != 0 {
		t.Errorf("Expected no resolved interface refs, got %d", len(cp.ResolvedInterfaceRefs))
	}
}

func TestResolveCPinterfaceRefs_Success(t *testing.T) {
	setupCpResolutionTest(t)
	className, methodName, methodSig := "some/Interface", "foo", "()V"
	cp := buildInterfaceRefCP(className, methodName, methodSig)

	err := ResolveCPinterfaceRefs(cp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cp.ResolvedInterfaceRefs) != 1 {
		t.Fatalf("Expected 1 resolved interface ref, got %d", len(cp.ResolvedInterfaceRefs))
	}

	resEntry := cp.ResolvedInterfaceRefs[0]
	if *stringPool.GetStringPointer(resEntry.ClassIndex) != className {
		t.Errorf("Expected class name %q, got %q", className, *stringPool.GetStringPointer(resEntry.ClassIndex))
	}
	if *stringPool.GetStringPointer(resEntry.NameIndex) != methodName {
		t.Errorf("Expected method name %q, got %q", methodName, *stringPool.GetStringPointer(resEntry.NameIndex))
	}
	if *stringPool.GetStringPointer(resEntry.TypeIndex) != methodSig {
		t.Errorf("Expected method signature %q, got %q", methodSig, *stringPool.GetStringPointer(resEntry.TypeIndex))
	}

	expectedFQN := className + "." + methodName + methodSig
	if *stringPool.GetStringPointer(resEntry.FQNameIndex) != expectedFQN {
		t.Errorf("Expected FQN %q, got %q", expectedFQN, *stringPool.GetStringPointer(resEntry.FQNameIndex))
	}
}

func TestResolveCPinterfaceRefs_MultipleEntries(t *testing.T) {
	setupCpResolutionTest(t)
	cp := buildInterfaceRefCP("some/Interface", "foo", "()V")

	// add a second interface entry pointing to the same class/name/type
	cp.InterfaceRefs = append(cp.InterfaceRefs, InterfaceRefEntry{ClassIndex: 1, NameAndType: 2})

	err := ResolveCPinterfaceRefs(cp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cp.ResolvedInterfaceRefs) != 2 {
		t.Fatalf("Expected 2 resolved interface refs, got %d", len(cp.ResolvedInterfaceRefs))
	}
}

// --- ResolveCPmethRefs ---

func TestResolveCPmethRefs_NilCpIndex(t *testing.T) {
	setupCpResolutionTest(t)
	cp := &CPool{CpIndex: nil, MethodRefs: []MethodRefEntry{{}}}

	err := ResolveCPmethRefs(cp)
	if err == nil {
		t.Fatalf("Expected error for nil CpIndex, got nil")
	}
}

func TestResolveCPmethRefs_NilMethodRefs(t *testing.T) {
	setupCpResolutionTest(t)
	cp := &CPool{CpIndex: []CpEntry{{}}, MethodRefs: nil}

	err := ResolveCPmethRefs(cp)
	if err == nil {
		t.Fatalf("Expected error for nil MethodRefs, got nil")
	}
	if !strings.Contains(err.Error(), "invalid constant pool") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestResolveCPmethRefs_EmptyMethodRefsSlice(t *testing.T) {
	setupCpResolutionTest(t)
	// MethodRefs is a non-nil, but empty, slice: the nil check passes,
	// but the resolution loop has nothing to do.
	cp := &CPool{
		CpIndex:    []CpEntry{{}},
		MethodRefs: []MethodRefEntry{},
	}

	err := ResolveCPmethRefs(cp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cp.ResolvedMethodRefs) != 0 {
		t.Errorf("Expected no resolved method refs, got %d", len(cp.ResolvedMethodRefs))
	}
}

func TestResolveCPmethRefs_Success(t *testing.T) {
	setupCpResolutionTest(t)
	className, methodName, methodSig := "some/Class", "bar", "(I)Z"
	cp := buildMethodRefCP(className, methodName, methodSig)

	err := ResolveCPmethRefs(cp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cp.ResolvedMethodRefs) != 1 {
		t.Fatalf("Expected 1 resolved method ref, got %d", len(cp.ResolvedMethodRefs))
	}

	resEntry := cp.ResolvedMethodRefs[0]
	if *stringPool.GetStringPointer(resEntry.ClassIndex) != className {
		t.Errorf("Expected class name %q, got %q", className, *stringPool.GetStringPointer(resEntry.ClassIndex))
	}
	if *stringPool.GetStringPointer(resEntry.NameIndex) != methodName {
		t.Errorf("Expected method name %q, got %q", methodName, *stringPool.GetStringPointer(resEntry.NameIndex))
	}
	if *stringPool.GetStringPointer(resEntry.TypeIndex) != methodSig {
		t.Errorf("Expected method signature %q, got %q", methodSig, *stringPool.GetStringPointer(resEntry.TypeIndex))
	}

	expectedFQN := className + "." + methodName + methodSig
	if *stringPool.GetStringPointer(resEntry.FQNameIndex) != expectedFQN {
		t.Errorf("Expected FQN %q, got %q", expectedFQN, *stringPool.GetStringPointer(resEntry.FQNameIndex))
	}
}

func TestResolveCPmethRefs_MultipleEntries(t *testing.T) {
	setupCpResolutionTest(t)
	cp := buildMethodRefCP("some/Class", "bar", "(I)Z")

	// add a second method ref pointing to the same class/name/type
	cp.MethodRefs = append(cp.MethodRefs, MethodRefEntry{ClassIndex: 1, NameAndType: 2})

	err := ResolveCPmethRefs(cp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(cp.ResolvedMethodRefs) != 2 {
		t.Fatalf("Expected 2 resolved method refs, got %d", len(cp.ResolvedMethodRefs))
	}
}
