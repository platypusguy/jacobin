package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestLinkedList_RemoveFirst_And_RemoveLast(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)
	_ = linkedlistAddLast([]interface{}{ll, strObj("a")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("b")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("c")})

	// removeFirst -> "a"
	vFirst := linkedlistRemoveFirst([]interface{}{ll}).(*object.Object)
	if object.GoStringFromStringObject(vFirst) != "a" {
		t.Fatalf("removeFirst mismatch: expected 'a', got %q", object.GoStringFromStringObject(vFirst))
	}

	// removeLast -> "c"
	vLast := linkedlistRemoveLast([]interface{}{ll}).(*object.Object)
	if object.GoStringFromStringObject(vLast) != "c" {
		t.Fatalf("removeLast mismatch: expected 'c', got %q", object.GoStringFromStringObject(vLast))
	}

	// Now list contains [b]; removeFirst on empty list should throw NoSuchElementException
	ll2 := newLinkedListObj(t)
	if err := linkedlistRemoveFirst([]interface{}{ll2}); err == nil {
		t.Fatalf("expected error for removeFirst on empty list")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.NoSuchElementException {
			t.Fatalf("expected NoSuchElementException, got %d", geb.ExceptionType)
		}
	}
	if err := linkedlistRemoveLast([]interface{}{ll2}); err == nil {
		t.Fatalf("expected error for removeLast on empty list")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.NoSuchElementException {
			t.Fatalf("expected NoSuchElementException, got %d", geb.ExceptionType)
		}
	}
}

func TestLinkedList_RemoveFirstOccurrence_And_LastOccurrence(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)
	// [a, b, a, c]
	a1 := strObj("a")
	b := strObj("b")
	a2 := strObj("a")
	c := strObj("c")
	_ = linkedlistAddLast([]interface{}{ll, a1})
	_ = linkedlistAddLast([]interface{}{ll, b})
	_ = linkedlistAddLast([]interface{}{ll, a2})
	_ = linkedlistAddLast([]interface{}{ll, c})

	// removeFirstOccurrence("a") => true; list becomes [b, a, c]
	assertJavaBoolLL(t, linkedlistRemoveFirstOccurrence([]interface{}{ll, strObj("a")}), types.JavaBoolTrue, "removeFirstOccurrence should return true")
	if s := object.GoStringFromStringObject(linkedlistGetFirst([]interface{}{ll}).(*object.Object)); s != "b" {
		t.Fatalf("after removeFirstOccurrence, expected head 'b', got %q", s)
	}

	// removeLastOccurrence("a") => true; list becomes [b, c]
	assertJavaBoolLL(t, linkedlistRemoveLastOccurrence([]interface{}{ll, strObj("a")}), types.JavaBoolTrue, "removeLastOccurrence should return true")
	if s := object.GoStringFromStringObject(linkedlistGetLast([]interface{}{ll}).(*object.Object)); s != "c" {
		t.Fatalf("after removeLastOccurrence, expected tail 'c', got %q", s)
	}

	// removing non-existent -> false
	assertJavaBoolLL(t, linkedlistRemoveFirstOccurrence([]interface{}{ll, strObj("z")}), types.JavaBoolFalse, "removeFirstOccurrence non-existent")
	assertJavaBoolLL(t, linkedlistRemoveLastOccurrence([]interface{}{ll, strObj("z")}), types.JavaBoolFalse, "removeLastOccurrence non-existent")
}

func TestLinkedList_Remove_Object_Variant_CurrentImpl(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)
	x := strObj("x")
	y := strObj("y")
	_ = linkedlistAddLast([]interface{}{ll, x})
	_ = linkedlistAddLast([]interface{}{ll, y})

	// Current impl of linkedlistRemove(self, element) removes first match and returns the element object
	ret := linkedlistRemove([]interface{}{ll, strObj("x")})
	obj, ok := ret.(*object.Object)
	if !ok || object.GoStringFromStringObject(obj) != "x" {
		t.Fatalf("linkedlistRemove returned wrong object: %T", ret)
	}
	// Now first should be 'y'
	if s := object.GoStringFromStringObject(linkedlistGetFirst([]interface{}{ll}).(*object.Object)); s != "y" {
		t.Fatalf("expected head 'y' after removal, got %q", s)
	}
}

func TestLinkedList_RemoveAtIndex_CurrentErrorBehavior(t *testing.T) {
	globals.InitGlobals("test")

	llobj := newLinkedListObj(t)
	_ = linkedlistAddLast([]interface{}{llobj, strObj("a")})

	// Per current impl, passing (list, index) causes an IllegalArgumentException since args[1] is not int64
	if err := linkedlistRemoveAtIndex([]interface{}{llobj, int32(86)}); err == nil {
		t.Fatalf("expected error for current removeAtIndex implementation")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
		}
	}
}

func TestLinkedList_Set_Size_ToString(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)
	_ = linkedlistAddLast([]interface{}{ll, strObj("a")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("b")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("c")})

	// set index 1 to "B"; returns old value
	old := linkedlistSet([]interface{}{ll, int64(1), strObj("B")}).(*object.Object)
	if object.GoStringFromStringObject(old) != "b" {
		t.Fatalf("set returned wrong old value: %q", object.GoStringFromStringObject(old))
	}

	// size should be 3
	if sz := linkedlistSize([]interface{}{ll}).(int64); sz != 3 {
		t.Fatalf("size mismatch: expected 3, got %d", sz)
	}

	// toString format: "LinkedList{a, B, c}"
	sObj := LinkedlistToString([]interface{}{ll}).(*object.Object)
	if s := object.GoStringFromStringObject(sObj); s != "LinkedList{a, B, c}" {
		t.Fatalf("toString mismatch: got %q", s)
	}

	// bounds errors for set
	if err := linkedlistSet([]interface{}{ll, int64(-1), strObj("x")}); err == nil {
		t.Fatalf("expected error for negative index in set")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}
	if err := linkedlistSet([]interface{}{ll, int64(99), strObj("x")}); err == nil {
		t.Fatalf("expected error for large index in set")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}
}

func TestLinkedList_Sort_Placeholder(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)
	_ = linkedlistAddLast([]interface{}{ll, strObj("b")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("a")})

	// Current implementation expects the list object as the sole arg and returns nil
	if ret := linkedlistSort([]interface{}{ll}); ret != nil {
		t.Fatalf("sort returned error: %v", ret)
	}
	// Order remains the same (no actual comparator logic)
	if s := object.GoStringFromStringObject(linkedlistGetFirst([]interface{}{ll}).(*object.Object)); s != "b" {
		t.Fatalf("sort placeholder changed order unexpectedly; head=%q", s)
	}
}

func TestLinkedList_ToArray_And_ToArrayTyped_Variants(t *testing.T) {
	globals.InitStringPool()

	// Build list [a, b]
	ll := newLinkedListObj(t)
	a := strObj("a")
	b := strObj("b")
	_ = linkedlistAddLast([]interface{}{ll, a})
	_ = linkedlistAddLast([]interface{}{ll, b})

	// toArray -> object with []interface{}
	arrObj := linkedlistToArray([]interface{}{ll}).(*object.Object)
	vals, ok := arrObj.FieldTable["value"].Fvalue.([]interface{})
	if !ok || len(vals) != 2 {
		t.Fatalf("toArray returned wrong type/length: %T, %d", arrObj.FieldTable["value"].Fvalue, len(vals))
	}
	if vals[0] != a || vals[1] != b {
		t.Fatalf("toArray element mismatch")
	}

	// toArrayTyped with object array (len >= listLen) populates and null-terminates
	objArr := &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: "[Ljava/lang/Object;", Fvalue: make([]*object.Object, 3)},
	}}
	ret1 := linkedlistToArrayTyped([]interface{}{ll, objArr}).(*object.Object)
	out1 := ret1.FieldTable["value"].Fvalue.([]*object.Object)
	if out1[0] != a || out1[1] != b || out1[2] != nil {
		t.Fatalf("toArrayTyped object array populate/null-terminate mismatch: %v", out1)
	}

	// toArrayTyped with object array (len < listLen) returns a new object
	objArrSmall := &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: "[Ljava/lang/Object;", Fvalue: make([]*object.Object, 1)},
	}}
	ret2 := linkedlistToArrayTyped([]interface{}{ll, objArrSmall}).(*object.Object)
	out2 := ret2.FieldTable["value"].Fvalue.([]*object.Object)
	if len(out2) != 2 || out2[0] != a || out2[1] != b {
		t.Fatalf("toArrayTyped new object array mismatch")
	}

	// Build an int64 list [10, 20, 30]
	llNums := newLinkedListObj(t)
	_ = linkedlistAddLast([]interface{}{llNums, int64(10)})
	_ = linkedlistAddLast([]interface{}{llNums, int64(20)})
	_ = linkedlistAddLast([]interface{}{llNums, int64(30)})

	// int64 array with sufficient length
	intArr := &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: "[I", Fvalue: make([]int64, 4)},
	}}
	_ = linkedlistToArrayTyped([]interface{}{llNums, intArr}).(*object.Object)
	ints := intArr.FieldTable["value"].Fvalue.([]int64)
	if ints[0] != 10 || ints[1] != 20 || ints[2] != 30 || ints[3] != 0 {
		t.Fatalf("toArrayTyped int64 len>= mismatch: %v", ints)
	}

	// int64 array with insufficient length -> returns new object with exact length
	intArrSmall := &object.Object{FieldTable: map[string]object.Field{
		"value": {Ftype: "[I", Fvalue: make([]int64, 2)},
	}}
	ret3 := linkedlistToArrayTyped([]interface{}{llNums, intArrSmall}).(*object.Object)
	ints2 := ret3.FieldTable["value"].Fvalue.([]int64)
	if len(ints2) != 3 || ints2[0] != 10 || ints2[1] != 20 || ints2[2] != 30 {
		t.Fatalf("toArrayTyped int64 new array mismatch: %v", ints2)
	}

	// ArrayStoreException: provide object list but int64 array target
	if err := linkedlistToArrayTyped([]interface{}{ll, intArr}); err == nil {
		t.Fatalf("expected ArrayStoreException for mismatched element/array types")
	} else if geb, ok := err.(*ghelpers.GErrBlk); ok {
		if geb.ExceptionType != excNames.ArrayStoreException {
			t.Fatalf("expected ArrayStoreException, got %d", geb.ExceptionType)
		}
	}
}
