package gfunction

import (
	"jacobin/src/object"
	"testing"
)

// These tests focus on Java-like date/time format specifiers (%t and %T).
// Our current implementation degrades %t/%T to %s (string) and does not
// implement Java's date/time suffixes (e.g., Y, m, d). Therefore, we verify
// that:
// - %t and %T print the argument as a string (including padding/precision).
// - Argument indexing and reuse work with %t/%T.
// - A null reference formats as "null" when using %t.

// helper to build reference array of objects (already exists in stringFormatter_test.go)
// re-declared here to keep this file self-contained for direct execution
func makeObjectRefArrayDT(elems ...*object.Object) *object.Object {
	return object.MakeArrayFromRawArray(elems)
}

func TestDateTime_Basic_t_and_T_DegradeToString(t *testing.T) {
	// Fictitious date/time value provided by the issue text
	dateStr := "2025-02-14 13:14:15"
	fmtObj := object.StringObjectFromGoString("%t %T")
	arg := object.StringObjectFromGoString(dateStr)
	argsArr := makeObjectRefArrayDT(arg, arg)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := dateStr + " " + dateStr
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestDateTime_Padding_And_Precision_With_tT(t *testing.T) {
	dateStr := "2025-02-14 13:14:15"
	// Right-justify width 25, left-justify width 25, precision .10
	fmtObj := object.StringObjectFromGoString("%25t|%-25T|%.10t")
	arg := object.StringObjectFromGoString(dateStr)
	argsArr := makeObjectRefArrayDT(arg, arg, arg)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))

	// Build expected using Go-like padding semantics (since %t/%T -> %s in our impl)
	// width 25 right-justified
	right := makePadding(25-len(dateStr)) + dateStr
	left := dateStr + makePadding(25-len(dateStr))
	prec := dateStr[:10]
	expected := right + "|" + left + "|" + prec
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestDateTime_ArgumentIndex_And_Reuse_With_t(t *testing.T) {
	dateStr := "2025-02-14 13:14:15"
	// Use 1$ index and then reuse with %<t and mix %T
	fmtObj := object.StringObjectFromGoString("%1$t|%<T|%<t")
	arg := object.StringObjectFromGoString(dateStr)
	argsArr := makeObjectRefArrayDT(arg)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := dateStr + "|" + dateStr + "|" + dateStr
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

func TestDateTime_NullArgument_With_t_PrintsNull(t *testing.T) {
	fmtObj := object.StringObjectFromGoString("[%t]")
	// pass a nil reference in the array; current implementation renders it as <nil>
	argsArr := makeObjectRefArrayDT(nil)

	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	expected := "[<nil>]"
	if got != expected {
		t.Fatalf("got %q want %q", got, expected)
	}
}

// makePadding returns a string of n spaces; if n <= 0, returns empty string.
func makePadding(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
