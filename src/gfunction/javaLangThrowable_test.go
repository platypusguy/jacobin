/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/statics"
	"strings"
	"testing"
)

func TestJavaLangThrowableClinit(t *testing.T) {
	statics.Statics = make(map[string]statics.Static)

	throwableClinit(nil)
	_, ok := statics.Statics["Throwable.UNASSIGNED_STACK"]
	if !ok {
		t.Error("JavaLangThrowableClinit: Throwable.UNASSIGNED_STACK not found")
	}

	_, ok = statics.Statics["Throwable.SUPPRESSED_SENTINEL"]
	if !ok {
		t.Error("JavaLangThrowableClinit: Throwable.SUPPRESSED_SENTINEL not found")
	}

	_, ok = statics.Statics["Throwable.EMPTY_THROWABLE_ARRAY"]
	if !ok {
		t.Error("Throwable.EMPTY_THROWABLE_ARRAY not found")
	}
}

func TestJavaLangThrowableFillInStackTraceWrongParmCount(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.SEVERE)

	params := []interface{}{1}
	err := fillInStackTrace(params)

	var retVal error
	switch err.(type) {
	case error:

		retVal = err.(error)
	default:

		t.Error("JavaLangThrowableFillInStack should have returned an error, but did not")
	}

	errMsg := retVal.Error()
	if !strings.HasPrefix(errMsg, "fillInStackTrace() expected two parameters") {
		t.Errorf("did not get expected error message, got: %s", errMsg)
	}
}

func TestJavaLangThrowableFillInStackTraceValid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.SEVERE)

	classloader.InitMethodArea()

	// we need to pass in a pointer to a valid JVM stack, so put one together here
	f := frames.CreateFrame(2) // create a new frame
	f.Thread = 1
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"

	clData := classloader.ClData{
		Name:        "",
		Superclass:  "",
		Module:      "test module",
		Pkg:         "",
		Interfaces:  nil,
		Fields:      nil,
		MethodTable: nil,
		Methods:     nil,
		Attributes:  nil,
		SourceFile:  "testClass.java",
		Bootstraps:  nil,
		CP:          classloader.CPool{},
		Access:      classloader.AccessFlags{},
		ClInit:      0,
	}
	klass := classloader.Klass{Loader: "testLoader", Data: &clData}
	classloader.MethAreaInsert("java/testClass", &klass)
	f.ClName = "java/testClass"
	f.MethName = "java/testClass.test"

	jvmStack := frames.CreateFrameStack()
	_ = frames.PushFrame(jvmStack, f)

	// create a Throwable object
	throw := object.MakeEmptyObject()
	klassType := "java/lang/Throwable"
	throw.Klass = &klassType

	// set up instantiate function
	globPtr := globals.GetGlobalRef()

	// Enable functions call InstantiateClass through a global function variable. (This avoids circularity issues.)
	globPtr.FuncInstantiateClass = InstantiateFillIn

	params := []interface{}{jvmStack, throw}
	retVal := fillInStackTrace(params)

	// var retVal error
	switch retVal.(type) {
	case error:
		t.Errorf("JavaLangThrowableFillInStack threw an unexpected error: %s", retVal.(error).Error())
	}

	// now, validate the fields in the stackTraceElementlement (ste)
	x := retVal.(*object.Field).Fvalue.(*object.Object)
	xt := x.Fields[0].Fvalue.(*[]*object.Object)
	xtt := *xt
	ste := xtt[0].FieldTable
	steDeclCl := ste["declaringClass"]
	if steDeclCl.Fvalue.(string) != "java/testClass" {
		t.Errorf("invalid STE entry for declaringClass: %s", steDeclCl)
	}

	steMethName := ste["methodName"]
	if steMethName.Fvalue.(string) != "java/testClass.test" {
		t.Errorf("invalid STE entry for methodName: %s", steMethName)
	}

	steFileName := ste["fileName"]
	if steFileName.Fvalue.(string) != "testClass.java" {
		t.Errorf("invalid STE entry for fileName: %s", steFileName)
	}

	steLoaderName := ste["classLoaderName"]
	if steLoaderName.Fvalue.(string) != "testLoader" {
		t.Errorf("invalid STE entry for classLoaderName: %s", steLoaderName)
	}
}

func InstantiateFillIn(name string, _ *list.List) (any, error) {
	o := object.MakeEmptyObject()
	o.Klass = &name
	return o, nil
}
