/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"container/list"
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

// newClassObjWithNameObj creates a java/lang/Class-style object whose
// "name" field is a *object.Object (a Java String object), as is produced
// by the real classloader for loaded classes.
func newClassObjWithNameObj(name string) *object.Object {
	o := object.MakeEmptyObject()
	o.KlassName = types.StringPoolJavaLangClassIndex
	o.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString(name)}
	return o
}

// newClassObjWithNameString creates a java/lang/Class-style object whose
// "name" field is a plain Golang string, as several functions in
// javaLangClass.go expect.
func newClassObjWithNameString(name string) *object.Object {
	o := object.MakeEmptyObject()
	o.KlassName = types.StringPoolJavaLangClassIndex
	o.FieldTable["name"] = object.Field{Ftype: types.GolangString, Fvalue: name}
	return o
}

// newClassObjWithKlass creates a java/lang/Class-style object with both a
// "name" field (as a Java string object) and a "$klass" field pointing to
// classloader.ClData, as used by classGetModifiers, classIsEnum, etc.
func newClassObjWithKlass(name string, cd *classloader.ClData) *object.Object {
	o := newClassObjWithNameObj(name)
	o.FieldTable["$klass"] = object.Field{Ftype: types.RawGoPointer, Fvalue: cd}
	return o
}

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

func TestIsPrimitive_False(t *testing.T) {
	setup()
	className := "java/lang/Integer"
	nonPrimitive := object.MakeEmptyObjectWithClassName(&className)
	ret := classIsPrimitive([]interface{}{nonPrimitive})
	if ret != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse for non-primitive class, got %v", ret)
	}
}

func TestIsPrimitive_True(t *testing.T) {
	setup()
	className := "int"
	nonPrimitive := object.MakeEmptyObjectWithClassName(&className)
	ret := classIsPrimitive([]interface{}{nonPrimitive})
	if ret == types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolTrue for isPrimitive() with class 'int', got %v", ret)
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

func TestIsNumeric(t *testing.T) {
	if !isNumeric("123") {
		t.Errorf("Expected true for '123'")
	}
	if isNumeric("") {
		t.Errorf("Expected false for empty string")
	}
	if isNumeric("12a") {
		t.Errorf("Expected false for '12a'")
	}
}

func TestClassClinitIsh(t *testing.T) {
	setup()
	unnamedModule = nil
	ClassClinitIsh()
	if unnamedModule == nil {
		t.Errorf("Expected unnamedModule to be initialized")
	}
	if object.GoStringFromStringPoolIndex(unnamedModule.KlassName) != classNameModule {
		t.Errorf("Expected KlassName %s, got %s", classNameModule,
			object.GoStringFromStringPoolIndex(unnamedModule.KlassName))
	}
}

func TestClassGetModule(t *testing.T) {
	setup()
	unnamedModule = nil
	result := classGetModule(nil)
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk when unnamedModule is nil, got %T", result)
	}

	ClassClinitIsh()
	result = classGetModule(nil)
	if result != unnamedModule {
		t.Errorf("Expected unnamedModule, got %v", result)
	}
}

// === classDescriptorString ===

func TestClassDescriptorString_InvalidObject(t *testing.T) {
	setup()
	result := classDescriptorString([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassDescriptorString_ArrayFormat(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[I")
	result := classDescriptorString([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "[I" {
		t.Errorf("Expected [I, got %s", got)
	}
}

func TestClassDescriptorString_Primitive(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("int")
	result := classDescriptorString([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "I" {
		t.Errorf("Expected I, got %s", got)
	}
}

func TestClassDescriptorString_Reference(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("java/lang/String")
	result := classDescriptorString([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "Ljava/lang/String;" {
		t.Errorf("Expected Ljava/lang/String;, got %s", got)
	}
}

// === classGetCanonicalName ===

func TestClassGetCanonicalName_InvalidObject(t *testing.T) {
	setup()
	result := classGetCanonicalName([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetCanonicalName_RegularClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("java/lang/String")
	result := classGetCanonicalName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java.lang.String" {
		t.Errorf("Expected java.lang.String, got %s", got)
	}
}

func TestClassGetCanonicalName_PrimitiveArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[I")
	result := classGetCanonicalName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "int[]" {
		t.Errorf("Expected int[], got %s", got)
	}
}

func TestClassGetCanonicalName_ObjectArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[Ljava/lang/String;")
	result := classGetCanonicalName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java/lang/String[]" {
		t.Errorf("Expected java/lang/String[], got %s", got)
	}
}

func TestClassGetCanonicalName_AnonymousClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("com/example/Outer$1")
	result := classGetCanonicalName([]interface{}{obj})
	if result != nil {
		t.Errorf("Expected nil for anonymous class, got %v", result)
	}
}

func TestClassGetCanonicalName_InnerClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("com/example/Outer$Inner")
	result := classGetCanonicalName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "com/example/Outer.Inner" {
		t.Errorf("Expected com/example/Outer.Inner, got %s", got)
	}
}

// === classGetDeclaringClass ===

func TestClassGetDeclaringClass_InvalidObject(t *testing.T) {
	setup()
	result := classGetDeclaringClass([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetDeclaringClass_NotInnerClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("java/lang/String")
	result := classGetDeclaringClass([]interface{}{obj})
	if result != nil {
		t.Errorf("Expected nil for non-inner class, got %v", result)
	}
}

func TestClassGetDeclaringClass_UnknownOuterClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("com/nonexistent/Outer$Inner")
	result := classGetDeclaringClass([]interface{}{obj})
	if result != nil {
		t.Errorf("Expected nil when outer class can't be loaded, got %v", result)
	}
}

// === classGetPackageName ===

func TestClassGetPackageName_InvalidObject(t *testing.T) {
	setup()
	result := classGetPackageName([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetPackageName_RegularClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("java/lang/String")
	result := classGetPackageName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java.lang" {
		t.Errorf("Expected java.lang, got %s", got)
	}
}

func TestClassGetPackageName_DefaultPackage(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("Foo")
	result := classGetPackageName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "" {
		t.Errorf("Expected empty string, got %s", got)
	}
}

func TestClassGetPackageName_PrimitiveArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[I")
	result := classGetPackageName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "" {
		t.Errorf("Expected empty string, got %s", got)
	}
}

func TestClassGetPackageName_ObjectArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[Ljava/lang/String;")
	result := classGetPackageName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java.lang" {
		t.Errorf("Expected java.lang, got %s", got)
	}
}

// === classGetTypeName ===

func TestClassGetTypeName_InvalidObject(t *testing.T) {
	setup()
	result := classGetTypeName([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetTypeName_RegularClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("java/lang/String")
	result := classGetTypeName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java.lang.String" {
		t.Errorf("Expected java.lang.String, got %s", got)
	}
}

func TestClassGetTypeName_PrimitiveArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[I")
	result := classGetTypeName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "int[]" {
		t.Errorf("Expected int[], got %s", got)
	}
}

func TestClassGetTypeName_ObjectArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("[Ljava/lang/String;")
	result := classGetTypeName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java.lang.String[]" {
		t.Errorf("Expected java.lang.String[], got %s", got)
	}
}

// === classToString ===

func TestClassToString_InvalidObject(t *testing.T) {
	setup()
	result := classToString([]any{nil})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "null" {
		t.Errorf("Expected null, got %s", got)
	}
}

func TestClassToString_ValidObject(t *testing.T) {
	setup()
	obj := newClassObjWithNameString("java/lang/String")
	result := classToString([]any{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "class java/lang/String" {
		t.Errorf("Expected 'class java/lang/String', got %s", got)
	}
}

// === ClassGetName ===

func TestClassGetName_InvalidObject(t *testing.T) {
	setup()
	result := ClassGetName([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetName_Valid(t *testing.T) {
	setup()
	obj := newClassObjWithNameObj("java/lang/String")
	result := ClassGetName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "java.lang.String" {
		t.Errorf("Expected java.lang.String, got %s", got)
	}
}

// === ClassGetSimpleName ===

func TestClassGetSimpleName_InvalidObject(t *testing.T) {
	setup()
	result := ClassGetSimpleName([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetSimpleName_RegularClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameObj("java/lang/String")
	result := ClassGetSimpleName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "String" {
		t.Errorf("Expected String, got %s", got)
	}
}

func TestClassGetSimpleName_InnerClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameObj("com/example/Outer$Inner")
	result := ClassGetSimpleName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "Inner" {
		t.Errorf("Expected Inner, got %s", got)
	}
}

func TestClassGetSimpleName_AnonymousClass(t *testing.T) {
	setup()
	obj := newClassObjWithNameObj("com/example/Outer$1")
	result := ClassGetSimpleName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "" {
		t.Errorf("Expected empty string, got %s", got)
	}
}

func TestClassGetSimpleName_PrimitiveArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameObj("[I")
	result := ClassGetSimpleName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "int[]" {
		t.Errorf("Expected int[], got %s", got)
	}
}

func TestClassGetSimpleName_ObjectArray(t *testing.T) {
	setup()
	obj := newClassObjWithNameObj("[Ljava/lang/String;")
	result := ClassGetSimpleName([]interface{}{obj})
	got := object.GoStringFromStringObject(result.(*object.Object))
	if got != "String[]" {
		t.Errorf("Expected String[], got %s", got)
	}
}

// === classCast & classIsInstance ===

func TestClassIsInstance_InvalidClassObject(t *testing.T) {
	setup()
	result := classIsInstance([]interface{}{nil, nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassIsInstance_NullObject(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	classObj := object.MakeEmptyObjectWithClassName(&classNm)
	result := classIsInstance([]interface{}{classObj, nil})
	if result != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse for null instance, got %v", result)
	}
}

func TestClassIsInstance_Match(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	classObj := object.MakeEmptyObjectWithClassName(&classNm)
	instance := object.MakeEmptyObjectWithClassName(&classNm)
	result := classIsInstance([]interface{}{classObj, instance})
	if result != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue, got %v", result)
	}
}

func TestClassIsInstance_NoMatch(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	otherNm := "java/lang/Integer"
	classObj := object.MakeEmptyObjectWithClassName(&classNm)
	instance := object.MakeEmptyObjectWithClassName(&otherNm)
	result := classIsInstance([]interface{}{classObj, instance})
	if result != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse, got %v", result)
	}
}

func TestClassCast_InvalidClassObject(t *testing.T) {
	setup()
	result := classCast([]interface{}{nil, nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassCast_NullObjectToCast(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	classObj := object.MakeEmptyObjectWithClassName(&classNm)
	result := classCast([]interface{}{classObj, object.Null})
	if result != nil {
		t.Errorf("Expected nil for casting null, got %v", result)
	}
}

func TestClassCast_ValidCast(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	classObj := object.MakeEmptyObjectWithClassName(&classNm)
	instance := object.MakeEmptyObjectWithClassName(&classNm)
	result := classCast([]interface{}{classObj, instance})
	if result != instance {
		t.Errorf("Expected same instance returned, got %v", result)
	}
}

func TestClassCast_InvalidCast(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	otherNm := "java/lang/Integer"
	classObj := object.MakeEmptyObjectWithClassName(&classNm)
	instance := object.MakeEmptyObjectWithClassName(&otherNm)
	result := classCast([]interface{}{classObj, instance})
	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	} else if errBlk.ExceptionType != excNames.ClassCastException {
		t.Errorf("Expected ClassCastException, got %v", errBlk.ExceptionType)
	}
}

// === classIsArray ===

func TestClassIsArray_InvalidObject(t *testing.T) {
	setup()
	result := classIsArray([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassIsArray_True(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: types.IntArray, Fvalue: []int64{}}
	result := classIsArray([]interface{}{obj})
	if result != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue, got %v", result)
	}
}

func TestClassIsArray_False(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	result := classIsArray([]interface{}{obj})
	if result != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse, got %v", result)
	}
}

// === classIsSealed ===

func TestClassIsSealed_InvalidObject(t *testing.T) {
	setup()
	result := classIsSealed([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassIsSealed_AlwaysFalse(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	result := classIsSealed([]interface{}{obj})
	if result != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse, got %v", result)
	}
}

// === $klass-based functions: classGetModifiers, classIsEnum, classIsInterface, classIsSynthetic, classGetPackage ===

func TestClassGetModifiers_InvalidObject(t *testing.T) {
	setup()
	result := classGetModifiers([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetModifiers_PublicFinal(t *testing.T) {
	setup()
	cd := &classloader.ClData{Access: classloader.AccessFlags{ClassIsPublic: true, ClassIsFinal: true}}
	obj := newClassObjWithKlass("some/Class", cd)
	result := classGetModifiers([]interface{}{obj}).(int64)
	if result != 0x0001|0x0010 {
		t.Errorf("Expected 0x11, got 0x%x", result)
	}
}

func TestClassIsEnum_InvalidObject(t *testing.T) {
	setup()
	result := classIsEnum([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassIsEnum_True(t *testing.T) {
	setup()
	cd := &classloader.ClData{Access: classloader.AccessFlags{ClassIsEnum: true}}
	obj := newClassObjWithKlass("some/Enum", cd)
	result := classIsEnum([]interface{}{obj})
	if result != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue, got %v", result)
	}
}

func TestClassIsEnum_False(t *testing.T) {
	setup()
	cd := &classloader.ClData{Access: classloader.AccessFlags{}}
	obj := newClassObjWithKlass("some/Class", cd)
	result := classIsEnum([]interface{}{obj})
	if result != types.JavaBoolFalse {
		t.Errorf("Expected JavaBoolFalse, got %v", result)
	}
}

func TestClassIsInterface_InvalidObject(t *testing.T) {
	setup()
	result := classIsInterface([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassIsInterface_True(t *testing.T) {
	setup()
	cd := &classloader.ClData{Access: classloader.AccessFlags{ClassIsInterface: true}}
	obj := newClassObjWithKlass("some/Iface", cd)
	result := classIsInterface([]interface{}{obj})
	if result != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue, got %v", result)
	}
}

func TestClassIsSynthetic_InvalidObject(t *testing.T) {
	setup()
	result := classIsSynthetic([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassIsSynthetic_True(t *testing.T) {
	setup()
	cd := &classloader.ClData{Access: classloader.AccessFlags{ClassIsSynthetic: true}}
	obj := newClassObjWithKlass("some/Class", cd)
	result := classIsSynthetic([]interface{}{obj})
	if result != types.JavaBoolTrue {
		t.Errorf("Expected JavaBoolTrue, got %v", result)
	}
}

func TestClassGetPackage_InvalidObject(t *testing.T) {
	setup()
	result := classGetPackage([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetPackage_NoPackage(t *testing.T) {
	setup()
	cd := &classloader.ClData{Pkg: ""}
	obj := newClassObjWithKlass("Foo", cd)
	result := classGetPackage([]interface{}{obj})
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestClassGetPackage_WithPackage(t *testing.T) {
	setup()
	cd := &classloader.ClData{Pkg: "java/lang"}
	obj := newClassObjWithKlass("java/lang/String", cd)
	result := classGetPackage([]interface{}{obj})
	// current implementation always returns nil (TODO in source)
	if result != nil {
		t.Errorf("Expected nil (per current TODO implementation), got %v", result)
	}
}

// === getComponentType / classComponentType ===

func TestGetComponentType_InvalidObject(t *testing.T) {
	setup()
	result := getComponentType([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestGetComponentType_NotArray(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	result := getComponentType([]interface{}{obj})
	if !object.IsNull(result) {
		t.Errorf("Expected null, got %v", result)
	}
}

func TestGetComponentType_PrimitiveArray(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: types.IntArray, Fvalue: []int64{}}
	result := getComponentType([]interface{}{obj})
	if _, ok := result.(*classloader.Klass); !ok {
		t.Errorf("Expected *classloader.Klass, got %T", result)
	}
}

func TestClassComponentType_DelegatesToGetComponentType(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	obj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	result := classComponentType([]interface{}{obj})
	if !object.IsNull(result) {
		t.Errorf("Expected null, got %v", result)
	}
}

// === classGetInterfaces ===

func TestClassGetInterfaces_InvalidObject(t *testing.T) {
	setup()
	result := classGetInterfaces([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetInterfaces_NoInterfaces(t *testing.T) {
	setup()
	// java/lang/Object has no interfaces
	obj := newClassObjWithNameObj("java/lang/Object")
	result := classGetInterfaces([]interface{}{obj})
	arrObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", result)
	}
	rawArray, ok := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if !ok {
		t.Fatalf("Expected []*object.Object, got %T", arrObj.FieldTable["value"].Fvalue)
	}
	if len(rawArray) != 0 {
		t.Errorf("Expected empty array, got length %d", len(rawArray))
	}
}

// === classGetSuperclass ===

func TestClassGetSuperclass_InvalidObject(t *testing.T) {
	setup()
	result := classGetSuperclass([]interface{}{nil})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetSuperclass_ObjectReturnsNil(t *testing.T) {
	setup()
	classNm := "java/lang/Object"
	cd := &classloader.ClData{}
	obj := object.MakeEmptyObjectWithClassName(&classNm)
	obj.FieldTable["$klass"] = object.Field{Ftype: types.RawGoPointer, Fvalue: cd}
	result := classGetSuperclass([]interface{}{obj})
	if result != nil {
		t.Errorf("Expected nil for java/lang/Object, got %v", result)
	}
}

func TestClassGetSuperclass_RegularClass(t *testing.T) {
	setup()
	k, err := simpleClassLoadByName("java/lang/String")
	if err != nil {
		t.Fatalf("Failed to load java/lang/String: %v", err)
	}
	result := classGetSuperclass([]interface{}{k.Data.ClassObject})
	if result == nil {
		t.Errorf("Expected non-nil superclass for java/lang/String")
	}
}

// === classGetField ===

func TestClassGetField_InvalidObject(t *testing.T) {
	setup()
	result := classGetField([]interface{}{nil, "someField"})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk, got %T", result)
	}
}

func TestClassGetField_NullFieldName(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	result := classGetField([]interface{}{obj, object.Null})
	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected GErrBlk, got %T", result)
	}
	if errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException, got %v", errBlk.ExceptionType)
	}
}

func TestClassGetField_FieldNotFound(t *testing.T) {
	setup()
	classNm := "java/lang/String"
	obj := object.MakeEmptyObjectWithClassName(&classNm)
	result := classGetField([]interface{}{obj, "nonExistentField"})
	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected GErrBlk, got %T", result)
	}
	if errBlk.ExceptionType != excNames.NoSuchFieldException {
		t.Errorf("Expected NoSuchFieldException, got %v", errBlk.ExceptionType)
	}
}

func TestClassGetField_FieldFound(t *testing.T) {
	setup()
	obj := object.MakeEmptyObject()
	obj.FieldTable["myField"] = object.Field{Ftype: types.Int, Fvalue: int64(42)}
	result := classGetField([]interface{}{obj, "myField"})
	fld, ok := result.(*Field)
	if !ok {
		t.Fatalf("Expected *Field, got %T", result)
	}
	if fld.Name != "myField" {
		t.Errorf("Expected field name 'myField', got %s", fld.Name)
	}
}

// === classForNameLZL ===

func TestClassForNameLZL_ClassNotFound(t *testing.T) {
	setup()
	classNameObj := object.StringObjectFromGoString("java.lang.No-Such-Class")

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	params := []interface{}{classNameObj, types.JavaBoolFalse, object.Null, list.New()}
	result := classForNameLZL(params)

	_ = w.Close()
	_, _ = io.ReadAll(r)
	os.Stderr = normalStderr

	errBlk, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("Expected GErrBlk, got %T", result)
	}
	if errBlk.ExceptionType != excNames.ClassNotFoundException {
		t.Errorf("Expected ClassNotFoundException, got %v", errBlk.ExceptionType)
	}
}

func TestClassForNameLZL_ClassFoundNoInit(t *testing.T) {
	setup()
	classNameObj := object.StringObjectFromGoString("java.lang.String")
	params := []interface{}{classNameObj, types.JavaBoolFalse, object.Null, list.New()}
	result := classForNameLZL(params)
	if _, ok := result.(*object.Object); !ok {
		t.Errorf("Expected *object.Object (java/lang/Class), got %T", result)
	}
}
