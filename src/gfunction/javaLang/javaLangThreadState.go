/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
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
	ghelpers.MethodSignatures["java/lang/Thread$State.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  threadStateCreateWithValue,
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
	UNDEFINED     = 6
)

var ThreadState = map[int]string{
	NEW:           "NEW",
	RUNNABLE:      "RUNNABLE",
	BLOCKED:       "BLOCKED",
	WAITING:       "WAITING",
	TIMED_WAITING: "TIMED_WAITING",
	TERMINATED:    "TERMINATED",
	UNDEFINED:     "UNDEFINED",
}

// synchronization and lazy init of enum singletons

var threadStateInstances []*object.Object // length 6, matches threadStateNames

// creates a thread state, but is not actually part of the OpenJDK API. Used internally by Thread class
func threadStateCreateWithValue(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.<init>(int): missing value")
	}
	state, ok := params[0].(int)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.<init>(int): invalid value")
	}
	if state < NEW || state > TERMINATED {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.<init>(int): invalid value")
	}

	ts := object.MakeEmptyObject()
	ts.KlassName = object.StringPoolIndexFromGoString("java/lang/Thread$State")
	ts.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: state}

	return ts
}

// threadStateToString implements Thread.State.toString(): String
// It expects an int value for the state and returns the corresponding string object
func threadStateToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): missing argument")
	}
	obj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): not an object")
	}
	state, ok := obj.FieldTable["value"].Fvalue.(int)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Thread$State.toString(): invalid field value")
	}

	stateString, ok := ThreadState[state]
	if ok {
		return object.StringObjectFromGoString(stateString)
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): invalid state")
}

// threadStateValueOfString implements Thread.State.valueOf(String): returns the int64 constant
func threadStateValueOf(params []interface{}) interface{} {
	if len(params) < 1 {
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
			return key // found it
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Thread$State.valueOf(String): no value found for "+name)
}

// threadStateValues implements Thread.State.values(): returns an array of all constants in decl order
func threadStateValues([]interface{}) interface{} {
	arr := object.Make1DimRefArray("Ljava/lang/Thread$State;", int64(len(threadStateInstances)))

	// Extract keys into a slice
	keys := make([]int, 0, len(ThreadState))
	for key := range ThreadState {
		keys = append(keys, key)
	}

	// Sort the slice
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	arr.FieldTable["value"] = object.Field{Ftype: types.RefArray + "Ljava/lang/Thread$State;", Fvalue: keys}
	return arr
}
