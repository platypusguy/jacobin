package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"testing"
)

// helper to create a java/lang/Double object with a given value
func makeDouble(val float64) *object.Object {
	return object.MakePrimitiveObject("java/lang/Double", types.Double, val)
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
	if geb, ok := out.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NumberFormatException {
		if !ok {
			t.Fatalf("parseDouble invalid: expected *ghelpers.GErrBlk, got %T", out)
		}
		t.Fatalf("parseDouble invalid: expected NumberFormatException, got %v", geb)
	}
	// empty -> NullPointerException according to implementation
	sEmpty := object.StringObjectFromGoString("")
	out = doubleParseDouble([]interface{}{sEmpty})
	if geb, ok := out.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
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

	// doubleToLongBits should normalize NaN
	nanVal := math.Float64frombits(0x7ff0000000000001) // a non-canonical NaN
	nanBits := doubleToLongBits([]interface{}{nanVal}).(int64)
	if uint64(nanBits) != 0x7ff8000000000000 {
		t.Fatalf("doubleToLongBits NaN normalization mismatch: got 0x%016X", uint64(nanBits))
	}

	// doubleToRawLongBits should NOT normalize NaN
	rawNanBits := doubleToRawLongBits([]interface{}{nanVal}).(int64)
	if uint64(rawNanBits) != 0x7ff0000000000001 {
		t.Fatalf("doubleToRawLongBits NaN normalization mismatch: got 0x%016X", uint64(rawNanBits))
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
	out := doubleToHexString([]interface{}{d})
	sObj := out.(*object.Object)
	got := object.GoStringFromStringObject(sObj)
	// Java's toHexString can be "0x1.8p0" or "0x1.8p+0" or "0x1.8p+00" depending on platform/version
	// Go's strconv.FormatFloat(1.5, 'x', -1, 64) produces "0x1.8p+00"
	if got != "0x1.8p+00" {
		t.Fatalf("toHexString got %q", got)
	}
}

func TestDouble_ToString_Variants(t *testing.T) {
	obj := makeDouble(123.25)
	out := doubleToString([]interface{}{obj}).(*object.Object)
	got := object.GoStringFromStringObject(out)
	if got != "123.25" {
		t.Fatalf("toString got %q", got)
	}
	// static variant
	out2 := doubleToStringStatic([]interface{}{123.25}).(*object.Object)
	got2 := object.GoStringFromStringObject(out2)
	if got2 != "123.25" {
		t.Fatalf("toStringStatic got %q", got2)
	}
}

func TestDouble_Compare_CompareTo_Equals(t *testing.T) {
	// Basic comparisons
	if doubleCompare([]interface{}{1.0, 2.0}).(int64) != -1 {
		t.Fail()
	}
	if doubleCompare([]interface{}{2.0, 1.0}).(int64) != 1 {
		t.Fail()
	}
	if doubleCompare([]interface{}{1.0, 1.0}).(int64) != 0 {
		t.Fail()
	}

	// NaN comparisons: NaN is equal to itself and greater than any other value (including +Inf)
	nan := math.NaN()
	inf := math.Inf(1)
	if doubleCompare([]interface{}{nan, nan}).(int64) != 0 {
		t.Errorf("NaN should equal NaN")
	}
	if doubleCompare([]interface{}{nan, inf}).(int64) != 1 {
		t.Errorf("NaN should be greater than Inf")
	}
	if doubleCompare([]interface{}{inf, nan}).(int64) != -1 {
		t.Errorf("Inf should be less than NaN")
	}

	// Zero comparisons: +0.0 is greater than -0.0
	pz := 0.0
	nz := math.Float64frombits(0x8000000000000000)
	if doubleCompare([]interface{}{pz, nz}).(int64) != 1 {
		t.Errorf("+0.0 should be greater than -0.0")
	}
	if doubleCompare([]interface{}{nz, pz}).(int64) != -1 {
		t.Errorf("-0.0 should be less than +0.0")
	}

	// Equals: matches bit pattern (like doubleCompare == 0)
	a := makeDouble(nan)
	if doubleEquals([]interface{}{a, makeDouble(nan)}) != types.JavaBoolTrue {
		t.Errorf("equals(NaN, NaN) should be true")
	}
	if doubleEquals([]interface{}{makeDouble(pz), makeDouble(nz)}) != types.JavaBoolFalse {
		t.Errorf("equals(+0.0, -0.0) should be false")
	}
}

func TestDouble_HashCode(t *testing.T) {
	v := 1.2345
	bits := math.Float64bits(v)
	expected := int64(int32(bits ^ (bits >> 32)))

	h1 := doubleHashCodeStatic([]interface{}{v}).(int64)
	if h1 != expected {
		t.Fatalf("hashCodeStatic mismatch: got %d want %d", h1, expected)
	}

	dobj := makeDouble(v)
	h2 := doubleHashCode([]interface{}{dobj}).(int64)
	if h2 != expected {
		t.Fatalf("hashCode mismatch: got %d want %d", h2, expected)
	}
}

func TestDouble_IsNaN_Static_And_Instance(t *testing.T) {
	nan := math.NaN()
	if doubleIsNaNStatic([]interface{}{nan}) != types.JavaBoolTrue {
		t.Fail()
	}
	if doubleIsNaNStatic([]interface{}{1.0}) != types.JavaBoolFalse {
		t.Fail()
	}

	if doubleIsNaN([]interface{}{makeDouble(nan)}) != types.JavaBoolTrue {
		t.Fail()
	}
	if doubleIsNaN([]interface{}{makeDouble(1.0)}) != types.JavaBoolFalse {
		t.Fail()
	}
}

func TestDouble_PrimitiveConversions(t *testing.T) {
	obj := makeDouble(65.9)
	if bv := doubleByteValue([]interface{}{obj}).(int64); bv != 65 { // byte cast then widen to int64
		t.Fatalf("byteValue expected 65, got %d", bv)
	}
	if iv := doubleIntValue([]interface{}{obj}).(int64); iv != 65 {
		t.Fatalf("intValue expected 65, got %d", iv)
	}
	if sv := doubleShortValue([]interface{}{obj}).(int64); sv != 65 {
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
	if v := doubleIsFiniteStatic([]interface{}{math.Inf(1)}); v != types.JavaBoolFalse {
		t.Fatalf("isFinite expected false for Inf, got %v", v)
	}
	if v := doubleIsFiniteStatic([]interface{}{0.0}); v != types.JavaBoolTrue {
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
