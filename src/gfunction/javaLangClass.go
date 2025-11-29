/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"strings"
)

// Implementation of some of the functions in Java/lang/Class. Note that a class
// implemented for reflection is referred to here a a Clazz. It is an object.

func Load_Lang_Class() {

	// There is no <clinit> for java/lang/Class.
	// The <clinit> type of code is executed in gfunction.go classClinitIsh().

	MethodSignatures["java/lang/Class.desiredAssertionStatus()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getAssertionsEnabledStatus,
		}

	MethodSignatures["java/lang/Class.desiredAssertionStatus0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getAssertionsEnabledStatus,
		}

	MethodSignatures["java/lang/Class.getComponentType()Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getComponentType,
		}

	MethodSignatures["java/lang/Class.getModule()Ljava/lang/Module;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  classGetModule,
		}

	MethodSignatures["java/lang/Class.getName()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getName,
		}

	MethodSignatures["java/lang/Class.getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  getPrimitiveClass,
		}

	MethodSignatures["java/lang/Class.isArray()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  classIsArray,
		}

	MethodSignatures["java/lang/Class.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

}

// getComponentType() returns a pointer to class of the type of an array.
// primitive arrays return the boxed class type, e.g. int[] returns Integer.class.
// multidimensional arrays return an array one dimension less, e.g. int[][] returns int[].class.
// Note: at present, this function returns a pointer to the loaded (but not instantiated) class.
func getComponentType(params []interface{}) interface{} {
	objPtr := params[0].(*object.Object)

	field := (*objPtr).FieldTable["value"]
	// If the object is not an array, return null.
	if !types.IsArray(field.Ftype) {
		return object.Null
	}

	componentType := field.Ftype[1:] // remove the leading '['
	if types.IsArray(componentType) {
		// If it's a multidimensional array, we return the pointer to the next dimension.
		return field.Fvalue
	}

	// If it's a primitive array, we return the boxed class type.
	if types.IsPrimitive(componentType) {
		// Convert the primitive type to its boxed class type.
		switch componentType {
		case types.Byte:
			componentType = "java/lang/Byte"
		case types.Char:
			componentType = "java/lang/Character"
		case types.Double:
			componentType = "java/lang/Double"
		case types.Float:
			componentType = "java/lang/Float"
		case types.Int:
			componentType = "java/lang/Integer"
		case types.Long:
			componentType = "java/lang/Long"
		case types.Short:
			componentType = "java/lang/Short"
		case types.Bool:
			componentType = "java/lang/Boolean"
		default:
			errMsg := fmt.Sprintf("getComponentType: unrecognized primitive type %s", componentType)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	} else {
		componentType = strings.TrimPrefix(componentType, types.Ref) // remove the leading 'L'
		componentType = strings.TrimSuffix(componentType, ";")       // remove the trailing ';'
	}

	// Load the class for the component type.
	cl, err := simpleClassLoadByName(componentType)
	if err != nil {
		errMsg := fmt.Sprintf("getComponentType: failed to load class %s: %s", componentType, err.Error())
		return getGErrBlk(excNames.ClassNotFoundException, errMsg)
	}

	return cl
}

// getPrimitiveClass() takes a one-word descriptor of a primitive and
// returns  apointer to the native primitive class that corresponds to it.
// This duplicates the behavior of OpenJDK JVMs.
// "java/lang/Class.getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;"
func getPrimitiveClass(params []interface{}) interface{} {
	primitive := params[0].(*object.Object)
	str := object.GoStringFromStringObject(primitive)

	var k *classloader.Klass
	var err error
	switch str {
	case "boolean":
		k, err = simpleClassLoadByName("java/lang/Boolean")
	case "byte":
		k, err = simpleClassLoadByName("java/lang/Byte")
	case "char":
		k, err = simpleClassLoadByName("java/lang/Character")
	case "double":
		k, err = simpleClassLoadByName("java/lang/Double")
	case "float":
		k, err = simpleClassLoadByName("java/lang/Float")
	case "int":
		k, err = simpleClassLoadByName("java/lang/Integer")
	case "long":
		k, err = simpleClassLoadByName("java/lang/Long")
	case "short":
		k, err = simpleClassLoadByName("java/lang/Short")
	case "void":
		k, err = simpleClassLoadByName("java/lang/Void")
	default:
		k = nil
		err = errors.New("getPrimitiveClass: unrecognized primitive")
	}

	if err == nil {
		return k
	} else {
		errMsg := fmt.Sprintf("getPrimitiveClass: %s: %s", err.Error(), str)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
}

// simpleClassLoadByName() just checks the MethodArea cache for the loaded
// class, and if it's not there, it loads it and returns a pointer to it.
// Logic basically duplicates similar functionality in instantiate.go
func simpleClassLoadByName(className string) (*classloader.Klass, error) {
	alreadyLoaded := classloader.MethAreaFetch(className)
	if alreadyLoaded != nil { // if the class is already loaded, skip the rest of this
		return alreadyLoaded, nil
	}

	// If not, try to load class by name
	err := classloader.LoadClassFromNameOnly(className)
	if err != nil {
		var errClassName = className
		if className == "" {
			errClassName = "<empty string>"
		}
		errMsg := fmt.Sprintf("simpleClassLoadByName: Failed to load class %s by name, reason: %s", errClassName, err.Error())
		trace.Error(errMsg)
		shutdown.Exit(shutdown.APP_EXCEPTION)
		return nil, errors.New(errMsg) // needed for testing, which does not cause an O/S exit on failure
	} else {
		return classloader.MethAreaFetch(className), nil
	}
}

// returns boolean indicating whether assertions are enabled or not.
// "java/lang/Class.desiredAssertionStatus()Z"
// "java/lang/Class.desiredAssertionStatus0()Z"
func getAssertionsEnabledStatus([]interface{}) interface{} {
	// note that statics have been preloaded before this function
	// can be called, and CLI processing has also occurred. So, we
	// know we have the latest assertion-enabled status.
	ste, ok := statics.QueryStatic("main", "$assertionsDisabled")
	if !ok {
		return types.JavaBoolFalse
	}
	if ste.Value.(int64) == int64(1) {
		return types.JavaBoolFalse
	} else {
		return types.JavaBoolTrue
	}
	// return 1 - x // return the 0 if disabled, 1 if not.
}

// "java/lang/Class.getName()Ljava/lang/String;"
func getName(params []interface{}) interface{} {
	primitive := params[0].(*object.Object)
	str := object.GoStringFromStringObject(primitive)
	return str
}

// "java/lang/Class.getName()Ljava/lang/String;"
func classIsArray(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fldType := obj.FieldTable["value"].Ftype
	if strings.HasPrefix(fldType, types.Array) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// classgetModule returns the unnamed module for any Class object
func classGetModule(params []interface{}) interface{} {
	if unnamedModule == nil {
		errMsg := "classGetModule: unnamed module not initialized"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}
	return unnamedModule
}

// Create a java/lang/Class instance -- that is a class ready for reflection from an existing object
func classCreateClassInstance(className string) (*object.Object, error) {
	cl, err := simpleClassLoadByName(className)
	if err != nil {
		return nil, err
	}
	if cl == nil {
		errMsg := fmt.Sprintf("classCreateClassInstance: failed to load class %s", className)
		return nil, errors.New(errMsg)
	}

	kl := object.MakeEmptyObject()
	if kl == nil {
		errMsg := fmt.Sprintf("classCreateClassInstance: failed to create new object of class %s", className)
		return nil, errors.New(errMsg)
	}

	return kl, nil
}

func classGetField(params []interface{}) interface{} {
	cl := params[0].(*object.Object)
	if object.IsNull(params[1]) {
		errMsg := "classGetField: null field name"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}
	fieldName := params[1].(string)
	_, ok := cl.FieldTable[fieldName]
	if !ok {
		errMsg := fmt.Sprintf("classGetField: field %s not found in %s",
			fieldName, *stringPool.GetStringPointer(cl.KlassName))
		return getGErrBlk(excNames.NoSuchFieldException, errMsg)
	}

	return NewField(cl, fieldName)
}
