/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/classloader"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"testing"
)

func TestLoad_Lang_Invoke_StringConcatFactory(t *testing.T) {
	globals.InitGlobals("test")
	Load_Lang_Invoke_StringConcatFactory()

	clinitSig := "java/lang/invoke/StringConcatFactory.<clinit>()V"
	if _, ok := ghelpers.MethodSignatures[clinitSig]; !ok {
		t.Errorf("Expected signature %s to be registered", clinitSig)
	}

	makeConcatSig := "java/lang/invoke/StringConcatFactory.makeConcatWithConstants(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"
	if gMeth, ok := ghelpers.MethodSignatures[makeConcatSig]; !ok {
		t.Errorf("Expected signature %s to be registered", makeConcatSig)
	} else if gMeth.ParamSlots != 5 {
		t.Errorf("Expected ParamSlots 5, got %d", gMeth.ParamSlots)
	}
}

func TestStringConcatFactory_MakeConcatWithConstants_IncorrectParamCount(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	params := []any{
		object.StringObjectFromGoString("caller/MyClass"),
		object.StringObjectFromGoString("makeConcatWithConstants"),
	}

	result := stringConcatFactory_makeConcatWithConstants(params)
	if result == nil {
		t.Errorf("Expected error block for incorrect parameter count, got nil")
	}
}

func TestStringConcatFactory_MakeConcatWithConstants_Success(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()

	// Build MethodType object with ptypes array
	ptypeObj := &object.Object{FieldTable: make(map[string]object.Field)}
	clData := &classloader.ClData{Name: "java/lang/String"}
	ptypeObj.FieldTable["$klass"] = object.Field{Ftype: "Lclassloader/ClData;", Fvalue: clData}

	ptypesArrayObj := &object.Object{FieldTable: make(map[string]object.Field)}
	ptypesArrayObj.FieldTable["value"] = object.Field{Ftype: "[Ljava/lang/Class;", Fvalue: []*object.Object{ptypeObj}}

	methTypeObj := &object.Object{FieldTable: make(map[string]object.Field)}
	methTypeObj.FieldTable["ptypes"] = object.Field{Ftype: "[Ljava/lang/Class;", Fvalue: ptypesArrayObj}

	params := []any{
		object.StringObjectFromGoString("caller/MyClass"),
		object.StringObjectFromGoString("makeConcatWithConstants"),
		methTypeObj,
		object.StringObjectFromGoString("recipe \x01"),
		[]*object.Object{object.StringObjectFromGoString("extraConstant")},
	}

	globals.TraceInst = true
	defer func() { globals.TraceInst = false }()

	result := stringConcatFactory_makeConcatWithConstants(params)
	callSiteObj, ok := result.(*object.Object)
	if !ok || callSiteObj == nil {
		t.Fatalf("Expected *object.Object return, got %T (%v)", result, result)
	}

	targetField, ok := callSiteObj.FieldTable["target"]
	if !ok || targetField.Fvalue == nil {
		t.Errorf("Expected 'target' field in CallSite object")
	}
}
