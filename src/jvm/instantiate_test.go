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
	"jacobin/log"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
	"os"
	"testing"
)

// Arrays are preloaded, so this should only confirm the presence of the class
// in the method area--and make sure it has no fields.
func TestInstantiateArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)
	classloader.InitMethodArea()

	anything, err := InstantiateClass(types.ByteArray, nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating array: %s", err.Error())
	}
	obj := anything.(*object.Object)
	if len(obj.Fields) != 0 {
		t.Errorf("Expected 0 fields in array class, got %d fields", len(obj.Fields))
	}
}

func TestInstantiateString1(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)
	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err := classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	gfunction.MTableLoadNatives(&classloader.MTable)
	statics.StaticsPreload()

	myobj, err := InstantiateClass("java/lang/String", nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating string: %s", err.Error())
	}

	obj := myobj.(*object.Object)
	if *obj.Klass != "java/lang/String" {
		t.Errorf("Expected 'java/lang/String', got %s", *obj.Klass)
	}

	if len(obj.Fields) < 5 {
		t.Errorf("Expected more than 4 fielsd in String object, got %d fields", len(obj.Fields))
	}
}

func TestInstantiateNonExistentClass(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

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
	gfunction.MTableLoadNatives(&classloader.MTable)
	statics.StaticsPreload()

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
