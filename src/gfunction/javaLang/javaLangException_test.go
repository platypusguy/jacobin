/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"testing"
)

func newExceptionObject() *object.Object {
	cls := "java/lang/Exception"
	return object.MakeEmptyObjectWithClassName(&cls)
}

func TestLoadJavaLangException_RegistersMethods(t *testing.T) {
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Java_Lang_Exception()

	expected := []string{
		"java/lang/Exception.<init>()V",
		"java/lang/Exception.<init>(Ljava/lang/String;)V",
		"java/lang/Exception.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V",
		"java/lang/Exception.<init>(Ljava/lang/Throwable;)V",
		"java/lang/Exception.<init>(Ljava/lang/String;Ljava/lang/Throwable;ZZ)V",
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
	}

	// the two implemented constructors should need context
	for _, sig := range []string{
		"java/lang/Exception.<init>()V",
		"java/lang/Exception.<init>(Ljava/lang/String;)V",
	} {
		if !ghelpers.MethodSignatures[sig].NeedsContext {
			t.Errorf("method signature %q should need context", sig)
		}
	}
}

func TestInitExceptionNull(t *testing.T) {
	jvmStack := setupNullPointerExceptionTest(t)
	obj := newExceptionObject()

	ret := ThrowableInitNull([]any{jvmStack, obj})
	if ret != nil {
		t.Errorf("ThrowableInitNull returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok || dm.Fvalue != object.Null {
		t.Errorf("detailMessage not null; got: %#v", dm)
	}
}

func TestInitExceptionString(t *testing.T) {
	jvmStack := setupNullPointerExceptionTest(t)
	obj := newExceptionObject()
	msg := object.StringObjectFromGoString("something failed")

	ret := throwableInitString([]any{jvmStack, obj, msg})
	if ret != nil {
		t.Errorf("throwableInitString returned non-nil: %v", ret)
	}

	dm, ok := obj.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	got := object.GoStringFromStringObject(dm.Fvalue.(*object.Object))
	if got != "something failed" {
		t.Errorf("expected detailMessage %q, got %q", "something failed", got)
	}
}
