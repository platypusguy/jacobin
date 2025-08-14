package gfunction

import (
    "jacobin/excNames"
    "jacobin/globals"
    "jacobin/object"
    "jacobin/types"
    "testing"
)

func newRandomObj() *object.Object {
    cn := "java/util/Random"
    return object.MakeEmptyObjectWithClassName(&cn)
}

func assertJavaBoolVal(t *testing.T, v interface{}) {
    t.Helper()
    if vi, ok := v.(int64); !ok || (vi != types.JavaBoolTrue && vi != types.JavaBoolFalse) {
        t.Fatalf("expected Java boolean int64 0/1, got %T (%v)", v, v)
    }
}

func TestRandom_Init_And_SetSeed_Determinism(t *testing.T) {
    globals.InitStringPool()

    // Two instances with the same seed should produce the same sequences
    r1 := newRandomObj()
    r2 := newRandomObj()
    _ = randomInitLong([]interface{}{r1, int64(12345)})
    _ = randomInitLong([]interface{}{r2, int64(12345)})

    // nextInt/nextLong/nextFloat/nextDouble/nextBoolean equality across instances
    if a, b := randomNextInt([]interface{}{r1}).(int64), randomNextInt([]interface{}{r2}).(int64); a != b { t.Fatalf("nextInt mismatch: %d vs %d", a, b) }
    if a, b := randomNextLong([]interface{}{r1}).(int64), randomNextLong([]interface{}{r2}).(int64); a != b { t.Fatalf("nextLong mismatch: %d vs %d", a, b) }
    if a, b := randomNextFloat([]interface{}{r1}).(float64), randomNextFloat([]interface{}{r2}).(float64); a != b { t.Fatalf("nextFloat mismatch: %v vs %v", a, b) }
    if a, b := randomNextDouble([]interface{}{r1}).(float64), randomNextDouble([]interface{}{r2}).(float64); a != b { t.Fatalf("nextDouble mismatch: %v vs %v", a, b) }
    if a, b := randomNextBoolean([]interface{}{r1}).(int64), randomNextBoolean([]interface{}{r2}).(int64); a != b { t.Fatalf("nextBoolean mismatch: %d vs %d", a, b) }

    // Gaussian: call twice and compare sequences
    g1a := randomNextGaussian([]interface{}{r1}).(float64)
    g2a := randomNextGaussian([]interface{}{r2}).(float64)
    if g1a != g2a { t.Fatalf("nextGaussian first mismatch: %v vs %v", g1a, g2a) }
    g1b := randomNextGaussian([]interface{}{r1}).(float64)
    g2b := randomNextGaussian([]interface{}{r2}).(float64)
    if g1b != g2b { t.Fatalf("nextGaussian second mismatch: %v vs %v", g1b, g2b) }

    // setSeed should reset stream deterministically
    _ = randomSetSeed([]interface{}{r1, int64(777)})
    _ = randomSetSeed([]interface{}{r2, int64(777)})
    if a, b := randomNextInt([]interface{}{r1}).(int64), randomNextInt([]interface{}{r2}).(int64); a != b { t.Fatalf("after setSeed, nextInt mismatch: %d vs %d", a, b) }
}

func TestRandom_Ranges_And_Types(t *testing.T) {
    globals.InitStringPool()

    r := newRandomObj()
    _ = randomInitVoid([]interface{}{r})

    // nextInt returns int64 in [0, 2^31)
    vi := randomNextInt([]interface{}{r})
    if _, ok := vi.(int64); !ok {
        t.Fatalf("nextInt did not return int64, got %T", vi)
    }
    // nextLong returns int64 in [0, 2^63)
    vl := randomNextLong([]interface{}{r})
    if _, ok := vl.(int64); !ok {
        t.Fatalf("nextLong did not return int64, got %T", vl)
    }

    // nextFloat/nextDouble return float64 in [0,1)
    vf := randomNextFloat([]interface{}{r}).(float64)
    if !(vf >= 0.0 && vf < 1.0) { t.Fatalf("nextFloat out of range: %v", vf) }
    vd := randomNextDouble([]interface{}{r}).(float64)
    if !(vd >= 0.0 && vd < 1.0) { t.Fatalf("nextDouble out of range: %v", vd) }

    // nextBoolean returns Java boolean 0/1 as int64
    vb := randomNextBoolean([]interface{}{r})
    assertJavaBoolVal(t, vb)

    // nextInt(bound) with invalid bound -> error
    if err := randomNextIntBound([]interface{}{r, int64(0)}); err == nil {
        t.Fatalf("expected error for non-positive bound")
    } else {
        if geb, ok := err.(*GErrBlk); ok {
            if geb.ExceptionType != excNames.IllegalArgumentException { t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType) }
        }
    }

    // nextInt(bound) with positive bound in [0, bound)
    for _, b := range []int64{1, 2, 17, 100} {
        x := randomNextIntBound([]interface{}{r, b}).(int64)
        if x < 0 || x >= b {
            t.Fatalf("nextInt(%d) out of range: %d", b, x)
        }
    }
}

func TestRandom_NextBytes_FillsArray(t *testing.T) {
    globals.InitStringPool()

    r := newRandomObj()
    _ = randomInitVoid([]interface{}{r})

    // Prepare an array of 16 zero bytes in String-object wrapper
    arr := make([]byte, 16)
    jb := object.JavaByteArrayFromGoByteArray(arr)
    arrObj := object.StringObjectFromJavaByteArray(jb)

    if ret := randomNextBytes([]interface{}{r, arrObj}); ret != nil {
        t.Fatalf("nextBytes returned error: %v", ret)
    }

    outJB := object.JavaByteArrayFromStringObject(arrObj)
    if len(outJB) != 16 { t.Fatalf("nextBytes length mismatch: %d", len(outJB)) }

    // Ensure not all zeros
    allZero := true
    for _, b := range outJB { if b != 0 { allZero = false; break } }
    if allZero { t.Fatalf("nextBytes produced all zeros (unlikely)") }
}
