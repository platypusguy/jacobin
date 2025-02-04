/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/object"
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
	result := getPrimitiveClass(params).(*GErrBlk)
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
