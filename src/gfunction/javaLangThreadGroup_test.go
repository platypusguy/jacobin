/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"jacobin/src/exceptions"
	"jacobin/src/stringPool"
	"sync"
	"testing"

	"jacobin/src/classloader"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

var initOnce sync.Once

// Ensure base init used by other gfunction tests
func ensureTGInit() {
	ensureInit()
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

// instantiateForThreadTest is a minimal stand-in for jvm.Instantiate to avoid circular imports.
// It satisfies the globals.FuncInstantiateClass signature used by gfunctions during tests.
func instantiateForThreadTest(name string, _ *list.List) (any, error) {
	o := object.MakeEmptyObject()
	o.KlassName = stringPool.GetStringIndex(&name)
	return o, nil
}

func TestInitializeGlobalThreadGroups_TestMode(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	// Force map to be nil to test map initialization
	gr.ThreadGroups = nil
	gr.JacobinName = "test"
	InitializeGlobalThreadGroups()

	sys := gr.ThreadGroups["system"].(*object.Object)
	mainTG := gr.ThreadGroups["main"].(*object.Object)
	if sys == nil || mainTG == nil {
		t.Errorf("expected system and main groups to be created; got %#v %#v", sys, mainTG)
	}
	// main's parent should be system
	p := mainTG.FieldTable["parent"].Fvalue.(*object.Object)
	if p != sys {
		t.Errorf("expected main.parent == system")
	}
	// system.subgroups should contain main
	lst := sys.FieldTable["subgroups"].Fvalue.(*list.List)
	found := false
	for e := lst.Front(); e != nil; e = e.Next() {
		if e.Value == mainTG {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected system.subgroups to include main")
	}
}

func TestInitializeGlobalThreadGroups_NonTestMode(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	gr.ThreadGroups = make(map[string]interface{})
	gr.JacobinName = "run"

	// Prepare method area with a minimal ThreadGroup Klass
	classloader.InitMethodArea()
	kg := &classloader.Klass{Status: 'N', Loader: "bootstrap", Data: &classloader.ClData{}}
	classloader.MethAreaInsert("java/lang/ThreadGroup", kg)

	// Use the test instantiator from ensureInit
	InitializeGlobalThreadGroups()

	if gr.ThreadGroups["system"] == nil || gr.ThreadGroups["main"] == nil {
		t.Errorf("expected both system and main keys present in globals.ThreadGroups")
	}
}

func TestThreadGroupClinit_Returns(t *testing.T) {
	if threadGroupClinit(nil) != nil {
		t.Errorf("clinit should just return nil")
	}
}

func TestThreadGroupInitWithParentNameMaxpriorityDaemon_ParamCount(t *testing.T) {
	if _, ok := ThreadGroupInitWithParentNameMaxpriorityDaemon([]any{1}).(*GErrBlk); !ok {
		t.Errorf("expected error for wrong param count")
	}
}

func TestThreadGroupInitWithParentNameMaxpriorityDaemon_TypeErrorsAndSuccess(t *testing.T) {
	ensureTGInit()
	gr := globals.GetGlobalRef()
	gr.ThreadGroups = make(map[string]interface{})

	// 1st param not object
	{
		res := ThreadGroupInitWithParentNameMaxpriorityDaemon([]any{123, object.Null, object.StringObjectFromGoString("tg"), int64(0), types.JavaBoolUninitialized})
		if _, ok := res.(*GErrBlk); !ok {
			t.Errorf("expected type error for 1st param")
		}
	}
	// 3rd param not String object
	{
		obj := threadGroupFake("x")
		res := ThreadGroupInitWithParentNameMaxpriorityDaemon([]any{obj, object.Null, 99, int64(0), types.JavaBoolUninitialized})
		if _, ok := res.(*GErrBlk); !ok {
			t.Errorf("expected type error for 3rd param")
		}
	}
	// maxPriority out of range
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		name := object.StringObjectFromGoString("bad")
		res := ThreadGroupInitWithParentNameMaxpriorityDaemon([]any{obj, object.Null, name, int64(9999), types.JavaBoolUninitialized})
		if _, ok := res.(*GErrBlk); !ok {
			t.Errorf("expected error for maxPriority out of range")
		}
	}
	// Success: daemon uninitialized, parent null
	{
		// tgName := "java/lang/ThreadGroup"
		obj := threadGroupFake("grpA")
		name := object.StringObjectFromGoString("grpA")
		res := ThreadGroupInitWithParentNameMaxpriorityDaemon([]any{obj, object.Null, name, int64(0), types.JavaBoolUninitialized})
		tg := res.(*object.Object)
		if tg.FieldTable["parent"].Fvalue != object.Null {
			t.Errorf("expected parent to remain null")
		}
		if tg.FieldTable["priority"].Ftype != types.Int {
			t.Errorf("priority not initialized")
		}
		if tg.FieldTable["subgroups"].Fvalue.(*list.List) == nil {
			t.Errorf("subgroups not initialized")
		}
		if gr.ThreadGroups["grpA"] == nil {
			t.Errorf("expected group to be registered globally")
		}
	}
	// Success: parent set, maxPriority in range, daemon true
	{
		parent := threadGroupFake("parent")
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		name := object.StringObjectFromGoString("grpB")
		res := ThreadGroupInitWithParentNameMaxpriorityDaemon([]any{obj, parent, name, int64(5), types.JavaBoolTrue})
		tg := res.(*object.Object)
		if tg.FieldTable["parent"].Fvalue.(*object.Object) != parent {
			t.Errorf("expected parent to be set")
		}
		if tg.FieldTable["maxpriority"].Fvalue.(int64) != 5 {
			t.Errorf("maxpriority not set")
		}
		if tg.FieldTable["daemon"].Fvalue != types.JavaBoolTrue {
			t.Errorf("daemon not set true")
		}
	}
}

func TestThreadGroupInitWithName_AllPaths(t *testing.T) {
	// wrong count
	if _, ok := threadGroupInitWithName([]any{1}).(*GErrBlk); !ok {
		t.Errorf("expected error for wrong count")
	}
	// first not object
	if _, ok := threadGroupInitWithName([]any{123, object.StringObjectFromGoString("n")}).(*GErrBlk); !ok {
		t.Errorf("expected error for first param type")
	}
	// second not string object
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		if _, ok := threadGroupInitWithName([]any{obj, 7}).(*GErrBlk); !ok {
			t.Errorf("expected error for second param type")
		}
	}
	// success
	{
		gr := globals.GetGlobalRef()
		gr.ThreadGroups = make(map[string]interface{})
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		name := object.StringObjectFromGoString("solo")
		res := threadGroupInitWithName([]any{obj, name})
		if _, ok := res.(*object.Object); !ok {
			t.Errorf("expected success object")
		}
		if gr.ThreadGroups["solo"] == nil {
			t.Errorf("expected group registered under name solo")
		}
	}
}

func TestThreadGroupInitWithParentAndName_AllPaths(t *testing.T) {
	// wrong count
	if _, ok := threadGroupInitWithParentAndName([]any{1}).(*GErrBlk); !ok {
		t.Errorf("expected error for wrong count")
	}
	// first not object
	if _, ok := threadGroupInitWithParentAndName([]any{123, object.MakeEmptyObject(), object.StringObjectFromGoString("x")}).(*GErrBlk); !ok {
		t.Errorf("expected error for 1st param type")
	}
	// second not object
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		if _, ok := threadGroupInitWithParentAndName([]any{obj, 99, object.StringObjectFromGoString("x")}).(*GErrBlk); !ok {
			t.Errorf("expected error for 2nd param type")
		}
	}
	// third not string object
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		parent := threadGroupFake("parent")
		if _, ok := threadGroupInitWithParentAndName([]any{obj, parent, 77}).(*GErrBlk); !ok {
			t.Errorf("expected error for 3rd param type")
		}
	}
	// name null
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		parent := threadGroupFake("parent")
		if _, ok := threadGroupInitWithParentAndName([]any{obj, parent, object.Null}).(*GErrBlk); !ok {
			t.Errorf("expected NPE for null name")
		}
	}
	// name not a String object
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		parent := threadGroupFake("parent")
		thName := "java/lang/Thread"
		notString := object.MakeEmptyObjectWithClassName(&thName)
		if _, ok := threadGroupInitWithParentAndName([]any{obj, parent, notString}).(*GErrBlk); !ok {
			t.Errorf("expected IAE for non-String name")
		}
	}
	// success adds to parent's subgroups
	{
		gr := globals.GetGlobalRef()
		gr.ThreadGroups = make(map[string]interface{})
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		parent := threadGroupFake("parent2")
		name := object.StringObjectFromGoString("child")
		res := threadGroupInitWithParentAndName([]any{obj, parent, name})
		tg := res.(*object.Object)
		lst := parent.FieldTable["subgroups"].Fvalue.(*list.List)
		if lst.Len() != 1 || lst.Front().Value != tg {
			t.Errorf("expected new group added to parent's subgroups")
		}
	}
}

func TestThreadGroupGetName_AllPaths(t *testing.T) {
	// wrong count
	if _, ok := threadGroupGetName([]any{}).(*GErrBlk); !ok {
		t.Errorf("expected error for wrong param count")
	}
	// wrong type
	if _, ok := threadGroupGetName([]any{123}).(*GErrBlk); !ok {
		t.Errorf("expected error for wrong type")
	}
	// name as Java String object
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		obj.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("alpha")}
		res := threadGroupGetName([]any{obj}).(*object.Object)
		if object.GoStringFromStringObject(res) != "alpha" {
			t.Errorf("expected alpha, got %s", object.GoStringFromStringObject(res))
		}
	}
	// name as Go string
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		obj.FieldTable["name"] = object.Field{Ftype: types.GolangString, Fvalue: "beta"}
		res := threadGroupGetName([]any{obj}).(*object.Object)
		if object.GoStringFromStringObject(res) != "beta" {
			t.Errorf("expected beta")
		}
	}
	// name as JavaByte array
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		jb := object.JavaByteArrayFromGoString("gamma")
		obj.FieldTable["name"] = object.Field{Ftype: types.ByteArray, Fvalue: jb}
		res := threadGroupGetName([]any{obj}).(*object.Object)
		if object.GoStringFromStringObject(res) != "gamma" {
			t.Errorf("expected gamma")
		}
	}
	// invalid type in name field
	{
		tgName := "java/lang/ThreadGroup"
		obj := object.MakeEmptyObjectWithClassName(&tgName)
		obj.FieldTable["name"] = object.Field{Ftype: types.Int, Fvalue: int64(7)}
		if _, ok := threadGroupGetName([]any{obj}).(*GErrBlk); !ok {
			t.Errorf("expected error for invalid name field type")
		}
	}
}
