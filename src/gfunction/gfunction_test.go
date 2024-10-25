/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/trace"
	"testing"
)

func f1([]interface{}) interface{} { return nil }
func f2([]interface{}) interface{} { return nil }
func f3([]interface{}) interface{} { return nil }

func TestMTableLoadLib(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	libMeths := make(map[string]GMeth)
	libMeths["test.f1()V"] = GMeth{ParamSlots: 0, GFunction: f1}
	libMeths["test.f2(I)V"] = GMeth{ParamSlots: 1, GFunction: f2}
	libMeths["test.f3(Ljava/lang/String;JZ)D"] = GMeth{ParamSlots: 3, GFunction: f3}
	mtbl := make(classloader.MT)
	loadlib(&mtbl, libMeths)
	if len(mtbl) != 3 {
		t.Errorf("ERROR, Expecting MTable with 3 entries, got: %d\n", len(mtbl))
	}
	mte := libMeths["test.f1()V"]
	if mte.ParamSlots != 0 {
		t.Errorf("ERROR, Expecting f1 MTable entry to have 0 param slots, got: %d\n",
			mte.ParamSlots)
	}
	mte = libMeths["test.f2(I)V"]
	if mte.ParamSlots != 1 {
		t.Errorf("ERROR, Expecting f2 MTable entry to have 1 param slots, got: %d\n",
			mte.ParamSlots)
	}
	mte = libMeths["test.f3(Ljava/lang/String;JZ)D"]
	if mte.ParamSlots != 3 {
		t.Errorf("ERROR, Expecting f3 MTable entry to have 3 param slots, got: %d\n",
			mte.ParamSlots)
	}

	if mte.NeedsContext {
		t.Errorf("ERROR, Expecting MTable entry's NeedContext to be false\n")
	}
}

// test loading of native functions

func TestMTableLoadGFunctions(t *testing.T) {
	classloader.MTable = make(map[string]classloader.MTentry)
	MTableLoadGFunctions(&classloader.MTable)
	mte, exists := classloader.MTable["java/lang/Object.<init>()V"]
	if !exists {
		t.Errorf("Expecting MTable entry for java/lang/Object.<init>()V, but it does not exist")
	}

	if mte.MType != 'G' {
		t.Errorf("Expecting java/lang/Object.<init>()V to be of type 'G', but got type: %c",
			mte.MType)
	}
}

// make sure that JustReturn in fact does nothing
func TestJustReturn(t *testing.T) {
	retVal := justReturn(nil)
	if retVal != nil {
		t.Errorf("Expecting nil return value, got: %v", retVal)
	}
}

func TestCheckKey(t *testing.T) {
	if checkKey("java/lang/Object") != false {
		t.Errorf("invalid key %s was allowed in gfunction", "java/lang/Object")
	}

	if checkKey("java/lang/Object.toString") != false {
		t.Errorf("invalid key %s was allowed in gfunction",
			"java/lang/Object.toString")
	}

	if checkKey("java/lang/Object.toString(") != false {
		t.Errorf("invalid key %s was allowed in gfunction",
			"java/lang/Object.toString(")
	}

	if checkKey("java/lang/Object.toString()") != false {
		t.Errorf("invalid key %s was allowed in gfunction",
			"java/lang/Object.toString()")
	}

	if checkKey("java/lang/Object.toString(I)Ljava/lang/String;") != true {
		t.Errorf("got unexpected error checking key %s",
			"java/lang/Object.toString(I)Ljava/lang/String;")
	}
}
