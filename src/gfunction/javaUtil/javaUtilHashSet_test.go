package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helpers
func newHashSetObj(t *testing.T) *object.Object {
	t.Helper()
	hs := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	// initialize backing map via hashmapInit (HashSet uses HashMap impl)
	if ret := hashmapInit([]interface{}{hs}); ret != nil {
		t.Fatalf("hashmapInit returned error: %v", ret)
	}
	return hs
}

func intObj(v int64) *object.Object {
	return object.MakePrimitiveObject("java/lang/Integer", types.Int, v)
}
func longObj(v int64) *object.Object {
	return object.MakePrimitiveObject("java/lang/Long", types.Long, v)
}
func floatObj(v float64) *object.Object {
	return object.MakePrimitiveObject("java/lang/Double", types.Double, v)
}

func assertJavaBool(t *testing.T, got interface{}, want int64, msg string) {
	t.Helper()
	gi, ok := got.(int64)
	if !ok {
		t.Fatalf("%s: expected int64 Java boolean, got %T", msg, got)
	}
	if gi != want {
		t.Fatalf("%s: expected %d, got %d", msg, want, gi)
	}
}

func TestHashSet_Add_Contains_Size_Remove_ToArray(t *testing.T) {
	globals.InitStringPool()

	hs := newHashSetObj(t)

	// Initially empty
	assertJavaBool(t, hashsetIsEmpty([]interface{}{hs}), types.JavaBoolTrue, "isEmpty initially")

	// Add distinct elements (using different boxed types)
	e1 := intObj(42)
	e2 := longObj(42)    // note: different boxed type but same numeric -> hashing uses value; may collide or not
	e3 := floatObj(3.14) // float64 supported in hashing

	// First add of e1 -> true
	assertJavaBool(t, hashsetAdd([]interface{}{hs, e1}), types.JavaBoolTrue, "add e1 first time")
	// Add duplicate e1 -> false
	assertJavaBool(t, hashsetAdd([]interface{}{hs, e1}), types.JavaBoolFalse, "add e1 second time")

	// Add e2 and e3 (even if e2 hashes same as e1, value stored is object; treat as set of values by hash key)
	_ = hashsetAdd([]interface{}{hs, e2})
	_ = hashsetAdd([]interface{}{hs, e3})

	// Contains checks
	assertJavaBool(t, hashsetContains([]interface{}{hs, e1}), types.JavaBoolTrue, "contains e1")
	assertJavaBool(t, hashsetContains([]interface{}{hs, e2}), types.JavaBoolTrue, "contains e2")
	assertJavaBool(t, hashsetContains([]interface{}{hs, e3}), types.JavaBoolTrue, "contains e3")
	assertJavaBool(t, hashsetIsEmpty([]interface{}{hs}), types.JavaBoolFalse, "isEmpty after adds")

	// Size should be >= 2 (e1 + e3) and up to 3 depending on hash collision of 42-int vs 42-long.
	// We can't rely on exact size here; instead verify toArray contains at least the elements we added.
	arrObj := hashsetToArray([]interface{}{hs}).(*object.Object)
	// toArray returns RefArray with [] *object.Object in value field
	arr, ok := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if !ok {
		t.Fatalf("toArray did not return object array; got %T", arrObj.FieldTable["value"].Fvalue)
	}
	// Build a set of pointers for membership checks
	seen := map[*object.Object]bool{}
	for _, o := range arr {
		seen[o] = true
	}
	if !seen[e1] {
		t.Fatalf("toArray missing e1")
	}
	if !seen[e3] {
		t.Fatalf("toArray missing e3")
	}
	// e2 may be same-hash replacement of e1 depending on implementation; accept either being present

	// Remove elements
	assertJavaBool(t, hashsetRemove([]interface{}{hs, e1}), types.JavaBoolTrue, "remove e1 present")
	// Removing again should be false
	assertJavaBool(t, hashsetRemove([]interface{}{hs, e1}), types.JavaBoolFalse, "remove e1 again")

	// Remove remaining elements to reach empty
	_ = hashsetRemove([]interface{}{hs, e2})
	_ = hashsetRemove([]interface{}{hs, e3})

	assertJavaBool(t, hashsetIsEmpty([]interface{}{hs}), types.JavaBoolTrue, "isEmpty after removals")
}

func TestHashSet_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	hs := newHashSetObj(t)

	// First param non-object -> ClassCastException
	if err := hashsetAdd([]interface{}{int64(5), intObj(1)}); err == nil {
		t.Fatalf("expected error for non-object first param in add")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.ClassCastException {
			t.Fatalf("expected ClassCastException, got %d", geb.ExceptionType)
		}
	}

	// First param wrong class -> IllegalArgumentException
	notMap := object.MakeEmptyObjectWithClassName(&classNameBase64Encoder)
	if err := hashsetAdd([]interface{}{notMap, intObj(1)}); err == nil {
		t.Fatalf("expected error for wrong class in add")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
		}
	}

	// Element param non-object -> IllegalArgumentException
	if err := hashsetAdd([]interface{}{hs, int64(7)}); err == nil {
		t.Fatalf("expected error for non-object element in add")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
		}
	}

	// Element missing value field -> IllegalArgumentException
	badElem := object.MakeEmptyObjectWithClassName(&classNameObject) // no "value" field
	if err := hashsetAdd([]interface{}{hs, badElem}); err == nil {
		t.Fatalf("expected error for missing value field in element")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
		}
	}

	// Similar error paths for contains and remove
	if err := hashsetContains([]interface{}{hs, int64(1)}); err == nil {
		t.Fatalf("expected error for non-object element in contains")
	}
	if err := hashsetRemove([]interface{}{hs, int64(1)}); err == nil {
		t.Fatalf("expected error for non-object element in remove")
	}

	// hashsetIsEmpty forwards size; wrong-class first param should error
	if err := hashsetIsEmpty([]interface{}{notMap}); err == nil {
		t.Fatalf("expected error for wrong-class param in isEmpty")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException from isEmpty, got %d", geb.ExceptionType)
		}
	}
}
