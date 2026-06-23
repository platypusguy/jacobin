/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025-6 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

func Load_Lang_Invoke_MethodHandle() {
	ghelpers.MethodSignatures["java/lang/invoke/MethodHandle.type()Ljava/lang/invoke/MethodType;"] =
		ghelpers.GMeth{ParamSlots: 0, NeedsContext: true, GFunction: mhType}
}

// Internal representation of java.lang.invoke.MethodHandle
// type MethodHandle struct {
// 	Kind          MethodHandleKind    // REF_getField, REF_invokeVirtual, etc.
// 	RefClass      string              // Declaring class
// 	RefName       string              // Method/field name
// 	RefDescriptor string              // Method/field descriptor
// 	DirectMethod  *classloader.Method // For direct method invocations
// 	IsVarArgs     bool
// }

// Method handle reference kinds (JVM spec §5.4.3.5)
type MethodHandleKind uint16

const (
	REF_getField         MethodHandleKind = 1
	REF_getStatic        MethodHandleKind = 2
	REF_putField         MethodHandleKind = 3
	REF_putStatic        MethodHandleKind = 4
	REF_invokeVirtual    MethodHandleKind = 5
	REF_invokeStatic     MethodHandleKind = 6
	REF_invokeSpecial    MethodHandleKind = 7
	REF_newInvokeSpecial MethodHandleKind = 8
	REF_invokeInterface  MethodHandleKind = 9
)

// CallSite represents a resolved invokedynamic call site
type CallSite struct {
	Target     *object.Object // The method handle object to invoke
	Type       *MethodType    // Expected signature
	IsVolatile bool           // MutableCallSite vs ConstantCallSite
}

// MethodType represents a method signature
type MethodType struct {
	ReturnType string
	ParamTypes []string
}

// createRawMethodHandleObject creates a methodHandle object with default values
func createRawMethodHandleObject() *object.Object {
	mhClassName := "java/lang/invoke/MethodHandle"
	mho := object.MakeEmptyObjectWithClassName(&mhClassName)

	mho.FieldTable["Kind"] = object.Field{Ftype: types.U16, Fvalue: 0}
	mho.FieldTable["RefClass"] = object.Field{Ftype: types.GolangString, Fvalue: ""}
	mho.FieldTable["RefName"] = object.Field{Ftype: types.GolangString, Fvalue: ""}
	mho.FieldTable["RefDescriptor"] = object.Field{Ftype: types.GolangString, Fvalue: ""}
	mho.FieldTable["DirectMethod"] = object.Field{Ftype: types.Ref, Fvalue: nil}
	mho.FieldTable["IsVarArgs"] = object.Field{Ftype: types.Bool, Fvalue: false}
	mho.FieldTable["type"] = object.Field{Ftype: types.Ref, Fvalue: nil}
	return mho
}

// type returns the type of the method handle as a MethodType object
func mhType(params []interface{}) interface{} {
	if params == nil {
		errMsg := fmt.Sprintf("mhType(): Invalid params array passed in")
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if len(params) != 1 {
		errMsg := fmt.Sprintf("mhType(): Expected 1 parameter, got %d", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	mho := params[0].(*object.Object)
	if *stringPool.GetStringPointer(mho.KlassName) != "java/lang/invoke/MethodHandle" {
		errMsg := fmt.Sprintf("mhType(): Expected MethodHandle object, got %s",
			*stringPool.GetStringPointer(mho.KlassName))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return mho.FieldTable["type"].Fvalue.(*object.Object)
}
