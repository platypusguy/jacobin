/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func TestCollection_Registration(t *testing.T) {
	Load_Util_Collection()

	methods := []string{
		"java/util/Collection.add(Ljava/lang/Object;)Z",
		"java/util/Collection.addAll(Ljava/util/Collection;)Z",
		"java/util/Collection.clear()V",
		"java/util/Collection.contains(Ljava/lang/Object;)Z",
		"java/util/Collection.containsAll(Ljava/util/Collection;)Z",
		"java/util/Collection.equals(Ljava/lang/Object;)Z",
		"java/util/Collection.hashCode()I",
		"java/util/Collection.isEmpty()Z",
		"java/util/Collection.iterator()Ljava/util/Iterator;",
		"java/util/Collection.remove(Ljava/lang/Object;)Z",
		"java/util/Collection.removeAll(Ljava/util/Collection;)Z",
		"java/util/Collection.retainAll(Ljava/util/Collection;)Z",
		"java/util/Collection.size()I",
		"java/util/Collection.spliterator()Ljava/util/Spliterator;",
		"java/util/Collection.toArray()[Ljava/lang/Object;",
		"java/util/Collection.toArray([Ljava/lang/Object;)[Ljava/lang/Object;",
		"java/util/Collection.parallelStream()Ljava/util/stream/Stream;",
		"java/util/Collection.removeIf(Ljava/util/function/Predicate;)Z",
		"java/util/Collection.stream()Ljava/util/stream/Stream;",
		"java/util/Collection.forEach(Ljava/util/function/Consumer;)V",
	}

	for _, m := range methods {
		_, ok := ghelpers.MethodSignatures[m]
		if !ok {
			t.Errorf("method %s not registered", m)
		}
	}
}

func TestCollection_Invoke(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Collection()
	Load_Util_Hash_Set()
	Load_Util_Set()

	// Use setOf to create a collection (which is also a Set)
	e1 := object.StringObjectFromGoString("one")
	res := setOf([]interface{}{e1})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("setOf failed: %v", geb.ErrMsg)
	}
	collObj := res.(*object.Object)

	// Test size()
	sizeRes := ghelpers.Invoke("java/util/Collection.size()I", []interface{}{collObj})
	if sizeRes.(int64) != 1 {
		t.Errorf("expected size 1, got %v", sizeRes)
	}

	// Test isEmpty()
	emptyRes := ghelpers.Invoke("java/util/Collection.isEmpty()Z", []interface{}{collObj})
	if emptyRes.(int64) != 0 { // In Jacobin, boolean often returns as int64 (0 or 1)
		t.Errorf("expected isEmpty false, got %v", emptyRes)
	}

	// Test contains()
	contRes := ghelpers.Invoke("java/util/Collection.contains(Ljava/lang/Object;)Z", []interface{}{collObj, e1})
	if contRes.(int64) != 1 {
		t.Errorf("expected contains true, got %v", contRes)
	}

	// Test iterator()
	iterRes := ghelpers.Invoke("java/util/Collection.iterator()Ljava/util/Iterator;", []interface{}{collObj})
	if geb, ok := iterRes.(*ghelpers.GErrBlk); ok {
		t.Errorf("Invoke iterator() failed: %v", geb.ExceptionType)
	} else if iterRes == nil {
		t.Errorf("Invoke iterator() returned nil")
	}

	// Test toArray() (should be trapped)
	arrRes := ghelpers.Invoke("java/util/Collection.toArray()[Ljava/lang/Object;", []interface{}{collObj})
	if geb, ok := arrRes.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected trap for toArray, got %T", arrRes)
	} else if geb.ExceptionType != excNames.UnsupportedOperationException {
		t.Errorf("expected UnsupportedOperationException, got %v", geb.ExceptionType)
	}
}
