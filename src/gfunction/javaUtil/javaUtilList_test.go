/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestList_Of(t *testing.T) {
	globals.InitStringPool()

	// List.of()
	res := listOf([]interface{}{})
	listObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	if className := object.GoStringFromStringPoolIndex(listObj.KlassName); className != "java/util/ArrayList" {
		t.Errorf("expected className java/util/ArrayList, got %s", className)
	}
	al, err := GetArrayListFromObject(listObj)
	if err != nil {
		t.Fatalf("getArrayListFromObject failed: %v", err)
	}
	if len(al) != 0 {
		t.Errorf("expected size 0, got %d", len(al))
	}

	// List.of(e1, e2)
	e1 := object.StringObjectFromGoString("one")
	e2 := object.StringObjectFromGoString("two")
	res = listOf([]interface{}{e1, e2})
	listObj = res.(*object.Object)
	al, _ = GetArrayListFromObject(listObj)
	if len(al) != 2 {
		t.Errorf("expected size 2, got %d", len(al))
	}
	if al[0] != e1 || al[1] != e2 {
		t.Errorf("elements mismatch")
	}

	// List.of(null) should return ghelpers.GErrBlk (NullPointerException)
	res = listOf([]interface{}{nil})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected *ghelpers.GErrBlk for null element, got %T", res)
	}
}

func TestList_OfVarargs(t *testing.T) {
	globals.InitStringPool()

	e1 := object.StringObjectFromGoString("one")
	e2 := object.StringObjectFromGoString("two")

	// Create a Java array object [Ljava/lang/Object;
	arrayObj := object.Make1DimRefArray("java/lang/Object;", 2)
	rawArray := arrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	rawArray[0] = e1
	rawArray[1] = e2

	res := listOfVarargs([]interface{}{arrayObj})
	listObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	al, _ := GetArrayListFromObject(listObj)
	if len(al) != 2 {
		t.Errorf("expected size 2, got %d", len(al))
	}
	if al[0] != e1 || al[1] != e2 {
		t.Errorf("elements mismatch")
	}
}

func TestList_Iterator(t *testing.T) {
	globals.InitStringPool()

	e1 := object.StringObjectFromGoString("one")
	e2 := object.StringObjectFromGoString("two")
	listObj := listOf([]interface{}{e1, e2}).(*object.Object)

	// Get iterator
	res := listIterator([]interface{}{listObj})
	iterObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object from listIterator, got %T", res)
	}

	// Test iterator
	if iteratorHasNext([]interface{}{iterObj}) != types.JavaBoolTrue {
		t.Error("expected hasNext to be true")
	}
	if next := iteratorNext([]interface{}{iterObj}); next != e1 {
		t.Errorf("expected 'one', got %v", next)
	}
	if iteratorHasNext([]interface{}{iterObj}) != types.JavaBoolTrue {
		t.Error("expected hasNext to be true")
	}
	if next := iteratorNext([]interface{}{iterObj}); next != e2 {
		t.Errorf("expected 'two', got %v", next)
	}
	if iteratorHasNext([]interface{}{iterObj}) != types.JavaBoolFalse {
		t.Error("expected hasNext to be false")
	}
}

func TestList_ListIterator(t *testing.T) {
	globals.InitStringPool()

	e1 := object.StringObjectFromGoString("one")
	e2 := object.StringObjectFromGoString("two")
	listObj := listOf([]interface{}{e1, e2}).(*object.Object)

	// List.listIterator()
	res := listListIterator([]interface{}{listObj})
	iterObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object from listListIterator, got %T", res)
	}

	if listiteratorHasNext([]interface{}{iterObj}) != types.JavaBoolTrue {
		t.Error("expected hasNext to be true")
	}
	if next := listiteratorNext([]interface{}{iterObj}); next != e1 {
		t.Errorf("expected 'one', got %v", next)
	}

	// List.listIterator(1)
	res = listListIteratorWithIndex([]interface{}{listObj, int64(1)})
	iterObj = res.(*object.Object)
	if next := listiteratorNext([]interface{}{iterObj}); next != e2 {
		t.Errorf("expected 'two', got %v", next)
	}
}
