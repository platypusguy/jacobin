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
	"jacobin/src/object"
	"jacobin/src/thread"
	"jacobin/src/types"
)

func Load_Lang_Thread_Group() {
	// <clinit>
	MethodSignatures["java/lang/ThreadGroup.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: threadGroupClinit}

	// Constructors
	MethodSignatures["java/lang/ThreadGroup.ThreadGroup(Ljava/lang/String;)Ljava/lang/ThreadGroup;"] =
		GMeth{ParamSlots: 1, GFunction: threadGroupCreateWithName}
	MethodSignatures["java/lang/ThreadGroup.ThreadGroup(Ljava/lang/ThreadGroup;Ljava/lang/String;)Ljava/lang/ThreadGroup;"] =
		GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.<init>(Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/lang/ThreadGroup.<init>(Ljava/lang/ThreadGroup;Ljava/lang/String;)V"] =
		GMeth{ParamSlots: 2, GFunction: trapFunction}

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

// java/lang/ThreadGroup.<clinit>()V
func threadGroupClinit(params []interface{}) any {
	return justReturn(nil)
}

// java/lang/ThreadGroup.ThreadGroup(Ljava/lang/String;)Ljava/lang/ThreadGroup;
// returns a new ThreadGroup object with the specified name and a null parent
func threadGroupCreateWithName(params []interface{}) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGroupCreateWithName: Expected thread group, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name, ok := params[0].(string)
	if !ok {
		errMsg := "threadGroupCreateWithName: Expected parameter to be a string"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	clName := "java/lang/ThreadGroup"
	obj := object.MakeEmptyObjectWithClassName(&clName)

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

	subgroups := object.Field{Ftype: types.LinkedList, Fvalue: list.List{}}
	obj.FieldTable["subgroups"] = subgroups

	return obj
}

// java/lang/ThreadGroup.getName()Ljava/lang/String;
func threadGroupGetName(params []interface{}) interface{} {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("threadGroupGetName: Expected thread group, got %d parameters", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	tg, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "threadGroupGetName: Expected parameter to be an object reference"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name := tg.FieldTable["name"].Fvalue.(string)
	return object.StringObjectFromGoString(name)
}
