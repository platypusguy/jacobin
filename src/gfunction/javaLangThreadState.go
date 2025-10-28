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
var threadStateNames = []string{
	"NEW",
	"RUNNABLE",
	"BLOCKED",
	"WAITING",
	"TIMED_WAITING",
	"TERMINATED",
}

// synchronization and lazy init of enum singletons
var threadStateOnce bool = false
var threadStateMutex = sync.Mutex{}
var threadStateClassName = "java/lang/Thread$State"
var threadStateInstances []*object.Object // length 6, matches threadStateNames
var threadStates map[string]*object.Object

func ensureThreadStateInited() {
	if threadStateOnce {
		return
	}
	threadStateInstances = make([]*object.Object, len(threadStateNames))
	for i, nm := range threadStateNames {
		obj := object.MakeEmptyObjectWithClassName(&threadStateClassName)
		// standard enum-like fields
		obj.FieldTable["value"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(nm)}
		threadStateInstances[i] = obj
		// expose as static so getstatic works
		_ = statics.AddStatic(threadStateClassName+"."+nm, statics.Static{Type: "Ljava/lang/Thread$State;", Value: obj})
	}
	threadStateOnce = true
}

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
}

// threadStateToString implements Thread.State.toString(): String
// Expected representations on the receiver object (in priority order):
// 1) Field "State" as a Go string with the state name.
// 2) Field "stateName" as a Go string with the state name.
// 3) Field "name" as a Go string (fallback if used by creator).
// 4) Field "name" as a Java String object.
// 5) Field "ordinal" (int) mapping to declaration-order names above.
// If none found, returns "UNKNOWN" to avoid throwing from toString().
func threadStateToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return getGErrBlk(excNames.IllegalArgumentException, "Thread$State.toString(): missing receiver")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return getGErrBlk(excNames.NullPointerException, "Thread$State.toString(): null receiver")
	}

	// 1) Go-string fields commonly used
	if f, ok := self.FieldTable["State"]; ok {
		if s, ok2 := f.Fvalue.(string); ok2 && s != "" {
			return object.StringObjectFromGoString(s)
		}
	}
	if f, ok := self.FieldTable["stateName"]; ok {
		if s, ok2 := f.Fvalue.(string); ok2 && s != "" {
			return object.StringObjectFromGoString(s)
		}
	}
	if f, ok := self.FieldTable["name"]; ok {
		// try Go string first
		if s, ok2 := f.Fvalue.(string); ok2 && s != "" {
			return object.StringObjectFromGoString(s)
		}
		// or a Java String object already
		if so, ok2 := f.Fvalue.(*object.Object); ok2 && so != nil && object.IsStringObject(so) {
			return so
		}
	}

	// 5) fallback: ordinal mapping
	if f, ok := self.FieldTable["ordinal"]; ok {
		switch vv := f.Fvalue.(type) {
		case int:
			if vv >= 0 && vv < len(threadStateNames) {
				return object.StringObjectFromGoString(threadStateNames[vv])
			}
		case int32:
			idx := int(vv)
			if idx >= 0 && idx < len(threadStateNames) {
				return object.StringObjectFromGoString(threadStateNames[idx])
			}
		case int64:
			idx := int(vv)
			if idx >= 0 && idx < len(threadStateNames) {
				return object.StringObjectFromGoString(threadStateNames[idx])
			}
		}
	}

	// Last resort: UNKNOWN
	return object.StringObjectFromGoString("UNKNOWN")
}

// threadStateValueOfString implements Thread.State.valueOf(String): returns the enum constant
func threadStateValueOf(params []interface{}) interface{} {
	ensureThreadStateInited()
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
