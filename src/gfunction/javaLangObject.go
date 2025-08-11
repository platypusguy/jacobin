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
	"unsafe"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_Object() {

	MethodSignatures["java/lang/Object.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Object.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// "java/lang/Object.clone(Ljava/lang/Object;)Ljava/lang/Object;" is PROTECTED

	MethodSignatures["java/lang/Object.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  objectEquals,
		}

	MethodSignatures["java/lang/Object.finalize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Object.getClass()Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectGetClass,
		}

	MethodSignatures["java/lang/Object.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectHashCode,
		}

	MethodSignatures["java/lang/Object.notify()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Object.notifyAll()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Object.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  objectToString,
		}

	MethodSignatures["java/lang/Object.wait()V"] = // wait until awakened, typically by being notified or interrupted
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Object.wait(J)V"] = // wait(long timeoutMillis)
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Object.wait(JI)V"] = // wait(long timeoutMillis, int nanos)
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
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
		errMsg := fmt.Sprintf("objectGetClass: Invalid object in objectGetClass(): %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jlc := javaLangClass{}
	jlc.name = object.GoStringFromStringPoolIndex(objPtr.KlassName)

	// get a pointer to the class contents from the method area
	o := classloader.MethAreaFetch(jlc.name)
	if o == nil {
		errMsg := fmt.Sprintf("objectGetClass: Class %s not loaded", jlc.name)
		return getGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}

	// syntactic sugar
	obj := *o

	// create the empty java.lang.Class structure
	jlc.loader = obj.Loader

	// fill in the jlc
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

	switch params[0].(type) {
	case *object.Object:
		inObj := params[0].(*object.Object)
		classNameSuffix := object.GetClassNameSuffix(inObj, false)
		if classNameSuffix == "LinkedList" {
			return linkedlistToString(params)
		}
		return object.StringifyAnythingJava(inObj)
	}

	errMsg := fmt.Sprintf("objectToString: Unsupported parameter type: %T", params[0])
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)
}

// "java/lang/Object.hashCode()I"
func objectHashCode(params []interface{}) interface{} {
	// params[0]: input Object
	switch params[0].(type) {
	case *object.Object:
		ptr := uintptr(unsafe.Pointer(params[0].(*object.Object)))
		hashCode := int64(ptr ^ (ptr >> 32))
		return hashCode
	}

	errMsg := fmt.Sprintf("objectHashCode: Unsupported parameter type: %T", params[0])
	return getGErrBlk(excNames.IllegalArgumentException, errMsg)
}

func objectEquals(params []interface{}) interface{} {
	this, ok := params[0].(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}
	that, ok := params[1].(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}

	// If they are the same object, even if null, return true.
	if this == that {
		return types.JavaBoolTrue
	}

	// Not the same object.
	return types.JavaBoolFalse
}
