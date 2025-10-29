/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"sync"
)

// Java: java.lang.Thread.State
// This gfunction file provides a minimal implementation focused on toString().
// States are represented as Go strings stored on the enum-like object.
// The toString method returns a Java String via object.StringObjectFromGoString().

// Declaration order matches Java's Thread.State

func Load_Lang_Thread_State() {
	MethodSignatures["java/lang/Thread$State.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadStateToString,
		}

	MethodSignatures["java/lang/Thread$State.valueOf(Ljava/lang/String;)Ljava/lang/Thread$State;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  threadStateValueOf,
		}

	MethodSignatures["java/lang/Thread$State.values()[Ljava/lang/Thread$State;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  threadStateValues,
		}

	ThreadStateInit()
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

// synchronization and lazy init of enum singletons
var threadStateOnce bool = false
var threadStateMutex = sync.Mutex{}
var threadStateClassName = "java/lang/Thread$State"
var threadStateInstances []*object.Object // length 6, matches threadStateNames
var threadStates map[string]*object.Object

// threadStateToString implements Thread.State.toString(): String
// It expects an int value for the state and returns the corresponding string object
func threadStateToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): missing object")
	}
	state, ok := params[0].(int)
	if !ok {
		return getGErrBlk(excNames.NullPointerException, "Thread$State.toString(): null object")
	}

	stateString, ok := ThreadState[state]
	if ok {
		return object.StringObjectFromGoString(stateString)
	} else {
		return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): invalid state")
	}
}

// threadStateValueOfString implements Thread.State.valueOf(String): returns the int64 constant
func threadStateValueOf(params []interface{}) interface{} {
	if len(params) < 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): missing argument")
	}
	strObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): argument is not a String")
	}
	if object.IsNull(strObj) {
		return getGErrBlk(excNames.NullPointerException, "Thread$State.valueOf(String): name is null")
	}
	if !object.IsStringObject(strObj) {
		return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): argument is not a String")
	}
	name := object.GoStringFromStringObject(strObj)
	for i, nm := range threadStateNames {
		if nm == name {
			return threadStateInstances[i]
		}
	}
	return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): no enum constant")
}

// threadStateValues implements Thread.State.values(): returns an array of all constants in decl order
func threadStateValues(params []interface{}) interface{} {
	ensureThreadStateInited()
	arr := object.Make1DimRefArray("Ljava/lang/Thread$State;", int64(len(threadStateInstances)))
	slot := arr.FieldTable["value"].Fvalue.([]*object.Object)
	copy(slot, threadStateInstances)
	arr.FieldTable["value"] = object.Field{Ftype: types.RefArray + "Ljava/lang/Thread$State;", Fvalue: slot}
	return arr
}

// threadStateClinit ensures constants are initialized when the class is initialized
func threadStateClinit([]interface{}) interface{} {
	ensureThreadStateInited()
	return nil
}
