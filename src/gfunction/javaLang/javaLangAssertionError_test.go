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

// setupAssertionErrorTest builds a minimal frame/stack/class environment and
// returns the resulting frame stack, ready to be used as params[0].
func setupAssertionErrorTest(t *testing.T) *list.List {
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

func newAssertionErrorObject() *object.Object {
	cls := "java/lang/AssertionError"
	return object.MakeEmptyObjectWithClassName(&cls)
}

func TestLoadJavaLangAssertionError_RegistersMethods(t *testing.T) {
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Java_Lang_AssertionError()

	expected := []string{
		"java/lang/AssertionError.<init>()V",
		"java/lang/AssertionError.<init>(Z)V",
		"java/lang/AssertionError.<init>(Ljava/lang/Object;)V",
		"java/lang/AssertionError.<init>(C)V",
		"java/lang/AssertionError.<init>(I)V",
		"java/lang/AssertionError.<init>(J)V",
		"java/lang/AssertionError.<init>(F)V",
		"java/lang/AssertionError.<init>(D)V",
		"java/lang/AssertionError.<init>(Ljava/lang/String;)V",
		"java/lang/AssertionError.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V",
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

func TestInitAssertionErrorNull(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorNull([]any{jvmStack, obj})
	if ret != nil {
		t.Errorf("initAssertionErrorNull returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok || dm.Fvalue != object.Null {
		t.Errorf("detailMessage not null; got: %#v", dm)
	}
}

func TestInitAssertionErrorBoolean_True(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorBoolean([]any{jvmStack, obj, types.JavaBoolTrue})
	if ret != nil {
		t.Errorf("initAssertionErrorBoolean returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "true" {
		t.Errorf("expected detailMessage %q, got %q", "true", got)
	}
}

func TestInitAssertionErrorBoolean_False(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorBoolean([]any{jvmStack, obj, types.JavaBoolFalse})
	if ret != nil {
		t.Errorf("initAssertionErrorBoolean returned non-nil: %v", ret)
	}

	dm := obj.FieldTable["detailMessage"]
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "false" {
		t.Errorf("expected detailMessage %q, got %q", "false", got)
	}
}

func TestInitAssertionErrorBoolean_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorBoolean([]any{jvmStack, obj, "not a bool"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorInt64(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorInt64([]any{jvmStack, obj, int64(42)})
	if ret != nil {
		t.Errorf("initAssertionErrorInt64 returned non-nil: %v", ret)
	}

	dm := obj.FieldTable["detailMessage"]
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "42" {
		t.Errorf("expected detailMessage %q, got %q", "42", got)
	}
}

func TestInitAssertionErrorInt64_Negative(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorInt64([]any{jvmStack, obj, int64(-7)})
	if ret != nil {
		t.Errorf("initAssertionErrorInt64 returned non-nil: %v", ret)
	}

	dm := obj.FieldTable["detailMessage"]
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "-7" {
		t.Errorf("expected detailMessage %q, got %q", "-7", got)
	}
}

func TestInitAssertionErrorInt64_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorInt64([]any{jvmStack, obj, "not an int64"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorChar(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorChar([]any{jvmStack, obj, int64('A')})
	if ret != nil {
		t.Errorf("initAssertionErrorChar returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "A" {
		t.Errorf("expected detailMessage %q, got %q", "A", got)
	}
}

func TestInitAssertionErrorChar_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorChar([]any{jvmStack, obj, "not a char"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorFloat(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorFloat([]any{jvmStack, obj, float64(3.5)})
	if ret != nil {
		t.Errorf("initAssertionErrorFloat returned non-nil: %v", ret)
	}

	dm := obj.FieldTable["detailMessage"]
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "3.5" {
		t.Errorf("expected detailMessage %q, got %q", "3.5", got)
	}
}

func TestInitAssertionErrorFloat_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorFloat([]any{jvmStack, obj, "not a float"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorDouble(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorDouble([]any{jvmStack, obj, float64(2.25)})
	if ret != nil {
		t.Errorf("initAssertionErrorDouble returned non-nil: %v", ret)
	}

	dm := obj.FieldTable["detailMessage"]
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "2.25" {
		t.Errorf("expected detailMessage %q, got %q", "2.25", got)
	}
}

func TestInitAssertionErrorDouble_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorDouble([]any{jvmStack, obj, "not a double"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorObject(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()
	inner := object.StringObjectFromGoString("hello object")

	ret := initAssertionErrorObject([]any{jvmStack, obj, inner})
	if ret != nil {
		t.Errorf("initAssertionErrorObject returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "hello object" {
		t.Errorf("expected detailMessage to contain %q, got %q", "hello object", got)
	}
}

func TestInitAssertionErrorObject_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorObject([]any{jvmStack, obj, "not an object"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorString(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()
	msg := object.StringObjectFromGoString("assertion failed")

	ret := initAssertionErrorString([]any{jvmStack, obj, msg})
	if ret != nil {
		t.Errorf("initAssertionErrorString returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "assertion failed" {
		t.Errorf("expected detailMessage %q, got %q", "assertion failed", got)
	}
}

func TestInitAssertionErrorString_WrongType(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()

	ret := initAssertionErrorString([]any{jvmStack, obj, "not a string object"})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestInitAssertionErrorString_NotAStringObject(t *testing.T) {
	jvmStack := setupAssertionErrorTest(t)
	obj := newAssertionErrorObject()
	// a non-String object, with a different KlassName
	notAString := newAssertionErrorObject()

	ret := initAssertionErrorString([]any{jvmStack, obj, notAString})
	errBlk, ok := ret.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected *ghelpers.GErrBlk, got %T: %v", ret, ret)
	}
	if errBlk.ErrMsg == "" {
		t.Errorf("expected non-empty error message")
	}
}
