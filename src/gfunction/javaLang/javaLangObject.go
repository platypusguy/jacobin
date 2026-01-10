/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
	"unsafe"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_Object() {

	ghelpers.MethodSignatures["java/lang/Object.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Object.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	// "java/lang/Object.clone(Ljava/lang/Object;)Ljava/lang/Object;" is PROTECTED

	ghelpers.MethodSignatures["java/lang/Object.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  objectEquals,
		}

	ghelpers.MethodSignatures["java/lang/Object.finalize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/lang/Object.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  objectGetClass,
		}

	ghelpers.MethodSignatures["java/lang/Object.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  objectHashCode,
		}

	ghelpers.MethodSignatures["java/lang/Object.notify()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Object.notifyAll()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Object.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  objectToString,
		}

	ghelpers.MethodSignatures["java/lang/Object.wait()V"] = // wait until awakened, typically by being notified or interrupted
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Object.wait(J)V"] = // wait(long timeoutMillis)
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Object.wait(JI)V"] = // wait(long timeoutMillis, int nanos)
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jlc := object.MakeEmptyObject()
	jlc.FieldTable = make(map[string]object.Field)
	name := object.GoStringFromStringPoolIndex(objPtr.KlassName)
	jlc.FieldTable["name"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: name,
	}

	if strings.HasPrefix(name, types.Array) { // arrays are handled differently
		arrClass := arrayGetClass(objPtr, name)
		return arrClass
	}

	// get a pointer to the class contents from the method area
	content := classloader.MethAreaFetch(name)
	if content == nil {
		errMsg := fmt.Sprintf("java/lang/Object.getClass: Class %s not loaded", name)
		return ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}

	// syntactic sugar
	obj := *content

	// create the empty java.lang.Class structure
	jlc.FieldTable["classLoader"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: obj.Loader,
	}

	// fill in the jlc
	objData := *obj.Data
	jlc.FieldTable["constantPool"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: objData.CP,
	}

	jlc.FieldTable["superClass"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: object.GoStringFromStringPoolIndex(objData.SuperclassIndex),
	}

	jlc.FieldTable["fields"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: objData.Fields,
	}

	jlc.FieldTable["interfaces"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: objData.Interfaces,
	}

	jlc.FieldTable["methods"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: objData.MethodTable,
	}

	jlc.FieldTable["methods"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: objData.MethodTable,
	}

	jlc.FieldTable["modifiers"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: objData.Access,
	}

	return jlc
}

// "java/lang/Object.toString()Ljava/lang/String;"
func objectToString(params []interface{}) interface{} {
	// params[0]: input Object

	switch params[0].(type) {
	case *object.Object:
		inObj := params[0].(*object.Object)
		classNameSuffix := object.GetClassNameSuffix(inObj, false)
		if classNameSuffix == "LinkedList" {
			return javaUtil.LinkedlistToString(params)
		}
		return object.StringifyAnythingJava(inObj)
	}

	errMsg := fmt.Sprintf("objectToString: Unsupported parameter type: %T", params[0])
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
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
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
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

// arrayGetClass creates a Class object for array types
// Arrays have special handling because they're not loaded from .class files
// Per JVM spec, all arrays have Object as their superclass
func arrayGetClass(objPtr *object.Object, arrayName string) *object.Object {
	jlc := object.MakeEmptyObject()
	jlc.FieldTable = make(map[string]object.Field)

	// Set the name field to the array type descriptor (e.g., "[Ljava/lang/String;" or "[I")
	jlc.FieldTable["name"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: arrayName,
	}

	// Determine the component type (the type of elements in the array)
	// For example: "[Ljava/lang/String;" -> "java/lang/String"
	//              "[I" -> "int"
	//              "[[I" -> "[I"
	componentType := ""
	if len(arrayName) > 1 {
		componentType = arrayName[1:] // Remove the leading '['

		// Convert internal format to readable format for object arrays
		// e.g., "Ljava/lang/String;" -> "java/lang/String"
		if strings.HasPrefix(componentType, "L") && strings.HasSuffix(componentType, ";") {
			componentType = componentType[1 : len(componentType)-1]
		}

		// Handle primitive types
		switch componentType {
		case "Z":
			componentType = "boolean"
		case "B":
			componentType = "byte"
		case "C":
			componentType = "char"
		case "D":
			componentType = "double"
		case "F":
			componentType = "float"
		case "I":
			componentType = "int"
		case "J":
			componentType = "long"
		case "S":
			componentType = "short"
		}
	}

	jlc.FieldTable["componentType"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: componentType,
	}

	// Arrays always have Object as their superclass
	jlc.FieldTable["superClass"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: "java/lang/Object",
	}

	// Arrays don't have fields (other than length, which is implicit)
	jlc.FieldTable["fields"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: []classloader.Field{},
	}

	// Arrays don't have methods
	jlc.FieldTable["methods"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: map[string]*classloader.Method{},
	}

	// Arrays don't have interfaces
	jlc.FieldTable["interfaces"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: []uint16{},
	}

	// Set modifiers - arrays are always public and final
	accessFlags := classloader.AccessFlags{
		ClassIsPublic: true,
		ClassIsFinal:  true,
	}
	jlc.FieldTable["modifiers"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: accessFlags,
	}

	// Arrays use the bootstrap classloader
	jlc.FieldTable["classLoader"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: "bootstrap",
	}

	return jlc
}
