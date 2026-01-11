package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
	"testing"
)

func TestFloat_Coverage(t *testing.T) {
	globals.InitGlobals("test")

	f1 := makeFloat(1.0)
	f2 := makeFloat(2.0)

	// floatCompare error paths
	if blk, ok := floatCompare([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatCompare(1) expected IllegalArgumentException")
	}
	if blk, ok := floatCompare([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatCompare(1, string) expected IllegalArgumentException")
	}

	// floatCompare branches
	if res := floatCompare([]interface{}{2.0, 1.0}).(int64); res != 1 {
		t.Errorf("floatCompare(2, 1) expected 1, got %d", res)
	}
	fZero := 0.0
	fNegZero := math.Copysign(0.0, -1.0)
	res0 := floatCompare([]interface{}{fZero, fNegZero}).(int64)
	if res0 != 1 {
		t.Errorf("floatCompare(0, -0) expected 1, got %d", res0)
	}
	resN0 := floatCompare([]interface{}{fNegZero, fZero}).(int64)
	if resN0 != -1 {
		t.Errorf("floatCompare(-0, 0) expected -1, got %d", resN0)
	}

	// floatCompareTo error paths
	if blk, ok := floatCompareTo([]interface{}{f1}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatCompareTo(1) expected IllegalArgumentException")
	}
	if blk, ok := floatCompareTo([]interface{}{1.0, 2.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatCompareTo(float, float) expected IllegalArgumentException")
	}

	// floatFloatToIntBits error paths
	if blk, ok := floatFloatToIntBits([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloatToIntBits() expected IllegalArgumentException")
	}
	if blk, ok := floatFloatToIntBits([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloatToIntBits(string) expected IllegalArgumentException")
	}
	// floatFloatToIntBits NaN normalization
	nan := math.NaN()
	bits := floatFloatToIntBits([]interface{}{nan}).(int64)
	if uint32(bits) != 0x7fc00000 {
		t.Errorf("floatFloatToIntBits(NaN) normalization failed, got %08x", uint32(bits))
	}

	// floatFloatToRawIntBits
	rawBits := floatFloatToRawIntBits([]interface{}{1.0}).(int64)
	if uint32(rawBits) != math.Float32bits(1.0) {
		t.Errorf("floatFloatToRawIntBits(1.0) failed")
	}
	// floatFloatToRawIntBits error paths
	if blk, ok := floatFloatToRawIntBits([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloatToRawIntBits() expected IllegalArgumentException")
	}
	if blk, ok := floatFloatToRawIntBits([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloatToRawIntBits(string) expected IllegalArgumentException")
	}

	// floatEquals error paths
	if blk, ok := floatEquals([]interface{}{f1}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatEquals(1) expected IllegalArgumentException")
	}
	if blk, ok := floatEquals([]interface{}{1.0, f2}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatEquals(float, obj) expected IllegalArgumentException")
	}
	if floatEquals([]interface{}{f1, object.Null}) != types.JavaBoolFalse {
		t.Errorf("floatEquals(obj, Null) expected false")
	}
	if floatEquals([]interface{}{f1, (*object.Object)(nil)}) != types.JavaBoolFalse {
		t.Errorf("floatEquals(obj, nil) expected false")
	}
	intObj := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(42))
	if floatEquals([]interface{}{f1, intObj}) != types.JavaBoolFalse {
		t.Errorf("floatEquals(Float, Integer) expected false")
	}

	// floatIsInfinite
	infObj := makeFloat(math.Inf(1))
	if floatIsInfinite([]interface{}{infObj}) != types.JavaBoolTrue {
		t.Errorf("floatIsInfinite(InfObj) expected true")
	}
	if floatIsInfinite([]interface{}{f1}) != types.JavaBoolFalse {
		t.Errorf("floatIsInfinite(1Obj) expected false")
	}

	// floatIsNaN / floatIsNaNStatic
	nanObj := makeFloat(math.NaN())
	if floatIsNaN([]interface{}{nanObj}) != types.JavaBoolTrue {
		t.Errorf("floatIsNaN(NaNObj) expected true")
	}
	if floatIsNaN([]interface{}{f1}) != types.JavaBoolFalse {
		t.Errorf("floatIsNaN(1Obj) expected false")
	}
	if floatIsNaNStatic([]interface{}{math.NaN()}) != types.JavaBoolTrue {
		t.Errorf("floatIsNaNStatic(NaN) expected true")
	}
	if floatIsNaNStatic([]interface{}{1.0}) != types.JavaBoolFalse {
		t.Errorf("floatIsNaNStatic(1) expected false")
	}

	// floatIsInfiniteStatic
	if floatIsInfiniteStatic([]interface{}{math.Inf(1)}) != types.JavaBoolTrue {
		t.Errorf("floatIsInfiniteStatic(Inf) expected true")
	}
	if floatIsInfiniteStatic([]interface{}{math.Inf(-1)}) != types.JavaBoolTrue {
		t.Errorf("floatIsInfiniteStatic(-Inf) expected true")
	}
	if floatIsInfiniteStatic([]interface{}{1.0}) != types.JavaBoolFalse {
		t.Errorf("floatIsInfiniteStatic(1) expected false")
	}

	// floatIntBitsToFloat error paths
	if blk, ok := floatIntBitsToFloat([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatIntBitsToFloat() expected IllegalArgumentException")
	}
	if blk, ok := floatIntBitsToFloat([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatIntBitsToFloat(string) expected IllegalArgumentException")
	}

	// floatMax / floatMin / floatSum error paths
	if blk, ok := floatMax([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatMax(1) expected IllegalArgumentException")
	}
	if blk, ok := floatMax([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatMax(1, string) expected IllegalArgumentException")
	}
	// floatMax actual logic
	if res := floatMax([]interface{}{1.0, 2.0}).(float64); res != 2.0 {
		t.Errorf("floatMax(1, 2) expected 2")
	}
	if res := floatMax([]interface{}{2.0, 1.0}).(float64); res != 2.0 {
		t.Errorf("floatMax(2, 1) expected 2")
	}
	if res := floatMax([]interface{}{math.NaN(), 1.0}).(float64); !math.IsNaN(res) {
		t.Errorf("floatMax(NaN, 1) expected NaN")
	}
	if res := floatMax([]interface{}{1.0, math.NaN()}).(float64); !math.IsNaN(res) {
		t.Errorf("floatMax(1, NaN) expected NaN")
	}

	if blk, ok := floatMin([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatMin(1) expected IllegalArgumentException")
	}
	if blk, ok := floatMin([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatMin(1, string) expected IllegalArgumentException")
	}
	// floatMin actual logic
	if res := floatMin([]interface{}{1.0, 2.0}).(float64); res != 1.0 {
		t.Errorf("floatMin(1, 2) expected 1")
	}
	if res := floatMin([]interface{}{2.0, 1.0}).(float64); res != 1.0 {
		t.Errorf("floatMin(2, 1) expected 1")
	}
	if res := floatMin([]interface{}{math.NaN(), 1.0}).(float64); !math.IsNaN(res) {
		t.Errorf("floatMin(NaN, 1) expected NaN")
	}
	if res := floatMin([]interface{}{1.0, math.NaN()}).(float64); !math.IsNaN(res) {
		t.Errorf("floatMin(1, NaN) expected NaN")
	}
	if blk, ok := floatSum([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatSum(1) expected IllegalArgumentException")
	}
	if blk, ok := floatSum([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatSum(1, string) expected IllegalArgumentException")
	}

	// floatParseFloat error paths
	if blk, ok := floatParseFloat([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatParseFloat() expected IllegalArgumentException")
	}
	if blk, ok := floatParseFloat([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatParseFloat(float) expected IllegalArgumentException")
	}
	sEmpty := object.StringObjectFromGoString("")
	if blk, ok := floatParseFloat([]interface{}{sEmpty}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.NullPointerException {
		t.Errorf("floatParseFloat(empty) expected NullPointerException")
	}

	// floatToHexString error paths
	if blk, ok := floatToHexString([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatToHexString() expected IllegalArgumentException")
	}
	if blk, ok := floatToHexString([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatToHexString(string) expected IllegalArgumentException")
	}

	// floatToString / floatToStringStatic
	sObj := floatToString([]interface{}{f1}).(*object.Object)
	if object.GoStringFromStringObject(sObj) != "1" {
		t.Errorf("floatToString(1) expected '1', got %s", object.GoStringFromStringObject(sObj))
	}
	// error paths
	if blk, ok := floatToStringStatic([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatToStringStatic() expected IllegalArgumentException")
	}
	if blk, ok := floatToStringStatic([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatToStringStatic(string) expected IllegalArgumentException")
	}

	// floatValueOf error paths
	if blk, ok := floatValueOf([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatValueOf() expected IllegalArgumentException")
	}
	if blk, ok := floatValueOf([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatValueOf(string) expected IllegalArgumentException")
	}

	// floatValueOfString error paths
	if blk, ok := floatValueOfString([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatValueOfString() expected IllegalArgumentException")
	}
	if blk, ok := floatValueOfString([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatValueOfString(float) expected IllegalArgumentException")
	}
	sInv := object.StringObjectFromGoString("abc")
	if blk, ok := floatValueOfString([]interface{}{sInv}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("floatValueOfString(abc) expected NumberFormatException")
	}

	// floatFloat16ToFloat / floatFloatToFloat16 error paths
	if blk, ok := floatFloat16ToFloat([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloat16ToFloat() expected IllegalArgumentException")
	}
	if blk, ok := floatFloat16ToFloat([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloat16ToFloat(float) expected IllegalArgumentException")
	}
	if blk, ok := floatFloatToFloat16([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloatToFloat16() expected IllegalArgumentException")
	}
	if blk, ok := floatFloatToFloat16([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatFloatToFloat16(string) expected IllegalArgumentException")
	}

	// floatEquals invalid self
	if blk, ok := floatEquals([]interface{}{1.0, f2}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("floatEquals(not-obj, obj) expected IllegalArgumentException")
	}

	// floatFloat16ToFloat special cases
	// Subnormal
	sub := int64(0x0001)
	floatFloat16ToFloat([]interface{}{sub})
	// 0
	zero := int64(0x0000)
	floatFloat16ToFloat([]interface{}{zero})
	// Neg Inf
	ninf := int64(0xFC00)
	floatFloat16ToFloat([]interface{}{ninf})
	// NaN
	nan16 := int64(0x7C01)
	floatFloat16ToFloat([]interface{}{nan16})

	// floatFloatToFloat16 special cases
	floatFloatToFloat16([]interface{}{0.0})
	floatFloatToFloat16([]interface{}{-0.0})
	floatFloatToFloat16([]interface{}{math.Inf(1)})
	floatFloatToFloat16([]interface{}{math.Inf(-1)})
	floatFloatToFloat16([]interface{}{math.NaN()})
	floatFloatToFloat16([]interface{}{1e30})  // overflow
	floatFloatToFloat16([]interface{}{1e-30}) // underflow
	floatFloatToFloat16([]interface{}{6e-8})  // subnormal result
}
