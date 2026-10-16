package javaIo

import (
	"bytes"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/gfunction/javaMath"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"strconv"
	"strings"
	"sync"
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

func TestPrintStreamErrorCases(t *testing.T) {
	invalidWriter := "not-a-writer"
	obj := makeStringObject("test")

	tests := []struct {
		name string
		f    func([]interface{}) interface{}
		args []interface{}
	}{
		{"PrintlnV", PrintlnV, []interface{}{invalidWriter}},
		{"PrintlnChar", PrintlnChar, []interface{}{invalidWriter, int64('A')}},
		{"PrintlnBIS", PrintlnBIS, []interface{}{invalidWriter, int64(1)}},
		{"PrintlnBoolean", PrintlnBoolean, []interface{}{invalidWriter, int64(1)}},
		{"PrintlnLong", PrintlnLong, []interface{}{invalidWriter, int64(1)}},
		{"PrintlnDouble", PrintlnDouble, []interface{}{invalidWriter, 1.0}},
		{"PrintlnFloat", PrintlnFloat, []interface{}{invalidWriter, 1.0}},
		{"PrintChar", PrintChar, []interface{}{invalidWriter, int64('A')}},
		{"PrintBIS", PrintBIS, []interface{}{invalidWriter, int64(1)}},
		{"PrintBoolean", PrintBoolean, []interface{}{invalidWriter, int64(1)}},
		{"PrintLong", PrintLong, []interface{}{invalidWriter, int64(1)}},
		{"PrintDouble", PrintDouble, []interface{}{invalidWriter, 1.0}},
		{"PrintFloat", PrintFloat, []interface{}{invalidWriter, 1.0}},
		{"Printf", Printf, []interface{}{invalidWriter, makeStringObject("%s"), nil}},
		{"PrintString_Writer", PrintString, []interface{}{invalidWriter, obj}},
		{"PrintString_Param", PrintString, []interface{}{new(bytes.Buffer), "not-an-object"}},
		{"PrintObject_Writer", PrintObject, []interface{}{invalidWriter, obj}},
		{"PrintlnObject_Writer", PrintlnObject, []interface{}{invalidWriter, obj}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret := tt.f(tt.args)
			if _, ok := ret.(*ghelpers.GErrBlk); !ok {
				t.Errorf("%s did not return GErrBlk, got %T", tt.name, ret)
			}
		})
	}
}

func TestPrintStreamBISInt8(t *testing.T) {
	buf := new(bytes.Buffer)

	// PrintlnBIS with int8
	buf.Reset()
	PrintlnBIS([]interface{}{buf, int8(42)})
	if got := buf.String(); got != "42\n" {
		t.Errorf("PrintlnBIS(int8(42)) = %q; want '42\\n'", got)
	}

	// PrintBIS with int8
	buf.Reset()
	PrintBIS([]interface{}{buf, int8(42)})
	if got := buf.String(); got != "42" {
		t.Errorf("PrintBIS(int8(42)) = %q; want '42'", got)
	}
}

func TestPrintStreamSpecialObjects(t *testing.T) {
	buf := new(bytes.Buffer)
	globals.InitStringPool()

	// 1. BigDecimal
	javaMath.Load_Math_Big_Decimal()
	bdObj := object.MakeEmptyObjectWithClassName(&types.ClassNameBigDecimal)
	ghelpers.Invoke("java/math/BigDecimal.<init>(Ljava/lang/String;)V", []interface{}{bdObj, makeStringObject("123.45")})
	buf.Reset()
	PrintObject([]interface{}{buf, bdObj})
	if got := buf.String(); got != "123.45" {
		t.Errorf("PrintObject(BigDecimal) = %q; want '123.45'", got)
	}

	// 2. BigInteger
	javaMath.Load_Math_Big_Integer()
	biObj := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ghelpers.Invoke("java/math/BigInteger.<init>(Ljava/lang/String;)V", []interface{}{biObj, makeStringObject("987654321")})
	buf.Reset()
	PrintObject([]interface{}{buf, biObj})
	if got := buf.String(); got != "987654321" {
		t.Errorf("PrintObject(BigInteger) = %q; want '987654321'", got)
	}

	// 4. ThreadState (needed for Thread)
	javaLang.Load_Lang_Thread_State()
	tsObj := object.MakeEmptyObjectWithClassName(&types.ClassNameThreadState)
	tsObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(1)} // RUNNABLE
	// To make it look like a real Thread.State object that ThreadStateToString expects
	tsObj.FieldTable["name"] = object.Field{Ftype: "Ljava/lang/String;", Fvalue: makeStringObject("RUNNABLE")}
	// ThreadState is an enum, toString should return its name.
	buf.Reset()
	PrintObject([]interface{}{buf, tsObj})
	if got := buf.String(); got != "RUNNABLE" {
		t.Errorf("PrintObject(ThreadState) = %q; want 'RUNNABLE'", got)
	}

	// 3. Thread
	javaLang.Load_Lang_Thread()
	statics.AddStatic("java/lang/Thread.MIN_PRIORITY", statics.Static{Type: types.Int, Value: int64(1)})
	statics.AddStatic("java/lang/Thread.MAX_PRIORITY", statics.Static{Type: types.Int, Value: int64(10)})
	statics.AddStatic("java/lang/Thread.NORM_PRIORITY", statics.Static{Type: types.Int, Value: int64(5)})

	threadObj := object.MakeEmptyObjectWithClassName(&types.ClassNameThread)
	threadObj.KlassName = types.StringPoolThreadIndex
	threadObj.ThMutex = new(sync.RWMutex)
	threadObj.FieldTable["ID"] = object.Field{Ftype: types.Long, Fvalue: int64(1)}
	threadObj.FieldTable["name"] = object.Field{Ftype: "Ljava/lang/String;", Fvalue: makeStringObject("TestThread")}
	threadObj.FieldTable["priority"] = object.Field{Ftype: types.Int, Fvalue: int64(5)}
	threadObj.FieldTable["state"] = object.Field{Ftype: types.Int, Fvalue: int64(1)} // RUNNABLE
	threadObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(0)} // Unused but might be checked

	buf.Reset()
	ret := PrintObject([]interface{}{buf, threadObj})
	if errBlk, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("PrintObject(Thread) returned error: %s", errBlk.ErrMsg)
	}
	got := buf.String()
	if !strings.Contains(got, "Thread[ID=1, Name=TestThread, Priority=5, State=RUNNABLE]") {
		t.Errorf("PrintObject(Thread) = %q; unexpected format", got)
	}
}

func TestPrintStreamPrintStringNull(t *testing.T) {
	buf := new(bytes.Buffer)

	// PrintString with null object
	buf.Reset()
	PrintString([]interface{}{buf, object.Null})
	if got := buf.String(); got != "null" {
		t.Errorf("PrintString(null) = %q; want 'null'", got)
	}

	// PrintlnString with null object
	buf.Reset()
	PrintlnString([]interface{}{buf, object.Null})
	if got := buf.String(); got != "null\n" {
		t.Errorf("PrintlnString(null) = %q; want 'null\\n'", got)
	}
}

func TestPrintStreamWriteFromByteArray(t *testing.T) {
	buf := new(bytes.Buffer)

	// Case 1: Success
	buf.Reset()
	str := "Hello World"
	strObj := makeStringObject(str)
	ret := printstreamWriteFromByteArray([]interface{}{buf, strObj})
	if ret != nil {
		t.Errorf("printstreamWriteFromByteArray returned error: %v", ret)
	}
	if got := buf.String(); got != str {
		t.Errorf("printstreamWriteFromByteArray output = %q; want %q", got, str)
	}

	// Case 2: Invalid writer
	ret = printstreamWriteFromByteArray([]interface{}{"not-a-writer", strObj})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("printstreamWriteFromByteArray with invalid writer did not return GErrBlk, got %T", ret)
	}

	// Case 3: Invalid byteArrayObj
	ret = printstreamWriteFromByteArray([]interface{}{buf, "not-an-object"})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("printstreamWriteFromByteArray with invalid object did not return GErrBlk, got %T", ret)
	}
}

func TestPrintStreamConstructors(t *testing.T) {
	Load_Io_PrintStream()

	// Check registration of all PrintStream.<init> signatures
	trappedSignatures := []struct {
		sig   string
		slots int
	}{
		{"java/io/PrintStream.<init>(Ljava/io/OutputStream;ZLjava/lang/String;)V", 3},
		{"java/io/PrintStream.<init>(Ljava/io/OutputStream;ZLjava/nio/charset/Charset;)V", 3},
		{"java/io/PrintStream.<init>(Ljava/io/File;Ljava/lang/String;)V", 2},
		{"java/io/PrintStream.<init>(Ljava/io/File;Ljava/nio/charset/Charset;)V", 2},
		{"java/io/PrintStream.<init>(Ljava/lang/String;Ljava/lang/String;)V", 2},
		{"java/io/PrintStream.<init>(Ljava/lang/String;Ljava/nio/charset/Charset;)V", 2},
	}

	for _, tt := range trappedSignatures {
		gm, ok := ghelpers.MethodSignatures[tt.sig]
		if !ok {
			t.Errorf("Signature %s not registered", tt.sig)
			continue
		}
		if gm.ParamSlots != tt.slots {
			t.Errorf("Signature %s ParamSlots = %d; want %d", tt.sig, gm.ParamSlots, tt.slots)
		}
		// Trapped constructors should return a trap block / error when invoked
		res := gm.GFunction(nil)
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Signature %s did not invoke TrapFunction, returned %v", tt.sig, res)
		}
	}

	implementedSignatures := []struct {
		sig   string
		slots int
	}{
		{"java/io/PrintStream.<init>(Ljava/io/OutputStream;)V", 1},
		{"java/io/PrintStream.<init>(Ljava/io/OutputStream;Z)V", 2},
		{"java/io/PrintStream.<init>(Ljava/io/File;)V", 1},
		{"java/io/PrintStream.<init>(Ljava/lang/String;)V", 1},
	}

	for _, tt := range implementedSignatures {
		gm, ok := ghelpers.MethodSignatures[tt.sig]
		if !ok {
			t.Errorf("Signature %s not registered", tt.sig)
			continue
		}
		if gm.ParamSlots != tt.slots {
			t.Errorf("Signature %s ParamSlots = %d; want %d", tt.sig, gm.ParamSlots, tt.slots)
		}
	}

	// Test initPrintStreamOutputStream & initPrintStreamOutputStreamBoolean
	psClassName := "java/io/PrintStream"
	fosClassName := "java/io/FileOutputStream"
	psObj := object.MakeEmptyObjectWithClassName(&psClassName)
	underlying := object.MakeEmptyObjectWithClassName(&fosClassName)
	underlying.FieldTable[ghelpers.FilePath] = object.Field{Ftype: "[B", Fvalue: []byte("/dev/null")}
	underlying.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: nil}

	res := initPrintStreamOutputStream([]any{psObj, underlying})
	if res != nil {
		t.Errorf("initPrintStreamOutputStream error: %v", res)
	}
	if psObj.FieldTable[ghelpers.FilePath].Fvalue == nil {
		t.Errorf("initPrintStreamOutputStream did not copy FilePath")
	}

	// Test with null underlying OutputStream
	res = initPrintStreamOutputStream([]any{psObj, nil})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("initPrintStreamOutputStream with nil did not return GErrBlk")
	}

	psObj2 := object.MakeEmptyObjectWithClassName(&psClassName)
	res = initPrintStreamOutputStreamBoolean([]any{psObj2, underlying, int64(1)})
	if res != nil {
		t.Errorf("initPrintStreamOutputStreamBoolean error: %v", res)
	}
	if psObj2.FieldTable[ghelpers.FilePath].Fvalue == nil {
		t.Errorf("initPrintStreamOutputStreamBoolean did not copy FilePath")
	}
	if fldOut, ok := psObj2.FieldTable["out"]; !ok || fldOut.Ftype != "java/io/OutputStream" || fldOut.Fvalue != underlying {
		t.Errorf("initPrintStreamOutputStreamBoolean out field = %+v; want Ftype java/io/OutputStream and Fvalue underlying", fldOut)
	}

	res = initPrintStreamOutputStreamBoolean([]any{psObj2, nil, int64(1)})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("initPrintStreamOutputStreamBoolean with nil did not return GErrBlk")
	}
}

func TestPrintStreamWrappingByteArrayOutputStream(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Io_PrintStream()
	Load_Io_ByteArrayOutputStream()

	baosClassName := "java/io/ByteArrayOutputStream"
	baosObj := object.MakeEmptyObjectWithClassName(&baosClassName)
	ByteArrayOutputStreamInit([]any{baosObj})

	psClassName := "java/io/PrintStream"
	psObj := object.MakeEmptyObjectWithClassName(&psClassName)

	// new PrintStream(byteArrayOutputStream, true)
	res := initPrintStreamOutputStreamBoolean([]any{psObj, baosObj, true})
	if res != nil {
		t.Fatalf("initPrintStreamOutputStreamBoolean failed: %v", res)
	}

	if fldOut, ok := psObj.FieldTable["out"]; !ok || fldOut.Ftype != "java/io/OutputStream" || fldOut.Fvalue != baosObj {
		t.Errorf("PrintStream out field = %+v; want Ftype java/io/OutputStream and Fvalue baosObj", fldOut)
	}

	// Print formatted text and string line to PrintStream
	fmtObj := makeStringObject("Count: %d, Message: %s")
	intArg := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(7))
	strArg := makeStringObject("ok")
	argsObj := object.MakeArrayFromRawArray([]*object.Object{intArg, strArg})

	Printf([]any{psObj, fmtObj, argsObj})
	PrintlnV([]any{psObj})
	PrintlnString([]any{psObj, makeStringObject("Second line")})

	// Check ByteArrayOutputStream content
	strResult := ByteArrayOutputStreamToString([]any{baosObj})
	strObj, ok := strResult.(*object.Object)
	if !ok {
		t.Fatalf("ByteArrayOutputStreamToString returned %T; want *object.Object", strResult)
	}
	gotStr := object.GoStringFromStringObject(strObj)
	expected := "Count: 7, Message: ok\nSecond line\n"
	if gotStr != expected {
		t.Errorf("Captured ByteArrayOutputStream output = %q; want %q", gotStr, expected)
	}

	// Also test initPrintStreamOutputStream (without boolean)
	baosObj2 := object.MakeEmptyObjectWithClassName(&baosClassName)
	ByteArrayOutputStreamInit([]any{baosObj2})
	psObj2 := object.MakeEmptyObjectWithClassName(&psClassName)
	res2 := initPrintStreamOutputStream([]any{psObj2, baosObj2})
	if res2 != nil {
		t.Fatalf("initPrintStreamOutputStream failed: %v", res2)
	}
	if fldOut, ok := psObj2.FieldTable["out"]; !ok || fldOut.Ftype != "java/io/OutputStream" || fldOut.Fvalue != baosObj2 {
		t.Errorf("PrintStream out field = %+v; want Ftype java/io/OutputStream and Fvalue baosObj2", fldOut)
	}

	PrintlnString([]any{psObj2, makeStringObject("Direct write")})
	strResult2 := ByteArrayOutputStreamToString([]any{baosObj2})
	strObj2 := strResult2.(*object.Object)
	if got := object.GoStringFromStringObject(strObj2); got != "Direct write\n" {
		t.Errorf("Captured output = %q; want 'Direct write\\n'", got)
	}
}

func TestPrintStreamWithObjectReceiver(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Io_PrintStream()

	// Create a buffer and wrap it in an os.File or custom writer in PrintStream object
	buf := new(bytes.Buffer)
	psClassName := "java/io/PrintStream"
	psObj := object.MakeEmptyObjectWithClassName(&psClassName)
	psObj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: buf}

	// Test Printf with *object.Object as params[0]
	fmtObj := makeStringObject("Hello %s %d")
	strArg := makeStringObject("world")
	intArg := object.MakePrimitiveObject("java/lang/Integer", types.Int, int64(42))

	argsObj := object.MakeArrayFromRawArray([]*object.Object{strArg, intArg})

	ret := Printf([]any{psObj, fmtObj, argsObj})
	if ret != psObj {
		t.Errorf("Printf should return params[0] (psObj), got %v", ret)
	}
	if got := buf.String(); got != "Hello world 42" {
		t.Errorf("Printf output = %q; want 'Hello world 42'", got)
	}

	// Test PrintString, PrintlnString, PrintlnV, PrintlnLong with psObj
	buf.Reset()
	PrintlnString([]any{psObj, makeStringObject("Line 1")})
	PrintString([]any{psObj, makeStringObject("Line 2")})
	PrintlnV([]any{psObj})
	PrintlnLong([]any{psObj, int64(99)})

	expected := "Line 1\nLine 2\n99\n"
	if got := buf.String(); got != expected {
		t.Errorf("Print operations output = %q; want %q", got, expected)
	}

	// Test PrintFlush
	ret = PrintFlush([]any{psObj})
	if ret != nil {
		t.Errorf("PrintFlush returned error: %v", ret)
	}

	// Test error when psObj has no writer or out field
	emptyPsObj := object.MakeEmptyObjectWithClassName(&psClassName)
	ret = Printf([]any{emptyPsObj, makeStringObject("test"), object.MakeArrayFromRawArray([]*object.Object{})})
	if _, isErr := ret.(*ghelpers.GErrBlk); !isErr {
		t.Errorf("Printf with emptyPsObj should return GErrBlk, got %v", ret)
	}
}
