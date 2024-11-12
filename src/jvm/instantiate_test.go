/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"jacobin/classloader"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"testing"
)

// Arrays are preloaded, so this should only confirm the presence of the class
// in the method area--and make sure it has no fields.
func TestInstantiateArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	anything, err := InstantiateClass(types.ByteArray, nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating array: %s", err.Error())
	}
	obj := anything.(*object.Object)
	if len(obj.FieldTable) != 0 {
		t.Errorf("Expected 0 fields in array class, got %d fields", len(obj.FieldTable))
	}
}

func TestInstantiateString1(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err := classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	myobj, err := InstantiateClass(types.StringClassName, nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating string: %s", err.Error())
	}

	obj := myobj.(*object.Object)
	klassType := stringPool.GetStringPointer(obj.KlassName)
	if *klassType != types.StringClassName {
		t.Errorf("Expected 'java/lang/String', got %s", *klassType)
	}

	if len(obj.FieldTable) < 2 {
		t.Errorf("Expected more than 1 field in String object, got %d fields", len(obj.FieldTable))
	}
}

func TestInstantiateNonExistentClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, err := os.Pipe()

	os.Stderr = werr

	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	statics.PreloadStatics()

	myobj, err := InstantiateClass("$nosuchclass", nil)

	// restore stderr
	_ = werr.Close()
	os.Stderr = normalStderr

	if err == nil {
		t.Errorf("Expected error message for nonexistent class, but got none")
	}

	if myobj != nil {
		t.Errorf("Expected nil object, got %v", myobj)
	}
}

func TestLoadValidClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, err := os.Pipe()
	if err != nil {
		t.Error("cannot create pipe for stderr")
	}
	os.Stderr = werr

	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	// we'll check that the class is loaded, then delete it, then load it and check again

	class := classloader.MethAreaFetch("java/lang/Integer")
	if class == nil {
		t.Errorf("Expected java.lang.Integer to be loaded in method area, but it wasn't")
	}

	classloader.MethAreaDelete("java/lang/Integer")
	class = classloader.MethAreaFetch("java/lang/Integer")
	if class != nil {
		t.Errorf("Expected java.lang.Integer to be absent from method area, but it wasn't")
	}

	// now load the class
	err = loadThisClass("java/lang/Integer")
	if err != nil {
		t.Errorf("Got unexpected error from loadThisClass(\"java/lang/Integer\"): %s", err.Error())
	}
	class = classloader.MethAreaFetch("java/lang/Integer")
	if class == nil {
		t.Errorf("Expected java.lang.Integer to be loaded in method area, but it wasn't")
	}

	// restore stderr
	_ = werr.Close()
	os.Stderr = normalStderr
}

// This should always work. java/lang/Object contains no instance or static fields,
// so this is about as simple a class instantiation as possible
func TestLoadClassJavaLangObject(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, err := os.Pipe()
	os.Stderr = werr

	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	err = loadThisClass(types.ObjectClassName)

	// this should always work. java/lang/Object contains no instance or static fields,
	// so this is about as simple a class instantiation as possible

	// restore stderr
	_ = werr.Close()
	os.Stderr = normalStderr

	if err != nil {
		t.Errorf("Got unexpected error from loadThisClass: %s", err.Error())
	}
}
