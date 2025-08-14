package gfunction

import (
    "jacobin/globals"
    "jacobin/object"
    "jacobin/types"
    "testing"
)

// Helper to extract class name string from an object
func classNameOf(obj *object.Object) string {
    return object.GoStringFromStringPoolIndex(obj.KlassName)
}

func TestShortDoubleValue(t *testing.T) {
    globals.InitStringPool()

    cases := []int64{0, 1, -1, 42, 127, -128, 32767, -32768}
    for _, v := range cases {
        shortObj := Populator("java/lang/Short", types.Short, v)
        res := shortDoubleValue([]interface{}{shortObj})
        d, ok := res.(float64)
        if !ok {
            t.Fatalf("expected float64 from shortDoubleValue, got %T", res)
        }
        if d != float64(v) {
            t.Fatalf("doubleValue mismatch: expected %v, got %v", float64(v), d)
        }
    }
}

func TestShortValueOf(t *testing.T) {
    globals.InitStringPool()

    cases := []int64{0, 1, -1, 12345, -12345, 32767, -32768}
    for _, v := range cases {
        res := shortValueOf([]interface{}{v})
        obj, ok := res.(*object.Object)
        if !ok {
            t.Fatalf("expected *object.Object from shortValueOf, got %T", res)
        }
        if cn := classNameOf(obj); cn != "java/lang/Short" {
            t.Fatalf("expected class java/lang/Short, got %s", cn)
        }
        // Check the boxed value
        val, ok := obj.FieldTable["value"].Fvalue.(int64)
        if !ok {
            t.Fatalf("expected int64 value field, got %T", obj.FieldTable["value"].Fvalue)
        }
        if val != v {
            t.Fatalf("valueOf mismatch: expected %d, got %d", v, val)
        }
    }
}

func TestShortRoundTrip_ValueOfThenDoubleValue(t *testing.T) {
    globals.InitStringPool()

    cases := []int64{7, -7, 30000, -30000}
    for _, v := range cases {
        obj := shortValueOf([]interface{}{v}).(*object.Object)
        res := shortDoubleValue([]interface{}{obj})
        d := res.(float64)
        if d != float64(v) {
            t.Fatalf("round-trip mismatch: expected %v, got %v", float64(v), d)
        }
    }
}
