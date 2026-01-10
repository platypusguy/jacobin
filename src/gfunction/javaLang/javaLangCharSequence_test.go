package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"reflect"
	"testing"
)

func TestLoad_Lang_CharSequence_RegistersMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_CharSequence()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/lang/CharSequence.compare(Ljava/lang/CharSequence;Ljava/lang/CharSequence;)I", 2, charSequenceCompare},
		{"java/lang/CharSequence.length()I", 0, charSequenceLength},
		{"java/lang/CharSequence.charAt(I)C", 1, charSequenceCharAt},
		{"java/lang/CharSequence.subSequence(II)Ljava/lang/CharSequence;", 2, charSequenceSubSequence},
		{"java/lang/CharSequence.toString()Ljava/lang/String;", 0, charSequenceToString},
	}

	for _, c := range checks {
		got, ok := ghelpers.MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing ghelpers.MethodSignatures entry for %s", c.key)
		}
		if got.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
		}
		if got.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}

func TestCharSequence_StringImplementation(t *testing.T) {
	globals.InitStringPool()

	str := "Hello"
	sObj := object.StringObjectFromGoString(str)

	// length()
	res := charSequenceLength([]interface{}{sObj})
	if res.(int64) != 5 {
		t.Errorf("expected length 5, got %v", res)
	}

	// charAt(1)
	res = charSequenceCharAt([]interface{}{sObj, int64(1)})
	if res.(int64) != int64('e') {
		t.Errorf("expected 'e', got %v", res)
	}

	// subSequence(1, 4) -> "ell"
	res = charSequenceSubSequence([]interface{}{sObj, int64(1), int64(4)})
	subObj := res.(*object.Object)
	if object.GoStringFromStringObject(subObj) != "ell" {
		t.Errorf("expected \"ell\", got %q", object.GoStringFromStringObject(subObj))
	}

	// toString()
	res = charSequenceToString([]interface{}{sObj})
	if res.(*object.Object) != sObj {
		t.Errorf("expected same object for toString on String")
	}
}

func TestCharSequence_StringBuilderImplementation(t *testing.T) {
	globals.InitStringPool()

	// Mocking a StringBuilder object
	sbClassName := "java/lang/StringBuilder"
	sbObj := object.MakeEmptyObjectWithClassName(&sbClassName)
	sbObj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []types.JavaByte{'H', 'i'}}
	sbObj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: int64(2)}

	// length()
	res := charSequenceLength([]interface{}{sbObj})
	if res.(int64) != 2 {
		t.Errorf("expected length 2, got %v", res)
	}

	// charAt(1)
	res = charSequenceCharAt([]interface{}{sbObj, int64(1)})
	if res.(int64) != int64('i') {
		t.Errorf("expected 'i', got %v", res)
	}
}

func TestCharSequence_Compare(t *testing.T) {
	globals.InitStringPool()

	s1 := object.StringObjectFromGoString("abc")
	s2 := object.StringObjectFromGoString("def")
	s3 := object.StringObjectFromGoString("abc")

	// compare(s1, s2) -> -1
	res := charSequenceCompare([]interface{}{s1, s2})
	if res.(int64) != -1 {
		t.Errorf("expected -1, got %v", res)
	}

	// compare(s2, s1) -> 1
	res = charSequenceCompare([]interface{}{s2, s1})
	if res.(int64) != 1 {
		t.Errorf("expected 1, got %v", res)
	}

	// compare(s1, s3) -> 0
	res = charSequenceCompare([]interface{}{s1, s3})
	if res.(int64) != 0 {
		t.Errorf("expected 0, got %v", res)
	}
}
