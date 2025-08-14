package gfunction

import (
    "jacobin/excNames"
    "jacobin/globals"
    "jacobin/object"
    "jacobin/types"
    "testing"
)

// Helpers
func newAtomicIntegerObj() *object.Object {
    return object.MakeEmptyObjectWithClassName(&atomicIntegerClassName)
}

func aiGet(t *testing.T, ai *object.Object) int64 {
    t.Helper()
    v := atomicIntegerGet([]interface{}{ai})
    vi, ok := v.(int64)
    if !ok {
        t.Fatalf("atomicIntegerGet did not return int64, got %T", v)
    }
    return vi
}

func TestAtomicInteger_Init_And_Set_Get(t *testing.T) {
    globals.InitStringPool()

    // Default init -> 0
    ai := newAtomicIntegerObj()
    if ret := atomicIntegerInitVoid([]interface{}{ai}); ret != nil {
        t.Fatalf("<init>() returned error: %v", ret)
    }
    if got := aiGet(t, ai); got != 0 {
        t.Fatalf("expected default 0, got %d", got)
    }

    // init(I) -> sets to provided
    ai2 := newAtomicIntegerObj()
    if ret := atomicIntegerInitInt([]interface{}{ai2, int64(42)}); ret != nil {
        t.Fatalf("<init>(I) returned error: %v", ret)
    }
    if got := aiGet(t, ai2); got != 42 {
        t.Fatalf("expected 42, got %d", got)
    }

    // set(I) (and other set variants map to same impl)
    if ret := atomicIntegerSet([]interface{}{ai2, int64(-7)}); ret != nil {
        t.Fatalf("set(I)V returned error: %v", ret)
    }
    if got := aiGet(t, ai2); got != -7 {
        t.Fatalf("expected -7 after set, got %d", got)
    }

    // setPlain -> same as set
    if ret := atomicIntegerSet([]interface{}{ai2, int64(100)}); ret != nil {
        t.Fatalf("setPlain(I)V returned error: %v", ret)
    }
    if got := aiGet(t, ai2); got != 100 {
        t.Fatalf("expected 100 after setPlain, got %d", got)
    }
}

func TestAtomicInteger_GetAndSet_And_CompareAndSet(t *testing.T) {
    globals.InitStringPool()

    ai := newAtomicIntegerObj()
    _ = atomicIntegerInitInt([]interface{}{ai, int64(5)})

    // getAndSet returns old value and updates to new
    old := atomicIntegerGetAndSet([]interface{}{ai, int64(9)}).(int64)
    if old != 5 {
        t.Fatalf("getAndSet old mismatch: expected 5, got %d", old)
    }
    if got := aiGet(t, ai); got != 9 {
        t.Fatalf("getAndSet new mismatch: expected 9, got %d", got)
    }

    // compareAndSet success
    res := atomicIntegerCompareAndSet([]interface{}{ai, int64(9), int64(11)}).(int64)
    if res != types.JavaBoolTrue {
        t.Fatalf("compareAndSet expected true, got %d", res)
    }
    if got := aiGet(t, ai); got != 11 {
        t.Fatalf("compareAndSet did not set new value, got %d", got)
    }

    // compareAndSet failure (mismatched expected)
    res2 := atomicIntegerCompareAndSet([]interface{}{ai, int64(99), int64(123)}).(int64)
    if res2 != types.JavaBoolFalse {
        t.Fatalf("compareAndSet expected false, got %d", res2)
    }
    if got := aiGet(t, ai); got != 11 {
        t.Fatalf("compareAndSet failure should not change value, got %d", got)
    }
}

func TestAtomicInteger_IncDec_Add_Variants(t *testing.T) {
    globals.InitStringPool()

    ai := newAtomicIntegerObj()
    _ = atomicIntegerInitInt([]interface{}{ai, int64(0)})

    // getAndIncrement: returns old, then increments
    old := atomicIntegerGetAndIncrement([]interface{}{ai}).(int64)
    if old != 0 || aiGet(t, ai) != 1 {
        t.Fatalf("getAndIncrement mismatch: old=%d cur=%d", old, aiGet(t, ai))
    }

    // getAndDecrement
    old = atomicIntegerGetAndDecrement([]interface{}{ai}).(int64)
    if old != 1 || aiGet(t, ai) != 0 {
        t.Fatalf("getAndDecrement mismatch: old=%d cur=%d", old, aiGet(t, ai))
    }

    // getAndAdd(+5)
    old = atomicIntegerGetAndAdd([]interface{}{ai, int64(5)}).(int64)
    if old != 0 || aiGet(t, ai) != 5 {
        t.Fatalf("getAndAdd(+5) mismatch: old=%d cur=%d", old, aiGet(t, ai))
    }

    // getAndAdd(-2)
    old = atomicIntegerGetAndAdd([]interface{}{ai, int64(-2)}).(int64)
    if old != 5 || aiGet(t, ai) != 3 {
        t.Fatalf("getAndAdd(-2) mismatch: old=%d cur=%d", old, aiGet(t, ai))
    }

    // incrementAndGet
    newv := atomicIntegerIncrementAndGet([]interface{}{ai}).(int64)
    if newv != 4 || aiGet(t, ai) != 4 {
        t.Fatalf("incrementAndGet mismatch: ret=%d cur=%d", newv, aiGet(t, ai))
    }

    // decrementAndGet
    newv = atomicIntegerDecrementAndGet([]interface{}{ai}).(int64)
    if newv != 3 || aiGet(t, ai) != 3 {
        t.Fatalf("decrementAndGet mismatch: ret=%d cur=%d", newv, aiGet(t, ai))
    }

    // addAndGet(+10)
    newv = atomicIntegerAddAndGet([]interface{}{ai, int64(10)}).(int64)
    if newv != 13 || aiGet(t, ai) != 13 {
        t.Fatalf("addAndGet(+10) mismatch: ret=%d cur=%d", newv, aiGet(t, ai))
    }
}

func TestAtomicInteger_ToString_And_ToFloat(t *testing.T) {
    globals.InitStringPool()

    ai := newAtomicIntegerObj()
    _ = atomicIntegerInitInt([]interface{}{ai, int64(-123)})

    // toString
    so := atomicIntegerToString([]interface{}{ai}).(*object.Object)
    s := object.GoStringFromStringObject(so)
    if s != "-123" {
        t.Fatalf("toString mismatch: expected -123, got %q", s)
    }

    // doubleValue()D and floatValue()F both map to atomicIntegerToFloat and return float64
    d1 := atomicIntegerToFloat([]interface{}{ai}).(float64)
    if d1 != -123.0 {
        t.Fatalf("doubleValue mismatch: expected -123.0, got %v", d1)
    }
    d2 := atomicIntegerToFloat([]interface{}{ai}).(float64)
    if d2 != -123.0 {
        t.Fatalf("floatValue mismatch: expected -123.0, got %v", d2)
    }
}

func TestAtomicInteger_ErrorPaths_In_AddHelper(t *testing.T) {
    globals.InitStringPool()

    // Wrong number of parameters
    if err := fnAtomicIntegerAdd([]interface{}{}, false); err == nil {
        t.Fatalf("expected error for wrong param count")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IllegalArgumentException {
            t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
        }
    }

    // First param null -> ClassCastException
    if err := fnAtomicIntegerAdd([]interface{}{object.Null, int64(1)}, false); err == nil {
        t.Fatalf("expected error for null object param")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.ClassCastException {
            t.Fatalf("expected ClassCastException, got %d", geb.ExceptionType)
        }
    }

    // First param wrong type -> ClassCastException
    if err := fnAtomicIntegerAdd([]interface{}{int64(5), int64(1)}, false); err == nil {
        t.Fatalf("expected error for non-object first param")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.ClassCastException {
            t.Fatalf("expected ClassCastException, got %d", geb.ExceptionType)
        }
    }

    // Second param wrong type -> ClassCastException
    ai := newAtomicIntegerObj()
    _ = atomicIntegerInitInt([]interface{}{ai, int64(0)})
    if err := fnAtomicIntegerAdd([]interface{}{ai, "not-int64"}, false); err == nil {
        t.Fatalf("expected error for non-int64 second param")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.ClassCastException {
            t.Fatalf("expected ClassCastException, got %d", geb.ExceptionType)
        }
    }

    // Missing value field -> NoSuchFieldException
    ai2 := newAtomicIntegerObj()
    // Intentionally do not init to create missing 'value'
    if err := fnAtomicIntegerAdd([]interface{}{ai2, int64(1)}, false); err == nil {
        t.Fatalf("expected error for missing value field")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.NoSuchFieldException {
            t.Fatalf("expected NoSuchFieldException, got %d", geb.ExceptionType)
        }
    }

    // Wrong field type -> IllegalArgumentException
    ai3 := newAtomicIntegerObj()
    ai3.FieldTable["value"] = object.Field{Ftype: types.Long, Fvalue: int64(0)} // wrong Ftype on purpose
    if err := fnAtomicIntegerAdd([]interface{}{ai3, int64(1)}, false); err == nil {
        t.Fatalf("expected error for wrong field type")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IllegalArgumentException {
            t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
        }
    }

    // Non-int64 value inside field -> IllegalArgumentException
    ai4 := newAtomicIntegerObj()
    ai4.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: "not-int64"}
    if err := fnAtomicIntegerAdd([]interface{}{ai4, int64(1)}, false); err == nil {
        t.Fatalf("expected error for non-int64 field value")
    } else if geb, ok := err.(*GErrBlk); ok {
        if geb.ExceptionType != excNames.IllegalArgumentException {
            t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
        }
    }
}
