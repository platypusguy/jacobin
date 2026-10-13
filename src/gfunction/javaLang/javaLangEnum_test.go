/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"testing"
)

func TestLoad_Lang_Enum_RegistersMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Enum()

	expected := []string{
		"java/lang/Enum.<clinit>()V",
		"java/lang/Enum.<init>(Ljava/lang/String;I)V",
		"java/lang/Enum.clone()Ljava/lang/Object;",
		"java/lang/Enum.compareTo(Ljava/lang/Enum;)I",
		"java/lang/Enum.compareTo(Ljava/lang/Object;)I",
		"java/lang/Enum.equals(Ljava/lang/Object;)Z",
		"java/lang/Enum.finalize()V",
		"java/lang/Enum.getDeclaringClass()Ljava/lang/Class;",
		"java/lang/Enum.hashCode()I",
		"java/lang/Enum.name()Ljava/lang/String;",
		"java/lang/Enum.ordinal()I",
		"java/lang/Enum.toString()Ljava/lang/String;",
		"java/lang/Enum.valueOf(Ljava/lang/Class;Ljava/lang/String;)Ljava/lang/Enum;",
	}
	for _, sig := range expected {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("expected %s to be registered", sig)
		}
	}
}

func TestLoad_Lang_Enum_Traps(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Enum()

	finalizeMeth := ghelpers.MethodSignatures["java/lang/Enum.finalize()V"]
	res := finalizeMeth.GFunction(nil)
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected finalize() to be trapped, got %v", res)
	}

	getDeclClass := ghelpers.MethodSignatures["java/lang/Enum.getDeclaringClass()Ljava/lang/Class;"]
	res = getDeclClass.GFunction(nil)
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected getDeclaringClass() to be trapped, got %v", res)
	}

	valueOfMeth := ghelpers.MethodSignatures["java/lang/Enum.valueOf(Ljava/lang/Class;Ljava/lang/String;)Ljava/lang/Enum;"]
	res = valueOfMeth.GFunction(nil)
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected valueOf() to be trapped, got %v", res)
	}
}

func makeEnumConstant(className, name string, ordinal int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&className)
	nameObj := object.StringObjectFromGoString(name)
	result := enumInit([]interface{}{obj, nameObj, ordinal})
	if result != nil {
		panic("enumInit failed in test helper")
	}
	return obj
}

func TestEnumInit_Success(t *testing.T) {
	className := "TestEnum"
	obj := object.MakeEmptyObjectWithClassName(&className)
	nameObj := object.StringObjectFromGoString("RED")
	result := enumInit([]interface{}{obj, nameObj, int64(0)})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
	fld, exists := obj.FieldTable["name"]
	if !exists {
		t.Fatalf("expected name field to be set")
	}
	strObj := fld.Fvalue.(*object.Object)
	if object.GoStringFromStringObject(strObj) != "RED" {
		t.Errorf("expected name RED, got %s", object.GoStringFromStringObject(strObj))
	}
	ordFld, exists := obj.FieldTable["ordinal"]
	if !exists {
		t.Fatalf("expected ordinal field to be set")
	}
	if ordFld.Fvalue.(int64) != 0 {
		t.Errorf("expected ordinal 0, got %v", ordFld.Fvalue)
	}
}

func TestEnumInit_NullReceiver(t *testing.T) {
	nameObj := object.StringObjectFromGoString("RED")
	result := enumInit([]interface{}{object.Null, nameObj, int64(0)})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumInit_InvalidName(t *testing.T) {
	className := "TestEnum"
	obj := object.MakeEmptyObjectWithClassName(&className)
	result := enumInit([]interface{}{obj, "not-a-string-object", int64(0)})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %d", gErr.ExceptionType)
	}
}

func TestEnumInit_InvalidOrdinal(t *testing.T) {
	className := "TestEnum"
	obj := object.MakeEmptyObjectWithClassName(&className)
	nameObj := object.StringObjectFromGoString("RED")
	result := enumInit([]interface{}{obj, nameObj, "not-an-int"})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %d", gErr.ExceptionType)
	}
}

func TestEnumName_Success(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "GREEN", 1)
	result := enumName([]interface{}{obj})
	strObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", result)
	}
	if object.GoStringFromStringObject(strObj) != "GREEN" {
		t.Errorf("expected GREEN, got %s", object.GoStringFromStringObject(strObj))
	}
}

func TestEnumName_NullReceiver(t *testing.T) {
	result := enumName([]interface{}{object.Null})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumName_MissingField(t *testing.T) {
	className := "TestEnum"
	obj := object.MakeEmptyObjectWithClassName(&className)
	result := enumName([]interface{}{obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
}

func TestEnumOrdinal_Success(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "BLUE", 2)
	result := enumOrdinal([]interface{}{obj})
	ordinal, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}
	if ordinal != 2 {
		t.Errorf("expected 2, got %d", ordinal)
	}
}

func TestEnumOrdinal_NullReceiver(t *testing.T) {
	result := enumOrdinal([]interface{}{object.Null})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumOrdinal_MissingField(t *testing.T) {
	className := "TestEnum"
	obj := object.MakeEmptyObjectWithClassName(&className)
	result := enumOrdinal([]interface{}{obj})
	if _, ok := result.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
}

func TestEnumToString_Success(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "YELLOW", 3)
	result := enumToString([]interface{}{obj})
	strObj, ok := result.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", result)
	}
	if object.GoStringFromStringObject(strObj) != "YELLOW" {
		t.Errorf("expected YELLOW, got %s", object.GoStringFromStringObject(strObj))
	}
}

func TestEnumEquals_SameObject(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "A", 0)
	result := enumEquals([]interface{}{obj, obj})
	if !isJavaTrue(result) {
		t.Errorf("expected true for same object, got %v", result)
	}
}

func TestEnumEquals_DifferentObjects(t *testing.T) {
	obj1 := makeEnumConstant("TestEnum", "A", 0)
	obj2 := makeEnumConstant("TestEnum", "B", 1)
	result := enumEquals([]interface{}{obj1, obj2})
	if isJavaTrue(result) {
		t.Errorf("expected false for different objects, got %v", result)
	}
}

func TestEnumEquals_NullOther(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "A", 0)
	result := enumEquals([]interface{}{obj, object.Null})
	if isJavaTrue(result) {
		t.Errorf("expected false for null other, got %v", result)
	}
}

func TestEnumEquals_NullReceiver(t *testing.T) {
	result := enumEquals([]interface{}{object.Null, object.Null})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumHashCode_Success(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "A", 0)
	result := enumHashCode([]interface{}{obj})
	if _, ok := result.(int64); !ok {
		t.Fatalf("expected int64, got %T", result)
	}
}

func TestEnumHashCode_NullReceiver(t *testing.T) {
	result := enumHashCode([]interface{}{object.Null})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumCompareTo_Equal(t *testing.T) {
	obj1 := makeEnumConstant("TestEnum", "A", 0)
	obj2 := makeEnumConstant("TestEnum", "A", 0)
	result := enumCompareTo([]interface{}{obj1, obj2})
	cmp, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}
	if cmp != 0 {
		t.Errorf("expected 0, got %d", cmp)
	}
}

func TestEnumCompareTo_Less(t *testing.T) {
	obj1 := makeEnumConstant("TestEnum", "A", 0)
	obj2 := makeEnumConstant("TestEnum", "B", 1)
	result := enumCompareTo([]interface{}{obj1, obj2})
	cmp, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}
	if cmp >= 0 {
		t.Errorf("expected negative, got %d", cmp)
	}
}

func TestEnumCompareTo_Greater(t *testing.T) {
	obj1 := makeEnumConstant("TestEnum", "B", 1)
	obj2 := makeEnumConstant("TestEnum", "A", 0)
	result := enumCompareTo([]interface{}{obj1, obj2})
	cmp, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}
	if cmp <= 0 {
		t.Errorf("expected positive, got %d", cmp)
	}
}

func TestEnumCompareTo_DifferentTypes(t *testing.T) {
	obj1 := makeEnumConstant("TestEnum1", "A", 0)
	obj2 := makeEnumConstant("TestEnum2", "A", 0)
	result := enumCompareTo([]interface{}{obj1, obj2})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.ClassCastException {
		t.Errorf("expected ClassCastException, got %d", gErr.ExceptionType)
	}
}

func TestEnumCompareTo_NullReceiver(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "A", 0)
	result := enumCompareTo([]interface{}{object.Null, obj})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumCompareTo_NullArgument(t *testing.T) {
	obj := makeEnumConstant("TestEnum", "A", 0)
	result := enumCompareTo([]interface{}{obj, object.Null})
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %d", gErr.ExceptionType)
	}
}

func TestEnumClone_AlwaysThrows(t *testing.T) {
	result := enumClone(nil)
	gErr, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk, got %T", result)
	}
	if gErr.ExceptionType != excNames.CloneNotSupportedException {
		t.Errorf("expected CloneNotSupportedException, got %d", gErr.ExceptionType)
	}
}

// isJavaTrue is a small helper to interpret the boolean-ish return values
// used across the codebase (types.JavaBoolTrue/False are int64-based).
func isJavaTrue(v interface{}) bool {
	iv, ok := v.(int64)
	if !ok {
		return false
	}
	return iv != 0
}
