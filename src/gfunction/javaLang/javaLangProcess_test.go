package javaLang

import (
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestProcess_BasicMethods(t *testing.T) {
	globals.InitStringPool()

	// Use current process ID for testing
	currentPid := int64(os.Getpid())

	// Create a mock Process object
	procObj := object.MakeOneFieldObject("java/lang/Process", "pid", types.Int, currentPid)

	// Test pid()
	res := processPid([]interface{}{procObj})
	if res.(int64) != currentPid {
		t.Errorf("pid() expected %d, got %v", currentPid, res)
	}

	// Test isAlive()
	res = processIsAlive([]interface{}{procObj})
	if res != types.JavaBoolTrue {
		t.Errorf("isAlive() expected true for current process")
	}

	// Test toHandle()
	res = processToHandle([]interface{}{procObj})
	handleObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("toHandle() did not return an object")
	}
	hPid := handleObj.FieldTable["pid"].Fvalue.(int64)
	if hPid != currentPid {
		t.Errorf("ProcessHandle PID expected %d, got %d", currentPid, hPid)
	}
}

func TestProcess_InvalidObject(t *testing.T) {
	globals.InitStringPool()

	// Object without pid field
	procObj := object.MakeEmptyObject()

	res := processPid([]interface{}{procObj})
	if _, ok := res.(*ghelpers.GErrBlk); !ok {
		t.Errorf("pid() should return error for object without pid field")
	}
}

func TestProcessGetInputStream(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process()
	Load_Lang_Process_Builder()

	sig := "java/lang/Process.getInputStream()Ljava/io/InputStream;"
	gmeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gmeth.ParamSlots != 0 {
		t.Errorf("ParamSlots expected 0, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(processGetInputStream).Pointer() {
		t.Errorf("GFunction does not point to processGetInputStream")
	}

	// Null receiver check
	if err, ok := processGetInputStream([]any{nil}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil receiver")
	}

	// Non-object receiver check
	if err, ok := processGetInputStream([]any{"not-an-object"}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for invalid receiver type")
	}

	// Mock Process object check
	procClassName := "java/lang/Process"
	mockProc := object.MakeEmptyObjectWithClassName(&procClassName)
	inRes := processGetInputStream([]any{mockProc})
	inObj, ok := inRes.(*object.Object)
	if !ok || inObj == nil {
		t.Fatalf("expected *object.Object from getInputStream(), got %T (%v)", inRes, inRes)
	}

	// Calling getInputStream() again returns cached stream
	inRes2 := processGetInputStream([]any{mockProc})
	if inRes2 != inObj {
		t.Errorf("expected getInputStream() to return cached stream object")
	}

	// Subprocess execution with ProcessBuilder and reading standard output
	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)
	processBuilderInitStringArray([]any{pbObj, []string{"echo", "hello jacobin"}})

	startRes := processBuilderStart([]any{pbObj})
	procStarted, ok := startRes.(*object.Object)
	if !ok || procStarted == nil {
		t.Fatalf("processBuilderStart expected *object.Object, got %T (%v)", startRes, startRes)
	}

	streamRes := processGetInputStream([]any{procStarted})
	streamObj, ok := streamRes.(*object.Object)
	if !ok || streamObj == nil {
		t.Fatalf("getInputStream expected *object.Object, got %T (%v)", streamRes, streamRes)
	}

	fhField, ok := streamObj.FieldTable[ghelpers.FileHandle]
	if !ok || fhField.Fvalue == nil {
		t.Fatalf("expected FileHandle in stream object")
	}
	reader, ok := fhField.Fvalue.(io.Reader)
	if !ok {
		t.Fatalf("expected io.Reader in stream FileHandle, got %T", fhField.Fvalue)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("error reading from process input stream: %v", err)
	}
	if !strings.Contains(string(data), "hello jacobin") {
		t.Errorf("expected output to contain 'hello jacobin', got %q", string(data))
	}

	// ProcessBuilder with redirectErrorStream(true)
	pbObj2 := object.MakeEmptyObjectWithClassName(&pbClassName)
	processBuilderInitStringArray([]any{pbObj2, []string{"sh", "-c", "echo err_msg >&2"}})
	processBuilderRedirectErrorStreamSet([]any{pbObj2, true})

	startRes2 := processBuilderStart([]any{pbObj2})
	procStarted2, ok := startRes2.(*object.Object)
	if !ok || procStarted2 == nil {
		t.Fatalf("processBuilderStart expected *object.Object, got %T (%v)", startRes2, startRes2)
	}

	streamRes2 := processGetInputStream([]any{procStarted2})
	streamObj2, ok := streamRes2.(*object.Object)
	if !ok || streamObj2 == nil {
		t.Fatalf("getInputStream expected *object.Object, got %T (%v)", streamRes2, streamRes2)
	}

	reader2, ok := streamObj2.FieldTable[ghelpers.FileHandle].Fvalue.(io.Reader)
	if !ok {
		t.Fatalf("expected io.Reader in stream FileHandle")
	}
	data2, err := io.ReadAll(reader2)
	if err != nil {
		t.Fatalf("error reading from process input stream with redirected error: %v", err)
	}
	if !strings.Contains(string(data2), "err_msg") {
		t.Errorf("expected output to contain 'err_msg', got %q", string(data2))
	}
}
