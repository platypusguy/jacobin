package gfunction

import (
    "jacobin/src/excNames"
    "jacobin/src/globals"
    "jacobin/src/object"
    "jacobin/src/types"
    "testing"
)

// helper: expect an error block with the given exception type
func assertErrUnsafe(t *testing.T, got interface{}, want int) {
    t.Helper()
    geb, ok := got.(*GErrBlk)
    if !ok {
        t.Fatalf("expected *GErrBlk, got %T", got)
    }
    if geb.ExceptionType != want {
        t.Fatalf("expected exception %d, got %d", want, geb.ExceptionType)
    }
}

func TestUnsafe_ArrayBaseOffset_And_IndexScale(t *testing.T) {
    globals.InitStringPool()

    // arrayBaseOffset: null -> NPE
    if res := unsafeArrayBaseOffset([]interface{}{object.Null}); res == nil {
        t.Fatalf("expected NPE for null param in arrayBaseOffset")
    } else {
        assertErrUnsafe(t, res, excNames.NullPointerException)
    }

    // arrayBaseOffset: non-null -> 0
    dummy := object.MakeEmptyObjectWithClassName(&classUnsafeName) // any object works
    if v := unsafeArrayBaseOffset([]interface{}{dummy}).(int64); v != 0 {
        t.Fatalf("arrayBaseOffset expected 0, got %d", v)
    }

    // arrayIndexScale0: build a faux "array class" object whose value field type encodes the array kind
    // byte[] => 1
    arrClassB := object.MakePrimitiveObject("java/lang/Class", types.ByteArray, nil) // Ftype == "[B"
    if v := unsafeArrayIndexScale0([]interface{}{arrClassB}).(int64); v != 1 {
        t.Fatalf("indexScale0 for [B expected 1, got %d", v)
    }

    // boolean[] => 1
    arrClassZ := object.MakePrimitiveObject("java/lang/Class", types.BoolArray, nil) // Ftype == "[Z"
    if v := unsafeArrayIndexScale0([]interface{}{arrClassZ}).(int64); v != 1 {
        t.Fatalf("indexScale0 for [Z expected 1, got %d", v)
    }

    // multi-dim (e.g., int[][]) => 8 (pointers)
    arrClass2D := object.MakePrimitiveObject("java/lang/Class", "[[I", nil)
    if v := unsafeArrayIndexScale0([]interface{}{arrClass2D}).(int64); v != 8 {
        t.Fatalf("indexScale0 for [[I expected 8, got %d", v)
    }

    // unsafeArrayIndexScale delegates to indexScale0 (also handles null check)
    if v := unsafeArrayIndexScale([]interface{}{arrClassB}).(int64); v != 1 {
        t.Fatalf("indexScale for [B expected 1, got %d", v)
    }
    if res := unsafeArrayIndexScale([]interface{}{object.Null}); res == nil {
        t.Fatalf("expected NPE for null param in arrayIndexScale")
    } else {
        assertErrUnsafe(t, res, excNames.NullPointerException)
    }
}

func TestUnsafe_CompareAndSetInt_GetUnsafe_ObjectFieldOffset(t *testing.T) {
    globals.InitStringPool()

    // compareAndSetInt returns Java true (1) per current implementation
    if v := unsafeCompareAndSetInt([]interface{}{}).(int64); v != types.JavaBoolTrue {
        t.Fatalf("compareAndSetInt expected true (1), got %d", v)
    }

    // getUnsafe returns an object of class jdk/internal/misc/Unsafe
    u := unsafeGetUnsafe([]interface{}{}).(*object.Object)
    if cn := object.GoStringFromStringPoolIndex(u.KlassName); cn != classUnsafeName {
        t.Fatalf("getUnsafe class mismatch: %q", cn)
    }

    // objectFieldOffset1 returns 0
    if v := unsafeObjectFieldOffset1([]interface{}{}).(int64); v != 0 {
        t.Fatalf("objectFieldOffset1 expected 0, got %d", v)
    }
}

func TestUnsafe_GetIntVolatile(t *testing.T) {
    globals.InitStringPool()

    // Behavior: when params[1] is nil, hash := 0, result = offset
    obj := object.MakeEmptyObjectWithClassName(&classUnsafeName)
    offset := int64(123)
    v := unsafeGetIntVolatile([]interface{}{obj, nil, offset}).(int64)
    if v != offset {
        t.Fatalf("getIntVolatile expected %d when second param is nil, got %d", offset, v)
    }

    // When params[1] is an object, the result is (obj.Mark.Hash + offset); verify type only
    v2 := unsafeGetIntVolatile([]interface{}{obj, obj, int64(0)})
    if _, ok := v2.(int64); !ok {
        t.Fatalf("getIntVolatile with object second param did not return int64, got %T", v2)
    }
}

func TestUnsafe_GetLong_ErrorPaths(t *testing.T) {
    globals.InitStringPool()

    // Second param not an object -> returns 0
    res1 := unsafeGetLong([]interface{}{nil, int64(5)})
    if res1.(int64) != 0 {
        t.Fatalf("getLong with non-object second param expected 0, got %v", res1)
    }

    // Second param object but third (offset) wrong type -> returns 0
    some := object.MakeEmptyObjectWithClassName(&classUnsafeName)
    res2 := unsafeGetLong([]interface{}{nil, some, "bad-offset"})
    if res2.(int64) != 0 {
        t.Fatalf("getLong with bad offset expected 0, got %v", res2)
    }
}
