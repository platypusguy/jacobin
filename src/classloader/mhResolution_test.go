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
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"strings"
	"testing"
)

// fakeGFuncInvoker returns a function suitable for globals.FuncInvokeGFunction
// that always succeeds by returning an empty object.Object, regardless of the
// gfunction name or the arguments passed to it.
func fakeGFuncInvoker(_ string, _ []interface{}) interface{} {
	return object.MakeEmptyObject()
}

// setupMhResolutionTest initializes globals, the method area, and the MTable
// for each test in this file.
func setupMhResolutionTest(t *testing.T) {
	globals.InitGlobals("test")
	globals.GetGlobalRef().FuncInvokeGFunction = fakeGFuncInvoker
	if err := Init(); err != nil {
		t.Fatalf("setupMhResolutionTest: Init() failed: %v", err)
	}
	LoadBaseClasses()
}

// buildMethodHandleCP builds a minimal CPool containing, at CP index 1, a
// MethodHandle entry with the given refKind pointing to a MethodRef at index 2,
// which in turn refers to className/methodName/methodSig.
func buildMethodHandleCP(refKind uint8, className, methodName, methodSig string) *CPool {
	cp := &CPool{}
	cp.CpIndex = make([]CpEntry, 3)
	cp.CpIndex[0] = CpEntry{Type: Dummy, Slot: 0}
	cp.CpIndex[1] = CpEntry{Type: MethodHandle, Slot: 0}
	cp.CpIndex[2] = CpEntry{Type: MethodRef, Slot: 0}

	cp.MethodHandles = []MethodHandleEntry{
		{RefKind: refKind, RefIndex: 2},
	}

	cp.ResolvedMethodRefs = []ResolvedMethodRefEntry{
		{
			ClassIndex:  stringPool.GetStringIndex(&className),
			NameIndex:   stringPool.GetStringIndex(&methodName),
			TypeIndex:   stringPool.GetStringIndex(&methodSig),
			FQNameIndex: stringPool.GetStringIndex(new(string)),
		},
	}
	return cp
}

// registerFakeMethod inserts a class into the method area and a matching
// native ('G') entry into the MTable, so that FetchMethodAndCP() finds it
// without needing to load a real class from disk.
func registerFakeMethod(className, methodName, methodSig string) {
	MethAreaInsert(className, &Klass{
		Status: 'N',
		Data:   &ClData{Name: className},
	})
	AddEntry(&MTable, className+"."+methodName+methodSig, MTentry{
		Meth:  "fake-native-method", // anything other than a JmEntry
		MType: 'G',
	})
}

func TestResolveMethodHandle_InvalidIndexTooLow(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{CpIndex: make([]CpEntry, 3)}
	fr := frames.CreateFrame(1)

	_, err := ResolveMethodHandle(cp, 0, fr)
	if err == nil {
		t.Errorf("Expected error for index 0, got nil")
	}
}

func TestResolveMethodHandle_InvalidIndexTooHigh(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{CpIndex: make([]CpEntry, 3)}
	fr := frames.CreateFrame(1)

	_, err := ResolveMethodHandle(cp, 10, fr)
	if err == nil {
		t.Errorf("Expected error for out-of-range index, got nil")
	}
}

func TestResolveMethodHandle_WrongEntryType(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0}, // not a MethodHandle
	}
	cp.ClassRefs = []uint32{types.StringPoolObjectIndex}
	fr := frames.CreateFrame(1)

	_, err := ResolveMethodHandle(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for non-MethodHandle CP entry, got nil")
	}
	if !strings.Contains(err.Error(), "is not a MethodHandle") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestResolveMethodHandle_InvalidRefKind(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/Class", "foo", "()V"
	cp := buildMethodHandleCP(0, className, methodName, methodSig) // refKind 0 is invalid
	fr := frames.CreateFrame(1)

	_, err := ResolveMethodHandle(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for invalid reference kind, got nil")
	}
	if !strings.Contains(err.Error(), "invalid reference kind") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestResolveMethodHandle_InvokeStaticSuccess(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/StaticHolder", "staticMeth", "()V"
	registerFakeMethod(className, methodName, methodSig)

	cp := buildMethodHandleCP(6, className, methodName, methodSig) // REF_invokeStatic
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	mho, err := ResolveMethodHandle(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if mho == nil {
		t.Fatalf("Expected a non-nil MethodHandle object")
	}
	if _, ok := mho.FieldTable["$target"]; !ok {
		t.Errorf("Expected MethodHandle object to have a $target field")
	}
}

func TestResolveMethodHandle_InvokeVirtualSuccess(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/InstanceHolder", "instMeth", "()V"
	registerFakeMethod(className, methodName, methodSig)

	cp := buildMethodHandleCP(5, className, methodName, methodSig) // REF_invokeVirtual
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	mho, err := ResolveMethodHandle(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if mho == nil {
		t.Fatalf("Expected a non-nil MethodHandle object")
	}
}

func TestResolveMethodHandle_NewInvokeSpecialConstructor(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/Constructed", "<init>", "()V"
	registerFakeMethod(className, methodName, methodSig)

	cp := buildMethodHandleCP(8, className, methodName, methodSig) // REF_newInvokeSpecial
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	mho, err := ResolveMethodHandle(cp, 1, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if mho == nil {
		t.Fatalf("Expected a non-nil MethodHandle object")
	}
}

func TestResolveMethodHandle_ConstructorWithWrongRefKind(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/BadConstructed", "<init>", "()V"
	registerFakeMethod(className, methodName, methodSig)

	// REF_invokeSpecial (7) used for a constructor should fail: must be REF_newInvokeSpecial (8)
	cp := buildMethodHandleCP(7, className, methodName, methodSig)
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	_, err := ResolveMethodHandle(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for <init> resolved with wrong reference kind, got nil")
	}
	if !strings.Contains(err.Error(), "REF_newInvokeSpecial") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestResolveMethodHandle_ClassInitializerRejected(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "some/HasClinit", "<clinit>", "()V"
	registerFakeMethod(className, methodName, methodSig)

	cp := buildMethodHandleCP(6, className, methodName, methodSig) // REF_invokeStatic
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	_, err := ResolveMethodHandle(cp, 1, fr)
	if err == nil {
		t.Fatalf("Expected error for <clinit> method handle, got nil")
	}
	if !strings.Contains(err.Error(), "class initializer") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestResolveMethodHandle_MethodNotFound(t *testing.T) {
	setupMhResolutionTest(t)
	className, methodName, methodSig := "no/such/Class", "missing", "()V"
	// deliberately do NOT register the method or the class

	cp := buildMethodHandleCP(6, className, methodName, methodSig)
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	_, err := ResolveMethodHandle(cp, 1, fr)
	if err == nil {
		t.Errorf("Expected error for method that cannot be found, got nil")
	}
}

// --- resolveFieldHandle (white-box) ---

func TestResolveFieldHandle_InvalidRefIndex(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{CpIndex: make([]CpEntry, 2)}
	fr := frames.CreateFrame(1)

	_, err := resolveFieldHandle(cp, 0, false, false, fr, 1)
	if err == nil {
		t.Errorf("Expected error for invalid field ref index, got nil")
	}
}

func TestResolveFieldHandle_WrongCPEntryType(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: MethodRef, Slot: 0}, // not FieldRef
	}
	fr := frames.CreateFrame(1)

	_, err := resolveFieldHandle(cp, 1, false, false, fr, 1)
	if err == nil {
		t.Errorf("Expected error for wrong CP entry type, got nil")
	}
}

func TestResolveFieldHandle_InstanceGetterSuccess(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: FieldRef, Slot: 0},
	}
	cp.FieldRefs = []ResolvedFieldEntry{
		{ClName: "some/Class", FldName: "count", FldType: "I"},
	}
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	mho, err := resolveFieldHandle(cp, 1, false, false, fr, 1) // REF_getField
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if mho == nil {
		t.Fatalf("Expected non-nil MethodHandle object")
	}
}

func TestResolveFieldHandle_StaticGetterMissingStatic(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: FieldRef, Slot: 0},
	}
	cp.FieldRefs = []ResolvedFieldEntry{
		{ClName: "some/StaticClass", FldName: "missingField", FldType: "I"},
	}
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	_, err := resolveFieldHandle(cp, 1, true, false, fr, 2) // REF_getStatic, static not registered
	if err == nil {
		t.Errorf("Expected error for missing static field, got nil")
	}
}

func TestResolveFieldHandle_StaticGetterSuccess(t *testing.T) {
	setupMhResolutionTest(t)
	_ = statics.AddStatic("some/StaticClass2.presentField", statics.Static{Type: types.Int, Value: int64(42)})

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: FieldRef, Slot: 0},
	}
	cp.FieldRefs = []ResolvedFieldEntry{
		{ClName: "some/StaticClass2", FldName: "presentField", FldType: "I"},
	}
	fr := frames.CreateFrame(1)
	fr.ClName = "caller/Class"

	mho, err := resolveFieldHandle(cp, 1, true, false, fr, 2) // REF_getStatic
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if mho == nil {
		t.Fatalf("Expected non-nil MethodHandle object")
	}
}

// --- GetPrimitiveClass ---

func TestGetPrimitiveClass_KnownTypes(t *testing.T) {
	setupMhResolutionTest(t)

	tests := map[string]string{
		"B": "java/lang/Byte",
		"C": "java/lang/Character",
		"D": "java/lang/Double",
		"F": "java/lang/Float",
		"I": "java/lang/Integer",
		"J": "java/lang/Long",
		"S": "java/lang/Short",
		"Z": "java/lang/Boolean",
		"V": "java/lang/Void",
	}

	for descriptor, wrapperClass := range tests {
		classObj := object.MakeEmptyObjectWithClassName(&wrapperClass)
		_ = statics.AddStatic(wrapperClass+".TYPE", statics.Static{Type: types.Ref, Value: classObj})

		result := GetPrimitiveClass(descriptor)
		if result == nil {
			t.Errorf("GetPrimitiveClass(%q): expected non-nil result", descriptor)
			continue
		}
		if result != classObj {
			t.Errorf("GetPrimitiveClass(%q): expected the registered object back", descriptor)
		}
	}
}

func TestGetPrimitiveClass_UnknownDescriptor(t *testing.T) {
	setupMhResolutionTest(t)

	result := GetPrimitiveClass("Q") // not a valid descriptor
	if result != nil {
		t.Errorf("Expected nil for unknown descriptor, got %v", result)
	}
}

func TestGetPrimitiveClass_StaticNotRegistered(t *testing.T) {
	setupMhResolutionTest(t)

	// Statics is a package-global map that other tests may have already
	// populated, so explicitly ensure this entry is absent before checking.
	delete(statics.Statics, "java/lang/Byte.TYPE")

	result := GetPrimitiveClass("B")
	if result != nil {
		t.Errorf("Expected nil when the TYPE static isn't registered, got %v", result)
	}
}

func TestGetPrimitiveClass_StaticWrongType(t *testing.T) {
	setupMhResolutionTest(t)
	_ = statics.AddStatic("java/lang/Short.TYPE", statics.Static{Type: types.Ref, Value: "not-an-object"})

	result := GetPrimitiveClass("S")
	if result != nil {
		t.Errorf("Expected nil when the TYPE static value isn't a *object.Object, got %v", result)
	}
}

// --- getClassObj (white-box) ---

func TestGetClassObj_Primitive(t *testing.T) {
	setupMhResolutionTest(t)
	wrapperClass := "java/lang/Integer"
	classObj := object.MakeEmptyObjectWithClassName(&wrapperClass)
	_ = statics.AddStatic(wrapperClass+".TYPE", statics.Static{Type: types.Ref, Value: classObj})

	fr := frames.CreateFrame(1)
	result, err := getClassObj("I", fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != classObj {
		t.Errorf("Expected the registered primitive class object to be returned")
	}
}

func TestGetClassObj_ObjectType(t *testing.T) {
	setupMhResolutionTest(t)
	fr := frames.CreateFrame(1)

	result, err := getClassObj("Ljava/lang/String;", fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("Expected a non-nil Class object")
	}
}

func TestGetClassObj_ArrayType(t *testing.T) {
	setupMhResolutionTest(t)
	fr := frames.CreateFrame(1)

	result, err := getClassObj("[I", fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("Expected a non-nil Class object")
	}
}

func TestGetClassObj_ForNameFails(t *testing.T) {
	setupMhResolutionTest(t)
	globals.GetGlobalRef().FuncInvokeGFunction = func(_ string, _ []interface{}) interface{} {
		return nil // simulate Class.forName failure
	}
	fr := frames.CreateFrame(1)

	_, err := getClassObj("Lno/such/Klass;", fr)
	if err == nil {
		t.Errorf("Expected error when Class.forName fails, got nil")
	}
}

// --- getMethodTypeObject / ResolveMethodType (white-box + exported) ---

func TestGetMethodTypeObject_Success(t *testing.T) {
	setupMhResolutionTest(t)
	fr := frames.CreateFrame(1)

	result, err := getMethodTypeObject("(I)V", fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("Expected a non-nil MethodType object")
	}
}

func TestGetMethodTypeObject_GFunctionFails(t *testing.T) {
	setupMhResolutionTest(t)
	globals.GetGlobalRef().FuncInvokeGFunction = func(_ string, _ []interface{}) interface{} {
		return nil
	}
	fr := frames.CreateFrame(1)

	_, err := getMethodTypeObject("(I)V", fr)
	if err == nil {
		t.Errorf("Expected error when gfunction returns nil, got nil")
	}
}

func TestResolveMethodType_Success(t *testing.T) {
	setupMhResolutionTest(t)
	descriptor := "(Ljava/lang/String;)V"

	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: UTF8, Slot: 0},
		{Type: MethodType, Slot: 0},
	}
	cp.Utf8Refs = []string{descriptor}
	cp.MethodTypes = []uint16{1} // points to CpIndex[1], the UTF8 descriptor

	fr := frames.CreateFrame(1)

	result, err := ResolveMethodType(cp, 2, fr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("Expected a non-nil MethodType object")
	}
}

func TestResolveMethodType_WrongEntryType(t *testing.T) {
	setupMhResolutionTest(t)
	cp := &CPool{}
	cp.CpIndex = []CpEntry{
		{Type: Dummy, Slot: 0},
		{Type: ClassRef, Slot: 0}, // not a MethodType
	}
	cp.ClassRefs = []uint32{types.StringPoolObjectIndex}

	fr := frames.CreateFrame(1)

	_, err := ResolveMethodType(cp, 1, fr)
	if err == nil {
		t.Errorf("Expected error for non-MethodType CP entry, got nil")
	}
}

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
	if !strings.Contains(err.Error(), "implementation pending") {
		t.Errorf("Unexpected error, expected 'implementation pending', got: %v", err)
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
