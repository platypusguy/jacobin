/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"testing"
)

func setRecordFieldOrder(t *testing.T, className string, names ...string) {
	t.Helper()
	if classloader.MethArea == nil {
		classloader.InitMethodArea()
	}
	previous := classloader.MethAreaFetch(className)
	fields := make([]classloader.Field, len(names))
	for i, name := range names {
		fields[i] = classloader.Field{NameStr: name}
	}
	klass := &classloader.Klass{
		Data: &classloader.ClData{
			Name:            className,
			Fields:          fields,
			SuperclassIndex: types.InvalidStringIndex,
		},
	}
	if previous == nil {
		classloader.MethAreaInsert(className, klass)
	} else {
		classloader.MethAreaUpdate(className, klass)
	}
	t.Cleanup(func() {
		if previous == nil {
			classloader.MethAreaDelete(className)
		} else {
			classloader.MethAreaUpdate(className, previous)
		}
	})
}

func TestLoad_Lang_Record_RegistersMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Record()

	expected := map[string]int{
		"java/lang/Record.<clinit>()V":                  0,
		"java/lang/Record.<init>()V":                    0,
		"java/lang/Record.equals(Ljava/lang/Object;)Z":  1,
		"java/lang/Record.hashCode()I":                  0,
		"java/lang/Record.toString()Ljava/lang/String;": 0,
	}

	for sig, slots := range expected {
		gmeth, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("expected %s to be registered", sig)
			continue
		}
		if gmeth.ParamSlots != slots {
			t.Errorf("for %s, expected %d slots, got %d", sig, slots, gmeth.ParamSlots)
		}
		if gmeth.GFunction == nil {
			t.Errorf("for %s, expected non-nil GFunction", sig)
		}
	}
}

func TestRecordInit(t *testing.T) {
	globals.InitGlobals("test")

	className := "com/example/Point"
	obj := object.MakeEmptyObjectWithClassName(&className)

	// Valid receiver
	res := recordInit([]any{obj})
	if res != nil {
		t.Errorf("expected nil return on recordInit, got %v", res)
	}

	// Missing receiver
	res = recordInit([]any{})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for empty params, got %v", res)
	}

	// Null receiver
	res = recordInit([]any{object.Null})
	errBlk, ok = res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for Null receiver, got %v", res)
	}

	// Invalid receiver type
	res = recordInit([]any{int64(42)})
	errBlk, ok = res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for invalid receiver type, got %v", res)
	}
}

func TestRecordEquals(t *testing.T) {
	globals.InitGlobals("test")

	className1 := "com/example/Point"
	className2 := "com/example/Point3D"

	p1 := object.MakeEmptyObjectWithClassName(&className1)
	p1.FieldTable["x"] = object.Field{Ftype: types.Int, Fvalue: int64(10)}
	p1.FieldTable["y"] = object.Field{Ftype: types.Int, Fvalue: int64(20)}

	p2 := object.MakeEmptyObjectWithClassName(&className1)
	p2.FieldTable["x"] = object.Field{Ftype: types.Int, Fvalue: int64(10)}
	p2.FieldTable["y"] = object.Field{Ftype: types.Int, Fvalue: int64(20)}

	p3 := object.MakeEmptyObjectWithClassName(&className1)
	p3.FieldTable["x"] = object.Field{Ftype: types.Int, Fvalue: int64(10)}
	p3.FieldTable["y"] = object.Field{Ftype: types.Int, Fvalue: int64(30)}

	pDiffClass := object.MakeEmptyObjectWithClassName(&className2)
	pDiffClass.FieldTable["x"] = object.Field{Ftype: types.Int, Fvalue: int64(10)}
	pDiffClass.FieldTable["y"] = object.Field{Ftype: types.Int, Fvalue: int64(20)}

	// Reflexive: p1.equals(p1) == true
	eq := recordEquals([]any{p1, p1})
	if eq != types.JavaBoolTrue {
		t.Errorf("expected reflexive equals to be true, got %v", eq)
	}

	// Symmetric: p1.equals(p2) == true && p2.equals(p1) == true
	eq = recordEquals([]any{p1, p2})
	if eq != types.JavaBoolTrue {
		t.Errorf("expected p1.equals(p2) to be true, got %v", eq)
	}
	eq = recordEquals([]any{p2, p1})
	if eq != types.JavaBoolTrue {
		t.Errorf("expected p2.equals(p1) to be true, got %v", eq)
	}

	// Different field values
	eq = recordEquals([]any{p1, p3})
	if eq != types.JavaBoolFalse {
		t.Errorf("expected p1.equals(p3) to be false, got %v", eq)
	}

	// Different class name
	eq = recordEquals([]any{p1, pDiffClass})
	if eq != types.JavaBoolFalse {
		t.Errorf("expected p1.equals(pDiffClass) to be false, got %v", eq)
	}

	// Compare with null or non-object
	eq = recordEquals([]any{p1, object.Null})
	if eq != types.JavaBoolFalse {
		t.Errorf("expected p1.equals(null) to be false, got %v", eq)
	}
	eq = recordEquals([]any{p1, nil})
	if eq != types.JavaBoolFalse {
		t.Errorf("expected p1.equals(nil) to be false, got %v", eq)
	}
	eq = recordEquals([]any{p1, int64(100)})
	if eq != types.JavaBoolFalse {
		t.Errorf("expected p1.equals(primitive) to be false, got %v", eq)
	}
	eq = recordEquals([]any{p1})
	if eq != types.JavaBoolFalse {
		t.Errorf("expected p1.equals() without arg to be false, got %v", eq)
	}

	// Null or invalid receiver
	errRes := recordEquals([]any{object.Null, p1})
	if _, ok := errRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk on null receiver, got %v", errRes)
	}
	errRes = recordEquals([]any{})
	if _, ok := errRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk on empty params, got %v", errRes)
	}
}

func TestRecordEquals_ReferenceComponentsUseJavaEquals(t *testing.T) {
	globals.InitGlobals("test")

	savedMethods := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = savedMethods }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	valueClass := "example/Value"
	ghelpers.MethodSignatures[valueClass+".equals(Ljava/lang/Object;)Z"] = ghelpers.GMeth{
		ParamSlots: 1,
		GFunction: func(params []any) any {
			if len(params) < 2 {
				return types.JavaBoolFalse
			}
			left, leftOK := params[0].(*object.Object)
			right, rightOK := params[1].(*object.Object)
			if !leftOK || !rightOK {
				return types.JavaBoolFalse
			}
			if left.FieldTable["value"].Fvalue == right.FieldTable["value"].Fvalue {
				return types.JavaBoolTrue
			}
			return types.JavaBoolFalse
		},
	}

	newValue := func(value int64) *object.Object {
		obj := object.MakeEmptyObjectWithClassName(&valueClass)
		obj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: value}
		return obj
	}

	recordClass := "example/Holder"
	left := object.MakeEmptyObjectWithClassName(&recordClass)
	left.FieldTable["value"] = object.Field{Ftype: valueClass, Fvalue: newValue(7)}
	equal := object.MakeEmptyObjectWithClassName(&recordClass)
	equal.FieldTable["value"] = object.Field{Ftype: valueClass, Fvalue: newValue(7)}
	different := object.MakeEmptyObjectWithClassName(&recordClass)
	different.FieldTable["value"] = object.Field{Ftype: valueClass, Fvalue: newValue(8)}

	if got := recordEquals([]any{left, equal}); got != types.JavaBoolTrue {
		t.Errorf("expected distinct reference components with equal Java values to compare equal, got %v", got)
	}
	if got := recordEquals([]any{left, different}); got != types.JavaBoolFalse {
		t.Errorf("expected distinct reference components with unequal Java values to compare unequal, got %v", got)
	}

	nullLeft := object.MakeEmptyObjectWithClassName(&recordClass)
	nullLeft.FieldTable["value"] = object.Field{Ftype: valueClass, Fvalue: object.Null}
	nullRight := object.MakeEmptyObjectWithClassName(&recordClass)
	nullRight.FieldTable["value"] = object.Field{Ftype: valueClass, Fvalue: object.Null}
	if got := recordEquals([]any{nullLeft, nullRight}); got != types.JavaBoolTrue {
		t.Errorf("expected two null reference components to compare equal, got %v", got)
	}
	if got := recordEquals([]any{left, nullLeft}); got != types.JavaBoolFalse {
		t.Errorf("expected null and non-null reference components to compare unequal, got %v", got)
	}

	ghelpers.MethodSignatures[valueClass+".equals(Ljava/lang/Object;)Z"] = ghelpers.GMeth{
		ParamSlots: 1,
		GFunction: func([]any) any {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "equals failed")
		},
	}
	if got, ok := recordEquals([]any{left, equal}).(*ghelpers.GErrBlk); !ok || got.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected equals GErrBlk to propagate, got %T (%v)", got, got)
	}
}

func TestRecordEquals_VariousTypes(t *testing.T) {
	globals.InitGlobals("test")

	className := "TestRecord"

	// Float NaN equality
	r1 := object.MakeEmptyObjectWithClassName(&className)
	r1.FieldTable["val"] = object.Field{Ftype: types.Double, Fvalue: math.NaN()}

	r2 := object.MakeEmptyObjectWithClassName(&className)
	r2.FieldTable["val"] = object.Field{Ftype: types.Double, Fvalue: math.NaN()}

	if eq := recordEquals([]any{r1, r2}); eq != types.JavaBoolTrue {
		t.Errorf("expected NaN equals NaN to be true in Record, got %v", eq)
	}

	// String objects equality
	strObj1 := object.StringObjectFromGoString("hello")
	strObj2 := object.StringObjectFromGoString("hello")
	strObj3 := object.StringObjectFromGoString("world")

	rs1 := object.MakeEmptyObjectWithClassName(&className)
	rs1.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: strObj1}

	rs2 := object.MakeEmptyObjectWithClassName(&className)
	rs2.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: strObj2}

	rs3 := object.MakeEmptyObjectWithClassName(&className)
	rs3.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: strObj3}

	if eq := recordEquals([]any{rs1, rs2}); eq != types.JavaBoolTrue {
		t.Errorf("expected string field equality to be true, got %v", eq)
	}
	if eq := recordEquals([]any{rs1, rs3}); eq != types.JavaBoolFalse {
		t.Errorf("expected string field inequality to be false, got %v", eq)
	}

	// Slices equality
	ra1 := object.MakeEmptyObjectWithClassName(&className)
	ra1.FieldTable["arr"] = object.Field{Ftype: "[I", Fvalue: []int64{1, 2, 3}}

	ra2 := object.MakeEmptyObjectWithClassName(&className)
	ra2.FieldTable["arr"] = object.Field{Ftype: "[I", Fvalue: []int64{1, 2, 3}}

	ra3 := object.MakeEmptyObjectWithClassName(&className)
	ra3.FieldTable["arr"] = object.Field{Ftype: "[I", Fvalue: []int64{1, 2, 4}}

	if eq := recordEquals([]any{ra1, ra2}); eq != types.JavaBoolTrue {
		t.Errorf("expected slice field equality to be true, got %v", eq)
	}
	if eq := recordEquals([]any{ra1, ra3}); eq != types.JavaBoolFalse {
		t.Errorf("expected slice field inequality to be false, got %v", eq)
	}
}

func TestRecordHashCode(t *testing.T) {
	globals.InitGlobals("test")

	className := "com/example/Person"
	setRecordFieldOrder(t, className, "name", "age")

	p1 := object.MakeEmptyObjectWithClassName(&className)
	p1.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Alice")}
	p1.FieldTable["age"] = object.Field{Ftype: types.Int, Fvalue: int64(30)}

	p2 := object.MakeEmptyObjectWithClassName(&className)
	p2.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Alice")}
	p2.FieldTable["age"] = object.Field{Ftype: types.Int, Fvalue: int64(30)}

	h1 := recordHashCode([]any{p1})
	h2 := recordHashCode([]any{p2})

	if _, ok := h1.(int64); !ok {
		t.Fatalf("expected int64 hash code, got %T (%v)", h1, h1)
	}
	if h1 != h2 {
		t.Errorf("expected equal records to have same hashCode: h1=%v, h2=%v", h1, h2)
	}
	expectedHash := int64(int32(31*(31+computeStringHashCode("Alice")) + 30))
	if h1 != expectedHash {
		t.Errorf("expected declaration-order hashCode %d, got %v", expectedHash, h1)
	}

	// Empty record
	pEmpty := object.MakeEmptyObjectWithClassName(&className)
	hEmpty := recordHashCode([]any{pEmpty})
	if hEmpty != int64(1) {
		t.Errorf("expected empty record hash code to be 1, got %v", hEmpty)
	}

	// Null or invalid receiver
	errRes := recordHashCode([]any{object.Null})
	if _, ok := errRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk on null receiver, got %v", errRes)
	}
	errRes = recordHashCode([]any{})
	if _, ok := errRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk on empty params, got %v", errRes)
	}
}

func TestRecordToString(t *testing.T) {
	globals.InitGlobals("test")

	className := "com/example/Person"
	setRecordFieldOrder(t, className, "name", "age", "active")
	p := object.MakeEmptyObjectWithClassName(&className)
	p.FieldTable["name"] = object.Field{Ftype: types.StringClassName, Fvalue: object.StringObjectFromGoString("Alice")}
	p.FieldTable["age"] = object.Field{Ftype: types.Int, Fvalue: int64(30)}
	p.FieldTable["active"] = object.Field{Ftype: types.Bool, Fvalue: types.JavaBoolTrue}

	res := recordToString([]any{p})
	strObj, ok := res.(*object.Object)
	if !ok || strObj == nil || object.IsNull(strObj) {
		t.Fatalf("expected non-null String object, got %T (%v)", res, res)
	}

	got := object.GoStringFromStringObject(strObj)
	expected := "Person[name=Alice, age=30, active=true]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}

	// Empty record
	pEmpty := object.MakeEmptyObjectWithClassName(&className)
	res = recordToString([]any{pEmpty})
	strObj = res.(*object.Object)
	got = object.GoStringFromStringObject(strObj)
	if got != "Person[]" {
		t.Errorf("expected %q, got %q", "Person[]", got)
	}

	// Null or invalid receiver
	errRes := recordToString([]any{object.Null})
	if _, ok := errRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk on null receiver, got %v", errRes)
	}
	errRes = recordToString([]any{})
	if _, ok := errRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected GErrBlk on empty params, got %v", errRes)
	}
}
