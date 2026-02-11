/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"testing"
)

func TestModifierIsPublicTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsPublic(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isPublic(PUBLIC) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsPublicFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PRIVATE)}
	result := modifierIsPublic(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isPublic(PRIVATE) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsPublicWithCombinedFlags(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC | STATIC | FINAL)}
	result := modifierIsPublic(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isPublic(PUBLIC|STATIC|FINAL) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsPrivateTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PRIVATE)}
	result := modifierIsPrivate(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isPrivate(PRIVATE) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsPrivateFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsPrivate(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isPrivate(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsProtectedTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PROTECTED)}
	result := modifierIsProtected(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isProtected(PROTECTED) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsProtectedFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsProtected(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isProtected(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsStaticTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(STATIC)}
	result := modifierIsStatic(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isStatic(STATIC) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsStaticFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsStatic(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isStatic(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsFinalTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(FINAL)}
	result := modifierIsFinal(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isFinal(FINAL) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsFinalFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsFinal(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isFinal(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsSynchronizedTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(SYNCHRONIZED)}
	result := modifierIsSynchronized(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isSynchronized(SYNCHRONIZED) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsSynchronizedFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsSynchronized(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isSynchronized(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsVolatileTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(VOLATILE)}
	result := modifierIsVolatile(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isVolatile(VOLATILE) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsVolatileFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsVolatile(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isVolatile(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsTransientTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(TRANSIENT)}
	result := modifierIsTransient(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isTransient(TRANSIENT) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsTransientFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsTransient(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isTransient(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsNativeTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(NATIVE)}
	result := modifierIsNative(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isNative(NATIVE) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsNativeFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsNative(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isNative(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsInterfaceTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(INTERFACE)}
	result := modifierIsInterface(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isInterface(INTERFACE) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsInterfaceFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsInterface(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isInterface(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsAbstractTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(ABSTRACT)}
	result := modifierIsAbstract(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isAbstract(ABSTRACT) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsAbstractFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsAbstract(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isAbstract(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierIsStrictTrue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(STRICT)}
	result := modifierIsStrict(params)
	if result.(int64) != 1 {
		t.Errorf("Expected isStrict(STRICT) to return 1, got %d", result.(int64))
	}
}

func TestModifierIsStrictFalse(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierIsStrict(params)
	if result.(int64) != 0 {
		t.Errorf("Expected isStrict(PUBLIC) to return 0, got %d", result.(int64))
	}
}

func TestModifierToStringEmpty(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(0)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "" {
		t.Errorf("Expected toString(0) to return empty string, got '%s'", resultStr)
	}
}

func TestModifierToStringPublic(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "public" {
		t.Errorf("Expected toString(PUBLIC) to return 'public', got '%s'", resultStr)
	}
}

func TestModifierToStringPrivate(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PRIVATE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "private" {
		t.Errorf("Expected toString(PRIVATE) to return 'private', got '%s'", resultStr)
	}
}

func TestModifierToStringProtected(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PROTECTED)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "protected" {
		t.Errorf("Expected toString(PROTECTED) to return 'protected', got '%s'", resultStr)
	}
}

func TestModifierToStringAbstract(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(ABSTRACT)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "abstract" {
		t.Errorf("Expected toString(ABSTRACT) to return 'abstract', got '%s'", resultStr)
	}
}

func TestModifierToStringStatic(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(STATIC)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "static" {
		t.Errorf("Expected toString(STATIC) to return 'static', got '%s'", resultStr)
	}
}

func TestModifierToStringFinal(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(FINAL)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "final" {
		t.Errorf("Expected toString(FINAL) to return 'final', got '%s'", resultStr)
	}
}

func TestModifierToStringTransient(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(TRANSIENT)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "transient" {
		t.Errorf("Expected toString(TRANSIENT) to return 'transient', got '%s'", resultStr)
	}
}

func TestModifierToStringVolatile(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(VOLATILE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "volatile" {
		t.Errorf("Expected toString(VOLATILE) to return 'volatile', got '%s'", resultStr)
	}
}

func TestModifierToStringSynchronized(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(SYNCHRONIZED)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "synchronized" {
		t.Errorf("Expected toString(SYNCHRONIZED) to return 'synchronized', got '%s'", resultStr)
	}
}

func TestModifierToStringNative(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(NATIVE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "native" {
		t.Errorf("Expected toString(NATIVE) to return 'native', got '%s'", resultStr)
	}
}

func TestModifierToStringStrict(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(STRICT)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "strictfp" {
		t.Errorf("Expected toString(STRICT) to return 'strictfp', got '%s'", resultStr)
	}
}

func TestModifierToStringInterface(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(INTERFACE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "interface" {
		t.Errorf("Expected toString(INTERFACE) to return 'interface', got '%s'", resultStr)
	}
}

func TestModifierToStringPublicStaticFinal(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC | STATIC | FINAL)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "public static final" {
		t.Errorf("Expected toString(PUBLIC|STATIC|FINAL) to return 'public static final', got '%s'", resultStr)
	}
}

func TestModifierToStringPublicAbstractInterface(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC | ABSTRACT | INTERFACE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "public abstract interface" {
		t.Errorf("Expected toString(PUBLIC|ABSTRACT|INTERFACE) to return 'public abstract interface', got '%s'", resultStr)
	}
}

func TestModifierToStringPrivateStaticFinalSynchronized(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PRIVATE | STATIC | FINAL | SYNCHRONIZED)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "private static final synchronized" {
		t.Errorf("Expected 'private static final synchronized', got '%s'", resultStr)
	}
}

func TestModifierToStringProtectedTransientVolatile(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PROTECTED | TRANSIENT | VOLATILE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "protected transient volatile" {
		t.Errorf("Expected 'protected transient volatile', got '%s'", resultStr)
	}
}

func TestModifierToStringAllModifiers(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC | PROTECTED | PRIVATE | ABSTRACT | STATIC | FINAL | TRANSIENT | VOLATILE | SYNCHRONIZED | NATIVE | STRICT | INTERFACE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	expected := "public protected private abstract static final transient volatile synchronized native strictfp interface"
	if resultStr != expected {
		t.Errorf("Expected '%s', got '%s'", expected, resultStr)
	}
}

func TestModifierToStringPublicFinalSynchronizedNative(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []interface{}{int64(PUBLIC | FINAL | SYNCHRONIZED | NATIVE)}
	result := modifierToString(params)
	strObj := result.(*object.Object)
	resultStr := object.GoStringFromStringObject(strObj)
	if resultStr != "public final synchronized native" {
		t.Errorf("Expected 'public final synchronized native', got '%s'", resultStr)
	}
}
