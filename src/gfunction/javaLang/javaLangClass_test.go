/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

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

func TestGetPrimitiveClass_Boolean(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("boolean")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_UnrecognizedPrimitive(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("unknown")
	params := []interface{}{obj}
	result := getPrimitiveClass(params).(*ghelpers.GErrBlk)
	if (*result).ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException, got different errror")
	}
}

func TestGetPrimitiveClass_Byte(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("byte")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Char(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("char")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Double(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("double")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Float(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("float")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Int(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("int")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Long(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("long")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Short(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("short")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestGetPrimitiveClass_Void(t *testing.T) {
	setup()

	obj := object.StringObjectFromGoString("void")
	params := []interface{}{obj}
	result := getPrimitiveClass(params)
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestSimpleClassLoadByName(t *testing.T) {
	setup()
	k, err := simpleClassLoadByName("java/lang/String")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if k == nil {
		t.Errorf("Expected *classloader.Klass for java/lang/String, got nil")
	}
}

func TestSimpleClassLoadByName_Error(t *testing.T) {
	setup()

	// ghelpers.Trap the error message written out due to not being able to find the class
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	k, err := simpleClassLoadByName("java/lang/No-Such-Class")

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "java.lang.ClassNotFoundException") {
		t.Errorf("Unexpected error message, got %s", errMsg)
	}

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if k != nil {
		t.Errorf("Expected nil data for non-existent class, got %v", k)
	}
}

func TestAssertionsEnabledStatus_Disabled(t *testing.T) {
	setup()
	statics.LoadProgramStatics()
	result := classGetAssertionsEnabledStatus(nil)
	if result != types.JavaBoolFalse {
		t.Errorf("Expected false, got %v", result)
	}
}

func TestAssertionsEnabledStatus_Enabled(t *testing.T) {
	setup()
	_ = statics.AddStatic("main.$assertionsDisabled",
		statics.Static{Type: types.Int, Value: types.JavaBoolFalse})
	result := classGetAssertionsEnabledStatus(nil)
	if result != types.JavaBoolTrue {
		t.Errorf("Expected true, got %v", result)
	}
}
