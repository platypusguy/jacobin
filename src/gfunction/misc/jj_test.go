package misc

import (
	"io"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

func TestJjStringifyScalar_BoolTrue(t *testing.T) {
	result := jjStringifyScalar(types.Bool, types.JavaBoolTrue)
	expected := object.StringObjectFromGoString("true")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_BoolFalse(t *testing.T) {
	result := jjStringifyScalar(types.Bool, types.JavaBoolFalse)
	expected := object.StringObjectFromGoString("false")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_Byte(t *testing.T) {
	result := object.GoStringFromStringObject(jjStringifyScalar(types.Byte, byte(0xAB)))
	expected := "0xab"
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_Char(t *testing.T) {
	result := jjStringifyScalar(types.Char, int64('A'))
	expected := object.StringObjectFromGoString("A")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_Double(t *testing.T) {
	result := jjStringifyScalar(types.Double, 3.14)
	expected := object.StringObjectFromGoString("3.14")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_Float(t *testing.T) {
	result := jjStringifyScalar(types.Float, 4.15)
	observed := object.GoStringFromStringObject(result)
	expected := "4.15"
	if observed != expected {
		t.Errorf("Expected: %v, observed: %v", expected, observed)
	}
}

func TestJjStringifyScalar_Int(t *testing.T) {
	result := jjStringifyScalar(types.Int, int64(42))
	expected := object.StringObjectFromGoString("42")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_Long(t *testing.T) {
	result := jjStringifyScalar(types.Long, int64(42))
	expected := object.StringObjectFromGoString("42")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_String(t *testing.T) {
	strObj := object.StringObjectFromGoString("test")
	result := jjStringifyScalar("Ljava/lang/String;", strObj)
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(strObj) {
		t.Errorf("Expected %v, got %v", strObj, result)
	}
}

func TestJjStringifyScalar_Short(t *testing.T) {
	result := jjStringifyScalar(types.Short, int64(42))
	expected := object.StringObjectFromGoString("42")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_RefNull(t *testing.T) {
	result := jjStringifyScalar(types.Ref, object.Null)
	expected := object.StringObjectFromGoString(types.NullString)
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyScalar_RefNonNull(t *testing.T) {
	globals.InitGlobals("test")
	expected := "ABC"
	obj := object.StringObjectFromGoString(expected)
	result := jjStringifyScalar(types.Ref, obj)
	observed := object.GoStringFromStringObject(result)
	if observed != expected {
		t.Errorf("Expected %s, got %s", expected, observed)
	}
}

func TestJjStringifyScalar_Default(t *testing.T) {
	result := jjStringifyScalar("UnknownType", 42)
	expected := object.StringObjectFromGoString("42")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyVector_ObjectArray(t *testing.T) {
	obj := &object.Object{
		FieldTable: map[string]object.Field{
			"value": {Fvalue: []int64{1, 2, 3}},
		},
	}
	result := jjStringifyVector(obj)
	expected := object.StringObjectFromGoString("1,2,3")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyVector_Slice(t *testing.T) {
	slice := []int{4, 5, 6}
	result := jjStringifyVector(slice)
	expected := object.StringObjectFromGoString("4,5,6")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyVector_EmptyArray(t *testing.T) {
	obj := &object.Object{
		FieldTable: map[string]object.Field{
			"value": {Fvalue: []int64{}},
		},
	}
	result := jjStringifyVector(obj)
	expected := object.StringObjectFromGoString("")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjStringifyVector_EmptySlice(t *testing.T) {
	slice := []int{}
	result := jjStringifyVector(slice)
	expected := object.StringObjectFromGoString("")
	if object.GoStringFromStringObject(result) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetStaticString_InvalidClassObject(t *testing.T) {
	params := []interface{}{nil, &object.Object{KlassName: types.InvalidStringIndex}}
	result := jjGetStaticString(params)
	expected := object.StringObjectFromGoString("jjGetStaticString: No class object")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetStaticString_InvalidFieldObject(t *testing.T) {
	params := []interface{}{&object.Object{KlassName: 1}, nil}
	result := jjGetStaticString(params)
	expected := object.StringObjectFromGoString("jjGetStaticString: Invalid field is missing or nil")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetStaticString_VectorType(t *testing.T) {
	classObj := object.StringObjectFromGoString("testClass")
	fieldObj := object.StringObjectFromGoString("testField")
	statics.AddStatic("testClass.testField",
		statics.Static{Type: types.Array, Value: &object.Object{FieldTable: map[string]object.Field{"value": {Fvalue: []int64{1, 2, 3}}}}})
	params := []interface{}{classObj, fieldObj}
	result := jjGetStaticString(params)
	expected := object.StringObjectFromGoString("1,2,3")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetStaticString_ScalarType(t *testing.T) {
	classObj := object.StringObjectFromGoString("testClass")
	fieldObj := object.StringObjectFromGoString("testField")
	statics.AddStatic("testClass.testField", statics.Static{Type: types.Int, Value: int64(42)})
	params := []interface{}{classObj, fieldObj}
	result := jjGetStaticString(params)
	expected := object.StringObjectFromGoString("42")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetFieldString_InvalidFieldObject(t *testing.T) {
	params := []interface{}{&object.Object{}, nil}
	result := jjGetFieldString(params)
	expected := object.StringObjectFromGoString("jjGetFieldString: Invalid field is missing or nil")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetFieldString_NoSuchFieldName(t *testing.T) {
	thisObj := &object.Object{FieldTable: make(map[string]object.Field)}
	fieldObj := object.StringObjectFromGoString("nonexistentField")
	params := []interface{}{thisObj, fieldObj}
	result := jjGetFieldString(params)
	expected := object.StringObjectFromGoString("jjGetFieldString: No such field name: nonexistentField")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetFieldString_JavaStringField(t *testing.T) {
	thisObj := &object.Object{
		FieldTable: map[string]object.Field{
			"javaStringField": {Ftype: "Ljava/lang/String;", Fvalue: []types.JavaByte{0x74, 0x65, 0x73, 0x74}},
		},
	}
	fieldObj := object.StringObjectFromGoString("javaStringField")
	params := []interface{}{thisObj, fieldObj}
	result := jjGetFieldString(params)
	expected := object.StringObjectFromGoString("test")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetFieldString_VectorField(t *testing.T) {
	thisObj := &object.Object{
		FieldTable: map[string]object.Field{
			"vectorField": {Ftype: types.Array, Fvalue: []int64{1, 2, 3}},
		},
	}
	fieldObj := object.StringObjectFromGoString("vectorField")
	params := []interface{}{thisObj, fieldObj}
	result := jjGetFieldString(params)
	expected := object.StringObjectFromGoString("1,2,3")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestJjGetFieldString_ScalarField(t *testing.T) {
	thisObj := &object.Object{
		FieldTable: map[string]object.Field{
			"scalarField": {Ftype: types.Int, Fvalue: int64(42)},
		},
	}
	fieldObj := object.StringObjectFromGoString("scalarField")
	params := []interface{}{thisObj, fieldObj}
	result := jjGetFieldString(params)
	expected := object.StringObjectFromGoString("42")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func fnTestJjDumpStatics(t *testing.T, selection int64, className string, threesome []string) {
	// Re-direct stderr.
	originalStderr := os.Stderr
	rerr, werr, _ := os.Pipe()
	os.Stderr = werr

	// Dump statics.
	objTitle := object.StringObjectFromGoString("TestDumpStatics")
	objClassName := object.StringObjectFromGoString(className)
	params := make([]interface{}, 3)
	params[0] = objTitle
	params[1] = selection
	params[2] = objClassName
	jjDumpStatics(params)

	// Close the working stderr, capture its contents, and restore the original stderr.
	_ = werr.Close()
	bytes, _ := io.ReadAll(rerr)
	contents := string(bytes[:])
	os.Stderr = originalStderr

	if !strings.Contains(contents, threesome[0]) || !strings.Contains(contents, threesome[1]) || !strings.Contains(contents, threesome[2]) {
		t.Errorf("fnTestDumpStatics(%d, \"%s\"): looking for these: %v", selection, className, threesome)
		t.Errorf("fnTestDumpStatics(%d, \"%s\"): didn't see them in DumpStatics output: %s", selection, className, contents)
	}

}

func TestJjDumpStatics(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	statics.Statics = make(map[string]statics.Static)

	err1 := statics.AddStatic("test.f1", statics.Static{Type: types.Byte, Value: 'B'})
	err2 := statics.AddStatic("test.f2", statics.Static{Type: types.Int, Value: int(42)})
	err3 := statics.AddStatic("test.f3", statics.Static{Type: types.Double, Value: 24.0})
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("TestIntConversions: got unexpected error adding statics for testing")
	}

	fnTestJjDumpStatics(t, statics.SelectAll, "", []string{"test.f1", "test.f2", "test.f3"})
	fnTestJjDumpStatics(t, statics.SelectUser, "", []string{"test.f1", "test.f2", "test.f3"})
	fnTestJjDumpStatics(t, statics.SelectClass, "test", []string{"test.f1", "test.f2", "test.f3"})
}

func TestJjDumpObject_InvalidObject(t *testing.T) {
	objTitle := object.StringObjectFromGoString("Test Object")
	params := []interface{}{nil, objTitle, int64(2)}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid object")
		}
	}()

	_ = jjDumpObject(params)
}

func TestJjDumpObject_InvalidTitle(t *testing.T) {
	obj := &object.Object{}
	params := []interface{}{obj, nil, int64(2)}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid title")
		}
	}()

	_ = jjDumpObject(params)
}

func TestJjDumpObject_InvalidIndent(t *testing.T) {
	obj := &object.Object{}
	objTitle := object.StringObjectFromGoString("Test Object")
	params := []interface{}{obj, objTitle, nil}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid indent")
		}
	}()

	_ = jjDumpObject(params)
}
