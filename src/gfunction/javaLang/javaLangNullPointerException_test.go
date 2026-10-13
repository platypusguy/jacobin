/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
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
	"jacobin/src/types"
	"testing"
)

// setupNullPointerExceptionTest builds a minimal frame/stack/class environment and
// returns the resulting frame stack, ready to be used as params[0].
func setupNullPointerExceptionTest(t *testing.T) *list.List {
	t.Helper()
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	f := frames.CreateFrame(2)
	f.Thread = 1
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"
	clData := classloader.ClData{
		Name:            "",
		SuperclassIndex: types.StringPoolObjectIndex,
		Module:          "test module",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Attributes:      nil,
		SourceFile:      "testClass.java",
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          0,
	}
	klass := classloader.Klass{Loader: "testLoader", Data: &clData}
	classloader.MethAreaInsert("java/testClass", &klass)
	f.ClName = "java/testClass"
	f.MethName = "java/testClass.test"

	jvmStack := frames.CreateFrameStack()
	_ = frames.PushFrame(jvmStack, f)

	globPtr := globals.GetGlobalRef()
	globPtr.FuncInstantiateClass = InstantiateFillIn

	return jvmStack
}

func newNullPointerExceptionObject() *object.Object {
	cls := "java/lang/NullPointerException"
	return object.MakeEmptyObjectWithClassName(&cls)
}

func TestLoadJavaLangNullPointerException_RegistersMethods(t *testing.T) {
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Java_Lang_NullPointerException()

	expected := []string{
		"java/lang/NullPointerException.<init>()V",
		"java/lang/NullPointerException.<init>(Ljava/lang/String;)V",
	}

	for _, sig := range expected {
		gmeth, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("expected method signature %q to be registered", sig)
			continue
		}
		if gmeth.GFunction == nil {
			t.Errorf("method signature %q has a nil GFunction", sig)
		}
		if !gmeth.NeedsContext {
			t.Errorf("method signature %q should need context", sig)
		}
	}
}

func TestInitNullPointerExceptionNull(t *testing.T) {
	jvmStack := setupNullPointerExceptionTest(t)
	obj := newNullPointerExceptionObject()

	ret := ThrowableInitNull([]any{jvmStack, obj})
	if ret != nil {
		t.Errorf("ThrowableInitNull returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok || dm.Fvalue != object.Null {
		t.Errorf("detailMessage not null; got: %#v", dm)
	}
}

func TestInitNullPointerExceptionString(t *testing.T) {
	jvmStack := setupNullPointerExceptionTest(t)
	obj := newNullPointerExceptionObject()
	msg := object.StringObjectFromGoString("null pointer")

	ret := throwableInitString([]any{jvmStack, obj, msg})
	if ret != nil {
		t.Errorf("throwableInitString returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "null pointer" {
		t.Errorf("expected detailMessage %q, got %q", "null pointer", got)
	}
}
