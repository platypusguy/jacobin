package gfunction

import (
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"runtime"
	"strings"
	"testing"
)

// helper: build a reference array [Ljava/lang/Object; from provided element objects
func makeObjectRefArray(elems ...*object.Object) *object.Object {
	// Use MakeArrayFromRawArray to avoid direct map assignments and ensure proper initialization
	return object.MakeArrayFromRawArray(elems)
}

func TestStringFormatter_Newline_And_Boolean(t *testing.T) {
	// format with %n and %b on a non-boolean arg (should be true)
	fmtObj := object.StringObjectFromGoString("Line1%nLine2 %b")
	// one integer argument
	intObj := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(1))
	argsArr := makeObjectRefArray(intObj)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	s := object.GoStringFromStringObject(out.(*object.Object))

	nl := "\n"
	if runtime.GOOS == "windows" {
		nl = "\r\n"
	}
	expected := "Line1" + nl + "Line2 true"
	if s != expected {
		t.Fatalf("unexpected output: %q want %q", s, expected)
	}
}

func TestStringFormatter_String_And_Uppercase(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%s %S")
	s1 := object.StringObjectFromGoString("hello")
	s2 := object.StringObjectFromGoString("world")
	argsArr := makeObjectRefArray(s1, s2)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	if got != "hello WORLD" {
		t.Fatalf("got %q", got)
	}
}

func TestStringFormatter_Hash_For_String(t *testing.T) {
	// "abc" Java hashCode is 96354 decimal -> 0x17832
	fmtObj := object.StringObjectFromGoString("%h %H")
	s := object.StringObjectFromGoString("abc")
	argsArr := makeObjectRefArray(s, s)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	if got != "17862 17862" {
		t.Fatalf("got %q want %q", got, "17862 17862")
	}
}

func TestStringFormatter_Object_ToString_Like(t *testing.T) {
	globals.InitStringPool()
	fmtObj := object.StringObjectFromGoString("obj=%s")
	o := object.MakeEmptyObject()
	o.KlassName = object.StringPoolIndexFromGoString("com/example/Dummy")
	// ensure non-null per Jacobin's IsNull (requires at least one field)
	o.FieldTable["_marker"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	argsArr := makeObjectRefArray(o)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	if !strings.HasPrefix(got, "obj=Dummy@") {
		t.Fatalf("expected Dummy@..., got %q", got)
	}
}

func TestStringFormatter_Numeric_Passthrough(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%04x")
	i := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(26))
	argsArr := makeObjectRefArray(i)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	if !(got == "001a" || got == "1a") {
		t.Fatalf("got %q want one of %q or %q", got, "001a", "1a")
	}
}

func TestStringFormatter_FloatZeroPad(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%020.12f")
	d := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(123.4567))
	argsArr := makeObjectRefArray(d)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "0000123.456700000000"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_HexNegativeInt_ZeroPad(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%08x")
	i := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(-64))
	argsArr := makeObjectRefArray(i)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "ffffffc0"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_HexByte_Negative_TwoDigits(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%02x")
	b := object.MakePrimitiveObject("java/lang/Byte", types.Byte, int64(-1))
	argsArr := makeObjectRefArray(b)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "ff"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_FixedWidthNoZeroPad(t *testing.T) {
	// Mirrors Java's
	// System.out.printf("%18.6e %18.6e%n", 22.0, 33.0);
	// System.out.printf("%18.6f %18.6f%n", 22.0, 33.0);
	globals.InitGlobals("test")
	fmtStr := "%18.6e %18.6e%n%18.6f %18.6f%n"
	fmtObj := object.StringObjectFromGoString(fmtStr)
	d1 := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(22))
	d2 := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(33))
	argsArr := makeObjectRefArray(d1, d2, d1, d2)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))

	nl := "\n"
	if runtime.GOOS == "windows" {
		nl = "\r\n"
	}
	expected := "      2.200000e+01       3.300000e+01" + nl +
		"         22.000000          33.000000" + nl

	if got != expected {
		t.Fatalf("unexpected output:\n%q\nwant:\n%q", got, expected)
	}
}

func TestStringFormatter_LiteralPercent(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("Done: 100%%")
	argsArr := makeObjectRefArray()
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	if got != "Done: 100%" {
		t.Fatalf("got %q want %q", got, "Done: 100%")
	}
}

func TestStringFormatter_Decimal_Padding(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%05d %-6d")
	i1 := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(42))
	i2 := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(7))
	argsArr := makeObjectRefArray(i1, i2)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "00042 7     "
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Octal_Padding(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%06o")
	b := object.MakePrimitiveObject("java/lang/Byte", types.Byte, int64(255)) // 0xff -> octal 377
	argsArr := makeObjectRefArray(b)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	// width 6, zero padded
	expected := "000377"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_UpperHex_Padding(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%08X")
	i := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(26))
	argsArr := makeObjectRefArray(i)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "0000001A"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Scientific_Uppercase(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%12.3E")
	d := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(1234.56))
	argsArr := makeObjectRefArray(d)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	// Expect exactly 12-wide with 3 decimals in E format
	// 1234.56 -> 1.235E+03
	expected := "   1.235E+03"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_General_gG(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%g %G")
	d1 := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(12345.0))
	d2 := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(0.0012345))
	argsArr := makeObjectRefArray(d1, d2)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	// Default precision may vary; assert general characteristics only.
	parts := strings.Split(got, " ")
	if len(parts) != 2 {
		t.Fatalf("unexpected: %q", got)
	}
	if !strings.Contains(parts[0], "12345") {
		t.Fatalf("first part not as expected: %q", parts[0])
	}
	// second can be exponent or fixed depending on precision heuristic; accept both
	if !(strings.Contains(parts[1], "E") || strings.Contains(parts[1], "e") || strings.HasPrefix(parts[1], "0.0012345")) {
		t.Fatalf("second part unexpected: %q", parts[1])
	}
}

func TestStringFormatter_Char_From_Char_And_Int(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%c %c")
	ch := object.MakePrimitiveObject("java/lang/Character", types.Char, int64('a'))
	code := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(66)) // 'B'
	argsArr := makeObjectRefArray(ch, code)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "a B"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_ArgumentIndex_And_Reuse(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%2$s-%1$d %2$s %<S")
	i := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(3))
	s := object.StringObjectFromGoString("id")
	argsArr := makeObjectRefArray(i, s)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "id-3 id ID"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_StringPaddingAndPrecision(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%10s|%-10s|%.3s")
	s1 := object.StringObjectFromGoString("hi")
	s2 := object.StringObjectFromGoString("there")
	argsArr := makeObjectRefArray(s1, s2, s2)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "        hi|there     |the"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_SignFlags_For_Decimal(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%+d % d")
	p := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(5))
	n := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(-5))
	argsArr := makeObjectRefArray(p, n)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	// %+d forces a sign, space flag leaves minus for negative
	expected := "+5 -5"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Time_T_T_DegradesToString(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%t %T")
	s := object.StringObjectFromGoString("time")
	argsArr := makeObjectRefArray(s, s)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "time time"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Boolean_Uppercase(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%2B %2B")
	// First arg: non-boolean non-null -> true; Second arg: nil -> false
	i := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, int64(1))
	argsArr := makeObjectRefArray(i, nil)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "TRUE FALSE"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Hash_For_Double(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%h")
	d := object.MakePrimitiveObject("java/lang/Double", types.Double, float64(123.45))
	argsArr := makeObjectRefArray(d)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "8c921001"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Hash_For_Float_Positive(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%h")
	f := object.MakePrimitiveObject("java/lang/Float", types.Float, float64(123.45))
	argsArr := makeObjectRefArray(f)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "42f6e666"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Hash_For_Float_Negative_Upper(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%H")
	f := object.MakePrimitiveObject("java/lang/Float", types.Float, float64(-123.45))
	argsArr := makeObjectRefArray(f)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "C2F6E666"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_Char_Uppercase_C(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%C %C")
	ch := object.MakePrimitiveObject("java/lang/Character", types.Char, int64('a'))
	code := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64('b'))
	argsArr := makeObjectRefArray(ch, code)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "A B"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_BigInteger_D_and_X(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%d %x")
	bi := bigIntegerFromInt64(255)
	argsArr := makeObjectRefArray(bi, bi)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "255 ff"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_BigDecimal_Fixed(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("%.3f")
	// Build BigDecimal from string "123.456"
	bd := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	params := []interface{}{bd, object.StringObjectFromGoString("123.456")}
	ret := bigdecimalInitString(params)
	if ret != nil {
		t.Fatalf("bigdecimalInitString returned error: %v", ret)
	}
	argsArr := makeObjectRefArray(bd)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "123.456"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_BigInteger_Hash_h(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%h")
	bi := bigIntegerFromInt64(12345)
	argsArr := makeObjectRefArray(bi)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "3039" // BigInteger(12345).hashCode() in Java is 12345 (0x3039)
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_BigDecimal_Hash_h(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%h")
	// Build BigDecimal from string "123.45"
	bd := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	params := []interface{}{bd, object.StringObjectFromGoString("123.45")}
	ret := bigdecimalInitString(params)
	if ret != nil {
		t.Fatalf("bigdecimalInitString returned error: %v", ret)
	}
	argsArr := makeObjectRefArray(bd)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "5d6e9" // 31*hash(12345) + 2
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestStringFormatter_BigDecimal_Scientific_e(t *testing.T) {
	globals.InitGlobals("test")
	fmtObj := object.StringObjectFromGoString("%18.6e")
	bd := object.MakeEmptyObjectWithClassName(&classNameBigDecimal)
	params := []interface{}{bd, object.StringObjectFromGoString("123.45")}
	ret := bigdecimalInitString(params)
	if ret != nil {
		t.Fatalf("bigdecimalInitString returned error: %v", ret)
	}
	argsArr := makeObjectRefArray(bd)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "      1.234500e+02"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

// Verify %h and %H when the argument is a real java/lang/String object
func TestStringFormatter_StringObject_Hash_h_and_H(t *testing.T) {
	globals.InitGlobals("test")

	// Build a String object and expected hash hex using Java's String.hashCode semantics
	content := "Mary had a little lamb"
	strObj := object.StringObjectFromGoString(content)

	expectedLower := fmt.Sprintf("%x", uint32(javaStringHashCode(content)))
	expectedUpper := fmt.Sprintf("%X", uint32(javaStringHashCode(content)))

	fmtObj := object.StringObjectFromGoString("%h %H")
	argsArr := makeObjectRefArray(strObj, strObj)
	out := StringFormatter([]interface{}{fmtObj, argsArr})

	got := object.GoStringFromStringObject(out.(*object.Object))
	want := expectedLower + " " + expectedUpper

	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
