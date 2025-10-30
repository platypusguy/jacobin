/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

func Load_Lang_Thread_Group() {
	// <clinit>
	MethodSignatures["java/lang/ThreadGroup.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: clinitGeneric}

	// Constructors
	MethodSignatures["java/lang/ThreadGroup.ThreadGroup(Ljava/lang/String;)Ljava/lang/ThreadGroup;"] =
		GMeth{ParamSlots: 1, GFunction: trapFunction}
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
		GMeth{ParamSlots: 0, GFunction: trapFunction}
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
