/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/thread"
	"jacobin/src/types"
	"testing"
)

// Ensure the ThreadGroups map exists for tests

func resetThreadGroups() {
	globals.GetGlobalRef().ThreadGroups = make(map[string]interface{})
}

// ---- threadGroupClinit() ----
func TestThreadGroupClinit_ReturnsNil(t *testing.T) {
	if got := threadGroupClinit(nil); got != nil {
		t.Fatalf("threadGroupClinit() = %v; want nil", got)
	}
}

// ---- threadGroupCreateWithName() ----
func TestThreadGroupCreateWithName_ArgCountError(t *testing.T) {
	resetThreadGroups()
	res := threadGroupCreateWithName([]interface{}{})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithName_WrongType(t *testing.T) {
	resetThreadGroups()
	res := threadGroupCreateWithName([]interface{}{123})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithName_HappyPath(t *testing.T) {
	resetThreadGroups()
	nameObj := object.StringObjectFromGoString("system")
	res := threadGroupCreateWithName([]interface{}{nameObj})
	grp, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}

	// parent should be null
	if p := grp.FieldTable["parent"]; p.Ftype != types.Ref || !object.IsNull(p.Fvalue) {
		t.Errorf("parent field not null ref: %+v", p)
	}
	// name stored as string object
	if n := grp.FieldTable["name"]; n.Ftype != types.Ref || !object.IsStringObject(n.Fvalue) || n.Fvalue.(*object.Object) != nameObj {
		t.Errorf("name field wrong: %+v", n)
	}
	// daemon default false
	if d := grp.FieldTable["daemon"]; d.Ftype != types.Int || d.Fvalue.(int64) != types.JavaBoolFalse {
		t.Errorf("daemon field wrong: %+v", d)
	}
	// threadgroup ref placeholder
	if tg := grp.FieldTable["threadgroup"]; tg.Ftype != types.Ref || tg.Fvalue != nil {
		t.Errorf("threadgroup field wrong: %+v", tg)
	}
	// priorities
	if pr := grp.FieldTable["priority"]; pr.Ftype != types.Int || pr.Fvalue.(int64) != int64(thread.NORM_PRIORITY) {
		t.Errorf("priority wrong: %+v", pr)
	}
	if mx := grp.FieldTable["maxpriority"]; mx.Ftype != types.Int || mx.Fvalue.(int64) != int64(thread.MAX_PRIORITY) {
		t.Errorf("maxpriority wrong: %+v", mx)
	}
	// subgroups list
	if sg := grp.FieldTable["subgroups"]; sg.Ftype != types.LinkedList {
		t.Errorf("subgroups Ftype wrong: %+v", sg)
	} else {
		if _, ok := sg.Fvalue.(*list.List); !ok {
			t.Errorf("subgroups Fvalue not *list.List: %T", sg.Fvalue)
		}
	}
	// added to globals map
	if got, ok := globals.GetGlobalRef().ThreadGroups[object.GoStringFromStringObject(nameObj)]; !ok || got != grp {
		t.Errorf("group not added to globals.ThreadGroups; ok=%v got=%T", ok, got)
	}
}

// ---- threadGroupCreateWithParentAndName() ----
func TestThreadGroupCreateWithParentAndName_ArgCountError(t *testing.T) {
	resetThreadGroups()
	res := threadGroupCreateWithParentAndName([]interface{}{object.StringObjectFromGoString("nameOnly")})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithParentAndName_WrongFirstType(t *testing.T) {
	resetThreadGroups()
	res := threadGroupCreateWithParentAndName([]interface{}{123, object.StringObjectFromGoString("x")})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithParentAndName_WrongSecondType(t *testing.T) {
	resetThreadGroups()
	parent := threadGroupCreateWithName([]interface{}{object.StringObjectFromGoString("p")}).(*object.Object)
	res := threadGroupCreateWithParentAndName([]interface{}{parent, 456})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithParentAndName_NullName(t *testing.T) {
	resetThreadGroups()
	parent := threadGroupCreateWithName([]interface{}{object.StringObjectFromGoString("p")}).(*object.Object)
	res := threadGroupCreateWithParentAndName([]interface{}{parent, object.Null})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithParentAndName_NotStringObject(t *testing.T) {
	resetThreadGroups()
	parent := threadGroupCreateWithName([]interface{}{object.StringObjectFromGoString("p")}).(*object.Object)
	notStr := object.MakeEmptyObject()
	res := threadGroupCreateWithParentAndName([]interface{}{parent, notStr})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupCreateWithParentAndName_HappyPath(t *testing.T) {
	resetThreadGroups()
	pName := object.StringObjectFromGoString("system")
	parent := threadGroupCreateWithName([]interface{}{pName}).(*object.Object)

	cName := object.StringObjectFromGoString("main")
	res := threadGroupCreateWithParentAndName([]interface{}{parent, cName})
	child, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", res)
	}
	// child parent field set
	if pf := child.FieldTable["parent"]; pf.Ftype != types.Ref || pf.Fvalue.(*object.Object) != parent {
		t.Fatalf("child.parent wrong: %+v", pf)
	}
	// parent subgroups includes child
	sg := parent.FieldTable["subgroups"].Fvalue.(*list.List)
	if sg.Len() != 1 {
		t.Fatalf("expected 1 subgroup, got %d", sg.Len())
	}
	if sg.Front().Value != child {
		t.Errorf("parent.subgroups first != child")
	}
	// child should also have been added to globals by its name (via createWithName)
	if got, ok := globals.GetGlobalRef().ThreadGroups[object.GoStringFromStringObject(cName)]; !ok || got != child {
		t.Errorf("child not present in globals map; ok=%v got=%T", ok, got)
	}
}

// ---- threadGroupGetName() ----
func TestThreadGroupGetName_ArgCountError(t *testing.T) {
	res := threadGroupGetName([]interface{}{})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupGetName_WrongType(t *testing.T) {
	res := threadGroupGetName([]interface{}{123})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

func TestThreadGroupGetName_ReturnsStringObject(t *testing.T) {
	resetThreadGroups()
	nameObj := object.StringObjectFromGoString("alpha")
	grp := threadGroupCreateWithName([]interface{}{nameObj}).(*object.Object)
	res := threadGroupGetName([]interface{}{grp})
	if !object.IsStringObject(res) {
		t.Fatalf("expected String object, got %T", res)
	}
	if object.GoStringFromStringObject(res.(*object.Object)) != "alpha" {
		t.Errorf("unexpected name value: %s", object.GoStringFromStringObject(res.(*object.Object)))
	}
}

func TestThreadGroupGetName_LegacyGoStringFallback(t *testing.T) {
	// Create a minimal ThreadGroup object with name stored as Go string
	clName := "java/lang/ThreadGroup"
	grp := object.MakeEmptyObjectWithClassName(&clName)
	grp.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: "legacy"}
	res := threadGroupGetName([]interface{}{grp})
	if !object.IsStringObject(res) {
		t.Fatalf("expected String object, got %T", res)
	}
	if s := object.GoStringFromStringObject(res.(*object.Object)); s != "legacy" {
		t.Errorf("expected 'legacy', got %q", s)
	}
}

func TestThreadGroupGetName_UnexpectedTypeError(t *testing.T) {
	clName := "java/lang/ThreadGroup"
	grp := object.MakeEmptyObjectWithClassName(&clName)
	grp.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: 123}
	res := threadGroupGetName([]interface{}{grp})
	gerr, ok := res.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", res)
	}
	if gerr.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException, got %v", gerr.ExceptionType)
	}
}

// ---- initializeGlobalThreadGroups() ----
func TestInitializeGlobalThreadGroups_SetsUpSystemAndMain(t *testing.T) {
	resetThreadGroups()
	initializeGlobalThreadGroups()
	sys, ok1 := globals.GetGlobalRef().ThreadGroups["system"].(*object.Object)
	main, ok2 := globals.GetGlobalRef().ThreadGroups["main"].(*object.Object)
	if !ok1 || !ok2 || sys == nil || main == nil {
		t.Fatalf("expected system and main groups present; got: system=%T main=%T", globals.GetGlobalRef().ThreadGroups["system"], globals.GetGlobalRef().ThreadGroups["main"])
	}
	// main's parent is system
	pf := main.FieldTable["parent"]
	if pf.Ftype != types.Ref || pf.Fvalue.(*object.Object) != sys {
		t.Errorf("main.parent not system: %+v", pf)
	}
	// system has main in its subgroups
	sg := sys.FieldTable["subgroups"].Fvalue.(*list.List)
	if sg.Len() != 1 || sg.Front().Value != main {
		t.Errorf("system.subgroups does not contain main")
	}
}
