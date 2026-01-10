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

func TestDouble_Coverage(t *testing.T) {
	globals.InitGlobals("test")

	// doubleCompareTo
	d1 := makeDouble(1.0)
	d2 := makeDouble(2.0)
	if res := doubleCompareTo([]interface{}{d1, d2}).(int64); res != -1 {
		t.Errorf("doubleCompareTo(1, 2) expected -1, got %d", res)
	}

	// doubleCompareTo error paths
	if blk, ok := doubleCompareTo([]interface{}{d1}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleCompareTo(1) expected IllegalArgumentException")
	}
	if blk, ok := doubleCompareTo([]interface{}{1.0, 2.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleCompareTo(float, float) expected IllegalArgumentException")
	}

	// doubleIsInfiniteStatic
	if doubleIsInfiniteStatic([]interface{}{math.Inf(1)}) != types.JavaBoolTrue {
		t.Errorf("doubleIsInfiniteStatic(Inf) expected true")
	}
	if doubleIsInfiniteStatic([]interface{}{1.0}) != types.JavaBoolFalse {
		t.Errorf("doubleIsInfiniteStatic(1) expected false")
	}

	// doubleToLongBits / doubleToRawLongBits error paths
	if blk, ok := doubleToLongBits([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToLongBits() expected IllegalArgumentException")
	}
	if blk, ok := doubleToLongBits([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToLongBits(string) expected IllegalArgumentException")
	}
	if blk, ok := doubleToRawLongBits([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToRawLongBits() expected IllegalArgumentException")
	}
	if blk, ok := doubleToRawLongBits([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToRawLongBits(string) expected IllegalArgumentException")
	}

	// doubleEquals error paths
	if blk, ok := doubleEquals([]interface{}{d1}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleEquals(1) expected IllegalArgumentException")
	}
	if blk, ok := doubleEquals([]interface{}{1.0, d2}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleEquals(float, obj) expected IllegalArgumentException")
	}
	if doubleEquals([]interface{}{d1, object.Null}) != types.JavaBoolFalse {
		t.Errorf("doubleEquals(obj, Null) expected false")
	}
	// Test with nil object (not just object.Null)
	if doubleEquals([]interface{}{d1, (*object.Object)(nil)}) != types.JavaBoolFalse {
		t.Errorf("doubleEquals(obj, nil) expected false")
	}

	classNameInt := "java/lang/Integer"
	intObj := object.MakePrimitiveObject(classNameInt, types.Int, int64(42))
	if doubleEquals([]interface{}{d1, intObj}) != types.JavaBoolFalse {
		t.Errorf("doubleEquals(Double, Integer) expected false")
	}

	// doubleLongBitsToDouble error paths
	if blk, ok := doubleLongBitsToDouble([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleLongBitsToDouble() expected IllegalArgumentException")
	}
	if blk, ok := doubleLongBitsToDouble([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleLongBitsToDouble(float) expected IllegalArgumentException")
	}

	// doubleMax / doubleMin / doubleSum error paths
	if blk, ok := doubleMax([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleMax(1) expected IllegalArgumentException")
	}
	if blk, ok := doubleMax([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleMax(1, string) expected IllegalArgumentException")
	}
	if blk, ok := doubleMin([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleMin(1) expected IllegalArgumentException")
	}
	if blk, ok := doubleMin([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleMin(1, string) expected IllegalArgumentException")
	}
	if blk, ok := doubleSum([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleSum(1) expected IllegalArgumentException")
	}
	if blk, ok := doubleSum([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleSum(1, string) expected IllegalArgumentException")
	}

	// doubleParseDouble error paths
	if blk, ok := doubleParseDouble([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleParseDouble() expected IllegalArgumentException")
	}
	if blk, ok := doubleParseDouble([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleParseDouble(float) expected IllegalArgumentException")
	}

	// doubleToHexString error paths
	if blk, ok := doubleToHexString([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToHexString() expected IllegalArgumentException")
	}
	if blk, ok := doubleToHexString([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToHexString(string) expected IllegalArgumentException")
	}

	// doubleToStringStatic error paths
	if blk, ok := doubleToStringStatic([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToStringStatic() expected IllegalArgumentException")
	}
	if blk, ok := doubleToStringStatic([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleToStringStatic(string) expected IllegalArgumentException")
	}

	// doubleValueOf error paths
	if blk, ok := doubleValueOf([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleValueOf() expected IllegalArgumentException")
	}
	if blk, ok := doubleValueOf([]interface{}{"abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleValueOf(string) expected IllegalArgumentException")
	}

	// doubleValueOfString error paths
	if blk, ok := doubleValueOfString([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleValueOfString() expected IllegalArgumentException")
	}
	if blk, ok := doubleValueOfString([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleValueOfString(float) expected IllegalArgumentException")
	}
	sInv := object.StringObjectFromGoString("abc")
	if blk, ok := doubleValueOfString([]interface{}{sInv}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.NumberFormatException {
		t.Errorf("doubleValueOfString(abc) expected NumberFormatException")
	}

	// doubleDoubleValue error paths
	if blk, ok := doubleDoubleValue([]interface{}{}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleDoubleValue() expected IllegalArgumentException")
	}
	if blk, ok := doubleDoubleValue([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleDoubleValue(float) expected IllegalArgumentException")
	}

	// Corrupted object for doubleDoubleValue
	objCorr := object.MakeEmptyObject()
	if blk, ok := doubleDoubleValue([]interface{}{objCorr}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleDoubleValue(corrupted) expected IllegalArgumentException")
	}

	// getFloat64ValueFromObject corrupted
	if _, ok := getFloat64ValueFromObject(objCorr); ok {
		t.Errorf("getFloat64ValueFromObject(corrupted) should return false")
	}
	objCorr.FieldTable = map[string]object.Field{"value": {Ftype: types.Double, Fvalue: "not a float"}}
	if _, ok := getFloat64ValueFromObject(objCorr); ok {
		t.Errorf("getFloat64ValueFromObject(bad type) should return false")
	}

	// doubleCompare error paths
	if blk, ok := doubleCompare([]interface{}{1.0}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleCompare(1) expected IllegalArgumentException")
	}
	if blk, ok := doubleCompare([]interface{}{1.0, "abc"}).(*ghelpers.GErrBlk); !ok || blk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("doubleCompare(1, string) expected IllegalArgumentException")
	}

	// doubleCompare more branches
	fZero := 0.0
	fNegZero := math.Copysign(0.0, -1.0)
	if res := doubleCompare([]interface{}{fZero, fNegZero}).(int64); res != 1 {
		t.Errorf("doubleCompare(0.0, -0.0) expected 1, got %d", res)
	}
	if res := doubleCompare([]interface{}{fNegZero, fZero}).(int64); res != -1 {
		t.Errorf("doubleCompare(-0.0, 0.0) expected -1, got %d", res)
	}
	if res := doubleCompare([]interface{}{math.NaN(), 1.0}).(int64); res != 1 {
		t.Errorf("doubleCompare(NaN, 1) expected 1")
	}
	if res := doubleCompare([]interface{}{1.0, math.NaN()}).(int64); res != -1 {
		t.Errorf("doubleCompare(1, NaN) expected -1")
	}
	if res := doubleCompare([]interface{}{math.NaN(), math.NaN()}).(int64); res != 0 {
		t.Errorf("doubleCompare(NaN, NaN) expected 0")
	}

	// doubleIsInfinite
	infObj := makeDouble(math.Inf(1))
	if doubleIsInfinite([]interface{}{infObj}) != types.JavaBoolTrue {
		t.Errorf("doubleIsInfinite(InfObj) expected true")
	}
	if doubleIsInfinite([]interface{}{d1}) != types.JavaBoolFalse {
		t.Errorf("doubleIsInfinite(1Obj) expected false")
	}
}
