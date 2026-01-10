/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestArrayList_Basic(t *testing.T) {
	globals.InitStringPool()

	// Initialize ArrayList
	alObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	arraylistInit([]interface{}{alObj})

	// add(Object)
	s1 := object.StringObjectFromGoString("one")
	res := arraylistAdd([]interface{}{alObj, s1})
	if res != types.JavaBoolTrue {
		t.Errorf("add expected true, got %v", res)
	}

	// size()
	sz := arraylistSize([]interface{}{alObj})
	if sz.(int64) != 1 {
		t.Errorf("size expected 1, got %v", sz)
	}

	// add(int, Object)
	s0 := object.StringObjectFromGoString("zero")
	arraylistAddAtIndex([]interface{}{alObj, int64(0), s0})

	sz = arraylistSize([]interface{}{alObj})
	if sz.(int64) != 2 {
		t.Errorf("size expected 2, got %v", sz)
	}

	// get(int)
	got := arraylistGet([]interface{}{alObj, int64(0)})
	if !object.EqualStringObjects(got.(*object.Object), s0) {
		t.Errorf("get(0) mismatch")
	}

	got = arraylistGet([]interface{}{alObj, int64(1)})
	if !object.EqualStringObjects(got.(*object.Object), s1) {
		t.Errorf("get(1) mismatch")
	}

	// set(int, Object)
	s1new := object.StringObjectFromGoString("one-new")
	old := arraylistSet([]interface{}{alObj, int64(1), s1new})
	if !object.EqualStringObjects(old.(*object.Object), s1) {
		t.Errorf("set(1) old value mismatch")
	}

	got = arraylistGet([]interface{}{alObj, int64(1)})
	if !object.EqualStringObjects(got.(*object.Object), s1new) {
		t.Errorf("get(1) after set mismatch")
	}

	// indexOf(Object)
	idx := arraylistIndexOf([]interface{}{alObj, s1new})
	if idx.(int64) != 1 {
		t.Errorf("indexOf mismatch: expected 1, got %v", idx)
	}

	// contains(Object)
	cont := arraylistContains([]interface{}{alObj, s0})
	if cont != types.JavaBoolTrue {
		t.Errorf("contains(s0) expected true")
	}

	// remove(int)
	rem := arraylistRemoveAtIndex([]interface{}{alObj, int64(0)})
	if !object.EqualStringObjects(rem.(*object.Object), s0) {
		t.Errorf("remove(0) mismatch")
	}

	sz = arraylistSize([]interface{}{alObj})
	if sz.(int64) != 1 {
		t.Errorf("size after remove expected 1, got %v", sz)
	}

	// isEmpty()
	empty := arraylistIsEmpty([]interface{}{alObj})
	if empty != types.JavaBoolFalse {
		t.Errorf("isEmpty expected false")
	}

	// clear()
	arraylistClear([]interface{}{alObj})
	sz = arraylistSize([]interface{}{alObj})
	if sz.(int64) != 0 {
		t.Errorf("size after clear expected 0, got %v", sz)
	}

	empty = arraylistIsEmpty([]interface{}{alObj})
	if empty != types.JavaBoolTrue {
		t.Errorf("isEmpty after clear expected true")
	}
}

func TestArrayList_RemoveObject(t *testing.T) {
	globals.InitStringPool()

	alObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	arraylistInit([]interface{}{alObj})

	s1 := object.StringObjectFromGoString("apple")
	s2 := object.StringObjectFromGoString("banana")
	arraylistAdd([]interface{}{alObj, s1})
	arraylistAdd([]interface{}{alObj, s2})

	// remove(Object)
	res := arraylistRemoveObject([]interface{}{alObj, s1})
	if res != types.JavaBoolTrue {
		t.Errorf("remove(Object) expected true")
	}

	sz := arraylistSize([]interface{}{alObj})
	if sz.(int64) != 1 {
		t.Errorf("size expected 1, got %v", sz)
	}

	got := arraylistGet([]interface{}{alObj, int64(0)})
	if !object.EqualStringObjects(got.(*object.Object), s2) {
		t.Errorf("remaining element mismatch")
	}
}

func TestArrayList_InitWithCapacity(t *testing.T) {
	globals.InitStringPool()

	alObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	arraylistInitWithCapacity([]interface{}{alObj, int64(50)})

	sz := arraylistSize([]interface{}{alObj})
	if sz.(int64) != 0 {
		t.Errorf("size expected 0")
	}
}

func TestArrayList_ToArray(t *testing.T) {
	globals.InitStringPool()

	alObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	arraylistInit([]interface{}{alObj})

	s1 := object.StringObjectFromGoString("a")
	s2 := object.StringObjectFromGoString("b")
	arraylistAdd([]interface{}{alObj, s1})
	arraylistAdd([]interface{}{alObj, s2})

	res := arraylistToArray([]interface{}{alObj})
	// In Jacobin, Populator for "[Ljava/lang/Object;" returns an *object.Object wrapping the array
	arrObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("toArray did not return *object.Object")
	}

	arr, ok := arrObj.FieldTable["value"].Fvalue.([]interface{})
	if !ok {
		t.Fatalf("toArray value is not []interface{}")
	}

	if len(arr) != 2 {
		t.Errorf("array length mismatch: %d", len(arr))
	}
}

func TestArrayList_Clone(t *testing.T) {
	globals.InitStringPool()

	alObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	arraylistInit([]interface{}{alObj})

	s1 := object.StringObjectFromGoString("clone-me")
	arraylistAdd([]interface{}{alObj, s1})

	cloneRes := arraylistClone([]interface{}{alObj})
	cloneObj := cloneRes.(*object.Object)

	sz := arraylistSize([]interface{}{cloneObj})
	if sz.(int64) != 1 {
		t.Errorf("clone size mismatch")
	}

	got := arraylistGet([]interface{}{cloneObj, int64(0)})
	if !object.EqualStringObjects(got.(*object.Object), s1) {
		t.Errorf("clone content mismatch")
	}

	// Modify original, clone should not change
	s2 := object.StringObjectFromGoString("extra")
	arraylistAdd([]interface{}{alObj, s2})

	sz = arraylistSize([]interface{}{alObj})
	if sz.(int64) != 2 {
		t.Errorf("original size mismatch")
	}

	sz = arraylistSize([]interface{}{cloneObj})
	if sz.(int64) != 1 {
		t.Errorf("clone size changed after original modification")
	}
}

func TestArrayList_Iterator(t *testing.T) {
	globals.InitStringPool()

	alObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	arraylistInit([]interface{}{alObj})

	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")
	arraylistAdd([]interface{}{alObj, s1})
	arraylistAdd([]interface{}{alObj, s2})

	// Get iterator
	res := arraylistIterator([]interface{}{alObj})
	iterObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object from arraylistIterator, got %T", res)
	}

	// Test iterator
	if iteratorHasNext([]interface{}{iterObj}) != types.JavaBoolTrue {
		t.Error("expected hasNext to be true")
	}
	if next := iteratorNext([]interface{}{iterObj}); next != s1 {
		t.Errorf("expected 'one', got %v", next)
	}
	if iteratorHasNext([]interface{}{iterObj}) != types.JavaBoolTrue {
		t.Error("expected hasNext to be true")
	}
	if next := iteratorNext([]interface{}{iterObj}); next != s2 {
		t.Errorf("expected 'two', got %v", next)
	}
	if iteratorHasNext([]interface{}{iterObj}) != types.JavaBoolFalse {
		t.Error("expected hasNext to be false")
	}
}
