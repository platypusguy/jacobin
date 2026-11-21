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
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"testing"
)

func createTestLookupObject() *object.Object {
	lookupClassName := "java/lang/invoke/MethodHandles$Lookup"
	lookupObj := object.MakeEmptyObjectWithClassName(&lookupClassName)
	callerClassObj := object.MakeEmptyObject()
	lookupObj.FieldTable["lookupClass"] = object.Field{Ftype: types.Ref, Fvalue: callerClassObj}
	return lookupObj
}

func createTestMethodTypeObject(paramCount int) *object.Object {
	mtClassName := "java/lang/invoke/MethodType"
	mtObj := object.MakeEmptyObjectWithClassName(&mtClassName)

	paramArray := object.Make1DimRefArray("java/lang/Class", int64(paramCount))
	mtObj.FieldTable["rtype"] = object.Field{Ftype: "Ljava/lang/Class;", Fvalue: object.MakeEmptyObject()}
	mtObj.FieldTable["ptypes"] = object.Field{Ftype: "[Ljava/lang/Class;", Fvalue: paramArray}
	return mtObj
}

func createTestDirectMethodHandle(kind int64) *object.Object {
	mho := createRawMethodHandleObject()
	mho.FieldTable["Kind"] = object.Field{Ftype: types.Int, Fvalue: kind}
	mho.FieldTable["RefClass"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("test/TargetClass")}
	mho.FieldTable["RefName"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("targetMethod")}
	mho.FieldTable["RefDescriptor"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("()V")}
	mho.FieldTable["DirectMethod"] = object.Field{Ftype: types.Bool, Fvalue: true}
	return mho
}

func TestLoad_Lang_Invoke_LambdaMetafactory_Registration(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Lang_Invoke_LambdaMetafactory()

	expectedSignatures := []string{
		"java/lang/invoke/LambdaMetafactory.<clinit>()V",
		"java/lang/invoke/LambdaMetafactory.metafactory(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodHandle;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/CallSite;",
		"java/lang/invoke/LambdaMetafactory.altMetafactory(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected method signature to be registered: %s", sig)
		}
	}
}

func TestLambdaMetafactory_Metafactory_Success(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	caller := createTestLookupObject()
	interfaceMethodName := object.StringObjectFromGoString("accept")
	factoryType := createTestMethodTypeObject(0)
	interfaceMethodType := createTestMethodTypeObject(1)
	implementation := createTestDirectMethodHandle(int64(REF_invokeStatic))
	dynamicMethodType := createTestMethodTypeObject(1)

	params := []any{
		caller,
		interfaceMethodName,
		factoryType,
		interfaceMethodType,
		implementation,
		dynamicMethodType,
	}

	res := lambdaMetafactory_metafactory(params)
	callSite, ok := res.(*object.Object)
	if !ok || callSite == nil {
		t.Fatalf("Expected *object.Object CallSite return, got %T (%v)", res, res)
	}

	targetField, ok := callSite.FieldTable["target"]
	if !ok || targetField.Fvalue == nil {
		t.Fatalf("Expected CallSite to contain 'target' field")
	}
	mh, ok := targetField.Fvalue.(*object.Object)
	if !ok || mh == nil {
		t.Fatalf("Expected target to be *object.Object MethodHandle")
	}

	if *stringPool.GetStringPointer(mh.KlassName) != "java/lang/invoke/MethodHandle" {
		t.Errorf("Expected target MethodHandle class name, got %s", *stringPool.GetStringPointer(mh.KlassName))
	}

	typeField, ok := callSite.FieldTable["type"]
	if !ok || typeField.Fvalue != factoryType {
		t.Errorf("Expected CallSite type to be factoryType")
	}
}

func TestLambdaMetafactory_Metafactory_ParamCountError(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	res := lambdaMetafactory_metafactory([]any{createTestLookupObject()})
	errBlk, ok := res.(*ghelpers.GErrBlk)
	if !ok || errBlk.ExceptionType != excNames.MalformedParametersException {
		t.Fatalf("Expected MalformedParametersException, got %v", res)
	}
}

func TestLambdaMetafactory_Metafactory_NullParamErrors(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	validCaller := createTestLookupObject()
	validName := object.StringObjectFromGoString("run")
	validFactoryType := createTestMethodTypeObject(0)
	validIfaceType := createTestMethodTypeObject(0)
	validImpl := createTestDirectMethodHandle(int64(REF_invokeStatic))
	validDynType := createTestMethodTypeObject(0)

	tests := []struct {
		name   string
		params []any
	}{
		{"null_caller", []any{nil, validName, validFactoryType, validIfaceType, validImpl, validDynType}},
		{"null_name", []any{validCaller, nil, validFactoryType, validIfaceType, validImpl, validDynType}},
		{"null_factory_type", []any{validCaller, validName, nil, validIfaceType, validImpl, validDynType}},
		{"null_iface_type", []any{validCaller, validName, validFactoryType, nil, validImpl, validDynType}},
		{"null_impl", []any{validCaller, validName, validFactoryType, validIfaceType, nil, validDynType}},
		{"null_dyn_type", []any{validCaller, validName, validFactoryType, validIfaceType, validImpl, nil}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := lambdaMetafactory_metafactory(tc.params)
			errBlk, ok := res.(*ghelpers.GErrBlk)
			if !ok || errBlk.ExceptionType != excNames.NullPointerException {
				t.Fatalf("%s: expected NullPointerException, got %v", tc.name, res)
			}
		})
	}
}

func TestLambdaMetafactory_Metafactory_ValidationErrors(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	validCaller := createTestLookupObject()
	validName := object.StringObjectFromGoString("apply")
	validFactoryType := createTestMethodTypeObject(0)
	validIfaceType := createTestMethodTypeObject(1)
	validImpl := createTestDirectMethodHandle(int64(REF_invokeStatic))
	validDynType := createTestMethodTypeObject(1)

	tests := []struct {
		name         string
		params       []any
		expectedExcp int
	}{
		{
			name: "invalid_caller_not_lookup",
			params: []any{
				object.MakeEmptyObject(), // empty object without lookupClass or Lookup class name
				validName,
				validFactoryType,
				validIfaceType,
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_interface_method_name_empty",
			params: []any{
				validCaller,
				object.StringObjectFromGoString(""),
				validFactoryType,
				validIfaceType,
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_interface_method_name_init",
			params: []any{
				validCaller,
				object.StringObjectFromGoString("<init>"),
				validFactoryType,
				validIfaceType,
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_interface_method_name_clinit",
			params: []any{
				validCaller,
				object.StringObjectFromGoString("<clinit>"),
				validFactoryType,
				validIfaceType,
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_interface_method_name_special_chars",
			params: []any{
				validCaller,
				object.StringObjectFromGoString("apply.method"),
				validFactoryType,
				validIfaceType,
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_factory_type_not_methodtype",
			params: []any{
				validCaller,
				validName,
				object.MakeEmptyObject(),
				validIfaceType,
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_iface_type_not_methodtype",
			params: []any{
				validCaller,
				validName,
				validFactoryType,
				object.MakeEmptyObject(),
				validImpl,
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_dyn_type_not_methodtype",
			params: []any{
				validCaller,
				validName,
				validFactoryType,
				validIfaceType,
				validImpl,
				object.MakeEmptyObject(),
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_impl_not_direct_method_handle",
			params: []any{
				validCaller,
				validName,
				validFactoryType,
				validIfaceType,
				object.MakeEmptyObject(), // Kind missing
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "invalid_impl_kind_zero",
			params: []any{
				validCaller,
				validName,
				validFactoryType,
				validIfaceType,
				createTestDirectMethodHandle(0),
				validDynType,
			},
			expectedExcp: excNames.LambdaConversionException,
		},
		{
			name: "param_count_mismatch_between_iface_and_dyn",
			params: []any{
				validCaller,
				validName,
				validFactoryType,
				createTestMethodTypeObject(1), // 1 parameter
				validImpl,
				createTestMethodTypeObject(2), // 2 parameters
			},
			expectedExcp: excNames.LambdaConversionException,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := lambdaMetafactory_metafactory(tc.params)
			errBlk, ok := res.(*ghelpers.GErrBlk)
			if !ok || errBlk.ExceptionType != tc.expectedExcp {
				t.Fatalf("%s: expected exception %d, got %v", tc.name, tc.expectedExcp, res)
			}
		})
	}
}
