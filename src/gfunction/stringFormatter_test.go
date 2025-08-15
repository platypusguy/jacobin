package gfunction

import (
    "jacobin/object"
    "jacobin/types"
    "runtime"
    "strings"
    "testing"
)

// helper: build a reference array [Ljava/lang/Object; from provided element objects
func makeObjectRefArray(elems ...*object.Object) *object.Object {
    arr := object.Make1DimRefArray("Ljava/lang/Object;", int64(len(elems)))
    fv := arr.FieldTable["value"]
    fv.Fvalue = elems
    arr.FieldTable["value"] = fv
    return arr
}

func TestStringFormatter_Newline_And_Boolean(t *testing.T) {
    // format with %n and %b on a non-boolean arg (should be true)
    fmtObj := object.StringObjectFromGoString("Line1%nLine2 %b")
    // one integer argument
    intObj := Populator("java/lang/Integer", types.Int, int64(1))
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
    fmtObj := object.StringObjectFromGoString("obj=%s")
    o := object.MakeEmptyObject()
    o.KlassName = object.StringPoolIndexFromGoString("com/example/Dummy")
    argsArr := makeObjectRefArray(o)

    out := StringFormatter([]interface{}{fmtObj, argsArr})
    got := object.GoStringFromStringObject(out.(*object.Object))
    if !strings.HasPrefix(got, "obj=Dummy@") {
        t.Fatalf("expected Dummy@..., got %q", got)
    }
}

func TestStringFormatter_Numeric_Passthrough(t *testing.T) {
    fmtObj := object.StringObjectFromGoString("%04x")
    i := Populator("java/lang/Integer", types.Int, int64(26))
    argsArr := makeObjectRefArray(i)
    out := StringFormatter([]interface{}{fmtObj, argsArr})
    got := object.GoStringFromStringObject(out.(*object.Object))
    if !(got == "001a" || got == "1a") {
        t.Fatalf("got %q want one of %q or %q", got, "001a", "1a")
    }
}

func TestStringFormatter_FloatZeroPad(t *testing.T) {
    fmtObj := object.StringObjectFromGoString("%020.12f")
    d := Populator("java/lang/Double", types.Double, float64(123.4567))
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
    i := Populator("java/lang/Integer", types.Int, int64(-64))
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
    b := Populator("java/lang/Byte", types.Byte, int64(-1))
    argsArr := makeObjectRefArray(b)
    out := StringFormatter([]interface{}{fmtObj, argsArr})
    got := object.GoStringFromStringObject(out.(*object.Object))
    expected := "ff"
    if got != expected {
        t.Fatalf("got %q want %q", got, expected)
    }
}
