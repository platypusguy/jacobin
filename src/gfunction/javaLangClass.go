/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/trace"
	"jacobin/types"
	"strings"
)

// Implementation of some of the functions in Java/lang/Class.

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
			GFunction:  trapFunction,
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
	x := statics.Statics["main.$assertionsDisabled"].Value.(int64)
	if x == 1 {
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

// Create a java/lang/Class instance -- that is a class ready for reflection
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
