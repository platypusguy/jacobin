/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestVector_Basic(t *testing.T) {
	Load_Util_Vector()
	Load_Util_Iterator()
	Load_Util_ListIterator()

	vObj := object.MakeEmptyObject()
	vObj.KlassName = 123 // dummy
	params := []interface{}{vObj}
	vectorInit(params)

	// add
	obj1 := object.StringObjectFromGoString("one")
	params = []interface{}{vObj, obj1}
	res := vectorAdd(params)
	if res != types.JavaBoolTrue {
		t.Errorf("vectorAdd failed")
	}

	// size
	res = vectorSize([]interface{}{vObj})
	if res.(int64) != 1 {
		t.Errorf("expected size 1, got %d", res)
	}

	// get
	res = vectorGet([]interface{}{vObj, int64(0)})
	if res != obj1 {
		t.Errorf("get(0) failed")
	}

	// add element
	obj2 := object.StringObjectFromGoString("two")
	vectorAddElement([]interface{}{vObj, obj2})

	if vectorSize([]interface{}{vObj}).(int64) != 2 {
		t.Errorf("expected size 2")
	}

	// contains
	if vectorContains([]interface{}{vObj, obj1}) != types.JavaBoolTrue {
		t.Errorf("should contain obj1")
	}

	// indexOf
	if vectorIndexOf([]interface{}{vObj, obj2}).(int64) != 1 {
		t.Errorf("indexOf obj2 should be 1")
	}

	// remove
	res = vectorRemoveAtIndex([]interface{}{vObj, int64(0)})
	if res != obj1 {
		t.Errorf("remove(0) should return obj1")
	}
	if vectorSize([]interface{}{vObj}).(int64) != 1 {
		t.Errorf("size should be 1 after remove")
	}

	// clear
	vectorClear([]interface{}{vObj})
	if vectorSize([]interface{}{vObj}).(int64) != 0 {
		t.Errorf("size should be 0 after clear")
	}
}

func TestVector_Capacity(t *testing.T) {
	Load_Util_Vector()
	vObj := object.MakeEmptyObject()
	vectorInitWithCapacity([]interface{}{vObj, int64(20)})
	cap := vectorCapacity([]interface{}{vObj})
	if cap.(int64) != 20 {
		t.Errorf("expected capacity 20, got %d", cap)
	}
}

func TestVector_ListIterator(t *testing.T) {
	Load_Util_Vector()
	Load_Util_ListIterator()

	vObj := object.MakePrimitiveObject("java/util/Vector", types.Vector, []interface{}{})
	vectorAdd([]interface{}{vObj, object.StringObjectFromGoString("a")})
	vectorAdd([]interface{}{vObj, object.StringObjectFromGoString("b")})
	vectorAdd([]interface{}{vObj, object.StringObjectFromGoString("c")})

	li := vectorListIterator([]interface{}{vObj}).(*object.Object)

	// hasNext, next
	if listiteratorHasNext([]interface{}{li}) != types.JavaBoolTrue {
		t.Errorf("li should have next")
	}
	res := listiteratorNext([]interface{}{li})
	if object.GoStringFromStringObject(res.(*object.Object)) != "a" {
		t.Errorf("expected a, got %v", res)
	}

	// hasPrevious, previous
	if listiteratorHasPrevious([]interface{}{li}) != types.JavaBoolTrue {
		t.Errorf("li should have previous")
	}
	res = listiteratorPrevious([]interface{}{li})
	if object.GoStringFromStringObject(res.(*object.Object)) != "a" {
		t.Errorf("expected a, got %v", res)
	}

	// set
	listiteratorSet([]interface{}{li, object.StringObjectFromGoString("A")})
	res = vectorGet([]interface{}{vObj, int64(0)})
	if object.GoStringFromStringObject(res.(*object.Object)) != "A" {
		t.Errorf("expected A after set")
	}

	// add
	listiteratorAdd([]interface{}{li, object.StringObjectFromGoString("new")})
	// list is [new, A, b, c], cursor is at index 1 (pointing to A)
	if vectorSize([]interface{}{vObj}).(int64) != 4 {
		t.Errorf("expected size 4 after add")
	}
	res = vectorGet([]interface{}{vObj, int64(0)})
	if object.GoStringFromStringObject(res.(*object.Object)) != "new" {
		t.Errorf("expected 'new' at index 0")
	}

	// remove
	listiteratorNext([]interface{}{li})   // return A, cursor 2, lastReturned 1
	listiteratorRemove([]interface{}{li}) // remove A, cursor 1, lastReturned -1
	if vectorSize([]interface{}{vObj}).(int64) != 3 {
		t.Errorf("expected size 3 after remove")
	}
}

func TestVector_ToString(t *testing.T) {
	Load_Util_Vector()
	vObj := object.MakeEmptyObject()
	vObj.KlassName = 123 // dummy
	vectorInit([]interface{}{vObj})

	// Empty vector
	res := vectorToString([]interface{}{vObj})
	str := object.GoStringFromStringObject(res.(*object.Object))
	if str != "[]" {
		t.Errorf("expected [], got %s", str)
	}

	// Vector with elements
	vectorAdd([]interface{}{vObj, object.StringObjectFromGoString("one")})
	vectorAdd([]interface{}{vObj, object.StringObjectFromGoString("two")})

	res = vectorToString([]interface{}{vObj})
	str = object.GoStringFromStringObject(res.(*object.Object))
	if str != "[one, two]" {
		t.Errorf("expected [one, two], got %s", str)
	}
}
