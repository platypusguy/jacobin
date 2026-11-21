/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"strings"
)

const (
	FLAG_SERIALIZABLE = 1 << 0
	FLAG_MARKERS      = 1 << 1
	FLAG_BRIDGES      = 1 << 2
)

func Load_Lang_Invoke_LambdaMetafactory() {
	ghelpers.MethodSignatures["java/lang/invoke/LambdaMetafactory.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  lambdaMetafactoryClinit,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaMetafactory.metafactory(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodHandle;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/CallSite;"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  lambdaMetafactory_metafactory,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaMetafactory.altMetafactory(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  lambdaMetafactory_altMetafactory,
		}
}

func lambdaMetafactoryClinit(_ []any) any {
	_ = statics.AddStatic("java/lang/invoke/LambdaMetafactory.FLAG_BRIDGES",
		statics.Static{
			Type:  types.Int,
			Value: 4,
		})

	_ = statics.AddStatic("java/lang/invoke/LambdaMetafactory.FLAG_MARKERS",
		statics.Static{
			Type:  types.Int,
			Value: 2,
		})

	_ = statics.AddStatic("java/lang/invoke/LambdaMetafactory.FLAG_SERIALIZABLE",
		statics.Static{
			Type:  types.Int,
			Value: 1,
		})

	return nil
}

// lambdaMetafactory_metafactory implements standard 6-parameter metafactory bootstrap method:
// metafactory(MethodHandles.Lookup caller, String interfaceMethodName, MethodType factoryType,
//
//	MethodType interfaceMethodType, MethodHandle implementation, MethodType dynamicMethodType)
func lambdaMetafactory_metafactory(params []any) any {
	if len(params) != 6 {
		errMsg := fmt.Sprintf("lambdaMetafactory_metafactory: incorrect parameter count: %d (expected 6)", len(params))
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	// 1. Null checks
	caller, ok0 := params[0].(*object.Object)
	if !ok0 || caller == nil || object.IsNull(caller) {
		errMsg := "lambdaMetafactory_metafactory: caller cannot be null"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	interfaceMethodNameObj, ok1 := params[1].(*object.Object)
	if !ok1 || interfaceMethodNameObj == nil || object.IsNull(interfaceMethodNameObj) {
		errMsg := "lambdaMetafactory_metafactory: interfaceMethodName cannot be null"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	factoryTypeObj, ok2 := params[2].(*object.Object)
	if !ok2 || factoryTypeObj == nil || object.IsNull(factoryTypeObj) {
		errMsg := "lambdaMetafactory_metafactory: factoryType cannot be null"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	interfaceMethodTypeObj, ok3 := params[3].(*object.Object)
	if !ok3 || interfaceMethodTypeObj == nil || object.IsNull(interfaceMethodTypeObj) {
		errMsg := "lambdaMetafactory_metafactory: interfaceMethodType cannot be null"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	implObj, ok4 := params[4].(*object.Object)
	if !ok4 || implObj == nil || object.IsNull(implObj) {
		errMsg := "lambdaMetafactory_metafactory: implementation cannot be null"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	dynamicMethodTypeObj, ok5 := params[5].(*object.Object)
	if !ok5 || dynamicMethodTypeObj == nil || object.IsNull(dynamicMethodTypeObj) {
		errMsg := "lambdaMetafactory_metafactory: dynamicMethodType cannot be null"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	// 2. Validate Lookup object
	if !isLookupObject(caller) {
		errMsg := "lambdaMetafactory_metafactory: caller is not a valid MethodHandles.Lookup object"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	// 3. Validate interface method name
	methName := object.GoStringFromStringObject(interfaceMethodNameObj)
	if methName == "" || methName == "<init>" || methName == "<clinit>" || strings.ContainsAny(methName, ".;/[]") {
		errMsg := fmt.Sprintf("lambdaMetafactory_metafactory: invalid interface method name: %q", methName)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	// 4. Validate MethodType objects
	if !isMethodTypeObject(factoryTypeObj) {
		errMsg := "lambdaMetafactory_metafactory: factoryType is not a valid MethodType object"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}
	if !isMethodTypeObject(interfaceMethodTypeObj) {
		errMsg := "lambdaMetafactory_metafactory: interfaceMethodType is not a valid MethodType object"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}
	if !isMethodTypeObject(dynamicMethodTypeObj) {
		errMsg := "lambdaMetafactory_metafactory: dynamicMethodType is not a valid MethodType object"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	// 5. Validate MethodHandle implementation object
	if !isDirectMethodHandle(implObj) {
		errMsg := "lambdaMetafactory_metafactory: implementation must be a direct MethodHandle"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	// 6. Arity and type compatibility checks
	ifaceParamCount := getMethodTypeParamCount(interfaceMethodTypeObj)
	dynParamCount := getMethodTypeParamCount(dynamicMethodTypeObj)
	if ifaceParamCount != dynParamCount {
		errMsg := fmt.Sprintf("lambdaMetafactory_metafactory: parameter count mismatch: interfaceMethodType has %d, dynamicMethodType has %d",
			ifaceParamCount, dynParamCount)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	if globals.TraceInst {
		trace.Raw(fmt.Sprintf("lambdaMetafactory: caller=%s, methName=%s, factoryType=%s, ifaceParamCount=%d",
			object.GoStringFromStringObject(caller), methName,
			object.StringifyAnythingGo(factoryTypeObj), ifaceParamCount))
	}

	// 7. Synthesize target MethodHandle and encapsulate in ConstantCallSite
	mh := createRawMethodHandleObject()
	mh.FieldTable["Kind"] = object.Field{Ftype: types.Int, Fvalue: int64(REF_invokeStatic)}
	mh.FieldTable["RefClass"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString("java/lang/invoke/LambdaMetafactory")}
	mh.FieldTable["RefName"] = object.Field{Ftype: types.Ref, Fvalue: interfaceMethodNameObj}
	mh.FieldTable["RefDescriptor"] = object.Field{Ftype: types.Ref, Fvalue: factoryTypeObj}
	mh.FieldTable["DirectMethod"] = object.Field{Ftype: types.Bool, Fvalue: true}
	mh.FieldTable["IsVarArgs"] = object.Field{Ftype: types.Bool, Fvalue: false}
	mh.FieldTable["type"] = object.Field{Ftype: types.Ref, Fvalue: factoryTypeObj}
	if targetField, ok := implObj.FieldTable["$target"]; ok {
		mh.FieldTable["$target"] = targetField
	}

	callSite := CreateRawConstantCallSiteObject()
	callSite.FieldTable["target"] = object.Field{Ftype: "Ljava/lang/invoke/MethodHandle;", Fvalue: mh}
	callSite.FieldTable["type"] = object.Field{Ftype: "Ljava/lang/invoke/MethodType;", Fvalue: factoryTypeObj}

	return callSite
}

// lambdaMetafactory_altMetafactory handles the variable-arity bootstrap method:
// altMetafactory(MethodHandles.Lookup caller, String interfaceMethodName, MethodType factoryType, Object... args)
func lambdaMetafactory_altMetafactory(params []any) any {
	if len(params) < 4 {
		errMsg := fmt.Sprintf("lambdaMetafactory_altMetafactory: insufficient parameter count: %d (expected at least 4)", len(params))
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	caller, ok0 := params[0].(*object.Object)
	if !ok0 || caller == nil || object.IsNull(caller) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lambdaMetafactory_altMetafactory: caller cannot be null")
	}

	interfaceMethodNameObj, ok1 := params[1].(*object.Object)
	if !ok1 || interfaceMethodNameObj == nil || object.IsNull(interfaceMethodNameObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lambdaMetafactory_altMetafactory: interfaceMethodName cannot be null")
	}

	factoryTypeObj, ok2 := params[2].(*object.Object)
	if !ok2 || factoryTypeObj == nil || object.IsNull(factoryTypeObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lambdaMetafactory_altMetafactory: factoryType cannot be null")
	}

	// Trailing args can arrive packed in params[3] as []*object.Object / []any or flat in params[3:]
	var staticArgs []any
	if len(params) == 4 {
		if arr, ok := params[3].([]any); ok {
			staticArgs = arr
		} else if arrObj, ok := params[3].([]*object.Object); ok {
			staticArgs = make([]any, len(arrObj))
			for i, o := range arrObj {
				staticArgs[i] = o
			}
		} else {
			staticArgs = []any{params[3]}
		}
	} else {
		staticArgs = params[3:]
	}

	if len(staticArgs) < 4 {
		errMsg := fmt.Sprintf("lambdaMetafactory_altMetafactory: expected at least 4 extra static arguments, got %d", len(staticArgs))
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	interfaceMethodTypeObj, ok3 := staticArgs[0].(*object.Object)
	implObj, ok4 := staticArgs[1].(*object.Object)
	dynamicMethodTypeObj, ok5 := staticArgs[2].(*object.Object)
	flagsVal := staticArgs[3]

	if !ok3 || interfaceMethodTypeObj == nil || object.IsNull(interfaceMethodTypeObj) ||
		!ok4 || implObj == nil || object.IsNull(implObj) ||
		!ok5 || dynamicMethodTypeObj == nil || object.IsNull(dynamicMethodTypeObj) ||
		flagsVal == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "lambdaMetafactory_altMetafactory: static arguments contain null")
	}

	flags, ok := ghelpers.GetInt64Val(flagsVal)
	if !ok {
		errMsg := fmt.Sprintf("lambdaMetafactory_altMetafactory: invalid flags value: %v", flagsVal)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}
	validFlags := int64(FLAG_SERIALIZABLE | FLAG_MARKERS | FLAG_BRIDGES)
	if (flags & ^validFlags) != 0 {
		errMsg := fmt.Sprintf("lambdaMetafactory_altMetafactory: invalid flags: %d", flags)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.LambdaConversionException, errMsg)
	}

	// Delegate core validation to metafactory
	return lambdaMetafactory_metafactory([]any{
		caller,
		interfaceMethodNameObj,
		factoryTypeObj,
		interfaceMethodTypeObj,
		implObj,
		dynamicMethodTypeObj,
	})
}

func isLookupObject(obj *object.Object) bool {
	if obj == nil {
		return false
	}
	if _, ok := obj.FieldTable["lookupClass"]; ok {
		return true
	}
	if ptr := stringPool.GetStringPointer(obj.KlassName); ptr != nil && *ptr == "java/lang/invoke/MethodHandles$Lookup" {
		return true
	}
	return false
}

func isMethodTypeObject(obj *object.Object) bool {
	if obj == nil {
		return false
	}
	if _, ok := obj.FieldTable["rtype"]; ok {
		return true
	}
	if _, ok := obj.FieldTable["ptypes"]; ok {
		return true
	}
	if ptr := stringPool.GetStringPointer(obj.KlassName); ptr != nil && *ptr == "java/lang/invoke/MethodType" {
		return true
	}
	return false
}

func isDirectMethodHandle(obj *object.Object) bool {
	if obj == nil {
		return false
	}
	kindField, ok := obj.FieldTable["Kind"]
	if !ok {
		return false
	}
	kind, ok := kindField.Fvalue.(int64)
	if !ok {
		if kInt, isInt := kindField.Fvalue.(int); isInt {
			kind = int64(kInt)
		} else {
			return false
		}
	}
	return kind >= 1 && kind <= 9
}

func getMethodTypeParamCount(mtObj *object.Object) int {
	if mtObj == nil {
		return 0
	}
	ptypesField, ok := mtObj.FieldTable["ptypes"]
	if !ok || ptypesField.Fvalue == nil || object.IsNull(ptypesField.Fvalue) {
		return 0
	}
	ptypesArr, ok := ptypesField.Fvalue.(*object.Object)
	if !ok || ptypesArr == nil {
		return 0
	}
	valField, ok := ptypesArr.FieldTable["value"]
	if !ok || valField.Fvalue == nil {
		return 0
	}
	if slice, ok := valField.Fvalue.([]*object.Object); ok {
		return len(slice)
	}
	return 0
}
