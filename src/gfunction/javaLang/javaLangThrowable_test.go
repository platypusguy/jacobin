/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"strings"
	"testing"
)

func TestJavaLangThrowableClinit(t *testing.T) {
	statics.Statics = make(map[string]statics.Static)
	globals.InitStringPool()

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
	trace.Init()

	params := []interface{}{1}
	err := FillInStackTrace(params)

	var retVal error
	switch err.(type) {
	case error:

		retVal = err.(error)
	default:

		t.Error("JavaLangThrowableFillInStack should have returned an error, but did not")
	}

	errMsg := retVal.Error()
	expPrefix := "FillInStackTrace: expected two parameters"
	if !strings.HasPrefix(errMsg, expPrefix) {
		t.Errorf("did not get expected error message: %s, observed: %s", expPrefix, errMsg)
	}
}

func TestJavaLangThrowableFillInStackTraceValid(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	classloader.InitMethodArea()

	// we need to pass in a pointer to a valid JVM stack, so put one together here
	f := frames.CreateFrame(2) // create a new frame
	f.Thread = 1
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"

	clData := classloader.ClData{
		Name:            "",
		SuperclassIndex: types.StringPoolObjectIndex,
		Module:          "test module",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Attributes:      nil,
		SourceFile:      "testClass.java",
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          0,
	}
	klass := classloader.Klass{Loader: "testLoader", Data: &clData}
	classloader.MethAreaInsert("java/testClass", &klass)
	f.ClName = "java/testClass"
	f.MethName = "java/testClass.test"

	jvmStack := frames.CreateFrameStack()
	_ = frames.PushFrame(jvmStack, f)

	// create a Throwable object
	str := "java/lang/Throwable"
	throw := object.MakeEmptyObjectWithClassName(&str)

	// set up instantiate function
	globPtr := globals.GetGlobalRef()

	// Enable functions call InstantiateClass through a global function variable. (This avoids circularity issues.)
	globPtr.FuncInstantiateClass = InstantiateFillIn

	params := []interface{}{jvmStack, throw}
	retVal := FillInStackTrace(params)

	// var retVal error
	switch retVal.(type) {
	case error:
		t.Errorf("JavaLangThrowableFillInStack threw an unexpected error: %s", retVal.(error).Error())
	}

	// now, validate the fields in the stackTraceElementlement (ste)
	x := retVal.(*object.Field).Fvalue.(*object.Object)
	xtt := x.FieldTable["value"].Fvalue.([]*object.Object)
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

/*
	func TestMinimalThrowEx(t *testing.T) {
		globals.InitGlobals("test")
		trace.Init()

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
		exceptions.ThrowEx(excNames.NullPointerException, "test of NPE", f)
	}
*/
func InstantiateFillIn(name string, _ *list.List) (any, error) {
	o := object.MakeEmptyObject()
	o.KlassName = stringPool.GetStringIndex(&name)
	return o, nil
}

// Additional tests for Throwable constructors
func TestThrowableInitNull_SetsFieldsAndStack(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	// Build a minimal frame and stack
	f := frames.CreateFrame(2)
	f.Thread = 1
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"
	clData := classloader.ClData{
		Name:            "",
		SuperclassIndex: types.StringPoolObjectIndex,
		Module:          "test module",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Attributes:      nil,
		SourceFile:      "testClass.java",
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          0,
	}
	klass := classloader.Klass{Loader: "testLoader", Data: &clData}
	classloader.MethAreaInsert("java/testClass", &klass)
	f.ClName = "java/testClass"
	f.MethName = "java/testClass.test"
	jvmStack := frames.CreateFrameStack()
	_ = frames.PushFrame(jvmStack, f)

	// Throwable instance
	cls := "java/lang/Throwable"
	th := object.MakeEmptyObjectWithClassName(&cls)

	// set up instantiate function required by StackTraceElement.of()
	globPtr := globals.GetGlobalRef()
	globPtr.FuncInstantiateClass = InstantiateFillIn

	// Call constructor
	ret := throwableInitNull([]interface{}{jvmStack, th})
	if ret != nil {
		t.Errorf("throwableInitNull returned non-nil: %v", ret)
	}

	// Validate fields
	dm, ok := th.FieldTable["detailMessage"]
	if !ok || dm.Fvalue != object.Null {
		t.Errorf("detailMessage not null; got: %#v", dm)
	}
	cause, ok := th.FieldTable["cause"]
	if !ok || cause.Fvalue != object.Null {
		t.Errorf("cause not null; got: %#v", cause)
	}

	// Validate stack trace populated
	steArrField, ok := th.FieldTable["stackTrace"]
	if !ok {
		t.Fatalf("stackTrace field missing")
	}
	steArrayObj := steArrField.Fvalue.(*object.Object)
	entries := steArrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(entries) == 0 {
		t.Fatalf("stackTrace array empty")
	}
	ste := entries[0].FieldTable
	if ste["declaringClass"].Fvalue.(string) != "java/testClass" {
		t.Errorf("declaringClass mismatch: %v", ste["declaringClass"].Fvalue)
	}
	if ste["methodName"].Fvalue.(string) != "java/testClass.test" {
		t.Errorf("methodName mismatch: %v", ste["methodName"].Fvalue)
	}
	if ste["fileName"].Fvalue.(string) != "testClass.java" {
		t.Errorf("fileName mismatch: %v", ste["fileName"].Fvalue)
	}
	if ste["classLoaderName"].Fvalue.(string) != "testLoader" {
		t.Errorf("classLoaderName mismatch: %v", ste["classLoaderName"].Fvalue)
	}
}

func TestThrowableInitString_SetsMessageAndStack(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	// Build frame and stack
	f := frames.CreateFrame(2)
	f.Thread = 1
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"
	clData := classloader.ClData{
		Name:            "",
		SuperclassIndex: types.StringPoolObjectIndex,
		Module:          "test module",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Attributes:      nil,
		SourceFile:      "testClass.java",
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          0,
	}
	klass := classloader.Klass{Loader: "testLoader", Data: &clData}
	classloader.MethAreaInsert("java/testClass", &klass)
	f.ClName = "java/testClass"
	f.MethName = "java/testClass.test"
	jvmStack := frames.CreateFrameStack()
	_ = frames.PushFrame(jvmStack, f)

	// Throwable and message
	cls := "java/lang/Throwable"
	th := object.MakeEmptyObjectWithClassName(&cls)
	msg := object.StringObjectFromGoString("hello")

	// set up instantiate function required by StackTraceElement.of()
	globPtr := globals.GetGlobalRef()
	globPtr.FuncInstantiateClass = InstantiateFillIn

	// Call constructor
	ret := throwableInitString([]interface{}{jvmStack, th, msg})
	if ret != nil {
		t.Errorf("throwableInitString returned non-nil: %v", ret)
	}

	// Validate message stored
	dm, ok := th.FieldTable["detailMessage"]
	if !ok {
		t.Fatalf("detailMessage field missing")
	}
	if object.GoStringFromStringObject(dm.Fvalue.(*object.Object)) != "hello" {
		t.Errorf("detailMessage content mismatch")
	}
	cause, ok := th.FieldTable["cause"]
	if !ok || cause.Fvalue != object.Null {
		t.Errorf("cause not null; got: %#v", cause)
	}

	// Validate stack trace populated
	steArrField, ok := th.FieldTable["stackTrace"]
	if !ok {
		t.Fatalf("stackTrace field missing")
	}
	steArrayObj := steArrField.Fvalue.(*object.Object)
	entries := steArrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(entries) == 0 {
		t.Fatalf("stackTrace array empty")
	}
	ste := entries[0].FieldTable
	if ste["declaringClass"].Fvalue.(string) != "java/testClass" {
		t.Errorf("declaringClass mismatch: %v", ste["declaringClass"].Fvalue)
	}
	if ste["methodName"].Fvalue.(string) != "java/testClass.test" {
		t.Errorf("methodName mismatch: %v", ste["methodName"].Fvalue)
	}
	if ste["fileName"].Fvalue.(string) != "testClass.java" {
		t.Errorf("fileName mismatch: %v", ste["fileName"].Fvalue)
	}
	if ste["classLoaderName"].Fvalue.(string) != "testLoader" {
		t.Errorf("classLoaderName mismatch: %v", ste["classLoaderName"].Fvalue)
	}
}
