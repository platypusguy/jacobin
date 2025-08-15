/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/trace"
	"jacobin/types"
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


func TestPopulator_PrimitivesAndString(t *testing.T) {
    globals.InitGlobals("test")
    globals.InitStringPool()

    // Integer primitive object
    iobj := Populator("java/lang/Integer", "I", int64(42))
    if iobj == nil {
        t.Fatalf("Populator returned nil for Integer")
    }
    fld, ok := iobj.FieldTable["value"]
    if !ok || fld.Ftype != "I" {
        t.Fatalf("Integer value field missing or wrong type: %#v", iobj.FieldTable["value"])
    }
    if v, ok := fld.Fvalue.(int64); !ok || v != 42 {
        t.Fatalf("Integer value mismatch: %v", fld.Fvalue)
    }

    // String via StringIndex path returns a proper String object
    sobj := Populator("java/lang/String", "T", "hello")
    if sobj == nil {
        t.Fatalf("Populator returned nil for String")
    }
    if !object.IsStringObject(sobj) {
        t.Fatalf("Populator did not create a String object")
    }
    if s := object.GoStringFromStringObject(sobj); s != "hello" {
        t.Fatalf("String content mismatch: %q", s)
    }
}

func TestReturnNullTrueFalse(t *testing.T) {
    if v := returnNull(nil); v != object.Null {
        t.Fatalf("returnNull did not return object.Null: %v", v)
    }
    if v := returnTrue(nil); v != types.JavaBoolTrue {
        t.Fatalf("returnTrue != true: %v", v)
    }
    if v := returnFalse(nil); v != types.JavaBoolFalse {
        t.Fatalf("returnFalse != false: %v", v)
    }
}

func TestEOFSetGet(t *testing.T) {
    obj := object.MakeEmptyObject()
    eofSet(obj, true)
    if !eofGet(obj) {
        t.Fatalf("eofGet expected true")
    }
    eofSet(obj, false)
    if eofGet(obj) {
        t.Fatalf("eofGet expected false")
    }
}

func TestReturnRandomLong_Type(t *testing.T) {
    v := returnRandomLong(nil)
    if _, ok := v.(int64); !ok {
        t.Fatalf("returnRandomLong did not return int64, got %T", v)
    }
}
