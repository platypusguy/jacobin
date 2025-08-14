package gfunction

import (
    "jacobin/excNames"
    "jacobin/globals"
    "jacobin/object"
    "testing"
)

// Helpers
func newPropertiesObj() *object.Object {
    return object.MakeEmptyObjectWithClassName(&classNameProperties)
}

func propInit(t *testing.T, p *object.Object) {
    t.Helper()
    if ret := propertiesInit([]interface{}{p}); ret != nil {
        t.Fatalf("propertiesInit returned error: %v", ret)
    }
}

func s(str string) *object.Object { return object.StringObjectFromGoString(str) }

func expectErrType(t *testing.T, got interface{}, expected int) {
    t.Helper()
    geb, ok := got.(*GErrBlk)
    if !ok {
        t.Fatalf("expected error block, got %T", got)
    }
    if geb.ExceptionType != expected {
        t.Fatalf("expected exception type %d, got %d", expected, geb.ExceptionType)
    }
}

func TestProperties_Set_Get_Size_Remove_ToString(t *testing.T) {
    globals.InitStringPool()

    p := newPropertiesObj()
    propInit(t, p)

    // Initially empty size
    if sz := propertiesSize([]interface{}{p}).(int64); sz != 0 {
        t.Fatalf("expected initial size 0, got %d", sz)
    }

    // getProperty missing -> null
    if v := propertiesGetProperty([]interface{}{p, s("missing")}); v != object.Null {
        t.Fatalf("expected null for missing key, got %T", v)
    }

    // getProperty with default should return default when missing (ensure no error from impl quirk)
    def := propertiesGetProperty([]interface{}{p, s("missing"), s("DEF")}).(*object.Object)
    if object.GoStringFromStringObject(def) != "DEF" {
        t.Fatalf("expected default DEF, got %q", object.GoStringFromStringObject(def))
    }

    // setProperty first time -> returns null, size increases
    if ret := propertiesSetProperty([]interface{}{p, s("alpha"), s("one")}); ret != object.Null {
        t.Fatalf("expected null from first set, got %T", ret)
    }
    if sz := propertiesSize([]interface{}{p}).(int64); sz != 1 {
        t.Fatalf("size after first set expected 1, got %d", sz)
    }

    // getProperty returns value
    g := propertiesGetProperty([]interface{}{p, s("alpha")} ).(*object.Object)
    if object.GoStringFromStringObject(g) != "one" {
        t.Fatalf("getProperty mismatch: expected 'one', got %q", object.GoStringFromStringObject(g))
    }

    // second set returns previous value
    prev := propertiesSetProperty([]interface{}{p, s("alpha"), s("uno")} ).(*object.Object)
    if object.GoStringFromStringObject(prev) != "one" {
        t.Fatalf("second set should return previous 'one', got %q", object.GoStringFromStringObject(prev))
    }

    // add another key to test toString ordering (case-insensitive sort)
    _ = propertiesSetProperty([]interface{}{p, s("Beta"), s("BVAL")})
    // toString format: {alpha=uno, Beta=BVAL} with case-insensitive ordering -> "alpha" then "Beta"
    ts := propertiesToString([]interface{}{p}).(*object.Object)
    if s := object.GoStringFromStringObject(ts); s != "{alpha=uno, Beta=BVAL}" {
        t.Fatalf("toString mismatch: got %q", s)
    }

    // remove existing key: current impl returns null even if present; size decreases
    rem := propertiesRemove([]interface{}{p, s("alpha")} )
    if rem != object.Null {
        t.Fatalf("remove expected to return null per current impl, got %T", rem)
    }
    if sz := propertiesSize([]interface{}{p}).(int64); sz != 1 {
        t.Fatalf("size after remove expected 1, got %d", sz)
    }
}

func TestProperties_Error_Paths(t *testing.T) {
    globals.InitStringPool()

    p := newPropertiesObj()
    propInit(t, p)

    // Non-object first param -> IllegalArgumentException
    if err := propertiesSize([]interface{}{int64(5)}); err == nil {
        t.Fatalf("expected error for non-object first param in size")
    } else { expectErrType(t, err, excNames.IllegalArgumentException) }

    // Missing map field (uninitialized) -> error
    raw := object.MakeEmptyObjectWithClassName(&classNameProperties)
    if err := propertiesSize([]interface{}{raw}); err == nil {
        t.Fatalf("expected error for missing map field in size")
    } else { expectErrType(t, err, excNames.IllegalArgumentException) }

    // setProperty with non-object key -> error
    if err := propertiesSetProperty([]interface{}{p, int64(1), s("x")}); err == nil {
        t.Fatalf("expected error for non-object key in setProperty")
    } else { expectErrType(t, err, excNames.IllegalArgumentException) }

    // setProperty with non-object value -> error
    if err := propertiesSetProperty([]interface{}{p, s("k"), int64(7)}); err == nil {
        t.Fatalf("expected error for non-object value in setProperty")
    } else { expectErrType(t, err, excNames.IllegalArgumentException) }

    // getProperty with non-object key -> error
    if err := propertiesGetProperty([]interface{}{p, int64(3)}); err == nil {
        t.Fatalf("expected error for non-object key in getProperty")
    } else { expectErrType(t, err, excNames.IllegalArgumentException) }

    // remove with non-object key -> error
    if err := propertiesRemove([]interface{}{p, int64(9)}); err == nil {
        t.Fatalf("expected error for non-object key in remove")
    } else { expectErrType(t, err, excNames.IllegalArgumentException) }
}
