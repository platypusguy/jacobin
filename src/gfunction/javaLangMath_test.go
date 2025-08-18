package gfunction

import (
    "jacobin/src/excNames"
    "math"
    "math/big"
    "testing"
)

func TestMath_Abs(t *testing.T) {
    if got := absFloat64([]interface{}{float64(-3.5)}).(float64); got != 3.5 {
        t.Fatalf("absFloat64(-3.5)=%v", got)
    }
    if got := absInt64([]interface{}{int64(-42)}).(int64); got != 42 {
        t.Fatalf("absInt64(-42)=%v", got)
    }
}

func TestMath_Trigonometry_Basics(t *testing.T) {
    if got := cosFloat64([]interface{}{float64(0)}).(float64); got != 1.0 {
        t.Fatalf("cos(0)=%v", got)
    }
    if got := sinFloat64([]interface{}{float64(0)}).(float64); got != 0.0 {
        t.Fatalf("sin(0)=%v", got)
    }
    if got := tanFloat64([]interface{}{float64(0)}).(float64); got != 0.0 {
        t.Fatalf("tan(0)=%v", got)
    }
    // atan2(y=0,x=-1) == pi
    if got := atan2Float64([]interface{}{float64(0), float64(-1)}).(float64); math.Abs(got-math.Pi) > 1e-12 {
        t.Fatalf("atan2(0,-1)=%v", got)
    }
}

func TestMath_Add_Sub_Mul_High(t *testing.T) {
    if got := addExactII([]interface{}{int64(2), int64(3)}).(int64); got != 5 {
        t.Fatalf("addExactII=%v", got)
    }
    if got := subtractExactJJ([]interface{}{int64(5), int64(8)}).(int64); got != -3 {
        t.Fatalf("subtractExactJJ=%v", got)
    }
    if got := multiplyExactII([]interface{}{int64(-7), int64(6)}).(int64); got != -42 {
        t.Fatalf("multiplyExactII=%v", got)
    }
    // multiplyHigh: check against big-int computation done by implementation with two known values
    // use positive operands that fit in int64 to avoid literal overflow
    aVal := uint64(0x0123456789abcdef)
    bVal := uint64(0x1111111122222222)
    hi := multiplyHighJJ([]interface{}{int64(aVal), int64(bVal)})
    // compute expected using big integers
    a := new(big.Int).SetUint64(aVal)
    b := new(big.Int).SetUint64(bVal)
    p := new(big.Int).Mul(a, b)
    expected := new(big.Int).Rsh(p, 64).Int64()
    if hi.(int64) != expected {
        t.Fatalf("multiplyHighJJ mismatch: got %x want %x", hi.(int64), expected)
    }
}

func TestMath_FloorDiv_Mod(t *testing.T) {
    // simple positive
    if got := floorDivII([]interface{}{int64(7), int64(3)}); got.(int64) != 2 {
        t.Fatalf("floorDiv 7/3=%v", got)
    }
    if got := floorModII([]interface{}{int64(7), int64(3)}).(int64); got != 1 {
        t.Fatalf("floorMod 7,3=%v", got)
    }
    // negative dividend: floorDiv should round toward -inf
    fd := floorDivII([]interface{}{int64(-7), int64(3)})
    if fd.(int64) != -3 { // since -7/3 = -2 with trunc, floor is -3
        t.Fatalf("floorDiv -7/3=%v", fd)
    }
    if got := floorModII([]interface{}{int64(-7), int64(3)}).(int64); got != 2 { // -7 = (-3)*3 + 2
        t.Fatalf("floorMod -7,3=%v", got)
    }
    // divide by zero -> ArithmeticException
    dz := floorDivII([]interface{}{int64(1), int64(0)})
    if geb, ok := dz.(*GErrBlk); !ok || geb.ExceptionType != excNames.ArithmeticException {
        t.Fatalf("floorDiv divide-by-zero expected ArithmeticException, got %T (%v)", dz, dz)
    }
}

func TestMath_Rounding(t *testing.T) {
    if got := floorFloat64([]interface{}{float64(3.9)}).(float64); got != 3.0 {
        t.Fatalf("floor 3.9=%v", got)
    }
    if got := ceilFloat64([]interface{}{float64(3.1)}).(float64); got != 4.0 {
        t.Fatalf("ceil 3.1=%v", got)
    }
    if got := rintFloat64([]interface{}{float64(2.3)}).(float64); got != 2.0 {
        t.Fatalf("rint 2.3=%v", got)
    }
    if got := roundInt64([]interface{}{float64(-2.3)}).(int64); got != -2 {
        t.Fatalf("round -2.3=%v", got)
    }
}

func TestMath_Exponent_Log_Pow(t *testing.T) {
    if got := expFloat64([]interface{}{float64(1)}).(float64); math.Abs(got-math.E) > 1e-12 {
        t.Fatalf("exp(1)=%v", got)
    }
    if got := expm1Float64([]interface{}{float64(1)}).(float64); math.Abs(got-(math.E-1)) > 1e-12 {
        t.Fatalf("expm1(1)=%v", got)
    }
    if got := logFloat64([]interface{}{math.E}).(float64); math.Abs(got-1.0) > 1e-12 {
        t.Fatalf("log(e)=%v", got)
    }
    if got := log10Float64([]interface{}{float64(1000)}).(float64); math.Abs(got-3.0) > 1e-12 {
        t.Fatalf("log10(1000)=%v", got)
    }
    if got := log1pFloat64([]interface{}{float64(0)}).(float64); got != 0.0 {
        t.Fatalf("log1p(0)=%v", got)
    }
    if got := powFloat64([]interface{}{float64(2), float64(10)}).(float64); got != 1024.0 {
        t.Fatalf("pow(2,10)=%v", got)
    }
}

func TestMath_Max_Min(t *testing.T) {
    if got := maxII([]interface{}{int64(2), int64(5)}).(int64); got != 5 {
        t.Fatalf("maxII=%v", got)
    }
    if got := minII([]interface{}{int64(2), int64(5)}).(int64); got != 2 {
        t.Fatalf("minII=%v", got)
    }
    if got := maxDD([]interface{}{float64(3.5), float64(3.6)}).(float64); got != 3.6 {
        t.Fatalf("maxDD=%v", got)
    }
    if got := minDD([]interface{}{float64(3.5), float64(3.6)}).(float64); got != 3.5 {
        t.Fatalf("minDD=%v", got)
    }
}

func TestMath_NextAfter_Up_Down(t *testing.T) {
    base := 1.0
    up := nextUpFloat64([]interface{}{base}).(float64)
    if !(up > base) {
        t.Fatalf("nextUp not greater than base: %v <= %v", up, base)
    }
    down := nextDownFloat64([]interface{}{base}).(float64)
    if !(down < base) {
        t.Fatalf("nextDown not less than base: %v >= %v", down, base)
    }
    na := nextAfterDD([]interface{}{float64(0), float64(1)}).(float64)
    if !(na > 0) {
        t.Fatalf("nextAfter(0->1) not > 0: %v", na)
    }
}

func TestMath_FMA_Hypot_Remainder(t *testing.T) {
    if got := fmaDDD([]interface{}{float64(2), float64(3), float64(4)}).(float64); got != 10.0 {
        t.Fatalf("fma(2,3,4)=%v", got)
    }
    if got := hypotFloat64([]interface{}{float64(3), float64(4)}).(float64); math.Abs(got-5.0) > 1e-12 {
        t.Fatalf("hypot(3,4)=%v", got)
    }
    if got := IEEEremainderFloat64([]interface{}{float64(5), float64(2)}).(float64); math.Abs(got-1.0) > 1e-12 {
        t.Fatalf("Remainder(5,2)=%v", got)
    }
}

func TestMath_Scalb_Signum_Sqrt(t *testing.T) {
    if got := scalbDI([]interface{}{float64(1.5), int64(1)}).(float64); math.Abs(got-3.0) > 1e-12 {
        t.Fatalf("scalb(1.5,1)=%v", got)
    }
    if got := signumFloat64([]interface{}{float64(-0.1)}).(float64); got != -1.0 {
        t.Fatalf("signum(-0.1)=%v", got)
    }
    if got := signumFloat64([]interface{}{float64(0.0)}).(float64); got != 0.0 {
        t.Fatalf("signum(0)=%v", got)
    }
    if got := sqrtFloat64([]interface{}{float64(144)}).(float64); got != 12.0 {
        t.Fatalf("sqrt(144)=%v", got)
    }
}

func TestMath_Degree_Radian_and_ToIntExact(t *testing.T) {
    // 180 deg -> pi radians; using PI constant in code
    if got := toRadiansFloat64([]interface{}{float64(180.0)}).(float64); math.Abs(got-math.Pi) > 1e-12 {
        t.Fatalf("toRadians(180)=%v", got)
    }
    if got := toDegreesFloat64([]interface{}{math.Pi}).(float64); math.Abs(got-180.0) > 1e-12 {
        t.Fatalf("toDegrees(pi)=%v", got)
    }
    if got := toIntExactInt64([]interface{}{int64(-123)}).(int64); got != -123 {
        t.Fatalf("toIntExact(-123)=%v", got)
    }
}

func TestMath_GetExponent_Ulp_Specials(t *testing.T) {
    // Normal number: exponent of 1.0 is 0
    {
            v := getExponentFloat64([]interface{}{float64(1.0)})
            var got int64
            switch x := v.(type) {
            case int64:
                got = x
            case int:
                got = int64(x)
            default:
                t.Fatalf("unexpected type for exponent: %T", v)
            }
            if got != 0 {
                t.Fatalf("getExponent(1.0)=%v", got)
            }
        }
    // NaN / Inf: returns MAX_DOUBLE_EXPONENT+1 == 1024
    {
        v := getExponentFloat64([]interface{}{math.NaN()})
        var got int64
        switch x := v.(type) {
        case int64:
            got = x
        case int:
            got = int64(x)
        default:
            t.Fatalf("unexpected type for exponent NaN: %T", v)
        }
        if got != MAX_DOUBLE_EXPONENT+1 {
            t.Fatalf("getExponent(NaN)=%v", got)
        }
    }
    {
        v := getExponentFloat64([]interface{}{math.Inf(1)})
        var got int64
        switch x := v.(type) {
        case int64:
            got = x
        case int:
            got = int64(x)
        default:
            t.Fatalf("unexpected type for exponent Inf: %T", v)
        }
        if got != MAX_DOUBLE_EXPONENT+1 {
            t.Fatalf("getExponent(Inf)=%v", got)
        }
    }

    // ULP: finite positive value should be positive and small
    if got := ulpFloat64([]interface{}{float64(1.0)}).(float64); !(got > 0 && got < 1e-15) {
        t.Fatalf("ulp(1.0)=%v", got)
    }
    // ULP of +/-Inf is +Inf
    if res := ulpFloat64([]interface{}{math.Inf(-1)}).(float64); !math.IsInf(res, 1) {
        t.Fatalf("ulp(-Inf) expected +Inf, got %v", res)
    }
    // NaN returns NaN
    if res := ulpFloat64([]interface{}{math.NaN()}).(float64); !math.IsNaN(res) {
        t.Fatalf("ulp(NaN) expected NaN, got %v", res)
    }
}

func TestMath_Random_Range(t *testing.T) {
    // random in [0,1)
    v := randomFloat64(nil).(float64)
    if !(v >= 0.0 && v < 1.0) {
        t.Fatalf("random out of range: %v", v)
    }
}
