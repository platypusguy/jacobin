/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
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
	"jacobin/src/types"
	"testing"
)

func TestLoad_Lang_Object_RegistersMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Object()

	expected := []string{
		"java/lang/Object.<clinit>()V",
		"java/lang/Object.<init>()V",
		"java/lang/Object.equals(Ljava/lang/Object;)Z",
		"java/lang/Object.finalize()V",
		"java/lang/Object.getClass()Ljava/lang/Class;",
		"java/lang/Object.getResourceAsStream(Ljava/lang/String;)Ljava/io/InputStream;",
		"java/lang/Object.hashCode()I",
		"java/lang/Object.notify()V",
		"java/lang/Object.notifyAll()V",
		"java/lang/Object.toString()Ljava/lang/String;",
		"java/lang/Object.wait()V",
		"java/lang/Object.wait(J)V",
		"java/lang/Object.wait(JI)V",
		"java/lang/Object.clone()Ljava/lang/Object;",
	}
	for _, sig := range expected {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("expected %s to be registered", sig)
		}
	}
}

func TestObjectInit_Success(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectInit([]interface{}{obj})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestObjectInit_NilObject(t *testing.T) {
	result := objectInit([]interface{}{object.Null})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
	if gErr.ErrMsg == "" {
		t.Errorf("expected a non-empty error message")
	}
}

func TestObjectInit_WrongType(t *testing.T) {
	result := objectInit([]interface{}{"not an object"})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectGetClass_Success(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()
	_ = classloader.Init()
	classloader.LoadBaseClasses()

	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)

	result := ObjectGetClass([]interface{}{obj})
	if _, ok := result.(*ghelpers.GErrBlk); ok {
		t.Fatalf("expected success, got error: %v", result)
	}

	k := classloader.MethAreaFetch("java/lang/Object")
	if k == nil || k.Data == nil {
		t.Fatalf("test setup failure: java/lang/Object not preloaded")
	}
	if result != k.Data.ClassObject {
		t.Errorf("expected result to be java/lang/Object's ClassObject")
	}
}

func TestObjectGetClass_NilObject(t *testing.T) {
	result := ObjectGetClass([]interface{}{object.Null})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectGetClass_WrongType(t *testing.T) {
	result := ObjectGetClass([]interface{}{42})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectGetClass_ClassNotLoaded(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	className := "test/NotLoadedClass"
	obj := object.MakeEmptyObjectWithClassName(&className)

	result := ObjectGetClass([]interface{}{obj})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
	if gErr.ErrMsg == "" {
		t.Errorf("expected a non-empty error message")
	}
}

func TestObjectToString_Success(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectToString([]interface{}{obj})
	if _, ok := result.(*ghelpers.GErrBlk); ok {
		t.Fatalf("expected success, got error: %v", result)
	}
	strObj, ok := result.(*object.Object)
	if !ok || object.IsNull(strObj) {
		t.Errorf("expected a non-null *object.Object, got %T", result)
	}
}

func TestObjectToString_NilObject(t *testing.T) {
	result := objectToString([]interface{}{object.Null})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectToString_UnsupportedType(t *testing.T) {
	result := objectToString([]interface{}{123})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectHashCode_Success(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectHashCode([]interface{}{obj})
	hc, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}
	_ = hc // any value is fine, just checking type
}

func TestObjectHashCode_NilObject(t *testing.T) {
	result := objectHashCode([]interface{}{object.Null})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectHashCode_UnsupportedType(t *testing.T) {
	result := objectHashCode([]interface{}{"nope"})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectEquals_SameObject(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectEquals([]interface{}{obj, obj})
	if result != types.JavaBoolTrue {
		t.Errorf("expected JavaBoolTrue, got %v", result)
	}
}

func TestObjectEquals_DifferentObjects(t *testing.T) {
	obj1 := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	obj2 := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectEquals([]interface{}{obj1, obj2})
	if result != types.JavaBoolFalse {
		t.Errorf("expected JavaBoolFalse, got %v", result)
	}
}

func TestObjectEquals_NilThis(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectEquals([]interface{}{object.Null, obj})
	if result != types.JavaBoolFalse {
		t.Errorf("expected JavaBoolFalse, got %v", result)
	}
}

func TestObjectEquals_ThatWrongType(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectEquals([]interface{}{obj, "not an object"})
	if result != types.JavaBoolFalse {
		t.Errorf("expected JavaBoolFalse, got %v", result)
	}
}

func makeFrameStackWithThread(threadID int) *list.List {
	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(1)
	f.Thread = threadID
	_ = frames.PushFrame(fs, f)
	return fs
}

func TestObjectWait_BadFrameStack(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectWait([]interface{}{"not a frame stack", obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectWait_NilObject(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	result := objectWait([]interface{}{fs, object.Null})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectWait_NotOwner(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	// obj is not locked by thread 1 at all
	result := objectWait([]interface{}{fs, obj})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
	if gErr.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestObjectWait_TimedOutSuccess(t *testing.T) {
	globals.InitGlobals("test")
	EnsureTGInit()
	gr := globals.GetGlobalRef()
	gr.Threads[1] = ThreadCreateObject(nil)

	fs := makeFrameStackWithThread(1)
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("ObjLock failed: %v", err)
	}
	// 1ms wait, should time out cleanly and return nil
	result := objectWait([]interface{}{fs, obj, int64(1)})
	if result != nil {
		t.Errorf("expected nil (timed-out wait), got %v", result)
	}
}

func TestObjectNotify_BadFrameStack(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectNotify([]interface{}{"not a frame stack", obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectNotify_NilObject(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	result := objectNotify([]interface{}{fs, object.Null})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectNotify_NotOwner(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectNotify([]interface{}{fs, obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectNotify_Success(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("ObjLock failed: %v", err)
	}
	result := objectNotify([]interface{}{fs, obj})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestObjectNotifyAll_BadFrameStack(t *testing.T) {
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectNotifyAll([]interface{}{123, obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectNotifyAll_NilObject(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	result := objectNotifyAll([]interface{}{fs, object.Null})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectNotifyAll_NotOwner(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	result := objectNotifyAll([]interface{}{fs, obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T", result)
	}
}

func TestObjectNotifyAll_Success(t *testing.T) {
	fs := makeFrameStackWithThread(1)
	obj := object.MakeEmptyObjectWithClassName(&types.ObjectClassName)
	if err := obj.ObjLock(1); err != nil {
		t.Fatalf("ObjLock failed: %v", err)
	}
	result := objectNotifyAll([]interface{}{fs, obj})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestArrayGetClass_PrimitiveIntArray(t *testing.T) {
	obj := object.MakeEmptyObject()
	result := arrayGetClass(obj, "[I")
	if result == nil {
		t.Fatalf("expected non-nil result")
	}
	nameFld := result.FieldTable["name"]
	if nameFld.Fvalue != "[I" {
		t.Errorf("expected name=[I, got %v", nameFld.Fvalue)
	}
	componentFld := result.FieldTable["componentType"]
	if componentFld.Fvalue != "int" {
		t.Errorf("expected componentType=int, got %v", componentFld.Fvalue)
	}
	superFld := result.FieldTable["superClass"]
	if superFld.Fvalue != "java/lang/Object" {
		t.Errorf("expected superClass=java/lang/Object, got %v", superFld.Fvalue)
	}
}

func TestArrayGetClass_ObjectArray(t *testing.T) {
	obj := object.MakeEmptyObject()
	result := arrayGetClass(obj, "[Ljava/lang/String;")
	componentFld := result.FieldTable["componentType"]
	if componentFld.Fvalue != "java/lang/String" {
		t.Errorf("expected componentType=java/lang/String, got %v", componentFld.Fvalue)
	}
}

func TestArrayGetClass_AllPrimitiveTypes(t *testing.T) {
	cases := map[string]string{
		"[Z": "boolean",
		"[B": "byte",
		"[C": "char",
		"[D": "double",
		"[F": "float",
		"[I": "int",
		"[J": "long",
		"[S": "short",
	}
	obj := object.MakeEmptyObject()
	for desc, expected := range cases {
		result := arrayGetClass(obj, desc)
		got := result.FieldTable["componentType"].Fvalue
		if got != expected {
			t.Errorf("descriptor %s: expected componentType=%s, got %v", desc, expected, got)
		}
	}
}

func TestArrayGetClass_NestedArray(t *testing.T) {
	obj := object.MakeEmptyObject()
	result := arrayGetClass(obj, "[[I")
	componentFld := result.FieldTable["componentType"]
	if componentFld.Fvalue != "[I" {
		t.Errorf("expected componentType=[I, got %v", componentFld.Fvalue)
	}
}

func TestArrayGetClass_ModifiersAndMisc(t *testing.T) {
	obj := object.MakeEmptyObject()
	result := arrayGetClass(obj, "[I")

	if _, ok := result.FieldTable["fields"]; !ok {
		t.Errorf("expected fields entry to be present")
	}
	if _, ok := result.FieldTable["methods"]; !ok {
		t.Errorf("expected methods entry to be present")
	}
	if _, ok := result.FieldTable["interfaces"]; !ok {
		t.Errorf("expected interfaces entry to be present")
	}
	accessFlags, ok := result.FieldTable["modifiers"].Fvalue.(classloader.AccessFlags)
	if !ok {
		t.Fatalf("expected modifiers Fvalue to be classloader.AccessFlags")
	}
	if !accessFlags.ClassIsPublic || !accessFlags.ClassIsFinal {
		t.Errorf("expected array class to be public and final")
	}
	loaderFld := result.FieldTable["classLoader"]
	if loaderFld.Fvalue != "bootstrap" {
		t.Errorf("expected classLoader=bootstrap, got %v", loaderFld.Fvalue)
	}
}
