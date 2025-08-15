package gfunction

import (
    "jacobin/excNames"
    "jacobin/object"
    "testing"
)

func TestLong_ValueOf_And_DoubleValue(t *testing.T) {
    // valueOf(J) -> Long object, then doubleValue() -> primitive double
    obj := longValueOf([]interface{}{int64(42)})
    lobj, ok := obj.(*object.Object)
    if !ok {
        t.Fatalf("valueOf did not return object, got %T", obj)
    }
    out := longDoubleValue([]interface{}{lobj})
    if got := out.(float64); got != 42.0 {
        t.Fatalf("doubleValue mismatch: got %v want %v", got, 42.0)
    }
}

func TestLong_ParseLong_Valid_And_Invalid(t *testing.T) {
    // NOTE: current implementation reads params[1] for the string argument
    s := object.StringObjectFromGoString("12345")
    out := longParseLong([]interface{}{nil, s})
    if got := out.(int64); got != 12345 {
        t.Fatalf("parseLong valid: got %d", got)
    }
    // invalid -> NumberFormatException
    sinv := object.StringObjectFromGoString("abc")
    out = longParseLong([]interface{}{nil, sinv})
    if geb, ok := out.(*GErrBlk); !ok || geb.ExceptionType != excNames.NumberFormatException {
        if !ok {
            t.Fatalf("parseLong invalid: expected *GErrBlk, got %T", out)
        }
        t.Fatalf("parseLong invalid: expected NumberFormatException, got %v", geb)
    }
}

func TestLong_RotateLeft_Right(t *testing.T) {
    // Use a known pattern and compare by unsigned bit pattern
    val := int64(0x0123456789abcdef)
    // rotate left by 8 -> 0x23456789abcdef01
    outL := longRotateLeft([]interface{}{val, int64(8)})
    gotL := uint64(outL.(int64))
    var expectedL uint64 = 0x23456789abcdef01
    if gotL != expectedL {
        t.Fatalf("rotateLeft mismatch: got 0x%016x want 0x%016x", gotL, expectedL)
    }
    // rotate right by 12 -> expected using uint64 rotation
    outR := longRotateRight([]interface{}{val, int64(12)})
    gotR := uint64(outR.(int64))
    // compute expected separately without constant overflow
    var base uint64 = 0x0123456789abcdef
    expectedR := (base >> 12) | (base << (64 - 12))
    if gotR != expectedR {
        t.Fatalf("rotateRight mismatch: got 0x%016x want 0x%016x", gotR, expectedR)
    }
}

func TestLong_ToHexString_And_ToString(t *testing.T) {
    // toHexString pads to 16 digits per implementation
    out := longToHexString([]interface{}{int64(1)})
    sObj := out.(*object.Object)
    if got := object.GoStringFromStringObject(sObj); got != "0000000000000001" {
        t.Fatalf("toHexString(1) got %q", got)
    }
    out = longToHexString([]interface{}{int64(-1)})
    sObj = out.(*object.Object)
    if got := object.GoStringFromStringObject(sObj); got != "ffffffffffffffff" {
        t.Fatalf("toHexString(-1) got %q", got)
    }

    // toString decimal
    out = longToString([]interface{}{int64(-123)})
    sObj = out.(*object.Object)
    if got := object.GoStringFromStringObject(sObj); got != "-123" {
        t.Fatalf("toString(-123) got %q", got)
    }
}
