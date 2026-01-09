package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helpers
func newLinkedListObj(t *testing.T) *object.Object {
	t.Helper()
	ll := object.MakeEmptyObjectWithClassName(&classNameLinkedList)
	if ret := linkedlistInit([]interface{}{ll}); ret != nil {
		t.Fatalf("linkedlistInit returned error: %v", ret)
	}
	return ll
}

func strObj(s string) *object.Object { return object.StringObjectFromGoString(s) }

func assertJavaBoolLL(t *testing.T, got interface{}, want int64, msg string) {
	t.Helper()
	gi, ok := got.(int64)
	if !ok {
		switch got.(type) {
		case *GErrBlk:
			errBlk := *got.(*GErrBlk)
			t.Fatalf("%s: expected Java boolean (int64), got GErrBlk %d (%s)", msg, errBlk.ExceptionType, errBlk.ErrMsg)
		}
		t.Fatalf("%s: expected Java boolean (int64), got %T", msg, got)
	}
	if gi != want {
		t.Fatalf("%s: expected %d, got %d", msg, want, gi)
	}
}

func TestLinkedList_Init_And_IsEmpty(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)

	// Initially empty
	assertJavaBoolLL(t, linkedlistIsEmpty([]interface{}{ll}), types.JavaBoolTrue, "isEmpty initially")

	// addLastRetBool returns true and makes list non-empty
	ret := linkedlistAddLastRetBool([]interface{}{ll, strObj("first")})
	assertJavaBoolLL(t, ret, types.JavaBoolTrue, "addLastRetBool should return true")
	assertJavaBoolLL(t, linkedlistIsEmpty([]interface{}{ll}), types.JavaBoolFalse, "isEmpty after add")

	// Clear should empty it
	if err := linkedlistClear([]interface{}{ll}); err != nil {
		t.Fatalf("clear returned error: %v", err)
	}
	assertJavaBoolLL(t, linkedlistIsEmpty([]interface{}{ll}), types.JavaBoolTrue, "isEmpty after clear")
}

func TestLinkedList_AddFirst_Last_And_Getters(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)

	// addLast then getFirst/getLast
	_ = linkedlistAddLast([]interface{}{ll, strObj("b")})
	// addFirst places at head
	_ = linkedlistAddFirst([]interface{}{ll, strObj("a")})

	// getFirst -> "a"
	vFirst := linkedlistGetFirst([]interface{}{ll}).(*object.Object)
	if object.GoStringFromStringObject(vFirst) != "a" {
		t.Fatalf("getFirst mismatch: expected 'a', got %q", object.GoStringFromStringObject(vFirst))
	}
	// getLast -> "b"
	vLast := linkedlistGetLast([]interface{}{ll}).(*object.Object)
	if object.GoStringFromStringObject(vLast) != "b" {
		t.Fatalf("getLast mismatch: expected 'b', got %q", object.GoStringFromStringObject(vLast))
	}

	// get by index
	v0 := linkedlistGet([]interface{}{ll, int64(0)}).(*object.Object)
	if object.GoStringFromStringObject(v0) != "a" {
		t.Fatalf("get(0) mismatch: expected 'a', got %q", object.GoStringFromStringObject(v0))
	}
	v1 := linkedlistGet([]interface{}{ll, int64(1)}).(*object.Object)
	if object.GoStringFromStringObject(v1) != "b" {
		t.Fatalf("get(1) mismatch: expected 'b', got %q", object.GoStringFromStringObject(v1))
	}

	// get out-of-bounds -> IndexOutOfBoundsException
	if err := linkedlistGet([]interface{}{ll, int64(-1)}); err == nil {
		t.Fatalf("expected error for negative index")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}
	if err := linkedlistGet([]interface{}{ll, int64(2)}); err == nil {
		t.Fatalf("expected error for index past end")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}

	// On an empty list, getFirst/getLast should throw NoSuchElementException
	ll2 := newLinkedListObj(t)
	if err := linkedlistGetFirst([]interface{}{ll2}); err == nil {
		t.Fatalf("expected error for getFirst on empty list")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.NoSuchElementException {
			t.Fatalf("expected NoSuchElementException, got %d", geb.ExceptionType)
		}
	}
	if err := linkedlistGetLast([]interface{}{ll2}); err == nil {
		t.Fatalf("expected error for getLast on empty list")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.NoSuchElementException {
			t.Fatalf("expected NoSuchElementException, got %d", geb.ExceptionType)
		}
	}
}

func TestLinkedList_AddAtIndex_Bounds_And_Order(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)

	// Start with [a, c]
	_ = linkedlistAddLast([]interface{}{ll, strObj("a")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("c")})

	// Insert b at index 1 -> [a, b, c]
	if err := linkedlistAddAtIndex([]interface{}{ll, int64(1), strObj("b")}); err != nil {
		t.Fatalf("addAtIndex returned error: %v", err)
	}
	if s := object.GoStringFromStringObject(linkedlistGet([]interface{}{ll, int64(1)}).(*object.Object)); s != "b" {
		t.Fatalf("addAtIndex failed: expected 'b' at index 1, got %q", s)
	}

	// Insert d at end index == len -> [a, b, c, d]
	if err := linkedlistAddAtIndex([]interface{}{ll, int64(3), strObj("d")}); err != nil {
		t.Fatalf("addAtIndex at end returned error: %v", err)
	}
	if s := object.GoStringFromStringObject(linkedlistGet([]interface{}{ll, int64(3)}).(*object.Object)); s != "d" {
		t.Fatalf("expected 'd' at index 3, got %q", s)
	}

	// Bounds: negative or > len
	if err := linkedlistAddAtIndex([]interface{}{ll, int64(-1), strObj("x")}); err == nil {
		t.Fatalf("expected error for negative add index")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}
	if err := linkedlistAddAtIndex([]interface{}{ll, int64(99), strObj("x")}); err == nil {
		t.Fatalf("expected error for too-large add index")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IndexOutOfBoundsException {
			t.Fatalf("expected IndexOutOfBoundsException, got %d", geb.ExceptionType)
		}
	}
}

func TestLinkedList_Contains_For_Strings_And_Ints(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)

	// Add String elements
	_ = linkedlistAddLast([]interface{}{ll, strObj("x")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("y")})

	assertJavaBoolLL(t, linkedlistContains([]interface{}{ll, strObj("y")}), types.JavaBoolTrue, "contains 'y'")
	assertJavaBoolLL(t, linkedlistContains([]interface{}{ll, strObj("z")}), types.JavaBoolFalse, "not contains 'z'")

	// Add int64 primitive and test contains with int64
	_ = linkedlistAddLast([]interface{}{ll, int64(42)})
	assertJavaBoolLL(t, linkedlistContains([]interface{}{ll, int64(42)}), types.JavaBoolTrue, "contains 42")

	// Non-String object should raise UnsupportedOperationException in comparator/equality
	nonStringObj := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(7))
	if err := linkedlistContains([]interface{}{ll, nonStringObj}); err == nil {
		t.Fatalf("expected error for non-String object contains check")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.UnsupportedOperationException {
			t.Fatalf("expected UnsupportedOperationException, got %d", geb.ExceptionType)
		}
	}
}

func TestLinkedList_IndexOf_And_LastIndexOf(t *testing.T) {
	globals.InitStringPool()

	ll := newLinkedListObj(t)

	// Build [a, b, a, c]
	_ = linkedlistAddLast([]interface{}{ll, strObj("a")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("b")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("a")})
	_ = linkedlistAddLast([]interface{}{ll, strObj("c")})

	// indexOf("a") should be 0
	idxA := linkedlistIndexOf([]interface{}{ll, strObj("a")}).(int64)
	if idxA != 0 {
		t.Fatalf("indexOf('a') expected 0, got %d", idxA)
	}

	// indexOf("z") -> -1
	idxZ := linkedlistIndexOf([]interface{}{ll, strObj("z")}).(int64)
	if idxZ != -1 {
		t.Fatalf("indexOf('z') expected -1, got %d", idxZ)
	}

	// lastIndexOf currently returns distance from end per implementation; for last 'a' at head-index 2 in len=4, distance is 1
	lastA := linkedlistLastIndexOf([]interface{}{ll, strObj("a")}).(int64)
	if lastA != 1 {
		t.Fatalf("lastIndexOf('a') expected distance-from-end 1 per current impl, got %d", lastA)
	}
}
