/*
TestThreadActiveCount/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
*/

package gfunction

import (
	"bytes"
	"container/list"
	"io"
	"jacobin/src/exceptions"
	"os"
	"sync"
	"testing"

	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/thread"
	"jacobin/src/types"
)

var initOnce sync.Once

// instantiateForThreadTest is a minimal stand-in for jvm.Instantiate to avoid circular imports.
// It satisfies the globals.FuncInstantiateClass signature used by gfunctions during tests.
func instantiateForThreadTest(name string, _ *list.List) (any, error) {
	o := object.MakeEmptyObject()
	o.KlassName = stringPool.GetStringIndex(&name)
	return o, nil
}

func ensureInit() {
	initOnce.Do(func() {
		globals.InitGlobals("test")
		globals.InitStringPool()
		gr := globals.GetGlobalRef()
		gr.Threads = make(map[int]interface{})
		gr.ThreadGroups = make(map[string]interface{})

		gr.FuncFillInStackTrace = FillInStackTrace
		gr.FuncInvokeGFunction = Invoke
		gr.FuncThrowException = exceptions.ThrowExNil
		// Set a local fake instantiator to avoid importing the jvm package in tests
		gr.FuncInstantiateClass = instantiateForThreadTest
		InitializeGlobalThreadGroups()
		Load_Lang_Thread_Group()
		Load_Lang_Thread()

	})
}

func init() {
	ensureInit()
}

// Helpers
func makeJavaString(s string) *object.Object {
	return object.StringObjectFromGoString(s)
}

func makeEmptyThreadWithIDName(id int64, name string) *object.Object {
	ensureInit()
	InitializeGlobalThreadGroups()
	if globals.GetGlobalRef().ThreadGroups["main"] == nil {
		Load_Lang_Thread_Group()
	}
	t := threadCreateNoarg(nil).(*object.Object)
	t.FieldTable["ID"] = object.Field{Ftype: types.Int, Fvalue: id}
	// Use Java byte array for name when needed by getters
	jbytes := object.JavaByteArrayFromGoString(name)
	t.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: jbytes}
	return t
}

// ---- threadActiveCount ----
func TestThreadActiveCount(t *testing.T) {
	globals.InitGlobals("test")
	gr := globals.GetGlobalRef()
	gr.Threads = make(map[int]interface{})
	if got := threadActiveCount(nil).(int64); got != 0 {
		t.Fatalf("activeCount = %d; want 0", got)
	}
	gr.Threads[1] = makeEmptyThreadWithIDName(1, "Thread-1")
	gr.Threads[2] = makeEmptyThreadWithIDName(2, "Thread-2")
	if got := threadActiveCount(nil).(int64); got != 2 {
		t.Fatalf("activeCount = %d; want 2", got)
	}
}

// ---- threadClinit ----
func TestThreadClinit_SetsStatics(t *testing.T) {
	threadClinit(nil)
	min := statics.GetStaticValue("java/lang/Thread", "MIN_PRIORITY").(int64)
	norm := statics.GetStaticValue("java/lang/Thread", "NORM_PRIORITY").(int64)
	max := statics.GetStaticValue("java/lang/Thread", "MAX_PRIORITY").(int64)
	if min != int64(thread.MIN_PRIORITY) || norm != int64(thread.NORM_PRIORITY) || max != int64(thread.MAX_PRIORITY) {
		t.Fatalf("priority statics wrong: %d %d %d", min, norm, max)
	}
}

// ---- threadCreateFromPackageConstructor ----
func TestThreadCreateFromPackageConstructor_ArgCountError(t *testing.T) {
	res := threadCreateFromPackageConstructor([]interface{}{1})
	if _, ok := res.(*GErrBlk); !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
}

func TestThreadCreateFromPackageConstructor_TypeErrors(t *testing.T) {
	name := makeJavaString("worker")
	runnable := NewRunnable(object.JavaByteArrayFromGoString("c"), object.JavaByteArrayFromGoString("m"), object.JavaByteArrayFromGoString("()V"))
	// 0 wrong type
	res0 := threadCreateFromPackageConstructor([]interface{}{123, name, int64(5), runnable, int64(0), nil})
	if _, ok := res0.(*GErrBlk); !ok {
		t.Fatalf("param0 wrong type: expected error, got %T", res0)
	}
	// 1 wrong type
	res1 := threadCreateFromPackageConstructor([]interface{}{nil, 42, int64(5), runnable, int64(0), nil})
	if _, ok := res1.(*GErrBlk); !ok {
		t.Fatalf("param1 wrong type: expected error, got %T", res1)
	}
	// 2 wrong type
	res2 := threadCreateFromPackageConstructor([]interface{}{nil, name, "x", runnable, int64(0), nil})
	if _, ok := res2.(*GErrBlk); !ok {
		t.Fatalf("param2 wrong type: expected error, got %T", res2)
	}
	// 3 wrong type
	res3 := threadCreateFromPackageConstructor([]interface{}{nil, name, int64(5), 13, int64(0), nil})
	if _, ok := res3.(*GErrBlk); !ok {
		t.Fatalf("param3 wrong type: expected error, got %T", res3)
	}
	// 4 wrong type
	res4 := threadCreateFromPackageConstructor([]interface{}{nil, name, int64(5), runnable, 0, nil})
	if _, ok := res4.(*GErrBlk); !ok {
		t.Fatalf("param4 wrong type: expected error, got %T", res4)
	}
	// 5 wrong type
	res5 := threadCreateFromPackageConstructor([]interface{}{nil, name, int64(5), runnable, int64(0), 7})
	if _, ok := res5.(*GErrBlk); !ok {
		t.Fatalf("param5 wrong type: expected error, got %T", res5)
	}
}

func TestThreadCreateFromPackageConstructor_Success(t *testing.T) {
	globals.InitGlobals("test")
	InitializeGlobalThreadGroups()

	parent := threadGroupFake("grp")
	name := makeJavaString("worker")
	runnable := NewRunnable(object.JavaByteArrayFromGoString("c"), object.JavaByteArrayFromGoString("m"), object.JavaByteArrayFromGoString("()V"))
	res := threadCreateFromPackageConstructor([]interface{}{parent, name, int64(5), runnable, int64(0), nil})
	th, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	if th.FieldTable["threadgroup"].Fvalue != parent {
		t.Fatalf("thread group not set from param0")
	}
}

// ---- threadCreateNoarg ----
func TestThreadCreateNoarg_Defaults(t *testing.T) {
	obj := threadCreateNoarg(nil).(*object.Object)
	if obj.FieldTable["ID"].Ftype != types.Int {
		t.Fatalf("ID type wrong")
	}

	state := obj.FieldTable["state"].Fvalue.(*object.Object)
	if state.FieldTable["value"].Fvalue != NEW {
		t.Fatalf("state default wrong: %v", obj.FieldTable["state"].Fvalue)
	}
	if obj.FieldTable["daemon"].Fvalue.(int64) != types.JavaBoolFalse {
		t.Fatalf("daemon default not false")
	}
	if obj.FieldTable["threadgroup"].Fvalue == nil {
		t.Fatalf("threadgroup not set")
	}
}

// ---- threadCreateWithName ----
func TestThreadCreateWithName_TypeError(t *testing.T) {
	res := threadCreateWithName([]interface{}{42})
	g, ok := res.(*GErrBlk)
	if !ok || g.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IAE, got %T %+v", res, res)
	}
}

func TestThreadCreateWithName_SetsName(t *testing.T) {
	nm := makeJavaString("A")
	obj := threadCreateWithName([]interface{}{nm}).(*object.Object)
	val := obj.FieldTable["name"].Fvalue.([]types.JavaByte)
	if object.GoStringFromJavaByteArray(val) != "A" {
		t.Fatalf("name not set from string")
	}
}

// ---- threadCreateWithRunnable ----
func TestThreadCreateWithRunnable_SetsTask(t *testing.T) {
	r := NewRunnable(object.JavaByteArrayFromGoString("C"), object.JavaByteArrayFromGoString("run"), object.JavaByteArrayFromGoString("()V"))
	obj := threadCreateWithRunnable([]interface{}{r}).(*object.Object)
	if obj.FieldTable["task"].Fvalue != r {
		t.Fatalf("task not set")
	}
}

// ---- threadCreateWithRunnableAndName ----
func TestThreadCreateWithRunnableAndName_TypeError(t *testing.T) {
	// Provide a valid runnable to avoid panic on nil deref; make name wrong type
	r := NewRunnable(object.JavaByteArrayFromGoString("C"), object.JavaByteArrayFromGoString("run"), object.JavaByteArrayFromGoString("()V"))
	res := threadCreateWithRunnableAndName([]interface{}{r, 7})
	if _, ok := res.(*GErrBlk); !ok {
		t.Fatalf("expected error for name type, got %T", res)
	}
}

func TestThreadCreateWithRunnableAndName_SetsBoth(t *testing.T) {
	r := NewRunnable(object.JavaByteArrayFromGoString("C"), object.JavaByteArrayFromGoString("run"), object.JavaByteArrayFromGoString("()V"))
	name := makeJavaString("Zed")
	obj := threadCreateWithRunnableAndName([]interface{}{r, name}).(*object.Object)
	if obj.FieldTable["task"].Fvalue != r {
		t.Fatalf("task not set")
	}
	// name stored as GolangString holding []JavaByte; verify bytes content
	got := obj.FieldTable["name"].Fvalue.([]types.JavaByte)
	if object.GoStringFromJavaByteArray(got) != "Zed" {
		t.Fatalf("name not set correctly")
	}
}

// ---- threadCurrentThread ----
func TestThreadCurrentThread_ParamErrors(t *testing.T) {
	if _, ok := threadCurrentThread(nil).(*GErrBlk); !ok {
		t.Fatalf("expected len error")
	}
	if _, ok := threadCurrentThread([]interface{}{123}).(*GErrBlk); !ok {
		t.Fatalf("expected type error")
	}
}

func TestThreadCurrentThread_Success(t *testing.T) {
	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(8)
	f.Thread = 7
	frames.PushFrame(fs, f)
	th := makeEmptyThreadWithIDName(7, "X")
	globals.GetGlobalRef().Threads[7] = th
	res := threadCurrentThread([]interface{}{fs})
	if res != th {
		t.Fatalf("currentThread returned %T, want thread object", res)
	}
}

// ---- threadDumpStack ----
func captureStderr(run func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	run()
	_ = w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func TestThreadDumpStack_ParamErrors(t *testing.T) {
	if _, ok := threadDumpStack(nil).(*GErrBlk); !ok {
		t.Fatalf("expected len error")
	}
	if _, ok := threadDumpStack([]interface{}{42}).(*GErrBlk); !ok {
		t.Fatalf("expected type error")
	}
}

func TestThreadDumpStack_StrictJDK(t *testing.T) {
	gr := globals.GetGlobalRef()
	old := gr.StrictJDK
	gr.StrictJDK = true
	defer func() { gr.StrictJDK = old }()
	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(2)
	f.ClName = "C"
	f.MethName = "m"
	frames.PushFrame(fs, f)
	// Need a thread in globals map for header path? In StrictJDK header is fixed message, not using thread
	out := captureStderr(func() { _ = threadDumpStack([]interface{}{fs}) })
	if out == "" || out[:len("java.lang.Exception: Stack trace")] != "java.lang.Exception: Stack trace" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestThreadDumpStack_NonStrict_PrintsThreadAndPC(t *testing.T) {
	gr := globals.GetGlobalRef()
	old := gr.StrictJDK
	gr.StrictJDK = false
	defer func() { gr.StrictJDK = old }()

	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(2)
	f.ClName = "C"
	f.MethName = "m"
	f.PC = 3
	f.Thread = 9
	frames.PushFrame(fs, f)
	th := makeEmptyThreadWithIDName(9, "T9")
	globals.GetGlobalRef().Threads[9] = th
	out := captureStderr(func() { _ = threadDumpStack([]interface{}{fs}) })
	if out == "" || !bytes.Contains([]byte(out), []byte("Stack trace (thread")) || !bytes.Contains([]byte(out), []byte("PC: 3")) {
		t.Fatalf("unexpected output: %q", out)
	}
}

// ---- threadEnumerate ----
func TestThreadEnumerate_ParamError(t *testing.T) {
	if _, ok := threadEnumerate(nil).(*GErrBlk); !ok {
		t.Fatalf("expected len error")
	}
}

func TestThreadEnumerate_FillsArray(t *testing.T) {
	gr := globals.GetGlobalRef()
	gr.Threads = make(map[int]interface{})
	th1, th2 := makeEmptyThreadWithIDName(1, "A"), makeEmptyThreadWithIDName(2, "B")
	gr.Threads[1] = th1
	gr.Threads[2] = th2
	arr := make([]*object.Object, 3)
	arrObj := object.MakePrimitiveObject("java/lang/Object", types.RefArray, arr)
	cnt := threadEnumerate([]interface{}{arrObj}).(int)
	if cnt != 2 {
		t.Fatalf("enumerate count = %d; want 2", cnt)
	}
	if arr[0] == nil || arr[1] == nil {
		t.Fatalf("enumerate did not fill array")
	}
}

// ---- threadGetId ----
func TestThreadGetId_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadGetId(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	if _, ok := threadGetId([]interface{}{123}).(*GErrBlk); !ok {
		t.Fatalf("type error expected")
	}
	// Missing ID field
	dummy := object.MakeEmptyObject()
	if _, ok := threadGetId([]interface{}{dummy}).(*GErrBlk); !ok {
		t.Fatalf("missing field error expected")
	}
	// Wrong ID type
	wrongID := object.MakeEmptyObject()
	wrongID.FieldTable = map[string]object.Field{"ID": {Ftype: types.Int, Fvalue: "x"}}
	if _, ok := threadGetId([]interface{}{wrongID}).(*GErrBlk); !ok {
		t.Fatalf("ID type error expected")
	}
	// Success
	good := makeEmptyThreadWithIDName(42, "n")
	if got := threadGetId([]interface{}{good}).(int64); got != 42 {
		t.Fatalf("getId = %d; want 42", got)
	}
}

// ---- threadGetName ----
func TestThreadGetName_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadGetName(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	th := threadCreateWithName([]interface{}{makeJavaString("Neo")}).(*object.Object)
	strObj := threadGetName([]interface{}{th}).(*object.Object)
	if object.GoStringFromStringObject(strObj) != "Neo" {
		t.Fatalf("getName wrong")
	}
}

// ---- threadGetPriority ----
func TestThreadGetPriority_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadGetPriority(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	th := threadCreateNoarg(nil).(*object.Object)
	val := threadGetPriority([]interface{}{th}).(int64)
	if val == 0 {
		t.Fatalf("getPriority returned 0 unexpectedly")
	}
}

// ---- threadGetStackTrace ----
func TestThreadGetStackTrace_ParamErrors(t *testing.T) {
	if _, ok := threadGetStackTrace(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	if _, ok := threadGetStackTrace([]interface{}{123, 456}).(*GErrBlk); !ok {
		t.Fatalf("type error expected")
	}
}

func TestThreadGetStackTrace_Success(t *testing.T) {
	globals.GetGlobalRef().FuncInstantiateClass = instantiateForThreadTest
	// Initialize a minimal Method Area entry and frame metadata so that
	// StackTraceElement.initStackTraceElement can fetch class info safely
	classloader.InitMethodArea()

	// Insert a minimal Klass for a test class name
	clData := classloader.ClData{
		Name:            "",
		SuperclassIndex: types.ObjectPoolStringIndex,
		Module:          "test module",
		Pkg:             "",
		Interfaces:      nil,
		Fields:          nil,
		MethodTable:     nil,
		Attributes:      nil,
		SourceFile:      "testClass.java",
		Bootstraps:      nil,
		CP:              classloader.CPool{},
		Access:          classloader.AccessFlags{},
		ClInit:          0,
	}
	klass := classloader.Klass{Loader: "testLoader", Data: &clData}
	classloader.MethAreaInsert("java/testClass", &klass)

	// Build frame stack with a frame that references the test class
	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(1)
	f.Thread = 1
	f.ClName = "java/testClass"
	f.MethName = "java/testClass.test" // treat as JDK-like to avoid needing method lookup
	f.MethType = "()V"
	frames.PushFrame(fs, f)

	obj := threadGetStackTrace([]interface{}{fs, nil}).(*object.Object)
	if obj == nil {
		t.Fatalf("expected stack trace object")
	}
}

// ---- threadGetState ----
func TestThreadGetState_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadGetState(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}

	th := threadCreateNoarg(nil).(*object.Object)
	st := threadGetState([]interface{}{th}).(*object.Object)
	state := st.FieldTable["value"].Fvalue.(int)
	if state != NEW {
		t.Fatalf("getState = %d; want NEW", state)
	}
}

// ---- threadGetThreadGroup ----
func TestThreadGetThreadGroup_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadGetThreadGroup(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	th := threadCreateNoarg(nil).(*object.Object)
	grp := threadGetThreadGroup([]interface{}{th})
	if _, ok := grp.(*object.Object); !ok {
		t.Fatalf("getThreadGroup did not return object")
	}
	// Corrupt the field to force internal error path
	th.FieldTable["threadgroup"] = object.Field{Ftype: types.Ref, Fvalue: 7}
	if _, ok := threadGetThreadGroup([]interface{}{th}).(*GErrBlk); !ok {
		t.Fatalf("expected internal error when threadgroup is wrong type")
	}
}

// ---- threadIsInterrupted ----
func TestThreadIsInterrupted_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadIsInterrupted(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	if _, ok := threadIsInterrupted([]interface{}{42}).(*GErrBlk); !ok {
		t.Fatalf("type error expected")
	}
	th := threadCreateNoarg(nil).(*object.Object)
	val := threadIsInterrupted([]interface{}{th}).(int64)
	if val != types.JavaBoolFalse {
		t.Fatalf("isInterrupted default wrong: %d", val)
	}
}

// ---- threadRun ----
func TestThreadRun_ParamLenError(t *testing.T) {
	if _, ok := threadRun(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
}

func TestThreadRun_NilTaskReturnsNil(t *testing.T) {
	th := threadCreateNoarg(nil).(*object.Object)
	th.FieldTable["task"] = object.Field{Ftype: types.Ref, Fvalue: nil}
	if threadRun([]interface{}{th}) != nil {
		t.Fatalf("nil task should return nil")
	}
}

func TestThreadRun_NoSuchMethodError(t *testing.T) {
	globals.InitGlobals("test")

	normalStderr := os.Stderr
	_, writer, _ := os.Pipe()
	os.Stderr = writer

	// Build a runnable pointing to a class/method that cannot be found by FetchMethodAndCP
	runObj := NewRunnable(object.JavaByteArrayFromGoString("no/such/Class"),
		object.JavaByteArrayFromGoString("run"),
		object.JavaByteArrayFromGoString("()V"))

	InitializeGlobalThreadGroups()
	th := threadCreateNoarg(nil).(*object.Object)
	th.FieldTable["task"] = object.Field{Ftype: types.Ref, Fvalue: runObj}

	// Ensure ID set for any preparatory steps
	th.FieldTable["ID"] = object.Field{Ftype: types.Int, Fvalue: int64(2)}

	res := threadRun([]interface{}{th})

	_ = writer.Close()
	os.Stderr = normalStderr

	if g, ok := res.(*GErrBlk); !ok || g.ExceptionType != excNames.NoSuchMethodError {
		t.Fatalf("expected NoSuchMethodError, got %T %+v", res, res)
	}

}

// ---- threadSetName ----
func TestThreadSetName_ErrorsAndSuccess(t *testing.T) {
	if _, ok := threadSetName(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	if _, ok := threadSetName([]interface{}{42, makeJavaString("x")}).(*GErrBlk); !ok {
		t.Fatalf("type error expected for thread")
	}
	th := threadCreateNoarg(nil).(*object.Object)
	if _, ok := threadSetName([]interface{}{th, object.Null}).(*GErrBlk); !ok {
		t.Fatalf("null name should cause NPE")
	}
	badStr := object.MakeEmptyObject() // missing 'value'
	if _, ok := threadSetName([]interface{}{th, badStr}).(*GErrBlk); !ok {
		t.Fatalf("missing value field should error")
	}
	// Wrong type for value
	wrong := object.MakeEmptyObject()
	wrong.FieldTable = map[string]object.Field{"value": {Ftype: types.Int, Fvalue: int64(1)}}
	if _, ok := threadSetName([]interface{}{th, wrong}).(*GErrBlk); !ok {
		t.Fatalf("wrong value type should error")
	}
	// Success
	nm := makeJavaString("Z")
	if got := threadSetName([]interface{}{th, nm}); got != nil {
		t.Fatalf("expected nil on success, got %T", got)
	}
}

// ---- threadSetPriority ----
func TestThreadSetPriority_AllPaths(t *testing.T) {
	if _, ok := threadSetPriority(nil).(*GErrBlk); !ok {
		t.Fatalf("len error expected")
	}
	if _, ok := threadSetPriority([]interface{}{42, int64(5)}).(*GErrBlk); !ok {
		t.Fatalf("type error for thread expected")
	}
	th := threadCreateNoarg(nil).(*object.Object)
	if _, ok := threadSetPriority([]interface{}{th, "x"}).(*GErrBlk); !ok {
		t.Fatalf("type error for priority expected")
	}
	minP := statics.GetStaticValue("java/lang/Thread", "MIN_PRIORITY").(int64)
	maxP := statics.GetStaticValue("java/lang/Thread", "MAX_PRIORITY").(int64)
	if _, ok := threadSetPriority([]interface{}{th, minP - 1}).(*GErrBlk); !ok {
		t.Fatalf("below min should error")
	}
	if _, ok := threadSetPriority([]interface{}{th, maxP + 1}).(*GErrBlk); !ok {
		t.Fatalf("above max should error")
	}
	if got := threadSetPriority([]interface{}{th, minP}); got != nil {
		t.Fatalf("setPriority success expected nil, got %T", got)
	}
}

// ---- threadSleep ----
func TestThreadSleep_TypeErrorAndSuccess(t *testing.T) {
	if _, ok := threadSleep([]interface{}{"x"}).(*GErrBlk); !ok {
		t.Fatalf("type error expected")
	}
	if threadSleep([]interface{}{int64(1)}) != nil {
		t.Fatalf("sleep should return nil")
	}
}

// ---- cloneNotSupportedException ----
func TestCloneNotSupportedException_ReturnsGErr(t *testing.T) {
	res := cloneNotSupportedException(nil)
	g, ok := res.(*GErrBlk)
	if !ok || g.ExceptionType != excNames.CloneNotSupportedException {
		t.Fatalf("expected CloneNotSupportedException, got %T %+v", res, res)
	}
}

// ---- threadNumbering ----
func TestThreadNumbering_InitAndNext(t *testing.T) {
	// Reset global counter for deterministic test (safe because package under test)
	thread.ThreadNumber = 0
	_ = threadNumbering(nil)
	n1 := threadNumberingNext(nil).(int64)
	n2 := threadNumberingNext(nil).(int64)
	if !(n1 == 1 && n2 == 2) {
		t.Fatalf("thread numbering wrong: %d %d", n1, n2)
	}
}

// Safety: ensure classloader symbol referenced so linter doesn’t drop imports in trimmed builds
var _ = classloader.MTentry{}
