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
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
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
func ClassClinitIsh() {
	// Initialize the unnamedModule singleton.
	if unnamedModule == nil {
		unnamedModule = object.MakeEmptyObject()
		unnamedModule.KlassName = object.StringPoolIndexFromGoString(classNameModule)
		unnamedModule.ThMutex.Lock()
		unnamedModule.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: nil,
		}
		unnamedModule.FieldTable["isNamed"] = object.Field{
			Ftype:  types.Bool,
			Fvalue: types.JavaBoolFalse,
		}
		unnamedModule.FieldTable["value"] = object.Field{
			Ftype:  types.ModuleClassRef,
			Fvalue: nil,
		}
		unnamedModule.ThMutex.Unlock()
	}
}

func Load_Lang_Class() {

	// Implemented methods
	ghelpers.MethodSignatures["java/lang/Class.desiredAssertionStatus()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classGetAssertionsEnabledStatus}
	ghelpers.MethodSignatures["java/lang/Class.desiredAssertionStatus0()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classGetAssertionsEnabledStatus}
	ghelpers.MethodSignatures["java/lang/Class.getCanonicalName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classGetCanonicalName}
	ghelpers.MethodSignatures["java/lang/Class.getComponentType()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: getComponentType}
	ghelpers.MethodSignatures["java/lang/Class.getModule()Ljava/lang/Module;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classGetModule}
	ghelpers.MethodSignatures["java/lang/Class.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ClassGetName}
	ghelpers.MethodSignatures["java/lang/Class.getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: getPrimitiveClass}
	ghelpers.MethodSignatures["java/lang/Class.getSimpleName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ClassGetSimpleName}
	ghelpers.MethodSignatures["java/lang/Class.getSuperclass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classGetSuperclass}
	ghelpers.MethodSignatures["java/lang/Class.isArray()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classIsArray}
	ghelpers.MethodSignatures["java/lang/Class.isInterface()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classIsInterface}
	ghelpers.MethodSignatures["java/lang/Class.isPrimitive()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classIsPrimitive}
	ghelpers.MethodSignatures["java/lang/Class.registerNatives()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ClinitGeneric}
	ghelpers.MethodSignatures["java/lang/Class.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: classToString}

	// --- Alphabetized trapped methods ---
	addTrap := func(signature string, slots int) {
		ghelpers.MethodSignatures[signature] =
			ghelpers.GMeth{ParamSlots: slots, GFunction: ghelpers.TrapFunction}
	}

	addTrap("java/lang/Class.accessFlags()Ljava/util/Set;", 0)
	addTrap("java/lang/Class.asSubclass(Ljava/lang/Class;)Ljava/lang/Class;", 1)
	addTrap("java/lang/Class.cast(Ljava/lang/Object;)Ljava/lang/Object;", 1)
	addTrap("java/lang/Class.componentType()Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.describeConstable()Ljava/util/Optional;", 0)
	addTrap("java/lang/Class.descriptorString()Ljava/lang/String;", 0)
	addTrap("java/lang/Class.forName(Ljava/lang/Module;Ljava/lang/String;)Ljava/lang/Class;", 2)
	addTrap("java/lang/Class.forName(Ljava/lang/String;)Ljava/lang/Class;", 1)
	addTrap("java/lang/Class.forName(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", 3)
	addTrap("java/lang/Class.getAnnotatedInterfaces()[Ljava/lang/reflect/AnnotatedType;", 0)
	addTrap("java/lang/Class.getAnnotatedSuperclass()Ljava/lang/reflect/AnnotatedType;", 0)
	addTrap("java/lang/Class.getAnnotation(Ljava/lang/Class;)Ljava/lang/annotation/Annotation;", 1)
	addTrap("java/lang/Class.getAnnotations()[Ljava/lang/annotation/Annotation;", 0)
	addTrap("java/lang/Class.getConstructor([Ljava/lang/Class;)Ljava/lang/reflect/Constructor;", 1)
	addTrap("java/lang/Class.getConstructors()[Ljava/lang/reflect/Constructor;", 0)
	addTrap("java/lang/Class.getDeclaredAnnotationsByType(Ljava/lang/Class;)[Ljava/lang/annotation/Annotation;", 1)
	addTrap("java/lang/Class.getDeclaredClasses()[Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getDeclaredConstructor([Ljava/lang/Class;)Ljava/lang/reflect/Constructor;", 1)
	addTrap("java/lang/Class.getDeclaredConstructors()[Ljava/lang/reflect/Constructor;", 0)
	addTrap("java/lang/Class.getDeclaredField(Ljava/lang/String;)Ljava/lang/reflect/Field;", 1)
	addTrap("java/lang/Class.getDeclaredFields()[Ljava/lang/reflect/Field;", 0)
	addTrap("java/lang/Class.getDeclaredMethod(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", 2)
	addTrap("java/lang/Class.getDeclaredMethods()[Ljava/lang/reflect/Method;", 0)
	addTrap("java/lang/Class.getDeclaringClass()Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getEnclosingClass()Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getEnclosingConstructor()Ljava/lang/reflect/Constructor;", 0)
	addTrap("java/lang/Class.getEnclosingMethod()Ljava/lang/reflect/Method;", 0)
	addTrap("java/lang/Class.getEnumConstants()[Ljava/lang/Object;", 0)
	addTrap("java/lang/Class.getField(Ljava/lang/String;)Ljava/lang/reflect/Field;", 1)
	addTrap("java/lang/Class.getFields()[Ljava/lang/reflect/Field;", 0)
	addTrap("java/lang/Class.getGenericInterfaces()[Ljava/lang/reflect/Type;", 0)
	addTrap("java/lang/Class.getGenericSuperclass()Ljava/lang/reflect/Type;", 0)
	addTrap("java/lang/Class.getInterfaces()[Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getModifiers()I", 0)
	addTrap("java/lang/Class.getNestHost()Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getNestMembers()[Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getPackage()Ljava/lang/Package;", 0)
	addTrap("java/lang/Class.getPackageName()Ljava/lang/String;", 0)
	addTrap("java/lang/Class.getPermittedSubclasses()[Ljava/lang/Class;", 0)
	addTrap("java/lang/Class.getProtectionDomain()Ljava/security/ProtectionDomain;", 0)
	addTrap("java/lang/Class.getRecordComponents()[Ljava/lang/reflect/RecordComponent;", 0)
	// addTrap("java/lang/Class.getSimpleName()Ljava/lang/String;", 0)
	addTrap("java/lang/Class.getTypeName()Ljava/lang/String;", 0)
	addTrap("java/lang/Class.isEnum()Z", 0)
	addTrap("java/lang/Class.isInstance(Ljava/lang/Object;)Z", 1)
	addTrap("java/lang/Class.isNestmateOf(Ljava/lang/Class;)Z", 1)
	addTrap("java/lang/Class.isSealed()Z", 0)
	addTrap("java/lang/Class.isSynthetic()Z", 0)
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

// to distinguish name from canonical name, see comments before classGetName() below for diiferences between the
// different names that java/lang/CLass can return
// java/lang/Class.getCanonicalName()Ljava/lang/String
func classGetCanonicalName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "classGetCanonicalName: invalid or null object")
	}

	rawName := object.GoStringFromStringPoolIndex(obj.KlassName)
	if strings.HasPrefix(rawName, types.Array) {
		switch rawName[1] {
		case 'B':
			return object.StringObjectFromGoString("byte[]")
		case 'C':
			return object.StringObjectFromGoString("char[]")
		case 'D':
			return object.StringObjectFromGoString("double[]")
		case 'F':
			return object.StringObjectFromGoString("float[]")
		case 'I':
			return object.StringObjectFromGoString("int[]")
		case 'J':
			return object.StringObjectFromGoString("long[]")
		case 'S':
			return object.StringObjectFromGoString("short[]")
		case 'Z':
			return object.StringObjectFromGoString("boolean[]")
		case 'L':
			arrayObject := rawName[2 : len(rawName)-1] // remove the leading '[L' and trailing ';'
			// Check if the component type is anonymous or local class
			if strings.Contains(arrayObject, "$") {
				lastDollar := strings.LastIndex(arrayObject, "$")
				afterDollar := arrayObject[lastDollar+1:]
				// Anonymous class (e.g., Outer$1) or local class (e.g., Outer$1Local) - return nil
				if len(afterDollar) > 0 && afterDollar[0] >= '0' && afterDollar[0] <= '9' {
					return nil
				}
				// Inner class - replace $ with .
				arrayObject = strings.ReplaceAll(arrayObject, "$", ".")
			}
			return object.StringObjectFromGoString(arrayObject + "[]") // so java.lang.String[], for example
		}
	}

	// Handle anonymous and local classes - return nil
	if strings.Contains(rawName, "$") {
		lastDollar := strings.LastIndex(rawName, "$")
		afterDollar := rawName[lastDollar+1:]
		// Check if it's anonymous or local class (starts with digit after $)
		if len(afterDollar) > 0 && afterDollar[0] >= '0' && afterDollar[0] <= '9' {
			return nil
		}
		// Inner class - replace $ with . for canonical name
		canonicalName := strings.ReplaceAll(rawName, "$", ".")
		return object.StringObjectFromGoString(canonicalName)
	}

	return object.StringObjectFromGoString(rawName)
}

// getComponentType() returns a pointer to class of the type of an array.
// primitive arrays return the boxed class type, e.g. int[] returns Integer.class.
// multidimensional arrays return an array one dimension less, e.g. int[][] returns int[].class.
// Note: at present, this function returns a pointer to the loaded (but not instantiated) class.
func getComponentType(params []interface{}) interface{} {
	objPtr, ok := params[0].(*object.Object)
	if !ok || object.IsNull(objPtr) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "getComponentType: invalid or null object")
	}

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

// TODO: Is this function needed? Its currently trapped.
func classGetField(params []interface{}) interface{} {
	cl, ok := params[0].(*object.Object)
	if !ok || object.IsNull(cl) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "classGetField: invalid or null object")
	}
	if object.IsNull(params[1]) {
		errMsg := "classGetField: null field name"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	fieldName := params[1].(string)
	_, ok = cl.FieldTable[fieldName]
	if !ok {
		errMsg := fmt.Sprintf("classGetField: field %s not found in %s",
			fieldName, *stringPool.GetStringPointer(cl.KlassName))
		return ghelpers.GetGErrBlk(excNames.NoSuchFieldException, errMsg)
	}

	return NewField(cl, fieldName)
}

// classGetModule returns the unnamed module for any Class object
func classGetModule([]interface{}) interface{} {
	if unnamedModule == nil {
		errMsg := "classGetModule: unnamed module not initialized"
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, errMsg)
	}
	return unnamedModule
}

// The various names are defined differently:
// Class type         getName()                getCanonicalName()
// ----------------   ---------------------   ---------------------
// Standard Class:    java.lang.String,        java.lang.String
// Inner Class:       com.example.Outer$Inner  com.example.Outer.Inner
// String Array:      [Ljava.lang.String;      java.lang.String[]
// Primitive Array:   [I                       int[]
// Anonymous Class:   com.example.Outer$1      null (Anonymous classes have no canonical name)
// Local Class:       com.example.Outer$1Local null
//
// "java/lang/Class.classGetName()Ljava/lang/String;"
func ClassGetName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "java/lang/class.GetName(): invalid or null object")
	}
	name := obj.FieldTable["name"].Fvalue.(string)
	name = util.ConvertInternalClassNameToUserFormat(name)
	nameObj := object.StringObjectFromGoString(name)
	return nameObj
}

// getPrimitiveClass() takes a one-word descriptor of a primitive and
// returns a pointer to the native primitive class that corresponds to it.
// This duplicates the behavior of OpenJDK JVMs.
// "java/lang/Class.getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;"
func getPrimitiveClass(params []interface{}) interface{} {
	primitive, ok := params[0].(*object.Object)
	if !ok || object.IsNull(primitive) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "getPrimitiveClass: invalid or null object")
	}
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

// | Class Type      | getName()               | getSimpleName()
// |-----------------|-------------------------|----------------
// | Regular class   | java.lang.String        | String
// | Inner class     | com.example.Outer$Inner | Inner
// | Anonymous class | com.example.MyClass$1   | (empty string)
// | Array           | [Ljava.lang.String;     | String[]
// | Primitive array | [I                      | int[]
//
// > getName() returns the fully qualified name with package,
// and it uses internal JVM notation for arrays (e.g., `[L...;`)
//
// > getSimpleName() returns just the class name as written in source code,
// and it returns empty string for anonymous classes and it uses
// readable array notation (e.g., `[]`)
//
// java/lang/Class.getSimpleName()Ljava/lang/String;
func ClassGetSimpleName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ClassGetSimpleName: invalid or null object")
	}

	name := obj.FieldTable["name"].Fvalue.(string)

	// Handle arrays - use canonical name logic but return simple name for component
	if strings.HasPrefix(name, types.Array) {
		switch name[1] {
		case 'B':
			return object.StringObjectFromGoString("byte[]")
		case 'C':
			return object.StringObjectFromGoString("char[]")
		case 'D':
			return object.StringObjectFromGoString("double[]")
		case 'F':
			return object.StringObjectFromGoString("float[]")
		case 'I':
			return object.StringObjectFromGoString("int[]")
		case 'J':
			return object.StringObjectFromGoString("long[]")
		case 'S':
			return object.StringObjectFromGoString("short[]")
		case 'Z':
			return object.StringObjectFromGoString("boolean[]")
		case 'L':
			// For object arrays, extract the simple name of the component type
			arrayObject := name[2 : len(name)-1] // remove '[L' and ';'
			lastSlash := strings.LastIndex(arrayObject, "/")
			if lastSlash >= 0 {
				arrayObject = arrayObject[lastSlash+1:] // get simple name after last /
			}
			// Handle inner/anonymous classes in array component type
			lastDollar := strings.LastIndex(arrayObject, "$")
			if lastDollar >= 0 {
				// Check if it's anonymous (ends with digits after $)
				afterDollar := arrayObject[lastDollar+1:]
				if isNumeric(afterDollar) {
					return object.StringObjectFromGoString("[]") // anonymous class array
				}
				arrayObject = arrayObject[lastDollar+1:] // inner class name
			}
			return object.StringObjectFromGoString(arrayObject + "[]")
		}
	}

	// Handle anonymous classes - they contain $ followed by digits
	if strings.Contains(name, "$") {
		lastDollar := strings.LastIndex(name, "$")
		afterDollar := name[lastDollar+1:]
		if isNumeric(afterDollar) {
			return object.StringObjectFromGoString("") // anonymous class returns empty string
		}
		// Inner class - return the name after the last $
		simpleName := afterDollar
		return object.StringObjectFromGoString(simpleName)
	}

	// Regular class - return the name after the last /
	lastSlash := strings.LastIndex(name, "/")
	if lastSlash >= 0 {
		name = name[lastSlash+1:]
	}

	return object.StringObjectFromGoString(name)
}

// Helper function to check if a string contains only digits
func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// java/lang/Class.getSuperclass()Ljava/lang/Class; return a java/lang/Class object representing the superclass.
// java/lang/Object, primitives, and interfaces return null.
// Consult https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Class.html#getSuperclass()
func classGetSuperclass(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "classGetSuperclass: invalid or null object")
	}

	// if the object is an array, return Object.class
	if classIsArray(params).(bool) {
		return globals.JLCmap["java/lang/Object"]
	}

	// if the object is an interface, the superclass is null
	if classIsInterface(params).(bool) {
		return nil
	}

	// java/lang/Object, primitives, and void all return null for the superclass
	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	switch className {
	case "java/lang/Object",
		"java/lang/Byte",
		"java/lang/Character",
		"java/lang/Double",
		"java/lang/Float",
		"java/lang/Integer",
		"java/lang/Long",
		"java/lang/Short",
		"java/lang/Boolean",
		"java/lang/Void":
		return nil
	}

	klassPtr := obj.FieldTable["$klass"].Fvalue.(*classloader.ClData)
	scNameIndex := klassPtr.SuperclassIndex
	scName := *stringPool.GetStringPointer(scNameIndex)

	scClass := globals.JLCmap[scName]
	return scClass
}

// java/lang/Class.isArray()Z
func classIsArray(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "classIsArray: invalid or null object")
	}
	fldType := obj.FieldTable["value"].Ftype
	if strings.HasPrefix(fldType, types.Array) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// java/lang/Class.isInterface()Z
func classIsInterface(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "classIsPrimitive: invalid or null object")
	}

	klass := obj.FieldTable["$klass"].Fvalue.(*classloader.ClData)
	return klass.Access.ClassIsInterface
}

// java/lang/Class.isPrimitive()Z
func classIsPrimitive(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "classIsPrimitive: invalid or null object")
	}
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
