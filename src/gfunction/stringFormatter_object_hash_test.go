package gfunction

import (
	"fmt"
	"jacobin/src/object"
	"testing"
	"unsafe"
)

func TestStringFormatter_Object_Hash_h_and_H(t *testing.T) {
	// Create a plain object with a class name and verify %h/%H use Object.hashCode() semantics
	className := "org/example/MyObject"
	obj := object.MakeEmptyObjectWithClassName(&className)
	// Expected is Integer.toHexString(obj.hashCode()), where our Object.hashCode is ptr^(ptr>>32)
	ptr := uintptr(unsafe.Pointer(obj))
	h := uint32(ptr ^ (ptr >> 32))
	expectedLower := fmt.Sprintf("%x", h)
	expectedUpper := fmt.Sprintf("%X", h)

	fmtObj := object.StringObjectFromGoString("%h %H")
	argsArr := makeObjectRefArray(obj, obj)
	out := StringFormatter([]interface{}{fmtObj, argsArr})
	got := object.GoStringFromStringObject(out.(*object.Object))
	want := expectedLower + " " + expectedUpper
	if got != want {
		t.Fatalf("got %q want %q (hash=%#x)", got, want, uint32(obj.Mark.Hash))
	}
}
