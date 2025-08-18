package gfunction

import (
    "fmt"
    "jacobin/src/excNames"
    "jacobin/src/object"
    "jacobin/src/types"
    "math"
    "testing"
)

// helper to create a java/lang/Float object with a given value (stored as float64)
// Note: internal helper getFloat64ValueFromObject expects field type == types.Double,
// so we construct the Float object with Double field type to interoperate with Float methods.
func makeFloat(val float64) *object.Object {
    return object.MakePrimitiveObject("java/lang/Float", types.Double, val)
}

func TestFloat_ValueOf_And_FloatValue(t *testing.T) {
    v := 12.75
    // valueOf(F) should return a Float object holding the value; its field type is types.Float
    obj := floatValueOf([]interface{}{v})
    fobj, ok := obj.(*object.Object)
    if !ok {
        t.Fatalf("valueOf did not return object, got %T", obj)
    }
    // verify stored value matches
    if got, ok2 := fobj.FieldTable["value"].Fvalue.(float64); !ok2 || got != v {
        t.Fatalf("valueOf stored value mismatch: got %v", fobj.FieldTable["value"].Fvalue)
    }
    // Now test floatValue() using a compatible object created via makeFloat
    out := floatFloatValue([]interface{}{makeFloat(v)})
    if got := out.(float64); got != v {
        t.Fatalf("floatFloatValue mismatch: got %v want %v", got, v)
    }
}

func TestFloat_ParseFloat_Valid_Invalid_Empty(t *testing.T) {
    // valid
    s := object.StringObjectFromGoString("3.5")
    out := floatParseFloat([]interface{}{s})
    if got := out.(float64); float32(got) != float32(3.5) {
        t.Fatalf("parseFloat valid: got %v", got)
    }
    // invalid -> NumberFormatException
    sInv := object.StringObjectFromGoString("abc")
    out = floatParseFloat([]interface{}{sInv})
    if geb, ok := out.(*GErrBlk); !ok || geb.ExceptionType != excNames.NumberFormatException {
        if !ok {
            t.Fatalf("parseFloat invalid: expected *GErrBlk, got %T", out)
        }
        t.Fatalf("parseFloat invalid: expected NumberFormatException, got %v", geb)
    }
    // empty -> NullPointerException according to implementation
    sEmpty := object.StringObjectFromGoString("")
    out = floatParseFloat([]interface{}{sEmpty})
    if geb, ok := out.(*GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
        t.Fatalf("parseFloat empty: expected NullPointerException, got %T (%v)", out, out)
    }
}

func TestFloat_IntBits_RoundTrip_And_NegZero(t *testing.T) {
    // round-trip a normal value via bits
    v := float32(-7.25)
    bits := math.Float32bits(v)
    gotF := floatIntBitsToFloat([]interface{}{int64(bits)}).(float64)
    if math.Float32bits(float32(gotF)) != bits {
        t.Fatalf("intBitsToFloat roundtrip mismatch: got %x want %x", math.Float32bits(float32(gotF)), bits)
    }

    // floatToIntBits should produce float32 bits of given float64 value
    v2 := 5.5
    bits2 := floatFloatToIntBits([]interface{}{v2}).(int64)
    if uint32(bits2) != math.Float32bits(float32(v2)) {
        t.Fatalf("floatToIntBits mismatch: got %08x want %08x", uint32(bits2), math.Float32bits(float32(v2)))
    }

    // negative zero: 0x80000000
    negZeroBits := uint32(0x80000000)
    gotNegZero := floatIntBitsToFloat([]interface{}{int64(negZeroBits)}).(float64)
    if !math.Signbit(gotNegZero) {
        t.Fatalf("expected negative zero sign bit")
    }
}

func TestFloat_ToHexString_And_ToString(t *testing.T) {
    v := 123.25
    obj := makeFloat(v)
    // toHexString prints Float64bits (per current implementation)
    sObj := floatToHexString([]interface{}{obj}).(*object.Object)
    gotHex := object.GoStringFromStringObject(sObj)
    expectedHex := fmt.Sprintf("0x%016X", math.Float64bits(v))
    if gotHex != expectedHex {
        t.Fatalf("toHexString got %q want %q", gotHex, expectedHex)
    }

    // instance toString -> %g
    sObj2 := floatToString([]interface{}{obj}).(*object.Object)
    gotStr := object.GoStringFromStringObject(sObj2)
    if gotStr != "123.25" {
        t.Fatalf("toString got %q", gotStr)
    }

    // static toString(F) -> %f default precision (6)
    sObj3 := floatToStringStatic([]interface{}{v}).(*object.Object)
    gotStr2 := object.GoStringFromStringObject(sObj3)
    if gotStr2 != "123.250000" {
        t.Fatalf("toStringStatic got %q", gotStr2)
    }
}

func TestFloat_Compare_CompareTo_Equals(t *testing.T) {
    a := makeFloat(1.0)
    b := makeFloat(2.0)

    cmp := floatCompare([]interface{}{a, b}).(int64)
    if cmp != -1 {
        t.Fatalf("compare expected -1, got %d", cmp)
    }

    cto := floatCompareTo([]interface{}{a, b}).(int64)
    if cto != -1 {
        t.Fatalf("compareTo expected -1, got %d", cto)
    }

    eq := floatEquals([]interface{}{a, makeFloat(1.0)})
    if eq != types.JavaBoolTrue {
        t.Fatalf("equals expected true, got %v", eq)
    }
}

func TestFloat_PrimitiveConversions(t *testing.T) {
    obj := makeFloat(65.9)
    if bv := floatByteValue([]interface{}{obj}).(int64); bv != 65 { // byte cast then widen
        t.Fatalf("byteValue expected 65, got %d", bv)
    }
    if iv := floatIntValue([]interface{}{obj}).(int64); iv != 65 { // returned as int64 of int32
        t.Fatalf("intValue expected 65, got %d", iv)
    }
    if sv := floatShortValue([]interface{}{obj}).(int16); sv != 65 {
        t.Fatalf("shortValue expected 65, got %d", sv)
    }
    if lv := floatLongValue([]interface{}{obj}).(int64); lv != 65 {
        t.Fatalf("longValue expected 65, got %d", lv)
    }
    if dv := floatDoubleValue([]interface{}{obj}).(float64); dv != 65.9 {
        t.Fatalf("doubleValue expected 65.9, got %v", dv)
    }
}

func TestFloat_Max_Min_Sum(t *testing.T) {
    if mx := floatMax([]interface{}{1.0, 2.5}).(float64); mx != 2.5 {
        t.Fatalf("max expected 2.5, got %v", mx)
    }
    if mn := floatMin([]interface{}{1.0, -2.5}).(float64); mn != -2.5 {
        t.Fatalf("min expected -2.5, got %v", mn)
    }
    if sm := floatSum([]interface{}{1.25, 2.75}).(float64); sm != 4.0 {
        t.Fatalf("sum expected 4.0, got %v", sm)
    }
}

func TestFloat_IsInfinite_IsFinite(t *testing.T) {
    inf := makeFloat(math.Inf(1))
    if v := floatIsInfinite([]interface{}{inf}); v != types.JavaBoolTrue {
        t.Fatalf("isInfinite expected true, got %v", v)
    }
    if v := floatIsFinite([]interface{}{inf}); v != types.JavaBoolFalse {
        t.Fatalf("isFinite expected false for Inf, got %v", v)
    }
    finite := makeFloat(0.0)
    if v := floatIsFinite([]interface{}{finite}); v != types.JavaBoolTrue {
        t.Fatalf("isFinite expected true for 0.0, got %v", v)
    }
}

func TestFloat_ValueOfString(t *testing.T) {
    s := object.StringObjectFromGoString("-10.5")
    out := floatValueOfString([]interface{}{s})
    fobj, ok := out.(*object.Object)
    if !ok {
        t.Fatalf("valueOf(String) did not return object, got %T", out)
    }
    if val, ok := fobj.FieldTable["value"].Fvalue.(float64); !ok || val != -10.5 {
        t.Fatalf("valueOfString value mismatch: %v", fobj.FieldTable["value"].Fvalue)
    }
}

func TestFloat_Float16_Conversions(t *testing.T) {
    // 0x3C00 is +1.0 in IEEE 754 half precision
    one16 := int64(0x3C00)
    out := floatFloat16ToFloat([]interface{}{one16}).(float64)
    if float32(out) != float32(1.0) {
        t.Fatalf("float16ToFloat for 1.0 failed: got %v", out)
    }

    // Convert 1.0f32 back to half
    back := floatFloatToFloat16([]interface{}{1.0}).(int64)
    if uint16(back) != 0x3C00 {
        t.Fatalf("floatToFloat16 for 1.0 failed: got 0x%04X", uint16(back))
    }

    // Infinity mapping
    inf16 := int64(0x7C00)
    outInf := floatFloat16ToFloat([]interface{}{inf16}).(float64)
    if !math.IsInf(outInf, 1) {
        t.Fatalf("float16ToFloat for +Inf failed: got %v", outInf)
    }
}
