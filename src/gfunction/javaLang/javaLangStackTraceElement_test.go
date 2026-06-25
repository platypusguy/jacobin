/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestStackTraceElement_Getters(t *testing.T) {
	obj := object.MakeEmptyObject()
	obj.FieldTable["classLoaderName"] = object.Field{Ftype: types.GolangString, Fvalue: "app"}
	obj.FieldTable["declaringClass"] = object.Field{Ftype: types.GolangString, Fvalue: "com.example.Test"}
	obj.FieldTable["fileName"] = object.Field{Ftype: types.GolangString, Fvalue: "Test.java"}
	obj.FieldTable["methodName"] = object.Field{Ftype: types.GolangString, Fvalue: "main"}
	obj.FieldTable["moduleName"] = object.Field{Ftype: types.GolangString, Fvalue: "test.module"}
	obj.FieldTable["sourceLine"] = object.Field{Ftype: types.GolangString, Fvalue: "42"}

	params := []interface{}{obj}

	// Test getClassLoaderName
	res := steGetClassLoaderName(params)
	if object.GoStringFromStringObject(res.(*object.Object)) != "app" {
		t.Errorf("Expected app, got %v", object.GoStringFromStringObject(res.(*object.Object)))
	}

	// Test getClassName
	res = steGetClassName(params)
	if object.GoStringFromStringObject(res.(*object.Object)) != "com.example.Test" {
		t.Errorf("Expected com.example.Test, got %v", object.GoStringFromStringObject(res.(*object.Object)))
	}

	// Test getFileName
	res = steGetFileName(params)
	if object.GoStringFromStringObject(res.(*object.Object)) != "Test.java" {
		t.Errorf("Expected Test.java, got %v", object.GoStringFromStringObject(res.(*object.Object)))
	}

	// Test getMethodName
	res = steGetMethodName(params)
	if object.GoStringFromStringObject(res.(*object.Object)) != "main" {
		t.Errorf("Expected main, got %v", object.GoStringFromStringObject(res.(*object.Object)))
	}

	// Test getModuleName
	res = steGetModuleName(params)
	if object.GoStringFromStringObject(res.(*object.Object)) != "test.module" {
		t.Errorf("Expected test.module, got %v", object.GoStringFromStringObject(res.(*object.Object)))
	}

	// Test getLineNumber
	resNum := steGetLineNumber(params)
	if resNum.(int64) != 42 {
		t.Errorf("Expected 42, got %v", resNum)
	}
}

func TestStackTraceElement_Getters_MissingFields(t *testing.T) {
	obj := object.MakeEmptyObject()
	params := []interface{}{obj}

	// Test getClassLoaderName missing
	res := steGetClassLoaderName(params)
	if object.GoStringFromStringObject(res.(*object.Object)) != "<missing>" {
		t.Errorf("Expected <missing>, got %v", object.GoStringFromStringObject(res.(*object.Object)))
	}

	// Test getLineNumber missing
	resNum := steGetLineNumber(params)
	if resNum.(int64) != -1 {
		t.Errorf("Expected -1, got %v", resNum)
	}
}

func TestStackTraceElement_Init(t *testing.T) {
	// Setup StringPool and GlobalRef
	globals.InitGlobals("test")
	globals.InitStringPool()

	// Setup Method Area
	classloader.InitMethodArea()
	klass := &classloader.Klass{
		Loader: "app",
		Data: &classloader.ClData{
			SourceFile: "Test.java",
			Module:     "test.module",
		},
	}
	classloader.MethAreaInsert("com/example/Test", klass)

	// Setup Frame
	frame := &frames.Frame{
		ClName:   "com/example/Test",
		MethName: "main",
		MethType: "([Ljava/lang/String;)V",
		PC:       10,
	}

	// Setup StackTraceElement object
	ste := object.MakeEmptyObject()

	// Call initStackTraceElement
	initStackTraceElement(ste, frame, true)

	// Verify fields
	if ste.FieldTable["declaringClass"].Fvalue.(string) != "com/example/Test" {
		t.Errorf("Expected com/example/Test, got %v", ste.FieldTable["declaringClass"].Fvalue)
	}
	if ste.FieldTable["methodName"].Fvalue.(string) != "main" {
		t.Errorf("Expected main, got %v", ste.FieldTable["methodName"].Fvalue)
	}
	if ste.FieldTable["classLoaderName"].Fvalue.(string) != "app" {
		t.Errorf("Expected app, got %v", ste.FieldTable["classLoaderName"].Fvalue)
	}
	if ste.FieldTable["fileName"].Fvalue.(string) != "Test.java" {
		t.Errorf("Expected Test.java, got %v", ste.FieldTable["fileName"].Fvalue)
	}
	if ste.FieldTable["moduleName"].Fvalue.(string) != "test.module" {
		t.Errorf("Expected test.module, got %v", ste.FieldTable["moduleName"].Fvalue)
	}
}

func TestStackTraceElement_Of(t *testing.T) {
	// Setup StringPool and GlobalRef
	globals.InitGlobals("test")
	globals.InitStringPool()

	// Setup Globals and mock FuncInstantiateClass
	g := globals.GetGlobalRef()
	g.FuncInstantiateClass = func(classname string, frameStack *list.List) (any, error) {
		return object.MakeEmptyObject(), nil
	}

	// Setup Throwable with frameStackRef
	jvmStack := list.New()
	jvmStack.PushBack(&frames.Frame{
		ClName:   "com/example/Test",
		MethName: "main",
		PC:       5,
	})

	throwable := object.MakeEmptyObject()
	throwable.FieldTable["frameStackRef"] = object.Field{Fvalue: jvmStack}

	// Setup Method Area for initStackTraceElement
	classloader.InitMethodArea()
	klass := &classloader.Klass{
		Loader: "platform",
		Data: &classloader.ClData{
			SourceFile: "Test.java",
		},
	}
	classloader.MethAreaInsert("com/example/Test", klass)

	// Call of
	params := []interface{}{throwable, int64(1)}
	res := of(params)

	// Verify result
	arrayObj := res.(*object.Object)
	rawArray := arrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(rawArray) != 1 {
		t.Fatalf("Expected array length 1, got %d", len(rawArray))
	}

	ste := rawArray[0]
	if ste.FieldTable["declaringClass"].Fvalue.(string) != "com/example/Test" {
		t.Errorf("Expected com/example/Test, got %v", ste.FieldTable["declaringClass"].Fvalue)
	}
}

func TestStackTraceElement_SearchLineNumberTable(t *testing.T) {
	// LineNumberTable format:
	// u2 line_number_table_length;
	// { u2 start_pc; u2 line_number; } line_number_table[line_number_table_length];

	// length = 2
	// entry 1: start_pc=0, line_number=10
	// entry 2: start_pc=5, line_number=20
	attrContent := []byte{
		0, 2, // length
		0, 0, 0, 10, // entry 1
		0, 5, 0, 20, // entry 2
	}

	// Test PC=0 -> line 10
	line := searchLineNumberTable(attrContent, 0)
	if line != 10 {
		t.Errorf("Expected 10, got %d", line)
	}

	// Test PC=3 -> line 10
	line = searchLineNumberTable(attrContent, 3)
	if line != 10 {
		t.Errorf("Expected 10, got %d", line)
	}

	// Test PC=5 -> line 20
	line = searchLineNumberTable(attrContent, 5)
	if line != 20 {
		t.Errorf("Expected 20, got %d", line)
	}

	// Test PC=10 -> line 20
	line = searchLineNumberTable(attrContent, 10)
	if line != 20 {
		t.Errorf("Expected 20, got %d", line)
	}
}
