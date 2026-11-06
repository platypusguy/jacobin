/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"math"
	"strings"
	"testing"
)

// --- ResolveCallSite ---

func TestResolveCallSite_ImplementationPending(t *testing.T) {
	setupMhResolutionTest(t)

	className, methodName, methodSig := "some/BsmHolder", "bsm", "()V"
	registerFakeMethod(className, methodName, methodSig)

	cp := &CPool{}
	cp.CpIndex = make([]CpEntry, 5)
	cp.CpIndex[0] = CpEntry{Type: Dummy, Slot: 0}
	cp.CpIndex[1] = CpEntry{Type: InvokeDynamic, Slot: 0}
	cp.CpIndex[2] = CpEntry{Type: MethodRef, Slot: 0}
	cp.CpIndex[3] = CpEntry{Type: UTF8, Slot: 0}
	cp.CpIndex[4] = CpEntry{Type: UTF8, Slot: 1}

	cp.InvokeDynamics = []InvokeDynamicEntry{
		{BootstrapIndex: 0, NameAndType: 0},
	}
	cp.NameAndTypes = []NameAndTypeEntry{
		{NameIndex: 3, DescIndex: 4},
	}
	cp.Utf8Refs = []string{"targetMethod", "()V"}
	cp.Bootstraps = []BootstrapMethod{
		{MethodRef: 2}, // points to a MethodHandle CP entry... but see below
	}

	// NOTE: the bootstrap method's MethodRef must itself point to a MethodHandle
	// CP entry for ResolveMethodHandle() to succeed; wire that up here.
	cp.CpIndex[2] = CpEntry{Type: MethodHandle, Slot: 0}
	cp.MethodHandles = []MethodHandleEntry{
		{RefKind: 6, RefIndex: 2}, // REF_invokeStatic -- but RefIndex must point to a MethodRef
	}
	// Since RefIndex 2 now holds the MethodHandle itself, adjust: point the handle's
	// RefIndex to a genuine MethodRef entry appended at the end of the CP.
	cp.CpIndex = append(cp.CpIndex, CpEntry{Type: MethodRef, Slot: 0})
	cp.MethodHandles[0].RefIndex = uint16(len(cp.CpIndex) - 1)
	cp.ResolvedMethodRefs = []ResolvedMethodRefEntry{
		{
			ClassIndex:  stringPool.GetStringIndex(&className),
			NameIndex:   stringPool.GetStringIndex(&methodName),
			TypeIndex:   stringPool.GetStringIndex(&methodSig),
			FQNameIndex: stringPool.GetStringIndex(new(string)),
		},
	}

	className2 := "caller/Class"
	MethAreaInsert(className2, &Klass{Status: 'N', Data: &ClData{Name: className2}})

	fr := frames.CreateFrame(1)
	fr.ClName = className2

	_, err := ResolveCallSite(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected 'implementation pending' error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected") {
		t.Errorf("Unexpected error, expected 'unexpected' in error message, got: %v", err)
	}
}

func TestResolveCallSite_ClassNotFound(t *testing.T) {
	setupMhResolutionTest(t)

	cp := &CPool{}
	cp.CpIndex = make([]CpEntry, 2)
	cp.CpIndex[0] = CpEntry{Type: Dummy, Slot: 0}
	cp.CpIndex[1] = CpEntry{Type: InvokeDynamic, Slot: 0}
	cp.InvokeDynamics = []InvokeDynamicEntry{
		{BootstrapIndex: 0, NameAndType: 0},
	}

	fr := frames.CreateFrame(1)
	fr.ClName = "class/not/registered/Anywhere"

	_, err := ResolveCallSite(cp, 1, fr)
	if err == nil {
		t.Errorf("Expected error when the frame's class cannot be found, got nil")
	}
}

// --- ResolveStaticArgument ---

func TestResolveStaticArgument_InvalidIndexZero(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{CpIndex: make([]CpEntry, 3)}
	fr := frames.CreateFrame(1)

	_, err := ResolveStaticArgument(cp, 0, fr)
	if err == nil {
		t.Fatalf("Expected error for index 0, got nil")
	}
	if !strings.Contains(err.Error(), "invalid CP index") {
		t.Errorf("Expected 'invalid CP index' error, got: %v", err)
	}
}

func TestResolveStaticArgument_InvalidIndexNegative(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{CpIndex: make([]CpEntry, 3)}
	fr := frames.CreateFrame(1)

	_, err := ResolveStaticArgument(cp, -1, fr)
	if err == nil {
		t.Fatalf("Expected error for negative index, got nil")
	}
	if !strings.Contains(err.Error(), "invalid CP index") {
		t.Errorf("Expected 'invalid CP index' error, got: %v", err)
	}
}

func TestResolveStaticArgument_InvalidIndexOutOfRange(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{CpIndex: make([]CpEntry, 3)}
	fr := frames.CreateFrame(1)

	_, err := ResolveStaticArgument(cp, 10, fr)
	if err == nil {
		t.Fatalf("Expected error for out-of-range index, got nil")
	}
	if !strings.Contains(err.Error(), "invalid CP index") {
		t.Errorf("Expected 'invalid CP index' error, got: %v", err)
	}
}

func TestResolveStaticArgument_NilCPool(t *testing.T) {
	setupMhResolutionTest(t)
	fr := frames.CreateFrame(1)

	_, err := ResolveStaticArgument(nil, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for nil CPool, got nil")
	}
	if !strings.Contains(err.Error(), "invalid CP index") {
		t.Errorf("Expected 'invalid CP index' error, got: %v", err)
	}
}

func TestResolveStaticArgument_IntConst(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: IntConst, Slot: 0},
	}
	cp.IntConsts = []int32{42}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil object")
	}
	if obj.FieldTable["value"].Fvalue != int64(42) {
		t.Errorf("Expected boxed integer 42, got %v", obj.FieldTable["value"].Fvalue)
	}
}

func TestResolveStaticArgument_LongConst(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: LongConst, Slot: 0},
	}
	cp.LongConsts = []int64{1234567890123}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil object")
	}
	if obj.FieldTable["value"].Fvalue != int64(1234567890123) {
		t.Errorf("Expected boxed long 1234567890123, got %v", obj.FieldTable["value"].Fvalue)
	}
}

func TestResolveStaticArgument_FloatConst(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: FloatConst, Slot: 0},
	}
	cp.Floats = []float32{3.14}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil object")
	}
	val, ok := obj.FieldTable["value"].Fvalue.(float64)
	if !ok || math.Abs(val-float64(float32(3.14))) > 0.001 {
		t.Errorf("Expected boxed float ~3.14, got %v", obj.FieldTable["value"].Fvalue)
	}
}

func TestResolveStaticArgument_DoubleConst(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: DoubleConst, Slot: 0},
	}
	cp.Doubles = []float64{2.71828}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil object")
	}
	val, ok := obj.FieldTable["value"].Fvalue.(float64)
	if !ok || math.Abs(val-2.71828) > 0.00001 {
		t.Errorf("Expected boxed double 2.71828, got %v", obj.FieldTable["value"].Fvalue)
	}
}

func TestResolveStaticArgument_StringConst(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: UTF8, Slot: 0},
		{Type: StringConst, Slot: 1},
	}
	cp.Utf8Refs = []string{"hello from StringConst"}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 2, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil object")
	}
	str := object.GoStringFromStringObject(obj)
	if str != "hello from StringConst" {
		t.Errorf("Expected 'hello from StringConst', got '%s'", str)
	}
}

func TestResolveStaticArgument_UTF8(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: UTF8, Slot: 0},
	}
	cp.Utf8Refs = []string{"direct utf8 string"}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil object")
	}
	str := object.GoStringFromStringObject(obj)
	if str != "direct utf8 string" {
		t.Errorf("Expected 'direct utf8 string', got '%s'", str)
	}
}

func TestResolveStaticArgument_ClassRef_Success(t *testing.T) {
	setupMhResolutionTest(t)
	className := "java/lang/Object"
	classIdx := stringPool.GetStringIndex(&className)

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0},
	}
	cp.ClassRefs = []uint32{classIdx}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil Class object")
	}
}

func TestResolveStaticArgument_ClassRef_LoadClassFails(t *testing.T) {
	setupMhResolutionTest(t)
	className := "nonexistent/class/ForTesting"
	classIdx := stringPool.GetStringIndex(&className)

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0},
	}
	cp.ClassRefs = []uint32{classIdx}

	fr := frames.CreateFrame(1)
	_, err := ResolveStaticArgument(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error when class loading fails, got nil")
	}
	if !strings.Contains(err.Error(), "could not load class") {
		t.Errorf("Expected error to contain 'could not load class', got: %v", err)
	}
}

func TestResolveStaticArgument_ClassRef_ArrayType_ObjectArray(t *testing.T) {
	setupMhResolutionTest(t)
	className := "[Ljava/lang/String;"
	classIdx := stringPool.GetStringIndex(&className)

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0},
	}
	cp.ClassRefs = []uint32{classIdx}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error resolving array class: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil Class object for array class")
	}
}

func TestResolveStaticArgument_ClassRef_ArrayType_PrimitiveArray(t *testing.T) {
	setupMhResolutionTest(t)
	className := "[I"
	classIdx := stringPool.GetStringIndex(&className)

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0},
	}
	cp.ClassRefs = []uint32{classIdx}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error resolving primitive array class: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil Class object for primitive array class")
	}
}

func TestResolveStaticArgument_ClassRef_ArrayType_MultiDimensionalArray(t *testing.T) {
	setupMhResolutionTest(t)
	className := "[[Ljava/lang/Object;"
	classIdx := stringPool.GetStringIndex(&className)

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0},
	}
	cp.ClassRefs = []uint32{classIdx}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error resolving multi-dimensional array class: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil Class object for multi-dimensional array class")
	}
}

func TestResolveStaticArgument_ClassRef_ArrayType_ResolutionFails(t *testing.T) {
	setupMhResolutionTest(t)
	globals.GetGlobalRef().FuncInvokeGFunction = func(_ string, _ []interface{}) interface{} {
		return nil // simulate Class.forName failure
	}
	className := "[Lsome/InvalidClass;"
	classIdx := stringPool.GetStringIndex(&className)

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0},
	}
	cp.ClassRefs = []uint32{classIdx}

	fr := frames.CreateFrame(1)
	_, err := ResolveStaticArgument(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error when array class resolution fails, got nil")
	}
	if !strings.Contains(err.Error(), "could not resolve array class") {
		t.Errorf("Expected error to contain 'could not resolve array class', got: %v", err)
	}
}

func TestResolveStaticArgument_MethodType_Success(t *testing.T) {
	setupMhResolutionTest(t)
	descriptor := "(Ljava/lang/String;)I"

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: UTF8, Slot: 0},
		{Type: MethodType, Slot: 0},
	}
	cp.Utf8Refs = []string{descriptor}
	cp.MethodTypes = []uint16{1}

	fr := frames.CreateFrame(1)
	obj, err := ResolveStaticArgument(cp, 2, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil MethodType object")
	}
}

func TestResolveStaticArgument_MethodType_Error(t *testing.T) {
	setupMhResolutionTest(t)
	globals.GetGlobalRef().FuncInvokeGFunction = func(_ string, _ []interface{}) interface{} {
		return nil
	}
	descriptor := "(Ljava/lang/String;)I"

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: UTF8, Slot: 0},
		{Type: MethodType, Slot: 0},
	}
	cp.Utf8Refs = []string{descriptor}
	cp.MethodTypes = []uint16{1}

	fr := frames.CreateFrame(1)
	_, err := ResolveStaticArgument(cp, 2, fr)
	if err == nil {
		t.Fatalf("Expected error when MethodType resolution fails, got nil")
	}
}

func TestResolveStaticArgument_MethodHandle_Success(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/StaticTargetClass", "staticTargetMethod", "()V"
	registerFakeMethod(className, methodName, methodSig)

	cp := buildMethodHandleCP(6, className, methodName, methodSig)
	fr := frames.CreateFrame(1)

	obj, err := ResolveStaticArgument(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("Expected non-nil MethodHandle object")
	}
}

func TestResolveStaticArgument_MethodHandle_Error(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: MethodHandle, Slot: 0},
	}
	cp.MethodHandles = []MethodHandleEntry{
		{RefKind: 6, RefIndex: 99}, // invalid RefIndex
	}

	fr := frames.CreateFrame(1)
	_, err := ResolveStaticArgument(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for invalid MethodHandle CP entry, got nil")
	}
}

func TestResolveStaticArgument_InvalidEntryType(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: InvokeDynamic, Slot: 0},
	}
	cp.InvokeDynamics = []InvokeDynamicEntry{
		{BootstrapIndex: 0, NameAndType: 0},
	}

	fr := frames.CreateFrame(1)
	_, err := ResolveStaticArgument(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for non-static-argument entry type, got nil")
	}
	if !strings.Contains(err.Error(), "not a valid bootstrap static argument") {
		t.Errorf("Expected error containing 'not a valid bootstrap static argument', got: %v", err)
	}
}
