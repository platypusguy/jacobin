/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"jacobin/src/classloader"
	"jacobin/src/exceptions"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// === NOTE ===
// This file contains tests for the java.lang.Object class. Due to circular dependencies,
// the java/lang/Class, while being a gfunction, must be tested inside the jvm package.
// ============

var setUpRun = false

func setup() {
	if !setUpRun {
		globals.InitGlobals("test")
		globals.GetGlobalRef().FuncThrowException = exceptions.ThrowExNil
		classloader.InitMethodArea()
		_ = classloader.Init()
		classloader.LoadBaseClasses()
		setUpRun = true
	}
}

func TestGetNameWithAStringObject(t *testing.T) {
	setup()
	InitGlobalFunctionPointers()
	rawCl := classloader.MethAreaFetch("java/lang/String")
	rawCl.Data.ClInit = types.ClInitRun

	// instantiate a String object and return a *object.Object
	str, ok := globals.GetGlobalRef().FuncInstantiateClass("java/lang/String", nil)
	if ok != nil {
		t.Fatalf("Failed to instantiate java/lang/String: %v", str)
	}

	// returns the skeletal mirrored java/lang/Class instance of the passed-in object
	parmsA := []interface{}{str}
	cl := javaLang.ObjectGetClass(parmsA)

	parmsB := []interface{}{cl}
	nameObj := javaLang.ClassGetName(parmsB).(*object.Object)

	observed := object.GoStringFromStringObject(nameObj)
	if observed == "" {
		t.Error("Expected java/lang/String, got \"\"")
	}
	if observed != "java/lang/String" {
		t.Errorf("Expected java/lang/String, got %s", observed)
	}
}
