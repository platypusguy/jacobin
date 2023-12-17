/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/shutdown"
	"strings"
)

// StackTraceElement is a class that is primarily used by Throwable to gather data about the
// entries in the JVM stack. Because this data is so tightly bound to the specific implementation
// of the JVM, the methods of this class are fairly faithfully reproduced in the golang-native
// methods in this file. Consult:
// https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/StackTraceElement.html

func Load_Lang_StackTraceELement() map[string]GMeth {

	MethodSignatures["java/lang/StackTraceElement.of(Ljava/lang/Throwable;I)[Ljava/lang/StackTraceElement;"] =
		GMeth{
			ParamSlots:   2,
			GFunction:    of,
			NeedsContext: true,
		}

	MethodSignatures["java/lang/StackTraceElement.initStackTraceElements([Ljava/lang/StackTraceElement;Ljava/lang/Throwable;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  initStackTraceElements,
		}

	return MethodSignatures
}

/*
	 From the Java code for this method:

	 Returns an array of StackTraceElements of the given depth
	 filled from the backtrace of a given Throwable.

	    static StackTraceElement[] of(Throwable x, int depth) {
			StackTraceElement[] stackTrace = new StackTraceElement[depth];
			for (int i = 0; i < depth; i++) {
				stackTrace[i] = new StackTraceElement();
		  }

		// VM to fill in StackTraceElement
		initStackTraceElements(stackTrace, x);

		// ensure the proper StackTraceElement initialization
		for (StackTraceElement ste : stackTrace) {
			ste.computeFormat();
		}
		return stackTrace;
	}
*/
func of(params []interface{}) interface{} {

	throwable := params[0].(*object.Object)
	depth := params[1].(int64)

	// get a pointer to the JVM stack
	jvmStackRef := throwable.FieldTable["frameStackRef"].Fvalue.(*list.List)
	if jvmStackRef == nil {
		errMsg := "nil parameter for 'frameStackRef' in Throwable, found in StackTraceElement.of()"
		_ = log.Log(errMsg, log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// create the 1-dimensional array of stackTraceElements
	stackTraceElementClassName := "java/lang/StackTraceElement"
	stackTrace := object.Make1DimRefArray(&stackTraceElementClassName, depth)

	// insert empty stackTraceElements into the array.
	rawArrayPtr := stackTrace.Fields[0].Fvalue.(*[]*object.Object)
	rawArray := *rawArrayPtr
	global := globals.GetGlobalRef()
	for i := int64(0); i < depth; i++ {
		ste, err := global.FuncInstantiateClass("java/lang/StackTraceElement", nil)
		if err == nil {
			rawArray[i] = ste.(*object.Object)
		}
	}

	argsToPass := []interface{}{stackTrace, throwable}
	initStackTraceElements(argsToPass)

	return stackTrace
}

// This is a native function in HotSpot that accepts an array of empty
// stackTraceElements and a Throwable and fills in the values in the array
// by repeated calls to initStackTraceElement() below.
// Returns nothing.
func initStackTraceElements(params []interface{}) interface{} {
	arrayObjPtr := params[0].(*object.Object) // the array of stackTraceElements we'll fill in
	arrayObj := *arrayObjPtr
	rawSteArrayPtr := arrayObj.Fields[0].Fvalue.(*[]*object.Object)
	rawSteArray := *rawSteArrayPtr

	throwable := params[1].(*object.Object) // pointer to the Throwable object
	jvmStack := throwable.FieldTable["frameStackRef"].Fvalue.(*list.List)

	var i = 0
	for e := jvmStack.Front(); e != nil; e = e.Next() {
		frame := e.Value.(*frames.Frame)
		ste := rawSteArray[i]
		i += 1
		initStackTraceElement(ste, frame)
	}

	return nil

	// Note: an improvement is to add the logic to the class parse showing
	// the Java code source line number. JACOBIN-224 refers to this.
}

// initStackTraceElement accepts a single stackTraceElement and JVM stack
// info and fills in the former with the latter. It's a private method and
// called only from initStackTraceElements(), so we don't need it to strictly
// follow the HotSpot way of implementing it. Official definition:
// initStackTraceElement(Ljava/lang/StackTraceElement;Ljava/lang/StackFrameInfo;)V
func initStackTraceElement(ste *object.Object, frm *frames.Frame) {
	frame := *frm
	stackTrace := *ste

	// helper function to facilitate subsequent field updates
	// thanks to JetBrains' AI Assistant for this suggestion
	addField := func(name, value string) {
		fld := object.Field{}
		fld.Fvalue = value
		stackTrace.FieldTable[name] = &fld
	}

	addField("declaringClass", frame.ClName)
	addField("methodName", frame.MethName)

	methClass := classloader.MethAreaFetch(frame.ClName)

	addField("classLoaderName", methClass.Loader)
	addField("fileName", methClass.Data.SourceFile)
	addField("moduleName", methClass.Data.Module)

	sourceLineNumber := ""
	// now get the source line number for any non-JDK files
	if strings.HasPrefix(frame.ClName, "java") || strings.HasPrefix(frame.ClName, "jdk") ||
		strings.HasPrefix(frame.ClName, "sun.") || strings.HasPrefix(frame.MethName, "<init>") {
		addField("sourceLine", "")
	} else {
		class := classloader.MethAreaFetch(frame.ClName)
		sourceMap := class.Data.Methods[0].CodeAttr.BytecodeSourceMap
		prev := uint16(0)
		for i := 0; i < len(sourceMap); i++ {
			entry := sourceMap[i]
			if entry.BytecodePos > uint16(frame.PC) {
				break
			} else {
				prev = entry.BytecodePos
			}
		}
		sourceLineNumber = fmt.Sprintf("%d", prev)
		addField("sourceLine", sourceLineNumber)
	}
}
