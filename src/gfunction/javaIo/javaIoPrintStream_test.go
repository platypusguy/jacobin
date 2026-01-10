package javaIo

import (
	"bytes"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
	"testing"
)

// Helper to create a Java String object
func makeStringObject(s string) *object.Object {
	return object.StringObjectFromGoString(s)
}

func TestPrintStreamPrintlnFunctions(t *testing.T) {
	buf := new(bytes.Buffer)

	// println()
	buf.Reset()
	PrintlnV([]interface{}{buf})
	if got := buf.String(); got != "\n" {
		t.Errorf("PrintlnV() = %q; want newline", got)
	}

	// println(String)
	buf.Reset()
	PrintlnString([]interface{}{buf, makeStringObject("hello")})
	if got := buf.String(); got != "hello\n" {
		t.Errorf("PrintlnString() = %q; want 'hello\\n'", got)
	}

	// println(byte), println(int), println(short) all use PrintlnBIS
	for _, val := range []int64{65, -1, 0} {
		buf.Reset()
		PrintlnBIS([]interface{}{buf, val})
		want := strconv.FormatInt(val, 10) + "\n"
		if got := buf.String(); got != want && val != 65 {
			// Actually, PrintlnBIS uses fmt.Fprintln on int64 so the output is a decimal string plus newline
			want = strconv.FormatInt(val, 10) + "\n"
		}
		if got := buf.String(); got != want {
			t.Errorf("PrintlnBIS(%d) = %q; want %q", val, got, want)
		}
	}

	// println(boolean)
	buf.Reset()
	PrintlnBoolean([]interface{}{buf, int64(1)})
	if got := buf.String(); got != "true\n" {
		t.Errorf("PrintlnBoolean(true) = %q; want 'true\\n'", got)
	}
	buf.Reset()
	PrintlnBoolean([]interface{}{buf, int64(0)})
	if got := buf.String(); got != "false\n" {
		t.Errorf("PrintlnBoolean(false) = %q; want 'false\\n'", got)
	}

	// println(long)
	buf.Reset()
	PrintlnLong([]interface{}{buf, int64(1234567890)})
	if got := buf.String(); got != "1234567890\n" {
		t.Errorf("PrintlnLong() = %q; want '1234567890\\n'", got)
	}

	// println(double)
	buf.Reset()
	PrintlnDouble([]interface{}{buf, 3.1415926535})
	if got := buf.String(); got != "3.1415926535\n" {
		t.Errorf("PrintlnDouble() = %q; want '3.1415926535\\n'", got)
	}

	// println(float)
	buf.Reset()
	PrintlnFloat([]interface{}{buf, 2.71828})
	if got := buf.String(); got != "2.71828\n" {
		t.Errorf("PrintlnFloat() = %q; want '2.71828\\n'", got)
	}

	// println(char)
	buf.Reset()
	PrintlnChar([]interface{}{buf, int64('Z')})
	if got := buf.String(); got != "Z\n" {
		t.Errorf("PrintlnChar() = %q; want 'Z\\n'", got)
	}
}

func TestPrintStreamPrintFunctions(t *testing.T) {
	buf := new(bytes.Buffer)

	// print(String)
	buf.Reset()
	PrintString([]interface{}{buf, makeStringObject("test string")})
	if got := buf.String(); got != "test string" {
		t.Errorf("PrintString() = %q; want 'test string'", got)
	}

	// print(byte), print(int), print(short) use PrintBIS
	for _, val := range []int64{65, -2, 0} {
		buf.Reset()
		PrintBIS([]interface{}{buf, val})
		want := strconv.FormatInt(val, 10)
		if got := buf.String(); got != want {
			t.Errorf("PrintBIS(%d) = %q; want %q", val, got, want)
		}
	}

	// print(boolean)
	buf.Reset()
	PrintBoolean([]interface{}{buf, int64(1)})
	if got := buf.String(); got != "true" {
		t.Errorf("PrintBoolean(true) = %q; want 'true'", got)
	}
	buf.Reset()
	PrintBoolean([]interface{}{buf, int64(0)})
	if got := buf.String(); got != "false" {
		t.Errorf("PrintBoolean(false) = %q; want 'false'", got)
	}

	// print(long)
	buf.Reset()
	PrintLong([]interface{}{buf, int64(9876543210)})
	if got := buf.String(); got != "9876543210" {
		t.Errorf("PrintLong() = %q; want '9876543210'", got)
	}

	// print(double)
	buf.Reset()
	PrintDouble([]interface{}{buf, 6.62607015})
	if got := buf.String(); got != "6.62607015" {
		t.Errorf("PrintDouble() = %q; want '6.62607015'", got)
	}

	// print(float)
	buf.Reset()
	PrintFloat([]interface{}{buf, 1.41421})
	if got := buf.String(); got != "1.41421" {
		t.Errorf("PrintFloat() = %q; want '1.41421'", got)
	}

	// print(char)
	buf.Reset()
	PrintChar([]interface{}{buf, int64('Y')})
	if got := buf.String(); got != "Y" {
		t.Errorf("PrintChar() = %q; want 'Y'", got)
	}
}

func TestPrintStreamPrintf(t *testing.T) {
	// Create a buffer to capture output instead of writing to os.Stdout
	var buf bytes.Buffer
	ps := &buf

	// Prepare format string object
	fmtStr := object.StringObjectFromGoString("Hello %d and %d!\n")

	// Prepare argument objects (Java Strings)
	globals.InitStringPool()
	arg1 := object.MakePrimitiveObject(types.Int, types.Int, int64(1))
	arg2 := object.MakePrimitiveObject(types.Int, types.Int, int64(2))
	args := []*object.Object{arg1, arg2}
	iArr := object.MakePrimitiveObject("java/lang/Object", types.RefArray, args)

	// Params: first PrintStream, then format string, then each arg separately
	params := []interface{}{ps, fmtStr, iArr}

	// Call the Printf function under test
	ret := Printf(params)
	if ret != ps {
		t.Errorf("Printf did not return the PrintStream object as expected, errMsg: %s", ret)
	}

	// Check the output captured in buf
	want := "Hello 1 and 2!\n"
	got := buf.String()
	if got != want {
		t.Errorf("Printf output mismatch:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestPrintStreamPrintObject(t *testing.T) {
	buf := new(bytes.Buffer)

	// Case 1: Print(String object) - should use _printObject
	buf.Reset()
	PrintObject([]interface{}{buf, makeStringObject("test object string")})
	if got := buf.String(); got != "test object string" {
		t.Errorf("PrintObject(String) = %q; want 'test object string'", got)
	}

	// Case 2: Println(String object)
	buf.Reset()
	PrintlnObject([]interface{}{buf, makeStringObject("hello object")})
	if got := buf.String(); got != "hello object\n" {
		t.Errorf("PrintlnObject(String) = %q; want 'hello object\\n'", got)
	}

	// Case 3: Print normal object with fields
	buf.Reset()
	globals.InitStringPool()
	className := "MyClass"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable = map[string]object.Field{
		"field1": {Ftype: types.Int, Fvalue: int64(123)},
		"field2": {Ftype: types.Int, Fvalue: int64(456)},
	}
	// className suffix will be "MyClass"
	PrintObject([]interface{}{buf, obj})
	got := buf.String()
	// field order in map is random, so we check for components
	if !(got == "MyClass{field1=123, field2=456}" || got == "MyClass{field2=456, field1=123}") {
		t.Errorf("PrintObject(Object) = %q; unexpected format", got)
	}

	// Case 4: Println normal object
	buf.Reset()
	PrintlnObject([]interface{}{buf, obj})
	got = buf.String()
	if !(got == "MyClass{field1=123, field2=456}\n" || got == "MyClass{field2=456, field1=123}\n") {
		t.Errorf("PrintlnObject(Object) = %q; unexpected format", got)
	}

	// Case 5: Error - non-io.Writer
	ret := PrintObject([]interface{}{"not-a-writer", obj})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("PrintObject with invalid writer did not return ghelpers.GErrBlk, got %T", ret)
	}

	// Case 6: Print Object that is actually a String
	buf.Reset()
	strObj := makeStringObject("internal string")
	PrintObject([]interface{}{buf, strObj})
	if got := buf.String(); got != "internal string" {
		t.Errorf("PrintObject(StringObj) = %q; want 'internal string'", got)
	}
}

func TestPrintStreamPrintLinkedList(t *testing.T) {
	buf := new(bytes.Buffer)
	globals.InitStringPool()

	// Load the necessary gfunctions
	javaUtil.Load_Util_LinkedList()

	// Create a LinkedList object
	llObj := object.MakeEmptyObjectWithClassName(&types.ClassNameLinkedList)
	llObj.FieldTable = make(map[string]object.Field)

	// Case 1: Empty LinkedList
	// If we don't put a "value" field, _printLinkedList should return an error
	ret := PrintObject([]interface{}{buf, llObj})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("PrintObject(Empty LL without value field) should return ghelpers.GErrBlk, got %T", ret)
	}

	// Correctly initialize LinkedList
	ghelpers.Invoke("java/util/LinkedList.<init>()V", []interface{}{llObj})

	// Case 2: LinkedList with elements
	ghelpers.Invoke("java/util/LinkedList.add(Ljava/lang/Object;)Z", []interface{}{llObj, makeStringObject("A")})
	ghelpers.Invoke("java/util/LinkedList.add(Ljava/lang/Object;)Z", []interface{}{llObj, makeStringObject("B")})

	buf.Reset()
	PrintObject([]interface{}{buf, llObj})
	if got := buf.String(); got != "[A, B]" {
		t.Errorf("PrintObject(LinkedList) = %q; want '[A, B]'", got)
	}

	buf.Reset()
	PrintlnObject([]interface{}{buf, llObj})
	if got := buf.String(); got != "[A, B]\n" {
		t.Errorf("PrintlnObject(LinkedList) = %q; want '[A, B]\\n'", got)
	}

	// Case 3: LinkedList error cases
	// non-io.Writer (already covered by PrintObject/PrintlnObject initial check, but _printLinkedList has its own check)
	// Actually PrintObject checks writer if params[1] is null, but if not null, it passes it to _printObject or _printLinkedList.
	// _printLinkedList checks writer too.

	ret = PrintObject([]interface{}{"not-a-writer", llObj})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("_printLinkedList with invalid writer did not return ghelpers.GErrBlk")
	}

	// Case 4: null LinkedList
	buf.Reset()
	// To trigger null branch in _printLinkedList, we need params[1] to be *object.Object but representing Null
	PrintObject([]interface{}{buf, object.Null})
	// Actually PrintObject/PrintlnObject check object.IsNull(params[1]) and handle it themselves.
	// So they only call _printLinkedList if it is NOT null and matches classNameLinkedList.
	// This means the null branch in _printLinkedList might be unreachable via PrintObject/PrintlnObject.
	// But we can test it if we call _printLinkedList directly from another gfunction file (if it was exported).
	// Since it's not exported, we can only reach it if someone else calls it.
}

func TestPrintStreamPrintObjectNull(t *testing.T) {
	buf := new(bytes.Buffer)

	// print(null Object)
	buf.Reset()
	err := PrintObject([]interface{}{buf, nil})
	if err != nil {
		t.Fatalf("PrintObject with nil returned error: %v", err)
	}
	if got := buf.String(); got != "null" {
		t.Errorf("PrintObject(nil) output = %q; want 'null'", got)
	}

	// println(null Object)
	buf.Reset()
	err = PrintlnObject([]interface{}{buf, nil})
	if err != nil {
		t.Fatalf("PrintlnObject with nil returned error: %v", err)
	}
	if got := buf.String(); got != "null\n" {
		t.Errorf("PrintlnObject(nil) output = %q; want 'null\\n'", got)
	}
}
