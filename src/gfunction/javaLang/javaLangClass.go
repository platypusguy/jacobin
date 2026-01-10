/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
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
var unnamedModule = object.Null
var classNameModule = "java/lang/Module"

// ClassClinitIsh is handled special because the natural <clinit> function is never called.
/*
JVM spec:
"The java.lang.Class class is automatically initialized when the JVM is started.
However, because it is so tightly integrated with the JVM itself, its static initializer
s not necessarily run in the same way as other classes."
*/

// TODO: What is this?
// Called from <clinit> of java/lang/Class
func ClassClinitIsh() {
	// Initialize the unnamedModule singleton.
	if unnamedModule == nil {
		unnamedModule = &object.Object{
			KlassName: object.StringPoolIndexFromGoString(classNameModule),
			FieldTable: map[string]object.Field{
				"name": {
					Ftype:  types.StringClassRef,
					Fvalue: nil,
				},
				"isNamed": {
					Ftype:  types.Bool,
					Fvalue: types.JavaBoolFalse,
				},
				"value": {
					Ftype:  types.ModuleClassRef,
					Fvalue: nil,
				},
			},
		}
	}
}

func Load_Lang_Class() {

	// There is no <clinit> for java/lang/Class.
	// The <clinit> type of code is executed in gfunction.go classClinitIsh().

	ghelpers.MethodSignatures["java/lang/Class.desiredAssertionStatus()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classGetAssertionsEnabledStatus,
		}

	ghelpers.MethodSignatures["java/lang/Class.desiredAssertionStatus0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classGetAssertionsEnabledStatus,
		}

	ghelpers.MethodSignatures["java/lang/Class.getComponentType()Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  getComponentType,
		}

	ghelpers.MethodSignatures["java/lang/Class.getModule()Ljava/lang/Module;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classGetModule,
		}

	ghelpers.MethodSignatures["java/lang/Class.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classGetName,
		}

	ghelpers.MethodSignatures["java/lang/Class.getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  getPrimitiveClass,
		}

	ghelpers.MethodSignatures["java/lang/Class.isArray()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classIsArray,
		}

	ghelpers.MethodSignatures["java/lang/Class.isPrimitive()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classIsPrimitive,
		}

	ghelpers.MethodSignatures["java/lang/Class.registerNatives()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Class.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  classToString,
		}
}

// returns boolean indicating whether assertions are enabled or not.
// "java/lang/Class.desiredAssertionStatus()Z"
// "java/lang/Class.desiredAssertionStatus0()Z"
func classGetAssertionsEnabledStatus([]interface{}) interface{} {
	// note that statics have been preloaded before this function
	// can be called, and CLI processing has also occurred. So, we
	// know we have the latest assertion status.
	ste, ok := statics.QueryStatic("main", "$assertionsDisabled")
	if !ok {
		return types.JavaBoolFalse
	}
	if ste.Value.(int64) == int64(1) {
		return types.JavaBoolFalse
	} else {
		return types.JavaBoolTrue
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
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	} else {
		componentType = strings.TrimPrefix(componentType, types.Ref) // remove the leading 'L'
		componentType = strings.TrimSuffix(componentType, ";")       // remove the trailing ';'
	}

	// Load the class for the component type.
	cl, err := simpleClassLoadByName(componentType)
	if err != nil {
		errMsg := fmt.Sprintf("getComponentType: failed to load class %s: %s", componentType, err.Error())
		return ghelpers.GetGErrBlk(excNames.ClassNotFoundException, errMsg)
	}
	return cl
}

func classGetField(params []interface{}) interface{} {
	cl := params[0].(*object.Object)
	if object.IsNull(params[1]) {
		errMsg := "classGetField: null field name"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	fieldName := params[1].(string)
	_, ok := cl.FieldTable[fieldName]
	if !ok {
		errMsg := fmt.Sprintf("classGetField: field %s not found in %s",
			fieldName, *stringPool.GetStringPointer(cl.KlassName))
		return ghelpers.GetGErrBlk(excNames.NoSuchFieldException, errMsg)
	}

	return NewField(cl, fieldName)
}

// classgetModule returns the unnamed module for any Class object
func classGetModule([]interface{}) interface{} {
	if unnamedModule == nil {
		errMsg := "classGetModule: unnamed module not initialized"
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, errMsg)
	}
	return unnamedModule
}

// "java/lang/Class.classGetName()Ljava/lang/String;"
func classGetName(params []interface{}) interface{} {
	class := params[0].(*object.Object)
	name := class.FieldTable["name"].Fvalue.(string)
	return object.StringObjectFromGoString(name)
}

// getPrimitiveClass() takes a one-word descriptor of a primitive and
// returns a pointer to the native primitive class that corresponds to it.
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
}

// "java/lang/Class.isArray()Ljava/lang/String;"
func classIsArray(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fldType := obj.FieldTable["value"].Ftype
	if strings.HasPrefix(fldType, types.Array) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func classIsPrimitive(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fldType := obj.FieldTable["value"].Ftype
	if types.IsPrimitive(fldType) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

/* JDK Javadoc:
 * Converts the object to a string. The string representation is the
 * string "class" or "interface", followed by a space, and then by the
 * name of the class in the format returned by classGetName.
 * If this Class object represents a primitive type,  this method
 * returns the name of the primitive type. If this Class object
 * represents void this method returns "void". If this Class object
 * represents an array type, this method returns "class " followed
 * by classGetName.
 *
 * TODO: handle interfaces
 */
func classToString(params []any) any {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return object.StringObjectFromGoString("null")
	}

	name := obj.FieldTable["name"].Fvalue.(string)
	str := fmt.Sprintf("class %s", name)
	return object.StringObjectFromGoString(str)
}

// === helper functions (not part of the javaLangClass class API) ===

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
