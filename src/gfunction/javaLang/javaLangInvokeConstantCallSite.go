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
	"jacobin/src/object"
	"jacobin/src/stringPool"
)

const (
	ConstantCallSiteClassName = "java/lang/invoke/ConstantCallSite"
	CallSiteClassName         = "java/lang/invoke/CallSite"
)

func Load_Lang_Invoke_ConstantCallSite() {
	// Constructors
	ghelpers.MethodSignatures["java/lang/invoke/ConstantCallSite.<init>(Ljava/lang/invoke/MethodHandle;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  constantCallSiteInit,
		}

	ghelpers.MethodSignatures["java/lang/invoke/ConstantCallSite.<init>(Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodHandle;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  constantCallSiteInitWithHook,
		}

	// Instance methods on ConstantCallSite
	ghelpers.MethodSignatures["java/lang/invoke/ConstantCallSite.getTarget()Ljava/lang/invoke/MethodHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  constantCallSiteGetTarget,
		}

	ghelpers.MethodSignatures["java/lang/invoke/ConstantCallSite.setTarget(Ljava/lang/invoke/MethodHandle;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  constantCallSiteSetTarget,
		}

	ghelpers.MethodSignatures["java/lang/invoke/ConstantCallSite.dynamicInvoker()Ljava/lang/invoke/MethodHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  constantCallSiteDynamicInvoker,
		}

	ghelpers.MethodSignatures["java/lang/invoke/ConstantCallSite.type()Ljava/lang/invoke/MethodType;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  callSiteType,
		}

	// Methods on base CallSite
	ghelpers.MethodSignatures["java/lang/invoke/CallSite.type()Ljava/lang/invoke/MethodType;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  callSiteType,
		}

	ghelpers.MethodSignatures["java/lang/invoke/CallSite.getTarget()Ljava/lang/invoke/MethodHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  constantCallSiteGetTarget,
		}

	ghelpers.MethodSignatures["java/lang/invoke/CallSite.setTarget(Ljava/lang/invoke/MethodHandle;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  constantCallSiteSetTarget,
		}

	ghelpers.MethodSignatures["java/lang/invoke/CallSite.dynamicInvoker()Ljava/lang/invoke/MethodHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  constantCallSiteDynamicInvoker,
		}
}

// CreateRawConstantCallSiteObject creates a ConstantCallSite object with uninitialized fields.
func CreateRawConstantCallSiteObject() *object.Object {
	className := ConstantCallSiteClassName
	ccso := object.MakeEmptyObjectWithClassName(&className)
	ccso.FieldTable["target"] = object.Field{Ftype: "Ljava/lang/invoke/MethodHandle;", Fvalue: nil}
	ccso.FieldTable["type"] = object.Field{Ftype: "Ljava/lang/invoke/MethodType;", Fvalue: nil}
	return ccso
}

// MakeConstantCallSiteObject creates an initialized ConstantCallSite object binding the given permanent target.
func MakeConstantCallSiteObject(target *object.Object) (*object.Object, any) {
	if target == nil || object.IsNull(target) {
		errMsg := "MakeConstantCallSiteObject: target cannot be null"
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	ccso := CreateRawConstantCallSiteObject()
	ccso.FieldTable["target"] = object.Field{Ftype: "Ljava/lang/invoke/MethodHandle;", Fvalue: target}

	target.ThMutex.RLock()
	if tField, ok := target.FieldTable["type"]; ok && tField.Fvalue != nil {
		ccso.FieldTable["type"] = tField
	} else if rField, ok := target.FieldTable["RefDescriptor"]; ok && rField.Fvalue != nil {
		ccso.FieldTable["type"] = rField
	}
	target.ThMutex.RUnlock()

	return ccso, nil
}

// constantCallSiteInit: ConstantCallSite.<init>(MethodHandle target)
// Creates a call site with a permanent target.
// Throws NullPointerException if target is null.
func constantCallSiteInit(params []any) any {
	if params == nil || len(params) < 2 {
		errMsg := fmt.Sprintf("constantCallSiteInit: expected at least 2 parameters (this, target), got %d", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil || object.IsNull(thisObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteInit: receiver object is null")
	}

	targetObj, ok := params[1].(*object.Object)
	if !ok || targetObj == nil || object.IsNull(targetObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteInit: target cannot be null")
	}

	thisObj.ThMutex.Lock()
	if thisObj.FieldTable == nil {
		thisObj.FieldTable = make(map[string]object.Field)
	}
	thisObj.FieldTable["target"] = object.Field{Ftype: "Ljava/lang/invoke/MethodHandle;", Fvalue: targetObj}

	targetObj.ThMutex.RLock()
	if tField, ok := targetObj.FieldTable["type"]; ok && tField.Fvalue != nil {
		thisObj.FieldTable["type"] = tField
	} else if rField, ok := targetObj.FieldTable["RefDescriptor"]; ok && rField.Fvalue != nil {
		thisObj.FieldTable["type"] = rField
	}
	targetObj.ThMutex.RUnlock()
	thisObj.ThMutex.Unlock()

	return nil
}

// constantCallSiteInitWithHook: ConstantCallSite.<init>(MethodType type, MethodHandle createTargetHook)
// Creates a call site with a permanent target, possibly bound to the call site itself.
// Throws NullPointerException if type or createTargetHook is null.
func constantCallSiteInitWithHook(params []any) any {
	if params == nil || len(params) < 3 {
		errMsg := fmt.Sprintf("constantCallSiteInitWithHook: expected at least 3 parameters, got %d", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil || object.IsNull(thisObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteInitWithHook: receiver object is null")
	}

	typeObj, ok := params[1].(*object.Object)
	if !ok || typeObj == nil || object.IsNull(typeObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteInitWithHook: type cannot be null")
	}

	hookObj, ok := params[2].(*object.Object)
	if !ok || hookObj == nil || object.IsNull(hookObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteInitWithHook: createTargetHook cannot be null")
	}

	thisObj.ThMutex.Lock()
	if thisObj.FieldTable == nil {
		thisObj.FieldTable = make(map[string]object.Field)
	}
	thisObj.FieldTable["type"] = object.Field{Ftype: "Ljava/lang/invoke/MethodType;", Fvalue: typeObj}
	thisObj.FieldTable["target"] = object.Field{Ftype: "Ljava/lang/invoke/MethodHandle;", Fvalue: hookObj}
	thisObj.ThMutex.Unlock()

	return nil
}

// constantCallSiteGetTarget: MethodHandle getTarget()
// Returns the target method of the call site, which behaves like a final field in ConstantCallSite.
// Throws IllegalStateException if the target is uninitialized.
func constantCallSiteGetTarget(params []any) any {
	if params == nil || len(params) < 1 {
		errMsg := "constantCallSiteGetTarget: expected receiver object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil || object.IsNull(thisObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteGetTarget: receiver object is null")
	}

	thisObj.ThMutex.RLock()
	targetField, ok := thisObj.FieldTable["target"]
	thisObj.ThMutex.RUnlock()

	if !ok || targetField.Fvalue == nil || object.IsNull(targetField.Fvalue) {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "constantCallSiteGetTarget: target is not initialized")
	}

	targetObj, ok := targetField.Fvalue.(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "constantCallSiteGetTarget: target field is not a MethodHandle object")
	}

	return targetObj
}

// constantCallSiteSetTarget: void setTarget(MethodHandle ignore)
// Always throws UnsupportedOperationException because this kind of call site cannot change its target.
func constantCallSiteSetTarget(params []any) any {
	if params == nil || len(params) < 1 {
		errMsg := "constantCallSiteSetTarget: expected receiver object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil || object.IsNull(thisObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "constantCallSiteSetTarget: receiver object is null")
	}

	return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "ConstantCallSite cannot change its target")
}

// constantCallSiteDynamicInvoker: MethodHandle dynamicInvoker()
// Returns this call site's permanent target.
func constantCallSiteDynamicInvoker(params []any) any {
	return constantCallSiteGetTarget(params)
}

// callSiteType: MethodType type()
// Returns the type of this call site's target.
func callSiteType(params []any) any {
	if params == nil || len(params) < 1 {
		errMsg := "callSiteType: expected receiver object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil || object.IsNull(thisObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "callSiteType: receiver object is null")
	}

	thisObj.ThMutex.RLock()
	typeField, hasType := thisObj.FieldTable["type"]
	targetField, hasTarget := thisObj.FieldTable["target"]
	thisObj.ThMutex.RUnlock()

	if hasType && typeField.Fvalue != nil && !object.IsNull(typeField.Fvalue) {
		if typeObj, ok := typeField.Fvalue.(*object.Object); ok && typeObj != nil {
			return typeObj
		}
	}

	if hasTarget && targetField.Fvalue != nil && !object.IsNull(targetField.Fvalue) {
		if targetObj, ok := targetField.Fvalue.(*object.Object); ok && targetObj != nil {
			targetObj.ThMutex.RLock()
			tField, ok := targetObj.FieldTable["type"]
			if !ok || tField.Fvalue == nil {
				tField, ok = targetObj.FieldTable["RefDescriptor"]
			}
			targetObj.ThMutex.RUnlock()

			if ok && tField.Fvalue != nil && !object.IsNull(tField.Fvalue) {
				if typeObj, ok := tField.Fvalue.(*object.Object); ok && typeObj != nil {
					return typeObj
				}
			}
		}
	}

	cls := ""
	if p := stringPool.GetStringPointer(thisObj.KlassName); p != nil {
		cls = *p
	}
	return ghelpers.GetGErrBlk(excNames.IllegalStateException, fmt.Sprintf("callSiteType: type is not initialized for %s", cls))
}
