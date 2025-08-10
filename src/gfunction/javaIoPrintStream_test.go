package gfunction_test

import (
	"bytes"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/types"
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
	gfunction.PrintlnV([]interface{}{buf})
	if got := buf.String(); got != "\n" {
		t.Errorf("PrintlnV() = %q; want newline", got)
	}

	// println(String)
	buf.Reset()
	gfunction.PrintlnString([]interface{}{buf, makeStringObject("hello")})
	if got := buf.String(); got != "hello\n" {
		t.Errorf("PrintlnString() = %q; want 'hello\\n'", got)
	}

	// println(byte), println(int), println(short) all use PrintlnBIS
	for _, val := range []int64{65, -1, 0} {
		buf.Reset()
		gfunction.PrintlnBIS([]interface{}{buf, val})
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
	gfunction.PrintlnBoolean([]interface{}{buf, int64(1)})
	if got := buf.String(); got != "true\n" {
		t.Errorf("PrintlnBoolean(true) = %q; want 'true\\n'", got)
	}
	buf.Reset()
	gfunction.PrintlnBoolean([]interface{}{buf, int64(0)})
	if got := buf.String(); got != "false\n" {
		t.Errorf("PrintlnBoolean(false) = %q; want 'false\\n'", got)
	}

	// println(long)
	buf.Reset()
	gfunction.PrintlnLong([]interface{}{buf, int64(1234567890)})
	if got := buf.String(); got != "1234567890\n" {
		t.Errorf("PrintlnLong() = %q; want '1234567890\\n'", got)
	}

	// println(double)
	buf.Reset()
	gfunction.PrintlnDouble([]interface{}{buf, 3.1415926535})
	if got := buf.String(); got != "3.1415926535\n" {
		t.Errorf("PrintlnDouble() = %q; want '3.1415926535\\n'", got)
	}

	// println(float)
	buf.Reset()
	gfunction.PrintlnFloat([]interface{}{buf, 2.71828})
	if got := buf.String(); got != "2.71828\n" {
		t.Errorf("PrintlnFloat() = %q; want '2.71828\\n'", got)
	}

	// println(char)
	buf.Reset()
	gfunction.PrintlnChar([]interface{}{buf, int64('Z')})
	if got := buf.String(); got != "Z\n" {
		t.Errorf("PrintlnChar() = %q; want 'Z\\n'", got)
	}
}

func TestPrintStreamPrintFunctions(t *testing.T) {
	buf := new(bytes.Buffer)

	// print(String)
	buf.Reset()
	gfunction.PrintString([]interface{}{buf, makeStringObject("test string")})
	if got := buf.String(); got != "test string" {
		t.Errorf("PrintString() = %q; want 'test string'", got)
	}

	// print(byte), print(int), print(short) use PrintBIS
	for _, val := range []int64{65, -2, 0} {
		buf.Reset()
		gfunction.PrintBIS([]interface{}{buf, val})
		want := strconv.FormatInt(val, 10)
		if got := buf.String(); got != want {
			t.Errorf("PrintBIS(%d) = %q; want %q", val, got, want)
		}
	}

	// print(boolean)
	buf.Reset()
	gfunction.PrintBoolean([]interface{}{buf, int64(1)})
	if got := buf.String(); got != "true" {
		t.Errorf("PrintBoolean(true) = %q; want 'true'", got)
	}
	buf.Reset()
	gfunction.PrintBoolean([]interface{}{buf, int64(0)})
	if got := buf.String(); got != "false" {
		t.Errorf("PrintBoolean(false) = %q; want 'false'", got)
	}

	// print(long)
	buf.Reset()
	gfunction.PrintLong([]interface{}{buf, int64(9876543210)})
	if got := buf.String(); got != "9876543210" {
		t.Errorf("PrintLong() = %q; want '9876543210'", got)
	}

	// print(double)
	buf.Reset()
	gfunction.PrintDouble([]interface{}{buf, 6.62607015})
	if got := buf.String(); got != "6.62607015" {
		t.Errorf("PrintDouble() = %q; want '6.62607015'", got)
	}

	// print(float)
	buf.Reset()
	gfunction.PrintFloat([]interface{}{buf, 1.41421})
	if got := buf.String(); got != "1.41421" {
		t.Errorf("PrintFloat() = %q; want '1.41421'", got)
	}

	// print(char)
	buf.Reset()
	gfunction.PrintChar([]interface{}{buf, int64('Y')})
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
	ret := gfunction.Printf(params)
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

func TestPrintStreamPrintObjectNull(t *testing.T) {
	buf := new(bytes.Buffer)

	// print(null Object)
	buf.Reset()
	err := gfunction.PrintObject([]interface{}{buf, nil})
	if err != nil {
		t.Fatalf("PrintObject with nil returned error: %v", err)
	}
	if got := buf.String(); got != "null" {
		t.Errorf("PrintObject(nil) output = %q; want 'null'", got)
	}

	// println(null Object)
	buf.Reset()
	err = gfunction.PrintlnObject([]interface{}{buf, nil})
	if err != nil {
		t.Fatalf("PrintlnObject with nil returned error: %v", err)
	}
	if got := buf.String(); got != "null\n" {
		t.Errorf("PrintlnObject(nil) output = %q; want 'null\\n'", got)
	}
}
