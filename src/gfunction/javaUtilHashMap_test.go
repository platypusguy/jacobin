package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helpers
func newHashMapObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&classNameHashMap)
}

func hmInit(t *testing.T, hm *object.Object) {
	t.Helper()
	if ret := hashmapInit([]interface{}{hm}); ret != nil {
		t.Fatalf("hashmapInit returned error: %v", ret)
	}
}

func strKey(s string) *object.Object { return object.StringObjectFromGoString(s) }
func intKey(v int64) *object.Object  { return Populator("java/lang/Integer", types.Int, v) }

func TestHashMap_Put_Get_Size_Contains_Remove_StringKeys(t *testing.T) {
	globals.InitStringPool()

	hm := newHashMapObj()
	hmInit(t, hm)

	k := strKey("alpha")
	v1 := object.StringObjectFromGoString("one")

	// Initial put returns null (no previous)
	prev := hashmapPut([]interface{}{hm, k, v1})
	if prev != object.Null {
		t.Fatalf("expected null previous value on first put, got %T", prev)
	}

	// Size should be 1
	if sz := hashmapSize([]interface{}{hm}).(int64); sz != 1 {
		t.Fatalf("expected size 1, got %d", sz)
	}

	// containsKey true
	if ck := hashmapContainsKey([]interface{}{hm, k}).(int64); ck != types.JavaBoolTrue {
		t.Fatalf("expected containsKey true, got %d", ck)
	}

	// get should return same object reference as v1
	got := hashmapGet([]interface{}{hm, k}).(*object.Object)
	if got != v1 {
		t.Fatalf("get returned different object than stored")
	}

	// Second put with same key returns previous value
	v2 := object.StringObjectFromGoString("uno")
	prev2 := hashmapPut([]interface{}{hm, k, v2})
	if prev2 != v1 {
		t.Fatalf("expected previous value object from second put")
	}

	// remove returns current value and decreases size
	rem := hashmapRemove([]interface{}{hm, k})
	if rem != v2 {
		t.Fatalf("remove returned wrong value")
	}
	if sz := hashmapSize([]interface{}{hm}).(int64); sz != 0 {
		t.Fatalf("expected size 0 after remove, got %d", sz)
	}

	// get/contains on missing returns null/false
	if val := hashmapGet([]interface{}{hm, k}); val != object.Null {
		t.Fatalf("expected null on get of missing key, got %T", val)
	}
	if ck2 := hashmapContainsKey([]interface{}{hm, k}).(int64); ck2 != types.JavaBoolFalse {
		t.Fatalf("expected containsKey false on missing key")
	}
}

func TestHashMap_IsEmpty(t *testing.T) {
	globals.InitStringPool()

	hm := newHashMapObj()
	hmInit(t, hm)

	// Initially empty
	if empty := hashmapIsEmpty([]interface{}{hm}).(int64); empty != types.JavaBoolTrue {
		t.Fatalf("expected isEmpty true initially")
	}

	// After put
	k := strKey("a")
	v := object.StringObjectFromGoString("A")
	_ = hashmapPut([]interface{}{hm, k, v})
	if empty := hashmapIsEmpty([]interface{}{hm}).(int64); empty != types.JavaBoolFalse {
		t.Fatalf("expected isEmpty false after put")
	}

	// After remove
	_ = hashmapRemove([]interface{}{hm, k})
	if empty := hashmapIsEmpty([]interface{}{hm}).(int64); empty != types.JavaBoolTrue {
		t.Fatalf("expected isEmpty true after remove")
	}
}

func TestHashMap_IntKeys_And_MixedValues(t *testing.T) {
	globals.InitStringPool()

	hm := newHashMapObj()
	hmInit(t, hm)

	k1 := intKey(42)
	v1 := Populator("java/lang/Long", types.Long, int64(9001))

	// put and get with int key
	_ = hashmapPut([]interface{}{hm, k1, v1})
	got := hashmapGet([]interface{}{hm, intKey(42)}).(*object.Object)
	if got != v1 {
		t.Fatalf("expected to retrieve same value object for int key")
	}
}

func TestHashMap_PutAll_MergesEntries(t *testing.T) {
	globals.InitStringPool()

	src := newHashMapObj()
	dst := newHashMapObj()
	hmInit(t, src)
	hmInit(t, dst)

	// Populate src with 2 entries
	_ = hashmapPut([]interface{}{src, strKey("a"), object.StringObjectFromGoString("A")})
	_ = hashmapPut([]interface{}{src, intKey(7), object.StringObjectFromGoString("seven")})

	// Merge into dst
	if ret := hashmapPutAll([]interface{}{dst, src}); ret != nil {
		t.Fatalf("putAll returned error: %v", ret)
	}

	if sz := hashmapSize([]interface{}{dst}).(int64); sz != 2 {
		t.Fatalf("expected size 2 in dst after putAll, got %d", sz)
	}

	if v := hashmapGet([]interface{}{dst, strKey("a")}).(*object.Object); object.GoStringFromStringObject(v) != "A" {
		t.Fatalf("expected key 'a' -> 'A' after putAll")
	}
	if v := hashmapGet([]interface{}{dst, intKey(7)}).(*object.Object); object.GoStringFromStringObject(v) != "seven" {
		t.Fatalf("expected key 7 -> 'seven' after putAll")
	}
}

func TestHashMap_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	hm := newHashMapObj()
	hmInit(t, hm)

	// Wrong param count
	if err := hashmapPut([]interface{}{hm, strKey("k")}); err == nil {
		t.Fatalf("expected error for wrong param count in put")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException")
		}
	}

	// First param not object -> ClassCastException
	if err := hashmapGet([]interface{}{int64(5), strKey("k")}); err == nil {
		t.Fatalf("expected error for non-object first param in get")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.ClassCastException {
			t.Fatalf("expected ClassCastException")
		}
	}

	// First param wrong class -> IllegalArgumentException
	notMap := object.MakeEmptyObjectWithClassName(&classNameBase64Encoder)
	if err := hashmapSize([]interface{}{notMap}); err == nil {
		t.Fatalf("expected error for wrong class in size")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException")
		}
	}

	// Invalid key param (not an object) -> IllegalArgumentException from _getKey
	if err := hashmapPut([]interface{}{hm, int64(1), object.StringObjectFromGoString("v")}); err == nil {
		t.Fatalf("expected error for invalid key param (non-object)")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException")
		}
	}

	// Uninitialized map field: create object of HashMap class but do not call init
	raw := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	if err := hashmapSize([]interface{}{raw}); err == nil {
		t.Fatalf("expected error for missing map field in size")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException")
		}
	}

	// containsKey wrong param count
	if err := hashmapContainsKey([]interface{}{hm}); err == nil {
		t.Fatalf("expected error for wrong param count in containsKey")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException")
		}
	}

	// remove wrong param count
	if err := hashmapRemove([]interface{}{hm}); err == nil {
		t.Fatalf("expected error for wrong param count in remove")
	} else if geb, ok := err.(*GErrBlk); ok {
		if geb.ExceptionType != excNames.IllegalArgumentException {
			t.Fatalf("expected IllegalArgumentException")
		}
	}
}
