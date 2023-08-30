/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/shutdown"
)

// Implementation of some of the functions in in Java/lang/Class.

func Load_Lang_Class() map[string]GMeth {

	MethodSignatures["java/lang/Class.getPrimitiveClass(Ljava/lang/String;)Ljava/lang/Class;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  getPrimitiveClass,
		}

	MethodSignatures["java/lang/Class.desiredAssertionStatus()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getAssertionsEnabledStatus,
		}

	MethodSignatures["java/lang/Class.desiredAssertionStatus0()Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  getAssertionsEnabledStatus0,
		}
	return MethodSignatures
}

// getPrimitiveClass() takes a one-word descriptor of a primitive and
// returns the native primitive class that corresponds to it. Not
// quite sure how this class is used, but it needs to implemented.
func getPrimitiveClass(params []interface{}) interface{} {
	primitive := params[0].(*object.Object)
	str := object.GetGoStringFromJavaStringPtr(primitive)
	if str == "int" {
		k, err := simpleClassLoadByName("java/lang/Integer")
		if err != nil {
			return errors.New("getPrimitiveClass() could not load java/lang/Integer")
		} else {
			return k
		}
	} else {
		errMsg := fmt.Sprintf("getPrimitiveClass() does not handle: %s", str)
		_ = log.Log(errMsg, log.SEVERE)
		return errors.New(errMsg)

	}
}

// simpleClassLoadByName() just checks the MethodArea cache for the loaded
// class, and if it's not there, it loads it and returns a pointer to it.
// Logic basically duplicates similar functionality in instantiate.go
func simpleClassLoadByName(className string) (*Klass, error) {
	alreadyLoaded := MethAreaFetch(className)
	if alreadyLoaded != nil { // if the class is already loaded, skip the rest of this
		return alreadyLoaded, nil
	}

	// If not, try to load class by name
	err := LoadClassFromNameOnly(className)
	if err != nil {
		errMsg := "instantiateClass()-getPrimitivelass(): Failed to load class " + className
		_ = log.Log(errMsg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		shutdown.Exit(shutdown.APP_EXCEPTION)
		return nil, errors.New(errMsg) // needed for testing, which does not shutdown on failure
	} else {
		return MethAreaFetch(className), nil
	}
}

// returns boolean indicating whether assertions are enabled or not.
func getAssertionsEnabledStatus(params []interface{}) interface{} {
	g := globals.GetGlobalRef()
	return g.AssertionsEnabled
}

// returns boolean indicating whether assertions are enabled or not.
// Effectively identical to getAsserionsEnabledStatus()
func getAssertionsEnabledStatus0(params []interface{}) interface{} {
	g := globals.GetGlobalRef()
	return g.AssertionsEnabled
}
