/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/thread"
	"jacobin/src/types"
)

func Load_Lang_Thread_Group() {
	// Initialize the initial global thread groups
	initializeGlobalThreadGroups()

	// <clinit>
	MethodSignatures["java/lang/ThreadGroup.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: threadGroupClinit}

	// Constructors
	MethodSignatures["java/lang/ThreadGroup.<init>(Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 1, GFunction: threadGroupInitWithName}
	MethodSignatures["java/lang/ThreadGroup.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 2, GFunction: threadGroupCreateWithParentAndName}

	// Public instance methods (alphabetical by JVM signature for consistency)
	MethodSignatures["java/lang/ThreadGroup.activeCount()I"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.activeGroupCount()I"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.allowThreadSuspension(Z)Z"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.checkAccess()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.destroy()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.enumerate([Ljava/lang/Thread;)I"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.enumerate([Ljava/lang/Thread;Z)I"] =
		GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.enumerate([Ljava/lang/ThreadGroup;)I"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.enumerate([Ljava/lang/ThreadGroup;Z)I"] =
		GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.getMaxPriority()I"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.getName()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: threadGroupGetName}
	MethodSignatures["java/lang/ThreadGroup.getParent()Ljava/lang/ThreadGroup;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.interrupt()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.isDaemon()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.isDestroyed()Z"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.list()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.parentOf(Ljava/lang/ThreadGroup;)Z"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.setDaemon(Z)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.setMaxPriority(I)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.stop()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.suspend()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.resume()V"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.toString()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.uncaughtException(Ljava/lang/Thread;Ljava/lang/Throwable;)V"] =
		GMeth{ParamSlots: 2, GFunction: trapFunction}
}

// Initialize global thread groups: create "system" and its child "main"
func initializeGlobalThreadGroups() {
	gr := globals.GetGlobalRef()
	if gr.ThreadGroups == nil {
		gr.ThreadGroups = make(map[string]interface{})
	}

	// Create "system" group
	sys := threadGroupInitWithName([]interface{}{object.StringObjectFromGoString("system")})
	gr.ThreadGroups["system"] = sys

	// Create "main" group as a child of "system"
	sysObj, _ := gr.ThreadGroups["system"].(*object.Object)
	mainGrp := threadGroupCreateWithParentAndName([]interface{}{sysObj,
		object.StringObjectFromGoString("main")})
	gr.ThreadGroups["main"] = mainGrp
}

// java/lang/ThreadGroup.<clinit>()V
func threadGroupClinit(_ []interface{}) any {
	return justReturn(nil)
}

func ThreadGroupInitWithParentNameMaxpriorityDaemon(initParams []interface{}) any {
	if len(initParams) != 5 { // the four named params + the ThreadGroup object itself
		return getGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("ThreadGroupInitWithParentNameMaxpriorityDaemon: Expected 5 parameters, got %d",
				len(initParams)))
	}

	obj, ok := initParams[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"ThreadGroupInitWithParentNameMaxpriorityDaemon: Expected first parameter to be an object reference")
	}

	parentObj, ok := initParams[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"ThreadGroupInitWithParentNameMaxpriorityDaemon: Expected second parameter to be an object reference")
	}

	nameObj, ok := initParams[2].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"ThreadGroupInitWithParentNameMaxpriorityDaemon: Expected third parameter to be a String object")
	}

	maxPriority := initParams[3].(int64)
	daemon := initParams[4].(int64)

	if parentObj != object.Null {
		obj.FieldTable["parent"] = object.Field{Ftype: types.Ref, Fvalue: parentObj}
	}
	if nameObj != object.Null {
		obj.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: nameObj}
	}

	if maxPriority != 0 { // 0 = uninitialized
		if maxPriority < thread.MIN_PRIORITY || maxPriority > thread.MAX_PRIORITY {
			return getGErrBlk(excNames.IllegalArgumentException,
				"ThreadGroupInitWithParentNameMaxpriorityDaemon: maxPriority out of range")
		}
		obj.FieldTable["maxpriority"] = object.Field{Ftype: types.Int, Fvalue: maxPriority}
	}

	if daemon == types.JavaBoolFalse || daemon == types.JavaBoolTrue {
		obj.FieldTable["daemon"] = object.Field{Ftype: types.Bool, Fvalue: daemon}
	}

	obj.FieldTable["parent"] = object.Field{Ftype: types.Ref, Fvalue: parentObj}

	// initialize the fields that are not passed as parameters
	obj.FieldTable["priority"] =
		object.Field{Ftype: types.Int, Fvalue: int64(thread.NORM_PRIORITY)}

	subgroups := object.Field{Ftype: types.LinkedList, Fvalue: list.New()}
	obj.FieldTable["subgroups"] = subgroups

	// add the thread group to the global list of thread groups
	globals.GetGlobalRef().ThreadGroups[object.GoStringFromStringObject(nameObj)] = obj

	return obj
}

// java/lang/ThreadGroup.ThreadGroup(Ljava/lang/String;)Ljava/lang/ThreadGroup;
// accepts a  ThreadGroup *object and adds the specified name and a null parent
func threadGroupInitWithName(params []interface{}) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("threadGroupInitWithName: Expected 2 parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGroupInitWithName: Expected 1st parameter to be an object reference, was not"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "threadGroupCreateWithName: Expected parameter to be a string object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// clName := "java/lang/ThreadGroup"
	// obj := object.MakeEmptyObjectWithClassName(&clName)

	nullField := object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable["parent"] = nullField

	nameField := object.Field{Ftype: types.Ref, Fvalue: name}
	obj.FieldTable["name"] = nameField

	daemonField := object.Field{Ftype: types.Int, Fvalue: types.JavaBoolFalse}
	obj.FieldTable["daemon"] = daemonField

	threadGroup := object.Field{Ftype: types.Ref, Fvalue: nil}
	obj.FieldTable["threadgroup"] = threadGroup

	priority := object.Field{Ftype: types.Int, Fvalue: int64(thread.NORM_PRIORITY)}
	obj.FieldTable["priority"] = priority

	maxPriority := object.Field{Ftype: types.Int, Fvalue: int64(thread.MAX_PRIORITY)}
	obj.FieldTable["maxpriority"] = maxPriority

	subgroups := object.Field{Ftype: types.LinkedList, Fvalue: list.New()}
	obj.FieldTable["subgroups"] = subgroups

	// add the thread group to the global list of thread groups
	globals.GetGlobalRef().ThreadGroups[object.GoStringFromStringObject(name)] = obj

	return obj
}

// java/lang/ThreadGroup.ThreadGroup(Ljava/lang/ThreadGroup;Ljava/lang/String;)Ljava/lang/ThreadGroup;
// returns a new ThreadGroup object with the specified name and parent
func threadGroupCreateWithParentAndName(params []interface{}) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf(
			"threadGroupCreateWithParentAndName: Expected 2 parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// First param: parent ThreadGroup (object reference)
	parentObj, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"threadGroupCreateWithParentAndName: Expected first parameter to be a ThreadGroup object")
	}

	// Second param: name (Java String object)
	nameObj, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"threadGroupCreateWithParentAndName: Expected second parameter to be a String object")
	}
	if object.IsNull(nameObj) {
		return getGErrBlk(excNames.NullPointerException,
			"threadGroupCreateWithParentAndName: name is null")
	}
	if !object.IsStringObject(nameObj) {
		return getGErrBlk(excNames.IllegalArgumentException,
			"threadGroupCreateWithParentAndName: second parameter is not a String")
	}

	// Create group with name
	created := threadGroupInitWithName([]interface{}{nameObj})

	// If creation returned an error block, pass it through
	if gerr, isErr := created.(*GErrBlk); isErr {
		return gerr
	}

	// Otherwise, set parent field and return the object
	tg, ok := created.(*object.Object)
	if !ok {
		// Shouldn’t happen, but fail gracefully
		return getGErrBlk(excNames.IllegalArgumentException,
			"threadGroupCreateWithParentAndName: factory returned non-object")
	}

	tg.FieldTable["parent"] = object.Field{Ftype: types.Ref, Fvalue: parentObj}

	// Now add this thread group to the parent's list of subgroups
	parentSubgroups := parentObj.FieldTable["subgroups"].Fvalue.(*list.List)
	parentSubgroups.PushBack(tg)

	return tg
}

// java/lang/ThreadGroup.getName()Ljava/lang/String;
func threadGroupGetName(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGroupGetName: Expected thread group, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	tg, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"threadGroupGetName: Expected parameter to be an object reference")
	}

	f := tg.FieldTable["name"]
	// If stored as Java String object, just return it
	if obj, ok := f.Fvalue.(*object.Object); ok && object.IsStringObject(obj) {
		return obj
	}

	// Fallback in case legacy code stored Go string
	if s, ok := f.Fvalue.(string); ok {
		return object.StringObjectFromGoString(s)
	}

	// Fallback in case legacy code stored Go string
	if s, ok := f.Fvalue.([]types.JavaByte); ok {
		return object.StringObjectFromJavaByteArray(s)
	}

	return getGErrBlk(excNames.IllegalArgumentException,
		"threadGroupGetName: name field is not a String")
}
