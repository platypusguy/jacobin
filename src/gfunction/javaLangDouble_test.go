package gfunction

import (
    "jacobin/src/object"
    "jacobin/src/types"
    "jacobin/src/excNames"
    "math"
    "testing"
    "fmt"
)

// helper to create a java/lang/Double object with a given value
func makeDouble(val float64) *object.Object {
    return Populator("java/lang/Double", types.Double, val)
}

func TestDouble_ValueOf_And_DoubleValue(t *testing.T) {
    // valueOf(D) -> Double, then doubleValue()D -> primitive
    d := 42.5
    obj := doubleValueOf([]interface{}{d})
    dobj, ok := obj.(*object.Object)
    if !ok {
        t.Fatalf("valueOf did not return object, got %T", obj)
    }
    out := doubleDoubleValue([]interface{}{dobj})
    if got := out.(float64); got != d {
        t.Fatalf("doubleDoubleValue mismatch: got %v want %v", got, d)
    }
}

func TestDouble_ParseDouble_Valid_Invalid_Empty(t *testing.T) {
    // valid
    s := object.StringObjectFromGoString("3.5")
    out := doubleParseDouble([]interface{}{s})
    if got := out.(float64); got != 3.5 {
        t.Fatalf("parseDouble valid: got %v", got)
    }
    // invalid -> NumberFormatException
    sInv := object.StringObjectFromGoString("abc")
    out = doubleParseDouble([]interface{}{sInv})
    if geb, ok := out.(*GErrBlk); !ok || geb.ExceptionType != excNames.NumberFormatException {
        if !ok {
            t.Fatalf("parseDouble invalid: expected *GErrBlk, got %T", out)
        }
        t.Fatalf("parseDouble invalid: expected NumberFormatException, got %v", geb)
    }
    // empty -> NullPointerException according to implementation
    sEmpty := object.StringObjectFromGoString("")
    out = doubleParseDouble([]interface{}{sEmpty})
    if geb, ok := out.(*GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
        t.Fatalf("parseDouble empty: expected NullPointerException, got %T (%v)", out, out)
    }
}

func TestDouble_ToLongBits_And_LongBitsToDouble(t *testing.T) {
    // normal value round-trip
    val := -123.75
    bits := doubleToLongBits([]interface{}{val}).(int64)
    back := doubleLongBitsToDouble([]interface{}{bits}).(float64)
    if math.Float64bits(back) != math.Float64bits(val) {
        t.Fatalf("round-trip bits mismatch: got %x want %x", math.Float64bits(back), math.Float64bits(val))
    }

    // NaN should map to canonical NaN pattern 0x7FF8000000000000 per implementation
    nanBits := doubleToLongBits([]interface{}{math.NaN()}).(int64)
    if uint64(nanBits) != 0x7FF8000000000000 {
        t.Fatalf("NaN bits mismatch: got 0x%016X", uint64(nanBits))
    }

    // Check signed zero via longBitsToDouble
    var one uint64 = 1
    negZeroBits := int64(one << 63)
    negZero := doubleLongBitsToDouble([]interface{}{negZeroBits}).(float64)
    if !math.Signbit(negZero) {
        t.Fatalf("expected negative zero sign bit")
    }
}

func TestDouble_ToHexString(t *testing.T) {
    d := 1.5
    obj := makeDouble(d)
    out := doubleToHexString([]interface{}{obj})
    sObj := out.(*object.Object)
    got := object.GoStringFromStringObject(sObj)
    expected := fmt.Sprintf("0x%016X", math.Float64bits(d))
    if got != expected {
        t.Fatalf("toHexString got %q want %q", got, expected)
    }
}

func TestDouble_ToString_Variants(t *testing.T) {
    obj := makeDouble(123.25)
    out := doubleToString([]interface{}{obj}).(*object.Object)
    got := object.GoStringFromStringObject(out)
    // %g formatting should produce "123.25" for this value
    if got != "123.25" {
        t.Fatalf("toString got %q", got)
    }
    // static variant uses %f default precision (6)
    out2 := doubleToStringStatic([]interface{}{123.25}).(*object.Object)
    got2 := object.GoStringFromStringObject(out2)
    if got2 != "123.250000" {
        t.Fatalf("toStringStatic got %q", got2)
    }
}

func TestDouble_Compare_CompareTo_Equals(t *testing.T) {
    a := makeDouble(1.0)
    b := makeDouble(2.0)

    cmp := doubleCompare([]interface{}{a, b}).(int64)
    if cmp != -1 {
        t.Fatalf("compare expected -1, got %d", cmp)
    }

    cto := doubleCompareTo([]interface{}{a, b}).(int64)
    if cto != -1 {
        t.Fatalf("compareTo expected -1, got %d", cto)
    }

    eq := doubleEquals([]interface{}{a, makeDouble(1.0)})
    if eq != types.JavaBoolTrue {
        t.Fatalf("equals expected true, got %v", eq)
    }
}

func TestDouble_PrimitiveConversions(t *testing.T) {
    obj := makeDouble(65.9)
    if bv := doubleByteValue([]interface{}{obj}).(int64); bv != 65 { // byte cast then widen to int64
        t.Fatalf("byteValue expected 65, got %d", bv)
    }
    if iv := doubleIntValue([]interface{}{obj}).(int32); iv != 65 {
        t.Fatalf("intValue expected 65, got %d", iv)
    }
    if sv := doubleShortValue([]interface{}{obj}).(int16); sv != 65 {
        t.Fatalf("shortValue expected 65, got %d", sv)
    }
    if lv := doubleLongValue([]interface{}{obj}).(int64); lv != 65 {
        t.Fatalf("longValue expected 65, got %d", lv)
    }
    if fv := doubleFloatValue([]interface{}{obj}).(float64); fv != 65.9 {
        t.Fatalf("floatValue expected 65.9, got %v", fv)
    }
}

func TestDouble_Max_Min_Sum(t *testing.T) {
    if mx := doubleMax([]interface{}{1.0, 2.5}).(float64); mx != 2.5 {
        t.Fatalf("max expected 2.5, got %v", mx)
    }
    if mn := doubleMin([]interface{}{1.0, -2.5}).(float64); mn != -2.5 {
        t.Fatalf("min expected -2.5, got %v", mn)
    }
    if sm := doubleSum([]interface{}{1.25, 2.75}).(float64); sm != 4.0 {
        t.Fatalf("sum expected 4.0, got %v", sm)
    }
}

func TestDouble_IsInfinite_IsFinite(t *testing.T) {
    inf := makeDouble(math.Inf(1))
    if v := doubleIsInfinite([]interface{}{inf}); v != types.JavaBoolTrue {
        t.Fatalf("isInfinite expected true, got %v", v)
    }
    if v := doubleIsFinite([]interface{}{inf}); v != types.JavaBoolFalse {
        t.Fatalf("isFinite expected false for Inf, got %v", v)
    }
    finite := makeDouble(0.0)
    if v := doubleIsFinite([]interface{}{finite}); v != types.JavaBoolTrue {
        t.Fatalf("isFinite expected true for 0.0, got %v", v)
    }
}

func TestDouble_ValueOfString(t *testing.T) {
    s := object.StringObjectFromGoString("-10.5")
    out := doubleValueOfString([]interface{}{s})
    dobj, ok := out.(*object.Object)
    if !ok {
        t.Fatalf("valueOf(String) did not return object, got %T", out)
    }
    if val := doubleDoubleValue([]interface{}{dobj}).(float64); val != -10.5 {
        t.Fatalf("valueOf(String) value mismatch: %v", val)
    }
}
