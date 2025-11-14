/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/thread"
	"jacobin/src/types"
)

func Load_Lang_Thread_Group() {

	// <clinit>
	MethodSignatures["java/lang/ThreadGroup.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: threadGroupClinit}

	// Constructors
	MethodSignatures["java/lang/ThreadGroup.<init>(Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 1, GFunction: threadGroupInitWithName}
	MethodSignatures["java/lang/ThreadGroup.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 2, GFunction: threadGroupInitWithParentAndName}

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

// Initialize global thread groups: create "system" group and its child "main"
func InitializeGlobalThreadGroups() {
	gr := globals.GetGlobalRef()
	if gr.ThreadGroups == nil {
		gr.ThreadGroups = make(map[string]interface{})
	}

	// in test mode, we don't need to create the system and main
	// thread groups manually
	if gr.JacobinName == "test" {
		fakeSystemTg := threadGroupFake("system")
		gr.ThreadGroups["system"] = fakeSystemTg

		fakeMainTg := threadGroupFake("main")
		gr.ThreadGroups["main"] = fakeMainTg

		fakeMainTg.FieldTable["parent"] =
			object.Field{Ftype: types.Ref, Fvalue: fakeSystemTg}

		// Now add this thread group to the parent's list of subgroups
		parentSubgroups := fakeSystemTg.FieldTable["subgroups"].Fvalue.(*list.List)
		parentSubgroups.PushBack(fakeMainTg)

		return
	}

	// java/lang/ThreadGroup has no clinit to run, so we can turn it off
	// likewise its parent is java/lang/Object, which does not need clinit
	// This is necessary because we call instantiateClass directly with no
	// JVM framestack
	k := classloader.MethAreaFetch("java/lang/ThreadGroup")
	k.Data.ClInit = types.ClInitRun

	// instantiate a new ThreadGroup object with the name "system"
	tg, err := gr.FuncInstantiateClass("java/lang/ThreadGroup", nil)
	if err != nil {
		errMsg :=
			"initializeGlobalThreadGroups: Failed to instantiate java/lang/ThreadGroup"
		gr.FuncThrowException(excNames.InternalError, errMsg)
	}

	systg := threadGroupInitWithName(
		[]interface{}{tg, object.StringObjectFromGoString("system")})
	gr.ThreadGroups["system"] = systg

	// Create "main" group as a child of "system"
	maintg, err := gr.FuncInstantiateClass("java/lang/ThreadGroup", nil)
	if err != nil {
		errMsg :=
			"initializeGlobalThreadGroups: Failed to instantiate java/lang/ThreadGroup"
		gr.FuncThrowException(excNames.InternalError, errMsg)
	}

	maintg = threadGroupInitWithName(
		[]interface{}{tg, object.StringObjectFromGoString("main")})

	gr.ThreadGroups["main"] = maintg
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
			"ThreadGroupInitWithParentNameMaxpriorityDaemon: Expect 1st parameter to be an object reference")
	}

	parentObj, ok := initParams[1].(*object.Object)
	if !ok && parentObj != object.Null {
		return getGErrBlk(excNames.IllegalArgumentException,
			"ThreadGroupInitWithParentNameMaxpriorityDaemon: Expect 2nd parameter to be an object reference")
	}

	nameObj, ok := initParams[2].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"ThreadGroupInitWithParentNameMaxpriorityDaemon: Expected third parameter to be a String object")
	}

	maxPriority := initParams[3].(int64)
	daemon := initParams[4].(types.JavaBool)

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

	if daemon != types.JavaBoolUninitialized {
		if daemon == types.JavaBoolFalse || daemon == types.JavaBoolTrue {
			obj.FieldTable["daemon"] = object.Field{Ftype: types.Bool, Fvalue: daemon}
		}
	}

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

	args := []interface{}{obj, object.Null, name, int64(0), types.JavaBoolUninitialized}
	updatedObj := ThreadGroupInitWithParentNameMaxpriorityDaemon(args)
	return updatedObj

}

// java/lang/ThreadGroup.ThreadGroup(Ljava/lang/ThreadGroup;Ljava/lang/String;)Ljava/lang/ThreadGroup;
// returns a new ThreadGroup object with the specified name and parent
// note: because the parent is specified, we add this thread group to the parent's list of subgroups'
func threadGroupInitWithParentAndName(params []interface{}) any {
	if len(params) != 3 {
		errMsg := fmt.Sprintf(
			"threadGroupCreateWithParentAndName: Expected 2 parameters, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// First param: parent ThreadGroup (object reference)
	obj, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGroupInitWithParentAndName: Expected 1st parameter to be an object reference, was not"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Second param: parent ThreadGroup (object reference)
	parentObj, ok := params[1].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException,
			"threadGroupCreateWithParentAndName: Expected first parameter to be a ThreadGroup object")
	}

	// Third param: name (Java String object)
	nameObj, ok := params[2].(*object.Object)
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
			"threadGroupCreateWithParentAndName: name is not a String")
	}

	args := []interface{}{obj, parentObj, nameObj, int64(0), types.JavaBoolUninitialized}
	tg := ThreadGroupInitWithParentNameMaxpriorityDaemon(args)

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

// == fake thread group for testing purposes ==
func threadGroupFake(name string) *object.Object {
	clName := "java/lang/ThreadGroup"
	obj := object.MakeEmptyObjectWithClassName(&clName)

	nullField := object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable["parent"] = nullField

	nameField := object.Field{
		Ftype:  types.Ref,
		Fvalue: object.JavaByteArrayFromGoString(name)}
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
	globals.GetGlobalRef().ThreadGroups[name] = obj

	return obj
}
