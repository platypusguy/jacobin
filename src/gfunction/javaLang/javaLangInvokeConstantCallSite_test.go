/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"jacobin/src/types"
	"testing"
)

func TestLoad_Lang_Invoke_ConstantCallSite_Registration(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Invoke_ConstantCallSite()

	expectedMethods := []string{
		"java/lang/invoke/ConstantCallSite.<init>(Ljava/lang/invoke/MethodHandle;)V",
		"java/lang/invoke/ConstantCallSite.<init>(Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodHandle;)V",
		"java/lang/invoke/ConstantCallSite.getTarget()Ljava/lang/invoke/MethodHandle;",
		"java/lang/invoke/ConstantCallSite.setTarget(Ljava/lang/invoke/MethodHandle;)V",
		"java/lang/invoke/ConstantCallSite.dynamicInvoker()Ljava/lang/invoke/MethodHandle;",
		"java/lang/invoke/ConstantCallSite.type()Ljava/lang/invoke/MethodType;",
		"java/lang/invoke/CallSite.type()Ljava/lang/invoke/MethodType;",
		"java/lang/invoke/CallSite.getTarget()Ljava/lang/invoke/MethodHandle;",
		"java/lang/invoke/CallSite.setTarget(Ljava/lang/invoke/MethodHandle;)V",
		"java/lang/invoke/CallSite.dynamicInvoker()Ljava/lang/invoke/MethodHandle;",
	}

	for _, sig := range expectedMethods {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected method signature to be registered: %s", sig)
		}
	}
}

func createTestDummyMethodHandle(mtObj *object.Object) *object.Object {
	mho := createRawMethodHandleObject()
	mho.FieldTable["Kind"] = object.Field{Ftype: types.Int, Fvalue: int64(REF_invokeStatic)}
	mho.FieldTable["RefClass"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("test/Caller")}
	mho.FieldTable["RefName"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("targetMethod")}
	mho.FieldTable["RefDescriptor"] = object.Field{Ftype: types.Ref, Fvalue: mtObj}
	mho.FieldTable["type"] = object.Field{Ftype: types.Ref, Fvalue: mtObj}
	return mho
}

func createTestDummyMethodType() *object.Object {
	mtClassName := "java/lang/invoke/MethodType"
	mtObj := object.MakeEmptyObjectWithClassName(&mtClassName)
	mtObj.FieldTable["rtype"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("Ljava/lang/String;")}
	return mtObj
}

func TestConstantCallSite_Init_Success(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	mtObj := createTestDummyMethodType()
	mhObj := createTestDummyMethodHandle(mtObj)
	callSiteObj := CreateRawConstantCallSiteObject()

	res := constantCallSiteInit([]any{callSiteObj, mhObj})
	if res != nil {
		t.Fatalf("Expected nil return on successful init, got %v", res)
	}

	targetField, ok := callSiteObj.FieldTable["target"]
	if !ok || targetField.Fvalue != mhObj {
		t.Errorf("Expected target field to be set to mhObj")
	}

	typeField, ok := callSiteObj.FieldTable["type"]
	if !ok || typeField.Fvalue != mtObj {
		t.Errorf("Expected type field to be extracted from mhObj and set to mtObj")
	}
}

func TestConstantCallSite_Init_NullParams(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	callSiteObj := CreateRawConstantCallSiteObject()

	// Null target
	res1 := constantCallSiteInit([]any{callSiteObj, nil})
	errBlk1, ok := res1.(*ghelpers.GErrBlk)
	if !ok || errBlk1.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null target, got %v", res1)
	}

	// Null receiver (this)
	res2 := constantCallSiteInit([]any{nil, callSiteObj})
	errBlk2, ok := res2.(*ghelpers.GErrBlk)
	if !ok || errBlk2.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null receiver, got %v", res2)
	}

	// Insufficient params
	res3 := constantCallSiteInit([]any{callSiteObj})
	errBlk3, ok := res3.(*ghelpers.GErrBlk)
	if !ok || errBlk3.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for insufficient parameters, got %v", res3)
	}
}

func TestConstantCallSite_InitWithHook_Success(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	mtObj := createTestDummyMethodType()
	hookObj := createTestDummyMethodHandle(mtObj)
	callSiteObj := CreateRawConstantCallSiteObject()

	res := constantCallSiteInitWithHook([]any{callSiteObj, mtObj, hookObj})
	if res != nil {
		t.Fatalf("Expected nil return on successful init with hook, got %v", res)
	}

	if callSiteObj.FieldTable["type"].Fvalue != mtObj {
		t.Errorf("Expected type field to match passed mtObj")
	}
	if callSiteObj.FieldTable["target"].Fvalue != hookObj {
		t.Errorf("Expected target field to match passed hookObj")
	}
}

func TestConstantCallSite_InitWithHook_NullParams(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	callSiteObj := CreateRawConstantCallSiteObject()
	mtObj := createTestDummyMethodType()
	hookObj := createTestDummyMethodHandle(mtObj)

	// Null type
	res1 := constantCallSiteInitWithHook([]any{callSiteObj, nil, hookObj})
	if errBlk, ok := res1.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null type, got %v", res1)
	}

	// Null hook
	res2 := constantCallSiteInitWithHook([]any{callSiteObj, mtObj, nil})
	if errBlk, ok := res2.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null hook, got %v", res2)
	}
}

func TestConstantCallSite_GetTarget_And_DynamicInvoker(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	mtObj := createTestDummyMethodType()
	mhObj := createTestDummyMethodHandle(mtObj)

	callSiteObj, err := MakeConstantCallSiteObject(mhObj)
	if err != nil || callSiteObj == nil {
		t.Fatalf("MakeConstantCallSiteObject failed: %v", err)
	}

	// getTarget
	resTarget := constantCallSiteGetTarget([]any{callSiteObj})
	if resTarget != mhObj {
		t.Errorf("Expected getTarget to return mhObj, got %v", resTarget)
	}

	// dynamicInvoker
	resInvoker := constantCallSiteDynamicInvoker([]any{callSiteObj})
	if resInvoker != mhObj {
		t.Errorf("Expected dynamicInvoker to return mhObj, got %v", resInvoker)
	}

	// Uninitialized target
	emptyCSO := CreateRawConstantCallSiteObject()
	resUninit := constantCallSiteGetTarget([]any{emptyCSO})
	if errBlk, ok := resUninit.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalStateException {
		t.Errorf("Expected IllegalStateException for uninitialized target, got %v", resUninit)
	}
}

func TestConstantCallSite_SetTarget_ThrowsUnsupportedOperationException(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	mtObj := createTestDummyMethodType()
	mhObj := createTestDummyMethodHandle(mtObj)
	callSiteObj, _ := MakeConstantCallSiteObject(mhObj)

	newMH := createTestDummyMethodHandle(mtObj)
	res := constantCallSiteSetTarget([]any{callSiteObj, newMH})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.UnsupportedOperationException {
		t.Fatalf("Expected UnsupportedOperationException on setTarget, got %v", res)
	}
}

func TestCallSite_Type(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	mtObj := createTestDummyMethodType()
	mhObj := createTestDummyMethodHandle(mtObj)

	callSiteObj, _ := MakeConstantCallSiteObject(mhObj)

	res := callSiteType([]any{callSiteObj})
	if res != mtObj {
		t.Errorf("Expected callSiteType to return mtObj, got %v", res)
	}

	// Uninitialized type and target
	emptyCSO := CreateRawConstantCallSiteObject()
	resEmpty := callSiteType([]any{emptyCSO})
	if errBlk, ok := resEmpty.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalStateException {
		t.Errorf("Expected IllegalStateException for uninitialized CallSite type, got %v", resEmpty)
	}
}

func TestMakeConstantCallSiteObject_NullTarget(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	cso, err := MakeConstantCallSiteObject(nil)
	if cso != nil {
		t.Errorf("Expected nil object on null target, got %v", cso)
	}
	if errBlk, ok := err.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException on null target, got %v", err)
	}
}
