/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestIterator_ArrayList(t *testing.T) {
	globals.InitStringPool()

	// Create an ArrayList and add some elements
	al := object.MakeEmptyObject()
	al.KlassName = object.StringPoolIndexFromGoString("java/util/ArrayList")
	arraylistInit([]interface{}{al})

	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")
	s3 := object.StringObjectFromGoString("three")

	arraylistAdd([]interface{}{al, s1})
	arraylistAdd([]interface{}{al, s2})
	arraylistAdd([]interface{}{al, s3})

	// Get iterator
	iter := arraylistIterator([]interface{}{al}).(*object.Object)

	// Test hasNext and next
	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolTrue {
		t.Fatal("expected hasNext to be true")
	}

	res := iteratorNext([]interface{}{iter}).(*object.Object)
	if object.GoStringFromStringObject(res) != "one" {
		t.Fatalf("expected 'one', got %s", object.GoStringFromStringObject(res))
	}

	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolTrue {
		t.Fatal("expected hasNext to be true")
	}

	res = iteratorNext([]interface{}{iter}).(*object.Object)
	if object.GoStringFromStringObject(res) != "two" {
		t.Fatalf("expected 'two', got %s", object.GoStringFromStringObject(res))
	}

	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolTrue {
		t.Fatal("expected hasNext to be true")
	}

	res = iteratorNext([]interface{}{iter}).(*object.Object)
	if object.GoStringFromStringObject(res) != "three" {
		t.Fatalf("expected 'three', got %s", object.GoStringFromStringObject(res))
	}

	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolFalse {
		t.Fatal("expected hasNext to be false")
	}
}

func TestIterator_Remove_ArrayList(t *testing.T) {
	globals.InitStringPool()

	// Create an ArrayList and add some elements
	al := object.MakeEmptyObject()
	al.KlassName = object.StringPoolIndexFromGoString("java/util/ArrayList")
	arraylistInit([]interface{}{al})

	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")
	s3 := object.StringObjectFromGoString("three")

	arraylistAdd([]interface{}{al, s1})
	arraylistAdd([]interface{}{al, s2})
	arraylistAdd([]interface{}{al, s3})

	// Get iterator
	iter := arraylistIterator([]interface{}{al}).(*object.Object)

	// Test IllegalStateException before next()
	res := iteratorRemove([]interface{}{iter})
	if _, ok := res.(*ghelpers.GErrBlk); !ok || res.(*ghelpers.GErrBlk).ExceptionType != excNames.IllegalStateException {
		t.Fatal("expected IllegalStateException when calling remove() before next()")
	}

	// next() then remove() s1
	iteratorNext([]interface{}{iter})
	iteratorRemove([]interface{}{iter})

	// Test IllegalStateException calling remove() twice
	res = iteratorRemove([]interface{}{iter})
	if _, ok := res.(*ghelpers.GErrBlk); !ok || res.(*ghelpers.GErrBlk).ExceptionType != excNames.IllegalStateException {
		t.Fatal("expected IllegalStateException when calling remove() twice")
	}

	// Verify s1 is removed
	sz := arraylistSize([]interface{}{al}).(int64)
	if sz != 2 {
		t.Fatalf("expected size 2, got %d", sz)
	}
	if object.GoStringFromStringObject(arraylistGet([]interface{}{al, int64(0)}).(*object.Object)) != "two" {
		t.Fatal("expected first element to be 'two'")
	}

	// next() then remove() s2
	iteratorNext([]interface{}{iter})
	iteratorRemove([]interface{}{iter})

	// Verify s2 is removed
	sz = arraylistSize([]interface{}{al}).(int64)
	if sz != 1 {
		t.Fatalf("expected size 1, got %d", sz)
	}
	if object.GoStringFromStringObject(arraylistGet([]interface{}{al, int64(0)}).(*object.Object)) != "three" {
		t.Fatal("expected first element to be 'three'")
	}

	// next() then remove() s3
	iteratorNext([]interface{}{iter})
	iteratorRemove([]interface{}{iter})

	// Verify s3 is removed
	sz = arraylistSize([]interface{}{al}).(int64)
	if sz != 0 {
		t.Fatalf("expected size 0, got %d", sz)
	}
}

func TestIterator_Remove_LinkedList(t *testing.T) {
	globals.InitStringPool()

	// Create a LinkedList and add some elements
	ll := object.MakeEmptyObject()
	ll.KlassName = object.StringPoolIndexFromGoString("java/util/LinkedList")
	linkedlistInit([]interface{}{ll})

	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")

	linkedlistAddLast([]interface{}{ll, s1})
	linkedlistAddLast([]interface{}{ll, s2})

	// Get iterator
	iter := linkedlistIterator([]interface{}{ll}).(*object.Object)

	// Test IllegalStateException before next()
	res := iteratorRemove([]interface{}{iter})
	if _, ok := res.(*ghelpers.GErrBlk); !ok || res.(*ghelpers.GErrBlk).ExceptionType != excNames.IllegalStateException {
		t.Fatal("expected IllegalStateException when calling remove() before next()")
	}

	// next() then remove() s1
	iteratorNext([]interface{}{iter})
	iteratorRemove([]interface{}{iter})

	// Verify s1 is removed
	llst, _ := ghelpers.GetLinkedListFromObject(ll)
	if llst.Len() != 1 {
		t.Fatalf("expected length 1, got %d", llst.Len())
	}
	if object.GoStringFromStringObject(llst.Front().Value.(*object.Object)) != "two" {
		t.Fatal("expected first element to be 'two'")
	}

	// next() then remove() s2
	iteratorNext([]interface{}{iter})
	iteratorRemove([]interface{}{iter})

	// Verify s2 is removed
	if llst.Len() != 0 {
		t.Fatalf("expected length 0, got %d", llst.Len())
	}
}

func TestIterator_Vector(t *testing.T) {
	globals.InitStringPool()

	// Create a Vector and add some elements
	v := object.MakeEmptyObject()
	v.KlassName = object.StringPoolIndexFromGoString("java/util/Vector")
	vectorInit([]interface{}{v})

	s1 := object.StringObjectFromGoString("one")
	s2 := object.StringObjectFromGoString("two")
	s3 := object.StringObjectFromGoString("three")

	vectorAdd([]interface{}{v, s1})
	vectorAdd([]interface{}{v, s2})
	vectorAdd([]interface{}{v, s3})

	// Get iterator
	iter := vectorIterator([]interface{}{v}).(*object.Object)

	// Test hasNext and next
	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolTrue {
		t.Fatal("expected hasNext to be true")
	}

	res := iteratorNext([]interface{}{iter}).(*object.Object)
	if object.GoStringFromStringObject(res) != "one" {
		t.Fatalf("expected 'one', got %s", object.GoStringFromStringObject(res))
	}

	// Test remove s1
	iteratorRemove([]interface{}{iter})

	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolTrue {
		t.Fatal("expected hasNext to be true")
	}

	res = iteratorNext([]interface{}{iter}).(*object.Object)
	if object.GoStringFromStringObject(res) != "two" {
		t.Fatalf("expected 'two', got %s", object.GoStringFromStringObject(res))
	}

	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolTrue {
		t.Fatal("expected hasNext to be true")
	}

	res = iteratorNext([]interface{}{iter}).(*object.Object)
	if object.GoStringFromStringObject(res) != "three" {
		t.Fatalf("expected 'three', got %s", object.GoStringFromStringObject(res))
	}

	if iteratorHasNext([]interface{}{iter}) != types.JavaBoolFalse {
		t.Fatal("expected hasNext to be false")
	}

	// Verify s1 was removed and others remain
	sz := vectorSize([]interface{}{v}).(int64)
	if sz != 2 {
		t.Fatalf("expected size 2, got %d", sz)
	}
	if object.GoStringFromStringObject(vectorGet([]interface{}{v, int64(0)}).(*object.Object)) != "two" {
		t.Fatal("expected first element to be 'two'")
	}
}

func TestIterator_Unsupported(t *testing.T) {
	globals.InitStringPool()

	// Create an unsupported collection object
	col := object.MakeEmptyObject()
	col.KlassName = object.StringPoolIndexFromGoString("java/util/HashSet") // Assume HashSet is unsupported by Iterator GFunction for now

	// Get iterator
	iter := NewIterator(col)

	// Test hasNext
	res := iteratorHasNext([]interface{}{iter})
	if _, ok := res.(*ghelpers.GErrBlk); !ok || res.(*ghelpers.GErrBlk).ExceptionType != excNames.UnsupportedOperationException {
		t.Fatal("expected UnsupportedOperationException for hasNext on unsupported collection")
	}

	// Test next
	res = iteratorNext([]interface{}{iter})
	if _, ok := res.(*ghelpers.GErrBlk); !ok || res.(*ghelpers.GErrBlk).ExceptionType != excNames.UnsupportedOperationException {
		t.Fatal("expected UnsupportedOperationException for next on unsupported collection")
	}

	// Test remove
	res = iteratorRemove([]interface{}{iter})
	if _, ok := res.(*ghelpers.GErrBlk); !ok || res.(*ghelpers.GErrBlk).ExceptionType != excNames.UnsupportedOperationException {
		t.Fatal("expected UnsupportedOperationException for remove on unsupported collection")
	}
}
