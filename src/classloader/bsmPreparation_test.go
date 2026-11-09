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

// --- marshalBSMargs ---

func TestMarshalBSMArgs_Structure_NonVarargs(t *testing.T) {
	setupMhResolutionTest(t)

	callerClass := "caller/MyClass"
	MethAreaInsert(callerClass, &Klass{Status: 'N', Data: &ClData{Name: callerClass}})

	fr := frames.CreateFrame(1)
	fr.ClName = callerClass

	staticArg1 := object.StringObjectFromGoString("recipe\x01")
	staticArgs := []*object.Object{staticArg1}

	bsmSig := "(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/String;)Ljava/lang/invoke/CallSite;"
	params, err := marshalBSMargs(fr, "testMethod", "(Ljava/lang/String;)Ljava/lang/String;", bsmSig, staticArgs)
	if err != nil {
		t.Fatalf("marshalBSMargs failed: %v", err)
	}

	if len(params) != 4 {
		t.Fatalf("Expected 4 parameters (lookup, name, type, staticArg), got %d", len(params))
	}

	// Param 0: MethodHandles.Lookup object
	lookupObj, ok := params[0].(*object.Object)
	if !ok || lookupObj == nil {
		t.Fatalf("Param 0 is not a non-nil *object.Object")
	}
	lookupClassName := stringPool.GetStringPointer(lookupObj.KlassName)
	if lookupClassName == nil || *lookupClassName != "java/lang/invoke/MethodHandles$Lookup" {
		t.Errorf("Param 0 class name: expected 'java/lang/invoke/MethodHandles$Lookup', got %v", lookupClassName)
	}
	lookupClassField, ok := lookupObj.FieldTable["lookupClass"]
	if !ok || lookupClassField.Fvalue == nil {
		t.Errorf("Param 0 Lookup object missing lookupClass field")
	}

	// Param 1: CallSite name (java.lang.String)
	nameObj, ok := params[1].(*object.Object)
	if !ok || nameObj == nil {
		t.Fatalf("Param 1 is not a non-nil *object.Object")
	}
	if gotName := object.GoStringFromStringObject(nameObj); gotName != "testMethod" {
		t.Errorf("Param 1 name: expected 'testMethod', got '%s'", gotName)
	}

	// Param 2: CallSite MethodType (java.lang.invoke.MethodType)
	typeObj, ok := params[2].(*object.Object)
	if !ok || typeObj == nil {
		t.Fatalf("Param 2 is not a non-nil *object.Object")
	}

	// Param 3: static argument
	if params[3] != staticArg1 {
		t.Errorf("Param 3: expected staticArg1, got %v", params[3])
	}
}

func TestMarshalBSMArgs_Varargs_SingleRecipe(t *testing.T) {
	setupMhResolutionTest(t)

	callerClass := "caller/MyClass"
	MethAreaInsert(callerClass, &Klass{Status: 'N', Data: &ClData{Name: callerClass}})

	fr := frames.CreateFrame(1)
	fr.ClName = callerClass

	recipeArg := object.StringObjectFromGoString("Hello \x01")
	staticArgs := []*object.Object{recipeArg}

	bsmSig := "(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"
	params, err := marshalBSMargs(fr, "makeConcatWithConstants", "(Ljava/lang/String;)Ljava/lang/String;", bsmSig, staticArgs)
	if err != nil {
		t.Fatalf("marshalBSMargs failed: %v", err)
	}

	if len(params) != 5 {
		t.Fatalf("Expected 5 parameters (lookup, name, type, recipe, constants), got %d", len(params))
	}

	// Param 0: Lookup
	if _, ok := params[0].(*object.Object); !ok {
		t.Errorf("Param 0 is not a non-nil *object.Object")
	}
	// Param 1: CallSite name
	if nameObj, ok := params[1].(*object.Object); !ok || object.GoStringFromStringObject(nameObj) != "makeConcatWithConstants" {
		t.Errorf("Param 1 name mismatch: %v", params[1])
	}
	// Param 2: CallSite MethodType
	if _, ok := params[2].(*object.Object); !ok {
		t.Errorf("Param 2 is not a non-nil *object.Object")
	}
	// Param 3: Recipe
	if params[3] != recipeArg {
		t.Errorf("Param 3: expected recipeArg, got %v", params[3])
	}
	// Param 4: Extra constants (empty slice)
	constantsSlice, ok := params[4].([]*object.Object)
	if !ok {
		t.Fatalf("Param 4 is not []*object.Object, got %T", params[4])
	}
	if len(constantsSlice) != 0 {
		t.Errorf("Param 4: expected empty slice, got %d elements", len(constantsSlice))
	}
}

func TestMarshalBSMArgs_Varargs_WithExtraConstants(t *testing.T) {
	setupMhResolutionTest(t)

	callerClass := "caller/MyClass"
	MethAreaInsert(callerClass, &Klass{Status: 'N', Data: &ClData{Name: callerClass}})

	fr := frames.CreateFrame(1)
	fr.ClName = callerClass

	recipeArg := object.StringObjectFromGoString("Answer: \x01, Code: \x02")
	const1 := object.StringObjectFromGoString("extraConst1")
	const2 := object.StringObjectFromGoString("extraConst2")
	staticArgs := []*object.Object{recipeArg, const1, const2}

	bsmSig := "(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"
	params, err := marshalBSMargs(fr, "makeConcatWithConstants", "(Ljava/lang/String;)Ljava/lang/String;", bsmSig, staticArgs)
	if err != nil {
		t.Fatalf("marshalBSMargs failed: %v", err)
	}

	if len(params) != 5 {
		t.Fatalf("Expected 5 parameters, got %d", len(params))
	}

	if params[3] != recipeArg {
		t.Errorf("Param 3: expected recipeArg, got %v", params[3])
	}

	constantsSlice, ok := params[4].([]*object.Object)
	if !ok {
		t.Fatalf("Param 4 is not []*object.Object, got %T", params[4])
	}
	if len(constantsSlice) != 2 {
		t.Fatalf("Param 4: expected 2 elements, got %d", len(constantsSlice))
	}
	if constantsSlice[0] != const1 || constantsSlice[1] != const2 {
		t.Errorf("Param 4 elements mismatch: got %v", constantsSlice)
	}
}

func TestMarshalBSMArgs_Varargs_NoStaticArgs(t *testing.T) {
	setupMhResolutionTest(t)

	callerClass := "caller/MyClass"
	MethAreaInsert(callerClass, &Klass{Status: 'N', Data: &ClData{Name: callerClass}})

	fr := frames.CreateFrame(1)
	fr.ClName = callerClass

	bsmSig := "(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"
	params, err := marshalBSMargs(fr, "makeConcatWithConstants", "()Ljava/lang/String;", bsmSig, nil)
	if err != nil {
		t.Fatalf("marshalBSMargs failed: %v", err)
	}

	if len(params) != 5 {
		t.Fatalf("Expected 5 parameters, got %d", len(params))
	}

	recipeObj, ok := params[3].(*object.Object)
	if !ok || object.GoStringFromStringObject(recipeObj) != "" {
		t.Errorf("Param 3: expected empty string object, got %v", params[3])
	}

	constantsSlice, ok := params[4].([]*object.Object)
	if !ok || len(constantsSlice) != 0 {
		t.Errorf("Param 4: expected empty slice, got %v", params[4])
	}
}

func TestMarshalBSMArgs_ClassNotFound(t *testing.T) {
	setupMhResolutionTest(t)

	globals.GetGlobalRef().FuncInvokeGFunction = func(_ string, _ []interface{}) interface{} {
		return nil
	}

	fr := frames.CreateFrame(1)
	fr.ClName = "unregistered/CallerClass"

	_, err := marshalBSMargs(fr, "testMethod", "()V", "()V", nil)
	if err == nil {
		t.Fatalf("Expected error when caller class object cannot be resolved, got nil")
	}
	if !strings.Contains(err.Error(), "could not get Class object for caller") {
		t.Errorf("Expected 'could not get Class object for caller' in error message, got: %v", err)
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
