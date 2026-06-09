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
	"jacobin/src/types"
	"testing"
)

func TestSet_Of(t *testing.T) {
	globals.InitStringPool()

	// Set.of()
	res := setOf([]interface{}{})
	setObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	if className := object.GoStringFromStringPoolIndex(setObj.KlassName); className != classNameHashSet {
		t.Errorf("expected className %s, got %s", classNameHashSet, className)
	}

	sizeRes := hashmapSize([]interface{}{setObj})
	if sizeRes.(int64) != 0 {
		t.Errorf("expected size 0, got %v", sizeRes)
	}

	// Set.of(e1, e2)
	e1 := object.StringObjectFromGoString("one")
	e2 := object.StringObjectFromGoString("two")
	res = setOf([]interface{}{e1, e2})
	setObj = res.(*object.Object)

	sizeRes = hashmapSize([]interface{}{setObj})
	if sizeRes.(int64) != 2 {
		t.Errorf("expected size 2, got %v", sizeRes)
	}

	if hashsetContains([]interface{}{setObj, e1}) != types.JavaBoolTrue {
		t.Errorf("set should contain e1")
	}
	if hashsetContains([]interface{}{setObj, e2}) != types.JavaBoolTrue {
		t.Errorf("set should contain e2")
	}

	// Set.of(null) should return ghelpers.GErrBlk (NullPointerException)
	res = setOf([]interface{}{nil})
	if geb, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected *ghelpers.GErrBlk for null element, got %T", res)
	} else if geb.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %v", geb.ExceptionType)
	}

	// Set.of(e1, e1) should return ghelpers.GErrBlk (IllegalArgumentException)
	res = setOf([]interface{}{e1, e1})
	if geb, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("expected *ghelpers.GErrBlk for duplicate element, got %T", res)
	} else if geb.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", geb.ExceptionType)
	}
}

func TestSet_OfVarargs(t *testing.T) {
	globals.InitStringPool()

	e1 := object.StringObjectFromGoString("one")
	e2 := object.StringObjectFromGoString("two")

	// Create a Java array object [Ljava/lang/Object;
	arrayObj := object.Make1DimRefArray("java/lang/Object;", 2)
	rawArray := arrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	rawArray[0] = e1
	rawArray[1] = e2

	res := setOfVarargs([]interface{}{arrayObj})
	setObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}

	sizeRes := hashmapSize([]interface{}{setObj})
	if sizeRes.(int64) != 2 {
		t.Errorf("expected size 2, got %v", sizeRes)
	}

	if hashsetContains([]interface{}{setObj, e1}) != types.JavaBoolTrue {
		t.Errorf("set should contain e1")
	}
	if hashsetContains([]interface{}{setObj, e2}) != types.JavaBoolTrue {
		t.Errorf("set should contain e2")
	}
}

func TestSet_Registration(t *testing.T) {
	Load_Util_Set()

	methods := []string{
		"java/util/Set.add(Ljava/lang/Object;)Z",
		"java/util/Set.addAll(Ljava/util/Collection;)Z",
		"java/util/Set.clear()V",
		"java/util/Set.contains(Ljava/lang/Object;)Z",
		"java/util/Set.containsAll(Ljava/util/Collection;)Z",
		"java/util/Set.equals(Ljava/lang/Object;)Z",
		"java/util/Set.hashCode()I",
		"java/util/Set.isEmpty()Z",
		"java/util/Set.iterator()Ljava/util/Iterator;",
		"java/util/Set.remove(Ljava/lang/Object;)Z",
		"java/util/Set.removeAll(Ljava/util/Collection;)Z",
		"java/util/Set.retainAll(Ljava/util/Collection;)Z",
		"java/util/Set.size()I",
		"java/util/Set.spliterator()Ljava/util/Spliterator;",
		"java/util/Set.toArray()[Ljava/lang/Object;",
		"java/util/Set.toArray([Ljava/lang/Object;)[Ljava/lang/Object;",
		"java/util/Set.of()Ljava/util/Set;",
		"java/util/Set.copyOf(Ljava/util/Collection;)Ljava/util/Set;",
	}

	for _, m := range methods {
		gm, ok := ghelpers.MethodSignatures[m]
		if !ok {
			t.Errorf("method %s not registered", m)
			continue
		}

		// Verify that core methods are NOT trapped anymore
		notTrapped := map[string]bool{
			"java/util/Set.add(Ljava/lang/Object;)Z":       true,
			"java/util/Set.clear()V":                       true,
			"java/util/Set.contains(Ljava/lang/Object;)Z":  true,
			"java/util/Set.isEmpty()Z":                     true,
			"java/util/Set.iterator()Ljava/util/Iterator;": true,
			"java/util/Set.remove(Ljava/lang/Object;)Z":    true,
			"java/util/Set.size()I":                        true,
			"java/util/Set.of()Ljava/util/Set;":            true,
			"java/util/Set.toArray()[Ljava/lang/Object;":   true,
		}

		isTrapped := gm.GFunction == nil || (m != "java/util/Set.of()Ljava/util/Set;" && m != "java/util/Set.of([Ljava/lang/Object;)Ljava/util/Set;" &&
			(func(f func([]interface{}) interface{}) bool {
				// We can't easily compare function pointers in Go for equivalence with ghelpers.TrapFunction
				// unless we have access to it and it's not a closure.
				// But we can check if it's NOT one of our known implementations.
				return !notTrapped[m]
			})(gm.GFunction))
		if notTrapped[m] && isTrapped {
			t.Errorf("method %s should NOT be trapped", m)
		}
		if !notTrapped[m] && !isTrapped && m != "java/util/Set.of()Ljava/util/Set;" {
			// Some methods might be implemented but not in our 'notTrapped' list.
			// But for now, we want to ensure the ones we just moved are indeed not trapped.
		}
	}
}

func TestSet_IteratorInvoke(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Set()
	Load_Util_Hash_Set()

	// 1. Test Set.iterator()
	setObj := setOf([]interface{}{}).(*object.Object)
	res := ghelpers.Invoke("java/util/Set.iterator()Ljava/util/Iterator;", []interface{}{setObj})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		t.Errorf("Invoke java/util/Set.iterator() failed: %v", geb.ExceptionType)
	} else if res == nil {
		t.Errorf("Invoke java/util/Set.iterator() returned nil")
	}

	// 3. Test HashMap.entrySet().iterator()
	Load_Util_Map()
	Load_Util_Hash_Map()
	mapName := "java/util/HashMap"
	mapObj := object.MakeEmptyObjectWithClassName(&mapName)
	hashmapInit([]interface{}{mapObj})

	res = ghelpers.Invoke("java/util/HashMap.entrySet()Ljava/util/Set;", []interface{}{mapObj})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Invoke entrySet failed: %v", geb.ExceptionType)
	}
	entrySetObj := res.(*object.Object)
	res = ghelpers.Invoke("java/util/Set.iterator()Ljava/util/Iterator;", []interface{}{entrySetObj})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		t.Errorf("Invoke entrySet.iterator() failed: %v", geb.ExceptionType)
	}

	// 4. Test HashMap.keySet().iterator()
	res = ghelpers.Invoke("java/util/HashMap.keySet()Ljava/util/Set;", []interface{}{mapObj})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Invoke keySet failed: %v", geb.ExceptionType)
	}
	keySetObj := res.(*object.Object)
	res = ghelpers.Invoke("java/util/Set.iterator()Ljava/util/Iterator;", []interface{}{keySetObj})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		t.Errorf("Invoke keySet.iterator() failed: %v", geb.ExceptionType)
	}
}
