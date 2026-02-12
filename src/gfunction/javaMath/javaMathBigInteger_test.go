package javaMath

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
	"testing"
)

// Helpers
func biFromInt64(v int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ghelpers.InitBigIntegerField(obj, v)
	return obj
}

func bigIntOf(obj *object.Object) *big.Int {
	return obj.FieldTable["value"].Fvalue.(*big.Int)
}

func asString(obj interface{}) string {
	return object.GoStringFromStringObject(obj.(*object.Object))
}

func asJavaBytesFromStringObject(obj *object.Object) []types.JavaByte {
	return object.JavaByteArrayFromStringObject(obj)
}

func TestBigInteger_ValueOf_And_ToString(t *testing.T) {
	globals.InitStringPool()

	// valueOf
	bi := bigIntegerValueOf([]interface{}{int64(-12345)}).(*object.Object)
	if bigIntOf(bi).Int64() != -12345 {
		t.Fatalf("valueOf mismatch: expected -12345, got %d", bigIntOf(bi).Int64())
	}

	// toString
	s := BigIntegerToString([]interface{}{bi}).(*object.Object)
	if asString(s) != "-12345" {
		t.Fatalf("toString mismatch: expected -12345, got %q", asString(s))
	}
}

func TestBigInteger_InitString_And_Radix(t *testing.T) {
	globals.InitStringPool()

	// valid base-10
	base := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ret := bigIntegerInitString([]interface{}{base, object.StringObjectFromGoString("98765")})
	if ret != nil {
		t.Fatalf("unexpected error initializing from string: %v", ret)
	}
	if bigIntOf(base).String() != "98765" {
		t.Fatalf("init string mismatch: expected 98765, got %s", bigIntOf(base).String())
	}

	// invalid base-10
	base = object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ret = bigIntegerInitString([]interface{}{base, object.StringObjectFromGoString("12AB")})
	if ret == nil {
		t.Fatalf("expected NumberFormatException for invalid decimal string")
	}

	// valid radix 16
	base = object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ret = bigIntegerInitStringRadix([]interface{}{base, object.StringObjectFromGoString("1a"), int64(16)})
	if ret != nil {
		t.Fatalf("unexpected error initializing from hex string: %v", ret)
	}
	if bigIntOf(base).Int64() != 26 {
		t.Fatalf("expected 26 from hex 1a, got %d", bigIntOf(base).Int64())
	}

	// invalid radix parse
	base = object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ret = bigIntegerInitStringRadix([]interface{}{base, object.StringObjectFromGoString("12Z"), int64(10)})
	if ret == nil {
		t.Fatalf("expected NumberFormatException for invalid radix parse")
	}
}

func TestBigInteger_Arithmetic_Add_Sub_Mul(t *testing.T) {
	globals.InitStringPool()

	a := biFromInt64(1234)
	b := biFromInt64(66)

	// add
	sum := bigIntegerAdd([]interface{}{a, b}).(*object.Object)
	if bigIntOf(sum).Int64() != 1300 {
		t.Fatalf("add mismatch: expected 1300, got %d", bigIntOf(sum).Int64())
	}

	// subtract
	diff := bigIntegerSubtract([]interface{}{a, b}).(*object.Object)
	if bigIntOf(diff).Int64() != 1168 {
		t.Fatalf("subtract mismatch: expected 1168, got %d", bigIntOf(diff).Int64())
	}

	// multiply by BigInteger
	prod := bigIntegerMultiply([]interface{}{a, b}).(*object.Object)
	if bigIntOf(prod).Int64() != 81444 {
		t.Fatalf("multiply mismatch: expected 81444, got %d", bigIntOf(prod).Int64())
	}

	// multiply by int64 shortcut
	prod2 := bigIntegerMultiply([]interface{}{a, int64(-2)}).(*object.Object)
	if bigIntOf(prod2).Int64() != -2468 {
		t.Fatalf("multiply by int mismatch: expected -2468, got %d", bigIntOf(prod2).Int64())
	}
}

func TestBigInteger_Divide_Remainder_Mod(t *testing.T) {
	globals.InitStringPool()

	a := biFromInt64(10)
	b := biFromInt64(3)

	// divide
	q := bigIntegerDivide([]interface{}{a, b}).(*object.Object)
	if bigIntOf(q).Int64() != 3 {
		t.Fatalf("divide mismatch: expected 3, got %d", bigIntOf(q).Int64())
	}

	// remainder
	r := bigIntegerRemainder([]interface{}{a, b}).(*object.Object)
	if bigIntOf(r).Int64() != 1 {
		t.Fatalf("remainder mismatch: expected 1, got %d", bigIntOf(r).Int64())
	}

	// divide by zero/nonpositive should error
	zero := biFromInt64(0)
	if err := bigIntegerDivide([]interface{}{a, zero}); err == nil {
		t.Fatalf("expected error for divide by zero")
	}
	if err := bigIntegerRemainder([]interface{}{a, zero}); err == nil {
		t.Fatalf("expected error for remainder by zero")
	}

	// mod requires positive modulus
	if err := bigIntegerMod([]interface{}{a, zero}); err == nil {
		t.Fatalf("expected error for modulus not positive")
	}
	mod := bigIntegerMod([]interface{}{a, biFromInt64(7)}).(*object.Object)
	if bigIntOf(mod).Int64() != 3 {
		t.Fatalf("mod mismatch: expected 3, got %d", bigIntOf(mod).Int64())
	}
}

func TestBigInteger_ModInverse_ModPow(t *testing.T) {
	globals.InitStringPool()

	mm := biFromInt64(11)
	xx := biFromInt64(3)

	inv := bigIntegerModInverse([]interface{}{xx, mm}).(*object.Object) // 3 * 4 ≡ 1 (mod 11)
	if bigIntOf(inv).Int64() != 4 {
		t.Fatalf("modInverse mismatch: expected 4, got %d", bigIntOf(inv).Int64())
	}

	// non-invertible when gcd != 1, e.g., 6 mod 12
	if err := bigIntegerModInverse([]interface{}{biFromInt64(6), biFromInt64(12)}); err == nil {
		t.Fatalf("expected error for non-invertible modInverse")
	}

	// modPow: 2^10 mod 11 = 1024 mod 11 = 1
	two := biFromInt64(2)
	ten := biFromInt64(10)
	res := bigIntegerModPow([]interface{}{two, ten, mm}).(*object.Object)
	if bigIntOf(res).Int64() != 1 {
		t.Fatalf("modPow mismatch: expected 1, got %d", bigIntOf(res).Int64())
	}

	// invalid modulus (<= 0)
	if err := bigIntegerModPow([]interface{}{two, ten, biFromInt64(0)}); err == nil {
		t.Fatalf("expected error for nonpositive modulus in modPow")
	}
}

func TestBigInteger_Compare_Signum_Abs_Negate(t *testing.T) {
	globals.InitStringPool()

	a := biFromInt64(-5)
	b := biFromInt64(7)

	cmp := bigIntegerCompareTo([]interface{}{a, b}).(int64)
	if cmp >= 0 {
		t.Fatalf("compareTo mismatch: expected negative, got %d", cmp)
	}

	sig := bigIntegerSignum([]interface{}{a}).(int64)
	if sig != -1 {
		t.Fatalf("signum mismatch: expected -1, got %d", sig)
	}

	abs := bigIntegerAbs([]interface{}{a}).(*object.Object)
	if bigIntOf(abs).Int64() != 5 {
		t.Fatalf("abs mismatch: expected 5, got %d", bigIntOf(abs).Int64())
	}

	neg := bigIntegerNegate([]interface{}{b}).(*object.Object)
	if bigIntOf(neg).Int64() != -7 {
		t.Fatalf("negate mismatch: expected -7, got %d", bigIntOf(neg).Int64())
	}
}

func TestBigInteger_BitProps_And_BitOps(t *testing.T) {
	globals.InitStringPool()

	// 0b1011_0001 = 177 (has 4 set bits)
	bi := biFromInt64(0b10110001)
	bc := bigIntegerBitCount([]interface{}{bi}).(int64)
	if bc != 4 { // bits set in 0b10110001
		t.Fatalf("bitCount mismatch: expected 4, got %d", bc)
	}

	bl := bigIntegerBitLength([]interface{}{bi}).(int64)
	if bl != 8 {
		t.Fatalf("bitLength mismatch: expected 8, got %d", bl)
	}

	// testBit
	if bigIntegerTestBit([]interface{}{bi, int64(0)}).(int64) != types.JavaBoolTrue {
		t.Fatalf("expected LSB bit set")
	}
	if bigIntegerTestBit([]interface{}{bi, int64(1)}).(int64) != types.JavaBoolFalse {
		t.Fatalf("expected bit 1 not set")
	}

	// setBit
	set := bigIntegerSetBit([]interface{}{bi, int64(1)}).(*object.Object)
	if bigIntOf(set).Int64() != 0b10110011 {
		t.Fatalf("setBit mismatch: expected 0b10110011, got %b", bigIntOf(set).Int64())
	}

	// shifts
	lsh := bigIntegerShiftLeft([]interface{}{bi, int64(2)}).(*object.Object)
	if bigIntOf(lsh).Int64() != (0b10110001 << 2) {
		t.Fatalf("shiftLeft mismatch: expected %d, got %d", 0b10110001<<2, bigIntOf(lsh).Int64())
	}
	rsh := bigIntegerShiftRight([]interface{}{bi, int64(3)}).(*object.Object)
	if bigIntOf(rsh).Int64() != (0b10110001 >> 3) {
		t.Fatalf("shiftRight mismatch: expected %d, got %d", 0b10110001>>3, bigIntOf(rsh).Int64())
	}
}

func TestBigInteger_Equals_Max_Min_GCD(t *testing.T) {
	globals.InitStringPool()

	a := biFromInt64(42)
	b := biFromInt64(42)
	c := biFromInt64(17)

	if bigIntegerEquals([]interface{}{a, b}).(int64) != types.JavaBoolTrue {
		t.Fatalf("equals mismatch: expected true")
	}
	if bigIntegerEquals([]interface{}{a, c}).(int64) != types.JavaBoolFalse {
		t.Fatalf("equals mismatch: expected false")
	}
	// wrong type argument should return error block (IllegalArgumentException)
	if err := bigIntegerEquals([]interface{}{a, int64(5)}); err == nil {
		t.Fatalf("expected error for equals with non-object argument")
	} else {
		if geb, ok := err.(*ghelpers.GErrBlk); ok {
			if geb.ExceptionType != excNames.IllegalArgumentException {
				t.Fatalf("expected IllegalArgumentException, got %d", geb.ExceptionType)
			}
		}
	}

	// max/min
	mx := bigIntegerMax([]interface{}{a, c}).(*object.Object)
	if bigIntOf(mx).Int64() != 42 {
		t.Fatalf("max mismatch: expected 42, got %d", bigIntOf(mx).Int64())
	}
	mn := bigIntegerMin([]interface{}{a, c}).(*object.Object)
	if bigIntOf(mn).Int64() != 17 {
		t.Fatalf("min mismatch: expected 17, got %d", bigIntOf(mn).Int64())
	}

	// gcd
	g := bigIntegerGCD([]interface{}{biFromInt64(48), biFromInt64(18)}).(*object.Object)
	if bigIntOf(g).Int64() != 6 {
		t.Fatalf("gcd mismatch: expected 6, got %d", bigIntOf(g).Int64())
	}
}

func TestBigInteger_IsProbablePrime_And_Pow_Sqrt(t *testing.T) {
	globals.InitStringPool()

	// isProbablePrime true for a known prime
	p := biFromInt64(101)
	if bigIntegerIsProbablePrime([]interface{}{p, int64(10)}).(int64) != 1 {
		t.Fatalf("isProbablePrime mismatch: expected true for 101")
	}

	// pow normal
	two := biFromInt64(2)
	eight := bigIntegerPow([]interface{}{two, int64(3)}).(*object.Object)
	if bigIntOf(eight).Int64() != 8 {
		t.Fatalf("pow mismatch: expected 8, got %d", bigIntOf(eight).Int64())
	}

	// pow negative exponent => error
	if err := bigIntegerPow([]interface{}{two, int64(-1)}); err == nil {
		t.Fatalf("expected error for negative exponent in pow")
	}

	// sqrt normal
	nine := biFromInt64(9)
	three := bigIntegerSqrt([]interface{}{nine}).(*object.Object)
	if bigIntOf(three).Int64() != 3 {
		t.Fatalf("sqrt mismatch: expected 3, got %d", bigIntOf(three).Int64())
	}

	// sqrt negative => error
	if err := bigIntegerSqrt([]interface{}{biFromInt64(-1)}); err == nil {
		t.Fatalf("expected error for negative sqrt input")
	}
}

func TestBigInteger_ByteArray_Construct_And_ToByteArray(t *testing.T) {
	globals.InitStringPool()

	// Construct from signed byte array representing -1 (0xFF)
	base := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	jb := object.JavaByteArrayFromGoByteArray([]byte{0xFF})
	byteArrObj := object.StringObjectFromJavaByteArray(jb)
	ret := bigIntegerInitByteArray([]interface{}{base, byteArrObj})
	if ret != nil {
		t.Fatalf("unexpected error initializing from byte array: %v", ret)
	}
	if bigIntOf(base).Int64() != -1 {
		t.Fatalf("byte-array init mismatch: expected -1, got %d", bigIntOf(base).Int64())
	}

	// toByteArray on positive number uses magnitude bytes in current implementation
	pos := biFromInt64(0x1234)
	arrObj := bigIntegerToByteArray([]interface{}{pos}).(*object.Object)
	got := asJavaBytesFromStringObject(arrObj)
	want := object.JavaByteArrayFromGoByteArray(big.NewInt(0x1234).Bytes())
	if !object.JavaByteArrayEquals(got, want) {
		t.Fatalf("toByteArray mismatch: expected %v, got %v", want, got)
	}
}

func TestBigInteger_HashCode(t *testing.T) {
	globals.InitStringPool()

	a := biFromInt64(12345)
	b := biFromInt64(12345)
	c := biFromInt64(54321)

	ha := bigIntegerHashCode([]interface{}{a}).(int64)
	hb := bigIntegerHashCode([]interface{}{b}).(int64)
	hc := bigIntegerHashCode([]interface{}{c}).(int64)

	if ha != hb {
		t.Errorf("hashCode mismatch for equal objects: %d != %d", ha, hb)
	}
	if ha == hc {
		t.Errorf("hashCode collision for different objects: %d == %d", ha, hc)
	}
}

func TestBigInteger_Conversions(t *testing.T) {
	globals.InitStringPool()

	val := int64(1234567890)
	bi := biFromInt64(val)

	if got := bigIntegerInt64Value([]interface{}{bi}).(int64); got != val {
		t.Errorf("intValue/longValue mismatch: expected %d, got %d", val, got)
	}

	if got := bigIntegerFloat64Value([]interface{}{bi}).(float64); got != float64(val) {
		t.Errorf("float64Value mismatch: expected %f, got %f", float64(val), got)
	}

	// byteValueExact
	biByte := biFromInt64(127)
	if got := bigIntegerByteValueExact([]interface{}{biByte}).(int64); got != 127 {
		t.Errorf("byteValueExact mismatch: expected 127, got %d", got)
	}

	biLarge := biFromInt64(128)
	if err := bigIntegerByteValueExact([]interface{}{biLarge}); err == nil {
		t.Error("expected ArithmeticException for byteValueExact out of range")
	}
}

func TestBigInteger_ToStringRadix(t *testing.T) {
	globals.InitStringPool()

	bi := biFromInt64(255)

	// Hex
	s := bigIntegerToStringRadix([]interface{}{bi, int64(16)}).(*object.Object)
	if asString(s) != "ff" {
		t.Errorf("toString(16) mismatch: expected ff, got %s", asString(s))
	}

	// Base 36
	s = bigIntegerToStringRadix([]interface{}{bi, int64(36)}).(*object.Object)
	if asString(s) != "73" {
		t.Errorf("toString(36) mismatch: expected 73, got %s", asString(s))
	}

	// Invalid radix
	if err := bigIntegerToStringRadix([]interface{}{bi, int64(1)}); err == nil {
		t.Error("expected IllegalArgumentException for radix < 2")
	}
	if err := bigIntegerToStringRadix([]interface{}{bi, int64(63)}); err == nil {
		t.Error("expected IllegalArgumentException for radix > 62")
	}
}

func TestBigInteger_Shift(t *testing.T) {
	globals.InitStringPool()

	bi := biFromInt64(1) // 0b1

	// shiftLeft
	res := bigIntegerShiftLeft([]interface{}{bi, int64(3)}).(*object.Object)
	if bigIntOf(res).Int64() != 8 {
		t.Errorf("shiftLeft(3) mismatch: expected 8, got %d", bigIntOf(res).Int64())
	}

	// shiftRight
	bi = biFromInt64(8)
	res = bigIntegerShiftRight([]interface{}{bi, int64(2)}).(*object.Object)
	if bigIntOf(res).Int64() != 2 {
		t.Errorf("shiftRight(2) mismatch: expected 2, got %d", bigIntOf(res).Int64())
	}
}

func TestBigInteger_Bitwise(t *testing.T) {
	globals.InitStringPool()

	a := biFromInt64(0b1100)
	b := biFromInt64(0b1010)

	// AND: 1100 & 1010 = 1000 (8)
	res := bigIntegerAnd([]interface{}{a, b}).(*object.Object)
	if bigIntOf(res).Int64() != 8 {
		t.Errorf("and mismatch: expected 8, got %d", bigIntOf(res).Int64())
	}

	// OR: 1100 | 1010 = 1110 (14)
	res = bigIntegerOr([]interface{}{a, b}).(*object.Object)
	if bigIntOf(res).Int64() != 14 {
		t.Errorf("or mismatch: expected 14, got %d", bigIntOf(res).Int64())
	}

	// XOR: 1100 ^ 1010 = 0110 (6)
	res = bigIntegerXor([]interface{}{a, b}).(*object.Object)
	if bigIntOf(res).Int64() != 6 {
		t.Errorf("xor mismatch: expected 6, got %d", bigIntOf(res).Int64())
	}

	// NOT: ~0 = -1
	res = bigIntegerNot([]interface{}{biFromInt64(0)}).(*object.Object)
	if bigIntOf(res).Int64() != -1 {
		t.Errorf("not mismatch: expected -1, got %d", bigIntOf(res).Int64())
	}

	// ANDNOT: 1100 & ~1010 = 1100 & 0101 = 0100 (4)
	res = bigIntegerAndNot([]interface{}{a, b}).(*object.Object)
	if bigIntOf(res).Int64() != 4 {
		t.Errorf("andNot mismatch: expected 4, got %d", bigIntOf(res).Int64())
	}
}

func TestBigInteger_StaticProbablyPrime(t *testing.T) {
	globals.InitStringPool()

	// bitLength 10
	res := bigIntegerProbablyPrime([]interface{}{int64(10), nil}).(*object.Object)
	bi := bigIntOf(res)
	if bi.BitLen() > 10 {
		t.Errorf("probablyPrime(10) bit length too large: %d", bi.BitLen())
	}
	if !bi.ProbablyPrime(10) {
		t.Errorf("probablyPrime(10) returned a non-prime: %s", bi.String())
	}
}
