/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_Object() {

	MethodSignatures["java/lang/Object.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Object.getClass()Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectGetClass,
		}

	MethodSignatures["java/lang/Object.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectToString,
		}

}

// === the internal representation of a java.lang.Class() instance ===
// this is not a faithful reproduction of the OpenJDK version, but rather
// the one we use in Jacobin
type javaLangClass struct {
	accessFlags    classloader.AccessFlags
	name           string
	superClassName string
	interfaceNames []string
	constantPool   classloader.CPool
	fields         []classloader.Field
	methods        map[string]*classloader.Method
	loader         string
	superClass     string
	interfaces     []uint16 // indices into UTF8Refs
	// instanceSlotCount uint
	// staticSlotCount   uint
	// staticVars        Slots
}

// "java/lang/Object.getClass()Ljava/lang/Class;"
func objectGetClass(params []interface{}) interface{} {
	objPtr := params[0].(*object.Object)
	if objPtr == nil || objPtr.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("Invalid object in objectGetClass(): %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jlc := javaLangClass{}
	jlc.name = object.GoStringFromStringPoolIndex(objPtr.KlassName)

	obj := *classloader.MethAreaFetch(jlc.name)
	jlc.loader = obj.Loader

	objData := *obj.Data
	jlc.constantPool = objData.CP
	jlc.superClass = object.GoStringFromStringPoolIndex(objData.SuperclassIndex)
	jlc.fields = objData.Fields
	jlc.interfaces = objData.Interfaces
	jlc.methods = objData.MethodTable
	jlc.accessFlags = objData.Access
	return &jlc
}

// "java/lang/Object.toString()Ljava/lang/String;"
func objectToString(params []interface{}) interface{} {
	// params[0]: input Object
	var str string

	switch params[0].(type) {
	case *object.Object:
		inObj := params[0].(*object.Object)
		str = object.ObjectFieldToString(inObj, "value")
		return object.StringObjectFromGoString(str)
	}

	errMsg := fmt.Sprintf("Unsupported parameter type: %T", params[0])
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)
}
