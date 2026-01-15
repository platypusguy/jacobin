/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"sort"
)

// Java: java.lang.Thread.State
// This gfunction file provides a minimal implementation focused on toString().
// States are represented as Go strings stored on the enum-like object.
// The toString method returns a Java String via object.StringObjectFromGoString().

// Declaration order matches Java's Thread.State

func Load_Lang_Thread_State() {
	ghelpers.MethodSignatures["java/lang/Thread$State.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Thread$State.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  threadStateToString,
		}

	ghelpers.MethodSignatures["java/lang/Thread$State.valueOf(Ljava/lang/String;)Ljava/lang/Thread$State;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  threadStateValueOf,
		}

	ghelpers.MethodSignatures["java/lang/Thread$State.values()[Ljava/lang/Thread$State;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  threadStateValues,
		}
}

const (
	NEW           = 0
	RUNNABLE      = 1
	BLOCKED       = 2
	WAITING       = 3
	TIMED_WAITING = 4
	TERMINATED    = 5
)

var ThreadState = map[int]string{
	NEW:           "NEW",
	RUNNABLE:      "RUNNABLE",
	BLOCKED:       "BLOCKED",
	WAITING:       "WAITING",
	TIMED_WAITING: "TIMED_WAITING",
	TERMINATED:    "TERMINATED",
}

// threadStateValueOf implements Thread.State.valueOf(String): returns the associated Thread.State object.
func threadStateValueOf(params []interface{}) interface{} {
	if len(params) != 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): missing argument")
	}
	strObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): argument is not a String")
	}
	if object.IsNull(strObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Thread$State.valueOf(String): name is null")
	}
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): argument is not a String")
	}
	name := object.GoStringFromStringObject(strObj)
	for key, value := range ThreadState {
		if value == name {
			// Success!
			obj := object.MakePrimitiveObject("java/lang/Thread$State", types.Int, key)
			return obj
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): no value found for "+name)
}

// threadStateValues implements Thread.State.values(): returns an array of all constants in decl order
func threadStateValues([]interface{}) interface{} {
	// Extract keys into a slice
	keys := make([]int, 0, len(ThreadState))
	for key := range ThreadState {
		keys = append(keys, key)
	}

	// Sort the slice
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// Create the output array.
	var objArray []*object.Object
	for key := range keys {
		obj := object.MakePrimitiveObject("java/lang/Thread$State", types.Int, key)
		objArray = append(objArray, obj)
	}
	objObjArray := object.MakePrimitiveObject("java/lang/Thread$State", types.RefArray+"java/lang/Thread$State", objArray)

	return objObjArray
}

// threadStateToString implements Thread.State.toString(): String
// It expects an int value for the state and returns the corresponding string object
func threadStateToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): missing argument")
	}
	thStateObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): not an object")
	}
	state, ok := thStateObj.FieldTable["value"].Fvalue.(int)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): missing/invalid field value")
	}
	if state > TERMINATED || state < NEW {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("Thread$State.toString(): invalid state: %d", state))
	}

	stateString, ok := ThreadState[state]
	if ok {
		return object.StringObjectFromGoString(stateString)
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): invalid state")
}
