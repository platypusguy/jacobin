/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"testing"
)

func TestClassFromInstance_NilObject(t *testing.T) {
	globals.InitGlobals("test")

	cl := ClassFromInstance(nil)
	if cl != nil {
		t.Errorf("expected nil result for nil object, got %v", cl)
	}
}

func TestClassFromInstance_ClassNotInMethodArea(t *testing.T) {
	globals.InitGlobals("test")
	InitMethodArea()

	className := "com/test/NotLoadedClass"
	obj := object.MakeEmptyObjectWithClassName(&className)

	cl := ClassFromInstance(obj)
	if cl != nil {
		t.Errorf("expected nil result when class is not in method area, got %v", cl)
	}
}

func TestClassFromInstance_Success(t *testing.T) {
	globals.InitGlobals("test")
	InitMethodArea()

	className := "com/test/SomeClass"
	superclassName := "java/lang/Object"
	superIndex := stringPool.GetStringIndex(&superclassName)

	obj := object.MakeEmptyObjectWithClassName(&className)

	fld := Field{
		AccessFlags: 0,
		NameStr:     "counter",
		DescStr:     "I",
		IsStatic:    false,
	}

	meth := &Method{}

	kd := &ClData{
		Name:            className,
		SuperclassIndex: superIndex,
		Fields:          []Field{fld},
		MethodTable:     map[string]*Method{"doSomething": meth},
	}
	k := &Klass{
		Status: 'L',
		Loader: "bootstrap",
		Data:   kd,
	}
	MethAreaInsert(className, k)

	// also register the superclass so SuperClass lookup succeeds cleanly
	superKD := &ClData{Name: superclassName}
	superK := &Klass{Status: 'L', Loader: "bootstrap", Data: superKD}
	MethAreaInsert(superclassName, superK)

	cl := ClassFromInstance(obj)
	if cl == nil {
		t.Fatalf("expected non-nil result")
	}
	if cl.Name == nil || *cl.Name != className {
		t.Errorf("expected Name to be %q, got %v", className, cl.Name)
	}
	if cl.SuperClass == nil || cl.SuperClass.Data.Name != superclassName {
		t.Errorf("expected SuperClass to point to %q, got %v", superclassName, cl.SuperClass)
	}
	if cl.IsArray {
		t.Errorf("expected IsArray to be false, since no 'value' field with array type present")
	}

	fldEntry, ok := cl.Fields["counter"]
	if !ok {
		t.Fatalf("expected Fields to contain 'counter'")
	}
	if fldEntry.Type != "I" {
		t.Errorf("expected field type 'I', got %q", fldEntry.Type)
	}
	if fldEntry.IsStatic {
		t.Errorf("expected IsStatic to be false")
	}

	methEntry, ok := cl.Methods["doSomething"]
	if !ok {
		t.Fatalf("expected Methods to contain 'doSomething'")
	}
	if methEntry.Name != "doSomething" {
		t.Errorf("expected method Name to be 'doSomething', got %q", methEntry.Name)
	}
}

func TestClassFromInstance_ArrayValueField(t *testing.T) {
	globals.InitGlobals("test")
	InitMethodArea()

	className := "com/test/ArrayHolderClass"
	superclassName := "java/lang/Object"
	superIndex := stringPool.GetStringIndex(&superclassName)

	obj := object.MakeEmptyObjectWithClassName(&className)
	// simulate an array-typed "value" field, as is used for array wrapper objects
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.IntArray,
		Fvalue: []int64{1, 2, 3},
	}

	kd := &ClData{
		Name:            className,
		SuperclassIndex: superIndex,
		Fields:          []Field{},
		MethodTable:     map[string]*Method{},
	}
	k := &Klass{Status: 'L', Loader: "bootstrap", Data: kd}
	MethAreaInsert(className, k)

	superKD := &ClData{Name: superclassName}
	superK := &Klass{Status: 'L', Loader: "bootstrap", Data: superKD}
	MethAreaInsert(superclassName, superK)

	cl := ClassFromInstance(obj)
	if cl == nil {
		t.Fatalf("expected non-nil result")
	}
	if !cl.IsArray {
		t.Errorf("expected IsArray to be true when 'value' field has an array type")
	}
}
