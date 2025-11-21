/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */
package gfunction

import (
	"container/list"
	"io"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

// Many tests rely on ensureTGInit() defined in javaLangThreadGroup_test.go.

func TestThreadClinitConstants(t *testing.T) {
	ensureTGInit()
	// Load_Lang_Thread invoked in ensureInit sets statics via threadClinit
	// Verify constants exist
	min := staticsGet("java/lang/Thread", "MIN_PRIORITY")
	norm := staticsGet("java/lang/Thread", "NORM_PRIORITY")
	max := staticsGet("java/lang/Thread", "MAX_PRIORITY")
	if min == nil || norm == nil || max == nil {
		t.Fatalf("expected Thread priority statics to be set")
	}
}

// helper to access statics safely in tests (avoids import cycle)
func staticsGet(cls, name string) any { return statics.GetStaticValue(cls, name) }

func TestThreadCreateNoarg_Defaults(t *testing.T) {
	ensureTGInit()
	obj := ThreadCreateNoarg(nil).(*object.Object)
	if obj.FieldTable["ID"].Fvalue.(int64) == 0 {
		t.Errorf("expected non-zero thread ID")
	}
	// Default name Thread-N
	nameObj := obj.FieldTable["name"].Fvalue.(*object.Object)
	name := object.GoStringFromStringObject(nameObj)
	if !strings.HasPrefix(name, "Thread-") {
		t.Errorf("expected default name 'Thread-N', got %s", name)
	}
	// Default state NEW
	st := obj.FieldTable["state"].Fvalue.(*object.Object)
	// sanity: ensure an enum-like object returned
	if st == nil {
		t.Errorf("expected non-nil state")
	}
	// Daemon false
	if obj.FieldTable["daemon"].Fvalue.(int64) != types.JavaBoolFalse {
		t.Errorf("expected daemon false")
	}
	// Interrupted false
	if obj.FieldTable["interrupted"].Fvalue.(int64) != types.JavaBoolFalse {
		t.Errorf("expected interrupted false")
	}
	// Thread group main exists
	tg := obj.FieldTable["threadgroup"].Fvalue.(*object.Object)
	if tg == nil {
		t.Errorf("expected threadgroup set")
	}
	// Priority is NORM_PRIORITY
	norm := staticsGet("java/lang/Thread", "NORM_PRIORITY").(int64)
	if obj.FieldTable["priority"].Fvalue.(int64) != norm {
		t.Errorf("expected priority %d", norm)
	}
	// Frame stack present
	if obj.FieldTable["framestack"].Ftype != types.LinkedList {
		t.Errorf("expected framestack LinkedList type")
	}
	// Task is nil
	if obj.FieldTable["task"].Fvalue != nil {
		t.Errorf("expected nil task")
	}
}

func TestThreadInitWithName_ErrWrongArity(t *testing.T) {
	ensureTGInit()
	ret := threadInitWithName([]any{ThreadCreateNoarg(nil)})
	g := ret.(*GErrBlk)
	if g.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IAE, got %d", g.ExceptionType)
	}
}

func TestThreadInitWithName_ErrWrongTypes(t *testing.T) {
	ensureTGInit()
	ret := threadInitWithName([]any{123, object.StringObjectFromGoString("n")})
	if ret.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IAE for non-thread first arg")
	}
	th := ThreadCreateNoarg(nil).(*object.Object)
	ret2 := threadInitWithName([]any{th, 5})
	if ret2.(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IAE for non-string name")
	}
}

func TestThreadInitWithName_Success(t *testing.T) {
	ensureTGInit()
	th := ThreadCreateNoarg(nil).(*object.Object)
	nm := object.StringObjectFromGoString("Alpha")
	_ = threadInitWithName([]any{th, nm})
	got := th.FieldTable["name"].Fvalue.(*object.Object)
	if object.GoStringFromStringObject(got) != object.GoStringFromStringObject(nm) {
		t.Errorf("name not set correctly")
	}
}

func TestThreadInitWithRunnableAndName_Paths(t *testing.T) {
	ensureTGInit()
	// wrong arity
	if threadInitWithRunnableAndName([]any{ThreadCreateNoarg(nil)}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for arity")
	}
	// wrong types each position
	th := ThreadCreateNoarg(nil).(*object.Object)
	nm := object.StringObjectFromGoString("B")
	runnable := makeRunnableDescriptor("C", "run", "()V")

	if threadInitWithRunnableAndName([]any{123, runnable, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-thread")
	}
	if threadInitWithRunnableAndName([]any{th, 456, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-runnable")
	}
	if threadInitWithRunnableAndName([]any{th, runnable, 789}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-string name")
	}
	// success
	threadInitWithRunnableAndName([]any{th, runnable, nm})
	if th.FieldTable["task"].Fvalue.(*object.Object) != runnable {
		t.Errorf("runnable not set")
	}
	if th.FieldTable["name"].Fvalue.(*object.Object) != nm {
		t.Errorf("name not set")
	}
}

func TestThreadInitWithThreadGroupAndName_Paths(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	mainTG := gr.ThreadGroups["main"].(*object.Object)
	th := ThreadCreateNoarg(nil).(*object.Object)
	nm := object.StringObjectFromGoString("D")

	if threadInitWithThreadGroupAndName([]any{th, mainTG}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE arity")
	}
	if threadInitWithThreadGroupAndName([]any{123, mainTG, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-thread")
	}
	if threadInitWithThreadGroupAndName([]any{th, 456, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-threadgroup")
	}
	if threadInitWithThreadGroupAndName([]any{th, mainTG, 789}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-name")
	}
	threadInitWithThreadGroupAndName([]any{th, mainTG, nm})
	if th.FieldTable["threadgroup"].Fvalue.(*object.Object) != mainTG {
		t.Errorf("threadgroup not set")
	}
}

func TestThreadInitWithThreadGroupRunnableAndName_Paths(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	mainTG := gr.ThreadGroups["main"].(*object.Object)
	th := ThreadCreateNoarg(nil).(*object.Object)
	nm := object.StringObjectFromGoString("E")
	runnable := makeRunnableDescriptor("C2", "run", "()V")

	if threadInitWithThreadGroupRunnableAndName([]any{th, mainTG, runnable}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE arity")
	}
	if threadInitWithThreadGroupRunnableAndName([]any{123, mainTG, runnable, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-thread")
	}
	if threadInitWithThreadGroupRunnableAndName([]any{th, 456, runnable, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-threadgroup")
	}
	if threadInitWithThreadGroupRunnableAndName([]any{th, mainTG, 789, nm}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-runnable")
	}
	if threadInitWithThreadGroupRunnableAndName([]any{th, mainTG, runnable, 999}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for non-name")
	}
	threadInitWithThreadGroupRunnableAndName([]any{th, mainTG, runnable, nm})
	if th.FieldTable["task"].Fvalue.(*object.Object) != runnable {
		t.Errorf("task not set")
	}
}

func TestThreadInitFromPackageConstructor_Paths(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	mainTG := gr.ThreadGroups["main"].(*object.Object)
	th := ThreadCreateNoarg(nil).(*object.Object)
	nm := object.StringObjectFromGoString("P")
	runnable := makeRunnableDescriptor("RC", "run", "()V")

	// arity error
	if threadInitFromPackageConstructor([]any{th}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE arity")
	}

	// wrong types by positions
	// 0 can be nil or thread; use non-thread
	if threadInitFromPackageConstructor([]any{123, mainTG, nm, int64(5), runnable, int64(0), nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for first thread param")
	}
	// 1 must be threadgroup or nil
	if threadInitFromPackageConstructor([]any{th, 456, nm, int64(5), runnable, int64(0), nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for threadgroup param")
	}
	// 2 must be name string
	if threadInitFromPackageConstructor([]any{th, mainTG, 789, int64(5), runnable, int64(0), nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for name param")
	}
	// 3 must be int
	if threadInitFromPackageConstructor([]any{th, mainTG, nm, "x", runnable, int64(0), nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for priority param")
	}
	// 4 must be runnable or nil
	if threadInitFromPackageConstructor([]any{th, mainTG, nm, int64(5), 111, int64(0), nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for runnable param")
	}
	// 5 must be int64
	if threadInitFromPackageConstructor([]any{th, mainTG, nm, int64(5), runnable, 3, nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for long param")
	}
	// 6 may be object or nil, BUT code checks index 5 when 6 is non-nil -> triggers error
	if threadInitFromPackageConstructor([]any{th, mainTG, nm, int64(5), runnable, int64(0), object.MakeEmptyObject()}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE for access control context param due to current implementation")
	}

	// success path (with 6 = nil)
	threadInitFromPackageConstructor([]any{th, mainTG, nm, int64(5), runnable, int64(0), nil})
	if th.FieldTable["name"].Fvalue.(*object.Object) != nm || th.FieldTable["task"].Fvalue.(*object.Object) != runnable {
		t.Errorf("expected runnable+name wired through")
	}
	if th.FieldTable["threadgroup"].Fvalue.(*object.Object) != mainTG {
		t.Errorf("threadgroup not set on thread")
	}
}

func TestThreadActiveCount(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	gr.Threads = map[int]interface{}{}
	gr.Threads[1] = ThreadCreateNoarg(nil)
	gr.Threads[2] = ThreadCreateNoarg(nil)
	if threadActiveCount(nil).(int64) != 2 {
		t.Errorf("expected 2 active threads")
	}
}

func TestThreadDumpStack_Paths(t *testing.T) {
	ensureTGInit()
	// arity
	if threadDumpStack(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE arity")
	}
	if threadDumpStack([]any{123}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("expected IAE type")
	}
	// capture stderr for both StrictJDK branches
	normal := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(1)
	f.ClName = "A"
	f.MethName = "m"
	f.PC = 3
	f.Thread = 5
	fs.PushFront(f)
	// register thread with name
	th := ThreadCreateNoarg(nil).(*object.Object)
	th.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("T")}
	globals.GetGlobalRef().Threads[5] = th

	globals.GetGlobalRef().StrictJDK = false
	_ = threadDumpStack([]any{fs})
	_ = w.Close()
	bytes1, _ := io.ReadAll(r)

	r2, w2, _ := os.Pipe()
	os.Stderr = w2
	globals.GetGlobalRef().StrictJDK = true
	_ = threadDumpStack([]any{fs})
	_ = w2.Close()
	bytes2, _ := io.ReadAll(r2)
	os.Stderr = normal

	if !strings.Contains(string(bytes1), "Stack trace (thread T)") {
		t.Errorf("expected custom header")
	}
	if !strings.Contains(string(bytes2), "java.lang.Exception: Stack trace") {
		t.Errorf("expected JDK header")
	}
}

func TestThreadGetId_Paths(t *testing.T) {
	ensureTGInit()
	if threadGetId(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity")
	}
	if threadGetId([]any{123}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("type")
	}
	th := ThreadCreateNoarg(nil).(*object.Object)
	// success
	id := threadGetId([]any{th}).(int64)
	if id != th.FieldTable["ID"].Fvalue.(int64) {
		t.Errorf("id mismatch")
	}
	// missing ID field
	th2 := object.MakeEmptyObject()
	th2.FieldTable = map[string]object.Field{}
	if threadGetId([]any{th2}).(*GErrBlk).ExceptionType != excNames.InternalException {
		t.Fatal("missing ID should be internal error")
	}
	// wrong ID type
	th3 := object.MakeEmptyObject()
	th3.FieldTable = map[string]object.Field{"ID": {Ftype: types.Int, Fvalue: "x"}}
	if threadGetId([]any{th3}).(*GErrBlk).ExceptionType != excNames.InternalException {
		t.Fatal("wrong ID type")
	}
}

func TestThreadGetNamePriorityStateGroupInterrupted(t *testing.T) {
	ensureTGInit()
	th := ThreadCreateNoarg(nil).(*object.Object)
	// getName
	if threadGetName(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity getName")
	}
	if threadGetName([]any{123}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("type getName")
	}
	nm := threadGetName([]any{th}).(*object.Object)
	if object.GoStringFromStringObject(nm) == "" {
		t.Errorf("expected name non-empty")
	}

	// getPriority
	if threadGetPriority(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity getPriority")
	}
	pr := threadGetPriority([]any{th}).(int64)
	if pr == 0 {
		t.Errorf("expected non-zero priority")
	}

	// getState
	if threadGetState(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity getState")
	}
	if threadGetState([]any{123}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("type getState")
	}
	thState, ok := th.FieldTable["state"].Fvalue.(*object.Object)
	if !ok {
		t.Errorf("state missing or is not an object")
	}
	stateValue := thState.FieldTable["value"].Fvalue.(int)
	if stateValue != 0 {
		t.Errorf("expected state value 0")
	}

	// getThreadGroup
	if threadGetThreadGroup(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity getTG")
	}
	tg := threadGetThreadGroup([]any{th}).(*object.Object)
	if tg == nil {
		t.Errorf("expected TG object")
	}

	// isInterrupted
	if threadIsInterrupted(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity isInt")
	}
	if threadIsInterrupted([]any{123}).(*GErrBlk).ExceptionType != excNames.InternalException {
		t.Fatal("type isInt")
	}
	if threadIsInterrupted([]any{th}).(int64) != types.JavaBoolFalse {
		t.Errorf("expected false")
	}
}

func TestThreadGetStackTrace_Paths(t *testing.T) {
	ensureTGInit()
	// arity/type
	if threadGetStackTrace([]any{list.New()}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity")
	}
	if threadGetStackTrace([]any{123, nil}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("type")
	}

	fs := frames.CreateFrameStack()
	// success relies on FillInStackTrace writing a stack trace
	traceArrObj := threadGetStackTrace([]any{fs, object.MakeEmptyObject()}).(*object.Object)
	if traceArrObj == nil {
		t.Errorf("expected non-nil stack trace array object")
	}
}

func TestThreadSetName_Paths(t *testing.T) {
	ensureTGInit()
	if threadSetName(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity")
	}
	if threadSetName([]any{123, object.StringObjectFromGoString("x")}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("type")
	}
	th := ThreadCreateNoarg(nil).(*object.Object)
	if threadSetName([]any{th, object.Null}).(*GErrBlk).ExceptionType != excNames.NullPointerException {
		t.Fatal("npe expected")
	}
	nm := object.StringObjectFromGoString("Zed")
	if threadSetName([]any{th, nm}) != nil {
		t.Fatal("expected nil return on success")
	}
	got := threadGetName([]any{th}).(*object.Object)
	if object.GoStringFromStringObject(got) != "Zed" {
		t.Errorf("name not updated")
	}
}

func TestThreadSetPriority_Paths(t *testing.T) {
	ensureTGInit()
	if threadSetPriority(nil).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("arity")
	}
	if threadSetPriority([]any{123, int64(5)}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("type")
	}
	th := ThreadCreateNoarg(nil).(*object.Object)
	if threadSetPriority([]any{th, "x"}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("priority type")
	}
	min := staticsGet("java/lang/Thread", "MIN_PRIORITY").(int64)
	max := staticsGet("java/lang/Thread", "MAX_PRIORITY").(int64)
	if threadSetPriority([]any{th, min - 1}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("below min")
	}
	if threadSetPriority([]any{th, max + 1}).(*GErrBlk).ExceptionType != excNames.IllegalArgumentException {
		t.Fatal("above max")
	}
	if threadSetPriority([]any{th, (min + max) / 2}) != nil {
		t.Fatal("expected success")
	}
}

func TestThreadSleep_Paths(t *testing.T) {
	ensureTGInit()
	if threadSleep([]any{"x"}).(*GErrBlk).ExceptionType != excNames.IOException {
		t.Fatal("type err")
	}
	if threadSleep([]any{int64(0)}) != nil {
		t.Fatal("expected nil on success")
	}
}

func TestCloneNotSupportedException(t *testing.T) {
	ensureTGInit()
	ret := cloneNotSupportedException(nil).(*GErrBlk)
	if ret.ExceptionType != excNames.CloneNotSupportedException {
		t.Fatalf("wrong exception type")
	}
}

func TestThreadNumbering_InitAndNext(t *testing.T) {
	ensureTGInit()
	// init
	_ = threadNumbering(nil)
	// next increments
	a := threadNumberingNext(nil).(int64)
	b := threadNumberingNext(nil).(int64)
	if b != a+1 {
		t.Errorf("expected monotonically increasing thread numbers")
	}
}

// makeRunnableDescriptor builds the Runnable-like descriptor object expected by threadRun.
func makeRunnableDescriptor(className, methodName, sig string) *object.Object {
	o := object.MakeEmptyObject()
	o.FieldTable = map[string]object.Field{}
	o.FieldTable["clName"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(className)}
	o.FieldTable["methName"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(methodName)}
	o.FieldTable["signature"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(sig)}
	return o
}
